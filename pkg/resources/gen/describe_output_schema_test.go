package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

// TODO [next PRs]: Production definitions should not be used in tests; replace them with fake objects or remove the test.
func TestDescribeOutputSchema_describeSqlOverride(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.CatalogIntegrationObjectStorage}, nil, nil)
	require.Equal(t, "CATALOG INTEGRATION", model.ObjectType)
	require.True(t, model.Hooks.DescribeOutputSchema.Generate)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, `schemas.DescribeOutputListSchema(sdk.ObjectTypeCatalogIntegration, schemas.DescribeCatalogIntegrationObjectStorageDetailsSchema),`)
	require.NotContains(t, got, "ObjectTypeCatalogIntegrationObjectStorage")
}

// TODO [next PRs]: Production definitions should not be used in tests; replace them with fake objects or remove the test.
func TestDescribeOutputSchema_nameWords(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.Equal(t, "COMPUTE POOL", model.ObjectType)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, `schemas.DescribeOutputListSchema(sdk.ObjectTypeComputePool, schemas.DescribeComputePoolDetailsSchema),`)
}

// TODO [next PRs]: Production definitions should not be used in tests; replace them with fake objects or remove the test.
func TestDescribeOutputSchema_authenticationPolicy(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.AuthenticationPolicy}, nil, nil)
	require.True(t, model.Hooks.DescribeOutputSchema.Generate)
	require.Equal(t, "AUTHENTICATION POLICY", model.ObjectType)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, `schemas.DescribeOutputListSchema(sdk.ObjectTypeAuthenticationPolicy, schemas.DescribeAuthenticationPolicyDetailsSchema),`)
}

// TODO [next PRs]: Production definitions should not be used in tests; replace them with fake objects or remove the test.
func TestDescribeOutputSchema_catalogIntegrationVariant(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.CatalogIntegrationAwsGlue}, nil, nil)
	require.True(t, model.Hooks.DescribeOutputSchema.Generate)
	require.Equal(t, "CATALOG INTEGRATION", model.ObjectType)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, `schemas.DescribeOutputListSchema(sdk.ObjectTypeCatalogIntegration, schemas.DescribeCatalogIntegrationAwsGlueDetailsSchema),`)
	require.NotContains(t, got, "ObjectTypeCatalogIntegrationAwsGlue")
}

func TestDescribeOutputSchema_additionalDescription(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Schema}, nil, nil)
	require.NotEmpty(t, model.AdditionalDescribeOutputDescription)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, `schemas.DescribeOutputListSchema(sdk.ObjectTypeSchema, schemas.DescribeSchemaDetailsSchema, "`+model.AdditionalDescribeOutputDescription+`"),`)
}

func TestFqnSchema(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.FqnSchema.Present)

	got := renderSchemaTemplate(t, "fqnSchema", model)
	require.Contains(t, got, "schemas.FullyQualifiedNameOutputSchema(),")
}

func TestShowOutputSchema(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.CatalogIntegrationAwsGlue}, nil, nil)
	require.True(t, model.Hooks.ShowOutputSchema.Generate)

	got := renderSchemaTemplate(t, "showOutputSchema", model)
	require.Contains(t, got, `schemas.ShowOutputListSchema(sdk.ObjectTypeCatalogIntegration, schemas.ShowCatalogIntegrationSchema),`)
}

func renderSchemaTemplate(t *testing.T, name string, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, SchemaTemplate.ExecuteTemplate(&b, name, model))
	return b.String()
}
