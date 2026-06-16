package sdk

import (
	"testing"
)

func TestFileFormatsCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("minimal", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeCsv,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = CSV`, id.FullyQualifiedName())
	})

	t.Run("complete CSV", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   new(true),
			Temporary:   new(true),
			name:        id,
			IfNotExists: new(true),
			Type:        FileFormatTypeCsv,

			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				CSVCompression:               new(CsvCompressionBz2),
				CSVRecordDelimiter:           new("-"),
				CSVFieldDelimiter:            new(":"),
				CSVFileExtension:             new("csv"),
				CSVSkipHeader:                new(5),
				CSVSkipBlankLines:            new(true),
				CSVDateFormat:                new("YYYY-MM-DD"),
				CSVTimeFormat:                new("HH:mm:SS"),
				CSVTimestampFormat:           new("time"),
				CSVBinaryFormat:              new(BinaryFormatUtf8),
				CSVEscape:                    new("\\"),
				CSVEscapeUnenclosedField:     new("§"),
				CSVTrimSpace:                 new(true),
				CSVFieldOptionallyEnclosedBy: new("\""),
				CSVNullIf: &[]NullString{
					{"nul"},
					{"nulll"},
				},
				CSVErrorOnColumnCountMismatch: new(true),
				CSVReplaceInvalidCharacters:   new(true),
				CSVEmptyFieldAsNull:           new(true),
				CSVSkipByteOrderMark:          new(true),
				CSVEncoding:                   new(CsvEncodingIso2022kr),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = CSV COMPRESSION = BZ2 RECORD_DELIMITER = '-' FIELD_DELIMITER = ':' FILE_EXTENSION = 'csv' SKIP_HEADER = 5 SKIP_BLANK_LINES = true DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH:mm:SS' TIMESTAMP_FORMAT = 'time' BINARY_FORMAT = UTF8 ESCAPE = '\\' ESCAPE_UNENCLOSED_FIELD = '§' TRIM_SPACE = true FIELD_OPTIONALLY_ENCLOSED_BY = '\"' NULL_IF = ('nul', 'nulll') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = true EMPTY_FIELD_AS_NULL = true SKIP_BYTE_ORDER_MARK = true ENCODING = 'ISO2022KR'`, id.FullyQualifiedName())
	})

	t.Run("complete JSON", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   new(true),
			Temporary:   new(true),
			name:        id,
			IfNotExists: new(true),
			Type:        FileFormatTypeJson,

			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				JSONCompression:     new(JsonCompressionBrotli),
				JSONDateFormat:      new("YYYY-MM-DD"),
				JSONTimeFormat:      new("HH:mm:SS"),
				JSONTimestampFormat: new("aze"),
				JSONBinaryFormat:    new(BinaryFormatHex),
				JSONTrimSpace:       new(true),
				JSONNullIf: []NullString{
					{"c1"},
					{"c2"},
				},
				JSONFileExtension:            new("json"),
				JSONEnableOctal:              new(true),
				JSONAllowDuplicate:           new(true),
				JSONStripOuterArray:          new(true),
				JSONStripNullValues:          new(true),
				JSONReplaceInvalidCharacters: new(true),
				JSONSkipByteOrderMark:        new(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = JSON COMPRESSION = BROTLI DATE_FORMAT = 'YYYY-MM-DD' TIME_FORMAT = 'HH:mm:SS' TIMESTAMP_FORMAT = 'aze' BINARY_FORMAT = HEX TRIM_SPACE = true NULL_IF = ('c1', 'c2') FILE_EXTENSION = 'json' ENABLE_OCTAL = true ALLOW_DUPLICATE = true STRIP_OUTER_ARRAY = true STRIP_NULL_VALUES = true REPLACE_INVALID_CHARACTERS = true SKIP_BYTE_ORDER_MARK = true`, id.FullyQualifiedName())
	})

	t.Run("complete Avro", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   new(true),
			Temporary:   new(true),
			name:        id,
			IfNotExists: new(true),
			Type:        FileFormatTypeAvro,

			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				AvroCompression:              new(AvroCompressionDeflate),
				AvroTrimSpace:                new(true),
				AvroReplaceInvalidCharacters: new(true),
				AvroNullIf:                   &[]NullString{{"nul"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = AVRO COMPRESSION = DEFLATE TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nul')`, id.FullyQualifiedName())
	})

	t.Run("complete ORC", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   new(true),
			Temporary:   new(true),
			name:        id,
			IfNotExists: new(true),
			Type:        FileFormatTypeOrc,

			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				ORCTrimSpace:                new(true),
				ORCReplaceInvalidCharacters: new(true),
				ORCNullIf:                   &[]NullString{{"nul"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = ORC TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nul')`, id.FullyQualifiedName())
	})

	t.Run("complete Parquet", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   new(true),
			Temporary:   new(true),
			name:        id,
			IfNotExists: new(true),
			Type:        FileFormatTypeParquet,

			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				ParquetCompression:              new(ParquetCompressionLzo),
				ParquetBinaryAsText:             new(true),
				ParquetTrimSpace:                new(true),
				ParquetReplaceInvalidCharacters: new(true),
				ParquetNullIf:                   &[]NullString{{"nil"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = PARQUET COMPRESSION = LZO BINARY_AS_TEXT = true TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`, id.FullyQualifiedName())
	})

	t.Run("complete XML", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			OrReplace:   new(true),
			Temporary:   new(true),
			name:        id,
			IfNotExists: new(true),
			Type:        FileFormatTypeXml,
			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				XMLCompression:          new(XmlCompressionZstd),
				XMLIgnoreUTF8Errors:     new(true),
				XMLPreserveSpace:        new(true),
				XMLStripOuterElement:    new(true),
				XMLDisableSnowflakeData: new(true),
				XMLDisableAutoConvert:   new(true),
				XMLSkipByteOrderMark:    new(true),
			},
			Comment: new("test file format"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TEMPORARY FILE FORMAT IF NOT EXISTS %s TYPE = XML COMPRESSION = ZSTD IGNORE_UTF8_ERRORS = true PRESERVE_SPACE = true STRIP_OUTER_ELEMENT = true DISABLE_SNOWFLAKE_DATA = true DISABLE_AUTO_CONVERT = true SKIP_BYTE_ORDER_MARK = true COMMENT = 'test file format'`, id.FullyQualifiedName())
	})

	t.Run("previous test", func(t *testing.T) {
		opts := &CreateFileFormatOptions{
			name: id,
			Type: FileFormatTypeCsv,
			LegacyFileFormatTypeOptions: LegacyFileFormatTypeOptions{
				CSVNullIf:                     &[]NullString{{"NULL"}},
				CSVSkipBlankLines:             new(false),
				CSVTrimSpace:                  new(false),
				CSVErrorOnColumnCountMismatch: new(true),
				CSVReplaceInvalidCharacters:   new(false),
				CSVEmptyFieldAsNull:           new(false),
				CSVSkipByteOrderMark:          new(false),
			},
			Comment: new("great comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE FILE FORMAT %s TYPE = CSV SKIP_BLANK_LINES = false TRIM_SPACE = false NULL_IF = ('NULL') ERROR_ON_COLUMN_COUNT_MISMATCH = true REPLACE_INVALID_CHARACTERS = false EMPTY_FIELD_AS_NULL = false SKIP_BYTE_ORDER_MARK = false COMMENT = 'great comment'`, id.FullyQualifiedName())
	})
}

func TestFileFormatsAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("rename", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()
		opts := &AlterFileFormatOptions{
			IfExists: new(true),
			name:     id,
			Rename: &AlterFileFormatRenameOptions{
				NewName: newId,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FILE FORMAT IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := &AlterFileFormatOptions{
			IfExists: new(true),
			name:     id,
			Set: &LegacyFileFormatTypeOptions{
				AvroCompression:              new(AvroCompressionBrotli),
				AvroTrimSpace:                new(true),
				AvroReplaceInvalidCharacters: new(true),
				AvroNullIf:                   &[]NullString{{"nil"}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER FILE FORMAT IF EXISTS %s SET COMPRESSION = BROTLI TRIM_SPACE = true REPLACE_INVALID_CHARACTERS = true NULL_IF = ('nil')`, id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		opts := &AlterFileFormatOptions{
			IfExists: new(true),
			name:     id,
			Set: &LegacyFileFormatTypeOptions{
				AvroCompression:              new(AvroCompressionBrotli),
				AvroTrimSpace:                new(true),
				AvroReplaceInvalidCharacters: new(true),
				AvroNullIf:                   &[]NullString{{"nil"}},
				Comment:                      new("some comment"),
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
			IfExists: new(true),
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
				Pattern: new("test"),
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
				Pattern: new(id.Name()),
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
