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
`, renderTemplate(t, "generatedCreateRequest", model.CreateRequestFunc()))

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

func Test_constructorArg(t *testing.T) {
	object := createScope{stem: "ComputePool"}
	block := createScope{stem: "WidgetPolicy", inBlock: true}
	family := ssot.NewEnum("ComputePoolInstanceFamily", "ComputePoolInstanceFamilies", "CPU_X64_XS")
	arg := func(scope createScope, attr ssot.Attribute) (*CreateConstructorArg, bool) {
		return constructorArg(scope, attr, applyResourceAttributeRules("", []ssot.Attribute{attr})[0])
	}

	got, ok := arg(object, ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt, Required: true})
	require.True(t, ok)
	require.Equal(t, &CreateConstructorArg{Var: "minNodes", Get: `d.Get("min_nodes").(int)`}, got)

	got, ok = arg(block, ssot.Attribute{Name: "instance_family", Type: ssot.AttrEnum, Required: true, Enum: family})
	require.True(t, ok)
	require.Equal(t, &CreateConstructorArg{Var: "instanceFamily", Get: `d.Get(prefix+"instance_family").(string)`, EnumName: "ComputePoolInstanceFamily"}, got)

	_, ok = arg(object, ssot.Attribute{Name: "comment", Type: ssot.AttrString})
	require.False(t, ok)
	_, ok = arg(object, ssot.Attribute{Name: "comment", Type: ssot.AttrString, Required: true, CreateMapping: ssot.StringAttributeCreateBuilder})
	require.False(t, ok)

	items := ssot.Attribute{Name: "items", Type: ssot.AttrList, Required: true, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	_, ok = arg(block, items)
	require.False(t, ok)
	require.PanicsWithValue(t, `object "ComputePool": required attribute "items" has no CreateMapping`, func() {
		arg(object, items)
	})
}

func Test_collectionMapperName(t *testing.T) {
	require.Equal(t, "toComputePoolBackupInstanceFamilies", collectionMapperName("ComputePool", defs.ComputePoolAttrs.BackupInstanceFamilies))
}

func Test_AuthenticationPolicySecurityIntegrationsCreateMapping(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.AuthenticationPolicy}, nil, nil)
	var got *CreateOptionalCall
	for _, attr := range model.Attributes {
		if attr.Name == "security_integrations" {
			got = attr.CreateOptionalCall
		}
	}
	require.Equal(t, &CreateOptionalCall{
		Helper: "attributeMappedValueCreateBuilder",
		Key:    `"security_integrations"`,
		Target: "request.WithSecurityIntegrations, ToSecurityIntegrationsRequest",
	}, got)
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
	require.True(t, model.Hooks.ParseIdFromConfig.Generate)
	require.False(t, model.Hooks.NewCreateRequest.Generate)
	require.False(t, model.Hooks.ApplyCreateOptionals.Generate)
	require.Equal(t, "id, err := schemaParseIdFromConfig(d)\n", renderTemplate(t, "parseIdFromConfig", model))
	require.Equal(t, `func schemaParseIdFromConfig(d *tfschema.ResourceData) (sdk.DatabaseObjectIdentifier, error) {
	return sdk.NewDatabaseObjectIdentifier(d.Get("database").(string), d.Get("name").(string)), nil
}
`, renderTemplate(t, "generatedParseIdFromConfig", model))
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

func Test_attachCreateMappings_unsupportedBlocksStayHandwritten(t *testing.T) {
	src := []ssot.Attribute{{
		Name:     "storage_location",
		Type:     ssot.AttrList,
		Required: true,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "url",
			Type: ssot.AttrString,
		}}),
	}, {
		Name: "policy",
		Type: ssot.AttrList,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "comment",
			Type: ssot.AttrString,
		}, {
			Name:     "scopes",
			Type:     ssot.AttrList,
			Required: true,
			Elem:     &ssot.Attribute{Type: ssot.AttrString},
		}}),
	}}
	dst := applyResourceAttributeRules("SchemaObjectIdentifier", src)
	require.NotPanics(t, func() {
		attachCreateMappings("Widget", src, dst)
	})
	for _, attr := range dst {
		require.Nil(t, attr.CreateConstructorArg, attr.Name)
		require.Nil(t, attr.CreateOptionalCall, attr.Name)
		require.Nil(t, attr.CreateBlock, attr.Name)
		for _, child := range attr.Children {
			require.Nil(t, child.CreateConstructorArg, child.Name)
			require.Nil(t, child.CreateOptionalCall, child.Name)
			require.Nil(t, child.CreateBlock, child.Name)
		}
	}
}

func Test_attachCreateMappings_collectionsUseElementMappers(t *testing.T) {
	mode := ssot.NewEnum("WidgetMode", "WidgetModes", "ON")
	model := createOptionalsModel("Widget", []ssot.Attribute{{
		Name: "modes",
		Type: ssot.AttrSet,
		Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: mode},
	}, {
		Name: "labels",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrString},
	}, {
		Name: "client_policy",
		Type: ssot.AttrSet,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:     "minimum_version",
			Type:     ssot.AttrString,
			Required: true,
		}}),
	}, {
		Name:     "policy",
		Type:     ssot.AttrList,
		MaxItems: new(1),
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "issuers",
			Type: ssot.AttrSet,
			Elem: &ssot.Attribute{Type: ssot.AttrString},
		}}),
	}})

	require.Equal(t, `func widgetApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateWidgetRequest) error {
	return errors.Join(
		setMappedValueCreateBuilder(d, "modes", request.WithModes, toWidgetModes),
		listValueCreate(d, "labels", &request.Labels, toWidgetLabels),
		setMappedValueCreateBuilder(d, "client_policy", request.WithClientPolicy, toWidgetClientPolicy),
		singleBlockAttributeCreate(d, "policy", request.WithPolicy, parseWidgetPolicy),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))
	blocks := model.CreateBlockFuncs()
	require.Len(t, blocks, 1)
	require.Equal(t, []CreateOptionalCall{{
		Helper: "setMappedValueCreateBuilder",
		Key:    `prefix+"issuers"`,
		Target: "request.WithIssuers, toWidgetPolicyIssuers",
	}}, blocks[0].Calls)
}

func Test_attachCreateMappings_nestedListBlocks(t *testing.T) {
	mode := ssot.NewEnum("WidgetMode", "WidgetModes", "ON")
	level := ssot.NewEnum("WidgetLevel", "WidgetLevels", "LOW")
	model := createOptionalsModel("Widget", []ssot.Attribute{{
		Name:     "policy",
		Type:     ssot.AttrList,
		MaxItems: new(1),
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:     "mode",
			Type:     ssot.AttrEnum,
			Required: true,
			Enum:     mode,
		}, {
			Name:     "label",
			Type:     ssot.AttrString,
			Required: true,
		}, {
			Name: "comment",
			Type: ssot.AttrString,
		}, {
			Name: "level",
			Type: ssot.AttrEnum,
			Enum: level,
		}, {
			Name: "kinds",
			Type: ssot.AttrList,
			Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: mode},
		}, {
			Name:           "integration",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "AccountObjectIdentifier",
		}, {
			Name:     "limit",
			Type:     ssot.AttrList,
			MaxItems: new(1),
			Elem: ssot.ObjectElem([]ssot.Attribute{{
				Name: "count",
				Type: ssot.AttrInt,
			}, {
				Name: "enabled",
				Type: ssot.AttrBool,
			}}),
		}}),
	}, {
		Name: "storage_location",
		Type: ssot.AttrList,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:     "name",
			Type:     ssot.AttrString,
			Required: true,
		}, {
			Name: "url",
			Type: ssot.AttrString,
		}}),
	}, {
		Name: "rule",
		Type: ssot.AttrList,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:     "priority",
			Type:     ssot.AttrInt,
			Required: true,
		}}),
	}})

	require.Equal(t, `func widgetApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateWidgetRequest) error {
	return errors.Join(
		singleBlockAttributeCreate(d, "policy", request.WithPolicy, parseWidgetPolicy),
		listBlockAttributeCreate(d, "storage_location", &request.StorageLocation, parseWidgetStorageLocation),
		listBlockAttributeCreate(d, "rule", &request.Rule, parseWidgetRule),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))

	blocks := model.CreateBlockFuncs()
	require.Len(t, blocks, 4)
	require.Equal(t, `func parseWidgetPolicy(d *tfschema.ResourceData, prefix string) (sdk.WidgetPolicyRequest, error) {
	mode, err := sdk.ToWidgetMode(d.Get(prefix+"mode").(string))
	if err != nil {
		return sdk.WidgetPolicyRequest{}, err
	}
	label := d.Get(prefix+"label").(string)
	request := sdk.NewWidgetPolicyRequest(mode, label)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, prefix+"comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, prefix+"level", request.WithLevel, sdk.ToWidgetLevel),
		listValueCreate(d, prefix+"kinds", &request.Kinds, toWidgetPolicyKinds),
		accountObjectIdentifierAttributeCreateBuilder(d, prefix+"integration", request.WithIntegration),
		singleBlockAttributeCreate(d, prefix+"limit", request.WithLimit, parseWidgetPolicyLimit),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[0]))
	require.Equal(t, `func parseWidgetPolicyLimit(d *tfschema.ResourceData, prefix string) (sdk.WidgetPolicyLimitRequest, error) {
	request := sdk.NewWidgetPolicyLimitRequest()
	errs := errors.Join(
		intAttributeCreateBuilder(d, prefix+"count", request.WithCount),
		booleanStringAttributeCreateBuilder(d, prefix+"enabled", request.WithEnabled),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[1]))
	require.Equal(t, `func parseWidgetStorageLocation(d *tfschema.ResourceData, prefix string) (sdk.WidgetStorageLocationRequest, error) {
	name := d.Get(prefix+"name").(string)
	request := sdk.NewWidgetStorageLocationRequest(name)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, prefix+"url", request.WithUrl),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[2]))
	require.Equal(t, `func parseWidgetRule(d *tfschema.ResourceData, prefix string) (sdk.WidgetRuleRequest, error) {
	priority := d.Get(prefix+"priority").(int)
	request := sdk.NewWidgetRuleRequest(priority)
	return *request, nil
}
`, renderTemplate(t, "generatedCreateRequest", blocks[3]))
}

func Test_attachCreateMappings_explicitBlockMapping(t *testing.T) {
	src := []ssot.Attribute{{
		Name:          "policy",
		Type:          ssot.AttrSet,
		CreateMapping: ssot.ListValueCreate("toWidgetPolicy"),
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "comment",
			Type: ssot.AttrString,
		}}),
	}}
	dst := applyResourceAttributeRules("", src)
	attachCreateMappings("Widget", src, dst)
	require.Nil(t, dst[0].CreateBlock)
	require.Nil(t, dst[0].Children[0].CreateOptionalCall)
	require.Equal(t, &CreateOptionalCall{
		Helper: "listValueCreate",
		Key:    `"policy"`,
		Target: "&request.Policy, toWidgetPolicy",
	}, dst[0].CreateOptionalCall)
}

func Test_attachCreateMappings_fourLevelBlocks(t *testing.T) {
	model := createOptionalsModel("Widget", []ssot.Attribute{fourLevelBlock(nil)})

	require.Equal(t, `func widgetApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateWidgetRequest) error {
	return errors.Join(
		singleBlockAttributeCreate(d, "a", request.WithA, parseWidgetA),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))

	blocks := model.CreateBlockFuncs()
	require.Len(t, blocks, 4)
	require.Equal(t, `func parseWidgetA(d *tfschema.ResourceData, prefix string) (sdk.WidgetARequest, error) {
	request := sdk.NewWidgetARequest()
	errs := errors.Join(
		listBlockAttributeCreate(d, prefix+"b", &request.B, parseWidgetAB),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[0]))
	require.Equal(t, `func parseWidgetAB(d *tfschema.ResourceData, prefix string) (sdk.WidgetABRequest, error) {
	request := sdk.NewWidgetABRequest()
	errs := errors.Join(
		singleBlockAttributeCreate(d, prefix+"c", request.WithC, parseWidgetABC),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[1]))
	require.Equal(t, `func parseWidgetABC(d *tfschema.ResourceData, prefix string) (sdk.WidgetABCRequest, error) {
	request := sdk.NewWidgetABCRequest()
	errs := errors.Join(
		listBlockAttributeCreate(d, prefix+"d", &request.D, parseWidgetABCD),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[2]))
	require.Equal(t, `func parseWidgetABCD(d *tfschema.ResourceData, prefix string) (sdk.WidgetABCDRequest, error) {
	leaf := d.Get(prefix+"leaf").(string)
	request := sdk.NewWidgetABCDRequest(leaf)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, prefix+"note", request.WithNote),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[3]))
}

func Test_attachCreateMappings_fourLevelBlocksWithDeepSet(t *testing.T) {
	model := createOptionalsModel("Widget", []ssot.Attribute{fourLevelBlock(func(d *ssot.Attribute) { d.Type = ssot.AttrSet })})

	blocks := model.CreateBlockFuncs()
	require.Len(t, blocks, 3)
	require.Equal(t, `func parseWidgetABC(d *tfschema.ResourceData, prefix string) (sdk.WidgetABCRequest, error) {
	request := sdk.NewWidgetABCRequest()
	errs := errors.Join(
		setMappedValueCreateBuilder(d, prefix+"d", request.WithD, toWidgetABCD),
	)
	return *request, errs
}
`, renderTemplate(t, "generatedCreateRequest", blocks[2]))
}

func Test_attachCreateMappings_fourLevelBlocksWithUnmappedLeafStayHandwritten(t *testing.T) {
	model := createOptionalsModel("Widget", []ssot.Attribute{fourLevelBlock(func(d *ssot.Attribute) {
		d.Elem.Children[1].Type = ssot.AttrIdentifier
		d.Elem.Children[1].IdentifierType = "SchemaObjectIdentifier"
		d.Elem.Children[1].Required = true
	})})

	require.Empty(t, model.CreateBlockFuncs())
	require.Empty(t, model.CreateRequestBuild().Calls)
	for attr := model.Attributes[0]; ; attr = attr.Children[0] {
		require.Nil(t, attr.CreateOptionalCall, attr.Name)
		require.Nil(t, attr.CreateBlock, attr.Name)
		if len(attr.Children) == 0 || len(attr.Children[0].Children) == 0 {
			break
		}
	}
}

// fourLevelBlock is a -> b -> c -> d, alternating single and multi-item lists. adjustDeepest changes d.
func fourLevelBlock(adjustDeepest func(d *ssot.Attribute)) ssot.Attribute {
	block := ssot.Attribute{
		Name:     "a",
		Type:     ssot.AttrList,
		MaxItems: new(1),
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "b",
			Type: ssot.AttrList,
			Elem: ssot.ObjectElem([]ssot.Attribute{{
				Name:     "c",
				Type:     ssot.AttrList,
				MaxItems: new(1),
				Elem: ssot.ObjectElem([]ssot.Attribute{{
					Name: "d",
					Type: ssot.AttrList,
					Elem: ssot.ObjectElem([]ssot.Attribute{{
						Name:     "leaf",
						Type:     ssot.AttrString,
						Required: true,
					}, {
						Name: "note",
						Type: ssot.AttrString,
					}}),
				}}),
			}}),
		}}),
	}
	if adjustDeepest != nil {
		adjustDeepest(&block.Elem.Children[0].Elem.Children[0].Elem.Children[0])
	}
	return block
}

func Test_attachCreateMappings_skipCreateBlockStaysHandwritten(t *testing.T) {
	src := []ssot.Attribute{{
		Name:            "schedule",
		Type:            ssot.AttrList,
		MaxItems:        new(1),
		SkipCreateBlock: true,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:     "expression",
			Type:     ssot.AttrString,
			Required: true,
		}}),
	}, {
		Name: "comment",
		Type: ssot.AttrString,
	}}
	model := createOptionalsModel("Widget", src)

	require.Empty(t, model.CreateBlockFuncs())
	require.Nil(t, model.Attributes[0].CreateOptionalCall)
	require.Nil(t, model.Attributes[0].CreateBlock)
	require.Nil(t, model.Attributes[0].Children[0].CreateConstructorArg)
	require.Equal(t, "stringAttributeCreateBuilder", model.Attributes[1].CreateOptionalCall.Helper)
}

func Test_AlertCreateMappings_skipsNestedParsers(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: defs.Alert}, nil, nil)

	require.Empty(t, model.CreateBlockFuncs())
	require.Equal(t, `func alertApplyCreateOptionals(d *tfschema.ResourceData, request *sdk.CreateAlertRequest) error {
	return errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}
`, renderTemplate(t, "generatedApplyCreateOptionals", model))
}

func Test_ApplyCreateOptionals_rendersBlockParseFuncs(t *testing.T) {
	model := createOptionalsModel("Widget", []ssot.Attribute{{
		Name:     "policy",
		Type:     ssot.AttrList,
		MaxItems: new(1),
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "comment",
			Type: ssot.AttrString,
		}}),
	}})
	var b strings.Builder
	require.NoError(t, LifecycleTemplate.Execute(&b, model))
	rendered := b.String()
	require.Contains(t, rendered, "func widgetApplyCreateOptionals(")
	require.Contains(t, rendered, "func parseWidgetPolicy(")
	require.Less(t, strings.Index(rendered, "func widgetApplyCreateOptionals("), strings.Index(rendered, "func parseWidgetPolicy("))
}

func createOptionalsModel(name string, attrs []ssot.Attribute) ResourceModel {
	return ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:   name,
		IdType: "AccountObjectIdentifier",
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithApplyCreateOptionals(ssot.Generated()),
		},
		Attributes: attrs,
	}}, nil, nil)
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

func renderTemplate(t *testing.T, name string, model any) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, LifecycleTemplate.ExecuteTemplate(&b, name, model))
	return b.String()
}
