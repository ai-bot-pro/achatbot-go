package frames

import (
	"fmt"

	pipelineframes "github.com/weedge/pipeline-go/pkg/frames"

	"achatbot/pkg/types"
)

// VADStateAudioRawFrame represents a VAD state audio frame
type VADStateAudioRawFrame struct {
	*pipelineframes.AudioRawFrame

	State    types.VADState `json:"state"`
	SpeechID int            `json:"speech_id"`
	IsFinal  bool           `json:"is_final"`
	StartAtS float64        `json:"start_at_s"`
	CurAtS   float64        `json:"cur_at_s"`
	EndAtS   float64        `json:"end_at_s"`
}

// NewVADStateAudioRawFrame creates a new VADStateAudioRawFrame with default values
func NewVADStateAudioRawFrame() *VADStateAudioRawFrame {
	return &VADStateAudioRawFrame{
		State:    types.Quiet,
		SpeechID: 0,
		IsFinal:  false,
		StartAtS: 0.0,
		CurAtS:   0.0,
		EndAtS:   0.0,
	}
}

// String implements string representation of VADStateAudioRawFrame
func (f *VADStateAudioRawFrame) String() string {
	return fmt.Sprintf("%s (state: %v speech_id: %d is_final: %t speech_id: %d start_at_s: %.2f cur_at_s: %.2f end_at_s: %.2f)",
		f.AudioRawFrame.String(),
		f.State,
		f.SpeechID,
		f.IsFinal,
		f.SpeechID,
		f.StartAtS,
		f.CurAtS,
		f.EndAtS,
	)
}
