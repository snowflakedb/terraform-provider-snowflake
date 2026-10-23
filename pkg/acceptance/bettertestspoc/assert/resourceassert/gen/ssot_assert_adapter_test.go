package gen

import (
	"strings"
	"testing"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_AssertionsFromSsotObject_mapsClassifiedAttributes(t *testing.T) {
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithResourceAssert(),
		},
		Attributes: []ssot.Attribute{
			{Name: "min_nodes", Type: ssot.AttrInt, Required: true},
			{Name: "enabled", Type: ssot.AttrBool, Required: true},
			{Name: "auto_resume", Type: ssot.AttrBool},
			{Name: "names", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}},
			{Name: "tags", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}},
		},
	}
	model := AssertionsFromSsotObject(obj)
	require.Equal(t, "Widget", model.Name)
	require.Equal(t, expectOrdered(
		classifiedNames(ssot.IdentityFields(obj)),
		classifiedNames(classifyAttributes(obj)),
		[]string{resources.FullyQualifiedNameAttributeName},
	), attributeNames(model))
	for _, field := range ssot.IdentityFields(obj) {
		requireAssertionMatches(t, assertionByName(t, model, field.Name), field)
	}
	for _, attr := range obj.ResourceAttributes() {
		requireAssertionMatches(t, assertionByName(t, model, attr.Name), ssot.ClassifyAttribute(obj, attr))
	}
	requireAssertionMatches(t, assertionByName(t, model, resources.FullyQualifiedNameAttributeName), ssot.ClassifiedField{
		Name:    resources.FullyQualifiedNameAttributeName,
		Storage: ssot.AttrString,
	})
}

func Test_AssertionsFromSsotObject_includesIdentityFields(t *testing.T) {
	for _, idType := range []string{"AccountObjectIdentifier", "DatabaseObjectIdentifier", "SchemaObjectIdentifier"} {
		t.Run(idType, func(t *testing.T) {
			obj := ssot.Object{
				Name:       "Widget",
				IdType:     idType,
				Generators: []ssot.GeneratorOption{ssot.WithResourceAssert()},
			}
			require.Equal(t, expectOrdered(
				classifiedNames(ssot.IdentityFields(obj)),
				[]string{resources.FullyQualifiedNameAttributeName},
			), attributeNames(AssertionsFromSsotObject(obj)))
		})
	}
}

func Test_ModelFromResourceSchemaDetails_emitsNothingExtra(t *testing.T) {
	model := ModelFromResourceSchemaDetails(genhelpers.ExtractResourceSchemaDetails("Warehouse", map[string]*schema.Schema{
		"warehouse_size": {Type: schema.TypeString, Required: true},
	}), nil, nil)
	assertions := renderAssertionTemplates(t, model)
	require.Contains(t, assertions, "func (w *WarehouseResourceAssert) HasWarehouseSize(expected string) *WarehouseResourceAssert {")
	require.Contains(t, assertions, "func (w *WarehouseResourceAssert) HasWarehouseSizeString(expected string) *WarehouseResourceAssert {")
	require.NotContains(t, assertions, "HasWarehouseSizeSdk")
	require.NotContains(t, assertions, "HasWarehouseSizeBool")
	require.NotContains(t, assertions, "strconv.FormatBool")
}

func Test_AssertionsFromSsotObject_joinsParameterCatalog(t *testing.T) {
	// The object name selects the parameter catalog level.
	obj := ssot.Object{
		Name:   "Schema",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithResourceAssert(),
		},
		ResourceHooks: []ssot.ResourceHookOption{
			ssot.WithParametersAttributesSchema(ssot.Generated()),
		},
		Attributes: []ssot.Attribute{
			{Name: "note", Type: ssot.AttrString},
		},
	}
	parameters := ssotdefs.ParameterDefsForLevel(ssot.ParameterLevelSchema)
	require.NotEmpty(t, parameters)
	parameterFields := collections.Map(parameters, func(parameter ssot.ParameterDef) ssot.ClassifiedField {
		return ssot.ClassifyParameter(parameter.FieldName(), parameter.Kind)
	})
	identityNames := classifiedNames(ssot.IdentityFields(obj))
	for _, field := range parameterFields {
		require.NotContains(t, identityNames, field.Name)
		require.NotEqual(t, "note", field.Name)
		require.NotEqual(t, resources.FullyQualifiedNameAttributeName, field.Name)
	}

	model := AssertionsFromSsotObject(obj)
	require.Equal(t, expectOrdered(
		identityNames,
		[]string{"note"},
		classifiedNames(parameterFields),
		[]string{resources.FullyQualifiedNameAttributeName},
	), attributeNames(model))
	requireAssertionMatches(t, assertionByName(t, model, "note"), ssot.ClassifyAttribute(obj, obj.Attributes[0]))
	for _, field := range parameterFields {
		requireAssertionMatches(t, assertionByName(t, model, field.Name), field)
	}
}

func Test_AssertionsFromSsotObject_panics(t *testing.T) {
	t.Run("generator off", func(t *testing.T) {
		require.Panics(t, func() {
			AssertionsFromSsotObject(ssot.Object{Name: "Widget", IdType: "AccountObjectIdentifier"})
		})
	})
	t.Run("parameter level empty", func(t *testing.T) {
		require.Panics(t, func() {
			AssertionsFromSsotObject(ssot.Object{
				Name:       "Widget",
				IdType:     "AccountObjectIdentifier",
				Generators: []ssot.GeneratorOption{ssot.WithResourceAssert()},
				ResourceHooks: []ssot.ResourceHookOption{
					ssot.WithParametersAttributesSchema(ssot.Generated()),
				},
			})
		})
	})
}

func requireAssertionMatches(t *testing.T, got ResourceAttributeAssertionModel, field ssot.ClassifiedField) {
	t.Helper()
	require.Equal(t, field.Name, got.Name)
	require.Equal(t, field.Required, got.IsRequired)
	require.Equal(t, field.SdkValueType, got.SdkValueType)
	switch field.Storage {
	case ssot.AttrBool:
		require.Equal(t, "Bool", got.Type)
		require.Equal(t, "bool", got.ExpectedType)
		require.Equal(t, "BoolValueSet", got.AssertionCreator)
		require.True(t, got.ShouldGenerateTypedAssertion)
		require.False(t, got.IsCollection)
		require.Empty(t, got.SubType)
	case ssot.AttrInt:
		require.Equal(t, "Int", got.Type)
		require.Equal(t, "int", got.ExpectedType)
		require.Equal(t, "IntValueSet", got.AssertionCreator)
		require.True(t, got.ShouldGenerateTypedAssertion)
		require.False(t, got.IsCollection)
		require.Empty(t, got.SubType)
	case ssot.AttrString:
		require.Equal(t, "String", got.Type)
		require.Equal(t, "string", got.ExpectedType)
		require.Equal(t, "StringValueSet", got.AssertionCreator)
		require.True(t, got.ShouldGenerateTypedAssertion)
		require.False(t, got.IsCollection)
		require.Empty(t, got.SubType)
	case ssot.AttrList:
		require.Equal(t, "List", got.Type)
		require.Equal(t, "String", got.SubType)
		require.Equal(t, "...string", got.ExpectedType)
		require.Equal(t, "ListContainsExactlyStringValuesInOrder", got.AssertionCreator)
		require.True(t, got.ShouldGenerateTypedAssertion)
		require.True(t, got.IsCollection)
	case ssot.AttrSet:
		require.Equal(t, "Set", got.Type)
		require.Equal(t, "String", got.SubType)
		require.Equal(t, "...string", got.ExpectedType)
		require.Equal(t, "SetContainsExactlyStringValues", got.AssertionCreator)
		require.True(t, got.ShouldGenerateTypedAssertion)
		require.True(t, got.IsCollection)
	default:
		t.Fatalf("attribute %q has unexpected storage %q", field.Name, field.Storage)
	}
}

func classifyAttributes(obj ssot.Object) []ssot.ClassifiedField {
	return collections.Map(obj.ResourceAttributes(), func(attr ssot.Attribute) ssot.ClassifiedField {
		return ssot.ClassifyAttribute(obj, attr)
	})
}

func classifiedNames(fields []ssot.ClassifiedField) []string {
	return collections.Map(fields, func(field ssot.ClassifiedField) string { return field.Name })
}

func expectOrdered(groups ...[]string) []string {
	var names []string
	for _, group := range groups {
		names = append(names, group...)
	}
	return genhelpers.OrderByNameWithPreferredFront(names, func(name string) string { return name })
}

func attributeNames(model ResourceAssertionsModel) []string {
	return collections.Map(model.Attributes, func(attr ResourceAttributeAssertionModel) string { return attr.Name })
}

func assertionByName(t *testing.T, model ResourceAssertionsModel, name string) ResourceAttributeAssertionModel {
	t.Helper()
	for _, attr := range model.Attributes {
		if attr.Name == name {
			return attr
		}
	}
	t.Fatalf("attribute %q not found in %s", name, strings.Join(attributeNames(model), ", "))
	return ResourceAttributeAssertionModel{}
}

func renderAssertionTemplates(t *testing.T, model ResourceAssertionsModel) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, DefinitionTemplate.Execute(&buf, model))
	require.NoError(t, AssertionsTemplate.Execute(&buf, model))
	return buf.String()
}
