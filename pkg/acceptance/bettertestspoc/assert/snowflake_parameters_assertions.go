package assert

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type (
	parametersProvider[I sdk.ObjectIdentifier] func(*testing.T, I) []*sdk.Parameter
)

// SnowflakeParametersAssert is an embeddable struct that should be used to construct new Snowflake parameters assertions.
// It implements both TestCheckFuncProvider and ImportStateCheckFuncProvider which makes it easy to create new resource assertions.
type SnowflakeParametersAssert[I sdk.ObjectIdentifier] struct {
	assertions []snowflakeParameterAssertion
	id         I
	objectType sdk.ObjectType
	provider   parametersProvider[I]
	parameters []*sdk.Parameter
}

type snowflakeParameterAssertion struct {
	parameterName string
	expectedValue string
}

// NewSnowflakeParametersAssertWithProvider creates a SnowflakeParametersAssert with id and the provider.
// Object to check is lazily fetched from Snowflake when the checks are being run.
func NewSnowflakeParametersAssertWithProvider[I sdk.ObjectIdentifier](id I, objectType sdk.ObjectType, provider parametersProvider[I]) *SnowflakeParametersAssert[I] {
	return &SnowflakeParametersAssert[I]{
		assertions: make([]snowflakeParameterAssertion, 0),
		id:         id,
		objectType: objectType,
		provider:   provider,
	}
}

// NewSnowflakeParametersAssertWithParameters creates a SnowflakeParametersAssert with parameters already fetched from Snowflake.
// All the checks are run against the given set of parameters.
func NewSnowflakeParametersAssertWithParameters[I sdk.ObjectIdentifier](id I, objectType sdk.ObjectType, parameters []*sdk.Parameter) *SnowflakeParametersAssert[I] {
	return &SnowflakeParametersAssert[I]{
		assertions: make([]snowflakeParameterAssertion, 0),
		id:         id,
		objectType: objectType,
		parameters: parameters,
	}
}

func snowflakeParameterBoolValueSet[T ~string](parameterName T, expected bool) snowflakeParameterAssertion {
	return snowflakeParameterValueSet(parameterName, strconv.FormatBool(expected))
}

func snowflakeParameterIntValueSet[T ~string](parameterName T, expected int) snowflakeParameterAssertion {
	return snowflakeParameterValueSet(parameterName, strconv.Itoa(expected))
}

func snowflakeParameterStringUnderlyingValueSet[T ~string, U ~string](parameterName T, expected U) snowflakeParameterAssertion {
	return snowflakeParameterValueSet(parameterName, string(expected))
}

func snowflakeParameterValueSet[T ~string](parameterName T, expected string) snowflakeParameterAssertion {
	return snowflakeParameterAssertion{parameterName: string(parameterName), expectedValue: expected}
}

// TODO: can we just replace all above with this one?
func snowflakeParameterValueSetGeneric[T ~string, U bool | int | ~string](parameterName T, expected U) snowflakeParameterAssertion {
	return snowflakeParameterAssertion{parameterName: string(parameterName), expectedValue: fmt.Sprintf("%s", expected)}
}

// VerifyAll implements InPlaceAssertionVerifier to allow easier creation of new Snowflake parameters assertions.
// It verifies all the assertions accumulated earlier and gathers the results of the checks.
func (s *SnowflakeParametersAssert[_]) VerifyAll(t *testing.T) {
	t.Helper()
	err := s.runSnowflakeParametersAssertions(t)
	require.NoError(t, err)
}

func (s *SnowflakeParametersAssert[_]) runSnowflakeParametersAssertions(t *testing.T) error {
	t.Helper()

	var parameters []*sdk.Parameter
	switch {
	case s.provider != nil:
		parameters = s.provider(t, s.id)
	case s.parameters != nil:
		parameters = s.parameters
	default:
		return fmt.Errorf("cannot proceed with parameters assertion for object %s[%s]: parameters or parameters provider must be specified", s.objectType, s.id.FullyQualifiedName())
	}

	var result []error

	for i, assertion := range s.assertions {
		if v := helpers.FindParameter(t, parameters, assertion.parameterName).Value; assertion.expectedValue != v {
			result = append(result, fmt.Errorf("parameter assertion for %s[%s][%s][%d/%d] failed: expected value %s, got %s", s.objectType, s.id.FullyQualifiedName(), assertion.parameterName, i+1, len(s.assertions), assertion.expectedValue, v))
		}
	}

	return errors.Join(result...)
}
