package util

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.Warnf("Retry failed with the following error: %v", err)

		<-time.After(timeout)
		log.Infof("Retrying... Attempt: %d. Left: %d", i+1, attempts-1-i)
	}

	return nil, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}
