package resourceshowoutputassert

func (f *FileFormatCsvDescribeOutputAssert) HasNullIf(expected ...string) *FileFormatCsvDescribeOutputAssert {
	f.ListContainsExactlyStringValuesInOrder("null_if", expected...)
	return f
}
