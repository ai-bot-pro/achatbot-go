package processors

import (
	"bytes"
	"encoding/binary"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"

	"achatbot/pkg/common"
	"achatbot/pkg/consts"
	"achatbot/pkg/params"
	achatbot_frames "achatbot/pkg/types/frames"
)

// WebsocketTransportWriter processes output for  WebSocket server
type WebsocketTransportWriter struct {
	websocket            common.IWebSocketConn
	params               *params.WebsocketServerParams
	websocketAudioBuffer []byte
}

// NewWebsocketTransportWriter creates a new WebsocketTransportWriter
func NewWebsocketTransportWriter(
	websocket common.IWebSocketConn,
	params *params.WebsocketServerParams,
) *WebsocketTransportWriter {
	p := &WebsocketTransportWriter{
		websocket:            websocket,
		params:               params,
		websocketAudioBuffer: make([]byte, 0),
	}

	return p
}

func (p *WebsocketTransportWriter) WriteRawAudio(data []byte) error {
	p.websocketAudioBuffer = append(p.websocketAudioBuffer, data...)

	for len(p.websocketAudioBuffer) >= p.params.AudioOutFrameSize {
		frame := frames.NewAudioRawFrame(
			p.websocketAudioBuffer[:p.params.AudioOutFrameSize],
			p.params.AudioOutSampleRate,
			p.params.AudioOutChannels,
			p.params.AudioOutSampleWidth,
		)

		if p.params.AudioOutAddWavHeader && len(frame.Audio) > 0 {
			wavData := p.AudioOutAddWavHeader(frame)
			wavFrame := frames.NewAudioRawFrame(
				wavData,
				frame.SampleRate,
				frame.NumChannels,
				frame.SampleWidth,
			)
			frame = wavFrame
		}

		// 安全地更新缓冲区
		p.websocketAudioBuffer = p.websocketAudioBuffer[p.params.AudioOutFrameSize:]

		err := p.SendPayload(frame)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *WebsocketTransportWriter) WriteFrame(frame frames.Frame) error {
	var err error
	switch f := frame.(type) {
	case *frames.TextFrame, *frames.StartInterruptionFrame:
		err = p.SendPayload(f)
	case *achatbot_frames.AnimationAudioRawFrame:
		err = p.WriteAnimationAudioFrame(f)
	}
	return err
}

// WriteAnimationAudioFrame writes an animation audio frame to the WebSocket
func (p *WebsocketTransportWriter) WriteAnimationAudioFrame(frame *achatbot_frames.AnimationAudioRawFrame) error {
	if p.params.AudioOutAddWavHeader && len(frame.Audio) > 0 {
		// Create a copy of the frame with WAV header
		wavData := p.AudioOutAddWavHeader(frame.AudioRawFrame)
		frame.Audio = wavData
	}

	return p.SendPayload(frame)
}

// SendPayload sends a payload to the WebSocket
func (p *WebsocketTransportWriter) SendPayload(frame frames.Frame) error {
	// Serialize the frame
	payload, err := p.params.Serializer.Serialize(frame)
	if err != nil {
		logger.Error("serialize frame error", "error", err, "frame", frame)
		return err
	}
	if len(payload) == 0 {
		logger.Warn("serialize frame produced no payload", "frame", frame)
		return nil
	}

	// Send the payload
	messageType := consts.BinaryMessage // BinaryMessage by default
	if isStringPayload(payload) {
		messageType = consts.TextMessage // TextMessage
	}

	err = p.websocket.WriteMessage(messageType, payload)
	if err != nil {
		logger.Error("send_payload error", "error", err)
		return err
	}

	return nil
}

// AudioOutAddWavHeader adds a WAV header to raw audio data
func (p *WebsocketTransportWriter) AudioOutAddWavHeader(frame *frames.AudioRawFrame) []byte {
	if len(frame.Audio) == 0 {
		return frame.Audio
	}

	// Create WAV header
	buf := new(bytes.Buffer)

	// RIFF header
	buf.WriteString("RIFF")
	binary.Write(buf, binary.LittleEndian, uint32(36+len(frame.Audio))) // ChunkSize
	buf.WriteString("WAVE")

	// fmt subchunk
	buf.WriteString("fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16))                // Subchunk1Size
	binary.Write(buf, binary.LittleEndian, uint16(1))                 // AudioFormat (PCM)
	binary.Write(buf, binary.LittleEndian, uint16(frame.NumChannels)) // NumChannels
	binary.Write(buf, binary.LittleEndian, uint32(frame.SampleRate))  // SampleRate
	byteRate := frame.SampleRate * frame.NumChannels * frame.SampleWidth
	binary.Write(buf, binary.LittleEndian, uint32(byteRate)) // ByteRate
	blockAlign := frame.NumChannels * frame.SampleWidth
	binary.Write(buf, binary.LittleEndian, uint16(blockAlign))          // BlockAlign
	binary.Write(buf, binary.LittleEndian, uint16(frame.SampleWidth*8)) // BitsPerSample

	// data subchunk
	buf.WriteString("data")
	binary.Write(buf, binary.LittleEndian, uint32(len(frame.Audio))) // Subchunk2Size
	buf.Write(frame.Audio)

	return buf.Bytes()
}

// isStringPayload checks if payload should be sent as text
func isStringPayload(payload []byte) bool {
	// Simple check: if payload is valid UTF-8 and looks like JSON, send as text
	// Otherwise send as binary
	if len(payload) == 0 {
		return false
	}

	// Check if first character suggests JSON ({ or [)
	firstChar := payload[0]
	return firstChar == '{' || firstChar == '['
}
