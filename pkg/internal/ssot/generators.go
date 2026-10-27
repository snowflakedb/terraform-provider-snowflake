package ssot

import (
	"fmt"
)

type SchemaGenerator struct {
	Show     bool
	Describe bool
	// ShowManualFields and DescribeManualFields are the same *Attribute pointers that sit in Columns.
	// The key stays in the generated schema map and ToSchema, but the mapping line is omitted: a
	// generated additionalSchemaMapper[T] stub is emitted, and *_ext.go must implement it.
	ShowManualFields     []*Attribute
	DescribeManualFields []*Attribute
}

type SdkGenerator struct {
	// DefineFromSsot defines this object's SDK interface from this ssot object alone.
	DefineFromSsot bool

	Show     bool
	Describe bool

	ShowDoc                string
	ShowSql                string
	ShowTerse              bool
	ShowOptionalSql        []string
	ShowOptionalLike       bool
	ShowOptionalExtendedIn bool
	ShowOptionalStartsWith bool
	ShowOptionalLimitFrom  bool
	ShowOptionalLimit      bool
	ShowByIDFiltering      []SdkShowByIDFilteringKind
	ShowQueryStructName    string

	DescribeDoc     string
	DescribeSql     string
	DescribeMapping SdkDescriptionMapping
	// DescribeDetails emits the plain details struct (<Name>Details) from Describe.Columns next to the DESCRIBE rows.
	// DescribeDetailsNestedFields are the same *Attribute pointers that sit in Describe.Columns. Each is a list or set
	// whose elements are a hand-written struct, so the column is not derived from its placeholder Elem.
	DescribeDetails             bool
	DescribeDetailsNestedFields []*Attribute
	// DescribeDetailsMethod emits the DescribeDetails method: Describe, Parse<Name>Described, and Id.
	// The parse function stays in the SDK ext file.
	DescribeDetailsMethod bool
	// DescribeDetailsParse emits Parse<Name>Described, which fills the details struct from the DESCRIBE property rows.
	// Every Describe column except id is mapped by one entry. The per-element parser of a prefix entry and the
	// validate<Name>Described function stay in the SDK ext file.
	DescribeDetailsParse    []DescribeDetailsMapping
	DescribeQueryStructName string
	// DescribePlainOnlyFields are fields on the DESCRIBE plain struct that are not DESCRIBE columns.
	// convert skips them. additionalConvert in the SDK ext file fills them.
	DescribePlainOnlyFields []SdkPlainOnlyField

	Drop                bool
	DropDoc             string
	DropSql             string
	DropQueryStructName string
	DropSafelyHook      bool
	DropSafelyForce     bool
	DropCascadeRestrict bool

	Create                      bool
	CreateDoc                   string
	CreateSql                   string
	CreateQueryStructName       string
	CreateWithoutIfNotExists    bool
	CreateWithoutOrReplace      bool
	CreateAdditionalValidations bool
	// CreateForVariants emits each variant CREATE as a custom method on this object's interface.
	// TODO [next PRs]: the same switch for variant SHOW/DESCRIBE.
	CreateForVariants bool

	Alter                          bool
	AlterDoc                       string
	AlterSql                       string
	AlterQueryStructName           string
	AlterWithoutIfExists           bool
	AlterSetList                   bool
	AlterSetAdditionalValidations  bool
	AlterAdditionalValidations     bool
	AlterWithoutExactlyOneValueSet bool
	// AlterForVariants emits each variant ALTER as a custom method on this object's interface.
	// TODO [next PRs]: the same switch for variant SHOW/DESCRIBE.
	AlterForVariants bool

	CustomOperations       bool
	CustomInterfaceMethods bool
	ShowParameters         bool
	ShowParametersDetails  bool
}

type ModelGenerator struct {
	Enabled bool
	// RequestTypeOverrides maps an attribute Name to the real Go type for a Required nested-block
	// attribute (a list or set of objects) that the generator cannot classify on its own. Mirrors the
	// old-list complexListAttributesOverrides table, scoped to this object instead of by resource name.
	RequestTypeOverrides map[string]string
	// Multiline lists string attributes whose config model builder emits MultilineWrapperVariable.
	// Mirrors the old-list multilineAttributesOverrides, scoped to this object.
	Multiline []*Attribute
}

// ResourceGenerator emits the Terraform resource. Enabled is set by WithResource.
// An object in All without it is not a resource-generation input.
type ResourceGenerator struct {
	Enabled bool
}

type ResourceAssertGenerator struct {
	Enabled bool
}

type ShowOutputAssertGenerator struct {
	Show     bool
	Describe bool
	// ShowSkipFields and DescribeSkipFields are the same *Attribute pointers that sit in Columns.
	// The column stays in Columns. Has* is not generated.
	ShowSkipFields     []*Attribute
	DescribeSkipFields []*Attribute
}

type ObjectAssertGenerator struct {
	Show                 bool
	Describe             bool
	NoIdentifiableObject bool
	// ShowSkipFields and DescribeSkipFields are the same *Attribute pointers that sit in Columns.
	// The column stays in Columns. Has* is not generated.
	ShowSkipFields     []*Attribute
	DescribeSkipFields []*Attribute
	// ShowProjection is the method on the VariantOf row, such as AsAdaptive.
	// The fetch is testClient.<VariantOf.Name>.Show, then that method.
	// The object type constant is sdk.ObjectType<VariantOf.Name>.
	ShowProjection string

	// ShowNestedFields and DescribeNestedFields hold the *Attribute pointers that sit in Columns.
	// The column gets the nested assertion adapter (Has<Field>With) instead of a plain Has*.
	ShowNestedFields     []ObjectAssertNestedField
	DescribeNestedFields []ObjectAssertNestedField
}

// ObjectAssertNestedField is a nested-block column. The nested assertion type is derived from the Go type of
// its SDK field. GoType is that type when set, e.g. "[]sdk.ExternalVolumeStorageLocationDetails".
// Empty GoType means the default for a list or set column: "[]sdk." + object Name + the column name in
// camel case with one trailing "s" dropped + "Details" (storage_locations on ExternalVolume is
// []sdk.ExternalVolumeStorageLocationDetails).
type ObjectAssertNestedField struct {
	Attr   *Attribute
	GoType string
}

// ObjectAssertNestedOption overrides the defaults of one nested-block column.
type ObjectAssertNestedOption func(*ObjectAssertNestedField)

// ObjectAssertNestedGoType sets the Go type of the SDK field when the default composition does not match.
func ObjectAssertNestedGoType(goType string) ObjectAssertNestedOption {
	return func(f *ObjectAssertNestedField) { f.GoType = goType }
}

// ObjectParametersAssertGenerator emits Snowflake object parameter assertions from the parameter catalog.
// Constants are sdk.<Name>Parameter.
type ObjectParametersAssertGenerator struct {
	Enabled bool
}

// ResourceParametersAssertGenerator emits resource parameter assertions from the parameter catalog.
// Constants are sdk.<Name>Parameter.
type ResourceParametersAssertGenerator struct {
	Enabled bool
}

type Generators struct {
	Schema                   SchemaGenerator
	Sdk                      SdkGenerator
	Resource                 ResourceGenerator
	Model                    ModelGenerator
	ResourceAssert           ResourceAssertGenerator
	ShowOutputAssert         ShowOutputAssertGenerator
	ObjectAssert             ObjectAssertGenerator
	ObjectParametersAssert   ObjectParametersAssertGenerator
	ResourceParametersAssert ResourceParametersAssertGenerator
}

type GeneratorOption func(*Generators)

type SchemaOption func(*SchemaGenerator)

type SchemaShowOption func(*SchemaGenerator)

type SchemaDescribeOption func(*SchemaGenerator)

func SchemaShow(opts ...SchemaShowOption) SchemaOption {
	return func(s *SchemaGenerator) {
		s.Show = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SchemaShowManualFields copies ManualFields from the old show schema def.
// Each argument is the *Attribute already listed in Show.Columns.
func SchemaShowManualFields(attrs ...*Attribute) SchemaShowOption {
	return func(s *SchemaGenerator) { s.ShowManualFields = append(s.ShowManualFields, attrs...) }
}

func SchemaDescribe(opts ...SchemaDescribeOption) SchemaOption {
	return func(s *SchemaGenerator) {
		s.Describe = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SchemaDescribeManualFields copies ManualFields from the old describe schema def.
// Each argument is the *Attribute already listed in Describe.Columns.
func SchemaDescribeManualFields(attrs ...*Attribute) SchemaDescribeOption {
	return func(s *SchemaGenerator) { s.DescribeManualFields = append(s.DescribeManualFields, attrs...) }
}

func WithSchema(opts ...SchemaOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.Schema)
		}
	}
}

type SdkOption func(*SdkGenerator)

// SdkDefineFromSsot defines this object's SDK interface from this ssot object alone.
// The interface is built from the operation parts whose flags are on. A part whose flag is off is absent.
// NamePlural must be unique among objects that are not variants. A variant shares this interface.
// An object without this option stays on the handwritten SDK def.
func SdkDefineFromSsot() SdkOption {
	return func(s *SdkGenerator) { s.DefineFromSsot = true }
}

type SdkShowOption func(*SdkGenerator)

type SdkDescribeOption func(*SdkGenerator)

func SdkShow(opts ...SdkShowOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Show = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

func SdkDescribe(opts ...SdkDescribeOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Describe = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SdkDescribeDetails emits the plain details struct from Describe.Columns. It needs DescribeRows, because the rows
// then take the <Name>Property name and <Name>Details is free. The Id field comes first unless a column is named id.
// Parsing the rows into the struct stays in the SDK ext file.
func SdkDescribeDetails(opts ...SdkDescribeDetailsOption) SdkDescribeOption {
	return func(s *SdkGenerator) {
		s.DescribeDetails = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

type SdkDescribeDetailsOption func(*SdkGenerator)

// DescribeDetailsMapping maps one Describe column to the property rows it is read from.
// Prefix is false for an exact property name and true for a nested list, whose elements are the rows whose property
// name contains Property.
type DescribeDetailsMapping struct {
	Attr     *Attribute
	Property string
	Prefix   bool
}

// DescribeProperty reads a string column from the property row with exactly this name.
func DescribeProperty(attr *Attribute, property string) DescribeDetailsMapping {
	return DescribeDetailsMapping{Attr: attr, Property: property}
}

// DescribePropertyPrefix reads a nested list column from every property row whose name contains prefix. The column must
// be marked with SdkDescribeDetailsNestedField. Each row value goes through the hand-written element parser
// Parse<Name><Column key in camel case without the plural s>Described(value string) (<element>, error).
func DescribePropertyPrefix(attr *Attribute, prefix string) DescribeDetailsMapping {
	return DescribeDetailsMapping{Attr: attr, Property: prefix, Prefix: true}
}

// SdkDescribeDetailsParse emits the function that parses the DESCRIBE property rows into the details struct.
// Exact properties are matched first, then prefixes, and any other property is an error. The generated function is
// Parse<Name>Described, which SdkDescribeDetailsMethod calls. validate<Name>Described stays in the SDK ext file.
func SdkDescribeDetailsParse(mappings ...DescribeDetailsMapping) SdkDescribeDetailsOption {
	return func(s *SdkGenerator) { s.DescribeDetailsParse = append(s.DescribeDetailsParse, mappings...) }
}

// SdkDescribeDetailsMethod emits the DescribeDetails method on the interface and its body: it calls Describe, parses
// the rows with the hand-written Parse<Name>Described(rows []<Name>Property) (<Name>Details, error), and sets Id.
func SdkDescribeDetailsMethod() SdkDescribeDetailsOption {
	return func(s *SdkGenerator) { s.DescribeDetailsMethod = true }
}

// SdkDescribeDetailsNestedField marks a Describe column as a list or set of a hand-written struct.
// The element type is NestedDetailsElem: <Name>, the column key in camel case without the plural s, and Details.
func SdkDescribeDetailsNestedField(attrs ...*Attribute) SdkDescribeDetailsOption {
	return func(s *SdkGenerator) { s.DescribeDetailsNestedFields = append(s.DescribeDetailsNestedFields, attrs...) }
}

type SdkDropOption func(*SdkGenerator)

func SdkDrop(opts ...SdkDropOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Drop = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

func SdkDropSafelyHook() SdkDropOption {
	return func(s *SdkGenerator) { s.DropSafelyHook = true }
}

func SdkDropSafelyForce() SdkDropOption {
	return func(s *SdkGenerator) { s.DropSafelyForce = true }
}

// SdkDropCascadeRestrict adds optional CASCADE and RESTRICT, and marks them as conflicting.
func SdkDropCascadeRestrict() SdkDropOption {
	return func(s *SdkGenerator) { s.DropCascadeRestrict = true }
}

// SdkDropDoc overrides the default DROP doc URL.
func SdkDropDoc(url string) SdkDropOption {
	return func(s *SdkGenerator) { s.DropDoc = url }
}

type SdkCreateOption func(*SdkGenerator)

func SdkCreate(opts ...SdkCreateOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Create = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SdkCreateDoc overrides the default CREATE doc URL.
func SdkCreateDoc(url string) SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateDoc = url }
}

// SdkCreateSql overrides the default CREATE command words.
func SdkCreateSql(sql string) SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateSql = sql }
}

// SdkCreateQueryStructName overrides the default CREATE query struct name.
func SdkCreateQueryStructName(name string) SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateQueryStructName = name }
}

// SdkCreateWithoutIfNotExists leaves IF NOT EXISTS out. CREATE includes it unless this option is set.
func SdkCreateWithoutIfNotExists() SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateWithoutIfNotExists = true }
}

// SdkCreateWithoutOrReplace leaves OR REPLACE out. CREATE includes it unless this option is set.
func SdkCreateWithoutOrReplace() SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateWithoutOrReplace = true }
}

func SdkCreateAdditionalValidations() SdkCreateOption {
	return func(s *SdkGenerator) { s.CreateAdditionalValidations = true }
}

// SdkCreateForVariants emits CREATE for each variant of this object as a custom method.
// The variant supplies the CREATE spec.
// TODO [next PRs]: the same switch for variant SHOW/DESCRIBE.
func SdkCreateForVariants() SdkOption {
	return func(s *SdkGenerator) { s.CreateForVariants = true }
}

// SdkAlterForVariants emits ALTER for each variant of this object as a custom method.
// The variant supplies the ALTER spec.
// TODO [next PRs]: the same switch for variant SHOW/DESCRIBE.
func SdkAlterForVariants() SdkOption {
	return func(s *SdkGenerator) { s.AlterForVariants = true }
}

type SdkAlterOption func(*SdkGenerator)

func SdkAlter(opts ...SdkAlterOption) SdkOption {
	return func(s *SdkGenerator) {
		s.Alter = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// SdkAlterDoc overrides the default ALTER doc URL.
func SdkAlterDoc(url string) SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterDoc = url }
}

// SdkAlterSql overrides the default ALTER command words.
func SdkAlterSql(sql string) SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterSql = sql }
}

// SdkAlterQueryStructName overrides the default ALTER query struct name.
func SdkAlterQueryStructName(name string) SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterQueryStructName = name }
}

// SdkAlterWithoutIfExists leaves IF EXISTS out. ALTER includes it unless this option is set.
func SdkAlterWithoutIfExists() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterWithoutIfExists = true }
}

// SdkAlterSetList makes SET a list with no parentheses. SET is a keyword unless this option is set.
// TODO [next PRs]: Check why SET is a keyword for some objects and a list for others.
func SdkAlterSetList() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterSetList = true }
}

// SdkAlterSetAdditionalValidations marks the SET struct with additional validations.
func SdkAlterSetAdditionalValidations() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterSetAdditionalValidations = true }
}

// SdkAlterAdditionalValidations marks the ALTER query struct with additional validations.
func SdkAlterAdditionalValidations() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterAdditionalValidations = true }
}

// SdkAlterWithoutExactlyOneValueSet leaves out the check that exactly one alter branch is set.
// ALTER includes that check unless this option is set.
func SdkAlterWithoutExactlyOneValueSet() SdkAlterOption {
	return func(s *SdkGenerator) { s.AlterWithoutExactlyOneValueSet = true }
}

// SdkCustomOperations emits Object.CustomOperations. Each entry is one custom operation.
// The query struct name is the kind plus Name. The object words are the uppercase words of Name.
// The doc URL is the kind slug, a hyphen, and the name slug, unless the spec sets Doc.
func SdkCustomOperations() SdkOption {
	return func(s *SdkGenerator) { s.CustomOperations = true }
}

// SdkCustomInterfaceMethods emits Object.CustomInterfaceMethods. Each entry is one interface declaration.
// The body stays in the SDK ext file. The default signature is id of IdType and a return of error.
// An empty doc emits no comment.
func SdkCustomInterfaceMethods() SdkOption {
	return func(s *SdkGenerator) { s.CustomInterfaceMethods = true }
}

// SdkShowParameters emits ShowParameters as an interface declaration.
// The doc is empty. The parameter is id of IdType. The returns are []*Parameter and error.
// The body stays in the SDK ext file.
func SdkShowParameters() SdkOption {
	return func(s *SdkGenerator) { s.ShowParameters = true }
}

// SdkShowParametersDetails emits the typed parameters-details struct and ShowParametersDetails.
// The fields are ParameterLevel() plus ShowParametersDetailsAdditionalLevels.
// The struct name is Name plus ParametersDetails. The generated body calls ShowParameters.
func SdkShowParametersDetails() SdkOption {
	return func(s *SdkGenerator) { s.ShowParametersDetails = true }
}

func WithSdk(opts ...SdkOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.Sdk)
		}
	}
}

type ResourceOption func(*ResourceGenerator)

// WithResource enables resource generation for this object.
// Zero options still sets Enabled. An object without this call is left out of make generate-resource.
func WithResource(opts ...ResourceOption) GeneratorOption {
	return func(g *Generators) {
		g.Resource.Enabled = true
		for _, opt := range opts {
			opt(&g.Resource)
		}
	}
}

type ModelOption func(*ModelGenerator)

// ModelRequestType overrides the constructor parameter type for a Required nested-block attribute
// (a list or set of objects). attr is the same *Attribute listed in Attributes. goType is the full type
// expression used as the constructor parameter type, e.g. "[]sdk.ExternalVolumeStorageLocationRequest".
func ModelRequestType(attr *Attribute, goType string) ModelOption {
	return func(m *ModelGenerator) {
		if m.RequestTypeOverrides == nil {
			m.RequestTypeOverrides = make(map[string]string)
		}
		m.RequestTypeOverrides[attr.Name] = goType
	}
}

// ModelMultiline marks string attributes whose config model builder emits MultilineWrapperVariable.
// Each argument is the same *Attribute listed in Attributes. Mirrors multilineAttributesOverrides.
func ModelMultiline(attrs ...*Attribute) ModelOption {
	return func(m *ModelGenerator) {
		for _, attr := range attrs {
			if attr == nil {
				panic("ModelMultiline: nil attribute")
			}
			m.Multiline = append(m.Multiline, attr)
		}
	}
}

func WithModel(opts ...ModelOption) GeneratorOption {
	return func(g *Generators) {
		g.Model.Enabled = true
		for _, opt := range opts {
			opt(&g.Model)
		}
	}
}

type ResourceAssertOption func(*ResourceAssertGenerator)

func WithResourceAssert(opts ...ResourceAssertOption) GeneratorOption {
	return func(g *Generators) {
		g.ResourceAssert.Enabled = true
		for _, opt := range opts {
			opt(&g.ResourceAssert)
		}
	}
}

type ShowOutputAssertOption func(*ShowOutputAssertGenerator)

type ShowOutputShowOption func(*ShowOutputAssertGenerator)

type ShowOutputDescribeOption func(*ShowOutputAssertGenerator)

func ShowOutputShow(opts ...ShowOutputShowOption) ShowOutputAssertOption {
	return func(s *ShowOutputAssertGenerator) {
		s.Show = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// ShowOutputShowSkipFields copies SkipFields from the old show output def.
// Each argument is the *Attribute already listed in Show.Columns.
func ShowOutputShowSkipFields(attrs ...*Attribute) ShowOutputShowOption {
	return func(s *ShowOutputAssertGenerator) { s.ShowSkipFields = append(s.ShowSkipFields, attrs...) }
}

func ShowOutputDescribe(opts ...ShowOutputDescribeOption) ShowOutputAssertOption {
	return func(s *ShowOutputAssertGenerator) {
		s.Describe = true
		for _, opt := range opts {
			opt(s)
		}
	}
}

// ShowOutputDescribeSkipFields copies SkipFields from the old describe show-output def.
// Each argument is the *Attribute already listed in Describe.Columns.
func ShowOutputDescribeSkipFields(attrs ...*Attribute) ShowOutputDescribeOption {
	return func(s *ShowOutputAssertGenerator) { s.DescribeSkipFields = append(s.DescribeSkipFields, attrs...) }
}

func WithShowOutputAssert(opts ...ShowOutputAssertOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.ShowOutputAssert)
		}
	}
}

type ObjectAssertOption func(*ObjectAssertGenerator)

type ObjectAssertShowOption func(*ObjectAssertGenerator)

type ObjectAssertDescribeOption func(*ObjectAssertGenerator)

func ObjectAssertShow(opts ...ObjectAssertShowOption) ObjectAssertOption {
	return func(a *ObjectAssertGenerator) {
		a.Show = true
		for _, opt := range opts {
			opt(a)
		}
	}
}

// ObjectAssertShowNestedField copies NestedAssertFields from the old show object-assertion def.
// attr is the *Attribute already listed in Show.Columns.
func ObjectAssertShowNestedField(attr *Attribute, opts ...ObjectAssertNestedOption) ObjectAssertShowOption {
	return func(a *ObjectAssertGenerator) {
		a.ShowNestedFields = append(a.ShowNestedFields, newObjectAssertNestedField(attr, opts))
	}
}

func newObjectAssertNestedField(attr *Attribute, opts []ObjectAssertNestedOption) ObjectAssertNestedField {
	field := ObjectAssertNestedField{Attr: attr}
	for _, opt := range opts {
		opt(&field)
	}
	return field
}

// ObjectAssertShowSkipFields copies SkipFields from the old show object-assertion def.
// Each argument is the *Attribute already listed in Show.Columns.
func ObjectAssertShowSkipFields(attrs ...*Attribute) ObjectAssertShowOption {
	return func(a *ObjectAssertGenerator) { a.ShowSkipFields = append(a.ShowSkipFields, attrs...) }
}

// ObjectAssertShowFromBase fetches through testClient.<VariantOf.Name>.<method>, such as ShowByIDAdaptive.
// The client name is VariantOf.Name, returned by Object.ShowFromBase.
func ObjectAssertShowFromBase(projection string) ObjectAssertShowOption {
	return func(a *ObjectAssertGenerator) {
		a.ShowProjection = projection
	}
}

// ShowFromBase is the SHOW fetch for a variant object assertion.
// ClientName is VariantOf.Name. Method is the ObjectAssertShowFromBase argument.
type ShowFromBase struct {
	ClientName string
	Method     string
}

// ShowFromBase returns the resolved SHOW projection. ok is false when ObjectAssertShowFromBase is not set.
func (obj Object) ShowFromBase() (ShowFromBase, bool) {
	method := ObjectGenerators(obj.Generators...).ObjectAssert.ShowProjection
	if method == "" {
		return ShowFromBase{}, false
	}
	if obj.VariantOf == nil || obj.VariantOf.Name == "" {
		panic(fmt.Sprintf("object %q: object assertion show projection %q requires VariantOf", obj.Name, method))
	}
	return ShowFromBase{ClientName: obj.VariantOf.Name, Method: method}, true
}

func ObjectAssertDescribe(opts ...ObjectAssertDescribeOption) ObjectAssertOption {
	return func(a *ObjectAssertGenerator) {
		a.Describe = true
		for _, opt := range opts {
			opt(a)
		}
	}
}

// ObjectAssertDescribeNestedField copies NestedAssertFields from the old describe object-assertion def.
// attr is the *Attribute already listed in Describe.Columns.
func ObjectAssertDescribeNestedField(attr *Attribute, opts ...ObjectAssertNestedOption) ObjectAssertDescribeOption {
	return func(a *ObjectAssertGenerator) {
		a.DescribeNestedFields = append(a.DescribeNestedFields, newObjectAssertNestedField(attr, opts))
	}
}

func WithNoIdentifiableObject() ObjectAssertDescribeOption {
	return func(a *ObjectAssertGenerator) { a.NoIdentifiableObject = true }
}

// ObjectAssertDescribeSkipFields copies SkipFields from the old describe object-assertion def.
// Each argument is the *Attribute already listed in Describe.Columns.
func ObjectAssertDescribeSkipFields(attrs ...*Attribute) ObjectAssertDescribeOption {
	return func(a *ObjectAssertGenerator) { a.DescribeSkipFields = append(a.DescribeSkipFields, attrs...) }
}

func WithObjectAssert(opts ...ObjectAssertOption) GeneratorOption {
	return func(g *Generators) {
		for _, opt := range opts {
			opt(&g.ObjectAssert)
		}
	}
}

type ObjectParametersAssertOption func(*ObjectParametersAssertGenerator)

func WithObjectParametersAssert(opts ...ObjectParametersAssertOption) GeneratorOption {
	return func(g *Generators) {
		g.ObjectParametersAssert.Enabled = true
		for _, opt := range opts {
			opt(&g.ObjectParametersAssert)
		}
	}
}

type ResourceParametersAssertOption func(*ResourceParametersAssertGenerator)

func WithResourceParametersAssert(opts ...ResourceParametersAssertOption) GeneratorOption {
	return func(g *Generators) {
		g.ResourceParametersAssert.Enabled = true
		for _, opt := range opts {
			opt(&g.ResourceParametersAssert)
		}
	}
}

func ObjectGenerators(opts ...GeneratorOption) Generators {
	var g Generators
	for _, opt := range opts {
		opt(&g)
	}
	return g
}
