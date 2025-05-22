package sdk

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

var (
	ErrNilOptions                    = NewError("options cannot be nil")
	ErrPatternRequiredForLikeKeyword = NewError("pattern must be specified for like keyword")

	// re-importing from internal package
	ErrObjectNotFound = collections.ErrObjectNotFound

	// go-snowflake errors.
	ErrDoesNotExistOrOperationCannotBePerformed = NewError("object does not exist, or operation cannot be performed")
	ErrObjectNotExistOrAuthorized               = NewError("object does not exist or not authorized")
	ErrAccountIsEmpty                           = NewError("account is empty")
	ErrGrantPartiallyExecuted                   = NewError("grant partially executed")

	// snowflake-sdk errors.
	ErrInvalidObjectIdentifier = NewError("invalid object identifier")
	ErrDifferentDatabase       = NewError("database must be the same")
	ErrDifferentSchema         = NewError("schema must be the same")
)

type IntErrType string

const (
	IntErrEqual          IntErrType = "equal to"
	IntErrGreaterOrEqual IntErrType = "greater than or equal to"
	IntErrGreater        IntErrType = "greater than"
	IntErrLessOrEqual    IntErrType = "less than or equal to"
	IntErrLess           IntErrType = "less than"
)

func errIntValue(structName string, fieldName string, intErrType IntErrType, limit int) error {
	return newError(fmt.Sprintf("%s field: %s must be %s %d", structName, fieldName, string(intErrType), limit), 2)
}

func errIntBetween(structName string, fieldName string, from int, to int) error {
	return newError(fmt.Sprintf("%s field: %s must be between %d and %d", structName, fieldName, from, to), 2)
}

func errInvalidIdentifier(structName string, identifierField string) error {
	return newError(fmt.Sprintf("invalid object identifier of %s field: %s", structName, identifierField), 2)
}

func errOneOf(structName string, fieldNames ...string) error {
	return newError(fmt.Sprintf("%v fields: %v are incompatible and cannot be set at the same time", structName, fieldNames), 2)
}

func errNotSet(structName string, fieldNames ...string) error {
	return newError(fmt.Sprintf("%v fields: %v should be set", structName, fieldNames), 2)
}

func errSet(structName string, fieldNames ...string) error {
	return newError(fmt.Sprintf("%v fields: %v should not be set", structName, fieldNames), 2)
}

func errExactlyOneOf(structName string, fieldNames ...string) error {
	return newError(fmt.Sprintf("exactly one of %s fields %v must be set", structName, fieldNames), 2)
}

func errAtLeastOneOf(structName string, fieldNames ...string) error {
	return newError(fmt.Sprintf("at least one of %s fields %v must be set", structName, fieldNames), 2)
}

func errMoreThanOneOf(structName string, fieldNames ...string) error {
	return newError(fmt.Sprintf("more than one field (%v) of %s cannot be set", fieldNames, structName), 2)
}

func errInvalidValue(structName string, fieldName string, invalidValue string) error {
	return newError(fmt.Sprintf("invalid value %s of struct %s field: %s", invalidValue, structName, fieldName), 2)
}

func decodeDriverError(err error) error {
	if err == nil {
		return nil
	}
	log.Printf("[DEBUG] err: %v", err)
	m := map[string]error{
		"Object does not exist, or operation cannot be performed": ErrDoesNotExistOrOperationCannotBePerformed,
		"does not exist or not authorized":                        ErrObjectNotExistOrAuthorized,
		"account is empty":                                        ErrAccountIsEmpty,
		"Grant partially executed":                                ErrGrantPartiallyExecuted,
	}
	for k, v := range m {
		if strings.Contains(err.Error(), k) {
			return v
		}
	}

	return err
}

const errorIndentRune = '›'

var errorFileInfoRegexp = regexp.MustCompile(`\[\w+\.\w+:\d+\] `)

type Error struct {
	file         string
	line         int
	message      string
	nestedErrors []error
}

func (e *Error) Error() string {
	builder := new(strings.Builder)
	writeTree(e, builder, 0)
	return builder.String()
}

// NewError creates new sdk.Error with information like filename or line number (depending on where NewError was called)
func NewError(message string, nestedErrors ...error) error {
	return newError(message, 2, nestedErrors...)
}

// JoinErrors returns an error that wraps the given errors. Any nil error values are discarded.
// JoinErrors returns nil if errs contains no non-nil values, otherwise returns sdk.Error with errs as its nested errors
func JoinErrors(errs ...error) error {
	notNilErrs := make([]error, 0)
	for _, err := range errs {
		if err != nil {
			notNilErrs = append(notNilErrs, err)
		}
	}
	if len(notNilErrs) == 0 {
		return nil
	}
	return newError("joined error", 2, notNilErrs...)
}

// newError is a function that is supposed to be used by other sdk.Error constructors like NewError or JoinErrors.
// First of all, it returns error implementation which is against Golang conventions, but it's convenient to use
// in other constructors, because then there's no need for casting and guessing which type of error it is.
// The second reason is that there is a mysterious skip parameter which is only useful for other constructor functions.
// It determines how many function stack calls have to be skipped to get the right filename and line information,
// which is too low-level for normal use.
func newError(message string, skip int, nested ...error) *Error {
	line, filename := getCallerInfo(skip)
	return &Error{
		file:         filename,
		line:         line,
		message:      message,
		nestedErrors: nested,
	}
}

func getCallerInfo(skip int) (int, string) {
	_, file, line, _ := runtime.Caller(skip + 1)
	fileSplit := strings.Split(file, "/")
	var filename string
	if len(fileSplit) > 1 {
		filename = fileSplit[len(fileSplit)-1]
	} else {
		filename = fileSplit[0]
	}
	return line, filename
}

func writeTree(e error, builder *strings.Builder, indent int) {
	var sdkErr *Error
	if joinedErr, ok := e.(interface{ Unwrap() []error }); ok { //nolint:all
		errs := joinedErr.Unwrap()
		for i, err := range errs {
			if i > 0 {
				builder.WriteByte('\n')
			}
			writeTree(err, builder, indent)
		}
	} else if errors.As(e, &sdkErr) {
		builder.WriteString(strings.Repeat(fmt.Sprintf("%b ", errorIndentRune), indent) + fmt.Sprintf("[%s:%d] %s", sdkErr.file, sdkErr.line, sdkErr.message))
		for _, err := range sdkErr.nestedErrors {
			builder.WriteByte('\n')
			writeTree(err, builder, indent+2)
		}
	} else {
		builder.WriteString(strings.Repeat(fmt.Sprintf("%b ", errorIndentRune), indent) + e.Error())
	}
}
