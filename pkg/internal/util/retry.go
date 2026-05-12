package util

import (
	"fmt"
	"log"
	"time"
)

// TODO(SNOW-3071484): Improve retry logic.
// - Provide sane default attempts and sleep duration.
// - Discuss if it should support exp backoff.
// - Retry on errors that are retriable, not all. Currently, the callers are responsible for this.
// - Add unit tests.
// - Handle error history.
func Retry(attempts int, sleepDuration time.Duration, f func() (error, bool)) error {
	for i := 0; i < attempts; i++ {
		err, done := f()
		if err != nil {
			return err
		}
		if done {
			return nil
		} else {
			log.Printf("[INFO] operation not finished yet, retrying in %v seconds", sleepDuration.Seconds())
			time.Sleep(sleepDuration)
		}
	}
	return fmt.Errorf("giving up after %v attempts", attempts)
}
