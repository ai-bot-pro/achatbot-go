package frames

import pipelineframes "github.com/weedge/pipeline-go/pkg/frames"

// Emitted by when the bot should be interrupted. This will mainly cause the
// same actions as if the user interrupted except that the
// UserStartedSpeakingFrame and UserStoppedSpeakingFrame won't be generated.
type BotInterruptionFrame struct {
	*pipelineframes.SystemFrame
}
