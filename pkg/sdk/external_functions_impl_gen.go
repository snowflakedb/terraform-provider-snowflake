package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ExternalFunctions = (*externalFunctions)(nil)
var _ convertibleRow[ExternalFunctionProperty] = new(externalFunctionPropertyRow)
var _ convertibleRow[ExternalFunction] = new(externalFunctionRow)

type externalFunctions struct {
	client *Client
}

func (v *externalFunctions) Create(ctx context.Context, request *CreateExternalFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalFunctions) Alter(ctx context.Context, request *AlterExternalFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *externalFunctions) Show(ctx context.Context, request *ShowExternalFunctionRequest) ([]ExternalFunction, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[externalFunctionRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[externalFunctionRow, ExternalFunction](dbRows)
}

func (v *externalFunctions) ShowByID(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*ExternalFunction, error) {
	request := NewShowExternalFunctionRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	externalFunctions, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(externalFunctions, func(r ExternalFunction) bool { return r.ID().FullyQualifiedName() == id.FullyQualifiedName() })
}

func (v *externalFunctions) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*ExternalFunction, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *externalFunctions) Describe(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]ExternalFunctionProperty, error) {
	opts := &DescribeFunctionOptions{
		name: id,
	}
	rows, err := validateAndQuery[externalFunctionPropertyRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[externalFunctionPropertyRow, ExternalFunctionProperty](rows)
}

func (r *CreateExternalFunctionRequest) toOpts() *CreateExternalFunctionOptions {
	opts := &CreateExternalFunctionOptions{
		OrReplace:             r.OrReplace,
		Secure:                r.Secure,
		name:                  r.name.WithoutArguments(),
		ResultDataType:        r.ResultDataType,
		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Comment:               r.Comment,
		ApiIntegration:        r.ApiIntegration,
		MaxBatchRows:          r.MaxBatchRows,
		Compression:           r.Compression,
		RequestTranslator:     r.RequestTranslator,
		ResponseTranslator:    r.ResponseTranslator,
		As:                    r.As,
	}
	if r.Arguments != nil {
		s := make([]ExternalFunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ExternalFunctionArgument(v)
		}
		opts.Arguments = s
	}
	if r.Headers != nil {
		s := make([]ExternalFunctionHeader, len(r.Headers))
		for i, v := range r.Headers {
			s[i] = ExternalFunctionHeader(v)
		}
		opts.Headers = s
	}
	if r.ContextHeaders != nil {
		s := make([]ExternalFunctionContextHeader, len(r.ContextHeaders))
		for i, v := range r.ContextHeaders {
			s[i] = ExternalFunctionContextHeader(v)
		}
		opts.ContextHeaders = s
	}
	return opts
}

func (r *AlterExternalFunctionRequest) toOpts() *AlterExternalFunctionOptions {
	opts := &AlterExternalFunctionOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.Set != nil {
		opts.Set = &ExternalFunctionSet{
			ApiIntegration: r.Set.ApiIntegration,

			MaxBatchRows:       r.Set.MaxBatchRows,
			Compression:        r.Set.Compression,
			RequestTranslator:  r.Set.RequestTranslator,
			ResponseTranslator: r.Set.ResponseTranslator,
		}
		if r.Set.Headers != nil {
			s := make([]ExternalFunctionHeader, len(r.Set.Headers))
			for i, v := range r.Set.Headers {
				s[i] = ExternalFunctionHeader(v)
			}
			opts.Set.Headers = s
		}
		if r.Set.ContextHeaders != nil {
			s := make([]ExternalFunctionContextHeader, len(r.Set.ContextHeaders))
			for i, v := range r.Set.ContextHeaders {
				s[i] = ExternalFunctionContextHeader(v)
			}
			opts.Set.ContextHeaders = s
		}
	}
	if r.Unset != nil {
		opts.Unset = &ExternalFunctionUnset{
			Comment:            r.Unset.Comment,
			Headers:            r.Unset.Headers,
			ContextHeaders:     r.Unset.ContextHeaders,
			MaxBatchRows:       r.Unset.MaxBatchRows,
			Compression:        r.Unset.Compression,
			Secure:             r.Unset.Secure,
			RequestTranslator:  r.Unset.RequestTranslator,
			ResponseTranslator: r.Unset.ResponseTranslator,
		}
	}
	return opts
}

func (r *ShowExternalFunctionRequest) toOpts() *ShowExternalFunctionOptions {
	opts := &ShowExternalFunctionOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r externalFunctionRow) convert() (*ExternalFunction, error) {
	e := &ExternalFunction{
		CreatedOn:          r.CreatedOn,
		Name:               r.Name,
		IsBuiltin:          r.IsBuiltin == "Y",
		IsAggregate:        r.IsAggregate == "Y",
		IsAnsi:             r.IsAnsi == "Y",
		MinNumArguments:    r.MinNumArguments,
		MaxNumArguments:    r.MaxNumArguments,
		ArgumentsRaw:       r.Arguments,
		Description:        r.Description,
		IsTableFunction:    r.IsTableFunction == "Y",
		ValidForClustering: r.ValidForClustering == "Y",
		IsExternalFunction: r.IsExternalFunction == "Y",
		Language:           r.Language,
	}
	arguments := strings.TrimLeft(r.Arguments, r.Name)
	returnIndex := strings.Index(arguments, ") RETURN ")
	parsedArguments, err := ParseFunctionAndProcedureArguments(arguments[:returnIndex+1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse external function arguments: %w", err)
	} else {
		e.Arguments = collections.Map(parsedArguments, func(a ParsedArgument) DataType {
			return DataType(a.ArgType)
		})
	}
	if r.SchemaName.Valid {
		e.SchemaName = strings.Trim(r.SchemaName.String, `"`)
	}
	if r.CatalogName.Valid {
		e.CatalogName = strings.Trim(r.CatalogName.String, `"`)
	}
	if r.IsSecure.Valid {
		e.IsSecure = r.IsSecure.String == "Y"
	}
	if r.IsMemoizable.Valid {
		e.IsMemoizable = r.IsMemoizable.String == "Y"
	}
	if r.IsDataMetric.Valid {
		e.IsDataMetric = r.IsDataMetric.String == "Y"
	}
	return e, nil
}

func (r *DescribeExternalFunctionRequest) toOpts() *DescribeExternalFunctionOptions {
	opts := &DescribeExternalFunctionOptions{
		name: r.name,
	}
	return opts
}

func (r externalFunctionPropertyRow) convert() (*ExternalFunctionProperty, error) {
	return &ExternalFunctionProperty{
		Property: r.Property,
		Value:    r.Value,
	}, nil
}
