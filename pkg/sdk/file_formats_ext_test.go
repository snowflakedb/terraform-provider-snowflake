package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	fileFormatsTests.CreateCsv.
		withExpectedSqlf(case_FileFormats_sql_CreateCsv_basic,
			`CREATE FILE FORMAT %s TYPE = CSV`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_CreateCsv_all,
			func(opts *CreateCsvFileFormatOptions) {
				opts.IfNotExists = new(true)
				opts.Compression = new(CsvCompressionGzip)
				opts.RecordDelimiter = &StageFileFormatStringOrNone{Value: new("\\n")}
				opts.FieldDelimiter = &StageFileFormatStringOrNone{Value: new(",")}
				opts.MultiLine = new(true)
				opts.FileExtension = new(".csv")
				opts.SkipHeader = new(2)
				opts.SkipBlankLines = new(true)
				opts.DateFormat = &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD")}
				opts.TimeFormat = &StageFileFormatStringOrAuto{Value: new("HH24:MI:SS")}
				opts.TimestampFormat = &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD HH24:MI:SS")}
				opts.BinaryFormat = new(BinaryFormatHex)
				opts.Escape = &StageFileFormatStringOrNone{Value: new("\\")}
				opts.EscapeUnenclosedField = &StageFileFormatStringOrNone{Value: new("\\")}
				opts.TrimSpace = new(true)
				opts.FieldOptionallyEnclosedBy = &StageFileFormatStringOrNone{Value: new("\"")}
				opts.NullIf = &NullIfList{NullIf: []NullString{{S: "NULL"}, {S: ""}}}
				opts.ErrorOnColumnCountMismatch = new(true)
				opts.ReplaceInvalidCharacters = new(true)
				opts.EmptyFieldAsNull = new(true)
				opts.SkipByteOrderMark = new(true)
				opts.Encoding = new(CsvEncodingUtf8)
				opts.Comment = new("some comment")
			},
			`CREATE FILE FORMAT IF NOT EXISTS %s TYPE = CSV COMPRESSION = GZIP RECORD_DELIMITER = '\\n' FIELD_DELIMITER = ',' MULTI_LINE = true FILE_EXTENSION = '.csv' SKIP_HEADER = 2 SKIP_BLANK_LINES = true DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH24:MI:SS' TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS' BINARY_FORMAT = HEX ESCAPE = '\\' ESCAPE_UNENCLOSED_FIELD = '\\' TRIM_SPACE = true FIELD_OPTIONALLY_ENCLOSED_BY = '\"' NULL_IF = ('NULL', '') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = true EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true ENCODING = UTF8 COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateCsv_orReplace",
			func(opts *CreateCsvFileFormatOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE FILE FORMAT %s TYPE = CSV`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.CreateJson.
		withExpectedSqlf(case_FileFormats_sql_CreateJson_basic,
			`CREATE FILE FORMAT %s TYPE = JSON`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_CreateJson_all,
			func(opts *CreateJsonFileFormatOptions) {
				opts.IfNotExists = new(true)
				opts.Compression = new(JsonCompressionGzip)
				opts.DateFormat = &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD")}
				opts.TimeFormat = &StageFileFormatStringOrAuto{Value: new("HH24:MI:SS")}
				opts.TimestampFormat = &StageFileFormatStringOrAuto{Value: new("YYYY-MM-DD HH24:MI:SS")}
				opts.BinaryFormat = new(BinaryFormatBase64)
				opts.TrimSpace = new(true)
				opts.MultiLine = new(true)
				opts.NullIf = &NullIfList{NullIf: []NullString{{S: "NULL"}}}
				opts.FileExtension = new(".json")
				opts.EnableOctal = new(true)
				opts.AllowDuplicate = new(true)
				opts.StripOuterArray = new(true)
				opts.StripNullValues = new(true)
				opts.IgnoreUtf8Errors = new(true)
				opts.SkipByteOrderMark = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE FILE FORMAT IF NOT EXISTS %s TYPE = JSON COMPRESSION = GZIP DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH24:MI:SS' TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS' BINARY_FORMAT = BASE64 TRIM_SPACE = true MULTI_LINE = true NULL_IF = ('NULL') FILE_EXTENSION = '.json' ENABLE_OCTAL = true ALLOW_DUPLICATE = true STRIP_OUTER_ARRAY = true STRIP_NULL_VALUES = true IGNORE_UTF8_ERRORS = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateJson_orReplace",
			func(opts *CreateJsonFileFormatOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE FILE FORMAT %s TYPE = JSON`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.CreateAvro.
		withExpectedSqlf(case_FileFormats_sql_CreateAvro_basic,
			`CREATE FILE FORMAT %s TYPE = AVRO`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_CreateAvro_all,
			func(opts *CreateAvroFileFormatOptions) {
				opts.IfNotExists = new(true)
				opts.Compression = new(AvroCompressionGzip)
				opts.TrimSpace = new(true)
				opts.ReplaceInvalidCharacters = new(true)
				opts.NullIf = &NullIfList{NullIf: []NullString{{S: "NULL"}, {S: ""}}}
				opts.Comment = new("some comment")
			},
			`CREATE FILE FORMAT IF NOT EXISTS %s TYPE = AVRO COMPRESSION = GZIP TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL', '') COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateAvro_orReplace",
			func(opts *CreateAvroFileFormatOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE FILE FORMAT %s TYPE = AVRO`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.CreateOrc.
		withExpectedSqlf(case_FileFormats_sql_CreateOrc_basic,
			`CREATE FILE FORMAT %s TYPE = ORC`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_CreateOrc_all,
			func(opts *CreateOrcFileFormatOptions) {
				opts.IfNotExists = new(true)
				opts.TrimSpace = new(true)
				opts.ReplaceInvalidCharacters = new(true)
				opts.NullIf = &NullIfList{NullIf: []NullString{{S: "NULL"}}}
				opts.Comment = new("some comment")
			},
			`CREATE FILE FORMAT IF NOT EXISTS %s TYPE = ORC TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL') COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOrc_orReplace",
			func(opts *CreateOrcFileFormatOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE FILE FORMAT %s TYPE = ORC`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.CreateParquet.
		withModify(case_FileFormats_validation_CreateParquet_opts_ConflictingFields_Compression_SnappyCompression, func(opts *CreateParquetFileFormatOptions) {
			opts.Compression = new(ParquetCompressionSnappy)
			opts.SnappyCompression = new(true)
		}).
		withExpectedSqlf(case_FileFormats_sql_CreateParquet_basic,
			`CREATE FILE FORMAT %s TYPE = PARQUET`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_CreateParquet_all,
			func(opts *CreateParquetFileFormatOptions) {
				opts.IfNotExists = new(true)
				opts.Compression = new(ParquetCompressionSnappy)
				opts.BinaryAsText = new(true)
				opts.UseLogicalType = new(true)
				opts.TrimSpace = new(true)
				opts.UseVectorizedScanner = new(true)
				opts.ReplaceInvalidCharacters = new(true)
				opts.NullIf = &NullIfList{NullIf: []NullString{{S: "NULL"}}}
				opts.Comment = new("some comment")
			},
			`CREATE FILE FORMAT IF NOT EXISTS %s TYPE = PARQUET COMPRESSION = SNAPPY BINARY_AS_TEXT = true USE_LOGICAL_TYPE = true TRIM_SPACE = true USE_VECTORIZED_SCANNER = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL') COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateParquet_orReplace",
			func(opts *CreateParquetFileFormatOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE FILE FORMAT %s TYPE = PARQUET`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.CreateXml.
		withExpectedSqlf(case_FileFormats_sql_CreateXml_basic,
			`CREATE FILE FORMAT %s TYPE = XML`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_CreateXml_all,
			func(opts *CreateXmlFileFormatOptions) {
				opts.IfNotExists = new(true)
				opts.Compression = new(XmlCompressionGzip)
				opts.IgnoreUtf8Errors = new(true)
				opts.PreserveSpace = new(true)
				opts.StripOuterElement = new(true)
				opts.DisableSnowflakeData = new(true)
				opts.DisableAutoConvert = new(true)
				opts.SkipByteOrderMark = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE FILE FORMAT IF NOT EXISTS %s TYPE = XML COMPRESSION = GZIP IGNORE_UTF8_ERRORS = true PRESERVE_SPACE = true STRIP_OUTER_ELEMENT = true DISABLE_SNOWFLAKE_DATA = true DISABLE_AUTO_CONVERT = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateXml_orReplace",
			func(opts *CreateXmlFileFormatOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE FILE FORMAT %s TYPE = XML`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsAlterCsvRenameTarget := randomSchemaObjectIdentifierInSchema(fileFormatsTestIdSchemaObjectIdentifier.SchemaId())
	fileFormatsTests.AlterCsv.
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterCsv_RenameTo,
			func(opts *AlterCsvFileFormatOptions) { opts.RenameTo = &fileFormatsAlterCsvRenameTarget },
			`ALTER FILE FORMAT %s RENAME TO %s`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(), fileFormatsAlterCsvRenameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterCsv_Set,
			func(opts *AlterCsvFileFormatOptions) {
				opts.IfExists = new(true)
				opts.Set = &AlterCsvFileFormatSet{
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
					Comment:                    new("some comment"),
				}
			},
			`ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = GZIP RECORD_DELIMITER = '\\n' FIELD_DELIMITER = ',' MULTI_LINE = true FILE_EXTENSION = '.csv' SKIP_HEADER = 2 SKIP_BLANK_LINES = true DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH24:MI:SS' TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS' BINARY_FORMAT = HEX ESCAPE = '\\' ESCAPE_UNENCLOSED_FIELD = '\\' TRIM_SPACE = true FIELD_OPTIONALLY_ENCLOSED_BY = '\"' NULL_IF = ('NULL', '') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = true EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true ENCODING = UTF8 COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterCsv_SetNullIfEmptyList",
			func(opts *AlterCsvFileFormatOptions) {
				opts.Set = &AlterCsvFileFormatSet{NullIf: &NullIfList{}}
			},
			`ALTER FILE FORMAT %s SET NULL_IF = ()`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsAlterJsonRenameTarget := randomSchemaObjectIdentifierInSchema(fileFormatsTestIdSchemaObjectIdentifier.SchemaId())
	fileFormatsTests.AlterJson.
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterJson_RenameTo,
			func(opts *AlterJsonFileFormatOptions) { opts.RenameTo = &fileFormatsAlterJsonRenameTarget },
			`ALTER FILE FORMAT %s RENAME TO %s`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(), fileFormatsAlterJsonRenameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterJson_Set,
			func(opts *AlterJsonFileFormatOptions) {
				opts.IfExists = new(true)
				opts.Set = &AlterJsonFileFormatSet{
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
					Comment:           new("some comment"),
				}
			},
			`ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = GZIP DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH24:MI:SS' TIMESTAMP_FORMAT = 'YYYY-MM-DD HH24:MI:SS' BINARY_FORMAT = BASE64 TRIM_SPACE = true MULTI_LINE = true NULL_IF = ('NULL') FILE_EXTENSION = '.json' ENABLE_OCTAL = true ALLOW_DUPLICATE = true STRIP_OUTER_ARRAY = true STRIP_NULL_VALUES = true IGNORE_UTF8_ERRORS = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterJson_SetNullIfEmptyList",
			func(opts *AlterJsonFileFormatOptions) {
				opts.Set = &AlterJsonFileFormatSet{NullIf: &NullIfList{}}
			},
			`ALTER FILE FORMAT %s SET NULL_IF = ()`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsAlterAvroRenameTarget := randomSchemaObjectIdentifierInSchema(fileFormatsTestIdSchemaObjectIdentifier.SchemaId())
	fileFormatsTests.AlterAvro.
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterAvro_RenameTo,
			func(opts *AlterAvroFileFormatOptions) { opts.RenameTo = &fileFormatsAlterAvroRenameTarget },
			`ALTER FILE FORMAT %s RENAME TO %s`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(), fileFormatsAlterAvroRenameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterAvro_Set,
			func(opts *AlterAvroFileFormatOptions) {
				opts.IfExists = new(true)
				opts.Set = &AlterAvroFileFormatSet{
					Compression:              new(AvroCompressionGzip),
					TrimSpace:                new(true),
					ReplaceInvalidCharacters: new(true),
					NullIf:                   &NullIfList{NullIf: []NullString{{S: "NULL"}, {S: ""}}},
					Comment:                  new("some comment"),
				}
			},
			`ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = GZIP TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL', '') COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterAvro_SetNullIfEmptyList",
			func(opts *AlterAvroFileFormatOptions) {
				opts.Set = &AlterAvroFileFormatSet{NullIf: &NullIfList{}}
			},
			`ALTER FILE FORMAT %s SET NULL_IF = ()`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsAlterOrcRenameTarget := randomSchemaObjectIdentifierInSchema(fileFormatsTestIdSchemaObjectIdentifier.SchemaId())
	fileFormatsTests.AlterOrc.
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterOrc_RenameTo,
			func(opts *AlterOrcFileFormatOptions) { opts.RenameTo = &fileFormatsAlterOrcRenameTarget },
			`ALTER FILE FORMAT %s RENAME TO %s`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(), fileFormatsAlterOrcRenameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterOrc_Set,
			func(opts *AlterOrcFileFormatOptions) {
				opts.IfExists = new(true)
				opts.Set = &AlterOrcFileFormatSet{
					TrimSpace:                new(true),
					ReplaceInvalidCharacters: new(true),
					NullIf:                   &NullIfList{NullIf: []NullString{{S: "NULL"}}},
					Comment:                  new("some comment"),
				}
			},
			`ALTER FILE FORMAT IF EXISTS %s SET TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL') COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterOrc_SetNullIfEmptyList",
			func(opts *AlterOrcFileFormatOptions) {
				opts.Set = &AlterOrcFileFormatSet{NullIf: &NullIfList{}}
			},
			`ALTER FILE FORMAT %s SET NULL_IF = ()`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsAlterParquetRenameTarget := randomSchemaObjectIdentifierInSchema(fileFormatsTestIdSchemaObjectIdentifier.SchemaId())
	fileFormatsTests.AlterParquet.
		withModify(case_FileFormats_validation_AlterParquet_opts_Set_ConflictingFields, func(opts *AlterParquetFileFormatOptions) {
			opts.Set = &AlterParquetFileFormatSet{
				Compression:       new(ParquetCompressionSnappy),
				SnappyCompression: new(true),
			}
		}).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterParquet_RenameTo,
			func(opts *AlterParquetFileFormatOptions) { opts.RenameTo = &fileFormatsAlterParquetRenameTarget },
			`ALTER FILE FORMAT %s RENAME TO %s`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(), fileFormatsAlterParquetRenameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterParquet_Set,
			func(opts *AlterParquetFileFormatOptions) {
				opts.IfExists = new(true)
				opts.Set = &AlterParquetFileFormatSet{
					Compression:              new(ParquetCompressionSnappy),
					BinaryAsText:             new(true),
					UseLogicalType:           new(true),
					TrimSpace:                new(true),
					UseVectorizedScanner:     new(true),
					ReplaceInvalidCharacters: new(true),
					NullIf:                   &NullIfList{NullIf: []NullString{{S: "NULL"}}},
					Comment:                  new("some comment"),
				}
			},
			`ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = SNAPPY BINARY_AS_TEXT = true USE_LOGICAL_TYPE = true TRIM_SPACE = true USE_VECTORIZED_SCANNER = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('NULL') COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterParquet_SetNullIfEmptyList",
			func(opts *AlterParquetFileFormatOptions) {
				opts.Set = &AlterParquetFileFormatSet{NullIf: &NullIfList{}}
			},
			`ALTER FILE FORMAT %s SET NULL_IF = ()`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsAlterXmlRenameTarget := randomSchemaObjectIdentifierInSchema(fileFormatsTestIdSchemaObjectIdentifier.SchemaId())
	fileFormatsTests.AlterXml.
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterXml_RenameTo,
			func(opts *AlterXmlFileFormatOptions) { opts.RenameTo = &fileFormatsAlterXmlRenameTarget },
			`ALTER FILE FORMAT %s RENAME TO %s`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(), fileFormatsAlterXmlRenameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_AlterXml_Set,
			func(opts *AlterXmlFileFormatOptions) {
				opts.IfExists = new(true)
				opts.Set = &AlterXmlFileFormatSet{
					Compression:          new(XmlCompressionGzip),
					IgnoreUtf8Errors:     new(true),
					PreserveSpace:        new(true),
					StripOuterElement:    new(true),
					DisableSnowflakeData: new(true),
					DisableAutoConvert:   new(true),
					SkipByteOrderMark:    new(true),
					Comment:              new("some comment"),
				}
			},
			`ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = GZIP IGNORE_UTF8_ERRORS = true PRESERVE_SPACE = true STRIP_OUTER_ELEMENT = true DISABLE_SNOWFLAKE_DATA = true DISABLE_AUTO_CONVERT = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'some comment'`,
			fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.Drop.
		withExpectedSqlf(case_FileFormats_sql_Drop_basic,
			`DROP FILE FORMAT %s`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_Drop_all,
			func(opts *DropFileFormatOptions) { opts.IfExists = new(true) },
			`DROP FILE FORMAT IF EXISTS %s`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	fileFormatsTests.Show.
		withExpectedSql(case_FileFormats_sql_Show_basic, `SHOW FILE FORMATS`).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_Show_all,
			func(opts *ShowFileFormatOptions) {
				opts.Like = &Like{Pattern: new("some_pattern")}
				opts.In = &In{Schema: NewDatabaseObjectIdentifier("db", "schema")}
			},
			`SHOW FILE FORMATS LIKE 'some_pattern' IN SCHEMA "db"."schema"`,
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_Show_Like,
			func(opts *ShowFileFormatOptions) { opts.Like = &Like{Pattern: new("some_pattern")} },
			`SHOW FILE FORMATS LIKE 'some_pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_FileFormats_sql_Show_In,
			func(opts *ShowFileFormatOptions) { opts.In = &In{Schema: NewDatabaseObjectIdentifier("db", "schema")} },
			`SHOW FILE FORMATS IN SCHEMA "db"."schema"`,
		)

	fileFormatsTests.Describe.
		withExpectedSqlf(case_FileFormats_sql_Describe_basic,
			`DESCRIBE FILE FORMAT %s`, fileFormatsTestIdSchemaObjectIdentifier.FullyQualifiedName())
}

func TestParseFileFormatCsv(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "COMPRESSION", Value: "GZIP"},
		{Name: "RECORD_DELIMITER", Value: "\\n"},
		{Name: "FIELD_DELIMITER", Value: ","},
		{Name: "FILE_EXTENSION", Value: "csv"},
		{Name: "SKIP_HEADER", Value: "1"},
		{Name: "PARSE_HEADER", Value: "false"},
		{Name: "SKIP_BLANK_LINES", Value: "true"},
		{Name: "DATE_FORMAT", Value: "YYYY-MM-DD"},
		{Name: "TIME_FORMAT", Value: "HH24:MI:SS"},
		{Name: "TIMESTAMP_FORMAT", Value: "YYYY-MM-DD HH24:MI:SS"},
		{Name: "BINARY_FORMAT", Value: "HEX"},
		{Name: "ESCAPE", Value: "\\"},
		{Name: "ESCAPE_UNENCLOSED_FIELD", Value: "\\"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "FIELD_OPTIONALLY_ENCLOSED_BY", Value: "'"},
		{Name: "NULL_IF", Value: "[NULL, ]"},
		{Name: "ERROR_ON_COLUMN_COUNT_MISMATCH", Value: "true"},
		{Name: "VALIDATE_UTF8", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "EMPTY_FIELD_AS_NULL", Value: "true"},
		{Name: "SKIP_BYTE_ORDER_MARK", Value: "true"},
		{Name: "ENCODING", Value: "UTF8"},
		{Name: "MULTI_LINE", Value: "true"},
	}

	csv, err := parseFileFormatCsv(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, csv.Id)
	require.Equal(t, FileFormatTypeCsv, csv.Type)
	require.Equal(t, CsvCompressionGzip, csv.Compression)
	require.Equal(t, "\\n", csv.RecordDelimiter)
	require.Equal(t, ",", csv.FieldDelimiter)
	require.Equal(t, "csv", csv.FileExtension)
	require.Equal(t, 1, csv.SkipHeader)
	require.False(t, csv.ParseHeader)
	require.True(t, csv.SkipBlankLines)
	require.Equal(t, "YYYY-MM-DD", csv.DateFormat)
	require.Equal(t, "HH24:MI:SS", csv.TimeFormat)
	require.Equal(t, "YYYY-MM-DD HH24:MI:SS", csv.TimestampFormat)
	require.Equal(t, BinaryFormatHex, csv.BinaryFormat)
	require.Equal(t, "\\", csv.Escape)
	require.Equal(t, "\\", csv.EscapeUnenclosedField)
	require.True(t, csv.TrimSpace)
	require.Equal(t, "'", csv.FieldOptionallyEnclosedBy)
	require.Equal(t, []string{"NULL", ""}, csv.NullIf)
	require.True(t, csv.ErrorOnColumnCountMismatch)
	require.True(t, csv.ValidateUtf8)
	require.True(t, csv.ReplaceInvalidCharacters)
	require.True(t, csv.EmptyFieldAsNull)
	require.True(t, csv.SkipByteOrderMark)
	require.Equal(t, CsvEncodingUtf8, csv.Encoding)
	require.True(t, csv.MultiLine)
}

func TestParseFileFormatCsv_invalidSkipHeader(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "SKIP_HEADER", Value: "not-a-number"},
	}

	_, err := parseFileFormatCsv(properties, id)
	require.ErrorContains(t, err, `cannot cast SKIP_HEADER value "not-a-number" to int`)
}

func TestParseFileFormatCsv_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"PARSE_HEADER"},
		{"TRIM_SPACE"},
		{"ERROR_ON_COLUMN_COUNT_MISMATCH"},
		{"VALIDATE_UTF8"},
		{"SKIP_BLANK_LINES"},
		{"REPLACE_INVALID_CHARACTERS"},
		{"EMPTY_FIELD_AS_NULL"},
		{"SKIP_BYTE_ORDER_MARK"},
		{"MULTI_LINE"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "CSV"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatCsv(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatCsv_multipleInvalidValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "SKIP_HEADER", Value: "not-a-number"},
		{Name: "TRIM_SPACE", Value: "not-a-bool"},
	}

	_, err := parseFileFormatCsv(properties, id)
	require.ErrorContains(t, err, `cannot cast SKIP_HEADER value "not-a-number" to int`)
	require.ErrorContains(t, err, `cannot cast TRIM_SPACE value "not-a-bool" to bool`)
}

func TestParseFileFormatCsv_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	csv, err := parseFileFormatCsv(properties, id)
	require.NoError(t, err)
	require.Equal(t, FileFormatTypeCsv, csv.Type)
}

func TestParseFileFormatJson(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "JSON"},
		{Name: "COMPRESSION", Value: "AUTO"},
		{Name: "DATE_FORMAT", Value: "YYYY-MM-DD"},
		{Name: "TIME_FORMAT", Value: "HH24:MI:SS"},
		{Name: "TIMESTAMP_FORMAT", Value: "YYYY-MM-DD HH24:MI:SS"},
		{Name: "BINARY_FORMAT", Value: "HEX"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "MULTI_LINE", Value: "true"},
		{Name: "STRIP_OUTER_ARRAY", Value: "true"},
		{Name: "NULL_IF", Value: "[]"},
		{Name: "FILE_EXTENSION", Value: "json"},
		{Name: "ENABLE_OCTAL", Value: "true"},
		{Name: "ALLOW_DUPLICATE", Value: "true"},
		{Name: "STRIP_NULL_VALUES", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "IGNORE_UTF8_ERRORS", Value: "true"},
		{Name: "SKIP_BYTE_ORDER_MARK", Value: "true"},
	}

	json, err := parseFileFormatJson(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, json.Id)
	require.Equal(t, FileFormatTypeJson, json.Type)
	require.Equal(t, JsonCompressionAuto, json.Compression)
	require.Equal(t, "YYYY-MM-DD", json.DateFormat)
	require.Equal(t, "HH24:MI:SS", json.TimeFormat)
	require.Equal(t, "YYYY-MM-DD HH24:MI:SS", json.TimestampFormat)
	require.Equal(t, BinaryFormatHex, json.BinaryFormat)
	require.True(t, json.TrimSpace)
	require.True(t, json.MultiLine)
	require.True(t, json.StripOuterArray)
	require.Equal(t, []string{}, json.NullIf)
	require.Equal(t, "json", json.FileExtension)
	require.True(t, json.EnableOctal)
	require.True(t, json.AllowDuplicate)
	require.True(t, json.StripNullValues)
	require.True(t, json.ReplaceInvalidCharacters)
	require.True(t, json.IgnoreUtf8Errors)
	require.True(t, json.SkipByteOrderMark)
}

func TestParseFileFormatJson_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"TRIM_SPACE"},
		{"MULTI_LINE"},
		{"ENABLE_OCTAL"},
		{"ALLOW_DUPLICATE"},
		{"STRIP_OUTER_ARRAY"},
		{"STRIP_NULL_VALUES"},
		{"REPLACE_INVALID_CHARACTERS"},
		{"IGNORE_UTF8_ERRORS"},
		{"SKIP_BYTE_ORDER_MARK"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "JSON"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatJson(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatJson_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "JSON"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	json, err := parseFileFormatJson(properties, id)
	require.NoError(t, err)
	require.Equal(t, FileFormatTypeJson, json.Type)
}

func TestParseFileFormatAvro(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "AVRO"},
		{Name: "COMPRESSION", Value: "GZIP"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "NULL_IF", Value: "[NULL, ]"},
	}

	avro, err := parseFileFormatAvro(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, avro.Id)
	require.Equal(t, FileFormatTypeAvro, avro.Type)
	require.Equal(t, AvroCompressionGzip, avro.Compression)
	require.True(t, avro.TrimSpace)
	require.True(t, avro.ReplaceInvalidCharacters)
	require.Equal(t, []string{"NULL", ""}, avro.NullIf)
}

func TestParseFileFormatAvro_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"TRIM_SPACE"},
		{"REPLACE_INVALID_CHARACTERS"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "AVRO"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatAvro(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatAvro_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "AVRO"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	avro, err := parseFileFormatAvro(properties, id)
	require.NoError(t, err)
	require.Equal(t, FileFormatTypeAvro, avro.Type)
}

func TestParseFileFormatOrc(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "ORC"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "false"},
		{Name: "NULL_IF", Value: "[NULL]"},
	}

	orc, err := parseFileFormatOrc(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, orc.Id)
	require.Equal(t, FileFormatTypeOrc, orc.Type)
	require.True(t, orc.TrimSpace)
	require.False(t, orc.ReplaceInvalidCharacters)
	require.Equal(t, []string{"NULL"}, orc.NullIf)
}

func TestParseFileFormatOrc_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"TRIM_SPACE"},
		{"REPLACE_INVALID_CHARACTERS"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "ORC"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatOrc(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatOrc_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "ORC"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	orc, err := parseFileFormatOrc(properties, id)
	require.NoError(t, err)
	require.Equal(t, FileFormatTypeOrc, orc.Type)
}

func TestParseFileFormatParquet(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "PARQUET"},
		{Name: "COMPRESSION", Value: "SNAPPY"},
		{Name: "BINARY_AS_TEXT", Value: "true"},
		{Name: "USE_LOGICAL_TYPE", Value: "true"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "USE_VECTORIZED_SCANNER", Value: "false"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "NULL_IF", Value: "[NULL]"},
	}

	parquet, err := parseFileFormatParquet(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, parquet.Id)
	require.Equal(t, FileFormatTypeParquet, parquet.Type)
	require.Equal(t, ParquetCompressionSnappy, parquet.Compression)
	require.True(t, parquet.BinaryAsText)
	require.True(t, parquet.UseLogicalType)
	require.True(t, parquet.TrimSpace)
	require.False(t, parquet.UseVectorizedScanner)
	require.True(t, parquet.ReplaceInvalidCharacters)
	require.Equal(t, []string{"NULL"}, parquet.NullIf)
}

func TestParseFileFormatParquet_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"BINARY_AS_TEXT"},
		{"USE_LOGICAL_TYPE"},
		{"TRIM_SPACE"},
		{"USE_VECTORIZED_SCANNER"},
		{"REPLACE_INVALID_CHARACTERS"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "PARQUET"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatParquet(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatParquet_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "PARQUET"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	parquet, err := parseFileFormatParquet(properties, id)
	require.NoError(t, err)
	require.Equal(t, FileFormatTypeParquet, parquet.Type)
}

func TestParseFileFormatXml(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "XML"},
		{Name: "COMPRESSION", Value: "GZIP"},
		{Name: "IGNORE_UTF8_ERRORS", Value: "true"},
		{Name: "PRESERVE_SPACE", Value: "true"},
		{Name: "STRIP_OUTER_ELEMENT", Value: "false"},
		{Name: "DISABLE_SNOWFLAKE_DATA", Value: "true"},
		{Name: "DISABLE_AUTO_CONVERT", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "SKIP_BYTE_ORDER_MARK", Value: "true"},
	}

	xml, err := parseFileFormatXml(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, xml.Id)
	require.Equal(t, FileFormatTypeXml, xml.Type)
	require.Equal(t, XmlCompressionGzip, xml.Compression)
	require.True(t, xml.IgnoreUtf8Errors)
	require.True(t, xml.PreserveSpace)
	require.False(t, xml.StripOuterElement)
	require.True(t, xml.DisableSnowflakeData)
	require.True(t, xml.DisableAutoConvert)
	require.True(t, xml.ReplaceInvalidCharacters)
	require.True(t, xml.SkipByteOrderMark)
}

func TestParseFileFormatXml_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"IGNORE_UTF8_ERRORS"},
		{"PRESERVE_SPACE"},
		{"STRIP_OUTER_ELEMENT"},
		{"DISABLE_SNOWFLAKE_DATA"},
		{"DISABLE_AUTO_CONVERT"},
		{"REPLACE_INVALID_CHARACTERS"},
		{"SKIP_BYTE_ORDER_MARK"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "XML"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatXml(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatXml_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "XML"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	xml, err := parseFileFormatXml(properties, id)
	require.NoError(t, err)
	require.Equal(t, FileFormatTypeXml, xml.Type)
}

func TestParseFileFormatAllDetails(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("csv", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "CSV"},
			{Name: "COMPRESSION", Value: "GZIP"},
			{Name: "TRIM_SPACE", Value: "true"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, id, details.Id)
		require.Equal(t, FileFormatTypeCsv, details.Type)
		require.NotNil(t, details.Csv)
		require.Equal(t, CsvCompressionGzip, details.Csv.Compression)
		require.True(t, details.Csv.TrimSpace)
		require.Nil(t, details.Json)
		require.Nil(t, details.Avro)
		require.Nil(t, details.Orc)
		require.Nil(t, details.Parquet)
		require.Nil(t, details.Xml)
	})

	t.Run("json", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "JSON"},
			{Name: "COMPRESSION", Value: "AUTO"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeJson, details.Type)
		require.NotNil(t, details.Json)
		require.Equal(t, JsonCompressionAuto, details.Json.Compression)
		require.Nil(t, details.Csv)
	})

	t.Run("avro", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "AVRO"},
			{Name: "TRIM_SPACE", Value: "true"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeAvro, details.Type)
		require.NotNil(t, details.Avro)
		require.True(t, details.Avro.TrimSpace)
	})

	t.Run("orc", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "ORC"},
			{Name: "TRIM_SPACE", Value: "true"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeOrc, details.Type)
		require.NotNil(t, details.Orc)
		require.True(t, details.Orc.TrimSpace)
	})

	t.Run("parquet", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "PARQUET"},
			{Name: "COMPRESSION", Value: "SNAPPY"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeParquet, details.Type)
		require.NotNil(t, details.Parquet)
		require.Equal(t, ParquetCompressionSnappy, details.Parquet.Compression)
	})

	t.Run("xml", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "XML"},
			{Name: "COMPRESSION", Value: "GZIP"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeXml, details.Type)
		require.NotNil(t, details.Xml)
		require.Equal(t, XmlCompressionGzip, details.Xml.Compression)
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := parseFileFormatAllDetails([]FileFormatProperty{{Name: "TYPE", Value: "NOT_A_TYPE"}}, id)
		require.Error(t, err)
	})

	t.Run("missing type", func(t *testing.T) {
		_, err := parseFileFormatAllDetails([]FileFormatProperty{{Name: "COMPRESSION", Value: "GZIP"}}, id)
		require.ErrorContains(t, err, "describe did not return a recognized file format type")
	})

	t.Run("propagates csv parsing error", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "CSV"},
			{Name: "SKIP_HEADER", Value: "not-a-number"},
		}
		_, err := parseFileFormatAllDetails(properties, id)
		require.ErrorContains(t, err, `cannot cast SKIP_HEADER value "not-a-number" to int`)
	})
}

// TestParseFileFormat_lowercaseEnumValues covers legacy objects whose DESCRIBE FILE FORMAT values are
// not all-uppercase. Every enum-typed property is normalized through its To<Enum> converter, so such
// objects parse (and, in the resources, import) instead of failing a case-sensitive comparison.
func TestParseFileFormat_lowercaseEnumValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("csv", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "csv"},
			{Name: "COMPRESSION", Value: "gzip"},
			{Name: "BINARY_FORMAT", Value: "base64"},
			{Name: "ENCODING", Value: "utf8"},
		}

		csv, err := parseFileFormatCsv(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeCsv, csv.Type)
		require.Equal(t, CsvCompressionGzip, csv.Compression)
		require.Equal(t, BinaryFormatBase64, csv.BinaryFormat)
		require.Equal(t, CsvEncodingUtf8, csv.Encoding)
	})

	t.Run("json", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "json"},
			{Name: "COMPRESSION", Value: "auto"},
			{Name: "BINARY_FORMAT", Value: "hex"},
		}

		json, err := parseFileFormatJson(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeJson, json.Type)
		require.Equal(t, JsonCompressionAuto, json.Compression)
		require.Equal(t, BinaryFormatHex, json.BinaryFormat)
	})

	t.Run("avro", func(t *testing.T) {
		avro, err := parseFileFormatAvro([]FileFormatProperty{
			{Name: "TYPE", Value: "avro"},
			{Name: "COMPRESSION", Value: "gzip"},
		}, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeAvro, avro.Type)
		require.Equal(t, AvroCompressionGzip, avro.Compression)
	})

	t.Run("orc", func(t *testing.T) {
		orc, err := parseFileFormatOrc([]FileFormatProperty{{Name: "TYPE", Value: "orc"}}, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeOrc, orc.Type)
	})

	t.Run("parquet", func(t *testing.T) {
		parquet, err := parseFileFormatParquet([]FileFormatProperty{
			{Name: "TYPE", Value: "parquet"},
			{Name: "COMPRESSION", Value: "snappy"},
		}, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeParquet, parquet.Type)
		require.Equal(t, ParquetCompressionSnappy, parquet.Compression)
	})

	t.Run("xml", func(t *testing.T) {
		xml, err := parseFileFormatXml([]FileFormatProperty{
			{Name: "TYPE", Value: "xml"},
			{Name: "COMPRESSION", Value: "gzip"},
		}, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeXml, xml.Type)
		require.Equal(t, XmlCompressionGzip, xml.Compression)
	})

	t.Run("all details", func(t *testing.T) {
		details, err := parseFileFormatAllDetails([]FileFormatProperty{
			{Name: "TYPE", Value: "parquet"},
			{Name: "COMPRESSION", Value: "snappy"},
		}, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeParquet, details.Type)
		require.NotNil(t, details.Parquet)
		require.Equal(t, ParquetCompressionSnappy, details.Parquet.Compression)
	})
}

// TestParseFileFormat_invalidEnumValues asserts that an unrecognized enum value is reported as an
// error rather than being silently carried through as a raw string.
func TestParseFileFormat_invalidEnumValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("csv collects every invalid enum value", func(t *testing.T) {
		_, err := parseFileFormatCsv([]FileFormatProperty{
			{Name: "TYPE", Value: "NOT_A_TYPE"},
			{Name: "COMPRESSION", Value: "NOT_A_COMPRESSION"},
			{Name: "BINARY_FORMAT", Value: "NOT_A_BINARY_FORMAT"},
			{Name: "ENCODING", Value: "NOT_AN_ENCODING"},
		}, id)
		require.ErrorContains(t, err, "NOT_A_TYPE")
		require.ErrorContains(t, err, "NOT_A_COMPRESSION")
		require.ErrorContains(t, err, "NOT_A_BINARY_FORMAT")
		require.ErrorContains(t, err, "NOT_AN_ENCODING")
	})

	t.Run("json invalid compression", func(t *testing.T) {
		_, err := parseFileFormatJson([]FileFormatProperty{
			{Name: "TYPE", Value: "JSON"},
			{Name: "COMPRESSION", Value: "NOT_A_COMPRESSION"},
		}, id)
		require.ErrorContains(t, err, "NOT_A_COMPRESSION")
	})

	t.Run("avro invalid compression", func(t *testing.T) {
		_, err := parseFileFormatAvro([]FileFormatProperty{
			{Name: "TYPE", Value: "AVRO"},
			{Name: "COMPRESSION", Value: "NOT_A_COMPRESSION"},
		}, id)
		require.ErrorContains(t, err, "NOT_A_COMPRESSION")
	})

	t.Run("parquet invalid compression", func(t *testing.T) {
		_, err := parseFileFormatParquet([]FileFormatProperty{
			{Name: "TYPE", Value: "PARQUET"},
			{Name: "COMPRESSION", Value: "NOT_A_COMPRESSION"},
		}, id)
		require.ErrorContains(t, err, "NOT_A_COMPRESSION")
	})

	t.Run("xml invalid compression", func(t *testing.T) {
		_, err := parseFileFormatXml([]FileFormatProperty{
			{Name: "TYPE", Value: "XML"},
			{Name: "COMPRESSION", Value: "NOT_A_COMPRESSION"},
		}, id)
		require.ErrorContains(t, err, "NOT_A_COMPRESSION")
	})
}

// TestParseFileFormatCsv_hyphenatedEncoding covers DESCRIBE FILE FORMAT values that Snowflake
// stores as IANA-style aliases (utf-8, UTF-16LE) instead of the canonical enum names.
func TestParseFileFormatCsv_hyphenatedEncoding(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	for _, tc := range []struct {
		raw  string
		want CsvEncoding
	}{
		{raw: "utf-8", want: CsvEncodingUtf8},
		{raw: "UTF-16LE", want: CsvEncodingUtf16le},
		{raw: "ISO-8859-1", want: CsvEncodingIso88591},
	} {
		t.Run(tc.raw, func(t *testing.T) {
			csv, err := parseFileFormatCsv([]FileFormatProperty{
				{Name: "TYPE", Value: "CSV"},
				{Name: "ENCODING", Value: tc.raw},
			}, id)
			require.NoError(t, err)
			require.Equal(t, tc.want, csv.Encoding)
		})
	}
}
