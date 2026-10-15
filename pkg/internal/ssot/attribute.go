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
	RawDescription        bool
	Enum                  *Enum
	Elem                  *Attribute
	OptionalInShow        bool // SHOW/DESCRIBE SDK field is a pointer; drives schema ToSchema nil-guards
	OptionalInDb          bool // DB row is nullable; live for the SDK paired-struct db and plain kinds
	StoredAsText          bool // the db column is text, and the plain kind stays the semantic type
	ShowName              string
	Readonly              bool
}

func (a Attribute) OutputKey() string {
	if a.ShowName != "" {
		return a.ShowName
	}
	return a.Name
}
