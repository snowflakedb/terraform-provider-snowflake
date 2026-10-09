package gen

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

const (
	RuleOptionalBoolIsBooleanString            = "optional_bool_is_boolean_string"
	RuleRequiredBoolIsNative                   = "required_bool_is_native"
	RuleOptionalIntMinAllowsZeroUsesIntDefault = "optional_int_min_allows_zero_uses_int_default"
	RuleIdentifierSuppressesQuoting            = "identifier_suppresses_quoting"
	RuleEnumValidates                          = "enum_validates"
	RuleEnumDescriptionListsValues             = "enum_description_lists_values"
	RuleEnumFieldNormalizes                    = "enum_field_normalizes"
	RuleCollectionOfEnumNormalizes             = "collection_of_enum_normalizes"
)

// attributeRuleContext distinguishes a top-level attribute from a collection
// elem (list or set) so rules that only apply to fields (enum DiffSuppress,
// descriptions) can skip elems.
type attributeRuleContext struct {
	AsCollectionElem bool
	ResourceName     string
}

type attributeRule struct {
	Name        string
	Description string
	When        func(Attribute, attributeRuleContext) bool
	Then        func(Attribute, attributeRuleContext, *ResolvedAttribute)
}

// attributeRules is the ordered policy applied by applyAttributeRules.
// When two rules could touch the same field they are disjoint in When.
var attributeRules = []attributeRule{
	{
		// That is not always true, there are at least two groups of optional bool attributes that don't use this pattern:
		// - Create-time-only flags, e.g. snowflake_warehouse.initially_suspended
		// - Discriminator flags inside blocks (snowflake_session_policy.allowed_secondary_roles.0.all).
		// TODO [next PRs]: adjust the logic if needed with next objects migrations
		Name:        RuleOptionalBoolIsBooleanString,
		Description: "Optional bools use boolean-string chrome because TF TypeBool cannot distinguish unset from false.",
		When: func(in Attribute, _ attributeRuleContext) bool {
			return in.Type == AttrBool && !in.Required
		},
		Then: func(in Attribute, _ attributeRuleContext, out *ResolvedAttribute) {
			out.TfType = SchemaTypeString
			out.DefaultExpr = "BooleanDefault"
			out.ValidateDiagExpr = "validateBooleanString"
			if !in.RawDescription && in.Description != "" {
				out.DescriptionExpr = fmt.Sprintf("booleanStringFieldDescription(%s)", strconv.Quote(in.Description))
			}
		},
	},
	{
		Name:        RuleRequiredBoolIsNative,
		Description: "Required bools use native TF TypeBool; the user must set them, so no sentinel default is needed.",
		When: func(in Attribute, _ attributeRuleContext) bool {
			return in.Type == AttrBool && in.Required
		},
		Then: func(_ Attribute, _ attributeRuleContext, out *ResolvedAttribute) {
			out.TfType = SchemaTypeBool
		},
	},
	{
		Name:        RuleOptionalIntMinAllowsZeroUsesIntDefault,
		Description: "Optional ints use IntDefault only when Min shows 0 is an acceptable value (otherwise 0 already means unset).",
		When: func(in Attribute, _ attributeRuleContext) bool {
			return in.Type == AttrInt && !in.Required && in.Min != nil && *in.Min <= 0
		},
		Then: func(_ Attribute, _ attributeRuleContext, out *ResolvedAttribute) {
			out.DefaultExpr = "IntDefault"
		},
	},
	{
		Name:        RuleIdentifierSuppressesQuoting,
		Description: "Identifier attributes suppress quoting differences in diffs.",
		When: func(in Attribute, _ attributeRuleContext) bool {
			return in.Type == AttrIdentifier
		},
		Then: func(_ Attribute, _ attributeRuleContext, out *ResolvedAttribute) {
			out.DiffSuppressExpr = "suppressIdentifierQuoting"
		},
	},
	{
		Name:        RuleEnumValidates,
		Description: "Enum attributes (including collection elems) validate through the SDK converter.",
		When: func(in Attribute, ctx attributeRuleContext) bool {
			return in.Type == AttrEnum && ctx.sdkEnumTo(in) != ""
		},
		Then: func(in Attribute, ctx attributeRuleContext, out *ResolvedAttribute) {
			out.ValidateDiagExpr = fmt.Sprintf("sdkValidation(%s)", ctx.sdkEnumTo(in))
		},
	},
	{
		Name:        RuleEnumDescriptionListsValues,
		Description: "Enum fields and collections of enums append enumValuesDescription to the description.",
		When: func(in Attribute, ctx attributeRuleContext) bool {
			if ctx.AsCollectionElem || in.Description == "" || ctx.sdkEnumAll(in) == "" {
				return false
			}
			return in.Type == AttrEnum || isCollectionOfEnum(in)
		},
		Then: func(in Attribute, ctx attributeRuleContext, out *ResolvedAttribute) {
			expr := fmt.Sprintf("%s + enumValuesDescription(%s)", strconv.Quote(in.Description+" "), ctx.sdkEnumAll(in))
			if in.AdditionalDescription != "" {
				expr += " + " + strconv.Quote(" "+in.AdditionalDescription)
			}
			out.DescriptionExpr = expr
		},
	},
	{
		Name:        RuleEnumFieldNormalizes,
		Description: "Top-level enum fields normalize values in diffs via SuppressIfAny(NormalizeAndCompare).",
		When: func(in Attribute, ctx attributeRuleContext) bool {
			return in.Type == AttrEnum && ctx.sdkEnumTo(in) != "" && !ctx.AsCollectionElem
		},
		Then: func(in Attribute, ctx attributeRuleContext, out *ResolvedAttribute) {
			out.DiffSuppressExpr = fmt.Sprintf("SuppressIfAny(NormalizeAndCompare(%s))", ctx.sdkEnumTo(in))
		},
	},
	{
		Name:        RuleCollectionOfEnumNormalizes,
		Description: "Collections of enums (list or set) normalize values in diffs via NormalizeAndCompare on the collection field.",
		When: func(in Attribute, ctx attributeRuleContext) bool {
			return isCollectionOfEnum(in) && ctx.sdkEnumTo(in) != "" && !ctx.AsCollectionElem
		},
		Then: func(in Attribute, ctx attributeRuleContext, out *ResolvedAttribute) {
			out.DiffSuppressExpr = fmt.Sprintf("NormalizeAndCompare(%s)", ctx.sdkEnumTo(in))
		},
	},
}

func isCollection(in Attribute) bool {
	return in.Type == AttrList || in.Type == AttrSet
}

func isCollectionOfEnum(in Attribute) bool {
	return isCollection(in) && in.Elem != nil && in.Elem.Type == AttrEnum
}

func enumName(in Attribute) string {
	if in.Enum != "" {
		return in.Enum
	}
	if in.Elem != nil && in.Elem.Enum != "" {
		return in.Elem.Enum
	}
	if in.Type == AttrEnum && in.Name != "" {
		return genhelpers.SnakeCaseToCamel(in.Name)
	}
	return ""
}

func enumPlural(in Attribute) string {
	if in.EnumPlural != "" {
		return in.EnumPlural
	}
	if in.Elem != nil && in.Elem.EnumPlural != "" {
		return in.Elem.EnumPlural
	}
	return ""
}

// TODO [next PRs]: move enum abstraction from SDK.
func (ctx attributeRuleContext) sdkEnumTo(in Attribute) string {
	name := enumName(in)
	if name == "" || ctx.ResourceName == "" {
		return ""
	}
	// sdk.To<Resource><Enum>
	return "sdk.To" + ctx.ResourceName + name
}

// TODO [next PRs]: move enum abstraction from SDK.
func (ctx attributeRuleContext) sdkEnumAll(in Attribute) string {
	plural := enumPlural(in)
	if plural == "" || ctx.ResourceName == "" {
		return ""
	}
	// sdk.All<Resource><EnumPlural>
	return "sdk.All" + ctx.ResourceName + plural
}

// applyAttributeRules is the dedicated processing step from written Attribute intent to ResolvedAttribute schema chrome.
func applyAttributeRules(resourceName string, in []Attribute) []ResolvedAttribute {
	ctx := attributeRuleContext{ResourceName: resourceName}
	return collections.Map(in, func(attr Attribute) ResolvedAttribute { return resolveAttribute(attr, ctx) })
}

func resolveAttribute(in Attribute, ctx attributeRuleContext) ResolvedAttribute {
	out := attributeBaseline(in, ctx)
	for _, rule := range attributeRules {
		if rule.When(in, ctx) {
			rule.Then(in, ctx, &out)
		}
	}
	if in.Elem != nil {
		elemCtx := ctx
		elemCtx.AsCollectionElem = true
		elem := resolveAttribute(*in.Elem, elemCtx)
		out.Elem = &elem
	}
	return out
}

func attributeBaseline(in Attribute, ctx attributeRuleContext) ResolvedAttribute {
	out := ResolvedAttribute{
		Name:   in.Name,
		TfType: tfTypeFor(in.Type),
	}
	if !ctx.AsCollectionElem {
		out.Required = in.Required
		out.Optional = !in.Required
		out.ForceNew = in.ForceNew
	}
	if in.Description != "" {
		out.DescriptionExpr = strconv.Quote(in.Description)
	}
	if in.Min != nil {
		out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntAtLeast(%d))", *in.Min)
	}
	return out
}

func tfTypeFor(t AttributeType) SchemaValueType {
	switch t {
	case AttrInt:
		return SchemaTypeInt
	case AttrBool:
		return SchemaTypeBool
	case AttrList:
		return SchemaTypeList
	case AttrSet:
		return SchemaTypeSet
	default:
		return SchemaTypeString
	}
}
