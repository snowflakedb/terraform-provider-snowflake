package helpers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	IDDelimiter = "|"
)

// ListContentToString strips list elements of double quotes or brackets.
func ListContentToString(listString string) string {
	re := regexp.MustCompile(`[\"\[\]]`)
	return re.ReplaceAllString(listString, "")
}

// StringToBool converts a string to a bool.
func StringToBool(s string) bool {
	return strings.ToLower(s) == "true"
}

// EncodeSnowflakeID generates a unique ID for a resource.
func EncodeSnowflakeID(attributes ...interface{}) string {
	// is attribute already an object identifier?
	if len(attributes) == 1 {
		if id, ok := attributes[0].(sdk.ObjectIdentifier); ok {
			if val := reflect.ValueOf(id); val.Kind() == reflect.Ptr && val.IsNil() {
				log.Panicf("Nil object identifier received")
			}
			parts := make([]string, 0)
			switch v := id.(type) {
			case sdk.AccountObjectIdentifier:
				parts = append(parts, v.Name())
			case *sdk.AccountObjectIdentifier:
				parts = append(parts, v.Name())
			case sdk.DatabaseObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.Name())
			case *sdk.DatabaseObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.Name())
			case sdk.SchemaObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.Name())
			case *sdk.SchemaObjectIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.Name())
			case sdk.TableColumnIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.TableName(), v.Name())
			case *sdk.TableColumnIdentifier:
				parts = append(parts, v.DatabaseName(), v.SchemaName(), v.TableName(), v.Name())
			default:
				log.Panicf("Unsupported object identifier: %v", id)
			}
			return strings.Join(parts, IDDelimiter)
		}
	}
	var parts []string
	for i, attr := range attributes {
		if attr == nil {
			attributes[i] = ""
		}
		switch reflect.TypeOf(attr).Kind() {
		case reflect.String:
			parts = append(parts, attr.(string))
		case reflect.Bool:
			parts = append(parts, strconv.FormatBool(attr.(bool)))
		case reflect.Slice:
			parts = append(parts, strings.Join(attr.([]string), ","))
		}
	}
	return strings.Join(parts, "|")
}

// DecodeSnowflakeIDLegacy decodes a Snowflake ID into an ObjectIdentifier.
//
// This is a legacy function. Use Parse<level>ObjectIdentifier instead.
// This function is still used in the preview resources. The end goal is to use new Parse<level>ObjectIdentifier functions instead.
// Still, this could be later considered as a test helper in a separate package.
func DecodeSnowflakeIDLegacy(id string) sdk.ObjectIdentifier {
	parts := strings.Split(id, IDDelimiter)
	switch len(parts) {
	case 1:
		return sdk.NewAccountObjectIdentifier(parts[0])
	case 2:
		return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1])
	case 3:
		return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2])
	case 4:
		return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3])
	default:
		return nil
	}
}

func DecodeSnowflakeIDErr[ID sdk.ObjectIdentifierConstraint](id string) (ID, error) {
	return DecodeSnowflakeIDLegacy(id).(ID), nil
}

// DecodeSnowflakeParameterID decodes identifier (usually passed as one of the parameter in tf configuration) into sdk.ObjectIdentifier.
// identifier can be specified in two ways: quoted and unquoted, e.g.
//
// quoted { "some_identifier": "\"database.name\".\"schema.name\".\"test.name\" }
// (note that here dots as part of the name are allowed)
//
// unquoted { "some_identifier": "database_name.schema_name.test_name" }
// (note that here dots as part of the name are NOT allowed, because they're treated in this case as dividers)
//
// The following configuration { "some_identifier": "db.name" } will be parsed as an object called "name" that lives
// inside database called "db", not a database called "db.name". In this case quotes should be used.
func DecodeSnowflakeParameterID(identifier string) (sdk.ObjectIdentifier, error) {
	parts, err := sdk.ParseIdentifierString(identifier)
	if err != nil {
		return nil, err
	}
	switch len(parts) {
	case 1:
		return sdk.NewAccountObjectIdentifier(parts[0]), nil
	case 2:
		return sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]), nil
	case 3:
		return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], parts[2]), nil
	case 4:
		return sdk.NewTableColumnIdentifier(parts[0], parts[1], parts[2], parts[3]), nil
	default:
		return nil, fmt.Errorf("unable to classify identifier: %s", identifier)
	}
}

// DecodeSnowflakeAccountIdentifier decodes account identifier (usually passed as one of the parameter in tf configuration) into sdk.AccountIdentifier.
// Check more in https://docs.snowflake.com/en/sql-reference/sql/create-account#required-parameters.
func DecodeSnowflakeAccountIdentifier(identifier string) (sdk.AccountIdentifier, error) {
	parts, err := sdk.ParseIdentifierString(identifier)
	if err != nil {
		return sdk.AccountIdentifier{}, err
	}
	switch len(parts) {
	case 1:
		return sdk.AccountIdentifier{}, fmt.Errorf("identifier: %s seems to be account locator and these are not allowed - please use <organization_name>.<account_name>", identifier)
	case 2:
		return sdk.NewAccountIdentifier(parts[0], parts[1]), nil
	default:
		return sdk.AccountIdentifier{}, fmt.Errorf("unable to classify account identifier: %s, expected format: <organization_name>.<account_name>", identifier)
	}
}

type StorageLocation struct {
	Name                 string `json:"NAME"`
	StorageProvider      string `json:"STORAGE_PROVIDER"`
	StorageBaseUrl       string `json:"STORAGE_BASE_URL"`
	StorageAwsRoleArn    string `json:"STORAGE_AWS_ROLE_ARN"`
	StorageAwsExternalId string `json:"STORAGE_AWS_EXTERNAL_ID"`
	EncryptionType       string `json:"ENCRYPTION_TYPE"`
	EncryptionKmsKeyId   string `json:"ENCRYPTION_KMS_KEY_ID"`
	AzureTenantId        string `json:"AZURE_TENANT_ID"`
}

func validateParsedExternalVolumeDescribed(p ParsedExternalVolumeDescribed) error {
	if len(p.StorageLocations) == 0 {
		return fmt.Errorf("No storage locations could be parsed from the external volume.")
	}
	if len(p.AllowWrites) == 0 {
		return fmt.Errorf("The external volume AllowWrites property could not be parsed.")
	}

	for _, s := range p.StorageLocations {
		if len(s.Name) == 0 {
			return fmt.Errorf("A storage location's Name in this volume could not be parsed.")
		}
		if !slices.Contains(sdk.AsStringList(sdk.AllStorageProviderValues), s.StorageProvider) {
			return fmt.Errorf("invalid storage provider parsed: %s", s)
		}
		if len(s.StorageBaseUrl) == 0 {
			return fmt.Errorf("A storage location's StorageBaseUrl in this volume could not be parsed.")
		}

		storageProvider, err := sdk.ToStorageProvider(s.StorageProvider)
		if err != nil {
			return err
		}

		switch storageProvider {
		case sdk.StorageProviderS3, sdk.StorageProviderS3GOV:
			if len(s.StorageAwsRoleArn) == 0 {
				return fmt.Errorf("An S3 storage location's StorageAwsRoleArn in this volume could not be parsed.")
			}
		case sdk.StorageProviderAzure:
			if len(s.AzureTenantId) == 0 {
				return fmt.Errorf("An Azure storage location's AzureTenantId in this volume could not be parsed.")
			}
		}
	}

	return nil
}

type ParsedExternalVolumeDescribed struct {
	StorageLocations []StorageLocation
	Active           string
	Comment          string
	AllowWrites      string
}

func ParseExternalVolumeDescribed(props []sdk.ExternalVolumeProperty) (ParsedExternalVolumeDescribed, error) {
	parsedExternalVolumeDescribed := ParsedExternalVolumeDescribed{}
	var storageLocations []StorageLocation
	for _, p := range props {
		switch {
		case p.Name == "COMMENT":
			parsedExternalVolumeDescribed.Comment = p.Value
		case p.Name == "ACTIVE":
			parsedExternalVolumeDescribed.Active = p.Value
		case p.Name == "ALLOW_WRITES":
			parsedExternalVolumeDescribed.AllowWrites = p.Value
		case strings.Contains(p.Name, "STORAGE_LOCATION_"):
			storageLocation := StorageLocation{}
			err := json.Unmarshal([]byte(p.Value), &storageLocation)
			if err != nil {
				return ParsedExternalVolumeDescribed{}, err
			}
			storageLocations = append(
				storageLocations,
				storageLocation,
			)
		default:
			return ParsedExternalVolumeDescribed{}, fmt.Errorf("Unrecognized external volume property: %s", p.Name)
		}
	}

	parsedExternalVolumeDescribed.StorageLocations = storageLocations
	err := validateParsedExternalVolumeDescribed(parsedExternalVolumeDescribed)
	if err != nil {
		return ParsedExternalVolumeDescribed{}, err
	}

	return parsedExternalVolumeDescribed, nil
}

// TODO(SNOW-1569530): address during identifiers rework follow-up
func ParseRootLocation(location string) (sdk.SchemaObjectIdentifier, string, error) {
	location = strings.TrimPrefix(location, "@")
	parts, err := sdk.ParseIdentifierStringWithOpts(location, func(r *csv.Reader) {
		r.Comma = '.'
		r.LazyQuotes = true
	})
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, "", err
	}
	if len(parts) < 3 {
		return sdk.SchemaObjectIdentifier{}, "", fmt.Errorf("expected 3 parts for location %s, got %d", location, len(parts))
	}
	parts[2] = strings.Join(parts[2:], ".")
	lastParts := strings.Split(parts[2], "/")
	return sdk.NewSchemaObjectIdentifier(parts[0], parts[1], lastParts[0]), path.Join(lastParts[1:]...), nil
}

// ContainsIdentifierIgnoringQuotes takes ids (a slice of Snowflake identifiers represented as strings), and
// id (a string representing Snowflake id). It checks if id is contained within ids ignoring quotes around identifier parts.
//
// The original quoting should be retrieved to avoid situations like "object" == "\"object\"" (true)
// where that should not be a truthful comparison (different ids). Right now, we assume this case won't happen because the quoting difference would only appear
// in cases where the identifier parts are upper-cased and returned without quotes by snowflake, e.g. "OBJECT" == "\"OBJECT\"" (true)
// which is correct (the same ids).
func ContainsIdentifierIgnoringQuotes(ids []string, id string) bool {
	if len(ids) == 0 || len(id) == 0 {
		return false
	}

	idToCompare, err := DecodeSnowflakeParameterID(id)
	if err != nil {
		return false
	}

	for _, stringId := range ids {
		objectIdentifier, err := DecodeSnowflakeParameterID(stringId)
		if err != nil {
			return false
		}
		if idToCompare.FullyQualifiedName() == objectIdentifier.FullyQualifiedName() {
			return true
		}
	}

	return false
}
