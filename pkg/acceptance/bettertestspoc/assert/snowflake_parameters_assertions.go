package assert

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

type (
	ParametersProvider[I sdk.ObjectIdentifier]           func(*testing.T, I) []*sdk.Parameter
	testClientParametersProvider[I sdk.ObjectIdentifier] func(client *helpers.TestClient) ParametersProvider[I]
)

// SnowflakeParametersAssert is an embeddable struct that should be used to construct new Snowflake parameters assertions.
// It implements both TestCheckFuncProvider and ImportStateCheckFuncProvider which makes it easy to create new resource assertions.
type SnowflakeParametersAssert[I sdk.ObjectIdentifier] struct {
	assertions                   []SnowflakeParameterAssertion
	id                           I
	objectType                   sdk.ObjectType
	testClientParametersProvider testClientParametersProvider[I]
	parameters                   []*sdk.Parameter
}

type snowflakeParameterAssertionType int

const (
	snowflakeParameterAssertionTypeExpectedValue snowflakeParameterAssertionType = iota
	snowflakeParameterAssertionTypeDefaultValue
	snowflakeParameterAssertionTypeDefaultValueOnLevel
	snowflakeParameterAssertionTypeLevel
)

type SnowflakeParameterAssertion struct {
	parameterName string
	expectedValue string
	parameterType sdk.ParameterType
	assertionType snowflakeParameterAssertionType
}

// NewSnowflakeParametersAssertWithTestClientParametersProvider creates a SnowflakeParametersAssert with id and the test client-varying parameters provider.
// Object to check is lazily fetched from Snowflake when the checks are being run.
func NewSnowflakeParametersAssertWithTestClientParametersProvider[I sdk.ObjectIdentifier](id I, objectType sdk.ObjectType, testClientParametersProvider testClientParametersProvider[I]) *SnowflakeParametersAssert[I] {
	return &SnowflakeParametersAssert[I]{
		assertions:                   make([]SnowflakeParameterAssertion, 0),
		id:                           id,
		objectType:                   objectType,
		testClientParametersProvider: testClientParametersProvider,
	}
}

// NewSnowflakeParametersAssertWithParameters creates a SnowflakeParametersAssert with parameters already fetched from Snowflake.
// All the checks are run against the given set of parameters.
func NewSnowflakeParametersAssertWithParameters[I sdk.ObjectIdentifier](id I, objectType sdk.ObjectType, parameters []*sdk.Parameter) *SnowflakeParametersAssert[I] {
	return &SnowflakeParametersAssert[I]{
		assertions: make([]SnowflakeParameterAssertion, 0),
		id:         id,
		objectType: objectType,
		parameters: parameters,
	}
}

func (s *SnowflakeParametersAssert[I]) AddAssertion(assertion SnowflakeParameterAssertion) {
	s.assertions = append(s.assertions, assertion)
}

func SnowflakeParameterBoolValueSet[T ~string](parameterName T, expected bool) SnowflakeParameterAssertion {
	return SnowflakeParameterValueSet(parameterName, strconv.FormatBool(expected))
}

func SnowflakeParameterIntValueSet[T ~string](parameterName T, expected int) SnowflakeParameterAssertion {
	return SnowflakeParameterValueSet(parameterName, strconv.Itoa(expected))
}

func SnowflakeParameterStringUnderlyingValueSet[T ~string, U ~string](parameterName T, expected U) SnowflakeParameterAssertion {
	return SnowflakeParameterValueSet(parameterName, string(expected))
}

func SnowflakeParameterValueSet[T ~string](parameterName T, expected string) SnowflakeParameterAssertion {
	return SnowflakeParameterAssertion{parameterName: string(parameterName), expectedValue: expected}
}

func SnowflakeParameterDefaultValueSet[T ~string](parameterName T) SnowflakeParameterAssertion {
	return SnowflakeParameterAssertion{parameterName: string(parameterName), assertionType: snowflakeParameterAssertionTypeDefaultValue}
}

func SnowflakeParameterDefaultValueOnLevelSet[T ~string](parameterName T, parameterType sdk.ParameterType) SnowflakeParameterAssertion {
	return SnowflakeParameterAssertion{parameterName: string(parameterName), parameterType: parameterType, assertionType: snowflakeParameterAssertionTypeDefaultValueOnLevel}
}

func SnowflakeParameterLevelSet[T ~string](parameterName T, parameterType sdk.ParameterType) SnowflakeParameterAssertion {
	return SnowflakeParameterAssertion{parameterName: string(parameterName), parameterType: parameterType, assertionType: snowflakeParameterAssertionTypeLevel}
}

// ToTerraformTestCheckFunc implements TestCheckFuncProvider to allow easier creation of new Snowflake object parameters assertions.
// It goes through all the assertion accumulated earlier and gathers the results of the checks.
func (s *SnowflakeParametersAssert[_]) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		return s.runSnowflakeParametersAssertions(t, testClient)
	}
}

// ToTerraformImportStateCheckFunc implements ImportStateCheckFuncProvider to allow easier creation of new Snowflake object parameters assertions.
// It goes through all the assertion accumulated earlier and gathers the results of the checks.
func (s *SnowflakeParametersAssert[_]) ToTerraformImportStateCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.ImportStateCheckFunc {
	t.Helper()
	return func(_ []*terraform.InstanceState) error {
		return s.runSnowflakeParametersAssertions(t, testClient)
	}
}

// VerifyAll implements InPlaceAssertionVerifier to allow easier creation of new Snowflake parameters assertions.
// It verifies all the assertions accumulated earlier and gathers the results of the checks.
func (s *SnowflakeParametersAssert[_]) VerifyAll(t *testing.T, testClient *helpers.TestClient) {
	t.Helper()
	err := s.runSnowflakeParametersAssertions(t, testClient)
	require.NoError(t, err)
}

func (s *SnowflakeParametersAssert[_]) runSnowflakeParametersAssertions(t *testing.T, testClient *helpers.TestClient) error {
	t.Helper()

	var parameters []*sdk.Parameter
	switch {
	case s.parameters != nil:
		parameters = s.parameters
	case s.testClientParametersProvider != nil:
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		parameters = s.testClientParametersProvider(testClient)(t, s.id)
	default:
		return fmt.Errorf("cannot proceed with parameters assertion for object %s[%s]: parameters or parameters provider must be specified", s.objectType, s.id.FullyQualifiedName())
	}

	var result []error

	for i, assertion := range s.assertions {
		switch assertion.assertionType {
		case snowflakeParameterAssertionTypeExpectedValue:
			if v := helpers.FindParameter(t, parameters, assertion.parameterName).Value; assertion.expectedValue != v {
				result = append(result, fmt.Errorf(
					"parameter assertion for %s[%s][%s][%d/%d] failed: expected value %s, got %s",
					s.objectType, s.id.FullyQualifiedName(), assertion.parameterName, i+1, len(s.assertions), assertion.expectedValue, v,
				))
			}
		case snowflakeParameterAssertionTypeDefaultValue:
			if p := helpers.FindParameter(t, parameters, assertion.parameterName); p.Default != p.Value {
				result = append(result, fmt.Errorf(
					"parameter assertion for %s[%s][%s][%d/%d] failed: expected default value %s, got %s",
					s.objectType, s.id.FullyQualifiedName(), assertion.parameterName, i+1, len(s.assertions), p.Default, p.Value,
				))
			}
		case snowflakeParameterAssertionTypeDefaultValueOnLevel:
			if p := helpers.FindParameter(t, parameters, assertion.parameterName); p.Default != p.Value || p.Level != assertion.parameterType {
				result = append(result, fmt.Errorf(
					"parameter assertion for %s[%s][%s][%d/%d] failed: expected default value %s on level %s, got %s and level %s",
					s.objectType, s.id.FullyQualifiedName(), assertion.parameterName, i+1, len(s.assertions), p.Default, assertion.parameterType, p.Value, p.Level,
				))
			}
		case snowflakeParameterAssertionTypeLevel:
			if p := helpers.FindParameter(t, parameters, assertion.parameterName); p.Level != assertion.parameterType {
				result = append(result, fmt.Errorf(
					"parameter assertion for %s[%s][%s][%d/%d] failed: expected level %s, got %s",
					s.objectType, s.id.FullyQualifiedName(), assertion.parameterName, i+1, len(s.assertions), assertion.parameterType, p.Level,
				))
			}
		default:
			return fmt.Errorf("cannot proceed with parameters assertion for object %s[%s]: assertion type must be specified", s.objectType, s.id.FullyQualifiedName())
		}
	}

	return errors.Join(result...)
}
