package resources

import (
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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
