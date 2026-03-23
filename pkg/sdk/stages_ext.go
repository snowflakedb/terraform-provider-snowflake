package sdk

import (
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (v *Stage) Location() string {
	return NewStageLocation(v.ID(), "").ToSql()
}

func (s *CreateInternalStageRequest) ID() SchemaObjectIdentifier {
	return s.name
}

func (s *CreateOnS3StageRequest) ID() SchemaObjectIdentifier {
	return s.name
}

// FileFormatCsv represents CSV file format properties from DESCRIBE STAGE
type FileFormatCsv struct {
	Type                       string
	RecordDelimiter            string
	FieldDelimiter             string
	FileExtension              string
	SkipHeader                 int
	ParseHeader                bool
	DateFormat                 string
	TimeFormat                 string
	TimestampFormat            string
	BinaryFormat               string
	Escape                     string
	EscapeUnenclosedField      string
	TrimSpace                  bool
	FieldOptionallyEnclosedBy  string
	NullIf                     []string
	Compression                string
	ErrorOnColumnCountMismatch bool
	ValidateUtf8               bool
	SkipBlankLines             bool
	ReplaceInvalidCharacters   bool
	EmptyFieldAsNull           bool
	SkipByteOrderMark          bool
	Encoding                   string
	MultiLine                  bool
}

// FileFormatJson represents JSON file format properties from DESCRIBE STAGE
type FileFormatJson struct {
	Type                     string
	Compression              string
	DateFormat               string
	TimeFormat               string
	TimestampFormat          string
	BinaryFormat             string
	TrimSpace                bool
	MultiLine                bool
	NullIf                   []string
	FileExtension            string
	EnableOctal              bool
	AllowDuplicate           bool
	StripOuterArray          bool
	StripNullValues          bool
	ReplaceInvalidCharacters bool
	IgnoreUtf8Errors         bool
	SkipByteOrderMark        bool
}

// FileFormatAvro represents AVRO file format properties from DESCRIBE STAGE
type FileFormatAvro struct {
	Type                     string
	Compression              string
	TrimSpace                bool
	ReplaceInvalidCharacters bool
	NullIf                   []string
}

// FileFormatOrc represents ORC file format properties from DESCRIBE STAGE
type FileFormatOrc struct {
	Type                     string
	TrimSpace                bool
	ReplaceInvalidCharacters bool
	NullIf                   []string
}

// FileFormatParquet represents Parquet file format properties from DESCRIBE STAGE
type FileFormatParquet struct {
	Type                     string
	Compression              string
	BinaryAsText             bool
	UseLogicalType           bool
	TrimSpace                bool
	UseVectorizedScanner     bool
	ReplaceInvalidCharacters bool
	NullIf                   []string
}

// FileFormatXml represents XML file format properties from DESCRIBE STAGE
type FileFormatXml struct {
	Type                     string
	Compression              string
	IgnoreUtf8Errors         bool
	PreserveSpace            bool
	StripOuterElement        bool
	DisableAutoConvert       bool
	ReplaceInvalidCharacters bool
	SkipByteOrderMark        bool
}

// StageDirectoryTable represents directory table properties from DESCRIBE STAGE
type StageDirectoryTable struct {
	Enable                       bool
	AutoRefresh                  bool
	DirectoryNotificationChannel *string
	LastRefreshedOn              *string
}

// StageDetails represents the parsed result of DESCRIBE STAGE
type StageDetails struct {
	FileFormatName    *SchemaObjectIdentifier
	FileFormatCsv     *FileFormatCsv
	FileFormatJson    *FileFormatJson
	FileFormatAvro    *FileFormatAvro
	FileFormatOrc     *FileFormatOrc
	FileFormatParquet *FileFormatParquet
	FileFormatXml     *FileFormatXml
	DirectoryTable    *StageDirectoryTable
	PrivateLink       *StagePrivateLink
	Location          *StageLocationDetails
	Credentials       *StageCredentials
}

type StagePrivateLink struct {
	UsePrivatelinkEndpoint bool
}

// StageLocationDetails represents location properties from DESCRIBE STAGE
type StageLocationDetails struct {
	Url               []string
	AwsAccessPointArn string
}

// StageCredentials represents credentials properties from DESCRIBE STAGE
type StageCredentials struct {
	AwsKeyId string
}

// ParseStageDetails parses []StageProperty into StageDetails
func ParseStageDetails(properties []StageProperty) (*StageDetails, error) {
	details := &StageDetails{}

	details.DirectoryTable = parseDirectoryTable(properties)
	details.PrivateLink = parsePrivateLink(properties)
	details.Location = parseStageLocationDetails(properties)
	details.Credentials = parseStageCredentials(properties)
	details.FileFormatName = parseFileFormatName(properties)

	fileFormatType, err := collections.FindFirst(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT" && prop.Name == "TYPE"
	})
	if err != nil {
		return details, nil
	}
	formatType, err := ToFileFormatType(fileFormatType.Value)
	if err != nil {
		return nil, err
	}
	switch formatType {
	case FileFormatTypeCSV:
		details.FileFormatCsv = parseCsvFileFormat(properties)
	case FileFormatTypeJSON:
		details.FileFormatJson = parseJsonFileFormat(properties)
	case FileFormatTypeAvro:
		details.FileFormatAvro = parseAvroFileFormat(properties)
	case FileFormatTypeORC:
		details.FileFormatOrc = parseOrcFileFormat(properties)
	case FileFormatTypeParquet:
		details.FileFormatParquet = parseParquetFileFormat(properties)
	case FileFormatTypeXML:
		details.FileFormatXml = parseXmlFileFormat(properties)
	}

	return details, nil
}

func parseFileFormatName(properties []StageProperty) *SchemaObjectIdentifier {
	fileFormatName, err := collections.FindFirst(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT" && prop.Name == "FORMAT_NAME"
	})
	if err != nil {
		return nil
	}
	idRaw := strings.ReplaceAll(fileFormatName.Value, `\"`, `"`)
	id, err := ParseSchemaObjectIdentifier(idRaw)
	if err != nil {
		return nil
	}
	return &id
}

func parseCsvFileFormat(properties []StageProperty) *FileFormatCsv {
	csv := &FileFormatCsv{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "TYPE":
			csv.Type = prop.Value
		case "RECORD_DELIMITER":
			csv.RecordDelimiter = prop.Value
		case "FIELD_DELIMITER":
			csv.FieldDelimiter = prop.Value
		case "FILE_EXTENSION":
			csv.FileExtension = prop.Value
		case "SKIP_HEADER":
			val, _ := strconv.Atoi(prop.Value)
			csv.SkipHeader = val
		case "PARSE_HEADER":
			csv.ParseHeader = prop.Value == "true"
		case "DATE_FORMAT":
			csv.DateFormat = prop.Value
		case "TIME_FORMAT":
			csv.TimeFormat = prop.Value
		case "TIMESTAMP_FORMAT":
			csv.TimestampFormat = prop.Value
		case "BINARY_FORMAT":
			csv.BinaryFormat = prop.Value
		case "ESCAPE":
			csv.Escape = prop.Value
		case "ESCAPE_UNENCLOSED_FIELD":
			csv.EscapeUnenclosedField = prop.Value
		case "TRIM_SPACE":
			csv.TrimSpace = prop.Value == "true"
		case "FIELD_OPTIONALLY_ENCLOSED_BY":
			csv.FieldOptionallyEnclosedBy = prop.Value
		case "NULL_IF":
			csv.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMPRESSION":
			csv.Compression = prop.Value
		case "ERROR_ON_COLUMN_COUNT_MISMATCH":
			csv.ErrorOnColumnCountMismatch = prop.Value == "true"
		case "VALIDATE_UTF8":
			csv.ValidateUtf8 = prop.Value == "true"
		case "SKIP_BLANK_LINES":
			csv.SkipBlankLines = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			csv.ReplaceInvalidCharacters = prop.Value == "true"
		case "EMPTY_FIELD_AS_NULL":
			csv.EmptyFieldAsNull = prop.Value == "true"
		case "SKIP_BYTE_ORDER_MARK":
			csv.SkipByteOrderMark = prop.Value == "true"
		case "ENCODING":
			csv.Encoding = prop.Value
		case "MULTI_LINE":
			csv.MultiLine = prop.Value == "true"
		}
	}

	return csv
}

func parseJsonFileFormat(properties []StageProperty) *FileFormatJson {
	json := &FileFormatJson{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "TYPE":
			json.Type = prop.Value
		case "COMPRESSION":
			json.Compression = prop.Value
		case "DATE_FORMAT":
			json.DateFormat = prop.Value
		case "TIME_FORMAT":
			json.TimeFormat = prop.Value
		case "TIMESTAMP_FORMAT":
			json.TimestampFormat = prop.Value
		case "BINARY_FORMAT":
			json.BinaryFormat = prop.Value
		case "TRIM_SPACE":
			json.TrimSpace = prop.Value == "true"
		case "MULTI_LINE":
			json.MultiLine = prop.Value == "true"
		case "NULL_IF":
			json.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		case "FILE_EXTENSION":
			json.FileExtension = prop.Value
		case "ENABLE_OCTAL":
			json.EnableOctal = prop.Value == "true"
		case "ALLOW_DUPLICATE":
			json.AllowDuplicate = prop.Value == "true"
		case "STRIP_OUTER_ARRAY":
			json.StripOuterArray = prop.Value == "true"
		case "STRIP_NULL_VALUES":
			json.StripNullValues = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			json.ReplaceInvalidCharacters = prop.Value == "true"
		case "IGNORE_UTF8_ERRORS":
			json.IgnoreUtf8Errors = prop.Value == "true"
		case "SKIP_BYTE_ORDER_MARK":
			json.SkipByteOrderMark = prop.Value == "true"
		}
	}

	return json
}

func parseAvroFileFormat(properties []StageProperty) *FileFormatAvro {
	avro := &FileFormatAvro{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "TYPE":
			avro.Type = prop.Value
		case "COMPRESSION":
			avro.Compression = prop.Value
		case "TRIM_SPACE":
			avro.TrimSpace = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			avro.ReplaceInvalidCharacters = prop.Value == "true"
		case "NULL_IF":
			avro.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}

	return avro
}

func parseOrcFileFormat(properties []StageProperty) *FileFormatOrc {
	orc := &FileFormatOrc{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "TYPE":
			orc.Type = prop.Value
		case "TRIM_SPACE":
			orc.TrimSpace = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			orc.ReplaceInvalidCharacters = prop.Value == "true"
		case "NULL_IF":
			orc.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}

	return orc
}

func parseParquetFileFormat(properties []StageProperty) *FileFormatParquet {
	parquet := &FileFormatParquet{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "TYPE":
			parquet.Type = prop.Value
		case "COMPRESSION":
			parquet.Compression = prop.Value
		case "BINARY_AS_TEXT":
			parquet.BinaryAsText = prop.Value == "true"
		case "USE_LOGICAL_TYPE":
			parquet.UseLogicalType = prop.Value == "true"
		case "TRIM_SPACE":
			parquet.TrimSpace = prop.Value == "true"
		case "USE_VECTORIZED_SCANNER":
			parquet.UseVectorizedScanner = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			parquet.ReplaceInvalidCharacters = prop.Value == "true"
		case "NULL_IF":
			parquet.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}

	return parquet
}

func parseXmlFileFormat(properties []StageProperty) *FileFormatXml {
	xml := &FileFormatXml{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "TYPE":
			xml.Type = prop.Value
		case "COMPRESSION":
			xml.Compression = prop.Value
		case "IGNORE_UTF8_ERRORS":
			xml.IgnoreUtf8Errors = prop.Value == "true"
		case "PRESERVE_SPACE":
			xml.PreserveSpace = prop.Value == "true"
		case "STRIP_OUTER_ELEMENT":
			xml.StripOuterElement = prop.Value == "true"
		case "DISABLE_AUTO_CONVERT":
			xml.DisableAutoConvert = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			xml.ReplaceInvalidCharacters = prop.Value == "true"
		case "SKIP_BYTE_ORDER_MARK":
			xml.SkipByteOrderMark = prop.Value == "true"
		}
	}

	return xml
}

func parseDirectoryTable(properties []StageProperty) *StageDirectoryTable {
	dt := &StageDirectoryTable{}
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "DIRECTORY"
	})
	for _, prop := range filtered {
		switch prop.Name {
		case "ENABLE":
			dt.Enable = prop.Value == "true"
		case "AUTO_REFRESH":
			dt.AutoRefresh = prop.Value == "true"
		case "DIRECTORY_NOTIFICATION_CHANNEL":
			dt.DirectoryNotificationChannel = &prop.Value
		case "LAST_REFRESHED_ON":
			if prop.Value != "" {
				dt.LastRefreshedOn = &prop.Value
			}
		}
	}

	return dt
}

func parsePrivateLink(properties []StageProperty) *StagePrivateLink {
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "PRIVATELINK"
	})
	if len(filtered) == 0 {
		return nil
	}
	pl := &StagePrivateLink{}
	for _, prop := range filtered {
		if prop.Name == "USE_PRIVATELINK_ENDPOINT" {
			pl.UsePrivatelinkEndpoint = prop.Value == "true"
		}
	}
	return pl
}

func parseStageLocationDetails(properties []StageProperty) *StageLocationDetails {
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_LOCATION"
	})
	if len(filtered) == 0 {
		return nil
	}
	loc := &StageLocationDetails{}
	for _, prop := range filtered {
		switch prop.Name {
		case "URL":
			loc.Url = ParseCommaSeparatedStringArray(prop.Value, true)
		case "AWS_ACCESS_POINT_ARN":
			loc.AwsAccessPointArn = prop.Value
		}
	}
	return loc
}

func parseStageCredentials(properties []StageProperty) *StageCredentials {
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_CREDENTIALS"
	})
	if len(filtered) == 0 {
		return nil
	}
	creds := &StageCredentials{}
	for _, prop := range filtered {
		if prop.Name == "AWS_KEY_ID" {
			creds.AwsKeyId = prop.Value
		}
	}
	return creds
}
