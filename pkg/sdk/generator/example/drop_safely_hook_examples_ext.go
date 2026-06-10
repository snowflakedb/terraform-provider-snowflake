package example

import "context"

func (v *dropSafelyHookExamples) dropSafelyHook(_ context.Context, _ AccountObjectIdentifier) error {
	return nil
}
