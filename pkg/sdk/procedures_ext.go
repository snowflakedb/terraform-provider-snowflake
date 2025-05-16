package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

const (
	DefaultProcedureComment     = "user-defined procedure"
	JavaSnowparkPackageString   = "com.snowflake:snowpark:"
	PythonSnowparkPackageString = "snowflake-snowpark-python=="
)

func (v *Procedure) ID() SchemaObjectIdentifierWithArguments {
	return NewSchemaObjectIdentifierWithArguments(v.CatalogName, v.SchemaName, v.Name, v.ArgumentsOld...)
}

// ProcedureDetails contains aggregated describe results for the given procedure.
type ProcedureDetails struct {
	Signature                  string  // present for all procedure types
	Returns                    string  // present for all procedure types
	Language                   string  // present for all procedure types
	NullHandling               *string // present for all procedure types but SQL
	Body                       *string // present for all procedure types (hidden when SECURE)
	Volatility                 *string // present for all procedure types but SQL
	ExternalAccessIntegrations *string // list present for python, java, and scala
	Secrets                    *string // map present for python, java, and scala
	Imports                    *string // list present for python, java, and scala (hidden when SECURE)
	Handler                    *string // present for python, java, and scala (hidden when SECURE)
	RuntimeVersion             *string // present for python, java, and scala (hidden when SECURE)
	Packages                   *string // list // present for python, java, and scala (hidden when SECURE)
	TargetPath                 *string // present for scala and java (hidden when SECURE)
	InstalledPackages          *string // list present for python (hidden when SECURE)
	ExecuteAs                  string  // present for all procedure types

	NormalizedImports                    []NormalizedPath
	NormalizedTargetPath                 *NormalizedPath
	ReturnDataType                       datatypes.DataType
	ReturnNotNull                        bool
	NormalizedArguments                  []NormalizedArgument
	NormalizedExternalAccessIntegrations []AccountObjectIdentifier
	NormalizedSecrets                    map[string]SchemaObjectIdentifier
	// NormalizedPackages does not contain a snowpark package - it is extracted only as a version in SnowparkVersion below
	NormalizedPackages []string
	SnowparkVersion    string
}

func procedureDetailsFromRows(rows []ProcedureDetail) (*ProcedureDetails, error) {
	v := &ProcedureDetails{}
	var errs []error
	for _, row := range rows {
		switch row.Property {
		case "signature":
			errs = append(errs, row.setStringValueOrError("signature", &v.Signature))
		case "returns":
			errs = append(errs, row.setStringValueOrError("returns", &v.Returns))
		case "language":
			errs = append(errs, row.setStringValueOrError("language", &v.Language))
		case "execute as":
			errs = append(errs, row.setStringValueOrError("execute as", &v.ExecuteAs))
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
			// TODO [SNOW-1850370]: merge these and unit test
			switch strings.ToUpper(v.Language) {
			case "JAVA", "SCALA":
				filtered := make([]string, 0)
				var found bool
				for _, o := range p {
					o := strings.TrimSpace(o)
					if strings.HasPrefix(o, JavaSnowparkPackageString) {
						v.SnowparkVersion = strings.TrimPrefix(o, JavaSnowparkPackageString)
						found = true
					} else {
						filtered = append(filtered, o)
					}
				}
				v.NormalizedPackages = filtered
				if !found {
					errs = append(errs, fmt.Errorf("could not parse package from Snowflake, expected at least snowpark package, got %v", filtered))
				}
			case "PYTHON":
				filtered := make([]string, 0)
				var found bool
				for _, o := range p {
					o := strings.TrimSpace(o)
					if strings.HasPrefix(o, PythonSnowparkPackageString) {
						v.SnowparkVersion = strings.TrimPrefix(o, PythonSnowparkPackageString)
						found = true
					} else {
						filtered = append(filtered, o)
					}
				}
				v.NormalizedPackages = filtered
				if !found {
					errs = append(errs, fmt.Errorf("could not parse package from Snowflake, expected at least snowpark package, got %v", filtered))
				}
			}
		}
	} else {
		switch strings.ToUpper(v.Language) {
		case "JAVA", "SCALA", "PYTHON":
			errs = append(errs, fmt.Errorf("could not parse package from Snowflake, expected at least snowpark package, got nil"))
		default:
			v.NormalizedPackages = []string{}
		}
	}

	return v, errors.Join(errs...)
}

func (d *ProcedureDetail) setStringValueOrError(property string, field *string) error {
	if d.Value == nil {
		return fmt.Errorf("value expected for field %s", property)
	} else {
		*field = *d.Value
	}
	return nil
}

func (d *ProcedureDetail) setOptionalBoolValueOrError(property string, field **bool) error {
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

func (v *procedures) DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*ProcedureDetails, error) {
	rows, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return procedureDetailsFromRows(rows)
}

func (v *procedures) ShowParameters(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Procedure: id,
		},
	})
}

func (s *CreateForJavaProcedureRequest) WithProcedureDefinitionWrapped(procedureDefinition string) *CreateForJavaProcedureRequest {
	s.ProcedureDefinition = String(fmt.Sprintf(`$$%s$$`, procedureDefinition))
	return s
}

func (s *CreateForPythonProcedureRequest) WithProcedureDefinitionWrapped(procedureDefinition string) *CreateForPythonProcedureRequest {
	s.ProcedureDefinition = String(fmt.Sprintf(`$$%s$$`, procedureDefinition))
	return s
}

func (s *CreateForScalaProcedureRequest) WithProcedureDefinitionWrapped(procedureDefinition string) *CreateForScalaProcedureRequest {
	s.ProcedureDefinition = String(fmt.Sprintf(`$$%s$$`, procedureDefinition))
	return s
}

func NewCreateForSQLProcedureRequestDefinitionWrapped(
	name SchemaObjectIdentifier,
	returns ProcedureSQLReturnsRequest,
	procedureDefinition string,
) *CreateForSQLProcedureRequest {
	s := CreateForSQLProcedureRequest{}
	s.name = name
	s.Returns = returns
	s.ProcedureDefinition = fmt.Sprintf(`$$%s$$`, procedureDefinition)
	return &s
}

func NewCreateForJavaScriptProcedureRequestDefinitionWrapped(
	name SchemaObjectIdentifier,
	resultDataType datatypes.DataType,
	procedureDefinition string,
) *CreateForJavaScriptProcedureRequest {
	s := CreateForJavaScriptProcedureRequest{}
	s.name = name
	s.ResultDataType = resultDataType
	s.ProcedureDefinition = fmt.Sprintf(`$$%s$$`, procedureDefinition)
	return &s
}
