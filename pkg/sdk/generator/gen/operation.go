package gen

import "fmt"

type OperationKind string

const (
	OperationKindCreate   OperationKind = "Create"
	OperationKindAlter    OperationKind = "Alter"
	OperationKindDrop     OperationKind = "Drop"
	OperationKindShow     OperationKind = "Show"
	OperationKindShowByID OperationKind = "ShowByID"
	OperationKindDescribe OperationKind = "Describe"
	OperationKindGrant    OperationKind = "Grant"
	OperationKindRevoke   OperationKind = "Revoke"
)

type DescriptionMappingKind string

const (
	DescriptionMappingKindSingleValue DescriptionMappingKind = "single_value"
	DescriptionMappingKindSlice       DescriptionMappingKind = "slice"
)

type ShowMappingKind string

const (
	ShowMappingKindSingleValue ShowMappingKind = "single_value"
	ShowMappingKindSlice       ShowMappingKind = "slice"
)

type InstanceMethodKind string

const (
	InstanceMethodKindSingleValue InstanceMethodKind = "single_value"
	InstanceMethodKindSlice       InstanceMethodKind = "slice"
)

// Operation defines a single operation for given object or objects family (e.g. CREATE DATABASE ROLE)
type Operation struct {
	// Name is the operation's name, e.g. "Create"
	Name string
	// ObjectInterface points to the containing interface
	ObjectInterface *Interface
	// Doc is the URL for the doc used to create given operation, e.g. https://docs.snowflake.com/en/sql-reference/sql/create-database-role
	Doc string
	// OptsField defines opts used to create SQL for given operation
	OptsField *Field
	// HelperStructs are struct definitions that are not tied to OptsField, but tied to the Operation itself, e.g. Show() return type
	HelperStructs []*Field
	// ShowKind defines a kind of mapping that needs to be performed in particular case of Show implementation
	// TODO(SNOW-2183036) This is a temporary solution to support single value and slice return types for Show operation.
	ShowKind *ShowMappingKind
	// ShowMapping is a definition of mapping needed by Operation kind of OperationKindShow
	ShowMapping *Mapping
	// DescribeKind defines a kind of mapping that needs to be performed in particular case of Describe implementation
	DescribeKind *DescriptionMappingKind
	// DescribeMapping is a definition of mapping needed by Operation kind of OperationKindDescribe
	DescribeMapping *Mapping
	// InstanceMethodMapping is a definition of mapping needed when an InstanceMethodOperation returns a result struct
	InstanceMethodMapping *Mapping
	// InstanceMethodKind defines the kind of result for an InstanceMethodOperation.
	// For single_value/slice it is set to the matching named constant and InstanceMethodMapping should also be set.
	// For scalar, the InstanceMethodScalarReturnType should be set and InstanceMethodMapping stays nil.
	InstanceMethodKind *InstanceMethodKind
	// InstanceMethodScalarReturnType should be set to the Go return type name (e.g. "int", "string") for scalar instance methods.
	InstanceMethodScalarReturnType string
	// ShowByIDFiltering defines a kind of filterings performed in ShowByID operation
	ShowByIDFiltering []ShowByIDFiltering
	// ShowResultFilterHook, when true, inserts a hook call in the generated Show implementation.
	// It allows filtering the rows by implementing excludeFromShow() method on the given row type.
	// The implementation should be provided in the _ext.go file.
	ShowResultFilterHook bool
	// DropSafelyHook, when true, inserts a hook call in the generated DropSafely implementation.
	// It allows running operation before dropping the object safely by implementing dropSafelyHook() method on the given object interface implementation.
	// The implementation should be provided in the _ext.go file.
	DropSafelyHook bool
	// DropSafelyForce, when true, appends .WithForce(true) to the Drop request in the generated DropSafely implementation.
	DropSafelyForce bool
	// RequestAdjust, when true, inserts a request.adjust() call in the generated impl before toOpts() is called.
	// The adjust() method must be implemented manually on the request type in a _ext.go file.
	RequestAdjust bool

	// TODO [SNOW-2324252]: Consider splitting the Operation into definition and generation model
	// new fields used to move the old template executors logic into simpler template generation based on prepared model

	// StructsToGenerate is a list of all newly introduced structs comprised of HelperStructs and OptsField; contains only unique structs
	StructsToGenerate []*Field
	// ObjectIdMethod is a model to generate the ID() method for an SDK object; replaces the old logic
	ObjectIdMethod *ShowObjectIdMethod
	// ObjectTypeMethod is a model to generate the ObjectType() method for an SDK object; replaces the old logic
	ObjectTypeMethod *ShowObjectTypeMethod
	// DtosToGenerate is a list of all newly introduced dto structs based on the operation opts; contains only unique structs
	DtosToGenerate []*Field
}

type Mapping struct {
	MappingFuncName string
	From            *Field
	To              *Field
	// FieldPairs carries the per-field conversion metadata.
	// The mapping needs to be built from a PairedStructs definition with WithConvertGeneration() enabled.
	// Otherwise, the old placeholder is used.
	FieldPairs []FieldPair
	// SkipConvert is set by preprocessDefinition when another Mapping with the same From.Name has
	// already been scheduled for emission in the same interface. Guards and convert bodies are suppressed.
	SkipConvert bool
	// AdditionalConvert enforces the additional conversion invocation (for cases where we have plain only field without manual conversion for any other field)
	AdditionalConvert bool
}

// AddFieldPairs converts the paired field definitions into a slice of FieldPair values used in conversion generation.
func (m *Mapping) AddFieldPairs(p *PairedStructs) {
	pairs := make([]FieldPair, 0, len(p.fields))
	for _, f := range p.fields {
		if f.plainOnly {
			m.AdditionalConvert = true
			continue
		}
		pairs = append(pairs, FieldPair{
			DbFieldName:    f.resolvedDbFieldName(),
			PlainFieldName: f.resolvedPlainFieldName(),
			DbKind:         f.dbKind,
			PlainKind:      f.plainKind,
			IsEnum:         f.isEnum,
			IsJson:         f.isJson,
			CustomParser:   f.customParser,
			ValueAdjuster:  f.valueAdjuster,
			BoolTrueValue:  f.boolTrueValue,
			BoolParsed:     f.boolParsed,
			manualConvert:  f.manualConvert,
		})
	}
	m.FieldPairs = pairs
}

// HasManualConvert reports whether any field in this Mapping is marked as manual convert.
// When true, the generated convert() will call r.additionalConvert(result) after all generated mappings.
func (m *Mapping) HasManualConvert() bool {
	if m.AdditionalConvert {
		return true
	}
	for _, f := range m.FieldPairs {
		if f.manualConvert {
			return true
		}
	}
	return false
}

func newOperation(kind string, doc string) *Operation {
	return &Operation{
		Name:          kind,
		Doc:           doc,
		HelperStructs: make([]*Field, 0),
	}
}

func newMapping(mappingFuncName string, from, to *Field) *Mapping {
	return &Mapping{
		MappingFuncName: mappingFuncName,
		From:            from,
		To:              to,
	}
}

func (s *Operation) withOptionsStruct(optsField *Field) *Operation {
	s.OptsField = optsField
	return s
}

func (s *Operation) withHelperStruct(helperStruct *Field) *Operation {
	s.HelperStructs = append(s.HelperStructs, helperStruct)
	return s
}

func (s *Operation) withHelperStructs(helperStructs ...*Field) *Operation {
	s.HelperStructs = append(s.HelperStructs, helperStructs...)
	return s
}

func (s *Operation) withScalarReturnType(scalarReturnType string) *Operation {
	s.InstanceMethodScalarReturnType = scalarReturnType
	return s
}

func (s *Operation) withObjectInterface(objectInterface *Interface) *Operation {
	s.ObjectInterface = objectInterface
	return s
}

func addShowMapping(op *Operation, from, to *Field) {
	op.ShowMapping = newMapping("convert", from, to)
}

func addDescriptionMapping(op *Operation, from, to *Field) {
	op.DescribeMapping = newMapping("convert", from, to)
}

func newNoSqlOperation(kind string) *Operation {
	operation := newOperation(kind, "placeholder").
		withOptionsStruct(nil)
	return operation
}

// TODO [next PRs]: add functional options to modify the operation on creation
func (i *Interface) newSimpleOperation(kind string, doc string, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	f := make([]*Field, len(helperStructs))
	if len(f) > 0 {
		for i, hs := range helperStructs {
			f[i] = hs.IntoField()
		}
	}
	operation := newOperation(kind, doc).
		withOptionsStruct(queryStruct.IntoField()).
		withHelperStructs(f...)
	i.Operations = append(i.Operations, operation)
	return i
}

func (i *Interface) newSimpleScalarOperation(kind string, doc string, queryStruct *QueryStruct, scalarReturnType string, helperStructs ...IntoField) *Interface {
	f := make([]*Field, len(helperStructs))
	if len(f) > 0 {
		for i, hs := range helperStructs {
			f[i] = hs.IntoField()
		}
	}
	operation := newOperation(kind, doc).
		withOptionsStruct(queryStruct.IntoField()).
		withHelperStructs(f...).
		withScalarReturnType(scalarReturnType)
	i.Operations = append(i.Operations, operation)
	return i
}

func (i *Interface) newOperationWithDBMapping(
	kind string,
	doc string,
	dbRepresentation *dbStruct,
	resourceRepresentation *plainStruct,
	queryStruct *QueryStruct,
	addMappingFunc func(op *Operation, from, to *Field),
	helperStructs ...IntoField,
) *Operation {
	db := dbRepresentation.IntoField()
	res := resourceRepresentation.IntoField()
	f := make([]*Field, len(helperStructs))
	if len(f) > 0 {
		for i, hs := range helperStructs {
			f[i] = hs.IntoField()
		}
	}
	op := newOperation(kind, doc).
		withHelperStruct(db).
		withHelperStruct(res).
		withOptionsStruct(queryStruct.IntoField()).
		withHelperStructs(f...)
	addMappingFunc(op, db, res)
	i.Operations = append(i.Operations, op)
	return op
}

type IntoField interface {
	IntoField() *Field
}

func (i *Interface) CreateOperation(doc string, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.newSimpleOperation(string(OperationKindCreate), doc, queryStruct, helperStructs...)
}

func (i *Interface) AlterOperation(doc string, queryStruct *QueryStruct) *Interface {
	return i.newSimpleOperation(string(OperationKindAlter), doc, queryStruct)
}

// DropOperationOption is a functional option for configuring DropSafely behavior.
type DropOperationOption func(*Operation)

// WithDropSafelyHook enables running action before dropping the objects safely.
func WithDropSafelyHook() DropOperationOption {
	return func(op *Operation) { op.DropSafelyHook = true }
}

// WithDropSafelyForce adds .WithForce(true) to the Drop request in the generated DropSafely implementation.
func WithDropSafelyForce() DropOperationOption {
	return func(op *Operation) { op.DropSafelyForce = true }
}

func (i *Interface) DropOperation(doc string, queryStruct *QueryStruct, opts ...DropOperationOption) *Interface {
	i.newSimpleOperation(string(OperationKindDrop), doc, queryStruct)
	// TODO [next PRs]: the operation is obtained in an ugly way; the new operation provate helpers should be made consistent and return operation instead of an interface
	op := i.Operations[len(i.Operations)-1]
	for _, opt := range opts {
		opt(op)
	}
	return i
}

func (i *Interface) GrantOperation(doc string, queryStruct *QueryStruct) *Interface {
	return i.newSimpleOperation(string(OperationKindGrant), doc, queryStruct)
}

func (i *Interface) RevokeOperation(doc string, queryStruct *QueryStruct) *Interface {
	return i.newSimpleOperation(string(OperationKindRevoke), doc, queryStruct)
}

func (i *Interface) appendShowByID(filtering []ShowByIDFilteringKind) *Interface {
	if len(filtering) == 1 && filtering[0] == ShowByIDSuppressed {
		return i
	}
	if len(filtering) == 1 && filtering[0] == ShowByIDNoFiltering {
		return i.ShowByIdOperationWithNoFiltering()
	}
	if len(filtering) == 0 {
		return i.ShowByIdOperationWithFiltering(ShowByIDLikeFiltering)
	}
	return i.ShowByIdOperationWithFiltering(filtering[0], filtering[1:]...)
}

func (i *Interface) ShowOperation(doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *QueryStruct, filtering ...ShowByIDFilteringKind) *Interface {
	return i.showOperation(doc, dbRepresentation, resourceRepresentation, queryStruct, addShowMapping, filtering...)
}

func (i *Interface) showOperation(doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *QueryStruct, addMappingFunc func(op *Operation, from, to *Field), filtering ...ShowByIDFilteringKind) *Interface {
	op := i.newOperationWithDBMapping(string(OperationKindShow), doc, dbRepresentation, resourceRepresentation, queryStruct, addMappingFunc)
	kind := ShowMappingKindSlice
	i.ShowObjectName = op.ShowMapping.To.Name
	op.ShowKind = &kind
	return i.appendShowByID(filtering)
}

func (i *Interface) CustomShowOperation(operationName string, showKind ShowMappingKind, doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *QueryStruct) *Interface {
	return i.customShowOperation(operationName, showKind, doc, dbRepresentation, resourceRepresentation, queryStruct, addShowMapping)
}

func (i *Interface) customShowOperation(operationName string, showKind ShowMappingKind, doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *QueryStruct, addMappingFunc func(op *Operation, from, to *Field), helperStructs ...IntoField) *Interface {
	op := i.newOperationWithDBMapping(operationName, doc, dbRepresentation, resourceRepresentation, queryStruct, addMappingFunc, helperStructs...)
	op.ShowKind = &showKind
	return i
}

// ShowByIdOperationWithNoFiltering adds a ShowByID operation to the interface without any filtering. Should be used for objects that do not implement any filtering options.
func (i *Interface) ShowByIdOperationWithNoFiltering() *Interface {
	op := newNoSqlOperation(string(OperationKindShowByID))
	i.Operations = append(i.Operations, op)
	return i
}

// ShowByIdOperationWithFiltering adds a ShowByID operation to the interface with filtering. Should be used for objects that implement filtering options e.g. Like or In.
func (i *Interface) ShowByIdOperationWithFiltering(filter ShowByIDFilteringKind, filtering ...ShowByIDFilteringKind) *Interface {
	op := newNoSqlOperation(string(OperationKindShowByID)).
		withObjectInterface(i).
		withFiltering(append(filtering, filter)...)
	i.Operations = append(i.Operations, op)
	return i
}

func (i *Interface) DescribeOperation(describeKind DescriptionMappingKind, doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.describeOperation(describeKind, doc, dbRepresentation, resourceRepresentation, queryStruct, addDescriptionMapping, helperStructs...)
}

func (i *Interface) describeOperation(describeKind DescriptionMappingKind, doc string, dbRepresentation *dbStruct, resourceRepresentation *plainStruct, queryStruct *QueryStruct, addMappingFunc func(op *Operation, from, to *Field), helperStructs ...IntoField) *Interface {
	op := i.newOperationWithDBMapping(string(OperationKindDescribe), doc, dbRepresentation, resourceRepresentation, queryStruct, addMappingFunc, helperStructs...)
	op.DescribeKind = &describeKind
	return i
}

// CustomOperationOption is a functional option for configuring CustomOperation behavior.
type CustomOperationOption func(*Operation)

// WithRequestAdjust enables a request.adjust() call in the generated impl before toOpts() and validateAndExec.
// The adjust() method must be implemented manually on the request type in a _ext.go file.
func WithRequestAdjust() CustomOperationOption {
	return func(op *Operation) { op.RequestAdjust = true }
}

func (i *Interface) CustomOperation(kind string, doc string, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.newSimpleOperation(kind, doc, queryStruct, helperStructs...)
}

// CustomOperationWithOpts creates a custom operation with functional options for configuring behavior.
// Use this instead of CustomOperation when you need to attach options like WithRequestAdjust.
func (i *Interface) CustomOperationWithOpts(kind string, doc string, queryStruct *QueryStruct, opts []CustomOperationOption, helperStructs ...IntoField) *Interface {
	i.newSimpleOperation(kind, doc, queryStruct, helperStructs...)
	op := i.Operations[len(i.Operations)-1]
	for _, opt := range opts {
		opt(op)
	}
	return i
}

func (i *Interface) ShowParameters(identifierKind string) *Interface {
	objectIdentifierKind, err := ToObjectIdentifierKind(identifierKind)
	if err != nil {
		panic(fmt.Errorf("invalid identifier kind: %s", identifierKind))
	}
	return i.WithCustomInterfaceMethod(
		"ShowParameters",
		"",
		[]*MethodParameter{NewMethodParameter("id", string(objectIdentifierKind))},
		"[]*Parameter", "error",
	)
}
