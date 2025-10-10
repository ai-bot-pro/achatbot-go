package common

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/weedge/pipeline-go/pkg"
)

// EventHandlerManager 事件处理器管理器
type EventHandlerManager struct {
	name          string
	eventHandlers map[string][]any
	mutex         sync.RWMutex
}

// NewEventHandlerManagerWithName 创建一个新的事件处理器管理器
func NewEventHandlerManagerWithName(name string) *EventHandlerManager {
	// 增加全局对象计数器

	manager := &EventHandlerManager{
		name:          name,
		eventHandlers: make(map[string][]any),
	}

	manager.name = fmt.Sprintf("%s#%d", name, pkg.CountForType(name))

	return manager
}

// Name 获取管理器名称
func (e *EventHandlerManager) Name() string {
	return e.name
}

// EventNames 获取所有已注册的事件名称
func (e *EventHandlerManager) EventNames() []string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	names := make([]string, 0, len(e.eventHandlers))
	for name := range e.eventHandlers {
		names = append(names, name)
	}
	return names
}

// RegisterEventHandler 注册一个事件处理器
func (e *EventHandlerManager) RegisterEventHandler(eventName string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if _, exists := e.eventHandlers[eventName]; exists {
		return fmt.Errorf("event handler %s already registered", eventName)
	}

	e.eventHandlers[eventName] = make([]any, 0)
	return nil
}

// AddEventHandler 添加单个事件处理函数
func (e *EventHandlerManager) AddEventHandler(eventName string, handler any) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	handlers, exists := e.eventHandlers[eventName]
	if !exists {
		return fmt.Errorf("event handler %s not registered", eventName)
	}

	// 验证handler是否为函数
	handlerValue := reflect.ValueOf(handler)
	if handlerValue.Kind() != reflect.Func {
		return fmt.Errorf("handler is not a function")
	}

	e.eventHandlers[eventName] = append(handlers, handler)
	return nil
}

// AddEventHandlers 添加多个事件处理函数
func (e *EventHandlerManager) AddEventHandlers(eventName string, handlers []any) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	existingHandlers, exists := e.eventHandlers[eventName]
	if !exists {
		return fmt.Errorf("event handler %s not registered", eventName)
	}

	// 验证所有handlers是否为函数
	for _, handler := range handlers {
		handlerValue := reflect.ValueOf(handler)
		if handlerValue.Kind() != reflect.Func {
			return fmt.Errorf("handler is not a function")
		}
	}

	e.eventHandlers[eventName] = append(existingHandlers, handlers...)
	return nil
}

// CallEventHandler 调用指定事件的所有处理函数
func (e *EventHandlerManager) CallEventHandler(eventName string, args ...any) error {
	e.mutex.RLock()
	handlers, exists := e.eventHandlers[eventName]
	e.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("event handler %s not registered", eventName)
	}

	for _, handler := range handlers {
		// 检查handler是否为函数类型
		handlerValue := reflect.ValueOf(handler)
		if handlerValue.Kind() != reflect.Func {
			log.Printf("Handler is not a function in event %s", eventName)
			continue
		}

		// 准备参数
		argValues := make([]reflect.Value, len(args)+1) // +1 是因为第一个参数是EventHandlerManager本身
		argValues[0] = reflect.ValueOf(e)

		for i, arg := range args {
			argValues[i+1] = reflect.ValueOf(arg)
		}

		// 调用函数
		tryCallHandler(handlerValue, argValues, eventName)
	}

	return nil
}

// AsyncCallEventHandler 异步调用指定事件的所有处理函数
func (e *EventHandlerManager) AsyncCallEventHandler(eventName string, args ...any) <-chan error {
	resultChan := make(chan error, 1)

	go func() {
		resultChan <- e.CallEventHandler(eventName, args...)
		close(resultChan)
	}()

	return resultChan
}

// tryCallHandler 尝试调用处理函数并捕获异常
func tryCallHandler(handlerValue reflect.Value, argValues []reflect.Value, eventName string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Exception in event handler %s: %v", eventName, r)
		}
	}()

	handlerValue.Call(argValues)
}
