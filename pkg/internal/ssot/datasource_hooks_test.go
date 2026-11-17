package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func dataSourceTestObject(opts ...GeneratorOption) Object {
	return Object{
		Name:       "Thing",
		NamePlural: "Things",
		Show:       OutputSpec{Columns: []*Attribute{{Name: "name", Readonly: true}}},
		Generators: append([]GeneratorOption{WithSdk(SdkShow())}, opts...),
	}
}

func Test_EffectiveDataSourceHooks_defaults(t *testing.T) {
	t.Run("an object with only SHOW exposes show output", func(t *testing.T) {
		hooks := dataSourceTestObject().EffectiveDataSourceHooks()
		require.Equal(t, Generated(), hooks.ShowInSdk)
		require.Equal(t, Generated(), hooks.ShowOutputSchema)
		require.Equal(t, Generated(), hooks.SetId)
		require.Equal(t, Omitted(), hooks.AdditionalSchema)
		require.Equal(t, Omitted(), hooks.ItemExtra)
		require.Equal(t, Omitted(), hooks.DescribeInSdk)
		require.Equal(t, Omitted(), hooks.DescribeOutputSchema)
		require.Equal(t, Omitted(), hooks.DescribeOutputSet)
		require.Equal(t, Omitted(), hooks.ShowParametersInSdk)
		require.Equal(t, Omitted(), hooks.ParametersOutputSchema)
		require.Equal(t, Omitted(), hooks.ParametersOutputSet)
		require.Equal(t, Omitted(), hooks.SwitchesSchema, "nothing to switch")
	})
	t.Run("describe is exposed when the object has it", func(t *testing.T) {
		obj := dataSourceTestObject(WithSdk(SdkDescribe()))
		obj.Describe = OutputSpec{Columns: []*Attribute{{Name: "created_on", Readonly: true}}}
		hooks := obj.EffectiveDataSourceHooks()
		require.Equal(t, Generated(), hooks.DescribeInSdk)
		require.Equal(t, Generated(), hooks.DescribeOutputSchema)
		require.Equal(t, Generated(), hooks.DescribeOutputSet)
		require.Equal(t, Generated(), hooks.SwitchesSchema)
		require.Equal(t, Omitted(), hooks.ShowParametersInSdk)
	})
	t.Run("SDK DESCRIBE without columns is not exposed", func(t *testing.T) {
		hooks := dataSourceTestObject(WithSdk(SdkDescribe())).EffectiveDataSourceHooks()
		require.Equal(t, Omitted(), hooks.DescribeInSdk)
	})
	t.Run("parameters are exposed when the SDK has them", func(t *testing.T) {
		hooks := dataSourceTestObject(WithSdk(SdkShowParameters())).EffectiveDataSourceHooks()
		require.Equal(t, Generated(), hooks.ShowParametersInSdk)
		require.Equal(t, Generated(), hooks.ParametersOutputSchema)
		require.Equal(t, Generated(), hooks.ParametersOutputSet)
		require.Equal(t, Generated(), hooks.SwitchesSchema)
		require.Equal(t, Omitted(), hooks.DescribeInSdk)
	})
}

func Test_EffectiveDataSourceHooks_optOut(t *testing.T) {
	describing := func(opts ...DataSourceHookOption) Object {
		obj := dataSourceTestObject(WithSdk(SdkDescribe()), WithSdk(SdkShowParameters()))
		obj.Describe = OutputSpec{Columns: []*Attribute{{Name: "created_on", Readonly: true}}}
		obj.DataSourceHooks = opts
		return obj
	}

	t.Run("omitting the describe call drops its schema and set", func(t *testing.T) {
		hooks := describing(WithDataSourceDescribeInSdk(Omitted())).EffectiveDataSourceHooks()
		require.Equal(t, Omitted(), hooks.DescribeInSdk)
		require.Equal(t, Omitted(), hooks.DescribeOutputSchema)
		require.Equal(t, Omitted(), hooks.DescribeOutputSet)
		require.Equal(t, Generated(), hooks.ShowParametersInSdk)
		require.Equal(t, Generated(), hooks.SwitchesSchema, "with_parameters is left")
	})
	t.Run("omitting both calls drops the switches", func(t *testing.T) {
		hooks := describing(WithDataSourceDescribeInSdk(Omitted()), WithDataSourceShowParametersInSdk(Omitted())).EffectiveDataSourceHooks()
		require.Equal(t, Omitted(), hooks.SwitchesSchema)
		require.Equal(t, Omitted(), hooks.ParametersOutputSet)
	})
	t.Run("a hand-written describe call keeps the generated schema and set", func(t *testing.T) {
		hooks := describing(WithDataSourceDescribeInSdk(ExtCall())).EffectiveDataSourceHooks()
		require.Equal(t, ExtCall(), hooks.DescribeInSdk)
		require.Equal(t, Generated(), hooks.DescribeOutputSchema)
		require.Equal(t, Generated(), hooks.DescribeOutputSet)
	})
	t.Run("an explicit ext call stays", func(t *testing.T) {
		hooks := describing(
			WithDataSourceDescribeInSdk(Omitted()),
			WithDataSourceDescribeOutputSet(ExtCall()),
			WithDataSourceSetId(ExtCall()),
		).EffectiveDataSourceHooks()
		require.Equal(t, ExtCall(), hooks.DescribeOutputSet)
		require.Equal(t, ExtCall(), hooks.SetId)
	})
	t.Run("additions", func(t *testing.T) {
		hooks := dataSourceTestObject().EffectiveDataSourceHooks()
		require.Equal(t, Omitted(), hooks.ItemExtra)
		obj := dataSourceTestObject()
		obj.DataSourceHooks = []DataSourceHookOption{WithDataSourceItemExtra(ExtCall()), WithDataSourceAdditionalSchema(ExtCall())}
		hooks = obj.EffectiveDataSourceHooks()
		require.Equal(t, ExtCall(), hooks.ItemExtra)
		require.Equal(t, ExtCall(), hooks.AdditionalSchema)
	})
}

func Test_Object_DataSourceFilters(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		require.Equal(t, DataSourceFilters{}, dataSourceTestObject().DataSourceFilters())
	})
	t.Run("follow the SDK show options", func(t *testing.T) {
		obj := dataSourceTestObject(WithSdk(SdkShow(
			SdkShowOptionalLike(),
			SdkShowOptionalStartsWith(),
			SdkShowOptionalLimitFrom(),
			SdkShowTerse(),
			SdkShowOptionalSql("HISTORY"),
			SdkShowOptionalExtendedIn(),
		)), WithDataSource(DataSourceWithTerse(), DataSourceWithHistory()))
		require.Equal(t, DataSourceFilters{
			Like: true, StartsWith: true, Limit: true, Terse: true, History: true, In: DataSourceInExtended,
		}, obj.DataSourceFilters())
	})
	t.Run("terse is left out unless opted in", func(t *testing.T) {
		obj := dataSourceTestObject(WithSdk(SdkShow(SdkShowTerse(), SdkShowOptionalLike())))
		require.Equal(t, DataSourceFilters{Like: true}, obj.DataSourceFilters())
		obj = dataSourceTestObject(WithSdk(SdkShow(SdkShowTerse(), SdkShowOptionalLike())), WithDataSource(DataSourceWithTerse()))
		require.Equal(t, DataSourceFilters{Like: true, Terse: true}, obj.DataSourceFilters())
	})
	t.Run("history is left out unless opted in", func(t *testing.T) {
		show := SdkShow(SdkShowOptionalSql("HISTORY"), SdkShowOptionalLike())
		require.Equal(t, DataSourceFilters{Like: true}, dataSourceTestObject(WithSdk(show)).DataSourceFilters())
		require.Equal(t, DataSourceFilters{Like: true, History: true}, dataSourceTestObject(WithSdk(show), WithDataSource(DataSourceWithHistory())).DataSourceFilters())
	})
	t.Run("an IN without a schema level", func(t *testing.T) {
		filters := dataSourceTestObject(WithSdk(SdkShow(SdkShowOptionalExtendedInWithoutSchema()))).DataSourceFilters()
		require.Equal(t, DataSourceInExtendedWithoutSchema, filters.In)
	})
	t.Run("starts with and limit can be left out", func(t *testing.T) {
		show := SdkShow(SdkShowOptionalLike(), SdkShowOptionalStartsWith(), SdkShowOptionalLimitFrom())
		obj := dataSourceTestObject(WithSdk(show), WithDataSource(DataSourceWithoutFilters(DataSourceFilterStartsWith)))
		require.Equal(t, DataSourceFilters{Like: true, Limit: true}, obj.DataSourceFilters())
		obj = dataSourceTestObject(WithSdk(show), WithDataSource(DataSourceWithoutFilters(DataSourceFilterStartsWith, DataSourceFilterLimit)))
		require.Equal(t, DataSourceFilters{Like: true}, obj.DataSourceFilters())
	})
	t.Run("in flavors", func(t *testing.T) {
		require.Equal(t, DataSourceInPlain, dataSourceTestObject(WithSdk(SdkShow(SdkShowOptionalIn()))).DataSourceFilters().In)
		require.Equal(t, DataSourceInAccount, dataSourceTestObject(WithSdk(SdkShow(SdkShowOptionalInAccount()))).DataSourceFilters().In)
	})
	t.Run("on and the other limit option", func(t *testing.T) {
		filters := dataSourceTestObject(WithSdk(SdkShow(SdkShowOptionalOn(), SdkShowOptionalLimit()))).DataSourceFilters()
		require.True(t, filters.On)
		require.True(t, filters.Limit)
	})
}

func Test_Object_DataSourceDescribe(t *testing.T) {
	require.Equal(t, DataSourceDescribe{Method: "Describe"}, dataSourceTestObject(WithSdk(SdkDescribe())).DataSourceDescribe())
	require.Equal(t,
		DataSourceDescribe{Method: "Describe", Slice: true},
		dataSourceTestObject(WithSdk(SdkDescribe(SdkDescribeMapping(SdkDescriptionMappingSlice)))).DataSourceDescribe(),
	)
	require.Equal(t,
		DataSourceDescribe{Method: "DescribeDetails"},
		dataSourceTestObject(WithSdk(SdkDescribe(SdkDescribeDetailsFromIdentifier()))).DataSourceDescribe(),
		"DescribeDetails wins when the SDK emits it",
	)
	require.Equal(t,
		DataSourceDescribe{Method: "DescribeDetails"},
		dataSourceTestObject(WithSdk(SdkDescribe(SdkDescribeMapping(SdkDescriptionMappingSlice), SdkDescribeDetails(SdkDescribeDetailsMethod())))).DataSourceDescribe(),
		"details is one result even when the rows are a list",
	)
}

func Test_ValidateObject_dataSource(t *testing.T) {
	withHooks := func(obj Object, opts ...DataSourceHookOption) Object {
		obj.DataSourceHooks = opts
		return obj
	}

	t.Run("a valid data source", func(t *testing.T) {
		require.NotPanics(t, func() { ValidateObject(dataSourceTestObject(WithDataSource())) })
	})
	t.Run("an enabled variant is skipped, not an error", func(t *testing.T) {
		base := dataSourceTestObject()
		variant := dataSourceTestObject(WithDataSource())
		variant.Name = "ThingSecondary"
		variant.VariantOf = &base
		require.NotPanics(t, func() { ValidateObject(variant) })
	})
	t.Run("hooks need WithDataSource", func(t *testing.T) {
		require.PanicsWithValue(t,
			`object "Thing": DataSourceHooks need WithDataSource()`,
			func() { ValidateObject(withHooks(dataSourceTestObject(), WithDataSourceItemExtra(ExtCall()))) },
		)
	})
	t.Run("a variant carries no hooks", func(t *testing.T) {
		base := dataSourceTestObject()
		variant := withHooks(dataSourceTestObject(WithDataSource()), WithDataSourceItemExtra(ExtCall()))
		variant.Name = "ThingSecondary"
		variant.VariantOf = &base
		require.Panics(t, func() { ValidateObject(variant) })
	})
	t.Run("generated describe needs DESCRIBE", func(t *testing.T) {
		require.Panics(t, func() {
			ValidateObject(withHooks(dataSourceTestObject(WithDataSource()), WithDataSourceDescribeInSdk(Generated())))
		})
	})
	t.Run("a generated describe schema needs DESCRIBE columns", func(t *testing.T) {
		require.Panics(t, func() {
			ValidateObject(withHooks(dataSourceTestObject(WithDataSource()), WithDataSourceDescribeInSdk(ExtCall())))
		})
	})
	t.Run("generated parameters need SHOW PARAMETERS", func(t *testing.T) {
		require.Panics(t, func() {
			ValidateObject(withHooks(dataSourceTestObject(WithDataSource()), WithDataSourceShowParametersInSdk(Generated())))
		})
	})
	t.Run("a hand-written describe with no generated schema is allowed", func(t *testing.T) {
		require.NotPanics(t, func() {
			ValidateObject(withHooks(
				dataSourceTestObject(WithDataSource()),
				WithDataSourceDescribeInSdk(ExtCall()),
				WithDataSourceDescribeOutputSchema(ExtCall()),
				WithDataSourceDescribeOutputSet(ExtCall()),
			))
		})
	})

	t.Run("a property-list describe needs the details method", func(t *testing.T) {
		propertyList := func(sdk ...SdkOption) Object {
			obj := dataSourceTestObject(WithDataSource(), WithSdk(sdk...))
			obj.DescribeRows = OutputSpec{Columns: []*Attribute{{Name: "property", Readonly: true}}}
			obj.Describe = OutputSpec{Columns: []*Attribute{{Name: "created_on", Readonly: true}}}
			return obj
		}
		require.Panics(t, func() { ValidateObject(propertyList(SdkDescribe())) })
		require.NotPanics(t, func() {
			ValidateObject(propertyList(SdkDescribe(SdkDescribeDetails(SdkDescribeDetailsMethod()))))
		})
		hooks := propertyList(SdkDescribe())
		hooks.DataSourceHooks = []DataSourceHookOption{
			WithDataSourceDescribeInSdk(ExtCall()),
			WithDataSourceDescribeOutputSet(ExtCall()),
		}
		require.NotPanics(t, func() { ValidateObject(hooks) })
	})
}
