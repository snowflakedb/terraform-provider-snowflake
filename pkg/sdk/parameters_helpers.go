package sdk

import (
	"fmt"
	"strconv"
)

// ClientTelemetryEnableSessionParameter is an undocumented session parameter, so we are keeping it out SessionParameter list for now.
const ClientTelemetryEnableSessionParameter = "CLIENT_TELEMETRY_ENABLED"

func parseBooleanParameter(parameter, value string) (_ *bool, err error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("boolean value (\"true\"/\"false\") expected for %v parameter, but got %v instead", parameter, value)
	}
	return &b, nil
}

type ParameterConstraint interface {
	AccountParameter | SessionParameter | UserParameter | TaskParameter | ObjectParameter
}

func setBooleanValue[P ParameterConstraint](parameter P, value string, setField **bool) error {
	b, err := parseBooleanParameter(string(parameter), value)
	if err != nil {
		return err
	}
	*setField = b
	return nil
}

func setIntegerValue[P ParameterConstraint](parameter P, value string, setField **int) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("failed to parse parameter %s, expected an integer, but got %v", parameter, value)
	}
	*setField = Pointer(v)
	return nil
}

// Render returns the SQL token for the value of a parameter that the parameter catalog does not support. It matches the SQL
// of the catalog write path for the same type. STRING is also the kind of enum parameters, because the catalog renders them as string literals.
func (k ParameterValueKind) Render(parameter string, value string) (string, error) {
	switch k {
	case ParameterValueKindString:
		return SingleQuotes.Modify(value), nil
	case ParameterValueKindIdentifier:
		id, err := ParseObjectIdentifierString(value)
		if err != nil {
			return "", fmt.Errorf("invalid identifier value for %s parameter: %w", parameter, err)
		}
		return id.FullyQualifiedName(), nil
	case ParameterValueKindInteger:
		i, err := strconv.Atoi(value)
		if err != nil {
			return "", fmt.Errorf("integer value expected for %s parameter, but got %v instead", parameter, value)
		}
		return strconv.Itoa(i), nil
	case ParameterValueKindFloat:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", fmt.Errorf("float value expected for %s parameter, but got %v instead", parameter, value)
		}
		return formatFloatLiteral(f), nil
	case ParameterValueKindBoolean:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "", fmt.Errorf("boolean value (\"true\"/\"false\") expected for %s parameter, but got %v instead", parameter, value)
		}
		return strconv.FormatBool(b), nil
	default:
		return "", fmt.Errorf("unknown parameter value kind %q", string(k))
	}
}
