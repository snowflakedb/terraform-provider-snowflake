package resources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// borrowed from https://github.com/terraform-providers/terraform-provider-aws/blob/master/aws/structure.go#L924:6

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(int); ok {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func ExpandObjectIdentifierSet(configured []any, objectType sdk.ObjectType) ([]sdk.ObjectIdentifier, error) {
	vs := expandStringList(configured)
	ids := make([]sdk.ObjectIdentifier, len(vs))
	for i, idRaw := range vs {
		var id sdk.ObjectIdentifier
		var err error
		// TODO(SNOW-1229218): Use a common mapper to get object id.
		if objectType == sdk.ObjectTypeAccount {
			id, err = sdk.ParseAccountIdentifier(idRaw)
			if err != nil {
				return nil, fmt.Errorf("invalid account id: %w", err)
			}
		} else {
			id, err = GetOnObjectIdentifier(objectType, idRaw)
			if err != nil {
				return nil, fmt.Errorf("invalid object id: %w", err)
			}
		}
		ids[i] = id
	}
	return ids, nil
}

func expandStringListAllowEmpty(configured []interface{}) []string {
	// Allow empty values during expansion
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok {
			vs = append(vs, val)
		} else {
			vs = append(vs, "")
		}
	}
	return vs
}

func expandObjectIdentifier(objectIdentifier interface{}) (string, string, string) {
	objectIdentifierMap := objectIdentifier.([]interface{})[0].(map[string]interface{})
	objectName := objectIdentifierMap["name"].(string)
	var objectSchema string
	if v := objectIdentifierMap["schema"]; v != nil {
		objectSchema = v.(string)
	}
	var objectDatabase string
	if v := objectIdentifierMap["database"]; v != nil {
		objectDatabase = v.(string)
	}
	return objectDatabase, objectSchema, objectName
}

// ADiffB takes all the elements of A that are not also present in B, A-B in set notation
func ADiffB(setA []interface{}, setB []interface{}) []string {
	res := make([]string, 0)
	sliceA := expandStringList(setA)
	sliceB := expandStringList(setB)
	for _, s := range sliceA {
		if !slices.Contains(sliceB, s) {
			res = append(res, s)
		}
	}
	return res
}

// quoteColumnNames wraps each column name in double quotes for SQL compatibility.
func quoteColumnNames(columns []string) []string {
	quoted := make([]string, len(columns))
	for i, col := range columns {
		quoted[i] = fmt.Sprintf(`"%s"`, col)
	}
	return quoted
}

// needsQuoting determines if an identifier needs to be quoted based on Snowflake rules.
// Identifiers need quoting if they:
// - Contain special characters (space, hyphen, @, etc.)
// - Have mixed case (e.g., "userId" contains both upper and lower)
// - Start with a digit
// - Are reserved words (simplified check)
func needsQuoting(identifier string) bool {
	if identifier == "" {
		return false
	}

	// Check for special characters (anything not alphanumeric or underscore)
	for _, r := range identifier {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return true
		}
	}

	// Check if starts with digit
	firstChar := rune(identifier[0])
	if firstChar >= '0' && firstChar <= '9' {
		return true
	}

	// Check for mixed case
	hasUpper := false
	hasLower := false
	for _, r := range identifier {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		}
		if r >= 'a' && r <= 'z' {
			hasLower = true
		}
		if hasUpper && hasLower {
			return true
		}
	}

	// If we get here, it's a simple identifier (all upper, all lower, or numeric with underscore)
	// These don't need quoting
	return false
}

// quoteIdentifierIfNeeded quotes an identifier only if it needs quoting
func quoteIdentifierIfNeeded(identifier string) string {
	if needsQuoting(identifier) {
		return fmt.Sprintf(`"%s"`, identifier)
	}
	return identifier
}

// quoteColumnNamesSelectively applies selective quoting to column names
func quoteColumnNamesSelectively(columns []string) []string {
	quoted := make([]string, len(columns))
	for i, col := range columns {
		quoted[i] = quoteIdentifierIfNeeded(col)
	}
	return quoted
}

// NormalizeIdentifier returns the normalized form of an identifier for comparison
// Unquoted identifiers are uppercased by Snowflake, quoted identifiers preserve case
func NormalizeIdentifier(identifier string) string {
	if needsQuoting(identifier) {
		// If it needs quoting, return as-is (quoted identifiers preserve case)
		return identifier
	}
	// Simple identifiers are uppercased by Snowflake
	return strings.ToUpper(identifier)
}

func reorderStringList(configured []string, actual []string) []string {
	// Reorder the actual list to match the configured list
	// This is necessary because the actual list may not be saved in the same order as the configured list
	// The actual list may not be the same size as the configured list and may contain items not in the configured list

	// Create a map of the actual list
	actualMap := make(map[string]bool)
	for _, v := range actual {
		actualMap[v] = true
	}
	reorderedList := make([]string, 0)
	for _, v := range configured {
		if _, ok := actualMap[v]; ok {
			reorderedList = append(reorderedList, v)
		}
	}
	// add any items in the actual list that are not in the configured list to the end
	for _, v := range actual {
		if _, ok := actualMap[v]; !ok {
			reorderedList = append(reorderedList, v)
		}
	}
	return reorderedList
}
