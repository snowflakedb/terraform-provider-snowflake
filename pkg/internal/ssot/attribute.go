package ssot

import "fmt"

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
	RawDescription        bool
	Enum                  *Enum
	Elem                  *Attribute
	OptionalInShow        bool // SHOW/DESCRIBE SDK field is a pointer; drives schema ToSchema nil-guards
	OptionalInDb          bool // DB row is nullable; live for the SDK paired-struct db and plain kinds
	StoredAsText          bool // the db column is text, and the plain kind stays the semantic type
	ShowName              string
	Readonly              bool
	CreateMapping         CreateMapping
	UpdateMapping         UpdateMapping
	// SuppressShowChange composes IgnoreChangeToCurrentSnowflakeValueInShow(OutputKey()) onto the diff suppressor.
	// TODO [next PRs]: Replace this bool with an explicit SHOW, DESCRIBE, or none choice.
	// Both outputs on an object do not imply SHOW. Some fields compare against DESCRIBE, some against neither,
	// and some need a mapping when the output column is not the attribute name.
	SuppressShowChange bool
	// ConflictsWith names other attributes on the same object. The generator emits schema.ConflictsWith.
	ConflictsWith []Attribute
	// CreateOnly keeps an optional bool as native TypeBool with IgnoreAfterCreation.
	// Optional bools otherwise become boolean strings.
	CreateOnly bool
	// RelatedResource is the SSOT object whose Terraform resource is linked from the description.
	// The generator emits resources.<Name>. Nil means no link.
	RelatedResource *Object
}

func (a Attribute) OutputKey() string {
	if a.ShowName != "" {
		return a.ShowName
	}
	return a.Name
}

// ConstructorGet is the d.Get expression for a required attribute with no CreateMapping.
// The bool is false when the attribute is not a constructor argument.
func (a Attribute) ConstructorGet() (string, bool) {
	if a.CreateMapping != nil || !a.Required {
		return "", false
	}
	var cast string
	switch a.Type {
	case AttrInt:
		cast = "int"
	case AttrString:
		cast = "string"
	case AttrBool:
		cast = "bool"
	case AttrEnum:
		if a.Enum == nil || a.Enum.Name == "" {
			panic(fmt.Sprintf("constructor attribute %q has no Enum", a.Name))
		}
		cast = "string"
	default:
		panic(fmt.Sprintf("required attribute %q has no CreateMapping", a.Name))
	}
	return fmt.Sprintf("d.Get(%q).(%s)", a.Name, cast), true
}
