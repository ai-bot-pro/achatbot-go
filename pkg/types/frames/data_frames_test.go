package frames

import (
	"testing"
)

func TestAnimationAudioRawFrame(t *testing.T) {
	// 创建测试数据
	audioData := []byte("test audio data")
	sampleRate := 16000
	numChannels := 1
	sampleWidth := 2
	animationJSON := `{"animation": "smile", "duration": 2.5}`
	avatarStatus := "speaking"

	// 创建AnimationAudioRawFrame实例
	frame := NewAnimationAudioRawFrame(audioData, sampleRate, numChannels, sampleWidth, animationJSON, avatarStatus)

	// 验证AudioRawFrame字段是否正确设置
	if string(frame.Audio) != string(audioData) {
		t.Errorf("Expected audio data %s, got %s", string(audioData), string(frame.Audio))
	}

	if frame.SampleRate != sampleRate {
		t.Errorf("Expected sample rate %d, got %d", sampleRate, frame.SampleRate)
	}

	if frame.NumChannels != numChannels {
		t.Errorf("Expected num channels %d, got %d", numChannels, frame.NumChannels)
	}

	if frame.SampleWidth != sampleWidth {
		t.Errorf("Expected sample width %d, got %d", sampleWidth, frame.SampleWidth)
	}

	// 验证新添加的字段是否正确设置
	if frame.AnimationJSON != animationJSON {
		t.Errorf("Expected animation JSON %s, got %s", animationJSON, frame.AnimationJSON)
	}

	if frame.AvatarStatus != avatarStatus {
		t.Errorf("Expected avatar status %s, got %s", avatarStatus, frame.AvatarStatus)
	}

	// 验证String方法
	str := frame.String()
	if len(str) == 0 {
		t.Error("String method should return a non-empty string")
	}

	// 检查字符串是否包含关键信息
	expectedContains := []string{
		"animation_json: " + animationJSON,
		"avatar_status: " + avatarStatus,
	}

	for _, expected := range expectedContains {
		if !contains(str, expected) {
			t.Errorf("String method output should contain '%s', got: %s", expected, str)
		}
	}
}

// 辅助函数检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) != -1
}

// 辅助函数查找子字符串的位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestAnimationAudioRawFrameWithDefaults(t *testing.T) {
	// 测试使用默认值创建帧
	audioData := []byte{}
	sampleRate := 0
	numChannels := 0
	sampleWidth := 0
	animationJSON := ""
	avatarStatus := ""

	frame := NewAnimationAudioRawFrame(audioData, sampleRate, numChannels, sampleWidth, animationJSON, avatarStatus)

	// 验证默认值
	if len(frame.Audio) != 0 {
		t.Errorf("Expected empty audio data, got length %d", len(frame.Audio))
	}

	if frame.SampleRate != 0 {
		t.Errorf("Expected sample rate 0, got %d", frame.SampleRate)
	}

	if frame.AnimationJSON != "" {
		t.Errorf("Expected empty animation JSON, got %s", frame.AnimationJSON)
	}

	if frame.AvatarStatus != "" {
		t.Errorf("Expected empty avatar status, got %s", frame.AvatarStatus)
	}
}
func TestFunctionCallFrame(t *testing.T) {
	// 创建测试数据
	toolCallID := "call_1234567890"
	functionName := "get_current_weather"
	arguments := `{"location": "San Francisco, CA", "format": "celsius"}`
	index := 0

	// 创建FunctionCallFrame实例
	frame := NewFunctionCallFrame(toolCallID, functionName, arguments, index)

	// 验证字段是否正确设置
	if frame.ToolCallID != toolCallID {
		t.Errorf("Expected toolCallID %s, got %s", toolCallID, frame.ToolCallID)
	}

	if frame.FunctionName != functionName {
		t.Errorf("Expected functionName %s, got %s", functionName, frame.FunctionName)
	}

	if frame.Arguments != arguments {
		t.Errorf("Expected arguments %s, got %s", arguments, frame.Arguments)
	}

	if frame.Index != index {
		t.Errorf("Expected index %d, got %d", index, frame.Index)
	}

	if frame.Type != "function" {
		t.Errorf("Expected type 'function', got %s", frame.Type)
	}

	// 验证String方法
	str := frame.String()
	if len(str) == 0 {
		t.Error("String method should return a non-empty string")
	}

	// 检查字符串是否包含关键信息
	expectedContains := []string{
		"function_name: " + functionName,
		"tool_call_id: " + toolCallID,
		"arguments: " + arguments,
		"index: 0",
		"type: function",
	}

	for _, expected := range expectedContains {
		if !contains(str, expected) {
			t.Errorf("String method output should contain '%s', got: %s", expected, str)
		}
	}
}

func TestFunctionCallFrameArgumentsDict(t *testing.T) {
	// 创建测试数据
	toolCallID := "call_1234567890"
	functionName := "get_current_weather"
	arguments := `{"location": "San Francisco, CA", "format": "celsius"}`
	index := 0

	// 创建FunctionCallFrame实例
	frame := NewFunctionCallFrame(toolCallID, functionName, arguments, index)

	// 测试ArgumentsDict方法
	args, err := frame.ArgumentsDict()
	if err != nil {
		t.Errorf("ArgumentsDict should not return an error, got: %v", err)
	}

	// 验证解析结果
	if location, ok := args["location"]; !ok || location != "San Francisco, CA" {
		t.Errorf("Expected location 'San Francisco, CA', got %v", location)
	}

	if format, ok := args["format"]; !ok || format != "celsius" {
		t.Errorf("Expected format 'celsius', got %v", format)
	}
}

func TestFunctionCallFrameArgumentsDictInvalidJSON(t *testing.T) {
	// 创建测试数据，带有无效的JSON
	toolCallID := "call_1234567890"
	functionName := "get_current_weather"
	arguments := `{"location": "San Francisco, CA", "format":}` // 无效的JSON
	index := 0

	// 创建FunctionCallFrame实例
	frame := NewFunctionCallFrame(toolCallID, functionName, arguments, index)

	// 测试ArgumentsDict方法应该返回错误
	_, err := frame.ArgumentsDict()
	if err == nil {
		t.Error("ArgumentsDict should return an error for invalid JSON")
	}
}
