package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowAccountParametersSchema = make(map[string]*schema.Schema)
	accountParameters           = []sdk.AccountParameter{
		// TODO [SNOW-2298247]: Only session parameters are present here; add the missing ones while adding them to the accounts data source.
		// session parameters
		sdk.AccountParameterAbortDetachedQuery,
		sdk.AccountParameterAutocommit,
		sdk.AccountParameterBinaryInputFormat,
		sdk.AccountParameterBinaryOutputFormat,
		sdk.AccountParameterClientMetadataRequestUseConnectionCtx,
		sdk.AccountParameterClientResultColumnCaseInsensitive,
		sdk.AccountParameterDateInputFormat,
		sdk.AccountParameterDateOutputFormat,
		sdk.AccountParameterErrorOnNondeterministicMerge,
		sdk.AccountParameterErrorOnNondeterministicUpdate,
		sdk.AccountParameterGeographyOutputFormat,
		sdk.AccountParameterLockTimeout,
		sdk.AccountParameterLogLevel,
		sdk.AccountParameterMultiStatementCount,
		sdk.AccountParameterQueryTag,
		sdk.AccountParameterQuotedIdentifiersIgnoreCase,
		sdk.AccountParameterRowsPerResultset,
		sdk.AccountParameterS3StageVpceDnsName,
		sdk.AccountParameterStatementQueuedTimeoutInSeconds,
		sdk.AccountParameterStatementTimeoutInSeconds,
		sdk.AccountParameterTimestampDayIsAlways24h,
		sdk.AccountParameterTimestampInputFormat,
		sdk.AccountParameterTimestampLtzOutputFormat,
		sdk.AccountParameterTimestampNtzOutputFormat,
		sdk.AccountParameterTimestampOutputFormat,
		sdk.AccountParameterTimestampTypeMapping,
		sdk.AccountParameterTimestampTzOutputFormat,
		sdk.AccountParameterTimezone,
		sdk.AccountParameterTimeInputFormat,
		sdk.AccountParameterTimeOutputFormat,
		sdk.AccountParameterTraceLevel,
		sdk.AccountParameterTransactionAbortOnError,
		sdk.AccountParameterTransactionDefaultIsolationLevel,
		sdk.AccountParameterTwoDigitCenturyStart,
		sdk.AccountParameterUnsupportedDdlAction,
		sdk.AccountParameterUseCachedResult,
		sdk.AccountParameterWeekOfYearPolicy,
		sdk.AccountParameterWeekStart,
	}
)

func init() {
	for _, param := range accountParameters {
		ShowAccountParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func AccountParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	accountParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(accountParameters, sdk.AccountParameter(param.Key)) {
			accountParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return accountParametersValue
}
