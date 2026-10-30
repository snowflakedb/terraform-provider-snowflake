package gen

import (
	"io"
	"log"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_ShowOperationFromObject_plainIn(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "SchemaObjectIdentifier",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow(
			ssot.SdkShowOptionalLike(),
			ssot.SdkShowOptionalIn(),
			ssot.SdkShowOptionalStartsWith(),
			ssot.SdkShowOptionalLimit(),
			ssot.SdkShowAdditionalValidations(),
			ssot.SdkShowByIDFiltering(ssot.SdkShowByIDIn, ssot.SdkShowByIDLike),
		))},
	}
	got := NewInterface("Fakes", "Fake", "SchemaObjectIdentifier").ShowOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "SchemaObjectIdentifier").ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
		ShowPairedStructs(obj),
		NewQueryStruct("ShowFakes").
			Show().
			SQL("FAKES").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimit().
			WithAdditionalValidations(),
		ShowByIDInFiltering,
		ShowByIDLikeFiltering,
	)
	requireSameOperations(t, want, got)
}

func Test_ShowOperationFromObject_inAccount(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "AccountObjectIdentifier",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow(
			ssot.SdkShowOptionalInAccount(),
			ssot.SdkShowByIDNoFiltering(),
		))},
	}
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
		ShowPairedStructs(obj),
		NewQueryStruct("ShowFakes").Show().SQL("FAKES").OptionalInAccount(),
		ShowByIDNoFiltering,
	)
	requireSameOperations(t, want, got)
}

func Test_ShowOperationFromObject_noFilteringAndPredicate(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "AccountObjectIdentifier",
		Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow(
			ssot.SdkShowByIDNoFiltering(),
			ssot.SdkShowByIDPredicate(ssot.SdkShowByIDPredicateNameAndLocator),
		))},
	}
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
		ShowOperationWithPairedStructs(
			"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
			ShowPairedStructs(obj),
			NewQueryStruct("ShowFakes").Show().SQL("FAKES"),
			ShowByIDNoFiltering,
		).
		WithShowByIDFindPredicateKind(ShowByIDFindPredicateNameAndLocator)
	requireSameOperations(t, want, got)
	require.Equal(t, ShowByIDFindPredicateNameAndLocator, got.ShowByIDFindPredicateKind)
}

func Test_ShowDescribeFromObject(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	column := ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}
	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "AccountObjectIdentifier",
		Show:       column,
		Describe:   column,
		Generators: []ssot.GeneratorOption{ssot.WithSdk(
			ssot.SdkShow(
				ssot.SdkShowTerse(),
				ssot.SdkShowOptionalSql("HISTORY"),
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalExtendedIn(),
				ssot.SdkShowOptionalOn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimitFrom(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDExtendedIn, ssot.SdkShowByIDLike),
			),
			ssot.SdkDescribe(
				ssot.SdkDescribeMapping(ssot.SdkDescriptionMappingSlice),
				ssot.SdkDescribeQueryStructName("DescribeFake"),
			),
		)},
	}
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
		ShowOperationFromObject(obj).
		DescribeOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
		ShowOperationWithPairedStructs(
			"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
			ShowPairedStructs(obj),
			NewQueryStruct("ShowFakes").
				Show().
				Terse().
				SQL("FAKES").
				OptionalSQL("HISTORY").
				OptionalLike().
				OptionalExtendedIn().
				OptionalOn().
				OptionalStartsWith().
				OptionalLimitFrom().
				OptionalLimit(),
			ShowByIDExtendedInFiltering,
			ShowByIDLikeFiltering,
		).
		DescribeOperationWithPairedStructs(
			DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-fake",
			DescribePairedStructs(obj),
			NewQueryStruct("DescribeFake").
				Describe().
				SQL("FAKE").
				Name().
				WithValidation(ValidIdentifier, "name"),
		)
	requireSameOperations(t, want, got)
}

func Test_DescribeOperationFromObject_identifierOverride(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "AccountObjectIdentifier",
		Describe:   ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(ssot.SdkDescribeIdentifierOverride("ObjectIdentifier")))},
	}
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
		DescribeOperationWithPairedStructs(
			DescriptionMappingKindSingleValue,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-fake",
			DescribePairedStructs(obj),
			NewQueryStruct("DescFake").
				Describe().
				SQL("FAKE").
				Identifier("name", "ObjectIdentifier", IdentifierOptions().Required()).
				WithValidation(ValidIdentifier, "name"),
		)
	requireSameOperations(t, want, got)
}

func Test_DescribeOperationFromObject_plainShowVariants(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	running := ssot.Attribute{Name: "running", Type: ssot.AttrInt, OptionalInShow: true, OptionalInDb: true}
	monitor := ssot.Attribute{
		Name:           "resource_monitor",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		CustomParser:   "ParseAccountObjectIdentifierExcludingExplicitNullString",
	}
	main := ssot.Object{
		Name:          "Widget",
		NamePlural:    "Widgets",
		IdType:        "AccountObjectIdentifier",
		Describe:      ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		ShowPlainName: "WidgetRegular",
		Show:          ssot.OutputSpec{Columns: []*ssot.Attribute{&name, &running}},
		ShowGeneric:   ssot.OutputSpec{Columns: []*ssot.Attribute{&name, &running, &monitor}},
		Generators:    []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())},
	}
	adaptive := ssot.Object{
		Name:      "WidgetAdaptive",
		VariantOf: &main,
		Show:      ssot.OutputSpec{Columns: []*ssot.Attribute{&name, &monitor}},
	}
	interactive := ssot.Object{
		Name:      "WidgetInteractive",
		VariantOf: &main,
		Show:      ssot.OutputSpec{Columns: []*ssot.Attribute{&running}},
	}
	main.ShowVariants = []*ssot.Object{&adaptive, &interactive}

	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(main)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationWithPairedStructs(
		DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-widget",
		DescribePairedStructs(main),
		NewQueryStruct("DescWidget").
			Describe().
			SQL("WIDGET").
			Name().
			WithValidation(ValidIdentifier, "name"),
		PlainStruct("WidgetRegular").Field("Name", "string").Field("Running", "*int"),
		PlainStruct("WidgetAdaptive").Field("Name", "string").Field("ResourceMonitor", "AccountObjectIdentifier"),
		PlainStruct("WidgetInteractive").Field("Running", "*int"),
	)
	requireSameOperations(t, want, got)
	require.Equal(t, []*ShowProjection{
		{Source: "Widget", Method: "AsRegular", Target: "WidgetRegular", Fields: []string{"Name", "Running"}},
		{Source: "Widget", Method: "AsAdaptive", Target: "WidgetAdaptive", Fields: []string{"Name", "ResourceMonitor"}},
		{Source: "Widget", Method: "AsInteractive", Target: "WidgetInteractive", Fields: []string{"Running"}},
	}, got.ShowProjections)
	require.Empty(t, got.ShowByIDProjections)
}

func Test_DescribeOperationFromObject_showByIDProjections(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	main := ssot.Object{
		Name:          "Widget",
		NamePlural:    "Widgets",
		IdType:        "AccountObjectIdentifier",
		Describe:      ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		ShowPlainName: "WidgetRegular",
		Show:          ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
		Generators:    []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShow(), ssot.SdkDescribe())},
	}
	adaptive := ssot.Object{
		Name:      "WidgetAdaptive",
		VariantOf: &main,
		Show:      ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
	}
	main.ShowVariants = []*ssot.Object{&adaptive}

	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").
		ShowOperationFromObject(main).
		DescribeOperationFromObject(main)
	require.Equal(t, []*ShowByIDProjection{
		{Suffix: "Regular", ProjectionMethod: "AsRegular", Target: "WidgetRegular"},
		{Suffix: "Adaptive", ProjectionMethod: "AsAdaptive", Target: "WidgetAdaptive"},
	}, got.ShowByIDProjections)
}

func Test_DescribeOperationFromObject_showProjectionPanics(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	column := ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}
	obj := ssot.Object{
		Name:          "Widget",
		ShowPlainName: "Other",
		Show:          column,
		Describe:      column,
		Generators:    []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())},
	}
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	})
}

func Test_ShowDescribeFromObject_panics(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	column := ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}

	tests := []struct {
		name string
		call func()
	}{
		{
			name: "show option off",
			call: func() {
				NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
					ShowOperationFromObject(ssot.Object{Name: "Fake", NamePlural: "Fakes"})
			},
		},
		{
			name: "show by id no filtering combined",
			call: func() {
				NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
					ShowOperationFromObject(showMethodObject(ssot.SdkShow(
						ssot.SdkShowByIDFiltering(ssot.SdkShowByIDNone, ssot.SdkShowByIDLike),
					), column, ssot.OutputSpec{}))
			},
		},
		{
			name: "unknown show by id predicate",
			call: func() {
				NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
					ShowOperationFromObject(showMethodObject(ssot.SdkShow(
						ssot.SdkShowByIDPredicate("nope"),
					), column, ssot.OutputSpec{}))
			},
		},
		{
			name: "unknown show by id filtering",
			call: func() {
				NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
					ShowOperationFromObject(showMethodObject(ssot.SdkShow(
						ssot.SdkShowDoc("https://show"),
						ssot.SdkShowSql("FAKES"),
						ssot.SdkShowByIDFiltering("nope"),
					), column, ssot.OutputSpec{}))
			},
		},
		{
			name: "describe option off",
			call: func() {
				NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
					DescribeOperationFromObject(ssot.Object{Name: "Fake"})
			},
		},
		{
			name: "unknown description mapping",
			call: func() {
				NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
					DescribeOperationFromObject(showMethodObject(ssot.SdkDescribe(
						ssot.SdkDescribeDoc("https://desc"),
						ssot.SdkDescribeSql("FAKE"),
						ssot.SdkDescribeMapping("rows"),
					), ssot.OutputSpec{}, column))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Panics(t, tt.call)
		})
	}
}

func Test_ShowDescribeFromObject_docURL(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	column := ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}

	t.Run("defaults", func(t *testing.T) {
		show := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
			ShowOperationFromObject(showMethodObject(ssot.SdkShow(ssot.SdkShowSql("FAKES")), column, ssot.OutputSpec{}))
		require.Equal(t, "https://docs.snowflake.com/en/sql-reference/sql/show-fakes", show.Operations[0].Doc)

		describe := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
			DescribeOperationFromObject(showMethodObject(ssot.SdkDescribe(
				ssot.SdkDescribeSql("FAKE"),
			), ssot.OutputSpec{}, column))
		require.Equal(t, "https://docs.snowflake.com/en/sql-reference/sql/desc-fake", describe.Operations[0].Doc)
	})

	t.Run("overrides", func(t *testing.T) {
		show := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
			ShowOperationFromObject(showMethodObject(ssot.SdkShow(
				ssot.SdkShowDoc("https://custom/show"),
				ssot.SdkShowSql("FAKES"),
			), column, ssot.OutputSpec{}))
		require.Equal(t, "https://custom/show", show.Operations[0].Doc)

		describe := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
			DescribeOperationFromObject(showMethodObject(ssot.SdkDescribe(
				ssot.SdkDescribeDoc("https://custom/desc"),
				ssot.SdkDescribeSql("FAKE"),
			), ssot.OutputSpec{}, column))
		require.Equal(t, "https://custom/desc", describe.Operations[0].Doc)
	})
}

func Test_ShowDescribeFromObject_sql(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	column := ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}

	t.Run("defaults", func(t *testing.T) {
		showObj := showMethodObject(ssot.SdkShow(), column, ssot.OutputSpec{})
		gotShow := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationFromObject(showObj)
		wantShow := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationWithPairedStructs(
			"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
			ShowPairedStructs(showObj),
			NewQueryStruct("ShowFakes").Show().SQL("FAKES"),
		)
		requireSameOperations(t, wantShow, gotShow)

		describeObj := showMethodObject(ssot.SdkDescribe(), ssot.OutputSpec{}, column)
		gotDescribe := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DescribeOperationFromObject(describeObj)
		wantDescribe := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DescribeOperationWithPairedStructs(
			DescriptionMappingKindSingleValue,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-fake",
			DescribePairedStructs(describeObj),
			NewQueryStruct("DescFake").Describe().SQL("FAKE").Name().WithValidation(ValidIdentifier, "name"),
		)
		requireSameOperations(t, wantDescribe, gotDescribe)
	})

	t.Run("overrides", func(t *testing.T) {
		showObj := showMethodObject(ssot.SdkShow(ssot.SdkShowSql("CUSTOM OBJECTS")), column, ssot.OutputSpec{})
		gotShow := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationFromObject(showObj)
		wantShow := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationWithPairedStructs(
			"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
			ShowPairedStructs(showObj),
			NewQueryStruct("ShowFakes").Show().SQL("CUSTOM OBJECTS"),
		)
		requireSameOperations(t, wantShow, gotShow)

		describeObj := showMethodObject(ssot.SdkDescribe(
			ssot.SdkDescribeSql("CUSTOM OBJECT"),
		), ssot.OutputSpec{}, column)
		gotDescribe := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DescribeOperationFromObject(describeObj)
		wantDescribe := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DescribeOperationWithPairedStructs(
			DescriptionMappingKindSingleValue,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-fake",
			DescribePairedStructs(describeObj),
			NewQueryStruct("DescFake").Describe().SQL("CUSTOM OBJECT").Name().WithValidation(ValidIdentifier, "name"),
		)
		requireSameOperations(t, wantDescribe, gotDescribe)
	})
}

func Test_ShowOperationFromObject_optionalSqlParts(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	column := ssot.OutputSpec{Columns: []*ssot.Attribute{&name}}
	obj := showMethodObject(ssot.SdkShow(ssot.SdkShowOptionalSql("HISTORY", "FUTURE")), column, ssot.OutputSpec{})
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-fakes",
		ShowPairedStructs(obj),
		NewQueryStruct("ShowFakes").Show().SQL("FAKES").OptionalSQL("HISTORY").OptionalSQL("FUTURE"),
	)
	requireSameOperations(t, want, got)
}

func Test_DropOperationFromObject(t *testing.T) {
	obj := dropMethodObject(ssot.SdkDrop(ssot.SdkDropSafelyHook(), ssot.SdkDropCascadeRestrict()))
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DropOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-fake",
		NewQueryStruct("DropFake").
			Drop().
			SQL("FAKE").
			IfExists().
			Name().
			WithValidation(ValidIdentifier, "name").
			OptionalSQL("CASCADE").
			OptionalSQL("RESTRICT").
			WithValidation(ConflictingFields, "Cascade", "Restrict"),
		WithDropSafelyHook(),
	)
	requireSameOperations(t, want, got)
}

func Test_DropOperationFromObject_panics(t *testing.T) {
	require.Panics(t, func() {
		NewInterface("Fakes", "Fake", "AccountObjectIdentifier").
			DropOperationFromObject(ssot.Object{Name: "Fake", NamePlural: "Fakes"})
	})
}

func Test_DropOperationFromObject_force(t *testing.T) {
	obj := dropMethodObject(ssot.SdkDrop(ssot.SdkDropSafelyForce()))
	got := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DropOperationFromObject(obj)
	want := NewInterface("Fakes", "Fake", "AccountObjectIdentifier").DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-fake",
		NewQueryStruct("DropFake").
			Drop().
			SQL("FAKE").
			IfExists().
			Name().
			WithValidation(ValidIdentifier, "name"),
		WithDropSafelyForce(),
	)
	requireSameOperations(t, want, got)
}

func dropMethodObject(opt ssot.SdkOption) ssot.Object {
	return ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(opt)},
	}
}

func showMethodObject(opt ssot.SdkOption, show, describe ssot.OutputSpec) ssot.Object {
	return ssot.Object{
		Name:       "Fake",
		NamePlural: "Fakes",
		IdType:     "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(opt)},
		Show:       show,
		Describe:   describe,
	}
}

func Test_CreateAlterFromObject(t *testing.T) {
	family := ssot.NewEnum("WidgetFamily", "WidgetFamilies", "SMALL")
	target := ssot.Attribute{Name: "target", Type: ssot.AttrIdentifier, CreateSql: ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("FOR APPLICATION"))}
	count := ssot.Attribute{Name: "count", Type: ssot.AttrInt, Required: true, CreateSql: ssot.IntParameter(), SetSql: ssot.IntParameter(ssot.Optional())}
	kind := ssot.Attribute{Name: "kind", Type: ssot.AttrEnum, Required: true, Enum: family, CreateSql: ssot.EnumParameter()}
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString, CreateSql: ssot.StringParameter(), SetSql: ssot.StringParameter(), UnsetSql: ssot.BoolKeyword()}
	items := ssot.Attribute{Name: "items", Type: ssot.AttrList, CreateSql: ssot.ListParameter("WidgetItem"), SetSql: ssot.ListParameter("WidgetItem"), UnsetSql: ssot.BoolKeyword()}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{items, label, kind, count, target},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&target, &count, &kind, ssot.Tags, &label, &items}},
		Set:        []*ssot.Attribute{&items, &count, &label},
		Unset:      []*ssot.Attribute{&label, &items},
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterKeyword("RESUME"),
			ssot.AlterKeyword("STOP ALL"),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(
			ssot.SdkCreate(ssot.SdkCreateWithoutOrReplace(), ssot.SdkCreateAdditionalValidations()),
			ssot.SdkAlter(ssot.SdkAlterSetAdditionalValidations()),
		)},
	}
	listItem := NewQueryStruct("WidgetItem").Enum("Value", family, KeywordOptions().SingleQuotes().Required())

	t.Run("create", func(t *testing.T) {
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj, listItem)
		want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
			NewQueryStruct("CreateWidget").
				Create().
				SQL("WIDGET").
				IfNotExists().
				Name().
				OptionalIdentifier("ForApplication", "AccountObjectIdentifier", IdentifierOptions().SQL("FOR APPLICATION")).
				NumberAssignment("COUNT", ParameterOptions().Required()).
				EnumAssignment("KIND", family, ParameterOptions().NoQuotes().Required()).
				OptionalTags().
				OptionalTextAssignment("LABEL", ParameterOptions().SingleQuotes()).
				ListAssignment("ITEMS", "WidgetItem", ParameterOptions().Parentheses()).
				WithValidation(ValidIdentifier, "name").
				WithAdditionalValidations(),
			listItem,
		)
		requireSameOperations(t, want, got)
	})

	t.Run("alter", func(t *testing.T) {
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
		want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
			NewQueryStruct("AlterWidget").
				Alter().
				SQL("WIDGET").
				IfExists().
				Name().
				OptionalSQL("RESUME").
				OptionalSQL("STOP ALL").
				OptionalQueryStructField(
					"Set",
					NewQueryStruct("WidgetSet").
						ListAssignment("ITEMS", "WidgetItem", ParameterOptions().Parentheses()).
						OptionalNumberAssignment("COUNT", ParameterOptions()).
						OptionalTextAssignment("LABEL", ParameterOptions().SingleQuotes()).
						WithValidation(AtLeastOneValueSet, "Items", "Count", "Label").
						WithAdditionalValidations(),
					KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Unset",
					NewQueryStruct("WidgetUnset").
						OptionalSQL("LABEL").
						OptionalSQL("ITEMS").
						WithValidation(AtLeastOneValueSet, "Label", "Items"),
					ListOptions().NoParentheses().SQL("UNSET"),
				).
				OptionalSetTags().
				OptionalUnsetTags().
				WithValidation(ValidIdentifier, "name").
				WithValidation(ExactlyOneValueSet, "Resume", "StopAll", "Set", "Unset", "SetTags", "UnsetTags"),
		)
		requireSameOperations(t, want, got)
	})
}

func Test_CreateOperationFromObject_atLeastOneOf(t *testing.T) {
	schedule := ssot.Attribute{Name: "schedule", Type: ssot.AttrString, CreateSql: ssot.StringParameter()}
	expire := ssot.Attribute{Name: "expire_after_days", Type: ssot.AttrInt, CreateSql: ssot.IntParameter()}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, CreateSql: ssot.StringParameter()}
	obj := ssot.Object{
		Name:         "Widget",
		Attributes:   []ssot.Attribute{schedule, expire, comment},
		AtLeastOneOf: [][]*ssot.Attribute{{&schedule, &expire}},
		Create:       ssot.CreateSpec{Clauses: []*ssot.Attribute{&schedule, &expire, &comment}},
		Generators:   []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate())},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			OrReplace().
			SQL("WIDGET").
			IfNotExists().
			Name().
			OptionalTextAssignment("SCHEDULE", ParameterOptions().SingleQuotes()).
			OptionalNumberAssignment("EXPIRE_AFTER_DAYS", ParameterOptions()).
			OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(AtLeastOneValueSet, "Schedule", "ExpireAfterDays"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateOperationFromObject_fakeObject(t *testing.T) {
	comment := ssot.Attribute{
		Name:      "comment",
		Type:      ssot.AttrString,
		CreateSql: ssot.StringParameter(),
	}
	count := ssot.Attribute{
		Name:      "count",
		Type:      ssot.AttrInt,
		Required:  true,
		CreateSql: ssot.IntParameter(),
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{comment, count},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&count, &comment, ssot.Tags}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(
			ssot.SdkCreateAdditionalValidations(),
		))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			OrReplace().
			SQL("WIDGET").
			IfNotExists().
			Name().
			NumberAssignment("COUNT", ParameterOptions().Required()).
			OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
			OptionalTags().
			WithValidation(ValidIdentifier, "name").
			WithValidation(ConflictingFields, "OrReplace", "IfNotExists").
			WithAdditionalValidations(),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterOperationFromObject_unsetConflicts(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	schedule := ssot.Attribute{Name: "schedule", Type: ssot.AttrString}
	expire := ssot.Attribute{Name: "expire_after_days", Type: ssot.AttrInt}
	obj := ssot.Object{
		Name:           "Widget",
		Attributes:     []ssot.Attribute{comment, schedule, expire},
		Unset:          []*ssot.Attribute{&comment, &schedule, &expire},
		UnsetConflicts: [][]*ssot.Attribute{{&schedule, &expire}},
		Alter:          ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterUnset()}},
		Generators:     []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(ssot.SdkAlterWithoutIfExists(), ssot.SdkAlterWithoutExactlyOneValueSet()))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField(
				"Unset",
				NewQueryStruct("WidgetUnset").
					OptionalSQL("COMMENT").
					OptionalSQL("SCHEDULE").
					OptionalSQL("EXPIRE_AFTER_DAYS").
					WithValidation(AtLeastOneValueSet, "Comment", "Schedule", "ExpireAfterDays").
					WithValidation(ConflictingFields, "Schedule", "ExpireAfterDays"),
				ListOptions().NoParentheses().SQL("UNSET"),
			).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterOperationFromObject_setList(t *testing.T) {
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString, SetSql: ssot.StringParameter()}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{label},
		Set:        []*ssot.Attribute{&label},
		Alter:      ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterSet()}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(
			ssot.SdkAlterSetList(),
			ssot.SdkAlterWithoutIfExists(),
			ssot.SdkAlterWithoutExactlyOneValueSet(),
		))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField(
				"Set",
				NewQueryStruct("WidgetSet").
					OptionalTextAssignment("LABEL", ParameterOptions().SingleQuotes()).
					WithValidation(AtLeastOneValueSet, "Label"),
				ListOptions().NoParentheses().SQL("SET"),
			).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterOperationFromObject_requiredIdentifierList(t *testing.T) {
	accounts := ssot.Attribute{
		Name:   "accounts",
		Type:   ssot.AttrList,
		Elem:   &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountIdentifier"},
		SetSql: ssot.PlainParameter(),
	}
	block := ssot.Attribute{
		Name:      "add",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children:  []ssot.Attribute{ssot.WithCreateSql(&accounts, ssot.PlainParameter(ssot.Required()))},
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{accounts},
		Alter:      ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterKeywordBlock("Add", "ADD", &block)}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(ssot.SdkAlterWithoutIfExists(), ssot.SdkAlterWithoutExactlyOneValueSet()))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField(
				"Add",
				NewQueryStruct("WidgetAdd").
					ListAssignment("ACCOUNTS", "AccountIdentifier", ParameterOptions().Required()).
					WithValidation(AtLeastOneValueSet, "Accounts").
					WithValidation(ValidIdentifier, "Accounts"),
				KeywordOptions().SQL("ADD"),
			).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateAlterFromObject_parameters(t *testing.T) {
	orig := parametersForLevel
	t.Cleanup(func() { parametersForLevel = orig })
	params := []ssot.ParameterDef{
		{SqlName: "DATA_RETENTION_TIME_IN_DAYS", Kind: string(ssot.AttrInt)},
		{SqlName: "EXTERNAL_VOLUME", Kind: "AccountObjectIdentifier"},
	}
	parametersForLevel = func(ssot.ParameterLevel) []ssot.ParameterDef { return params }

	comment := ssot.Attribute{
		Name:      "comment",
		Type:      ssot.AttrString,
		CreateSql: ssot.StringParameter(),
		SetSql:    ssot.StringParameter(),
		UnsetSql:  ssot.BoolKeyword(),
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{comment},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{ssot.Parameters, &comment, ssot.Tags}},
		Set:        []*ssot.Attribute{ssot.Parameters, &comment},
		Unset:      []*ssot.Attribute{ssot.Parameters, &comment},
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(), ssot.SdkAlter())},
	}

	t.Run("create", func(t *testing.T) {
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
		want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
			NewQueryStruct("CreateWidget").
				Create().
				OrReplace().
				SQL("WIDGET").
				IfNotExists().
				Name().
				WithParameters(params...).
				OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
				OptionalTags().
				WithValidation(ValidIdentifier, "name").
				WithValidation(ConflictingFields, "OrReplace", "IfNotExists").
				WithValidation(ValidIdentifierIfSet, "ExternalVolume"),
		)
		requireSameOperations(t, want, got)
	})

	t.Run("alter", func(t *testing.T) {
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
		want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
			NewQueryStruct("AlterWidget").
				Alter().
				SQL("WIDGET").
				IfExists().
				Name().
				OptionalQueryStructField(
					"Set",
					NewQueryStruct("WidgetSet").
						WithParameters(params...).
						OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
						WithValidation(ValidIdentifierIfSet, "ExternalVolume").
						WithValidation(AtLeastOneValueSet, "DataRetentionTimeInDays", "ExternalVolume", "Comment"),
					KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Unset",
					NewQueryStruct("WidgetUnset").
						WithParametersUnset(params...).
						OptionalSQL("COMMENT").
						WithValidation(AtLeastOneValueSet, "DataRetentionTimeInDays", "ExternalVolume", "Comment"),
					ListOptions().NoParentheses().SQL("UNSET"),
				).
				OptionalSetTags().
				OptionalUnsetTags().
				WithValidation(ValidIdentifier, "name").
				WithValidation(ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
		)
		requireSameOperations(t, want, got)
	})

	t.Run("parameters stay first in AtLeastOneValueSet when they follow attributes", func(t *testing.T) {
		later := obj
		later.Set = []*ssot.Attribute{&comment, ssot.Parameters}
		later.Unset = []*ssot.Attribute{&comment, ssot.Parameters}
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(later)
		want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
			NewQueryStruct("AlterWidget").
				Alter().
				SQL("WIDGET").
				IfExists().
				Name().
				OptionalQueryStructField(
					"Set",
					NewQueryStruct("WidgetSet").
						OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
						WithParameters(params...).
						WithValidation(ValidIdentifierIfSet, "ExternalVolume").
						WithValidation(AtLeastOneValueSet, "DataRetentionTimeInDays", "ExternalVolume", "Comment"),
					KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Unset",
					NewQueryStruct("WidgetUnset").
						OptionalSQL("COMMENT").
						WithParametersUnset(params...).
						WithValidation(AtLeastOneValueSet, "DataRetentionTimeInDays", "ExternalVolume", "Comment"),
					ListOptions().NoParentheses().SQL("UNSET"),
				).
				OptionalSetTags().
				OptionalUnsetTags().
				WithValidation(ValidIdentifier, "name").
				WithValidation(ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags"),
		)
		requireSameOperations(t, want, got)
	})
}

func Test_parametersForStatement_addsAlterLevels(t *testing.T) {
	orig := parametersForLevel
	t.Cleanup(func() { parametersForLevel = orig })
	base := []ssot.ParameterDef{{SqlName: "STATEMENT_TIMEOUT_IN_SECONDS", Kind: string(ssot.AttrInt)}}
	extra := []ssot.ParameterDef{
		{SqlName: "FALLBACK_WAREHOUSE", Kind: "AccountObjectIdentifier"},
		{SqlName: "STATEMENT_TIMEOUT_IN_SECONDS", Kind: string(ssot.AttrInt)},
	}
	parametersForLevel = func(level ssot.ParameterLevel) []ssot.ParameterDef {
		if level == ssot.ParameterLevelWarehouseInteractive {
			return extra
		}
		return base
	}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:                                  "Widget",
		ShowParametersDetailsAdditionalLevels: []ssot.ParameterLevel{ssot.ParameterLevelWarehouseInteractive},
		Attributes:                            []ssot.Attribute{comment},
		Create:                                ssot.CreateSpec{Clauses: []*ssot.Attribute{ssot.Parameters}},
		Set:                                   []*ssot.Attribute{ssot.Parameters},
		Unset:                                 []*ssot.Attribute{ssot.Parameters},
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterSet(),
			ssot.AlterUnset(),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(
			ssot.SdkCreate(ssot.SdkCreateWithoutOrReplace(), ssot.SdkCreateWithoutIfNotExists()),
			ssot.SdkAlter(ssot.SdkAlterWithoutIfExists(), ssot.SdkAlterWithoutExactlyOneValueSet()),
		)},
	}

	gotCreate := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	wantCreate := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			WithParameters(base...).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, wantCreate, gotCreate)
}

func Test_AlterOperationFromObject_clauseOrder(t *testing.T) {
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString, SetSql: ssot.StringParameter(), UnsetSql: ssot.BoolKeyword()}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{label},
		Set:        []*ssot.Attribute{&label},
		Unset:      []*ssot.Attribute{&label},
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterKeyword("ENABLE MANAGED ACCESS"),
			ssot.AlterKeyword("DISABLE MANAGED ACCESS"),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(ssot.SdkAlterWithoutIfExists()))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField(
				"Set",
				NewQueryStruct("WidgetSet").
					OptionalTextAssignment("LABEL", ParameterOptions().SingleQuotes()).
					WithValidation(AtLeastOneValueSet, "Label"),
				KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				NewQueryStruct("WidgetUnset").
					OptionalSQL("LABEL").
					WithValidation(AtLeastOneValueSet, "Label"),
				ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalSQL("ENABLE MANAGED ACCESS").
			OptionalSQL("DISABLE MANAGED ACCESS").
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags", "EnableManagedAccess", "DisableManagedAccess"),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterOperationFromObject_listAndOptionalKeyword(t *testing.T) {
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{label},
		Set:        []*ssot.Attribute{&label},
		Unset:      []*ssot.Attribute{&label},
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterKeyword("RESUME"),
			ssot.AlterKeywordNotExactlyOne("IF SUSPENDED"),
			ssot.AlterRenameTo(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterList("ADD TABLES", "SchemaObjectIdentifier"),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(
			ssot.SdkAlterWithoutIfExists(),
			ssot.SdkAlterAdditionalValidations(),
		))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			OptionalSQL("RESUME").
			OptionalSQL("IF SUSPENDED").
			RenameTo().
			OptionalQueryStructField(
				"Set",
				NewQueryStruct("WidgetSet").
					OptionalTextAssignment("LABEL", ParameterOptions().SingleQuotes()).
					WithValidation(AtLeastOneValueSet, "Label"),
				KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				NewQueryStruct("WidgetUnset").
					OptionalSQL("LABEL").
					WithValidation(AtLeastOneValueSet, "Label"),
				ListOptions().NoParentheses().SQL("UNSET"),
			).
			ListAssignment("ADD TABLES", "SchemaObjectIdentifier", ParameterOptions().Parentheses().NoEquals()).
			OptionalSetTags().
			OptionalUnsetTags().
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "Resume", "RenameTo", "Set", "Unset", "AddTables", "SetTags", "UnsetTags").
			WithValidation(ValidIdentifierIfSet, "RenameTo").
			WithValidation(ValidIdentifier, "AddTables").
			WithAdditionalValidations(),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterOperationFromObject_renameSwap(t *testing.T) {
	label := ssot.Attribute{Name: "label", Type: ssot.AttrString, SetSql: ssot.StringParameter(), UnsetSql: ssot.BoolKeyword()}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{label},
		Set:        []*ssot.Attribute{&label},
		Unset:      []*ssot.Attribute{&label},
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSwapWith(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterKeyword("ENABLE MANAGED ACCESS"),
			ssot.AlterKeyword("DISABLE MANAGED ACCESS"),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(ssot.SdkAlterWithoutIfExists()))},
	}
	got := NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			RenameTo().
			Identifier("SwapWith", InterfaceIdentifierPointerKind, IdentifierOptions().SQL("SWAP WITH")).
			OptionalQueryStructField(
				"Set",
				NewQueryStruct("WidgetSet").
					OptionalTextAssignment("LABEL", ParameterOptions().SingleQuotes()).
					WithValidation(AtLeastOneValueSet, "Label"),
				KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				NewQueryStruct("WidgetUnset").
					OptionalSQL("LABEL").
					WithValidation(AtLeastOneValueSet, "Label"),
				ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalSQL("ENABLE MANAGED ACCESS").
			OptionalSQL("DISABLE MANAGED ACCESS").
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "RenameTo", "SwapWith", "Set", "Unset", "SetTags", "UnsetTags", "EnableManagedAccess", "DisableManagedAccess").
			WithValidation(ValidIdentifierIfSet, "RenameTo").
			WithValidation(ValidIdentifierIfSet, "SwapWith"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateOperationFromObject_beforeObject(t *testing.T) {
	transient := ssot.Attribute{Name: "is_transient", Type: ssot.AttrBool, CreateSql: ssot.BoolKeyword(ssot.Sql("TRANSIENT"))}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, CreateSql: ssot.StringParameter()}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{comment, transient},
		Create: ssot.CreateSpec{
			KeywordsBefore: []*ssot.Attribute{&transient},
			Clauses:        []*ssot.Attribute{&comment},
		},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate())},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			OrReplace().
			OptionalSQL("TRANSIENT").
			SQL("WIDGET").
			IfNotExists().
			Name().
			OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ConflictingFields, "OrReplace", "IfNotExists"),
	)
	requireSameOperations(t, want, got)

	require.PanicsWithValue(t, "KeywordsBefore list contains tags", func() {
		tagged := obj
		tagged.Create.KeywordsBefore = []*ssot.Attribute{ssot.Tags}
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(tagged)
	})
}

func Test_CreateOperationFromObject_enumListItem(t *testing.T) {
	family := ssot.NewEnum("WidgetFamily", "WidgetFamilies", "SMALL")
	items := ssot.Attribute{
		Name: "items",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family},
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{items},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&items}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(
			ssot.SdkCreateWithoutOrReplace(),
			ssot.SdkCreateWithoutIfNotExists(),
		))},
	}
	item := NewQueryStruct("WidgetFamilyListItem").Enum("Value", family, KeywordOptions().SingleQuotes().Required())
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			ListAssignment("ITEMS", "WidgetFamilyListItem", ParameterOptions().Parentheses()).
			WithValidation(ValidIdentifier, "name"),
		item,
	)
	requireSameOperations(t, want, got)
}

func Test_CreateVariantOperationsFromObject_fakeObject(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	kind := ssot.Attribute{
		Name:      "kind",
		SqlOnly:   true,
		CreateSql: ssot.StaticSQL("kind", "KIND = 'MINI'"),
	}
	base := ssot.Object{
		Name: "Widget",
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkCreateForVariants()),
		},
	}
	variant := ssot.Object{
		Name:                       "WidgetMini",
		VariantOf:                  &base,
		VariantEntityCreateName:    "MINI WIDGET",
		SdkCreateRequestObjectName: "MiniWidget",
		Create: ssot.CreateSpec{
			Doc:                   "https://docs.example.com/mini",
			AdditionalValidations: true,
			Clauses:               []*ssot.Attribute{&kind, &comment},
		},
	}
	orig := variantsOf
	t.Cleanup(func() { variantsOf = orig })
	variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{variant} }
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateVariantOperationsFromObject(base)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CustomOperation(
		"CreateMini",
		"https://docs.example.com/mini",
		NewQueryStruct("CreateMiniWidget").
			Create().
			OrReplace().
			SQL("MINI WIDGET").
			IfNotExists().
			Name().
			SQLWithCustomFieldName("kind", "KIND = 'MINI'").
			OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ConflictingFields, "OrReplace", "IfNotExists").
			WithAdditionalValidations(),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateVariantOperationsFromObject_withoutOrReplace(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	base := ssot.Object{
		Name: "Widget",
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkCreateForVariants()),
		},
	}
	variant := ssot.Object{
		Name:                       "WidgetMini",
		VariantOf:                  &base,
		VariantEntityCreateName:    "MINI WIDGET",
		SdkCreateRequestObjectName: "MiniWidget",
		Create: ssot.CreateSpec{
			Doc:              "https://docs.example.com/mini",
			WithoutOrReplace: true,
			Clauses:          []*ssot.Attribute{&comment},
		},
	}
	orig := variantsOf
	t.Cleanup(func() { variantsOf = orig })
	variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{variant} }
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateVariantOperationsFromObject(base)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CustomOperation(
		"CreateMini",
		"https://docs.example.com/mini",
		NewQueryStruct("CreateMiniWidget").
			Create().
			SQL("MINI WIDGET").
			IfNotExists().
			Name().
			OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateVariantOperationsFromObject_panics(t *testing.T) {
	base := ssot.Object{Name: "Widget"}
	variant := ssot.Object{Name: "WidgetMini", VariantOf: &base}
	orig := variantsOf
	t.Cleanup(func() { variantsOf = orig })
	variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{variant} }
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateVariantOperationsFromObject(base)
	})
	base.Generators = []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreateForVariants())}
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateVariantOperationsFromObject(base)
	})
	variantsOf = func(ssot.Object) []ssot.Object { return nil }
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateVariantOperationsFromObject(base)
	})
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(variant)
	})
}

func Test_AlterVariantOperationsFromObject_fakeObject(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	base := ssot.Object{
		Name: "Widget",
		Generators: []ssot.GeneratorOption{
			ssot.WithSdk(ssot.SdkAlterForVariants()),
		},
	}
	variant := ssot.Object{
		Name:                       "WidgetMini",
		VariantOf:                  &base,
		SdkCreateRequestObjectName: "MiniWidget",
		Set:                        []*ssot.Attribute{&comment},
		Unset:                      []*ssot.Attribute{&comment},
		Alter: ssot.AlterSpec{
			Doc:                   "https://docs.example.com/alter-mini",
			AdditionalValidations: true,
			Clauses: []ssot.AlterClause{
				ssot.AlterRenameTo(),
				ssot.AlterSet(),
				ssot.AlterUnset(),
			},
		},
	}
	orig := variantsOf
	t.Cleanup(func() { variantsOf = orig })
	variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{variant} }
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterVariantOperationsFromObject(base)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CustomOperation(
		"AlterMini",
		"https://docs.example.com/alter-mini",
		NewQueryStruct("AlterMiniWidget").
			Alter().
			SQL("WIDGET").
			IfExists().
			Name().
			RenameTo().
			OptionalQueryStructField(
				"Set",
				NewQueryStruct("WidgetMiniSet").
					OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
					WithValidation(AtLeastOneValueSet, "Comment"),
				KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				NewQueryStruct("WidgetMiniUnset").
					OptionalSQL("COMMENT").
					WithValidation(AtLeastOneValueSet, "Comment"),
				ListOptions().NoParentheses().SQL("UNSET"),
			).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "RenameTo", "Set", "Unset").
			WithValidation(ValidIdentifierIfSet, "RenameTo").
			WithAdditionalValidations(),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterVariantOperationsFromObject_panics(t *testing.T) {
	base := ssot.Object{Name: "Widget"}
	variant := ssot.Object{Name: "WidgetMini", VariantOf: &base}
	orig := variantsOf
	t.Cleanup(func() { variantsOf = orig })
	variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{variant} }
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterVariantOperationsFromObject(base)
	})
	base.Generators = []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlterForVariants())}
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterVariantOperationsFromObject(base)
	})
	variantsOf = func(ssot.Object) []ssot.Object { return nil }
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterVariantOperationsFromObject(base)
	})
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(variant)
	})
}

func Test_CustomOperationsFromObject_fakeObject(t *testing.T) {
	transient := ssot.Attribute{
		Name:      "is_transient",
		Type:      ssot.AttrBool,
		CreateSql: ssot.BoolKeyword(ssot.Sql("TRANSIENT")),
	}
	obj := ssot.Object{
		Name:       "Fake",
		IdType:     "DatabaseObjectIdentifier",
		Attributes: []ssot.Attribute{transient},
		Create:     ssot.CreateSpec{KeywordsBefore: []*ssot.Attribute{&transient}},
		CustomOperations: []ssot.CustomOperationSpec{
			{
				Kind:                  "Clone",
				Doc:                   "https://docs.snowflake.com/en/sql-reference/sql/create-fake",
				AdditionalValidations: true,
				Clauses: []ssot.CustomOperationClause{
					ssot.CustomStaticSQL("CREATE"),
					ssot.CustomOrReplace(),
					ssot.CustomAttribute(&transient),
					ssot.CustomObjectWords(),
					ssot.CustomIfNotExists(),
					ssot.CustomName(),
					ssot.CustomPredefined("Clone"),
				},
			},
			{
				Kind: "Undrop",
				Clauses: []ssot.CustomOperationClause{
					ssot.CustomStaticSQL("UNDROP"),
					ssot.CustomObjectWords(),
					ssot.CustomName(),
				},
			},
		},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCustomOperations())},
	}
	got := NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomOperationsFromObject(obj)
	want := NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").
		CustomOperation(
			"Clone",
			"https://docs.snowflake.com/en/sql-reference/sql/create-fake",
			NewQueryStruct("CloneFake").
				Create().
				OrReplace().
				OptionalSQL("TRANSIENT").
				SQL("FAKE").
				IfNotExists().
				Name().
				PredefinedQueryStructField("Clone", "Clone", KeywordOptions()).
				WithValidation(ValidIdentifier, "name").
				WithValidation(ConflictingFields, "OrReplace", "IfNotExists").
				WithAdditionalValidations(),
		).
		CustomOperation(
			"Undrop",
			"https://docs.snowflake.com/en/sql-reference/sql/undrop-fake",
			NewQueryStruct("UndropFake").
				SQL("UNDROP").
				SQL("FAKE").
				Name().
				WithValidation(ValidIdentifier, "name"),
		)
	requireSameOperations(t, want, got)
}

func Test_CustomOperationsFromObject_panics(t *testing.T) {
	transient := ssot.Attribute{
		Name:      "is_transient",
		Type:      ssot.AttrBool,
		CreateSql: ssot.BoolKeyword(ssot.Sql("TRANSIENT")),
	}
	enabled := ssot.Object{
		Name:       "Fake",
		Attributes: []ssot.Attribute{transient},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCustomOperations())},
	}
	require.Panics(t, func() {
		NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").
			CustomOperationsFromObject(ssot.Object{Name: "Fake"})
	})
	require.Panics(t, func() {
		bad := enabled
		bad.CustomOperations = []ssot.CustomOperationSpec{{
			Kind:    "Clone",
			Clauses: []ssot.CustomOperationClause{ssot.CustomStaticSQL("")},
		}}
		NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomOperationsFromObject(bad)
	})
	require.Panics(t, func() {
		bad := enabled
		bad.CustomOperations = []ssot.CustomOperationSpec{{
			Kind:    "Clone",
			Clauses: []ssot.CustomOperationClause{ssot.CustomAttribute(nil)},
		}}
		NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomOperationsFromObject(bad)
	})
	require.Panics(t, func() {
		bad := enabled
		bad.CustomOperations = []ssot.CustomOperationSpec{{
			Kind:    "Clone",
			Clauses: []ssot.CustomOperationClause{{Kind: "nope"}},
		}}
		NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomOperationsFromObject(bad)
	})
}

func Test_CustomInterfaceMethodsFromObject_fakeObject(t *testing.T) {
	obj := ssot.Object{
		Name:       "Fake",
		IdType:     "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCustomInterfaceMethods())},
		CustomInterfaceMethods: []ssot.CustomInterfaceMethodSpec{{
			Name: "Use",
		}},
	}
	got := NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomInterfaceMethodsFromObject(obj)
	require.Equal(t, []*CustomInterfaceMethod{{
		Name:        "Use",
		Parameters:  []*MethodParameter{NewMethodParameter("id", "DatabaseObjectIdentifier")},
		ReturnTypes: []string{"error"},
	}}, got.CustomMethods)

	overridden := obj
	overridden.CustomInterfaceMethods = []ssot.CustomInterfaceMethodSpec{{
		Name: "Use",
		Doc:  "Use is based on https://custom/use",
		Parameters: []ssot.CustomInterfaceMethodParameter{{
			Name: "request",
			Kind: "*UseFakeRequest",
		}},
		Returns: []string{"bool", "error"},
	}}
	got = NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomInterfaceMethodsFromObject(overridden)
	require.Equal(t, []*CustomInterfaceMethod{{
		Name:        "Use",
		Doc:         "Use is based on https://custom/use",
		Parameters:  []*MethodParameter{NewMethodParameter("request", "*UseFakeRequest")},
		ReturnTypes: []string{"bool", "error"},
	}}, got.CustomMethods)

	emptyDoc := obj
	emptyDoc.CustomInterfaceMethods = []ssot.CustomInterfaceMethodSpec{{
		Name:    "ShowByIDExperimentalSafely",
		Returns: []string{"*Fake", "error"},
	}}
	got = NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomInterfaceMethodsFromObject(emptyDoc)
	require.Equal(t, []*CustomInterfaceMethod{{
		Name:        "ShowByIDExperimentalSafely",
		Parameters:  []*MethodParameter{NewMethodParameter("id", "DatabaseObjectIdentifier")},
		ReturnTypes: []string{"*Fake", "error"},
	}}, got.CustomMethods)
}

func Test_CustomInterfaceMethodsFromObject_panics(t *testing.T) {
	require.Panics(t, func() {
		NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").
			CustomInterfaceMethodsFromObject(ssot.Object{Name: "Fake"})
	})
	require.Panics(t, func() {
		NewInterface("Fakes", "Fake", "DatabaseObjectIdentifier").CustomInterfaceMethodsFromObject(ssot.Object{
			Name:                   "Fake",
			Generators:             []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCustomInterfaceMethods())},
			CustomInterfaceMethods: []ssot.CustomInterfaceMethodSpec{{Name: "Use"}},
		})
	})
}

func Test_ShowParametersFromObject_fakeObject(t *testing.T) {
	orig := parametersForLevel
	t.Cleanup(func() { parametersForLevel = orig })
	parametersForLevel = func(ssot.ParameterLevel) []ssot.ParameterDef {
		return []ssot.ParameterDef{
			{SqlName: "DATA_RETENTION_TIME_IN_DAYS", Kind: "int"},
			{SqlName: "LOG_LEVEL", Kind: "LogLevel"},
		}
	}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(
			ssot.SdkShowParameters(),
			ssot.SdkShowParametersDetails(),
		)},
	}
	got := NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").
		ShowParametersFromObject(obj).
		ShowParametersDetailsFromObject(obj)
	require.Equal(t, []*CustomInterfaceMethod{
		{
			Name:        "ShowParameters",
			Parameters:  []*MethodParameter{NewMethodParameter("id", "DatabaseObjectIdentifier")},
			ReturnTypes: []string{"[]*Parameter", "error"},
		},
		{
			Name:        "ShowParametersDetails",
			Parameters:  []*MethodParameter{NewMethodParameter("id", "DatabaseObjectIdentifier")},
			ReturnTypes: []string{"*WidgetParametersDetails", "error"},
		},
	}, got.CustomMethods)
	require.Equal(t, &ParametersDetailsConfig{
		Fields: []ParameterDetailField{
			{FieldName: "DataRetentionTimeInDays", Key: "DATA_RETENTION_TIME_IN_DAYS", GoType: "int", Parser: "strconv.Atoi"},
			{FieldName: "LogLevel", Key: "LOG_LEVEL", GoType: "LogLevel", Parser: "ToLogLevel"},
		},
	}, got.ParametersDetails)
}

func Test_ShowParametersDetailsFromObject_includesAlterLevels(t *testing.T) {
	obj := ssot.Object{
		Name:                                  "Widget",
		IdType:                                "DatabaseObjectIdentifier",
		ShowParametersDetailsAdditionalLevels: []ssot.ParameterLevel{ssot.ParameterLevelWarehouseInteractive},
		Generators:                            []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkShowParametersDetails())},
	}
	got := NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").ShowParametersDetailsFromObject(obj)
	keys := make([]string, len(got.ParametersDetails.Fields))
	for i, field := range got.ParametersDetails.Fields {
		keys[i] = field.Key
	}
	require.Equal(t, []string{
		"FALLBACK_WAREHOUSE",
		"MAX_CONCURRENCY_LEVEL",
		"STATEMENT_QUEUED_TIMEOUT_IN_SECONDS",
		"STATEMENT_TIMEOUT_IN_SECONDS",
	}, keys)
}

func Test_ShowParametersFromObject_panics(t *testing.T) {
	obj := ssot.Object{Name: "Widget", IdType: "DatabaseObjectIdentifier"}
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").ShowParametersFromObject(obj)
	})
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").ShowParametersDetailsFromObject(obj)
	})
}

func Test_InterfaceFromObject_enabledOperationsOnly(t *testing.T) {
	obj := ssot.Object{
		Name:       "Widget",
		NamePlural: "Widgets",
		IdType:     "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDefineFromSsot(), ssot.SdkDrop())},
	}
	got := InterfaceFromObject(obj)
	require.Equal(t, "Widgets", got.Name)
	require.Equal(t, "Widget", got.NameSingular)
	require.Equal(t, "AccountObjectIdentifier", got.IdentifierKind)
	require.Len(t, got.Operations, 1)
	require.Equal(t, "Drop", got.Operations[0].Name)
	require.Empty(t, got.Enums)
}

func Test_SdkInterfacesFromSsot_uniquePluralAndFlag(t *testing.T) {
	base := ssot.Object{
		Name:       "Warehouse",
		NamePlural: "Warehouses",
		IdType:     "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDefineFromSsot(), ssot.SdkDrop())},
	}
	variant := ssot.Object{
		Name:       "WarehouseAdaptive",
		NamePlural: "Warehouses",
		IdType:     "AccountObjectIdentifier",
		VariantOf:  &base,
	}
	objects := []ssot.Object{
		{
			Name:       "VariantA",
			NamePlural: "Family",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())},
		},
		{Name: "VariantB", NamePlural: "Family", IdType: "AccountObjectIdentifier"},
		{
			Name:       "Widget",
			NamePlural: "Widgets",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDefineFromSsot(), ssot.SdkDrop())},
		},
		{
			Name:       "DescribedOnly",
			NamePlural: "DescribedOnlys",
			IdType:     "AccountObjectIdentifier",
			Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())},
		},
		{Name: "Quiet", NamePlural: "Quiets", IdType: "AccountObjectIdentifier"},
		base,
		variant,
	}
	got := SdkInterfacesFromSsot(objects, nil)
	require.Len(t, got, 2)
	require.Equal(t, "Widgets", got[0].Name)
	require.Equal(t, "Warehouses", got[1].Name)
}

func Test_SdkInterfacesFromSsot_override(t *testing.T) {
	objects := []ssot.Object{{
		Name:       "Widget",
		NamePlural: "Widgets",
		IdType:     "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDefineFromSsot(), ssot.SdkDrop())},
	}}
	got := SdkInterfacesFromSsot(objects, map[string]SdkInterfaceOverride{
		"Widget": func(iface *Interface) {
			iface.WithShowObjectType("WidgetAlias")
		},
	})
	require.Equal(t, "WidgetAlias", got[0].ShowObjectTypeName)
	require.Panics(t, func() {
		SdkInterfacesFromSsot(objects, map[string]SdkInterfaceOverride{
			"Missing": func(*Interface) {},
		})
	})
}

func Test_MergeSdkDefinitions_skipsOldNameSsotEmits(t *testing.T) {
	old := []*Interface{{Name: "Widgets"}, {Name: "OldOnes"}}
	fromSsot := []*Interface{{Name: "Widgets"}, {Name: "Fresh"}}
	got := MergeSdkDefinitions(old, fromSsot)
	names := make([]string, len(got))
	for i, iface := range got {
		names[i] = iface.Name
	}
	require.Equal(t, []string{"OldOnes", "Widgets", "Fresh"}, names)
}

func Test_CreateAlterFromObject_panics(t *testing.T) {
	obj := ssot.Object{Name: "Widget", NamePlural: "Widgets"}
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	})
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	})
}

func requireSameOperations(t *testing.T, want, got *Interface) {
	t.Helper()
	require.Equal(t, want.ShowObjectName, got.ShowObjectName)
	require.Len(t, got.Operations, len(want.Operations))
	for i := range want.Operations {
		want.Operations[i].ObjectInterface = nil
		got.Operations[i].ObjectInterface = nil
		require.Equal(t, want.Operations[i], got.Operations[i])
	}
}

func Test_AlterOperationFromObject_keywordListAndPlainText(t *testing.T) {
	obj := ssot.Object{
		Name: "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterKeywordList("ModifyCondition", "MODIFY CONDITION EXISTS", "string"),
			ssot.AlterPlainText("ModifyAction", "MODIFY ACTION"),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter(ssot.SdkAlterWithoutIfExists()))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			Name().
			UnnamedList("ModifyCondition", "string", KeywordOptions().SQL("MODIFY CONDITION EXISTS").Parentheses().NoComma()).
			OptionalAssignmentWithFieldName("MODIFY ACTION", "string", ParameterOptions().NoEquals(), "ModifyAction").
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "ModifyCondition", "ModifyAction"),
	)
	requireSameOperations(t, want, got)
}

func Test_AlterOperationFromObject_textAndBlocks(t *testing.T) {
	widget := newWidgetLocationAttributes()
	locations := widget.list
	update := ssot.Attribute{
		Name:      "update_location",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			{Name: "location_name", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", CreateSql: ssot.StringParameter(ssot.Sql("LOCATION"), ssot.NoEquals(), ssot.Required())},
		},
	}
	obj := ssot.Object{
		Name: "Widget",
		Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
			ssot.AlterText("REMOVE LOCATION"),
			ssot.AlterParameterBlock("AddLocation", "ADD LOCATION", &locations),
			ssot.AlterKeywordBlock("UpdateLocation", "UPDATE", &update),
		}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter())},
	}

	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-widget",
		NewQueryStruct("AlterWidget").
			Alter().
			SQL("WIDGET").
			IfExists().
			Name().
			OptionalTextAssignment("REMOVE LOCATION", ParameterOptions().SingleQuotes().NoEquals()).
			OptionalQueryStructField("AddLocation", queryStructFromObjectList("Widget", locations), ParameterOptions().SQL("ADD LOCATION")).
			OptionalQueryStructField("UpdateLocation", NewQueryStruct("WidgetUpdateLocation").
				TextAssignment("LOCATION", ParameterOptions().SingleQuotes().NoEquals().Required()), KeywordOptions().SQL("UPDATE")).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "RemoveLocation", "AddLocation", "UpdateLocation"),
	)
	requireSameOperations(t, want, got)
}

func Test_DescribeOperationFromObject_propertyList(t *testing.T) {
	parent := ssot.Attribute{Name: "parent_property", Type: ssot.AttrString, Readonly: true}
	property := ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true}
	value := ssot.Attribute{Name: "property_value", Type: ssot.AttrString, Readonly: true}
	active := ssot.Attribute{Name: "active", Type: ssot.AttrString, Readonly: true}
	obj := ssot.Object{
		Name:         "Widget",
		NamePlural:   "Widgets",
		IdType:       "AccountObjectIdentifier",
		Describe:     ssot.OutputSpec{Columns: []*ssot.Attribute{&active}},
		DescribeRows: ssot.OutputSpec{Columns: []*ssot.Attribute{&parent, &property, &value}},
		Generators:   []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe())},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationWithPairedStructs(
		DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-widget",
		StructPair("widgetDescRow", "WidgetProperty").
			Text("parent_property").
			Text("property").
			Text("property_value"),
		NewQueryStruct("DescWidget").
			Describe().
			SQL("WIDGET").
			Name().
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_DescribeOperationFromObject_details(t *testing.T) {
	property := ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true}
	locations := ssot.Attribute{Name: "storage_locations", Type: ssot.AttrList, Readonly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	active := ssot.Attribute{Name: "active", Type: ssot.AttrString, Readonly: true}
	size := ssot.Attribute{Name: "size", Type: ssot.AttrInt, Readonly: true, OptionalInShow: true}
	obj := ssot.Object{
		Name:         "Widget",
		NamePlural:   "Widgets",
		IdType:       "AccountObjectIdentifier",
		Describe:     ssot.OutputSpec{Columns: []*ssot.Attribute{&locations, &active, &size}},
		DescribeRows: ssot.OutputSpec{Columns: []*ssot.Attribute{&property}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeDetails(ssot.SdkDescribeDetailsNestedField(&locations)),
		))},
	}
	want := PlainStruct("WidgetDetails").
		Field("Id", "AccountObjectIdentifier").
		Field("StorageLocations", "[]WidgetStorageLocationDetails").
		Field("Active", "string").
		Field("Size", "*int")
	require.Equal(t, want.IntoField(), DescribeDetailsPlain(obj).IntoField())

	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	helperNames := collections.Map(got.Operations[0].HelperStructs, func(f *Field) string { return f.Name })
	require.Contains(t, helperNames, "WidgetDetails")
	require.Contains(t, helperNames, "WidgetProperty")

	id := ssot.Attribute{Name: "id", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Readonly: true}
	obj.Describe = ssot.OutputSpec{Columns: []*ssot.Attribute{&id, &active}}
	obj.Generators = []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(ssot.SdkDescribeDetails()))}
	require.Equal(t,
		PlainStruct("WidgetDetails").Field("Id", "AccountObjectIdentifier").Field("Active", "string").IntoField(),
		DescribeDetailsPlain(obj).IntoField(),
		"a column named id replaces the leading Id field",
	)
}

func Test_DescribeOperationFromObject_detailsMethod(t *testing.T) {
	property := ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true}
	active := ssot.Attribute{Name: "active", Type: ssot.AttrString, Readonly: true}
	obj := ssot.Object{
		Name:         "Widget",
		NamePlural:   "Widgets",
		IdType:       "AccountObjectIdentifier",
		Describe:     ssot.OutputSpec{Columns: []*ssot.Attribute{&active}},
		DescribeRows: ssot.OutputSpec{Columns: []*ssot.Attribute{&property}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeDetails(ssot.SdkDescribeDetailsMethod()),
		))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	require.Equal(t, []*DescribeDetailsConfig{{MethodName: "DescribeDetails", DetailsType: "WidgetDetails", ParseFunc: "ParseWidgetDescribed"}}, got.DescribeDetailsMethods)
	require.Len(t, got.CustomMethods, 1)
	require.Equal(t, "DescribeDetails", got.CustomMethods[0].Name)
	require.Equal(t, []string{"*WidgetDetails", "error"}, got.CustomMethods[0].ReturnTypes)
	require.Equal(t, "AccountObjectIdentifier", got.CustomMethods[0].Parameters[0].Kind)
}

func Test_DescribeOperationFromObject_detailsFromIdentifier(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true}
	obj := func(idType string, opts ...ssot.SdkDescribeOption) ssot.Object {
		return ssot.Object{
			Name:       "Widget",
			NamePlural: "Widgets",
			IdType:     idType,
			Describe:   ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
			Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(opts...))},
		}
	}
	schema := obj("SchemaObjectIdentifier",
		ssot.SdkDescribePlainOnlyField("SchemaName", "string"),
		ssot.SdkDescribePlainOnlyField("DatabaseName", "string"),
		ssot.SdkDescribeDetailsFromIdentifier(),
	)
	got := NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").DescribeOperationFromObject(schema)
	require.Equal(t, &DescribeDetailsFromIdentifierConfig{
		DetailsType: "WidgetDetails",
		Fields:      []string{"DatabaseName", "SchemaName"},
	}, got.DescribeDetailsFromIdentifier)
	require.Empty(t, got.DescribeDetailsMethods)
	require.Len(t, got.CustomMethods, 1)
	require.Equal(t, "DescribeDetails", got.CustomMethods[0].Name)
	require.Equal(t, "DescribeDetails calls Describe and sets database and schema from the identifier. DESCRIBE does not return those columns.", got.CustomMethods[0].Doc)
	require.Equal(t, []string{"*WidgetDetails", "error"}, got.CustomMethods[0].ReturnTypes)
	require.Equal(t, "SchemaObjectIdentifier", got.CustomMethods[0].Parameters[0].Kind)
	preprocessDefinition(got, log.New(io.Discard, "", 0))
	require.False(t, got.Operations[0].DescribeMapping.HasManualConvert())
	require.Equal(t, NewShowObjectIDMethod("WidgetDetails", SchemaObjectIdentifier), got.Operations[0].ObjectIdMethod)

	var idOut strings.Builder
	require.NoError(t, ShowObjectIdMethodTemplate.Execute(&idOut, got.Operations[0].ObjectIdMethod))
	require.Contains(t, idOut.String(), "func (v *WidgetDetails) ID() SchemaObjectIdentifier {")
	require.Contains(t, idOut.String(), "return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name, )")

	var ifaceOut strings.Builder
	require.NoError(t, InterfaceTemplate.Execute(&ifaceOut, got))
	require.Contains(t, ifaceOut.String(), "DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*WidgetDetails, error)")

	var implOut strings.Builder
	require.NoError(t, ImplementationTemplate.Execute(&implOut, got))
	require.Contains(t, implOut.String(), `func (v *widgets) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*WidgetDetails, error) {
	details, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	details.DatabaseName = id.DatabaseName()
	details.SchemaName = id.SchemaName()
	return details, nil
}`)
	require.NotContains(t, implOut.String(), "ParseWidgetDescribed")
	require.NotContains(t, implOut.String(), "details.Id")
	require.NotContains(t, implOut.String(), "additionalConvert")

	database := obj("DatabaseObjectIdentifier",
		ssot.SdkDescribePlainOnlyField("DatabaseName", "string"),
		ssot.SdkDescribeDetailsFromIdentifier(),
	)
	got = NewInterface("Widgets", "Widget", "DatabaseObjectIdentifier").DescribeOperationFromObject(database)
	require.Equal(t, []string{"DatabaseName"}, got.DescribeDetailsFromIdentifier.Fields)
	require.Equal(t, "DescribeDetails calls Describe and sets database from the identifier. DESCRIBE does not return that column.", got.CustomMethods[0].Doc)
	preprocessDefinition(got, log.New(io.Discard, "", 0))
	require.False(t, got.Operations[0].DescribeMapping.HasManualConvert())
	require.Equal(t, NewShowObjectIDMethod("WidgetDetails", DatabaseObjectIdentifier), got.Operations[0].ObjectIdMethod)
	implOut.Reset()
	require.NoError(t, ImplementationTemplate.Execute(&implOut, got))
	require.Contains(t, implOut.String(), "details.DatabaseName = id.DatabaseName()")
	require.NotContains(t, implOut.String(), "SchemaName")
	require.NotContains(t, implOut.String(), "additionalConvert")

	schemaOnly := obj("SchemaObjectIdentifier",
		ssot.SdkDescribePlainOnlyField("SchemaName", "string"),
		ssot.SdkDescribeDetailsFromIdentifier(),
	)
	got = NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").DescribeOperationFromObject(schemaOnly)
	require.Equal(t, []string{"SchemaName"}, got.DescribeDetailsFromIdentifier.Fields)
	require.Equal(t, "DescribeDetails calls Describe and sets schema from the identifier. DESCRIBE does not return that column.", got.CustomMethods[0].Doc)
	preprocessDefinition(got, log.New(io.Discard, "", 0))
	require.Nil(t, got.Operations[0].ObjectIdMethod)

	require.PanicsWithValue(t, `object "Widget": SdkDescribeDetailsFromIdentifier needs IdType SchemaObjectIdentifier or DatabaseObjectIdentifier`, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(obj("AccountObjectIdentifier",
			ssot.SdkDescribePlainOnlyField("DatabaseName", "string"),
			ssot.SdkDescribeDetailsFromIdentifier(),
		))
	})
}

func Test_preprocessDefinition_describeIDFromColumns(t *testing.T) {
	iface := NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").
		DescribeOperationWithPairedStructs(
			DescriptionMappingKindSingleValue,
			"https://example.com",
			StructPair("widgetDescRow", "WidgetDetails").
				Text("name").
				Text("database_name").
				Text("schema_name"),
			NewQueryStruct("DescWidget").Describe().SQL("WIDGET").Name().WithValidation(ValidIdentifier, "name"),
		)
	preprocessDefinition(iface, log.New(io.Discard, "", 0))
	require.Equal(t, NewShowObjectIDMethod("WidgetDetails", SchemaObjectIdentifier), iface.Operations[0].ObjectIdMethod)

	optional := NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").
		DescribeOperationWithPairedStructs(
			DescriptionMappingKindSingleValue,
			"https://example.com",
			StructPair("widgetDescRow", "WidgetDetails").
				Text("name").
				OptionalText("database_name").
				Text("schema_name"),
			NewQueryStruct("DescWidget").Describe().SQL("WIDGET").Name().WithValidation(ValidIdentifier, "name"),
		)
	preprocessDefinition(optional, log.New(io.Discard, "", 0))
	require.Nil(t, optional.Operations[0].ObjectIdMethod)
}

func Test_preprocessDefinition_identifierPlainOnlySkipsAdditionalConvert(t *testing.T) {
	describe := func(fields *PairedStructs, fromIdentifier []string) *Interface {
		iface := NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").
			DescribeOperationWithPairedStructs(
				DescriptionMappingKindSingleValue,
				"https://example.com",
				fields,
				NewQueryStruct("DescWidget").Describe().SQL("WIDGET").Name().WithValidation(ValidIdentifier, "name"),
			)
		if fromIdentifier != nil {
			iface = iface.WithDescribeDetailsFromIdentifier("WidgetDetails", fromIdentifier)
		}
		preprocessDefinition(iface, log.New(io.Discard, "", 0))
		return iface
	}
	covered := describe(
		StructPair("widgetDescRow", "WidgetDetails").Text("name").PlainOnlyField("DatabaseName", "string").PlainOnlyField("SchemaName", "string"),
		[]string{"DatabaseName", "SchemaName"},
	)
	require.False(t, covered.Operations[0].DescribeMapping.HasManualConvert())
	require.Equal(t, NewShowObjectIDMethod("WidgetDetails", SchemaObjectIdentifier), covered.Operations[0].ObjectIdMethod)

	extra := describe(
		StructPair("widgetDescRow", "WidgetDetails").Text("name").PlainOnlyField("DatabaseName", "string").PlainOnlyField("Comment", "string"),
		[]string{"DatabaseName"},
	)
	require.True(t, extra.Operations[0].DescribeMapping.HasManualConvert())
	require.Nil(t, extra.Operations[0].ObjectIdMethod)

	manualColumn := describe(
		StructPair("widgetDescRow", "WidgetDetails").Text("comment", WithManualConvert()).PlainOnlyField("DatabaseName", "string"),
		[]string{"DatabaseName"},
	)
	require.True(t, manualColumn.Operations[0].DescribeMapping.HasManualConvert())

	withoutOption := describe(
		StructPair("widgetDescRow", "WidgetDetails").Text("name").PlainOnlyField("DatabaseName", "string"),
		nil,
	)
	require.True(t, withoutOption.Operations[0].DescribeMapping.HasManualConvert())
}

func Test_DescribeOperationFromObject_detailsParse(t *testing.T) {
	property := ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true}
	locations := ssot.Attribute{Name: "storage_locations", Type: ssot.AttrList, Readonly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	active := ssot.Attribute{Name: "active", Type: ssot.AttrString, Readonly: true}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, Readonly: true}
	obj := ssot.Object{
		Name:         "ExternalWidget",
		NamePlural:   "ExternalWidgets",
		IdType:       "AccountObjectIdentifier",
		Describe:     ssot.OutputSpec{Columns: []*ssot.Attribute{&locations, &active, &comment}},
		DescribeRows: ssot.OutputSpec{Columns: []*ssot.Attribute{&property}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeDetails(
				ssot.SdkDescribeDetailsNestedField(&locations),
				ssot.SdkDescribeDetailsParse(
					ssot.DescribeProperty(&comment, "COMMENT"),
					ssot.DescribeProperty(&active, "ACTIVE"),
					ssot.DescribePropertyPrefix(&locations, "LOCATION_"),
				),
			),
		))},
	}
	got := NewInterface("ExternalWidgets", "ExternalWidget", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	require.Equal(t, &DescribeParseConfig{
		RowsType:      "ExternalWidgetProperty",
		PropertyField: "Property",
		ValueField:    "PropertyValue",
		DetailsType:   "ExternalWidgetDetails",
		ErrorLabel:    "external widget",
		ValidateFunc:  "validateExternalWidgetDescribed",
		Properties: []DescribeParseProperty{
			{Property: "COMMENT", Field: "Comment"},
			{Property: "ACTIVE", Field: "Active"},
		},
		Prefixes: []DescribeParsePrefix{{
			Prefix:        "LOCATION_",
			Field:         "StorageLocations",
			Var:           "storageLocations",
			ElementType:   "ExternalWidgetStorageLocationDetails",
			ElementParser: "ParseExternalWidgetStorageLocationDescribed",
		}},
	}, got.DescribeParse)
}

func Test_Interface_WithDescribeDetailsMethod_manyMethods(t *testing.T) {
	i := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").
		WithDescribeDetails("WidgetDetails", "ParseWidgetDescribed").
		WithDescribeDetailsMethod("DescribeMiniDetails", "WidgetMiniDetails", "ParseWidgetMiniDescribed")

	require.Equal(t, []*DescribeDetailsConfig{
		{MethodName: "DescribeDetails", DetailsType: "WidgetDetails", ParseFunc: "ParseWidgetDescribed"},
		{MethodName: "DescribeMiniDetails", DetailsType: "WidgetMiniDetails", ParseFunc: "ParseWidgetMiniDescribed"},
	}, i.DescribeDetailsMethods)
	require.Equal(t, []string{"DescribeDetails", "DescribeMiniDetails"}, collections.Map(i.CustomMethods, func(m *CustomInterfaceMethod) string { return m.Name }))
}

func Test_DescribeOperationFromObject_detailsForVariants(t *testing.T) {
	property := ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true}
	active := ssot.Attribute{Name: "active", Type: ssot.AttrString, Readonly: true}
	mode := ssot.NewEnum("WidgetMode", "WidgetModes", "FAST", "SLOW")
	config := ssot.Attribute{
		Name: "config", Type: ssot.AttrObject, Required: true,
		Children: []ssot.Attribute{
			{Name: "uri", Type: ssot.AttrString, Required: true},
			{Name: "mode", Type: ssot.AttrEnum, Enum: mode},
			{Name: "scopes", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}},
			{Name: "fixed", SqlOnly: true, CreateSql: ssot.StaticSQL("fixed", "A = B")},
		},
	}
	auth := ssot.Attribute{Name: "o_auth_config", Type: ssot.AttrObject, Children: []ssot.Attribute{{Name: "token", Type: ssot.AttrString}}}
	base := ssot.Object{
		Name:         "Widget",
		NamePlural:   "Widgets",
		IdType:       "AccountObjectIdentifier",
		DescribeRows: ssot.OutputSpec{Columns: []*ssot.Attribute{&property}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeDetailsForVariants(),
		))},
	}
	mini := ssot.Object{
		Name:      "WidgetMini",
		IdType:    "AccountObjectIdentifier",
		VariantOf: &base,
		Describe:  ssot.OutputSpec{Columns: []*ssot.Attribute{&active}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeSql("WIDGET"),
			ssot.SdkDescribeDetails(),
		))},
	}
	maxi := ssot.Object{
		Name:      "WidgetMaxi",
		IdType:    "AccountObjectIdentifier",
		VariantOf: &base,
		Describe:  ssot.OutputSpec{Columns: []*ssot.Attribute{&active, &config, &auth}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeSql("WIDGET"),
			ssot.SdkDescribeDetails(),
		))},
	}
	other := ssot.Object{Name: "WidgetOther", IdType: "AccountObjectIdentifier", VariantOf: &base, Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{&active}}}
	orig := variantsOf
	t.Cleanup(func() { variantsOf = orig })
	variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{mini, other, maxi} }

	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(base)

	require.Equal(t, []*DescribeDetailsConfig{
		{MethodName: "DescribeMiniDetails", DetailsType: "WidgetMiniDetails", ParseFunc: "ParseWidgetMiniDescribed"},
		{MethodName: "DescribeMaxiDetails", DetailsType: "WidgetMaxiDetails", ParseFunc: "ParseWidgetMaxiDescribed"},
	}, got.DescribeDetailsMethods, "a variant without SdkDescribeDetails has no details")

	helpers := collections.Map(got.Operations[0].HelperStructs, func(f *Field) string { return f.Name })
	require.Equal(t, []string{"WidgetMiniDetails", "WidgetMaxiDetails", "WidgetMaxiConfigDetails", "WidgetMaxiOAuthConfigDetails"}, helpers[len(helpers)-4:])

	require.Equal(t,
		PlainStruct("WidgetMaxiDetails").
			Field("Id", "AccountObjectIdentifier").
			Field("Active", "string").
			Field("Config", "WidgetMaxiConfigDetails").
			Field("OAuthConfig", "*WidgetMaxiOAuthConfigDetails").
			IntoField(),
		DescribeDetailsPlain(maxi).IntoField(),
		"an object column is a struct of its Children, a pointer when it is not required",
	)

	nested := DescribeNestedPlains(maxi)
	require.Len(t, nested, 2)
	require.Equal(t,
		PlainStruct("WidgetMaxiConfigDetails").
			Field("Uri", "string").
			Field("Mode", "WidgetMode").
			Field("Scopes", "[]string").
			IntoField(),
		nested[0].IntoField(),
		"a Child that is only SQL is not a field",
	)

	t.Run("no variant with details", func(t *testing.T) {
		variantsOf = func(ssot.Object) []ssot.Object { return []ssot.Object{other} }
		require.Panics(t, func() {
			NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(base)
		})
	})
}

func Test_DescribeOperationFromObject_detailsParseParsedAndRows(t *testing.T) {
	property := ssot.Attribute{Name: "property", Type: ssot.AttrString, Readonly: true}
	value := ssot.Attribute{Name: "value", Type: ssot.AttrString, Readonly: true}
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true}
	mode := ssot.Attribute{Name: "mode", Type: ssot.AttrEnum, Enum: ssot.NewEnum("WidgetMode", "WidgetModes", "FAST"), Readonly: true}
	obj := ssot.Object{
		Name:         "Widget",
		NamePlural:   "Widgets",
		IdType:       "AccountObjectIdentifier",
		Describe:     ssot.OutputSpec{Columns: []*ssot.Attribute{&name, &mode}},
		DescribeRows: ssot.OutputSpec{Columns: []*ssot.Attribute{&property, &value}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkDescribe(
			ssot.SdkDescribeDetails(
				ssot.SdkDescribeDetailsRowsField("descriptions", "WidgetDescriptions"),
				ssot.SdkDescribeDetailsParseIgnoreUnknown(),
				ssot.SdkDescribeDetailsParse(
					ssot.DescribeProperty(&name, "NAME"),
					ssot.DescribePropertyParsed(&mode, "MODE"),
				),
			),
		))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").DescribeOperationFromObject(obj)
	require.Equal(t, &DescribeParseConfig{
		RowsType:      "WidgetProperty",
		PropertyField: "Property",
		ValueField:    "Value",
		DetailsType:   "WidgetDetails",
		ErrorLabel:    "widget",
		ValidateFunc:  "validateWidgetDescribed",
		IgnoreUnknown: true,
		RowsField:     "descriptions",
		RowsFieldType: "WidgetDescriptions",
		Properties: []DescribeParseProperty{
			{Property: "NAME", Field: "Name"},
			{Property: "MODE", Field: "Mode", Parser: "ParseWidgetModeDescribed"},
		},
	}, got.DescribeParse)
}
