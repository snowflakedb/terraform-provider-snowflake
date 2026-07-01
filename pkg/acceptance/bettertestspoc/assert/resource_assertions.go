package assert

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	describeOutputPath       = "describe_output.0."
	describeOutputCollection = "describe_output.#"
	showOutputPath           = "show_output.0."
	showOutputCollection     = "show_output.#"
	parametersPath           = "parameters.0."
	parametersCollection     = "parameters.#"
)

var (
	_ TestCheckFuncProvider        = (*ResourceAssert)(nil)
	_ ImportStateCheckFuncProvider = (*ResourceAssert)(nil)
)

// ResourceAssert is an embeddable struct that should be used to construct new resource assertions (for resource, show output, parameters, etc.).
// It implements both TestCheckFuncProvider and ImportStateCheckFuncProvider which makes it easy to create new resource assertions.
type ResourceAssert struct {
	name       string
	id         string
	prefix     string
	assertions []ResourceAssertion

	assertionPath string
}

// NewResourceAssert creates a ResourceAssert where the resource name should be used as a key for assertions.
func NewResourceAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name:       name,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewResourceShowOutputAssert creates a ResourceAssert for show output assertions with the resource name as a key.
func NewResourceShowOutputAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name: name,
		assertions: []ResourceAssertion{
			ValueSetFullPath(showOutputCollection, "1"),
		},
		assertionPath: showOutputPath,
	}
}

// NewResourceDescribeOutputAssert creates a ResourceAssert for describe output assertions with the resource name as a key.
func NewResourceDescribeOutputAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name: name,
		assertions: []ResourceAssertion{
			ValueSetFullPath(describeOutputCollection, "1"),
		},
		assertionPath: describeOutputPath,
	}
}

// NewResourceDescribeOutputAssertAtRow creates a ResourceAssert for describe output assertions at a specific row index.
func NewResourceDescribeOutputAssertAtRow(name string, rowIndex int) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: fmt.Sprintf("describe_output.%d.", rowIndex),
	}
}

// NewResourceParametersAssert creates a ResourceAssert for parameters assertions with the resource name as a key.
func NewResourceParametersAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name: name,
		assertions: []ResourceAssertion{
			ValueSetFullPath(parametersCollection, "1"),
		},
		assertionPath: parametersPath,
	}
}

// NewImportedResourceAssert creates a ResourceAssert where the resource id should be used as a key for assertions.
func NewImportedResourceAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id:         id,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewImportedResourceShowOutputAssert creates a ResourceAssert for show output assertions with the resource id as a key.
func NewImportedResourceShowOutputAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id: id,
		assertions: []ResourceAssertion{
			ValueSetFullPath(showOutputCollection, "1"),
		},
		assertionPath: showOutputPath,
	}
}

// NewImportedResourceDescribeOutputAssert creates a ResourceAssert for describe output assertions with the resource id as a key.
func NewImportedResourceDescribeOutputAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id: id,
		assertions: []ResourceAssertion{
			ValueSetFullPath(describeOutputCollection, "1"),
		},
		assertionPath: describeOutputPath,
	}
}

// NewImportedResourceParametersAssert creates a ResourceAssert for parameters assertions with the resource id as a key.
func NewImportedResourceParametersAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id: id,
		assertions: []ResourceAssertion{
			ValueSetFullPath(parametersCollection, "1"),
		},
		assertionPath: parametersPath,
	}
}

// NewDatasourceShowOutputAssert creates a ResourceAssert for show output assertions on a datasource at the given index.
func NewDatasourceShowOutputAssert(name string, objectsPath string, idx int) *ResourceAssert {
	return &ResourceAssert{
		name: name,
		assertions: []ResourceAssertion{
			ValueSetFullPath(fmt.Sprintf("%s.%d.%s", objectsPath, idx, showOutputCollection), "1"),
		},
		assertionPath: fmt.Sprintf("%s.%d.%s", objectsPath, idx, showOutputPath),
	}
}

// NewDatasourceDescribeOutputAssert creates a ResourceAssert for describe output assertions on a datasource at the given index.
func NewDatasourceDescribeOutputAssert(name string, objectsPath string, idx int) *ResourceAssert {
	return &ResourceAssert{
		name: name,
		assertions: []ResourceAssertion{
			ValueSetFullPath(fmt.Sprintf("%s.%d.%s", objectsPath, idx, describeOutputCollection), "1"),
		},
		assertionPath: fmt.Sprintf("%s.%d.%s", objectsPath, idx, describeOutputPath),
	}
}

// NewDatasourceParametersAssert creates a ResourceAssert for parameters assertions on a datasource at the given index.
func NewDatasourceParametersAssert(name string, objectsPath string, idx int) *ResourceAssert {
	return &ResourceAssert{
		name: name,
		assertions: []ResourceAssertion{
			ValueSetFullPath(fmt.Sprintf("%s.%d.%s", objectsPath, idx, parametersCollection), "1"),
		},
		assertionPath: fmt.Sprintf("%s.%d.%s", objectsPath, idx, parametersPath),
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

	fullPath string
}

func (r *ResourceAssert) AddAssertion(assertion ResourceAssertion) {
	assertion.fullPath = r.assertionPath + assertion.fieldName
	r.assertions = append(r.assertions, assertion)
}

func SetElem(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName + ".*", expectedValue: expected, resourceAssertionType: resourceAssertionTypeSetElem}
}

func ValuePresent(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, resourceAssertionType: resourceAssertionTypeValuePresent}
}

func ValueSet(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ValueSetFullPath(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet, fullPath: fieldName}
}

func ValueNotSet(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: fieldName, resourceAssertionType: resourceAssertionTypeValueNotSet}
}

func (r *ResourceAssert) BoolValueSet(fieldName string, expected bool) {
	r.AddAssertion(ValueSet(fieldName, strconv.FormatBool(expected)))
}

func (r *ResourceAssert) IntValueSet(fieldName string, expected int) {
	r.AddAssertion(ValueSet(fieldName, strconv.Itoa(expected)))
}

func (r *ResourceAssert) FloatValueSet(fieldName string, expected float64) {
	r.AddAssertion(ValueSet(fieldName, strconv.FormatFloat(expected, 'f', -1, 64)))
}

func (r *ResourceAssert) StringValueSet(fieldName string, expected string) {
	r.AddAssertion(ValueSet(fieldName, expected))
}

func (r *ResourceAssert) ValueSet(fieldName string, expected string) {
	r.AddAssertion(ValueSet(fieldName, expected))
}

func (r *ResourceAssert) ValueNotSet(fieldName string) {
	r.AddAssertion(ValueNotSet(fieldName))
}

func (r *ResourceAssert) ValuePresent(fieldName string) {
	r.AddAssertion(ValuePresent(fieldName))
}

// TODO [SNOW-3113138]: do we want to generate assertions for the length only?
func (r *ResourceAssert) CollectionLength(fieldName string, expected int) {
	r.AddAssertion(ValueSet(fieldName+".#", strconv.Itoa(expected)))
}

func (r *ResourceAssert) SetContainsElem(fieldName string, expected string) {
	r.AddAssertion(SetElem(fieldName, expected))
}

func (r *ResourceAssert) ListContainsElem(fieldName string, index int, expected string) {
	r.AddAssertion(ValueSet(fmt.Sprintf("%s.%d", fieldName, index), expected))
}

func (r *ResourceAssert) SetContainsExactlyBoolValues(fieldName string, expectedValues ...bool) {
	r.SetContainsExactlyStringValues(fieldName, collections.Map(expectedValues, strconv.FormatBool)...)
}

func (r *ResourceAssert) SetContainsExactlyIntValues(fieldName string, expectedValues ...int) {
	r.SetContainsExactlyStringValues(fieldName, collections.Map(expectedValues, strconv.Itoa)...)
}

// TODO [SNOW-3113138]: extract common conversions
func (r *ResourceAssert) SetContainsExactlyFloatValues(fieldName string, expectedValues ...float64) {
	r.SetContainsExactlyStringValues(fieldName, collections.Map(expectedValues, func(v float64) string {
		return strconv.FormatFloat(v, 'f', -1, 64)
	})...)
}

func (r *ResourceAssert) SetContainsExactlyStringValues(fieldName string, expectedValues ...string) {
	r.CollectionLength(fieldName, len(expectedValues))
	for _, value := range expectedValues {
		r.SetContainsElem(fieldName, value)
	}
}

func (r *ResourceAssert) ListContainsExactlyBoolValuesInOrder(fieldName string, expectedValues ...bool) {
	r.ListContainsExactlyStringValuesInOrder(fieldName, collections.Map(expectedValues, strconv.FormatBool)...)
}

func (r *ResourceAssert) ListContainsExactlyIntValuesInOrder(fieldName string, expectedValues ...int) {
	r.ListContainsExactlyStringValuesInOrder(fieldName, collections.Map(expectedValues, strconv.Itoa)...)
}

// TODO [SNOW-3113138]: extract common conversions
func (r *ResourceAssert) ListContainsExactlyFloatValuesInOrder(fieldName string, expectedValues ...float64) {
	r.ListContainsExactlyStringValuesInOrder(fieldName, collections.Map(expectedValues, func(v float64) string {
		return strconv.FormatFloat(v, 'f', -1, 64)
	})...)
}

func (r *ResourceAssert) ListContainsExactlyStringValuesInOrder(fieldName string, expectedValues ...string) {
	r.CollectionLength(fieldName, len(expectedValues))
	for idx, value := range expectedValues {
		r.ListContainsElem(fieldName, idx, value)
	}
}

const (
	parametersValueSuffix       = ".0.value"
	parametersLevelSuffix       = ".0.level"
	parametersKeySuffix         = ".0.key"
	parametersDefaultSuffix     = ".0.default"
	parametersDescriptionSuffix = ".0.description"
)

func (r *ResourceAssert) ParameterValueSet(parameterName string, expected string) {
	r.ValueSet(strings.ToLower(parameterName)+parametersValueSuffix, expected)
}

func (r *ResourceAssert) ParameterBoolValueSet(parameterName string, expected bool) {
	r.ValueSet(strings.ToLower(parameterName)+parametersValueSuffix, strconv.FormatBool(expected))
}

func (r *ResourceAssert) ParameterIntValueSet(parameterName string, expected int) {
	r.ValueSet(strings.ToLower(parameterName)+parametersValueSuffix, strconv.Itoa(expected))
}

func (r *ResourceAssert) ParameterLevelSet(parameterName string, expected sdk.ParameterType) {
	r.ValueSet(strings.ToLower(parameterName)+parametersLevelSuffix, string(expected))
}

func (r *ResourceAssert) ParameterKeySet(parameterName string, expected string) {
	r.ValueSet(strings.ToLower(parameterName)+parametersKeySuffix, expected)
}

func (r *ResourceAssert) ParameterDefaultSet(parameterName string, expected string) {
	r.ValueSet(strings.ToLower(parameterName)+parametersDefaultSuffix, expected)
}

func (r *ResourceAssert) ParameterDescriptionSet(parameterName string, expected string) {
	r.ValueSet(strings.ToLower(parameterName)+parametersDescriptionSuffix, expected)
}

func (r *ResourceAssert) ParameterDescriptionPresent(parameterName string) {
	r.ValuePresent(strings.ToLower(parameterName) + parametersDescriptionSuffix)
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
				if err := resource.TestCheckTypeSetElemAttr(r.name, a.fullPath, a.expectedValue)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, a.fullPath, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValueSet:
				if err := resource.TestCheckResourceAttr(r.name, a.fullPath, a.expectedValue)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, a.fullPath, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValueNotSet:
				if err := resource.TestCheckNoResourceAttr(r.name, a.fullPath)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, a.fullPath, i+1, len(r.assertions), errCut))
				}
			case resourceAssertionTypeValuePresent:
				if err := resource.TestCheckResourceAttrSet(r.name, a.fullPath)(s); err != nil {
					errCut, _ := strings.CutPrefix(err.Error(), fmt.Sprintf("%s: ", r.name))
					result = append(result, fmt.Errorf("%s %s assertion [%d/%d]: failed with error: %s", r.name, a.fullPath, i+1, len(r.assertions), errCut))
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
				if err := importchecks.TestCheckResourceAttrInstanceState(r.id, a.fullPath, a.expectedValue)(s); err != nil {
					result = append(result, fmt.Errorf("imported %s assertion (path: %s) [%d/%d]: failed with error: %w", r.id, a.fullPath, i+1, len(r.assertions), err))
				}
			case resourceAssertionTypeValueNotSet:
				if err := importchecks.TestCheckResourceAttrNotInInstanceState(r.id, a.fullPath)(s); err != nil {
					result = append(result, fmt.Errorf("imported %s assertion (path: %s) [%d/%d]: failed with error: %w", r.id, a.fullPath, i+1, len(r.assertions), err))
				}
			case resourceAssertionTypeValuePresent:
				if err := importchecks.TestCheckResourceAttrInstanceStateSet(r.id, a.fullPath)(s); err != nil {
					result = append(result, fmt.Errorf("imported %s assertion (path: %s) [%d/%d]: failed with error: %w", r.id, a.fullPath, i+1, len(r.assertions), err))
				}
			}
		}

		return errors.Join(result...)
	}
}
