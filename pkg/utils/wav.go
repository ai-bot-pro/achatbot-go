package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"achatbot/pkg/consts"
)

// SaveAudioToFile 将音频数据保存为 WAV 文件
// 参数:
//   - audioData: 音频数据字节
//   - fileName: 文件名
//   - audioDir: 音频目录（可选，默认为 RECORDS_DIR）
//   - channels: 声道数（可选，默认为 1）
//   - sampleWidth: 采样宽度（可选，默认为 2）
//   - sampleRate: 采样率（可选，默认为 16000）
//
// 返回:
//   - string: 文件路径
//   - error: 错误信息
func SaveAudioToFile(audioData []byte, fileName string, opts ...WAVOption) (string, error) {
	// 设置默认选项
	options := &WAVOptions{
		AudioDir:    consts.RECORDS_DIR,
		Channels:    consts.DefaultChannels,
		SampleWidth: consts.DefaultSampleWidth,
		SampleRate:  consts.DefaultRate,
	}

	// 应用传入选项
	for _, opt := range opts {
		opt(options)
	}

	// 创建目录
	if err := os.MkdirAll(options.AudioDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", options.AudioDir, err)
	}

	// 构建文件路径
	filePath := filepath.Join(options.AudioDir, fileName)

	// 创建 WAV 文件
	wavFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer wavFile.Close()

	// 写入 WAV 头部
	header := createWAVHeader(options.Channels, options.SampleWidth, options.SampleRate, len(audioData))
	if _, err := wavFile.Write(header); err != nil {
		return "", fmt.Errorf("failed to write WAV header: %w", err)
	}

	// 写入音频数据
	if _, err := wavFile.Write(audioData); err != nil {
		return "", fmt.Errorf("failed to write audio data: %w", err)
	}

	return filePath, nil
}

// ReadAudioFile 从 WAV 文件读取音频数据
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - []byte: 音频数据
//   - error: 错误信息
func ReadAudioFile(filePath string) ([]byte, error) {
	// 打开文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 检查文件大小
	if len(data) < 44 {
		return nil, fmt.Errorf("invalid WAV file: file too small")
	}

	// 返回音频数据部分（跳过 44 字节的头部）
	return data[44:], nil
}

// ReadWAVToBytes 读取 WAV 文件的原始字节和采样率
// 参数:
//   - filePath: 文件路径
//
// 返回:
//   - []byte: 音频数据
//   - int: 采样率
//   - error: 错误信息
func ReadWAVToBytes(filePath string) ([]byte, int, error) {
	// 打开文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 检查文件大小
	if len(data) < 44 {
		return nil, 0, fmt.Errorf("invalid WAV file: file too small")
	}

	// 解析 WAV 头部获取采样率
	sampleRate := parseSampleRateFromWAVHeader(data[:44])

	// 返回音频数据部分（跳过 44 字节的头部）和采样率
	return data[44:], sampleRate, nil
}

// WAVOptions WAV 文件选项
type WAVOptions struct {
	AudioDir    string
	Channels    int
	SampleWidth int
	SampleRate  int
}

// WAVOption WAV 选项函数类型
type WAVOption func(*WAVOptions)

// WithAudioDir 设置音频目录
func WithAudioDir(dir string) WAVOption {
	return func(o *WAVOptions) {
		o.AudioDir = dir
	}
}

// WithChannels 设置声道数
func WithChannels(channels int) WAVOption {
	return func(o *WAVOptions) {
		o.Channels = channels
	}
}

// WithSampleWidth 设置采样宽度
func WithSampleWidth(width int) WAVOption {
	return func(o *WAVOptions) {
		o.SampleWidth = width
	}
}

// WithSampleRate 设置采样率
func WithSampleRate(rate int) WAVOption {
	return func(o *WAVOptions) {
		o.SampleRate = rate
	}
}

// createWAVHeader 创建 WAV 文件头部
func createWAVHeader(channels, sampleWidth, sampleRate, dataLen int) []byte {
	header := make([]byte, 44)

	// RIFF header
	copy(header[0:4], "RIFF")
	// 文件长度（整个文件大小-8）
	fileSize := 36 + dataLen
	header[4] = byte(fileSize)
	header[5] = byte(fileSize >> 8)
	header[6] = byte(fileSize >> 16)
	header[7] = byte(fileSize >> 24)
	// WAVE header
	copy(header[8:12], "WAVE")
	// fmt chunk
	copy(header[12:16], "fmt ")
	// fmt chunk size (16 for PCM)
	header[16] = 16
	header[17] = 0
	header[18] = 0
	header[19] = 0
	// Audio format (1 for PCM)
	header[20] = 1
	header[21] = 0
	// Number of channels
	header[22] = byte(channels)
	header[23] = byte(channels >> 8)
	// Sample rate
	header[24] = byte(sampleRate)
	header[25] = byte(sampleRate >> 8)
	header[26] = byte(sampleRate >> 16)
	header[27] = byte(sampleRate >> 24)
	// Byte rate (SampleRate * NumChannels * BitsPerSample/8)
	byteRate := sampleRate * channels * sampleWidth
	header[28] = byte(byteRate)
	header[29] = byte(byteRate >> 8)
	header[30] = byte(byteRate >> 16)
	header[31] = byte(byteRate >> 24)
	// Block align (NumChannels * BitsPerSample/8)
	blockAlign := channels * sampleWidth
	header[32] = byte(blockAlign)
	header[33] = byte(blockAlign >> 8)
	// Bits per sample
	bitsPerSample := sampleWidth * 8
	header[34] = byte(bitsPerSample)
	header[35] = byte(bitsPerSample >> 8)
	// data chunk
	copy(header[36:40], "data")
	// Data chunk size
	header[40] = byte(dataLen)
	header[41] = byte(dataLen >> 8)
	header[42] = byte(dataLen >> 16)
	header[43] = byte(dataLen >> 24)

	return header
}

// parseSampleRateFromWAVHeader 从 WAV 头部解析采样率
func parseSampleRateFromWAVHeader(header []byte) int {
	// 采样率位于第 24-27 字节
	return int(header[24]) |
		int(header[25])<<8 |
		int(header[26])<<16 |
		int(header[27])<<24
}
