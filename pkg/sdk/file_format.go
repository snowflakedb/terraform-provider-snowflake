package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateFileFormatOptions)
	_ validatable = new(AlterFileFormatOptions)
	_ validatable = new(DropFileFormatOptions)
	_ validatable = new(ShowFileFormatsOptions)
	_ validatable = new(describeFileFormatOptions)
)

type FileFormats interface {
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateFileFormatOptions) error
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterFileFormatOptions) error
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropFileFormatOptions) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, opts *ShowFileFormatsOptions) ([]FileFormat, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*FileFormat, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*FileFormat, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetails, error)
}

var _ FileFormats = (*fileFormats)(nil)

type fileFormats struct {
	client *Client
}

type FileFormat struct {
	Name          SchemaObjectIdentifier
	CreatedOn     time.Time
	Type          FileFormatType
	Owner         string
	Comment       string
	OwnerRoleType string
	Options       FileFormatTypeOptions
}

func (v *FileFormat) ID() SchemaObjectIdentifier {
	return v.Name
}

func (v *FileFormat) ObjectType() ObjectType {
	return ObjectTypeFileFormat
}

type FileFormatRow struct {
	FormatOptions string    `db:"format_options"`
	CreatedOn     time.Time `db:"created_on"`
	Name          string    `db:"name"`
	DatabaseName  string    `db:"database_name"`
	SchemaName    string    `db:"schema_name"`
	FormatType    string    `db:"type"`
	Owner         string    `db:"owner"`
	Comment       string    `db:"comment"`
	OwnerRoleType string    `db:"owner_role_type"`
}

type showFileFormatsOptionsResult struct {
	// CSV + shared fields
	Type                       string   `json:"TYPE"`
	RecordDelimiter            string   `json:"RECORD_DELIMITER"`
	FieldDelimiter             string   `json:"FIELD_DELIMITER"`
	FileExtension              string   `json:"FILE_EXTENSION"`
	SkipHeader                 int      `json:"SKIP_HEADER"`
	ParseHeader                bool     `json:"PARSE_HEADER"`
	DateFormat                 string   `json:"DATE_FORMAT"`
	TimeFormat                 string   `json:"TIME_FORMAT"`
	TimestampFormat            string   `json:"TIMESTAMP_FORMAT"`
	BinaryFormat               string   `json:"BINARY_FORMAT"`
	Escape                     string   `json:"ESCAPE"`
	EscapeUnenclosedField      string   `json:"ESCAPE_UNENCLOSED_FIELD"`
	TrimSpace                  bool     `json:"TRIM_SPACE"`
	FieldOptionallyEnclosedBy  string   `json:"FIELD_OPTIONALLY_ENCLOSED_BY"`
	NullIf                     []string `json:"NULL_IF"`
	Compression                string   `json:"COMPRESSION"`
	ErrorOnColumnCountMismatch bool     `json:"ERROR_ON_COLUMN_COUNT_MISMATCH"`
	ValidateUTF8               bool     `json:"VALIDATE_UTF8"`
	SkipBlankLines             bool     `json:"SKIP_BLANK_LINES"`
	ReplaceInvalidCharacters   bool     `json:"REPLACE_INVALID_CHARACTERS"`
	EmptyFieldAsNull           bool     `json:"EMPTY_FIELD_AS_NULL"`
	SkipByteOrderMark          bool     `json:"SKIP_BYTE_ORDER_MARK"`
	Encoding                   string   `json:"ENCODING"`

	// JSON fields
	EnableOctal      bool `json:"ENABLE_OCTAL"`
	AllowDuplicate   bool `json:"ALLOW_DUPLICATE"`
	StripOuterArray  bool `json:"STRIP_OUTER_ARRAY"`
	StripNullValues  bool `json:"STRIP_NULL_VALUES"`
	IgnoreUTF8Errors bool `json:"IGNORE_UTF8_ERRORS"`

	// Parquet fields
	BinaryAsText bool `json:"BINARY_AS_TEXT"`

	// XML fields
	PreserveSpace        bool `json:"PRESERVE_SPACE"`
	StripOuterElement    bool `json:"STRIP_OUTER_ELEMENT"`
	DisableSnowflakeData bool `json:"DISABLE_SNOWFLAKE_DATA"`
	DisableAutoConvert   bool `json:"DISABLE_AUTO_CONVERT"`
}

func (row FileFormatRow) convert() *FileFormat {
	inputOptions := showFileFormatsOptionsResult{}
	err := json.Unmarshal([]byte(row.FormatOptions), &inputOptions)
	if err != nil {
		fmt.Printf("%s", err)
		panic("cannot parse options json")
	}

	ff := &FileFormat{
		Name:          NewSchemaObjectIdentifier(row.DatabaseName, row.SchemaName, row.Name),
		CreatedOn:     row.CreatedOn,
		Type:          FileFormatType(row.FormatType),
		Owner:         row.Owner,
		Comment:       row.Comment,
		OwnerRoleType: row.OwnerRoleType,
		Options:       FileFormatTypeOptions{},
	}

	newNullIf := make([]NullString, len(inputOptions.NullIf))
	for i, s := range inputOptions.NullIf {
		newNullIf[i] = NullString{s}
	}

	switch ff.Type {
	case FileFormatTypeCSV:
		ff.Options.CSVCompression = (*CSVCompression)(&inputOptions.Compression)
		ff.Options.CSVRecordDelimiter = &inputOptions.RecordDelimiter
		ff.Options.CSVFieldDelimiter = &inputOptions.FieldDelimiter
		ff.Options.CSVFileExtension = &inputOptions.FileExtension
		ff.Options.CSVParseHeader = &inputOptions.ParseHeader
		ff.Options.CSVSkipHeader = &inputOptions.SkipHeader
		ff.Options.CSVSkipBlankLines = &inputOptions.SkipBlankLines
		ff.Options.CSVDateFormat = &inputOptions.DateFormat
		ff.Options.CSVTimeFormat = &inputOptions.TimeFormat
		ff.Options.CSVTimestampFormat = &inputOptions.TimestampFormat
		ff.Options.CSVBinaryFormat = (*BinaryFormat)(&inputOptions.BinaryFormat)
		ff.Options.CSVEscape = &inputOptions.Escape
		ff.Options.CSVEscapeUnenclosedField = &inputOptions.EscapeUnenclosedField
		ff.Options.CSVTrimSpace = &inputOptions.TrimSpace
		ff.Options.CSVFieldOptionallyEnclosedBy = &inputOptions.FieldOptionallyEnclosedBy
		ff.Options.CSVNullIf = &newNullIf
		ff.Options.CSVErrorOnColumnCountMismatch = &inputOptions.ErrorOnColumnCountMismatch
		ff.Options.CSVReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.CSVEmptyFieldAsNull = &inputOptions.EmptyFieldAsNull
		ff.Options.CSVSkipByteOrderMark = &inputOptions.SkipByteOrderMark
		ff.Options.CSVEncoding = (*CSVEncoding)(&inputOptions.Encoding)
	case FileFormatTypeJSON:
		ff.Options.JSONCompression = (*JSONCompression)(&inputOptions.Compression)
		ff.Options.JSONDateFormat = &inputOptions.DateFormat
		ff.Options.JSONTimeFormat = &inputOptions.TimeFormat
		ff.Options.JSONTimestampFormat = &inputOptions.TimestampFormat
		ff.Options.JSONBinaryFormat = (*BinaryFormat)(&inputOptions.BinaryFormat)
		ff.Options.JSONTrimSpace = &inputOptions.TrimSpace
		ff.Options.JSONNullIf = newNullIf
		ff.Options.JSONFileExtension = &inputOptions.FileExtension
		ff.Options.JSONEnableOctal = &inputOptions.EnableOctal
		ff.Options.JSONAllowDuplicate = &inputOptions.AllowDuplicate
		ff.Options.JSONStripOuterArray = &inputOptions.StripOuterArray
		ff.Options.JSONStripNullValues = &inputOptions.StripNullValues
		ff.Options.JSONReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.JSONIgnoreUTF8Errors = &inputOptions.IgnoreUTF8Errors
		ff.Options.JSONSkipByteOrderMark = &inputOptions.SkipByteOrderMark
	case FileFormatTypeAvro:
		ff.Options.AvroTrimSpace = &inputOptions.TrimSpace
		ff.Options.AvroNullIf = &newNullIf
		ff.Options.AvroCompression = (*AvroCompression)(&inputOptions.Compression)
		ff.Options.AvroReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
	case FileFormatTypeORC:
		ff.Options.ORCTrimSpace = &inputOptions.TrimSpace
		ff.Options.ORCReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.ORCNullIf = &newNullIf
	case FileFormatTypeParquet:
		ff.Options.ParquetTrimSpace = &inputOptions.TrimSpace
		ff.Options.ParquetNullIf = &newNullIf
		ff.Options.ParquetCompression = (*ParquetCompression)(&inputOptions.Compression)
		ff.Options.ParquetBinaryAsText = &inputOptions.BinaryAsText
		ff.Options.ParquetReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
	case FileFormatTypeXML:
		ff.Options.XMLCompression = (*XMLCompression)(&inputOptions.Compression)
		ff.Options.XMLIgnoreUTF8Errors = &inputOptions.IgnoreUTF8Errors
		ff.Options.XMLPreserveSpace = &inputOptions.PreserveSpace
		ff.Options.XMLStripOuterElement = &inputOptions.StripOuterElement
		ff.Options.XMLDisableSnowflakeData = &inputOptions.DisableSnowflakeData
		ff.Options.XMLDisableAutoConvert = &inputOptions.DisableAutoConvert
		ff.Options.XMLReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.XMLSkipByteOrderMark = &inputOptions.SkipByteOrderMark
	}

	return ff
}

type FileFormatType string

var (
	FileFormatTypeCSV     FileFormatType = "CSV"
	FileFormatTypeJSON    FileFormatType = "JSON"
	FileFormatTypeAvro    FileFormatType = "AVRO"
	FileFormatTypeORC     FileFormatType = "ORC"
	FileFormatTypeParquet FileFormatType = "PARQUET"
	FileFormatTypeXML     FileFormatType = "XML"
)

type BinaryFormat string

var (
	BinaryFormatHex    BinaryFormat = "HEX"
	BinaryFormatBase64 BinaryFormat = "BASE64"
	BinaryFormatUTF8   BinaryFormat = "UTF8"
)

type CSVCompression string

var (
	CSVCompressionAuto       CSVCompression = "AUTO"
	CSVCompressionGzip       CSVCompression = "GZIP"
	CSVCompressionBz2        CSVCompression = "BZ2"
	CSVCompressionBrotli     CSVCompression = "BROTLI"
	CSVCompressionZstd       CSVCompression = "ZSTD"
	CSVCompressionDeflate    CSVCompression = "DEFLATE"
	CSVCompressionRawDeflate CSVCompression = "RAW_DEFLATE"
	CSVCompressionNone       CSVCompression = "NONE"
)

type CSVEncoding string

var (
	CSVEncodingBIG5        CSVEncoding = "BIG5"
	CSVEncodingEUCJP       CSVEncoding = "EUCJP"
	CSVEncodingEUCKR       CSVEncoding = "EUCKR"
	CSVEncodingGB18030     CSVEncoding = "GB18030"
	CSVEncodingIBM420      CSVEncoding = "IBM420"
	CSVEncodingIBM424      CSVEncoding = "IBM424"
	CSVEncodingISO2022CN   CSVEncoding = "ISO2022CN"
	CSVEncodingISO2022JP   CSVEncoding = "ISO2022JP"
	CSVEncodingISO2022KR   CSVEncoding = "ISO2022KR"
	CSVEncodingISO88591    CSVEncoding = "ISO88591"
	CSVEncodingISO88592    CSVEncoding = "ISO88592"
	CSVEncodingISO88595    CSVEncoding = "ISO88595"
	CSVEncodingISO88596    CSVEncoding = "ISO88596"
	CSVEncodingISO88597    CSVEncoding = "ISO88597"
	CSVEncodingISO88598    CSVEncoding = "ISO88598"
	CSVEncodingISO88599    CSVEncoding = "ISO88599"
	CSVEncodingISO885915   CSVEncoding = "ISO885915"
	CSVEncodingKOI8R       CSVEncoding = "KOI8R"
	CSVEncodingSHIFTJIS    CSVEncoding = "SHIFTJIS"
	CSVEncodingUTF8        CSVEncoding = "UTF8"
	CSVEncodingUTF16       CSVEncoding = "UTF16"
	CSVEncodingUTF16BE     CSVEncoding = "UTF16BE"
	CSVEncodingUTF16LE     CSVEncoding = "UTF16LE"
	CSVEncodingUTF32       CSVEncoding = "UTF32"
	CSVEncodingUTF32BE     CSVEncoding = "UTF32BE"
	CSVEncodingUTF32LE     CSVEncoding = "UTF32LE"
	CSVEncodingWINDOWS1250 CSVEncoding = "WINDOWS1250"
	CSVEncodingWINDOWS1251 CSVEncoding = "WINDOWS1251"
	CSVEncodingWINDOWS1252 CSVEncoding = "WINDOWS1252"
	CSVEncodingWINDOWS1253 CSVEncoding = "WINDOWS1253"
	CSVEncodingWINDOWS1254 CSVEncoding = "WINDOWS1254"
	CSVEncodingWINDOWS1255 CSVEncoding = "WINDOWS1255"
	CSVEncodingWINDOWS1256 CSVEncoding = "WINDOWS1256"
)

type JSONCompression string

var (
	JSONCompressionAuto       JSONCompression = "AUTO"
	JSONCompressionGzip       JSONCompression = "GZIP"
	JSONCompressionBz2        JSONCompression = "BZ2"
	JSONCompressionBrotli     JSONCompression = "BROTLI"
	JSONCompressionZstd       JSONCompression = "ZSTD"
	JSONCompressionDeflate    JSONCompression = "DEFLATE"
	JSONCompressionRawDeflate JSONCompression = "RAW_DEFLATE"
	JSONCompressionNone       JSONCompression = "NONE"
)

type AvroCompression string

var (
	AvroCompressionAuto       AvroCompression = "AUTO"
	AvroCompressionGzip       AvroCompression = "GZIP"
	AvroCompressionBrotli     AvroCompression = "BROTLI"
	AvroCompressionZstd       AvroCompression = "ZSTD"
	AvroCompressionDeflate    AvroCompression = "DEFLATE"
	AvroCompressionRawDeflate AvroCompression = "RAW_DEFLATE"
	AvroCompressionNone       AvroCompression = "NONE"
)

type ParquetCompression string

var (
	ParquetCompressionAuto   ParquetCompression = "AUTO"
	ParquetCompressionLzo    ParquetCompression = "LZO"
	ParquetCompressionSnappy ParquetCompression = "SNAPPY"
	ParquetCompressionNone   ParquetCompression = "NONE"
)

type XMLCompression string

var (
	XMLCompressionAuto       XMLCompression = "AUTO"
	XMLCompressionGzip       XMLCompression = "GZIP"
	XMLCompressionBz2        XMLCompression = "BZ2"
	XMLCompressionBrotli     XMLCompression = "BROTLI"
	XMLCompressionZstd       XMLCompression = "ZSTD"
	XMLCompressionDeflate    XMLCompression = "DEFLATE"
	XMLCompressionRawDeflate XMLCompression = "RAW_DEFLATE"
	XMLCompressionNone       XMLCompression = "NONE"
)

type NullString struct {
	S string `ddl:"parameter,no_equals,single_quotes"`
}

// CreateFileFormatOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-file-format.
type CreateFileFormatOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary   *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	fileFormat  bool                   `ddl:"static" sql:"FILE FORMAT"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	Type        FileFormatType         `ddl:"parameter" sql:"TYPE"`
	FileFormatTypeOptions
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateFileFormatOptions) validate() error {
	fields := opts.FileFormatTypeOptions.fieldsByType()

	for formatType := range fields {
		if opts.Type == formatType {
			continue
		}
		if anyValueSet(fields[formatType]...) {
			return fmt.Errorf("cannot set %s fields when TYPE = %s", formatType, opts.Type)
		}
	}

	err := opts.FileFormatTypeOptions.validate()
	if err != nil {
		return err
	}

	return nil
}

func (v *fileFormats) Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateFileFormatOptions) error {
	if opts == nil {
		opts = &CreateFileFormatOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}

	_, err = v.client.exec(ctx, sql)
	return err
}

// AlterFileFormatOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-file-format.
type AlterFileFormatOptions struct {
	alter      bool                   `ddl:"static" sql:"ALTER"`
	fileFormat bool                   `ddl:"static" sql:"FILE FORMAT"`
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`

	Rename *AlterFileFormatRenameOptions
	Set    *FileFormatTypeOptions `ddl:"list,no_comma" sql:"SET"`
}

func (opts *AlterFileFormatOptions) validate() error {
	if !exactlyOneValueSet(opts.Rename, opts.Set) {
		return fmt.Errorf("only one of Rename or Set can be set at once.")
	}
	if valueSet(opts.Set) {
		err := opts.Set.validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type AlterFileFormatRenameOptions struct {
	NewName SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
}

type FileFormatTypeOptions struct {
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	// CSV type options
	CSVCompression                *CSVCompression `ddl:"parameter" sql:"COMPRESSION"`
	CSVRecordDelimiter            *string         `ddl:"parameter,single_quotes" sql:"RECORD_DELIMITER"`
	CSVFieldDelimiter             *string         `ddl:"parameter,single_quotes" sql:"FIELD_DELIMITER"`
	CSVFileExtension              *string         `ddl:"parameter,single_quotes" sql:"FILE_EXTENSION"`
	CSVParseHeader                *bool           `ddl:"parameter" sql:"PARSE_HEADER"`
	CSVSkipHeader                 *int            `ddl:"parameter" sql:"SKIP_HEADER"`
	CSVSkipBlankLines             *bool           `ddl:"parameter" sql:"SKIP_BLANK_LINES"`
	CSVDateFormat                 *string         `ddl:"parameter,single_quotes" sql:"DATE_FORMAT"`
	CSVTimeFormat                 *string         `ddl:"parameter,single_quotes" sql:"TIME_FORMAT"`
	CSVTimestampFormat            *string         `ddl:"parameter,single_quotes" sql:"TIMESTAMP_FORMAT"`
	CSVBinaryFormat               *BinaryFormat   `ddl:"parameter" sql:"BINARY_FORMAT"`
	CSVEscape                     *string         `ddl:"parameter,single_quotes" sql:"ESCAPE"`
	CSVEscapeUnenclosedField      *string         `ddl:"parameter,single_quotes" sql:"ESCAPE_UNENCLOSED_FIELD"`
	CSVTrimSpace                  *bool           `ddl:"parameter" sql:"TRIM_SPACE"`
	CSVFieldOptionallyEnclosedBy  *string         `ddl:"parameter,single_quotes" sql:"FIELD_OPTIONALLY_ENCLOSED_BY"`
	CSVNullIf                     *[]NullString   `ddl:"parameter,parentheses" sql:"NULL_IF"`
	CSVErrorOnColumnCountMismatch *bool           `ddl:"parameter" sql:"ERROR_ON_COLUMN_COUNT_MISMATCH"`
	CSVReplaceInvalidCharacters   *bool           `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	CSVEmptyFieldAsNull           *bool           `ddl:"parameter" sql:"EMPTY_FIELD_AS_NULL"`
	CSVSkipByteOrderMark          *bool           `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`
	CSVEncoding                   *CSVEncoding    `ddl:"parameter,single_quotes" sql:"ENCODING"`

	// JSON type options
	JSONCompression              *JSONCompression `ddl:"parameter" sql:"COMPRESSION"`
	JSONDateFormat               *string          `ddl:"parameter,single_quotes" sql:"DATE_FORMAT"`
	JSONTimeFormat               *string          `ddl:"parameter,single_quotes" sql:"TIME_FORMAT"`
	JSONTimestampFormat          *string          `ddl:"parameter,single_quotes" sql:"TIMESTAMP_FORMAT"`
	JSONBinaryFormat             *BinaryFormat    `ddl:"parameter" sql:"BINARY_FORMAT"`
	JSONTrimSpace                *bool            `ddl:"parameter" sql:"TRIM_SPACE"`
	JSONNullIf                   []NullString     `ddl:"parameter,parentheses" sql:"NULL_IF"`
	JSONFileExtension            *string          `ddl:"parameter,single_quotes" sql:"FILE_EXTENSION"`
	JSONEnableOctal              *bool            `ddl:"parameter" sql:"ENABLE_OCTAL"`
	JSONAllowDuplicate           *bool            `ddl:"parameter" sql:"ALLOW_DUPLICATE"`
	JSONStripOuterArray          *bool            `ddl:"parameter" sql:"STRIP_OUTER_ARRAY"`
	JSONStripNullValues          *bool            `ddl:"parameter" sql:"STRIP_NULL_VALUES"`
	JSONReplaceInvalidCharacters *bool            `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	JSONIgnoreUTF8Errors         *bool            `ddl:"parameter" sql:"IGNORE_UTF8_ERRORS"`
	JSONSkipByteOrderMark        *bool            `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`

	// AVRO type options
	AvroCompression              *AvroCompression `ddl:"parameter" sql:"COMPRESSION"`
	AvroTrimSpace                *bool            `ddl:"parameter" sql:"TRIM_SPACE"`
	AvroReplaceInvalidCharacters *bool            `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	AvroNullIf                   *[]NullString    `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// ORC type options
	ORCTrimSpace                *bool         `ddl:"parameter" sql:"TRIM_SPACE"`
	ORCReplaceInvalidCharacters *bool         `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	ORCNullIf                   *[]NullString `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// PARQUET type options
	ParquetCompression              *ParquetCompression `ddl:"parameter" sql:"COMPRESSION"`
	ParquetSnappyCompression        *bool               `ddl:"parameter" sql:"SNAPPY_COMPRESSION"`
	ParquetBinaryAsText             *bool               `ddl:"parameter" sql:"BINARY_AS_TEXT"`
	ParquetTrimSpace                *bool               `ddl:"parameter" sql:"TRIM_SPACE"`
	ParquetReplaceInvalidCharacters *bool               `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	ParquetNullIf                   *[]NullString       `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// XML type options
	XMLCompression              *XMLCompression `ddl:"parameter" sql:"COMPRESSION"`
	XMLIgnoreUTF8Errors         *bool           `ddl:"parameter" sql:"IGNORE_UTF8_ERRORS"`
	XMLPreserveSpace            *bool           `ddl:"parameter" sql:"PRESERVE_SPACE"`
	XMLStripOuterElement        *bool           `ddl:"parameter" sql:"STRIP_OUTER_ELEMENT"`
	XMLDisableSnowflakeData     *bool           `ddl:"parameter" sql:"DISABLE_SNOWFLAKE_DATA"`
	XMLDisableAutoConvert       *bool           `ddl:"parameter" sql:"DISABLE_AUTO_CONVERT"`
	XMLReplaceInvalidCharacters *bool           `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	XMLSkipByteOrderMark        *bool           `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`
}

func (opts *FileFormatTypeOptions) fieldsByType() map[FileFormatType][]any {
	return map[FileFormatType][]any{
		FileFormatTypeCSV: {
			opts.CSVCompression,
			opts.CSVRecordDelimiter,
			opts.CSVFieldDelimiter,
			opts.CSVFileExtension,
			opts.CSVParseHeader,
			opts.CSVSkipHeader,
			opts.CSVSkipBlankLines,
			opts.CSVDateFormat,
			opts.CSVTimeFormat,
			opts.CSVTimestampFormat,
			opts.CSVBinaryFormat,
			opts.CSVEscape,
			opts.CSVEscapeUnenclosedField,
			opts.CSVTrimSpace,
			opts.CSVFieldOptionallyEnclosedBy,
			opts.CSVNullIf,
			opts.CSVErrorOnColumnCountMismatch,
			opts.CSVReplaceInvalidCharacters,
			opts.CSVEmptyFieldAsNull,
			opts.CSVSkipByteOrderMark,
			opts.CSVEncoding,
		},
		FileFormatTypeJSON: {
			opts.JSONCompression,
			opts.JSONDateFormat,
			opts.JSONTimeFormat,
			opts.JSONTimestampFormat,
			opts.JSONBinaryFormat,
			opts.JSONTrimSpace,
			opts.JSONNullIf,
			opts.JSONFileExtension,
			opts.JSONEnableOctal,
			opts.JSONAllowDuplicate,
			opts.JSONStripOuterArray,
			opts.JSONStripNullValues,
			opts.JSONReplaceInvalidCharacters,
			opts.JSONIgnoreUTF8Errors,
			opts.JSONSkipByteOrderMark,
		},
		FileFormatTypeAvro: {
			opts.AvroCompression,
			opts.AvroTrimSpace,
			opts.AvroReplaceInvalidCharacters,
			opts.AvroNullIf,
		},
		FileFormatTypeORC: {
			opts.ORCTrimSpace,
			opts.ORCReplaceInvalidCharacters,
			opts.ORCNullIf,
		},
		FileFormatTypeParquet: {
			opts.ParquetCompression,
			opts.ParquetSnappyCompression,
			opts.ParquetBinaryAsText,
			opts.ParquetTrimSpace,
			opts.ParquetReplaceInvalidCharacters,
			opts.ParquetNullIf,
		},
		FileFormatTypeXML: {
			opts.XMLCompression,
			opts.XMLIgnoreUTF8Errors,
			opts.XMLPreserveSpace,
			opts.XMLStripOuterElement,
			opts.XMLDisableSnowflakeData,
			opts.XMLDisableAutoConvert,
			opts.XMLReplaceInvalidCharacters,
			opts.XMLSkipByteOrderMark,
		},
	}
}

func (opts *FileFormatTypeOptions) validate() error {
	fields := opts.fieldsByType()
	count := 0

	for formatType := range fields {
		if anyValueSet(fields[formatType]...) {
			count += 1
			if count > 1 {
				return fmt.Errorf("Cannot set options for different format types")
			}
		}
	}

	if everyValueSet(opts.CSVParseHeader, opts.CSVSkipHeader) && *opts.CSVParseHeader {
		return fmt.Errorf("ParseHeader and SkipHeader cannot be set simultaneously")
	}

	if everyValueSet(opts.JSONIgnoreUTF8Errors, opts.JSONReplaceInvalidCharacters) && *opts.JSONIgnoreUTF8Errors && *opts.JSONReplaceInvalidCharacters {
		return fmt.Errorf("IgnoreUTF8Errors and ReplaceInvalidCharacters cannot be set simultaneously")
	}

	if everyValueSet(opts.ParquetCompression, opts.ParquetSnappyCompression) && *opts.ParquetSnappyCompression {
		return fmt.Errorf("Compression and SnappyCompression cannot be set simultaneously")
	}

	if everyValueSet(opts.XMLIgnoreUTF8Errors, opts.XMLReplaceInvalidCharacters) && *opts.XMLIgnoreUTF8Errors && *opts.XMLReplaceInvalidCharacters {
		return fmt.Errorf("IgnoreUTF8Errors and ReplaceInvalidCharacters cannot be set simultaneously")
	}

	validEnclosedBy := []string{"NONE", "'", `"`}
	if valueSet(opts.CSVFieldOptionallyEnclosedBy) && !slices.Contains(validEnclosedBy, *opts.CSVFieldOptionallyEnclosedBy) {
		return fmt.Errorf("CSVFieldOptionallyEnclosedBy must be one of %v", validEnclosedBy)
	}
	return nil
}

func (v *fileFormats) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterFileFormatOptions) error {
	if opts == nil {
		opts = &AlterFileFormatOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// DropFileFormatOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-file-format.
type DropFileFormatOptions struct {
	drop       bool                   `ddl:"static" sql:"DROP"`
	fileFormat string                 `ddl:"static" sql:"FILE FORMAT"`
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropFileFormatOptions) validate() error {
	return nil
}

func (v *fileFormats) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropFileFormatOptions) error {
	if opts == nil {
		opts = &DropFileFormatOptions{}
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

func (v *fileFormats) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropFileFormatOptions{IfExists: Bool(true)}) }, ctx, id)
}

// ShowFileFormatsOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-file-formats.
type ShowFileFormatsOptions struct {
	show        bool  `ddl:"static" sql:"SHOW"`
	fileFormats bool  `ddl:"static" sql:"FILE FORMATS"`
	Like        *Like `ddl:"keyword" sql:"LIKE"`
	In          *In   `ddl:"keyword" sql:"IN"`
}

func (opts *ShowFileFormatsOptions) validate() error {
	return nil
}

func (v *fileFormats) Show(ctx context.Context, opts *ShowFileFormatsOptions) ([]FileFormat, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[FileFormatRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[FileFormatRow, FileFormat](dbRows)
	return resultList, nil
}

func (v *fileFormats) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*FileFormat, error) {
	fileFormats, err := v.client.FileFormats.Show(ctx, &ShowFileFormatsOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: id.SchemaId(),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(fileFormats, func(format FileFormat) bool {
		return format.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
}

func (v *fileFormats) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*FileFormat, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

type FileFormatDetails struct {
	Type    FileFormatType
	Options FileFormatTypeOptions
}

type FileFormatDetailsRow struct {
	Property         string
	Property_Type    string
	Property_Value   string
	Property_Default string
}

// describeFileFormatOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-file-format.
type describeFileFormatOptions struct {
	describe   bool                   `ddl:"static" sql:"DESCRIBE"`
	fileFormat string                 `ddl:"static" sql:"FILE FORMAT"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *describeFileFormatOptions) validate() error {
	return nil
}

func (v *fileFormats) Describe(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetails, error) {
	opts := &describeFileFormatOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []FileFormatDetailsRow
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	details := FileFormatDetails{}
	for _, row := range rows {
		if row.Property == "TYPE" {
			details.Type = FileFormatType(row.Property_Value)
			break
		}
	}

	switch details.Type {
	case FileFormatTypeCSV:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "RECORD_DELIMITER":
				details.Options.CSVRecordDelimiter = &v
			case "FIELD_DELIMITER":
				details.Options.CSVFieldDelimiter = &v
			case "FILE_EXTENSION":
				details.Options.CSVFileExtension = &v
			case "SKIP_HEADER":
				i, err := strconv.ParseInt(v, 10, 0)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to int: %w`, v, err)
				}
				i0 := int(i)
				details.Options.CSVSkipHeader = &i0
			case "PARSE_HEADER":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast PARSE_HEADER value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVParseHeader = &b
			case "DATE_FORMAT":
				details.Options.CSVDateFormat = &v
			case "TIME_FORMAT":
				details.Options.CSVTimeFormat = &v
			case "TIMESTAMP_FORMAT":
				details.Options.CSVTimestampFormat = &v
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.Options.CSVBinaryFormat = &bf
			case "ESCAPE":
				details.Options.CSVEscape = &v
			case "ESCAPE_UNENCLOSED_FIELD":
				details.Options.CSVEscapeUnenclosedField = &v
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVTrimSpace = &b
			case "FIELD_OPTIONALLY_ENCLOSED_BY":
				details.Options.CSVFieldOptionallyEnclosedBy = &v
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range ParseCommaSeparatedStringArray(v, false) {
					if s != "" {
						newNullIf = append(newNullIf, NullString{s})
					}
				}
				details.Options.CSVNullIf = &newNullIf
			case "COMPRESSION":
				comp := CSVCompression(v)
				details.Options.CSVCompression = &comp
			case "ERROR_ON_COLUMN_COUNT_MISMATCH":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ERROR_ON_COLUMN_COUNT_MISMATCH value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVErrorOnColumnCountMismatch = &b
			// case "VALIDATE_UTF8":
			// 	details.Options.C = &v
			case "SKIP_BLANK_LINES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BLANK_LINES value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVSkipBlankLines = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVReplaceInvalidCharacters = &b
			case "EMPTY_FIELD_AS_NULL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast EMPTY_FIELD_AS_NULL value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVEmptyFieldAsNull = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.CSVSkipByteOrderMark = &b
			case "ENCODING":
				enc := CSVEncoding(v)
				details.Options.CSVEncoding = &enc
			}
		}
	case FileFormatTypeJSON:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "FILE_EXTENSION":
				details.Options.JSONFileExtension = &v
			case "DATE_FORMAT":
				details.Options.JSONDateFormat = &v
			case "TIME_FORMAT":
				details.Options.JSONTimeFormat = &v
			case "TIMESTAMP_FORMAT":
				details.Options.JSONTimestampFormat = &v
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.Options.JSONBinaryFormat = &bf
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range ParseCommaSeparatedStringArray(v, false) {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.JSONNullIf = newNullIf
			case "COMPRESSION":
				comp := JSONCompression(v)
				details.Options.JSONCompression = &comp
			case "ENABLE_OCTAL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ENABLE_OCTAL value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONEnableOctal = &b
			case "ALLOW_DUPLICATE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ALLOW_DUPLICATE value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONAllowDuplicate = &b
			case "STRIP_OUTER_ARRAY":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_OUTER_ARRAY value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONStripOuterArray = &b
			case "STRIP_NULL_VALUES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_NULL_VALUES value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONStripNullValues = &b
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONIgnoreUTF8Errors = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONReplaceInvalidCharacters = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.JSONSkipByteOrderMark = &b
			}
		}
	case FileFormatTypeAvro:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.AvroTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range ParseCommaSeparatedStringArray(v, false) {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.AvroNullIf = &newNullIf
			case "COMPRESSION":
				comp := AvroCompression(v)
				details.Options.AvroCompression = &comp
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.AvroReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeORC:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.ORCTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range ParseCommaSeparatedStringArray(v, false) {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.ORCNullIf = &newNullIf
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.ORCReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeParquet:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetTrimSpace = &b
			case "NULL_IF":
				newNullIf := []NullString{}
				for _, s := range ParseCommaSeparatedStringArray(v, false) {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.ParquetNullIf = &newNullIf
			case "COMPRESSION":
				comp := ParquetCompression(v)
				details.Options.ParquetCompression = &comp
			case "BINARY_AS_TEXT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast BINARY_AS_TEXT value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetBinaryAsText = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeXML:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "COMPRESSION":
				comp := XMLCompression(v)
				details.Options.XMLCompression = &comp
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLIgnoreUTF8Errors = &b
			case "PRESERVE_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast PRESERVE_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLPreserveSpace = &b
			case "STRIP_OUTER_ELEMENT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_OUTER_ELEMENT value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLStripOuterElement = &b
			case "DISABLE_SNOWFLAKE_DATA":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast DISABLE_SNOWFLAKE_DATA value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLDisableSnowflakeData = &b
			case "DISABLE_AUTO_CONVERT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast DISABLE_AUTO_CONVERT value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLDisableAutoConvert = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLReplaceInvalidCharacters = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.XMLSkipByteOrderMark = &b
			}
		}
	default:
		return nil, fmt.Errorf("Describe did not return format type")
	}

	return &details, nil
}
