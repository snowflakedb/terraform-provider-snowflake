package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (a *ProcedureAssert) HasCreatedOnNotEmpty() *ProcedureAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return a
}

func (a *ProcedureAssert) HasExternalAccessIntegrationsNil() *ProcedureAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
		t.Helper()
		if o.ExternalAccessIntegrations != nil {
			return fmt.Errorf("expected external_access_integrations to be nil but was: %v", *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return a
}

func (a *ProcedureAssert) HasSecretsNil() *ProcedureAssert {
	a.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
		t.Helper()
		if o.Secrets != nil {
			return fmt.Errorf("expected secrets to be nil but was: %v", *o.Secrets)
		}
		return nil
	})
	return a
}

func (f *ProcedureAssert) HasExactlyExternalAccessIntegrations(integrations ...sdk.AccountObjectIdentifier) *ProcedureAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
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

func (p *ProcedureAssert) HasArgumentsRawContains(substring string) *ProcedureAssert {
	p.AddAssertion(func(t *testing.T, o *sdk.Procedure) error {
		t.Helper()
		if !strings.Contains(o.ArgumentsRaw, substring) {
			return fmt.Errorf("expected arguments raw contain: %v, to contain: %v", o.ArgumentsRaw, substring)
		}
		return nil
	})
	return p
}
