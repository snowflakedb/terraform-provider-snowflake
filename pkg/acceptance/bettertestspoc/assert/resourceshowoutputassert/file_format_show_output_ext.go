package resourceshowoutputassert

func (f *FileFormatShowOutputAssert) HasCreatedOnNotEmpty() *FileFormatShowOutputAssert {
	f.ValuePresent("created_on")
	return f
}

func (f *FileFormatShowOutputAssert) HasFormatOptionsNotEmpty() *FileFormatShowOutputAssert {
	f.ValuePresent("format_options")
	return f
}
