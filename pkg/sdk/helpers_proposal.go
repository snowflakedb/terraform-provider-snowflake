package sdk

import (
	"context"
	"fmt"
)

// validatable is our sdk interface for anything that can be validated (e.g. CreateXxxOptions).
type validatable interface {
	// validate is a method which should contain all rules for checking validity of the implementing type.
	validate() error
}

// validateAndExec is just a proposal how we can remove some of the boilerplate.
func validateAndExec(client *Client, ctx context.Context, opts validatable) error {
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = client.exec(ctx, sql)
	return err
}

// validateAndQuery is just a proposal how we can remove some of the boilerplate.
func validateAndQuery[T any](client *Client, ctx context.Context, opts validatable) ([]T, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}

	var dest []T
	err = client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

// validateAndQueryOne is just a proposal how we can remove some of the boilerplate.
func validateAndQueryOne[T any](client *Client, ctx context.Context, opts validatable) (*T, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}

	var dest T
	err = client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func createIfNil[T any](t *T) *T {
	if t == nil {
		return new(T)
	}
	return t
}

type convertibleRowDeprecated[T any] interface {
	convert() *T
}

// TODO: Used in Option 1 for convert (remove if not chosen)
type ConvertibleRowDeprecated[T any] interface {
	Convert() *T
}

func convertRows[T convertibleRowDeprecated[U], U any](dbRows []T) []U {
	resultList := make([]U, len(dbRows))
	for i, row := range dbRows {
		resultList[i] = *(row.convert())
	}
	return resultList
}

type convertibleRow[T any] interface {
	// TODO [SNOW-2259477]: rename to convert
	convertErr() (*T, error)
}

// TODO [SNOW-2259477]: rename to convertRows
func convertRowsErr[T convertibleRow[U], U any](dbRows []T) ([]U, error) {
	resultList := make([]U, len(dbRows))
	for i, row := range dbRows {
		converted, err := conversionErrorWrapped(row.convertErr())
		if err != nil {
			return nil, err
		}
		resultList[i] = *converted
	}
	return resultList, nil
}

func conversionErrorWrapped[U any](converted *U, err error) (*U, error) {
	if err != nil {
		return nil, fmt.Errorf("conversion from Snowflake failed with error: %w", err)
	} else {
		return converted, nil
	}
}

type optionsProvider[T any] interface {
	toOpts() *T
}
