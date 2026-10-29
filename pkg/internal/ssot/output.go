package ssot

// OutputSpec is SHOW or DESCRIBE columns. Attributes is CREATE/ALTER sequence;
// Columns is the full output list (membership and Snowflake order). They often
// differ. Columns holds pointers to attribute vars so names are not restated as
// strings (those vars are also how object variants reuse fields). Readonly on
// Attribute marks a non-writable output column.
type OutputSpec struct {
	Columns []*Attribute
}

// PairedShowColumns is the paired SHOW. ShowGeneric wins when it is set. Otherwise it is Show.Columns.
func (obj Object) PairedShowColumns() []*Attribute {
	if len(obj.ShowGeneric.Columns) > 0 {
		return obj.ShowGeneric.Columns
	}
	return obj.Show.Columns
}

// PairedDescribeColumns is the paired DESCRIBE row. DescribeRows wins when it is set. Otherwise it is Describe.Columns.
func (obj Object) PairedDescribeColumns() []*Attribute {
	if len(obj.DescribeRows.Columns) > 0 {
		return obj.DescribeRows.Columns
	}
	return obj.Describe.Columns
}

// DescribePlainName is the plain DESCRIBE struct name. A DESCRIBE with a raw row (DescribeRows) is a property list,
// so its rows are <Name>Property and <Name>Details stays free for the parsed details. Otherwise it is <Name>Details.
func (obj Object) DescribePlainName() string {
	if len(obj.DescribeRows.Columns) > 0 {
		return obj.Name + "Property"
	}
	return obj.Name + "Details"
}

func (obj Object) ShowOutputKeys() []string {
	return outputKeys(obj.PairedShowColumns())
}

func (obj Object) DescribeOutputKeys() []string {
	return outputKeys(obj.Describe.Columns)
}

func (obj Object) ShowOutputAttributes() []*Attribute {
	return obj.PairedShowColumns()
}

func (obj Object) DescribeOutputAttributes() []*Attribute {
	return obj.Describe.Columns
}

func outputKeys(attrs []*Attribute) []string {
	keys := make([]string, 0, len(attrs))
	for _, a := range attrs {
		if a == nil {
			continue
		}
		keys = append(keys, a.OutputKey())
	}
	return keys
}
