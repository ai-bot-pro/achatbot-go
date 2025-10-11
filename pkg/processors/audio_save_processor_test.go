package processors

import (
	"os"
	"path/filepath"
	"testing"

	"achatbot/pkg/consts"

	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/processors"
)

func TestAudioSaveProcessor(t *testing.T) {
	// 创建测试目录
	testDir := filepath.Join(consts.RECORDS_DIR, "test_audio_save")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// 创建 AudioSaveProcessor
	processor := NewAudioSaveProcessor("test", testDir, true)

	// 创建测试音频帧
	audioData := []byte("test audio data")
	audioFrame := frames.NewAudioRawFrame(audioData, 16000, 1, 2)

	// 处理音频帧
	processor.ProcessFrame(audioFrame, processors.FrameDirectionDownstream)

	// 检查文件是否创建
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No audio file was created")
	}

	// 验证至少创建了一个文件
	foundWavFile := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".wav" {
			foundWavFile = true
			break
		}
	}

	if !foundWavFile {
		t.Error("No WAV file was created")
	}
}

func TestSaveAllAudioProcessor(t *testing.T) {
	// 创建测试目录
	testDir := filepath.Join(consts.RECORDS_DIR, "test_save_all")
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// 创建 SaveAllAudioProcessor
	processor := NewSaveAllAudioProcessor("test_all", testDir, 16000, 1, 2, 0)

	// 发送开始帧
	startFrame := frames.NewStartFrame()
	processor.ProcessFrame(startFrame, processors.FrameDirectionDownstream)

	// 创建测试音频帧
	audioData := []byte("test audio data for all")
	audioFrame := frames.NewAudioRawFrame(audioData, 16000, 1, 2)
	processor.ProcessFrame(audioFrame, processors.FrameDirectionDownstream)

	// 发送结束帧
	endFrame := frames.NewEndFrame()
	processor.ProcessFrame(endFrame, processors.FrameDirectionDownstream)

	// 检查文件是否创建
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("No audio file was created")
	}

	// 验证至少创建了一个文件
	foundWavFile := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".wav" {
			foundWavFile = true
			break
		}
	}

	if !foundWavFile {
		t.Error("No WAV file was created")
	}
}
