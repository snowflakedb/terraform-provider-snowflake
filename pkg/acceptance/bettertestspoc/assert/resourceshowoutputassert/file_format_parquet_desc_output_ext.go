package resourceshowoutputassert

func (f *FileFormatParquetDescribeOutputAssert) HasNullIf(expected ...string) *FileFormatParquetDescribeOutputAssert {
	f.ListContainsExactlyStringValuesInOrder("null_if", expected...)
	return f
}
