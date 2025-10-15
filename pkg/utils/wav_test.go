package utils

import (
	"os"
	"path/filepath"
	"testing"

	"achatbot/pkg/consts"

	"github.com/stretchr/testify/assert"
)

func TestSaveAudioToFile(t *testing.T) {
	// 准备测试数据
	audioData := []byte("test audio data")
	fileName := "test.wav"
	testDir := filepath.Join(consts.RECORDS_DIR, "test")

	// 清理之前的测试文件
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// 测试基本功能
	filePath, err := SaveAudioToFile(audioData, fileName, WithAudioDir(testDir))
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testDir, fileName), filePath)

	// 检查文件是否存在
	_, err = os.Stat(filePath)
	assert.NoError(t, err)

	// 检查文件内容
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Len(t, data, 44+len(audioData)) // 44字节头部 + 音频数据

	// 验证WAV头部
	assert.Equal(t, "RIFF", string(data[0:4]))
	assert.Equal(t, "WAVE", string(data[8:12]))
	assert.Equal(t, "fmt ", string(data[12:16]))
	assert.Equal(t, "data", string(data[36:40]))
}

func TestSaveAudioToFileWithOptions(t *testing.T) {
	// 准备测试数据
	audioData := []byte("test audio data for reading with options")
	fileName := "test_options.wav"
	testDir := filepath.Join(consts.RECORDS_DIR, "test_options")

	// 清理之前的测试文件
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// 使用自定义选项测试
	filePath, err := SaveAudioToFile(
		audioData,
		fileName,
		WithAudioDir(testDir),
		WithChannels(2),
		WithSampleWidth(2),
		WithSampleRate(44100),
	)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(testDir, fileName), filePath)

	// 检查文件内容
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	// 验证自定义参数是否正确写入头部
	// 验证声道数 (22-23字节)
	assert.Equal(t, byte(2), data[22])
	assert.Equal(t, byte(0), data[23])

	// 验证采样率 (24-27字节)
	assert.Equal(t, byte(0x44), data[24]) // 44100 & 0xFF
	assert.Equal(t, byte(0xac), data[25]) // (44100 >> 8) & 0xFF
	assert.Equal(t, byte(0x0), data[26])  // (44100 >> 16) & 0xFF
	assert.Equal(t, byte(0x0), data[27])  // (44100 >> 24) & 0xFF
}

func TestReadAudioFile(t *testing.T) {
	// 准备测试数据
	audioData := []byte("test audio data for reading")
	fileName := "test_read.wav"
	testDir := filepath.Join(consts.RECORDS_DIR, "test_read")

	// 清理之前的测试文件
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// 先创建一个WAV文件
	filePath, err := SaveAudioToFile(audioData, fileName, WithAudioDir(testDir))
	assert.NoError(t, err)

	// 测试读取功能
	readData, err := ReadAudioFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, audioData, readData)

	// 测试读取不存在的文件
	_, err = ReadAudioFile(filepath.Join(testDir, "nonexistent.wav"))
	assert.Error(t, err)
}

func TestReadWAVToBytes(t *testing.T) {
	// 准备测试数据
	audioData := []byte("test audio data for reading with sample rate")
	fileName := "test_read_with_rate.wav"
	testDir := filepath.Join(consts.RECORDS_DIR, "test_read_rate")

	// 清理之前的测试文件
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// 先创建一个WAV文件
	sampleRate := 22050
	filePath, err := SaveAudioToFile(
		audioData,
		fileName,
		WithAudioDir(testDir),
		WithSampleRate(sampleRate),
	)
	assert.NoError(t, err)

	// 测试读取功能
	readData, rate, err := ReadWAVToBytes(filePath)
	assert.NoError(t, err)
	assert.Equal(t, audioData, readData)
	assert.Equal(t, sampleRate, rate)

	// 测试读取不存在的文件
	_, _, err = ReadWAVToBytes(filepath.Join(testDir, "nonexistent.wav"))
	assert.Error(t, err)
}

func TestCreateWAVHeader(t *testing.T) {
	// 测试默认参数的头部创建
	header := createWAVHeader(1, 2, 16000, 100)
	assert.Len(t, header, 44)
	assert.Equal(t, "RIFF", string(header[0:4]))
	assert.Equal(t, "WAVE", string(header[8:12]))
	assert.Equal(t, "fmt ", string(header[12:16]))
	assert.Equal(t, "data", string(header[36:40]))

	// 验证声道数
	assert.Equal(t, byte(1), header[22])
	assert.Equal(t, byte(0), header[23])

	// 验证采样率 (16000 = 0x3E80)
	assert.Equal(t, byte(0x80), header[24]) // 16000的低位
	assert.Equal(t, byte(0x3E), header[25]) // 16000的次低位
	assert.Equal(t, byte(0x0), header[26])  // 16000的次高位
	assert.Equal(t, byte(0x0), header[27])  // 16000的高位
}

func TestParseSampleRateFromWAVHeader(t *testing.T) {
	// 创建一个测试头部
	header := createWAVHeader(1, 2, 44100, 100)

	// 解析采样率
	rate := parseSampleRateFromWAVHeader(header)
	assert.Equal(t, 44100, rate)
}
