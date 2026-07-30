package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ validatable = new(CreateFileFormatOptionsLegacy)
	_ validatable = new(AlterFileFormatOptionsLegacy)
	_ validatable = new(DropFileFormatOptionsLegacy)
	_ validatable = new(ShowFileFormatsOptionsLegacy)
	_ validatable = new(describeFileFormatOptionsLegacy)

	_ convertibleRow[FileFormatLegacy] = new(FileFormatRowLegacy)
)

type FileFormatsLegacy interface {
	Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateFileFormatOptionsLegacy) error
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterFileFormatOptionsLegacy) error
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropFileFormatOptionsLegacy) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, opts *ShowFileFormatsOptionsLegacy) ([]FileFormatLegacy, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatLegacy, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatLegacy, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetailsLegacy, error)
}

var _ FileFormatsLegacy = (*fileFormatsLegacy)(nil)

type fileFormatsLegacy struct {
	client *Client
}

type FileFormatLegacy struct {
	Name          SchemaObjectIdentifier
	CreatedOn     time.Time
	Type          FileFormatType
	Owner         string
	Comment       string
	OwnerRoleType string
	Options       FileFormatTypeOptionsLegacy
}

func (v *FileFormatLegacy) ID() SchemaObjectIdentifier {
	return v.Name
}

func (v *FileFormatLegacy) ObjectType() ObjectType {
	return ObjectTypeFileFormat
}

type FileFormatRowLegacy struct {
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

type showFileFormatsOptionsResultLegacy struct {
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

func (row FileFormatRowLegacy) convert() (*FileFormatLegacy, error) {
	inputOptions := showFileFormatsOptionsResultLegacy{}
	err := json.Unmarshal([]byte(row.FormatOptions), &inputOptions)
	if err != nil {
		return nil, fmt.Errorf("cannot parse format options: %w", err)
	}

	ff := &FileFormatLegacy{
		Name:          NewSchemaObjectIdentifier(row.DatabaseName, row.SchemaName, row.Name),
		CreatedOn:     row.CreatedOn,
		Type:          FileFormatType(row.FormatType),
		Owner:         row.Owner,
		Comment:       row.Comment,
		OwnerRoleType: row.OwnerRoleType,
		Options:       FileFormatTypeOptionsLegacy{},
	}

	newNullIf := make([]NullString, len(inputOptions.NullIf))
	for i, s := range inputOptions.NullIf {
		newNullIf[i] = NullString{s}
	}

	switch ff.Type {
	case FileFormatTypeCsv:
		ff.Options.CSVCompression = (*CsvCompression)(&inputOptions.Compression)
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
		ff.Options.CSVEncoding = (*CsvEncoding)(&inputOptions.Encoding)
	case FileFormatTypeJson:
		ff.Options.JSONCompression = (*JsonCompression)(&inputOptions.Compression)
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
	case FileFormatTypeOrc:
		ff.Options.ORCTrimSpace = &inputOptions.TrimSpace
		ff.Options.ORCReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.ORCNullIf = &newNullIf
	case FileFormatTypeParquet:
		ff.Options.ParquetTrimSpace = &inputOptions.TrimSpace
		ff.Options.ParquetNullIf = &newNullIf
		ff.Options.ParquetCompression = (*ParquetCompression)(&inputOptions.Compression)
		ff.Options.ParquetBinaryAsText = &inputOptions.BinaryAsText
		ff.Options.ParquetReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
	case FileFormatTypeXml:
		ff.Options.XMLCompression = (*XmlCompression)(&inputOptions.Compression)
		ff.Options.XMLIgnoreUTF8Errors = &inputOptions.IgnoreUTF8Errors
		ff.Options.XMLPreserveSpace = &inputOptions.PreserveSpace
		ff.Options.XMLStripOuterElement = &inputOptions.StripOuterElement
		ff.Options.XMLDisableSnowflakeData = &inputOptions.DisableSnowflakeData
		ff.Options.XMLDisableAutoConvert = &inputOptions.DisableAutoConvert
		ff.Options.XMLReplaceInvalidCharacters = &inputOptions.ReplaceInvalidCharacters
		ff.Options.XMLSkipByteOrderMark = &inputOptions.SkipByteOrderMark
	}

	return ff, nil
}

// TODO (next PRs): Rename it to QuotedItem, move to the def file, and update callers.
type NullString struct {
	S string `ddl:"parameter,no_equals,single_quotes"`
}

// CreateFileFormatOptionsLegacy is based on https://docs.snowflake.com/en/sql-reference/sql/create-file-format.
type CreateFileFormatOptionsLegacy struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary   *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	fileFormat  bool                   `ddl:"static" sql:"FILE FORMAT"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`
	Type        FileFormatType         `ddl:"parameter" sql:"TYPE"`
	FileFormatTypeOptionsLegacy
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

func (opts *CreateFileFormatOptionsLegacy) validate() error {
	fields := opts.FileFormatTypeOptionsLegacy.fieldsByType()

	for formatType := range fields {
		if opts.Type == formatType {
			continue
		}
		if anyValueSet(fields[formatType]...) {
			return fmt.Errorf("cannot set %s fields when TYPE = %s", formatType, opts.Type)
		}
	}

	err := opts.FileFormatTypeOptionsLegacy.validate()
	if err != nil {
		return err
	}

	return nil
}

func (v *fileFormatsLegacy) Create(ctx context.Context, id SchemaObjectIdentifier, opts *CreateFileFormatOptionsLegacy) error {
	if opts == nil {
		opts = &CreateFileFormatOptionsLegacy{}
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

// AlterFileFormatOptionsLegacy is based on https://docs.snowflake.com/en/sql-reference/sql/alter-file-format.
type AlterFileFormatOptionsLegacy struct {
	alter      bool                   `ddl:"static" sql:"ALTER"`
	fileFormat bool                   `ddl:"static" sql:"FILE FORMAT"`
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`

	Rename *AlterFileFormatRenameOptions
	Set    *FileFormatTypeOptionsLegacy `ddl:"list,no_comma" sql:"SET"`
}

func (opts *AlterFileFormatOptionsLegacy) validate() error {
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

type FileFormatTypeOptionsLegacy struct {
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	// CSV type options
	CSVCompression                *CsvCompression `ddl:"parameter" sql:"COMPRESSION"`
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
	CSVEncoding                   *CsvEncoding    `ddl:"parameter,single_quotes" sql:"ENCODING"`

	// JSON type options
	JSONCompression              *JsonCompression `ddl:"parameter" sql:"COMPRESSION"`
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
	XMLCompression              *XmlCompression `ddl:"parameter" sql:"COMPRESSION"`
	XMLIgnoreUTF8Errors         *bool           `ddl:"parameter" sql:"IGNORE_UTF8_ERRORS"`
	XMLPreserveSpace            *bool           `ddl:"parameter" sql:"PRESERVE_SPACE"`
	XMLStripOuterElement        *bool           `ddl:"parameter" sql:"STRIP_OUTER_ELEMENT"`
	XMLDisableSnowflakeData     *bool           `ddl:"parameter" sql:"DISABLE_SNOWFLAKE_DATA"`
	XMLDisableAutoConvert       *bool           `ddl:"parameter" sql:"DISABLE_AUTO_CONVERT"`
	XMLReplaceInvalidCharacters *bool           `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	XMLSkipByteOrderMark        *bool           `ddl:"parameter" sql:"SKIP_BYTE_ORDER_MARK"`
}

func (opts *FileFormatTypeOptionsLegacy) fieldsByType() map[FileFormatType][]any {
	return map[FileFormatType][]any{
		FileFormatTypeCsv: {
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
		FileFormatTypeJson: {
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
		FileFormatTypeOrc: {
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
		FileFormatTypeXml: {
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

func (opts *FileFormatTypeOptionsLegacy) validate() error {
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

func (v *fileFormatsLegacy) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterFileFormatOptionsLegacy) error {
	if opts == nil {
		opts = &AlterFileFormatOptionsLegacy{}
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

// DropFileFormatOptionsLegacy is based on https://docs.snowflake.com/en/sql-reference/sql/drop-file-format.
type DropFileFormatOptionsLegacy struct {
	drop       bool                   `ddl:"static" sql:"DROP"`
	fileFormat string                 `ddl:"static" sql:"FILE FORMAT"`
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropFileFormatOptionsLegacy) validate() error {
	return nil
}

func (v *fileFormatsLegacy) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropFileFormatOptionsLegacy) error {
	if opts == nil {
		opts = &DropFileFormatOptionsLegacy{}
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

func (v *fileFormatsLegacy) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropFileFormatOptionsLegacy{IfExists: Bool(true)}) }, ctx, id)
}

// ShowFileFormatsOptionsLegacy is based on https://docs.snowflake.com/en/sql-reference/sql/show-file-formats.
type ShowFileFormatsOptionsLegacy struct {
	show        bool  `ddl:"static" sql:"SHOW"`
	fileFormats bool  `ddl:"static" sql:"FILE FORMATS"`
	Like        *Like `ddl:"keyword" sql:"LIKE"`
	In          *In   `ddl:"keyword" sql:"IN"`
}

func (opts *ShowFileFormatsOptionsLegacy) validate() error {
	return nil
}

func (v *fileFormatsLegacy) Show(ctx context.Context, opts *ShowFileFormatsOptionsLegacy) ([]FileFormatLegacy, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[FileFormatRowLegacy](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[FileFormatRowLegacy, FileFormatLegacy](dbRows)
}

func (v *fileFormatsLegacy) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatLegacy, error) {
	fileFormats, err := v.client.FileFormatsLegacy.Show(ctx, &ShowFileFormatsOptionsLegacy{
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
	return collections.FindFirst(fileFormats, func(format FileFormatLegacy) bool {
		return format.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
}

func (v *fileFormatsLegacy) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatLegacy, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

type FileFormatDetailsLegacy struct {
	Type    FileFormatType
	Options FileFormatTypeOptionsLegacy
}

type FileFormatDetailsRowLegacy struct {
	Property         string
	Property_Type    string
	Property_Value   string
	Property_Default string
}

// describeFileFormatOptionsLegacy is based on https://docs.snowflake.com/en/sql-reference/sql/desc-file-format.
type describeFileFormatOptionsLegacy struct {
	describe   bool                   `ddl:"static" sql:"DESCRIBE"`
	fileFormat string                 `ddl:"static" sql:"FILE FORMAT"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *describeFileFormatOptionsLegacy) validate() error {
	return nil
}

func (v *fileFormatsLegacy) Describe(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetailsLegacy, error) {
	opts := &describeFileFormatOptionsLegacy{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	var rows []FileFormatDetailsRowLegacy
	err = v.client.query(ctx, &rows, sql)
	if err != nil {
		return nil, err
	}
	details := FileFormatDetailsLegacy{}
	for _, row := range rows {
		if row.Property == "TYPE" {
			details.Type = FileFormatType(row.Property_Value)
			break
		}
	}

	switch details.Type {
	case FileFormatTypeCsv:
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
				for s := range strings.SplitSeq(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.CSVNullIf = &newNullIf
			case "COMPRESSION":
				comp := CsvCompression(v)
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
				enc := CsvEncoding(v)
				details.Options.CSVEncoding = &enc
			}
		}
	case FileFormatTypeJson:
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
				for s := range strings.SplitSeq(strings.Trim(v, "[]"), ", ") {
					newNullIf = append(newNullIf, NullString{s})
				}
				details.Options.JSONNullIf = newNullIf
			case "COMPRESSION":
				comp := JsonCompression(v)
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
				for s := range strings.SplitSeq(strings.Trim(v, "[]"), ", ") {
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
	case FileFormatTypeOrc:
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
				for s := range strings.SplitSeq(strings.Trim(v, "[]"), ", ") {
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
				for s := range strings.SplitSeq(strings.Trim(v, "[]"), ", ") {
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
	case FileFormatTypeXml:
		for _, row := range rows {
			if row.Property_Value == "" {
				continue
			}
			v := row.Property_Value
			switch row.Property {
			case "COMPRESSION":
				comp := XmlCompression(v)
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
