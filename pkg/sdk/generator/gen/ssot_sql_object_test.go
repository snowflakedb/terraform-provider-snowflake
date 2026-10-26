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
