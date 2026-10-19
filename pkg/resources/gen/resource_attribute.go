package gen

// SchemaValueType is the plugin-sdk schema.ValueType name without the package prefix.
type SchemaValueType string

const (
	SchemaTypeString SchemaValueType = "TypeString"
	SchemaTypeInt    SchemaValueType = "TypeInt"
	SchemaTypeBool   SchemaValueType = "TypeBool"
	SchemaTypeList   SchemaValueType = "TypeList"
	SchemaTypeSet    SchemaValueType = "TypeSet"
)

// ResourceAttribute is schema chrome after applying resource attribute rules.
type ResourceAttribute struct {
	Name                 string
	TfType               SchemaValueType
	Required             bool
	Optional             bool
	ForceNew             bool
	DescriptionExpr      string
	DefaultExpr          string
	ValidateDiagExpr     string
	DiffSuppressExpr     string
	Elem                 *ResourceAttribute
	EnumName             string // SDK enum type name; empty when the attribute is not an enum
	CreateConstructorArg *CreateConstructorArg
	CreateOptionalCall   *CreateOptionalCall
	UpdateCall           *UpdateCall
}
