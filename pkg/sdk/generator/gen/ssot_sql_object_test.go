package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
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

	got := createQueryStructFromConfig(obj, obj.SdkConfig(), nil)

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

func Test_appendSqlSequence_wrappedString(t *testing.T) {
	condition := ssot.Attribute{
		Name:     "condition",
		Type:     ssot.AttrString,
		Required: true,
		CreateSql: ssot.WrapString(
			ssot.SQLBefore("ifExists", "IF (EXISTS("),
			ssot.SQLAfter("closeIfExists", "))"),
		),
	}
	got := NewQueryStruct("CreateAlert")
	got.owner = "Alert"
	got = got.appendSqlSequence([]*ssot.Attribute{&condition}, attributeSqlCreate)
	want := NewQueryStruct("CreateAlert").QueryStructField(
		"Condition",
		NewQueryStruct("AlertCondition").
			SQLWithCustomFieldName("ifExists", "IF (EXISTS(").
			Text("Condition", KeywordOptions().NoQuotes().Required()).
			SQLWithCustomFieldName("closeIfExists", "))"),
		ListOptions().NoComma().Required(),
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
		Alter: ssot.AlterSpec{
			Clauses: []ssot.AlterClause{ssot.AlterSet(
				ssot.AlterGroup("Oauth", &setOauth),
				ssot.AlterGroup(ssot.MainAlterGroup, &setBearer, &comment),
			)},
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
					WithConflictingGroups([]string{"SetOAuth"}, []string{"SetBearer", "Comment"}).
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

func Test_CreateOperationFromObject_keywordStringList(t *testing.T) {
	condition := ssot.Attribute{
		Name:      "condition",
		Type:      ssot.AttrObject,
		Required:  true,
		CreateSql: ssot.ObjectParameter(ssot.Sql("IF"), ssot.Parentheses(), ssot.Required(), ssot.FieldName("condition")),
		Children: []ssot.Attribute{{
			Name:      "condition",
			Type:      ssot.AttrList,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: ssot.StringListKeyword(ssot.Parentheses(), ssot.NoComma(), ssot.Sql("EXISTS"), ssot.FieldName("Condition")),
		}},
	}
	action := ssot.Attribute{
		Name:      "then_action",
		Type:      ssot.AttrString,
		Required:  true,
		CreateSql: ssot.PlainParameter(ssot.Sql("THEN"), ssot.NoEquals(), ssot.Required(), ssot.FieldName("action")),
	}
	obj := ssot.Object{
		Name:       "Widget",
		Create:     ssot.CreateSpec{Clauses: []*ssot.Attribute{&condition, &action}},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate(ssot.SdkCreateWithoutOrReplace(), ssot.SdkCreateWithoutIfNotExists()))},
	}
	got := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperationFromObject(obj)
	want := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			SQL("WIDGET").
			Name().
			QueryStructField("condition", NewQueryStruct("WidgetCondition").
				UnnamedList("Condition", "string", KeywordOptions().Parentheses().NoComma().SQL("EXISTS")),
				ListOptions().NoComma().Parentheses().SQL("IF").Required()).
			AssignmentWithFieldName("THEN", "string", ParameterOptions().NoEquals().Required(), "action").
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

func Test_CreateOperationFromObject_stringListAndPreexistingStructs(t *testing.T) {
	scope := ssot.NewEnum("WidgetScope", "WidgetScopes", "GLOBAL", "LOCAL")
	scopeAttr := ssot.Attribute{Name: "scope", Type: ssot.AttrEnum, Enum: scope, SqlOnly: true, CreateSql: ssot.EnumKeyword()}
	widgetColumns := ssot.Attribute{
		Name:      "create_columns",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.StructName("CreateWidgetColumns")),
		Children: []ssot.Attribute{
			{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.Required())},
		},
	}
	columns := ssot.Attribute{Name: "columns_and_constraints", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&widgetColumns, ssot.Parentheses(), ssot.Required())}
	clusterBy := ssot.Attribute{Name: "cluster_by", Type: ssot.AttrList, SqlOnly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}, CreateSql: ssot.StringListKeyword(ssot.Parentheses(), ssot.Sql("CLUSTER BY"))}
	widgetFormat := ssot.Attribute{
		Name:      "file_format",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.StructName("WidgetFileFormat")),
		Children: []ssot.Attribute{
			{Name: "format_name", Type: ssot.AttrString, SqlOnly: true},
		},
	}
	format := ssot.Attribute{Name: "stage_file_format", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&widgetFormat, ssot.Parentheses(), ssot.NoComma(), ssot.Sql("STAGE_FILE_FORMAT ="))}
	policy := ssot.Attribute{
		Name:      "row_access_policy",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			{
				Name:           "name",
				Type:           ssot.AttrIdentifier,
				IdentifierType: "SchemaObjectIdentifier",
				Required:       true,
				SqlOnly:        true,
				CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("ROW ACCESS POLICY"), ssot.FieldName("Name")),
			},
			{
				Name:      "on",
				Type:      ssot.AttrList,
				Required:  true,
				SqlOnly:   true,
				Elem:      &ssot.Attribute{Type: ssot.AttrString},
				CreateSql: ssot.StringListKeyword(ssot.Parentheses(), ssot.Sql("ON")),
			},
		},
	}
	obj := ssot.Object{
		Name:       "Widget",
		NamePlural: "Widgets",
		IdType:     "SchemaObjectIdentifier",
		Create: ssot.CreateSpec{
			KeywordsBefore: []*ssot.Attribute{&scopeAttr},
			Clauses:        []*ssot.Attribute{&columns, &clusterBy, &format, &policy},
		},
		Generators: []ssot.GeneratorOption{ssot.WithSdk(ssot.SdkCreate())},
	}
	ssot.ValidateObject(obj)

	newColumns := func() *QueryStruct {
		return NewQueryStruct("CreateWidgetColumns").Text("Name", KeywordOptions().Required())
	}
	newFormat := func() *QueryStruct {
		return NewQueryStruct("WidgetFileFormat").OptionalTextAssignment("FORMAT_NAME", ParameterOptions().SingleQuotes())
	}
	unrelated := NewQueryStruct("WidgetUnrelated").Text("Name", KeywordOptions())

	got := NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").CreateOperationFromObject(obj, unrelated)

	want := NewInterface("Widgets", "Widget", "SchemaObjectIdentifier").CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-widget",
		NewQueryStruct("CreateWidget").
			Create().
			OrReplace().
			OptionalEnum("Scope", scope, KeywordOptions()).
			SQL("WIDGET").
			IfNotExists().
			Name().
			QueryStructField("ColumnsAndConstraints", newColumns(), ListOptions().Parentheses().Required()).
			NamedList("CLUSTER BY", "string", KeywordOptions().Parentheses().SQL("CLUSTER BY")).
			OptionalQueryStructField("StageFileFormat", newFormat(), ListOptions().Parentheses().NoComma().SQL("STAGE_FILE_FORMAT =")).
			OptionalQueryStructField("RowAccessPolicy", NewQueryStruct("WidgetRowAccessPolicy").
				Identifier("Name", "SchemaObjectIdentifier", IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
				WithValidation(ValidIdentifier, "Name").
				NamedList("ON", "string", KeywordOptions().Parentheses().SQL("ON").Required()),
				ListOptions().NoComma()).
			WithValidation(ValidIdentifier, "name").
			WithValidation(ConflictingFields, "OrReplace", "IfNotExists"),
		unrelated,
	)
	requireSameOperations(t, want, got)
}

func Test_queryStructFromObject_tableClusteringAction(t *testing.T) {
	attr := ssotdefs.TableClusteringAction
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	want := NewQueryStruct("TableClusteringAction").
		PredefinedQueryStructField("ClusterBy", "[]string", KeywordOptions().Parentheses().SQL("CLUSTER BY")).
		OptionalQueryStructField("Recluster",
			NewQueryStruct("TableReclusterAction").
				OptionalNumberAssignment("MAX_SIZE", ParameterOptions()).
				OptionalAssignmentWithFieldName("WHERE", "*string", ParameterOptions().NoEquals(), "Condition"),
			ListOptions().NoComma().SQL("RECLUSTER")).
		OptionalQueryStructField("ChangeReclusterState",
			NewQueryStruct("TableReclusterChangeState").
				OptionalEnum("State", ssotdefs.ReclusterState, KeywordOptions()).
				SQL("RECLUSTER"),
			ListOptions().NoComma()).
		OptionalSQL("DROP CLUSTERING KEY").
		WithValidation(ExactlyOneValueSet, "ClusterBy", "Recluster", "ChangeReclusterState", "DropClusteringKey")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableSearchOptimizationAction(t *testing.T) {
	attr := ssotdefs.TableSearchOptimizationAction
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	searchMethodWithTarget := func() *QueryStruct {
		return NewQueryStruct("TableSearchMethodWithTarget").
			Enum("Method", ssotdefs.TableSearchMethod, KeywordOptions().Required()).
			QueryStructField("Args",
				NewQueryStruct("TableSearchMethodArgs").
					PredefinedQueryStructField("Targets", "[]string", KeywordOptions()).
					OptionalTextAssignment("ANALYZER", ParameterOptions().ArrowEquals().SingleQuotes()),
				ListOptions().Parentheses().Required())
	}
	want := NewQueryStruct("TableSearchOptimizationAction").
		OptionalQueryStructField("Add",
			NewQueryStruct("TableAddSearchOptimization").
				SQL("ADD SEARCH OPTIMIZATION").
				ListQueryStructField("On", searchMethodWithTarget(), KeywordOptions().SQL("ON")),
			ListOptions().NoComma()).
		OptionalQueryStructField("Drop",
			NewQueryStruct("TableDropSearchOptimization").
				SQL("DROP SEARCH OPTIMIZATION").
				ListQueryStructField("On",
					NewQueryStruct("TableDropSearchOptimizationOn").
						OptionalQueryStructField("SearchMethodWithTarget", searchMethodWithTarget(), ListOptions().NoComma()).
						OptionalText("ColumnName", KeywordOptions()).
						OptionalText("ExpressionId", KeywordOptions()).
						WithValidation(ExactlyOneValueSet, "SearchMethodWithTarget", "ColumnName", "ExpressionId"),
					KeywordOptions().SQL("ON")),
			ListOptions().NoComma()).
		WithValidation(ExactlyOneValueSet, "Add", "Drop")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableAggregationPolicy(t *testing.T) {
	attr := ssotdefs.TableCreateAttrs.AggregationPolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	want := NewQueryStruct("TableAggregationPolicy").
		Identifier("AggregationPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
		WithValidation(ValidIdentifier, "AggregationPolicy").
		ListAssignment("ENTITY KEY", "Column", ParameterOptions().NoEquals().Parentheses())
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableJoinPolicy(t *testing.T) {
	attr := ssotdefs.TableCreateAttrs.JoinPolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	want := NewQueryStruct("TableJoinPolicy").
		Identifier("JoinPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("JOIN POLICY").Required()).
		WithValidation(ValidIdentifier, "JoinPolicy").
		ListAssignment("ALLOWED JOIN KEYS", "Column", ParameterOptions().NoEquals().Parentheses())
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableSetAggregationPolicy(t *testing.T) {
	attr := ssotdefs.TableAggregationPolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttributeUse("Table", attr, ssot.ClauseSet)
	want := NewQueryStruct("TableSetAggregationPolicy").
		SQL("SET").
		Identifier("AggregationPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("AGGREGATION POLICY").Required()).
		WithValidation(ValidIdentifier, "AggregationPolicy").
		ListAssignment("ENTITY KEY", "Column", ParameterOptions().NoEquals().Parentheses()).
		OptionalSQL("FORCE")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableUnsetAggregationPolicy(t *testing.T) {
	attr := ssotdefs.TableAggregationPolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttributeUse("Table", attr, ssot.ClauseUnset)
	want := NewQueryStruct("TableUnsetAggregationPolicy").
		SQL("UNSET AGGREGATION POLICY")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableSetJoinPolicy(t *testing.T) {
	attr := ssotdefs.TableJoinPolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttributeUse("Table", attr, ssot.ClauseSet)
	want := NewQueryStruct("TableSetJoinPolicy").
		SQL("SET").
		Identifier("JoinPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("JOIN POLICY").Required()).
		WithValidation(ValidIdentifier, "JoinPolicy").
		OptionalSQL("FORCE")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableUnsetJoinPolicy(t *testing.T) {
	attr := ssotdefs.TableJoinPolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttributeUse("Table", attr, ssot.ClauseUnset)
	want := NewQueryStruct("TableUnsetJoinPolicy").
		SQL("UNSET JOIN POLICY")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_alterSetQueryStruct_table(t *testing.T) {
	require.NotPanics(t, func() { ssot.ValidateObject(ssotdefs.Table) })
	got := AlterSetQueryStruct(ssotdefs.Table, ssot.SdkGenerator{})
	want := NewQueryStruct("TableSet").
		OptionalBooleanAssignment("ENABLE_SCHEMA_EVOLUTION", ParameterOptions()).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", ParameterOptions()).
		OptionalBooleanAssignment("CHANGE_TRACKING", ParameterOptions()).
		OptionalAssignment("DEFAULT_DDL_COLLATION", "StringAllowEmpty", ParameterOptions()).
		OptionalAssignment("ICEBERG_DEFAULT_DDL_COLLATION", "StringAllowEmpty", ParameterOptions()).
		OptionalBooleanAssignment("ERROR_LOGGING", ParameterOptions()).
		PredefinedQueryStructField("Contact", "[]TableContact", KeywordOptions().SQL("CONTACT")).
		OptionalComment().
		OptionalBooleanAssignment("ROW_TIMESTAMP", ParameterOptions()).
		WithValidation(AtLeastOneValueSet, "EnableSchemaEvolution", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "DefaultDdlCollation", "IcebergDefaultDdlCollation", "ErrorLogging", "Contact", "Comment", "RowTimestamp")
	require.Equal(t, want.IntoField(), got.IntoField())

	unset := AlterUnsetQueryStruct(ssotdefs.Table)
	wantUnset := NewQueryStruct("TableUnset").
		OptionalAssignment("CONTACT", "string", ParameterOptions().NoEquals()).
		OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
		OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
		OptionalSQL("CHANGE_TRACKING").
		OptionalSQL("DEFAULT_DDL_COLLATION").
		OptionalSQL("ICEBERG_DEFAULT_DDL_COLLATION").
		OptionalSQL("ENABLE_SCHEMA_EVOLUTION").
		OptionalSQL("ERROR_LOGGING").
		OptionalSQL("COMMENT").
		OptionalSQL("ROW_TIMESTAMP").
		WithValidation(AtLeastOneValueSet, "Contact", "DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "DefaultDdlCollation", "IcebergDefaultDdlCollation", "EnableSchemaEvolution", "ErrorLogging", "Comment", "RowTimestamp").
		WithConflictingGroups([]string{"Contact"}, []string{"DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ChangeTracking", "DefaultDdlCollation", "IcebergDefaultDdlCollation", "EnableSchemaEvolution", "ErrorLogging", "Comment", "RowTimestamp"})
	require.Equal(t, wantUnset.IntoField(), unset.IntoField())
}

func Test_queryStructFromObject_tableStorageLifecyclePolicy(t *testing.T) {
	attr := ssotdefs.TableStorageLifecyclePolicy
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	want := NewQueryStruct("TableStorageLifecyclePolicy").
		Identifier("StorageLifecyclePolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("ADD STORAGE LIFECYCLE POLICY").Required()).
		WithValidation(ValidIdentifier, "StorageLifecyclePolicy").
		ListAssignment("ON", "Column", ParameterOptions().NoEquals().Parentheses().Required()).
		WithValidation(ValidateValueSet, "On")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableRowAccessPolicy(t *testing.T) {
	for _, attr := range []ssot.Attribute{
		ssotdefs.TableAddRowAccessPolicy,
		ssotdefs.TableDropRowAccessPolicy,
		ssotdefs.TableDropAndAddRowAccessPolicy,
	} {
		require.NotPanics(t, func() {
			ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
		})
	}
	add := QueryStructFromAttribute("Table", ssotdefs.TableAddRowAccessPolicy)
	require.Equal(t, NewQueryStruct("TableAddRowAccessPolicy").
		SQLWithCustomFieldName("add", "ADD").
		Identifier("RowAccessPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
		WithValidation(ValidIdentifier, "RowAccessPolicy").
		ListAssignmentWithFieldName("ON", "Column", ParameterOptions().NoEquals().Parentheses().Required(), "On").
		WithValidation(ValidateValueSet, "On").
		IntoField(), add.IntoField())
	drop := QueryStructFromAttribute("Table", ssotdefs.TableDropRowAccessPolicy)
	require.Equal(t, NewQueryStruct("TableDropRowAccessPolicy").
		SQLWithCustomFieldName("drop", "DROP").
		Identifier("RowAccessPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
		WithValidation(ValidIdentifier, "RowAccessPolicy").
		IntoField(), drop.IntoField())
	got := QueryStructFromAttribute("Table", ssotdefs.TableDropAndAddRowAccessPolicy)
	want := NewQueryStruct("TableDropAndAddRowAccessPolicy").
		QueryStructField("Drop", drop, ListOptions().NoComma().Required()).
		QueryStructField("Add", add, ListOptions().NoComma().Required())
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableColumn(t *testing.T) {
	attr := ssotdefs.TableColumn
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)

	// Table emits the structs the column statements embed, and each one carries its standalone toOpts.
	virtual := QueryStructFromAttribute("Table", ssotdefs.TableColumnVirtual)
	inlineConstraint := QueryStructFromAttribute("Table", ssotdefs.TableColumnInlineConstraint)
	maskingPolicy := QueryStructFromAttribute("Table", ssotdefs.TableColumnMaskingPolicy)
	projectionPolicy := QueryStructFromAttribute("Table", ssotdefs.TableColumnProjectionPolicy)
	renameColumn := QueryStructFromAttribute("Table", ssotdefs.TableRenameColumnAction)
	setMasking := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnMaskingPolicy, ssot.ClauseSet)
	unsetMasking := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnMaskingPolicy, ssot.ClauseUnset)
	setProjection := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnProjectionPolicy, ssot.ClauseSet)
	unsetProjection := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnProjectionPolicy, ssot.ClauseUnset)
	setTags := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnTags, ssot.ClauseSet)
	unsetTags := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnTags, ssot.ClauseUnset)
	dropColumns := QueryStructFromAttribute("Table", ssotdefs.TableDropColumnAction)

	add := NewQueryStruct("TableColumnAddAction").
		SQL("COLUMN").
		IfNotExists().
		Text("Name", KeywordOptions().Required().DoubleQuotes()).
		PredefinedQueryStructField("ColumnType", "datatypes.DataType", ParameterOptions().NoQuotes().NoEquals().Required()).
		OptionalQueryStructField("Virtual", virtual, ListOptions().NoComma()).
		OptionalTextAssignment("COLLATE", ParameterOptions().NoEquals().SingleQuotes()).
		OptionalQueryStructField("DefaultValue",
			NewQueryStruct("ColumnDefaultValue").
				OptionalAssignmentWithFieldName("DEFAULT", "*string", ParameterOptions().NoEquals(), "Expression").
				OptionalQueryStructField("Identity",
					NewQueryStruct("ColumnIdentity").
						NumberAssignment("START", ParameterOptions().NoQuotes().NoEquals().Required()).
						NumberAssignment("INCREMENT", ParameterOptions().NoQuotes().NoEquals().Required()).
						OptionalSQL("ORDER").
						OptionalSQL("NOORDER").
						WithValidation(ConflictingFields, "Order", "Noorder"),
					ListOptions().NoComma().SQL("IDENTITY")).
				WithValidation(ExactlyOneValueSet, "Expression", "Identity"),
			ListOptions().NoComma()).
		OptionalSQL("NOT NULL").
		OptionalQueryStructField("InlineConstraint", inlineConstraint, ListOptions().NoComma()).
		OptionalQueryStructField("MaskingPolicy", maskingPolicy, ListOptions().NoComma()).
		OptionalQueryStructField("ProjectionPolicy", projectionPolicy, ListOptions().NoComma()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", ParameterOptions().NoEquals().SingleQuotes())
	alter := NewQueryStruct("TableColumnAlterAction").
		SQL("COLUMN").
		Text("Name", KeywordOptions().Required().DoubleQuotes()).
		OptionalSQL("DROP DEFAULT").
		PredefinedQueryStructField("SetDefault", "*SequenceName", ParameterOptions().NoEquals().SQL("SET DEFAULT")).
		OptionalQueryStructField("NotNullConstraint",
			NewQueryStruct("TableColumnNotNullConstraint").
				OptionalSQLWithCustomFieldName("Set", "SET NOT NULL").
				OptionalSQLWithCustomFieldName("Drop", "DROP NOT NULL"),
			ListOptions().NoComma()).
		OptionalQueryStructField("DataType",
			NewQueryStruct("TableColumnAlterDataType").
				PredefinedQueryStructField("DataType", "datatypes.DataType", ParameterOptions().NoQuotes().NoEquals().Required()).
				OptionalTextAssignment("COLLATE", ParameterOptions().NoEquals().SingleQuotes()),
			ListOptions().NoComma().SQL("SET DATA TYPE")).
		OptionalTextAssignment("COMMENT", ParameterOptions().NoEquals().SingleQuotes()).
		OptionalSQL("UNSET COMMENT").
		// The name is Required, so it is not one of the exactly-one actions.
		WithValidation(ExactlyOneValueSet, "DropDefault", "SetDefault", "NotNullConstraint", "DataType", "Comment", "UnsetComment")
	want := NewQueryStruct("TableColumnAction").
		OptionalQueryStructField("Add", add, ListOptions().NoComma().SQL("ADD")).
		OptionalQueryStructField("Rename", renameColumn, ListOptions().NoComma()).
		ListQueryStructField("Alter", alter, KeywordOptions().SQL("ALTER")).
		OptionalQueryStructField("SetMaskingPolicy", setMasking, ListOptions().NoComma()).
		OptionalQueryStructField("UnsetMaskingPolicy", unsetMasking, ListOptions().NoComma()).
		OptionalQueryStructField("SetProjectionPolicy", setProjection, ListOptions().NoComma()).
		OptionalQueryStructField("UnsetProjectionPolicy", unsetProjection, ListOptions().NoComma()).
		OptionalQueryStructField("SetTags", setTags, ListOptions().NoComma()).
		OptionalQueryStructField("UnsetTags", unsetTags, ListOptions().NoComma()).
		OptionalQueryStructField("DropColumns", dropColumns, ListOptions().NoComma()).
		WithValidation(ExactlyOneValueSet, "Add", "Rename", "Alter", "SetMaskingPolicy", "UnsetMaskingPolicy", "SetProjectionPolicy", "UnsetProjectionPolicy", "SetTags", "UnsetTags", "DropColumns")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_sqlOnlyConflictsWith(t *testing.T) {
	order := ssot.Attribute{Name: "order", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("ORDER"))}
	noorder := ssot.Attribute{Name: "noorder", Type: ssot.AttrBool, SqlOnly: true, ConflictsWith: []ssot.Attribute{order}, CreateSql: ssot.BoolKeyword(ssot.Sql("NOORDER"))}
	// Both sides name each other, and the pair is still added once, in the order of the children.
	order.ConflictsWith = []ssot.Attribute{noorder}
	attr := ssot.Attribute{Name: "identity", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.ObjectParameter(), Children: []ssot.Attribute{noorder, order}}
	got := QueryStructFromAttribute("Widget", attr)
	want := NewQueryStruct("WidgetIdentity").
		OptionalSQL("NOORDER").
		OptionalSQL("ORDER").
		WithValidation(ConflictingFields, "Noorder", "Order")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableColumnActionStructs(t *testing.T) {
	for _, attr := range []ssot.Attribute{
		ssotdefs.TableRenameColumnAction,
		ssotdefs.TableDropColumnAction,
	} {
		require.NotPanics(t, func() {
			ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
		})
	}
	// The set and unset column actions have no CREATE clause. They are SET and UNSET projections of one attribute.
	for _, attr := range []ssot.Attribute{
		ssotdefs.TableAlterColumnMaskingPolicy,
		ssotdefs.TableAlterColumnTags,
	} {
		require.NotPanics(t, func() {
			ssot.ValidateObject(ssot.Object{Name: "Table", Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{ssot.AlterSet(&attr), ssot.AlterUnset(&attr)}}})
		})
	}
	rename := QueryStructFromAttribute("Table", ssotdefs.TableRenameColumnAction)
	require.Equal(t, NewQueryStruct("TableRenameColumnAction").
		SQL("RENAME COLUMN").
		Text("OldName", KeywordOptions().Required().DoubleQuotes()).
		AssignmentWithFieldName("TO", "string", ParameterOptions().NoEquals().DoubleQuotes().Required(), "NewName").
		WithSharedToOpts().IntoField(), rename.IntoField())
	setMasking := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnMaskingPolicy, ssot.ClauseSet)
	require.Equal(t, NewQueryStruct("TableSetColumnMaskingPolicy").
		SQL("ALTER COLUMN").
		Text("Name", KeywordOptions().Required().DoubleQuotes()).
		SQL("SET").
		Identifier("MaskingPolicy", "SchemaObjectIdentifier", IdentifierOptions().SQL("MASKING POLICY").Required()).
		ListAssignment("USING", "Column", ParameterOptions().NoEquals().Parentheses()).
		OptionalSQL("FORCE").
		WithValidation(ValidIdentifier, "MaskingPolicy").
		WithSharedToOpts().IntoField(), setMasking.IntoField())
	unsetMasking := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnMaskingPolicy, ssot.ClauseUnset)
	require.Equal(t, NewQueryStruct("TableUnsetColumnMaskingPolicy").
		SQL("ALTER COLUMN").
		Text("Name", KeywordOptions().Required().DoubleQuotes()).
		SQL("UNSET").
		SQL("MASKING POLICY").
		WithSharedToOpts().IntoField(), unsetMasking.IntoField())
	setTags := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnTags, ssot.ClauseSet)
	require.Equal(t, NewQueryStruct("TableSetColumnTags").
		SQL("ALTER COLUMN").
		Text("Name", KeywordOptions().Required().DoubleQuotes()).
		SetTags().
		WithSharedToOpts().IntoField(), setTags.IntoField())
	unsetTags := QueryStructFromAttributeUse("Table", ssotdefs.TableAlterColumnTags, ssot.ClauseUnset)
	require.Equal(t, NewQueryStruct("TableUnsetColumnTags").
		SQL("ALTER COLUMN").
		Text("Name", KeywordOptions().Required().DoubleQuotes()).
		UnsetTags().
		WithSharedToOpts().IntoField(), unsetTags.IntoField())
	drop := QueryStructFromAttribute("Table", ssotdefs.TableDropColumnAction)
	require.Equal(t, NewQueryStruct("TableDropColumnAction").
		SQL("DROP COLUMN").
		OptionalSQL("IF EXISTS").
		PredefinedQueryStructField("Columns", "[]Column", ParameterOptions().NoEquals().Required()).
		WithSharedToOpts().IntoField(), drop.IntoField())
}

func Test_queryStructFromObject_tableConstraints(t *testing.T) {
	for _, tc := range []struct {
		name    string
		wrapper ssot.Attribute
		kinds   []ssot.Attribute
	}{
		{"out of line", ssotdefs.TableOutOfLineConstraint, []ssot.Attribute{ssotdefs.TableOutOfLineUniquePK, ssotdefs.TableOutOfLineFK, ssotdefs.TableOutOfLineCH}},
		{"inline", ssotdefs.TableColumnInlineConstraint, []ssot.Attribute{ssotdefs.TableColumnInlineUniquePK, ssotdefs.TableColumnInlineFK, ssotdefs.TableColumnInlineCH}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&tc.wrapper}}})
			})
			got := QueryStructFromAttribute("Table", tc.wrapper)
			want := NewQueryStruct(tc.wrapper.CreateSql.StructName).WithSharedToOpts()
			for _, kind := range tc.kinds {
				want = want.OptionalQueryStructField(kind.CreateSql.Field, QueryStructFromAttribute("Table", kind), ListOptions().NoComma())
			}
			want = want.WithValidation(ExactlyOneValueSet, "UniquePK", "FK", "CH")
			require.Equal(t, want.IntoField(), got.IntoField())
		})
	}
}

func Test_queryStructFromObject_tableConstraintAction(t *testing.T) {
	attr := ssotdefs.TableConstraintAction
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	rename := NewQueryStruct("TableConstraintRenameAction").
		Text("OldName", KeywordOptions().DoubleQuotes().Required()).
		AssignmentWithFieldName("TO", "string", ParameterOptions().DoubleQuotes().NoEquals().Required(), "NewName")
	alter := NewQueryStruct("TableConstraintAlterAction").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", ParameterOptions().DoubleQuotes().NoEquals(), "ConstraintName").
		OptionalSQL("PRIMARY KEY").
		OptionalSQL("UNIQUE").
		OptionalSQL("FOREIGN KEY").
		PredefinedQueryStructField("Columns", "[]string", ParameterOptions().NoEquals().Parentheses()).
		OptionalSQL("ENFORCED").
		OptionalSQL("NOT ENFORCED").
		OptionalSQL("VALIDATE").
		OptionalSQL("NOVALIDATE").
		OptionalSQL("RELY").
		OptionalSQL("NORELY").
		WithValidation(ExactlyOneValueSet, "ConstraintName", "PrimaryKey", "Unique", "ForeignKey")
	drop := NewQueryStruct("TableConstraintDropAction").
		OptionalAssignmentWithFieldName("CONSTRAINT", "*string", ParameterOptions().DoubleQuotes().NoEquals(), "ConstraintName").
		OptionalSQL("PRIMARY KEY").
		OptionalSQL("UNIQUE").
		OptionalSQL("FOREIGN KEY").
		PredefinedQueryStructField("Columns", "[]string", ParameterOptions().NoEquals().Parentheses()).
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(ExactlyOneValueSet, "ConstraintName", "PrimaryKey", "Unique", "ForeignKey")
	want := NewQueryStruct("TableConstraintAction").
		OptionalQueryStructField("Add", QueryStructFromAttribute("Table", ssotdefs.TableOutOfLineConstraint), ListOptions().NoComma().SQL("ADD")).
		OptionalQueryStructField("Rename", rename, ListOptions().NoComma().SQL("RENAME CONSTRAINT")).
		OptionalQueryStructField("Alter", alter, ListOptions().NoComma().SQL("ALTER")).
		OptionalQueryStructField("Drop", drop, ListOptions().NoComma().SQL("DROP")).
		WithValidation(ExactlyOneValueSet, "Add", "Rename", "Alter", "Drop")
	require.Equal(t, want.IntoField(), got.IntoField())
}

func Test_queryStructFromObject_tableExternalTableAction(t *testing.T) {
	attr := ssotdefs.TableExternalTableAction
	require.NotPanics(t, func() {
		ssot.ValidateObject(ssot.Object{Name: "Table", Create: ssot.CreateSpec{Clauses: []*ssot.Attribute{&attr}}})
	})
	got := QueryStructFromAttribute("Table", attr)
	add := NewQueryStruct("TableExternalTableColumnAddAction").
		SQL("ADD COLUMN").
		IfNotExists().
		Text("Name", KeywordOptions().Required()).
		PredefinedQueryStructField("ColumnType", "datatypes.DataType", ParameterOptions().NoQuotes().NoEquals().Required()).
		ListAssignmentWithFieldName("AS", "string", ParameterOptions().NoEquals().Parentheses(), "Expression").
		OptionalTextAssignment("COMMENT", ParameterOptions().NoEquals().SingleQuotes())
	want := NewQueryStruct("TableExternalTableAction").
		OptionalQueryStructField("Add", add, ListOptions().NoComma())
	require.Equal(t, want.IntoField(), got.IntoField())
}
