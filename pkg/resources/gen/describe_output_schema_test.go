package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func TestDescribeOutputSchema_describeSqlOverride(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.CatalogIntegrationObjectStorage}, nil, nil)
	require.Equal(t, "CATALOG INTEGRATION", model.ObjectType)
	require.True(t, model.Hooks.DescribeOutputSchema.Generate)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, "DESCRIBE CATALOG INTEGRATION")
	require.Contains(t, got, "schemas.DescribeCatalogIntegrationObjectStorageDetailsSchema")
	require.NotContains(t, got, "DESCRIBE CATALOG INTEGRATION OBJECT STORAGE")
}

func TestDescribeOutputSchema_nameWords(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.Equal(t, "COMPUTE POOL", model.ObjectType)

	got := renderSchemaTemplate(t, "describeOutputSchema", model)
	require.Contains(t, got, "DESCRIBE COMPUTE POOL")
	require.Contains(t, got, "schemas.DescribeComputePoolDetailsSchema")
}

func renderSchemaTemplate(t *testing.T, name string, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, SchemaTemplate.ExecuteTemplate(&b, name, model))
	return b.String()
}
