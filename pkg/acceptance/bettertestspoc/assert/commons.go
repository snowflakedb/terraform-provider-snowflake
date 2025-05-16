package assert

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"golang.org/x/exp/maps"
)

// TestCheckFuncProvider is an interface with just one method providing resource.TestCheckFunc.
// It allows using it as input the "Check:" in resource.TestStep.
// It should be used with AssertThat.
type TestCheckFuncProvider interface {
	ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc
}

// AssertThat should be used for "Check:" input in resource.TestStep instead of e.g. resource.ComposeTestCheckFunc.
// It allows performing all the checks implementing the TestCheckFuncProvider interface.
func AssertThat(t *testing.T, testClient *helpers.TestClient, fs ...TestCheckFuncProvider) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		var result []error

		for i, f := range fs {
			if err := f.ToTerraformTestCheckFunc(t, testClient)(s); err != nil {
				result = append(result, fmt.Errorf("check %d/%d error:\n%w", i+1, len(fs), err))
			}
		}

		return errors.Join(result...)
	}
}

var _ TestCheckFuncProvider = (*testCheckFuncWrapper)(nil)

type testCheckFuncWrapper struct {
	f resource.TestCheckFunc
}

func (w *testCheckFuncWrapper) ToTerraformTestCheckFunc(_ *testing.T, _ *helpers.TestClient) resource.TestCheckFunc {
	return w.f
}

// Check allows using the basic terraform checks while using AssertThat.
// To use, just simply wrap the check in Check.
func Check(f resource.TestCheckFunc) TestCheckFuncProvider {
	return &testCheckFuncWrapper{f}
}

// ImportStateCheckFuncProvider is an interface with just one method providing resource.ImportStateCheckFunc.
// It allows using it as input the "ImportStateCheck:" in resource.TestStep for import tests.
// It should be used with AssertThatImport.
type ImportStateCheckFuncProvider interface {
	ToTerraformImportStateCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.ImportStateCheckFunc
}

// AssertThatImport should be used for "ImportStateCheck:" input in resource.TestStep instead of e.g. importchecks.ComposeImportStateCheck.
// It allows performing all the checks implementing the ImportStateCheckFuncProvider interface.
func AssertThatImport(t *testing.T, testClient *helpers.TestClient, fs ...ImportStateCheckFuncProvider) resource.ImportStateCheckFunc {
	t.Helper()
	return func(s []*terraform.InstanceState) error {
		var result []error

		for i, f := range fs {
			if err := f.ToTerraformImportStateCheckFunc(t, testClient)(s); err != nil {
				result = append(result, fmt.Errorf("check %d/%d error:\n%w", i+1, len(fs), err))
			}
		}

		return errors.Join(result...)
	}
}

var _ ImportStateCheckFuncProvider = (*importStateCheckFuncWrapper)(nil)

type importStateCheckFuncWrapper struct {
	f resource.ImportStateCheckFunc
}

func (w *importStateCheckFuncWrapper) ToTerraformImportStateCheckFunc(_ *testing.T, _ *helpers.TestClient) resource.ImportStateCheckFunc {
	return w.f
}

// CheckImport allows using the basic terraform import checks while using AssertThatImport.
// To use, just simply wrap the check in CheckImport.
func CheckImport(f resource.ImportStateCheckFunc) ImportStateCheckFuncProvider {
	return &importStateCheckFuncWrapper{f}
}

// InPlaceAssertionVerifier is an interface providing a method allowing verifying all the prepared assertions in place.
// It does not return function like TestCheckFuncProvider or ImportStateCheckFuncProvider; it runs all the assertions in place instead.
type InPlaceAssertionVerifier interface {
	VerifyAll(t *testing.T, testClient *helpers.TestClient)
}

// AssertThatObject should be used in the SDK tests for created object validation.
// It verifies all the prepared assertions in place.
func AssertThatObject(t *testing.T, objectAssert InPlaceAssertionVerifier, testClient *helpers.TestClient) {
	t.Helper()
	objectAssert.VerifyAll(t, testClient)
}

func ContainsExactlyInAnyOrder(resourceKey string, attributePath string, expectedItems []map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		var actualItems []map[string]string
		var resourceValue *terraform.ResourceState

		if value, ok := state.RootModule().Resources[resourceKey]; ok {
			resourceValue = value
		} else {
			return fmt.Errorf("resource %s not found", resourceKey)
		}

		// Allocate space for actualItems and assert length
		for attrKey, attrValue := range resourceValue.Primary.Attributes {
			if strings.HasPrefix(attrKey, attributePath) {
				attr := strings.TrimPrefix(attrKey, attributePath+".")

				if attr == "#" {
					attrValueLen, err := strconv.Atoi(attrValue)
					if err != nil {
						return fmt.Errorf("failed to convert length of the attribute %s: %w", attrKey, err)
					}
					if len(expectedItems) != attrValueLen {
						return fmt.Errorf("expected to find %d items in %s, but found %d", len(expectedItems), attributePath, attrValueLen)
					}

					actualItems = make([]map[string]string, attrValueLen)
					for i := range actualItems {
						actualItems[i] = make(map[string]string)
					}
				}
			}
		}

		// Gather all actual items
		for attrKey, attrValue := range resourceValue.Primary.Attributes {
			if strings.HasPrefix(attrKey, attributePath) {
				attr := strings.TrimPrefix(attrKey, attributePath+".")

				if strings.HasSuffix(attr, "%") || strings.HasSuffix(attr, "#") {
					continue
				}

				attrParts := strings.SplitN(attr, ".", 2)
				index, indexErr := strconv.Atoi(attrParts[0])
				isIndex := indexErr == nil

				if len(attrParts) > 1 && isIndex {
					itemKey := attrParts[1]
					actualItems[index][itemKey] = attrValue
				}
			}
		}

		errs := make([]error, 0)
		for _, actualItem := range actualItems {
			found := false
			if slices.ContainsFunc(expectedItems, func(expected map[string]string) bool { return maps.Equal(expected, actualItem) }) {
				found = true
			}
			if !found {
				errs = append(errs, fmt.Errorf("unexpected item found: %s", actualItem))
			}
		}

		for _, expectedItem := range expectedItems {
			found := false
			if slices.ContainsFunc(actualItems, func(actual map[string]string) bool { return maps.Equal(actual, expectedItem) }) {
				found = true
			}
			if !found {
				errs = append(errs, fmt.Errorf("expected item to be found, but it wasn't: %s", expectedItem))
			}
		}

		return errors.Join(errs...)
	}
}
