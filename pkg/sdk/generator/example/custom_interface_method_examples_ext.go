package example

import (
	"context"
)

// This file contains manual implementations of custom interface methods declared
// with WithCustomInterfaceMethod in the definition file. These methods appear in
// the generated interface but have no generated implementation.
//
// In a real SDK object, this file would live alongside the generated files and
// provide the full implementation.

func (v *customInterfaceMethodExamples) UnsetAll(_ context.Context) error {
	return nil
}

func (v *customInterfaceMethodExamples) ShowParameters(_ context.Context, _ AccountObjectIdentifier) ([]*Parameter, error) {
	return nil, nil
}

func (v *customInterfaceMethodExamples) SuspendRootTasks(_ context.Context, _ SchemaObjectIdentifier, _ SchemaObjectIdentifier) ([]SchemaObjectIdentifier, error) {
	return nil, nil
}

func (v *customInterfaceMethodExamples) Refresh(_ context.Context, _ AccountObjectIdentifier) {
}

func (v *customInterfaceMethodExamples) ResumeTasks(_ context.Context, _ []SchemaObjectIdentifier) error {
	return nil
}
