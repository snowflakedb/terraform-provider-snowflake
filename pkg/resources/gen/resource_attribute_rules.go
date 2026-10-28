package gen

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

const (
	RuleOptionalBoolIsBooleanString            = "optional_bool_is_boolean_string"
	RuleCreateOnlyBoolIgnoresChanges           = "create_only_bool_ignores_changes"
	RuleCreateOnlyNativeBoolIsNative           = "create_only_native_bool_is_native"
	RuleRequiredBoolIsNative                   = "required_bool_is_native"
	RuleOptionalIntMinAllowsZeroUsesIntDefault = "optional_int_min_allows_zero_uses_int_default"
	RuleIdentifierSuppressesQuoting            = "identifier_suppresses_quoting"
	RuleIdentifierValidates                    = "identifier_validates"
	RuleRelatedResourceDescription             = "related_resource_description"
	RuleEnumValidates                          = "enum_validates"
	RuleEnumDescriptionListsValues             = "enum_description_lists_values"
	RuleEnumFieldNormalizes                    = "enum_field_normalizes"
	RuleCollectionOfEnumNormalizes             = "collection_of_enum_normalizes"
	RuleSetOfIdentifiersNormalizes             = "set_of_identifiers_normalizes"
	RuleConfiguredDiffSuppress                 = "configured_diff_suppress"
	RuleSingleObjectIsListOfOne                = "single_object_is_list_of_one"
)

// resourceAttributeRuleContext distinguishes a top-level attribute from a collection
// elem (list or set) and from a field inside a nested block, so rules that only apply
// to fields (enum DiffSuppress, descriptions, top-level set normalization) can skip the others.
// RelationPrefix is the Terraform path of this attribute's parent, without a trailing dot.
// Empty at the resource root. A list with MaxItems 1 uses parent.0.
type resourceAttributeRuleContext struct {
	AsCollectionElem bool
	InBlock          bool
	IdType           string
	RelationPrefix   string
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
		// NativeBool skips this rule. Create-only bools still use boolean strings, plus IgnoreAfterCreation.
		// CreateOnly together with NativeBool is RuleCreateOnlyNativeBoolIsNative (warehouse initially_suspended).
		// Discriminator flags inside blocks (snowflake_session_policy.allowed_secondary_roles.0.all)
		// are still boolean strings.
		// TODO [next PRs]: adjust the logic if needed with next objects migrations
		Name:        RuleOptionalBoolIsBooleanString,
		Description: "Optional bools use boolean-string chrome because TF TypeBool cannot distinguish unset from false.",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrBool && !in.Required && !in.NativeBool
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
		Name:        RuleCreateOnlyBoolIgnoresChanges,
		Description: "Create-only optional bools ignore changes after creation.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrBool && !in.Required && in.CreateOnly && !ctx.AsCollectionElem
		},
		Then: func(_ ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DiffSuppressExpr = "IgnoreAfterCreation"
		},
	},
	{
		// Warehouse initially_suspended was native TypeBool because CreateOnly used to select that type.
		// NativeBool marks the attributes that keep that schema.
		Name:        RuleCreateOnlyNativeBoolIsNative,
		Description: "Create-only optional bools with NativeBool stay native TypeBool.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrBool && !in.Required && in.CreateOnly && in.NativeBool && !ctx.AsCollectionElem
		},
		Then: func(_ ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.TfType = SchemaTypeBool
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
		Name:        RuleIdentifierValidates,
		Description: "Identifier fields validate with IsValidIdentifier. A collection element uses its IdentifierType. A top-level field uses the object id type.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrIdentifier && identifierSchemaType(in, ctx) != ""
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
			out.ValidateDiagExpr = fmt.Sprintf("IsValidIdentifier[sdk.%s]()", identifierSchemaType(in, ctx))
		},
	},
	{
		Name:        RuleRelatedResourceDescription,
		Description: "Attributes that name another resource append a docs link via relatedResourceDescription.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.RelatedResource != nil && in.RelatedResource.Name != "" && in.Description != "" && !ctx.AsCollectionElem
		},
		Then: func(in ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DescriptionExpr = fmt.Sprintf("relatedResourceDescription(%s, resources.%s)", strconv.Quote(in.Description), in.RelatedResource.Name)
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
		Description: "A top-level set of enums uses NormalizeAndCompareEnumsInSet. A list, and a nested set, use NormalizeAndCompare on the collection field.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return isCollectionOfEnum(in) && ctx.sdkEnumTo(in) != "" && !ctx.AsCollectionElem
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
			if in.Type == ssot.AttrSet && !ctx.InBlock {
				out.DiffSuppressExpr = fmt.Sprintf("NormalizeAndCompareEnumsInSet(%q, %s)", in.Name, ctx.sdkEnumTo(in))
				return
			}
			out.DiffSuppressExpr = fmt.Sprintf("NormalizeAndCompare(%s)", ctx.sdkEnumTo(in))
		},
	},
	{
		Name:        RuleSetOfIdentifiersNormalizes,
		Description: "A top-level set of identifiers uses NormalizeAndCompareIdentifiersInSet.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return isSetOfIdentifiers(in) && !ctx.AsCollectionElem && !ctx.InBlock
		},
		Then: func(in ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.DiffSuppressExpr = fmt.Sprintf("NormalizeAndCompareIdentifiersInSet(%q)", in.Name)
		},
	},
	{
		Name:        RuleSingleObjectIsListOfOne,
		Description: "A single object is a list of at most one element, because SDKv2 cannot declare a nested object: MaxItems 1, and MinItems 1 when it is required.",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrObject
		},
		Then: func(in ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.TfType = SchemaTypeList
			out.MaxItems = new(1)
			if in.Required {
				out.MinItems = new(1)
			}
		},
	},
	{
		Name:        RuleConfiguredDiffSuppress,
		Description: "Configured suppressors compose with inferred normalization.",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.DiffSuppress != nil
		},
		Then: func(in ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			for _, suppress := range in.DiffSuppress.Suppressors() {
				key := suppress.Key
				if key == "" {
					key = in.OutputKey()
				}
				var expression string
				switch suppress.Kind {
				case ssot.SuppressStatement:
					expression = "DiffSuppressStatement"
				case ssot.SuppressShow:
					expression = fmt.Sprintf("IgnoreChangeToCurrentSnowflakeValueInShow(%q)", key)
				case ssot.SuppressDescribe:
					expression = fmt.Sprintf("IgnoreChangeToCurrentSnowflakeValueInDescribe(%q)", key)
				default:
					panic(fmt.Sprintf("unknown diff suppression %q", suppress.Kind))
				}
				out.DiffSuppressExpr = composeDiffSuppress(out.DiffSuppressExpr, expression)
			}
		},
	},
}

func isCollection(in ssot.Attribute) bool {
	return in.Type == ssot.AttrList || in.Type == ssot.AttrSet
}

func isCollectionOfEnum(in ssot.Attribute) bool {
	return isCollection(in) && in.Elem != nil && in.Elem.Type == ssot.AttrEnum
}

func isSetOfIdentifiers(in ssot.Attribute) bool {
	return in.Type == ssot.AttrSet && in.Elem != nil && in.Elem.Type == ssot.AttrIdentifier
}

// identifierSchemaType is the SDK identifier type used for schema validation.
// Collection elements use their own IdentifierType. Nested block fields are left alone.
// A top-level field uses the object id type, then IdentifierType.
func identifierSchemaType(in ssot.Attribute, ctx resourceAttributeRuleContext) string {
	if ctx.AsCollectionElem {
		return in.IdentifierType
	}
	if ctx.InBlock {
		return ""
	}
	if ctx.IdType != "" {
		return ctx.IdType
	}
	return in.IdentifierType
}

func relationPaths(prefix string, attrs []ssot.Attribute) []string {
	if len(attrs) == 0 {
		return nil
	}
	return collections.Map(attrs, func(attr ssot.Attribute) string {
		if prefix == "" {
			return attr.Name
		}
		return prefix + "." + attr.Name
	})
}

// listBlockRelationPrefix is the Terraform path of children of a list with MaxItems 1.
// Other blocks cannot name child relations, so their children keep an empty prefix.
func listBlockRelationPrefix(parent ssot.Attribute, ctx resourceAttributeRuleContext) string {
	if parent.Type != ssot.AttrList || parent.MaxItems == nil || *parent.MaxItems != 1 {
		return ""
	}
	segment := parent.Name + ".0"
	if ctx.RelationPrefix == "" {
		return segment
	}
	return ctx.RelationPrefix + "." + segment
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
	symbol := "sdk.All" + e.NamePlural
	if in.IncludeAliasedInTheDescription {
		symbol += "String"
	}
	return symbol
}

func composeDiffSuppress(existing, ignore string) string {
	if existing == "" {
		return ignore
	}
	const prefix = "SuppressIfAny("
	if strings.HasPrefix(existing, prefix) && strings.HasSuffix(existing, ")") {
		inner := existing[len(prefix) : len(existing)-1]
		return fmt.Sprintf("SuppressIfAny(%s, %s)", inner, ignore)
	}
	return fmt.Sprintf("SuppressIfAny(%s, %s)", existing, ignore)
}

// applyResourceAttributeRules is the dedicated processing step from written Attribute intent to ResourceAttribute schema chrome.
func applyResourceAttributeRules(idType string, in []ssot.Attribute) []ResourceAttribute {
	ctx := resourceAttributeRuleContext{IdType: idType}
	return collections.Map(in, func(attr ssot.Attribute) ResourceAttribute { return resolveResourceAttribute(attr, ctx) })
}

func resolveResourceAttribute(in ssot.Attribute, ctx resourceAttributeRuleContext) ResourceAttribute {
	out := resourceAttributeBaseline(in, ctx)
	for _, rule := range resourceAttributeRules {
		if rule.When(in, ctx) {
			rule.Then(in, ctx, &out)
		}
	}
	if in.Elem != nil && !in.IsObjectList() {
		elemCtx := ctx
		elemCtx.AsCollectionElem = true
		elem := resolveResourceAttribute(*in.Elem, elemCtx)
		out.Elem = &elem
	}
	if fields := in.NestedFields(); len(fields) > 0 {
		nestedCtx := ctx
		nestedCtx.IdType = ""
		nestedCtx.InBlock = true
		nestedCtx.AsCollectionElem = false
		nestedCtx.RelationPrefix = listBlockRelationPrefix(in, ctx)
		out.Children = collections.Map(fields, func(child ssot.Attribute) ResourceAttribute {
			return resolveResourceAttribute(child, nestedCtx)
		})
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
		out.Sensitive = in.Sensitive
		out.MinItems = in.MinItems
		out.MaxItems = in.MaxItems
		out.ConflictsWith = relationPaths(ctx.RelationPrefix, in.ConflictsWith)
		out.AtLeastOneOf = relationPaths(ctx.RelationPrefix, in.AtLeastOneOf)
	}
	if in.Description != "" {
		out.DescriptionExpr = strconv.Quote(in.Description)
	}
	switch {
	case in.Min != nil && in.Max != nil:
		out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntBetween(%d, %d))", *in.Min, *in.Max)
	case in.Min != nil:
		out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntAtLeast(%d))", *in.Min)
	case in.Max != nil:
		out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntAtMost(%d))", *in.Max)
	}
	if in.Type == ssot.AttrEnum && in.Enum != nil {
		out.EnumName = in.Enum.Name
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

// applyExactlyOneOf sets ExactlyOneOf on the attributes of each group: the names of the whole group, itself included.
// The attributes are the top-level ones of the object, in the order of the groups' own attributes.
func applyExactlyOneOf(groups [][]*ssot.Attribute, attributes []ResourceAttribute) {
	for _, group := range groups {
		names := collections.Map(group, func(attr *ssot.Attribute) string { return attr.Name })
		for _, name := range names {
			i := slices.IndexFunc(attributes, func(attr ResourceAttribute) bool { return attr.Name == name })
			if i < 0 {
				panic(fmt.Sprintf("ExactlyOneOf attribute %q is not a resource attribute", name))
			}
			attributes[i].ExactlyOneOf = names
		}
	}
}
