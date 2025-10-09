package common

import (
	"github.com/weedge/pipeline-go/pkg/frames"

	achatbot_frames "achatbot/pkg/types/frames"
)

// IVADAnalyzer 定义了语音活动检测（VAD）分析器的接口。
type IVADAnalyzer interface {
	// AnalyzeAudio 对输入的音频缓冲区进行分析，返回 VAD Frame
	AnalyzeAudio(buffer []byte) *achatbot_frames.VADStateAudioRawFrame

	// Reset 重置 VAD 的统计信息和模型状态。
	Reset()

	// GetSampleRate 返回 VAD 的采样率。
	GetSampleRate() int

	// Release 释放资源。
	Release()
}

// IVoiceConfidenceProvider 语音置信度提供者接口
type IVoiceConfidenceProvider interface {
	// IsActiveSpeech 判断当前音频是否是活跃的语音。
	IsActiveSpeech(audio []byte) bool

	// Release 释放资源。
	Release()

	// GetSampleInfo 返回语音置信度提供者的采样率信息和采样窗口大小。
	GetSampleInfo() (int, int)
}

// We'll use the standard net/http package for WebSocket support
// You may need to add the gorilla/websocket dependency or use standard library
// For now, we'll define a generic interface
// IWebSocketConn defines the interface for WebSocket connections
type IWebSocketConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type ITransportWriter interface {
	WriteRawAudio(data []byte) error

	WriteFrame(frame frames.Frame) error

	//WriteAnimationAudioFrame(frame *achatbot_frames.AnimationAudioRawFrame) error
	//SendMessage(frame *achatbot_frames.TransportMessageFrame) error
	//SendText(frame *frames.TextFrame) error
	//WriteImageFrame(frame *frames.ImageRawFrame) error
}
