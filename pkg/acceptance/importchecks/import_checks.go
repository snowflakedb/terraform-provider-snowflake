package importchecks

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ComposeAggregateImportStateCheck does the same as ComposeImportStateCheck, but it aggregates all the occurred errors,
// instead of returning the first encountered one.
func ComposeAggregateImportStateCheck(fs ...resource.ImportStateCheckFunc) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		var result []error

		for i, f := range fs {
			if err := f(s); err != nil {
				result = append(result, fmt.Errorf("check %d/%d error: %w", i+1, len(fs), err))
			}
		}

		return errors.Join(result...)
	}
}

// ComposeImportStateCheck is based on unexported composeImportStateCheck from teststep_providers_test.go
func ComposeImportStateCheck(fs ...resource.ImportStateCheckFunc) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		for i, f := range fs {
			if err := f(s); err != nil {
				return fmt.Errorf("check %d/%d error: %w", i+1, len(fs), err)
			}
		}
		return nil
	}
}

// TestCheckResourceAttrInstanceState is based on unexported testCheckResourceAttrInstanceState from teststep_providers_test.go
func TestCheckResourceAttrInstanceState(id string, attributeName, attributeValue string) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if attrVal, ok := v.Attributes[attributeName]; ok {
				if attrVal != attributeValue {
					return fmt.Errorf("invalid value for attribute %s -  expected: %s, got: %s", attributeName, attributeValue, attrVal)
				}

				return nil
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}

// TestCheckResourceAttrNotInInstanceState is based on unexported testCheckResourceAttrInstanceState from teststep_providers_test.go,
// but instead of comparing values, it only checks if the attribute is present in the InstanceState.
func TestCheckResourceAttrNotInInstanceState(id string, attributeName string) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if _, ok := v.Attributes[attributeName]; ok {
				return fmt.Errorf("attribute %s found in instance state, but expected not to be there", attributeName)
			}
		}

		return nil
	}
}

// TestCheckResourceAttrInstanceStateSet is based on unexported testCheckResourceAttrInstanceState from teststep_providers_test.go,
// but instead of comparing values, it only checks if the value is set.
func TestCheckResourceAttrInstanceStateSet(id string, attributeName string) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if _, ok := v.Attributes[attributeName]; ok {
				return nil
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}

func TestCheckResourceAttrInstanceStateWith(id, attributeName string, checkValueFunc resource.CheckResourceAttrWithFunc) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		for _, v := range is {
			if v.ID != id {
				continue
			}

			if attrVal, ok := v.Attributes[attributeName]; ok {
				return checkValueFunc(attrVal)
			}
		}

		return fmt.Errorf("attribute %s not found in instance state", attributeName)
	}
}
