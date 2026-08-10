package sdk

func init() {
	id := stagesTestIdSchemaObjectIdentifier
	ffId := randomSchemaObjectIdentifier()
	integrationId := NewAccountObjectIdentifier("integration")
	tagId := NewAccountObjectIdentifier("tag-name")
	tagId2 := NewAccountObjectIdentifier("tag-name2")
	renameTarget := randomSchemaObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifier()

	stagesTests.CreateInternal.
		withModify(
			case_Stages_validation_CreateInternal_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_FileFormatOptions_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FileFormatOptions: &FileFormatOptions{}}
			},
			errExactlyOneOf("FileFormat", "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_FileFormatOptions_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions:  &FileFormatCsvOptions{},
						JsonOptions: &FileFormatJsonOptions{},
					},
				}
			},
			errExactlyOneOf("FileFormat", "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_SkipHeader_ParseHeader_ConflictingFields",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							SkipHeader:  new(1),
							ParseHeader: new(true),
						},
					},
				}
			},
			errOneOf("CsvOptions", "SkipHeader", "ParseHeader"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_RecordDelimiter_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{RecordDelimiter: &StageFileFormatStringOrNone{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.RecordDelimiter", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_RecordDelimiter_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							RecordDelimiter: &StageFileFormatStringOrNone{Value: new("\\n"), None: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.RecordDelimiter", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_FieldDelimiter_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{FieldDelimiter: &StageFileFormatStringOrNone{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.FieldDelimiter", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_FieldDelimiter_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							FieldDelimiter: &StageFileFormatStringOrNone{Value: new(","), None: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.FieldDelimiter", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_DateFormat_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{DateFormat: &StageFileFormatStringOrAuto{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.DateFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_DateFormat_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							DateFormat: &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD"), Auto: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.DateFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_TimeFormat_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{TimeFormat: &StageFileFormatStringOrAuto{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.TimeFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_TimeFormat_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							TimeFormat: &StageFileFormatStringOrAuto{Value: new("HH24:MI:SS"), Auto: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.TimeFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_TimestampFormat_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{TimestampFormat: &StageFileFormatStringOrAuto{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.TimestampFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_TimestampFormat_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							TimestampFormat: &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD HH24:MI:SS"), Auto: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.TimestampFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_Escape_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{Escape: &StageFileFormatStringOrNone{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.Escape", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_Escape_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							Escape: &StageFileFormatStringOrNone{Value: new("\\"), None: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.Escape", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_EscapeUnenclosedField_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{EscapeUnenclosedField: &StageFileFormatStringOrNone{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.EscapeUnenclosedField", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_EscapeUnenclosedField_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							EscapeUnenclosedField: &StageFileFormatStringOrNone{Value: new("\\"), None: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.EscapeUnenclosedField", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_FieldOptionallyEnclosedBy_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{FieldOptionallyEnclosedBy: &StageFileFormatStringOrNone{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.FieldOptionallyEnclosedBy", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Csv_FieldOptionallyEnclosedBy_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							FieldOptionallyEnclosedBy: &StageFileFormatStringOrNone{Value: new("\""), None: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.CsvOptions.FieldOptionallyEnclosedBy", "Value", "None"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_IgnoreUtf8Errors_ReplaceInvalidCharacters_ConflictingFields",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{IgnoreUtf8Errors: new(true), ReplaceInvalidCharacters: new(true)},
					},
				}
			},
			errOneOf("JsonOptions", "IgnoreUtf8Errors", "ReplaceInvalidCharacters"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_DateFormat_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{DateFormat: &StageFileFormatStringOrAuto{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.JsonOptions.DateFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_DateFormat_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{
							DateFormat: &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD"), Auto: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.JsonOptions.DateFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_TimeFormat_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{TimeFormat: &StageFileFormatStringOrAuto{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.JsonOptions.TimeFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_TimeFormat_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{
							TimeFormat: &StageFileFormatStringOrAuto{Value: new("HH24:MI:SS"), Auto: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.JsonOptions.TimeFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_TimestampFormat_ExactlyOneValueSet_NoneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{TimestampFormat: &StageFileFormatStringOrAuto{}},
					},
				}
			},
			errExactlyOneOf("FileFormat.JsonOptions.TimestampFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Json_TimestampFormat_ExactlyOneValueSet_MoreThanOneSet",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{
							TimestampFormat: &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD HH24:MI:SS"), Auto: new(true)},
						},
					},
				}
			},
			errExactlyOneOf("FileFormat.JsonOptions.TimestampFormat", "Value", "Auto"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Parquet_Compression_SnappyCompression_ConflictingFields",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						ParquetOptions: &FileFormatParquetOptions{Compression: new(ParquetCompressionSnappy), SnappyCompression: new(true)},
					},
				}
			},
			errOneOf("ParquetOptions", "Compression", "SnappyCompression"),
		).
		withAdditionalValidationCase(
			"validation_CreateInternal_Xml_IgnoreUtf8Errors_ReplaceInvalidCharacters_ConflictingFields",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						XmlOptions: &FileFormatXmlOptions{IgnoreUtf8Errors: new(true), ReplaceInvalidCharacters: new(true)},
					},
				}
			},
			errOneOf("XmlOptions", "IgnoreUtf8Errors", "ReplaceInvalidCharacters"),
		).
		withExpectedSqlf(
			case_Stages_sql_CreateInternal_basic,
			"CREATE STAGE %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_CreateInternal_all,
			func(opts *CreateInternalStageOptions) {
				opts.Temporary = new(true)
				opts.IfNotExists = new(true)
				opts.Encryption = &InternalStageEncryption{SnowflakeFull: &InternalStageEncryptionSnowflakeFull{}}
				opts.DirectoryTableOptions = &InternalDirectoryTableOptions{Enable: true, AutoRefresh: new(true)}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = new("some comment")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "tag-value"}}
			},
			`CREATE TEMPORARY STAGE IF NOT EXISTS %s ENCRYPTION = (TYPE = 'SNOWFLAKE_FULL') DIRECTORY = (ENABLE = true AUTO_REFRESH = true) FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment' TAG (%s = 'tag-value')`,
			id.FullyQualifiedName(), ffId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_encryptionSnowflakeSse",
			func(opts *CreateInternalStageOptions) {
				opts.Encryption = &InternalStageEncryption{SnowflakeSse: &InternalStageEncryptionSnowflakeSse{}}
			},
			`CREATE STAGE %s ENCRYPTION = (TYPE = 'SNOWFLAKE_SSE')`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_withFormatName",
			func(opts *CreateInternalStageOptions) { opts.FileFormat = &StageFileFormat{FormatName: new(ffId)} },
			`CREATE STAGE %s FILE_FORMAT = (FORMAT_NAME = %s)`, id.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_csvBasic",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FileFormatOptions: &FileFormatOptions{CsvOptions: &FileFormatCsvOptions{}}}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = CSV)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_csvAllFieldsWithSkipHeader",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							Compression:                new(CsvCompressionGzip),
							RecordDelimiter:            &StageFileFormatStringOrNone{Value: new("\\n")},
							FieldDelimiter:             &StageFileFormatStringOrNone{Value: new(",")},
							MultiLine:                  new(true),
							FileExtension:              new(".csv"),
							SkipHeader:                 new(2),
							SkipBlankLines:             new(true),
							DateFormat:                 &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD")},
							TimeFormat:                 &StageFileFormatStringOrAuto{Value: new("HH24:MI:SS")},
							TimestampFormat:            &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD HH24:MI:SS")},
							BinaryFormat:               new(BinaryFormatHex),
							Escape:                     &StageFileFormatStringOrNone{Value: new("\\")},
							EscapeUnenclosedField:      &StageFileFormatStringOrNone{Value: new("\\")},
							TrimSpace:                  new(true),
							FieldOptionallyEnclosedBy:  &StageFileFormatStringOrNone{Value: new("\"")},
							NullIf:                     &NullIfList{NullIf: []NullString{{S: "NULL"}, {S: ""}}},
							ErrorOnColumnCountMismatch: new(true),
							ReplaceInvalidCharacters:   new(true),
							EmptyFieldAsNull:           new(true),
							SkipByteOrderMark:          new(true),
							Encoding:                   new(CsvEncodingUtf8),
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = CSV COMPRESSION = GZIP RECORD_DELIMITER = '\\n' FIELD_DELIMITER = ',' MULTI_LINE = true FILE_EXTENSION = '.csv' SKIP_HEADER = 2 SKIP_BLANK_LINES = true DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH24:MI:SS' TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS' BINARY_FORMAT = HEX ESCAPE = '\\' ESCAPE_UNENCLOSED_FIELD = '\\' TRIM_SPACE = true FIELD_OPTIONALLY_ENCLOSED_BY = '\"' NULL_IF = ('NULL', '') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = true EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true ENCODING = UTF8)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_csvStringOrNoneWithNone",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							SkipHeader:      new(1),
							RecordDelimiter: &StageFileFormatStringOrNone{None: new(true)},
							FieldDelimiter:  &StageFileFormatStringOrNone{None: new(true)},
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = CSV RECORD_DELIMITER = NONE FIELD_DELIMITER = NONE SKIP_HEADER = 1)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_csvStringOrAutoWithAuto",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						CsvOptions: &FileFormatCsvOptions{
							SkipHeader:      new(1),
							DateFormat:      &StageFileFormatStringOrAuto{Auto: new(true)},
							TimeFormat:      &StageFileFormatStringOrAuto{Auto: new(true)},
							TimestampFormat: &StageFileFormatStringOrAuto{Auto: new(true)},
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = CSV SKIP_HEADER = 1 DATE_FORMAT = AUTO TIME_FORMAT = AUTO TIMESTAMP_FORMAT = AUTO)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_jsonBasic",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FileFormatOptions: &FileFormatOptions{JsonOptions: &FileFormatJsonOptions{}}}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = JSON)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_jsonAllFieldsWithIgnoreUtf8Errors",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{
							Compression:       new(JsonCompressionGzip),
							DateFormat:        &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD")},
							TimeFormat:        &StageFileFormatStringOrAuto{Value: new("HH24:MI:SS")},
							TimestampFormat:   &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD HH24:MI:SS")},
							BinaryFormat:      new(BinaryFormatBase64),
							TrimSpace:         new(true),
							MultiLine:         new(true),
							NullIf:            &NullIfList{NullIf: []NullString{{S: "NULL"}}},
							FileExtension:     new(".json"),
							EnableOctal:       new(true),
							AllowDuplicate:    new(true),
							StripOuterArray:   new(true),
							StripNullValues:   new(true),
							IgnoreUtf8Errors:  new(true),
							SkipByteOrderMark: new(true),
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = JSON COMPRESSION = GZIP DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH24:MI:SS' TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS' BINARY_FORMAT = BASE64 TRIM_SPACE = true MULTI_LINE = true NULL_IF = ('NULL') FILE_EXTENSION = '.json' ENABLE_OCTAL = true ALLOW_DUPLICATE = true STRIP_OUTER_ARRAY = true STRIP_NULL_VALUES = true IGNORE_UTF8_ERRORS = true SKIP_BYTE_ORDER_MARK = true)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_jsonAllFieldsWithReplaceInvalidCharacters",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						JsonOptions: &FileFormatJsonOptions{
							Compression:              new(JsonCompressionBrotli),
							DateFormat:               &StageFileFormatStringOrAuto{Auto: new(true)},
							TimeFormat:               &StageFileFormatStringOrAuto{Auto: new(true)},
							TimestampFormat:          &StageFileFormatStringOrAuto{Auto: new(true)},
							BinaryFormat:             new(BinaryFormatUtf8),
							TrimSpace:                new(false),
							MultiLine:                new(false),
							NullIf:                   &NullIfList{NullIf: []NullString{{S: ""}}},
							FileExtension:            new(".jsonl"),
							EnableOctal:              new(false),
							AllowDuplicate:           new(false),
							StripOuterArray:          new(false),
							StripNullValues:          new(false),
							ReplaceInvalidCharacters: new(true),
							SkipByteOrderMark:        new(false),
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = JSON COMPRESSION = BROTLI DATE_FORMAT = AUTO TIME_FORMAT = AUTO TIMESTAMP_FORMAT = AUTO BINARY_FORMAT = UTF8 TRIM_SPACE = false MULTI_LINE = false NULL_IF = ('') FILE_EXTENSION = '.jsonl' ENABLE_OCTAL = false ALLOW_DUPLICATE = false STRIP_OUTER_ARRAY = false STRIP_NULL_VALUES = false REPLACE_INVALID_CHARACTERS = true SKIP_BYTE_ORDER_MARK = false)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_avroBasic",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FileFormatOptions: &FileFormatOptions{AvroOptions: &FileFormatAvroOptions{}}}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = AVRO)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_avroAllFields",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						AvroOptions: &FileFormatAvroOptions{
							Compression:              new(AvroCompressionGzip),
							TrimSpace:                new(true),
							ReplaceInvalidCharacters: new(true),
							NullIf:                   &NullIfList{NullIf: []NullString{{S: "NULL"}, {S: ""}}},
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = AVRO COMPRESSION = GZIP TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL', ''))`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_orcBasic",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FileFormatOptions: &FileFormatOptions{OrcOptions: &FileFormatOrcOptions{}}}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = ORC)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_orcAllFields",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						OrcOptions: &FileFormatOrcOptions{
							TrimSpace:                new(true),
							ReplaceInvalidCharacters: new(true),
							NullIf:                   &NullIfList{NullIf: []NullString{{S: "NULL"}}},
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = ORC TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL'))`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_parquetAllFieldsWithCompression",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						ParquetOptions: &FileFormatParquetOptions{
							Compression:              new(ParquetCompressionSnappy),
							BinaryAsText:             new(true),
							UseLogicalType:           new(true),
							TrimSpace:                new(true),
							UseVectorizedScanner:     new(true),
							ReplaceInvalidCharacters: new(true),
							NullIf:                   &NullIfList{NullIf: []NullString{{S: "NULL"}}},
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = PARQUET COMPRESSION = SNAPPY BINARY_AS_TEXT = true USE_LOGICAL_TYPE = true TRIM_SPACE = true USE_VECTORIZED_SCANNER = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL'))`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_parquetAllFieldsWithSnappyCompression",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						ParquetOptions: &FileFormatParquetOptions{
							SnappyCompression:        new(true),
							BinaryAsText:             new(false),
							UseLogicalType:           new(false),
							TrimSpace:                new(false),
							UseVectorizedScanner:     new(false),
							ReplaceInvalidCharacters: new(false),
							NullIf:                   &NullIfList{NullIf: []NullString{{S: ""}}},
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = PARQUET SNAPPY_COMPRESSION = true BINARY_AS_TEXT = false USE_LOGICAL_TYPE = false TRIM_SPACE = false USE_VECTORIZED_SCANNER = false REPLACE_INVALID_CHARACTERS = false NULL_IF = (''))`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_xmlBasic",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{FileFormatOptions: &FileFormatOptions{XmlOptions: &FileFormatXmlOptions{}}}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = XML)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_xmlAllFieldsWithIgnoreUtf8Errors",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						XmlOptions: &FileFormatXmlOptions{
							Compression:        new(XmlCompressionGzip),
							IgnoreUtf8Errors:   new(true),
							PreserveSpace:      new(true),
							StripOuterElement:  new(true),
							DisableAutoConvert: new(true),
							SkipByteOrderMark:  new(true),
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = XML COMPRESSION = GZIP IGNORE_UTF8_ERRORS = true PRESERVE_SPACE = true STRIP_OUTER_ELEMENT = true DISABLE_AUTO_CONVERT = true SKIP_BYTE_ORDER_MARK = true)`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateInternal_xmlAllFieldsWithReplaceInvalidCharacters",
			func(opts *CreateInternalStageOptions) {
				opts.FileFormat = &StageFileFormat{
					FileFormatOptions: &FileFormatOptions{
						XmlOptions: &FileFormatXmlOptions{
							Compression:              new(XmlCompressionBz2),
							PreserveSpace:            new(false),
							StripOuterElement:        new(false),
							DisableAutoConvert:       new(false),
							ReplaceInvalidCharacters: new(true),
							SkipByteOrderMark:        new(false),
						},
					},
				}
			},
			`CREATE STAGE %s FILE_FORMAT = (TYPE = XML COMPRESSION = BZ2 PRESERVE_SPACE = false STRIP_OUTER_ELEMENT = false DISABLE_AUTO_CONVERT = false REPLACE_INVALID_CHARACTERS = true SKIP_BYTE_ORDER_MARK = false)`,
			id.FullyQualifiedName(),
		)

	stagesTests.CreateOnS3.
		withDefaultOpts(func() *CreateOnS3StageOptions {
			return &CreateOnS3StageOptions{
				name:                id,
				ExternalStageParams: ExternalS3StageParams{Url: "s3://example.com"},
			}
		}).
		withModify(
			case_Stages_validation_CreateOnS3_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateOnS3StageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withExpectedSqlf(
			case_Stages_sql_CreateOnS3_basic,
			"CREATE STAGE %s URL = 's3://example.com'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_CreateOnS3_all,
			func(opts *CreateOnS3StageOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.ExternalStageParams = ExternalS3StageParams{
					Url:                "some url",
					AwsAccessPointArn:  new("aws-access-point-arn"),
					StorageIntegration: &integrationId,
					Encryption: &ExternalStageS3Encryption{
						AwsCse: &ExternalStageS3EncryptionAwsCse{MasterKey: "master-key"},
					},
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = new("some comment")
			},
			`CREATE OR REPLACE TEMPORARY STAGE %s URL = 'some url' AWS_ACCESS_POINT_ARN = 'aws-access-point-arn' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'AWS_CSE' MASTER_KEY = 'master-key') FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnS3_allOptionsDirectoryTableAndCredentials",
			func(opts *CreateOnS3StageOptions) {
				opts.Temporary = new(true)
				opts.IfNotExists = new(true)
				opts.ExternalStageParams = ExternalS3StageParams{
					Url:               "some url",
					AwsAccessPointArn: new("aws-access-point-arn"),
					Credentials: &ExternalStageS3Credentials{
						AwsKeyId:     new("aws-key-id"),
						AwsSecretKey: new("aws-secret-key"),
						AwsToken:     new("aws-token"),
					},
					Encryption: &ExternalStageS3Encryption{
						AwsSseKms: &ExternalStageS3EncryptionAwsSseKms{KmsKeyId: new("kms-key-id")},
					},
					UsePrivatelinkEndpoint: new(true),
				}
				opts.DirectoryTableOptions = &StageS3DirectoryTableOptions{
					Enable:          true,
					RefreshOnCreate: new(true),
					AutoRefresh:     new(true),
					AwsSnsTopic:     new("arn:aws:sns:us-west-2:123456789012:my-sns-topic"),
				}
			},
			`CREATE TEMPORARY STAGE IF NOT EXISTS %s URL = 'some url' AWS_ACCESS_POINT_ARN = 'aws-access-point-arn' CREDENTIALS = (AWS_KEY_ID = 'aws-key-id' AWS_SECRET_KEY = 'aws-secret-key' AWS_TOKEN = 'aws-token') ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = 'kms-key-id') USE_PRIVATELINK_ENDPOINT = true DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true AWS_SNS_TOPIC = 'arn:aws:sns:us-west-2:123456789012:my-sns-topic')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnS3_encryptionAwsSseS3",
			func(opts *CreateOnS3StageOptions) {
				opts.ExternalStageParams.Encryption = &ExternalStageS3Encryption{AwsSseS3: &ExternalStageS3EncryptionAwsSseS3{}}
			},
			`CREATE STAGE %s URL = 's3://example.com' ENCRYPTION = (TYPE = 'AWS_SSE_S3')`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnS3_encryptionNone",
			func(opts *CreateOnS3StageOptions) {
				opts.ExternalStageParams.Encryption = &ExternalStageS3Encryption{None: &ExternalStageS3EncryptionNone{}}
			},
			`CREATE STAGE %s URL = 's3://example.com' ENCRYPTION = (TYPE = 'NONE')`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnS3_credentialsAwsRole",
			func(opts *CreateOnS3StageOptions) {
				opts.ExternalStageParams.Credentials = &ExternalStageS3Credentials{AwsRole: new("arn:aws:iam::123456789012:role/MyRole")}
			},
			`CREATE STAGE %s URL = 's3://example.com' CREDENTIALS = (AWS_ROLE = 'arn:aws:iam::123456789012:role/MyRole')`, id.FullyQualifiedName(),
		)

	stagesTests.CreateOnGCS.
		withDefaultOpts(func() *CreateOnGCSStageOptions {
			return &CreateOnGCSStageOptions{
				name:                id,
				ExternalStageParams: ExternalGCSStageParams{Url: "gcs://example.com"},
			}
		}).
		withModify(
			case_Stages_validation_CreateOnGCS_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateOnGCSStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withExpectedSqlf(
			case_Stages_sql_CreateOnGCS_basic,
			"CREATE STAGE %s URL = 'gcs://example.com'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_CreateOnGCS_all,
			func(opts *CreateOnGCSStageOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.ExternalStageParams = ExternalGCSStageParams{
					Url:                "some url",
					StorageIntegration: integrationId,
					Encryption: &ExternalStageGCSEncryption{
						GcsSseKms: &ExternalStageGCSEncryptionGcsSseKms{KmsKeyId: new("kms-key-id")},
					},
				}
				opts.DirectoryTableOptions = &ExternalGCSDirectoryTableOptions{
					Enable:                  true,
					RefreshOnCreate:         new(true),
					AutoRefresh:             new(true),
					NotificationIntegration: new("notification-integration"),
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = new("some comment")
			},
			`CREATE OR REPLACE TEMPORARY STAGE %s URL = 'some url' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'kms-key-id') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true NOTIFICATION_INTEGRATION = 'notification-integration') FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnGCS_encryptionNone",
			func(opts *CreateOnGCSStageOptions) {
				opts.ExternalStageParams.Encryption = &ExternalStageGCSEncryption{None: &ExternalStageGCSEncryptionNone{}}
			},
			`CREATE STAGE %s URL = 'gcs://example.com' ENCRYPTION = (TYPE = 'NONE')`, id.FullyQualifiedName(),
		)

	stagesTests.CreateOnAzure.
		withDefaultOpts(func() *CreateOnAzureStageOptions {
			return &CreateOnAzureStageOptions{
				name:                id,
				ExternalStageParams: ExternalAzureStageParams{Url: "azure://example.com"},
			}
		}).
		withModify(
			case_Stages_validation_CreateOnAzure_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateOnAzureStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withExpectedSqlf(
			case_Stages_sql_CreateOnAzure_basic,
			"CREATE STAGE %s URL = 'azure://example.com'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_CreateOnAzure_all,
			func(opts *CreateOnAzureStageOptions) {
				opts.OrReplace = new(true)
				opts.Temporary = new(true)
				opts.ExternalStageParams = ExternalAzureStageParams{
					Url:                "some url",
					StorageIntegration: &integrationId,
					Encryption: &ExternalStageAzureEncryption{
						AzureCse: &ExternalStageAzureEncryptionAzureCse{MasterKey: "master-key"},
					},
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = new("some comment")
			},
			`CREATE OR REPLACE TEMPORARY STAGE %s URL = 'some url' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'AZURE_CSE' MASTER_KEY = 'master-key') FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnAzure_allOptionsDirectoryTableAndCredentials",
			func(opts *CreateOnAzureStageOptions) {
				opts.IfNotExists = new(true)
				opts.DirectoryTableOptions = &ExternalAzureDirectoryTableOptions{
					Enable:                  true,
					RefreshOnCreate:         new(true),
					AutoRefresh:             new(true),
					NotificationIntegration: new("notification-integration"),
				}
				opts.ExternalStageParams = ExternalAzureStageParams{
					Url:         "some url",
					Credentials: &ExternalStageAzureCredentials{AzureSasToken: "azure-sas-token"}, //nolint:gosec
					Encryption: &ExternalStageAzureEncryption{
						AzureCse: &ExternalStageAzureEncryptionAzureCse{MasterKey: "master-key"},
					},
					UsePrivatelinkEndpoint: new(true),
				}
			},
			`CREATE STAGE IF NOT EXISTS %s URL = 'some url' CREDENTIALS = (AZURE_SAS_TOKEN = 'azure-sas-token') ENCRYPTION = (TYPE = 'AZURE_CSE' MASTER_KEY = 'master-key') USE_PRIVATELINK_ENDPOINT = true DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true NOTIFICATION_INTEGRATION = 'notification-integration')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnAzure_encryptionNone",
			func(opts *CreateOnAzureStageOptions) {
				opts.ExternalStageParams.Encryption = &ExternalStageAzureEncryption{None: &ExternalStageAzureEncryptionNone{}}
			},
			`CREATE STAGE %s URL = 'azure://example.com' ENCRYPTION = (TYPE = 'NONE')`, id.FullyQualifiedName(),
		)

	stagesTests.CreateOnS3Compatible.
		withDefaultOpts(func() *CreateOnS3CompatibleStageOptions {
			return &CreateOnS3CompatibleStageOptions{
				name: id,
				ExternalStageParams: ExternalS3CompatibleStageParams{
					Url:      "s3://example.com",
					Endpoint: "some endpoint",
				},
			}
		}).
		withModify(
			case_Stages_validation_CreateOnS3Compatible_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateOnS3CompatibleStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withExpectedSqlf(
			case_Stages_sql_CreateOnS3Compatible_basic,
			"CREATE STAGE %s URL = 's3://example.com' ENDPOINT = 'some endpoint'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_CreateOnS3Compatible_all,
			func(opts *CreateOnS3CompatibleStageOptions) {
				opts.Temporary = new(true)
				opts.IfNotExists = new(true)
				opts.ExternalStageParams = ExternalS3CompatibleStageParams{
					Url:      "some url",
					Endpoint: "some endpoint",
					Credentials: &ExternalStageS3CompatibleCredentials{
						AwsKeyId:     "aws-key-id",
						AwsSecretKey: "aws-secret-key",
					},
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = new("some comment")
				opts.DirectoryTableOptions = &StageS3CompatibleDirectoryTableOptions{
					Enable:          true,
					RefreshOnCreate: new(true),
					AutoRefresh:     new(true),
				}
			},
			`CREATE TEMPORARY STAGE IF NOT EXISTS %s URL = 'some url' ENDPOINT = 'some endpoint' CREDENTIALS = (AWS_KEY_ID = 'aws-key-id' AWS_SECRET_KEY = 'aws-secret-key') DIRECTORY = (ENABLE = true REFRESH_ON_CREATE = true AUTO_REFRESH = true) FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), ffId.FullyQualifiedName(),
		)

	stagesTests.Alter.
		withModify(
			case_Stages_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterStageOptions) {
				opts.RenameTo = new(SchemaObjectIdentifier)
				opts.SetTags = []TagAssociation{}
			},
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_Alter_RenameTo,
			func(opts *AlterStageOptions) {
				opts.IfExists = new(true)
				opts.RenameTo = &renameTarget
			},
			"ALTER STAGE IF EXISTS %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_Alter_SetTags,
			func(opts *AlterStageOptions) {
				opts.IfExists = new(true)
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "tag-value"},
					{Name: tagId2, Value: "tag-value2"},
				}
			},
			`ALTER STAGE IF EXISTS %s SET TAG %s = 'tag-value', %s = 'tag-value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_Alter_UnsetTags,
			func(opts *AlterStageOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER STAGE %s UNSET TAG %s, %s`, id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	stagesTests.AlterInternalStage.
		withModify(
			case_Stages_validation_AlterInternalStage_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterInternalStageStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterInternalStage_basic,
			func(opts *AlterInternalStageStageOptions) {
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			"ALTER STAGE %s SET COMMENT = 'some comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterInternalStage_all,
			func(opts *AlterInternalStageStageOptions) {
				opts.IfExists = new(true)
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			"ALTER STAGE IF EXISTS %s SET FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'",
			id.FullyQualifiedName(), ffId.FullyQualifiedName(),
		)

	stagesTests.AlterExternalS3Stage.
		withModify(
			case_Stages_validation_AlterExternalS3Stage_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterExternalS3StageStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterExternalS3Stage_basic,
			func(opts *AlterExternalS3StageStageOptions) {
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			"ALTER STAGE %s SET COMMENT = 'some comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterExternalS3Stage_all,
			func(opts *AlterExternalS3StageStageOptions) {
				opts.IfExists = new(true)
				opts.ExternalStageParams = &ExternalS3StageParams{
					Url:                "some url",
					AwsAccessPointArn:  new("aws-access-point-arn"),
					StorageIntegration: &integrationId,
					Encryption:         &ExternalStageS3Encryption{None: &ExternalStageS3EncryptionNone{}},
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			`ALTER STAGE IF EXISTS %s SET URL = 'some url' AWS_ACCESS_POINT_ARN = 'aws-access-point-arn' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'NONE') FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalS3Stage_credentials",
			func(opts *AlterExternalS3StageStageOptions) {
				opts.ExternalStageParams = &ExternalS3StageParams{
					Url: "s3://example.com",
					Credentials: &ExternalStageS3Credentials{
						AwsKeyId:     new("aws-key-id"),
						AwsSecretKey: new("aws-secret-key"),
					},
				}
			},
			"ALTER STAGE %s SET URL = 's3://example.com' CREDENTIALS = (AWS_KEY_ID = 'aws-key-id' AWS_SECRET_KEY = 'aws-secret-key')",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalS3Stage_encryptionAwsSseS3",
			func(opts *AlterExternalS3StageStageOptions) {
				opts.ExternalStageParams = &ExternalS3StageParams{
					Url:        "s3://example.com",
					Encryption: &ExternalStageS3Encryption{AwsSseS3: &ExternalStageS3EncryptionAwsSseS3{}},
				}
			},
			`ALTER STAGE %s SET URL = 's3://example.com' ENCRYPTION = (TYPE = 'AWS_SSE_S3')`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalS3Stage_encryptionAwsSseKms",
			func(opts *AlterExternalS3StageStageOptions) {
				opts.ExternalStageParams = &ExternalS3StageParams{
					Url: "s3://example.com",
					Encryption: &ExternalStageS3Encryption{
						AwsSseKms: &ExternalStageS3EncryptionAwsSseKms{KmsKeyId: new("kms-key-id")},
					},
				}
			},
			`ALTER STAGE %s SET URL = 's3://example.com' ENCRYPTION = (TYPE = 'AWS_SSE_KMS' KMS_KEY_ID = 'kms-key-id')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalS3Stage_credentialsAwsRole",
			func(opts *AlterExternalS3StageStageOptions) {
				opts.ExternalStageParams = &ExternalS3StageParams{
					Url:         "s3://example.com",
					Credentials: &ExternalStageS3Credentials{AwsRole: new("arn:aws:iam::123456789012:role/MyRole")},
				}
			},
			`ALTER STAGE %s SET URL = 's3://example.com' CREDENTIALS = (AWS_ROLE = 'arn:aws:iam::123456789012:role/MyRole')`,
			id.FullyQualifiedName(),
		)

	stagesTests.AlterExternalGCSStage.
		withModify(
			case_Stages_validation_AlterExternalGCSStage_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterExternalGCSStageStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterExternalGCSStage_basic,
			func(opts *AlterExternalGCSStageStageOptions) {
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			"ALTER STAGE %s SET COMMENT = 'some comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterExternalGCSStage_all,
			func(opts *AlterExternalGCSStageStageOptions) {
				opts.IfExists = new(true)
				opts.ExternalStageParams = &ExternalGCSStageParams{
					Url:                "some url",
					StorageIntegration: integrationId,
					Encryption:         &ExternalStageGCSEncryption{None: &ExternalStageGCSEncryptionNone{}},
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			`ALTER STAGE IF EXISTS %s SET URL = 'some url' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'NONE') FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalGCSStage_encryptionGcsSseKms",
			func(opts *AlterExternalGCSStageStageOptions) {
				opts.ExternalStageParams = &ExternalGCSStageParams{
					Url:                "gcs://example.com",
					StorageIntegration: integrationId,
					Encryption: &ExternalStageGCSEncryption{
						GcsSseKms: &ExternalStageGCSEncryptionGcsSseKms{KmsKeyId: new("kms-key-id")},
					},
				}
			},
			`ALTER STAGE %s SET URL = 'gcs://example.com' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'GCS_SSE_KMS' KMS_KEY_ID = 'kms-key-id')`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(),
		)

	stagesTests.AlterExternalAzureStage.
		withModify(
			case_Stages_validation_AlterExternalAzureStage_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterExternalAzureStageStageOptions) {
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId), FileFormatOptions: &FileFormatOptions{}}
			},
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterExternalAzureStage_basic,
			func(opts *AlterExternalAzureStageStageOptions) {
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			"ALTER STAGE %s SET COMMENT = 'some comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterExternalAzureStage_all,
			func(opts *AlterExternalAzureStageStageOptions) {
				opts.IfExists = new(true)
				opts.ExternalStageParams = &ExternalAzureStageParams{
					Url:                "some url",
					StorageIntegration: &integrationId,
					Encryption:         &ExternalStageAzureEncryption{None: &ExternalStageAzureEncryptionNone{}},
				}
				opts.FileFormat = &StageFileFormat{FormatName: new(ffId)}
				opts.Comment = &StringAllowEmpty{Value: "some comment"}
			},
			`ALTER STAGE IF EXISTS %s SET URL = 'some url' STORAGE_INTEGRATION = %s ENCRYPTION = (TYPE = 'NONE') FILE_FORMAT = (FORMAT_NAME = %s) COMMENT = 'some comment'`,
			id.FullyQualifiedName(), integrationId.FullyQualifiedName(), ffId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalAzureStage_credentials",
			func(opts *AlterExternalAzureStageStageOptions) {
				opts.ExternalStageParams = &ExternalAzureStageParams{
					Url:         "azure://example.com",
					Credentials: &ExternalStageAzureCredentials{AzureSasToken: "azure-sas-token"}, //nolint:gosec
				}
			},
			"ALTER STAGE %s SET URL = 'azure://example.com' CREDENTIALS = (AZURE_SAS_TOKEN = 'azure-sas-token')",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterExternalAzureStage_encryptionAzureCse",
			func(opts *AlterExternalAzureStageStageOptions) {
				opts.ExternalStageParams = &ExternalAzureStageParams{
					Url: "azure://example.com",
					Encryption: &ExternalStageAzureEncryption{
						AzureCse: &ExternalStageAzureEncryptionAzureCse{MasterKey: "master-key"},
					},
				}
			},
			`ALTER STAGE %s SET URL = 'azure://example.com' ENCRYPTION = (TYPE = 'AZURE_CSE' MASTER_KEY = 'master-key')`,
			id.FullyQualifiedName(),
		)

	stagesTests.AlterDirectoryTable.
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterDirectoryTable_SetDirectory,
			func(opts *AlterDirectoryTableStageOptions) {
				opts.IfExists = new(true)
				opts.SetDirectory = &DirectoryTableSet{Enable: true}
			},
			`ALTER STAGE IF EXISTS %s SET DIRECTORY = (ENABLE = true)`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_AlterDirectoryTable_Refresh,
			func(opts *AlterDirectoryTableStageOptions) {
				opts.IfExists = new(true)
				opts.Refresh = &DirectoryTableRefresh{}
			},
			`ALTER STAGE IF EXISTS %s REFRESH`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterDirectoryTable_refreshAllOptions",
			func(opts *AlterDirectoryTableStageOptions) {
				opts.IfExists = new(true)
				opts.Refresh = &DirectoryTableRefresh{Subpath: new("subpath")}
			},
			`ALTER STAGE IF EXISTS %s REFRESH SUBPATH = 'subpath'`, id.FullyQualifiedName(),
		)

	stagesTests.Drop.
		withExpectedSqlf(case_Stages_sql_Drop_basic, "DROP STAGE %s", id.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_Stages_sql_Drop_all,
			func(opts *DropStageOptions) { opts.IfExists = new(true) },
			"DROP STAGE IF EXISTS %s", id.FullyQualifiedName(),
		)

	stagesTests.Describe.
		withExpectedSqlf(case_Stages_sql_Describe_basic, "DESCRIBE STAGE %s", id.FullyQualifiedName())

	stagesTests.Show.
		withExpectedSql(case_Stages_sql_Show_basic, "SHOW STAGES").
		withModifyAndExpectedSqlf(
			case_Stages_sql_Show_all,
			func(opts *ShowStageOptions) {
				opts.Like = &Like{Pattern: new("some pattern")}
				opts.In = &ExtendedIn{In: In{Schema: schemaId}}
			},
			`SHOW STAGES LIKE 'some pattern' IN SCHEMA %s`, schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_Show_Like,
			func(opts *ShowStageOptions) { opts.Like = &Like{Pattern: new("stage_pattern")} },
			`SHOW STAGES LIKE 'stage_pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Stages_sql_Show_In,
			func(opts *ShowStageOptions) { opts.In = &ExtendedIn{In: In{Schema: schemaId}} },
			`SHOW STAGES IN SCHEMA %s`, schemaId.FullyQualifiedName(),
		)
}
