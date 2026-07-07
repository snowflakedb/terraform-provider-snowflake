package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	userParametersSchema     = make(map[string]*schema.Schema)
	userParametersCustomDiff = ParametersCustomDiff(
		userParametersProvider,
		parameter[sdk.UserParameter]{sdk.UserParameterAbortDetachedQuery, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterAutocommit, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterBinaryInputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterBinaryOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientMemoryLimit, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientMetadataRequestUseConnectionCtx, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientPrefetchThreads, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientResultChunkSize, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientResultColumnCaseInsensitive, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientSessionKeepAlive, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterClientTimestampTypeMapping, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterDateInputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterDateOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterEnableUnloadPhysicalTypeOptimization, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterErrorOnNondeterministicMerge, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterErrorOnNondeterministicUpdate, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterGeographyOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterGeometryOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterJdbcTreatDecimalAsInt, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterJdbcTreatTimestampNtzAsUtc, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterJdbcUseSessionTimezone, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterJsonIndent, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterLockTimeout, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterLogLevel, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterLogEventLevel, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterMultiStatementCount, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterNoorderSequenceAsDefault, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterOdbcTreatDecimalAsInt, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterQueryTag, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterRowsPerResultset, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterS3StageVpceDnsName, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterSearchPath, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterSimulatedDataSharingConsumer, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterStatementQueuedTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterStatementTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterStrictJsonOutput, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampDayIsAlways24h, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampInputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampLtzOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampNtzOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampTypeMapping, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimestampTzOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimezone, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimeInputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTimeOutputFormat, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTraceLevel, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTransactionAbortOnError, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTransactionDefaultIsolationLevel, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterTwoDigitCenturyStart, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterUnsupportedDdlAction, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterUseCachedResult, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterWeekOfYearPolicy, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterWeekStart, valueTypeInt, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterEnableUnredactedQuerySyntaxError, valueTypeBool, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterNetworkPolicy, valueTypeString, sdk.ParameterTypeUser},
		parameter[sdk.UserParameter]{sdk.UserParameterPreventUnloadToInternalStages, valueTypeBool, sdk.ParameterTypeUser},
	)
)

type parameterDef[T ~string] struct {
	Name          T
	Type          schema.ValueType
	Description   string
	DiffSuppress  schema.SchemaDiffSuppressFunc
	ValidateDiag  schema.SchemaValidateDiagFunc
	ConflictsWith []string
}

func init() {
	// TODO [SNOW-1645342]: move to the SDK
	userParameterFields := []parameterDef[sdk.UserParameter]{
		// session params
		{Name: sdk.UserParameterAbortDetachedQuery, Type: schema.TypeBool, Description: "Specifies the action that Snowflake performs for in-progress queries if connectivity is lost due to abrupt termination of a session (e.g. network outage, browser termination, service interruption)."},
		{Name: sdk.UserParameterAutocommit, Type: schema.TypeBool, Description: "Specifies whether autocommit is enabled for the session. Autocommit determines whether a DML statement, when executed without an active transaction, is automatically committed after the statement successfully completes. For more information, see [Transactions](https://docs.snowflake.com/en/sql-reference/transactions)."},
		{Name: sdk.UserParameterBinaryInputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToBinaryInputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToBinaryInputFormat), Description: "The format of VARCHAR values passed as input to VARCHAR-to-BINARY conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output)."},
		{Name: sdk.UserParameterBinaryOutputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToBinaryOutputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToBinaryOutputFormat), Description: "The format for VARCHAR values returned as output by BINARY-to-VARCHAR conversion functions. For more information, see [Binary input and output](https://docs.snowflake.com/en/sql-reference/binary-input-output)."},
		{Name: sdk.UserParameterClientMemoryLimit, Type: schema.TypeInt, Description: "Parameter that specifies the maximum amount of memory the JDBC driver or ODBC driver should use for the result set from queries (in MB)."},
		{Name: sdk.UserParameterClientMetadataRequestUseConnectionCtx, Type: schema.TypeBool, Description: "For specific ODBC functions and JDBC methods, this parameter can change the default search scope from all databases/schemas to the current database/schema. The narrower search typically returns fewer rows and executes more quickly."},
		{Name: sdk.UserParameterClientPrefetchThreads, Type: schema.TypeInt, Description: "Parameter that specifies the number of threads used by the client to pre-fetch large result sets. The driver will attempt to honor the parameter value, but defines the minimum and maximum values (depending on your system’s resources) to improve performance."},
		{Name: sdk.UserParameterClientResultChunkSize, Type: schema.TypeInt, Description: "Parameter that specifies the maximum size of each set (or chunk) of query results to download (in MB). The JDBC driver downloads query results in chunks."},
		{Name: sdk.UserParameterClientResultColumnCaseInsensitive, Type: schema.TypeBool, Description: "Parameter that indicates whether to match column name case-insensitively in ResultSet.get* methods in JDBC."},
		{Name: sdk.UserParameterClientSessionKeepAlive, Type: schema.TypeBool, Description: "Parameter that indicates whether to force a user to log in again after a period of inactivity in the session."},
		{Name: sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, Type: schema.TypeInt, Description: "Number of seconds in-between client attempts to update the token for the session."},
		{Name: sdk.UserParameterClientTimestampTypeMapping, Type: schema.TypeString, Description: "Specifies the [TIMESTAMP_* variation](https://docs.snowflake.com/en/sql-reference/data-types-datetime.html#label-datatypes-timestamp-variations) to use when binding timestamp variables for JDBC or ODBC applications that use the bind API to load data."},
		{Name: sdk.UserParameterDateInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.UserParameterDateOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the DATE data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.UserParameterEnableUnloadPhysicalTypeOptimization, Type: schema.TypeBool, Description: "Specifies whether to set the schema for unloaded Parquet files based on the logical column data types (i.e. the types in the unload SQL query or source table) or on the unloaded column values (i.e. the smallest data types and precision that support the values in the output columns of the unload SQL statement or source table)."},
		{Name: sdk.UserParameterErrorOnNondeterministicMerge, Type: schema.TypeBool, Description: "Specifies whether to return an error when the [MERGE](https://docs.snowflake.com/en/sql-reference/sql/merge) command is used to update or delete a target row that joins multiple source rows and the system cannot determine the action to perform on the target row."},
		{Name: sdk.UserParameterErrorOnNondeterministicUpdate, Type: schema.TypeBool, Description: "Specifies whether to return an error when the [UPDATE](https://docs.snowflake.com/en/sql-reference/sql/update) command is used to update a target row that joins multiple source rows and the system cannot determine the action to perform on the target row."},
		{Name: sdk.UserParameterGeographyOutputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToGeographyOutputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToGeographyOutputFormat), Description: "Display format for [GEOGRAPHY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geography)."},
		{Name: sdk.UserParameterGeometryOutputFormat, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToGeometryOutputFormat), DiffSuppress: NormalizeAndCompare(sdk.ToGeometryOutputFormat), Description: "Display format for [GEOMETRY values](https://docs.snowflake.com/en/sql-reference/data-types-geospatial.html#label-data-types-geometry)."},
		{Name: sdk.UserParameterJdbcTreatDecimalAsInt, Type: schema.TypeBool, Description: "Specifies how JDBC processes columns that have a scale of zero (0)."},
		{Name: sdk.UserParameterJdbcTreatTimestampNtzAsUtc, Type: schema.TypeBool, Description: "Specifies how JDBC processes TIMESTAMP_NTZ values."},
		{Name: sdk.UserParameterJdbcUseSessionTimezone, Type: schema.TypeBool, Description: "Specifies whether the JDBC Driver uses the time zone of the JVM or the time zone of the session (specified by the [TIMEZONE](https://docs.snowflake.com/en/sql-reference/parameters#label-timezone) parameter) for the getDate(), getTime(), and getTimestamp() methods of the ResultSet class."},
		{Name: sdk.UserParameterJsonIndent, Type: schema.TypeInt, Description: "Specifies the number of blank spaces to indent each new element in JSON output in the session. Also specifies whether to insert newline characters after each element."},
		{Name: sdk.UserParameterLockTimeout, Type: schema.TypeInt, Description: "Number of seconds to wait while trying to lock a resource, before timing out and aborting the statement."},
		{Name: sdk.UserParameterLogLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToLogLevel), DiffSuppress: NormalizeAndCompare(sdk.ToLogLevel), Description: "Specifies the severity level of messages that should be ingested and made available in the active event table. Messages at the specified level (and at more severe levels) are ingested. For more information about log levels, see [Setting log level](https://docs.snowflake.com/en/developer-guide/logging-tracing/logging-log-level)."},
		{Name: sdk.UserParameterLogEventLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToLogLevel), DiffSuppress: NormalizeAndCompare(sdk.ToLogLevel), Description: "Specifies the severity level of log events (rows with record type EVENT) that should be ingested and made available in the active event table. Log events at the specified level (and at more severe levels) are ingested. For more information, see [LOG_EVENT_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#log_event_level). " + enumValuesDescription(sdk.AllLogLevels)},
		{Name: sdk.UserParameterMultiStatementCount, Type: schema.TypeInt, Description: "Number of statements to execute when using the multi-statement capability."},
		{Name: sdk.UserParameterNoorderSequenceAsDefault, Type: schema.TypeBool, Description: "Specifies whether the ORDER or NOORDER property is set by default when you create a new sequence or add a new table column. The ORDER and NOORDER properties determine whether or not the values are generated for the sequence or auto-incremented column in [increasing or decreasing order](https://docs.snowflake.com/en/user-guide/querying-sequences.html#label-querying-sequences-increasing-values)."},
		{Name: sdk.UserParameterOdbcTreatDecimalAsInt, Type: schema.TypeBool, Description: "Specifies how ODBC processes columns that have a scale of zero (0)."},
		{Name: sdk.UserParameterQueryTag, Type: schema.TypeString, Description: "Optional string that can be used to tag queries and other SQL statements executed within a session. The tags are displayed in the output of the [QUERY_HISTORY, QUERY_HISTORY_BY_*](https://docs.snowflake.com/en/sql-reference/functions/query_history) functions."},
		{Name: sdk.UserParameterQuotedIdentifiersIgnoreCase, Type: schema.TypeBool, Description: "Specifies whether letters in double-quoted object identifiers are stored and resolved as uppercase letters. By default, Snowflake preserves the case of alphabetic characters when storing and resolving double-quoted identifiers (see [Identifier resolution](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing)). You can use this parameter in situations in which [third-party applications always use double quotes around identifiers](https://docs.snowflake.com/en/sql-reference/identifiers-syntax.html#label-identifier-casing-parameter)."},
		{Name: sdk.UserParameterRowsPerResultset, Type: schema.TypeInt, Description: "Specifies the maximum number of rows returned in a result set. A value of 0 specifies no maximum."},
		{Name: sdk.UserParameterS3StageVpceDnsName, Type: schema.TypeString, Description: "Specifies the DNS name of an Amazon S3 interface endpoint. Requests sent to the internal stage of an account via [AWS PrivateLink for Amazon S3](https://docs.aws.amazon.com/AmazonS3/latest/userguide/privatelink-interface-endpoints.html) use this endpoint to connect. For more information, see [Accessing Internal stages with dedicated interface endpoints](https://docs.snowflake.com/en/user-guide/private-internal-stages-aws.html#label-aws-privatelink-internal-stage-network-isolation)."},
		{Name: sdk.UserParameterSearchPath, Type: schema.TypeString, Description: "Specifies the path to search to resolve unqualified object names in queries. For more information, see [Name resolution in queries](https://docs.snowflake.com/en/sql-reference/name-resolution.html#label-object-name-resolution-search-path). Comma-separated list of identifiers. An identifier can be a fully or partially qualified schema name."},
		{Name: sdk.UserParameterSimulatedDataSharingConsumer, Type: schema.TypeString, Description: "Specifies the name of a consumer account to simulate for testing/validating shared data, particularly shared secure views. When this parameter is set in a session, shared views return rows as if executed in the specified consumer account rather than the provider account. For more information, see [Introduction to Secure Data Sharing](https://docs.snowflake.com/en/user-guide/data-sharing-intro) and [Working with shares](https://docs.snowflake.com/en/user-guide/data-sharing-provider)."},
		{Name: sdk.UserParameterStatementQueuedTimeoutInSeconds, Type: schema.TypeInt, Description: "Amount of time, in seconds, a SQL statement (query, DDL, DML, etc.) remains queued for a warehouse before it is canceled by the system. This parameter can be used in conjunction with the [MAX_CONCURRENCY_LEVEL](https://docs.snowflake.com/en/sql-reference/parameters#label-max-concurrency-level) parameter to ensure a warehouse is never backlogged."},
		{Name: sdk.UserParameterStatementTimeoutInSeconds, Type: schema.TypeInt, Description: "Amount of time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system."},
		{Name: sdk.UserParameterStrictJsonOutput, Type: schema.TypeBool, Description: "This parameter specifies whether JSON output in a session is compatible with the general standard (as described by [http://json.org](http://json.org)). By design, Snowflake allows JSON input that contains non-standard values; however, these non-standard values might result in Snowflake outputting JSON that is incompatible with other platforms and languages. This parameter, when enabled, ensures that Snowflake outputs valid/compatible JSON."},
		{Name: sdk.UserParameterTimestampDayIsAlways24h, Type: schema.TypeBool, Description: "Specifies whether the [DATEADD](https://docs.snowflake.com/en/sql-reference/functions/dateadd) function (and its aliases) always consider a day to be exactly 24 hours for expressions that span multiple days."},
		{Name: sdk.UserParameterTimestampInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported timestamp format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of timestamps stored in the system during the session)."},
		{Name: sdk.UserParameterTimestampLtzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_LTZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.UserParameterTimestampNtzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_NTZ data type."},
		{Name: sdk.UserParameterTimestampOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP data type alias. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.UserParameterTimestampTypeMapping, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToTimestampTypeMapping), DiffSuppress: NormalizeAndCompare(sdk.ToTimestampTypeMapping), Description: "Specifies the TIMESTAMP_* variation that the TIMESTAMP data type alias maps to."},
		{Name: sdk.UserParameterTimestampTzOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIMESTAMP_TZ data type. If no format is specified, defaults to [TIMESTAMP_OUTPUT_FORMAT](https://docs.snowflake.com/en/sql-reference/parameters#label-timestamp-output-format). For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.UserParameterTimezone, Type: schema.TypeString, Description: "Specifies the time zone for the session. You can specify a [time zone name](https://data.iana.org/time-zones/tzdb-2021a/zone1970.tab) or a [link name](https://data.iana.org/time-zones/tzdb-2021a/backward) from release 2021a of the [IANA Time Zone Database](https://www.iana.org/time-zones) (e.g. America/Los_Angeles, Europe/London, UTC, Etc/GMT, etc.)."},
		{Name: sdk.UserParameterTimeInputFormat, Type: schema.TypeString, Description: "Specifies the input format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output). Any valid, supported time format or AUTO (AUTO specifies that Snowflake attempts to automatically detect the format of times stored in the system during the session)."},
		{Name: sdk.UserParameterTimeOutputFormat, Type: schema.TypeString, Description: "Specifies the display format for the TIME data type. For more information, see [Date and time input and output formats](https://docs.snowflake.com/en/sql-reference/date-time-input-output)."},
		{Name: sdk.UserParameterTraceLevel, Type: schema.TypeString, ValidateDiag: sdkValidation(sdk.ToTraceLevel), DiffSuppress: NormalizeAndCompare(sdk.ToTraceLevel), Description: "Controls how trace events are ingested into the event table. For more information about trace levels, see [Setting trace level](https://docs.snowflake.com/en/developer-guide/logging-tracing/tracing-trace-level)."},
		{Name: sdk.UserParameterTransactionAbortOnError, Type: schema.TypeBool, Description: "Specifies the action to perform when a statement issued within a non-autocommit transaction returns with an error."},
		{Name: sdk.UserParameterTransactionDefaultIsolationLevel, Type: schema.TypeString, Description: "Specifies the isolation level for transactions in the user session."},
		{Name: sdk.UserParameterTwoDigitCenturyStart, Type: schema.TypeInt, Description: "Specifies the “century start” year for 2-digit years (i.e. the earliest year such dates can represent). This parameter prevents ambiguous dates when importing or converting data with the `YY` date format component (i.e. years represented as 2 digits)."},
		{Name: sdk.UserParameterUnsupportedDdlAction, Type: schema.TypeString, Description: "Determines if an unsupported (i.e. non-default) value specified for a constraint property returns an error."},
		{Name: sdk.UserParameterUseCachedResult, Type: schema.TypeBool, Description: "Specifies whether to reuse persisted query results, if available, when a matching query is submitted."},
		{Name: sdk.UserParameterWeekOfYearPolicy, Type: schema.TypeInt, Description: "Specifies how the weeks in a given year are computed. `0`: The semantics used are equivalent to the ISO semantics, in which a week belongs to a given year if at least 4 days of that week are in that year. `1`: January 1 is included in the first week of the year and December 31 is included in the last week of the year."},
		{Name: sdk.UserParameterWeekStart, Type: schema.TypeInt, Description: "Specifies the first day of the week (used by week-related date functions). `0`: Legacy Snowflake behavior is used (i.e. ISO-like semantics). `1` (Monday) to `7` (Sunday): All the week-related functions use weeks that start on the specified day of the week."},
		{Name: sdk.UserParameterEnableUnredactedQuerySyntaxError, Type: schema.TypeBool, Description: "Controls whether query text is redacted if a SQL query fails due to a syntax or parsing error. If `FALSE`, the content of a failed query is redacted in the views, pages, and functions that provide a query history. Only users with a role that is granted or inherits the AUDIT privilege can set the ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR parameter. When using the ALTER USER command to set the parameter to `TRUE` for a particular user, modify the user that you want to see the query text, not the user who executed the query (if those are different users)."},
		{Name: sdk.UserParameterNetworkPolicy, Type: schema.TypeString, ValidateDiag: IsValidIdentifier[sdk.AccountObjectIdentifier](), DiffSuppress: suppressIdentifierQuoting, Description: "Specifies the network policy to enforce for your account. Network policies enable restricting access to your account based on users’ IP address. For more details, see [Controlling network traffic with network policies](https://docs.snowflake.com/en/user-guide/network-policies). Any existing network policy (created using [CREATE NETWORK POLICY](https://docs.snowflake.com/en/sql-reference/sql/create-network-policy))."},
		{Name: sdk.UserParameterPreventUnloadToInternalStages, Type: schema.TypeBool, Description: "Specifies whether to prevent data unload operations to internal (Snowflake) stages using [COPY INTO <location>](https://docs.snowflake.com/en/sql-reference/sql/copy-into-location) statements."},
	}

	// TODO [SNOW-1645342]: extract this method after moving to SDK
	for _, field := range userParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		userParametersSchema[fieldName] = &schema.Schema{
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

func userParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), userParametersProviderFunc, sdk.ParseAccountObjectIdentifier)
}

func userParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.AccountObjectIdentifier] {
	return c.Users.ShowParameters
}

// TODO [SNOW-1645342]: make generic based on type definition
func handleUserParameterRead(d *schema.ResourceData, userParameters []*sdk.Parameter) error {
	for _, p := range userParameters {
		switch p.Key {
		case
			string(sdk.UserParameterClientMemoryLimit),
			string(sdk.UserParameterClientPrefetchThreads),
			string(sdk.UserParameterClientResultChunkSize),
			string(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency),
			string(sdk.UserParameterJsonIndent),
			string(sdk.UserParameterLockTimeout),
			string(sdk.UserParameterMultiStatementCount),
			string(sdk.UserParameterRowsPerResultset),
			string(sdk.UserParameterStatementQueuedTimeoutInSeconds),
			string(sdk.UserParameterStatementTimeoutInSeconds),
			string(sdk.UserParameterTwoDigitCenturyStart),
			string(sdk.UserParameterWeekOfYearPolicy),
			string(sdk.UserParameterWeekStart):
			value, err := strconv.Atoi(p.Value)
			if err != nil {
				return err
			}
			if err := d.Set(strings.ToLower(p.Key), value); err != nil {
				return err
			}
		case
			string(sdk.UserParameterBinaryInputFormat),
			string(sdk.UserParameterBinaryOutputFormat),
			string(sdk.UserParameterClientTimestampTypeMapping),
			string(sdk.UserParameterDateInputFormat),
			string(sdk.UserParameterDateOutputFormat),
			string(sdk.UserParameterGeographyOutputFormat),
			string(sdk.UserParameterGeometryOutputFormat),
			string(sdk.UserParameterLogLevel),
			string(sdk.UserParameterLogEventLevel),
			string(sdk.UserParameterQueryTag),
			string(sdk.UserParameterS3StageVpceDnsName),
			string(sdk.UserParameterSearchPath),
			string(sdk.UserParameterSimulatedDataSharingConsumer),
			string(sdk.UserParameterTimestampInputFormat),
			string(sdk.UserParameterTimestampLtzOutputFormat),
			string(sdk.UserParameterTimestampNtzOutputFormat),
			string(sdk.UserParameterTimestampOutputFormat),
			string(sdk.UserParameterTimestampTypeMapping),
			string(sdk.UserParameterTimestampTzOutputFormat),
			string(sdk.UserParameterTimezone),
			string(sdk.UserParameterTimeInputFormat),
			string(sdk.UserParameterTimeOutputFormat),
			string(sdk.UserParameterTraceLevel),
			string(sdk.UserParameterTransactionDefaultIsolationLevel),
			string(sdk.UserParameterUnsupportedDdlAction),
			string(sdk.UserParameterNetworkPolicy):
			if err := d.Set(strings.ToLower(p.Key), p.Value); err != nil {
				return err
			}
		case
			string(sdk.UserParameterAbortDetachedQuery),
			string(sdk.UserParameterAutocommit),
			string(sdk.UserParameterClientMetadataRequestUseConnectionCtx),
			string(sdk.UserParameterClientResultColumnCaseInsensitive),
			string(sdk.UserParameterClientSessionKeepAlive),
			string(sdk.UserParameterEnableUnloadPhysicalTypeOptimization),
			string(sdk.UserParameterErrorOnNondeterministicMerge),
			string(sdk.UserParameterErrorOnNondeterministicUpdate),
			string(sdk.UserParameterJdbcTreatDecimalAsInt),
			string(sdk.UserParameterJdbcTreatTimestampNtzAsUtc),
			string(sdk.UserParameterJdbcUseSessionTimezone),
			string(sdk.UserParameterNoorderSequenceAsDefault),
			string(sdk.UserParameterOdbcTreatDecimalAsInt),
			string(sdk.UserParameterQuotedIdentifiersIgnoreCase),
			string(sdk.UserParameterStrictJsonOutput),
			string(sdk.UserParameterTimestampDayIsAlways24h),
			string(sdk.UserParameterTransactionAbortOnError),
			string(sdk.UserParameterUseCachedResult),
			string(sdk.UserParameterEnableUnredactedQuerySyntaxError),
			string(sdk.UserParameterPreventUnloadToInternalStages):
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
func handleUserParametersCreate(d *schema.ResourceData, sessionParameters *sdk.SessionParameters, objectParameters *sdk.UserObjectParametersRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.UserParameterAbortDetachedQuery, &sessionParameters.AbortDetachedQuery),
		handleParameterCreate(d, sdk.UserParameterAutocommit, &sessionParameters.Autocommit),
		handleParameterCreateWithMapping(d, sdk.UserParameterBinaryInputFormat, &sessionParameters.BinaryInputFormat, stringToStringEnumProvider(sdk.ToBinaryInputFormat)),
		handleParameterCreateWithMapping(d, sdk.UserParameterBinaryOutputFormat, &sessionParameters.BinaryOutputFormat, stringToStringEnumProvider(sdk.ToBinaryOutputFormat)),
		handleParameterCreate(d, sdk.UserParameterClientMemoryLimit, &sessionParameters.ClientMemoryLimit),
		handleParameterCreate(d, sdk.UserParameterClientMetadataRequestUseConnectionCtx, &sessionParameters.ClientMetadataRequestUseConnectionCtx),
		handleParameterCreate(d, sdk.UserParameterClientPrefetchThreads, &sessionParameters.ClientPrefetchThreads),
		handleParameterCreate(d, sdk.UserParameterClientResultChunkSize, &sessionParameters.ClientResultChunkSize),
		handleParameterCreate(d, sdk.UserParameterClientResultColumnCaseInsensitive, &sessionParameters.ClientResultColumnCaseInsensitive),
		handleParameterCreate(d, sdk.UserParameterClientSessionKeepAlive, &sessionParameters.ClientSessionKeepAlive),
		handleParameterCreate(d, sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, &sessionParameters.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterCreateWithMapping(d, sdk.UserParameterClientTimestampTypeMapping, &sessionParameters.ClientTimestampTypeMapping, stringToStringEnumProvider(sdk.ToClientTimestampTypeMapping)),
		handleParameterCreate(d, sdk.UserParameterDateInputFormat, &sessionParameters.DateInputFormat),
		handleParameterCreate(d, sdk.UserParameterDateOutputFormat, &sessionParameters.DateOutputFormat),
		handleParameterCreate(d, sdk.UserParameterEnableUnloadPhysicalTypeOptimization, &sessionParameters.EnableUnloadPhysicalTypeOptimization),
		handleParameterCreate(d, sdk.UserParameterErrorOnNondeterministicMerge, &sessionParameters.ErrorOnNondeterministicMerge),
		handleParameterCreate(d, sdk.UserParameterErrorOnNondeterministicUpdate, &sessionParameters.ErrorOnNondeterministicUpdate),
		handleParameterCreateWithMapping(d, sdk.UserParameterGeographyOutputFormat, &sessionParameters.GeographyOutputFormat, stringToStringEnumProvider(sdk.ToGeographyOutputFormat)),
		handleParameterCreateWithMapping(d, sdk.UserParameterGeometryOutputFormat, &sessionParameters.GeometryOutputFormat, stringToStringEnumProvider(sdk.ToGeometryOutputFormat)),
		handleParameterCreate(d, sdk.UserParameterJdbcTreatDecimalAsInt, &sessionParameters.JdbcTreatDecimalAsInt),
		handleParameterCreate(d, sdk.UserParameterJdbcTreatTimestampNtzAsUtc, &sessionParameters.JdbcTreatTimestampNtzAsUtc),
		handleParameterCreate(d, sdk.UserParameterJdbcUseSessionTimezone, &sessionParameters.JdbcUseSessionTimezone),
		handleParameterCreate(d, sdk.UserParameterJsonIndent, &sessionParameters.JsonIndent),
		handleParameterCreate(d, sdk.UserParameterLockTimeout, &sessionParameters.LockTimeout),
		handleParameterCreateWithMapping(d, sdk.UserParameterLogLevel, &sessionParameters.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreateWithMapping(d, sdk.UserParameterLogEventLevel, &sessionParameters.LogEventLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreate(d, sdk.UserParameterMultiStatementCount, &sessionParameters.MultiStatementCount),
		handleParameterCreate(d, sdk.UserParameterNoorderSequenceAsDefault, &sessionParameters.NoorderSequenceAsDefault),
		handleParameterCreate(d, sdk.UserParameterOdbcTreatDecimalAsInt, &sessionParameters.OdbcTreatDecimalAsInt),
		handleParameterCreate(d, sdk.UserParameterQueryTag, &sessionParameters.QueryTag),
		handleParameterCreate(d, sdk.UserParameterQuotedIdentifiersIgnoreCase, &sessionParameters.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.UserParameterRowsPerResultset, &sessionParameters.RowsPerResultset),
		handleParameterCreate(d, sdk.UserParameterS3StageVpceDnsName, &sessionParameters.S3StageVpceDnsName),
		handleParameterCreate(d, sdk.UserParameterSearchPath, &sessionParameters.SearchPath),
		handleParameterCreate(d, sdk.UserParameterSimulatedDataSharingConsumer, &sessionParameters.SimulatedDataSharingConsumer),
		handleParameterCreate(d, sdk.UserParameterStatementQueuedTimeoutInSeconds, &sessionParameters.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, sdk.UserParameterStatementTimeoutInSeconds, &sessionParameters.StatementTimeoutInSeconds),
		handleParameterCreate(d, sdk.UserParameterStrictJsonOutput, &sessionParameters.StrictJsonOutput),
		handleParameterCreate(d, sdk.UserParameterTimestampDayIsAlways24h, &sessionParameters.TimestampDayIsAlways24h),
		handleParameterCreate(d, sdk.UserParameterTimestampInputFormat, &sessionParameters.TimestampInputFormat),
		handleParameterCreate(d, sdk.UserParameterTimestampLtzOutputFormat, &sessionParameters.TimestampLTZOutputFormat),
		handleParameterCreate(d, sdk.UserParameterTimestampNtzOutputFormat, &sessionParameters.TimestampNTZOutputFormat),
		handleParameterCreate(d, sdk.UserParameterTimestampOutputFormat, &sessionParameters.TimestampOutputFormat),
		handleParameterCreateWithMapping(d, sdk.UserParameterTimestampTypeMapping, &sessionParameters.TimestampTypeMapping, stringToStringEnumProvider(sdk.ToTimestampTypeMapping)),
		handleParameterCreate(d, sdk.UserParameterTimestampTzOutputFormat, &sessionParameters.TimestampTZOutputFormat),
		handleParameterCreate(d, sdk.UserParameterTimezone, &sessionParameters.Timezone),
		handleParameterCreate(d, sdk.UserParameterTimeInputFormat, &sessionParameters.TimeInputFormat),
		handleParameterCreate(d, sdk.UserParameterTimeOutputFormat, &sessionParameters.TimeOutputFormat),
		handleParameterCreateWithMapping(d, sdk.UserParameterTraceLevel, &sessionParameters.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
		handleParameterCreate(d, sdk.UserParameterTransactionAbortOnError, &sessionParameters.TransactionAbortOnError),
		handleParameterCreateWithMapping(d, sdk.UserParameterTransactionDefaultIsolationLevel, &sessionParameters.TransactionDefaultIsolationLevel, stringToStringEnumProvider(sdk.ToTransactionDefaultIsolationLevel)),
		handleParameterCreate(d, sdk.UserParameterTwoDigitCenturyStart, &sessionParameters.TwoDigitCenturyStart),
		handleParameterCreateWithMapping(d, sdk.UserParameterUnsupportedDdlAction, &sessionParameters.UnsupportedDDLAction, stringToStringEnumProvider(sdk.ToUnsupportedDDLAction)),
		handleParameterCreate(d, sdk.UserParameterUseCachedResult, &sessionParameters.UseCachedResult),
		handleParameterCreate(d, sdk.UserParameterWeekOfYearPolicy, &sessionParameters.WeekOfYearPolicy),
		handleParameterCreate(d, sdk.UserParameterWeekStart, &sessionParameters.WeekStart),
		handleParameterCreate(d, sdk.UserParameterEnableUnredactedQuerySyntaxError, &objectParameters.EnableUnredactedQuerySyntaxError),
		handleParameterCreateWithMapping(d, sdk.UserParameterNetworkPolicy, &objectParameters.NetworkPolicy, stringToAccountObjectIdentifier),
		handleParameterCreate(d, sdk.UserParameterPreventUnloadToInternalStages, &objectParameters.PreventUnloadToInternalStages),
	)
}

func handleUserParametersUpdate(d *schema.ResourceData, setSession *sdk.SessionParameters, setObject *sdk.UserObjectParametersRequest, unsetSession *sdk.SessionParametersUnset, unsetObject *sdk.UserObjectParametersUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.UserParameterAbortDetachedQuery, &setSession.AbortDetachedQuery, &unsetSession.AbortDetachedQuery),
		handleParameterUpdate(d, sdk.UserParameterAutocommit, &setSession.Autocommit, &unsetSession.Autocommit),
		handleParameterUpdateWithMapping(d, sdk.UserParameterBinaryInputFormat, &setSession.BinaryInputFormat, &unsetSession.BinaryInputFormat, stringToStringEnumProvider(sdk.ToBinaryInputFormat)),
		handleParameterUpdateWithMapping(d, sdk.UserParameterBinaryOutputFormat, &setSession.BinaryOutputFormat, &unsetSession.BinaryOutputFormat, stringToStringEnumProvider(sdk.ToBinaryOutputFormat)),
		handleParameterUpdate(d, sdk.UserParameterClientMemoryLimit, &setSession.ClientMemoryLimit, &unsetSession.ClientMemoryLimit),
		handleParameterUpdate(d, sdk.UserParameterClientMetadataRequestUseConnectionCtx, &setSession.ClientMetadataRequestUseConnectionCtx, &unsetSession.ClientMetadataRequestUseConnectionCtx),
		handleParameterUpdate(d, sdk.UserParameterClientPrefetchThreads, &setSession.ClientPrefetchThreads, &unsetSession.ClientPrefetchThreads),
		handleParameterUpdate(d, sdk.UserParameterClientResultChunkSize, &setSession.ClientResultChunkSize, &unsetSession.ClientResultChunkSize),
		handleParameterUpdate(d, sdk.UserParameterClientResultColumnCaseInsensitive, &setSession.ClientResultColumnCaseInsensitive, &unsetSession.ClientResultColumnCaseInsensitive),
		handleParameterUpdate(d, sdk.UserParameterClientSessionKeepAlive, &setSession.ClientSessionKeepAlive, &unsetSession.ClientSessionKeepAlive),
		handleParameterUpdate(d, sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, &setSession.ClientSessionKeepAliveHeartbeatFrequency, &unsetSession.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterUpdateWithMapping(d, sdk.UserParameterClientTimestampTypeMapping, &setSession.ClientTimestampTypeMapping, &unsetSession.ClientTimestampTypeMapping, stringToStringEnumProvider(sdk.ToClientTimestampTypeMapping)),
		handleParameterUpdate(d, sdk.UserParameterDateInputFormat, &setSession.DateInputFormat, &unsetSession.DateInputFormat),
		handleParameterUpdate(d, sdk.UserParameterDateOutputFormat, &setSession.DateOutputFormat, &unsetSession.DateOutputFormat),
		handleParameterUpdate(d, sdk.UserParameterEnableUnloadPhysicalTypeOptimization, &setSession.EnableUnloadPhysicalTypeOptimization, &unsetSession.EnableUnloadPhysicalTypeOptimization),
		handleParameterUpdate(d, sdk.UserParameterErrorOnNondeterministicMerge, &setSession.ErrorOnNondeterministicMerge, &unsetSession.ErrorOnNondeterministicMerge),
		handleParameterUpdate(d, sdk.UserParameterErrorOnNondeterministicUpdate, &setSession.ErrorOnNondeterministicUpdate, &unsetSession.ErrorOnNondeterministicUpdate),
		handleParameterUpdateWithMapping(d, sdk.UserParameterGeographyOutputFormat, &setSession.GeographyOutputFormat, &unsetSession.GeographyOutputFormat, stringToStringEnumProvider(sdk.ToGeographyOutputFormat)),
		handleParameterUpdateWithMapping(d, sdk.UserParameterGeometryOutputFormat, &setSession.GeometryOutputFormat, &unsetSession.GeometryOutputFormat, stringToStringEnumProvider(sdk.ToGeometryOutputFormat)),
		handleParameterUpdate(d, sdk.UserParameterJdbcTreatDecimalAsInt, &setSession.JdbcTreatDecimalAsInt, &unsetSession.JdbcTreatDecimalAsInt),
		handleParameterUpdate(d, sdk.UserParameterJdbcTreatTimestampNtzAsUtc, &setSession.JdbcTreatTimestampNtzAsUtc, &unsetSession.JdbcTreatTimestampNtzAsUtc),
		handleParameterUpdate(d, sdk.UserParameterJdbcUseSessionTimezone, &setSession.JdbcUseSessionTimezone, &unsetSession.JdbcUseSessionTimezone),
		handleParameterUpdate(d, sdk.UserParameterJsonIndent, &setSession.JsonIndent, &unsetSession.JsonIndent),
		handleParameterUpdate(d, sdk.UserParameterLockTimeout, &setSession.LockTimeout, &unsetSession.LockTimeout),
		handleParameterUpdateWithMapping(d, sdk.UserParameterLogLevel, &setSession.LogLevel, &unsetSession.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdateWithMapping(d, sdk.UserParameterLogEventLevel, &setSession.LogEventLevel, &unsetSession.LogEventLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdate(d, sdk.UserParameterMultiStatementCount, &setSession.MultiStatementCount, &unsetSession.MultiStatementCount),
		handleParameterUpdate(d, sdk.UserParameterNoorderSequenceAsDefault, &setSession.NoorderSequenceAsDefault, &unsetSession.NoorderSequenceAsDefault),
		handleParameterUpdate(d, sdk.UserParameterOdbcTreatDecimalAsInt, &setSession.OdbcTreatDecimalAsInt, &unsetSession.OdbcTreatDecimalAsInt),
		handleParameterUpdate(d, sdk.UserParameterQueryTag, &setSession.QueryTag, &unsetSession.QueryTag),
		handleParameterUpdate(d, sdk.UserParameterQuotedIdentifiersIgnoreCase, &setSession.QuotedIdentifiersIgnoreCase, &unsetSession.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, sdk.UserParameterRowsPerResultset, &setSession.RowsPerResultset, &unsetSession.RowsPerResultset),
		handleParameterUpdate(d, sdk.UserParameterS3StageVpceDnsName, &setSession.S3StageVpceDnsName, &unsetSession.S3StageVpceDnsName),
		handleParameterUpdate(d, sdk.UserParameterSearchPath, &setSession.SearchPath, &unsetSession.SearchPath),
		handleParameterUpdate(d, sdk.UserParameterSimulatedDataSharingConsumer, &setSession.SimulatedDataSharingConsumer, &unsetSession.SimulatedDataSharingConsumer),
		handleParameterUpdate(d, sdk.UserParameterStatementQueuedTimeoutInSeconds, &setSession.StatementQueuedTimeoutInSeconds, &unsetSession.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, sdk.UserParameterStatementTimeoutInSeconds, &setSession.StatementTimeoutInSeconds, &unsetSession.StatementTimeoutInSeconds),
		handleParameterUpdate(d, sdk.UserParameterStrictJsonOutput, &setSession.StrictJsonOutput, &unsetSession.StrictJsonOutput),
		handleParameterUpdate(d, sdk.UserParameterTimestampDayIsAlways24h, &setSession.TimestampDayIsAlways24h, &unsetSession.TimestampDayIsAlways24h),
		handleParameterUpdate(d, sdk.UserParameterTimestampInputFormat, &setSession.TimestampInputFormat, &unsetSession.TimestampInputFormat),
		handleParameterUpdate(d, sdk.UserParameterTimestampLtzOutputFormat, &setSession.TimestampLTZOutputFormat, &unsetSession.TimestampLTZOutputFormat),
		handleParameterUpdate(d, sdk.UserParameterTimestampNtzOutputFormat, &setSession.TimestampNTZOutputFormat, &unsetSession.TimestampNTZOutputFormat),
		handleParameterUpdate(d, sdk.UserParameterTimestampOutputFormat, &setSession.TimestampOutputFormat, &unsetSession.TimestampOutputFormat),
		handleParameterUpdateWithMapping(d, sdk.UserParameterTimestampTypeMapping, &setSession.TimestampTypeMapping, &unsetSession.TimestampTypeMapping, stringToStringEnumProvider(sdk.ToTimestampTypeMapping)),
		handleParameterUpdate(d, sdk.UserParameterTimestampTzOutputFormat, &setSession.TimestampTZOutputFormat, &unsetSession.TimestampTZOutputFormat),
		handleParameterUpdate(d, sdk.UserParameterTimezone, &setSession.Timezone, &unsetSession.Timezone),
		handleParameterUpdate(d, sdk.UserParameterTimeInputFormat, &setSession.TimeInputFormat, &unsetSession.TimeInputFormat),
		handleParameterUpdate(d, sdk.UserParameterTimeOutputFormat, &setSession.TimeOutputFormat, &unsetSession.TimeOutputFormat),
		handleParameterUpdateWithMapping(d, sdk.UserParameterTraceLevel, &setSession.TraceLevel, &unsetSession.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
		handleParameterUpdate(d, sdk.UserParameterTransactionAbortOnError, &setSession.TransactionAbortOnError, &unsetSession.TransactionAbortOnError),
		handleParameterUpdateWithMapping(d, sdk.UserParameterTransactionDefaultIsolationLevel, &setSession.TransactionDefaultIsolationLevel, &unsetSession.TransactionDefaultIsolationLevel, stringToStringEnumProvider(sdk.ToTransactionDefaultIsolationLevel)),
		handleParameterUpdate(d, sdk.UserParameterTwoDigitCenturyStart, &setSession.TwoDigitCenturyStart, &unsetSession.TwoDigitCenturyStart),
		handleParameterUpdateWithMapping(d, sdk.UserParameterUnsupportedDdlAction, &setSession.UnsupportedDDLAction, &unsetSession.UnsupportedDDLAction, stringToStringEnumProvider(sdk.ToUnsupportedDDLAction)),
		handleParameterUpdate(d, sdk.UserParameterUseCachedResult, &setSession.UseCachedResult, &unsetSession.UseCachedResult),
		handleParameterUpdate(d, sdk.UserParameterWeekOfYearPolicy, &setSession.WeekOfYearPolicy, &unsetSession.WeekOfYearPolicy),
		handleParameterUpdate(d, sdk.UserParameterWeekStart, &setSession.WeekStart, &unsetSession.WeekStart),
		handleParameterUpdate(d, sdk.UserParameterEnableUnredactedQuerySyntaxError, &setObject.EnableUnredactedQuerySyntaxError, &unsetObject.EnableUnredactedQuerySyntaxError),
		handleParameterUpdateWithMapping(d, sdk.UserParameterNetworkPolicy, &setObject.NetworkPolicy, &unsetObject.NetworkPolicy, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, sdk.UserParameterPreventUnloadToInternalStages, &setObject.PreventUnloadToInternalStages, &unsetObject.PreventUnloadToInternalStages),
	)
}
