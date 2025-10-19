package utils

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/weedge/pipeline-go/pkg/logger"
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

func WaitGroupTaskTimeOut(ctx context.Context, task *sync.WaitGroup, timeout time.Duration) {
	// Wait for push audioTask to finish with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	done := make(chan struct{})
	go func() {
		task.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Task completed successfully
	case <-ctx.Done():
		// Timeout occurred
		logger.Warn("Timeout occurred while waiting for push audioTask task to finish")
	}
}
