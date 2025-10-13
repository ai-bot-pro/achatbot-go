package main

import (
	"fmt"
	"time"

	"achatbot/pkg/common"
)

func main() {
	// 创建事件处理器管理器
	manager := common.NewEventHandlerManagerWithName("MyEventManager")

	fmt.Printf("Created event manager: %s \n", manager.Name())

	// 注册事件
	err := manager.RegisterEventHandler("user_login")
	if err != nil {
		fmt.Printf("Failed to register event: %v\n", err)
		return
	}

	err = manager.RegisterEventHandler("user_logout")
	if err != nil {
		fmt.Printf("Failed to register event: %v\n", err)
		return
	}

	// 添加事件处理函数
	err = manager.AddEventHandler("user_login", func(manager *common.EventHandlerManager, username string) {
		fmt.Printf("User %s logged in at %s\n", username, time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Printf("Failed to add event handler: %v\n", err)
		return
	}

	// 添加多个事件处理函数
	logoutHandlers := []interface{}{
		func(manager *common.EventHandlerManager, username string) {
			fmt.Printf("Logging logout event for user %s\n", username)
		},
		func(manager *common.EventHandlerManager, username string) {
			fmt.Printf("Sending notification for user %s logout\n", username)
		},
	}

	err = manager.AddEventHandlers("user_logout", logoutHandlers)
	if err != nil {
		fmt.Printf("Failed to add event handlers: %v\n", err)
		return
	}

	// 显示已注册的事件
	fmt.Println("Registered events:", manager.EventNames())

	// 触发事件
	fmt.Println("\nTriggering user_login event:")
	err = manager.CallEventHandler("user_login", "Alice")
	if err != nil {
		fmt.Printf("Failed to call event handler: %v\n", err)
	}

	fmt.Println("\nTriggering user_logout event:")
	err = manager.CallEventHandler("user_logout", "Alice")
	if err != nil {
		fmt.Printf("Failed to call event handler: %v\n", err)
	}

	// 异步触发事件
	fmt.Println("\nTriggering async user_login event:")
	resultChan := manager.AsyncCallEventHandler("user_login", "Bob")

	// 等待异步操作完成
	err = <-resultChan
	if err != nil {
		fmt.Printf("Failed to call event handler: %v\n", err)
	} else {
		fmt.Println("Async event handled successfully")
	}
}
