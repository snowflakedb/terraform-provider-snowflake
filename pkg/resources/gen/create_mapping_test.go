package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_ComputePoolCreateMappings_renderCreateFuncs(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ComputePool}, nil, nil)
	require.True(t, model.Hooks.ParseIdFromConfig.Generate)
	require.True(t, model.Hooks.NewCreateRequest.Generate)
	require.True(t, model.Hooks.ApplyCreateOptionals.Generate)
	require.Equal(t, []string{"name"}, model.IdConfigFields())
	require.Equal(t, "for_application", model.Attributes[0].Name)
	require.Equal(t, "accountObjectIdentifierAttributeCreateBuilder", model.Attributes[0].CreateOptionalCall.Helper)

	require.Equal(t, "id, err := computePoolParseIdFromConfig(d)\n", renderTemplate(t, "parseIdFromConfig", model))
	require.Equal(t, `func computePoolParseIdFromConfig(d *tfschema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(d.Get("name").(string)), nil
}
`, renderTemplate(t, "generatedParseIdFromConfig", model))

	require.Equal(t, "req, err := computePoolNewCreateRequest(d, id)\n", renderTemplate(t, "newCreateRequest", model))
	require.Equal(t, `if err := computePoolApplyCreateOptionals(d, req); err != nil {
		return diag.FromErr(err)
	}
`, renderTemplate(t, "applyCreateOptionals", model))

	require.Equal(t, `func computePoolNewCreateRequest(d *tfschema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateComputePoolRequest, error) {
	minNodes := d.Get("min_nodes").(int)
	maxNodes := d.Get("max_nodes").(int)
	instanceFamily, err := sdk.ToComputePoolInstanceFamily(d.Get("instance_family").(string))
	if err != nil {
		return nil, err
	}
	return sdk.NewCreateComputePoolRequest(id, minNodes, maxNodes, instanceFamily), nil
}
`, renderTemplate(t, "generatedNewCreateRequest", model))

	require.Equal(t, `func computePoolApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateComputePoolRequest) error {
	return errors.Join(
		accountObjectIdentifierAttributeCreateBuilder(d, "for_application", request.WithForApplication),
		listValueCreate(d, "backup_instance_families", &request.BackupInstanceFamilies, toComputePoolBackupInstanceFamilies),
		booleanStringAttributeCreateBuilder(d, "auto_resume", request.WithAutoResume),
		booleanStringAttributeCreateBuilder(d, "initially_suspended", request.WithInitiallySuspended),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend_secs", request.WithAutoSuspendSecs),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))
}

func Test_listEnumMapperName(t *testing.T) {
	require.Equal(t, "toComputePoolBackupInstanceFamilies", listEnumMapperName("ComputePool", defs.ComputePoolAttrs.BackupInstanceFamilies))
}

func Test_ExternalVolumeCreateMappings_renderCreateFuncs(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.ExternalVolumeObject}, nil, nil)
	require.True(t, model.Hooks.ParseIdFromConfig.Generate)
	require.False(t, model.Hooks.NewCreateRequest.Generate)
	require.True(t, model.Hooks.ApplyCreateOptionals.Generate)
	require.Equal(t, []string{"name", "allow_writes", "comment"}, defs.ExternalVolumeObject.ShowOutputKeys())
	require.Equal(t, `func externalVolumeParseIdFromConfig(d *tfschema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(d.Get("name").(string)), nil
}
`, renderTemplate(t, "generatedParseIdFromConfig", model))
	require.Equal(t, "req, err := externalVolumeNewCreateRequestExt(d, id)\n", renderTemplate(t, "newCreateRequest", model))
	require.Equal(t, `func externalVolumeApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateExternalVolumeRequest) error {
	return errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		booleanStringAttributeCreateBuilder(d, "allow_writes", request.WithAllowWrites),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))
}

func Test_SchemaCreateMappings_stayExt(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Schema}, nil, nil)
	require.False(t, model.Hooks.ParseIdFromConfig.Generate)
	require.False(t, model.Hooks.NewCreateRequest.Generate)
	require.False(t, model.Hooks.ApplyCreateOptionals.Generate)
	require.Equal(t, "id, err := schemaParseIdFromConfigExt(d)\n", renderTemplate(t, "parseIdFromConfig", model))
	require.Equal(t, "req, err := schemaNewCreateRequestExt(d, id)\n", renderTemplate(t, "newCreateRequest", model))
	require.Equal(t, `if err := schemaApplyCreateOptionalsExt(d, req); err != nil {
		return diag.FromErr(err)
	}
`, renderTemplate(t, "applyCreateOptionals", model))
	for _, attr := range model.Attributes {
		require.Nil(t, attr.CreateConstructorArg)
		require.Nil(t, attr.CreateOptionalCall)
		require.Nil(t, attr.UpdateCall)
	}
}

func Test_ParseIdFromConfig_databaseAndSchemaIdentifiers(t *testing.T) {
	database := parseIdModel("Schema", "DatabaseObjectIdentifier")
	require.Equal(t, []string{"database", "name"}, database.IdConfigFields())
	require.Equal(t, `func schemaParseIdFromConfig(d *tfschema.ResourceData) (sdk.DatabaseObjectIdentifier, error) {
	return sdk.NewDatabaseObjectIdentifier(d.Get("database").(string), d.Get("name").(string)), nil
}
`, renderTemplate(t, "generatedParseIdFromConfig", database))

	nested := parseIdModel("Service", "SchemaObjectIdentifier")
	require.Equal(t, []string{"database", "schema", "name"}, nested.IdConfigFields())
	require.Equal(t, `func serviceParseIdFromConfig(d *tfschema.ResourceData) (sdk.SchemaObjectIdentifier, error) {
	return sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string)), nil
}
`, renderTemplate(t, "generatedParseIdFromConfig", nested))
}

func Test_attachCreateMappings_skipsNested(t *testing.T) {
	src := []ssot.Attribute{{
		Name:          "client_policy",
		Type:          ssot.AttrSet,
		Required:      true,
		CreateMapping: ssot.StringAttributeCreateBuilder,
		Children: []ssot.Attribute{{
			Name:          "client_type",
			Type:          ssot.AttrString,
			Required:      true,
			CreateMapping: ssot.StringAttributeCreateBuilder,
			Children: []ssot.Attribute{{
				Name:          "version",
				Type:          ssot.AttrString,
				CreateMapping: ssot.StringAttributeCreateBuilder,
			}},
		}},
	}}
	dst := applyResourceAttributeRules("SchemaObjectIdentifier", src)
	require.NotPanics(t, func() {
		attachCreateMappings("AuthenticationPolicy", "SchemaObjectIdentifier", src, dst)
	})
	require.Nil(t, dst[0].CreateConstructorArg)
	require.Nil(t, dst[0].CreateOptionalCall)
	require.Nil(t, dst[0].Children[0].CreateOptionalCall)
	require.Nil(t, dst[0].Children[0].Children[0].CreateOptionalCall)
}

func parseIdModel(name, idType string) ResourceModel {
	return ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:   name,
		IdType: idType,
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithParseIdFromConfig(ssot.Generated()),
		},
	}}, nil, nil)
}

func renderTemplate(t *testing.T, name string, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, LifecycleTemplate.ExecuteTemplate(&b, name, model))
	return b.String()
}
