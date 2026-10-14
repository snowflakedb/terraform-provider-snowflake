package gen

import (
	"strings"
	"testing"

	sdkdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_AssertionsFromSsotObject_emitsStorageAndSdkMethods(t *testing.T) {
	family := ssot.NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	obj := ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithResourceAssert(),
		},
		Attributes: []ssot.Attribute{
			{Name: "instance_family", Type: ssot.AttrEnum, Enum: family, Required: true},
			{Name: "auto_resume", Type: ssot.AttrBool},
			{Name: "for_application", Type: ssot.AttrIdentifier},
			{Name: "min_nodes", Type: ssot.AttrInt, Required: true},
			{Name: "comment", Type: ssot.AttrString},
		},
	}
	logLevel := parameterdefs.ParameterDef{SqlName: "LOG_LEVEL", Kind: "LogLevel"}
	model := assertionsFromSsotObject(obj, []parameterdefs.ParameterDef{logLevel})
	names := attributeNames(model)
	identityFields := 1
	attributeFields := len(obj.Attributes)
	parameterFields := 1
	fullyQualifiedNameFields := 1
	require.Len(t, names, identityFields+attributeFields+parameterFields+fullyQualifiedNameFields)
	require.Equal(t, []string{
		"name",
		"auto_resume",
		"comment",
		"for_application",
		"fully_qualified_name",
		"instance_family",
		"log_level",
		"min_nodes",
	}, names)

	assertions := renderAssertionTemplates(t, model)
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasInstanceFamily(expected string) *PoolResourceAssert {")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasInstanceFamilyString(expected string) *PoolResourceAssert {")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasInstanceFamilySdk(expected sdk.PoolFamily) *PoolResourceAssert {")
	require.Contains(t, assertions, `p.StringValueSet("instance_family", string(expected))`)
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasAutoResume(expected string) *PoolResourceAssert {")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasAutoResumeBool(expected bool) *PoolResourceAssert {")
	require.Contains(t, assertions, `p.StringValueSet("auto_resume", strconv.FormatBool(expected))`)
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasForApplication(expected string) *PoolResourceAssert {")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasForApplicationSdk(expected sdk.AccountObjectIdentifier) *PoolResourceAssert {")
	require.Contains(t, assertions, `p.StringValueSet("for_application", expected.Name())`)
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasMinNodes(expected int) *PoolResourceAssert {")
	require.NotContains(t, assertions, "HasMinNodesSdk")
	require.NotContains(t, assertions, "HasMinNodesBool")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasFullyQualifiedNameString(expected string) *PoolResourceAssert {")
	require.NotContains(t, assertions, "HasFullyQualifiedNameSdk")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasComment(expected string) *PoolResourceAssert {")
	require.NotContains(t, assertions, "HasCommentSdk")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasName(expected string) *PoolResourceAssert {")
	require.NotContains(t, assertions, "HasNameSdk")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasLogLevel(expected string) *PoolResourceAssert {")
	require.Contains(t, assertions, "func (p *PoolResourceAssert) HasLogLevelSdk(expected sdk.LogLevel) *PoolResourceAssert {")
	require.Contains(t, assertions, `p.StringValueSet("log_level", string(expected))`)
}

func Test_AssertionsFromSsotObject_databaseAndSchemaIdentity(t *testing.T) {
	databaseObj := ssot.Object{
		Name:       "Child",
		IdType:     "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithResourceAssert()},
	}
	require.Equal(t, []string{"database", "name", "fully_qualified_name"}, attributeNames(AssertionsFromSsotObject(databaseObj)))

	schemaObj := ssot.Object{
		Name:       "Nested",
		IdType:     "SchemaObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithResourceAssert()},
	}
	require.Equal(t, []string{"database", "schema", "name", "fully_qualified_name"}, attributeNames(AssertionsFromSsotObject(schemaObj)))
}

func Test_ModelFromResourceSchemaDetails_emitsNothingExtra(t *testing.T) {
	model := ModelFromResourceSchemaDetails(genhelpers.ExtractResourceSchemaDetails("Warehouse", map[string]*schema.Schema{
		"warehouse_size": {Type: schema.TypeString, Required: true},
	}), nil, nil)
	assertions := renderAssertionTemplates(t, model)
	require.Contains(t, assertions, "func (w *WarehouseResourceAssert) HasWarehouseSize(expected string) *WarehouseResourceAssert {")
	require.Contains(t, assertions, "func (w *WarehouseResourceAssert) HasWarehouseSizeString(expected string) *WarehouseResourceAssert {")
	require.NotContains(t, assertions, "HasWarehouseSizeSdk")
	require.NotContains(t, assertions, "HasWarehouseSizeBool")
	require.NotContains(t, assertions, "strconv.FormatBool")
}

func Test_AssertionsFromSsotObject_refusesDisabledObject(t *testing.T) {
	require.Panics(t, func() {
		AssertionsFromSsotObject(ssot.Object{Name: "Pool", IdType: "AccountObjectIdentifier"})
	})
}

func Test_AssertionsFromSsotObject_joinsParameterCatalog(t *testing.T) {
	obj := ssot.Object{
		Name:   "Schema",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithResourceAssert(),
		},
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithParametersAttributesSchema(ssot.Generated()),
		},
		Attributes: []ssot.Attribute{
			{Name: "comment", Type: ssot.AttrString},
		},
	}
	model := AssertionsFromSsotObject(obj)
	names := attributeNames(model)
	identityFields := 2
	attributeFields := len(obj.Attributes)
	parameterFields := len(sdkdefs.ParameterDefsForLevel(parameterdefs.ParameterLevelSchema))
	fullyQualifiedNameFields := 1
	require.Len(t, names, identityFields+attributeFields+parameterFields+fullyQualifiedNameFields)
	require.Contains(t, names, "database")
	require.Contains(t, names, "name")
	require.Contains(t, names, "comment")
	require.Contains(t, names, "fully_qualified_name")
	require.Contains(t, names, "log_level")
	require.Contains(t, names, "enable_console_output")
	require.Contains(t, names, "data_retention_time_in_days")
	require.Contains(t, names, "external_volume")

	assertions := renderAssertionTemplates(t, model)
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasLogLevel(expected string) *SchemaResourceAssert {")
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasLogLevelSdk(expected sdk.LogLevel) *SchemaResourceAssert {")
	require.Contains(t, assertions, `s.StringValueSet("log_level", string(expected))`)
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasExternalVolume(expected string) *SchemaResourceAssert {")
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasExternalVolumeSdk(expected sdk.AccountObjectIdentifier) *SchemaResourceAssert {")
	require.Contains(t, assertions, `s.StringValueSet("external_volume", expected.Name())`)
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasEnableConsoleOutput(expected bool) *SchemaResourceAssert {")
	require.NotContains(t, assertions, "HasEnableConsoleOutputBool")
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasDataRetentionTimeInDays(expected int) *SchemaResourceAssert {")
	require.NotContains(t, assertions, "HasDataRetentionTimeInDaysSdk")
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasDefaultDdlCollation(expected string) *SchemaResourceAssert {")
	require.NotContains(t, assertions, "HasDefaultDdlCollationSdk")
	require.Contains(t, assertions, "func (s *SchemaResourceAssert) HasFullyQualifiedNameString(expected string) *SchemaResourceAssert {")
	require.NotContains(t, names, "parameters")
}

func attributeNames(model ResourceAssertionsModel) []string {
	return collections.Map(model.Attributes, func(attr ResourceAttributeAssertionModel) string { return attr.Name })
}

func renderAssertionTemplates(t *testing.T, model ResourceAssertionsModel) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, DefinitionTemplate.Execute(&buf, model))
	require.NoError(t, AssertionsTemplate.Execute(&buf, model))
	return buf.String()
}
