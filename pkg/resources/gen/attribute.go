package gen

// AttributeType is the semantic type of a resource attribute, not the Terraform schema type.
type AttributeType string

const (
	AttrString     AttributeType = "string"
	AttrInt        AttributeType = "int"
	AttrBool       AttributeType = "bool"
	AttrIdentifier AttributeType = "identifier"
	AttrEnum       AttributeType = "enum"
	AttrList       AttributeType = "list"
	AttrSet        AttributeType = "set"
)

// SchemaValueType is the plugin-sdk schema.ValueType name without the package prefix.
type SchemaValueType string

const (
	SchemaTypeString SchemaValueType = "TypeString"
	SchemaTypeInt    SchemaValueType = "TypeInt"
	SchemaTypeBool   SchemaValueType = "TypeBool"
	SchemaTypeList   SchemaValueType = "TypeList"
	SchemaTypeSet    SchemaValueType = "TypeSet"
)

// Attribute is the written intent for one resource attribute.
type Attribute struct {
	Name                  string
	Type                  AttributeType
	Required              bool
	Description           string
	AdditionalDescription string
	ForceNew              bool
	Min                   *int
	RawDescription        bool
	// TODO [next PRs]: move enum abstraction from SDK (provides both name and plural).
	Enum       string
	EnumPlural string
	Elem       *Attribute
}

// ResolvedAttribute is schema chrome after applying attribute rules.
type ResolvedAttribute struct {
	Name             string
	TfType           SchemaValueType
	Required         bool
	Optional         bool
	ForceNew         bool
	DescriptionExpr  string
	DefaultExpr      string
	ValidateDiagExpr string
	DiffSuppressExpr string
	Elem             *ResolvedAttribute
}
