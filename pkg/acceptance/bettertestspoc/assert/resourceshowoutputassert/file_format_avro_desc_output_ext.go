package resourceshowoutputassert

func (f *FileFormatAvroDescribeOutputAssert) HasNullIf(expected ...string) *FileFormatAvroDescribeOutputAssert {
	f.ListContainsExactlyStringValuesInOrder("null_if", expected...)
	return f
}
