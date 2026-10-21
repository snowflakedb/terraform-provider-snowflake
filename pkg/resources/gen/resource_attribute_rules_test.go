package gen

import (
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
	require.Equal(t, "NormalizeAndCompare(sdk.ToComputePoolInstanceFamily)", got.DiffSuppressExpr)
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
		Name:               "query_acceleration_max_scale_factor",
		Type:               ssot.AttrInt,
		Description:        "Specifies the maximum scale factor.",
		Min:                new(0),
		Max:                new(100),
		SuppressShowChange: true,
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
		SuppressShowChange:    true,
	}, resourceAttributeRuleContext{})

	require.Equal(t, "sdkValidation(sdk.ToWarehouseTypeUserSettable)", got.ValidateDiagExpr)
	require.Equal(t, `SuppressIfAny(NormalizeAndCompare(sdk.ToWarehouseTypeUserSettable), IgnoreChangeToCurrentSnowflakeValueInShow("type"))`, got.DiffSuppressExpr)
	require.Equal(t, `"Specifies warehouse type. " + enumValuesDescription(sdk.AllWarehouseTypesUserSettable) + " Warehouse needs to be suspended."`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_enumAliasesUseStringList(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:        "warehouse_size",
		Type:        ssot.AttrEnum,
		Description: "Specifies the size of the virtual warehouse.",
		Enum: ssot.NewEnum("WarehouseSize", "WarehouseSizes", "XSMALL").
			WithAliases("XSMALL", "X-SMALL"),
	}, resourceAttributeRuleContext{})

	require.Equal(t, `"Specifies the size of the virtual warehouse. " + enumValuesDescription(sdk.AllWarehouseSizesString)`, got.DescriptionExpr)
}

func Test_applyResourceAttributeRules_identifierValidationRelatedResourceAndShowSuppress(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:               "resource_monitor",
		Type:               ssot.AttrIdentifier,
		Description:        "Specifies the name of a resource monitor that is explicitly assigned to the warehouse.",
		RelatedResource:    &ssot.Object{Name: "ResourceMonitor"},
		SuppressShowChange: true,
		ConflictsWith:      []ssot.Attribute{{Name: "generation"}},
	}, resourceAttributeRuleContext{IdType: "AccountObjectIdentifier"})

	require.Equal(t, "IsValidIdentifier[sdk.AccountObjectIdentifier]()", got.ValidateDiagExpr)
	require.Equal(t, `SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("resource_monitor"))`, got.DiffSuppressExpr)
	require.Equal(t, `relatedResourceDescription("Specifies the name of a resource monitor that is explicitly assigned to the warehouse.", resources.ResourceMonitor)`, got.DescriptionExpr)
	require.Equal(t, []string{"generation"}, got.ConflictsWith)
}

func Test_applyResourceAttributeRules_showSuppressOnPlainField(t *testing.T) {
	got := resolveResourceAttribute(ssot.Attribute{
		Name:               "max_cluster_count",
		Type:               ssot.AttrInt,
		Min:                new(1),
		SuppressShowChange: true,
	}, resourceAttributeRuleContext{})

	require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(1))", got.ValidateDiagExpr)
	require.Empty(t, got.DefaultExpr)
	require.Equal(t, `IgnoreChangeToCurrentSnowflakeValueInShow("max_cluster_count")`, got.DiffSuppressExpr)
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
