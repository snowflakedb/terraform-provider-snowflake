package gen

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
}

// HasManualConvert reports whether any field in this Mapping is marked as manual convert.
// When true, the generated convert() will call r.additionalConvert(result) after all generated mappings.
func (m *Mapping) HasManualConvert() bool {
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
	if queryStruct.identifierField != nil {
		queryStruct.identifierField.Kind = i.IdentifierKind
	}
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
	if queryStruct.identifierField != nil {
		queryStruct.identifierField.Kind = i.IdentifierKind
	}
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
	if queryStruct.identifierField != nil {
		queryStruct.identifierField.Kind = i.IdentifierKind
	}
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

func (i *Interface) DropOperation(doc string, queryStruct *QueryStruct) *Interface {
	return i.newSimpleOperation(string(OperationKindDrop), doc, queryStruct)
}

func (i *Interface) GrantOperation(doc string, queryStruct *QueryStruct) *Interface {
	return i.newSimpleOperation(string(OperationKindGrant), doc, queryStruct)
}

func (i *Interface) RevokeOperation(doc string, queryStruct *QueryStruct) *Interface {
	return i.newSimpleOperation(string(OperationKindRevoke), doc, queryStruct)
}

func (i *Interface) appendShowByID(filtering []ShowByIDFilteringKind) *Interface {
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

func (i *Interface) CustomOperation(kind string, doc string, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.newSimpleOperation(kind, doc, queryStruct, helperStructs...)
}
