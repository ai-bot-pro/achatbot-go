package params

import (
	"achatbot/pkg/common"
	"fmt"
	"reflect"
)

// AudioCameraParams represents audio and camera parameters configuration
type AudioCameraParams struct {
	*AudioVADParams

	// Camera output parameters
	CameraOutEnabled   bool `json:"camera_out_enabled"`
	CameraOutWidth     int  `json:"camera_out_width"`
	CameraOutHeight    int  `json:"camera_out_height"`
	CameraOutFramerate int  `json:"camera_out_framerate"`
	CameraOutIsLive    bool `json:"camera_out_is_live"`

	// Transport writer
	TransportWriter common.ITransportWriter
}

// NewAudioCameraParams creates a new AudioCameraParams with default values
func NewAudioCameraParams() *AudioCameraParams {
	return &AudioCameraParams{
		AudioVADParams:     NewAudioVADParams(),
		CameraOutEnabled:   false,
		CameraOutWidth:     640,
		CameraOutHeight:    480,
		CameraOutFramerate: 30,
		CameraOutIsLive:    false,
		TransportWriter:    nil,
	}
}

// WithAudioVADParams sets audio VAD parameters
func (p *AudioCameraParams) WithAudioVADParams(AudioVADParams *AudioVADParams) *AudioCameraParams {
	p.AudioVADParams = AudioVADParams
	return p
}

// WithCameraOutEnabled sets camera output enabled
func (p *AudioCameraParams) WithCameraOutEnabled(enabled bool) *AudioCameraParams {
	p.CameraOutEnabled = enabled
	return p
}

// WithCameraOutWidth sets camera output width
func (p *AudioCameraParams) WithCameraOutWidth(width int) *AudioCameraParams {
	p.CameraOutWidth = width
	return p
}

// WithCameraOutHeight sets camera output height
func (p *AudioCameraParams) WithCameraOutHeight(height int) *AudioCameraParams {
	p.CameraOutHeight = height
	return p
}

// WithCameraOutFramerate sets camera output framerate
func (p *AudioCameraParams) WithCameraOutFramerate(framerate int) *AudioCameraParams {
	p.CameraOutFramerate = framerate
	return p
}

// WithCameraOutIsLive sets camera output is live
func (p *AudioCameraParams) WithCameraOutIsLive(isLive bool) *AudioCameraParams {
	p.CameraOutIsLive = isLive
	return p
}

// WithTransportWriter sets transport writer
func (p *AudioCameraParams) WithTransportWriter(TransportWriter common.ITransportWriter) *AudioCameraParams {
	p.TransportWriter = TransportWriter
	return p
}

// GetCameraOutWidth returns camera output width
func (p *AudioCameraParams) GetCameraOutWidth() int {
	return p.CameraOutWidth
}

// GetCameraOutHeight returns camera output height
func (p *AudioCameraParams) GetCameraOutHeight() int {
	return p.CameraOutHeight
}

// GetCameraOutFramerate returns camera output framerate
func (p *AudioCameraParams) GetCameraOutFramerate() int {
	return p.CameraOutFramerate
}

// IsCameraOutEnabled checks if camera output is enabled
func (p *AudioCameraParams) IsCameraOutEnabled() bool {
	return p.CameraOutEnabled
}

// IsCameraOutIsLive checks if camera output is live
func (p *AudioCameraParams) IsCameraOutIsLive() bool {
	return p.CameraOutIsLive
}

// GetTransportWriter returns the transport writer
func (p *AudioCameraParams) GetTransportWriter() common.ITransportWriter {
	return p.TransportWriter
}

func (p *AudioCameraParams) String() string {
	return fmt.Sprintf("AudioCameraParams{AudioVADParams: %s, CameraOutEnabled: %t, CameraOutFramerate: %d, CameraOutHeight: %d, CameraOutIsLive: %t, CameraOutWidth: %d, TransportWriter: %s}",
		p.AudioVADParams.String(),
		p.CameraOutEnabled,
		p.CameraOutFramerate,
		p.CameraOutHeight,
		p.CameraOutIsLive,
		p.CameraOutWidth,
		reflect.TypeOf(p.TransportWriter),
	)
}
