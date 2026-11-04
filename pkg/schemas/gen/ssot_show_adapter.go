package gen

import (
	"fmt"

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
			skip = append(skip, attr.OutputKey())
		}
		// The mapper sees only the SDK slice, so a set column needs an explicit override.
		// Skipped, manual, and list-kept columns are left to the mapper.
		_, keep := keepList[attr]
		_, manualAttr := manualSet[attr]
		if attr.Type == ssot.AttrSet && !attr.SkipInSchema && !keep && !manualAttr {
			if overrides == nil {
				overrides = make(map[string]schema.ValueType)
			}
			overrides[attr.OutputKey()] = schema.TypeSet
		}
	}
	var manual []string
	if len(manualAttrs) > 0 {
		manual = collections.Map(manualAttrs, func(attr *ssot.Attribute) string { return attr.OutputKey() })
	}
	return ShowResultSchemaDetails{
		IsDescribe:    isDescribe,
		SkipFields:    skip,
		ManualFields:  manual,
		TypeOverrides: overrides,
		StructDetails: genhelpers.StructDetails{Name: name, Fields: fields},
	}
}
