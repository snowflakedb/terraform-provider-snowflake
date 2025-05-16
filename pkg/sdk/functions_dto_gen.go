package sdk

import "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateForJavaFunctionOptions]       = new(CreateForJavaFunctionRequest)
	_ optionsProvider[CreateForJavascriptFunctionOptions] = new(CreateForJavascriptFunctionRequest)
	_ optionsProvider[CreateForPythonFunctionOptions]     = new(CreateForPythonFunctionRequest)
	_ optionsProvider[CreateForScalaFunctionOptions]      = new(CreateForScalaFunctionRequest)
	_ optionsProvider[CreateForSQLFunctionOptions]        = new(CreateForSQLFunctionRequest)
	_ optionsProvider[AlterFunctionOptions]               = new(AlterFunctionRequest)
	_ optionsProvider[DropFunctionOptions]                = new(DropFunctionRequest)
	_ optionsProvider[ShowFunctionOptions]                = new(ShowFunctionRequest)
	_ optionsProvider[DescribeFunctionOptions]            = new(DescribeFunctionRequest)
)

type CreateForJavaFunctionRequest struct {
	OrReplace                  *bool
	Temporary                  *bool
	Secure                     *bool
	IfNotExists                *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []FunctionArgumentRequest
	CopyGrants                 *bool
	Returns                    FunctionReturnsRequest // required
	ReturnNullValues           *ReturnNullValues
	NullInputBehavior          *NullInputBehavior
	ReturnResultsBehavior      *ReturnResultsBehavior
	RuntimeVersion             *string
	Comment                    *string
	Imports                    []FunctionImportRequest
	Packages                   []FunctionPackageRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []SecretReference
	TargetPath                 *string
	EnableConsoleOutput        *bool
	LogLevel                   *LogLevel
	MetricLevel                *MetricLevel
	TraceLevel                 *TraceLevel
	FunctionDefinition         *string
}

type FunctionArgumentRequest struct {
	ArgName        string // required
	ArgDataTypeOld DataType
	ArgDataType    datatypes.DataType // required
	DefaultValue   *string
}

type FunctionReturnsRequest struct {
	ResultDataType *FunctionReturnsResultDataTypeRequest
	Table          *FunctionReturnsTableRequest
}

type FunctionReturnsResultDataTypeRequest struct {
	ResultDataTypeOld DataType
	ResultDataType    datatypes.DataType // required
}

type FunctionReturnsTableRequest struct {
	Columns []FunctionColumnRequest
}

type FunctionColumnRequest struct {
	ColumnName        string // required
	ColumnDataTypeOld DataType
	ColumnDataType    datatypes.DataType // required
}

type FunctionImportRequest struct {
	Import string
}

type FunctionPackageRequest struct {
	Package string
}

type CreateForJavascriptFunctionRequest struct {
	OrReplace             *bool
	Temporary             *bool
	Secure                *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []FunctionArgumentRequest
	CopyGrants            *bool
	Returns               FunctionReturnsRequest // required
	ReturnNullValues      *ReturnNullValues
	NullInputBehavior     *NullInputBehavior
	ReturnResultsBehavior *ReturnResultsBehavior
	Comment               *string
	EnableConsoleOutput   *bool
	LogLevel              *LogLevel
	MetricLevel           *MetricLevel
	TraceLevel            *TraceLevel
	FunctionDefinition    string // required
}

type CreateForPythonFunctionRequest struct {
	OrReplace                  *bool
	Temporary                  *bool
	Secure                     *bool
	Aggregate                  *bool
	IfNotExists                *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []FunctionArgumentRequest
	CopyGrants                 *bool
	Returns                    FunctionReturnsRequest // required
	ReturnNullValues           *ReturnNullValues
	NullInputBehavior          *NullInputBehavior
	ReturnResultsBehavior      *ReturnResultsBehavior
	RuntimeVersion             string // required
	Comment                    *string
	Imports                    []FunctionImportRequest
	Packages                   []FunctionPackageRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []SecretReference
	EnableConsoleOutput        *bool
	LogLevel                   *LogLevel
	MetricLevel                *MetricLevel
	TraceLevel                 *TraceLevel
	FunctionDefinition         *string
}

type CreateForScalaFunctionRequest struct {
	OrReplace                  *bool
	Temporary                  *bool
	Secure                     *bool
	IfNotExists                *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []FunctionArgumentRequest
	CopyGrants                 *bool
	ResultDataTypeOld          DataType
	ResultDataType             datatypes.DataType // required
	ReturnNullValues           *ReturnNullValues
	NullInputBehavior          *NullInputBehavior
	ReturnResultsBehavior      *ReturnResultsBehavior
	RuntimeVersion             string // required
	Comment                    *string
	Imports                    []FunctionImportRequest
	Packages                   []FunctionPackageRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []SecretReference
	TargetPath                 *string
	EnableConsoleOutput        *bool
	LogLevel                   *LogLevel
	MetricLevel                *MetricLevel
	TraceLevel                 *TraceLevel
	FunctionDefinition         *string
}

type CreateForSQLFunctionRequest struct {
	OrReplace             *bool
	Temporary             *bool
	Secure                *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []FunctionArgumentRequest
	CopyGrants            *bool
	Returns               FunctionReturnsRequest // required
	ReturnNullValues      *ReturnNullValues
	ReturnResultsBehavior *ReturnResultsBehavior
	Memoizable            *bool
	Comment               *string
	EnableConsoleOutput   *bool
	LogLevel              *LogLevel
	MetricLevel           *MetricLevel
	TraceLevel            *TraceLevel
	FunctionDefinition    string // required
}

type AlterFunctionRequest struct {
	IfExists    *bool
	name        SchemaObjectIdentifierWithArguments // required
	RenameTo    *SchemaObjectIdentifier
	Set         *FunctionSetRequest
	Unset       *FunctionUnsetRequest
	SetSecure   *bool
	UnsetSecure *bool
	SetTags     []TagAssociation
	UnsetTags   []ObjectIdentifier
}

type FunctionSetRequest struct {
	Comment                    *string
	ExternalAccessIntegrations []AccountObjectIdentifier
	SecretsList                *SecretsListRequest
	EnableConsoleOutput        *bool
	LogLevel                   *LogLevel
	MetricLevel                *MetricLevel
	TraceLevel                 *TraceLevel
}

type SecretsListRequest struct {
	SecretsList []SecretReference // required
}

type FunctionUnsetRequest struct {
	Comment                    *bool
	ExternalAccessIntegrations *bool
	EnableConsoleOutput        *bool
	LogLevel                   *bool
	MetricLevel                *bool
	TraceLevel                 *bool
}

type DropFunctionRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifierWithArguments // required
}

type ShowFunctionRequest struct {
	Like *Like
	In   *ExtendedIn
}

type DescribeFunctionRequest struct {
	name SchemaObjectIdentifierWithArguments // required
}
