package gen

import "fmt"

// DescribeDetailsConfig holds the data needed to generate one Describe details accessor.
type DescribeDetailsConfig struct {
	// MethodName is the interface method, e.g. "DescribeDetails" or "DescribeAwsGlueDetails".
	MethodName string
	// DetailsType is the parsed details struct, e.g. "ExternalVolumeDetails".
	DetailsType string
	// ParseFunc parses the Describe rows into the details: func(rows []<NameSingular>Property) (<DetailsType>, error).
	// It is hand-written in the *_ext.go file.
	ParseFunc string
}

// WithDescribeDetails declares a DescribeDetails accessor on the interface and generates its body:
// it calls Describe, parses the rows with parseFunc, and sets Id on the result.
func (i *Interface) WithDescribeDetails(detailsType, parseFunc string) *Interface {
	return i.WithDescribeDetailsMethod("DescribeDetails", detailsType, parseFunc)
}

// WithDescribeDetailsMethod is WithDescribeDetails for a method with another name, so that one interface can have the
// details of several kinds, such as a method per variant. Each method has its own details type and parse function.
func (i *Interface) WithDescribeDetailsMethod(methodName, detailsType, parseFunc string) *Interface {
	i.DescribeDetailsMethods = append(i.DescribeDetailsMethods, &DescribeDetailsConfig{MethodName: methodName, DetailsType: detailsType, ParseFunc: parseFunc})
	return i.WithCustomInterfaceMethod(
		methodName,
		"",
		[]*MethodParameter{NewMethodParameter("id", i.IdentifierKind)},
		"*"+detailsType, "error",
	)
}

// DescribeParseConfig holds the data needed to generate Parse<NameSingular>Described.
type DescribeParseConfig struct {
	// RowsType is the Describe row struct, e.g. "ExternalVolumeProperty".
	RowsType string
	// PropertyField and ValueField are the row struct fields that hold the property name and its value.
	// Empty defaults to "Property" and "PropertyValue", which are the camel-cased property and property_value columns.
	PropertyField string
	ValueField    string
	// DetailsType is the parsed details struct, e.g. "ExternalVolumeDetails".
	DetailsType string
	// ErrorLabel names the object in the unrecognized property error, e.g. "external volume".
	ErrorLabel string
	// ValidateFunc is the hand-written check of the parsed details: func(<DetailsType>) error.
	ValidateFunc string
	// IgnoreUnknown skips a property with no mapping and logs a warning, instead of returning an error.
	IgnoreUnknown bool
	// RowsField and RowsFieldType keep the original rows in the details: details.<RowsField> = <RowsFieldType>(rows). Empty is off.
	RowsField     string
	RowsFieldType string
	// Properties are the exact property matches, in order.
	Properties []DescribeParseProperty
	// Prefixes are the nested list matches, in order, after the exact ones.
	Prefixes []DescribeParsePrefix
}

// DescribeParseProperty reads one string field from the row whose property name is Property.
type DescribeParseProperty struct {
	Property string
	Field    string
	// Parser is the hand-written func(value string) (<field type>, error). Empty reads the value as it is.
	Parser string
}

// DescribeParsePrefix reads a list field from every row whose property name contains Prefix.
type DescribeParsePrefix struct {
	Prefix string
	Field  string
	// Var is the local slice variable, the field name with a lowercase first letter.
	Var         string
	ElementType string
	// ElementParser is the hand-written func(value string) (<ElementType>, error).
	ElementParser string
}

// DescribeDetailsFromIdentifierConfig drives the columnar DescribeDetails body.
// It calls Describe and copies identifier fields. It does not parse property rows or set Id.
type DescribeDetailsFromIdentifierConfig struct {
	// DetailsType is the DESCRIBE plain struct, e.g. "ProjectionPolicyDetails".
	DetailsType string
	// Fields are plain-only identifier fields copied from id, in assignment order.
	// Each field name is also the identifier method: DatabaseName calls id.DatabaseName().
	Fields []string
}

// WithDescribeDetailsFromIdentifier declares DescribeDetails and generates its body:
// it calls Describe and copies each field from the identifier.
func (i *Interface) WithDescribeDetailsFromIdentifier(detailsType string, fields []string) *Interface {
	if detailsType == "" || len(fields) == 0 {
		panic("DescribeDetails from identifier needs a details type and a field")
	}
	i.DescribeDetailsFromIdentifier = &DescribeDetailsFromIdentifierConfig{
		DetailsType: detailsType,
		Fields:      fields,
	}
	return i.WithCustomInterfaceMethod(
		"DescribeDetails",
		describeDetailsFromIdentifierDoc(fields),
		[]*MethodParameter{NewMethodParameter("id", i.IdentifierKind)},
		"*"+detailsType, "error",
	)
}

func describeDetailsFromIdentifierDoc(fields []string) string {
	switch {
	case len(fields) == 2 && fields[0] == "DatabaseName" && fields[1] == "SchemaName":
		return "DescribeDetails calls Describe and sets database and schema from the identifier. DESCRIBE does not return those columns."
	case len(fields) == 1 && fields[0] == "DatabaseName":
		return "DescribeDetails calls Describe and sets database from the identifier. DESCRIBE does not return that column."
	case len(fields) == 1 && fields[0] == "SchemaName":
		return "DescribeDetails calls Describe and sets schema from the identifier. DESCRIBE does not return that column."
	default:
		panic(fmt.Sprintf("DescribeDetails from identifier has fields %v", fields))
	}
}

// WithDescribeParse generates Parse<NameSingular>Described from the config. It is the function WithDescribeDetails calls.
func (i *Interface) WithDescribeParse(config DescribeParseConfig) *Interface {
	if config.PropertyField == "" {
		config.PropertyField = "Property"
	}
	if config.ValueField == "" {
		config.ValueField = "PropertyValue"
	}
	i.DescribeParse = &config
	return i
}
