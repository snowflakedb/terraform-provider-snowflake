package gen

import (
	"strings"
	"testing"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ShowOutputAssertionsFromSsot_usesShowAndDescribeColumns(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	errorCode := ssot.Attribute{Name: "error_code", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithShowOutputAssert(ssot.ShowOutputShow(), ssot.ShowOutputDescribe()),
		},
		Show:     ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
		Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{&errorCode}},
	}

	showDetails := showOutputDetailsFromSsot(obj, false)
	require.Equal(t, SdkObjectShowOutputDetails{
		dataSourceDef: &dataSourceDef{},
		SdkObjectDetails: genhelpers.SdkObjectDetails{
			StructDetails: genhelpers.StructDetails{
				Name:   "sdk.Widget",
				Fields: []genhelpers.Field{ssot.ColumnField(obj, comment)},
			},
		},
	}, showDetails)
	show := ShowOutputAssertionsFromSsot(obj, false)
	require.Equal(t, "Widget", show.Name)
	require.False(t, show.IsDescribeOutput)
	require.Empty(t, show.DataSourceName)
	require.Equal(t, []string{"Comment"}, showOutputAssertionNames(show))

	descDetails := showOutputDetailsFromSsot(obj, true)
	require.Equal(t, SdkObjectShowOutputDetails{
		dataSourceDef: &dataSourceDef{},
		SdkObjectDetails: genhelpers.SdkObjectDetails{
			IsDataSourceOutput: true,
			StructDetails: genhelpers.StructDetails{
				Name:   "sdk.WidgetDetails",
				Fields: []genhelpers.Field{ssot.ColumnField(obj, errorCode)},
			},
		},
	}, descDetails)
	desc := ShowOutputAssertionsFromSsot(obj, true)
	require.Equal(t, "WidgetDetails", desc.Name)
	require.True(t, desc.IsDescribeOutput)
	require.Equal(t, []string{"ErrorCode"}, showOutputAssertionNames(desc))
}

func Test_showOutputDetailsFromSsot_passesSkipFields(t *testing.T) {
	tables := ssot.Attribute{Name: "tables", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithShowOutputAssert(ssot.ShowOutputShow(ssot.ShowOutputShowSkipFields(&tables))),
		},
		Show: ssot.OutputSpec{Columns: []*ssot.Attribute{&tables}},
	}

	show := showOutputDetailsFromSsot(obj, false)
	require.Equal(t, []string{"Tables"}, show.SkipFields)
	model := ShowOutputAssertionsFromSsot(obj, false)
	require.True(t, model.Attributes[0].SkipGeneration)
}

func Test_ShowOutputAssertionsFromSsot_refusesDisabledOrEmpty(t *testing.T) {
	column := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:   "Widget",
			IdType: "AccountObjectIdentifier",
			Show:   ssot.OutputSpec{Columns: []*ssot.Attribute{&column}},
		}, false)
	})
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:       "Widget",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithShowOutputAssert(ssot.ShowOutputShow())},
		}, false)
	})
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:       "Widget",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithShowOutputAssert(ssot.ShowOutputShow())},
			Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{nil}},
		}, false)
	})
	require.Panics(t, func() {
		ShowOutputAssertionsFromSsot(ssot.Object{
			Name:       "Widget",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithSchema(ssot.SchemaShow())},
			Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&column}},
		}, false)
	})
}

func Test_MapToResourceShowOutputAssertion_reflectedField(t *testing.T) {
	field := genhelpers.Field{Name: "WarehouseSize", ConcreteType: "string", UnderlyingType: "string"}
	model := ResourceShowOutputAssertionsModel{
		Name:       "Warehouse",
		Attributes: []ResourceShowOutputAssertionModel{MapToResourceShowOutputAssertion(field)},
	}
	assertions := renderShowOutputAssertions(t, model)
	require.Contains(t, assertions, "func (w *WarehouseShowOutputAssert) HasWarehouseSize(expected string) *WarehouseShowOutputAssert {")
	require.Contains(t, assertions, `w.StringValueSet("warehouse_size", expected)`)
	require.NotContains(t, assertions, "HasWarehouseSizeSdk")
	require.NotContains(t, assertions, "HasWarehouseSizeBool")
}

func showOutputAssertionNames(model ResourceShowOutputAssertionsModel) []string {
	names := make([]string, len(model.Attributes))
	for i, attr := range model.Attributes {
		names[i] = attr.Name
	}
	return names
}

func renderShowOutputAssertions(t *testing.T, model ResourceShowOutputAssertionsModel) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, AssertionsTemplate.Execute(&buf, model))
	return buf.String()
}

func Test_ShowOutputAssertionsFromSsot_showRowFromBase(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	base := ssot.Object{
		Name:       "Widget",
		IdType:     "AccountObjectIdentifier",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
		Generators: []ssot.GeneratorOption{ssot.WithShowOutputAssert(ssot.ShowOutputShow())},
	}
	variant := func(opts ...ssot.ShowOutputShowOption) ssot.Object {
		return ssot.Object{
			Name:       "WidgetGlue",
			IdType:     "AccountObjectIdentifier",
			VariantOf:  &base,
			Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
			Generators: []ssot.GeneratorOption{ssot.WithShowOutputAssert(ssot.ShowOutputShow(opts...))},
		}
	}
	render := func(t *testing.T, tmpl *template.Template, model ResourceShowOutputAssertionsModel) string {
		t.Helper()
		var buf strings.Builder
		require.NoError(t, tmpl.Execute(&buf, model))
		return buf.String()
	}

	t.Run("a reusing variant gets only constructors returning the base type", func(t *testing.T) {
		model := ShowOutputAssertionsFromSsot(variant(ssot.ShowOutputShowRowFromBase()), false)
		require.Equal(t, "Widget", model.ReusedBase)
		definition := render(t, DefinitionTemplate, model)
		require.Contains(t, definition, "func WidgetGlueShowOutput(t *testing.T, name string) *WidgetShowOutputAssert {")
		require.Contains(t, definition, "return WidgetShowOutput(t, name)")
		require.Contains(t, definition, "func ImportedWidgetGlueShowOutput(t *testing.T, id string) *WidgetShowOutputAssert {")
		require.Contains(t, definition, "return ImportedWidgetShowOutput(t, id)")
		require.NotContains(t, definition, "type WidgetGlueShowOutputAssert")
		require.Empty(t, strings.TrimSpace(render(t, AssertionsTemplate, model)), "no attribute methods on a reusing variant")
	})
	t.Run("a variant that does not reuse keeps its own assertion", func(t *testing.T) {
		model := ShowOutputAssertionsFromSsot(variant(), false)
		require.Empty(t, model.ReusedBase)
		require.Contains(t, render(t, DefinitionTemplate, model), "type WidgetGlueShowOutputAssert struct")
		require.Contains(t, render(t, AssertionsTemplate, model), "HasComment")
	})
	t.Run("describe output of a reusing variant is still generated", func(t *testing.T) {
		obj := variant(ssot.ShowOutputShowRowFromBase())
		obj.Describe = ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}}
		obj.Generators = []ssot.GeneratorOption{ssot.WithShowOutputAssert(ssot.ShowOutputShow(ssot.ShowOutputShowRowFromBase()), ssot.ShowOutputDescribe())}
		model := ShowOutputAssertionsFromSsot(obj, true)
		require.Empty(t, model.ReusedBase)
		require.Contains(t, render(t, DefinitionTemplate, model), "type WidgetGlueDescribeOutputAssert struct")
	})
}
