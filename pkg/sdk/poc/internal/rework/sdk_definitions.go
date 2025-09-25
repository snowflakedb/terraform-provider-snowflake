package rework

import (
	"context"
	"errors"
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

type convertibleRow[T any] interface {
	convert() (*T, error)
}

type Client = sdk.Client

type (
	ObjectIdentifier         = sdk.ObjectIdentifier
	AccountObjectIdentifier  = sdk.AccountObjectIdentifier
	DatabaseObjectIdentifier = sdk.DatabaseObjectIdentifier
	ExternalObjectIdentifier = sdk.ExternalObjectIdentifier
	SchemaObjectIdentifier   = sdk.SchemaObjectIdentifier
	TableColumnIdentifier    = sdk.TableColumnIdentifier

	In   = sdk.In
	Like = sdk.Like
)

type ValuesBehavior = sdk.ValuesBehavior
type ObjectType = sdk.ObjectType

const ObjectTypeSequence = sdk.ObjectTypeSequence

func NewSchemaObjectIdentifier(_, _, _ string) SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifier("", "", "")
}

func randomSchemaObjectIdentifier() SchemaObjectIdentifier {
	return SchemaObjectIdentifier{}
}

func assertOptsInvalidJoinedErrors(t *testing.T, _ validatable, _ ...error) {
	t.Helper()
}

func assertOptsValidAndSQLEquals(t *testing.T, _ validatable, _ string, _ ...any) {
	t.Helper()
}

func ValidObjectIdentifier(objectIdentifier ObjectIdentifier) bool {
	return sdk.ValidObjectIdentifier(objectIdentifier)
}

func valueSet(_ any) bool {
	return true
}

func everyValueSet(_ ...any) bool {
	return true
}

func exactlyOneValueSet(_ ...any) bool {
	return true
}

func JoinErrors(errs ...error) error {
	return sdk.JoinErrors(errs...)
}

var (
	ErrNilOptions              = sdk.ErrNilOptions
	ErrInvalidObjectIdentifier = sdk.ErrInvalidObjectIdentifier
)

func errOneOf(_ ...string) error {
	return errors.New("")
}

func errExactlyOneOf(_ ...string) error {
	return errors.New("")
}

func validateAndExec(_ *Client, _ context.Context, _ validatable) error {
	return nil
}

func validateAndQuery[T any](_ *Client, _ context.Context, _ validatable) ([]T, error) {
	return nil, nil
}

func validateAndQueryOne[T any](_ *Client, _ context.Context, _ validatable) (*T, error) {
	return nil, nil
}

func convertRows[T convertibleRow[U], U any](_ []T) ([]U, error) {
	return []U{}, nil
}

type ObjectIdentifierConstraint = sdk.ObjectIdentifierConstraint

func SafeShowById[T any, ID ObjectIdentifierConstraint](
	c *Client,
	f func(context.Context, ID) (T, error),
	con context.Context,
	id ID,
) (T, error) {
	return sdk.SafeShowById(c, f, con, id)
}

func SafeDrop[ID ObjectIdentifierConstraint](
	c *Client,
	f func() error,
	con context.Context,
	id ID,
) error {
	return sdk.SafeDrop(c, f, con, id)
}
