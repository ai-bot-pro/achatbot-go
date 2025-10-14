package main

import (
	"achatbot/pkg/common"
	"encoding/json"
	"fmt"
)

func printList(name string, list []map[string]any) {
	fmt.Printf("%s:\n", name)
	for i, item := range list {
		// 安全地访问map元素
		role, roleOk := item["role"]
		content, contentOk := item["content"]

		if roleOk && contentOk {
			fmt.Printf("  %d: role=%s, content=%s\n", i, role, content)
		} else {
			// 如果无法获取role和content，打印整个map
			fmt.Printf("  %d: %v\n", i, item)
		}
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== ChatHistory 示例 ===")

	// 1. 创建一个无限制的聊天历史
	fmt.Println("1. 无限制聊天历史:")
	unlimitedHistory := common.NewChatHistory(nil, nil, nil)
	unlimitedHistory.Append(map[string]any{"role": "user", "content": "Hello"})
	unlimitedHistory.Append(map[string]any{"role": "assistant", "content": "Hi there!"})
	unlimitedHistory.Append(map[string]any{"role": "user", "content": "How are you?"})
	unlimitedHistory.Append(map[string]any{"role": "assistant", "content": "I'm doing well, thank you!"})
	printList("无限制历史", unlimitedHistory.ToList())

	// 2. 创建一个有限制的聊天历史（大小为2）
	fmt.Println("2. 限制大小的聊天历史（大小=2）:")
	size := 2
	limitedHistory := common.NewChatHistory(&size, nil, nil)

	// 添加3对对话，应该只保留最后2对
	limitedHistory.Append(map[string]any{"role": "user", "content": "First question"})
	limitedHistory.Append(map[string]any{"role": "assistant", "content": "First answer"})
	limitedHistory.Append(map[string]any{"role": "user", "content": "Second question"})
	limitedHistory.Append(map[string]any{"role": "assistant", "content": "Second answer"})
	limitedHistory.Append(map[string]any{"role": "user", "content": "Third question"})
	limitedHistory.Append(map[string]any{"role": "assistant", "content": "Third answer"})
	printList("限制大小历史", limitedHistory.ToList())

	// 3. 创建带初始消息的聊天历史
	fmt.Println("3. 带初始消息的聊天历史:")
	initMsg := map[string]any{"role": "system", "content": "You are a helpful assistant"}
	initTools := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        "get_weather",
			"description": "Get weather information",
		},
	}
	historyWithInit := common.NewChatHistory(&size, initMsg, initTools)
	historyWithInit.Append(map[string]any{"role": "user", "content": "What's the weather like?"})
	historyWithInit.Append(map[string]any{"role": "assistant", "content": "It's sunny today!"})
	printList("带初始消息历史", historyWithInit.ToList())

	// 4. 演示JSON序列化/反序列化
	fmt.Println("4. JSON序列化/反序列化:")
	data, err := json.Marshal(limitedHistory)
	if err != nil {
		fmt.Printf("序列化错误: %v\n", err)
	} else {
		fmt.Printf("序列化结果: %s\n", string(data))

		// 反序列化
		var restoredHistory common.ChatHistory
		err = json.Unmarshal(data, &restoredHistory)
		if err != nil {
			fmt.Printf("反序列化错误: %v\n", err)
		} else {
			printList("反序列化后的历史", restoredHistory.ToList())
		}
	}

	// 5. 演示Pop操作
	fmt.Println("5. Pop操作演示:")
	history := common.NewChatHistory(nil, nil, nil)
	history.Append(map[string]any{"role": "user", "content": "Message 1"})
	history.Append(map[string]any{"role": "assistant", "content": "Response 1"})
	history.Append(map[string]any{"role": "user", "content": "Message 2"})
	history.Append(map[string]any{"role": "assistant", "content": "Response 2"})

	fmt.Println("Pop前:")
	printList("历史", history.ToList())

	history.Pop(-1) // 移除最后一个元素
	fmt.Println("Pop(-1)后:")
	printList("历史", history.ToList())

	history.Pop(0) // 移除第一个元素
	fmt.Println("Pop(0)后:")
	printList("历史", history.ToList())
}
