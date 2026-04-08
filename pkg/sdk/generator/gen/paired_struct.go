package gen

import "strings"

// PairedFieldOption is a functional option for customizing a paired field definition.
// It enables flexible and extensible field configuration without proliferating method variants.
type PairedFieldOption func(*pairedField)

// WithDbFieldName overrides the Go field name in the db row struct for this field.
// By default the Go field name is derived from the db column name via sqlToFieldName.
// The db: tag always uses the original dbColumnName regardless of this override.
//
// Example:
//
//	Text("type", WithDbFieldName("StorageType"), WithPlainFieldName("StorageType"))
//	// db:    StorageType string `db:"type"`
//	// plain: StorageType string
func WithDbFieldName(name string) PairedFieldOption {
	return func(f *pairedField) {
		f.dbFieldName = name
	}
}

// WithPlainFieldName overrides the plain struct field name that would otherwise be
// auto-derived from the db column name. Use this when the plain struct field name
// should differ from the CamelCase conversion of the db column name.
//
// Example:
//
//	Text("organization_name", WithPlainFieldName("OrganizationBasicName"))
//	// db:    OrganizationName string `db:"organization_name"`
//	// plain: OrganizationBasicName string
func WithPlainFieldName(name string) PairedFieldOption {
	return func(f *pairedField) {
		f.plainFieldName = name
	}
}

// WithRequiredInPlain strips the pointer from the plain kind, making the plain struct field
// non-nullable even when the db field uses a nullable type. Applied after the base kind is
// set by the builder method.
//
// Example:
//
//	OptionalText("comment", WithRequiredInPlain())
//	// db:    Comment sql.NullString `db:"comment"`
//	// plain: Comment string
//
//	OptionalBool("enabled", WithRequiredInPlain())
//	// db:    Enabled sql.NullBool `db:"enabled"`
//	// plain: Enabled bool
func WithRequiredInPlain() PairedFieldOption {
	return func(f *pairedField) {
		f.plainKind = strings.TrimPrefix(f.plainKind, "*")
	}
}

// pairedField holds the definition for a single field in both the DB row struct and the plain SDK struct.
type pairedField struct {
	// dbColumnName is the snake_case column name used for the db: tag and to auto-derive the plain field name.
	dbColumnName string
	// dbFieldName is an optional override for the Go field name in the db row struct.
	// When empty, sqlToFieldName(dbColumnName, true) is used.
	dbFieldName string
	// plainFieldName is an optional override for the plain struct field name.
	// When empty, sqlToFieldName(dbColumnName, true) is used.
	plainFieldName string
	// dbKind is the Go type used in the db row struct (e.g. "string", "sql.NullString", "bool").
	dbKind string
	// plainKind is the Go type used in the plain SDK struct (e.g. "string", "*string", "bool").
	plainKind string
}

// resolvedPlainFieldName returns the explicit override or the CamelCase conversion of dbColumnName.
func (f *pairedField) resolvedPlainFieldName() string {
	if f.plainFieldName != "" {
		return f.plainFieldName
	}
	return sqlToFieldName(f.dbColumnName, true)
}

// PairedStructs defines the DB row struct and the plain SDK struct in a single unified definition.
// Each field is added once with explicit db and plain types, eliminating the need to maintain
// two separate DbStruct and PlainStruct definitions in sync.
//
// Usage:
//
//	g.StructPair("connectionRow", "Connection").
//	    Text("name").
//	    OptionalText("comment").
//	    Time("created_on").
//	    Bool("is_primary").
//	    Text("organization_name", g.WithPlainFieldName("OrganizationBasicName")).
//	    OptionalText("region_group", g.WithRequiredInPlain()).
//	    Field("primary", "string", "ExternalObjectIdentifier").
//	    PlainField("failover_allowed_to_accounts", "[]AccountIdentifier")
type PairedStructs struct {
	dbName    string
	plainName string
	fields    []pairedField
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

// AccountObjectIdentifier adds an AccountObjectIdentifier field. The db kind is string and
// the plain kind is AccountObjectIdentifier. The plain field name defaults to "Id" to match
// the convention of plainStruct.AccountObjectIdentifier(), but can be overridden with
// WithPlainFieldName.
//
//	db:    <GoFieldName> string `db:"<dbColumnName>"`
//	plain: Id AccountObjectIdentifier
func (p *PairedStructs) AccountObjectIdentifier(dbColumnName string, opts ...PairedFieldOption) *PairedStructs {
	// Pre-apply the "Id" default; caller-supplied WithPlainFieldName will override it.
	allOpts := append([]PairedFieldOption{WithPlainFieldName("Id")}, opts...)
	return p.addField(dbColumnName, "string", "AccountObjectIdentifier", allOpts)
}

// asDbStruct materialises the definition as a *dbStruct compatible with ShowOperation and
// DescribeOperation. When a field has a dbFieldName override, FieldWithGoName is used so
// that the generated db row struct uses the explicit Go identifier while the db: tag still
// holds the original column name.
func (p *PairedStructs) asDbStruct() *dbStruct {
	s := DbStruct(p.dbName)
	for _, f := range p.fields {
		if f.dbFieldName != "" {
			s.FieldWithGoName(f.dbColumnName, f.dbFieldName, f.dbKind)
		} else {
			s.Field(f.dbColumnName, f.dbKind)
		}
	}
	return s
}

// asPlainStruct materialises the definition as a *plainStruct compatible with ShowOperation and
// DescribeOperation. Each field uses the resolved plain field name (explicit override or
// CamelCase conversion of dbColumnName) and the plain kind.
func (p *PairedStructs) asPlainStruct() *plainStruct {
	s := PlainStruct(p.plainName)
	for _, f := range p.fields {
		s.Field(f.resolvedPlainFieldName(), f.plainKind)
	}
	return s
}

// ShowOperationWithPairedStructs is equivalent to ShowOperation but accepts a single PairedStructs
// definition instead of separate DbStruct and PlainStruct arguments. The PairedStructs is
// materialised into both structs and forwarded to the existing ShowOperation unchanged.
func (i *Interface) ShowOperationWithPairedStructs(doc string, pairedStructs *PairedStructs, queryStruct *QueryStruct) *Interface {
	return i.ShowOperation(doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), queryStruct)
}

// DescribeOperationWithPairedStructs is equivalent to DescribeOperation but accepts a single
// PairedStructs definition instead of separate DbStruct and PlainStruct arguments. The
// PairedStructs is materialised into both structs and forwarded to the existing DescribeOperation.
func (i *Interface) DescribeOperationWithPairedStructs(describeKind DescriptionMappingKind, doc string, pairedStructs *PairedStructs, queryStruct *QueryStruct, helperStructs ...IntoField) *Interface {
	return i.DescribeOperation(describeKind, doc, pairedStructs.asDbStruct(), pairedStructs.asPlainStruct(), queryStruct, helperStructs...)
}
