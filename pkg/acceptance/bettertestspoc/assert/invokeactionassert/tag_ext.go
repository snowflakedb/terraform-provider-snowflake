package invokeactionassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TagDoesNotExist(t *testing.T, id sdk.SchemaObjectIdentifier) assert.TestCheckFuncProvider {
	t.Helper()
	return newNonExistenceCheck(
		sdk.ObjectTypeTag,
		id,
		func(testClient *helpers.TestClient) showByIDFunc[*sdk.Tag, sdk.SchemaObjectIdentifier] {
			return testClient.Tag.Show
		})
}

// tagValueOnObjectCheck verifies that a tag has the expected value on a given object.
// The objectIdProvider is called at check time, allowing lazy resolution of IDs
// that are set up after test step definitions (e.g., in PreConfig callbacks).
type tagValueOnObjectCheck struct {
	tagId            sdk.SchemaObjectIdentifier
	objectIdProvider func() sdk.ObjectIdentifier
	objectType       sdk.ObjectType
	expectedValue    string
}

func (c *tagValueOnObjectCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		objectId := c.objectIdProvider()
		val, err := testClient.Tag.GetForObject(t, c.tagId, objectId, c.objectType)
		if err != nil {
			return fmt.Errorf("error getting tag %s for %s %s: %w", c.tagId.FullyQualifiedName(), c.objectType, objectId.FullyQualifiedName(), err)
		}
		if val == nil {
			return fmt.Errorf("expected tag value %q on %s %s, got nil", c.expectedValue, c.objectType, objectId.FullyQualifiedName())
		}
		if *val != c.expectedValue {
			return fmt.Errorf("expected tag value %q on %s %s, got %q", c.expectedValue, c.objectType, objectId.FullyQualifiedName(), *val)
		}
		return nil
	}
}

// TagValueOnObject asserts that the given tag has the expected value on the specified object.
// The objectId is resolved lazily via a provider function, so it can reference variables
// that are set up after test step definitions (e.g., in PreConfig callbacks).
func TagValueOnObject(t *testing.T, tagId sdk.SchemaObjectIdentifier, objectIdProvider func() sdk.ObjectIdentifier, objectType sdk.ObjectType, expectedValue string) assert.TestCheckFuncProvider {
	t.Helper()
	return &tagValueOnObjectCheck{
		tagId:            tagId,
		objectIdProvider: objectIdProvider,
		objectType:       objectType,
		expectedValue:    expectedValue,
	}
}
