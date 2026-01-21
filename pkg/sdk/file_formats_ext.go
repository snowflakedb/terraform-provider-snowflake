package sdk

func (opts FileFormatOptions) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.CsvOptions, opts.JsonOptions, opts.AvroOptions, opts.OrcOptions, opts.ParquetOptions, opts.XmlOptions) {
		errs = append(errs, errExactlyOneOf("FileFormat", "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions"))
	}
	if valueSet(opts.CsvOptions) {
		if !exactlyOneValueSet(opts.CsvOptions.SkipHeader, opts.CsvOptions.ParseHeader) {
			errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions", "SkipHeader", "ParseHeader"))
		}
		if valueSet(opts.CsvOptions.RecordDelimiter) {
			if !exactlyOneValueSet(opts.CsvOptions.RecordDelimiter.Value, opts.CsvOptions.RecordDelimiter.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.RecordDelimiter", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.FieldDelimiter) {
			if !exactlyOneValueSet(opts.CsvOptions.FieldDelimiter.Value, opts.CsvOptions.FieldDelimiter.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.FieldDelimiter", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.DateFormat) {
			if !exactlyOneValueSet(opts.CsvOptions.DateFormat.Value, opts.CsvOptions.DateFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.DateFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.CsvOptions.TimeFormat) {
			if !exactlyOneValueSet(opts.CsvOptions.TimeFormat.Value, opts.CsvOptions.TimeFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.TimeFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.CsvOptions.TimestampFormat) {
			if !exactlyOneValueSet(opts.CsvOptions.TimestampFormat.Value, opts.CsvOptions.TimestampFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.TimestampFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.CsvOptions.Escape) {
			if !exactlyOneValueSet(opts.CsvOptions.Escape.Value, opts.CsvOptions.Escape.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.Escape", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.EscapeUnenclosedField) {
			if !exactlyOneValueSet(opts.CsvOptions.EscapeUnenclosedField.Value, opts.CsvOptions.EscapeUnenclosedField.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.EscapeUnenclosedField", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.FieldOptionallyEnclosedBy) {
			if !exactlyOneValueSet(opts.CsvOptions.FieldOptionallyEnclosedBy.Value, opts.CsvOptions.FieldOptionallyEnclosedBy.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.FieldOptionallyEnclosedBy", "Value", "None"))
			}
		}
	}
	if valueSet(opts.JsonOptions) {
		if !exactlyOneValueSet(opts.JsonOptions.IgnoreUtf8Errors, opts.JsonOptions.ReplaceInvalidCharacters) {
			errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions", "IgnoreUtf8Errors", "ReplaceInvalidCharacters"))
		}
		if valueSet(opts.JsonOptions.DateFormat) {
			if !exactlyOneValueSet(opts.JsonOptions.DateFormat.Value, opts.JsonOptions.DateFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions.DateFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.JsonOptions.TimeFormat) {
			if !exactlyOneValueSet(opts.JsonOptions.TimeFormat.Value, opts.JsonOptions.TimeFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions.TimeFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.JsonOptions.TimestampFormat) {
			if !exactlyOneValueSet(opts.JsonOptions.TimestampFormat.Value, opts.JsonOptions.TimestampFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions.TimestampFormat", "Value", "Auto"))
			}
		}
	}
	if valueSet(opts.ParquetOptions) {
		if !exactlyOneValueSet(opts.ParquetOptions.Compression, opts.ParquetOptions.SnappyCompression) {
			errs = append(errs, errExactlyOneOf("FileFormat.ParquetOptions", "Compression", "SnappyCompression"))
		}
	}
	if valueSet(opts.XmlOptions) {
		if !exactlyOneValueSet(opts.XmlOptions.IgnoreUtf8Errors, opts.XmlOptions.ReplaceInvalidCharacters) {
			errs = append(errs, errExactlyOneOf("FileFormat.XmlOptions", "IgnoreUtf8Errors", "ReplaceInvalidCharacters"))
		}
	}
	return JoinErrors(errs...)
}
