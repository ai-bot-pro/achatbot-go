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

// BotSpeakingFrame indicates that the bot is currently speaking
type BotSpeakingFrame struct {
	*pipelineframes.ControlFrame
}

// NewBotSpeakingFrame creates a new BotSpeakingFrame
func NewBotSpeakingFrame() *BotSpeakingFrame {
	return &BotSpeakingFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("BotSpeakingFrame"),
		},
	}
}

// BotStartedSpeakingFrame indicates that the bot has started speaking
type BotStartedSpeakingFrame struct {
	*pipelineframes.ControlFrame
}

// NewBotStartedSpeakingFrame creates a new BotStartedSpeakingFrame
func NewBotStartedSpeakingFrame() *BotStartedSpeakingFrame {
	return &BotStartedSpeakingFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("BotStartedSpeakingFrame"),
		},
	}
}

// BotStoppedSpeakingFrame indicates that the bot has stopped speaking
type BotStoppedSpeakingFrame struct {
	*pipelineframes.ControlFrame
}

// NewBotStoppedSpeakingFrame creates a new BotStoppedSpeakingFrame
func NewBotStoppedSpeakingFrame() *BotStoppedSpeakingFrame {
	return &BotStoppedSpeakingFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("BotStoppedSpeakingFrame"),
		},
	}
}

// TTSStartedFrame indicates that TTS has started
type TTSStartedFrame struct {
	*pipelineframes.ControlFrame
}

// NewTTSStartedFrame creates a new TTSStartedFrame
func NewTTSStartedFrame() *TTSStartedFrame {
	return &TTSStartedFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("TTSStartedFrame"),
		},
	}
}

// TTSStoppedFrame indicates that TTS has stopped
type TTSStoppedFrame struct {
	*pipelineframes.ControlFrame
}

// NewTTSStoppedFrame creates a new TTSStoppedFrame
func NewTTSStoppedFrame() *TTSStoppedFrame {
	return &TTSStoppedFrame{
		ControlFrame: &pipelineframes.ControlFrame{
			BaseFrame: pipelineframes.NewBaseFrameWithName("TTSStoppedFrame"),
		},
	}
}
