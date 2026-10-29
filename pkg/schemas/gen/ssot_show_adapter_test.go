package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_showResultSchemaDetailsFromSsot_usesShowAndDescribeColumns(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	errorCode := ssot.Attribute{Name: "error_code", Type: ssot.AttrString}
	obj := ssot.Object{
		Name: "Widget",
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{&comment},
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{&errorCode},
		},
	}

	require.Equal(t, ShowResultSchemaDetails{
		StructDetails: genhelpers.StructDetails{
			Name:   "sdk.Widget",
			Fields: []genhelpers.Field{ssot.ColumnField(obj, comment)},
		},
	}, showResultSchemaDetailsFromSsot(obj, false))
	require.Equal(t, ShowResultSchemaDetails{
		IsDescribe: true,
		StructDetails: genhelpers.StructDetails{
			Name:   "sdk.WidgetDetails",
			Fields: []genhelpers.Field{ssot.ColumnField(obj, errorCode)},
		},
	}, showResultSchemaDetailsFromSsot(obj, true))
}

func Test_showResultSchemaDetailsFromSsot_collectsSchemaSkips(t *testing.T) {
	actives := ssot.Attribute{Name: "actives", Type: ssot.AttrString, SkipInSchema: true}
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{Name: "Widget", Show: ssot.OutputSpec{Columns: []*ssot.Attribute{&name, &actives}}}

	got := showResultSchemaDetailsFromSsot(obj, false)
	require.Equal(t, []string{"actives"}, got.SkipFields)
	require.Equal(t, []genhelpers.Field{ssot.ColumnField(obj, name), ssot.ColumnField(obj, actives)}, got.Fields)
}

func Test_showResultSchemaDetailsFromSsot_showFromBase(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	base := ssot.Object{Name: "Widget"}
	obj := ssot.Object{
		Name:      "WidgetGlue",
		VariantOf: &base,
		Generators: []ssot.GeneratorOption{
			ssot.WithSchema(ssot.SchemaShow(ssot.SchemaShowFromBase())),
		},
		Show: ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
	}

	got := showResultSchemaDetailsFromSsot(obj, false)
	require.Equal(t, "sdk.WidgetGlue", got.Name)
	require.Equal(t, "sdk.Widget", got.SdkType)
	model := ModelFromStructDetails(got, nil, nil)
	require.Equal(t, "WidgetGlue", model.Name)
	require.Equal(t, "sdk.Widget", model.SdkType)
}

func Test_showResultSchemaDetailsFromSsot_panicsOnNilColumn(t *testing.T) {
	require.Panics(t, func() {
		showResultSchemaDetailsFromSsot(ssot.Object{
			Name: "Widget",
			Show: ssot.OutputSpec{Columns: []*ssot.Attribute{nil}},
		}, false)
	})
}
