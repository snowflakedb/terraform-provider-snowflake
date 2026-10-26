package gen

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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
	nested := g.ObjectAssert.ShowNestedFields
	if isDescribe {
		nested = g.ObjectAssert.DescribeNestedFields
	}
	fields := make([]genhelpers.Field, 0, len(attrs))
	for _, attr := range attrs {
		if attr == nil {
			panic(fmt.Sprintf("object %q: nil %s column", obj.Name, kind))
		}
		field := ssot.ColumnField(obj, *attr)
		if nestedField, ok := findNestedField(nested, attr); ok {
			field.ConcreteType, field.UnderlyingType = nestedFieldTypes(nestedFieldGoType(obj, nestedField))
		}
		fields = append(fields, field)
	}
	noIdentifiableObject := false
	skip := objectAssertSkipNames(obj, g.ObjectAssert.ShowSkipFields)
	var projection *genhelpers.ShowProjectionDef
	objectTypeName := ""
	if isDescribe {
		noIdentifiableObject = g.ObjectAssert.NoIdentifiableObject
		skip = objectAssertSkipNames(obj, g.ObjectAssert.DescribeSkipFields)
	} else if fromBase, ok := obj.ShowFromBase(); ok {
		projection = &genhelpers.ShowProjectionDef{
			ClientName: fromBase.ClientName,
			Method:     fromBase.Method,
		}
		objectTypeName = fromBase.ClientName
	}
	return genhelpers.SdkObjectDetails{
		IdType:               "sdk." + obj.IdType,
		IsDataSourceOutput:   isDescribe,
		ObjectTypeName:       objectTypeName,
		NoIdentifiableObject: noIdentifiableObject,
		ShowProjection:       projection,
		SkipFields:           skip,
		NestedAssertFields:   objectAssertNestedNames(obj, nested),
		StructDetails: genhelpers.StructDetails{
			Name:   name,
			Fields: fields,
		},
	}
}

func objectAssertSkipNames(obj ssot.Object, skipped []*ssot.Attribute) []string {
	if len(skipped) == 0 {
		return nil
	}
	names := make([]string, len(skipped))
	for i, attr := range skipped {
		names[i] = ssot.ColumnField(obj, *attr).Name
	}
	return names
}

func objectAssertNestedNames(obj ssot.Object, nested []ssot.ObjectAssertNestedField) []string {
	if len(nested) == 0 {
		return nil
	}
	return collections.Map(nested, func(n ssot.ObjectAssertNestedField) string { return ssot.ColumnField(obj, *n.Attr).Name })
}

func findNestedField(nested []ssot.ObjectAssertNestedField, attr *ssot.Attribute) (ssot.ObjectAssertNestedField, bool) {
	found, err := collections.FindFirst(nested, func(n ssot.ObjectAssertNestedField) bool { return n.Attr == attr })
	if err != nil {
		return ssot.ObjectAssertNestedField{}, false
	}
	return *found, true
}

// nestedFieldGoType is the explicit GoType, or the default composition for a list or set column.
func nestedFieldGoType(obj ssot.Object, nested ssot.ObjectAssertNestedField) string {
	if nested.GoType != "" {
		return nested.GoType
	}
	if nested.Attr.Type != ssot.AttrList && nested.Attr.Type != ssot.AttrSet {
		panic(fmt.Sprintf("object %q: nested column %q is not a list or set, set ObjectAssertNestedGoType", obj.Name, nested.Attr.Name))
	}
	item := strings.TrimSuffix(genhelpers.SnakeCaseToCamel(nested.Attr.OutputKey()), "s")
	return "[]sdk." + obj.Name + item + "Details"
}

// nestedFieldTypes derives the underlying kind from the Go type of a nested-block field.
func nestedFieldTypes(goType string) (concrete, underlying string) {
	pointer := ""
	base := goType
	if rest, ok := strings.CutPrefix(base, "*"); ok {
		pointer, base = "*", rest
	}
	if strings.HasPrefix(base, "[]") {
		return goType, pointer + "slice"
	}
	return goType, pointer + "struct"
}
