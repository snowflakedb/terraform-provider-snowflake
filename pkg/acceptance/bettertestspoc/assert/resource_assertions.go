package assert

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	_ TestCheckFuncProvider        = (*ResourceAssert)(nil)
	_ ImportStateCheckFuncProvider = (*ResourceAssert)(nil)
)

// ResourceAssert is an embeddable struct that should be used to construct new resource assertions (for resource, show output, parameters, etc.).
// It implements both TestCheckFuncProvider and ImportStateCheckFuncProvider which makes it easy to create new resource assertions.
type ResourceAssert struct {
	name             string
	id               string
	prefix           string
	assertions       []ResourceAssertion
	additionalPrefix string
}

// NewResourceAssert creates a ResourceAssert where the resource name should be used as a key for assertions.
func NewResourceAssert(name string, prefix string) *ResourceAssert {
	return &ResourceAssert{
		name:       name,
		prefix:     prefix,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewImportedResourceAssert creates a ResourceAssert where the resource id should be used as a key for assertions.
func NewImportedResourceAssert(id string, prefix string) *ResourceAssert {
	return &ResourceAssert{
		id:         id,
		prefix:     prefix,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewDatasourceAssert creates a ResourceAssert for data sources.
func NewDatasourceAssert(name string, prefix string, additionalPrefix string) *ResourceAssert {
	return &ResourceAssert{
		name:             name,
		prefix:           prefix,
		assertions:       make([]ResourceAssertion, 0),
		additionalPrefix: additionalPrefix,
	}
}

type resourceAssertionType string

const (
	resourceAssertionTypeValuePresent = "VALUE_PRESENT"
	resourceAssertionTypeValueSet     = "VALUE_SET"
	resourceAssertionTypeValueNotSet  = "VALUE_NOT_SET"
	resourceAssertionTypeSetElem      = "SET_ELEM"
)

type ResourceAssertion struct {
	fieldName             string
	expectedValue         string
	resourceAssertionType resourceAssertionType
}

func (r *ResourceAssert) AddAssertion(assertion ResourceAssertion) {
	assertion.fieldName = r.additionalPrefix + assertion.fieldName
	r.assertions = append(r.assertions, assertion)
}

func SetElem(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeSetElem}
}

func ValuePresent(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, resourceAssertionType: resourceAssertionTypeValuePresent}
}

func ValueSet(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ValueNotSet(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, resourceAssertionType: resourceAssertionTypeValueNotSet}
}

const showOutputPrefix = "show_output.0."

func ResourceShowOutputBoolValueSet(fieldName string, expected bool) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, strconv.FormatBool(expected))
}

func ResourceShowOutputIntValueSet(fieldName string, expected int) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, strconv.Itoa(expected))
}

func ResourceShowOutputFloatValueSet(fieldName string, expected float64) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, strconv.FormatFloat(expected, 'f', -1, 64))
}

func ResourceShowOutputStringUnderlyingValueSet[U ~string](fieldName string, expected U) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, string(expected))
}

func ResourceShowOutputValueSet(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: showOutputPrefix + fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ResourceShowOutputValueNotSet(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: showOutputPrefix + fieldName, resourceAssertionType: resourceAssertionTypeValueNotSet}
}

func ResourceShowOutputValuePresent(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: showOutputPrefix + fieldName, resourceAssertionType: resourceAssertionTypeValuePresent}
}

const (
	parametersPrefix      = "parameters.0."
	parametersValueSuffix = ".0.value"
	parametersLevelSuffix = ".0.level"
)

func ResourceParameterBoolValueSet[T ~string](parameterName T, expected bool) ResourceAssertion {
	return ResourceParameterValueSet(parameterName, strconv.FormatBool(expected))
}

func ResourceParameterIntValueSet[T ~string](parameterName T, expected int) ResourceAssertion {
	return ResourceParameterValueSet(parameterName, strconv.Itoa(expected))
}

func ResourceParameterStringUnderlyingValueSet[T ~string, U ~string](parameterName T, expected U) ResourceAssertion {
	return ResourceParameterValueSet(parameterName, string(expected))
}

func ResourceParameterValueSet[T ~string](parameterName T, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPrefix + strings.ToLower(string(parameterName)) + parametersValueSuffix, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ResourceParameterLevelSet[T ~string](parameterName T, parameterType sdk.ParameterType) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPrefix + strings.ToLower(string(parameterName)) + parametersLevelSuffix, expectedValue: string(parameterType), resourceAssertionType: resourceAssertionTypeValueSet}
}

// ToTerraformTestCheckFunc implements TestCheckFuncProvider to allow easier creation of new resource assertions.
// It goes through all the assertion accumulated earlier and gathers the results of the checks.
func (r *ResourceAssert) ToTerraformTestCheckFunc(t *testing.T, _ *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		var result []error

		for i, a := range r.assertions {
			switch a.resourceAssertionType {
			case resourceAssertionTypeSetElem:
				if err := resource.TestCheckTypeSetElemAttr(r.name, a.fieldName, a.expectedValue)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, r.prefix, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValueSet:
				if err := resource.TestCheckResourceAttr(r.name, a.fieldName, a.expectedValue)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, r.prefix, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValueNotSet:
				if err := resource.TestCheckNoResourceAttr(r.name, a.fieldName)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, r.prefix, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValuePresent:
				if err := resource.TestCheckResourceAttrSet(r.name, a.fieldName)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, r.prefix, i+1, len(r.assertions), errCut))
				}
			}
		}

		return errors.Join(result...)
	}
}

// ToTerraformImportStateCheckFunc implements ImportStateCheckFuncProvider to allow easier creation of new resource assertions.
// It goes through all the assertion accumulated earlier and gathers the results of the checks.
func (r *ResourceAssert) ToTerraformImportStateCheckFunc(t *testing.T, _ *helpers.TestClient) resource.ImportStateCheckFunc {
	t.Helper()
	return func(s []*terraform.InstanceState) error {
		var result []error

		for i, a := range r.assertions {
			switch a.resourceAssertionType {
			case resourceAssertionTypeValueSet:
				if err := importchecks.TestCheckResourceAttrInstanceState(r.id, a.fieldName, a.expectedValue)(s); err != nil {
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %w", r.id, r.prefix, i+1, len(r.assertions), err))
				}
			case resourceAssertionTypeValueNotSet:
				if err := importchecks.TestCheckResourceAttrNotInInstanceState(r.id, a.fieldName)(s); err != nil {
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %w", r.id, r.prefix, i+1, len(r.assertions), err))
				}
			case resourceAssertionTypeValuePresent:
				if err := importchecks.TestCheckResourceAttrInstanceStateSet(r.id, a.fieldName)(s); err != nil {
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %w", r.id, r.prefix, i+1, len(r.assertions), err))
				}
			}
		}

		return errors.Join(result...)
	}
}
