package params

import (
	"achatbot/pkg/common"
	"fmt"
)

// AudioParams represents audio parameters configuration
type AudioParams struct {
	AudioOutEnabled           bool `json:"audio_out_enabled"`
	AudioOutSampleRate        int  `json:"audio_out_sample_rate"`
	AudioOutChannels          int  `json:"audio_out_channels"`
	AudioOutSampleWidth       int  `json:"audio_out_sample_width"`
	AudioOut10msChunks        int  `json:"audio_out_10ms_chunks"`
	AudioInEnabled            bool `json:"audio_in_enabled"`
	AudioInParticipantEnabled bool `json:"audio_in_participant_enabled"`
	AudioInSampleRate         int  `json:"audio_in_sample_rate"`
	AudioInChannels           int  `json:"audio_in_channels"`
	AudioInSampleWidth        int  `json:"audio_in_sample_width"`
}

// NewAudioParams creates a new AudioParams with default values
func NewAudioParams() *AudioParams {
	return &AudioParams{
		AudioOutSampleRate:  16000, // Default RATE
		AudioOutChannels:    1,     // Default CHANNELS
		AudioOutSampleWidth: 2,     // Default SAMPLE_WIDTH
		AudioOut10msChunks:  2,
		AudioInSampleRate:   16000, // Default RATE
		AudioInChannels:     1,     // Default CHANNELS
		AudioInSampleWidth:  2,     // Default SAMPLE_WIDTH
	}
}

// AudioVADParams extends AudioParams with VAD-specific parameters
type AudioVADParams struct {
	AudioParams
	VADEnabled          bool `json:"vad_enabled"`
	VADAudioPassthrough bool `json:"vad_audio_passthrough"`
	VADAnalyzer         common.IVADAnalyzer
}

// NewAudioVADParams creates a new AudioVADParams with default values
func NewAudioVADParams() *AudioVADParams {
	return &AudioVADParams{
		AudioParams:         *NewAudioParams(),
		VADEnabled:          false,
		VADAudioPassthrough: false,
		VADAnalyzer:         nil,
	}
}

// WithAudioOutEnabled sets audio output enabled
func (p *AudioParams) WithAudioOutEnabled(enabled bool) *AudioParams {
	p.AudioOutEnabled = enabled
	return p
}

// WithAudioOutSampleRate sets audio output sample rate
func (p *AudioParams) WithAudioOutSampleRate(rate int) *AudioParams {
	p.AudioOutSampleRate = rate
	return p
}

// WithAudioOutChannels sets audio output channels
func (p *AudioParams) WithAudioOutChannels(channels int) *AudioParams {
	p.AudioOutChannels = channels
	return p
}

// WithAudioOutSampleWidth sets audio output sample width
func (p *AudioParams) WithAudioOutSampleWidth(width int) *AudioParams {
	p.AudioOutSampleWidth = width
	return p
}

// WithAudioOut10msChunks sets audio output 10ms chunks
func (p *AudioParams) WithAudioOut10msChunks(chunks int) *AudioParams {
	p.AudioOut10msChunks = chunks
	return p
}

// WithAudioInEnabled sets audio input enabled
func (p *AudioParams) WithAudioInEnabled(enabled bool) *AudioParams {
	p.AudioInEnabled = enabled
	return p
}

// WithAudioInParticipantEnabled sets audio input participant enabled
func (p *AudioParams) WithAudioInParticipantEnabled(enabled bool) *AudioParams {
	p.AudioInParticipantEnabled = enabled
	return p
}

// WithAudioInSampleRate sets audio input sample rate
func (p *AudioParams) WithAudioInSampleRate(rate int) *AudioParams {
	p.AudioInSampleRate = rate
	return p
}

// WithAudioInChannels sets audio input channels
func (p *AudioParams) WithAudioInChannels(channels int) *AudioParams {
	p.AudioInChannels = channels
	return p
}

// WithAudioInSampleWidth sets audio input sample width
func (p *AudioParams) WithAudioInSampleWidth(width int) *AudioParams {
	p.AudioInSampleWidth = width
	return p
}

// WithVADEnabled sets VAD enabled
func (p *AudioVADParams) WithVADEnabled(enabled bool) *AudioVADParams {
	p.VADEnabled = enabled
	return p
}

// WithVADAudioPassthrough sets VAD audio passthrough
func (p *AudioVADParams) WithVADAudioPassthrough(enabled bool) *AudioVADParams {
	p.VADAudioPassthrough = enabled
	return p
}

// WithVADAnalyzer sets VAD analyzer
func (p *AudioVADParams) WithVADAnalyzer(analyzer common.IVADAnalyzer) *AudioVADParams {
	p.VADAnalyzer = analyzer
	return p
}

// GetAudioOutSampleRate returns audio output sample rate
func (p *AudioParams) GetAudioOutSampleRate() int {
	return p.AudioOutSampleRate
}

// GetAudioOutChannels returns audio output channels
func (p *AudioParams) GetAudioOutChannels() int {
	return p.AudioOutChannels
}

// GetAudioOutSampleWidth returns audio output sample width
func (p *AudioParams) GetAudioOutSampleWidth() int {
	return p.AudioOutSampleWidth
}

// GetAudioInSampleRate returns audio input sample rate
func (p *AudioParams) GetAudioInSampleRate() int {
	return p.AudioInSampleRate
}

// GetAudioInChannels returns audio input channels
func (p *AudioParams) GetAudioInChannels() int {
	return p.AudioInChannels
}

// GetAudioInSampleWidth returns audio input sample width
func (p *AudioParams) GetAudioInSampleWidth() int {
	return p.AudioInSampleWidth
}

// IsAudioOutEnabled checks if audio output is enabled
func (p *AudioParams) IsAudioOutEnabled() bool {
	return p.AudioOutEnabled
}

// IsAudioInEnabled checks if audio input is enabled
func (p *AudioParams) IsAudioInEnabled() bool {
	return p.AudioInEnabled
}

// IsAudioInParticipantEnabled checks if audio input participant is enabled
func (p *AudioParams) IsAudioInParticipantEnabled() bool {
	return p.AudioInParticipantEnabled
}

// IsVADEnabled checks if VAD is enabled
func (p *AudioVADParams) IsVADEnabled() bool {
	return p.VADEnabled
}

// IsVADAudioPassthrough checks if VAD audio passthrough is enabled
func (p *AudioVADParams) IsVADAudioPassthrough() bool {
	return p.VADAudioPassthrough
}

// GetVADAnalyzer returns the VAD analyzer
func (p *AudioVADParams) GetVADAnalyzer() common.IVADAnalyzer {
	return p.VADAnalyzer
}

// Validate validates the audio parameters
func (p *AudioParams) Validate() error {
	if p.AudioOutSampleRate <= 0 {
		return fmt.Errorf("audio_out_sample_rate must be positive, got %d", p.AudioOutSampleRate)
	}
	if p.AudioOutChannels <= 0 {
		return fmt.Errorf("audio_out_channels must be positive, got %d", p.AudioOutChannels)
	}
	if p.AudioOutSampleWidth <= 0 {
		return fmt.Errorf("audio_out_sample_width must be positive, got %d", p.AudioOutSampleWidth)
	}
	if p.AudioOut10msChunks <= 0 {
		return fmt.Errorf("audio_out_10ms_chunks must be positive, got %d", p.AudioOut10msChunks)
	}
	if p.AudioInSampleRate <= 0 {
		return fmt.Errorf("audio_in_sample_rate must be positive, got %d", p.AudioInSampleRate)
	}
	if p.AudioInChannels <= 0 {
		return fmt.Errorf("audio_in_channels must be positive, got %d", p.AudioInChannels)
	}
	if p.AudioInSampleWidth <= 0 {
		return fmt.Errorf("audio_in_sample_width must be positive, got %d", p.AudioInSampleWidth)
	}
	return nil
}

// Validate validates the VAD parameters
func (p *AudioVADParams) Validate() error {
	if err := p.AudioParams.Validate(); err != nil {
		return err
	}

	if p.VADEnabled && p.VADAnalyzer == nil {
		return fmt.Errorf("VAD is enabled but no VAD analyzer provided")
	}

	return nil
}

// Clone creates a copy of AudioParams
func (p *AudioParams) Clone() *AudioParams {
	clone := *p
	return &clone
}

// Clone creates a copy of AudioVADParams
func (p *AudioVADParams) Clone() *AudioVADParams {
	clone := *p
	// Note: VADAnalyzer might need deep copy depending on implementation
	return &clone
}

// String returns string representation of AudioParams
func (p *AudioParams) String() string {
	return fmt.Sprintf("AudioParams{OutEnabled: %t, OutRate: %d, OutCh: %d, OutWidth: %d, OutChunks: %d, InEnabled: %t, InPartEnabled: %t, InRate: %d, InCh: %d, InWidth: %d}",
		p.AudioOutEnabled, p.AudioOutSampleRate, p.AudioOutChannels, p.AudioOutSampleWidth,
		p.AudioOut10msChunks, p.AudioInEnabled, p.AudioInParticipantEnabled,
		p.AudioInSampleRate, p.AudioInChannels, p.AudioInSampleWidth)
}

// String returns string representation of AudioVADParams
func (p *AudioVADParams) String() string {
	vadAnalyzerStr := "nil"
	if p.VADAnalyzer != nil {
		vadAnalyzerStr = fmt.Sprintf("%T", p.VADAnalyzer)
	}
	return fmt.Sprintf("AudioVADParams{%s, VADEnabled: %t, VADPassThrough: %t, VADAnalyzer: %s}",
		p.AudioParams.String(), p.VADEnabled, p.VADAudioPassthrough, vadAnalyzerStr)
}
