package main

import (
	"fmt"
	"strings"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	UserTypePerson        = "PERSON"
	UserTypeService       = "SERVICE"
	UserTypeLegacyService = "LEGACY_SERVICE"
)

func HandleUsers(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[UserCsvRow, UserRepresentation](config, csvInput, MapUserToModel)
}

func MapUserToModel(user UserRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	userType := strings.ToUpper(user.Type)

	switch userType {
	case UserTypeService:
		return MapToServiceUser(user)
	case UserTypeLegacyService:
		return MapToLegacyServiceUser(user)
	case UserTypePerson, "":
		return MapToUser(user)
	default:
		return nil, nil, fmt.Errorf("unsupported user type: %s", user.Type)
	}
}

func MapToUser(user UserRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	userId := sdk.NewAccountObjectIdentifier(user.Name)
	resourceId := ResourceId(resources.User, userId.FullyQualifiedName())
	resourceModel := model.User(resourceId, user.Name)

	// Basic attributes
	handleIfNotEmpty(user.Comment, resourceModel.WithComment)
	handleIfNotEmpty(user.LoginName, resourceModel.WithLoginName)
	handleIfNotEmpty(user.DisplayName, resourceModel.WithDisplayName)
	handleIfNotEmpty(user.FirstName, resourceModel.WithFirstName)
	handleIfNotEmpty(user.MiddleName, resourceModel.WithMiddleName)
	handleIfNotEmpty(user.LastName, resourceModel.WithLastName)
	handleIfNotEmpty(user.Email, resourceModel.WithEmail)
	handleIfNotEmpty(user.RsaPublicKey, resourceModel.WithRsaPublicKey)
	handleIfNotEmpty(user.RsaPublicKey2, resourceModel.WithRsaPublicKey2)
	handleIfNotEmpty(user.DefaultWarehouse, resourceModel.WithDefaultWarehouse)
	handleIfNotEmpty(user.DefaultNamespace, resourceModel.WithDefaultNamespace)
	handleIfNotEmpty(user.DefaultRole, resourceModel.WithDefaultRole)
	// Use handleBoolAsString for disabled and must_change_password to avoid Terraform default "default" value
	handleBoolAsString(user.Disabled, resourceModel.WithDisabled)
	handleBoolAsString(user.MustChangePassword, resourceModel.WithMustChangePassword)

	// Secondary roles
	secondaryRolesOption := user.GetSecondaryRolesOption()
	if secondaryRolesOption != sdk.SecondaryRolesOptionDefault {
		resourceModel.WithDefaultSecondaryRolesOption(string(secondaryRolesOption))
	}

	// Parameters
	handleOptionalFieldWithBuilder(user.AbortDetachedQuery, resourceModel.WithAbortDetachedQuery)
	handleOptionalFieldWithBuilder(user.Autocommit, resourceModel.WithAutocommit)
	handleOptionalFieldWithBuilder(user.BinaryInputFormat, resourceModel.WithBinaryInputFormat)
	handleOptionalFieldWithBuilder(user.BinaryOutputFormat, resourceModel.WithBinaryOutputFormat)
	handleOptionalFieldWithBuilder(user.ClientMemoryLimit, resourceModel.WithClientMemoryLimit)
	handleOptionalFieldWithBuilder(user.ClientMetadataRequestUseConnectionCtx, resourceModel.WithClientMetadataRequestUseConnectionCtx)
	handleOptionalFieldWithBuilder(user.ClientPrefetchThreads, resourceModel.WithClientPrefetchThreads)
	handleOptionalFieldWithBuilder(user.ClientResultChunkSize, resourceModel.WithClientResultChunkSize)
	handleOptionalFieldWithBuilder(user.ClientResultColumnCaseInsensitive, resourceModel.WithClientResultColumnCaseInsensitive)
	handleOptionalFieldWithBuilder(user.ClientSessionKeepAlive, resourceModel.WithClientSessionKeepAlive)
	handleOptionalFieldWithBuilder(user.ClientSessionKeepAliveHeartbeatFrequency, resourceModel.WithClientSessionKeepAliveHeartbeatFrequency)
	handleOptionalFieldWithBuilder(user.ClientTimestampTypeMapping, resourceModel.WithClientTimestampTypeMapping)
	handleOptionalFieldWithBuilder(user.DateInputFormat, resourceModel.WithDateInputFormat)
	handleOptionalFieldWithBuilder(user.DateOutputFormat, resourceModel.WithDateOutputFormat)
	handleOptionalFieldWithBuilder(user.EnableUnloadPhysicalTypeOptimization, resourceModel.WithEnableUnloadPhysicalTypeOptimization)
	handleOptionalFieldWithBuilder(user.EnableUnredactedQuerySyntaxError, resourceModel.WithEnableUnredactedQuerySyntaxError)
	handleOptionalFieldWithBuilder(user.ErrorOnNondeterministicMerge, resourceModel.WithErrorOnNondeterministicMerge)
	handleOptionalFieldWithBuilder(user.ErrorOnNondeterministicUpdate, resourceModel.WithErrorOnNondeterministicUpdate)
	handleOptionalFieldWithBuilder(user.GeographyOutputFormat, resourceModel.WithGeographyOutputFormat)
	handleOptionalFieldWithBuilder(user.GeometryOutputFormat, resourceModel.WithGeometryOutputFormat)
	handleOptionalFieldWithBuilder(user.JdbcTreatDecimalAsInt, resourceModel.WithJdbcTreatDecimalAsInt)
	handleOptionalFieldWithBuilder(user.JdbcTreatTimestampNtzAsUtc, resourceModel.WithJdbcTreatTimestampNtzAsUtc)
	handleOptionalFieldWithBuilder(user.JdbcUseSessionTimezone, resourceModel.WithJdbcUseSessionTimezone)
	handleOptionalFieldWithBuilder(user.JsonIndent, resourceModel.WithJsonIndent)
	handleOptionalFieldWithBuilder(user.LockTimeout, resourceModel.WithLockTimeout)
	handleOptionalFieldWithBuilder(user.LogLevel, resourceModel.WithLogLevel)
	handleOptionalFieldWithBuilder(user.MultiStatementCount, resourceModel.WithMultiStatementCount)
	handleOptionalFieldWithBuilder(user.NetworkPolicy, resourceModel.WithNetworkPolicy)
	handleOptionalFieldWithBuilder(user.NoorderSequenceAsDefault, resourceModel.WithNoorderSequenceAsDefault)
	handleOptionalFieldWithBuilder(user.OdbcTreatDecimalAsInt, resourceModel.WithOdbcTreatDecimalAsInt)
	handleOptionalFieldWithBuilder(user.PreventUnloadToInternalStages, resourceModel.WithPreventUnloadToInternalStages)
	handleOptionalFieldWithBuilder(user.QueryTag, resourceModel.WithQueryTag)
	handleOptionalFieldWithBuilder(user.QuotedIdentifiersIgnoreCase, resourceModel.WithQuotedIdentifiersIgnoreCase)
	handleOptionalFieldWithBuilder(user.RowsPerResultset, resourceModel.WithRowsPerResultset)
	handleOptionalFieldWithBuilder(user.S3StageVpceDnsName, resourceModel.WithS3StageVpceDnsName)
	handleOptionalFieldWithBuilder(user.SearchPath, resourceModel.WithSearchPath)
	handleOptionalFieldWithBuilder(user.SimulatedDataSharingConsumer, resourceModel.WithSimulatedDataSharingConsumer)
	handleOptionalFieldWithBuilder(user.StatementQueuedTimeoutInSeconds, resourceModel.WithStatementQueuedTimeoutInSeconds)
	handleOptionalFieldWithBuilder(user.StatementTimeoutInSeconds, resourceModel.WithStatementTimeoutInSeconds)
	handleOptionalFieldWithBuilder(user.StrictJsonOutput, resourceModel.WithStrictJsonOutput)
	handleOptionalFieldWithBuilder(user.TimeInputFormat, resourceModel.WithTimeInputFormat)
	handleOptionalFieldWithBuilder(user.TimeOutputFormat, resourceModel.WithTimeOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampDayIsAlways24h, resourceModel.WithTimestampDayIsAlways24h)
	handleOptionalFieldWithBuilder(user.TimestampInputFormat, resourceModel.WithTimestampInputFormat)
	handleOptionalFieldWithBuilder(user.TimestampLTZOutputFormat, resourceModel.WithTimestampLtzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampNTZOutputFormat, resourceModel.WithTimestampNtzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampOutputFormat, resourceModel.WithTimestampOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampTZOutputFormat, resourceModel.WithTimestampTzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampTypeMapping, resourceModel.WithTimestampTypeMapping)
	handleOptionalFieldWithBuilder(user.Timezone, resourceModel.WithTimezone)
	handleOptionalFieldWithBuilder(user.TraceLevel, resourceModel.WithTraceLevel)
	handleOptionalFieldWithBuilder(user.TransactionAbortOnError, resourceModel.WithTransactionAbortOnError)
	handleOptionalFieldWithBuilder(user.TransactionDefaultIsolationLevel, resourceModel.WithTransactionDefaultIsolationLevel)
	handleOptionalFieldWithBuilder(user.TwoDigitCenturyStart, resourceModel.WithTwoDigitCenturyStart)
	handleOptionalFieldWithBuilder(user.UnsupportedDDLAction, resourceModel.WithUnsupportedDdlAction)
	handleOptionalFieldWithBuilder(user.UseCachedResult, resourceModel.WithUseCachedResult)
	handleOptionalFieldWithBuilder(user.WeekOfYearPolicy, resourceModel.WithWeekOfYearPolicy)
	handleOptionalFieldWithBuilder(user.WeekStart, resourceModel.WithWeekStart)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		userId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}

func MapToServiceUser(user UserRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	userId := sdk.NewAccountObjectIdentifier(user.Name)
	resourceId := ResourceId(resources.ServiceUser, userId.FullyQualifiedName())
	resourceModel := model.ServiceUser(resourceId, user.Name)

	// Basic attributes (SERVICE users cannot have: FirstName, LastName, MiddleName, Password, MustChangePassword, MinsToBypassMfa)
	handleIfNotEmpty(user.Comment, resourceModel.WithComment)
	handleIfNotEmpty(user.LoginName, resourceModel.WithLoginName)
	handleIfNotEmpty(user.DisplayName, resourceModel.WithDisplayName)
	handleIfNotEmpty(user.Email, resourceModel.WithEmail)
	handleIfNotEmpty(user.DefaultWarehouse, resourceModel.WithDefaultWarehouse)
	handleIfNotEmpty(user.DefaultNamespace, resourceModel.WithDefaultNamespace)
	handleIfNotEmpty(user.DefaultRole, resourceModel.WithDefaultRole)
	// Use handleBoolAsString for disabled to avoid Terraform default "default" value
	handleBoolAsString(user.Disabled, resourceModel.WithDisabled)
	handleIfNotEmpty(user.RsaPublicKey, resourceModel.WithRsaPublicKey)
	handleIfNotEmpty(user.RsaPublicKey2, resourceModel.WithRsaPublicKey2)

	// Secondary roles
	secondaryRolesOption := user.GetSecondaryRolesOption()
	if secondaryRolesOption != sdk.SecondaryRolesOptionDefault {
		resourceModel.WithDefaultSecondaryRolesOption(string(secondaryRolesOption))
	}

	// Parameters
	handleOptionalFieldWithBuilder(user.AbortDetachedQuery, resourceModel.WithAbortDetachedQuery)
	handleOptionalFieldWithBuilder(user.Autocommit, resourceModel.WithAutocommit)
	handleOptionalFieldWithBuilder(user.BinaryInputFormat, resourceModel.WithBinaryInputFormat)
	handleOptionalFieldWithBuilder(user.BinaryOutputFormat, resourceModel.WithBinaryOutputFormat)
	handleOptionalFieldWithBuilder(user.ClientMemoryLimit, resourceModel.WithClientMemoryLimit)
	handleOptionalFieldWithBuilder(user.ClientMetadataRequestUseConnectionCtx, resourceModel.WithClientMetadataRequestUseConnectionCtx)
	handleOptionalFieldWithBuilder(user.ClientPrefetchThreads, resourceModel.WithClientPrefetchThreads)
	handleOptionalFieldWithBuilder(user.ClientResultChunkSize, resourceModel.WithClientResultChunkSize)
	handleOptionalFieldWithBuilder(user.ClientResultColumnCaseInsensitive, resourceModel.WithClientResultColumnCaseInsensitive)
	handleOptionalFieldWithBuilder(user.ClientSessionKeepAlive, resourceModel.WithClientSessionKeepAlive)
	handleOptionalFieldWithBuilder(user.ClientSessionKeepAliveHeartbeatFrequency, resourceModel.WithClientSessionKeepAliveHeartbeatFrequency)
	handleOptionalFieldWithBuilder(user.ClientTimestampTypeMapping, resourceModel.WithClientTimestampTypeMapping)
	handleOptionalFieldWithBuilder(user.DateInputFormat, resourceModel.WithDateInputFormat)
	handleOptionalFieldWithBuilder(user.DateOutputFormat, resourceModel.WithDateOutputFormat)
	handleOptionalFieldWithBuilder(user.EnableUnloadPhysicalTypeOptimization, resourceModel.WithEnableUnloadPhysicalTypeOptimization)
	handleOptionalFieldWithBuilder(user.EnableUnredactedQuerySyntaxError, resourceModel.WithEnableUnredactedQuerySyntaxError)
	handleOptionalFieldWithBuilder(user.ErrorOnNondeterministicMerge, resourceModel.WithErrorOnNondeterministicMerge)
	handleOptionalFieldWithBuilder(user.ErrorOnNondeterministicUpdate, resourceModel.WithErrorOnNondeterministicUpdate)
	handleOptionalFieldWithBuilder(user.GeographyOutputFormat, resourceModel.WithGeographyOutputFormat)
	handleOptionalFieldWithBuilder(user.GeometryOutputFormat, resourceModel.WithGeometryOutputFormat)
	handleOptionalFieldWithBuilder(user.JdbcTreatDecimalAsInt, resourceModel.WithJdbcTreatDecimalAsInt)
	handleOptionalFieldWithBuilder(user.JdbcTreatTimestampNtzAsUtc, resourceModel.WithJdbcTreatTimestampNtzAsUtc)
	handleOptionalFieldWithBuilder(user.JdbcUseSessionTimezone, resourceModel.WithJdbcUseSessionTimezone)
	handleOptionalFieldWithBuilder(user.JsonIndent, resourceModel.WithJsonIndent)
	handleOptionalFieldWithBuilder(user.LockTimeout, resourceModel.WithLockTimeout)
	handleOptionalFieldWithBuilder(user.LogLevel, resourceModel.WithLogLevel)
	handleOptionalFieldWithBuilder(user.MultiStatementCount, resourceModel.WithMultiStatementCount)
	handleOptionalFieldWithBuilder(user.NetworkPolicy, resourceModel.WithNetworkPolicy)
	handleOptionalFieldWithBuilder(user.NoorderSequenceAsDefault, resourceModel.WithNoorderSequenceAsDefault)
	handleOptionalFieldWithBuilder(user.OdbcTreatDecimalAsInt, resourceModel.WithOdbcTreatDecimalAsInt)
	handleOptionalFieldWithBuilder(user.PreventUnloadToInternalStages, resourceModel.WithPreventUnloadToInternalStages)
	handleOptionalFieldWithBuilder(user.QueryTag, resourceModel.WithQueryTag)
	handleOptionalFieldWithBuilder(user.QuotedIdentifiersIgnoreCase, resourceModel.WithQuotedIdentifiersIgnoreCase)
	handleOptionalFieldWithBuilder(user.RowsPerResultset, resourceModel.WithRowsPerResultset)
	handleOptionalFieldWithBuilder(user.S3StageVpceDnsName, resourceModel.WithS3StageVpceDnsName)
	handleOptionalFieldWithBuilder(user.SearchPath, resourceModel.WithSearchPath)
	handleOptionalFieldWithBuilder(user.SimulatedDataSharingConsumer, resourceModel.WithSimulatedDataSharingConsumer)
	handleOptionalFieldWithBuilder(user.StatementQueuedTimeoutInSeconds, resourceModel.WithStatementQueuedTimeoutInSeconds)
	handleOptionalFieldWithBuilder(user.StatementTimeoutInSeconds, resourceModel.WithStatementTimeoutInSeconds)
	handleOptionalFieldWithBuilder(user.StrictJsonOutput, resourceModel.WithStrictJsonOutput)
	handleOptionalFieldWithBuilder(user.TimeInputFormat, resourceModel.WithTimeInputFormat)
	handleOptionalFieldWithBuilder(user.TimeOutputFormat, resourceModel.WithTimeOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampDayIsAlways24h, resourceModel.WithTimestampDayIsAlways24h)
	handleOptionalFieldWithBuilder(user.TimestampInputFormat, resourceModel.WithTimestampInputFormat)
	handleOptionalFieldWithBuilder(user.TimestampLTZOutputFormat, resourceModel.WithTimestampLtzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampNTZOutputFormat, resourceModel.WithTimestampNtzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampOutputFormat, resourceModel.WithTimestampOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampTZOutputFormat, resourceModel.WithTimestampTzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampTypeMapping, resourceModel.WithTimestampTypeMapping)
	handleOptionalFieldWithBuilder(user.Timezone, resourceModel.WithTimezone)
	handleOptionalFieldWithBuilder(user.TraceLevel, resourceModel.WithTraceLevel)
	handleOptionalFieldWithBuilder(user.TransactionAbortOnError, resourceModel.WithTransactionAbortOnError)
	handleOptionalFieldWithBuilder(user.TransactionDefaultIsolationLevel, resourceModel.WithTransactionDefaultIsolationLevel)
	handleOptionalFieldWithBuilder(user.TwoDigitCenturyStart, resourceModel.WithTwoDigitCenturyStart)
	handleOptionalFieldWithBuilder(user.UnsupportedDDLAction, resourceModel.WithUnsupportedDdlAction)
	handleOptionalFieldWithBuilder(user.UseCachedResult, resourceModel.WithUseCachedResult)
	handleOptionalFieldWithBuilder(user.WeekOfYearPolicy, resourceModel.WithWeekOfYearPolicy)
	handleOptionalFieldWithBuilder(user.WeekStart, resourceModel.WithWeekStart)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		userId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}

func MapToLegacyServiceUser(user UserRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	userId := sdk.NewAccountObjectIdentifier(user.Name)
	resourceId := ResourceId(resources.LegacyServiceUser, userId.FullyQualifiedName())
	resourceModel := model.LegacyServiceUser(resourceId, user.Name)

	// Basic attributes (LEGACY_SERVICE users cannot have: FirstName, LastName, MiddleName, MinsToBypassMfa)
	// But CAN have: Password, MustChangePassword
	handleIfNotEmpty(user.Comment, resourceModel.WithComment)
	handleIfNotEmpty(user.LoginName, resourceModel.WithLoginName)
	handleIfNotEmpty(user.DisplayName, resourceModel.WithDisplayName)
	handleIfNotEmpty(user.Email, resourceModel.WithEmail)
	handleIfNotEmpty(user.DefaultWarehouse, resourceModel.WithDefaultWarehouse)
	handleIfNotEmpty(user.DefaultNamespace, resourceModel.WithDefaultNamespace)
	handleIfNotEmpty(user.DefaultRole, resourceModel.WithDefaultRole)
	// Use handleBoolAsString for disabled and must_change_password to avoid Terraform default "default" value
	handleBoolAsString(user.Disabled, resourceModel.WithDisabled)
	handleBoolAsString(user.MustChangePassword, resourceModel.WithMustChangePassword)
	handleIfNotEmpty(user.RsaPublicKey, resourceModel.WithRsaPublicKey)
	handleIfNotEmpty(user.RsaPublicKey2, resourceModel.WithRsaPublicKey2)

	// Secondary roles
	secondaryRolesOption := user.GetSecondaryRolesOption()
	if secondaryRolesOption != sdk.SecondaryRolesOptionDefault {
		resourceModel.WithDefaultSecondaryRolesOption(string(secondaryRolesOption))
	}

	// Parameters
	handleOptionalFieldWithBuilder(user.AbortDetachedQuery, resourceModel.WithAbortDetachedQuery)
	handleOptionalFieldWithBuilder(user.Autocommit, resourceModel.WithAutocommit)
	handleOptionalFieldWithBuilder(user.BinaryInputFormat, resourceModel.WithBinaryInputFormat)
	handleOptionalFieldWithBuilder(user.BinaryOutputFormat, resourceModel.WithBinaryOutputFormat)
	handleOptionalFieldWithBuilder(user.ClientMemoryLimit, resourceModel.WithClientMemoryLimit)
	handleOptionalFieldWithBuilder(user.ClientMetadataRequestUseConnectionCtx, resourceModel.WithClientMetadataRequestUseConnectionCtx)
	handleOptionalFieldWithBuilder(user.ClientPrefetchThreads, resourceModel.WithClientPrefetchThreads)
	handleOptionalFieldWithBuilder(user.ClientResultChunkSize, resourceModel.WithClientResultChunkSize)
	handleOptionalFieldWithBuilder(user.ClientResultColumnCaseInsensitive, resourceModel.WithClientResultColumnCaseInsensitive)
	handleOptionalFieldWithBuilder(user.ClientSessionKeepAlive, resourceModel.WithClientSessionKeepAlive)
	handleOptionalFieldWithBuilder(user.ClientSessionKeepAliveHeartbeatFrequency, resourceModel.WithClientSessionKeepAliveHeartbeatFrequency)
	handleOptionalFieldWithBuilder(user.ClientTimestampTypeMapping, resourceModel.WithClientTimestampTypeMapping)
	handleOptionalFieldWithBuilder(user.DateInputFormat, resourceModel.WithDateInputFormat)
	handleOptionalFieldWithBuilder(user.DateOutputFormat, resourceModel.WithDateOutputFormat)
	handleOptionalFieldWithBuilder(user.EnableUnloadPhysicalTypeOptimization, resourceModel.WithEnableUnloadPhysicalTypeOptimization)
	handleOptionalFieldWithBuilder(user.EnableUnredactedQuerySyntaxError, resourceModel.WithEnableUnredactedQuerySyntaxError)
	handleOptionalFieldWithBuilder(user.ErrorOnNondeterministicMerge, resourceModel.WithErrorOnNondeterministicMerge)
	handleOptionalFieldWithBuilder(user.ErrorOnNondeterministicUpdate, resourceModel.WithErrorOnNondeterministicUpdate)
	handleOptionalFieldWithBuilder(user.GeographyOutputFormat, resourceModel.WithGeographyOutputFormat)
	handleOptionalFieldWithBuilder(user.GeometryOutputFormat, resourceModel.WithGeometryOutputFormat)
	handleOptionalFieldWithBuilder(user.JdbcTreatDecimalAsInt, resourceModel.WithJdbcTreatDecimalAsInt)
	handleOptionalFieldWithBuilder(user.JdbcTreatTimestampNtzAsUtc, resourceModel.WithJdbcTreatTimestampNtzAsUtc)
	handleOptionalFieldWithBuilder(user.JdbcUseSessionTimezone, resourceModel.WithJdbcUseSessionTimezone)
	handleOptionalFieldWithBuilder(user.JsonIndent, resourceModel.WithJsonIndent)
	handleOptionalFieldWithBuilder(user.LockTimeout, resourceModel.WithLockTimeout)
	handleOptionalFieldWithBuilder(user.LogLevel, resourceModel.WithLogLevel)
	handleOptionalFieldWithBuilder(user.MultiStatementCount, resourceModel.WithMultiStatementCount)
	handleOptionalFieldWithBuilder(user.NetworkPolicy, resourceModel.WithNetworkPolicy)
	handleOptionalFieldWithBuilder(user.NoorderSequenceAsDefault, resourceModel.WithNoorderSequenceAsDefault)
	handleOptionalFieldWithBuilder(user.OdbcTreatDecimalAsInt, resourceModel.WithOdbcTreatDecimalAsInt)
	handleOptionalFieldWithBuilder(user.PreventUnloadToInternalStages, resourceModel.WithPreventUnloadToInternalStages)
	handleOptionalFieldWithBuilder(user.QueryTag, resourceModel.WithQueryTag)
	handleOptionalFieldWithBuilder(user.QuotedIdentifiersIgnoreCase, resourceModel.WithQuotedIdentifiersIgnoreCase)
	handleOptionalFieldWithBuilder(user.RowsPerResultset, resourceModel.WithRowsPerResultset)
	handleOptionalFieldWithBuilder(user.S3StageVpceDnsName, resourceModel.WithS3StageVpceDnsName)
	handleOptionalFieldWithBuilder(user.SearchPath, resourceModel.WithSearchPath)
	handleOptionalFieldWithBuilder(user.SimulatedDataSharingConsumer, resourceModel.WithSimulatedDataSharingConsumer)
	handleOptionalFieldWithBuilder(user.StatementQueuedTimeoutInSeconds, resourceModel.WithStatementQueuedTimeoutInSeconds)
	handleOptionalFieldWithBuilder(user.StatementTimeoutInSeconds, resourceModel.WithStatementTimeoutInSeconds)
	handleOptionalFieldWithBuilder(user.StrictJsonOutput, resourceModel.WithStrictJsonOutput)
	handleOptionalFieldWithBuilder(user.TimeInputFormat, resourceModel.WithTimeInputFormat)
	handleOptionalFieldWithBuilder(user.TimeOutputFormat, resourceModel.WithTimeOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampDayIsAlways24h, resourceModel.WithTimestampDayIsAlways24h)
	handleOptionalFieldWithBuilder(user.TimestampInputFormat, resourceModel.WithTimestampInputFormat)
	handleOptionalFieldWithBuilder(user.TimestampLTZOutputFormat, resourceModel.WithTimestampLtzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampNTZOutputFormat, resourceModel.WithTimestampNtzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampOutputFormat, resourceModel.WithTimestampOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampTZOutputFormat, resourceModel.WithTimestampTzOutputFormat)
	handleOptionalFieldWithBuilder(user.TimestampTypeMapping, resourceModel.WithTimestampTypeMapping)
	handleOptionalFieldWithBuilder(user.Timezone, resourceModel.WithTimezone)
	handleOptionalFieldWithBuilder(user.TraceLevel, resourceModel.WithTraceLevel)
	handleOptionalFieldWithBuilder(user.TransactionAbortOnError, resourceModel.WithTransactionAbortOnError)
	handleOptionalFieldWithBuilder(user.TransactionDefaultIsolationLevel, resourceModel.WithTransactionDefaultIsolationLevel)
	handleOptionalFieldWithBuilder(user.TwoDigitCenturyStart, resourceModel.WithTwoDigitCenturyStart)
	handleOptionalFieldWithBuilder(user.UnsupportedDDLAction, resourceModel.WithUnsupportedDdlAction)
	handleOptionalFieldWithBuilder(user.UseCachedResult, resourceModel.WithUseCachedResult)
	handleOptionalFieldWithBuilder(user.WeekOfYearPolicy, resourceModel.WithWeekOfYearPolicy)
	handleOptionalFieldWithBuilder(user.WeekStart, resourceModel.WithWeekStart)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		userId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
