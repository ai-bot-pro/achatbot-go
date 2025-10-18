package main

import (
	"achatbot/pkg/utils"
	"fmt"
	"math"
)

func main() {
	// Example: Convert 24KHz audio samples to 16KHz
	// This is a common requirement when working with speech processing
	// where some systems use 24KHz and others use 16KHz

	// Create sample audio data at 24KHz
	// For example, a simple sine wave at 440Hz (A4 note)
	inputSamples := make([]float32, 240) // 0.01 seconds of audio at 24KHz
	for i := range inputSamples {
		// Generate a sine wave: sin(2π * frequency * time)
		frequency := 440.0 // A4 note
		time := float64(i) / 24000.0
		inputSamples[i] = float32(math.Sin(2 * 3.14159265359 * frequency * time))
	}

	fmt.Printf("Original samples (24KHz): %d samples\n", len(inputSamples))

	// Convert to 16KHz
	outputSamples := utils.Resample24KTo16K(inputSamples)

	fmt.Printf("Resampled samples (16KHz): %d samples\n", len(outputSamples))

	// Show first 10 samples of each
	fmt.Println("\nFirst 10 samples:")
	fmt.Println("Original (24KHz):", inputSamples[:10])
	fmt.Println("Resampled (16KHz):", outputSamples[:10])

	// Example with real-world scenario: converting byte array of int16 samples
	// First, let's create some int16 samples
	int16Samples := make([]byte, 480) // 240 samples * 2 bytes per sample
	for i := range 240 {
		// Convert float to int16
		sample := int16(inputSamples[i] * 32767)
		int16Samples[i*2] = byte(sample & 0xFF)
		int16Samples[i*2+1] = byte((sample >> 8) & 0xFF)
	}

	fmt.Printf("\nByte array length: %d bytes (%d int16 samples)\n", len(int16Samples), len(int16Samples)/2)

	// Convert to float samples
	floatSamples := utils.SamplesInt16ToFloat(int16Samples)
	fmt.Printf("Float samples: %d\n", len(floatSamples))

	// Resample
	resampledFloatSamples := utils.Resample24KTo16K(floatSamples)
	fmt.Printf("Resampled float samples: %d\n", len(resampledFloatSamples))

	// Convert back to int16 bytes if needed
	resampledInt16Samples := utils.SamplesFloatToInt16(resampledFloatSamples)
	fmt.Printf("Resampled int16 bytes: %d bytes (%d samples)\n", len(resampledInt16Samples), len(resampledInt16Samples)/2)
}
