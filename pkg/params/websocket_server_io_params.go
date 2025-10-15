package params

import (
	"fmt"

	"github.com/weedge/pipeline-go/pkg/serializers"
)

// WebsocketServerParams represents parameters for the  WebSocket server
type WebsocketServerParams struct {
	*AudioCameraParams
	Serializer           serializers.Serializer
	AudioOutAddWavHeader bool `json:"audio_out_add_wav_header"`
	AudioOutFrameMS      int  `json:"audio_out_frame_ms"`
}

// NewWebsocketServerParams creates a new WebsocketServerParams with default values
func NewWebsocketServerParams() *WebsocketServerParams {
	return &WebsocketServerParams{
		AudioCameraParams:    NewAudioCameraParams(),
		Serializer:           serializers.NewProtobufSerializer(),
		AudioOutAddWavHeader: false,
		AudioOutFrameMS:      200, // 200ms
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

// WithAudioOutFrameMS sets the audio frame size
func (p *WebsocketServerParams) WithAudioOutFrameMS(audioOutFrameMS int) *WebsocketServerParams {
	p.AudioOutFrameMS = audioOutFrameMS
	return p
}

func (p *WebsocketServerParams) String() string {
	return fmt.Sprintf("WebsocketServerParams{AudioCameraParams: %s, AudioOutAddWavHeader: %t, AudioOutFrameMS: %d, Serializer: %s}", p.AudioCameraParams, p.AudioOutAddWavHeader, p.AudioOutFrameMS, p.Serializer)
}
