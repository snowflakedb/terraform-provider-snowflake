package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type AccountCreateResponse struct {
	AccountLocator    string `json:"accountLocator,omitempty"`
	AccountLocatorUrl string `json:"accountLocatorUrl,omitempty"`
	OrganizationName  string
	AccountName       string         `json:"accountName,omitempty"`
	Url               string         `json:"url,omitempty"`
	Edition           AccountEdition `json:"edition,omitempty"`
	RegionGroup       string         `json:"regionGroup,omitempty"`
	Cloud             string         `json:"cloud,omitempty"`
	Region            string         `json:"region,omitempty"`
}

func ToAccountCreateResponse(v string) (*AccountCreateResponse, error) {
	var res AccountCreateResponse
	err := json.Unmarshal([]byte(v), &res)
	if err != nil {
		return nil, err
	}
	if len(res.Url) > 0 {
		url := strings.TrimPrefix(res.Url, `https://`)
		url = strings.TrimPrefix(url, `http://`)
		parts := strings.SplitN(url, "-", 2)
		if len(parts) == 2 {
			res.OrganizationName = strings.ToUpper(parts[0])
		}
	}
	return &res, nil
}

func (opts *AccountSet) additionalValidations() error {
	var errs []error
	if valueSet(opts.Force) && !valueSet(opts.PackagesPolicy) && !valueSet(opts.FeaturePolicySet) {
		errs = append(errs, NewError("force can only be set with PackagesPolicy and FeaturePolicy"))
	}
	if valueSet(opts.LegacyParameters) {
		if err := opts.LegacyParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *AccountUnset) additionalValidations() error {
	var errs []error
	if valueSet(opts.LegacyParameters) {
		if err := opts.LegacyParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *AlterAccountOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.Set) {
		if valueSet(opts.Set.ConsumptionBillingEntity) {
			if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
	}
	if valueSet(opts.Unset) {
		if valueSet(opts.Unset.ConsumptionBillingEntity) {
			if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
				errs = append(errs, ErrInvalidObjectIdentifier)
			}
		}
	}
	if valueSet(opts.Drop) || valueSet(opts.Rename) {
		if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

func (row accountDBRow) additionalConvert(result *Account) *Account {
	if row.ManagedAccounts.Valid {
		result.ManagedAccounts = Int(int(row.ManagedAccounts.Int32))
	}
	return result
}

func (v *Account) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.AccountName)
}

func (v *Account) AccountID() AccountIdentifier {
	return NewAccountIdentifier(v.OrganizationName, v.AccountName)
}

func (c *accounts) ShowParameters(ctx context.Context) ([]*Parameter, error) {
	return c.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Account: Bool(true),
		},
	})
}

func (c *accounts) UnsetAllParameters(ctx context.Context) error {
	return c.Alter(ctx, NewAlterAccountRequest().WithUnset(AccountUnset{
		Parameters: &AccountParametersUnset{
			AbortDetachedQuery:                                       Bool(true),
			ActivePythonProfiler:                                     Bool(true),
			AllowBindValuesAccess:                                    Bool(true),
			AllowClientMFACaching:                                    Bool(true),
			AllowedSpcsWorkloadTypes:                                 Bool(true),
			AllowIDToken:                                             Bool(true),
			Autocommit:                                               Bool(true),
			BaseLocationPrefix:                                       Bool(true),
			BinaryInputFormat:                                        Bool(true),
			BinaryOutputFormat:                                       Bool(true),
			Catalog:                                                  Bool(true),
			CatalogSync:                                              Bool(true),
			ClientEnableLogInfoStatementParameters:                   Bool(true),
			ClientEncryptionKeySize:                                  Bool(true),
			ClientMemoryLimit:                                        Bool(true),
			ClientMetadataRequestUseConnectionCtx:                    Bool(true),
			ClientMetadataUseSessionDatabase:                         Bool(true),
			ClientPrefetchThreads:                                    Bool(true),
			ClientResultChunkSize:                                    Bool(true),
			ClientResultColumnCaseInsensitive:                        Bool(true),
			ClientSessionKeepAlive:                                   Bool(true),
			ClientSessionKeepAliveHeartbeatFrequency:                 Bool(true),
			ClientTimestampTypeMapping:                               Bool(true),
			CortexCodeCliDailyEstCreditLimitPerUser:                  Bool(true),
			CortexCodeDesktopDailyEstCreditLimitPerUser:              Bool(true),
			CortexCodeSnowsightDailyEstCreditLimitPerUser:            Bool(true),
			CortexEnabledCrossRegion:                                 Bool(true),
			CortexModelsAllowlist:                                    Bool(true),
			CsvTimestampFormat:                                       Bool(true),
			DataMetricSchedule:                                       Bool(true),
			DataRetentionTimeInDays:                                  Bool(true),
			DateInputFormat:                                          Bool(true),
			DateOutputFormat:                                         Bool(true),
			DefaultDbtVersion:                                        Bool(true),
			DefaultDDLCollation:                                      Bool(true),
			DefaultNotebookComputePoolCpu:                            Bool(true),
			DefaultNotebookComputePoolGpu:                            Bool(true),
			DefaultNullOrdering:                                      Bool(true),
			DefaultStreamlitNotebookWarehouse:                        Bool(true),
			DisableUiDownloadButton:                                  Bool(true),
			DisallowedSpcsWorkloadTypes:                              Bool(true),
			DisableUserPrivilegeGrants:                               Bool(true),
			EnableAutomaticSensitiveDataClassificationLog:            Bool(true),
			EnableBudgetEventLogging:                                 Bool(true),
			EnableCortexAnalyst:                                      Bool(true),
			EnableDataCompaction:                                     Bool(true),
			EnableEgressCostOptimizer:                                Bool(true),
			EnableGetDdlUseDataTypeAlias:                             Bool(true),
			EnableIcebergMergeOnRead:                                 Bool(true),
			EnableNotebookCreationInPersonalDb:                       Bool(true),
			EnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement: Bool(true),
			EnableTagPropagationEventLogging:                         Bool(true),
			EnableIdentifierFirstLogin:                               Bool(true),
			EnableInternalStagesPrivatelink:                          Bool(true),
			EnableTriSecretAndRekeyOptOutForImageRepository:          Bool(true),
			EnableTriSecretAndRekeyOptOutForSpcsBlockStorage:         Bool(true),
			EnableUnhandledExceptionsReporting:                       Bool(true),
			EnableUnloadPhysicalTypeOptimization:                     Bool(true),
			EnableUnredactedQuerySyntaxError:                         Bool(true),
			EnableUnredactedSecureObjectError:                        Bool(true),
			EnforceNetworkRulesForInternalStages:                     Bool(true),
			ErrorOnNondeterministicMerge:                             Bool(true),
			ErrorOnNondeterministicUpdate:                            Bool(true),
			EventTable:                                               Bool(true),
			ExternalOAuthAddPrivilegedRolesToBlockedList:             Bool(true),
			ExternalVolume:                                           Bool(true),
			GeographyOutputFormat:                                    Bool(true),
			GeometryOutputFormat:                                     Bool(true),
			HybridTableLockTimeout:                                   Bool(true),
			IcebergVersionDefault:                                    Bool(true),
			InitialReplicationSizeLimitInTB:                          Bool(true),
			JdbcTreatDecimalAsInt:                                    Bool(true),
			JdbcTreatTimestampNtzAsUtc:                               Bool(true),
			JdbcUseSessionTimezone:                                   Bool(true),
			JsonIndent:                                               Bool(true),
			JsTreatIntegerAsBigInt:                                   Bool(true),
			ListingAutoFulfillmentReplicationRefreshSchedule:         Bool(true),
			LockTimeout:                                              Bool(true),
			LogLevel:                                                 Bool(true),
			LogEventLevel:                                            Bool(true),
			MaxConcurrencyLevel:                                      Bool(true),
			MaxDataExtensionTimeInDays:                               Bool(true),
			MetricLevel:                                              Bool(true),
			MinDataRetentionTimeInDays:                               Bool(true),
			MultiStatementCount:                                      Bool(true),
			NetworkPolicy:                                            Bool(true),
			NoorderSequenceAsDefault:                                 Bool(true),
			OAuthAddPrivilegedRolesToBlockedList:                     Bool(true),
			OdbcTreatDecimalAsInt:                                    Bool(true),
			PeriodicDataRekeying:                                     Bool(true),
			PipeExecutionPaused:                                      Bool(true),
			PreventUnloadToInlineURL:                                 Bool(true),
			PreventUnloadToInternalStages:                            Bool(true),
			PythonProfilerModules:                                    Bool(true),
			PythonProfilerTargetStage:                                Bool(true),
			QueryTag:                                                 Bool(true),
			QuotedIdentifiersIgnoreCase:                              Bool(true),
			ReadConsistencyMode:                                      Bool(true),
			ReplaceInvalidCharacters:                                 Bool(true),
			RequireStorageIntegrationForStageCreation:                Bool(true),
			RequireStorageIntegrationForStageOperation:               Bool(true),
			RowTimestampDefault:                                      Bool(true),
			RowsPerResultset:                                         Bool(true),
			S3StageVpceDnsName:                                       Bool(true),
			SearchPath:                                               Bool(true),
			ServerlessTaskMaxStatementSize:                           Bool(true),
			ServerlessTaskMinStatementSize:                           Bool(true),
			SimulatedDataSharingConsumer:                             Bool(true),
			SsoLoginPage:                                             Bool(true),
			SqlTraceQueryText:                                        Bool(true),
			StatementQueuedTimeoutInSeconds:                          Bool(true),
			StatementTimeoutInSeconds:                                Bool(true),
			StorageSerializationPolicy:                               Bool(true),
			StrictJsonOutput:                                         Bool(true),
			SuspendTaskAfterNumFailures:                              Bool(true),
			TaskAutoRetryAttempts:                                    Bool(true),
			TimestampDayIsAlways24h:                                  Bool(true),
			TimestampInputFormat:                                     Bool(true),
			TimestampLtzOutputFormat:                                 Bool(true),
			TimestampNtzOutputFormat:                                 Bool(true),
			TimestampOutputFormat:                                    Bool(true),
			TimestampTypeMapping:                                     Bool(true),
			TimestampTzOutputFormat:                                  Bool(true),
			Timezone:                                                 Bool(true),
			TimeInputFormat:                                          Bool(true),
			TimeOutputFormat:                                         Bool(true),
			TraceLevel:                                               Bool(true),
			TransactionAbortOnError:                                  Bool(true),
			TransactionDefaultIsolationLevel:                         Bool(true),
			TwoDigitCenturyStart:                                     Bool(true),
			UnsupportedDdlAction:                                     Bool(true),
			UserTaskManagedInitialWarehouseSize:                      Bool(true),
			UserTaskMinimumTriggerIntervalInSeconds:                  Bool(true),
			UserTaskTimeoutMs:                                        Bool(true),
			UseCachedResult:                                          Bool(true),
			UseWorkspacesForSql:                                      Bool(true),
			WeekOfYearPolicy:                                         Bool(true),
			WeekStart:                                                Bool(true),
		},
	}))
}

func (c *accounts) UnsetAllPoliciesSafely(ctx context.Context) error {
	return errors.Join(
		c.UnsetPolicySafely(ctx, PolicyKindAuthenticationPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindFeaturePolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPackagesPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPasswordPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindSessionPolicy),
	)
}

func (c *accounts) UnsetPolicySafely(ctx context.Context, kind PolicyKind) error {
	var unset *AccountUnset
	switch kind {
	case PolicyKindAuthenticationPolicy:
		unset = &AccountUnset{AuthenticationPolicy: Bool(true)}
	case PolicyKindFeaturePolicy:
		unset = &AccountUnset{FeaturePolicyUnset: &AccountFeaturePolicyUnset{FeaturePolicy: Bool(true)}}
	case PolicyKindPackagesPolicy:
		unset = &AccountUnset{PackagesPolicy: Bool(true)}
	case PolicyKindPasswordPolicy:
		unset = &AccountUnset{PasswordPolicy: Bool(true)}
	case PolicyKindSessionPolicy:
		unset = &AccountUnset{SessionPolicy: Bool(true)}
	default:
		return fmt.Errorf("policy kind %s is not supported for account policies", kind)
	}
	// If the policy is not attached to the account, Snowflake returns an error.
	err := c.Alter(ctx, NewAlterAccountRequest().WithUnset(*unset))
	if errors.Is(err, ErrPolicyNotAttachedToAccount) {
		return nil
	}
	return err
}

func (c *accounts) UnsetAll(ctx context.Context) error {
	return errors.Join(
		c.UnsetAllParameters(ctx),
		c.UnsetAllPoliciesSafely(ctx),
		c.Alter(ctx, NewAlterAccountRequest().WithUnset(AccountUnset{ResourceMonitor: Bool(true)})),
	)
}
