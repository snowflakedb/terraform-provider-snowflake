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

// ChildrenRule is a validation over the fields of the SQL struct of an object attribute.
type ChildrenRule string

const (
	// ChildrenAtLeastOne requires at least one field to be set (AtLeastOneValueSet).
	ChildrenAtLeastOne ChildrenRule = "at_least_one"
	// ChildrenExactlyOne requires exactly one field to be set (ExactlyOneValueSet).
	ChildrenExactlyOne ChildrenRule = "exactly_one"
)

// ObjectElem is the Elem of a list or set of objects: an AttrObject whose Children are the fields of each element.
func ObjectElem(children []Attribute) *Attribute {
	return &Attribute{Type: AttrObject, Children: children}
}

// Attribute is the written intent for one object attribute.
type Attribute struct {
	Name                  string
	Type                  AttributeType
	Required              bool
	Description           string
	AdditionalDescription string
	ForceNew              bool
	Min                   *int
	Max                   *int
	// MinItems and MaxItems are Terraform list or set bounds. Nil means unset.
	// A list with MinItems and no MaxItems holds multiple items, as storage_location does.
	MinItems       *int
	MaxItems       *int
	RawDescription bool
	Enum           *Enum
	// ShowEnum is the SHOW, DESCRIBE, and CREATE/SET enum when it differs from Enum.
	// Empty uses Enum. Resource schema and model builders keep Enum.
	ShowEnum *Enum
	// IncludeAliasedInTheDescription lists aliases in the resource description.
	// Empty lists canonical values.
	IncludeAliasedInTheDescription bool
	// Elem is the element of a list or set: a scalar attribute, or an AttrObject whose Children are the fields of a nested
	// block. A field of an Elem object may carry its own CREATE clause, which is how SDK query structs are generated.
	Elem *Attribute
	// Children are the fields of an AttrObject: a single block, an Elem, or a SqlOnly nested struct. They are valid only
	// there: a list or set of objects declares its fields on an Elem object (ObjectElem), never on the list itself.
	// A single block sets no MinItems or MaxItems: Terraform gets MaxItems 1, and MinItems 1 when it is Required.
	// A SqlOnly child is fixed text of the SQL struct, such as TYPE = OAUTH: a StaticSQL clause, valid when the object has an
	// object clause. It is in NestedSqlFields and not in NestedFields, so it is not a Terraform field.
	// A nested field may itself be a block. Generated create parses an optional list block; see the README.
	// Resource alter generation skips a nested attribute and every field under it.
	Children []Attribute
	// ChildrenRule adds a validation over the Children that are fields of the SQL struct of an object, not fixed texts:
	// at least one of them is set, or exactly one. The order is the Children order. It is valid only on an object with at
	// least two such Children. The schema-side AtLeastOneOf is not generated from it.
	ChildrenRule   ChildrenRule
	OptionalInShow bool // SHOW/DESCRIBE SDK field is a pointer; drives schema ToSchema nil-guards
	OptionalInDb   bool // DB row is nullable; live for the SDK paired-struct db and plain kinds
	StoredAsText   bool // the db column is text, and the plain kind stays the semantic type
	// BoolParsed, with StoredAsText, parses a bool column with strconv.ParseBool.
	// Without it, a text bool is true when the column is Y.
	BoolParsed    bool
	ManualConvert bool // paired-struct convert() skips this column; ext additionalConvert fills it
	CustomParser  string
	SkipInSchema  bool // column stays on the SDK struct and assertions, and is omitted from the public schema
	// SqlOnly keeps the attribute on a SQL sequence and out of schema, resource, model, and assertion generation.
	// ResourceAttributes drops these. A static CREATE value sets CreateSql to StaticSQL.
	SqlOnly        bool
	IdentifierType string // identifier Go type without the sdk prefix; required when Type is AttrIdentifier
	ShowName       string
	Readonly       bool
	// CreateMapping overrides the create helper deduced from Type, the schema default, and CreateOnly.
	// Nil means deduce one. A required attribute with a nil mapping is a create constructor argument.
	// SkipCreate leaves the attribute out of that deduction and out of constructor arguments.
	CreateMapping CreateMapping
	// SkipCreate leaves the attribute out of generated create.
	// There is no constructor argument, no ApplyCreateOptionals call, and no block parse function.
	// Ext create code builds the value when CREATE still needs it. The update helper stays.
	SkipCreate bool
	// UpdateMapping overrides the update helper deduced the same way.
	// Nil means deduce one. ForceNew and CreateOnly attributes are left out of set/unset.
	// SkipCreate does not leave the attribute out of set/unset.
	UpdateMapping UpdateMapping
	CreateSql     SqlClause // CREATE SQL clause. Zero is absent from that statement.
	SetSql        SqlClause // SET SQL clause. Zero is absent from that statement.
	UnsetSql      SqlClause // UNSET SQL clause. Zero is absent from that statement.
	DiffSuppress  DiffSuppress
	// Computed marks a field the resource schema ext turns into a computed attribute.
	// The generator still emits Optional. Drift tracking and output suppression skip it.
	Computed bool
	// ServerDefault marks a describe-only field that Snowflake fills in when the configuration omits it.
	// The generator drift-tracks that field and suppresses a diff against describe_output.
	// A SHOW column does not need this: every optional scalar SHOW column is drift-tracked.
	// Leave it unset when the omitted value is empty, as catalog_namespace is.
	ServerDefault bool
	// ConflictsWith names other attributes on the same level. The generator emits schema.ConflictsWith
	// as a path from the resource root: the name at the top, and parent.0.name under a single block
	// (an object, or a list of objects with MaxItems 1).
	// A set, or a list without that bound, cannot name these relations.
	ConflictsWith []Attribute
	// AtLeastOneOf names attributes of which at least one must be set. Paths follow ConflictsWith.
	// The list may include this attribute.
	AtLeastOneOf []Attribute
	// Sensitive marks the Terraform schema field sensitive.
	Sensitive bool
	// CreateOnly marks an optional bool that is applied only on create.
	// The generator emits IgnoreAfterCreation and a create helper, not an update helper.
	// The Terraform type stays a boolean string unless NativeBool is set.
	CreateOnly bool
	// NativeBool keeps an optional bool as Terraform TypeBool.
	// Optional bools otherwise use boolean-string schema chrome, including create-only ones.
	NativeBool bool
	// RelatedResource is the SSOT object whose Terraform resource is linked from the description.
	// The generator emits resources.<Name>. Nil means no link.
	RelatedResource *Object
	Tags            bool // optional TAG clause. Only valid on a CREATE list, not on Attributes.
	Parameters      bool // parameter catalog spread. Valid on a CREATE, SET, or UNSET list, not on Attributes.
}

// Tags is the optional TAG clause. Put it on a CREATE list. It is not an Attributes entry.
var Tags = &Attribute{Name: "tags", Tags: true}

// Parameters is the parameter catalog spread. Put it on a CREATE, SET, or UNSET list. It is not an Attributes entry.
var Parameters = &Attribute{Name: "parameters", Parameters: true}

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

// ColumnEnum is the enum used for a SHOW, DESCRIBE, or CREATE/SET column.
func (a Attribute) ColumnEnum() *Enum {
	if a.ShowEnum != nil {
		return a.ShowEnum
	}
	return a.Enum
}
