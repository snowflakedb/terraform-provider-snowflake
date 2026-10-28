package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/stretchr/testify/require"
)

func Test_resourceAttributeRules_namesAreUniqueAndDocumented(t *testing.T) {
	seen := make(map[string]struct{}, len(resourceAttributeRules))
	for _, rule := range resourceAttributeRules {
		require.NotEmpty(t, rule.Name)
		require.NotEmpty(t, rule.Description, rule.Name)
		_, dup := seen[rule.Name]
		require.False(t, dup, "duplicate rule name %q", rule.Name)
		seen[rule.Name] = struct{}{}
	}
}

func Test_applyResourceAttributeRules_optionalBoolIsBooleanString(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "auto_resume",
		Type:        ssot.AttrBool,
		Description: "Specifies whether to automatically resume.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, "BooleanDefault", got.DefaultExpr)
	require.Equal(t, "validateBooleanString", got.ValidateDiagExpr)
	require.Equal(t, `booleanStringFieldDescription("Specifies whether to automatically resume.")`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_optionalBoolRawDescription(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:           "initially_suspended",
		Type:           ssot.AttrBool,
		Description:    "Specifies whether the compute pool is created initially in the suspended state.",
		RawDescription: true,
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.Equal(t, "BooleanDefault", got.DefaultExpr)
	require.Equal(t, `"Specifies whether the compute pool is created initially in the suspended state."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_nativeBoolStaysTypeBool(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "initially_suspended",
		Type:        ssot.AttrBool,
		NativeBool:  true,
		Description: "Specifies whether the warehouse is created initially in the suspended state.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeBool, got.TfType)
	require.True(t, got.Optional)
	require.Empty(t, got.DefaultExpr)
	require.Empty(t, got.ValidateDiagExpr)
	require.Equal(t, `"Specifies whether the warehouse is created initially in the suspended state."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_sqlStatementNormalizes(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:         "condition",
		Type:         ssot.AttrString,
		Required:     true,
		Description:  "The SQL statement.",
		DiffSuppress: ssot.SuppressStatement,
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.Equal(t, "DiffSuppressStatement", got.DiffSuppressExpr)
	require.Equal(t, `"The SQL statement."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_requiredBoolIsNative(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "enabled",
		Type:        ssot.AttrBool,
		Required:    true,
		Description: "Specifies whether it is enabled.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeBool, got.TfType)
	require.True(t, got.Required)
	require.Empty(t, got.DefaultExpr)
	require.Empty(t, got.ValidateDiagExpr)
	require.Equal(t, `"Specifies whether it is enabled."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_optionalIntMinAllowsZeroUsesIntDefault(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "auto_suspend_secs",
		Type:        ssot.AttrInt,
		Description: "Number of seconds.",
		Min:         new(0),
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeInt, got.TfType)
	require.Equal(t, "IntDefault", got.DefaultExpr)
	require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(0))", got.ValidateDiagExpr)
}

func Test_applyResourceAttributeRules_optionalIntWithoutZeroMinHasNoIntDefault(t *testing.T) {
	t.Run("min zero", func(t *testing.T) {
		got := resolveResourceAttribute(ssot.Attribute{Name: "count", Type: ssot.AttrInt, Min: new(0)}, resourceAttributeRuleContext{})
		require.Equal(t, "IntDefault", got.DefaultExpr)
		require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(0))", got.ValidateDiagExpr)
	})
	t.Run("no min", func(t *testing.T) {
		got := resolveResourceAttribute(ssot.Attribute{Name: "n", Type: ssot.AttrInt}, resourceAttributeRuleContext{})
		require.Empty(t, got.DefaultExpr)
		require.Empty(t, got.ValidateDiagExpr)
	})
	t.Run("min greater than zero", func(t *testing.T) {
		got := resolveResourceAttribute(ssot.Attribute{Name: "n", Type: ssot.AttrInt, Min: new(1)}, resourceAttributeRuleContext{})
		require.Empty(t, got.DefaultExpr)
		require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(1))", got.ValidateDiagExpr)
	})
	t.Run("required with min zero", func(t *testing.T) {
		got := resolveResourceAttribute(ssot.Attribute{Name: "n", Type: ssot.AttrInt, Required: true, Min: new(0)}, resourceAttributeRuleContext{})
		require.Empty(t, got.DefaultExpr)
	})
}

func Test_applyResourceAttributeRules_identifierSuppressesQuoting(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "for_application",
		Type:        ssot.AttrIdentifier,
		ForceNew:    true,
		Description: "Specifies the Snowflake Native App name.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.ForceNew)
	require.Equal(t, "suppressIdentifierQuoting", got.DiffSuppressExpr)
}

func Test_applyResourceAttributeRules_enumField(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:                  "instance_family",
		Type:                  ssot.AttrEnum,
		Required:              true,
		ForceNew:              true,
		Description:           "Identifies the type of machine.",
		AdditionalDescription: "Not all instance families are supported.",
		Enum:                  defs.ComputePoolInstanceFamily,
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Required)
	require.Equal(t, "ComputePoolInstanceFamily", got.EnumName)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", got.ValidateDiagExpr)
	require.Equal(t, "SuppressIfAny(NormalizeAndCompare(sdk.ToComputePoolInstanceFamily))", got.DiffSuppressExpr)
	require.Equal(t, `"Identifies the type of machine. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies) + " Not all instance families are supported."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_listOfEnum(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "backup_instance_families",
		Type:        ssot.AttrList,
		Description: "Specifies an ordered list of instance families.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: defs.ComputePoolInstanceFamily,
		},
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeList, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, "NormalizeAndCompare(sdk.ToComputePoolInstanceFamily)", got.DiffSuppressExpr)
	require.Equal(t, `"Specifies an ordered list of instance families. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies)`, got.DescriptionExpr)
	require.NotNil(t, got.Elem)
	require.Equal(t, SchemaTypeString, got.Elem.TfType)
	require.Equal(t, "ComputePoolInstanceFamily", got.Elem.EnumName)
	require.False(t, got.Elem.Required)
	require.False(t, got.Elem.Optional)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", got.Elem.ValidateDiagExpr)
	require.Empty(t, got.Elem.DiffSuppressExpr)
	require.Empty(t, got.Elem.DescriptionExpr)
}

func Test_applyResourceAttributeRules_setOfEnum(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "allowed_values",
		Type:        ssot.AttrSet,
		Description: "Specifies a set of allowed values.",
		Elem: &ssot.Attribute{
			Type: ssot.AttrEnum,
			Enum: defs.ComputePoolInstanceFamily,
		},
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeSet, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, `NormalizeAndCompareEnumsInSet("allowed_values", sdk.ToComputePoolInstanceFamily)`, got.DiffSuppressExpr)
	require.Equal(t, `"Specifies a set of allowed values. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies)`, got.DescriptionExpr)
	require.NotNil(t, got.Elem)
	require.Equal(t, SchemaTypeString, got.Elem.TfType)
	require.Equal(t, "ComputePoolInstanceFamily", got.Elem.EnumName)
	require.False(t, got.Elem.Required)
	require.False(t, got.Elem.Optional)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", got.Elem.ValidateDiagExpr)
	require.Empty(t, got.Elem.DiffSuppressExpr)
	require.Empty(t, got.Elem.DescriptionExpr)
}

func Test_applyResourceAttributeRules_setOfString(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "allowed_accounts",
		Type:        ssot.AttrSet,
		Description: "Specifies a set of accounts.",
		Elem:        &ssot.Attribute{Type: ssot.AttrString},
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeSet, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, `"Specifies a set of accounts."`, got.DescriptionExpr)
	require.Empty(t, got.DiffSuppressExpr)
	require.NotNil(t, got.Elem)
	require.Equal(t, SchemaTypeString, got.Elem.TfType)
	require.False(t, got.Elem.Required)
	require.False(t, got.Elem.Optional)
	require.Empty(t, got.Elem.ValidateDiagExpr)
}

func Test_isCollection(t *testing.T) {
	require.True(t, isCollection(ssot.Attribute{Type: ssot.AttrList}))
	require.True(t, isCollection(ssot.Attribute{Type: ssot.AttrSet}))
	require.False(t, isCollection(ssot.Attribute{Type: ssot.AttrString}))
	require.False(t, isCollection(ssot.Attribute{Type: ssot.AttrEnum}))

	require.True(t, isCollectionOfEnum(ssot.Attribute{Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrEnum}}))
	require.True(t, isCollectionOfEnum(ssot.Attribute{Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrEnum}}))
	require.False(t, isCollectionOfEnum(ssot.Attribute{Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}}))
	require.False(t, isCollectionOfEnum(ssot.Attribute{Type: ssot.AttrSet}))
	require.False(t, isCollectionOfEnum(ssot.Attribute{Type: ssot.AttrEnum}))
}

func Test_sdkEnumSymbols(t *testing.T) {
	ctx := resourceAttributeRuleContext{}
	attr := ssot.Attribute{Name: "instance_family", Type: ssot.AttrEnum, Enum: defs.ComputePoolInstanceFamily}

	require.Equal(t, "ComputePoolInstanceFamily", enumDef(attr).Name)
	require.Equal(t, "ComputePoolInstanceFamilies", enumDef(attr).NamePlural)
	require.Equal(t, "sdk.ToComputePoolInstanceFamily", ctx.sdkEnumTo(attr))
	require.Equal(t, "sdk.AllComputePoolInstanceFamilies", ctx.sdkEnumAll(attr))
}

func Test_applyResourceAttributeRules_createOnlyBoolIgnoresChanges(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:           "initially_suspended",
		Type:           ssot.AttrBool,
		CreateOnly:     true,
		RawDescription: true,
		Description:    "Specifies whether the compute pool is created initially in the suspended state.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, "BooleanDefault", got.DefaultExpr)
	require.Equal(t, "validateBooleanString", got.ValidateDiagExpr)
	require.Equal(t, "IgnoreAfterCreation", got.DiffSuppressExpr)
	require.Equal(t, `"Specifies whether the compute pool is created initially in the suspended state."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_createOnlyNativeBoolIsNative(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "initially_suspended",
		Type:        ssot.AttrBool,
		CreateOnly:  true,
		NativeBool:  true,
		Description: "Specifies whether the warehouse is created initially in the ‘Suspended’ state.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeBool, got.TfType)
	require.True(t, got.Optional)
	require.Empty(t, got.DefaultExpr)
	require.Empty(t, got.ValidateDiagExpr)
	require.Equal(t, "IgnoreAfterCreation", got.DiffSuppressExpr)
	require.Equal(t, `"Specifies whether the warehouse is created initially in the ‘Suspended’ state."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_intBetweenKeepsIntDefault(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:         "query_acceleration_max_scale_factor",
		Type:         ssot.AttrInt,
		Description:  "Specifies the maximum scale factor.",
		Min:          new(0),
		Max:          new(100),
		DiffSuppress: ssot.SuppressShow,
	}, resourceAttributeRuleContext{})

	require.Equal(t, "validation.ToDiagFunc(validation.IntBetween(0, 100))", got.ValidateDiagExpr)
	require.Equal(t, "IntDefault", got.DefaultExpr)
	require.Equal(t, `IgnoreChangeToCurrentSnowflakeValueInShow("query_acceleration_max_scale_factor")`, got.DiffSuppressExpr)
}

func Test_applyResourceAttributeRules_enumNamePluralAndShowSuppress(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:                  "warehouse_type",
		Type:                  ssot.AttrEnum,
		Description:           "Specifies warehouse type.",
		AdditionalDescription: "Warehouse needs to be suspended.",
		Enum:                  ssot.NewEnum("WarehouseTypeUserSettable", "WarehouseTypesUserSettable", "STANDARD", "SNOWPARK-OPTIMIZED"),
		ShowName:              "type",
		DiffSuppress:          ssot.SuppressShow,
	}, resourceAttributeRuleContext{})

	require.Equal(t, "sdkValidation(sdk.ToWarehouseTypeUserSettable)", got.ValidateDiagExpr)
	require.Equal(t, `SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseTypeUserSettable), IgnoreChangeToCurrentSnowflakeValueInShow("type"))`, got.DiffSuppressExpr)
	require.Equal(t, `"Specifies warehouse type. " + enumValuesDescription(sdk.AllWarehouseTypesUserSettable) + " Warehouse needs to be suspended."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_enumAliasesUseCanonicalList(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "warehouse_size",
		Type:        ssot.AttrEnum,
		Description: "Specifies the size of the virtual warehouse.",
		Enum: ssot.NewEnum("WarehouseSize", "WarehouseSizes", "XSMALL").
			WithAliases("XSMALL", "X-SMALL"),
	}, resourceAttributeRuleContext{})

	require.Equal(t, `"Specifies the size of the virtual warehouse. " + enumValuesDescription(sdk.AllWarehouseSizes)`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_includeAliasedInTheDescription(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:                           "warehouse_size",
		Type:                           ssot.AttrEnum,
		Description:                    "Specifies the size of the virtual warehouse.",
		IncludeAliasedInTheDescription: true,
		Enum: ssot.NewEnum("WarehouseSize", "WarehouseSizes", "XSMALL").
			WithAliases("XSMALL", "X-SMALL"),
	}, resourceAttributeRuleContext{})

	require.Equal(t, `"Specifies the size of the virtual warehouse. " + enumValuesDescription(sdk.AllWarehouseSizesString)`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_identifierValidationRelatedResourceAndShowSuppress(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:            "resource_monitor",
		Type:            ssot.AttrIdentifier,
		Description:     "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
		RelatedResource: &ssot.Object{Name: "ResourceMonitor"},
		DiffSuppress:    ssot.SuppressShow,
		ConflictsWith:   []ssot.Attribute{{Name: "generation"}},
	}, resourceAttributeRuleContext{IdType: "AccountObjectIdentifier"})

	require.Equal(t, "IsValidIdentifier[sdk.AccountObjectIdentifier]()", got.ValidateDiagExpr)
	require.Equal(t, `SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor"))`, got.DiffSuppressExpr)
	require.Equal(t, `relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the warehouse.", resources.ResourceMonitor)`, got.DescriptionExpr)
	require.Equal(t, []string{"generation"}, got.ConflictsWith)
}

func Test_applyResourceAttributeRules_showSuppressOnPlainField(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:         "max_cluster_count",
		Type:         ssot.AttrInt,
		Min:          new(1),
		DiffSuppress: ssot.SuppressShow,
	}, resourceAttributeRuleContext{})

	require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(1))", got.ValidateDiagExpr)
	require.Empty(t, got.DefaultExpr)
	require.Equal(t, `IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count")`, got.DiffSuppressExpr)
}

func Test_applyResourceAttributeRules_diffSuppressComposition(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name: "condition", Type: ssot.AttrString,
		DiffSuppress: ssot.SuppressIfAny(ssot.SuppressStatement, ssot.SuppressDescribe.WithKey("definition")),
	}, resourceAttributeRuleContext{})
	require.Equal(t, `SuppressIfAny(DiffSuppressStatement, IgnoreChangeToCurrentSnowflakeValueInDescribe("definition"))`, got.DiffSuppressExpr)
	got = resolveResourceAttribute(ssot.Attribute{
		Name: "mode", Type: ssot.AttrEnum, Enum: ssot.NewEnum("Mode", "Modes", "A"), ShowName: "output_mode",
		DiffSuppress: ssot.SuppressIfAny(ssot.SuppressShow, ssot.SuppressDescribe),
	}, resourceAttributeRuleContext{})
	require.Equal(t, `SuppressIfAny(NormalizeAndCompare(sdk.ToMode), IgnoreChangeToCurrentSnowflakeValueInShow("output_mode"), IgnoreChangeToCurrentSnowflakeValueInDescribe("output_mode"))`, got.DiffSuppressExpr)
}

func Test_applyResourceAttributeRules_composableSuppressors(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		suppressor ssot.DiffSuppress
		expression string
	}{
		{name: "none"},
		{name: "empty composition", suppressor: ssot.SuppressIfAny()},
		{name: "describe", suppressor: ssot.SuppressDescribe, expression: `IgnoreChangeToCurrentSnowflakeValueInDescribe("output_mode")`},
		{name: "key override", suppressor: ssot.SuppressShow.WithKey("other"), expression: `IgnoreChangeToCurrentSnowflakeValueInShow("other")`},
		{
			name:       "nested composition",
			suppressor: ssot.SuppressIfAny(ssot.SuppressStatement, ssot.SuppressIfAny(ssot.SuppressShow, nil, ssot.SuppressDescribe.WithKey("definition"))),
			expression: `SuppressIfAny(DiffSuppressStatement, IgnoreChangeToCurrentSnowflakeValueInShow("output_mode"), IgnoreChangeToCurrentSnowflakeValueInDescribe("definition"))`,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			got := resolveResourceAttribute(ssot.Attribute{
				Name: "mode", Type: ssot.AttrString, ShowName: "output_mode", DiffSuppress: testCase.suppressor,
			}, resourceAttributeRuleContext{})
			require.Equal(t, testCase.expression, got.DiffSuppressExpr)
		})
	}
}

func Test_applyResourceAttributeRules_comment(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the compute pool.",
	}, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, `"Specifies a comment for the compute pool."`, got.DescriptionExpr)
	require.Empty(t, got.DefaultExpr)
	require.Empty(t, got.ValidateDiagExpr)
	require.Empty(t, got.DiffSuppressExpr)
}

func Test_applyResourceAttributeRules_nestedBlocksRecurse(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "outer",
		Type:        ssot.AttrSet,
		Description: "Outer block.",
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:     "inner",
			Type:     ssot.AttrList,
			Required: true,
			Elem: ssot.ObjectElem([]ssot.Attribute{{
				Name:        "leaf",
				Type:        ssot.AttrEnum,
				Required:    true,
				Description: "A leaf.",
				Enum:        defs.ComputePoolInstanceFamily,
			}, {
				Name: "target",
				Type: ssot.AttrIdentifier,
			}}),
		}}),
	}, resourceAttributeRuleContext{IdType: "SchemaObjectIdentifier"})

	require.Equal(t, SchemaTypeSet, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, `"Outer block."`, got.DescriptionExpr)
	require.Nil(t, got.Elem)
	require.Len(t, got.Children, 1)

	inner := got.Children[0]
	require.Equal(t, "inner", inner.Name)
	require.Equal(t, SchemaTypeList, inner.TfType)
	require.True(t, inner.Required)
	require.False(t, inner.Optional)
	require.Len(t, inner.Children, 2)

	leaf := inner.Children[0]
	require.Equal(t, "leaf", leaf.Name)
	require.Equal(t, SchemaTypeString, leaf.TfType)
	require.True(t, leaf.Required)
	require.Equal(t, "ComputePoolInstanceFamily", leaf.EnumName)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", leaf.ValidateDiagExpr)
	require.Equal(t, "SuppressIfAny(NormalizeAndCompare(sdk.ToComputePoolInstanceFamily))", leaf.DiffSuppressExpr)
	require.Equal(t, `"A leaf. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies)`, leaf.DescriptionExpr)

	target := inner.Children[1]
	require.Equal(t, "suppressIdentifierQuoting", target.DiffSuppressExpr)
	require.Empty(t, target.ValidateDiagExpr)

	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{gotSourceNested()},
	}}, nil, nil)
	var b strings.Builder
	require.NoError(t, SchemaTemplate.ExecuteTemplate(&b, "attributesSchemaDefinition", model))
	rendered := b.String()
	require.Equal(t, 2, strings.Count(rendered, "&tfschema.Resource{"))
	require.Contains(t, rendered, `"outer"`)
	require.Contains(t, rendered, `"inner"`)
	require.Contains(t, rendered, `"leaf"`)
}

func gotSourceNested() ssot.Attribute {
	return ssot.Attribute{
		Name: "outer",
		Type: ssot.AttrSet,
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name: "inner",
			Type: ssot.AttrList,
			Elem: ssot.ObjectElem([]ssot.Attribute{{
				Name: "leaf",
				Type: ssot.AttrString,
			}}),
		}}),
	}
}

func Test_applyResourceAttributeRules_nestedBlock(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "client_policy",
		Type:        ssot.AttrSet,
		Description: "Allows to set policies per-client type.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:        "client_type",
				Type:        ssot.AttrEnum,
				Required:    true,
				Description: "The client or driver type.",
				Enum:        defs.ClientPolicyDriverType,
			},
			{
				Name:        "target",
				Type:        ssot.AttrIdentifier,
				Description: "An identifier inside the block.",
			},
			{
				Name:        "minimum_version",
				Type:        ssot.AttrString,
				Required:    true,
				Description: "Minimum allowed version.",
			},
		}),
	}, resourceAttributeRuleContext{IdType: "SchemaObjectIdentifier"})

	require.Equal(t, SchemaTypeSet, got.TfType)
	require.True(t, got.Optional)
	require.Nil(t, got.Elem)
	require.Len(t, got.Children, 3)

	clientType := got.Children[0]
	require.Equal(t, "client_type", clientType.Name)
	require.True(t, clientType.Required)
	require.False(t, clientType.Optional)
	require.Equal(t, "ClientPolicyDriverType", clientType.EnumName)
	require.Equal(t, "sdkValidation(sdk.ToClientPolicyDriverType)", clientType.ValidateDiagExpr)
	require.Equal(t, "SuppressIfAny(NormalizeAndCompare(sdk.ToClientPolicyDriverType))", clientType.DiffSuppressExpr)
	require.Equal(t, `"The client or driver type. " + enumValuesDescription(sdk.AllClientPolicyDriverTypes)`, clientType.DescriptionExpr)

	target := got.Children[1]
	require.Equal(t, "suppressIdentifierQuoting", target.DiffSuppressExpr)
	require.Empty(t, target.ValidateDiagExpr)

	version := got.Children[2]
	require.True(t, version.Required)
	require.Equal(t, `"Minimum allowed version."`, version.DescriptionExpr)
}

func Test_applyResourceAttributeRules_listOfBlocksAllowsMultipleItems(t *testing.T) {
	attr := ssot.Attribute{
		Name:        "storage_location",
		Type:        ssot.AttrList,
		Required:    true,
		MinItems:    new(1),
		Description: "List of named cloud storage locations.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{Name: "storage_location_name", Type: ssot.AttrString, Required: true, Description: "Name of the storage location."},
			{Name: "storage_base_url", Type: ssot.AttrString, Required: true, Description: "Specifies the base URL."},
			{Name: "storage_aws_role_arn", Type: ssot.AttrString, Description: "Specifies the role ARN."},
			{Name: "allowed_locations", Type: ssot.AttrList, MinItems: new(1), Description: "Extra locations.", Elem: &ssot.Attribute{Type: ssot.AttrString}},
		}),
	}
	got := resolveResourceAttribute(attr, resourceAttributeRuleContext{})

	require.Equal(t, SchemaTypeList, got.TfType)
	require.True(t, got.Required)
	require.False(t, got.Optional)
	require.NotNil(t, got.MinItems)
	require.Equal(t, 1, *got.MinItems)
	require.Nil(t, got.MaxItems)
	require.Len(t, got.Children, 4)
	require.NotNil(t, got.Children[3].MinItems)
	require.Equal(t, 1, *got.Children[3].MinItems)
	require.Equal(t, SchemaTypeList, got.Children[3].TfType)
	require.Equal(t, SchemaTypeString, got.Children[3].Elem.TfType)

	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name:       "ExternalVolume",
		Attributes: []ssot.Attribute{attr},
	}}, nil, nil)
	var b strings.Builder
	require.NoError(t, SchemaTemplate.ExecuteTemplate(&b, "attributesSchemaDefinition", model))
	rendered := b.String()
	require.Contains(t, rendered, "MinItems: 1,")
	require.Equal(t, 2, strings.Count(rendered, "MinItems: 1,"))
	require.NotContains(t, rendered, "MaxItems:")
	require.Contains(t, rendered, `"storage_location_name"`)
	require.Contains(t, rendered, `"storage_base_url"`)
	require.Contains(t, rendered, `"storage_aws_role_arn"`)
	require.Contains(t, rendered, `"allowed_locations"`)
	require.Equal(t, 1, strings.Count(rendered, "&tfschema.Resource{"))
	require.Contains(t, rendered, "&tfschema.Schema{")
}

func Test_applyResourceAttributeRules_setOfIdentifiers(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "tables",
		Type:        ssot.AttrSet,
		Description: "Specifies the tables.",
		Elem: &ssot.Attribute{
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
		},
	}, resourceAttributeRuleContext{IdType: "AccountObjectIdentifier"})

	require.Equal(t, SchemaTypeSet, got.TfType)
	require.Equal(t, `NormalizeAndCompareIdentifiersInSet("tables")`, got.DiffSuppressExpr)
	require.NotNil(t, got.Elem)
	require.Equal(t, "IsValidIdentifier[sdk.SchemaObjectIdentifier]()", got.Elem.ValidateDiagExpr)
	require.Empty(t, got.Elem.DescriptionExpr)
}

func Test_applyResourceAttributeRules_nestedSetOfEnumsStaysElementWise(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:     "mfa_policy",
		Type:     ssot.AttrList,
		MaxItems: new(1),
		Elem: ssot.ObjectElem([]ssot.Attribute{{
			Name:        "allowed_methods",
			Type:        ssot.AttrSet,
			Description: "Specifies the allowed methods.",
			Elem: &ssot.Attribute{
				Type: ssot.AttrEnum,
				Enum: defs.MfaPolicyAllowedMethodsOption,
			},
		}}),
	}, resourceAttributeRuleContext{})

	require.Equal(t, "NormalizeAndCompare(sdk.ToMfaPolicyAllowedMethodsOption)", got.Children[0].DiffSuppressExpr)
}

func Test_applyResourceAttributeRules_nestedRelationsAndSensitive(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:     "alert_schedule",
		Type:     ssot.AttrList,
		MaxItems: new(1),
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:          "cron",
				Type:          ssot.AttrList,
				MaxItems:      new(1),
				ConflictsWith: []ssot.Attribute{{Name: "interval"}},
			},
			{
				Name:          "interval",
				Type:          ssot.AttrInt,
				ConflictsWith: []ssot.Attribute{{Name: "cron"}},
				AtLeastOneOf:  []ssot.Attribute{{Name: "interval"}, {Name: "cron"}},
			},
		}),
	}, resourceAttributeRuleContext{})

	require.Equal(t, []string{"alert_schedule.0.interval"}, got.Children[0].ConflictsWith)
	require.Equal(t, []string{"alert_schedule.0.cron"}, got.Children[1].ConflictsWith)
	require.Equal(t, []string{"alert_schedule.0.interval", "alert_schedule.0.cron"}, got.Children[1].AtLeastOneOf)

	secret := resolveResourceAttribute(ssot.Attribute{
		Name:      "storage_aws_secret_key",
		Type:      ssot.AttrString,
		Sensitive: true,
	}, resourceAttributeRuleContext{})
	require.True(t, secret.Sensitive)

	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name: "Widget",
		Attributes: []ssot.Attribute{
			{
				Name:     "alert_schedule",
				Type:     ssot.AttrList,
				MaxItems: new(1),
				Elem: ssot.ObjectElem([]ssot.Attribute{
					{Name: "cron", Type: ssot.AttrString, ConflictsWith: []ssot.Attribute{{Name: "interval"}}},
					{Name: "interval", Type: ssot.AttrInt, AtLeastOneOf: []ssot.Attribute{{Name: "cron"}, {Name: "interval"}}},
				}),
			},
			{Name: "storage_aws_secret_key", Type: ssot.AttrString, Sensitive: true},
		},
	}}, nil, nil)
	var b strings.Builder
	require.NoError(t, SchemaTemplate.ExecuteTemplate(&b, "attributesSchemaDefinition", model))
	rendered := b.String()
	require.Contains(t, rendered, `"alert_schedule.0.interval"`)
	require.Contains(t, rendered, "AtLeastOneOf: []string{")
	require.Contains(t, rendered, "Sensitive: true,")
}

func Test_ModelFromInputObject_skipsSqlOnly(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: ssot.Object{
		Name: "WarehouseAdaptive",
		Attributes: []ssot.Attribute{
			{Name: "comment", Type: ssot.AttrString},
			{
				Name:      "warehouse_type_static",
				SqlOnly:   true,
				CreateSql: ssot.StaticSQL("warehouseType", "WAREHOUSE_TYPE = 'ADAPTIVE'"),
			},
		},
	}}, nil, nil)
	require.Len(t, model.Attributes, 1)
	require.Equal(t, "comment", model.Attributes[0].Name)
}

func Test_applyResourceAttributeRules_singleObjectIsListOfOne(t *testing.T) {
	leaf := ssot.Attribute{Name: "leaf", Type: ssot.AttrString, Required: true, ForceNew: true}
	inner := ssot.Attribute{Name: "inner", Type: ssot.AttrObject, Children: []ssot.Attribute{leaf}}

	t.Run("optional", func(t *testing.T) {
		got := resolveResourceAttribute(ssot.Attribute{Name: "block", Type: ssot.AttrObject, ForceNew: true, Children: []ssot.Attribute{leaf, inner}}, resourceAttributeRuleContext{})

		require.Equal(t, SchemaTypeList, got.TfType)
		require.True(t, got.Optional)
		require.True(t, got.ForceNew)
		require.Nil(t, got.MinItems)
		require.Equal(t, new(1), got.MaxItems)
		require.Len(t, got.Children, 2)
		require.Equal(t, "leaf", got.Children[0].Name)
		require.Equal(t, SchemaTypeList, got.Children[1].TfType, "a nested object is a list of one too")
		require.Equal(t, new(1), got.Children[1].MaxItems)
	})

	t.Run("required", func(t *testing.T) {
		got := resolveResourceAttribute(ssot.Attribute{Name: "block", Type: ssot.AttrObject, Required: true, Children: []ssot.Attribute{leaf}}, resourceAttributeRuleContext{})

		require.Equal(t, SchemaTypeList, got.TfType)
		require.True(t, got.Required)
		require.Equal(t, new(1), got.MinItems)
		require.Equal(t, new(1), got.MaxItems)
	})

	t.Run("equals the list declaration", func(t *testing.T) {
		object := resolveResourceAttribute(ssot.Attribute{Name: "block", Type: ssot.AttrObject, Required: true, ForceNew: true, Children: []ssot.Attribute{leaf}}, resourceAttributeRuleContext{})
		list := resolveResourceAttribute(ssot.Attribute{
			Name: "block", Type: ssot.AttrList, Required: true, ForceNew: true, MinItems: new(1), MaxItems: new(1),
			Elem: ssot.ObjectElem([]ssot.Attribute{leaf}),
		}, resourceAttributeRuleContext{})
		require.Equal(t, list, object)
	})
}

func Test_applyExactlyOneOf(t *testing.T) {
	one := ssot.Attribute{Name: "one"}
	two := ssot.Attribute{Name: "two"}
	three := ssot.Attribute{Name: "three"}
	attributes := []ResourceAttribute{{Name: "one"}, {Name: "two"}, {Name: "three"}, {Name: "other"}}

	applyExactlyOneOf([][]*ssot.Attribute{{&one, &three}}, attributes)

	require.Equal(t, []string{"one", "three"}, attributes[0].ExactlyOneOf)
	require.Nil(t, attributes[1].ExactlyOneOf)
	require.Equal(t, []string{"one", "three"}, attributes[2].ExactlyOneOf)
	require.Nil(t, attributes[3].ExactlyOneOf)

	require.Panics(t, func() { applyExactlyOneOf([][]*ssot.Attribute{{&one, &two, {Name: "missing"}}}, attributes) })
}
