package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func Test_selectsForSchemaShow(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	base := ssot.Object{
		Name:       "Widget",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
		Generators: []ssot.GeneratorOption{ssot.WithSchema(ssot.SchemaShow())},
	}
	variant := func(opts ...ssot.SchemaShowOption) ssot.Object {
		return ssot.Object{
			Name:       "WidgetGlue",
			VariantOf:  &base,
			Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
			Generators: []ssot.GeneratorOption{ssot.WithSchema(ssot.SchemaShow(opts...))},
		}
	}

	require.True(t, selectsForSchemaShow(base))
	require.True(t, selectsForSchemaShow(variant()))
	require.False(t, selectsForSchemaShow(variant(ssot.SchemaShowFromBase())), "a variant reusing the base schema is not generated")
	require.False(t, selectsForSchemaShow(ssot.Object{Name: "Empty"}), "no SHOW columns")
}

func Test_showResultSchemaDetailsFromSsot_panicsOnNilColumn(t *testing.T) {
	require.Panics(t, func() {
		showResultSchemaDetailsFromSsot(ssot.Object{
			Name: "Widget",
			Show: ssot.OutputSpec{Columns: []*ssot.Attribute{nil}},
		}, false)
	})
}

func Test_showResultSchemaDetailsFromSsot_overridesSetColumnsToTypeSet(t *testing.T) {
	accounts := ssot.Attribute{Name: "allowed_accounts", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountIdentifier"}}
	tags := ssot.Attribute{Name: "tags", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	tables := ssot.Attribute{Name: "tables", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	obj := ssot.Object{
		Name:       "Widget",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&accounts, &tags, &tables}},
		Generators: []ssot.GeneratorOption{ssot.WithSchema(ssot.SchemaShow(ssot.SchemaShowListFields(&tables)))},
	}

	got := showResultSchemaDetailsFromSsot(obj, false)
	require.Equal(t, map[string]schema.ValueType{"allowed_accounts": schema.TypeSet}, got.TypeOverrides)
}
