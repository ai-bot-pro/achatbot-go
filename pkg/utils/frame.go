package utils

import (
	"github.com/weedge/pipeline-go/pkg/frames"

	achatbot_frames "achatbot/pkg/types/frames"
)

func GetAudioRawFrame(frame frames.Frame) *frames.AudioRawFrame {
	switch f := frame.(type) {
	case *frames.AudioRawFrame:
		return f
	case *achatbot_frames.VADStateAudioRawFrame:
		return f.AudioRawFrame
	case *achatbot_frames.AnimationAudioRawFrame:
		return f.AudioRawFrame
	default:
		return nil
	}
}
