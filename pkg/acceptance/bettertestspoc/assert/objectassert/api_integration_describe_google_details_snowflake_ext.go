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

type ApiIntegrationGoogleDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ApiIntegrationGoogleDetails, sdk.AccountObjectIdentifier]
}

func ApiIntegrationGoogleDetails(t *testing.T, id sdk.AccountObjectIdentifier) *ApiIntegrationGoogleDetailsAssert {
	t.Helper()
	return &ApiIntegrationGoogleDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("ApiIntegrationGoogleDetails"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ApiIntegrationGoogleDetails, sdk.AccountObjectIdentifier] {
			return testClient.ApiIntegration.DescribeGoogle
		}),
	}
}

func ApiIntegrationGoogleDetailsFromObject(t *testing.T, apiIntegrationGoogleDetails *sdk.ApiIntegrationGoogleDetails) *ApiIntegrationGoogleDetailsAssert {
	t.Helper()
	return &ApiIntegrationGoogleDetailsAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectType("ApiIntegrationGoogleDetails"), apiIntegrationGoogleDetails.ID(), apiIntegrationGoogleDetails),
	}
}

func (a *ApiIntegrationGoogleDetailsAssert) HasId(expected sdk.AccountObjectIdentifier) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.Id.Name() != expected.Name() {
			return fmt.Errorf("expected id: %v; got: %v", expected.Name(), o.Id.Name())
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasEnabled(expected bool) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.Enabled != expected {
			return fmt.Errorf("expected enabled: %v; got: %v", expected, o.Enabled)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasApiKeyEmpty() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.ApiKey != "" {
			return fmt.Errorf("expected api key: %v to be empty", o.ApiKey)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasApiKeyNotEmpty() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.ApiKey == "" {
			return fmt.Errorf("expected api key not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasApiProvider(expected sdk.ApiIntegrationGoogleApiProviderType) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.ApiProvider != strings.ToUpper(string(expected)) {
			return fmt.Errorf("expected api provider: %v; got: %v", expected, o.ApiProvider)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasGoogleAudience(expected string) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.GoogleAudience != expected {
			return fmt.Errorf("expected google audience: %v; got: %v", expected, o.GoogleAudience)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasAllowedPrefixes(expected ...string) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
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

func (a *ApiIntegrationGoogleDetailsAssert) HasBlockedPrefixes(expected ...string) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
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

func (a *ApiIntegrationGoogleDetailsAssert) HasComment(expected string) *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasGoogleApiServiceAccountNotEmpty() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if o.GoogleApiServiceAccount == "" {
			return fmt.Errorf("expected google api service account not empty; got empty")
		}
		return nil
	})
	return a
}

func (a *ApiIntegrationGoogleDetailsAssert) HasNoBlockedPrefixes() *ApiIntegrationGoogleDetailsAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.ApiIntegrationGoogleDetails) error {
		t.Helper()
		if len(o.BlockedPrefixes) != 0 {
			return fmt.Errorf("expected no blocked prefixes; got: %v", o.BlockedPrefixes)
		}
		return nil
	})
	return a
}
