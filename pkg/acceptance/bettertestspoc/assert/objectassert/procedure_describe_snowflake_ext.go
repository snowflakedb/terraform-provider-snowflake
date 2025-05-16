package objectassert

import (
	"fmt"
	"strings"
	"testing"

	assert2 "github.com/stretchr/testify/assert"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

// TODO [SNOW-1501905]: this file should be fully regenerated when adding and option to assert the results of describe
type ProcedureDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ProcedureDetails, sdk.SchemaObjectIdentifierWithArguments]
}

func ProcedureDetails(t *testing.T, id sdk.SchemaObjectIdentifierWithArguments) *ProcedureDetailsAssert {
	t.Helper()
	return &ProcedureDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectType("PROCEDURE_DETAILS"), id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ProcedureDetails, sdk.SchemaObjectIdentifierWithArguments] {
			return testClient.Procedure.DescribeDetails
		}),
	}
}

func (f *ProcedureDetailsAssert) HasSignature(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Signature != expected {
			return fmt.Errorf("expected signature: %v; got: %v", expected, o.Signature)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasReturns(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Returns != expected {
			return fmt.Errorf("expected returns: %v; got: %v", expected, o.Returns)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasLanguage(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Language != expected {
			return fmt.Errorf("expected language: %v; got: %v", expected, o.Language)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasBody(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Body == nil {
			return fmt.Errorf("expected body to have value; got: nil")
		}
		if *o.Body != expected {
			return fmt.Errorf("expected body: %v; got: %v", expected, *o.Body)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasNullHandling(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NullHandling == nil {
			return fmt.Errorf("expected null handling to have value; got: nil")
		}
		if *o.NullHandling != expected {
			return fmt.Errorf("expected null handling: %v; got: %v", expected, *o.NullHandling)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasVolatility(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Volatility == nil {
			return fmt.Errorf("expected volatility to have value; got: nil")
		}
		if *o.Volatility != expected {
			return fmt.Errorf("expected volatility: %v; got: %v", expected, *o.Volatility)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasExternalAccessIntegrations(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.ExternalAccessIntegrations == nil {
			return fmt.Errorf("expected external access integrations to have value; got: nil")
		}
		if *o.ExternalAccessIntegrations != expected {
			return fmt.Errorf("expected external access integrations: %v; got: %v", expected, *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasSecrets(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Secrets == nil {
			return fmt.Errorf("expected secrets to have value; got: nil")
		}
		if *o.Secrets != expected {
			return fmt.Errorf("expected secrets: %v; got: %v", expected, *o.Secrets)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasImports(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Imports == nil {
			return fmt.Errorf("expected imports to have value; got: nil")
		}
		if *o.Imports != expected {
			return fmt.Errorf("expected imports: %v; got: %v", expected, *o.Imports)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasHandler(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Handler == nil {
			return fmt.Errorf("expected handler to have value; got: nil")
		}
		if *o.Handler != expected {
			return fmt.Errorf("expected handler: %v; got: %v", expected, *o.Handler)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasRuntimeVersion(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.RuntimeVersion == nil {
			return fmt.Errorf("expected runtime version to have value; got: nil")
		}
		if *o.RuntimeVersion != expected {
			return fmt.Errorf("expected runtime version: %v; got: %v", expected, *o.RuntimeVersion)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasPackages(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Packages == nil {
			return fmt.Errorf("expected packages to have value; got: nil")
		}
		if *o.Packages != expected {
			return fmt.Errorf("expected packages: %v; got: %v", expected, *o.Packages)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasTargetPath(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.TargetPath == nil {
			return fmt.Errorf("expected target path to have value; got: nil")
		}
		if *o.TargetPath != expected {
			return fmt.Errorf("expected target path: %v; got: %v", expected, *o.TargetPath)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasInstalledPackages(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.InstalledPackages == nil {
			return fmt.Errorf("expected installed packages to have value; got: nil")
		}
		if *o.InstalledPackages != expected {
			return fmt.Errorf("expected installed packages: %v; got: %v", expected, *o.InstalledPackages)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasExecuteAs(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.ExecuteAs != expected {
			return fmt.Errorf("expected execute as: %v; got: %v", expected, o.ExecuteAs)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasBodyNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Body != nil {
			return fmt.Errorf("expected body to be nil, was %v", *o.Body)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasNullHandlingNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.NullHandling != nil {
			return fmt.Errorf("expected null handling to be nil, was %v", *o.NullHandling)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasVolatilityNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Volatility != nil {
			return fmt.Errorf("expected volatility to be nil, was %v", *o.Volatility)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasExternalAccessIntegrationsNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.ExternalAccessIntegrations != nil {
			return fmt.Errorf("expected external access integrations to be nil, was %v", *o.ExternalAccessIntegrations)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasSecretsNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Secrets != nil {
			return fmt.Errorf("expected secrets to be nil, was %v", *o.Secrets)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasImportsNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Imports != nil {
			return fmt.Errorf("expected imports to be nil, was %v", *o.Imports)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasHandlerNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Handler != nil {
			return fmt.Errorf("expected handler to be nil, was %v", *o.Handler)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasRuntimeVersionNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.RuntimeVersion != nil {
			return fmt.Errorf("expected runtime version to be nil, was %v", *o.RuntimeVersion)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasPackagesNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.Packages != nil {
			return fmt.Errorf("expected packages to be nil, was %v", *o.Packages)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasTargetPathNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.TargetPath != nil {
			return fmt.Errorf("expected target path to be nil, was %v", *o.TargetPath)
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasInstalledPackagesNil() *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.InstalledPackages != nil {
			return fmt.Errorf("expected installed packages to be nil, was %v", *o.InstalledPackages)
		}
		return nil
	})
	return f
}

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

func (f *ProcedureDetailsAssert) HasReturnDataType(expectedDataType datatypes.DataType) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.ReturnDataType == nil {
			return fmt.Errorf("expected return data type to have value; got: nil")
		}
		if !datatypes.AreTheSame(o.ReturnDataType, expectedDataType) {
			return fmt.Errorf("expected %s return data type, got %v", expectedDataType, o.ReturnDataType.ToSql())
		}
		return nil
	})
	return f
}

func (f *ProcedureDetailsAssert) HasReturnNotNull(expected bool) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.ReturnNotNull != expected {
			return fmt.Errorf("expected return not null %t; got: %t", expected, o.ReturnNotNull)
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

func (f *ProcedureDetailsAssert) HasSnowparkVersion(expected string) *ProcedureDetailsAssert {
	f.AddAssertion(func(t *testing.T, o *sdk.ProcedureDetails) error {
		t.Helper()
		if o.SnowparkVersion != expected {
			return fmt.Errorf("expected snowpark version %s; got: %s", expected, o.SnowparkVersion)
		}
		return nil
	})
	return f
}
