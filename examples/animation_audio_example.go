package main

import (
	"log"

	"achatbot/pkg/types/frames"
)

func main() {
	// 创建测试数据
	audioData := []byte("test audio data")
	sampleRate := 16000
	numChannels := 1
	sampleWidth := 2
	animationJSON := `{"animation": "smile", "duration": 2.5}`
	avatarStatus := "speaking"

	// 创建AnimationAudioRawFrame实例
	frame := frames.NewAnimationAudioRawFrame(
		audioData,
		sampleRate,
		numChannels,
		sampleWidth,
		animationJSON,
		avatarStatus,
	)

	// 输出帧的信息
	log.Printf("AnimationAudioRawFrame created: %s", frame.String())

	// 验证字段
	log.Printf("Audio data length: %d", len(frame.Audio))
	log.Printf("Sample rate: %d", frame.SampleRate)
	log.Printf("Number of channels: %d", frame.NumChannels)
	log.Printf("Sample width: %d", frame.SampleWidth)
	log.Printf("Animation JSON: %s", frame.AnimationJSON)
	log.Printf("Avatar status: %s", frame.AvatarStatus)
}
