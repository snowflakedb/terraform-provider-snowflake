package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

// The list has the shape of the External Volume storage location: parenthesized elements with a name and exactly one of the
// variants, a variant with a nested struct, and a variant with a static field. Each variant refers to attributes that
// already have their CREATE clause.
type widgetLocationAttributes struct {
	list           ssot.Attribute
	kind           *ssot.Enum
	encryptionKind *ssot.Enum
}

func newWidgetLocationAttributes() widgetLocationAttributes {
	kind := ssot.NewEnum("WidgetKind", "WidgetKinds", "A", "B")
	encryptionKind := ssot.NewEnum("WidgetEncryptionKind", "WidgetEncryptionKinds", "NONE", "SSE")
	name := ssot.Attribute{Name: "location_name", Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier", Required: true, CreateSql: ssot.StringParameter(ssot.Sql("NAME"), ssot.Required())}
	kindAttr := ssot.Attribute{Name: "kind", Type: ssot.AttrString, CreateSql: ssot.StringParameter(ssot.EnumOf(kind), ssot.SingleQuotes(), ssot.Required())}
	baseUrl := ssot.Attribute{Name: "base_url", Type: ssot.AttrString, Required: true, CreateSql: ssot.Clause()}
	privatelink := ssot.Attribute{Name: "use_privatelink", Type: ssot.AttrBool, CreateSql: ssot.Clause()}
	encryptionType := ssot.Attribute{Name: "encryption_type", Type: ssot.AttrString, CreateSql: ssot.Clause(ssot.Sql("TYPE"), ssot.FieldName("EncryptionType"), ssot.EnumOf(encryptionKind), ssot.SingleQuotes(), ssot.Required())}
	keyId := ssot.Attribute{Name: "encryption_key_id", Type: ssot.AttrString, CreateSql: ssot.Clause(ssot.Sql("KEY_ID"))}
	secondKind := ssot.Attribute{Name: "second_kind", SqlOnly: true, CreateSql: ssot.StaticSQL("secondKind", "KIND = 'SECOND'")}
	encryption := ssot.Attribute{
		Name:      "first_encryption",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("ENCRYPTION ="), ssot.Parentheses(), ssot.FieldName("Encryption")),
		Children:  []ssot.Attribute{encryptionType, keyId},
	}
	return widgetLocationAttributes{
		kind:           kind,
		encryptionKind: encryptionKind,
		list: ssot.Attribute{
			Name:      "location",
			Type:      ssot.AttrList,
			Required:  true,
			CreateSql: ssot.ListOfObjects(ssot.Sql("LOCATIONS"), ssot.Parentheses()),
			Elem: &ssot.Attribute{
				Type: ssot.AttrObject,
				CreateSql: ssot.ObjectParameter(ssot.Parentheses(), ssot.WithVariants(
					ssot.Variant("First", &kindAttr, &baseUrl, &privatelink, &encryption),
					ssot.Variant("Second", &secondKind, &baseUrl),
				)),
				Children: []ssot.Attribute{name, kindAttr, baseUrl, privatelink, encryptionType, keyId},
			},
		},
	}
}

func Test_createQueryStruct_objectListWithVariants(t *testing.T) {
	widget := newWidgetLocationAttributes()
	locations := widget.list
	obj := ssot.Object{
		Name:       "Widget",
		NamePlural: "Widgets",
		IdType:     "AccountObjectIdentifier",
		Attributes: []ssot.Attribute{locations},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&locations}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate())},
	}
	ssot.ValidateObject(obj)

	got := createQueryStructFromConfig(obj, obj.SdkConfig())

	want := NewQueryStruct("CreateWidget").
		Create().
		OrReplace().
		SQL("WIDGET").
		IfNotExists().
		Name().
		ListQueryStructField(
			"Locations",
			NewQueryStruct("WidgetLocationItem").
				QueryStructField(
					"WidgetLocation",
					NewQueryStruct("WidgetLocation").
						TextAssignment("NAME", ParameterOptions().SingleQuotes().Required()).
						OptionalQueryStructField(
							"FirstLocationParams",
							NewQueryStruct("FirstLocationParams").
								EnumAssignmentWithFieldName("KIND", widget.kind, ParameterOptions().SingleQuotes().Required(), "Kind").
								TextAssignment("BASE_URL", ParameterOptions().SingleQuotes().Required()).
								OptionalBooleanAssignment("USE_PRIVATELINK", ParameterOptions()).
								OptionalQueryStructField(
									"Encryption",
									NewQueryStruct("WidgetFirstEncryption").
										EnumAssignmentWithFieldName("TYPE", widget.encryptionKind, ParameterOptions().SingleQuotes().Required(), "EncryptionType").
										OptionalTextAssignment("KEY_ID", ParameterOptions().SingleQuotes()),
									ListOptions().Parentheses().NoComma().SQL("ENCRYPTION ="),
								),
							ListOptions().NoComma(),
						).
						OptionalQueryStructField(
							"SecondLocationParams",
							NewQueryStruct("SecondLocationParams").
								SQLWithCustomFieldName("secondKind", "KIND = 'SECOND'").
								TextAssignment("BASE_URL", ParameterOptions().SingleQuotes().Required()),
							ListOptions().NoComma(),
						).
						WithValidation(ExactlyOneValueSet, "FirstLocationParams", "SecondLocationParams"),
					ListOptions().Parentheses().NoComma().Required(),
				),
			ParameterOptions().Parentheses().Required().SQL("LOCATIONS"),
		).
		WithValidation(ValidateValueSet, "Locations").
		WithValidation(ValidIdentifier, "name").
		WithValidation(ConflictingFields, "OrReplace", "IfNotExists")

	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObjectList_buildsFreshStructs(t *testing.T) {
	locations := newWidgetLocationAttributes().list
	first := queryStructFromObjectList("Widget", locations)
	second := queryStructFromObjectList("Widget", locations)
	require.NotSame(t, first, second)
	require.NotPanics(t, func() {
		NewQueryStruct("A").QueryStructField("Locations", first, ListOptions())
		NewQueryStruct("B").QueryStructField("Locations", second, ListOptions())
	})
}

func Test_queryStructFromObjectList_withoutParenthesesIsNotWrapped(t *testing.T) {
	locations := newWidgetLocationAttributes().list
	elem := *locations.Elem
	elem.CreateSql = ssot.ObjectParameter(ssot.WithVariants(elem.CreateSql.Variants...))
	locations.Elem = &elem
	require.Equal(t, "WidgetLocation", queryStructFromObjectList("Widget", locations).IntoField().Name)
}

func Test_referencedEnums_walksObjectElemAndVariants(t *testing.T) {
	widget := newWidgetLocationAttributes()
	obj := ssot.Object{Name: "Widget", Attributes: []ssot.Attribute{widget.list}}
	require.ElementsMatch(t, []*ssot.Enum{widget.kind, widget.encryptionKind}, referencedEnums(obj))
}

func Test_referencedEnums_includesAdditionalEnums(t *testing.T) {
	widget := newWidgetLocationAttributes()
	extra := ssot.NewEnum("WidgetExtraKind", "WidgetExtraKinds", "NONE")
	obj := ssot.Object{Name: "Widget", Attributes: []ssot.Attribute{widget.list}, AdditionalEnums: []*ssot.Enum{extra, widget.kind}}
	require.ElementsMatch(t, []*ssot.Enum{widget.kind, widget.encryptionKind, extra}, referencedEnums(obj))
}

func Test_appendSqlSequence_singleObject(t *testing.T) {
	uri := ssot.Attribute{Name: "catalog_uri", Type: ssot.AttrString, Required: true}
	name := ssot.Attribute{Name: "catalog_name", Type: ssot.AttrString}
	config := ssot.Attribute{
		Name:      "rest_config",
		Type:      ssot.AttrObject,
		Required:  true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("REST_CONFIG ="), ssot.Parentheses()),
		Children:  []ssot.Attribute{uri, name},
	}
	got := NewQueryStruct("CreateWidget").appendSqlSequence([]*ssot.Attribute{&config}, attributeSqlCreate)
	want := NewQueryStruct("CreateWidget").QueryStructField(
		"RestConfig",
		NewQueryStruct("RestConfig").
			TextAssignment("CATALOG_URI", ParameterOptions().SingleQuotes().Required()).
			OptionalTextAssignment("CATALOG_NAME", ParameterOptions().SingleQuotes()),
		ListOptions().SQL("REST_CONFIG =").Parentheses().NoComma().Required(),
	)
	require.Equal(t, want.fields, got.fields)
}

func Test_appendSqlSequence_singleObjectWithStaticField(t *testing.T) {
	typeOauth := ssot.Attribute{Name: "type_oauth", SqlOnly: true, CreateSql: ssot.StaticSQL("restAuthType", "TYPE = OAUTH")}
	secret := ssot.Attribute{Name: "client_secret", Type: ssot.AttrString, Required: true}
	auth := ssot.Attribute{
		Name:      "rest_authentication",
		Type:      ssot.AttrObject,
		CreateSql: ssot.ObjectParameter(ssot.Sql("REST_AUTHENTICATION ="), ssot.Parentheses()),
		Children:  []ssot.Attribute{typeOauth, secret},
	}
	got := NewQueryStruct("CreateWidget").appendSqlSequence([]*ssot.Attribute{&auth}, attributeSqlCreate)
	want := NewQueryStruct("CreateWidget").OptionalQueryStructField(
		"RestAuthentication",
		NewQueryStruct("RestAuthentication").
			SQLWithCustomFieldName("restAuthType", "TYPE = OAUTH").
			TextAssignment("CLIENT_SECRET", ParameterOptions().SingleQuotes().Required()),
		ListOptions().SQL("REST_AUTHENTICATION =").Parentheses().NoComma(),
	)
	require.Equal(t, want.fields, got.fields)
}

func Test_createQueryStruct_exactlyOneOf(t *testing.T) {
	typeOauth := ssot.Attribute{Name: "type_oauth", SqlOnly: true, CreateSql: ssot.StaticSQL("restAuthType", "TYPE = OAUTH")}
	secret := ssot.Attribute{Name: "client_secret", Type: ssot.AttrString, Required: true}
	token := ssot.Attribute{Name: "bearer_token", Type: ssot.AttrString, Required: true}
	oauth := ssot.Attribute{
		Name:      "oauth_authentication",
		Type:      ssot.AttrObject,
		CreateSql: ssot.ObjectParameter(ssot.Sql("AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("OAuthAuthentication")),
		Children:  []ssot.Attribute{typeOauth, secret},
	}
	bearer := ssot.Attribute{
		Name:      "bearer_authentication",
		Type:      ssot.AttrObject,
		CreateSql: ssot.ObjectParameter(ssot.Sql("AUTHENTICATION ="), ssot.Parentheses()),
		Children:  []ssot.Attribute{token},
	}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	newObject := func(group []*ssot.Attribute, clauses ...*ssot.Attribute) ssot.Object {
		return ssot.Object{
			Name:         "Widget",
			NamePlural:   "Widgets",
			IdType:       "AccountObjectIdentifier",
			Attributes:   []ssot.Attribute{oauth, bearer, comment},
			ExactlyOneOf: [][]*ssot.Attribute{group},
			Create:       ssot.CreateSpec{Clauses: clauses},
			Generators:   []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate())},
		}
	}

	t.Run("all members in the statement", func(t *testing.T) {
		obj := newObject([]*ssot.Attribute{&oauth, &bearer}, &oauth, &bearer, &comment)
		ssot.ValidateObject(obj)
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
		want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
			NewQueryStruct("CreateWidget").
				Create().
				OrReplace().
				SQL("WIDGET").
				IfNotExists().
				Name().
				OptionalQueryStructField(
					"OAuthAuthentication",
					NewQueryStruct("WidgetOauthAuthentication").
						SQLWithCustomFieldName("restAuthType", "TYPE = OAUTH").
						TextAssignment("CLIENT_SECRET", ParameterOptions().SingleQuotes().Required()),
					ListOptions().SQL("AUTHENTICATION =").Parentheses().NoComma(),
				).
				OptionalQueryStructField(
					"BearerAuthentication",
					NewQueryStruct("WidgetBearerAuthentication").
						TextAssignment("BEARER_TOKEN", ParameterOptions().SingleQuotes().Required()),
					ListOptions().SQL("AUTHENTICATION =").Parentheses().NoComma(),
				).
				OptionalTextAssignment("COMMENT", ParameterOptions().SingleQuotes()).
				WithValidation(ValidIdentifier, "name").
				WithValidation(ConflictingFields, "OrReplace", "IfNotExists").
				WithValidation(ExactlyOneValueSet, "OAuthAuthentication", "BearerAuthentication"),
		)
		requireSameOperations(t, want, got)
	})

	t.Run("no member in the statement", func(t *testing.T) {
		obj := newObject([]*ssot.Attribute{&oauth, &bearer}, &comment)
		got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
		for _, validation := range got.Operations[0].OptsField.Validations {
			require.NotEqual(t, ExactlyOneValueSet, validation.Type)
		}
	})

	t.Run("some members in the statement", func(t *testing.T) {
		obj := newObject([]*ssot.Attribute{&oauth, &bearer}, &oauth, &comment)
		require.Panics(t, func() {
			NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
		})
	})
}

func Test_alterQueryStruct_setBlocksAllowEmptyAndConflicts(t *testing.T) {
	secret := ssot.Attribute{Name: "client_secret", Type: ssot.AttrString, Required: true}
	token := ssot.Attribute{Name: "bearer_token", Type: ssot.AttrString, Required: true}
	setOauth := ssot.Attribute{
		Name:     "set_oauth",
		Type:     ssot.AttrObject,
		SqlOnly:  true,
		SetSql:   ssot.ObjectParameter(ssot.Sql("AUTHENTICATION ="), ssot.Parentheses(), ssot.FieldName("SetOAuth")),
		Children: []ssot.Attribute{secret},
	}
	setBearer := ssot.Attribute{
		Name:     "set_bearer",
		Type:     ssot.AttrObject,
		SqlOnly:  true,
		SetSql:   ssot.ObjectParameter(ssot.Sql("AUTHENTICATION ="), ssot.Parentheses()),
		Children: []ssot.Attribute{token},
	}
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString, SetSql: ssot.StringParameter(ssot.AllowEmpty())}
	obj := ssot.Object{
		Name:       "Widget",
		NamePlural: "Widgets",
		IdType:     "AccountObjectIdentifier",
		Set:        []*ssot.Attribute{&setOauth, &setBearer, &comment},
		Alter: ssot.AlterSpec{
			Clauses:              []ssot.AlterClause{ssot.AlterSet()},
			SetConflictingFields: [][]*ssot.Attribute{{&setOauth, &setBearer}},
		},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkAlter())},
	}
	ssot.ValidateObject(obj)

	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").AlterOperationFromObject(obj)
	auth := ListOptions().SQL("AUTHENTICATION =").Parentheses().NoComma()
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
					OptionalQueryStructField("SetOAuth", NewQueryStruct("WidgetSetOauth").TextAssignment("CLIENT_SECRET", ParameterOptions().SingleQuotes().Required()), auth).
					OptionalQueryStructField("SetBearer", NewQueryStruct("WidgetSetBearer").TextAssignment("BEARER_TOKEN", ParameterOptions().SingleQuotes().Required()), auth).
					OptionalAssignment("COMMENT", "StringAllowEmpty", ParameterOptions()).
					WithValidation(ConflictingFields, "SetOAuth", "SetBearer").
					WithValidation(AtLeastOneValueSet, "SetOAuth", "SetBearer", "Comment"),
				KeywordOptions().SQL("SET"),
			).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ExactlyOneValueSet, "Set"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateOperationFromObject_singleBlockWithEnumSetsAndAtLeastOne(t *testing.T) {
	family := ssot.NewEnum("WidgetFamily", "WidgetFamilies", "SMALL")
	policy := ssot.Attribute{
		Name:         "policy",
		Type:         ssot.AttrObject,
		CreateSql:    ssot.ObjectParameter(ssot.Sql("POLICY ="), ssot.Parentheses()),
		ChildrenRule: ssot.ChildrenAtLeastOne,
		Children: []ssot.Attribute{
			{Name: "limit", Type: ssot.AttrInt},
			{Name: "families", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family}},
			{Name: "names", Type: ssot.AttrSet, Elem: &ssot.Attribute{Type: ssot.AttrString}, CreateSql: ssot.ListParameter("StringListItemWrapper")},
		},
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{policy},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&policy}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(
			ssot.SdkCreateWithoutOrReplace(),
			ssot.SdkCreateWithoutIfNotExists(),
		))},
	}
	nested := NewQueryStruct("WidgetPolicy").
		OptionalNumberAssignment("LIMIT", ParameterOptions()).
		ListAssignment("FAMILIES", "WidgetFamilyListItem", ParameterOptions().Parentheses()).
		ListAssignment("NAMES", "StringListItemWrapper", ParameterOptions().Parentheses()).
		WithValidation(AtLeastOneValueSet, "Limit", "Families", "Names")
	item := NewQueryStruct("WidgetFamilyListItem").Enum("Value", family, KeywordOptions().SingleQuotes().Required())
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField("Policy", nested, ListOptions().SQL("POLICY =").Parentheses().NoComma()).
			WithValidation(ValidIdentifier, "name"),
		item,
	)
	requireSameOperations(t, want, got)
}

func Test_CreateOperationFromObject_enumKeywordAndContainedChild(t *testing.T) {
	driver := ssot.NewEnum("WidgetDriver", "WidgetDrivers", "JDBC")
	version := ssot.Attribute{Name: "minimum_version", Type: ssot.AttrString, Required: true}
	params := ssot.Attribute{
		Name:      "entry_params",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("="), ssot.Parentheses(), ssot.FieldName("Params"), ssot.Required()),
		Children:  []ssot.Attribute{version},
	}
	entries := ssot.Attribute{
		Name:      "entries",
		Type:      ssot.AttrSet,
		CreateSql: ssot.ListOfObjects(ssot.Sql("ENTRIES"), ssot.Parentheses()),
		Elem: &ssot.Attribute{
			Type:      ssot.AttrObject,
			CreateSql: ssot.ObjectParameter(),
			Children: []ssot.Attribute{
				{Name: "driver", Type: ssot.AttrEnum, Required: true, Enum: driver, CreateSql: ssot.EnumKeyword(ssot.NoQuotes(), ssot.Required())},
				version,
				params,
			},
		},
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{entries},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&entries}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(
			ssot.SdkCreateWithoutOrReplace(),
			ssot.SdkCreateWithoutIfNotExists(),
		))},
	}
	nestedParams := NewQueryStruct("WidgetEntryParams").
		TextAssignment("MINIMUM_VERSION", ParameterOptions().SingleQuotes().Required())
	element := NewQueryStruct("WidgetEntries").
		Enum("Driver", driver, KeywordOptions().NoQuotes().Required()).
		QueryStructField("Params", nestedParams, ListOptions().SQL("=").Parentheses().NoComma().Required())
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			ListQueryStructField("Entries", element, ParameterOptions().Parentheses().SQL("ENTRIES")).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateOperationFromObject_exactlyOneKeywordAndIdentifierList(t *testing.T) {
	names := ssot.Attribute{
		Name:      "names",
		Type:      ssot.AttrSet,
		Elem:      &ssot.Attribute{Type: ssot.AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		CreateSql: ssot.IdentifierListKeyword(),
	}
	option := ssot.Attribute{
		Name:      "names_option",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("NAMES ="), ssot.FieldName("Names")),
		Children: []ssot.Attribute{
			{Name: "names_all", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("('ALL')"), ssot.FieldName("All"))},
			names,
		},
		ChildrenRule: ssot.ChildrenExactlyOne,
	}
	obj := ssot.Object{
		Name:       "Widget",
		Attributes: []ssot.Attribute{names},
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&option}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(
			ssot.SdkCreateWithoutOrReplace(),
			ssot.SdkCreateWithoutIfNotExists(),
		))},
	}
	nested := NewQueryStruct("WidgetNamesOption").
		OptionalSQLWithCustomFieldName("All", "('ALL')").
		UnnamedList("Names", "AccountObjectIdentifier", KeywordOptions().Parentheses()).
		WithValidation(ValidIdentifier, "Names").
		WithValidation(ExactlyOneValueSet, "All", "Names")
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField("Names", nested, ListOptions().SQL("NAMES =").NoComma()).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}

func Test_CreateOperationFromObject_objectAsParameter(t *testing.T) {
	option := ssot.Attribute{
		Name:      "option",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("OPTION"), ssot.FieldName("Option"), ssot.AsParameter()),
		Children: []ssot.Attribute{
			{Name: "all", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("('ALL')"), ssot.FieldName("All"))},
		},
	}
	obj := ssot.Object{
		Name:       "Widget",
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&option}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(ssot.SdkCreateWithoutOrReplace(), ssot.SdkCreateWithoutIfNotExists()))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			OptionalQueryStructField("Option", NewQueryStruct("WidgetOption").OptionalSQLWithCustomFieldName("All", "('ALL')"), ParameterOptions().SQL("OPTION")).
			WithValidation(ValidIdentifier, "name"),
	)
	requireSameOperations(t, want, got)
}
