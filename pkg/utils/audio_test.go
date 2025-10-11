package utils

import (
	"math"
	"testing"
)

const floatEqualityThreshold = 1e-6

func floatEquals(a, b float32) bool {
	return math.Abs(float64(a-b)) <= floatEqualityThreshold
}

func TestSamplesInt16ToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []float32
	}{
		{
			name:     "Empty input",
			input:    []byte{},
			expected: []float32{},
		},
		{
			name:     "Single sample - zero",
			input:    []byte{0x00, 0x00},
			expected: []float32{0.0},
		},
		{
			name:     "Single sample - positive max",
			input:    []byte{0xFF, 0x7F},           // 0x7FFF = 32767
			expected: []float32{32767.0 / 32768.0}, // Actual value is 0.999969...
		},
		{
			name:     "Single sample - negative max",
			input:    []byte{0x00, 0x80}, // 0x8000 = -32768
			expected: []float32{-1.0},
		},
		{
			name:     "Multiple samples",
			input:    []byte{0x00, 0x00, 0xFF, 0x7F, 0x00, 0x80},
			expected: []float32{0.0, 32767.0 / 32768.0, -1.0},
		},
		{
			name:     "Sample with negative value",
			input:    []byte{0xD2, 0x04}, // 0x04D2 = 1234
			expected: []float32{1234.0 / 32768.0},
		},
		{
			name:     "Sample with another negative value",
			input:    []byte{0x2E, 0xFB}, // 0xFB2E = -1234
			expected: []float32{-1234.0 / 32768.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SamplesInt16ToFloat(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("SamplesInt16ToFloat() length = %d, expected %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if !floatEquals(v, tt.expected[i]) {
					t.Errorf("SamplesInt16ToFloat() result[%d] = %f, expected %f", i, v, tt.expected[i])
				}
			}
		})
	}
}
