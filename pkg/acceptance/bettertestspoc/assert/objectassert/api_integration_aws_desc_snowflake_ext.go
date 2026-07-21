package objectassert

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationAwsDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationAwsApiProviderType) *ApiIntegrationAwsDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}
