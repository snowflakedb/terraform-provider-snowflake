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
	tagId            sdk.SchemaObjectIdentifier
	objectIdProvider func() sdk.ObjectIdentifier
	objectType       sdk.ObjectType
	expectedValue    string
	expectedLevel    *sdk.TagReferenceObjectDomain
	expectedApply    *sdk.TagReferenceApplyMethod
}

func (c *tagValueOnObjectCheck) WithLevel(level sdk.TagReferenceObjectDomain) *tagValueOnObjectCheck {
	c.expectedLevel = &level
	return c
}

func (c *tagValueOnObjectCheck) WithApplyMethod(method sdk.TagReferenceApplyMethod) *tagValueOnObjectCheck {
	c.expectedApply = &method
	return c
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

		if c.expectedLevel != nil || c.expectedApply != nil {
			ref, err := findTagReference(t, testClient, c.tagId, objectId, c.objectType)
			if err != nil {
				return err
			}
			if c.expectedLevel != nil && ref.Level != *c.expectedLevel {
				return fmt.Errorf("expected tag level %v on %s %s, got %v", *c.expectedLevel, c.objectType, objectId.FullyQualifiedName(), ref.Level)
			}
			if c.expectedApply != nil && ref.ApplyMethod != *c.expectedApply {
				return fmt.Errorf("expected tag apply method %v on %s %s, got %v", *c.expectedApply, c.objectType, objectId.FullyQualifiedName(), ref.ApplyMethod)
			}
		}

		return nil
	}
}

// TagValueOnObject asserts that the given tag has the expected value on the specified object.
// The objectId is resolved lazily via a provider function, so it can reference variables
// that are set up after test step definitions (e.g., in PreConfig callbacks).
// Use WithLevel and WithApplyMethod to additionally verify how the tag was applied.
func TagValueOnObject(t *testing.T, tagId sdk.SchemaObjectIdentifier, objectIdProvider func() sdk.ObjectIdentifier, objectType sdk.ObjectType, expectedValue string) *tagValueOnObjectCheck {
	t.Helper()
	return &tagValueOnObjectCheck{
		tagId:            tagId,
		objectIdProvider: objectIdProvider,
		objectType:       objectType,
		expectedValue:    expectedValue,
	}
}

// objectTypeToTagReferenceDomain maps sdk.ObjectType to the TagReferenceObjectDomain
// used by INFORMATION_SCHEMA.TAG_REFERENCES. Snowflake normalizes some types
// (e.g., VIEW -> TABLE for tag references).
func objectTypeToTagReferenceDomain(objectType sdk.ObjectType) sdk.TagReferenceObjectDomain {
	switch objectType {
	case sdk.ObjectTypeView, sdk.ObjectTypeMaterializedView:
		return sdk.TagReferenceObjectDomainTable
	default:
		return sdk.TagReferenceObjectDomain(objectType)
	}
}

func findTagReference(t *testing.T, testClient *helpers.TestClient, tagId sdk.SchemaObjectIdentifier, objectId sdk.ObjectIdentifier, objectType sdk.ObjectType) (*sdk.TagReference, error) {
	t.Helper()
	domain := objectTypeToTagReferenceDomain(objectType)
	refs, err := testClient.Tag.GetReferencesForObject(t, objectId, domain)
	if err != nil {
		return nil, fmt.Errorf("error getting tag references for %s %s: %w", objectType, objectId.FullyQualifiedName(), err)
	}
	return collections.FindFirst(refs, func(ref sdk.TagReference) bool {
		return ref.TagDatabase == tagId.DatabaseName() &&
			ref.TagSchema == tagId.SchemaName() &&
			ref.TagName == tagId.Name()
	})
}
