package sdk

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteTree(t *testing.T) {
	fileInfoRegTemplate := func(filename string) string {
		return fmt.Sprintf("\\[%s:.*\\]", filename)
	}
	errorsTestFileInfoReg := fileInfoRegTemplate("errors_test.go")

	testCases := map[string]struct {
		Error         error
		Indent        int
		Expected      string
		MatchContains []string
	}{
		"basic error - no indent": {
			Error:    errors.New("some error"),
			Indent:   0,
			Expected: "some error",
		},
		"basic error - indent": {
			Error:    errors.New("some error"),
			Indent:   1,
			Expected: fmt.Sprintf("%b some error", errorIndentRune),
		},
		"basic error - double indent": {
			Error:    errors.New("some error"),
			Indent:   2,
			Expected: fmt.Sprintf("%b %b some error", errorIndentRune, errorIndentRune),
		},
		"joined error - no indent": {
			Error:    errors.Join(errors.New("err one"), errors.New("err two")),
			Indent:   0,
			Expected: "err one\nerr two",
		},
		"joined error - indent": {
			Error:    errors.Join(errors.New("err one"), errors.New("err two")),
			Indent:   1,
			Expected: fmt.Sprintf("%b err one\n%b err two", errorIndentRune, errorIndentRune),
		},
		"joined error - double indent": {
			Error:    errors.Join(errors.New("err one"), errors.New("err two")),
			Indent:   2,
			Expected: fmt.Sprintf("%b %b err one\n%b %b err two", errorIndentRune, errorIndentRune, errorIndentRune, errorIndentRune),
		},
		"custom error - no indent": {
			Error:         NewError("some error"),
			Indent:        0,
			MatchContains: []string{fmt.Sprintf("%s some error", errorsTestFileInfoReg)},
		},
		"custom error - indent": {
			Error:  errors.Join(NewError("err one"), NewError("err two")),
			Indent: 1,
			MatchContains: []string{
				fmt.Sprintf("%b %s err one", errorIndentRune, errorsTestFileInfoReg),
				fmt.Sprintf("%b %s err two", errorIndentRune, errorsTestFileInfoReg),
			},
		},
		"custom error - double indent": {
			Error:  errors.Join(NewError("err one"), NewError("err two")),
			Indent: 2,
			MatchContains: []string{
				fmt.Sprintf("%b %b %s err one", errorIndentRune, errorIndentRune, errorsTestFileInfoReg),
				fmt.Sprintf("%b %b %s err two", errorIndentRune, errorIndentRune, errorsTestFileInfoReg),
			},
		},
		"nested errors - custom errors combined with std errors": {
			Error: NewError("root error",
				errors.New("regular error"),
				errors.Join(
					errors.New("regular nested error"),
					NewError("custom nested error"),
					JoinErrors(
						errors.New("regular nested nested error"),
						NewError("custom nested nested error"),
					),
				),
				NewError("custom error"),
			),
			Indent: 0,
			MatchContains: []string{
				fmt.Sprintf("%s root error", errorsTestFileInfoReg),
				fmt.Sprintf("%b %b regular error", errorIndentRune, errorIndentRune),
				// Nested errors (errors.Join-ed) are on the same level, because there's no root error there
				// we could make another indent here by introducing a root error for every errors.Join-ed error
				fmt.Sprintf("%b %b regular nested error", errorIndentRune, errorIndentRune),
				fmt.Sprintf("%b %b %s custom nested error", errorIndentRune, errorIndentRune, errorsTestFileInfoReg),
				// Here, for example, I've made a root error inside JoinErrors function
				fmt.Sprintf("%b %b %s joined error", errorIndentRune, errorIndentRune, errorsTestFileInfoReg),
				fmt.Sprintf("%b %b %b %b regular nested nested error", errorIndentRune, errorIndentRune, errorIndentRune, errorIndentRune),
				fmt.Sprintf("%b %b %b %b %s custom nested nested error", errorIndentRune, errorIndentRune, errorIndentRune, errorIndentRune, errorsTestFileInfoReg),
				fmt.Sprintf("%b %b %s custom error", errorIndentRune, errorIndentRune, errorsTestFileInfoReg),
			},
		},
		"custom error - predefined errors": {
			Error: errors.Join(
				ErrInvalidObjectIdentifier,
				errNotSet("Struct", "Field"),
			),
			Indent: 2,
			MatchContains: []string{
				// Predefined errors are pointing to the file where they're declared
				fmt.Sprintf("%s invalid object identifier", fileInfoRegTemplate("errors.go")),
				fmt.Sprintf("%s Struct fields: \\[Field\\] should be set", errorsTestFileInfoReg),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			builder := new(strings.Builder)
			writeTree(tc.Error, builder, tc.Indent)
			if len(tc.Expected) == 0 && len(tc.MatchContains) == 0 {
				t.Fatal("expected or contains should be specified on a test case")
			}
			if len(tc.Expected) > 0 {
				require.Equal(t, tc.Expected, builder.String())
			}
			if len(tc.MatchContains) > 0 {
				for _, regex := range tc.MatchContains {
					require.Regexpf(t, regex, builder.String(), "regex %s not in: %s", regex, builder.String())
				}
			}
		})
	}
}

func TestDecodeDriverError(t *testing.T) {
	unrecognizedErr := errors.New("some other error")
	errWithNotFound := errors.New("foo bar not found")
	tests := []struct {
		name  string
		input error
		want  error
	}{
		{
			name:  "nil error returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Match ErrDoesNotExistOrOperationCannotBePerformed",
			input: errors.New("Object does not exist, or operation cannot be performed"),
			want:  ErrDoesNotExistOrOperationCannotBePerformed,
		},
		{
			name:  "Match ErrObjectNotExistOrAuthorized",
			input: errors.New("does not exist or not authorized"),
			want:  ErrObjectNotExistOrAuthorized,
		},
		{
			name:  "Match ErrAccountIsEmpty",
			input: errors.New("account is empty"),
			want:  ErrAccountIsEmpty,
		},
		{
			name:  "Match ErrGrantPartiallyExecuted",
			input: errors.New("Grant partially executed"),
			want:  ErrGrantPartiallyExecuted,
		},
		{
			name:  "Match ErrPatNotFound",
			input: errors.New("Programmatic access token foo not found"),
			want:  ErrPatNotFound,
		},
		{
			name:  "Unmatched error returns original",
			input: unrecognizedErr,
			want:  unrecognizedErr,
		},
		{
			name:  "Unmatched error returns original",
			input: errWithNotFound,
			want:  errWithNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeDriverError(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
