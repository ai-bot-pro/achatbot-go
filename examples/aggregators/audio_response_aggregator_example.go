package main

import (
	"fmt"
	"reflect"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"

	"achatbot/pkg/processors/aggregators"
)

// ExampleFrameProcessor 示例帧处理器
type ExampleFrameProcessor struct {
	*processors.FrameProcessor
}

// NewExampleFrameProcessor 创建一个新的示例帧处理器
func NewExampleFrameProcessor(name string) *ExampleFrameProcessor {
	return &ExampleFrameProcessor{
		FrameProcessor: processors.NewFrameProcessor(name),
	}
}

// ProcessFrame 处理帧
func (p *ExampleFrameProcessor) ProcessFrame(frame frames.Frame, direction processors.FrameDirection) {
	fmt.Printf("[%s] Received frame: %s\n", p.Name(), frame.Name())

	// 如果是聚合的音频帧，打印其内容
	if audioFrame, ok := frame.(*frames.AudioRawFrame); ok {
		fmt.Printf("[%s] Aggregated audio length: %d bytes\n", p.Name(), len(audioFrame.Audio))
		fmt.Printf("[%s] Aggregated audio content: %s\n", p.Name(), string(audioFrame.Audio))
	}
}

func main() {
	logger.InitLoggerWithConfig(logger.NewDefaultLoggerConfig())
	fmt.Println("Audio Response Aggregator Example")

	// 创建开始帧和结束帧
	startFrame := frames.NewStartFrame()
	endFrame := frames.NewEndFrame()

	// 创建音频响应聚合器
	aggregator := aggregators.NewAudioResponseAggregator(
		reflect.TypeOf(&frames.StartFrame{}),
		reflect.TypeOf(&frames.EndFrame{}),
		reflect.TypeOf(&frames.AudioRawFrame{}), // 累加帧类型
		nil,                                     // 中间累加帧类型（可选）
	)

	// 创建示例下游处理器
	downstreamProcessor := NewExampleFrameProcessor("DownstreamProcessor")

	// 创建一个简单的管道
	p := pipeline.NewPipelineWithVerbose(
		[]processors.IFrameProcessor{
			aggregator,
			downstreamProcessor,
		},
		nil, nil,
		true,
	)

	fmt.Println("\nProcessing frames...")

	// 发送开始帧
	fmt.Println("\n1. Sending StartFrame")
	p.ProcessFrame(startFrame, processors.FrameDirectionDownstream)

	// 发送几个音频帧
	fmt.Println("\n2. Sending AudioRawFrames")
	audioFrame1 := frames.NewAudioRawFrame([]byte("Hello "), 16000, 1, 2)
	p.ProcessFrame(audioFrame1, processors.FrameDirectionDownstream)

	audioFrame2 := frames.NewAudioRawFrame([]byte("World! "), 16000, 1, 2)
	p.ProcessFrame(audioFrame2, processors.FrameDirectionDownstream)

	audioFrame3 := frames.NewAudioRawFrame([]byte("This is "), 16000, 1, 2)
	p.ProcessFrame(audioFrame3, processors.FrameDirectionDownstream)

	audioFrame4 := frames.NewAudioRawFrame([]byte("a test."), 16000, 1, 2)
	p.ProcessFrame(audioFrame4, processors.FrameDirectionDownstream)

	// 发送结束帧，这将触发聚合
	fmt.Println("\n3. Sending EndFrame (will trigger aggregation)")
	p.ProcessFrame(endFrame, processors.FrameDirectionDownstream)

	fmt.Println("\nExample completed.")
}
