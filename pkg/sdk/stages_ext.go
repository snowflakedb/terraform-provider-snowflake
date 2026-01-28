package sdk

import (
	"strconv"
	"strings"
	"time"

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

// StageDirectoryTable represents directory table properties from DESCRIBE STAGE
type StageDirectoryTable struct {
	Enable                       bool
	AutoRefresh                  bool
	DirectoryNotificationChannel *string
	LastRefreshedOn              *time.Time
}

// StageDetails represents the parsed result of DESCRIBE STAGE
type StageDetails struct {
	FileFormatName *SchemaObjectIdentifier
	FileFormatCsv  *FileFormatCsv
	DirectoryTable *StageDirectoryTable
	PrivateLink    *StagePrivateLink
	Location       *StageLocationDetails
	Credentials    *StageCredentials
}

type StagePrivateLink struct {
	UsePrivatelinkEndpoint bool
}

// StageLocationDetails represents location properties from DESCRIBE STAGE
type StageLocationDetails struct {
	Url               string
	AwsAccessPointArn string
}

// StageCredentials represents credentials properties from DESCRIBE STAGE
type StageCredentials struct {
	AwsKeyId string
}

// ParseStageDetails parses []StageProperty into StageDetails
func ParseStageDetails(properties []StageProperty) (*StageDetails, error) {
	details := &StageDetails{}

	directoryTable, err := parseDirectoryTable(properties)
	if err != nil {
		return nil, err
	}
	details.DirectoryTable = directoryTable
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
	if fileFormatType.Value == "CSV" {
		details.FileFormatCsv = parseCsvFileFormat(properties)
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

func parseDirectoryTable(properties []StageProperty) (*StageDirectoryTable, error) {
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
				t, err := time.Parse("2006-01-02 15:04:05.000 -0700", prop.Value)
				if err != nil {
					return nil, err
				}
				dt.LastRefreshedOn = &t
			}
		}
	}

	return dt, nil
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
		switch prop.Name {
		case "USE_PRIVATELINK_ENDPOINT":
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
			loc.Url = prop.Value
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
		switch prop.Name {
		case "AWS_KEY_ID":
			creds.AwsKeyId = prop.Value
		}
	}
	return creds
}
