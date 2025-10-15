package utils

import (
	"testing"
)

func TestResample24KTo16K(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		expected []float32
	}{
		{
			name:     "Empty input",
			input:    []float32{},
			expected: []float32{},
		},
		{
			name:     "Single sample",
			input:    []float32{0.5},
			expected: []float32{0.5},
		},
		{
			name:     "Three samples to two",
			input:    []float32{1.0, 2.0, 3.0},
			expected: []float32{1.0, 2.5},
		},
		{
			name:     "Six samples to four",
			input:    []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
			expected: []float32{1.0, 2.5, 4.0, 5.5},
		},
		{
			name:     "Simple sine wave",
			input:    []float32{0.0, 1.0, 0.0, -1.0, 0.0},
			expected: []float32{0.0, 0.5, -1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Resample24KTo16K(tt.input)

			// Check length
			if len(result) != len(tt.expected) {
				t.Errorf("Resample24KTo16K() length = %d, expected %d", len(result), len(tt.expected))
				return
			}

			// Check values (with some tolerance for floating point comparison)
			for i, v := range result {
				if !floatEquals(v, tt.expected[i]) {
					t.Errorf("Resample24KTo16K() result[%d] = %f, expected %f", i, v, tt.expected[i])
				}
			}
		})
	}
}

func BenchmarkResample24KTo16K(b *testing.B) {
	// Create a sample input of 24000 samples (1 second at 24kHz)
	input := make([]float32, 24000)
	for i := range input {
		input[i] = float32(i%1000) / 1000.0
	}

	
	for b.Loop() {
		Resample24KTo16K(input)
	}
}
