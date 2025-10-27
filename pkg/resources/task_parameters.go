package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	taskParametersSchema     = make(map[string]*schema.Schema)
	taskParametersCustomDiff = ParametersCustomDiff(
		taskParametersProvider,
		// task parameters
		parameter[sdk.TaskParameter]{sdk.TaskParameterServerlessTaskMaxStatementSize, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterServerlessTaskMinStatementSize, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterSuspendTaskAfterNumFailures, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTaskAutoRetryAttempts, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterUserTaskManagedInitialWarehouseSize, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterUserTaskTimeoutMs, valueTypeInt, sdk.ParameterTypeTask},
		// session parameters
		parameter[sdk.TaskParameter]{sdk.TaskParameterAbortDetachedQuery, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterAutocommit, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterBinaryInputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterBinaryOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientMemoryLimit, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientMetadataRequestUseConnectionCtx, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientPrefetchThreads, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientResultChunkSize, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientResultColumnCaseInsensitive, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientSessionKeepAlive, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterClientTimestampTypeMapping, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterDateInputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterDateOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterEnableUnloadPhysicalTypeOptimization, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterErrorOnNondeterministicMerge, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterErrorOnNondeterministicUpdate, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterGeographyOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterGeometryOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterJdbcTreatTimestampNtzAsUtc, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterJdbcUseSessionTimezone, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterJsonIndent, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterLockTimeout, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterLogLevel, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterMultiStatementCount, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterNoorderSequenceAsDefault, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterOdbcTreatDecimalAsInt, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterQueryTag, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterRowsPerResultset, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterS3StageVpceDnsName, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterSearchPath, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterStatementQueuedTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterStatementTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterStrictJsonOutput, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampDayIsAlways24h, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampInputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampLtzOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampNtzOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampTypeMapping, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimestampTzOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimezone, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimeInputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTimeOutputFormat, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTraceLevel, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTransactionAbortOnError, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTransactionDefaultIsolationLevel, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterTwoDigitCenturyStart, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterUnsupportedDdlAction, valueTypeString, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterUseCachedResult, valueTypeBool, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterWeekOfYearPolicy, valueTypeInt, sdk.ParameterTypeTask},
		parameter[sdk.TaskParameter]{sdk.TaskParameterWeekStart, valueTypeInt, sdk.ParameterTypeTask},
	)
)

func init() {
	// TODO [SNOW-1645342]: move to the SDK
	TaskParameterFields := []parameterDef[sdk.TaskParameter]{
		// task parameters
		{Name: sdk.TaskParameterServerlessTaskMaxStatementSize, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToWarehouseSize), DiffSuppress: NormalizeAndCompare(sdk.ToWarehouseSize), Description: fmt.Sprintf("Specifies the maximum warehouse size for serverless tasks. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseSizesString))},
		{Name: sdk.TaskParameterServerlessTaskMinStatementSize, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToWarehouseSize), DiffSuppress: NormalizeAndCompare(sdk.ToWarehouseSize), Description: fmt.Sprintf("Specifies the minimum warehouse size for serverless tasks. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.ValidWarehouseSizesString))},
		{Name: sdk.TaskParameterSuspendTaskAfterNumFailures, Type: schema.TypeInt, ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)), Description: "Specifies the number of consecutive failed task runs after which the current task is suspended automatically. The default is 0 (no automatic suspension)."},
		{Name: sdk.TaskParameterTaskAutoRetryAttempts, Type: schema.TypeInt, ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)), Description: "Specifies the number of automatic task graph retry attempts. If any task graphs complete in a FAILED state, Snowflake can automatically retry the task graphs from the last task in the graph that failed."},
		{Name: sdk.TaskParameterUserTaskManagedInitialWarehouseSize, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToWarehouseSize), DiffSuppress: NormalizeAndCompare(sdk.ToWarehouseSize), ConflictsWith: []string{"warehouse"}, Description: "Specifies the size of the compute resources to provision for the first run of the task, before a task history is available for Snowflake to determine an ideal size. Once a task has successfully completed a few runs, Snowflake ignores this parameter setting. Valid values are (case-insensitive): %s. (Conflicts with warehouse). For more information about warehouses, see [docs](./warehouse)."},
		{Name: sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds, Type: schema.TypeInt, ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)), Description: "Minimum amount of time between Triggered Task executions in seconds"},
		{Name: sdk.TaskParameterUserTaskTimeoutMs, Type: schema.TypeInt, ValidateDiag: validation.ToDiagFunc(validation.IntAtLeast(0)), Description: "Specifies the time limit on a single run of the task before it times out (in milliseconds)."},
		// session params
		{Name: sdk.TaskParameterAbortDetachedQuery, Type: schema.TypeBool, Description: "Specifies the action that Snowflake performs for in-progress queries if connectivity is lost due to abrupt termination of a session (e.g. network outage, browser termination, service interruption)."},
		{Name: sdk.TaskParameterAutocommit, Type: schema.TypeBool, Description: "Specifies whether autocommit is enabled for the session. Autocommit determines whether a DML statement, when executed without an active transaction, is automatically committed after the statement successfully completes. For more information, see [Transactions](https://docs.snowflake.com/en/sql-reference/transactions)."},
		{Name: sdk.TaskParameterBinaryInputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToBinaryInputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToBinaryInputFormat), Description: "The format of VARCHAR values passed as input to VARCHAR-to-BINARY conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output)."},
		{Name: sdk.TaskParameterBinaryOutputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToBinaryOutputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToBinaryOutputFormat), Description: "The format for VARCHAR values returned as output by BINARY-to-VARCHAR conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output)."},
		{Name: sdk.TaskParameterClientMemoryLimit, Type: schema.TypeInt, Description: "Parameter that specifies the maximum amount of memory the JDBC driver or ODBC driver should use for the result set from queries (in MB)."},
		{Name: sdk.TaskParameterClientMetadataRequestUseConnectionCtx, Type: schema.TypeBool, Description: "For specific ODBC functions and JDBC methods, this parameter can change the default search scope from all databases/schemas to the current database/schema. The narrower search typically returns fewer rows and executes more quickly."},
		{Name: sdk.TaskParameterClientPrefetchThreads, Type: schema.TypeInt, Description: "Parameter that specifies the number of threads used by the client to pre-fetch large result sets. The driver will attempt to honor the parameter value, but defines the minimum and maximum values (depending on your system’s resources) to improve performance."},
		{Name: sdk.TaskParameterClientResultChunkSize, Type: schema.TypeInt, Description: "Parameter that specifies the maximum size of each set (or chunk) of query results to download (in MB). The JDBC driver downloads query results in chunks."},
		{Name: sdk.TaskParameterClientResultColumnCaseInsensitive, Type: schema.TypeBool, Description: "Parameter that indicates whether to match column name case-insensitively in ResultSet.get* methods in JDBC."},
		{Name: sdk.TaskParameterClientSessionKeepAlive, Type: schema.TypeBool, Description: "Parameter that indicates whether to force a user to log in again after a period of inactivity in the session."},
		{Name: sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency, Type: schema.TypeInt, Description: "Number of seconds in-between client attempts to update the token for the session."},
		{Name: sdk.TaskParameterClientTimestampTypeMapping, Type: schema.TypeString, Description: "Specifies the [TIMESTAMP_* variation](https://docs.snowflake.com/en/sql-reference/data-types-datetime.html#label-datatypes-timestamp-variations) to use when binding timestamp variables for JDBC or ODBC applications that use the bind API to load data."},
		{Name: sdk.TaskParameterDateInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.TaskParameterDateOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.TaskParameterEnableUnloadPhysicalTypeOptimization, Type: schema.TypeBool, Description: "Specifies whether to set the schema for unloaded Parquet files based on the logical column data types (i.e. the types in the unload SQL query or source table) or on the unloaded column values (i.e. the smallest data types and precision that support the values in the output columns of the unload SQL statement or source table)."},
		{Name: sdk.TaskParameterErrorOnNondeterministicMerge, Type: schema.TypeBool, Description: "Specifies whether to return an error when the [MERGE](https://docs.snowflake.com/en/sql-reference/sql/merge) command is used to update or delete a target row that joins multiple source rows and the system cannot determine the action to perform on the target row."},
		{Name: sdk.TaskParameterErrorOnNondeterministicUpdate, Type: schema.TypeBool, Description: "Specifies whether to return an error when the [UPDATE](https://docs.snowflake.com/en/sql-reference/sql/update) command is used to update a target row that joins multiple source rows and the system cannot determine the action to perform on the target row."},
		{Name: sdk.TaskParameterGeographyOutputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToGeographyOutputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToGeographyOutputFormat), Description: "Display format for [GEOGRAPHY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geography)."},
		{Name: sdk.TaskParameterGeometryOutputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToGeometryOutputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToGeometryOutputFormat), Description: "Display format for [GEOMETRY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geometry)."},
		{Name: sdk.TaskParameterJdbcTreatTimestampNtzAsUtc, Type: schema.TypeBool, Description: "Specifies how JDBC processes TIMESTAMP_NTZ values."},
		{Name: sdk.TaskParameterJdbcUseSessionTimezone, Type: schema.TypeBool, Description: "Specifies whether the JDBC Driver uses the time zone of the JVM or the time zone of the session (specified by the [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#label-timezone) parameter) for the getDate(), getTime(), and getTimestamp() methods of the ResultSet class."},
		{Name: sdk.TaskParameterJsonIndent, Type: schema.TypeInt, Description: "Specifies the number of blank spaces to indent each new element in JSON output in the session. Also specifies whether to insert newline characters after each element."},
		{Name: sdk.TaskParameterLockTimeout, Type: schema.TypeInt, Description: "Number of seconds to wait while trying to lock a resource, before timing out and aborting the statement."},
		{Name: sdk.TaskParameterLogLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToLogLevel), DiffSuppress: NormalizeAndCompare(sdk.ToLogLevel), Description: "Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested. For more information about log levels, see [Setting log level](https://docs.snowflake.com/en/developer-guide/logging-tracing/logging-log-level)."},
		{Name: sdk.TaskParameterMultiStatementCount, Type: schema.TypeInt, Description: "Number of statements to execute when using the multi-statement capability."},
		{Name: sdk.TaskParameterNoorderSequenceAsDefault, Type: schema.TypeBool, Description: "Specifies whether the ORDER or NOORDER property is set by default when you create a new sequence or add a new table column. The ORDER and NOORDER properties determine whether or not the values are generated for the sequence or auto-incremented column in [increasing or decreasing order](https://docs.snowflake.com/en/user-guide/querying-sequences.html#label-querying-sequences-increasing-values)."},
		{Name: sdk.TaskParameterOdbcTreatDecimalAsInt, Type: schema.TypeBool, Description: "Specifies how ODBC processes columns that have a scale of zero (0)."},
		{Name: sdk.TaskParameterQueryTag, Type: schema.TypeString, Description: "Optional string that can be used to tag queries and other SQL statements executed within a session. The tags are displayed in the output of the [QUERY_HISTORY, QUERY_HISTORY_BY_*](https://docs.snowflake.com/en/sql-reference/functions/query_history) functions."},
		{Name: sdk.TaskParameterQuotedIdentifiersIgnoreCase, Type: schema.TypeBool, Description: "Specifies whether letters in double-quoted object identifiers are stored and resolved as uppercase letters. By default, Snowflake preserves the case of alphabetic characters when storing and resolving double-quoted identifiers (see [Identifier resolution](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing)). You can use this parameter in situations in which [third-party applications always use double quotes around identifiers](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing-parameter)."},
		{Name: sdk.TaskParameterRowsPerResultset, Type: schema.TypeInt, Description: "Specifies the maximum number of rows returned in a result set. A value of 0 specifies no maximum."},
		{Name: sdk.TaskParameterS3StageVpceDnsName, Type: schema.TypeString, Description: "Specifies the DNS name of an Amazon S3 interface endpoint. Requests sent to the internal stage of an account via [AWS PrivateLink for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html) use this endpoint to connect. For more information, see [Accessing Internal stages with dedicated interface endpoints](https://docs.snowflake.com/en/user-guide/private-internal-stages-aws.html#label-aws-privatelink-internal-stage-network-isolation)."},
		{Name: sdk.TaskParameterSearchPath, Type: schema.TypeString, Description: "Specifies the path to search to resolve unqualified object names in queries. For more information, see [Name resolution in queries](https://docs.snowflake.com/en/sql-reference/name-resolution.html#label-object-name-resolution-search-path). Comma-separated list of identifiers. An identifier can be a fully or partially qualified schema name."},
		{Name: sdk.TaskParameterStatementQueuedTimeoutInSeconds, Type: schema.TypeInt, Description: "Amount of time, in seconds, a SQL statement (query, DDL, DML, etc.) remains queued for a warehouse before it is canceled by the system. This parameter can be used in conjunction with the [MAX_CONCURRENCY_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#label-max-concurrency-level) parameter to ensure a warehouse is never backlogged."},
		{Name: sdk.TaskParameterStatementTimeoutInSeconds, Type: schema.TypeInt, Description: "Amount of time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system."},
		{Name: sdk.TaskParameterStrictJsonOutput, Type: schema.TypeBool, Description: "This parameter specifies whether JSON output in a session is compatible with the general standard (as described by [http://json.org](http://json.org)). By design, Snowflake allows JSON input that contains non-standard values; however, these non-standard values might result in Snowflake outputting JSON that is incompatible with other platforms and languages. This parameter, when enabled, ensures that Snowflake outputs valid/compatible JSON."},
		{Name: sdk.TaskParameterTimestampDayIsAlways24h, Type: schema.TypeBool, Description: "Specifies whether the [DATEADD](https://docs.snowflake.com/en/sql-reference/functions/dateadd) function (and its aliases) always consider a day to be exactly 24 hours for expressions that span multiple days."},
		{Name: sdk.TaskParameterTimestampInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported timestamp format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of timestamps stored in the system during the session)."},
		{Name: sdk.TaskParameterTimestampLtzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_LTZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.TaskParameterTimestampNtzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_NTZ data type."},
		{Name: sdk.TaskParameterTimestampOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.TaskParameterTimestampTypeMapping, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToTimestampTypeMapping), DiffSuppress: NormalizeAndCompare(sdk.ToTimestampTypeMapping), Description: "Specifies the TIMESTAMP_* variation that the TIMESTAMP data type alias maps to."},
		{Name: sdk.TaskParameterTimestampTzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_TZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.TaskParameterTimezone, Type: schema.TypeString, Description: "Specifies the time zone for the session. You can specify a [time zone name](https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab) or a [link name](https://data.iana.org/time-zones/tzdb-2021a/backward) from release 2021a of the [IANA Time Zone Database](https://www.iana.org/time-zones) (e.g. America/Los_Angeles, Europe/London, UTC, Etc/GMT, etc.)."},
		{Name: sdk.TaskParameterTimeInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported time format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of times stored in the system during the session)."},
		{Name: sdk.TaskParameterTimeOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.TaskParameterTraceLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToTraceLevel), DiffSuppress: NormalizeAndCompare(sdk.ToTraceLevel), Description: "Controls how trace events are ingested into the event table. For more information about trace levels, see [Setting trace level](https://docs.snowflake.com/en/developer-guide/logging-tracing/tracing-trace-level)."},
		{Name: sdk.TaskParameterTransactionAbortOnError, Type: schema.TypeBool, Description: "Specifies the action to perform when a statement issued within a non-autocommit transaction returns with an error."},
		{Name: sdk.TaskParameterTransactionDefaultIsolationLevel, Type: schema.TypeString, Description: "Specifies the isolation level for transactions in the user session."},
		{Name: sdk.TaskParameterTwoDigitCenturyStart, Type: schema.TypeInt, Description: "Specifies the “century start” year for 2-digit years (i.e. the earliest year such dates can represent). This parameter prevents ambiguous dates when importing or converting data with the `YY` date format component (i.e. years represented as 2 digits)."},
		{Name: sdk.TaskParameterUnsupportedDdlAction, Type: schema.TypeString, Description: "Determines if an unsupported (i.e. non-default) value specified for a constraint property returns an error."},
		{Name: sdk.TaskParameterUseCachedResult, Type: schema.TypeBool, Description: "Specifies whether to reuse persisted query results, if available, when a matching query is submitted."},
		{Name: sdk.TaskParameterWeekOfYearPolicy, Type: schema.TypeInt, Description: "Specifies how the weeks in a given year are computed. `0`: The semantics used are equivalent to the ISO semantics, in which a week belongs to a given year if at least 4 days of that week are in that year. `1`: January 1 is included in the first week of the year and December 31 is included in the last week of the year."},
		{Name: sdk.TaskParameterWeekStart, Type: schema.TypeInt, Description: "Specifies the first day of the week (used by week-related date functions). `0`: Legacy Snowflake behavior is used (i.e. ISO-like semantics). `1` (Monday) to `7` (Sunday): All the week-related functions use weeks that start on the specified day of the week."},
	}

	// TODO [SNOW-1645342]: extract this method after moving to SDK
	for _, field := range TaskParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		taskParametersSchema[fieldName] = &schema.Schema{
			Type:             field.Type,
			Description:      enrichWithReferenceToParameterDocs(field.Name, field.Description),
			Computed:         true,
			Optional:         true,
			ValidateDiagFunc: field.ValidateDiag,
			DiffSuppressFunc: field.DiffSuppress,
			ConflictsWith:    field.ConflictsWith,
		}
	}
}

func taskParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), taskParametersProviderFunc, sdk.ParseSchemaObjectIdentifier)
}

func taskParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.SchemaObjectIdentifier] {
	return c.Tasks.ShowParameters
}

// TODO [SNOW-1645342]: make generic based on type definition
func handleTaskParameterRead(d *schema.ResourceData, taskParameters []*sdk.Parameter) error {
	for _, p := range taskParameters {
		switch p.Key {
		case
			string(sdk.TaskParameterSuspendTaskAfterNumFailures),
			string(sdk.TaskParameterTaskAutoRetryAttempts),
			string(sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds),
			string(sdk.TaskParameterUserTaskTimeoutMs),
			string(sdk.TaskParameterClientMemoryLimit),
			string(sdk.TaskParameterClientPrefetchThreads),
			string(sdk.TaskParameterClientResultChunkSize),
			string(sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency),
			string(sdk.TaskParameterJsonIndent),
			string(sdk.TaskParameterLockTimeout),
			string(sdk.TaskParameterMultiStatementCount),
			string(sdk.TaskParameterRowsPerResultset),
			string(sdk.TaskParameterStatementQueuedTimeoutInSeconds),
			string(sdk.TaskParameterStatementTimeoutInSeconds),
			string(sdk.TaskParameterTwoDigitCenturyStart),
			string(sdk.TaskParameterWeekOfYearPolicy),
			string(sdk.TaskParameterWeekStart):
			value, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			if err := d.Set(strings.ToLower(p.Key), value); err != nil {
				return err
			}
		case
			string(sdk.TaskParameterUserTaskManagedInitialWarehouseSize),
			string(sdk.TaskParameterBinaryInputFormat),
			string(sdk.TaskParameterBinaryOutputFormat),
			string(sdk.TaskParameterClientTimestampTypeMapping),
			string(sdk.TaskParameterDateInputFormat),
			string(sdk.TaskParameterDateOutputFormat),
			string(sdk.TaskParameterGeographyOutputFormat),
			string(sdk.TaskParameterGeometryOutputFormat),
			string(sdk.TaskParameterLogLevel),
			string(sdk.TaskParameterQueryTag),
			string(sdk.TaskParameterS3StageVpceDnsName),
			string(sdk.TaskParameterSearchPath),
			string(sdk.TaskParameterTimestampInputFormat),
			string(sdk.TaskParameterTimestampLtzOutputFormat),
			string(sdk.TaskParameterTimestampNtzOutputFormat),
			string(sdk.TaskParameterTimestampOutputFormat),
			string(sdk.TaskParameterTimestampTypeMapping),
			string(sdk.TaskParameterTimestampTzOutputFormat),
			string(sdk.TaskParameterTimezone),
			string(sdk.TaskParameterTimeInputFormat),
			string(sdk.TaskParameterTimeOutputFormat),
			string(sdk.TaskParameterTraceLevel),
			string(sdk.TaskParameterTransactionDefaultIsolationLevel),
			string(sdk.TaskParameterUnsupportedDdlAction):
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		case
			string(sdk.TaskParameterAbortDetachedQuery),
			string(sdk.TaskParameterAutocommit),
			string(sdk.TaskParameterClientMetadataRequestUseConnectionCtx),
			string(sdk.TaskParameterClientResultColumnCaseInsensitive),
			string(sdk.TaskParameterClientSessionKeepAlive),
			string(sdk.TaskParameterEnableUnloadPhysicalTypeOptimization),
			string(sdk.TaskParameterErrorOnNondeterministicMerge),
			string(sdk.TaskParameterErrorOnNondeterministicUpdate),
			string(sdk.TaskParameterJdbcTreatTimestampNtzAsUtc),
			string(sdk.TaskParameterJdbcUseSessionTimezone),
			string(sdk.TaskParameterNoorderSequenceAsDefault),
			string(sdk.TaskParameterOdbcTreatDecimalAsInt),
			string(sdk.TaskParameterQuotedIdentifiersIgnoreCase),
			string(sdk.TaskParameterStrictJsonOutput),
			string(sdk.TaskParameterTimestampDayIsAlways24h),
			string(sdk.TaskParameterTransactionAbortOnError),
			string(sdk.TaskParameterUseCachedResult):
			value, err := strconv.ParseBool(p.Value)
			if err != nil {
				return err
			}
			if err := d.Set(strings.ToLower(p.Key), value); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO [SNOW-1348330]: consider using SessionParameters#setParam during parameters rework
// (because currently setParam already is able to set the right parameter based on the string value input,
// but GetConfigPropertyAsPointerAllowingZeroValue receives typed value,
// so this would be unnecessary running in circles)
// TODO [SNOW-1645342]: include mappers in the param definition (after moving it to the SDK: identity versus concrete)
func handleTaskParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateTaskRequest) diag.Diagnostics {
	createOpts.WithSessionParameters(sdk.SessionParameters{})
	if v, ok := d.GetOk("user_task_managed_initial_warehouse_size"); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(size))
	}
	diags := JoinDiags(
		// task parameters
		handleParameterCreate(d, sdk.TaskParameterUserTaskTimeoutMs, &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, sdk.TaskParameterSuspendTaskAfterNumFailures, &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, sdk.TaskParameterTaskAutoRetryAttempts, &createOpts.TaskAutoRetryAttempts),
		handleParameterCreate(d, sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds, &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		// session parameters
		handleParameterCreate(d, sdk.TaskParameterAbortDetachedQuery, &createOpts.SessionParameters.AbortDetachedQuery),
		handleParameterCreate(d, sdk.TaskParameterAutocommit, &createOpts.SessionParameters.Autocommit),
		handleParameterCreateWithMapping(d, sdk.TaskParameterBinaryInputFormat, &createOpts.SessionParameters.BinaryInputFormat, stringToStringEnumProvider(sdk.ToBinaryInputFormat)),
		handleParameterCreateWithMapping(d, sdk.TaskParameterBinaryOutputFormat, &createOpts.SessionParameters.BinaryOutputFormat, stringToStringEnumProvider(sdk.ToBinaryOutputFormat)),
		handleParameterCreate(d, sdk.TaskParameterClientMemoryLimit, &createOpts.SessionParameters.ClientMemoryLimit),
		handleParameterCreate(d, sdk.TaskParameterClientMetadataRequestUseConnectionCtx, &createOpts.SessionParameters.ClientMetadataRequestUseConnectionCtx),
		handleParameterCreate(d, sdk.TaskParameterClientPrefetchThreads, &createOpts.SessionParameters.ClientPrefetchThreads),
		handleParameterCreate(d, sdk.TaskParameterClientResultChunkSize, &createOpts.SessionParameters.ClientResultChunkSize),
		handleParameterCreate(d, sdk.TaskParameterClientResultColumnCaseInsensitive, &createOpts.SessionParameters.ClientResultColumnCaseInsensitive),
		handleParameterCreate(d, sdk.TaskParameterClientSessionKeepAlive, &createOpts.SessionParameters.ClientSessionKeepAlive),
		handleParameterCreate(d, sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency, &createOpts.SessionParameters.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterCreateWithMapping(d, sdk.TaskParameterClientTimestampTypeMapping, &createOpts.SessionParameters.ClientTimestampTypeMapping, stringToStringEnumProvider(sdk.ToClientTimestampTypeMapping)),
		handleParameterCreate(d, sdk.TaskParameterDateInputFormat, &createOpts.SessionParameters.DateInputFormat),
		handleParameterCreate(d, sdk.TaskParameterDateOutputFormat, &createOpts.SessionParameters.DateOutputFormat),
		handleParameterCreate(d, sdk.TaskParameterEnableUnloadPhysicalTypeOptimization, &createOpts.SessionParameters.EnableUnloadPhysicalTypeOptimization),
		handleParameterCreate(d, sdk.TaskParameterErrorOnNondeterministicMerge, &createOpts.SessionParameters.ErrorOnNondeterministicMerge),
		handleParameterCreate(d, sdk.TaskParameterErrorOnNondeterministicUpdate, &createOpts.SessionParameters.ErrorOnNondeterministicUpdate),
		handleParameterCreateWithMapping(d, sdk.TaskParameterGeographyOutputFormat, &createOpts.SessionParameters.GeographyOutputFormat, stringToStringEnumProvider(sdk.ToGeographyOutputFormat)),
		handleParameterCreateWithMapping(d, sdk.TaskParameterGeometryOutputFormat, &createOpts.SessionParameters.GeometryOutputFormat, stringToStringEnumProvider(sdk.ToGeometryOutputFormat)),
		handleParameterCreate(d, sdk.TaskParameterJdbcTreatTimestampNtzAsUtc, &createOpts.SessionParameters.JdbcTreatTimestampNtzAsUtc),
		handleParameterCreate(d, sdk.TaskParameterJdbcUseSessionTimezone, &createOpts.SessionParameters.JdbcUseSessionTimezone),
		handleParameterCreate(d, sdk.TaskParameterJsonIndent, &createOpts.SessionParameters.JsonIndent),
		handleParameterCreate(d, sdk.TaskParameterLockTimeout, &createOpts.SessionParameters.LockTimeout),
		handleParameterCreateWithMapping(d, sdk.TaskParameterLogLevel, &createOpts.SessionParameters.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreate(d, sdk.TaskParameterMultiStatementCount, &createOpts.SessionParameters.MultiStatementCount),
		handleParameterCreate(d, sdk.TaskParameterNoorderSequenceAsDefault, &createOpts.SessionParameters.NoorderSequenceAsDefault),
		handleParameterCreate(d, sdk.TaskParameterOdbcTreatDecimalAsInt, &createOpts.SessionParameters.OdbcTreatDecimalAsInt),
		handleParameterCreate(d, sdk.TaskParameterQueryTag, &createOpts.SessionParameters.QueryTag),
		handleParameterCreate(d, sdk.TaskParameterQuotedIdentifiersIgnoreCase, &createOpts.SessionParameters.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.TaskParameterRowsPerResultset, &createOpts.SessionParameters.RowsPerResultset),
		handleParameterCreate(d, sdk.TaskParameterS3StageVpceDnsName, &createOpts.SessionParameters.S3StageVpceDnsName),
		handleParameterCreate(d, sdk.TaskParameterSearchPath, &createOpts.SessionParameters.SearchPath),
		handleParameterCreate(d, sdk.TaskParameterStatementQueuedTimeoutInSeconds, &createOpts.SessionParameters.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, sdk.TaskParameterStatementTimeoutInSeconds, &createOpts.SessionParameters.StatementTimeoutInSeconds),
		handleParameterCreate(d, sdk.TaskParameterStrictJsonOutput, &createOpts.SessionParameters.StrictJsonOutput),
		handleParameterCreate(d, sdk.TaskParameterTimestampDayIsAlways24h, &createOpts.SessionParameters.TimestampDayIsAlways24h),
		handleParameterCreate(d, sdk.TaskParameterTimestampInputFormat, &createOpts.SessionParameters.TimestampInputFormat),
		handleParameterCreate(d, sdk.TaskParameterTimestampLtzOutputFormat, &createOpts.SessionParameters.TimestampLTZOutputFormat),
		handleParameterCreate(d, sdk.TaskParameterTimestampNtzOutputFormat, &createOpts.SessionParameters.TimestampNTZOutputFormat),
		handleParameterCreate(d, sdk.TaskParameterTimestampOutputFormat, &createOpts.SessionParameters.TimestampOutputFormat),
		handleParameterCreateWithMapping(d, sdk.TaskParameterTimestampTypeMapping, &createOpts.SessionParameters.TimestampTypeMapping, stringToStringEnumProvider(sdk.ToTimestampTypeMapping)),
		handleParameterCreate(d, sdk.TaskParameterTimestampTzOutputFormat, &createOpts.SessionParameters.TimestampTZOutputFormat),
		handleParameterCreate(d, sdk.TaskParameterTimezone, &createOpts.SessionParameters.Timezone),
		handleParameterCreate(d, sdk.TaskParameterTimeInputFormat, &createOpts.SessionParameters.TimeInputFormat),
		handleParameterCreate(d, sdk.TaskParameterTimeOutputFormat, &createOpts.SessionParameters.TimeOutputFormat),
		handleParameterCreateWithMapping(d, sdk.TaskParameterTraceLevel, &createOpts.SessionParameters.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
		handleParameterCreate(d, sdk.TaskParameterTransactionAbortOnError, &createOpts.SessionParameters.TransactionAbortOnError),
		handleParameterCreateWithMapping(d, sdk.TaskParameterTransactionDefaultIsolationLevel, &createOpts.SessionParameters.TransactionDefaultIsolationLevel, stringToStringEnumProvider(sdk.ToTransactionDefaultIsolationLevel)),
		handleParameterCreate(d, sdk.TaskParameterTwoDigitCenturyStart, &createOpts.SessionParameters.TwoDigitCenturyStart),
		handleParameterCreateWithMapping(d, sdk.TaskParameterUnsupportedDdlAction, &createOpts.SessionParameters.UnsupportedDDLAction, stringToStringEnumProvider(sdk.ToUnsupportedDDLAction)),
		handleParameterCreate(d, sdk.TaskParameterUseCachedResult, &createOpts.SessionParameters.UseCachedResult),
		handleParameterCreate(d, sdk.TaskParameterWeekOfYearPolicy, &createOpts.SessionParameters.WeekOfYearPolicy),
		handleParameterCreate(d, sdk.TaskParameterWeekStart, &createOpts.SessionParameters.WeekStart),
	)
	if *createOpts.SessionParameters == (sdk.SessionParameters{}) {
		createOpts.SessionParameters = nil
	}
	return diags
}

func handleTaskParametersUpdate(d *schema.ResourceData, set *sdk.TaskSetRequest, unset *sdk.TaskUnsetRequest) diag.Diagnostics {
	set.WithSessionParameters(sdk.SessionParameters{})
	unset.WithSessionParametersUnset(sdk.SessionParametersUnset{})
	diags := JoinDiags(
		// task parameters
		handleParameterUpdateWithMapping(d, sdk.TaskParameterUserTaskManagedInitialWarehouseSize, &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, stringToStringEnumProvider(sdk.ToWarehouseSize)),
		handleParameterUpdate(d, sdk.TaskParameterUserTaskTimeoutMs, &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleParameterUpdate(d, sdk.TaskParameterSuspendTaskAfterNumFailures, &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, sdk.TaskParameterTaskAutoRetryAttempts, &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdate(d, sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds, &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		// session parameters
		handleParameterUpdate(d, sdk.TaskParameterAbortDetachedQuery, &set.SessionParameters.AbortDetachedQuery, &unset.SessionParametersUnset.AbortDetachedQuery),
		handleParameterUpdate(d, sdk.TaskParameterAutocommit, &set.SessionParameters.Autocommit, &unset.SessionParametersUnset.Autocommit),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterBinaryInputFormat, &set.SessionParameters.BinaryInputFormat, &unset.SessionParametersUnset.BinaryInputFormat, stringToStringEnumProvider(sdk.ToBinaryInputFormat)),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterBinaryOutputFormat, &set.SessionParameters.BinaryOutputFormat, &unset.SessionParametersUnset.BinaryOutputFormat, stringToStringEnumProvider(sdk.ToBinaryOutputFormat)),
		handleParameterUpdate(d, sdk.TaskParameterClientMemoryLimit, &set.SessionParameters.ClientMemoryLimit, &unset.SessionParametersUnset.ClientMemoryLimit),
		handleParameterUpdate(d, sdk.TaskParameterClientMetadataRequestUseConnectionCtx, &set.SessionParameters.ClientMetadataRequestUseConnectionCtx, &unset.SessionParametersUnset.ClientMetadataRequestUseConnectionCtx),
		handleParameterUpdate(d, sdk.TaskParameterClientPrefetchThreads, &set.SessionParameters.ClientPrefetchThreads, &unset.SessionParametersUnset.ClientPrefetchThreads),
		handleParameterUpdate(d, sdk.TaskParameterClientResultChunkSize, &set.SessionParameters.ClientResultChunkSize, &unset.SessionParametersUnset.ClientResultChunkSize),
		handleParameterUpdate(d, sdk.TaskParameterClientResultColumnCaseInsensitive, &set.SessionParameters.ClientResultColumnCaseInsensitive, &unset.SessionParametersUnset.ClientResultColumnCaseInsensitive),
		handleParameterUpdate(d, sdk.TaskParameterClientSessionKeepAlive, &set.SessionParameters.ClientSessionKeepAlive, &unset.SessionParametersUnset.ClientSessionKeepAlive),
		handleParameterUpdate(d, sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency, &set.SessionParameters.ClientSessionKeepAliveHeartbeatFrequency, &unset.SessionParametersUnset.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterClientTimestampTypeMapping, &set.SessionParameters.ClientTimestampTypeMapping, &unset.SessionParametersUnset.ClientTimestampTypeMapping, stringToStringEnumProvider(sdk.ToClientTimestampTypeMapping)),
		handleParameterUpdate(d, sdk.TaskParameterDateInputFormat, &set.SessionParameters.DateInputFormat, &unset.SessionParametersUnset.DateInputFormat),
		handleParameterUpdate(d, sdk.TaskParameterDateOutputFormat, &set.SessionParameters.DateOutputFormat, &unset.SessionParametersUnset.DateOutputFormat),
		handleParameterUpdate(d, sdk.TaskParameterEnableUnloadPhysicalTypeOptimization, &set.SessionParameters.EnableUnloadPhysicalTypeOptimization, &unset.SessionParametersUnset.EnableUnloadPhysicalTypeOptimization),
		handleParameterUpdate(d, sdk.TaskParameterErrorOnNondeterministicMerge, &set.SessionParameters.ErrorOnNondeterministicMerge, &unset.SessionParametersUnset.ErrorOnNondeterministicMerge),
		handleParameterUpdate(d, sdk.TaskParameterErrorOnNondeterministicUpdate, &set.SessionParameters.ErrorOnNondeterministicUpdate, &unset.SessionParametersUnset.ErrorOnNondeterministicUpdate),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterGeographyOutputFormat, &set.SessionParameters.GeographyOutputFormat, &unset.SessionParametersUnset.GeographyOutputFormat, stringToStringEnumProvider(sdk.ToGeographyOutputFormat)),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterGeometryOutputFormat, &set.SessionParameters.GeometryOutputFormat, &unset.SessionParametersUnset.GeometryOutputFormat, stringToStringEnumProvider(sdk.ToGeometryOutputFormat)),
		handleParameterUpdate(d, sdk.TaskParameterJdbcTreatTimestampNtzAsUtc, &set.SessionParameters.JdbcTreatTimestampNtzAsUtc, &unset.SessionParametersUnset.JdbcTreatTimestampNtzAsUtc),
		handleParameterUpdate(d, sdk.TaskParameterJdbcUseSessionTimezone, &set.SessionParameters.JdbcUseSessionTimezone, &unset.SessionParametersUnset.JdbcUseSessionTimezone),
		handleParameterUpdate(d, sdk.TaskParameterJsonIndent, &set.SessionParameters.JsonIndent, &unset.SessionParametersUnset.JsonIndent),
		handleParameterUpdate(d, sdk.TaskParameterLockTimeout, &set.SessionParameters.LockTimeout, &unset.SessionParametersUnset.LockTimeout),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterLogLevel, &set.SessionParameters.LogLevel, &unset.SessionParametersUnset.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdate(d, sdk.TaskParameterMultiStatementCount, &set.SessionParameters.MultiStatementCount, &unset.SessionParametersUnset.MultiStatementCount),
		handleParameterUpdate(d, sdk.TaskParameterNoorderSequenceAsDefault, &set.SessionParameters.NoorderSequenceAsDefault, &unset.SessionParametersUnset.NoorderSequenceAsDefault),
		handleParameterUpdate(d, sdk.TaskParameterOdbcTreatDecimalAsInt, &set.SessionParameters.OdbcTreatDecimalAsInt, &unset.SessionParametersUnset.OdbcTreatDecimalAsInt),
		handleParameterUpdate(d, sdk.TaskParameterQueryTag, &set.SessionParameters.QueryTag, &unset.SessionParametersUnset.QueryTag),
		handleParameterUpdate(d, sdk.TaskParameterQuotedIdentifiersIgnoreCase, &set.SessionParameters.QuotedIdentifiersIgnoreCase, &unset.SessionParametersUnset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, sdk.TaskParameterRowsPerResultset, &set.SessionParameters.RowsPerResultset, &unset.SessionParametersUnset.RowsPerResultset),
		handleParameterUpdate(d, sdk.TaskParameterS3StageVpceDnsName, &set.SessionParameters.S3StageVpceDnsName, &unset.SessionParametersUnset.S3StageVpceDnsName),
		handleParameterUpdate(d, sdk.TaskParameterSearchPath, &set.SessionParameters.SearchPath, &unset.SessionParametersUnset.SearchPath),
		handleParameterUpdate(d, sdk.TaskParameterStatementQueuedTimeoutInSeconds, &set.SessionParameters.StatementQueuedTimeoutInSeconds, &unset.SessionParametersUnset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, sdk.TaskParameterStatementTimeoutInSeconds, &set.SessionParameters.StatementTimeoutInSeconds, &unset.SessionParametersUnset.StatementTimeoutInSeconds),
		handleParameterUpdate(d, sdk.TaskParameterStrictJsonOutput, &set.SessionParameters.StrictJsonOutput, &unset.SessionParametersUnset.StrictJsonOutput),
		handleParameterUpdate(d, sdk.TaskParameterTimestampDayIsAlways24h, &set.SessionParameters.TimestampDayIsAlways24h, &unset.SessionParametersUnset.TimestampDayIsAlways24h),
		handleParameterUpdate(d, sdk.TaskParameterTimestampInputFormat, &set.SessionParameters.TimestampInputFormat, &unset.SessionParametersUnset.TimestampInputFormat),
		handleParameterUpdate(d, sdk.TaskParameterTimestampLtzOutputFormat, &set.SessionParameters.TimestampLTZOutputFormat, &unset.SessionParametersUnset.TimestampLTZOutputFormat),
		handleParameterUpdate(d, sdk.TaskParameterTimestampNtzOutputFormat, &set.SessionParameters.TimestampNTZOutputFormat, &unset.SessionParametersUnset.TimestampNTZOutputFormat),
		handleParameterUpdate(d, sdk.TaskParameterTimestampOutputFormat, &set.SessionParameters.TimestampOutputFormat, &unset.SessionParametersUnset.TimestampOutputFormat),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterTimestampTypeMapping, &set.SessionParameters.TimestampTypeMapping, &unset.SessionParametersUnset.TimestampTypeMapping, stringToStringEnumProvider(sdk.ToTimestampTypeMapping)),
		handleParameterUpdate(d, sdk.TaskParameterTimestampTzOutputFormat, &set.SessionParameters.TimestampTZOutputFormat, &unset.SessionParametersUnset.TimestampTZOutputFormat),
		handleParameterUpdate(d, sdk.TaskParameterTimezone, &set.SessionParameters.Timezone, &unset.SessionParametersUnset.Timezone),
		handleParameterUpdate(d, sdk.TaskParameterTimeInputFormat, &set.SessionParameters.TimeInputFormat, &unset.SessionParametersUnset.TimeInputFormat),
		handleParameterUpdate(d, sdk.TaskParameterTimeOutputFormat, &set.SessionParameters.TimeOutputFormat, &unset.SessionParametersUnset.TimeOutputFormat),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterTraceLevel, &set.SessionParameters.TraceLevel, &unset.SessionParametersUnset.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
		handleParameterUpdate(d, sdk.TaskParameterTransactionAbortOnError, &set.SessionParameters.TransactionAbortOnError, &unset.SessionParametersUnset.TransactionAbortOnError),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterTransactionDefaultIsolationLevel, &set.SessionParameters.TransactionDefaultIsolationLevel, &unset.SessionParametersUnset.TransactionDefaultIsolationLevel, stringToStringEnumProvider(sdk.ToTransactionDefaultIsolationLevel)),
		handleParameterUpdate(d, sdk.TaskParameterTwoDigitCenturyStart, &set.SessionParameters.TwoDigitCenturyStart, &unset.SessionParametersUnset.TwoDigitCenturyStart),
		handleParameterUpdateWithMapping(d, sdk.TaskParameterUnsupportedDdlAction, &set.SessionParameters.UnsupportedDDLAction, &unset.SessionParametersUnset.UnsupportedDDLAction, stringToStringEnumProvider(sdk.ToUnsupportedDDLAction)),
		handleParameterUpdate(d, sdk.TaskParameterUseCachedResult, &set.SessionParameters.UseCachedResult, &unset.SessionParametersUnset.UseCachedResult),
		handleParameterUpdate(d, sdk.TaskParameterWeekOfYearPolicy, &set.SessionParameters.WeekOfYearPolicy, &unset.SessionParametersUnset.WeekOfYearPolicy),
		handleParameterUpdate(d, sdk.TaskParameterWeekStart, &set.SessionParameters.WeekStart, &unset.SessionParametersUnset.WeekStart),
	)
	if *set.SessionParameters == (sdk.SessionParameters{}) {
		set.SessionParameters = nil
	}
	if *unset.SessionParametersUnset == (sdk.SessionParametersUnset{}) {
		unset.SessionParametersUnset = nil
	}
	return diags
}
