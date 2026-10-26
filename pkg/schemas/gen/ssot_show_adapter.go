package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

func showResultSchemaDetailsFromSsot(obj ssot.Object, isDescribe bool) ShowResultSchemaDetails {
	g := ssot.ObjectGenerators(obj.Generators...)
	name := "sdk." + obj.Name
	attrs := obj.ShowOutputAttributes()
	manualAttrs := g.Schema.ShowManualFields
	if isDescribe {
		name += "Details"
		attrs = obj.DescribeOutputAttributes()
		manualAttrs = g.Schema.DescribeManualFields
	}
	fields := make([]genhelpers.Field, 0, len(attrs))
	var skip []string
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
	}
	var manual []string
	if len(manualAttrs) > 0 {
		manual = collections.Map(manualAttrs, func(attr *ssot.Attribute) string { return attr.OutputKey() })
	}
	return ShowResultSchemaDetails{
		IsDescribe:    isDescribe,
		SkipFields:    skip,
		ManualFields:  manual,
		StructDetails: genhelpers.StructDetails{Name: name, Fields: fields},
	}
}
