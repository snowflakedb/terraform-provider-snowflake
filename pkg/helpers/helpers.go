package helpers

import (
	"encoding/csv"
	"fmt"
	"log"
	"path"
	"reflect"
	"regexp"
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

// DecodeSnowflakeIDErrLegacy returns a decoded ObjectIdentifier and an error.
//
// This is a legacy function. Use Parse<level>ObjectIdentifier instead.
// This function is still used in the preview resources. The end goal is to use new Parse<level>ObjectIdentifier functions instead.
func DecodeSnowflakeIDErrLegacy[ID sdk.ObjectIdentifierConstraint](id string) (ID, error) {
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
	return containsValue(ids, id, DecodeSnowflakeParameterID, func(a, b sdk.ObjectIdentifier) bool { return a.FullyQualifiedName() == b.FullyQualifiedName() })
}

// ContainsEnum checks if the needle is contained within the haystack using the enumNormalizer to normalize the values.
func ContainsEnum[T ~string](haystack []string, needle string, enumNormalizer func(string) (T, error)) bool {
	return containsValue(haystack, needle, enumNormalizer, func(a, b T) bool { return a == b })
}

func containsValue[T comparable](haystack []string, needle string, normalize func(string) (T, error), cmp func(T, T) bool) bool {
	if len(haystack) == 0 || len(needle) == 0 {
		return false
	}

	needleNormalized, err := normalize(needle)
	if err != nil {
		return false
	}

	for _, haystackItem := range haystack {
		haystackItemNormalized, err := normalize(haystackItem)
		if err != nil {
			return false
		}
		if cmp(haystackItemNormalized, needleNormalized) {
			return true
		}
	}
	return false
}
