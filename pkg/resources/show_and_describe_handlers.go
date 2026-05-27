package resources

import (
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ShowOutputAttributeName        = "show_output"
	DescribeOutputAttributeName    = "describe_output"
	ParametersAttributeName        = "parameters"
	RelatedParametersAttributeName = "related_parameters"
)

func handleExternalChangesToObject(d *schema.ResourceData, outputAttributeName string, mappings ...outputMapping) error {
	return handleExternalChangesToObjectCmp(d, outputAttributeName, func(a, b any) bool { return a == b }, mappings...)
}

func handleExternalChangesToObjectDeepEqual(d *schema.ResourceData, outputAttributeName string, mappings ...outputMapping) error {
	return handleExternalChangesToObjectCmp(d, outputAttributeName, reflect.DeepEqual, mappings...)
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
