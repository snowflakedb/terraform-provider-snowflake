package sdk

import "context"

// dropSafelyHook runs before SafeDrop. DROP CODE BUNDLE has no documented pre-drop step.
func (*codeBundles) dropSafelyHook(context.Context, SchemaObjectIdentifier) error {
	return nil
}
