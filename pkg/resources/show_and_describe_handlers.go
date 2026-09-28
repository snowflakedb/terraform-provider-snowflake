package resources

import (
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ShowOutputAttributeName        = "show_output"
	DescribeOutputAttributeName    = "describe_output"
	ParametersAttributeName        = "parameters"
	RelatedParametersAttributeName = "related_parameters"
	ShowKeysOutputAttributeName    = "show_keys_output"
	ShowIndexesAttributeName       = "show_indexes"
)

func handleExternalChangesToObject(d *schema.ResourceData, outputAttributeName string, mappings ...outputMapping) error {
	return handleExternalChangesToObjectCmp(d, outputAttributeName, func(a, b any) bool { return a == b }, mappings...)
}

// handleExternalChangesToObjectDeepEqual compares previous show/describe state with a fresh
// mapper value. SDKv2 TypeList round-trips as []any on Get, while generated ToSchema emits
// native []string; type-strict DeepEqual would then mark every refresh as an external change.
// Both sides are projected to []any first (nil and empty slices compare equal).
// normalizeFunc stays for semantic projections (SHOW options → TRANSIENT), not type coercion.
func handleExternalChangesToObjectDeepEqual(d *schema.ResourceData, outputAttributeName string, mappings ...outputMapping) error {
	return handleExternalChangesToObjectCmp(d, outputAttributeName, sliceAwareDeepEqual, mappings...)
}

func handleExternalChangesToObjectCmp(d *schema.ResourceData, outputAttributeName string, cmpFunc func(any, any) bool, mappings ...outputMapping) error {
	if output, ok := d.GetOk(outputAttributeName); ok {
		outputList := output.([]any)
		if len(outputList) == 1 {
			result := outputList[0].(map[string]any)
			for _, mapping := range mappings {
				valueToCompareFrom := result[mapping.nameInOutput]
				if mapping.normalizeFunc != nil {
					valueToCompareFrom = mapping.normalizeFunc(valueToCompareFrom)
				}
				if !cmpFunc(valueToCompareFrom, mapping.valueToCompare) {
					if err := d.Set(mapping.nameInConfig, mapping.valueToSet); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// sliceAwareDeepEqual is reflect.DeepEqual with TypeList-friendly slice comparison:
// []string vs []any with the same elements are equal, and nil/empty slices are equal.
func sliceAwareDeepEqual(a, b any) bool {
	as, aSlice := asAnySlice(a)
	bs, bSlice := asAnySlice(b)
	switch {
	case aSlice && bSlice:
		return reflect.DeepEqual(as, bs)
	case aSlice && b == nil:
		return len(as) == 0
	case bSlice && a == nil:
		return len(bs) == 0
	default:
		return reflect.DeepEqual(a, b)
	}
}

func asAnySlice(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return nil, false
	}
	if rv.IsNil() || rv.Len() == 0 {
		return []any{}, true
	}
	out := make([]any, rv.Len())
	for i := range out {
		out[i] = rv.Index(i).Interface()
	}
	return out, true
}

func handleExternalChangesToObjectInFlatDescribeDeepEqual(d *schema.ResourceData, mappings ...outputMapping) error {
	return handleExternalChangesToObjectDeepEqual(d, DescribeOutputAttributeName, mappings...)
}

// handleExternalChangesToObjectInShow assumes that show output is kept in ShowOutputAttributeName attribute
func handleExternalChangesToObjectInShow(d *schema.ResourceData, mappings ...outputMapping) error {
	return handleExternalChangesToObject(d, ShowOutputAttributeName, mappings...)
}

// handleExternalChangesToObjectInFlatDescribe assumes that describe output is kept in DescribeOutputAttributeName attribute
// It is to be used with flat - (show-like) describe_output schemas
// To handle external changes to describe with properties like collections use `handleExternalChangesToObjectInDescribe()`
func handleExternalChangesToObjectInFlatDescribe(d *schema.ResourceData, mappings ...outputMapping) error {
	return handleExternalChangesToObject(d, DescribeOutputAttributeName, mappings...)
}

type outputMapping struct {
	nameInOutput   string
	nameInConfig   string
	valueToCompare any
	valueToSet     any
	normalizeFunc  func(any) any
}

// handleExternalChangesToObjectInDescribe assumes that describe output is kept in DescribeOutputAttributeName attribute
func handleExternalChangesToObjectInDescribe(d *schema.ResourceData, mappings ...describeMapping) error {
	if describeOutput, ok := d.GetOk(DescribeOutputAttributeName); ok {
		describeOutputList := describeOutput.([]any)
		if len(describeOutputList) == 1 {
			result := describeOutputList[0].(map[string]any)

			for _, mapping := range mappings {
				if result[mapping.nameInDescribe] == nil {
					continue
				}

				valueToCompareFromList := result[mapping.nameInDescribe].([]any)
				if len(valueToCompareFromList) != 1 {
					continue
				}

				valueToCompareFrom := valueToCompareFromList[0].(map[string]any)["value"]
				if mapping.normalizeFunc != nil {
					valueToCompareFrom = mapping.normalizeFunc(valueToCompareFrom)
				}
				if valueToCompareFrom != mapping.valueToCompare {
					if err := d.Set(mapping.nameInConfig, mapping.valueToSet); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

type describeMapping struct {
	nameInDescribe string
	nameInConfig   string
	valueToCompare any
	valueToSet     any
	normalizeFunc  func(any) any
}

// setStateToValuesFromConfig currently handles only int, float, and string types.
// It's needed for the case where:
// - previous config was empty (therefore Snowflake defaults had been used)
// - new config have the same values that are already in SF
func setStateToValuesFromConfig(d *schema.ResourceData, resourceSchema map[string]*schema.Schema, fields []string) error {
	if !d.GetRawConfig().IsNull() {
		vMap := d.GetRawConfig().AsValueMap()
		for _, field := range fields {
			if v, ok := vMap[field]; ok && !v.IsNull() {
				if schemaField, ok := resourceSchema[field]; ok {
					switch schemaField.Type {
					case schema.TypeInt:
						intVal, _ := v.AsBigFloat().Int64()
						if err := d.Set(field, intVal); err != nil {
							return err
						}
					case schema.TypeFloat:
						if err := d.Set(field, v.AsBigFloat()); err != nil {
							return err
						}
					case schema.TypeString:
						if err := d.Set(field, v.AsString()); err != nil {
							return err
						}
					case schema.TypeSet:
						if err := d.Set(field, ctyValToSliceString(v.AsValueSlice())); err != nil {
							return err
						}
					default:
						log.Printf("[DEBUG] field %s has unsupported schema type %v not found", field, schemaField.Type)
					}
				} else {
					log.Printf("[DEBUG] schema field %s not found", field)
				}
			}
		}
	}
	return nil
}
