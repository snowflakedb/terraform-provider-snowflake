package resources

import (
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// borrowed from https://github.com/terraform-providers/terraform-provider-aws/blob/master/aws/structure.go#L924:6

func expandIntList(configured []any) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(int); ok {
			vs = append(vs, val)
		}
	}
	return vs
}

func expandStringList(configured []any) []string {
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

func expandStringListAllowEmpty(configured []any) []string {
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

func expandObjectIdentifier(objectIdentifier any) (string, string, string) {
	objectIdentifierMap := objectIdentifier.([]any)[0].(map[string]any)
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
func ADiffB(setA []any, setB []any) []string {
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
