package processors

import (
	"log/slog"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/common"
	"achatbot/pkg/params"
)

// WebsocketServerOutputProcessor processes output for  WebSocket server
type WebsocketServerOutputProcessor struct {
	*AudioCameraOutputProcessor
}

// NewWebsocketServerOutputProcessor creates a new WebsocketServerOutputProcessor
func NewWebsocketServerOutputProcessor(
	name string,
	websocket common.IWebSocketConn,
	params *params.WebsocketServerParams,
) *WebsocketServerOutputProcessor {
	p := &WebsocketServerOutputProcessor{
		AudioCameraOutputProcessor: NewAudioCameraOutputProcessor(name, params.AudioCameraParams),
	}

	return p
}

// ProcessFrame processes a frame
func (p *WebsocketServerOutputProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// Call parent implementation
	p.AudioCameraOutputProcessor.ProcessFrame(frame, direction)

	// Handle specific frame types
	switch f := frame.(type) {
	case *frames.StartInterruptionFrame:
		err := p.transportWriter.WriteFrame(f)
		if err != nil {
			slog.Error("Error send StartInterruptionFrame", "error", err)
		}
	}
}
