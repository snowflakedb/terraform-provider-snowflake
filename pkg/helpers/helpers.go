package helpers

import (
	"encoding/csv"
	"fmt"
	"path"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	IDDelimiter = "|"
)

// StringToBool converts a string to a bool.
func StringToBool(s string) bool {
	return strings.ToLower(s) == "true"
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
