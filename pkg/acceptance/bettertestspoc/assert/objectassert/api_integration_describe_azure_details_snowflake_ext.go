package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ApiIntegrationAzureDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ApiIntegrationAzureDetails, sdk.AccountObjectIdentifier]
}

func ApiIntegrationAzureDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ApiIntegrationAzureDetailsAssert {
	t.Helper()
	return &ApiIntegrationAzureDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("ApiIntegrationAzureDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ApiIntegrationAzureDetails, sdk.AccountObjectIdentifier] {
			return testClient.ApiIntegration.DescribeAzure
		}),
	}
}

func ApiIntegrationAzureDetailsFromObject(t *testing.T, apiIntegrationAzureDetails *sdk.ApiIntegrationAzureDetails) *ApiIntegrationAzureDetailsAssert {
	t.Helper()
	return &ApiIntegrationAzureDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("ApiIntegrationAzureDetails"), apiIntegrationAzureDetails.ID(), apiIntegrationAzureDetails),
	}
}

func (a *ApiIntegrationAzureDetailsAssert) HasId(expected sdk.AccountObjectIdentifier) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.Id.Name() != expected.Name() {
			return fmt.Errorf("expected id: %v; got: %v", expected.Name(), o.Id.Name())
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasEnabled(expected bool) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasApiKey(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.ApiKey != expected {
			return fmt.Errorf("expected api key: %v; got: %v", expected, o.ApiKey)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasApiProvider(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.ApiProvider != expected {
			return fmt.Errorf("expected api provider: %v; got: %v", expected, o.ApiProvider)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureTenantId(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureTenantId != expected {
			return fmt.Errorf("expected azure tenant id: %v; got: %v", expected, o.AzureTenantId)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureAdApplicationId(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureAdApplicationId != expected {
			return fmt.Errorf("expected azure ad application id: %v; got: %v", expected, o.AzureAdApplicationId)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureMultiTenantAppName(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureMultiTenantAppName != expected {
			return fmt.Errorf("expected azure multi tenant app name: %v; got: %v", expected, o.AzureMultiTenantAppName)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureConsentUrl(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureConsentUrl != expected {
			return fmt.Errorf("expected azure consent url: %v; got: %v", expected, o.AzureConsentUrl)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
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

func (a *ApiIntegrationAzureDetailsAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
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

func (a *ApiIntegrationAzureDetailsAssert) HasComment(expected string) *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureMultiTenantAppNameNotEmpty() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureMultiTenantAppName == "" {
			return fmt.Errorf("expected azure multi tenant app name not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasAzureConsentUrlNotEmpty() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.AzureConsentUrl == "" {
			return fmt.Errorf("expected azure consent url not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasApiProviderNotEmpty() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if o.ApiProvider == "" {
			return fmt.Errorf("expected api provider not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationAzureDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationAzureDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationAzureDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
