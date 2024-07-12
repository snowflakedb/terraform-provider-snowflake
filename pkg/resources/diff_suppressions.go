package resources

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NormalizeAndCompare[T comparable](normalize func(string) (T, error)) schema.SchemaDiffSuppressFunc {
	return func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
		oldNormalized, err := normalize(oldValue)
		if err != nil {
			return false
		}
		newNormalized, err := normalize(newValue)
		if err != nil {
			return false
		}
		return oldNormalized == newNormalized
	}
}

// IgnoreAfterCreation should be used to ignore changes to the given attribute post creation.
func IgnoreAfterCreation(_, _, _ string, d *schema.ResourceData) bool {
	// For new resources always show the diff and in every other case we do not want to use this attribute
	return d.Id() != ""
}

func IgnoreChangeToCurrentSnowflakeValueInShow(keyInShowOutput string) schema.SchemaDiffSuppressFunc {
	return func(_, _, new string, d *schema.ResourceData) bool {
		if d.Id() == "" {
			return false
		}

		if queryOutput, ok := d.GetOk(ShowOutputAttributeName); ok {
			queryOutputList := queryOutput.([]any)
			if len(queryOutputList) == 1 {
				result := queryOutputList[0].(map[string]any)
				log.Printf("[DEBUG] IgnoreChangeToCurrentSnowflakeValueInShow: value for key %s is %v, new value is %s, comparison result is: %t", keyInShowOutput, result[keyInShowOutput], new, new == fmt.Sprintf("%v", result[keyInShowOutput]))
				if new == fmt.Sprintf("%v", result[keyInShowOutput]) {
					return true
				}
			}
		}
		return false
	}
}

func IgnoreChangeToCurrentSnowflakeValueInDescribe(keyInDescribeOutput string) schema.SchemaDiffSuppressFunc {
	return func(_, _, new string, d *schema.ResourceData) bool {
		if d.Id() == "" {
			return false
		}

		if queryOutput, ok := d.GetOk(DescribeOutputAttributeName); ok {
			queryOutputList := queryOutput.([]any)
			if len(queryOutputList) == 1 {
				result := queryOutputList[0].(map[string]any)
				newValueInDescribeList := result[keyInDescribeOutput].([]any)
				if len(newValueInDescribeList) == 1 {
					newValueInDescribe := newValueInDescribeList[0].(map[string]any)["value"]
					log.Printf("[DEBUG] IgnoreChangeToCurrentSnowflakeValueInDescribe: value for key %s is %v, new value is %s, comparison result is: %t", keyInDescribeOutput, newValueInDescribe, new, new == fmt.Sprintf("%v", newValueInDescribe))
					if new == fmt.Sprintf("%v", newValueInDescribe) {
						return true
					}
				}
			}
		}
		return false
	}
}

func SuppressIfAny(diffSuppressFunctions ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		var suppress bool
		for _, f := range diffSuppressFunctions {
			suppress = suppress || f(k, old, new, d)
		}
		return suppress
	}
}

func IgnoreValuesFromSetIfParamSet(key, param string, values []string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		params := d.Get(RelatedParametersAttributeName).([]any)
		if len(params) == 0 {
			return false
		}
		result := params[0].(map[string]any)
		param := result[strings.ToLower(param)].([]any)
		value := param[0].(map[string]any)["value"]
		if !helpers.StringToBool(value.(string)) {
			return false
		}
		if k == key+".#" {
			old, new := d.GetChange(key)
			var numOld, numNew int
			oldList := expandStringList(old.(*schema.Set).List())
			newList := expandStringList(new.(*schema.Set).List())
			for _, v := range oldList {
				if !slices.Contains(values, v) {
					numOld++
				}
			}
			for _, v := range newList {
				if !slices.Contains(values, v) {
					numNew++
				}
			}
			return numOld == numNew
		}
		return slices.Contains(values, old)
	}
}

// NormalizeAndCompareIdentifiersInSet is a diff suppression function that should be used at top-level TypeSet fields that
// hold identifiers to avoid diffs like:
// - "DATABASE"."SCHEMA"."OBJECT"
// + DATABASE.SCHEMA.OBJECT
// where both identifiers are pointing to the same object, but have different structure. When a diff occurs in the
// list or set, we have to handle two suppressions (one that prevents adding and one that prevents the removal).
// It's handled by the two statements with the help of helpers.ContainsIdentifierIgnoreQuotes and by getting
// the current state of ids to compare against. The dff suppressions for lists and sets are running for each element one by one
// and the first diff is usually .# referring to the collection length (we skip those).
func NormalizeAndCompareIdentifiersInSet(key string) schema.SchemaDiffSuppressFunc {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		if strings.HasSuffix(k, ".#") {
			return false
		}

		if oldValue == "" && !d.GetRawState().IsNull() {
			if helpers.ContainsIdentifierIgnoreQuotes(ctyValToSliceString(d.GetRawState().AsValueMap()[key].AsValueSet().Values()), newValue) {
				return true
			}
		}

		if newValue == "" {
			if helpers.ContainsIdentifierIgnoreQuotes(expandStringList(d.Get(key).(*schema.Set).List()), oldValue) {
				return true
			}
		}

		return false
	}
}
