package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_attributeRules_namesAreUniqueAndDocumented(t *testing.T) {
	seen := make(map[string]struct{}, len(attributeRules))
	for _, rule := range attributeRules {
		require.NotEmpty(t, rule.Name)
		require.NotEmpty(t, rule.Description, rule.Name)
		_, dup := seen[rule.Name]
		require.False(t, dup, "duplicate rule name %q", rule.Name)
		seen[rule.Name] = struct{}{}
	}
}

func Test_applyAttributeRules_optionalBoolIsBooleanString(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "auto_resume",
		Type:        AttrBool,
		Description: "Specifies whether to automatically resume.",
	}, attributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, "BooleanDefault", got.DefaultExpr)
	require.Equal(t, "validateBooleanString", got.ValidateDiagExpr)
	require.Equal(t, `booleanStringFieldDescription("Specifies whether to automatically resume.")`, got.DescriptionExpr)
}

func Test_applyAttributeRules_optionalBoolRawDescription(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:           "initially_suspended",
		Type:           AttrBool,
		Description:    "Specifies whether the compute pool is created initially in the suspended state.",
		RawDescription: true,
	}, attributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.Equal(t, "BooleanDefault", got.DefaultExpr)
	require.Equal(t, `"Specifies whether the compute pool is created initially in the suspended state."`, got.DescriptionExpr)
}

func Test_applyAttributeRules_requiredBoolIsNative(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "enabled",
		Type:        AttrBool,
		Required:    true,
		Description: "Specifies whether it is enabled.",
	}, attributeRuleContext{})

	require.Equal(t, SchemaTypeBool, got.TfType)
	require.True(t, got.Required)
	require.Empty(t, got.DefaultExpr)
	require.Empty(t, got.ValidateDiagExpr)
	require.Equal(t, `"Specifies whether it is enabled."`, got.DescriptionExpr)
}

func Test_applyAttributeRules_optionalIntMinAllowsZeroUsesIntDefault(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "auto_suspend_secs",
		Type:        AttrInt,
		Description: "Number of seconds.",
		Min:         new(0),
	}, attributeRuleContext{})

	require.Equal(t, SchemaTypeInt, got.TfType)
	require.Equal(t, "IntDefault", got.DefaultExpr)
	require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(0))", got.ValidateDiagExpr)
}

func Test_applyAttributeRules_optionalIntWithoutZeroMinHasNoIntDefault(t *testing.T) {
	t.Run("no min", func(t *testing.T) {
		got := resolveAttribute(Attribute{Name: "n", Type: AttrInt}, attributeRuleContext{})
		require.Empty(t, got.DefaultExpr)
		require.Empty(t, got.ValidateDiagExpr)
	})
	t.Run("min greater than zero", func(t *testing.T) {
		got := resolveAttribute(Attribute{Name: "n", Type: AttrInt, Min: new(1)}, attributeRuleContext{})
		require.Empty(t, got.DefaultExpr)
		require.Equal(t, "validation.ToDiagFunc(validation.IntAtLeast(1))", got.ValidateDiagExpr)
	})
	t.Run("required with min zero", func(t *testing.T) {
		got := resolveAttribute(Attribute{Name: "n", Type: AttrInt, Required: true, Min: new(0)}, attributeRuleContext{})
		require.Empty(t, got.DefaultExpr)
	})
}

func Test_applyAttributeRules_identifierSuppressesQuoting(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "for_application",
		Type:        AttrIdentifier,
		ForceNew:    true,
		Description: "Specifies the Snowflake Native App name.",
	}, attributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.ForceNew)
	require.Equal(t, "suppressIdentifierQuoting", got.DiffSuppressExpr)
}

func Test_applyAttributeRules_enumField(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:                  "instance_family",
		Type:                  AttrEnum,
		Required:              true,
		ForceNew:              true,
		Description:           "Identifies the type of machine.",
		AdditionalDescription: "Not all instance families are supported.",
		EnumPlural:            "InstanceFamilies",
	}, attributeRuleContext{ResourceName: "ComputePool"})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Required)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", got.ValidateDiagExpr)
	require.Equal(t, "SuppressIfAny(NormalizeAndCompare(sdk.ToComputePoolInstanceFamily))", got.DiffSuppressExpr)
	require.Equal(t, `"Identifies the type of machine. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies) + " Not all instance families are supported."`, got.DescriptionExpr)
}

func Test_applyAttributeRules_listOfEnum(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "backup_instance_families",
		Type:        AttrList,
		Description: "Specifies an ordered list of instance families.",
		Elem: &Attribute{
			Type:       AttrEnum,
			Enum:       "InstanceFamily",
			EnumPlural: "InstanceFamilies",
		},
	}, attributeRuleContext{ResourceName: "ComputePool"})

	require.Equal(t, SchemaTypeList, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, "NormalizeAndCompare(sdk.ToComputePoolInstanceFamily)", got.DiffSuppressExpr)
	require.Equal(t, `"Specifies an ordered list of instance families. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies)`, got.DescriptionExpr)
	require.NotNil(t, got.Elem)
	require.Equal(t, SchemaTypeString, got.Elem.TfType)
	require.False(t, got.Elem.Required)
	require.False(t, got.Elem.Optional)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", got.Elem.ValidateDiagExpr)
	require.Empty(t, got.Elem.DiffSuppressExpr)
	require.Empty(t, got.Elem.DescriptionExpr)
}

func Test_applyAttributeRules_setOfEnum(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "allowed_values",
		Type:        AttrSet,
		Description: "Specifies a set of allowed values.",
		Elem: &Attribute{
			Type:       AttrEnum,
			Enum:       "InstanceFamily",
			EnumPlural: "InstanceFamilies",
		},
	}, attributeRuleContext{ResourceName: "ComputePool"})

	require.Equal(t, SchemaTypeSet, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, "NormalizeAndCompare(sdk.ToComputePoolInstanceFamily)", got.DiffSuppressExpr)
	require.Equal(t, `"Specifies a set of allowed values. " + enumValuesDescription(sdk.AllComputePoolInstanceFamilies)`, got.DescriptionExpr)
	require.NotNil(t, got.Elem)
	require.Equal(t, SchemaTypeString, got.Elem.TfType)
	require.False(t, got.Elem.Required)
	require.False(t, got.Elem.Optional)
	require.Equal(t, "sdkValidation(sdk.ToComputePoolInstanceFamily)", got.Elem.ValidateDiagExpr)
	require.Empty(t, got.Elem.DiffSuppressExpr)
	require.Empty(t, got.Elem.DescriptionExpr)
}

func Test_applyAttributeRules_setOfString(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "allowed_accounts",
		Type:        AttrSet,
		Description: "Specifies a set of accounts.",
		Elem:        &Attribute{Type: AttrString},
	}, attributeRuleContext{})

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
	require.True(t, isCollection(Attribute{Type: AttrList}))
	require.True(t, isCollection(Attribute{Type: AttrSet}))
	require.False(t, isCollection(Attribute{Type: AttrString}))
	require.False(t, isCollection(Attribute{Type: AttrEnum}))

	require.True(t, isCollectionOfEnum(Attribute{Type: AttrList, Elem: &Attribute{Type: AttrEnum}}))
	require.True(t, isCollectionOfEnum(Attribute{Type: AttrSet, Elem: &Attribute{Type: AttrEnum}}))
	require.False(t, isCollectionOfEnum(Attribute{Type: AttrList, Elem: &Attribute{Type: AttrString}}))
	require.False(t, isCollectionOfEnum(Attribute{Type: AttrSet}))
	require.False(t, isCollectionOfEnum(Attribute{Type: AttrEnum}))
}

func Test_sdkEnumSymbols(t *testing.T) {
	ctx := attributeRuleContext{ResourceName: "ComputePool"}
	attr := Attribute{Name: "instance_family", Type: AttrEnum, EnumPlural: "InstanceFamilies"}

	require.Equal(t, "InstanceFamily", enumName(attr))
	require.Equal(t, "InstanceFamilies", enumPlural(attr))
	require.Equal(t, "sdk.ToComputePoolInstanceFamily", ctx.sdkEnumTo(attr))
	require.Equal(t, "sdk.AllComputePoolInstanceFamilies", ctx.sdkEnumAll(attr))
}

func Test_applyAttributeRules_comment(t *testing.T) {
	got := resolveAttribute(Attribute{
		Name:        "comment",
		Type:        AttrString,
		Description: "Specifies a comment for the compute pool.",
	}, attributeRuleContext{})

	require.Equal(t, SchemaTypeString, got.TfType)
	require.True(t, got.Optional)
	require.Equal(t, `"Specifies a comment for the compute pool."`, got.DescriptionExpr)
	require.Empty(t, got.DefaultExpr)
	require.Empty(t, got.ValidateDiagExpr)
	require.Empty(t, got.DiffSuppressExpr)
}
