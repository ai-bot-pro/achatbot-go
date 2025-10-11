package utils

import "container/list"

// RingBuffer 是一个固定大小的环形缓冲区, 简单实现, 基于 container/list
// Note: 当前场景下, 用于每个连接会话音频数据缓存, 暂无共享并发场景 (use LFRB)
type RingBuffer struct {
	buffer *list.List
	maxLen int
}

// newRingBuffer 创建一个新的环形缓冲区
func NewRingBuffer(maxLen int) *RingBuffer {
	return &RingBuffer{
		buffer: list.New(),
		maxLen: maxLen,
	}
}

// PushBack 向缓冲区末尾添加元素，如果超过最大长度则移除最旧的元素
func (rb *RingBuffer) PushBack(value any) {
	if rb.maxLen <= 0 {
		return
	}

	// 如果已达到最大长度，则移除最旧的元素
	if rb.buffer.Len() >= rb.maxLen {
		rb.buffer.Remove(rb.buffer.Front())
	}

	rb.buffer.PushBack(value)
}

// Push 向缓冲区添加多个元素
func (rb *RingBuffer) Push(values []any) {
	for _, value := range values {
		rb.PushBack(value)
	}
}

// PushBytes 向缓冲区添加[]bytes
func (rb *RingBuffer) PushBytes(values []byte) {
	for _, value := range values {
		rb.PushBack(value)
	}
}

// Head 获取环形缓冲区头位置
func (rb *RingBuffer) Head() int {
	return 0 // 由于使用list实现，头部始终是0
}

// Pop 从缓冲区弹出n个元素
func (rb *RingBuffer) Pop(n int) []any {
	if n <= 0 {
		return []any{}
	}

	result := make([]any, 0, n)
	count := 0

	for rb.buffer.Len() > 0 && count < n {
		element := rb.buffer.Front()
		if element != nil {
			result = append(result, element.Value)
			rb.buffer.Remove(element)
			count++
		} else {
			break
		}
	}

	return result
}

// Pop 从缓冲区弹出n个byte
func (rb *RingBuffer) PopBytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}

	result := make([]byte, 0, n)
	count := 0

	for rb.buffer.Len() > 0 && count < n {
		element := rb.buffer.Front()
		if element != nil {
			result = append(result, element.Value.(byte))
			rb.buffer.Remove(element)
			count++
		} else {
			break
		}
	}

	return result
}

// Get 从开始位置获取n个元素
func (rb *RingBuffer) Get(start int, n int) []any {
	if start < 0 || n <= 0 || start >= rb.buffer.Len() {
		return []any{}
	}

	result := make([]any, 0, n)
	count := 0

	// 找到起始位置
	current := rb.buffer.Front()
	for i := 0; i < start && current != nil; i++ {
		current = current.Next()
	}

	// 从起始位置获取n个元素
	for current != nil && count < n {
		result = append(result, current.Value)
		current = current.Next()
		count++
	}

	return result
}

// Size 获取缓冲区大小
func (rb *RingBuffer) Size() int {
	return rb.buffer.Len()
}

// Reset 重置缓冲区
func (rb *RingBuffer) Reset() {
	rb.buffer.Init()
}

// Front 返回缓冲区的第一个元素
func (rb *RingBuffer) Front() *list.Element {
	return rb.buffer.Front()
}

// Clear 清空缓冲区
func (rb *RingBuffer) Clear() {
	rb.buffer.Init()
}

// Len 返回缓冲区长度
func (rb *RingBuffer) Len() int {
	return rb.buffer.Len()
}
