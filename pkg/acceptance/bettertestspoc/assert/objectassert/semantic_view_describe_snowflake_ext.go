package objectassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SemanticViewDetailsAssert struct {
	*assert.SnowflakeObjectAssert[[]sdk.SemanticViewDetails, sdk.SchemaObjectIdentifier]
}
