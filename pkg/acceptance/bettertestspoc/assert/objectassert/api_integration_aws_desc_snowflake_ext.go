package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *ApiIntegrationAwsDetailsAssert) HasApiProviderType(expected sdk.ApiIntegrationAwsApiProviderType) *ApiIntegrationAwsDetailsAssert {
	return a.HasApiProvider(strings.ToUpper(string(expected)))
}

func (a *ApiIntegrationAwsDetailsAssert) HasApiAwsIamUserArnNotEmpty() *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.ApiAwsIamUserArn == "" {
			return fmt.Errorf("expected api aws iam user arn not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasApiAwsExternalIdNotEmpty() *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.ApiAwsExternalId == "" {
			return fmt.Errorf("expected api aws external id not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasApiKeyNotEmpty() *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.ApiKey == "" {
			return fmt.Errorf("expected api key not empty; got empty")
		}
		return nil
	})
	return a
}
