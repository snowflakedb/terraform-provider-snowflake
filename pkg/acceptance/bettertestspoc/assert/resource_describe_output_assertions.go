package assert

import (
	"strconv"
)

// TODO [SNOW-1501905]: this file duplicates resource_show_output_assertions.go file as a quick workaround for generated describe output assertions; it should be reworked with proper generation and deduplicated
const describeOutputPrefix = "describe_output.0."

func ResourceDescribeOutputBoolValueSet(fieldName string, expected bool) ResourceAssertion {
	return ResourceDescribeOutputValueSet(fieldName, strconv.FormatBool(expected))
}

func ResourceDescribeOutputBoolValueNotSet(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValueNotSet(fieldName)
}

func ResourceDescribeOutputBoolValuePresent(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValuePresent(fieldName)
}

func ResourceDescribeOutputIntValueSet(fieldName string, expected int) ResourceAssertion {
	return ResourceDescribeOutputValueSet(fieldName, strconv.Itoa(expected))
}

func ResourceDescribeOutputIntValueNotSet(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValueNotSet(fieldName)
}

func ResourceDescribeOutputIntValuePresent(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValuePresent(fieldName)
}

func ResourceDescribeOutputFloatValueSet(fieldName string, expected float64) ResourceAssertion {
	return ResourceDescribeOutputValueSet(fieldName, strconv.FormatFloat(expected, 'f', -1, 64))
}

func ResourceDescribeOutputFloatValueNotSet(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValueNotSet(fieldName)
}

func ResourceDescribeOutputFloatValuePresent(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValuePresent(fieldName)
}

func ResourceDescribeOutputStringUnderlyingValueSet[U ~string](fieldName string, expected U) ResourceAssertion {
	return ResourceDescribeOutputValueSet(fieldName, string(expected))
}

func ResourceDescribeOutputStringUnderlyingValueNotSet(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValueNotSet(fieldName)
}

func ResourceDescribeOutputStringUnderlyingValuePresent(fieldName string) ResourceAssertion {
	return ResourceDescribeOutputValuePresent(fieldName)
}

func ResourceDescribeOutputValueSet(fieldName string, expected string) ResourceAssertion {
	return ResourceAssertion{fieldName: describeOutputPrefix + fieldName, expectedValue: expected, resourceAssertionType: resourceAssertionTypeValueSet}
}

func ResourceDescribeOutputValueNotSet(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: describeOutputPrefix + fieldName, resourceAssertionType: resourceAssertionTypeValueNotSet}
}

func ResourceDescribeOutputValuePresent(fieldName string) ResourceAssertion {
	return ResourceAssertion{fieldName: describeOutputPrefix + fieldName, resourceAssertionType: resourceAssertionTypeValuePresent}
}

func ResourceDescribeOutputSetElem(fieldName string, expected string) ResourceAssertion {
	return SetElem(describeOutputPrefix+fieldName, expected)
}
