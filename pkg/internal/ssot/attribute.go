package ssot

// AttributeType is the semantic type of an object attribute, not the Terraform schema type.
type AttributeType string

const (
	AttrString     AttributeType = "string"
	AttrInt        AttributeType = "int"
	AttrBool       AttributeType = "bool"
	AttrFloat      AttributeType = "float64"
	AttrIdentifier AttributeType = "identifier"
	AttrEnum       AttributeType = "enum"
	AttrList       AttributeType = "list"
	AttrSet        AttributeType = "set"
	AttrTime       AttributeType = "time"
	// AttrObject is an object with Children. Alone it is a single block: one struct on the SDK side, and for Terraform a
	// list of at most one element, which the resource attribute rules derive. It is also the Elem of a list or set of
	// objects, and a SqlOnly nested struct that is not a Terraform field.
	AttrObject AttributeType = "object"
)

// ObjectElem is the Elem of a list or set of objects: an AttrObject whose Children are the fields of each element.
func ObjectElem(children []Attribute) *Attribute {
	return &Attribute{Type: AttrObject, Children: children}
}

// Attribute is the written intent for one object attribute.
type Attribute struct {
	// --- Identity and type: what the attribute is.

	// Name is the snake_case key of the attribute: the Terraform name, and the SHOW/DESCRIBE column key unless ShowName is set.
	Name string
	Type AttributeType
	// IdentifierType is the SDK identifier type without the sdk prefix, such as AccountObjectIdentifier. Required for AttrIdentifier.
	IdentifierType string
	Enum           *Enum
	// ShowEnum is the SHOW and CREATE/SET enum when it differs from Enum.
	// Empty uses Enum. Resource schema and model builders keep Enum. DESCRIBE uses DescribeEnum or Enum.
	ShowEnum *Enum
	// Elem is the element of a list or set: a scalar attribute, or an AttrObject whose Children are the fields of a nested
	// block. A field of an Elem object may carry its own CREATE clause, which is how SDK query structs are generated.
	Elem *Attribute
	// Children are the fields of an AttrObject: a single block, an Elem, or a SqlOnly nested struct. They are valid only
	// there: a list or set of objects declares its fields on an Elem object (ObjectElem), never on the list itself.
	// A single block sets no Min or Max: Terraform gets at most one item, and one item when it is Required.
	// A SqlOnly child is fixed text of the SQL struct, such as TYPE = OAUTH: a StaticSQL clause, valid when the object has an
	// object clause. It is in NestedSqlFields and not in NestedFields, so it is not a Terraform field.
	// A nested field may itself be a block. Generated create parses an optional list block; see the README.
	// Resource alter generation skips a nested attribute and every field under it.
	Children []Attribute
	// Deprecated is the schema deprecation message.
	Deprecated string

	// --- Snowflake constraints: facts about the attribute that Snowflake defines. The Terraform schema reflects them, but they
	// are not Terraform-specific.

	Required bool
	ForceNew bool
	// CreateOnly marks an optional bool that is applied only on create.
	CreateOnly bool
	Sensitive  bool
	// Min and Max are inclusive bounds. An int bounds its value. A list or set bounds its number of items. Nil means
	// unset. A required AttrList deduces MinItems 1 when Min is nil, so omit Min on a required list unless setting a
	// higher bound (Min of 0 or 1 is rejected).
	Min *int
	Max *int
	// ServerDefault marks a describe-only field that Snowflake fills in when the configuration omits it.
	// The generator drift-tracks that field and suppresses a diff against describe_output.
	// A SHOW column does not need this: every optional scalar SHOW column is drift-tracked.
	// Leave it unset when the omitted value is empty, as catalog_namespace is.
	ServerDefault bool
	// Validations are rules over the Children of an object: a single block, an Elem, or a SqlOnly nested struct. Each names
	// Children that are fields of the SQL struct, not fixed texts, tags, or Required. See Validation and the README.
	Validations []Validation

	// --- Snowflake SQL: how the attribute is written in CREATE, SET, and UNSET.

	// CreateSql, SetSql, and UnsetSql are the clause the attribute adds to that statement. Zero leaves it out.
	CreateSql SqlClause
	SetSql    SqlClause
	UnsetSql  SqlClause
	// SqlOnly keeps the attribute on a SQL sequence and out of schema, resource, model, and assertion generation.
	// ResourceAttributes drops these. A static CREATE value sets CreateSql to StaticSQL.
	SqlOnly bool
	// Tags marks the optional TAG clause of an SDK query struct. It is not an attribute: use the Tags var in a clause
	// list, never in Attributes.
	Tags bool
	// TagForm selects SET TAG or UNSET TAG when Tags is set. Empty emits optional TAG.
	TagForm TagForm
	// Parameters marks the spread of the parameter catalog into a CREATE, SET, or UNSET clause list. It is not an
	// attribute: use the Parameters var in a clause list, never in Attributes.
	Parameters bool

	// --- Snowflake SHOW and DESCRIBE: how Snowflake reports the attribute back.

	// Readonly marks a SHOW or DESCRIBE column that Snowflake reports and the resource never writes.
	// A Readonly column is not in Attributes, and every other column must be in Attributes: validation enforces both.
	Readonly bool
	// ShowName is the SHOW/DESCRIBE column key when it differs from Name.
	ShowName string
	// OptionalInShow makes the SHOW/DESCRIBE SDK field a pointer, nil when Snowflake reports no value.
	// The generated schema ToSchema then guards against nil. Independent of OptionalInDb.
	OptionalInShow bool
	// ShowString makes the SHOW column a string. The resource, CREATE, and SET keep Type.
	// A string attribute leaves it unset.
	ShowString bool
	// DescribeString makes the DESCRIBE column a string. The resource, CREATE, SET, and SHOW keep Type.
	// A string attribute leaves it unset. It is not combined with DescribeEnum.
	DescribeString bool
	// DescribeEnum is the DESCRIBE enum when it differs from Enum.
	// Empty uses Enum. SHOW and CREATE/SET keep ShowEnum or Enum. The resource schema keeps Enum.
	// Leave it unset when it is Enum. It is not combined with DescribeString.
	DescribeEnum *Enum
	// StoredAsText means the db column is text while the plain kind stays the semantic Type, such as a bool read from Y.
	StoredAsText bool
	// BoolTrueValue, with StoredAsText, is the text that converts to true, such as Y, YES, or ON.
	// Empty uses strconv.ParseBool.
	BoolTrueValue string

	// --- Terraform: how the attribute is represented in the resource schema and its outputs.

	Computed bool
	// NativeBool keeps an optional bool as Terraform TypeBool. Optional bools otherwise use boolean-string.
	NativeBool   bool
	DiffSuppress DiffSuppress
	// SchemaName is the Terraform key of a SHOW/DESCRIBE column when it differs from the snake_case of the Go field,
	// such as is_nullable for the null? column. It does not change the Go field, the db tag, or OutputKey.
	SchemaName string
	// SkipInSchema omits a SHOW/DESCRIBE column from the generated public schema only.
	// The column stays on the SDK struct and in assertions.
	SkipInSchema bool

	// --- Docs: the generated description.

	// Description is the Terraform schema description.
	Description string
	// AdditionalDescription is appended after the list of allowed values of an enum, or of a collection of enums.
	// It has no effect on other types.
	AdditionalDescription string
	// RawDescription keeps Description as written for an optional bool. Without it, the generator appends the boolean
	// string suffix: Available options are "true" or "false", and the "default" value means the Snowflake default.
	RawDescription bool
	// IncludeAliasedInTheDescription lists enum aliases next to the canonical values in the description.
	IncludeAliasedInTheDescription bool
	// RelatedResource is the SSOT object whose Terraform resource is linked from the description.
	// The generator emits resources.<Name>. Nil means no link.
	RelatedResource *Object

	// --- Provider SDK mapping: how the SDK stores and converts SHOW and DESCRIBE values.

	// OptionalInDb marks the database row column nullable (sql.NullString and the like).
	// Used by the SDK paired-struct db and plain kinds. Independent of OptionalInShow.
	OptionalInDb bool
	// SdkString marks an optional AttrObject that is one nested block in Terraform and one string in SQL, SHOW, and DESCRIBE,
	// such as the schedule of a backup policy. The resource, the model, and the resource assertions read the Children as a
	// block. CREATE, SET, UNSET, the SHOW and DESCRIBE columns, and their schemas and assertions read it as a string, see
	// SdkType. The generated create (SkipCreate) and update skip it: ext code converts between the block and the string.
	// Set OptionalInShow or OptionalInDb for the string column.
	SdkString bool
	// CustomParser names the SDK function that converts the db value to the plain field, such as a CSV column to a slice.
	CustomParser string
	// ManualConvert leaves the column out of the generated paired-struct convert(). The SDK ext additionalConvert fills it.
	ManualConvert bool
	// PlainOnlyKind is a field on the plain struct only. It is not a database column.
	PlainOnlyKind string
	// PairedPlainKind overrides the plain kind of a required SHOW or DESCRIBE column.
	// The database column stays a string. Optional columns use the type mapping.
	// Go names come from the column name, with characters other than letters, digits, and underscores removed.
	PairedPlainKind string

	// --- Provider generation: how the generated create and update treat the attribute.

	// SkipCreate leaves the attribute out of generated create.
	SkipCreate bool
	// SkipUpdate leaves the attribute out of generated update.
	SkipUpdate bool
	// SetAfterCreate applies a SkipCreate attribute with ALTER after the object is created, through
	// WithAfterCreate(Generated()). A failure is only a warning. Only NativeBool attributes are supported for now.
	SetAfterCreate bool
	// CreateMapping overrides the create helper deduced from Type, the schema default, and CreateOnly.
	// Nil means deduce one.
	CreateMapping CreateMapping
	// UpdateMapping overrides the update helper deduced the same way.
	// Nil means deduce one.
	UpdateMapping UpdateMapping
}

// TagForm selects the TAG clause a Tags child emits.
type TagForm string

const (
	// TagFormSet emits SET TAG.
	TagFormSet TagForm = "set"
	// TagFormUnset emits UNSET TAG.
	TagFormUnset TagForm = "unset"
)

// Tags is the optional TAG clause. Put it on a CREATE list. It is not an Attributes entry.
var Tags = &Attribute{Name: "tags", Tags: true}

// Parameters is the parameter catalog spread. Put it on a CREATE, SET, or UNSET list. It is not an Attributes entry.
var Parameters = &Attribute{Name: "parameters", Parameters: true}

// SdkType is the type the SDK holds the attribute as: a string for an SdkString object, the Type otherwise.
// SQL clauses, SHOW and DESCRIBE columns, and the code generated from them read it. The resource reads Type.
func (a Attribute) SdkType() AttributeType {
	if a.SdkString {
		return AttrString
	}
	return a.Type
}

func (a *Attribute) IsTags() bool {
	return a != nil && a.Tags
}

func (a *Attribute) IsParameters() bool {
	return a != nil && a.Parameters
}

func (a *Attribute) OutputKey() string {
	if a.ShowName != "" {
		return a.ShowName
	}
	return a.Name
}

// ChangesInPlace reports a resource attribute whose change Update applies: it is not ForceNew, Readonly, Computed,
// CreateOnly, or SqlOnly.
func (a Attribute) ChangesInPlace() bool {
	return !a.ForceNew && !a.Readonly && !a.Computed && !a.CreateOnly && !a.SqlOnly
}

// ColumnEnum returns ShowEnum when it is set, and Enum otherwise.
// SHOW and CREATE/SET use that result. DESCRIBE calls it on ForDescribe, which has cleared ShowEnum.
func (a Attribute) ColumnEnum() *Enum {
	if a.ShowEnum != nil {
		return a.ShowEnum
	}
	return a.Enum
}

// ForShow is the attribute as a SHOW column reads it. ShowString makes that column a string.
func (a Attribute) ForShow() Attribute {
	if !a.ShowString {
		return a
	}
	a.Type = AttrString
	a.Enum = nil
	a.ShowEnum = nil
	a.Elem = nil
	a.Children = nil
	return a
}

// ForDescribe is the attribute as a DESCRIBE column reads it.
// DescribeString makes that column a string. DescribeEnum replaces the enum for that column.
func (a Attribute) ForDescribe() Attribute {
	if a.DescribeString {
		a.Type = AttrString
		a.Enum = nil
		a.ShowEnum = nil
		a.DescribeEnum = nil
		a.Elem = nil
		a.Children = nil
		return a
	}
	if a.DescribeEnum != nil {
		a.Enum = a.DescribeEnum
	}
	// DESCRIBE uses DescribeEnum or Enum. ShowEnum stays on SHOW and CREATE/SET.
	a.ShowEnum = nil
	return a
}

// IsEnumCollection reports a list or set whose element is an enum.
func (a Attribute) IsEnumCollection() bool {
	return (a.Type == AttrList || a.Type == AttrSet) && a.Elem != nil && a.Elem.Type == AttrEnum
}
