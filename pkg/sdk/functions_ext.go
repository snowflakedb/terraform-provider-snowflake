package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

const DefaultFunctionComment = "user-defined function"

func (v *Function) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}

// FunctionDetails contains aggregated describe results for the given function.
type FunctionDetails struct {
	Signature                  string  // present for all function types
	Returns                    string  // present for all function types
	Language                   string  // present for all function types
	Body                       *string // present for all function types (hidden when SECURE)
	NullHandling               *string // present for all function types but SQL
	Volatility                 *string // present for all function types but SQL
	ExternalAccessIntegrations *string // list present for python, java, and scala
	Secrets                    *string // map present for python, java, and scala
	Imports                    *string // list present for python, java, and scala (hidden when SECURE)
	Handler                    *string // present for python, java, and scala (hidden when SECURE)
	RuntimeVersion             *string // present for python, java, and scala (hidden when SECURE)
	Packages                   *string // list // present for python, java, and scala
	TargetPath                 *string // present for scala and java (hidden when SECURE)
	InstalledPackages          *string // list present for python (hidden when SECURE)
	IsAggregate                *bool   // present for python

	NormalizedImports                    []NormalizedPath
	NormalizedTargetPath                 *NormalizedPath
	ReturnDataType                       datatypes.DataType
	ReturnNotNull                        bool
	NormalizedArguments                  []NormalizedArgument
	NormalizedExternalAccessIntegrations []AccountObjectIdentifier
	NormalizedSecrets                    map[string]SchemaObjectIdentifier
	NormalizedPackages                   []string
}

func functionDetailsFromRows(rows []FunctionDetail) (*FunctionDetails, error) {
	v := &FunctionDetails{}
	var errs []error
	for _, row := range rows {
		switch row.Property {
		case "signature":
			errs = append(errs, row.setStringValueOrError("signature", &v.Signature))
		case "returns":
			errs = append(errs, row.setStringValueOrError("returns", &v.Returns))
		case "language":
			errs = append(errs, row.setStringValueOrError("language", &v.Language))
		case "null handling":
			v.NullHandling = row.Value
		case "volatility":
			v.Volatility = row.Value
		case "body":
			v.Body = row.Value
		case "external_access_integrations":
			v.ExternalAccessIntegrations = row.Value
		case "secrets":
			v.Secrets = row.Value
		case "imports":
			v.Imports = row.Value
		case "handler":
			v.Handler = row.Value
		case "runtime_version":
			v.RuntimeVersion = row.Value
		case "packages":
			v.Packages = row.Value
		case "installed_packages":
			v.InstalledPackages = row.Value
		case "is_aggregate":
			errs = append(errs, row.setOptionalBoolValueOrError("is_aggregate", &v.IsAggregate))
		case "target_path":
			v.TargetPath = row.Value
		}
	}
	if e := errors.Join(errs...); e != nil {
		return nil, e
	}

	if normalizedImports, err := parseFunctionOrProcedureImports(v.Imports); err != nil {
		errs = append(errs, err)
	} else {
		v.NormalizedImports = normalizedImports
	}

	if v.TargetPath != nil {
		if p, err := parseFunctionOrProcedureStageLocationPath(*v.TargetPath); err != nil {
			errs = append(errs, err)
		} else {
			v.NormalizedTargetPath = p
		}
	}

	if dt, returnNotNull, err := parseFunctionOrProcedureReturns(v.Returns); err != nil {
		errs = append(errs, err)
	} else {
		v.ReturnDataType = dt
		v.ReturnNotNull = returnNotNull
	}

	if args, err := parseFunctionOrProcedureSignature(v.Signature); err != nil {
		errs = append(errs, err)
	} else {
		v.NormalizedArguments = args
	}

	if v.ExternalAccessIntegrations != nil {
		if p, err := parseFunctionOrProcedureExternalAccessIntegrations(*v.ExternalAccessIntegrations); err != nil {
			errs = append(errs, err)
		} else {
			v.NormalizedExternalAccessIntegrations = p
		}
	} else {
		v.NormalizedExternalAccessIntegrations = []AccountObjectIdentifier{}
	}

	if v.Secrets != nil {
		if p, err := parseFunctionOrProcedureSecrets(*v.Secrets); err != nil {
			errs = append(errs, err)
		} else {
			v.NormalizedSecrets = p
		}
	} else {
		v.NormalizedSecrets = map[string]SchemaObjectIdentifier{}
	}

	if v.Packages != nil {
		if p, err := parseFunctionOrProcedurePackages(*v.Packages); err != nil {
			errs = append(errs, err)
		} else {
			v.NormalizedPackages = p
		}
	} else {
		v.NormalizedPackages = []string{}
	}

	return v, errors.Join(errs...)
}

func (v *functions) DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error) {
	rows, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return functionDetailsFromRows(rows)
}

func (v *functions) ShowParameters(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Function: id,
		},
	})
}

func (d *FunctionDetail) setStringValueOrError(property string, field *string) error {
	if d.Value == nil {
		return fmt.Errorf("value expected for field %s", property)
	} else {
		*field = *d.Value
	}
	return nil
}

func (d *FunctionDetail) setOptionalBoolValueOrError(property string, field **bool) error {
	if d.Value != nil && *d.Value != "" {
		v, err := strconv.ParseBool(*d.Value)
		if err != nil {
			return fmt.Errorf("invalid value for field %s, err: %w", property, err)
		} else {
			*field = Bool(v)
		}
	}
	return nil
}

func (s *CreateForJavaFunctionRequest) WithFunctionDefinitionWrapped(functionDefinition string) *CreateForJavaFunctionRequest {
	s.FunctionDefinition = String(fmt.Sprintf(`$$%s$$`, functionDefinition))
	return s
}

func (s *CreateForPythonFunctionRequest) WithFunctionDefinitionWrapped(functionDefinition string) *CreateForPythonFunctionRequest {
	s.FunctionDefinition = String(fmt.Sprintf(`$$%s$$`, functionDefinition))
	return s
}

func (s *CreateForScalaFunctionRequest) WithFunctionDefinitionWrapped(functionDefinition string) *CreateForScalaFunctionRequest {
	s.FunctionDefinition = String(fmt.Sprintf(`$$%s$$`, functionDefinition))
	return s
}

func NewCreateForSQLFunctionRequestDefinitionWrapped(
	name SchemaObjectIdentifier,
	returns FunctionReturnsRequest,
	functionDefinition string,
) *CreateForSQLFunctionRequest {
	s := CreateForSQLFunctionRequest{}
	s.name = name
	s.Returns = returns
	s.FunctionDefinition = fmt.Sprintf(`$$%s$$`, functionDefinition)
	return &s
}

func NewCreateForJavascriptFunctionRequestDefinitionWrapped(
	name SchemaObjectIdentifier,
	returns FunctionReturnsRequest,
	functionDefinition string,
) *CreateForJavascriptFunctionRequest {
	s := CreateForJavascriptFunctionRequest{}
	s.name = name
	s.Returns = returns
	s.FunctionDefinition = fmt.Sprintf(`$$%s$$`, functionDefinition)
	return &s
}
