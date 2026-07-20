package sdk

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (f *StageFileFormat) additionalValidations() error {
	if f == nil {
		return nil
	}
	if valueSet(f.FileFormatOptions) {
		return f.FileFormatOptions.validate()
	}
	return nil
}

var AcceptableStageTypes = map[StageType][]StageType{
	StageTypeInternal: {StageTypeInternal, StageTypeInternalNoCse},
	StageTypeExternal: {StageTypeExternal},
}

func (v *Stage) Location() string {
	return NewStageLocation(v.ID(), "").ToSql()
}

func (s *CreateInternalStageRequest) ID() SchemaObjectIdentifier {
	return s.name
}

func (s *CreateOnS3StageRequest) ID() SchemaObjectIdentifier {
	return s.name
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

	fileFormatProperties := stageFileFormatProperties(properties)
	switch formatType {
	case FileFormatTypeCsv:
		details.FileFormatCsv = parseFileFormatCsv(fileFormatProperties, SchemaObjectIdentifier{})
	case FileFormatTypeJson:
		details.FileFormatJson = parseFileFormatJson(fileFormatProperties, SchemaObjectIdentifier{})
	case FileFormatTypeAvro:
		details.FileFormatAvro = parseFileFormatAvro(fileFormatProperties, SchemaObjectIdentifier{})
	case FileFormatTypeOrc:
		details.FileFormatOrc = parseFileFormatOrc(fileFormatProperties, SchemaObjectIdentifier{})
	case FileFormatTypeParquet:
		details.FileFormatParquet = parseFileFormatParquet(fileFormatProperties, SchemaObjectIdentifier{})
	case FileFormatTypeXml:
		details.FileFormatXml = parseFileFormatXml(fileFormatProperties, SchemaObjectIdentifier{})
	}

	return details, nil
}

// stageFileFormatProperties filters the STAGE_FILE_FORMAT-scoped properties out of a DESCRIBE
// STAGE result and converts them to FileFormatProperty so they can be parsed by the shared
// parseFileFormatCsv/Json/Avro/Orc/Parquet/Xml functions in file_formats_ext.go.
func stageFileFormatProperties(properties []StageProperty) []FileFormatProperty {
	filtered := collections.Filter(properties, func(prop StageProperty) bool {
		return prop.Parent == "STAGE_FILE_FORMAT"
	})
	return collections.Map(filtered, func(prop StageProperty) FileFormatProperty {
		return FileFormatProperty{Name: prop.Name, Type: prop.Type, Value: prop.Value, Default: prop.Default}
	})
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
