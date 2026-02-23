package objectassert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ExternalVolumeDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ExternalVolumeDetails, sdk.AccountObjectIdentifier]
}

func ExternalVolumeDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ExternalVolumeDetailsAssert {
	t.Helper()
	return &ExternalVolumeDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(
			sdk.ObjectType("EXTERNAL_VOLUME_DETAILS"),
			id,
			func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ExternalVolumeDetails, sdk.AccountObjectIdentifier] {
				return testClient.ExternalVolume.Describe
			}),
	}
}

func (e *ExternalVolumeDetailsAssert) HasActive(expected string) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if o.Active != expected {
			return fmt.Errorf("expected active: %v; got: %v", expected, o.Active)
		}
		return nil
	})
	return e
}

func (e *ExternalVolumeDetailsAssert) HasComment(expected string) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return e
}

func (e *ExternalVolumeDetailsAssert) HasAllowWrites(expected string) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if o.AllowWrites != expected {
			return fmt.Errorf("expected allow writes: %v; got: %v", expected, o.AllowWrites)
		}
		return nil
	})
	return e
}

func (e *ExternalVolumeDetailsAssert) HasStorageLocations(expected ...sdk.ExternalVolumeStorageLocationDetails) *ExternalVolumeDetailsAssert {
	e.AddAssertion(func(t *testing.T, o *sdk.ExternalVolumeDetails) error {
		t.Helper()
		if len(o.StorageLocations) != len(expected) {
			return fmt.Errorf("expected %d storage locations; got: %d\nexpected: %v\ngot: %v", len(expected), len(o.StorageLocations), expected, o.StorageLocations)
		}
		for i := range expected {
			if !reflect.DeepEqual(o.StorageLocations[i], expected[i]) {
				return fmt.Errorf("storage location at index %d differs:\nexpected: %v\ngot: %v", i, expected[i], o.StorageLocations[i])
			}
		}
		return nil
	})
	return e
}
