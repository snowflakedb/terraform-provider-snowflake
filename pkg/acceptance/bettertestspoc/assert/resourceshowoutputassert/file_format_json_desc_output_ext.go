package resourceshowoutputassert

func (f *FileFormatJsonDescribeOutputAssert) HasNullIf(expected ...string) *FileFormatJsonDescribeOutputAssert {
	f.ListContainsExactlyStringValuesInOrder("null_if", expected...)
	return f
}
