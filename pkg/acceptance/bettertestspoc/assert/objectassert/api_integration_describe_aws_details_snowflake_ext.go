package objectassert

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApiIntegrationAwsDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ApiIntegrationAwsDetails, sdk.AccountObjectIdentifier]
}

func ApiIntegrationAwsDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ApiIntegrationAwsDetailsAssert {
	t.Helper()
	return &ApiIntegrationAwsDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("ApiIntegrationAwsDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ApiIntegrationAwsDetails, sdk.AccountObjectIdentifier] {
			return testClient.ApiIntegration.DescribeAws
		}),
	}
}

func ApiIntegrationAwsDetailsFromObject(t *testing.T, apiIntegrationAwsDetails *sdk.ApiIntegrationAwsDetails) *ApiIntegrationAwsDetailsAssert {
	t.Helper()
	return &ApiIntegrationAwsDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("ApiIntegrationAwsDetails"), apiIntegrationAwsDetails.ID(), apiIntegrationAwsDetails),
	}
}

func (a *ApiIntegrationAwsDetailsAssert) HasEnabled(expected bool) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasApiKey(expected string) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.ApiKey != expected {
			return fmt.Errorf("expected api key: %v; got: %v", expected, o.ApiKey)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasApiProvider(expected sdk.ApiIntegrationAwsApiProviderType) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.ApiProvider != strings.ToUpper(string(expected)) {
			return fmt.Errorf("expected api provider: %v; got: %v", expected, o.ApiProvider)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasApiAwsRoleArn(expected string) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.ApiAwsRoleArn != expected {
			return fmt.Errorf("expected api aws role arn: %v; got: %v", expected, o.ApiAwsRoleArn)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		mapped := collections.Map(o.AllowedPrefixes, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected allowed prefixes: %v; got: %v", expected, o.AllowedPrefixes)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		mapped := collections.Map(o.BlockedPrefixes, func(item string) any { return item })
		mappedExpected := collections.Map(expected, func(item string) any { return item })
		if !slices.Equal(mapped, mappedExpected) {
			return fmt.Errorf("expected blocked prefixes: %v; got: %v", expected, o.BlockedPrefixes)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAwsDetailsAssert) HasComment(expected string) *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return a
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

func (a *ApiIntegrationAwsDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationAwsDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAwsDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
