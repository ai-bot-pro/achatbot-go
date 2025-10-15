// Package aggregators provides functionality for aggregating audio response frames.
package aggregators

import (
	"reflect"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/processors"

	achatbot_frames "achatbot/pkg/types/frames"
)

// AudioResponseAggregator 聚合音频响应帧的处理器
// 该处理器用于聚合音频帧，直到收到开始和结束帧为止
type AudioResponseAggregator struct {
	*processors.FrameProcessor
	startFrame                  reflect.Type
	curFrame                    frames.Frame
	endFrame                    reflect.Type
	accumulatorFrameType        reflect.Type
	interimAccumulatorFrameType reflect.Type

	aggregation        []byte
	aggregating        bool
	seenStartFrame     bool
	seenEndFrame       bool
	seenInterimResults bool
}

// NewAudioResponseAggregator 创建一个新的 AudioResponseAggregator 实例
// 参数:
//   - startFrame: 开始帧类型
//   - endFrame: 结束帧类型
//   - accumulatorFrameType: 累加帧类型
//   - interimAccumulatorFrameType: 中间累加帧类型
func NewAudioResponseAggregator(
	startFrame reflect.Type,
	endFrame reflect.Type,
	accumulatorFrameType reflect.Type,
	interimAccumulatorFrameType reflect.Type,
) *AudioResponseAggregator {
	aggregator := &AudioResponseAggregator{
		FrameProcessor:              processors.NewFrameProcessor("AudioResponseAggregator"),
		startFrame:                  startFrame,
		endFrame:                    endFrame,
		accumulatorFrameType:        accumulatorFrameType,
		interimAccumulatorFrameType: interimAccumulatorFrameType,
	}

	aggregator.reset()

	return aggregator
}

// NewAudioResponseAggregatorWithAccumulate 创建一个不需要中间累加帧的的 AudioResponseAggregator 实例
func NewAudioResponseAggregatorWithAccumulate(
	startFrame reflect.Type,
	endFrame reflect.Type,
	accumulatorFrameType reflect.Type,
) *AudioResponseAggregator {
	return NewAudioResponseAggregator(startFrame, endFrame, accumulatorFrameType, nil)
}

// reset 重置聚合状态
func (a *AudioResponseAggregator) reset() {
	a.aggregation = make([]byte, 0)
	a.aggregating = false
	a.seenStartFrame = false
	a.seenEndFrame = false
	a.seenInterimResults = false
	a.curFrame = nil
}

// ProcessFrame 处理帧
func (a *AudioResponseAggregator) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	// 调用父类方法
	a.FrameProcessor.ProcessFrame(frame, direction)

	sendAggregation := false

	// 检查帧类型并处理
	if a.isFrameType(frame, a.startFrame) {
		a.aggregating = true
		a.seenStartFrame = true
		a.seenEndFrame = false
		a.seenInterimResults = false

		a.PushFrame(frame, direction)
	} else if a.isFrameType(frame, a.endFrame) {
		a.seenEndFrame = true
		a.seenStartFrame = false

		// 我们可能已经收到了结束帧，但我们可能仍在聚合（即我们看到了中间结果但不是最终音频）
		a.aggregating = a.seenInterimResults || len(a.aggregation) == 0

		// 如果我们不再聚合（即没有更多中间结果接收到），则发送聚合
		sendAggregation = !a.aggregating
		a.PushFrame(frame, direction)
	} else if a.isFrameType(frame, a.accumulatorFrameType) {
		var curAudioFrame *frames.AudioRawFrame
		switch f := frame.(type) {
		case *frames.AudioRawFrame:
			curAudioFrame = f
		case *achatbot_frames.VADStateAudioRawFrame:
			curAudioFrame = f.AudioRawFrame
		case *achatbot_frames.AnimationAudioRawFrame:
			curAudioFrame = f.AudioRawFrame
		default:
			logger.Warnf("Frame is %T don't support aggregate", frame)
		}
		if a.aggregating && curAudioFrame != nil {
			a.curFrame = frame
			a.aggregation = append(a.aggregation, curAudioFrame.Audio...)

			// 我们已经收到了一个完整的句子，所以如果我们已经看到了结束帧并且我们仍在聚合，
			// 这意味着我们应该发送聚合。
			sendAggregation = a.seenEndFrame
		}

		// 我们刚刚得到了最终结果，所以让我们重置中间结果。
		a.seenInterimResults = false
	} else if a.interimAccumulatorFrameType != nil && a.isFrameType(frame, a.interimAccumulatorFrameType) {
		a.seenInterimResults = true
	} else {
		a.PushFrame(frame, direction)
	}

	if sendAggregation {
		a.pushAggregation(direction)
	}
}

// isFrameType 检查帧是否为指定类型的 Frame
func (a *AudioResponseAggregator) isFrameType(frame frames.Frame, frameType reflect.Type) bool {
	if frame == nil || frameType == nil {
		return false
	}

	return reflect.TypeOf(frame) == frameType
}

// pushAggregation 推送聚合结果
func (a *AudioResponseAggregator) pushAggregation(direction processors.FrameDirection) {
	// 确保我们有聚合数据
	if len(a.aggregation) > 0 {
		var frame *frames.AudioRawFrame

		if a.curFrame == nil {
			logger.Warn("curFrame is nil, don't to push aggregation")
			return
		}

		// 如果有当前音频帧，使用其属性
		switch f := a.curFrame.(type) {
		case *frames.AudioRawFrame:
			frame = f
		case *achatbot_frames.VADStateAudioRawFrame:
			frame = f.AudioRawFrame
		case *achatbot_frames.AnimationAudioRawFrame:
			frame = f.AudioRawFrame
		default:
			logger.Warnf("Frame is %T don't support aggregate to push", a.curFrame)
			return
		}
		frame.Audio = make([]byte, len(a.aggregation))
		copy(frame.Audio, a.aggregation)

		a.PushFrame(frame, direction)
		a.reset()
	}
}
