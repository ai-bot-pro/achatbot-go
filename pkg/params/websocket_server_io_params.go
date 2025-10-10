package params

import (
	"github.com/weedge/pipeline-go/pkg/serializers"
)

// WebsocketServerParams represents parameters for the  WebSocket server
type WebsocketServerParams struct {
	*AudioCameraParams
	Serializer           serializers.Serializer
	AudioOutAddWavHeader bool `json:"audio_out_add_wav_header"`
	AudioOutFrameSize    int  `json:"audio_frame_size"`
}

// NewWebsocketServerParams creates a new WebsocketServerParams with default values
func NewWebsocketServerParams() *WebsocketServerParams {
	return &WebsocketServerParams{
		AudioCameraParams:    NewAudioCameraParams(),
		Serializer:           serializers.NewProtobufSerializer(),
		AudioOutAddWavHeader: false,
		AudioOutFrameSize:    6400, // 200ms with 16K hz 1 channel 2 sample_width
	}
}

// WithSerializer sets the serializer
func (p *WebsocketServerParams) WithSerializer(serializer serializers.Serializer) *WebsocketServerParams {
	p.Serializer = serializer
	return p
}

// WithAudioOutAddWavHeader sets whether to add WAV header
func (p *WebsocketServerParams) WithAudioOutAddWavHeader(AudioOutAddWavHeader bool) *WebsocketServerParams {
	p.AudioOutAddWavHeader = AudioOutAddWavHeader
	return p
}

// WithAudioOutFrameSize sets the audio frame size
func (p *WebsocketServerParams) WithAudioOutFrameSize(size int) *WebsocketServerParams {
	p.AudioOutFrameSize = size
	return p
}
