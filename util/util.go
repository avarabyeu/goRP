package util

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Retry executes callback func until it executes successfully
func Retry(attempts int, timeout time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	var err error
	for i := 0; i < attempts; i++ {
		var res interface{}
		res, err = callback()
		if err == nil {
			return res, nil
		}
		zap.S().Warnf("Retry failed with the following error: %v", err)

		<-time.After(timeout)
		zap.S().Infof("Retrying... Attempt: %d. Left: %d", i+1, attempts-1-i)
	}

	return nil, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}
