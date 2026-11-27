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
	RuleRequiredListHasMinOne                  = "required_list_has_min_one"
)

// resourceAttributeRuleContext distinguishes a top-level attribute from a collection
// elem (list or set) and from a field inside a nested block, so rules that only apply
// to fields (enum DiffSuppress, descriptions, top-level set normalization) can skip the others.
// RelationPrefix is the Terraform path of this attribute's parent, without a trailing dot.
// Empty at the resource root. An AttrObject uses parent.0.
// InListBlock marks a field under an element of a list or set of objects, which a schema path cannot name.
type resourceAttributeRuleContext struct {
	AsCollectionElem bool
	InBlock          bool
	InListBlock      bool
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
		Description: "Identifier fields validate with IsValidIdentifier. A collection element uses its IdentifierType. A top-level field uses its IdentifierType when that is set, and the object id type otherwise.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrIdentifier && identifierSchemaType(in, ctx) != ""
		},
		Then: func(in ssot.Attribute, ctx resourceAttributeRuleContext, out *ResourceAttribute) {
			if idType := identifierSchemaType(in, ctx); idType == "AccountIdentifier" {
				out.ValidateDiagExpr = "IsValidAccountIdentifier()"
				return
			}
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
			return in.Type == ssot.AttrEnum || in.IsEnumCollection()
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
			return in.IsEnumCollection() && ctx.sdkEnumTo(in) != "" && !ctx.AsCollectionElem
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
			// SDKv2 ignores MinItems on Optional attributes (core_schema.go resets it to 0), so only Required sets MinItems: 1.
			if in.Required {
				out.MinItems = new(1)
			}
		},
	},
	{
		Name:        RuleRequiredListHasMinOne,
		Description: "Required lists require at least one item unless a higher Min is set.",
		When: func(in ssot.Attribute, ctx resourceAttributeRuleContext) bool {
			return in.Type == ssot.AttrList && in.Required && in.Min == nil && !ctx.AsCollectionElem
		},
		Then: func(_ ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			out.MinItems = new(1)
		},
	},
	{
		Name:        RuleConfiguredDiffSuppress,
		Description: "Configured suppressors compose with inferred normalization. Statement suppression documents SQL normalization. Ignore-always stands alone and replaces inferred normalization.",
		When: func(in ssot.Attribute, _ resourceAttributeRuleContext) bool {
			return in.DiffSuppress != nil
		},
		Then: func(in ssot.Attribute, _ resourceAttributeRuleContext, out *ResourceAttribute) {
			suppressors := in.DiffSuppress.Suppressors()
			if len(suppressors) == 1 && suppressors[0].Kind == ssot.SuppressAlways {
				out.DiffSuppressExpr = "ignoreAlwaysSuppressFunc"
				return
			}
			for _, suppress := range suppressors {
				key := suppress.Key
				if key == "" {
					key = in.EffectiveShowName()
				}
				var expression string
				switch suppress.Kind {
				case ssot.SuppressStatement:
					expression = "DiffSuppressStatement"
					if out.DescriptionExpr != "" {
						out.DescriptionExpr = fmt.Sprintf("diffSuppressStatementFieldDescription(%s)", out.DescriptionExpr)
					}
				case ssot.SuppressShow:
					expression = fmt.Sprintf("IgnoreChangeToCurrentSnowflakeValueInShow(%q)", key)
				case ssot.SuppressDescribe:
					expression = fmt.Sprintf("IgnoreChangeToCurrentSnowflakeValueInDescribe(%q)", key)
				case ssot.SuppressAlways:
					panic(fmt.Sprintf("ignore-always suppression for attribute %q stands alone", in.Name))
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

func isSetOfIdentifiers(in ssot.Attribute) bool {
	return in.Type == ssot.AttrSet && in.Elem != nil && in.Elem.Type == ssot.AttrIdentifier
}

// identifierSchemaType is the SDK identifier type used for schema validation.
// Collection elements use their own IdentifierType. Nested block fields are left alone.
// A top-level field uses its IdentifierType when that is set, so a schema object can
// point at a database or account object. Otherwise it uses the object id type.
func identifierSchemaType(in ssot.Attribute, ctx resourceAttributeRuleContext) string {
	// External identifiers are not parsed by IsValidIdentifier. Account identifiers use IsValidAccountIdentifier.
	if in.IdentifierType == "ExternalObjectIdentifier" {
		return ""
	}
	if ctx.AsCollectionElem {
		return in.IdentifierType
	}
	if ctx.InBlock {
		return ""
	}
	if in.IdentifierType != "" {
		return in.IdentifierType
	}
	if ctx.IdType != "" {
		return ctx.IdType
	}
	return ""
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

// listBlockRelationPrefix is the Terraform path of children of a single object (AttrObject).
// List and set blocks cannot name child relations, so their children keep an empty prefix.
func listBlockRelationPrefix(parent ssot.Attribute, ctx resourceAttributeRuleContext) string {
	if !parent.IsSingleBlock() {
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
		nestedCtx.InListBlock = ctx.InListBlock || !in.IsSingleBlock()
		nestedCtx.RelationPrefix = listBlockRelationPrefix(in, ctx)
		out.Children = collections.Map(fields, func(child ssot.Attribute) ResourceAttribute {
			return resolveResourceAttribute(child, nestedCtx)
		})
		if !nestedCtx.InListBlock {
			applyValidations(in.Validations, func(v ssot.Validation) []ssot.Attribute {
				return in.ValidationFieldsFor(v, ssot.ClauseCreate)
			}, nestedCtx.RelationPrefix, out.Children)
		}
	}
	return out
}

func resourceAttributeBaseline(in ssot.Attribute, ctx resourceAttributeRuleContext) ResourceAttribute {
	out := ResourceAttribute{
		Name:      in.Name,
		TfType:    tfTypeFor(in.Type),
		SdkString: in.SdkString,
	}
	if !ctx.AsCollectionElem {
		out.Required = in.Required
		out.Optional = !in.Required && !in.Computed
		out.Computed = in.Computed
		out.ForceNew = in.ForceNew
		out.Sensitive = in.Sensitive
		if in.Type == ssot.AttrList || in.Type == ssot.AttrSet {
			out.MinItems = in.Min
			out.MaxItems = in.Max
		}
	}
	if in.Description != "" {
		out.DescriptionExpr = strconv.Quote(in.Description)
	}
	if in.Deprecated != "" {
		out.DeprecatedExpr = strconv.Quote(in.Deprecated)
	}
	if in.Type == ssot.AttrInt {
		switch {
		case in.Min != nil && in.Max != nil:
			out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntBetween(%d, %d))", *in.Min, *in.Max)
		case in.Min != nil:
			out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntAtLeast(%d))", *in.Min)
		case in.Max != nil:
			out.ValidateDiagExpr = fmt.Sprintf("validation.ToDiagFunc(validation.IntAtMost(%d))", *in.Max)
		}
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

// applyValidations sets the schema relations of each validation on the attributes of its fields, as paths from the
// resource root under prefix. AtLeastOne and ExactlyOne give each field the whole list, itself included. Conflicting gives
// each field the others. A validation that names a SqlOnly field has no schema form, and only the SDK checks it.
// fieldsOf returns the fields of a validation in the order of the parent's fields.
func applyValidations(validations []ssot.Validation, fieldsOf func(ssot.Validation) []ssot.Attribute, prefix string, attributes []ResourceAttribute) {
	for _, validation := range validations {
		fields := fieldsOf(validation)
		if slices.ContainsFunc(fields, func(field ssot.Attribute) bool { return field.SqlOnly }) {
			continue
		}
		paths := relationPaths(prefix, fields)
		for i, field := range fields {
			j := slices.IndexFunc(attributes, func(attr ResourceAttribute) bool { return attr.Name == field.Name })
			if j < 0 {
				panic(fmt.Sprintf("%s validation field %q is not a resource attribute", validation.Kind, field.Name))
			}
			switch validation.Kind {
			case ssot.ValidationAtLeastOne:
				attributes[j].AtLeastOneOf = paths
			case ssot.ValidationExactlyOne:
				attributes[j].ExactlyOneOf = paths
			case ssot.ValidationConflicting:
				attributes[j].ConflictsWith = append(attributes[j].ConflictsWith, slices.Delete(slices.Clone(paths), i, i+1)...)
			default:
				panic(fmt.Sprintf("unknown validation kind %q", validation.Kind))
			}
		}
	}
}
