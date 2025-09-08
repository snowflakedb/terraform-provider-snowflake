package objectassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SemanticViewDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.SemanticViewDetails, sdk.SchemaObjectIdentifier]
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
