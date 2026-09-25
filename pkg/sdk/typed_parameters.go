package sdk

import "fmt"

// TypedParameter is a single SHOW PARAMETERS row with its Value/Default parsed into the parameter's
// Go type T (e.g. int, bool, LogLevel, AccountObjectIdentifier), so callers need no casting. It is
// produced by the generated <Object>ParametersDetails accessors. The generic Parameter (string-valued)
// remains the raw form returned by ShowParameters.
type TypedParameter[T any] struct {
	Key         string
	Value       T
	Default     T
	Level       ParameterType
	Description string
}

// newTypedParameter parses a raw string-valued Parameter into a TypedParameter[T] using the supplied
// parser (the generated accessors pass the parser matching the parameter's kind). Empty Value/Default
// are left as the zero value rather than treated as parse errors (SHOW may report empty strings).
func newTypedParameter[T any](raw *Parameter, parse func(string) (T, error)) (TypedParameter[T], error) {
	if raw == nil {
		return TypedParameter[T]{}, nil
	}
	tp := TypedParameter[T]{Key: raw.Key, Level: raw.Level, Description: raw.Description}
	if raw.Value != "" {
		v, err := parse(raw.Value)
		if err != nil {
			return tp, err
		}
		tp.Value = v
	}
	if raw.Default != "" {
		d, err := parse(raw.Default)
		if err != nil {
			return tp, err
		}
		tp.Default = d
	}
	return tp, nil
}

// fillTypedParameter parses raw into *target using parse, returning only an error (nil on success).
// The generated ShowParametersDetails accessors call one per field and errors.Join the results, so a
// caller sees every parse failure at once rather than just the first.
func fillTypedParameter[T any](raw *Parameter, parse func(string) (T, error), target *TypedParameter[T]) error {
	tp, err := newTypedParameter(raw, parse)
	if err != nil {
		return fmt.Errorf("parsing parameter %q: %w", raw.Key, err)
	}
	*target = tp
	return nil
}

// identityParse is the parser for string-valued parameters.
func identityParse(s string) (string, error) { return s, nil }

// parametersByKey indexes a SHOW PARAMETERS result by parameter key (SQL name).
func parametersByKey(params []*Parameter) map[string]*Parameter {
	m := make(map[string]*Parameter, len(params))
	for _, p := range params {
		m[p.Key] = p
	}
	return m
}
