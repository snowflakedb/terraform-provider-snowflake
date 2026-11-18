package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ModelFromInputObject(t *testing.T) {
	describedColumn := &ssot.Attribute{Name: "created_on", Readonly: true}
	input := DataSourceObject{Object: ssot.Object{
		Name:           "Thing",
		NamePlural:     "Things",
		PreviewFeature: true,
		Show:           ssot.OutputSpec{Columns: []*ssot.Attribute{{Name: "name", Readonly: true}}},
		Describe:       ssot.OutputSpec{Columns: []*ssot.Attribute{describedColumn}},
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkShow(ssot.SdkShowOptionalLike(), ssot.SdkShowOptionalExtendedIn()), ssot.SdkDescribe()),
		},
		DataSourceHooks: []ssot.DataSourceHookOption{ssot.WithDataSourceSetId(ssot.ExtCall())},
	}}

	model := ModelFromInputObject(input, genhelpers.NewPreambleModel("datasource", "0.1.0"), nil)

	require.Equal(t, "Thing", model.Name)
	require.Equal(t, "Things", model.NamePlural)
	require.Equal(t, "things", model.Prefix())
	require.True(t, model.PreviewFeature)
	require.Equal(t, ssot.DataSourceFilters{Like: true, In: ssot.DataSourceInExtended}, model.Filters)
	require.Equal(t, ssot.DataSourceDescribe{Method: "Describe"}, model.Describe)
	require.Equal(t, ssot.ExtCall(), model.Hooks.SetId)
	require.Equal(t, ssot.Generated(), model.Hooks.DescribeInSdk)
	require.Equal(t, "datasource", model.GeneratorName)
}

func Test_ModelFromInputObject_ExtDescribeWithoutSdkDescribe(t *testing.T) {
	input := DataSourceObject{Object: ssot.Object{
		Name:       "Thing",
		NamePlural: "Things",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{{Name: "name", Readonly: true}}},
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkShow()),
			ssot.WithDataSource(),
		},
		DataSourceHooks: []ssot.DataSourceHookOption{
			ssot.WithDataSourceDescribeInSdk(ssot.ExtCall()),
			ssot.WithDataSourceDescribeOutputSchema(ssot.ExtCall()),
			ssot.WithDataSourceDescribeOutputSet(ssot.ExtCall()),
		},
	}}

	model := ModelFromInputObject(input, genhelpers.NewPreambleModel("datasource", "0.1.0"), nil)

	require.Equal(t, "THING", model.DescribeSql)
	require.Equal(t, ssot.ExtCall(), model.Hooks.DescribeInSdk)
	require.Equal(t, ssot.Generated(), model.Hooks.SwitchesSchema)
	require.Equal(t, "Holds the aggregated output of all thing details queries.", model.ResultDescription())
}

func Test_ModelFromInputObject_CustomDescribeMethod(t *testing.T) {
	input := DataSourceObject{Object: ssot.Object{
		Name:       "Thing",
		NamePlural: "Things",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{{Name: "name", Readonly: true}}},
		CustomOperations: []ssot.CustomOperationSpec{{
			Kind:   "DescribeColumns",
			Paired: &ssot.CustomPairedShow{Plain: "ThingColumnDetails"},
		}},
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkShow()),
			ssot.WithDataSource(ssot.DataSourceDescribeMethod("DescribeColumns")),
		},
	}}

	model := ModelFromInputObject(input, genhelpers.NewPreambleModel("datasource", "0.1.0"), nil)

	require.Equal(t, ssot.DataSourceDescribe{
		Method:             "DescribeColumns",
		Slice:              true,
		RequestConstructor: "NewDescribeColumnsThingRequest",
		ResultName:         "ThingColumnDetails",
	}, model.Describe)
	require.Equal(t, "ThingColumnDetails", model.DescribeSchemaName())
	require.Equal(t, ssot.Generated(), model.Hooks.DescribeInSdk)
	require.Equal(t, "THING", model.DescribeSql)
}
