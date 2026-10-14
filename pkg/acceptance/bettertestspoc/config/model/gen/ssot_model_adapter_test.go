package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	sdkdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/stretchr/testify/require"
)

func Test_ModelFromSsotObject_emitsStorageAndSdkMethods(t *testing.T) {
	family := ssot.NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	obj := ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithModel(),
		},
		Attributes: []ssot.Attribute{
			{Name: "instance_family", Type: ssot.AttrEnum, Enum: family, Required: true},
			{Name: "auto_resume", Type: ssot.AttrBool},
			{Name: "for_application", Type: ssot.AttrIdentifier},
			{Name: "min_nodes", Type: ssot.AttrInt, Required: true},
			{Name: "comment", Type: ssot.AttrString},
		},
	}
	model := ModelFromSsotObject(obj)
	require.Equal(t, []string{
		"name",
		"auto_resume",
		"comment",
		"for_application",
		"instance_family",
		"min_nodes",
	}, attributeNames(model))

	builders := renderModelTemplates(t, model)
	require.Contains(t, builders, "func (p *PoolModel) WithInstanceFamily(instanceFamily string) *PoolModel {")
	require.Contains(t, builders, "func (p *PoolModel) WithInstanceFamilySdk(instanceFamily sdk.PoolFamily) *PoolModel {")
	require.Contains(t, builders, "p.InstanceFamily = tfconfig.StringVariable(string(instanceFamily))")
	require.Contains(t, builders, "func (p *PoolModel) WithAutoResume(autoResume string) *PoolModel {")
	require.Contains(t, builders, "func (p *PoolModel) WithAutoResumeBool(autoResume bool) *PoolModel {")
	require.Contains(t, builders, "p.AutoResume = tfconfig.StringVariable(strconv.FormatBool(autoResume))")
	require.Contains(t, builders, "func (p *PoolModel) WithForApplication(forApplication string) *PoolModel {")
	require.Contains(t, builders, "func (p *PoolModel) WithForApplicationSdk(forApplication sdk.AccountObjectIdentifier) *PoolModel {")
	require.Contains(t, builders, "p.ForApplication = tfconfig.StringVariable(forApplication.Name())")
	require.Contains(t, builders, "func (p *PoolModel) WithMinNodes(minNodes int) *PoolModel {")
	require.NotContains(t, builders, "WithMinNodesSdk")
	require.Contains(t, builders, "func (p *PoolModel) WithName(name string) *PoolModel {")
	require.NotContains(t, builders, "WithNameSdk")
	require.NotContains(t, builders, "FullyQualifiedName")
	require.Contains(t, builders, "func (p *PoolModel) WithComment(comment string) *PoolModel {")
	require.NotContains(t, builders, "WithCommentSdk")
	require.Contains(t, builders, "func Pool(")
	require.Contains(t, builders, "instanceFamily string,")
	require.Contains(t, builders, "minNodes int,")
	require.NotContains(t, builders, "instanceFamily sdk.PoolFamily,")
}

func Test_ModelFromSsotObject_databaseAndSchemaIdentity(t *testing.T) {
	databaseObj := ssot.Object{
		Name:       "Child",
		IdType:     "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithModel()},
	}
	require.Equal(t, []string{"database", "name"}, attributeNames(ModelFromSsotObject(databaseObj)))

	schemaObj := ssot.Object{
		Name:       "Nested",
		IdType:     "SchemaObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithModel()},
	}
	require.Equal(t, []string{"database", "schema", "name"}, attributeNames(ModelFromSsotObject(schemaObj)))
}

func Test_ModelFromSsotObject_joinsParameterCatalog(t *testing.T) {
	obj := ssot.Object{
		Name:   "Schema",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithModel(),
		},
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithParametersAttributesSchema(ssot.Generated()),
		},
		Attributes: []ssot.Attribute{
			{Name: "comment", Type: ssot.AttrString},
		},
	}
	model := ModelFromSsotObject(obj)
	names := attributeNames(model)
	// DatabaseObjectIdentifier contributes database and name. Attributes contributes comment.
	identityFields := 2
	attributeFields := len(obj.Attributes)
	parameterFields := len(sdkdefs.ParameterDefsForLevel(parameterdefs.ParameterLevelSchema))
	require.Len(t, names, identityFields+attributeFields+parameterFields)
	require.Contains(t, names, "database")
	require.Contains(t, names, "name")
	require.Contains(t, names, "comment")
	require.Contains(t, names, "log_level")
	require.Contains(t, names, "enable_console_output")
	require.Contains(t, names, "data_retention_time_in_days")
	require.NotContains(t, names, "fully_qualified_name")

	builders := renderModelTemplates(t, model)
	require.Contains(t, builders, "func (s *SchemaModel) WithLogLevel(logLevel string) *SchemaModel {")
	require.Contains(t, builders, "func (s *SchemaModel) WithLogLevelSdk(logLevel sdk.LogLevel) *SchemaModel {")
	require.Contains(t, builders, "func (s *SchemaModel) WithEnableConsoleOutput(enableConsoleOutput bool) *SchemaModel {")
	require.NotContains(t, builders, "WithEnableConsoleOutputBool")
	require.Contains(t, builders, "func (s *SchemaModel) WithDataRetentionTimeInDays(dataRetentionTimeInDays int) *SchemaModel {")
	require.NotContains(t, builders, "WithDataRetentionTimeInDaysSdk")
}

func Test_ModelFromSsotObject_refusesDisabledObject(t *testing.T) {
	require.Panics(t, func() {
		ModelFromSsotObject(ssot.Object{Name: "Pool", IdType: "AccountObjectIdentifier"})
	})
}

func attributeNames(model ResourceConfigBuilderModel) []string {
	names := make([]string, len(model.Attributes))
	for i, attr := range model.Attributes {
		names[i] = attr.Name
	}
	return names
}

func renderModelTemplates(t *testing.T, model ResourceConfigBuilderModel) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, DefinitionTemplate.Execute(&buf, model))
	require.NoError(t, BuildersTemplate.Execute(&buf, model))
	return buf.String()
}
