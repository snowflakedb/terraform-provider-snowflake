package objectassert

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationAzureDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationAzureApiProviderType) *ApiIntegrationAzureDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}
