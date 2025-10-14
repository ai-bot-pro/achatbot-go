package common

import (
	"encoding/json"
	"testing"
)

func TestChatHistoryNewChatHistory(t *testing.T) {
	ch := NewChatHistory(nil, nil, nil)
	if ch == nil {
		t.Error("Expected new ChatHistory instance, got nil")
	}
}

func TestChatHistorySetSize(t *testing.T) {
	ch := NewChatHistory(nil, nil, nil)
	size := 5
	ch.SetSize(&size)

	if ch.size == nil || *ch.size != size {
		t.Errorf("Expected size %d, got %v", size, ch.size)
	}
}

func TestChatHistoryClear(t *testing.T) {
	ch := NewChatHistory(nil, nil, nil)
	ch.Append(map[string]any{"role": "user", "content": "test"})
	ch.Clear()

	if len(ch.buffer) != 0 {
		t.Errorf("Expected empty buffer, got length %d", len(ch.buffer))
	}
}

func TestChatHistoryAppend(t *testing.T) {
	// Test with nil size (no limit)
	ch := NewChatHistory(nil, nil, nil)
	item := map[string]any{"role": "user", "content": "test"}
	ch.Append(item)

	if len(ch.buffer) != 1 {
		t.Errorf("Expected buffer length 1, got %d", len(ch.buffer))
	}

	// Test with negative size (no history)
	size := -1
	ch2 := NewChatHistory(&size, nil, nil)
	ch2.Append(item)

	if len(ch2.buffer) != 0 {
		t.Errorf("Expected buffer length 0 for negative size, got %d", len(ch2.buffer))
	}

	// Test with limited size
	size = 2
	ch3 := NewChatHistory(&size, nil, nil)

	// Add 6 items (3 pairs)
	for i := 0; i < 6; i++ {
		ch3.Append(map[string]any{"role": "user", "content": "test"})
		ch3.Append(map[string]any{"role": "assistant", "content": "response"})
	}

	// Should keep only last 2 pairs (4 items)
	expectedLen := 4
	if len(ch3.buffer) != expectedLen {
		t.Errorf("Expected buffer length %d, got %d", expectedLen, len(ch3.buffer))
	}
}

func TestChatHistoryPop(t *testing.T) {
	ch := NewChatHistory(nil, nil, nil)

	// Add some items
	items := []map[string]any{
		{"role": "user", "content": "test1"},
		{"role": "assistant", "content": "response1"},
		{"role": "user", "content": "test2"},
	}

	for _, item := range items {
		ch.Append(item)
	}

	// Pop last item
	ch.Pop(-1)
	if len(ch.buffer) != 2 {
		t.Errorf("Expected buffer length 2 after popping, got %d", len(ch.buffer))
	}

	// Pop first item
	ch.Pop(0)
	if len(ch.buffer) != 1 {
		t.Errorf("Expected buffer length 1 after popping, got %d", len(ch.buffer))
	}

	if ch.buffer[0]["content"] != "response1" {
		t.Errorf("Expected content 'response1', got '%s'", ch.buffer[0]["content"])
	}
}

func TestChatHistoryInit(t *testing.T) {
	ch := NewChatHistory(nil, nil, nil)
	msg := map[string]any{"role": "system", "content": "You are a helpful assistant"}
	ch.Init(msg)

	if ch.initChatMessage == nil || ch.initChatMessage["content"] != "You are a helpful assistant" {
		t.Error("Init message was not set correctly")
	}
}

func TestChatHistoryInitTools(t *testing.T) {
	ch := NewChatHistory(nil, nil, nil)
	tools := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        "get_weather",
			"description": "Get weather information",
		},
	}
	ch.InitTools(tools)

	if ch.initChatTools == nil || ch.initChatTools["type"] != "function" {
		t.Error("Init tools were not set correctly")
	}
}

func TestChatHistoryToList(t *testing.T) {
	// Test with init message and tools
	initMsg := map[string]any{"role": "system", "content": "You are a helpful assistant"}
	initTools := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        "get_weather",
			"description": "Get weather information",
		},
	}

	ch := NewChatHistory(nil, initMsg, initTools)

	// Add some conversation
	ch.Append(map[string]any{"role": "user", "content": "Hello"})
	ch.Append(map[string]any{"role": "assistant", "content": "Hi there!"})

	list := ch.ToList()

	// Should have init message, init tools, and 2 conversation items
	expectedLen := 4
	if len(list) != expectedLen {
		t.Errorf("Expected list length %d, got %d", expectedLen, len(list))
	}

	// Test without init tools
	ch2 := NewChatHistory(nil, initMsg, nil)
	ch2.Append(map[string]any{"role": "user", "content": "Hello"})
	ch2.Append(map[string]any{"role": "assistant", "content": "Hi there!"})

	list2 := ch2.ToList()

	// Should have init message and 2 conversation items
	expectedLen2 := 3
	if len(list2) != expectedLen2 {
		t.Errorf("Expected list length %d, got %d", expectedLen2, len(list2))
	}

	// Test without init message
	ch3 := NewChatHistory(nil, nil, nil)
	ch3.Append(map[string]any{"role": "user", "content": "Hello"})
	ch3.Append(map[string]any{"role": "assistant", "content": "Hi there!"})

	list3 := ch3.ToList()

	// Should have only 2 conversation items
	expectedLen3 := 2
	if len(list3) != expectedLen3 {
		t.Errorf("Expected list length %d, got %d", expectedLen3, len(list3))
	}
}

func TestChatHistoryJSONMarshalUnmarshal(t *testing.T) {
	initMsg := map[string]any{"role": "system", "content": "You are a helpful assistant"}
	initTools := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        "get_weather",
			"description": "Get weather information",
		},
	}
	size := 3
	ch := NewChatHistory(&size, initMsg, initTools)

	// Add some conversation
	ch.Append(map[string]any{"role": "user", "content": "Hello"})
	ch.Append(map[string]any{"role": "assistant", "content": "Hi there!"})

	// Marshal
	data, err := json.Marshal(ch)
	if err != nil {
		t.Errorf("Error marshaling ChatHistory: %v", err)
	}

	// Unmarshal
	var ch2 ChatHistory
	err = json.Unmarshal(data, &ch2)
	if err != nil {
		t.Errorf("Error unmarshaling ChatHistory: %v", err)
	}

	// Check values
	if ch2.size == nil || *ch2.size != size {
		t.Errorf("Expected size %d, got %v", size, ch2.size)
	}

	if ch2.initChatMessage == nil || ch2.initChatMessage["content"] != "You are a helpful assistant" {
		t.Error("Init message was not unmarshaled correctly")
	}
}
