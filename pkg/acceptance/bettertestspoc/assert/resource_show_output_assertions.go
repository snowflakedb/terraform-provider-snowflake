package assert

import (
	"strconv"
)

const showOutputPrefix = "show_output.0."

func ResourceShowOutputBoolValueSet(fieldName string, expected bool) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, strconv.FormatBool(expected))
}

func ResourceShowOutputBoolValueNotSet(fieldName string) ResourceAssertion {
	return ResourceShowOutputValueNotSet(fieldName)
}

func ResourceShowOutputBoolValuePresent(fieldName string) ResourceAssertion {
	return ResourceShowOutputValuePresent(fieldName)
}

func ResourceShowOutputIntValueSet(fieldName string, expected int) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, strconv.Itoa(expected))
}

func ResourceShowOutputIntValueNotSet(fieldName string) ResourceAssertion {
	return ResourceShowOutputValueNotSet(fieldName)
}

func ResourceShowOutputIntValuePresent(fieldName string) ResourceAssertion {
	return ResourceShowOutputValuePresent(fieldName)
}

func ResourceShowOutputFloatValueSet(fieldName string, expected float64) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, strconv.FormatFloat(expected, 'f', -1, 64))
}

func ResourceShowOutputFloatValueNotSet(fieldName string) ResourceAssertion {
	return ResourceShowOutputValueNotSet(fieldName)
}

func ResourceShowOutputFloatValuePresent(fieldName string) ResourceAssertion {
	return ResourceShowOutputValuePresent(fieldName)
}

func ResourceShowOutputStringUnderlyingValueSet[U ~string](fieldName string, expected U) ResourceAssertion {
	return ResourceShowOutputValueSet(fieldName, string(expected))
}

func ResourceShowOutputStringUnderlyingValueNotSet(fieldName string) ResourceAssertion {
	return ResourceShowOutputValueNotSet(fieldName)
}

func ResourceShowOutputStringUnderlyingValuePresent(fieldName string) ResourceAssertion {
	return ResourceShowOutputValuePresent(fieldName)
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
