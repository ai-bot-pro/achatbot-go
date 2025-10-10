package processors

import (
	"context"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	"achatbot/pkg/consts"
	"achatbot/pkg/params"
)

// WebsocketServerCallbacks defines callback functions for WebSocket events
type WebsocketServerCallbacks struct {
	OnClientConnected    func(ws common.IWebSocketConn)
	OnClientDisconnected func(ws common.IWebSocketConn)
}

// WebsocketServerInputProcessor processes audio input from WebSocket connections
type WebsocketServerInputProcessor struct {
	*AudioVADInputProcessor
	websocket  common.IWebSocketConn
	params     *params.WebsocketServerParams
	callbacks  *WebsocketServerCallbacks
	receiveCtx context.Context
	cancelRecv context.CancelFunc
}

// NewWebsocketServerInputProcessor creates a new WebsocketServerInputProcessor
func NewWebsocketServerInputProcessor(
	name string,
	websocket common.IWebSocketConn,
	params *params.WebsocketServerParams,
	callbacks *WebsocketServerCallbacks,
) *WebsocketServerInputProcessor {
	// Create the base audio VAD processor
	audioVADProcessor := NewAudioVADInputProcessor(name, params.AudioVADParams)

	// Create context for receive loop
	receiveCtx, cancelRecv := context.WithCancel(context.Background())

	return &WebsocketServerInputProcessor{
		AudioVADInputProcessor: audioVADProcessor,
		websocket:              websocket,
		params:                 params,
		callbacks:              callbacks,
		receiveCtx:             receiveCtx,
		cancelRecv:             cancelRecv,
	}
}

// Start starts the WebSocket processor
func (p *WebsocketServerInputProcessor) Start(frame *frames.StartFrame) {
	// Notify client connected
	if p.callbacks.OnClientConnected != nil {
		p.callbacks.OnClientConnected(p.websocket)
	}

	// Start receiving messages in a goroutine
	go p.receiveMessages()
	logger.Info("WebsocketServerInputProcessor Start")
}

// Stop stops the WebSocket processor
func (p *WebsocketServerInputProcessor) Stop(frame *frames.EndFrame) {
	logger.Info("WebsocketServerInputProcessor Stopping")

	// Cancel receive loop
	p.cancelRecv()

	// Close WebSocket connection if it's not already closed
	if p.websocket != nil {
		// Note: we typically don't check the connection state before closing
		// The Close method handles this internally
		p.websocket.Close()
	}

}

// Cancel cancels the WebSocket processor
func (p *WebsocketServerInputProcessor) Cancel(frame *frames.CancelFrame) {
	logger.Info("WebsocketServerInputProcessor Cancelling")

	// Cancel receive loop
	p.cancelRecv()

	// Close WebSocket connection if it's not already closed
	if p.websocket != nil {
		p.websocket.Close()
	}

	logger.Info("WebsocketServerInputProcessor Cancel Done")
}

// ProcessFrame processes a frame
func (p *WebsocketServerInputProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// call frame processor to init start/end/cancel/interruption frame and push/queue frame
	p.AudioVADInputProcessor.ProcessFrame(frame, direction)

	// process self start/end/cancel frame logic,
	// NOTE: don't to push frame again
	switch f := frame.(type) {
	case *frames.StartFrame:
		p.Start(f)
	case *frames.EndFrame:
		p.Stop(f)
	case *frames.CancelFrame:
		p.Cancel(f)
	}

}

// receiveMessages receives messages from the WebSocket connection
func (p *WebsocketServerInputProcessor) receiveMessages() {
	defer func() {
		logger.Info("WebSocket connection disconnected")
		if p.callbacks.OnClientDisconnected != nil {
			p.callbacks.OnClientDisconnected(p.websocket)
		}
	}()

	for {
		select {
		case <-p.receiveCtx.Done():
			// Context cancelled, exit the loop
			return
		default:
			// Read message from WebSocket
			messageType, message, err := p.websocket.ReadMessage()
			if err != nil {
				logger.Error("Error reading WebSocket message", "error", err)
				return
			}

			// We're only interested in binary messages (similar to iter_bytes() in Python)
			if messageType != consts.BinaryMessage {
				logger.Warn("Warning only interested in binary messages(serialized)", "messageType", messageType)
				continue
			}

			// Deserialize the message
			frame, err := p.params.Serializer.Deserialize(message)
			if err != nil {
				logger.Error("Error deserializing WebSocket message", "error", err)
				continue
			}
			if frame == nil {
				logger.Warn("Warning deserializing frame is nil")
				continue
			}

			// Process audio raw frames
			switch frame := frame.(type) {
			case *frames.AudioRawFrame:
				err := p.PushAudioFrame(frame)
				if err != nil {
					logger.Error("Error pushing audio frame", "error", err)
				}
			}

		}
	}
}
