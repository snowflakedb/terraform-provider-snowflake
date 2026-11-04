package gen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func showResultSchemaDetailsFromSsot(obj ssot.Object, isDescribe bool) ShowResultSchemaDetails {
	g := ssot.ObjectGenerators(obj.Generators...)
	name := "sdk." + obj.Name
	attrs := obj.ShowOutputAttributes()
	manualAttrs := g.Schema.ShowManualFields
	listAttrs := g.Schema.ShowListFields
	if isDescribe {
		name += "Details"
		attrs = obj.DescribeOutputAttributes()
		manualAttrs = g.Schema.DescribeManualFields
		listAttrs = g.Schema.DescribeListFields
	}
	fields := make([]genhelpers.Field, 0, len(attrs))
	var skip []string
	var overrides map[string]schema.ValueType
	keepList := make(map[*ssot.Attribute]struct{}, len(listAttrs))
	for _, attr := range listAttrs {
		keepList[attr] = struct{}{}
	}
	manualSet := make(map[*ssot.Attribute]struct{}, len(manualAttrs))
	for _, attr := range manualAttrs {
		manualSet[attr] = struct{}{}
	}
	for _, attr := range attrs {
		if attr == nil {
			kind := "SHOW"
			if isDescribe {
				kind = "DESCRIBE"
			}
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		fields = append(fields, ssot.ColumnField(obj, *attr))
		if attr.SkipInSchema {
			skip = append(skip, emittedSchemaKey(attr))
		}
		// The mapper sees only the SDK slice, so a set column needs an explicit override.
		// Skipped, manual, and list-kept columns are left to the mapper.
		_, keep := keepList[attr]
		_, manualAttr := manualSet[attr]
		if attr.Type == ssot.AttrSet && !attr.SkipInSchema && !keep && !manualAttr {
			if overrides == nil {
				overrides = make(map[string]schema.ValueType)
			}
			overrides[emittedSchemaKey(attr)] = schema.TypeSet
		}
	}
	var manual []string
	if len(manualAttrs) > 0 {
		manual = collections.Map(manualAttrs, emittedSchemaKey)
	}
	return ShowResultSchemaDetails{
		IsDescribe:    isDescribe,
		SkipFields:    skip,
		ManualFields:  manual,
		TypeOverrides: overrides,
		StructDetails: genhelpers.StructDetails{Name: name, Fields: fields},
	}
}

// emittedSchemaKey is the schema key ModelFromStructDetails matches.
// SchemaName is that key. Otherwise the column keeps OutputKey, which is the generated key for a snake_case column.
func emittedSchemaKey(attr *ssot.Attribute) string {
	if attr.SchemaName != "" {
		return attr.SchemaName
	}
	return attr.OutputKey()
}

// stampSchemaNames copies SchemaName onto reflected fields. The index key is the SDK type name.
// A SchemaName whose Go field is missing from the struct panics.
func stampSchemaNames(details genhelpers.StructDetails, index map[string]map[string]string) genhelpers.StructDetails {
	byField := index[details.Name]
	if len(byField) == 0 {
		return details
	}
	found := make(map[string]struct{}, len(byField))
	for i := range details.Fields {
		schemaName, ok := byField[details.Fields[i].Name]
		if !ok {
			continue
		}
		details.Fields[i].SchemaName = schemaName
		found[details.Fields[i].Name] = struct{}{}
	}
	missing := make([]string, 0)
	for goName := range byField {
		if _, ok := found[goName]; !ok {
			missing = append(missing, goName)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		panic(fmt.Sprintf("%s: schema name for %s is not on the struct", details.Name, strings.Join(missing, ", ")))
	}
	return details
}
