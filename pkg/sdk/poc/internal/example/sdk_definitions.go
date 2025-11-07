package example

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// This file's only purpose is to make generated objects compile (or close to compile).
// Later code will be generated inside sdk package, so the objects will be accessible there.

type optionsProvider[T any] interface {
	toOpts() *T
}

type validatable interface {
	validate() error
}

type Client struct{}

type (
	ObjectIdentifier         interface{}
	AccountObjectIdentifier  struct{}
	DatabaseObjectIdentifier struct{}
	ExternalObjectIdentifier struct{}
	SchemaObjectIdentifier   struct{}
	TableColumnIdentifier    struct{}
)

func RandomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	_ = t
	return AccountObjectIdentifier{}
}

func RandomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	_ = t
	return DatabaseObjectIdentifier{}
}

func RandomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	_ = t
	return SchemaObjectIdentifier{}
}

func ValidObjectIdentifier(objectIdentifier ObjectIdentifier) bool {
	_ = objectIdentifier
	return true
}

func valueSet(value interface{}) bool {
	_ = value
	return true
}

func anyValueSet(values ...interface{}) bool {
	_ = values
	return true
}

func everyValueSet(values ...interface{}) bool {
	_ = values
	return true
}

func exactlyOneValueSet(values ...interface{}) bool {
	_ = values
	return true
}

func errOneOf(fieldNames ...string) error {
	return fmt.Errorf("fields %v are incompatible and cannot be set at once", fieldNames)
}

func errExactlyOneOf(fieldNames ...string) error {
	return fmt.Errorf("exactly one of %v must be set", fieldNames)
}

func errAtLeastOneOf(fieldNames ...string) error {
	return fmt.Errorf("at least one of %v must be set", fieldNames)
}

var (
	ErrNilOptions              = errors.New("options cannot be nil")
	ErrInvalidObjectIdentifier = errors.New("invalid object identifier")
)

func validateAndExec(client *Client, ctx context.Context, opts validatable) error {
	_, _, _ = client, ctx, opts
	return nil
}

func validateAndQuery[T any](client *Client, ctx context.Context, opts validatable) (*[]T, error) {
	_, _, _ = client, ctx, opts
	return nil, nil
}

func validateAndQueryOne[T any](client *Client, ctx context.Context, opts validatable) (*T, error) {
	_, _, _ = client, ctx, opts
	return nil, nil
}

func assertOptsInvalid(t *testing.T, opts validatable, expectedError error) {
	t.Helper()
	_ = t
	_ = opts
	_ = expectedError
}

func assertOptsInvalidJoinedErrors(t *testing.T, opts validatable, expectedErrors ...error) {
	t.Helper()
	_ = t
	_ = opts
	_ = expectedErrors
}

func assertOptsValid(t *testing.T, opts validatable) {
	t.Helper()
	_ = t
	_ = opts
}

func assertSQLEquals(t *testing.T, opts any, format string, args ...any) {
	t.Helper()
	_ = t
	_ = opts
	_ = format
	_ = args
}

func assertOptsValidAndSQLEquals(t *testing.T, opts validatable, format string, args ...any) {
	t.Helper()
	_ = t
	_ = opts
	_ = format
	_ = args
}

func JoinErrors(errs ...error) error {
	return sdk.JoinErrors(errs...)
}

func randomDatabaseObjectIdentifier() DatabaseObjectIdentifier {
	return DatabaseObjectIdentifier{}
}
