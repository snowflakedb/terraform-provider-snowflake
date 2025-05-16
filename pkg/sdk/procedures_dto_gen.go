package sdk

// imports added manually
import "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateForJavaProcedureOptions]              = new(CreateForJavaProcedureRequest)
	_ optionsProvider[CreateForJavaScriptProcedureOptions]        = new(CreateForJavaScriptProcedureRequest)
	_ optionsProvider[CreateForPythonProcedureOptions]            = new(CreateForPythonProcedureRequest)
	_ optionsProvider[CreateForScalaProcedureOptions]             = new(CreateForScalaProcedureRequest)
	_ optionsProvider[CreateForSQLProcedureOptions]               = new(CreateForSQLProcedureRequest)
	_ optionsProvider[AlterProcedureOptions]                      = new(AlterProcedureRequest)
	_ optionsProvider[DropProcedureOptions]                       = new(DropProcedureRequest)
	_ optionsProvider[ShowProcedureOptions]                       = new(ShowProcedureRequest)
	_ optionsProvider[DescribeProcedureOptions]                   = new(DescribeProcedureRequest)
	_ optionsProvider[CallProcedureOptions]                       = new(CallProcedureRequest)
	_ optionsProvider[CreateAndCallForJavaProcedureOptions]       = new(CreateAndCallForJavaProcedureRequest)
	_ optionsProvider[CreateAndCallForScalaProcedureOptions]      = new(CreateAndCallForScalaProcedureRequest)
	_ optionsProvider[CreateAndCallForJavaScriptProcedureOptions] = new(CreateAndCallForJavaScriptProcedureRequest)
	_ optionsProvider[CreateAndCallForPythonProcedureOptions]     = new(CreateAndCallForPythonProcedureRequest)
	_ optionsProvider[CreateAndCallForSQLProcedureOptions]        = new(CreateAndCallForSQLProcedureRequest)
)

type CreateForJavaProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    ProcedureReturnsRequest   // required
	RuntimeVersion             string                    // required
	Packages                   []ProcedurePackageRequest // required
	Imports                    []ProcedureImportRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []SecretReference
	TargetPath                 *string
	NullInputBehavior          *NullInputBehavior
	ReturnResultsBehavior      *ReturnResultsBehavior
	Comment                    *string
	ExecuteAs                  *ExecuteAs
	ProcedureDefinition        *string
}

type ProcedureArgumentRequest struct {
	ArgName        string // required
	ArgDataTypeOld DataType
	ArgDataType    datatypes.DataType // required
	DefaultValue   *string
}

type ProcedureReturnsRequest struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
}

type ProcedureReturnsResultDataTypeRequest struct {
	ResultDataTypeOld DataType
	ResultDataType    datatypes.DataType // required
	Null              *bool
	NotNull           *bool
}

type ProcedureReturnsTableRequest struct {
	Columns []ProcedureColumnRequest
}

type ProcedureColumnRequest struct {
	ColumnName        string // required
	ColumnDataTypeOld DataType
	ColumnDataType    datatypes.DataType // required
}

type ProcedurePackageRequest struct {
	Package string // required
}

type ProcedureImportRequest struct {
	Import string // required
}

type CreateForJavaScriptProcedureRequest struct {
	OrReplace             *bool
	Secure                *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []ProcedureArgumentRequest
	CopyGrants            *bool
	ResultDataTypeOld     DataType
	ResultDataType        datatypes.DataType // required
	NotNull               *bool
	NullInputBehavior     *NullInputBehavior
	ReturnResultsBehavior *ReturnResultsBehavior
	Comment               *string
	ExecuteAs             *ExecuteAs
	ProcedureDefinition   string // required
}

type CreateForPythonProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    ProcedureReturnsRequest   // required
	RuntimeVersion             string                    // required
	Packages                   []ProcedurePackageRequest // required
	Imports                    []ProcedureImportRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []SecretReference
	NullInputBehavior          *NullInputBehavior
	ReturnResultsBehavior      *ReturnResultsBehavior
	Comment                    *string
	ExecuteAs                  *ExecuteAs
	ProcedureDefinition        *string
}

type CreateForScalaProcedureRequest struct {
	OrReplace                  *bool
	Secure                     *bool
	name                       SchemaObjectIdentifier // required
	Arguments                  []ProcedureArgumentRequest
	CopyGrants                 *bool
	Returns                    ProcedureReturnsRequest   // required
	RuntimeVersion             string                    // required
	Packages                   []ProcedurePackageRequest // required
	Imports                    []ProcedureImportRequest
	Handler                    string // required
	ExternalAccessIntegrations []AccountObjectIdentifier
	Secrets                    []SecretReference
	TargetPath                 *string
	NullInputBehavior          *NullInputBehavior
	ReturnResultsBehavior      *ReturnResultsBehavior
	Comment                    *string
	ExecuteAs                  *ExecuteAs
	ProcedureDefinition        *string
}

type CreateForSQLProcedureRequest struct {
	OrReplace             *bool
	Secure                *bool
	name                  SchemaObjectIdentifier // required
	Arguments             []ProcedureArgumentRequest
	CopyGrants            *bool
	Returns               ProcedureSQLReturnsRequest // required
	NullInputBehavior     *NullInputBehavior
	ReturnResultsBehavior *ReturnResultsBehavior
	Comment               *string
	ExecuteAs             *ExecuteAs
	ProcedureDefinition   string // required
}

type ProcedureSQLReturnsRequest struct {
	ResultDataType *ProcedureReturnsResultDataTypeRequest
	Table          *ProcedureReturnsTableRequest
	NotNull        *bool
}

type AlterProcedureRequest struct {
	IfExists  *bool
	name      SchemaObjectIdentifierWithArguments // required
	RenameTo  *SchemaObjectIdentifier
	Set       *ProcedureSetRequest
	Unset     *ProcedureUnsetRequest
	SetTags   []TagAssociation
	UnsetTags []ObjectIdentifier
	ExecuteAs *ExecuteAs
}

type ProcedureSetRequest struct {
	Comment                    *string
	ExternalAccessIntegrations []AccountObjectIdentifier
	SecretsList                *SecretsListRequest
	AutoEventLogging           *AutoEventLogging
	EnableConsoleOutput        *bool
	LogLevel                   *LogLevel
	MetricLevel                *MetricLevel
	TraceLevel                 *TraceLevel
}

// SecretsListRequest removed manually - redeclaration with function

type ProcedureUnsetRequest struct {
	Comment                    *bool
	ExternalAccessIntegrations *bool
	AutoEventLogging           *bool
	EnableConsoleOutput        *bool
	LogLevel                   *bool
	MetricLevel                *bool
	TraceLevel                 *bool
}

type DropProcedureRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifierWithArguments // required
}

type ShowProcedureRequest struct {
	Like *Like
	In   *ExtendedIn
}

type DescribeProcedureRequest struct {
	name SchemaObjectIdentifierWithArguments // required
}

type CallProcedureRequest struct {
	name              SchemaObjectIdentifier // required
	CallArguments     []string
	ScriptingVariable *string
}

type CreateAndCallForJavaProcedureRequest struct {
	Name                AccountObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	Returns             ProcedureReturnsRequest   // required
	RuntimeVersion      string                    // required
	Packages            []ProcedurePackageRequest // required
	Imports             []ProcedureImportRequest
	Handler             string // required
	NullInputBehavior   *NullInputBehavior
	ProcedureDefinition *string
	WithClause          *ProcedureWithClauseRequest
	ProcedureName       AccountObjectIdentifier // required
	CallArguments       []string
	ScriptingVariable   *string
}

type ProcedureWithClauseRequest struct {
	CteName    AccountObjectIdentifier // required
	CteColumns []string
	Statement  string // required
}

type CreateAndCallForScalaProcedureRequest struct {
	Name                AccountObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	Returns             ProcedureReturnsRequest   // required
	RuntimeVersion      string                    // required
	Packages            []ProcedurePackageRequest // required
	Imports             []ProcedureImportRequest
	Handler             string // required
	NullInputBehavior   *NullInputBehavior
	ProcedureDefinition *string
	WithClauses         []ProcedureWithClauseRequest
	ProcedureName       AccountObjectIdentifier // required
	CallArguments       []string
	ScriptingVariable   *string
}

type CreateAndCallForJavaScriptProcedureRequest struct {
	Name                AccountObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	ResultDataTypeOld   DataType
	ResultDataType      datatypes.DataType // required
	NotNull             *bool
	NullInputBehavior   *NullInputBehavior
	ProcedureDefinition string // required
	WithClauses         []ProcedureWithClauseRequest
	ProcedureName       AccountObjectIdentifier // required
	CallArguments       []string
	ScriptingVariable   *string
}

type CreateAndCallForPythonProcedureRequest struct {
	Name                AccountObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	Returns             ProcedureReturnsRequest   // required
	RuntimeVersion      string                    // required
	Packages            []ProcedurePackageRequest // required
	Imports             []ProcedureImportRequest
	Handler             string // required
	NullInputBehavior   *NullInputBehavior
	ProcedureDefinition *string
	WithClauses         []ProcedureWithClauseRequest
	ProcedureName       AccountObjectIdentifier // required
	CallArguments       []string
	ScriptingVariable   *string
}

type CreateAndCallForSQLProcedureRequest struct {
	Name                AccountObjectIdentifier // required
	Arguments           []ProcedureArgumentRequest
	Returns             ProcedureReturnsRequest // required
	NullInputBehavior   *NullInputBehavior
	ProcedureDefinition string // required
	WithClauses         []ProcedureWithClauseRequest
	ProcedureName       AccountObjectIdentifier // required
	CallArguments       []string
	ScriptingVariable   *string
}
