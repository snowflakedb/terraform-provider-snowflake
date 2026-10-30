package gen

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// CustomDiffEntry is one ComputedIfAnyAttributeChanged call.
// Key is the output attribute constant. Attrs are the resource attributes that mark it computed.
type CustomDiffEntry struct {
	Key   string
	Attrs []string
}

// CustomDiffEntries is the generated custom diff, in output order:
// show_output, describe_output, then fully_qualified_name.
// An output is included only when that schema hook is present and it has changed keys.
func (m ResourceModel) CustomDiffEntries() []CustomDiffEntry {
	entries := make([]CustomDiffEntry, 0, 3)
	entries = appendDiffEntry(entries, m.Hooks.ShowOutputSchema.Present, "ShowOutputAttributeName", m.changedAttributeKeys(m.ShowColumns))
	entries = appendDiffEntry(entries, m.Hooks.DescribeOutputSchema.Present, "DescribeOutputAttributeName", m.changedAttributeKeys(m.DescribeColumns))
	entries = appendDiffEntry(entries, m.Hooks.FqnSchema.Present, "FullyQualifiedNameAttributeName", m.fqnCustomDiffKeys())
	return entries
}

func appendDiffEntry(entries []CustomDiffEntry, present bool, key string, attrs []string) []CustomDiffEntry {
	if !present || len(attrs) == 0 {
		return entries
	}
	return append(entries, CustomDiffEntry{Key: key, Attrs: attrs})
}

// CustomDiffCalls is every generated customdiff.All argument, in call order:
// show_output, describe_output, parameters computed-if, fully_qualified_name,
// the parameters custom diff, then AdditionalCustomizeDiffs in list order.
func (m ResourceModel) CustomDiffCalls() []string {
	schema := m.Prefix() + "Schema"
	calls := make([]string, 0, 6)
	var fqn string
	for _, entry := range m.CustomDiffEntries() {
		call := formatComputedIf(schema, entry.Key, entry.Attrs)
		if entry.Key == "FullyQualifiedNameAttributeName" {
			fqn = call
			continue
		}
		calls = append(calls, call)
	}
	if call := m.parametersComputedIfCall(); call != "" {
		calls = append(calls, call)
	}
	if fqn != "" {
		calls = append(calls, fqn)
	}
	if m.Hooks.ParametersCustomDiff.Present {
		calls = append(calls, m.hookSymbol("ParametersCustomDiff", m.Hooks.ParametersCustomDiff))
	}
	for _, diff := range m.AdditionalCustomizeDiffs {
		calls = append(calls, string(diff))
	}
	return calls
}

func formatComputedIf(schema, key string, attrs []string) string {
	parts := make([]string, 0, 2+len(attrs))
	parts = append(parts, schema, key)
	for _, attr := range attrs {
		parts = append(parts, fmt.Sprintf("%q", attr))
	}
	return "ComputedIfAnyAttributeChanged(" + strings.Join(parts, ", ") + ")"
}

// parametersComputedIfCall watches the parameters block when both parameters schema hooks are present.
func (m ResourceModel) parametersComputedIfCall() string {
	if !m.Hooks.ParametersAttributesSchema.Present || !m.Hooks.ParametersFieldNames.Present {
		return ""
	}
	schema := m.hookSymbol("ParametersAttributesSchema", m.Hooks.ParametersAttributesSchema)
	fields := m.hookSymbol("ParametersFieldNames", m.Hooks.ParametersFieldNames)
	return fmt.Sprintf("ComputedIfAnyAttributeChanged(%s, ParametersAttributeName, %s...)", schema, fields)
}

func (m ResourceModel) hookSymbol(base string, hook ssot.ResourceHook) string {
	name := m.Prefix() + base
	if !hook.Generate {
		name += "Ext"
	}
	return name
}

func (m ResourceModel) identityFields() []ssot.ClassifiedField {
	idType := m.IdType
	if idType == "" {
		idType = "AccountObjectIdentifier"
	}
	return ssot.IdentityFields(ssot.Object{
		Name:                      m.Name,
		IdType:                    idType,
		IdentityNameForceNew:      m.IdentityNameForceNew,
		IdentityHierarchyForceNew: m.IdentityHierarchyForceNew,
	})
}

// changedAttributeKeys keeps show or describe columns that are top-level resource attributes.
// Nested columns are excluded. A column whose name is not a resource attribute is excluded, except for a read-only column
// whose ShowName is one (a DESCRIBE-only column that reads the property of that attribute in another type).
// The key is the resource attribute name, which is what ComputedIfAnyAttributeChanged watches.
func (m ResourceModel) changedAttributeKeys(columns []*ssot.Attribute) []string {
	names := m.diffAttributeNames()
	keys := make([]string, 0, len(columns))
	for _, col := range columns {
		if col == nil || (len(col.Children) > 0 && !col.SdkString) {
			continue
		}
		// For now, the set/list fields have to be excluded.
		// TODO [SNOW-1648997]: address the above comment
		if col.Type == ssot.AttrList || col.Type == ssot.AttrSet {
			continue
		}
		key := col.Name
		if _, ok := names[key]; !ok && col.Readonly && col.ShowName != "" {
			// A read-only column with another name reads the same property as the resource attribute of its output key.
			key = col.ShowName
		}
		if _, ok := names[key]; !ok {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}

// diffAttributeNames is the resource schema keys a diff can change:
// identifier attributes, plus attributes that are not nested blocks or sets or lists.
func (m ResourceModel) diffAttributeNames() map[string]struct{} {
	names := make(map[string]struct{}, len(m.Attributes)+3)
	for _, field := range m.identityFields() {
		names[field.Name] = struct{}{}
	}
	for _, attr := range m.Attributes {
		if !attr.SdkString && (len(attr.Children) > 0 || attr.TfType == SchemaTypeList || attr.TfType == SchemaTypeSet) {
			continue
		}
		names[attr.Name] = struct{}{}
	}
	return names
}

// fqnCustomDiffKeys is the identifier attributes that are not ForceNew, in IdentityFields order.
func (m ResourceModel) fqnCustomDiffKeys() []string {
	keys := make([]string, 0, 3)
	for _, field := range m.identityFields() {
		if !field.ForceNew {
			keys = append(keys, field.Name)
		}
	}
	return keys
}
