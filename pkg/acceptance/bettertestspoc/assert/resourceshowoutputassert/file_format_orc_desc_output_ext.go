package resourceshowoutputassert

func (f *FileFormatOrcDescribeOutputAssert) HasNullIf(expected ...string) *FileFormatOrcDescribeOutputAssert {
	f.ListContainsExactlyStringValuesInOrder("null_if", expected...)
	return f
}
