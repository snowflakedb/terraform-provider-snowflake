package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

func showResultSchemaDetailsFromSsot(obj ssot.Object, isDescribe bool) ShowResultSchemaDetails {
	name := "sdk." + obj.Name
	attrs := obj.ShowOutputAttributes()
	if isDescribe {
		name += "Details"
		attrs = obj.DescribeOutputAttributes()
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
	return ShowResultSchemaDetails{
		IsDescribe:    isDescribe,
		SkipFields:    skip,
		StructDetails: genhelpers.StructDetails{Name: name, Fields: fields},
	}
}
