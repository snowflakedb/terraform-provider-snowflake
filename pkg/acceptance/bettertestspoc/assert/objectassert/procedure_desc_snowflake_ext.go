package objectassert

import (
	"fmt"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (f *ProcedureDetailsAssert) HasInstalledPackagesNotEmpty() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.InstalledPackages == nil {
			return fmt.Errorf("expected installed packages to not be nil")
		}
		if *o.InstalledPackages == "" {
			return fmt.Errorf("expected installed packages to not be empty")
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasExactlyExternalAccessIntegrations(integrations ...sdk.AccountObjectIdentifier) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
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

func (f *ProcedureDetailsAssert) HasExactlySecrets(expectedSecrets map[string]sdk.SchemaObjectIdentifier) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
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

func (f *ProcedureDetailsAssert) HasExactlyImportsNormalizedInAnyOrder(imports ...sdk.NormalizedPath) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NormalizedImports == nil {
			return fmt.Errorf("expected imports to have value; got: nil")
		}
		if !assert2.ElementsMatch(t, imports, o.NormalizedImports) {
			return fmt.Errorf("expected %v imports, got %v", imports, o.NormalizedImports)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasNormalizedTargetPath(expectedStageLocation string, expectedPathOnStage string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NormalizedTargetPath == nil {
			return fmt.Errorf("expected normalized target path to have value; got: nil")
		}
		if o.NormalizedTargetPath.StageLocation != expectedStageLocation {
			return fmt.Errorf("expected %s stage location for target path, got %v", expectedStageLocation, o.NormalizedTargetPath.StageLocation)
		}
		if o.NormalizedTargetPath.PathOnStage != expectedPathOnStage {
			return fmt.Errorf("expected %s path on stage for target path, got %v", expectedPathOnStage, o.NormalizedTargetPath.PathOnStage)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasNormalizedTargetPathNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NormalizedTargetPath != nil {
			return fmt.Errorf("expected normalized target path to be nil, got: %s", *o.NormalizedTargetPath)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasExactlyExternalAccessIntegrationsNormalizedInAnyOrder(integrations ...sdk.AccountObjectIdentifier) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NormalizedExternalAccessIntegrations == nil {
			return fmt.Errorf("expected normalized external access integrations to have value; got: nil")
		}
		fullyQualifiedNamesExpected := collections.Map(integrations, func(id sdk.AccountObjectIdentifier) string { return id.FullyQualifiedName() })
		fullyQualifiedNamesGot := collections.Map(o.NormalizedExternalAccessIntegrations, func(id sdk.AccountObjectIdentifier) string { return id.FullyQualifiedName() })
		if !assert2.ElementsMatch(t, fullyQualifiedNamesExpected, fullyQualifiedNamesGot) {
			return fmt.Errorf("expected %v normalized external access integrations, got %v", integrations, o.NormalizedExternalAccessIntegrations)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) ContainsExactlySecrets(secrets map[string]sdk.SchemaObjectIdentifier) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NormalizedSecrets == nil {
			return fmt.Errorf("expected normalized secrets to have value; got: nil")
		}
		for k, v := range secrets {
			if s, ok := o.NormalizedSecrets[k]; !ok {
				return fmt.Errorf("expected normalized secrets to have a secret associated with key %s", k)
			} else if s.FullyQualifiedName() != v.FullyQualifiedName() {
				return fmt.Errorf("expected secret with key %s to have id %s, got %s", k, v.FullyQualifiedName(), s.FullyQualifiedName())
			}
		}
		for k := range o.NormalizedSecrets {
			if _, ok := secrets[k]; !ok {
				return fmt.Errorf("normalized secrets have unexpected key: %s", k)
			}
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasExactlyPackagesInAnyOrder(packages ...string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NormalizedPackages == nil {
			return fmt.Errorf("expected packages to have value; got: nil")
		}
		if !assert2.ElementsMatch(t, packages, o.NormalizedPackages) {
			return fmt.Errorf("expected %v packages, got %v", packages, o.NormalizedPackages)
		}
		return nil
	})
	return f
}
