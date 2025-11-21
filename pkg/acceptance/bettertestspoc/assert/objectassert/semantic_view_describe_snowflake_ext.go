package objectassert

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
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
				t.Helper()
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
			return detail.Property == expected.Property &&
				detail.PropertyValue == expected.PropertyValue &&
				reflect.DeepEqual(detail.ObjectName, expected.ObjectName) &&
				reflect.DeepEqual(detail.ObjectKind, expected.ObjectKind) &&
				reflect.DeepEqual(detail.ParentEntity, expected.ParentEntity)
		})
		if !found {
			return fmt.Errorf("expected semantic view to contain a detail row matching %s", semanticViewDetailsString(expected))
		}
		return nil
	})
	return s
}

func semanticViewDetailsString(d sdk.SemanticViewDetails) string {
	stringBuilder := new(strings.Builder)
	stringBuilder.WriteString("[")
	if d.ObjectKind != nil {
		stringBuilder.WriteString(fmt.Sprintf("object_kind=%s ", *d.ObjectKind))
	}
	if d.ObjectName != nil {
		stringBuilder.WriteString(fmt.Sprintf("object_name=%s ", *d.ObjectName))
	}
	if d.ParentEntity != nil {
		stringBuilder.WriteString(fmt.Sprintf("parent_entity=%s ", *d.ParentEntity))
	}
	stringBuilder.WriteString(fmt.Sprintf("property=%s property_value=%s", d.Property, d.PropertyValue))
	stringBuilder.WriteString("]")
	return stringBuilder.String()
}

func NewSemanticViewDetails(
	objectKind *string,
	objectName *string,
	parentEntity *string,
	property string,
	propertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		Property:      property,
		PropertyValue: propertyValue,
	}
	if objectKind != nil {
		details.ObjectKind = objectKind
	}
	if objectName != nil {
		details.ObjectName = objectName
	}
	if parentEntity != nil {
		details.ParentEntity = parentEntity
	}

	return details
}

func NewSemanticViewDetailsTable(
	tableAlias string,
	property string,
	propertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		Property:      property,
		PropertyValue: propertyValue,
	}
	details.ObjectKind = sdk.Pointer("TABLE")
	details.ObjectName = sdk.Pointer(tableAlias)
	return details
}

func NewSemanticViewDetailsDimension(
	dimensionName string,
	tableAlias string,
	property string,
	propertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		Property:      property,
		PropertyValue: propertyValue,
	}
	details.ObjectKind = sdk.Pointer("DIMENSION")
	details.ObjectName = sdk.Pointer(dimensionName)
	details.ParentEntity = sdk.Pointer(tableAlias)
	return details
}

func NewSemanticViewDetailsFact(
	factName string,
	tableAlias string,
	property string,
	propertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		Property:      property,
		PropertyValue: propertyValue,
	}
	details.ObjectKind = sdk.Pointer("FACT")
	details.ObjectName = sdk.Pointer(factName)
	details.ParentEntity = sdk.Pointer(tableAlias)
	return details
}

func NewSemanticViewDetailsMetric(
	metricName string,
	tableAlias string,
	property string,
	propertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		Property:      property,
		PropertyValue: propertyValue,
	}
	details.ObjectKind = sdk.Pointer("METRIC")
	details.ObjectName = sdk.Pointer(metricName)
	details.ParentEntity = sdk.Pointer(tableAlias)
	return details
}

func NewSemanticViewDetailsRelationship(
	relationshipName string,
	tableAlias string,
	property string,
	propertyValue string,
) sdk.SemanticViewDetails {
	details := sdk.SemanticViewDetails{
		Property:      property,
		PropertyValue: propertyValue,
	}
	details.ObjectKind = sdk.Pointer("RELATIONSHIP")
	details.ObjectName = sdk.Pointer(relationshipName)
	details.ParentEntity = sdk.Pointer(tableAlias)
	return details
}
