package gen

import (
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
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

// resourceAttributeRuleContext distinguishes a top-level attribute from a collection
// elem (list or set) so rules that only apply to fields (enum DiffSuppress,
// descriptions) can skip elems.
type resourceAttributeRuleContext struct {
	AsCollectionElem bool
}

type resourceAttributeRule struct {
	Name        string
	Description string
	When        func(ssot.Attribute, resourceAttributeRuleContext) bool
	Then        func(ssot.Attribute, resourceAttributeRuleContext, *ResourceAttribute)
}

// resourceAttributeRules is the ordered policy applied by applyResourceAttributeRules.
// When two rules could touch the same field they are disjoint in When.
var resourceAttributeRules = []resourceAttributeRule{
	{
		// That is not always true, there are at least two groups of optional bool attributes that don't use this pattern:
		// - Create-time-only flags, e.g. snowflake_warehouse.initially_suspended
		// - Discriminator flags inside blocks (snowflake_session_policy.allowed_secondary_roles.0.all).
		// TODO [next PRs]: adjust the logic if needed with next objects migrations
		Name:        RuleOptionalBoolIsBooleanString,
		Description: "Optional bools use boolean-string chrome because TF TypeBool cannot distinguish unset from false.",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrBool && !in.Required
		},
		Then: func(in ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
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
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrBool && in.Required
		},
		Then: func(_ ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.TfType = SchemaTypeBool
		},
	},
	{
		Name:        RuleOptionalIntMinAllowsZeroUsesIntDefault,
		Description: "Optional ints use IntDefault only when Min shows 0 is an acceptable value (otherwise 0 already means unset).",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrInt && !in.Required && in.Min != nil && *in.Min <= 0
		},
		Then: func(_ ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DefaultExpr = "IntDefault"
		},
	},
	{
		Name:        RuleIdentifierSuppressesQuoting,
		Description: "Identifier attributes suppress quoting differences in diffs.",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrIdentifier
		},
		Then: func(_ ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DiffSuppressExpr = "suppressIdentifierQuoting"
		},
	},
	{
		Name:        RuleEnumValidates,
		Description: "Enum attributes (including collection elems) validate through the SDK converter.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrEnum && ctx.sdkEnumTo(in) != ""
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
			out.ValidateDiagExpr = fmt.Sprintf("sdkValidation(%s)", ctx.sdkEnumTo(in))
		},
	},
	{
		Name:        RuleEnumDescriptionListsValues,
		Description: "Enum fields and collections of enums append enumValuesDescription to the description.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			if ctx.AsCollectionElem || in.Description == "" || ctx.sdkEnumAll(in) == "" {
				return false
			}
			return in.Type == ssot.AttrEnum || isCollectionOfEnum(in)
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
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
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrEnum && ctx.sdkEnumTo(in) != "" && !ctx.AsCollectionElem
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DiffSuppressExpr = fmt.Sprintf("SuppressIfAny(NormalizeAndCompare(%s))", ctx.sdkEnumTo(in))
		},
	},
	{
		Name:        RuleCollectionOfEnumNormalizes,
		Description: "Collections of enums (list or set) normalize values in diffs via NormalizeAndCompare on the collection field.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return isCollectionOfEnum(in) && ctx.sdkEnumTo(in) != "" && !ctx.AsCollectionElem
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DiffSuppressExpr = fmt.Sprintf("NormalizeAndCompare(%s)", ctx.sdkEnumTo(in))
		},
	},
}

func isCollection(in ssot.Attribute) bool {
	return in.Type == ssot.AttrList || in.Type == ssot.AttrSet
}

func isCollectionOfEnum(in ssot.Attribute) bool {
	return isCollection(in) && in.Elem != nil && in.Elem.Type == ssot.AttrEnum
}

func enumDef(in ssot.Attribute) *ssot.Enum {
	if in.Enum != nil {
		return in.Enum
	}
	if in.Elem != nil && in.Elem.Enum != nil {
		return in.Elem.Enum
	}
	return nil
}

func (ctx resourceAttributeRuleContext) sdkEnumTo(in ssot.Attribute) string {
	e := enumDef(in)
	if e == nil || e.Name == "" {
		return ""
	}
	return "sdk.To" + e.Name
}

func (ctx resourceAttributeRuleContext) sdkEnumAll(in ssot.Attribute) string {
	e := enumDef(in)
	if e == nil || e.NamePlural == "" {
		return ""
	}
	return "sdk.All" + e.NamePlural
}

// applyResourceAttributeRules is the dedicated processing step from written Attribute intent to ResourceAttribute schema chrome.
func applyResourceAttributeRules(in []ssot.Attribute) []ResourceAttribute {
	ctx := resourceAttributeRuleContext{}
	return collections.Map(in, func(attr ssot.Attribute) ResourceAttribute { return resolveResourceAttribute(attr, ctx) })
}

func resolveResourceAttribute(in ssot.Attribute, ctx resourceAttributeRuleContext) ResourceAttribute {
	out := resourceAttributeBaseline(in, ctx)
	for _, rule := range resourceAttributeRules {
		if rule.When(in, ctx) {
			rule.Then(in, ctx, &out)
		}
	}
	if in.Elem != nil {
		elemCtx := ctx
		elemCtx.AsCollectionElem = true
		elem := resolveResourceAttribute(*in.Elem, elemCtx)
		out.Elem = &elem
	}
	return out
}

func resourceAttributeBaseline(in ssot.Attribute, ctx resourceAttributeRuleContext) ResourceAttribute {
	out := ResourceAttribute{
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

func tfTypeFor(t ssot.AttributeType) SchemaValueType {
	switch t {
	case ssot.AttrInt:
		return SchemaTypeInt
	case ssot.AttrBool:
		return SchemaTypeBool
	case ssot.AttrList:
		return SchemaTypeList
	case ssot.AttrSet:
		return SchemaTypeSet
	default:
		return SchemaTypeString
	}
}
