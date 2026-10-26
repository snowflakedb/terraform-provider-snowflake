package gen

// DescribeDetailsConfig holds the data needed to generate the DescribeDetails accessor.
type DescribeDetailsConfig struct {
	// DetailsType is the parsed details struct, e.g. "ExternalVolumeDetails".
	DetailsType string
	// ParseFunc parses the Describe rows into the details: func(rows []<NameSingular>Property) (<DetailsType>, error).
	// It is hand-written in the *_ext.go file.
	ParseFunc string
}

// WithDescribeDetails declares a DescribeDetails accessor on the interface and generates its body:
// it calls Describe, parses the rows with parseFunc, and sets Id on the result.
func (i *Interface) WithDescribeDetails(detailsType, parseFunc string) *Interface {
	i.DescribeDetailsMethod = &DescribeDetailsConfig{DetailsType: detailsType, ParseFunc: parseFunc}
	return i.WithCustomInterfaceMethod(
		"DescribeDetails",
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
	// Properties are the exact property matches, in order.
	Properties []DescribeParseProperty
	// Prefixes are the nested list matches, in order, after the exact ones.
	Prefixes []DescribeParsePrefix
}

// DescribeParseProperty reads one string field from the row whose property name is Property.
type DescribeParseProperty struct {
	Property string
	Field    string
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
