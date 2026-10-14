package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// ShowOutputAssertionsFromSsot builds the show or describe assertion model for one opted-in object.
// Columns become genhelpers.Field values, then MapToResourceShowOutputAssertion, the same mapper the reflected path uses.
func ShowOutputAssertionsFromSsot(obj ssot.Object, isDescribe bool) ResourceShowOutputAssertionsModel {
	return ModelFromSdkObjectDetails(showOutputDetailsFromSsot(obj, isDescribe), nil)
}

func showOutputDetailsFromSsot(obj ssot.Object, isDescribe bool) SdkObjectShowOutputDetails {
	g := ssot.ObjectGenerators(obj.Generators...)
	kind := "SHOW"
	enabled := g.ShowOutputAssert.Show
	attrs := obj.ShowOutputAttributes()
	if isDescribe {
		kind = "DESCRIBE"
		enabled = g.ShowOutputAssert.Describe
		attrs = obj.DescribeOutputAttributes()
	}
	if !enabled {
		panic(fmt.Sprintf("object %q is not enabled for show output assertion %s generation", obj.Name, kind))
	}
	if len(attrs) == 0 {
		panic(fmt.Sprintf("object %q is enabled for show output assertion %s but has no %s fields", obj.Name, kind, kind))
	}
	name := "sdk." + obj.Name
	if isDescribe {
		name += "Details"
	}
	fields := make([]genhelpers.Field, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		fields = append(fields, ssot.ColumnField(obj, *attr))
	}
	dataSource, _ := dataSourceMappingNormalized[name]
	return SdkObjectShowOutputDetails{
		dataSourceDef: &dataSource,
		SdkObjectDetails: genhelpers.SdkObjectDetails{
			IsDataSourceOutput: isDescribe,
			StructDetails: genhelpers.StructDetails{
				Name:   name,
				Fields: fields,
			},
		},
	}
}
