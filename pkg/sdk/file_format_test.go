package sdk

import (
	"testing"
)

func TestFileFormatsCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("minimal", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeCSV,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = CSV`, id.FullyQualifiedName())
	})

	t.Run("CSV with empty string null_if", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeCSV,
			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVNullIf: &[]NullString{{""}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = CSV NULL_IF = ('')`, id.FullyQualifiedName())
	})

	t.Run("complete CSV", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			Type:        FileFormatTypeCSV,

			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVCompression:               &CSVCompressionBz2,
				CSVRecordDelimiter:           String("-"),
				CSVFieldDelimiter:            String(":"),
				CSVFileExtension:             String("csv"),
				CSVSkipHeader:                Int(5),
				CSVSkipBlankLines:            Bool(true),
				CSVDateFormat:                String("YYYY-MM-DD"),
				CSVTimeFormat:                String("HH:mm:SS"),
				CSVTimestampFormat:           String("time"),
				CSVBinaryFormat:              &BinaryFormatUTF8,
				CSVEscape:                    String("\\"),
				CSVEscapeUnenclosedField:     String("§"),
				CSVTrimSpace:                 Bool(true),
				CSVFieldOptionallyEnclosedBy: String("\""),
				CSVNullIf: &[]NullString{
					{"nul"},
					{"nulll"},
				},
				CSVErrorOnColumnCountMismatch: Bool(true),
				CSVReplaceInvalidCharacters:   Bool(true),
				CSVEmptyFieldAsNull:           Bool(true),
				CSVSkipByteOrderMark:          Bool(true),
				CSVEncoding:                   &CSVEncodingISO2022KR,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = CSV COMPRESSION = BZ2 RECORD_DELIMITER = '-' FIELD_DELIMITER = ':' FILE_EXTENSION = 'csv' SKIP_HEADER = 5 SKIP_BLANK_LINES = true DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH:mm:SS' TIMESTAMP_FORMAT = 'time' BINARY_FORMAT = UTF8 ESCAPE = '\\' ESCAPE_UNENCLOSED_FIELD = '§' TRIM_SPACE = true FIELD_OPTIONALLY_ENCLOSED_BY = '\"' NULL_IF = ('nul', 'nulll') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = true EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true ENCODING = 'ISO2022KR'`, id.FullyQualifiedName())
	})

	t.Run("complete JSON", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			Type:        FileFormatTypeJSON,

			FileFormatTypeOptions: FileFormatTypeOptions{
				JSONCompression:     &JSONCompressionBrotli,
				JSONDateFormat:      String("YYYY-MM-DD"),
				JSONTimeFormat:      String("HH:mm:SS"),
				JSONTimestampFormat: String("aze"),
				JSONBinaryFormat:    &BinaryFormatHex,
				JSONTrimSpace:       Bool(true),
				JSONNullIf: &[]NullString{
					{"c1"},
					{"c2"},
				},
				JSONFileExtension:            String("json"),
				JSONEnableOctal:              Bool(true),
				JSONAllowDuplicate:           Bool(true),
				JSONStripOuterArray:          Bool(true),
				JSONStripNullValues:          Bool(true),
				JSONReplaceInvalidCharacters: Bool(true),
				JSONSkipByteOrderMark:        Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = JSON COMPRESSION = BROTLI DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH:mm:SS' TIMESTAMP_FORMAT = 'aze' BINARY_FORMAT = HEX TRIM_SPACE = true NULL_IF = ('c1', 'c2') FILE_EXTENSION = 'json' ENABLE_OCTAL = true ALLOW_DUPLICATE = true STRIP_OUTER_ARRAY = true STRIP_NULL_VALUES = true REPLACE_INVALID_CHARACTERS = true SKIP_BYTE_ORDER_MARK = true`, id.FullyQualifiedName())
	})

	t.Run("JSON with empty string null_if", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeJSON,
			FileFormatTypeOptions: FileFormatTypeOptions{
				JSONNullIf: &[]NullString{{""}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = JSON NULL_IF = ('')`, id.FullyQualifiedName())
	})

	t.Run("complete Avro", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			Type:        FileFormatTypeAvro,

			FileFormatTypeOptions: FileFormatTypeOptions{
				AvroCompression:              &AvroCompressionDeflate,
				AvroTrimSpace:                Bool(true),
				AvroReplaceInvalidCharacters: Bool(true),
				AvroNullIf:                   &[]NullString{{"nul"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = AVRO COMPRESSION = DEFLATE TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nul')`, id.FullyQualifiedName())
	})

	t.Run("AVRO with empty string null_if", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeAvro,
			FileFormatTypeOptions: FileFormatTypeOptions{
				AvroNullIf: &[]NullString{{""}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = AVRO NULL_IF = ('')`, id.FullyQualifiedName())
	})

	t.Run("complete ORC", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			Type:        FileFormatTypeORC,

			FileFormatTypeOptions: FileFormatTypeOptions{
				ORCTrimSpace:                Bool(true),
				ORCReplaceInvalidCharacters: Bool(true),
				ORCNullIf:                   &[]NullString{{"nul"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = ORC TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nul')`, id.FullyQualifiedName())
	})

	t.Run("ORC with empty string null_if", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeORC,
			FileFormatTypeOptions: FileFormatTypeOptions{
				ORCNullIf: &[]NullString{{""}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = ORC NULL_IF = ('')`, id.FullyQualifiedName())
	})

	t.Run("complete Parquet", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			Type:        FileFormatTypeParquet,

			FileFormatTypeOptions: FileFormatTypeOptions{
				ParquetCompression:              &ParquetCompressionLzo,
				ParquetBinaryAsText:             Bool(true),
				ParquetTrimSpace:                Bool(true),
				ParquetReplaceInvalidCharacters: Bool(true),
				ParquetNullIf:                   &[]NullString{{"nil"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = PARQUET COMPRESSION = LZO BINARY_AS_TEXT = true TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`, id.FullyQualifiedName())
	})

	t.Run("PARQUET with empty string null_if", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeParquet,
			FileFormatTypeOptions: FileFormatTypeOptions{
				ParquetNullIf: &[]NullString{{""}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = PARQUET NULL_IF = ('')`, id.FullyQualifiedName())
	})

	t.Run("complete XML", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   Bool(true),
			Temporary:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			Type:        FileFormatTypeXML,
			FileFormatTypeOptions: FileFormatTypeOptions{
				XMLCompression:          &XMLCompressionZstd,
				XMLIgnoreUTF8Errors:     Bool(true),
				XMLPreserveSpace:        Bool(true),
				XMLStripOuterElement:    Bool(true),
				XMLDisableSnowflakeData: Bool(true),
				XMLDisableAutoConvert:   Bool(true),
				XMLSkipByteOrderMark:    Bool(true),
			},
			Comment: String("test file format"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = XML COMPRESSION = ZSTD IGNORE_UTF8_ERRORS = true PRESERVE_SPACE = true STRIP_OUTER_ELEMENT = true DISABLE_SNOWFLAKE_DATA = true DISABLE_AUTO_CONVERT = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'test file format'`, id.FullyQualifiedName())
	})

	t.Run("previous test", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeCSV,
			FileFormatTypeOptions: FileFormatTypeOptions{
				CSVNullIf:                     &[]NullString{{"NULL"}},
				CSVSkipBlankLines:             Bool(false),
				CSVTrimSpace:                  Bool(false),
				CSVErrorOnColumnCountMismatch: Bool(true),
				CSVReplaceInvalidCharacters:   Bool(false),
				CSVEmptyFieldAsNull:           Bool(false),
				CSVSkipByteOrderMark:          Bool(false),
			},
			Comment: String("great comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = CSV SKIP_BLANK_LINES = false TRIM_SPACE = false NULL_IF = ('NULL') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = false EMPTY_FIELD_AS_NULL = false SKIP_BYTE_ORDER_MARK = false COMMENT = 'great comment'`, id.FullyQualifiedName())
	})
}

func TestFileFormatsAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("rename", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()
		opts := &AlterFileFormatOptions{
			IfExists: Bool(true),
			name:     id,
			Rename: &AlterFileFormatRenameOptions{
				NewName: newId,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FILE FORMAT IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := &AlterFileFormatOptions{
			IfExists: Bool(true),
			name:     id,
			Set: &FileFormatTypeOptions{
				AvroCompression:              &AvroCompressionBrotli,
				AvroTrimSpace:                Bool(true),
				AvroReplaceInvalidCharacters: Bool(true),
				AvroNullIf:                   &[]NullString{{"nil"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = BROTLI TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`, id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := &AlterFileFormatOptions{
			IfExists: Bool(true),
			name:     id,
			Set: &FileFormatTypeOptions{
				AvroCompression:              &AvroCompressionBrotli,
				AvroTrimSpace:                Bool(true),
				AvroReplaceInvalidCharacters: Bool(true),
				AvroNullIf:                   &[]NullString{{"nil"}},
				Comment:                      String("some comment"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FILE FORMAT IF EXISTS %s SET COMMENT = 'some comment' COMPRESSION = BROTLI TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`, id.FullyQualifiedName())
	})
}

func TestFileFormatsDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("only name", func(t *testing.T) {
		opts := &DropFileFormatOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP FILE FORMAT %s`, id.FullyQualifiedName())
	})

	t.Run("with IfExists", func(t *testing.T) {
		opts := &DropFileFormatOptions{
			name:     id,
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP FILE FORMAT IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestFileFormatsShow(t *testing.T) {
	t.Run("without show options", func(t *testing.T) {
		opts := &ShowFileFormatsOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW FILE FORMATS`)
	})

	t.Run("with show options", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &ShowFileFormatsOptions{
			Like: &Like{
				Pattern: String("test"),
			},
			In: &In{
				Schema: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW FILE FORMATS LIKE 'test' IN SCHEMA %s`, id.FullyQualifiedName())
	})
}

func TestFileFormatsShowById(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &ShowFileFormatsOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: id.SchemaId(),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW FILE FORMATS LIKE '%s' IN SCHEMA %s`, id.Name(), id.SchemaId().FullyQualifiedName())
	})
}

func TestFileFormatsDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	opts := &describeFileFormatOptions{
		name: id,
	}
	assertOptsValidAndSQLEquals(t, opts, `DESCRIBE FILE FORMAT %s`, id.FullyQualifiedName())
}
