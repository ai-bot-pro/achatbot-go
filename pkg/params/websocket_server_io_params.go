package params

import (
	"github.com/weedge/pipeline-go/pkg/serializers"
)

// WebsocketServerParams represents parameters for the  WebSocket server
type WebsocketServerParams struct {
	*AudioCameraParams
	Serializer         serializers.Serializer
	AddWavHeader       bool `json:"add_wav_header"`
	AudioFrameSize     int  `json:"audio_frame_size"`
	AudioOutSampleRate int  `json:"audio_out_sample_rate"`
	AudioOutChannels   int  `json:"audio_out_channels"`
}

// NewWebsocketServerParams creates a new WebsocketServerParams with default values
func NewWebsocketServerParams() *WebsocketServerParams {
	return &WebsocketServerParams{
		AudioCameraParams:  NewAudioCameraParams(),
		Serializer:         serializers.NewJsonSerializer(),
		AddWavHeader:       false,
		AudioFrameSize:     16000 * 2 * 10 / 1000, // 10ms chunks at 16kHz 16-bit mono
		AudioOutSampleRate: 16000,
		AudioOutChannels:   1,
	}
}

// WithSerializer sets the serializer
func (p *WebsocketServerParams) WithSerializer(serializer serializers.Serializer) *WebsocketServerParams {
	p.Serializer = serializer
	return p
}

// WithAddWavHeader sets whether to add WAV header
func (p *WebsocketServerParams) WithAddWavHeader(addWavHeader bool) *WebsocketServerParams {
	p.AddWavHeader = addWavHeader
	return p
}

// WithAudioFrameSize sets the audio frame size
func (p *WebsocketServerParams) WithAudioFrameSize(size int) *WebsocketServerParams {
	p.AudioFrameSize = size
	return p
}

// WithAudioOutSampleRate sets the audio output sample rate
func (p *WebsocketServerParams) WithAudioOutSampleRate(rate int) *WebsocketServerParams {
	p.AudioOutSampleRate = rate
	return p
}

// WithAudioOutChannels sets the audio output channels
func (p *WebsocketServerParams) WithAudioOutChannels(channels int) *WebsocketServerParams {
	p.AudioOutChannels = channels
	return p
}
