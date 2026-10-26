package gen

import (
	"strings"
	"testing"

	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_ModelFromSsotObject_mapsClassifiedAttributes(t *testing.T) {
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithModel(),
		},
		Attributes: []ssot.Attribute{
			{Name: "min_nodes", Type: ssot.AttrInt, Required: true},
			{Name: "enabled", Type: ssot.AttrBool, Required: true},
			{Name: "auto_resume", Type: ssot.AttrBool},
			{Name: "names", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}},
			{Name: "tags", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}},
		},
	}
	model := ModelFromSsotObject(obj)
	require.Equal(t, "Widget", model.Name)
	require.Equal(t, expectOrdered(classifiedNames(ssot.IdentityFields(obj)), classifiedNames(classifyAttributes(obj))), attributeNames(model))
	for _, field := range ssot.IdentityFields(obj) {
		requireModelAttributeMatches(t, modelAttributeByName(t, model, field.Name), field)
	}
	for _, attr := range obj.ResourceAttributes() {
		requireModelAttributeMatches(t, modelAttributeByName(t, model, attr.Name), ssot.ClassifyAttribute(obj, attr))
	}
}

func Test_ModelFromSsotObject_includesIdentityFields(t *testing.T) {
	for _, idType := range []string{"AccountObjectIdentifier", "DatabaseObjectIdentifier", "SchemaObjectIdentifier"} {
		t.Run(idType, func(t *testing.T) {
			obj := ssot.Object{
				Name:       "Widget",
				IdType:     idType,
				Generators: []ssot.GeneratorOption{ssot.WithModel()},
			}
			require.Equal(t, expectOrdered(classifiedNames(ssot.IdentityFields(obj))), attributeNames(ModelFromSsotObject(obj)))
		})
	}
}

func Test_ModelFromSsotObject_joinsParameterCatalog(t *testing.T) {
	// The object name selects the parameter catalog level.
	obj := ssot.Object{
		Name:   "Schema",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithModel(),
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
	}

	model := ModelFromSsotObject(obj)
	require.Equal(t, expectOrdered(identityNames, []string{"note"}, classifiedNames(parameterFields)), attributeNames(model))
	requireModelAttributeMatches(t, modelAttributeByName(t, model, "note"), ssot.ClassifyAttribute(obj, obj.Attributes[0]))
	for _, field := range parameterFields {
		requireModelAttributeMatches(t, modelAttributeByName(t, model, field.Name), field)
	}
}

func Test_ModelFromSsotObject_panics(t *testing.T) {
	t.Run("generator off", func(t *testing.T) {
		require.Panics(t, func() {
			ModelFromSsotObject(ssot.Object{Name: "Widget", IdType: "AccountObjectIdentifier"})
		})
	})
	t.Run("parameter level empty", func(t *testing.T) {
		require.Panics(t, func() {
			ModelFromSsotObject(ssot.Object{
				Name:       "Widget",
				IdType:     "AccountObjectIdentifier",
				Generators: []ssot.GeneratorOption{ssot.WithModel()},
				ResourceHooks: []ssot.ResourceHookOption{
					ssot.WithParametersAttributesSchema(ssot.Generated()),
				},
			})
		})
	})
}

func requireModelAttributeMatches(t *testing.T, got ResourceConfigBuilderAttributeModel, field ssot.ClassifiedField) {
	t.Helper()
	require.Equal(t, field.Name, got.Name)
	require.Equal(t, field.Name, got.JsonName)
	require.Equal(t, field.Required, got.Required)
	require.Equal(t, "tfconfig", got.MethodImport)
	require.Equal(t, field.SdkValueType, got.SdkValueType)
	switch field.Storage {
	case ssot.AttrBool:
		require.Equal(t, "bool", got.AttributeType)
		require.Equal(t, "BoolVariable", got.VariableMethod)
		require.Equal(t, schema.TypeBool, got.OriginalType)
	case ssot.AttrInt:
		require.Equal(t, "int", got.AttributeType)
		require.Equal(t, "IntegerVariable", got.VariableMethod)
		require.Equal(t, schema.TypeInt, got.OriginalType)
	case ssot.AttrString:
		require.Equal(t, "string", got.AttributeType)
		require.Equal(t, "StringVariable", got.VariableMethod)
		require.Equal(t, schema.TypeString, got.OriginalType)
	case ssot.AttrList:
		require.Empty(t, got.AttributeType)
		require.Empty(t, got.VariableMethod)
		require.Equal(t, schema.TypeList, got.OriginalType)
	case ssot.AttrSet:
		require.Empty(t, got.AttributeType)
		require.Empty(t, got.VariableMethod)
		require.Equal(t, schema.TypeSet, got.OriginalType)
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

func attributeNames(model ResourceConfigBuilderModel) []string {
	return collections.Map(model.Attributes, func(attr ResourceConfigBuilderAttributeModel) string { return attr.Name })
}

func modelAttributeByName(t *testing.T, model ResourceConfigBuilderModel, name string) ResourceConfigBuilderAttributeModel {
	t.Helper()
	for _, attr := range model.Attributes {
		if attr.Name == name {
			return attr
		}
	}
	t.Fatalf("attribute %q not found in %s", name, strings.Join(attributeNames(model), ", "))
	return ResourceConfigBuilderAttributeModel{}
}

func Test_ModelFromSsotObject_nativeBoolAndNested(t *testing.T) {
	obj := ssot.Object{
		Name: "Widget", IdType: "AccountObjectIdentifier", Generators: []ssot.GeneratorOption{ssot.WithModel()},
		Attributes: []ssot.Attribute{
			{Name: "enabled", Type: ssot.AttrBool, NativeBool: true},
			{Name: "schedule", Type: ssot.AttrList, Elem: ssot.ObjectElem([]ssot.Attribute{{Name: "interval", Type: ssot.AttrInt}})},
		},
	}
	model := ModelFromSsotObject(obj)
	for _, field := range model.Attributes {
		switch field.Name {
		case "enabled":
			require.Equal(t, "BoolVariable", field.VariableMethod)
		case "schedule":
			require.Empty(t, field.AttributeType)
			require.False(t, field.Required)
		}
	}
}
