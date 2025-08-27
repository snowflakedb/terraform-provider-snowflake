package objectassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SemanticViewDetailsAssert struct {
	*assert.SnowflakeObjectAssert[[]sdk.SemanticViewDetails, sdk.SchemaObjectIdentifier]
}

//func SemanticViewDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *SemanticViewDetailsAssert {
//	t.Helper()
//	return &SemanticViewDetailsAssert{
//		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeSemanticView, id, func(testClient *helpers.TestClient) assert.ObjectProvider[[]sdk.SemanticViewDetails, sdk.SchemaObjectIdentifier] {
//			return testClient.SemanticView.Describe
//		}),
//	}
//}
