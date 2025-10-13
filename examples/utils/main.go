package main

import (
	"fmt"
	"os"

	"achatbot/pkg/consts"
	"achatbot/pkg/utils"
)

func main() {
	fmt.Println("WAV 处理函数使用示例")
	fmt.Println("==================")

	// 示例1: 基本保存音频文件
	audioData := []byte("这是一个测试音频数据")
	filePath, err := utils.SaveAudioToFile(audioData, "test.wav")
	if err != nil {
		fmt.Printf("保存音频文件失败: %v\n", err)
		return
	}
	fmt.Printf("音频文件已保存到: %s\n", filePath)

	// 示例2: 使用自定义选项保存音频文件
	customFilePath, err := utils.SaveAudioToFile(
		audioData,
		"custom_test.wav",
		utils.WithAudioDir(consts.RECORDS_DIR),
		utils.WithChannels(2),
		utils.WithSampleWidth(2),
		utils.WithSampleRate(44100),
	)
	if err != nil {
		fmt.Printf("保存自定义音频文件失败: %v\n", err)
		return
	}
	fmt.Printf("自定义音频文件已保存到: %s\n", customFilePath)

	// 示例3: 读取音频文件
	readData, err := utils.ReadAudioFile(filePath)
	if err != nil {
		fmt.Printf("读取音频文件失败: %v\n", err)
		return
	}
	fmt.Printf("读取到的音频数据: %s\n", string(readData))

	// 示例4: 读取WAV文件的原始字节和采样率
	rawBytes, sampleRate, err := utils.ReadWAVToBytes(filePath)
	if err != nil {
		fmt.Printf("读取WAV文件信息失败: %v\n", err)
		return
	}
	fmt.Printf("读取到的原始字节长度: %d, 采样率: %d\n", len(rawBytes), sampleRate)

	// 清理示例文件
	os.Remove(filePath)
	os.Remove(customFilePath)
	fmt.Println("示例完成，临时文件已清理")
}
