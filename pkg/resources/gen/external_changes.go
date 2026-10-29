package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ExternalChangeMapping is one outputMapping argument.
// OutputKey is the show or describe state key. ConfigKey is the resource attribute.
// Compare matches the value stored in that output. Set is the config value.
type ExternalChangeMapping struct {
	OutputKey string
	ConfigKey string
	Compare   string
	Set       string
}

// ExternalChangeBuild is the generated external-change function body.
// Prepare is the local declarations the mappings read.
type ExternalChangeBuild struct {
	Prepare  []string
	Mappings []ExternalChangeMapping
}

// ShowExternalChanges is the SHOW mappings for a generated HandleExternalChanges.
// An object that keeps the ext call gets an empty build.
func (m ResourceModel) ShowExternalChanges() ExternalChangeBuild {
	if !m.Hooks.HandleExternalChanges.Generate {
		return ExternalChangeBuild{}
	}
	if !m.Hooks.ShowByIdSafelyInSdk.Present {
		panic(fmt.Sprintf("object %q: HandleExternalChanges requires ShowByIdSafelyInSdk", m.Name))
	}
	return m.externalChanges(m.ShowColumns, nil, false, m.Prefix())
}

// DescribeExternalChanges is the DESCRIBE mappings for a generated HandleExternalChangesInDescribe.
// A column is included when it is describe-only and ServerDefault is set.
// An object that keeps the ext call gets an empty build.
func (m ResourceModel) DescribeExternalChanges() ExternalChangeBuild {
	if !m.Hooks.HandleExternalChangesInDescribe.Generate {
		return ExternalChangeBuild{}
	}
	if !m.Hooks.DescribeInSdk.Present {
		panic(fmt.Sprintf("object %q: HandleExternalChangesInDescribe requires DescribeInSdk", m.Name))
	}
	covered := make(map[string]struct{}, len(m.ShowColumns))
	for _, col := range m.ShowColumns {
		if col != nil {
			covered[col.Name] = struct{}{}
		}
	}
	return m.externalChanges(m.DescribeColumns, covered, true, "details")
}

func (m ResourceModel) externalChanges(columns []*ssot.Attribute, covered map[string]struct{}, describe bool, receiver string) ExternalChangeBuild {
	attrs := make(map[string]ResourceAttribute, len(m.Attributes))
	for _, attr := range m.Attributes {
		attrs[attr.Name] = attr
	}
	identity := m.identityNames()
	var build ExternalChangeBuild
	for _, col := range columns {
		if col == nil || !m.includeExternalChange(col, covered, identity, attrs, describe) {
			continue
		}
		build.Mappings = append(build.Mappings, m.mappingFor(col, attrs[col.Name], receiver, &build))
	}
	return build
}

// applyDerivedOutputDiffSuppress adds the show or describe output suppressor for every
// drift-tracked scalar. SHOW wins. A describe-only field needs ServerDefault.
// An explicit SuppressShow or SuppressDescribe on the attribute is left as written.
func (m *ResourceModel) applyDerivedOutputDiffSuppress(attrs []ssot.Attribute) {
	if len(attrs) != len(m.Attributes) {
		panic(fmt.Sprintf("object %q: %d resource attributes and %d schema attributes", m.Name, len(attrs), len(m.Attributes)))
	}
	show := columnByName(m.ShowColumns)
	describe := columnByName(m.DescribeColumns)
	identity := m.identityNames()
	for i, attr := range attrs {
		if !scalarDriftCandidate(attr, identity) {
			continue
		}
		if col, ok := show[attr.Name]; ok {
			if hasSuppressKind(attr.DiffSuppress, ssot.SuppressShow) {
				continue
			}
			expr := fmt.Sprintf("IgnoreChangeToCurrentSnowflakeValueInShow(%q)", col.OutputKey())
			m.Attributes[i].DiffSuppressExpr = composeDiffSuppress(m.Attributes[i].DiffSuppressExpr, expr)
			continue
		}
		col, ok := describe[attr.Name]
		if !ok || !attr.ServerDefault || hasSuppressKind(attr.DiffSuppress, ssot.SuppressDescribe) {
			continue
		}
		expr := fmt.Sprintf("IgnoreChangeToCurrentSnowflakeValueInDescribe(%q)", col.OutputKey())
		m.Attributes[i].DiffSuppressExpr = composeDiffSuppress(m.Attributes[i].DiffSuppressExpr, expr)
	}
}

func (m ResourceModel) identityNames() map[string]struct{} {
	fields := m.identityFields()
	names := make(map[string]struct{}, len(fields))
	for _, field := range fields {
		names[field.Name] = struct{}{}
	}
	return names
}

func columnByName(columns []*ssot.Attribute) map[string]*ssot.Attribute {
	byName := make(map[string]*ssot.Attribute, len(columns))
	for _, col := range columns {
		if col != nil {
			byName[col.Name] = col
		}
	}
	return byName
}

func scalarDriftCandidate(attr ssot.Attribute, identity map[string]struct{}) bool {
	if attr.Required || attr.Computed || attr.CreateOnly || attr.Readonly || attr.SqlOnly || attr.Name == "comment" {
		return false
	}
	if _, skip := identity[attr.Name]; skip {
		return false
	}
	return len(attr.Children) == 0 && attr.Type != ssot.AttrList && attr.Type != ssot.AttrSet && attr.Type != ssot.AttrObject
}

func hasSuppressKind(suppress ssot.DiffSuppress, kind ssot.DiffSuppressKind) bool {
	if suppress == nil {
		return false
	}
	for _, item := range suppress.Suppressors() {
		if item.Kind == kind {
			return true
		}
	}
	return false
}

func (m ResourceModel) includeExternalChange(col *ssot.Attribute, covered, identity map[string]struct{}, attrs map[string]ResourceAttribute, describe bool) bool {
	if _, skip := covered[col.Name]; skip {
		return false
	}
	attr, ok := attrs[col.Name]
	if !ok || attr.Required || attr.Computed {
		return false
	}
	if _, skip := identity[col.Name]; skip || col.Name == "comment" || col.CreateOnly {
		return false
	}
	if describe && !col.ServerDefault {
		return false
	}
	if len(col.Children) > 0 || col.Type == ssot.AttrList || col.Type == ssot.AttrSet {
		m.panicExternal(col.Name, "nested blocks are not generated")
	}
	return true
}

func (m ResourceModel) mappingFor(col *ssot.Attribute, resolved ResourceAttribute, receiver string, build *ExternalChangeBuild) ExternalChangeMapping {
	goField := receiver + "." + genhelpers.SnakeCaseToCamel(col.OutputKey())
	compare, set := m.outputValues(col, resolved, goField, build)
	return ExternalChangeMapping{OutputKey: col.OutputKey(), ConfigKey: col.Name, Compare: compare, Set: set}
}

// outputValues is the compare and set expressions for one outputMapping.
// An OptionalInShow column prepares a local from the nullable SDK field.
func (m ResourceModel) outputValues(col *ssot.Attribute, resolved ResourceAttribute, goField string, build *ExternalChangeBuild) (string, string) {
	local := genhelpers.FirstLetterLowercase(genhelpers.SnakeCaseToCamel(col.Name))
	switch col.Type {
	case ssot.AttrString, ssot.AttrEnum:
		if col.OptionalInShow {
			val, str := local+"Val", local+"Str"
			build.Prepare = append(build.Prepare, fmt.Sprintf("%s, %s := optionalStringOutputMapping(%s)", val, str, goField))
			return str, val
		}
		if col.Type == ssot.AttrEnum {
			return "string(" + goField + ")", goField
		}
		return goField, goField
	case ssot.AttrBool:
		if col.OptionalInShow {
			if resolved.TfType != SchemaTypeString {
				m.panicExternal(col.Name, "optional native bool is not generated")
			}
			val, str := local+"Val", local+"Str"
			build.Prepare = append(build.Prepare, fmt.Sprintf("%s, %s := optionalBooleanStringOutputMapping(%s)", val, str, goField))
			return val, str
		}
		if resolved.TfType == SchemaTypeString {
			return goField, "booleanStringFromBool(" + goField + ")"
		}
		return goField, goField
	case ssot.AttrInt:
		if !col.OptionalInShow {
			return goField, goField
		}
		helper := "optionalIntOutputMapping"
		if resolved.DefaultExpr == "IntDefault" {
			helper = "optionalIntOutputMappingIntDefault"
		}
		build.Prepare = append(build.Prepare, fmt.Sprintf("%s := %s(%s)", local, helper, goField))
		return local, local
	case ssot.AttrIdentifier:
		if col.OptionalInShow {
			m.panicExternal(col.Name, "optional %q is not generated", col.Type)
		}
		expr := goField + ".Name()"
		if col.IdentifierType != "AccountObjectIdentifier" {
			expr = goField + ".FullyQualifiedName()"
		}
		return expr, expr
	default:
		if col.OptionalInShow {
			m.panicExternal(col.Name, "optional %q is not generated", col.Type)
		}
		m.panicExternal(col.Name, "unsupported type %q", col.Type)
	}
	return "", ""
}

func (m ResourceModel) panicExternal(attr, format string, args ...any) {
	panic(fmt.Sprintf("object %q: external change attribute %q: "+format, append([]any{m.Name, attr}, args...)...))
}
