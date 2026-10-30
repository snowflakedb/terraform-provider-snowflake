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
	Name             string
	TfType           SchemaValueType
	Required         bool
	Optional         bool
	Computed         bool
	ForceNew         bool
	Sensitive        bool
	MinItems         *int
	MaxItems         *int
	DescriptionExpr  string
	DefaultExpr      string
	ValidateDiagExpr string
	DiffSuppressExpr string
	ConflictsWith    []string
	AtLeastOneOf     []string
	ExactlyOneOf     []string
	Elem             *ResourceAttribute
	Children         []ResourceAttribute
	// SdkString marks a block that SQL and the outputs hold as one string. The recompute lists still watch it.
	SdkString            bool
	EnumName             string // SDK enum type name; empty when the attribute is not an enum
	CreateConstructorArg *CreateConstructorArg
	CreateOptionalCall   *CreateOptionalCall
	CreateBlock          *CreateRequestFunc
	UpdateCall           *UpdateCall
}
