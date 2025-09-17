package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SemanticViewDetailsCollection struct {
	Details []sdk.SemanticViewDetails
}
type SemanticViewDetailsAssert struct {
	*assert.SnowflakeObjectAssert[SemanticViewDetailsCollection, sdk.SchemaObjectIdentifier]
}

func SemanticViewDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *SemanticViewDetailsAssert {
	t.Helper()
	return &SemanticViewDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeSemanticView, id, func(testClient *helpers.TestClient) assert.ObjectProvider[SemanticViewDetailsCollection, sdk.SchemaObjectIdentifier] {
			return func(t *testing.T, id sdk.SchemaObjectIdentifier) (*SemanticViewDetailsCollection, error) {
				details, err := testClient.SemanticView.Describe(t, id)
				if err != nil {
					return nil, err
				}
				return &SemanticViewDetailsCollection{Details: details}, nil
			}
		}),
	}
}

func (s *SemanticViewDetailsAssert) HasDetailsCount(expected int) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		if len(o.Details) != expected {
			return fmt.Errorf("expected %d semantic view details; got: %d", expected, len(o.Details))
		}
		return nil
	})
	return s
}

func (s *SemanticViewDetailsAssert) ContainsDetail(expected sdk.SemanticViewDetails) *SemanticViewDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *SemanticViewDetailsCollection) error {
		t.Helper()
		found := slices.ContainsFunc(o.Details, func(detail sdk.SemanticViewDetails) bool {
			return detail.ObjectKind == expected.ObjectKind &&
				detail.ObjectName == expected.ObjectName &&
				detail.Property == expected.Property &&
				detail.PropertyValue == expected.PropertyValue
		})
		if !found {
			return fmt.Errorf("expected semantic view to contain a detail row matching %v", expected)
		}
		return nil
	})
	return s
}

func NewSemanticViewDetails(
	ObjectKind string,
	ObjectName string,
	ParentEntity *string,
	Property string,
	PropertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		ObjectKind:    ObjectKind,
		ObjectName:    ObjectName,
		Property:      Property,
		PropertyValue: PropertyValue,
	}
	if ParentEntity != nil {
		details.ParentEntity = ParentEntity
	}

	return details
}
