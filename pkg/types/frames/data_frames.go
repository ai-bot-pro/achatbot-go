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

// SpriteFrame represents a sprite frame (collection of images)
type SpriteFrame struct {
	*pipelineframes.DataFrame
	Images []*pipelineframes.ImageRawFrame
}

// NewSpriteFrame creates a new SpriteFrame
func NewSpriteFrame(images []*pipelineframes.ImageRawFrame) *SpriteFrame {
	return &SpriteFrame{
		DataFrame: pipelineframes.NewDataFrameWithName("SpriteFrame"),
		Images:    images,
	}
}

// AnimationAudioRawFrame represents an animation audio frame
type AnimationAudioRawFrame struct {
	*pipelineframes.AudioRawFrame
	AnimationJSON string `json:"animation_json"`
	AvatarStatus  string `json:"avatar_status"`
}

// NewAnimationAudioRawFrame creates a new AnimationAudioRawFrame
func NewAnimationAudioRawFrame(audio []byte, sampleRate, numChannels, sampleWidth int, animationJSON, avatarStatus string) *AnimationAudioRawFrame {
	return &AnimationAudioRawFrame{
		AudioRawFrame: pipelineframes.NewAudioRawFrame(audio, sampleRate, numChannels, sampleWidth),
		AnimationJSON: animationJSON,
		AvatarStatus:  avatarStatus,
	}
}

// String implements string representation of AnimationAudioRawFrame
func (f *AnimationAudioRawFrame) String() string {
	return fmt.Sprintf("%s animation_json: %s avatar_status: %s",
		f.AudioRawFrame.String(),
		f.AnimationJSON,
		f.AvatarStatus)
}

// TransportMessageFrame represents a transport message frame
type TransportMessageFrame struct {
	*pipelineframes.DataFrame
	Message []byte
}

// NewTransportMessageFrame creates a new TransportMessageFrame
func NewTransportMessageFrame(message []byte) *TransportMessageFrame {
	return &TransportMessageFrame{
		DataFrame: pipelineframes.NewDataFrameWithName("TransportMessageFrame"),
		Message:   message,
	}
}

type PathAudioRawFrame struct {
	*pipelineframes.AudioRawFrame
	Path string `json:"path"`
}

// NewPathAudioRawFrame creates a new PathAudioRawFrame
func NewPathAudioRawFrame(audio []byte, sampleRate, numChannels, sampleWidth int, path string) *PathAudioRawFrame {
	return &PathAudioRawFrame{
		AudioRawFrame: pipelineframes.NewAudioRawFrame(audio, sampleRate, numChannels, sampleWidth),
		Path:          path,
	}
}

// String implements string representation of PathAudioRawFrame
func (f *PathAudioRawFrame) String() string {
	return fmt.Sprintf("%s path: %s", f.AudioRawFrame.String(), f.Path)
}
