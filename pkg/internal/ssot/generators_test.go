package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ObjectGenerators_options(t *testing.T) {
	g := ObjectGenerators(WithSchema(SchemaShow(), SchemaDescribe()))
	require.Equal(t, Generators{Schema: SchemaGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithSchema(SchemaShow()))
	require.Equal(t, Generators{Schema: SchemaGenerator{Show: true}}, g)

	g = ObjectGenerators(WithSchema(SchemaDescribe()))
	require.Equal(t, Generators{Schema: SchemaGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithSchema())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators(WithResource())
	require.Equal(t, Generators{Resource: ResourceGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithModel())
	require.Equal(t, Generators{Model: ModelGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithResourceAssert())
	require.Equal(t, Generators{ResourceAssert: ResourceAssertGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(ShowOutputShow(), ShowOutputDescribe()))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(ShowOutputShow()))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Show: true}}, g)

	g = ObjectGenerators(WithShowOutputAssert(ShowOutputDescribe()))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{Describe: true}}, g)

	tables := Attribute{Name: "tables"}
	g = ObjectGenerators(WithShowOutputAssert(
		ShowOutputShow(ShowOutputShowSkipFields(&tables)),
		ShowOutputDescribe(ShowOutputDescribeSkipFields(&tables)),
	))
	require.Equal(t, Generators{ShowOutputAssert: ShowOutputAssertGenerator{
		Show:               true,
		Describe:           true,
		ShowSkipFields:     []*Attribute{&tables},
		DescribeSkipFields: []*Attribute{&tables},
	}}, g)

	g = ObjectGenerators(WithShowOutputAssert())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertShow(), ObjectAssertDescribe()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertShow()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Show: true}}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertDescribe()))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithObjectAssert(ObjectAssertDescribe(WithNoIdentifiableObject())))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{Describe: true, NoIdentifiableObject: true}}, g)

	password := Attribute{Name: "password"}
	g = ObjectGenerators(WithObjectAssert(
		ObjectAssertShow(ObjectAssertShowSkipFields(&tables)),
		ObjectAssertDescribe(ObjectAssertDescribeSkipFields(&password)),
	))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{
		Show:               true,
		Describe:           true,
		ShowSkipFields:     []*Attribute{&tables},
		DescribeSkipFields: []*Attribute{&password},
	}}, g)

	g = ObjectGenerators(WithObjectAssert())
	require.Equal(t, Generators{}, g)
}

func Test_ShowFromBase(t *testing.T) {
	base := Object{Name: "Widget"}
	obj := Object{
		Name:      "WidgetAdaptive",
		VariantOf: &base,
		Generators: []GeneratorOption{
			WithObjectAssert(ObjectAssertShow(ObjectAssertShowFromBase("ShowByIDAdaptive"))),
		},
	}
	got, ok := obj.ShowFromBase()
	require.True(t, ok)
	require.Equal(t, ShowFromBase{ClientName: "Widget", Method: "ShowByIDAdaptive"}, got)

	_, ok = Object{Name: "Widget"}.ShowFromBase()
	require.False(t, ok)

	require.Panics(t, func() {
		Object{
			Name: "WidgetAdaptive",
			Generators: []GeneratorOption{
				WithObjectAssert(ObjectAssertShow(ObjectAssertShowFromBase("ShowByIDAdaptive"))),
			},
		}.ShowFromBase()
	})
}

func Test_ParameterAssertionNames(t *testing.T) {
	require.Equal(t, ParameterAssertionNames{}, Object{Name: "Schema"}.ParameterAssertionNames())

	base := Object{Name: "Warehouse"}
	got := Object{Name: "WarehouseAdaptive", VariantOf: &base}.ParameterAssertionNames()
	require.Equal(t, ParameterAssertionNames{ConstantPrefix: "Warehouse", ObjectTypeName: "Warehouse"}, got)

	require.Panics(t, func() {
		Object{Name: "WarehouseAdaptive", VariantOf: &Object{}}.ParameterAssertionNames()
	})
}

func Test_ObjectGenerators_sdk(t *testing.T) {
	g := ObjectGenerators(WithSdk(SdkShow(), SdkDescribe()))
	require.Equal(t, Generators{Sdk: SdkGenerator{Show: true, Describe: true}}, g)

	g = ObjectGenerators(WithSdk(SdkShow()))
	require.Equal(t, Generators{Sdk: SdkGenerator{Show: true}}, g)

	g = ObjectGenerators(WithSdk(SdkDescribe()))
	require.Equal(t, Generators{Sdk: SdkGenerator{Describe: true}}, g)

	g = ObjectGenerators(WithSdk(SdkShow(
		SdkShowDoc("https://show"),
		SdkShowSql("THINGS"),
		SdkShowTerse(),
		SdkShowOptionalSql("HISTORY", "FUTURE"),
		SdkShowOptionalLike(),
		SdkShowOptionalExtendedIn(),
		SdkShowOptionalStartsWith(),
		SdkShowOptionalLimitFrom(),
		SdkShowOptionalLimit(),
		SdkShowByIDFiltering(SdkShowByIDExtendedIn, SdkShowByIDLike),
	), SdkDescribe(
		SdkDescribeDoc("https://desc"),
		SdkDescribeSql("THING"),
		SdkDescribeMapping(SdkDescriptionMappingSlice),
		SdkDescribeQueryStructName("DescThing"),
	)))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Show:                    true,
		Describe:                true,
		ShowDoc:                 "https://show",
		ShowSql:                 "THINGS",
		ShowTerse:               true,
		ShowOptionalSql:         []string{"HISTORY", "FUTURE"},
		ShowOptionalLike:        true,
		ShowOptionalExtendedIn:  true,
		ShowOptionalStartsWith:  true,
		ShowOptionalLimitFrom:   true,
		ShowOptionalLimit:       true,
		ShowByIDFiltering:       []SdkShowByIDFilteringKind{SdkShowByIDExtendedIn, SdkShowByIDLike},
		DescribeDoc:             "https://desc",
		DescribeSql:             "THING",
		DescribeMapping:         SdkDescriptionMappingSlice,
		DescribeQueryStructName: "DescThing",
	}}, g)

	g = ObjectGenerators(WithSdk(SdkDrop(
		SdkDropSafelyHook(),
		SdkDropSafelyForce(),
		SdkDropCascadeRestrict(),
	)))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Drop:                true,
		DropSafelyHook:      true,
		DropSafelyForce:     true,
		DropCascadeRestrict: true,
	}}, g)

	g = ObjectGenerators(WithSdk(SdkCustomOperations()))
	require.Equal(t, Generators{Sdk: SdkGenerator{CustomOperations: true}}, g)

	g = ObjectGenerators(WithSdk(SdkCustomInterfaceMethods()))
	require.Equal(t, Generators{Sdk: SdkGenerator{CustomInterfaceMethods: true}}, g)

	g = ObjectGenerators(WithSdk(SdkShowParameters(), SdkShowParametersDetails()))
	require.Equal(t, Generators{Sdk: SdkGenerator{ShowParameters: true, ShowParametersDetails: true}}, g)

	g = ObjectGenerators(WithSdk(SdkDefineFromSsot()))
	require.Equal(t, Generators{Sdk: SdkGenerator{DefineFromSsot: true}}, g)

	g = ObjectGenerators(WithSdk())
	require.Equal(t, Generators{}, g)

	g = ObjectGenerators(WithObjectParametersAssert())
	require.Equal(t, Generators{ObjectParametersAssert: ObjectParametersAssertGenerator{Enabled: true}}, g)

	g = ObjectGenerators(WithResourceParametersAssert())
	require.Equal(t, Generators{ResourceParametersAssert: ResourceParametersAssertGenerator{Enabled: true}}, g)

	g = ObjectGenerators()
	require.Equal(t, Generators{}, g)
}
