package util

import (
	"context"
	"fmt"
	"log"
	"time"
)

// TODO(SNOW-3071484): Improve retry logic.
// - Provide sane default attempts and sleep duration.
// - Discuss if it should support exp backoff.
// - Retry on errors that are retriable, not all. Currently, the callers are responsible for this.
// - Handle error history.
func Retry(attempts int, sleepDuration time.Duration, f func() (error, bool)) error {
	for range attempts {
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

// RetryWithContext retries f until it reports done, returns an error, or ctx is canceled.
// Sleep between attempts respects ctx. Callers should set a deadline on ctx (for Terraform
// resources, CreateContext already carries the resource create timeout).
func RetryWithContext(ctx context.Context, sleepDuration time.Duration, f func() (error, bool)) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("giving up: %w", err)
		}
		err, done := f()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		log.Printf("[INFO] operation not finished yet, retrying in %v seconds", sleepDuration.Seconds())
		timer := time.NewTimer(sleepDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("giving up: %w", ctx.Err())
		case <-timer.C:
		}
	}
}
