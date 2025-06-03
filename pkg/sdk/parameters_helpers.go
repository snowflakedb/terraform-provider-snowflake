package sdk

import (
	"fmt"
	"strconv"
)

func parseBooleanParameter(parameter, value string) (_ *bool, err error) {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("boolean value (\"true\"/\"false\") expected for %v parameter, but got %v instead", parameter, value)
	}
	return &b, nil
}

func setBooleanValue[P AccountParameter | SessionParameter | UserParameter | TaskParameter | ObjectParameter](parameter P, value string, setField **bool) error {
	b, err := parseBooleanParameter(string(parameter), value)
	if err != nil {
		return err
	}
	*setField = b
	return nil
}

func setIntegerValue[P AccountParameter | SessionParameter | UserParameter | TaskParameter | ObjectParameter](parameter P, value string, setField **int) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("failed to parse parameter %s, expected an integer, but got %v", parameter, value)
	}
	*setField = Pointer(v)
	return nil
}
