package sdk

// imports edited manually
import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

type Functions interface {
	CreateForJava(ctx context.Context, request *CreateForJavaFunctionRequest) error
	CreateForJavascript(ctx context.Context, request *CreateForJavascriptFunctionRequest) error
	CreateForPython(ctx context.Context, request *CreateForPythonFunctionRequest) error
	CreateForScala(ctx context.Context, request *CreateForScalaFunctionRequest) error
	CreateForSQL(ctx context.Context, request *CreateForSQLFunctionRequest) error
	Alter(ctx context.Context, request *AlterFunctionRequest) error
	Drop(ctx context.Context, request *DropFunctionRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) error
	Show(ctx context.Context, request *ShowFunctionRequest) ([]Function, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Function, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*Function, error)
	Describe(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]FunctionDetail, error)

	// DescribeDetails is added manually; it returns aggregated describe results for the given function.
	DescribeDetails(ctx context.Context, id SchemaObjectIdentifierWithArguments) (*FunctionDetails, error)
	ShowParameters(ctx context.Context, id SchemaObjectIdentifierWithArguments) ([]*Parameter, error)
}

// CreateForJavaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#java-handler.
type CreateForJavaFunctionOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                     `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	function                   bool                      `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []FunctionArgument        `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    FunctionReturns           `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues           *ReturnNullValues         `ddl:"keyword"`
	languageJava               bool                      `ddl:"static" sql:"LANGUAGE JAVA"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             *string                   `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImport          `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackage         `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	EnableConsoleOutput        *bool                     `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *LogLevel                 `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel                *MetricLevel              `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel                 *TraceLevel               `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	FunctionDefinition         *string                   `ddl:"parameter,no_equals" sql:"AS"`
}

type FunctionArgument struct {
	ArgName        string             `ddl:"keyword,double_quotes"`
	ArgDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ArgDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
	DefaultValue   *string            `ddl:"parameter,no_equals" sql:"DEFAULT"`
}

type FunctionReturns struct {
	ResultDataType *FunctionReturnsResultDataType `ddl:"keyword"`
	Table          *FunctionReturnsTable          `ddl:"keyword" sql:"TABLE"`
}

type FunctionReturnsResultDataType struct {
	ResultDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ResultDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
}

type FunctionReturnsTable struct {
	Columns []FunctionColumn `ddl:"parameter,parentheses,no_equals"`
}

type FunctionColumn struct {
	ColumnName        string             `ddl:"keyword,double_quotes"`
	ColumnDataTypeOld DataType           `ddl:"keyword,no_quotes"`
	ColumnDataType    datatypes.DataType `ddl:"parameter,no_quotes,no_equals"`
}

type FunctionImport struct {
	Import string `ddl:"keyword,single_quotes"`
}

type FunctionPackage struct {
	Package string `ddl:"keyword,single_quotes"`
}

// CreateForJavascriptFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#javascript-handler.
type CreateForJavascriptFunctionOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	function              bool                   `ddl:"static" sql:"FUNCTION"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []FunctionArgument     `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               FunctionReturns        `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *ReturnNullValues      `ddl:"keyword"`
	languageJavascript    bool                   `ddl:"static" sql:"LANGUAGE JAVASCRIPT"`
	NullInputBehavior     *NullInputBehavior     `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	EnableConsoleOutput   *bool                  `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel              *LogLevel              `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel           *MetricLevel           `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel            *TraceLevel            `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	FunctionDefinition    string                 `ddl:"parameter,no_equals" sql:"AS"`
}

// CreateForPythonFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#python-handler.
type CreateForPythonFunctionOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                     `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	Aggregate                  *bool                     `ddl:"keyword" sql:"AGGREGATE"`
	function                   bool                      `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []FunctionArgument        `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	Returns                    FunctionReturns           `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues           *ReturnNullValues         `ddl:"keyword"`
	languagePython             bool                      `ddl:"static" sql:"LANGUAGE PYTHON"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImport          `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackage         `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	EnableConsoleOutput        *bool                     `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *LogLevel                 `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel                *MetricLevel              `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel                 *TraceLevel               `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	FunctionDefinition         *string                   `ddl:"parameter,no_equals" sql:"AS"`
}

// CreateForScalaFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#scala-handler.
type CreateForScalaFunctionOptions struct {
	create                     bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	Temporary                  *bool                     `ddl:"keyword" sql:"TEMPORARY"`
	Secure                     *bool                     `ddl:"keyword" sql:"SECURE"`
	function                   bool                      `ddl:"static" sql:"FUNCTION"`
	IfNotExists                *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier    `ddl:"identifier"`
	Arguments                  []FunctionArgument        `ddl:"list,must_parentheses"`
	CopyGrants                 *bool                     `ddl:"keyword" sql:"COPY GRANTS"`
	returns                    bool                      `ddl:"static" sql:"RETURNS"`
	ResultDataTypeOld          DataType                  `ddl:"parameter,no_equals"`
	ResultDataType             datatypes.DataType        `ddl:"parameter,no_quotes,no_equals"`
	ReturnNullValues           *ReturnNullValues         `ddl:"keyword"`
	languageScala              bool                      `ddl:"static" sql:"LANGUAGE SCALA"`
	NullInputBehavior          *NullInputBehavior        `ddl:"keyword"`
	ReturnResultsBehavior      *ReturnResultsBehavior    `ddl:"keyword"`
	RuntimeVersion             string                    `ddl:"parameter,single_quotes" sql:"RUNTIME_VERSION"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Imports                    []FunctionImport          `ddl:"parameter,parentheses" sql:"IMPORTS"`
	Packages                   []FunctionPackage         `ddl:"parameter,parentheses" sql:"PACKAGES"`
	Handler                    string                    `ddl:"parameter,single_quotes" sql:"HANDLER"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Secrets                    []SecretReference         `ddl:"parameter,parentheses" sql:"SECRETS"`
	TargetPath                 *string                   `ddl:"parameter,single_quotes" sql:"TARGET_PATH"`
	EnableConsoleOutput        *bool                     `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *LogLevel                 `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel                *MetricLevel              `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel                 *TraceLevel               `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	FunctionDefinition         *string                   `ddl:"parameter,no_equals" sql:"AS"`
}

// CreateForSQLFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-function#sql-handler.
type CreateForSQLFunctionOptions struct {
	create                bool                   `ddl:"static" sql:"CREATE"`
	OrReplace             *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	Temporary             *bool                  `ddl:"keyword" sql:"TEMPORARY"`
	Secure                *bool                  `ddl:"keyword" sql:"SECURE"`
	function              bool                   `ddl:"static" sql:"FUNCTION"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	Arguments             []FunctionArgument     `ddl:"list,must_parentheses"`
	CopyGrants            *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Returns               FunctionReturns        `ddl:"keyword" sql:"RETURNS"`
	ReturnNullValues      *ReturnNullValues      `ddl:"keyword"`
	ReturnResultsBehavior *ReturnResultsBehavior `ddl:"keyword"`
	Memoizable            *bool                  `ddl:"keyword" sql:"MEMOIZABLE"`
	Comment               *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	EnableConsoleOutput   *bool                  `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel              *LogLevel              `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel           *MetricLevel           `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel            *TraceLevel            `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
	FunctionDefinition    string                 `ddl:"parameter,no_equals" sql:"AS"`
}

// AlterFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-function.
type AlterFunctionOptions struct {
	alter       bool                                `ddl:"static" sql:"ALTER"`
	function    bool                                `ddl:"static" sql:"FUNCTION"`
	IfExists    *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	name        SchemaObjectIdentifierWithArguments `ddl:"identifier"`
	RenameTo    *SchemaObjectIdentifier             `ddl:"identifier" sql:"RENAME TO"`
	Set         *FunctionSet                        `ddl:"list" sql:"SET"`
	Unset       *FunctionUnset                      `ddl:"list" sql:"UNSET"`
	SetSecure   *bool                               `ddl:"keyword" sql:"SET SECURE"`
	UnsetSecure *bool                               `ddl:"keyword" sql:"UNSET SECURE"`
	SetTags     []TagAssociation                    `ddl:"keyword" sql:"SET TAG"`
	UnsetTags   []ObjectIdentifier                  `ddl:"keyword" sql:"UNSET TAG"`
}

type FunctionSet struct {
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	SecretsList                *SecretsList              `ddl:"parameter,parentheses" sql:"SECRETS"`
	EnableConsoleOutput        *bool                     `ddl:"parameter" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *LogLevel                 `ddl:"parameter,single_quotes" sql:"LOG_LEVEL"`
	MetricLevel                *MetricLevel              `ddl:"parameter,single_quotes" sql:"METRIC_LEVEL"`
	TraceLevel                 *TraceLevel               `ddl:"parameter,single_quotes" sql:"TRACE_LEVEL"`
}

type SecretsList struct {
	SecretsList []SecretReference `ddl:"list,must_parentheses"`
}

type FunctionUnset struct {
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
	ExternalAccessIntegrations *bool `ddl:"keyword" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	EnableConsoleOutput        *bool `ddl:"keyword" sql:"ENABLE_CONSOLE_OUTPUT"`
	LogLevel                   *bool `ddl:"keyword" sql:"LOG_LEVEL"`
	MetricLevel                *bool `ddl:"keyword" sql:"METRIC_LEVEL"`
	TraceLevel                 *bool `ddl:"keyword" sql:"TRACE_LEVEL"`
}

// DropFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-function.
type DropFunctionOptions struct {
	drop     bool                                `ddl:"static" sql:"DROP"`
	function bool                                `ddl:"static" sql:"FUNCTION"`
	IfExists *bool                               `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifierWithArguments `ddl:"identifier"`
}

// ShowFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-user-functions.
type ShowFunctionOptions struct {
	show          bool        `ddl:"static" sql:"SHOW"`
	userFunctions bool        `ddl:"static" sql:"USER FUNCTIONS"`
	Like          *Like       `ddl:"keyword" sql:"LIKE"`
	In            *ExtendedIn `ddl:"keyword" sql:"IN"`
}

type functionRow struct {
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
	IsExternalFunction         string         `db:"is_external_function"`
	Language                   string         `db:"language"`
	IsMemoizable               sql.NullString `db:"is_memoizable"`
	IsDataMetric               sql.NullString `db:"is_data_metric"`
}

type Function struct {
	CreatedOn                  string
	Name                       string
	SchemaName                 string
	IsBuiltin                  bool
	IsAggregate                bool
	IsAnsi                     bool
	MinNumArguments            int
	MaxNumArguments            int
	ArgumentsOld               []DataType
	ReturnTypeOld              DataType
	ArgumentsRaw               string
	Description                string
	CatalogName                string
	IsTableFunction            bool
	ValidForClustering         bool
	IsSecure                   bool
	Secrets                    *string
	ExternalAccessIntegrations *string
	IsExternalFunction         bool
	Language                   string
	IsMemoizable               bool
	IsDataMetric               bool
}

// DescribeFunctionOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-function.
type DescribeFunctionOptions struct {
	describe bool                                `ddl:"static" sql:"DESCRIBE"`
	function bool                                `ddl:"static" sql:"FUNCTION"`
	name     SchemaObjectIdentifierWithArguments `ddl:"identifier"`
}

type functionDetailRow struct {
	Property string         `db:"property"`
	Value    sql.NullString `db:"value"`
}

type FunctionDetail struct {
	Property string
	Value    *string
}
