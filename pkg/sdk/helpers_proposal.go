package sdk

import (
	"context"
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

func convertRows[T convertibleRowDeprecated[U], U any](dbRows []T) []U {
	resultList := make([]U, len(dbRows))
	for i, row := range dbRows {
		resultList[i] = *(row.convert())
	}
	return resultList
}

type optionsProvider[T any] interface {
	toOpts() *T
}
