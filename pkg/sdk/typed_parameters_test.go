package sdk

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTypedParameter_int(t *testing.T) {
	tests := []struct {
		name string
		raw  *Parameter
		want TypedParameter[int]
	}{
		{
			name: "parses value and default",
			raw: &Parameter{
				Key:         "DATA_RETENTION_TIME_IN_DAYS",
				Value:       "7",
				Default:     "1",
				Level:       ParameterTypeDatabase,
				Description: "retention",
			},
			want: TypedParameter[int]{
				Key:         "DATA_RETENTION_TIME_IN_DAYS",
				Value:       7,
				Default:     1,
				Level:       ParameterTypeDatabase,
				Description: "retention",
			},
		},
		{
			name: "nil raw is zero",
			raw:  nil,
			want: TypedParameter[int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTypedParameter(tt.raw, strconv.Atoi)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewTypedParameter_string(t *testing.T) {
	tests := []struct {
		name string
		raw  *Parameter
		want TypedParameter[string]
	}{
		{
			name: "empty strings are zero",
			raw:  &Parameter{Key: "QUERY_TAG"},
			want: TypedParameter[string]{Key: "QUERY_TAG"},
		},
		{
			name: "parses value",
			raw:  &Parameter{Key: "QUERY_TAG", Value: "job"},
			want: TypedParameter[string]{Key: "QUERY_TAG", Value: "job"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTypedParameter(tt.raw, identityParse)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFillTypedParameter(t *testing.T) {
	tests := []struct {
		name string
		raw  *Parameter
		want string
	}{
		{
			name: "sets parsed value",
			raw:  &Parameter{Key: "LOG_LEVEL", Value: "INFO"},
			want: "INFO",
		},
		{
			name: "nil raw is zero",
			raw:  nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target TypedParameter[string]
			require.NoError(t, fillTypedParameter(tt.raw, identityParse, &target))
			require.Equal(t, tt.want, target.Value)
		})
	}
}

func TestParametersByKey(t *testing.T) {
	logLevel := &Parameter{Key: "LOG_LEVEL", Value: "INFO"}
	queryTag := &Parameter{Key: "QUERY_TAG", Value: "job"}
	tests := []struct {
		name string
		in   []*Parameter
		key  string
		want *Parameter
	}{
		{name: "finds key", in: []*Parameter{logLevel, queryTag}, key: "LOG_LEVEL", want: logLevel},
		{name: "missing key", in: []*Parameter{logLevel}, key: "QUERY_TAG", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, parametersByKey(tt.in)[tt.key])
		})
	}
}
