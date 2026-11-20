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

	g = ObjectGenerators(WithSchema(SchemaUsedAsListEntry()))
	require.Equal(t, Generators{Schema: SchemaGenerator{UsedAsListEntry: true}}, g)

	g = ObjectGenerators(WithObjectAssert(
		ObjectAssertShow(
			ObjectAssertNoShowById(),
			ObjectAssertShowNoIdentifiableObject(),
			ObjectAssertShowByParentId("sdk.AccountObjectIdentifier", "User", "ShowProgrammaticAccessToken"),
			ObjectAssertShowFromObjectIDExpr("token.ID().Name()"),
			ObjectAssertShowObjectTypeName("Account"),
		),
		ObjectAssertDescribe(ObjectAssertDescribeFromObjectIDExpr("userDetails.Name")),
	))
	require.Equal(t, Generators{ObjectAssert: ObjectAssertGenerator{
		Show:                     true,
		Describe:                 true,
		NoShowById:               true,
		ShowNoIdentifiableObject: true,
		ShowByParentId: &ObjectAssertParentId{
			ParentIdType:   "sdk.AccountObjectIdentifier",
			ClientName:     "User",
			ShowMethodName: "ShowProgrammaticAccessToken",
		},
		ShowFromObjectIDExpr:     "token.ID().Name()",
		DescribeFromObjectIDExpr: "userDetails.Name",
		ShowObjectTypeName:       "Account",
	}}, g)
	require.Panics(t, func() {
		ObjectGenerators(WithObjectAssert(ObjectAssertShow(ObjectAssertShowByParentId("", "User", "Show"))))
	})
	require.Panics(t, func() {
		ObjectGenerators(WithObjectAssert(ObjectAssertShow(ObjectAssertShowFromObjectIDExpr(""))))
	})
	require.Panics(t, func() {
		ObjectGenerators(WithObjectAssert(ObjectAssertDescribe(ObjectAssertDescribeFromObjectIDExpr(""))))
	})
	require.Panics(t, func() {
		ObjectGenerators(WithObjectAssert(ObjectAssertShow(ObjectAssertShowObjectTypeName(""))))
	})
}

func Test_ShowFromBase(t *testing.T) {
	base := Object{Name: "Widget"}
	obj := Object{
		Name:    "WidgetAdaptive",
		Variant: &ObjectVariant{Of: &base, ObjectAssertShowFromBase: "ShowByIDAdaptive"},
		Generators: []GeneratorOption{
			WithObjectAssert(ObjectAssertShow()),
		},
	}
	got, ok := obj.ShowFromBase()
	require.True(t, ok)
	require.Equal(t, ShowFromBase{ClientName: "Widget", Method: "ShowByIDAdaptive"}, got)

	_, ok = Object{Name: "Widget"}.ShowFromBase()
	require.False(t, ok)

	require.Panics(t, func() {
		Object{
			Name:    "WidgetAdaptive",
			Variant: &ObjectVariant{ObjectAssertShowFromBase: "ShowByIDAdaptive"},
		}.ShowFromBase()
	})
}

func Test_ParameterAssertionNames(t *testing.T) {
	require.Equal(t, ParameterAssertionNames{}, Object{Name: "Schema"}.ParameterAssertionNames())

	base := Object{Name: "Warehouse"}
	got := Object{Name: "WarehouseAdaptive", Variant: &ObjectVariant{Of: &base}}.ParameterAssertionNames()
	require.Equal(t, ParameterAssertionNames{ConstantPrefix: "Warehouse", ObjectTypeName: "Warehouse"}, got)

	require.Panics(t, func() {
		Object{Name: "WarehouseAdaptive", Variant: &ObjectVariant{Of: &Object{}}}.ParameterAssertionNames()
	})
}

func Test_ObjectGenerators_sdk(t *testing.T) {
	scopes := Attribute{Name: "scopes", Type: AttrList}
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
		SdkShowOptionalExtendedInWithoutSchema(),
		SdkShowOptionalOn(),
		SdkShowOptionalStartsWith(),
		SdkShowOptionalLimitFrom(),
		SdkShowOptionalLimit(),
		SdkShowByIDFiltering(SdkShowByIDExtendedIn, SdkShowByIDLike),
		SdkShowByIDPredicate(SdkShowByIDPredicateNameAndLocator),
		SdkShowResultFilterHook(),
		SdkShowPlainOnly(&scopes, "[]Scope"),
	), SdkDescribe(
		SdkDescribeDoc("https://desc"),
		SdkDescribeSql("THING"),
		SdkDescribeMapping(SdkDescriptionMappingSlice),
		SdkDescribeQueryStructName("DescThing"),
	)))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Show:                                true,
		Describe:                            true,
		ShowDoc:                             "https://show",
		ShowSql:                             "THINGS",
		ShowTerse:                           true,
		ShowOptionalSql:                     []string{"HISTORY", "FUTURE"},
		ShowOptionalLike:                    true,
		ShowOptionalExtendedIn:              true,
		ShowOptionalExtendedInWithoutSchema: true,
		ShowOptionalOn:                      true,
		ShowOptionalStartsWith:              true,
		ShowOptionalLimitFrom:               true,
		ShowOptionalLimit:                   true,
		ShowByIDFiltering:                   []SdkShowByIDFilteringKind{SdkShowByIDExtendedIn, SdkShowByIDLike},
		ShowByIDPredicate:                   SdkShowByIDPredicateNameAndLocator,
		ShowResultFilterHook:                true,
		ShowPlainOnlyFields:                 []SdkShowPlainOnlyField{{Column: &scopes, Kind: "[]Scope"}},
		DescribeDoc:                         "https://desc",
		DescribeSql:                         "THING",
		DescribeMapping:                     SdkDescriptionMappingSlice,
		DescribeQueryStructName:             "DescThing",
	}}, g)

	g = ObjectGenerators(WithSdk(SdkDescribe(
		SdkDescribePlainOnlyField("DatabaseName", "string"),
		SdkDescribePlainOnlyField("SchemaName", "string"),
	)))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Describe: true,
		DescribePlainOnlyFields: []SdkPlainOnlyField{
			{Name: "DatabaseName", Kind: "string"},
			{Name: "SchemaName", Kind: "string"},
		},
	}}, g)

	g = ObjectGenerators(WithSdk(SdkDescribe(SdkDescribeDetailsFromIdentifier())))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Describe:                      true,
		DescribeDetailsFromIdentifier: true,
	}}, g)

	column := Attribute{Name: "column", Type: AttrString}
	g = ObjectGenerators(WithSdk(SdkDescribe(SdkDescribeDetails(
		SdkDescribeDetailsRowsField("descriptions", "Descriptions"),
		SdkDescribeDetailsParseIgnoreUnknown(),
		SdkDescribeDetailsParse(DescribePropertyParsed(&column, "COLUMN")),
	))))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Describe:                          true,
		DescribeDetails:                   true,
		DescribeDetailsRowsField:          &SdkPlainOnlyField{Name: "descriptions", Kind: "Descriptions"},
		DescribeDetailsParseIgnoreUnknown: true,
		DescribeDetailsParse:              []DescribeDetailsMapping{{Attr: &column, Property: "COLUMN", Parsed: true}},
	}}, g)

	g = ObjectGenerators(WithSdk(SdkDrop(
		SdkDropDoc("https://drop"),
		SdkDropSafelyHook(),
		SdkDropSafelyForce(),
		SdkDropCascadeRestrict(),
	)))
	require.Equal(t, Generators{Sdk: SdkGenerator{
		Drop:                true,
		DropDoc:             "https://drop",
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
