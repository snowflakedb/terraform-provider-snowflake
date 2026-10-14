package gen

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

// SnowflakeObjectAssertionsFromSsot builds the show or describe object assertion model for one opted-in object.
// Columns become genhelpers.Field values, then ModelFromSdkObjectDetails, the same mapper the reflected path uses.
func SnowflakeObjectAssertionsFromSsot(obj ssot.Object, isDescribe bool) SnowflakeObjectAssertionsModel {
	return ModelFromSdkObjectDetails(objectAssertionDetailsFromSsot(obj, isDescribe), nil, nil)
}

func objectAssertionDetailsFromSsot(obj ssot.Object, isDescribe bool) genhelpers.SdkObjectDetails {
	g := ssot.ObjectGenerators(obj.Generators...)
	kind := "SHOW"
	enabled := g.ObjectAssert.Show
	attrs := obj.ShowOutputAttributes()
	if isDescribe {
		kind = "DESCRIBE"
		enabled = g.ObjectAssert.Describe
		attrs = obj.DescribeOutputAttributes()
	}
	if !enabled {
		panic(fmt.Sprintf("object %q is not enabled for snowflake object assertion %s generation", obj.Name, kind))
	}
	if len(attrs) == 0 {
		panic(fmt.Sprintf("object %q is enabled for snowflake object assertion %s but has no %s fields", obj.Name, kind, kind))
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
	noIdentifiableObject := false
	if isDescribe {
		noIdentifiableObject = g.ObjectAssert.NoIdentifiableObject
	}
	return genhelpers.SdkObjectDetails{
		IdType:               "sdk." + obj.IdType,
		IsDataSourceOutput:   isDescribe,
		NoIdentifiableObject: noIdentifiableObject,
		StructDetails: genhelpers.StructDetails{
			Name:   name,
			Fields: fields,
		},
	}
}
