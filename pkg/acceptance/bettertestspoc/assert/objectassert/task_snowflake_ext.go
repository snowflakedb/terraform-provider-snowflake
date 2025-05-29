package objectassert

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TaskAssert) HasNotEmptyCreatedOn() *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasNotEmptyId() *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.Id == "" {
			return fmt.Errorf("expected id not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasPredecessorsInAnyOrder(ids ...sdk.SchemaObjectIdentifier) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if !assert.ElementsMatch(t, ids, o.Predecessors) {
			return fmt.Errorf("expected %v predecessors in task relations, got %v", ids, o.TaskRelations.Predecessors)
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasTaskRelations(expected sdk.TaskRelations) *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		errs := make([]error, 0)
		if !assert.ElementsMatch(t, o.TaskRelations.Predecessors, expected.Predecessors) {
			errs = append(errs, fmt.Errorf("expected %v predecessors in task relations, got %v", expected.Predecessors, o.TaskRelations.Predecessors))
		}
		if !reflect.DeepEqual(expected.FinalizerTask, o.TaskRelations.FinalizerTask) {
			errs = append(errs, fmt.Errorf("expected finalizer task: %v; got: %v", expected.FinalizerTask, o.TaskRelations.FinalizerTask))
		}
		if !reflect.DeepEqual(expected.FinalizedRootTask, o.TaskRelations.FinalizedRootTask) {
			errs = append(errs, fmt.Errorf("expected finalized root task: %v; got: %v", expected.FinalizedRootTask, o.TaskRelations.FinalizedRootTask))
		}
		return errors.Join(errs...)
	})
	return t
}

func (t *TaskAssert) HasWarehouseNil() *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.Warehouse != nil {
			return fmt.Errorf("expected warehouse to have no value; got: %s", o.Warehouse.Name())
		}
		return nil
	})
	return t
}

func (t *TaskAssert) HasErrorIntegrationNil() *TaskAssert {
	t.AddAssertion(func(t *testing.T, o *sdk.Task) error {
		t.Helper()
		if o.ErrorIntegration != nil {
			return fmt.Errorf("expected error integration to have no value; got: %s", o.ErrorIntegration.Name())
		}
		return nil
	})
	return t
}
