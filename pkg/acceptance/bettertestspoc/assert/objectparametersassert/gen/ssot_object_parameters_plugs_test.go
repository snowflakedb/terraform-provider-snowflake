package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_mergeSnowflakeObjectParameters_skipsHandWrittenNames(t *testing.T) {
	handWritten := []SnowflakeObjectParameters{
		{Name: "Alpha"},
		{Name: "Beta"},
		{Name: "Gamma", ParameterConstantPrefix: "HandWritten"},
	}
	ssotParameters := []SnowflakeObjectParameters{
		{Name: "Gamma", ParameterConstantPrefix: "FromSsot"},
	}

	merged := mergeSnowflakeObjectParameters(handWritten, ssotParameters)
	require.Equal(t, []string{"Alpha", "Beta", "Gamma"}, parameterObjectNames(merged))
	require.Equal(t, "FromSsot", merged[2].ParameterConstantPrefix)
}

func Test_mergeSnowflakeObjectParameters_panicsOnDuplicateSsotName(t *testing.T) {
	require.Panics(t, func() {
		mergeSnowflakeObjectParameters(nil, []SnowflakeObjectParameters{{Name: "Alpha"}, {Name: "Alpha"}})
	})
}

func Test_objectParametersFromSsot_schema(t *testing.T) {
	g := ssot.ObjectGenerators(defs.Schema.Generators...)

	objectAssertions := objectParametersFromSsot(defs.Schema, g.ObjectParametersAssert.Enabled)
	resourceAssertions := objectParametersFromSsot(defs.Schema, g.ResourceParametersAssert.Enabled)

	require.Equal(t, objectAssertions, resourceAssertions)
	require.Equal(t, "Schema", objectAssertions.Name)
	require.Equal(t, "sdk.DatabaseObjectIdentifier", objectAssertions.IdType)
	require.Empty(t, objectAssertions.ParameterConstantPrefix)
	require.Empty(t, objectAssertions.ObjectTypeName)
	require.Len(t, objectAssertions.Parameters, len(defs.ParameterDefsForLevel(ssot.ParameterLevelSchema)))
}

func Test_objectParametersFromSsot_variantUsesBaseNames(t *testing.T) {
	base := ssot.Object{Name: "Warehouse"}
	obj := ssot.Object{
		Name:       "WarehouseAdaptive",
		IdType:     "AccountObjectIdentifier",
		VariantOf:  &base,
		Generators: []ssot.GeneratorOption{ssot.WithObjectParametersAssert()},
	}
	got := objectParametersFromSsot(obj, true)
	require.Equal(t, "WarehouseAdaptive", got.Name)
	require.Equal(t, "Warehouse", got.ParameterConstantPrefix)
	require.Equal(t, "Warehouse", got.ObjectTypeName)
}

func Test_objectParametersFromSsot_panicsWhenNotEnabled(t *testing.T) {
	require.Panics(t, func() {
		objectParametersFromSsot(ssot.Object{Name: "Schema", IdType: "DatabaseObjectIdentifier"}, false)
	})
}

func Test_schemaIsOptedIntoBothParameterAssertionGenerators(t *testing.T) {
	require.Contains(t, parameterObjectNames(SnowflakeObjectParametersForAssertions()), "Schema")
	require.Contains(t, parameterObjectNames(SnowflakeResourceParametersForAssertions()), "Schema")
}

func parameterObjectNames(parameters []SnowflakeObjectParameters) []string {
	return collections.Map(parameters, func(p SnowflakeObjectParameters) string { return p.ObjectName() })
}
