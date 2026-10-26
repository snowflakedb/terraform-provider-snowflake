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
)

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
	Elem                           *Attribute
	// Children is the named fields of a list or set block. Elem stays a scalar collection element.
	// A nested field may itself be a block. Generated create parses an optional list block; see the README.
	// Alter generation skips a nested attribute and every field under it.
	Children       []Attribute
	OptionalInShow bool // SHOW/DESCRIBE SDK field is a pointer; drives schema ToSchema nil-guards
	OptionalInDb   bool // DB row is nullable; live for the SDK paired-struct db and plain kinds
	StoredAsText   bool // the db column is text, and the plain kind stays the semantic type
	ManualConvert  bool // paired-struct convert() skips this column; ext additionalConvert fills it
	CustomParser   string
	SkipInSchema   bool // column stays on the SDK struct and assertions, and is omitted from the public schema
	// SqlOnly keeps the attribute on a SQL sequence and out of schema, resource, model, and assertion generation.
	// ResourceAttributes drops these. A static CREATE value sets CreateSql to StaticSQL.
	SqlOnly        bool
	IdentifierType string // identifier Go type without the sdk prefix; required when Type is AttrIdentifier
	ShowName       string
	Readonly       bool
	// CreateMapping overrides the create helper deduced from Type, the schema default, and CreateOnly.
	// Nil means deduce one. A required attribute with a nil mapping is a create constructor argument.
	CreateMapping CreateMapping
	// SkipCreateBlock leaves a list block out of generated create.
	// Ext create code builds the value. The generator does not emit a parse function.
	SkipCreateBlock bool
	// UpdateMapping overrides the update helper deduced the same way.
	// Nil means deduce one. ForceNew and CreateOnly attributes are left out of set/unset.
	UpdateMapping UpdateMapping
	CreateSql     SqlClause // CREATE SQL clause. Zero is absent from that statement.
	SetSql        SqlClause // SET SQL clause. Zero is absent from that statement.
	UnsetSql      SqlClause // UNSET SQL clause. Zero is absent from that statement.
	DiffSuppress  DiffSuppress
	// ConflictsWith names other attributes on the same object. The generator emits schema.ConflictsWith.
	ConflictsWith []Attribute
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
