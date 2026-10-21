package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

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
						WithAdditionalValidations().
						WithValidation(AtLeastOneValueSet, "Items", "Count", "Label"),
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

func Test_CreateAlterFromObject_parameters(t *testing.T) {
	orig := parametersForObject
	t.Cleanup(func() { parametersForObject = orig })
	params := []ssot.ParameterDef{
		{SqlName: "DATA_RETENTION_TIME_IN_DAYS", Kind: string(ssot.AttrInt)},
		{SqlName: "EXTERNAL_VOLUME", Kind: "AccountObjectIdentifier"},
	}
	parametersForObject = func(ssot.Object) []ssot.ParameterDef { return params }

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
