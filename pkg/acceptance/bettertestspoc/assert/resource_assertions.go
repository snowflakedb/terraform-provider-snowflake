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
	describeOutputPath = "describe_output.0."
	showOutputPath     = "show_output.0."
	parametersPath     = "parameters.0."
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

	assertionPath string
}

// NewResourceAssert creates a ResourceAssert where the resource name should be used as a key for assertions.
func NewResourceAssert(name string, prefix string) *ResourceAssert {
	return &ResourceAssert{
		name:       name,
		prefix:     prefix,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewResourceAssertTmp creates a ResourceAssert where the resource name should be used as a key for assertions.
// TODO [next PRs]: rename to NewResourceAssert, remove the old NewResourceAssert when all objects are migrated
func NewResourceAssertTmp(name string) *ResourceAssert {
	return &ResourceAssert{
		name:       name,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewResourceShowOutputAssert creates a ResourceAssert for show output assertions with the resource name as a key.
func NewResourceShowOutputAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: showOutputPath,
	}
}

// NewResourceDescribeOutputAssert creates a ResourceAssert for describe output assertions with the resource name as a key.
func NewResourceDescribeOutputAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: describeOutputPath,
	}
}

// NewResourceParametersAssert creates a ResourceAssert for parameters assertions with the resource name as a key.
func NewResourceParametersAssert(name string) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: parametersPath,
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

// NewImportedResourceAssertTmp creates a ResourceAssert where the resource id should be used as a key for assertions.
// TODO [next PR]: rename to NewImportedResourceAssert, remove the old NewImportedResourceAssert when all objects are migrated
func NewImportedResourceAssertTmp(id string) *ResourceAssert {
	return &ResourceAssert{
		id:         id,
		assertions: make([]ResourceAssertion, 0),
	}
}

// NewImportedResourceShowOutputAssert creates a ResourceAssert for show output assertions with the resource id as a key.
func NewImportedResourceShowOutputAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id:            id,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: showOutputPath,
	}
}

// NewImportedResourceDescribeOutputAssert creates a ResourceAssert for describe output assertions with the resource id as a key.
func NewImportedResourceDescribeOutputAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id:            id,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: describeOutputPath,
	}
}

// NewImportedResourceParametersAssert creates a ResourceAssert for parameters assertions with the resource id as a key.
func NewImportedResourceParametersAssert(id string) *ResourceAssert {
	return &ResourceAssert{
		id:            id,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: parametersPath,
	}
}

// NewDatasourceAssert creates a ResourceAssert for data sources.
// TODO [next PRs]: remove this method entirely when all invocations replaced with NewDatasourceShowOutputAssert, NewDatasourceDescribeOutputAssert, and NewDatasourceParametersAssert
func NewDatasourceAssert(name string, prefix string, additionalPrefix string) *ResourceAssert {
	return &ResourceAssert{
		name:             name,
		prefix:           prefix,
		assertions:       make([]ResourceAssertion, 0),
		additionalPrefix: additionalPrefix,
	}
}

// NewDatasourceShowOutputAssert creates a ResourceAssert for show output assertions on a datasource at the given index.
func NewDatasourceShowOutputAssert(name string, objectsPath string, idx int) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: fmt.Sprintf("%s.%d.%s", objectsPath, idx, showOutputPath),
	}
}

// NewDatasourceDescribeOutputAssert creates a ResourceAssert for describe output assertions on a datasource at the given index.
func NewDatasourceDescribeOutputAssert(name string, objectsPath string, idx int) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
		assertionPath: fmt.Sprintf("%s.%d.%s", objectsPath, idx, describeOutputPath),
	}
}

// NewDatasourceParametersAssert creates a ResourceAssert for parameters assertions on a datasource at the given index.
func NewDatasourceParametersAssert(name string, objectsPath string, idx int) *ResourceAssert {
	return &ResourceAssert{
		name:          name,
		assertions:    make([]ResourceAssertion, 0),
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
	// TODO [next PRs]: remove additionalPrefix logic when all the objects are migrated
	assertion.fieldName = r.additionalPrefix + assertion.fieldName
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
	return ResourceAssertion{fieldName: parametersPath + strings.ToLower(string(parameterName)) + parametersValueSuffix, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ResourceParameterLevelSet[T ~string](parameterName T, parameterType sdk.ParameterType) ResourceAssertion {
	return ResourceAssertion{fieldName: parametersPath + strings.ToLower(string(parameterName)) + parametersLevelSuffix, expectedValue: string(parameterType), resourceAssertionType: resourceAssertionTypeValueSet}
}

func ResourceParameterKeySet[T ~string](parameterName T, expected string) ResourceAssertion {
	return ValueSet(parametersPath+strings.ToLower(string(parameterName))+parametersKeySuffix, expected)
}

func ResourceParameterDefaultSet[T ~string](parameterName T, expected string) ResourceAssertion {
	return ValueSet(parametersPath+strings.ToLower(string(parameterName))+parametersDefaultSuffix, expected)
}

func ResourceParameterDescriptionSet[T ~string](parameterName T, expected string) ResourceAssertion {
	return ValueSet(parametersPath+strings.ToLower(string(parameterName))+parametersDescriptionSuffix, expected)
}

func ResourceParameterDescriptionPresent[T ~string](parameterName T) ResourceAssertion {
	return ValuePresent(parametersPath + strings.ToLower(string(parameterName)) + parametersDescriptionSuffix)
}

func (r *ResourceAssert) ParameterValueSet(parameterName string, expected string) {
	r.AddAssertion(ValueSet(strings.ToLower(parameterName)+parametersValueSuffix, expected))
}

func (r *ResourceAssert) ParameterBoolValueSet(parameterName string, expected bool) {
	r.AddAssertion(ValueSet(strings.ToLower(parameterName)+parametersValueSuffix, strconv.FormatBool(expected)))
}

func (r *ResourceAssert) ParameterIntValueSet(parameterName string, expected int) {
	r.AddAssertion(ValueSet(strings.ToLower(parameterName)+parametersValueSuffix, strconv.Itoa(expected)))
}

func (r *ResourceAssert) ParameterLevelSet(parameterName string, expected sdk.ParameterType) {
	r.AddAssertion(ValueSet(strings.ToLower(parameterName)+parametersLevelSuffix, string(expected)))
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
