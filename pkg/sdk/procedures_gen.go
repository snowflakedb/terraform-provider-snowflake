package sdk

import (
	"context"
	"database/sql"

	// import added manually
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
)

type Procedures interface {
	CreateForJava(ctx context.Context, request *CreateForJavaProcedureRequest) error
	CreateForJavaScript(ctx context.Context, request *CreateForJavaScriptProcedureRequest) error
	CreateForPython(ctx context.Context, request *CreateForPythonProcedureRequest) error
	CreateForScala(ctx context.Context, request *CreateForScalaProcedureRequest) error
	CreateForSQL(ctx context.Context, request *CreateForSQLProcedureRequest) error
	Alter(ctx context.Context, request *AlterProcedureRequest) error
	Drop(ctx context.Context, request *DropProcedureRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) error
	Show(ctx context.Context, request *ShowProcedureRequest) ([]Procedure, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Procedure, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Procedure, error)
	Describe(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]ProcedureDetail, error)
	Call(ctx context.Context, request *CallProcedureRequest) error
	CreateAndCallForJava(ctx context.Context, request *CreateAndCallForJavaProcedureRequest) error
	CreateAndCallForScala(ctx context.Context, request *CreateAndCallForScalaProcedureRequest) error
	CreateAndCallForJavaScript(ctx context.Context, request *CreateAndCallForJavaScriptProcedureRequest) error
	CreateAndCallForPython(ctx context.Context, request *CreateAndCallForPythonProcedureRequest) error
	CreateAndCallForSQL(ctx context.Context, request *CreateAndCallForSQLProcedureRequest) error

	// DescribeDetails is added manually; it returns aggregated describe results for the given procedure.
	DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*ProcedureDetails, error)
	ShowParameters(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]*Parameter, error)
}

// CreateForJavaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#java-handler.
type CreateForJavaProcedureOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                      `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []ProcedureArgument       `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    ProcedureReturns          `ddl:"keyword" sql:"RETURNS"`
	languageJava               bool                      `ddl:"static" sql:"LANGUAGE JAVA"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage        `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport         `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ExecuteAs                `ddl:"parameter,no_quotes,no_equals" sql:"EXECUTE AS"`
	ProcedureDefinition        *string                   `ddl:"parameter,no_equals" sql:"AS"`
}

type ProcedureArgument struct {
	ArgName        string             `ddl:"keyword,double_quotes"`
	ArgDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ArgDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
	DefaultValue   *string            `ddl:"parameter,no_equals" sql:"DEFAULT"`
}

type ProcedureReturns struct {
	ResultDataType *ProcedureReturnsResultDataType `ddl:"keyword"`
	Table          *ProcedureReturnsTable          `ddl:"keyword" sql:"TABLE"`
}

type ProcedureReturnsResultDataType struct {
	ResultDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ResultDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
	Null              *bool              `ddl:"keyword" sql:"NULL"`
	NotNull           *bool              `ddl:"keyword" sql:"NOT NULL"`
}

type ProcedureReturnsTable struct {
	Columns []ProcedureColumn `ddl:"list,must_parentheses"`
}

type ProcedureColumn struct {
	ColumnName        string             `ddl:"keyword,double_quotes"`
	ColumnDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ColumnDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
}

type ProcedurePackage struct {
	Package string `ddl:"keyword,single_quotes"`
}

type ProcedureImport struct {
	Import string `ddl:"keyword,single_quotes"`
}

// CreateForJavaScriptProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#javascript-handler.
type CreateForJavaScriptProcedureOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	procedure             bool                   `ddl:"static" sql:"PROCEDURE"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []ProcedureArgument    `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	returns               bool                   `ddl:"static" sql:"RETURNS"`
	ResultDataTypeOld     DataType               `ddl:"parameter,no_equals"`
	ResultDataType        datatypes.DataType     `ddl:"parameter,no_quotes,no_equals"`
	NotNull               *bool                  `ddl:"keyword" sql:"NOT NULL"`
	languageJavascript    bool                   `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior     *NullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs             *ExecuteAs             `ddl:"parameter,no_quotes,no_equals" sql:"EXECUTE AS"`
	ProcedureDefinition   string                 `ddl:"parameter,no_equals" sql:"AS"`
}

// CreateForPythonProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#python-handler.
type CreateForPythonProcedureOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                      `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []ProcedureArgument       `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    ProcedureReturns          `ddl:"keyword" sql:"RETURNS"`
	languagePython             bool                      `ddl:"static" sql:"LANGUAGE PYTHON"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage        `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport         `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ExecuteAs                `ddl:"parameter,no_quotes,no_equals" sql:"EXECUTE AS"`
	ProcedureDefinition        *string                   `ddl:"parameter,no_equals" sql:"AS"`
}

// CreateForScalaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#scala-handler.
type CreateForScalaProcedureOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	procedure                  bool                      `ddl:"static" sql:"PROCEDURE"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []ProcedureArgument       `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    ProcedureReturns          `ddl:"keyword" sql:"RETURNS"`
	languageScala              bool                      `ddl:"static" sql:"LANGUAGE SCALA"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages                   []ProcedurePackage        `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports                    []ProcedureImport         `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs                  *ExecuteAs                `ddl:"parameter,no_quotes,no_equals" sql:"EXECUTE AS"`
	ProcedureDefinition        *string                   `ddl:"parameter,no_equals" sql:"AS"`
}

// CreateForSQLProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-procedure#snowflake-scripting-handler.
type CreateForSQLProcedureOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	procedure             bool                   `ddl:"static" sql:"PROCEDURE"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []ProcedureArgument    `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               ProcedureSQLReturns    `ddl:"keyword" sql:"RETURNS"`
	languageSql           bool                   `ddl:"static" sql:"LANGUAGE SQL"`
	NullInputBehavior     *NullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExecuteAs             *ExecuteAs             `ddl:"parameter,no_quotes,no_equals" sql:"EXECUTE AS"`
	ProcedureDefinition   string                 `ddl:"parameter,no_equals" sql:"AS"`
}

type ProcedureSQLReturns struct {
	ResultDataType *ProcedureReturnsResultDataType `ddl:"keyword"`
	Table          *ProcedureReturnsTable          `ddl:"keyword" sql:"TABLE"`
	NotNull        *bool                           `ddl:"keyword" sql:"NOT NULL"`
}

// AlterProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-procedure.
type AlterProcedureOptions struct {
	alter     bool                                `ddl:"static" sql:"ALTER"`
	procedure bool                                `ddl:"static" sql:"PROCEDURE"`
	IfExists  *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	name      SchemaObjectIdentifierWithArguments `ddl:"identifier"`
	RenameTo  *SchemaObjectIdentifier             `ddl:"identifier" sql:"RENAME TO"`
	Set       *ProcedureSet                       `ddl:"list" sql:"SET"`
	Unset     *ProcedureUnset                     `ddl:"list" sql:"UNSET"`
	SetTags   []TagAssociation                    `ddl:"keyword" sql:"SET TAG"`
	UnsetTags []ObjectIdentifier                  `ddl:"keyword" sql:"UNSET TAG"`
	ExecuteAs *ExecuteAs                          `ddl:"parameter,no_quotes,no_equals" sql:"EXECUTE AS"`
}

type ProcedureSet struct {
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	SecretsList                *SecretsList              `ddl:"parameter,parentheses" sql:"SECRETS"`
	AutoEventLogging           *AutoEventLogging         `ddl:"parameter,single_quotes" sql:"AUTO_EVENT_LOGGING"`
	EnableConsoleOutput        *bool                     `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *LogLevel                 `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel                *MetricLevel              `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel                 *TraceLevel               `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
}

// SecretsList removed manually - redeclared in functions

type ProcedureUnset struct {
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
	ExternalAccessIntegrations *bool `ddl:"keyword" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	AutoEventLogging           *bool `ddl:"keyword" sql:"AUTO_EVENT_LOGGING"`
	EnableConsoleOutput        *bool `ddl:"keyword" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	MetricLevel                *bool `ddl:"keyword" sql:"METRIC_LEVEL"`
	TraceLevel                 *bool `ddl:"keyword" sql:"TRACE_LEVEL"`
}

// DropProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-procedure.
type DropProcedureOptions struct {
	drop      bool                                `ddl:"static" sql:"DROP"`
	procedure bool                                `ddl:"static" sql:"PROCEDURE"`
	IfExists  *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	name      SchemaObjectIdentifierWithArguments `ddl:"identifier"`
}

// ShowProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-procedures.
type ShowProcedureOptions struct {
	show       bool        `ddl:"static" sql:"SHOW"`
	procedures bool        `ddl:"static" sql:"PROCEDURES"`
	Like       *Like       `ddl:"keyword" sql:"LIKE"`
	In         *ExtendedIn `ddl:"keyword" sql:"IN"`
}

type procedureRow struct {
	CreatedOn                  string         `db:"created_on"`
	Name                       string         `db:"name"`
	SchemaName                 string         `db:"schema_name"`
	IsBuiltin                  string         `db:"is_builtin"`
	IsAggregate                string         `db:"is_aggregate"`
	IsAnsi                     string         `db:"is_ansi"`
	MinNumArguments            int            `db:"min_num_arguments"`
	MaxNumArguments            int            `db:"max_num_arguments"`
	Arguments                  string         `db:"arguments"`
	Description                string         `db:"description"`
	CatalogName                string         `db:"catalog_name"`
	IsTableFunction            string         `db:"is_table_function"`
	ValidForClustering         string         `db:"valid_for_clustering"`
	IsSecure                   sql.NullString `db:"is_secure"`
	Secrets                    sql.NullString `db:"secrets"`
	ExternalAccessIntegrations sql.NullString `db:"external_access_integrations"`
}

type Procedure struct {
	CreatedOn                  string
	Name                       string
	SchemaName                 string
	IsBuiltin                  bool
	IsAggregate                bool
	IsAnsi                     bool
	MinNumArguments            int
	MaxNumArguments            int
	ArgumentsOld               []DataType
	ArgumentsRaw               string
	Description                string
	CatalogName                string
	IsTableFunction            bool
	ValidForClustering         bool
	IsSecure                   bool
	Secrets                    *string
	ExternalAccessIntegrations *string
}

// DescribeProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-procedure.
type DescribeProcedureOptions struct {
	describe  bool                                `ddl:"static" sql:"DESCRIBE"`
	procedure bool                                `ddl:"static" sql:"PROCEDURE"`
	name      SchemaObjectIdentifierWithArguments `ddl:"identifier"`
}

type procedureDetailRow struct {
	Property string         `db:"property"`
	Value    sql.NullString `db:"value"`
}

type ProcedureDetail struct {
	Property string
	Value    *string
}

// CallProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/call.
type CallProcedureOptions struct {
	call              bool                   `ddl:"static" sql:"CALL"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	CallArguments     []string               `ddl:"keyword,must_parentheses"`
	ScriptingVariable *string                `ddl:"parameter,no_quotes,no_equals" sql:"INTO"`
}

// CreateAndCallForJavaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/call-with#java-and-scala.
type CreateAndCallForJavaProcedureOptions struct {
	with                bool                    `ddl:"static" sql:"WITH"`
	Name                AccountObjectIdentifier `ddl:"identifier"`
	asProcedure         bool                    `ddl:"static" sql:"AS PROCEDURE"`
	Arguments           []ProcedureArgument     `ddl:"list,must_parentheses"`
	Returns             ProcedureReturns        `ddl:"keyword" sql:"RETURNS"`
	languageJava        bool                    `ddl:"static" sql:"LANGUAGE JAVA"`
	NullInputBehavior   *NullInputBehavior      `ddl:"keyword"`
	RuntimeVersion      string                  `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages            []ProcedurePackage      `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports             []ProcedureImport       `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler             string                  `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ProcedureDefinition *string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
	WithClause          *ProcedureWithClause    `ddl:"keyword"`
	call                bool                    `ddl:"static" sql:"CALL"`
	ProcedureName       AccountObjectIdentifier `ddl:"identifier"`
	CallArguments       []string                `ddl:"keyword,must_parentheses"`
	ScriptingVariable   *string                 `ddl:"parameter,no_quotes,no_equals" sql:"INTO"`
}

type ProcedureWithClause struct {
	prefix     bool                    `ddl:"static" sql:","`
	CteName    AccountObjectIdentifier `ddl:"identifier"`
	CteColumns []string                `ddl:"keyword,parentheses"`
	Statement  string                  `ddl:"parameter,no_quotes,no_equals" sql:"AS"`
}

// CreateAndCallForScalaProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/call-with#java-and-scala.
type CreateAndCallForScalaProcedureOptions struct {
	with                bool                    `ddl:"static" sql:"WITH"`
	Name                AccountObjectIdentifier `ddl:"identifier"`
	asProcedure         bool                    `ddl:"static" sql:"AS PROCEDURE"`
	Arguments           []ProcedureArgument     `ddl:"list,must_parentheses"`
	Returns             ProcedureReturns        `ddl:"keyword" sql:"RETURNS"`
	languageScala       bool                    `ddl:"static" sql:"LANGUAGE SCALA"`
	NullInputBehavior   *NullInputBehavior      `ddl:"keyword"`
	RuntimeVersion      string                  `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages            []ProcedurePackage      `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports             []ProcedureImport       `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler             string                  `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ProcedureDefinition *string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
	WithClauses         []ProcedureWithClause   `ddl:"keyword"`
	call                bool                    `ddl:"static" sql:"CALL"`
	ProcedureName       AccountObjectIdentifier `ddl:"identifier"`
	CallArguments       []string                `ddl:"keyword,must_parentheses"`
	ScriptingVariable   *string                 `ddl:"parameter,no_quotes,no_equals" sql:"INTO"`
}

// CreateAndCallForJavaScriptProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/call-with#javascript.
type CreateAndCallForJavaScriptProcedureOptions struct {
	with                bool                    `ddl:"static" sql:"WITH"`
	Name                AccountObjectIdentifier `ddl:"identifier"`
	asProcedure         bool                    `ddl:"static" sql:"AS PROCEDURE"`
	Arguments           []ProcedureArgument     `ddl:"list,must_parentheses"`
	returns             bool                    `ddl:"static" sql:"RETURNS"`
	ResultDataTypeOld   DataType                `ddl:"parameter,no_equals"`
	ResultDataType      datatypes.DataType      `ddl:"parameter,no_quotes,no_equals"`
	NotNull             *bool                   `ddl:"keyword" sql:"NOT NULL"`
	languageJavascript  bool                    `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior   *NullInputBehavior      `ddl:"keyword"`
	ProcedureDefinition string                  `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
	WithClauses         []ProcedureWithClause   `ddl:"keyword"`
	call                bool                    `ddl:"static" sql:"CALL"`
	ProcedureName       AccountObjectIdentifier `ddl:"identifier"`
	CallArguments       []string                `ddl:"keyword,must_parentheses"`
	ScriptingVariable   *string                 `ddl:"parameter,no_quotes,no_equals" sql:"INTO"`
}

// CreateAndCallForPythonProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/call-with#python.
type CreateAndCallForPythonProcedureOptions struct {
	with                bool                    `ddl:"static" sql:"WITH"`
	Name                AccountObjectIdentifier `ddl:"identifier"`
	asProcedure         bool                    `ddl:"static" sql:"AS PROCEDURE"`
	Arguments           []ProcedureArgument     `ddl:"list,must_parentheses"`
	Returns             ProcedureReturns        `ddl:"keyword" sql:"RETURNS"`
	languagePython      bool                    `ddl:"static" sql:"LANGUAGE PYTHON"`
	NullInputBehavior   *NullInputBehavior      `ddl:"keyword"`
	RuntimeVersion      string                  `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Packages            []ProcedurePackage      `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Imports             []ProcedureImport       `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Handler             string                  `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ProcedureDefinition *string                 `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
	WithClauses         []ProcedureWithClause   `ddl:"keyword"`
	call                bool                    `ddl:"static" sql:"CALL"`
	ProcedureName       AccountObjectIdentifier `ddl:"identifier"`
	CallArguments       []string                `ddl:"keyword,must_parentheses"`
	ScriptingVariable   *string                 `ddl:"parameter,no_quotes,no_equals" sql:"INTO"`
}

// CreateAndCallForSQLProcedureOptions is based on https://docs.snowflake.com/en/sql-reference/sql/call-with#snowflake-scripting.
type CreateAndCallForSQLProcedureOptions struct {
	with                bool                    `ddl:"static" sql:"WITH"`
	Name                AccountObjectIdentifier `ddl:"identifier"`
	asProcedure         bool                    `ddl:"static" sql:"AS PROCEDURE"`
	Arguments           []ProcedureArgument     `ddl:"list,must_parentheses"`
	Returns             ProcedureReturns        `ddl:"keyword" sql:"RETURNS"`
	languageSql         bool                    `ddl:"static" sql:"LANGUAGE SQL"`
	NullInputBehavior   *NullInputBehavior      `ddl:"keyword"`
	ProcedureDefinition string                  `ddl:"parameter,single_quotes,no_equals" sql:"AS"`
	WithClauses         []ProcedureWithClause   `ddl:"keyword"`
	call                bool                    `ddl:"static" sql:"CALL"`
	ProcedureName       AccountObjectIdentifier `ddl:"identifier"`
	CallArguments       []string                `ddl:"keyword,must_parentheses"`
	ScriptingVariable   *string                 `ddl:"parameter,no_quotes,no_equals" sql:"INTO"`
}
