package gen

import "strings"

// PairedFieldOption is a functional option for customizing a paired field definition.
// It enables flexible and extensible field configuration without proliferating method variants.
type PairedFieldOption func(*pairedField)

// WithDbFieldName overrides the Go field name in the db row struct for this field.
// By default, the Go field name is derived from the db column name via sqlToFieldName.
// The db tag always uses the original dbColumnName regardless of this override.
// Example:
//
//	Text("type", WithDbFieldName("StorageType"))
//	// db:    StorageType string `db:"type"`
func WithDbFieldName(name string) PairedFieldOption {
	return func(f *pairedField) {
		f.dbFieldName = name
	}
}

// WithPlainFieldName overrides the plain struct field name that would otherwise be
// auto-derived from the db column name. Use this when the plain struct field name
// should differ from the CamelCase conversion of the db column name.
// Example:
//
//	Text("organization_name", WithPlainFieldName("OrganizationBasicName"))
//	// plain: OrganizationBasicName string
func WithPlainFieldName(name string) PairedFieldOption {
	return func(f *pairedField) {
		f.plainFieldName = name
	}
}

// WithRequiredInPlain strips the pointer from the plain kind, making the plain struct field
// non-nullable even when the db field uses a nullable type. Applied after the base kind is
// set by the builder method.
// Example:
//
//	OptionalText("comment", WithRequiredInPlain())
//	// db:    Comment sql.NullString `db:"comment"`
//	// plain: Comment string
func WithRequiredInPlain() PairedFieldOption {
	return func(f *pairedField) {
		f.plainKind = strings.TrimPrefix(f.plainKind, "*")
	}
}

// WithCustomParser sets a custom parse function name to use when converting the db field to the plain field.
// The function must have signature func(string) (T, error) where T matches the plain kind.
// Example:
//
//	Field("signature", "string", "[]TableColumnSignature", WithCustomParser("ParseTableColumnSignature"))
func WithCustomParser(funcName string) PairedFieldOption {
	return func(f *pairedField) {
		f.customParser = funcName
	}
}

// pairedField holds the definition for a single field in both the DB row struct and the plain SDK struct.
type pairedField struct {
	// dbColumnName is the snake_case column name used for the db: tag and to auto-derive the field names.
	dbColumnName string
	// dbFieldName is an optional override for the Go field name in the db row struct.
	dbFieldName string
	// plainFieldName is an optional override for the plain struct field name.
	plainFieldName string
	// dbKind is the Go type used in the db row struct (e.g. "string", "sql.NullString", "bool").
	dbKind string
	// plainKind is the Go type used in the plain SDK struct (e.g. "string", "*string", "bool").
	plainKind string
	// isEnum marks that the plain field is an enum type.
	isEnum bool
	// isJson marks that the db string column should be JSON-unmarshaled into the plain field.
	isJson bool
	// customParser is the name of a custom parse function to use for conversion.
	customParser string
}

// resolvedPlainFieldName returns the explicit override or the CamelCase conversion of dbColumnName.
func (f *pairedField) resolvedPlainFieldName() string {
	if f.plainFieldName != "" {
		return f.plainFieldName
	}
	return sqlToFieldName(f.dbColumnName, true)
}

// resolvedDbFieldName returns the explicit override or the CamelCase conversion of dbColumnName.
func (f *pairedField) resolvedDbFieldName() string {
	if f.dbFieldName != "" {
		return f.dbFieldName
	}
	return sqlToFieldName(f.dbColumnName, true)
}

// PairedStructs defines the DB row struct and the plain SDK struct in a single unified definition.
// Each field is added once with explicit db and plain types, eliminating the need to maintain
// two separate DbStruct and PlainStruct definitions in sync and simplify the conversion generation.
type PairedStructs struct {
	dbName    string
	plainName string
	fields    []pairedField
	// generateConvert controls whether convert() body generation is enabled for this pair.
	generateConvert bool
}

// StructPair creates a new PairedStructs with the given DB row struct name and plain struct name.
func StructPair(dbName, plainName string) *PairedStructs {
	return &PairedStructs{
		dbName:    dbName,
		plainName: plainName,
		fields:    make([]pairedField, 0),
	}
}

func (p *PairedStructs) addField(dbColumnName, dbKind, plainKind string, opts []PairedFieldOption) *PairedStructs {
	f := pairedField{
		dbColumnName: dbColumnName,
		dbKind:       dbKind,
		plainKind:    plainKind,
	}
	for _, opt := range opts {
		opt(&f)
	}
	p.fields = append(p.fields, f)
	return p
}

// Field adds a field with fully explicit db and plain kinds. Use this for custom type pairs
// that do not fit the convenience methods (e.g. db string → plain ExternalObjectIdentifier).
func (p *PairedStructs) Field(dbColumnName, dbKind, plainKind string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, dbKind, plainKind, opts)
}

// Text adds a non-nullable string field to both the db row struct and the plain struct.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: <FieldName> string
func (p *PairedStructs) Text(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "string", "string", opts)
}

// OptionalText adds a nullable string field. The db kind is sql.NullString and the plain kind
// is *string by default. Use WithRequiredInPlain() to make the plain field non-nullable.
//
//	db:    <FieldName> sql.NullString `db:"<dbColumnName>"`
//	plain: <FieldName> *string
func (p *PairedStructs) OptionalText(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "sql.NullString", "*string", opts)
}

// Bool adds a non-nullable boolean field to both the db row struct and the plain struct.
//
//	db:    <FieldName> bool `db:"<dbColumnName>"`
//	plain: <FieldName> bool
func (p *PairedStructs) Bool(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "bool", "bool", opts)
}

// OptionalBool adds a nullable boolean field. The db kind is sql.NullBool and the plain kind
// is *bool by default. Use WithRequiredInPlain() to make the plain field non-nullable.
//
//	db:    <FieldName> sql.NullBool `db:"<dbColumnName>"`
//	plain: <FieldName> *bool
func (p *PairedStructs) OptionalBool(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "sql.NullBool", "*bool", opts)
}

// Number adds a non-nullable integer field to both the db row struct and the plain struct.
//
//	db:    <FieldName> int `db:"<dbColumnName>"`
//	plain: <FieldName> int
func (p *PairedStructs) Number(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "int", "int", opts)
}

// OptionalNumber adds a nullable integer field. The db kind is sql.NullInt64 and the plain kind
// is *int by default. Use WithRequiredInPlain() to make the plain field non-nullable.
//
//	db:    <FieldName> sql.NullInt64 `db:"<dbColumnName>"`
//	plain: <FieldName> *int
func (p *PairedStructs) OptionalNumber(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "sql.NullInt64", "*int", opts)
}

// Time adds a non-nullable time.Time field to both the db row struct and the plain struct.
//
//	db:    <FieldName> time.Time `db:"<dbColumnName>"`
//	plain: <FieldName> time.Time
func (p *PairedStructs) Time(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "time.Time", "time.Time", opts)
}

// OptionalTime adds a nullable time field. The db kind is sql.NullTime and the plain kind
// is *time.Time by default. Use WithRequiredInPlain() to make the plain field non-nullable.
//
//	db:    <FieldName> sql.NullTime `db:"<dbColumnName>"`
//	plain: <FieldName> *time.Time
func (p *PairedStructs) OptionalTime(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "sql.NullTime", "*time.Time", opts)
}

// PlainField adds a field where the db kind is string and the plain kind is a custom type.
// Use this for fields that are stored as raw strings in the db but represented as a custom
// SDK type in the plain struct (e.g. enums, identifiers, or slices parsed from the db string).
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: <FieldName> <plainKind>
func (p *PairedStructs) PlainField(dbColumnName, plainKind string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "string", plainKind, opts)
}

// StringList adds a field where the db kind is string and the plain kind is []string.
// Use this for fields that are stored as a single string in the db (e.g. comma-separated)
// but represented as a string slice in the plain struct.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: <FieldName> []string
func (p *PairedStructs) StringList(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "string", "[]string", opts)
}

// AccountObjectIdentifier adds an AccountObjectIdentifier field. The db kind is string and the plain kind
// is AccountObjectIdentifier. The plain field name defaults to "Id", but can be overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id AccountObjectIdentifier
func (p *PairedStructs) AccountObjectIdentifier(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	// Pre-apply the "Id" default; caller-supplied WithPlainFieldName will override it.
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "string", "AccountObjectIdentifier", allOpts)
}

// OptionalAccountObjectIdentifier adds a nullable AccountObjectIdentifier field. The db kind is string
// and the plain kind is *AccountObjectIdentifier. The plain field name defaults to "Id", but can be
// overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id *AccountObjectIdentifier
func (p *PairedStructs) OptionalAccountObjectIdentifier(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "sql.NullString", "*AccountObjectIdentifier", allOpts)
}

// DatabaseObjectIdentifier adds a DatabaseObjectIdentifier field. The db kind is string and the plain kind
// is DatabaseObjectIdentifier. The plain field name defaults to "Id", but can be overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id DatabaseObjectIdentifier
func (p *PairedStructs) DatabaseObjectIdentifier(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "string", "DatabaseObjectIdentifier", allOpts)
}

// SchemaObjectIdentifier adds a SchemaObjectIdentifier field. The db kind is string and the plain kind
// is SchemaObjectIdentifier. The plain field name defaults to "Id", but can be overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id SchemaObjectIdentifier
func (p *PairedStructs) SchemaObjectIdentifier(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "string", "SchemaObjectIdentifier", allOpts)
}

// OptionalSchemaObjectIdentifier adds a nullable SchemaObjectIdentifier field. The db kind is string
// and the plain kind is *SchemaObjectIdentifier. The plain field name defaults to "Id", but can be
// overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id *SchemaObjectIdentifier
func (p *PairedStructs) OptionalSchemaObjectIdentifier(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "sql.NullString", "*SchemaObjectIdentifier", allOpts)
}

// NullableSchemaObjectIdentifierArray adds a nullable SchemaObjectIdentifier slice field. The db kind is
// sql.NullString and the plain kind is []SchemaObjectIdentifier. The plain field name defaults to the
// camel-cased column name unless overridden with WithPlainFieldName.
//
//	db:    <FieldName> sql.NullString `db:"<dbColumnName>"`
//	plain: <FieldName> []SchemaObjectIdentifier
func (p *PairedStructs) NullableSchemaObjectIdentifierArray(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "sql.NullString", "[]SchemaObjectIdentifier", opts)
}

// AccountIdentifierArray adds a required AccountIdentifier slice field. The db kind is string and the
// plain kind is []AccountIdentifier.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: <FieldName> []AccountIdentifier
func (p *PairedStructs) AccountIdentifierArray(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	return p.addField(dbColumnName, "string", "[]AccountIdentifier", opts)
}

// SchemaObjectIdentifierWithArguments adds a SchemaObjectIdentifierWithArguments field. The db kind is string and the plain kind
// is SchemaObjectIdentifierWithArguments. The plain field name defaults to "Id", but can be overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id SchemaObjectIdentifierWithArguments
func (p *PairedStructs) SchemaObjectIdentifierWithArguments(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "string", "SchemaObjectIdentifierWithArguments", allOpts)
}

// OptionalSchemaObjectIdentifierWithArguments adds a nullable SchemaObjectIdentifierWithArguments field. The db kind is string
// and the plain kind is *SchemaObjectIdentifierWithArguments. The plain field name defaults to "Id", but can be
// overridden with WithPlainFieldName.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: Id *SchemaObjectIdentifierWithArguments
func (p *PairedStructs) OptionalSchemaObjectIdentifierWithArguments(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "sql.NullString", "*SchemaObjectIdentifierWithArguments", allOpts)
}

// Enum adds a field where the db kind is string and the plain kind is a custom enum type.
// Use this when the plain struct field is an SDK enum backed by a string column in the db.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: <FieldName> <enumType>
func (p *PairedStructs) Enum(dbColumnName string, enum *Enum, opts ...PairedFieldOption) *PairedStructs {
	f := pairedField{
		dbColumnName: dbColumnName,
		dbKind:       "string",
		plainKind:    enum.Kind(),
		isEnum:       true,
	}
	for _, opt := range opts {
		opt(&f)
	}
	p.fields = append(p.fields, f)
	return p
}

// OptionalEnum adds a nullable enum field. The db kind is sql.NullString and the plain kind is *<enumType>.
//
//	db:    <FieldName> sql.NullString `db:"<dbColumnName>"`
//	plain: <FieldName> *<enumType>
func (p *PairedStructs) OptionalEnum(dbColumnName string, enum *Enum, opts ...PairedFieldOption) *PairedStructs {
	f := pairedField{
		dbColumnName: dbColumnName,
		dbKind:       "sql.NullString",
		plainKind:    enum.KindPtr(),
		isEnum:       true,
	}
	for _, opt := range opts {
		opt(&f)
	}
	p.fields = append(p.fields, f)
	return p
}

// JsonField adds a field where the db kind is string and the plain kind is a custom type that will be populated via JSON unmarshalling.
//
//	db:    <FieldName> string `db:"<dbColumnName>"`
//	plain: <FieldName> <kind>
func (p *PairedStructs) JsonField(dbColumnName, kind string, opts ...PairedFieldOption) *PairedStructs {
	f := pairedField{
		dbColumnName: dbColumnName,
		dbKind:       "string",
		plainKind:    kind,
		isJson:       true,
	}
	for _, opt := range opts {
		opt(&f)
	}
	p.fields = append(p.fields, f)
	return p
}

// WithConvertGeneration opts this PairedStructs into automatic convert() body generation.
// By default, convert generation is disabled so existing PairedStructs usages in production defs continue to emit the placeholder.
func (p *PairedStructs) WithConvertGeneration() *PairedStructs {
	p.generateConvert = true
	return p
}

// toFieldPairs converts the paired field definitions into a slice of FieldPair values used in conversion generation.
func (p *PairedStructs) toFieldPairs() []FieldPair {
	pairs := make([]FieldPair, len(p.fields))
	for i, f := range p.fields {
		pairs[i] = FieldPair{
			DbFieldName:    f.resolvedDbFieldName(),
			PlainFieldName: f.resolvedPlainFieldName(),
			DbKind:         f.dbKind,
			PlainKind:      f.plainKind,
			IsEnum:         f.isEnum,
			IsJson:         f.isJson,
			CustomParser:   f.customParser,
		}
	}
	return pairs
}

// asDbStruct materializes the definition as a *dbStruct following the old implementation.
func (p *PairedStructs) asDbStruct() *dbStruct {
	s := DbStruct(p.dbName)
	for _, f := range p.fields {
		s.FieldWithName(f.dbColumnName, f.dbKind, f.resolvedDbFieldName())
	}
	return s
}

// asPlainStruct materializes the definition as a *plainStruct following the old implementation.
func (p *PairedStructs) asPlainStruct() *plainStruct {
	s := PlainStruct(p.plainName)
	for _, f := range p.fields {
		s.Field(f.resolvedPlainFieldName(), f.plainKind)
	}
	return s
}

func (p *PairedStructs) addMappingFunc() func(op *Operation, from, to *Field) {
	return func(op *Operation, from, to *Field) {
		addShowMapping(op, from, to)
		if p.generateConvert {
			op.ShowMapping.FieldPairs = p.toFieldPairs()
		}
	}
}

func (p *PairedStructs) addDescribeMappingFunc() func(op *Operation, from, to *Field) {
	return func(op *Operation, from, to *Field) {
		addDescriptionMapping(op, from, to)
		if p.generateConvert {
			op.DescribeMapping.FieldPairs = p.toFieldPairs()
		}
	}
}

func (p *PairedStructs) instanceMethodMappingFunc(kind InstanceMethodKind) func(op *Operation, from, to *Field) {
	return func(op *Operation, from, to *Field) {
		op.InstanceMethodMapping = newMapping("convert", from, to)
		op.InstanceMethodKind = &kind
		if p.generateConvert {
			op.InstanceMethodMapping.FieldPairs = p.toFieldPairs()
		}
	}
}

// ShowOperationWithPairedStructs is equivalent to ShowOperation but accepts a single PairedStructs
// definition instead of separate DbStruct and PlainStruct arguments.
func (i *Interface) ShowOperationWithPairedStructs(doc string, pairedStructs *PairedStructs, queryStruct *QueryStruct, filtering ...ShowByIDFilteringKind) *Interface {
	return i.showOperation(doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), queryStruct, pairedStructs.addMappingFunc(), filtering...)
}

// DescribeOperationWithPairedStructs is equivalent to DescribeOperation but accepts a single
// PairedStructs definition instead of separate DbStruct and PlainStruct arguments.
func (i *Interface) DescribeOperationWithPairedStructs(describeKind DescriptionMappingKind, doc string, pairedStructs *PairedStructs, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.describeOperation(describeKind, doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), queryStruct, pairedStructs.addDescribeMappingFunc(), helperStructs...)
}

// CustomShowOperationWithPairedStructs is equivalent to CustomShowOperation but accepts a
// single PairedStructs definition instead of separate DbStruct and PlainStruct arguments.
func (i *Interface) CustomShowOperationWithPairedStructs(operationName string, showKind ShowMappingKind, doc string, pairedStructs *PairedStructs, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.customShowOperation(operationName, showKind, doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), queryStruct, pairedStructs.addMappingFunc(), helperStructs...)
}
