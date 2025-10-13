package utils

// Resample24KTo16K converts audio samples from 24KHz to 16KHz sampling rate
func Resample24KTo16K(input []float32) []float32 {
	// Calculate output length
	// 16000/24000 = 2/3, so we need 2 samples for every 3 input samples
	return Resample(input, 24000, 16000) // 2/3
}

// Resample uses linear interpolation for resampling
func Resample(input []float32, inputRate, outputRate int) []float32 {
	if len(input) == 0 {
		return []float32{}
	}
	if inputRate == outputRate {
		return input
	}

	outputLength := int(float64(len(input)) * float64(outputRate) / float64(inputRate))

	// Handle edge case where output length is 0 but input is not empty
	if outputLength == 0 && len(input) > 0 {
		return []float32{input[0]}
	}

	output := make([]float32, outputLength)

	// Resampling ratio
	ratio := float64(inputRate) / float64(outputRate) // 24000/16000 = 1.5

	for i := range outputLength {
		// Calculate the exact position in the input
		pos := float64(i) * ratio

		// Get the indices of the two samples to interpolate between
		index := int(pos)
		nextIndex := index + 1

		// Handle edge case where we're at the last sample
		if nextIndex >= len(input) {
			output[i] = input[index]
			continue
		}

		// Linear interpolation
		fraction := pos - float64(index)
		output[i] = input[index] + float32(fraction)*(input[nextIndex]-input[index])
	}

	return output
}

func Resample24KTo16KBytes(input []byte) []byte {
	return ResampleBytes(input, 24000, 16000)
}

func ResampleBytes(input []byte, inputRate, outputRate int) []byte {
	inputFloat32 := SamplesInt16ToFloat(input)
	outputFloat32 := Resample(inputFloat32, inputRate, outputRate)
	outputBytes := SamplesFloatToInt16(outputFloat32)

	return outputBytes
}
