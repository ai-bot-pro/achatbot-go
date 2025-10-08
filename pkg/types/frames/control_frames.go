package frames

import pipelineframes "github.com/weedge/pipeline-go/pkg/frames"

// Emitted by VAD to indicate that a user has started speaking. This can be
// used for interruptions or other times when detecting that someone is
// speaking is more important than knowing what they're saying (as you will with a TranscriptionFrame)
type UserStartedSpeakingFrame struct {
	*pipelineframes.ControlFrame
}

func NewUserStartedSpeakingFrame() *UserStartedSpeakingFrame {
	return &UserStartedSpeakingFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("UserStartedSpeakingFrame"),
		},
	}
}

// UserStoppedSpeakingFrame is emitted by the VAD to indicate that a user stopped speaking.
type UserStoppedSpeakingFrame struct {
	*pipelineframes.ControlFrame
}

func NewUserStoppedSpeakingFrame() *UserStoppedSpeakingFrame {
	return &UserStoppedSpeakingFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("UserStoppedSpeakingFrame"),
		},
	}
}
