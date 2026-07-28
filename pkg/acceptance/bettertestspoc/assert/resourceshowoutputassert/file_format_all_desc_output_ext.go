package resourceshowoutputassert

// The describe output of any file format is a union of all the per-type describe outputs, so the type-specific fields
// are asserted with the per-type assertions below. Note that a given field can be checked on a file format of any type
// (e.g. to verify that a CSV-only field is not filled for a JSON file format).

func (f *FileFormatAllDescribeOutputAssert) Csv() *FileFormatCsvDescribeOutputAssert {
	return &FileFormatCsvDescribeOutputAssert{ResourceAssert: f.ResourceAssert}
}

func (f *FileFormatAllDescribeOutputAssert) Json() *FileFormatJsonDescribeOutputAssert {
	return &FileFormatJsonDescribeOutputAssert{ResourceAssert: f.ResourceAssert}
}

func (f *FileFormatAllDescribeOutputAssert) Avro() *FileFormatAvroDescribeOutputAssert {
	return &FileFormatAvroDescribeOutputAssert{ResourceAssert: f.ResourceAssert}
}

func (f *FileFormatAllDescribeOutputAssert) Orc() *FileFormatOrcDescribeOutputAssert {
	return &FileFormatOrcDescribeOutputAssert{ResourceAssert: f.ResourceAssert}
}

func (f *FileFormatAllDescribeOutputAssert) Parquet() *FileFormatParquetDescribeOutputAssert {
	return &FileFormatParquetDescribeOutputAssert{ResourceAssert: f.ResourceAssert}
}

func (f *FileFormatAllDescribeOutputAssert) Xml() *FileFormatXmlDescribeOutputAssert {
	return &FileFormatXmlDescribeOutputAssert{ResourceAssert: f.ResourceAssert}
}
