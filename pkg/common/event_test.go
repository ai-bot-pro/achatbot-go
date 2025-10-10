package common

import (
	"testing"
	"time"
)

// TestNewEventHandlerManager 测试创建事件处理器管理器
func TestNewEventHandlerManager(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	if manager == nil {
		t.Error("Failed to create EventHandlerManager")
	}
	
	if manager.Name() == "" {
		t.Error("EventHandlerManager name should not be empty")
	}
	
}

// TestRegisterEventHandler 测试注册事件处理器
func TestRegisterEventHandler(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	
	err := manager.RegisterEventHandler("test_event")
	if err != nil {
		t.Errorf("Failed to register event handler: %v", err)
	}
	
	// 尝试重复注册应该失败
	err = manager.RegisterEventHandler("test_event")
	if err == nil {
		t.Error("Should not be able to register the same event handler twice")
	}
}

// TestAddEventHandler 测试添加事件处理函数
func TestAddEventHandler(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	
	// 先注册事件
	err := manager.RegisterEventHandler("test_event")
	if err != nil {
		t.Fatalf("Failed to register event handler: %v", err)
	}
	
	// 添加处理函数
	err = manager.AddEventHandler("test_event", func(manager *EventHandlerManager, data string) {
		// 处理逻辑
	})
	if err != nil {
		t.Errorf("Failed to add event handler: %v", err)
	}
	
	// 尝试添加到未注册的事件应该失败
	err = manager.AddEventHandler("nonexistent_event", func() {})
	if err == nil {
		t.Error("Should not be able to add handler to unregistered event")
	}
}

// TestAddEventHandlers 测试添加多个事件处理函数
func TestAddEventHandlers(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	
	// 先注册事件
	err := manager.RegisterEventHandler("test_event")
	if err != nil {
		t.Fatalf("Failed to register event handler: %v", err)
	}
	
	// 添加多个处理函数
	handlers := []interface{}{
		func(manager *EventHandlerManager, data string) {},
		func(manager *EventHandlerManager, data string) {},
	}
	
	err = manager.AddEventHandlers("test_event", handlers)
	if err != nil {
		t.Errorf("Failed to add event handlers: %v", err)
	}
}

// TestCallEventHandler 测试调用事件处理函数
func TestCallEventHandler(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	
	// 先注册事件
	err := manager.RegisterEventHandler("test_event")
	if err != nil {
		t.Fatalf("Failed to register event handler: %v", err)
	}
	
	// 记录是否被调用
	called := false
	
	// 添加处理函数
	err = manager.AddEventHandler("test_event", func(manager *EventHandlerManager, data string) {
		called = true
		if data != "test_data" {
			t.Errorf("Expected 'test_data', got '%s'", data)
		}
	})
	if err != nil {
		t.Fatalf("Failed to add event handler: %v", err)
	}
	
	// 调用事件处理函数
	err = manager.CallEventHandler("test_event", "test_data")
	if err != nil {
		t.Errorf("Failed to call event handler: %v", err)
	}
	
	if !called {
		t.Error("Event handler was not called")
	}
}

// TestAsyncCallEventHandler 测试异步调用事件处理函数
func TestAsyncCallEventHandler(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	
	// 先注册事件
	err := manager.RegisterEventHandler("test_event")
	if err != nil {
		t.Fatalf("Failed to register event handler: %v", err)
	}
	
	// 记录是否被调用
	called := false
	
	// 添加处理函数
	err = manager.AddEventHandler("test_event", func(manager *EventHandlerManager, data string) {
		time.Sleep(100 * time.Millisecond) // 模拟耗时操作
		called = true
	})
	if err != nil {
		t.Fatalf("Failed to add event handler: %v", err)
	}
	
	// 异步调用事件处理函数
	resultChan := manager.AsyncCallEventHandler("test_event", "test_data")
	
	// 等待结果
	select {
	case err := <-resultChan:
		if err != nil {
			t.Errorf("Failed to call event handler: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Async call timed out")
	}
	
	if !called {
		t.Error("Event handler was not called")
	}
}

// TestEventNames 测试获取事件名称列表
func TestEventNames(t *testing.T) {
	manager := NewEventHandlerManagerWithName("test")
	
	// 注册几个事件
	events := []string{"event1", "event2", "event3"}
	for _, event := range events {
		err := manager.RegisterEventHandler(event)
		if err != nil {
			t.Fatalf("Failed to register event handler: %v", err)
		}
	}
	
	// 获取事件名称列表
	names := manager.EventNames()
	if len(names) != len(events) {
		t.Errorf("Expected %d event names, got %d", len(events), len(names))
	}
}