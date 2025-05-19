package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (a *FunctionAssert) HasCreatedOnNotEmpty() *FunctionAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Function) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *FunctionAssert) HasExternalAccessIntegrationsNil() *FunctionAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Function) error {
		t.Helper()
		if o.ExternalAccessIntegrations != nil {
			return fmt.Errorf("expected external_access_integrations to be nil but was: %v", *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return a
}

func (f *FunctionAssert) HasExactlyExternalAccessIntegrations(integrations ...sdk.AccountObjectIdentifier) *FunctionAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.Function) error {
		t.Helper()
		if o.ExternalAccessIntegrations == nil {
			return fmt.Errorf("expected external access integrations to have value; got: nil")
		}
		joined := strings.Join(collections.Map(integrations, func(ex sdk.AccountObjectIdentifier) string { return ex.FullyQualifiedName() }), ",")
		expected := fmt.Sprintf(`[%s]`, joined)
		if *o.ExternalAccessIntegrations != expected {
			return fmt.Errorf("expected external access integrations: %v; got: %v", expected, *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return f
}

func (f *FunctionAssert) HasExactlySecrets(expectedSecrets map[string]sdk.SchemaObjectIdentifier) *FunctionAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.Function) error {
		t.Helper()
		if o.Secrets == nil {
			return fmt.Errorf("expected secrets to have value; got: nil")
		}
		var parts []string
		for k, v := range expectedSecrets {
			parts = append(parts, fmt.Sprintf(`"%s":"\"%s\".\"%s\".%s"`, k, v.DatabaseName(), v.SchemaName(), v.Name()))
		}
		expected := fmt.Sprintf(`{%s}`, strings.Join(parts, ","))
		if *o.Secrets != expected {
			return fmt.Errorf("expected secrets: %v; got: %v", expected, *o.Secrets)
		}
		return nil
	})
	return f
}

func (f *FunctionAssert) HasArgumentsRawContains(substring string) *FunctionAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.Function) error {
		t.Helper()
		if !strings.Contains(o.ArgumentsRaw, substring) {
			return fmt.Errorf("expected arguments raw contain: %v, to contain: %v", o.ArgumentsRaw, substring)
		}
		return nil
	})
	return f
}
