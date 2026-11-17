package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validateCreateSql(t *testing.T) {
	top := attributePlace{}
	nested := attributePlace{nested: true}

	require.True(t, Attribute{Name: "label", Type: AttrString, CreateSql: StringParameter()}.validateCreateSql(top))
	require.False(t, Attribute{Name: "label", Type: AttrString, CreateSql: StringParameter(NoQuotes())}.validateCreateSql(top))

	// Clause() with no options is the type default. Single quotes are already the string default.
	require.True(t, Attribute{Name: "label", Type: AttrString, CreateSql: Clause()}.validateCreateSql(top))
	require.True(t, Attribute{Name: "label", Type: AttrString, CreateSql: Clause(SingleQuotes())}.validateCreateSql(top))
	require.False(t, Attribute{Name: "family", Type: AttrEnum, CreateSql: Clause(SingleQuotes())}.validateCreateSql(top))

	// A nested CREATE clause that only restates the default can be omitted.
	// The same clause next to a SET clause is what keeps the attribute on CREATE.
	require.True(t, Attribute{Name: "label", Type: AttrString, CreateSql: StringParameter()}.validateCreateSql(nested))
	both := Attribute{Name: "label", Type: AttrString, CreateSql: StringParameter(), SetSql: StringParameter(NoQuotes())}
	require.False(t, both.validateCreateSql(nested))

	require.False(t, Attribute{Name: "label", Type: AttrString}.validateCreateSql(top))
	require.False(t, Attribute{Name: "body", Type: AttrObject, CreateSql: ObjectParameter()}.validateCreateSql(top))
}

func Test_validateUnsetSql(t *testing.T) {
	top := attributePlace{}
	nested := attributePlace{nested: true}
	flag := Attribute{Name: "flag", Type: AttrBool, UnsetSql: BoolKeyword()}

	// A top-level UNSET clause that restates the bool keyword is the omission default.
	// On a nested attribute the same clause is membership in the UNSET projection.
	require.True(t, flag.validateUnsetSql(top))
	require.False(t, flag.validateUnsetSql(nested))
}

func Test_validateClauseSql(t *testing.T) {
	comment := Attribute{Name: "comment", Type: AttrString, CreateSql: StringParameter(Sql("COMMENT"), NoQuotes())}
	require.True(t, comment.validateClauseSql(ClauseCreate, attributePlace{}))

	renamed := Attribute{Name: "comment", Type: AttrString, CreateSql: StringParameter(Sql("NOTE"))}
	require.False(t, renamed.validateClauseSql(ClauseCreate, attributePlace{}))

	// The whole clause is already the type default, so the SQL text is not reported on its own.
	whole := Attribute{Name: "comment", Type: AttrString, CreateSql: StringParameter()}
	require.False(t, whole.validateClauseSql(ClauseCreate, attributePlace{}))

	// Static SQL and a list parameter with FieldName print clause.SQL, not the uppercase name.
	static := Attribute{Name: "set", CreateSql: StaticSQL("set", "SET")}
	require.False(t, static.validateClauseSql(ClauseCreate, attributePlace{}))
	listed := Attribute{Name: "on", Type: AttrList, Elem: &Attribute{Type: AttrString}, CreateSql: ListParameter("Column", Sql("ON"), FieldName("On"))}
	require.False(t, listed.validateClauseSql(ClauseCreate, attributePlace{}))
}

func Test_validateClauseField(t *testing.T) {
	object := Attribute{Name: "comment", Type: AttrObject, CreateSql: ObjectParameter(FieldName("Comment"))}
	require.True(t, object.validateClauseField(ClauseCreate))

	renamed := Attribute{Name: "comment", Type: AttrObject, CreateSql: ObjectParameter(FieldName("Note"))}
	require.False(t, renamed.validateClauseField(ClauseCreate))

	// A parameter derives its empty field name from the SQL text, not the attribute name.
	parameter := Attribute{Name: "comment", Type: AttrString, CreateSql: StringParameter(FieldName("Comment"))}
	require.False(t, parameter.validateClauseField(ClauseCreate))
}

func Test_validateUpdateField(t *testing.T) {
	place := attributePlace{stem: "Widget", checkUpdate: true, setOnly: true}
	body := Attribute{Name: "body", Type: AttrString, UpdateMapping: StringAttributeUpdateSetOnlyNotEmpty("Body")}
	require.True(t, body.validateUpdateField(place))

	renamed := Attribute{Name: "body", Type: AttrString, UpdateMapping: StringAttributeUpdateSetOnlyNotEmpty("SetBody")}
	require.False(t, renamed.validateUpdateField(place))
}

func Test_redundantObjectDefaults(t *testing.T) {
	base := Object{Name: "Warehouse"}
	require.Empty(t, Object{Name: "Warehouse"}.redundantObjectDefaults())

	details := Object{Name: "Warehouse", DescribeDetailsName: "WarehouseDetails"}
	require.Equal(t, []string{"Warehouse DescribeDetailsName restates the omitted default"}, details.redundantObjectDefaults())

	plain := Object{Name: "Warehouse", ShowPlainName: "Warehouse"}
	require.Equal(t, []string{"Warehouse ShowPlainName restates the omitted default"}, plain.redundantObjectDefaults())
	require.Empty(t, Object{Name: "Warehouse", ShowPlainName: "WarehouseRegular"}.redundantObjectDefaults())

	stem := Object{Name: "WarehouseAdaptive", SdkCreateRequestObjectName: "WarehouseAdaptive"}
	require.Equal(t, []string{"WarehouseAdaptive SdkCreateRequestObjectName restates the omitted default"}, stem.redundantObjectDefaults())
	require.Empty(t, Object{Name: "WarehouseAdaptive", SdkCreateRequestObjectName: "AdaptiveWarehouse"}.redundantObjectDefaults())

	create := Object{Name: "WarehouseAdaptive", VariantOf: &base, VariantEntityCreateName: "WAREHOUSE"}
	require.Equal(t, []string{"WarehouseAdaptive VariantEntityCreateName restates the omitted default"}, create.redundantObjectDefaults())
	alter := Object{Name: "WarehouseAdaptive", VariantOf: &base, VariantEntityAlterName: "WAREHOUSE"}
	require.Equal(t, []string{"WarehouseAdaptive VariantEntityAlterName restates the omitted default"}, alter.redundantObjectDefaults())
	require.Empty(t, Object{Name: "WarehouseInteractive", VariantOf: &base, VariantEntityCreateName: "INTERACTIVE WAREHOUSE"}.redundantObjectDefaults())
}

func Test_validateShowName(t *testing.T) {
	require.True(t, Attribute{Name: "comment", ShowName: "comment"}.validateShowName())
	require.False(t, Attribute{Name: "for_application", ShowName: "application"}.validateShowName())
}

func Test_validateSchemaName(t *testing.T) {
	require.True(t, Attribute{Name: "comment", Type: AttrString, SchemaName: "comment"}.validateSchemaName())
	require.False(t, Attribute{Name: "null?", Type: AttrBool, SchemaName: "is_nullable"}.validateSchemaName())
}

func Test_validateShowEnum(t *testing.T) {
	family := NewEnum("Family", "Families", "SMALL")
	other := NewEnum("Other", "Others", "SMALL")

	require.True(t, Attribute{Name: "family", Type: AttrEnum, Enum: family, ShowEnum: family}.validateShowEnum())
	require.False(t, Attribute{Name: "family", Type: AttrEnum, Enum: family, ShowEnum: other}.validateShowEnum())
}

func Test_validateExplicitDefaults_message(t *testing.T) {
	attr := Attribute{Name: "comment", Type: AttrString, ShowName: "comment", CreateSql: StringParameter()}
	require.Equal(t, []string{
		"Alert.comment CREATE restates the omitted default",
		"Alert.comment ShowName restates the omitted default",
	}, attr.validateExplicitDefaults("Alert.comment", attributePlace{}))
}
