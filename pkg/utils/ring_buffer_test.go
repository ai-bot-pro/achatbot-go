package utils

import (
	"reflect"
	"testing"
)

func TestRingBuffer_Push(t *testing.T) {
	rb := NewRingBuffer(5)

	// 测试Push方法
	values := []any{1, 2, 3, 4, 5}
	rb.Push(values)

	if rb.Len() != 5 {
		t.Errorf("Expected length 5, got %d", rb.Len())
	}
}

func TestRingBuffer_Head(t *testing.T) {
	rb := NewRingBuffer(5)

	// 测试Head方法
	head := rb.Head()
	if head != 0 {
		t.Errorf("Expected head position 0, got %d", head)
	}

	// 添加一些元素后再次测试
	rb.PushBack(1)
	rb.PushBack(2)
	head = rb.Head()
	if head != 0 {
		t.Errorf("Expected head position 0, got %d", head)
	}
}

func TestRingBuffer_Pop(t *testing.T) {
	rb := NewRingBuffer(5)

	// 添加一些元素
	values := []any{1, 2, 3, 4, 5}
	rb.Push(values)

	// 测试Pop方法
	popped := rb.Pop(3)
	expected := []any{1, 2, 3}

	if !reflect.DeepEqual(popped, expected) {
		t.Errorf("Expected popped values %v, got %v", expected, popped)
	}

	if rb.Len() != 2 {
		t.Errorf("Expected length 2 after pop, got %d", rb.Len())
	}
}
func TestRingBuffer_PopBytes(t *testing.T) {
	rb := NewRingBuffer(5)

	// 添加一些元素
	values := []byte{1, 2, 3, 4, 5}
	rb.PushBytes(values)

	// 测试Pop方法
	popped := rb.PopBytes(3)
	expected := []byte{1, 2, 3}

	if !reflect.DeepEqual(popped, expected) {
		t.Errorf("Expected popped values %v, got %v", expected, popped)
	}

	if rb.Len() != 2 {
		t.Errorf("Expected length 2 after pop, got %d", rb.Len())
	}

	popped = rb.PopBytes(3)
	expected = []byte{4, 5}

	if !reflect.DeepEqual(popped, expected) {
		t.Errorf("Expected popped values %v, got %v", expected, popped)
	}
	if rb.Len() != 0 {
		t.Errorf("Expected length 0 after pop, got %d", rb.Len())
	}
}

func TestRingBuffer_Get(t *testing.T) {
	rb := NewRingBuffer(5)

	// 添加一些元素
	values := []any{1, 2, 3, 4, 5}
	rb.Push(values)

	// 测试Get方法
	got := rb.Get(1, 3)
	expected := []any{2, 3, 4}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected values %v, got %v", expected, got)
	}

	// 测试边界情况
	got = rb.Get(0, 10) // 请求超过缓冲区大小的元素
	expected = []any{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected values %v, got %v", expected, got)
	}
}

func TestRingBuffer_Size(t *testing.T) {
	rb := NewRingBuffer(5)

	// 测试Size方法
	size := rb.Size()
	if size != 0 {
		t.Errorf("Expected size 0, got %d", size)
	}

	// 添加一些元素后再次测试
	rb.PushBack(1)
	rb.PushBack(2)
	size = rb.Size()
	if size != 2 {
		t.Errorf("Expected size 2, got %d", size)
	}
}

func TestRingBuffer_Reset(t *testing.T) {
	rb := NewRingBuffer(5)

	// 添加一些元素
	values := []any{1, 2, 3, 4, 5}
	rb.Push(values)

	// 测试Reset方法
	rb.Reset()

	if rb.Len() != 0 {
		t.Errorf("Expected length 0 after reset, got %d", rb.Len())
	}

	if rb.Size() != 0 {
		t.Errorf("Expected size 0 after reset, got %d", rb.Size())
	}
}
