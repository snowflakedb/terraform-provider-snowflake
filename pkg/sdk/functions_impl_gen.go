package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Functions = (*functions)(nil)
var _ convertibleRow[FunctionDetail] = new(functionDetailRow)
var _ convertibleRow[Function] = new(functionRow)

type functions struct {
	client *Client
}

func (v *functions) CreateForJava(ctx context.Context, request *CreateForJavaFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForJavascript(ctx context.Context, request *CreateForJavascriptFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForPython(ctx context.Context, request *CreateForPythonFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForScala(ctx context.Context, request *CreateForScalaFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) CreateForSQL(ctx context.Context, request *CreateForSQLFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) Alter(ctx context.Context, request *AlterFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) Drop(ctx context.Context, request *DropFunctionRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *functions) DropSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropFunctionRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *functions) Show(ctx context.Context, request *ShowFunctionRequest) ([]Function, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[functionRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[functionRow, Function](dbRows)
}

func (v *functions) ShowByID(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Function, error) {
	request := NewShowFunctionRequest().
		WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}}).
		WithLike(Like{Pattern: String(id.Name())})
	functions, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(functions, func(r Function) bool { return r.ID().FullyQualifiedName() == id.FullyQualifiedName() })
}

func (v *functions) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Function, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *functions) Describe(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]FunctionDetail, error) {
	opts := &DescribeFunctionOptions{
		name: id,
	}
	rows, err := validateAndQuery[functionDetailRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[functionDetailRow, FunctionDetail](rows)
}

func (r *CreateForJavaFunctionRequest) toOpts() *CreateForJavaFunctionOptions {
	opts := &CreateForJavaFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		RuntimeVersion:        r.RuntimeVersion,
		Comment:               r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		TargetPath:                 r.TargetPath,
		EnableConsoleOutput:        r.EnableConsoleOutput,
		LogLevel:                   r.LogLevel,
		MetricLevel:                r.MetricLevel,
		TraceLevel:                 r.TraceLevel,
		FunctionDefinition:         r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &FunctionReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]FunctionColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = FunctionColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Imports != nil {
		s := make([]FunctionImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	return opts
}

func (r *CreateForJavascriptFunctionRequest) toOpts() *CreateForJavascriptFunctionOptions {
	opts := &CreateForJavascriptFunctionOptions{
		OrReplace: r.OrReplace,
		Temporary: r.Temporary,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Comment:               r.Comment,
		EnableConsoleOutput:   r.EnableConsoleOutput,
		LogLevel:              r.LogLevel,
		MetricLevel:           r.MetricLevel,
		TraceLevel:            r.TraceLevel,
		FunctionDefinition:    r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &FunctionReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]FunctionColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = FunctionColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	return opts
}

func (r *CreateForPythonFunctionRequest) toOpts() *CreateForPythonFunctionOptions {
	opts := &CreateForPythonFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		Aggregate:   r.Aggregate,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		RuntimeVersion:        r.RuntimeVersion,
		Comment:               r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		EnableConsoleOutput:        r.EnableConsoleOutput,
		LogLevel:                   r.LogLevel,
		MetricLevel:                r.MetricLevel,
		TraceLevel:                 r.TraceLevel,
		FunctionDefinition:         r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &FunctionReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]FunctionColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = FunctionColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Imports != nil {
		s := make([]FunctionImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	return opts
}

func (r *CreateForScalaFunctionRequest) toOpts() *CreateForScalaFunctionOptions {
	opts := &CreateForScalaFunctionOptions{
		OrReplace:   r.OrReplace,
		Temporary:   r.Temporary,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants:            r.CopyGrants,
		ResultDataTypeOld:     r.ResultDataTypeOld,
		ResultDataType:        r.ResultDataType,
		ReturnNullValues:      r.ReturnNullValues,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		RuntimeVersion:        r.RuntimeVersion,
		Comment:               r.Comment,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		TargetPath:                 r.TargetPath,
		EnableConsoleOutput:        r.EnableConsoleOutput,
		LogLevel:                   r.LogLevel,
		MetricLevel:                r.MetricLevel,
		TraceLevel:                 r.TraceLevel,
		FunctionDefinition:         r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	if r.Imports != nil {
		s := make([]FunctionImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = FunctionImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.Packages != nil {
		s := make([]FunctionPackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = FunctionPackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	return opts
}

func (r *CreateForSQLFunctionRequest) toOpts() *CreateForSQLFunctionOptions {
	opts := &CreateForSQLFunctionOptions{
		OrReplace: r.OrReplace,
		Temporary: r.Temporary,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		ReturnNullValues:      r.ReturnNullValues,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Memoizable:            r.Memoizable,
		Comment:               r.Comment,
		EnableConsoleOutput:   r.EnableConsoleOutput,
		LogLevel:              r.LogLevel,
		MetricLevel:           r.MetricLevel,
		TraceLevel:            r.TraceLevel,
		FunctionDefinition:    r.FunctionDefinition,
	}
	if r.Arguments != nil {
		s := make([]FunctionArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = FunctionArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = FunctionReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &FunctionReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &FunctionReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]FunctionColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = FunctionColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	return opts
}

func (r *AlterFunctionRequest) toOpts() *AlterFunctionOptions {
	opts := &AlterFunctionOptions{
		IfExists:    r.IfExists,
		name:        r.name,
		RenameTo:    r.RenameTo,
		SetSecure:   r.SetSecure,
		UnsetSecure: r.UnsetSecure,
		SetTags:     r.SetTags,
		UnsetTags:   r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &FunctionSet{
			Comment:                    r.Set.Comment,
			ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations,

			EnableConsoleOutput: r.Set.EnableConsoleOutput,
			LogLevel:            r.Set.LogLevel,
			MetricLevel:         r.Set.MetricLevel,
			TraceLevel:          r.Set.TraceLevel,
		}
		if r.Set.SecretsList != nil {
			opts.Set.SecretsList = &SecretsList{
				SecretsList: r.Set.SecretsList.SecretsList,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &FunctionUnset{
			Comment:                    r.Unset.Comment,
			ExternalAccessIntegrations: r.Unset.ExternalAccessIntegrations,
			EnableConsoleOutput:        r.Unset.EnableConsoleOutput,
			LogLevel:                   r.Unset.LogLevel,
			MetricLevel:                r.Unset.MetricLevel,
			TraceLevel:                 r.Unset.TraceLevel,
		}
	}
	return opts
}

func (r *DropFunctionRequest) toOpts() *DropFunctionOptions {
	opts := &DropFunctionOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowFunctionRequest) toOpts() *ShowFunctionOptions {
	opts := &ShowFunctionOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r functionRow) convertErr() (*Function, error) {
	e := &Function{
		CreatedOn:          r.CreatedOn,
		Name:               r.Name,
		SchemaName:         strings.Trim(r.SchemaName, `"`),
		IsBuiltin:          r.IsBuiltin == "Y",
		IsAggregate:        r.IsAggregate == "Y",
		IsAnsi:             r.IsAnsi == "Y",
		MinNumArguments:    r.MinNumArguments,
		MaxNumArguments:    r.MaxNumArguments,
		ArgumentsRaw:       r.Arguments,
		Description:        r.Description,
		CatalogName:        strings.Trim(r.CatalogName, `"`),
		IsTableFunction:    r.IsTableFunction == "Y",
		ValidForClustering: r.ValidForClustering == "Y",
		IsExternalFunction: r.IsExternalFunction == "Y",
		Language:           r.Language,
	}
	arguments := strings.TrimLeft(r.Arguments, r.Name)
	returnIndex := strings.Index(arguments, ") RETURN ")
	e.ReturnTypeOld = DataType(arguments[returnIndex+len(") RETURN "):])
	parsedArguments, err := ParseFunctionAndProcedureArguments(arguments[:returnIndex+1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse function arguments: %w", err)
	} else {
		e.ArgumentsOld = collections.Map(parsedArguments, func(a ParsedArgument) DataType {
			return DataType(a.ArgType)
		})
	}

	if r.IsSecure.Valid {
		e.IsSecure = r.IsSecure.String == "Y"
	}
	if r.Secrets.Valid {
		e.Secrets = String(r.Secrets.String)
	}
	if r.ExternalAccessIntegrations.Valid {
		e.ExternalAccessIntegrations = String(r.ExternalAccessIntegrations.String)
	}
	if r.IsMemoizable.Valid {
		e.IsMemoizable = r.IsMemoizable.String == "Y"
	}
	if r.IsDataMetric.Valid {
		e.IsDataMetric = r.IsDataMetric.String == "Y"
	}
	return e, nil
}

func (r *DescribeFunctionRequest) toOpts() *DescribeFunctionOptions {
	opts := &DescribeFunctionOptions{
		name: r.name,
	}
	return opts
}

func (r functionDetailRow) convertErr() (*FunctionDetail, error) {
	e := &FunctionDetail{
		Property: r.Property,
	}
	if r.Value.Valid && r.Value.String != "null" {
		e.Value = String(r.Value.String)
	}
	return e, nil
}
