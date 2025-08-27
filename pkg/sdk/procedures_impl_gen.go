package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Procedures                      = (*procedures)(nil)
	_ convertibleRow[ProcedureDetail] = new(procedureDetailRow)
	_ convertibleRow[Procedure]       = new(procedureRow)
)

type procedures struct {
	client *Client
}

func (v *procedures) CreateForJava(ctx context.Context, request *CreateForJavaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForJavaScript(ctx context.Context, request *CreateForJavaScriptProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForPython(ctx context.Context, request *CreateForPythonProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForScala(ctx context.Context, request *CreateForScalaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateForSQL(ctx context.Context, request *CreateForSQLProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) Alter(ctx context.Context, request *AlterProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) Drop(ctx context.Context, request *DropProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) DropSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropProcedureRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *procedures) Show(ctx context.Context, request *ShowProcedureRequest) ([]Procedure, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[procedureRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[procedureRow, Procedure](dbRows)
}

func (v *procedures) ShowByID(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Procedure, error) {
	request := NewShowProcedureRequest().
		WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}}).
		WithLike(Like{Pattern: String(id.Name())})
	procedures, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(procedures, func(r Procedure) bool { return r.ID().FullyQualifiedName() == id.FullyQualifiedName() })
}

func (v *procedures) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Procedure, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *procedures) Describe(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]ProcedureDetail, error) {
	opts := &DescribeProcedureOptions{
		name: id,
	}
	rows, err := validateAndQuery[procedureDetailRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[procedureDetailRow, ProcedureDetail](rows)
}

func (v *procedures) Call(ctx context.Context, request *CallProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateAndCallForJava(ctx context.Context, request *CreateAndCallForJavaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateAndCallForScala(ctx context.Context, request *CreateAndCallForScalaProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateAndCallForJavaScript(ctx context.Context, request *CreateAndCallForJavaScriptProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateAndCallForPython(ctx context.Context, request *CreateAndCallForPythonProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *procedures) CreateAndCallForSQL(ctx context.Context, request *CreateAndCallForSQLProcedureRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateForJavaProcedureRequest) toOpts() *CreateForJavaProcedureOptions {
	opts := &CreateForJavaProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		TargetPath:                 r.TargetPath,
		NullInputBehavior:          r.NullInputBehavior,
		ReturnResultsBehavior:      r.ReturnResultsBehavior,
		Comment:                    r.Comment,
		ExecuteAs:                  r.ExecuteAs,
		ProcedureDefinition:        r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	return opts
}

func (r *CreateForJavaScriptProcedureRequest) toOpts() *CreateForJavaScriptProcedureOptions {
	opts := &CreateForJavaScriptProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants:            r.CopyGrants,
		ResultDataTypeOld:     r.ResultDataTypeOld,
		ResultDataType:        r.ResultDataType,
		NotNull:               r.NotNull,
		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Comment:               r.Comment,
		ExecuteAs:             r.ExecuteAs,
		ProcedureDefinition:   r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	return opts
}

func (r *CreateForPythonProcedureRequest) toOpts() *CreateForPythonProcedureOptions {
	opts := &CreateForPythonProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		NullInputBehavior:          r.NullInputBehavior,
		ReturnResultsBehavior:      r.ReturnResultsBehavior,
		Comment:                    r.Comment,
		ExecuteAs:                  r.ExecuteAs,
		ProcedureDefinition:        r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	return opts
}

func (r *CreateForScalaProcedureRequest) toOpts() *CreateForScalaProcedureOptions {
	opts := &CreateForScalaProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		RuntimeVersion: r.RuntimeVersion,

		Handler:                    r.Handler,
		ExternalAccessIntegrations: r.ExternalAccessIntegrations,
		Secrets:                    r.Secrets,
		TargetPath:                 r.TargetPath,
		NullInputBehavior:          r.NullInputBehavior,
		ReturnResultsBehavior:      r.ReturnResultsBehavior,
		Comment:                    r.Comment,
		ExecuteAs:                  r.ExecuteAs,
		ProcedureDefinition:        r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	return opts
}

func (r *CreateForSQLProcedureRequest) toOpts() *CreateForSQLProcedureOptions {
	opts := &CreateForSQLProcedureOptions{
		OrReplace: r.OrReplace,
		Secure:    r.Secure,
		name:      r.name,

		CopyGrants: r.CopyGrants,

		NullInputBehavior:     r.NullInputBehavior,
		ReturnResultsBehavior: r.ReturnResultsBehavior,
		Comment:               r.Comment,
		ExecuteAs:             r.ExecuteAs,
		ProcedureDefinition:   r.ProcedureDefinition,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureSQLReturns{
		NotNull: r.Returns.NotNull,
	}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
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

func (r *AlterProcedureRequest) toOpts() *AlterProcedureOptions {
	opts := &AlterProcedureOptions{
		IfExists: r.IfExists,
		name:     r.name,
		RenameTo: r.RenameTo,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
		ExecuteAs: r.ExecuteAs,
	}
	if r.Set != nil {
		opts.Set = &ProcedureSet{
			Comment:                    r.Set.Comment,
			ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations,

			AutoEventLogging:    r.Set.AutoEventLogging,
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
		opts.Unset = &ProcedureUnset{
			Comment:                    r.Unset.Comment,
			ExternalAccessIntegrations: r.Unset.ExternalAccessIntegrations,
			AutoEventLogging:           r.Unset.AutoEventLogging,
			EnableConsoleOutput:        r.Unset.EnableConsoleOutput,
			LogLevel:                   r.Unset.LogLevel,
			MetricLevel:                r.Unset.MetricLevel,
			TraceLevel:                 r.Unset.TraceLevel,
		}
	}
	return opts
}

func (r *DropProcedureRequest) toOpts() *DropProcedureOptions {
	opts := &DropProcedureOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowProcedureRequest) toOpts() *ShowProcedureOptions {
	opts := &ShowProcedureOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r procedureRow) convert() (*Procedure, error) {
	e := &Procedure{
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
	}
	arguments := strings.TrimLeft(r.Arguments, r.Name)
	returnIndex := strings.Index(arguments, ") RETURN ")
	e.ReturnTypeOld = DataType(arguments[returnIndex+len(") RETURN "):])
	parsedArguments, err := ParseFunctionAndProcedureArguments(arguments[:returnIndex+1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse procedure arguments: %w", err)
	} else {
		e.ArgumentsOld = collections.Map(parsedArguments, func(a ParsedArgument) DataType {
			return DataType(a.ArgType)
		})
	}
	if r.IsSecure.Valid {
		e.IsSecure = r.IsSecure.String == "Y"
	}
	return e, nil
}

func (r *DescribeProcedureRequest) toOpts() *DescribeProcedureOptions {
	opts := &DescribeProcedureOptions{
		name: r.name,
	}
	return opts
}

func (r procedureDetailRow) convert() (*ProcedureDetail, error) {
	e := &ProcedureDetail{
		Property: r.Property,
	}
	if r.Value.Valid && r.Value.String != "null" {
		e.Value = String(r.Value.String)
	}
	return e, nil
}

func (r *CallProcedureRequest) toOpts() *CallProcedureOptions {
	opts := &CallProcedureOptions{
		call:              false,
		name:              r.name,
		CallArguments:     r.CallArguments,
		ScriptingVariable: r.ScriptingVariable,
	}
	return opts
}

func (r *CreateAndCallForJavaProcedureRequest) toOpts() *CreateAndCallForJavaProcedureOptions {
	opts := &CreateAndCallForJavaProcedureOptions{
		Name: r.Name,

		RuntimeVersion: r.RuntimeVersion,

		Handler:             r.Handler,
		NullInputBehavior:   r.NullInputBehavior,
		ProcedureDefinition: r.ProcedureDefinition,

		ProcedureName:     r.ProcedureName,
		CallArguments:     r.CallArguments,
		ScriptingVariable: r.ScriptingVariable,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.WithClause != nil {
		opts.WithClause = &ProcedureWithClause{
			CteName:    r.WithClause.CteName,
			CteColumns: r.WithClause.CteColumns,
			Statement:  r.WithClause.Statement,
		}
	}
	return opts
}

func (r *CreateAndCallForScalaProcedureRequest) toOpts() *CreateAndCallForScalaProcedureOptions {
	opts := &CreateAndCallForScalaProcedureOptions{
		Name: r.Name,

		RuntimeVersion: r.RuntimeVersion,

		Handler:             r.Handler,
		NullInputBehavior:   r.NullInputBehavior,
		ProcedureDefinition: r.ProcedureDefinition,

		ProcedureName:     r.ProcedureName,
		CallArguments:     r.CallArguments,
		ScriptingVariable: r.ScriptingVariable,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.WithClauses != nil {
		s := make([]ProcedureWithClause, len(r.WithClauses))
		for i, v := range r.WithClauses {
			s[i] = ProcedureWithClause{
				CteName:    v.CteName,
				CteColumns: v.CteColumns,
				Statement:  v.Statement,
			}
		}
		opts.WithClauses = s
	}
	return opts
}

func (r *CreateAndCallForJavaScriptProcedureRequest) toOpts() *CreateAndCallForJavaScriptProcedureOptions {
	opts := &CreateAndCallForJavaScriptProcedureOptions{
		Name: r.Name,

		ResultDataTypeOld:   r.ResultDataTypeOld,
		ResultDataType:      r.ResultDataType,
		NotNull:             r.NotNull,
		NullInputBehavior:   r.NullInputBehavior,
		ProcedureDefinition: r.ProcedureDefinition,

		ProcedureName:     r.ProcedureName,
		CallArguments:     r.CallArguments,
		ScriptingVariable: r.ScriptingVariable,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	if r.WithClauses != nil {
		s := make([]ProcedureWithClause, len(r.WithClauses))
		for i, v := range r.WithClauses {
			s[i] = ProcedureWithClause{
				CteName:    v.CteName,
				CteColumns: v.CteColumns,
				Statement:  v.Statement,
			}
		}
		opts.WithClauses = s
	}
	return opts
}

func (r *CreateAndCallForPythonProcedureRequest) toOpts() *CreateAndCallForPythonProcedureOptions {
	opts := &CreateAndCallForPythonProcedureOptions{
		Name: r.Name,

		RuntimeVersion: r.RuntimeVersion,

		Handler:             r.Handler,
		NullInputBehavior:   r.NullInputBehavior,
		ProcedureDefinition: r.ProcedureDefinition,

		ProcedureName:     r.ProcedureName,
		CallArguments:     r.CallArguments,
		ScriptingVariable: r.ScriptingVariable,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.Packages != nil {
		s := make([]ProcedurePackage, len(r.Packages))
		for i, v := range r.Packages {
			s[i] = ProcedurePackage{
				Package: v.Package,
			}
		}
		opts.Packages = s
	}
	if r.Imports != nil {
		s := make([]ProcedureImport, len(r.Imports))
		for i, v := range r.Imports {
			s[i] = ProcedureImport{
				Import: v.Import,
			}
		}
		opts.Imports = s
	}
	if r.WithClauses != nil {
		s := make([]ProcedureWithClause, len(r.WithClauses))
		for i, v := range r.WithClauses {
			s[i] = ProcedureWithClause{
				CteName:    v.CteName,
				CteColumns: v.CteColumns,
				Statement:  v.Statement,
			}
		}
		opts.WithClauses = s
	}
	return opts
}

func (r *CreateAndCallForSQLProcedureRequest) toOpts() *CreateAndCallForSQLProcedureOptions {
	opts := &CreateAndCallForSQLProcedureOptions{
		Name: r.Name,

		NullInputBehavior:   r.NullInputBehavior,
		ProcedureDefinition: r.ProcedureDefinition,

		ProcedureName:     r.ProcedureName,
		CallArguments:     r.CallArguments,
		ScriptingVariable: r.ScriptingVariable,
	}
	if r.Arguments != nil {
		s := make([]ProcedureArgument, len(r.Arguments))
		for i, v := range r.Arguments {
			s[i] = ProcedureArgument{
				ArgName:        v.ArgName,
				ArgDataTypeOld: v.ArgDataTypeOld,
				ArgDataType:    v.ArgDataType,
				DefaultValue:   v.DefaultValue,
			}
		}
		opts.Arguments = s
	}
	opts.Returns = ProcedureReturns{}
	if r.Returns.ResultDataType != nil {
		opts.Returns.ResultDataType = &ProcedureReturnsResultDataType{
			ResultDataTypeOld: r.Returns.ResultDataType.ResultDataTypeOld,
			ResultDataType:    r.Returns.ResultDataType.ResultDataType,
			Null:              r.Returns.ResultDataType.Null,
			NotNull:           r.Returns.ResultDataType.NotNull,
		}
	}
	if r.Returns.Table != nil {
		opts.Returns.Table = &ProcedureReturnsTable{}
		if r.Returns.Table.Columns != nil {
			s := make([]ProcedureColumn, len(r.Returns.Table.Columns))
			for i, v := range r.Returns.Table.Columns {
				s[i] = ProcedureColumn{
					ColumnName:        v.ColumnName,
					ColumnDataTypeOld: v.ColumnDataTypeOld,
					ColumnDataType:    v.ColumnDataType,
				}
			}
			opts.Returns.Table.Columns = s
		}
	}
	if r.WithClauses != nil {
		s := make([]ProcedureWithClause, len(r.WithClauses))
		for i, v := range r.WithClauses {
			s[i] = ProcedureWithClause{
				CteName:    v.CteName,
				CteColumns: v.CteColumns,
				Statement:  v.Statement,
			}
		}
		opts.WithClauses = s
	}
	return opts
}
