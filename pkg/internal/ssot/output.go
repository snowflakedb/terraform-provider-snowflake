package ssot

// OutputSpec is SHOW or DESCRIBE columns. Attributes is CREATE/ALTER sequence;
// Columns is the full output list (membership and Snowflake order). They often
// differ. Columns holds pointers to attribute vars so names are not restated as
// strings (those vars are also how object variants reuse fields). Readonly on
// Attribute marks a non-writable output column.
type OutputSpec struct {
	Columns []*Attribute
}

func (o Object) ShowOutputKeys() []string {
	return outputKeys(o.Show.Columns)
}

func (o Object) DescribeOutputKeys() []string {
	return outputKeys(o.Describe.Columns)
}

func (o Object) ShowOutputAttributes() []*Attribute {
	return o.Show.Columns
}

func (o Object) DescribeOutputAttributes() []*Attribute {
	return o.Describe.Columns
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
