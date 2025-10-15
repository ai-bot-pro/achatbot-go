package utils

import "math"

const floatEqualityThreshold = 1e-6

func floatEquals(a, b float32) bool {
	return math.Abs(float64(a-b)) <= floatEqualityThreshold
}

// SamplesInt16ToFloat converts a byte slice of int16 values to a slice of float32 samples
func SamplesInt16ToFloat(inSamples []byte) []float32 {
	numSamples := len(inSamples) / 2
	outSamples := make([]float32, numSamples)

	for i := 0; i != numSamples; i++ {
		// Decode two bytes into an int16 using bit manipulation
		s16 := int16(inSamples[2*i]) | int16(inSamples[2*i+1])<<8
		outSamples[i] = float32(s16) / 32768
	}

	return outSamples
}

// SamplesFloatToInt16 converts a slice of float32 samples to a byte slice of int16 values
func SamplesFloatToInt16(inSamples []float32) []byte {
	numSamples := len(inSamples)
	outSamples := make([]byte, numSamples*2)

	for i := 0; i != numSamples; i++ {
		// Convert float32 to int16, clamping values outside [-1.0, 1.0] range
		var s16 int16
		if inSamples[i] >= 1.0 {
			s16 = 32767 // Max int16 value
		} else if inSamples[i] <= -1.0 {
			s16 = -32768 // Min int16 value
		} else {
			s16 = int16(inSamples[i] * 32768)
		}

		// Encode int16 into two bytes using bit manipulation
		outSamples[2*i] = byte(s16)
		outSamples[2*i+1] = byte(s16 >> 8)
	}

	return outSamples
}
