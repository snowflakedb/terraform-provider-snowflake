package example

import "context"

func (v *dropSafelySuppressedExamples) DropSafely(ctx context.Context, request *DropDropSafelySuppressedExampleRequest) error {
	return v.Drop(ctx, request)
}
