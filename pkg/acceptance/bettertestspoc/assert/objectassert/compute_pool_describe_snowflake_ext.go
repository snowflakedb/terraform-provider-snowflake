package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// TODO [SNOW-1501905]: this file should be fully regenerated when adding and option to assert the results of describe
type ComputePoolDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ComputePoolDetails, sdk.AccountObjectIdentifier]
}

func ComputePoolDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ComputePoolDetailsAssert {
	t.Helper()
	return &ComputePoolDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("COMPUTE_POOL_DETAILS"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ComputePoolDetails, sdk.AccountObjectIdentifier] {
			return testClient.ComputePool.Describe
		}),
	}
}

func (f *ComputePoolDetailsAssert) HasName(expected string) *ComputePoolDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return f
}

func (c *ComputePoolDetailsAssert) HasState(expected sdk.ComputePoolState) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.State != expected {
			return fmt.Errorf("expected state: %v; got: %v", expected, o.State)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasMinNodes(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.MinNodes != expected {
			return fmt.Errorf("expected min nodes: %v; got: %v", expected, o.MinNodes)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasMaxNodes(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.MaxNodes != expected {
			return fmt.Errorf("expected max nodes: %v; got: %v", expected, o.MaxNodes)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasInstanceFamily(expected sdk.ComputePoolInstanceFamily) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.InstanceFamily != expected {
			return fmt.Errorf("expected instance family: %v; got: %v", expected, o.InstanceFamily)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasNumServices(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.NumServices != expected {
			return fmt.Errorf("expected num services: %v; got: %v", expected, o.NumServices)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasNumJobs(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.NumJobs != expected {
			return fmt.Errorf("expected num jobs: %v; got: %v", expected, o.NumJobs)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasAutoSuspendSecs(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.AutoSuspendSecs != expected {
			return fmt.Errorf("expected auto suspend secs: %v; got: %v", expected, o.AutoSuspendSecs)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasAutoResume(expected bool) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.AutoResume != expected {
			return fmt.Errorf("expected auto resume: %v; got: %v", expected, o.AutoResume)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasActiveNodes(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.ActiveNodes != expected {
			return fmt.Errorf("expected active nodes: %v; got: %v", expected, o.ActiveNodes)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasIdleNodes(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.IdleNodes != expected {
			return fmt.Errorf("expected idle nodes: %v; got: %v", expected, o.IdleNodes)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasTargetNodes(expected int) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.TargetNodes != expected {
			return fmt.Errorf("expected target nodes: %v; got: %v", expected, o.TargetNodes)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasCreatedOn(expected time.Time) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasResumedOn(expected time.Time) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.ResumedOn != expected {
			return fmt.Errorf("expected resumed on: %v; got: %v", expected, o.ResumedOn)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasUpdatedOn(expected time.Time) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.UpdatedOn != expected {
			return fmt.Errorf("expected updated on: %v; got: %v", expected, o.UpdatedOn)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasOwner(expected string) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasComment(expected string) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.Comment == nil {
			return fmt.Errorf("expected comment to have value; got: nil")
		}
		if *o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, *o.Comment)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasIsExclusive(expected bool) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.IsExclusive != expected {
			return fmt.Errorf("expected is exclusive: %v; got: %v", expected, o.IsExclusive)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasApplication(expected *sdk.AccountObjectIdentifier) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.Application == nil {
			return fmt.Errorf("expected application to have value; got: nil")
		}
		if o.Application.FullyQualifiedName() != expected.FullyQualifiedName() {
			return fmt.Errorf("expected application: %v; got: %v", expected, *o.Application)
		}
		return nil
	})
	return c
}

func (a *ComputePoolDetailsAssert) HasCreatedOnNotEmpty() *ComputePoolDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.CreatedOn == (time.Time{}) {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ComputePoolDetailsAssert) HasResumedOnNotEmpty() *ComputePoolDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.ResumedOn == (time.Time{}) {
			return fmt.Errorf("expected resumed_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ComputePoolDetailsAssert) HasUpdatedOnNotEmpty() *ComputePoolDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.UpdatedOn == (time.Time{}) {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return a
}

func (c *ComputePoolDetailsAssert) HasNoComment() *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to have nil; got: %s", *o.Comment)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasNoApplication() *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.Application != nil {
			return fmt.Errorf("expected application to have nil; got: %s", *o.Application)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasErrorCode(expected string) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.ErrorCode != expected {
			return fmt.Errorf("expected error_code: %v; got: %v", expected, o.ErrorCode)
		}
		return nil
	})
	return c
}

func (c *ComputePoolDetailsAssert) HasStatusMessage(expected string) *ComputePoolDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.ComputePoolDetails) error {
		t.Helper()
		if o.StatusMessage != expected {
			return fmt.Errorf("expected status_message: %v; got: %v", expected, o.StatusMessage)
		}
		return nil
	})
	return c
}
