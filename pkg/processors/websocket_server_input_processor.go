package processors

import (
	"context"
	"log/slog"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/serializers"

	"achatbot/pkg/common"
	"achatbot/pkg/params"
	"achatbot/pkg/types/networks"
)

// WebsocketServerParams represents parameters for the WebSocket server
type WebsocketServerParams struct {
	*params.AudioVADParams
	Serializer serializers.Serializer
}

// WebsocketServerCallbacks defines callback functions for WebSocket events
type WebsocketServerCallbacks struct {
	OnClientConnected    func(ws common.WebSocketConn)
	OnClientDisconnected func(ws common.WebSocketConn)
}

// WebsocketServerInputProcessor processes audio input from WebSocket connections
type WebsocketServerInputProcessor struct {
	*AudioVADInputProcessor
	websocket  common.WebSocketConn
	params     *WebsocketServerParams
	callbacks  *WebsocketServerCallbacks
	receiveCtx context.Context
	cancelRecv context.CancelFunc
}

// NewWebsocketServerInputProcessor creates a new WebsocketServerInputProcessor
func NewWebsocketServerInputProcessor(
	name string,
	websocket common.WebSocketConn,
	params *WebsocketServerParams,
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
	slog.Info("Starting WebSocket server input processor")

	// Call parent start
	p.AudioVADInputProcessor.Start(frame)

	// Notify client connected
	if p.callbacks.OnClientConnected != nil {
		p.callbacks.OnClientConnected(p.websocket)
	}

	// Start receiving messages in a goroutine
	go p.receiveMessages()
}

// Stop stops the WebSocket processor
func (p *WebsocketServerInputProcessor) Stop() {
	slog.Info("Stopping WebSocket server input processor")

	// Cancel receive loop
	p.cancelRecv()

	// Close WebSocket connection if it's not already closed
	if p.websocket != nil {
		// Note: we typically don't check the connection state before closing
		// The Close method handles this internally
		p.websocket.Close()
	}

	// Call parent stop
	p.AudioVADInputProcessor.Stop()
}

// Cancel cancels the WebSocket processor
func (p *WebsocketServerInputProcessor) Cancel(frame *frames.CancelFrame) {
	slog.Info("Cancelling WebSocket server input processor")

	// Cancel receive loop
	p.cancelRecv()

	// Close WebSocket connection if it's not already closed
	if p.websocket != nil {
		p.websocket.Close()
	}

	// Call parent cancel
	p.AudioVADInputProcessor.Cancel(frame)
}

// receiveMessages receives messages from the WebSocket connection
func (p *WebsocketServerInputProcessor) receiveMessages() {
	defer func() {
		slog.Info("WebSocket connection disconnected")
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
				slog.Error("Error reading WebSocket message", "error", err)
				return
			}

			// We're only interested in binary messages (similar to iter_bytes() in Python)
			if messageType != networks.BinaryMessage {
				slog.Warn("Warning only interested in binary messages(serialized)", "messageType", messageType)
				continue
			}

			// Deserialize the message
			frame, err := p.params.Serializer.Deserialize(message)
			if err != nil {
				slog.Error("Error deserializing WebSocket message", "error", err)
				continue
			}
			if frame == nil {
				slog.Warn("Warning deserializing frame is nil")
				continue
			}

			// Process audio raw frames
			switch frame := frame.(type) {
			case *frames.AudioRawFrame:
				err := p.PushAudioFrame(frame)
				if err != nil {
					slog.Error("Error pushing audio frame", "error", err)
				}
			}

		}
	}
}
