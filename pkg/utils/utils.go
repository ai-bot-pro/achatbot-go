package utils

import (
	"os"
	"time"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// abs returns the absolute value of a duration
func Abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
