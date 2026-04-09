package invokeactionassert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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
// Optionally checks level and apply method via TAG_REFERENCES when set with WithLevel/WithApplyMethod.
type tagValueOnObjectCheck struct {
	tagId               sdk.SchemaObjectIdentifier
	objectIdProvider    func() sdk.ObjectIdentifier
	objectDomain        sdk.TagReferenceObjectDomain
	expectedValue       string
	expectedLevel       *sdk.TagReferenceObjectDomain
	expectedApplyMethod *sdk.TagReferenceApplyMethod
}

func (c *tagValueOnObjectCheck) WithLevel(level sdk.TagReferenceObjectDomain) *tagValueOnObjectCheck {
	c.expectedLevel = &level
	return c
}

func (c *tagValueOnObjectCheck) WithApplyMethod(method sdk.TagReferenceApplyMethod) *tagValueOnObjectCheck {
	c.expectedApplyMethod = &method
	return c
}

func (c *tagValueOnObjectCheck) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		objectId := c.objectIdProvider()

		refs, err := testClient.Tag.GetReferencesForObject(t, objectId, c.objectDomain)
		if err != nil {
			return fmt.Errorf("error getting tag references for %s %s: %w", c.objectDomain, objectId.FullyQualifiedName(), err)
		}
		ref, err := collections.FindFirst(refs, func(ref sdk.TagReference) bool {
			return ref.TagId().FullyQualifiedName() == c.tagId.FullyQualifiedName()
		})
		if err != nil {
			return fmt.Errorf("tag reference %s not found on %s %s: %w", c.tagId.FullyQualifiedName(), c.objectDomain, objectId.FullyQualifiedName(), err)
		}
		if ref.TagValue != c.expectedValue {
			return fmt.Errorf("expected tag value %q on %s %s, got %q", c.expectedValue, c.objectDomain, objectId.FullyQualifiedName(), ref.TagValue)
		}
		if c.expectedLevel != nil && ref.Level != *c.expectedLevel {
			return fmt.Errorf("expected tag level %v on %s %s, got %v", *c.expectedLevel, c.objectDomain, objectId.FullyQualifiedName(), ref.Level)
		}
		if c.expectedApplyMethod != nil && ref.ApplyMethod != *c.expectedApplyMethod {
			return fmt.Errorf("expected tag apply method %v on %s %s, got %v", *c.expectedApplyMethod, c.objectDomain, objectId.FullyQualifiedName(), ref.ApplyMethod)
		}

		return nil
	}
}

// TagValueOnObject asserts that the given tag has the expected value on the specified object.
// The objectId is resolved lazily via a provider function, so it can reference variables
// that are set up after test step definitions (e.g., in PreConfig callbacks).
// Use WithLevel and WithApplyMethod to additionally verify how the tag was applied.
func TagValueOnObject(t *testing.T, tagId sdk.SchemaObjectIdentifier, objectIdProvider func() sdk.ObjectIdentifier, objectDomain sdk.TagReferenceObjectDomain, expectedValue string) *tagValueOnObjectCheck {
	t.Helper()

	return &tagValueOnObjectCheck{
		tagId:            tagId,
		objectIdProvider: objectIdProvider,
		objectDomain:     objectDomain,
		expectedValue:    expectedValue,
	}
}
