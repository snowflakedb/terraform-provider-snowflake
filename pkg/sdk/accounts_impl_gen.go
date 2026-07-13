package sdk

import (
	"context"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/snowflakedb/gosnowflake/v2"
)

var (
	_ Accounts                = (*accounts)(nil)
	_ convertibleRow[Account] = new(accountDBRow)
)

type accounts struct {
	client *Client
}

func (c *accounts) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateAccountOptions) (*AccountCreateResponse, error) {
	if opts == nil {
		opts = &CreateAccountOptions{}
	}
	opts.name = id
	queryChanId := make(chan string, 1)
	err := validateAndExec(c.client, gosnowflake.WithQueryIDChan(ctx, queryChanId), opts)
	if err != nil {
		return nil, err
	}

	queryId := <-queryChanId
	rows, err := c.client.QueryUnsafe(gosnowflake.WithFetchResultByID(ctx, queryId), "")
	if err != nil {
		log.Printf("[WARN] Unable to retrieve create account output, err = %v", err)
	}

	response, err := getAccountCreateResponse(rows)
	if err != nil {
		return nil, fmt.Errorf("converting response from Snowflake: %w", err)
	}
	return response, nil
}

func getAccountCreateResponse(rows []map[string]*any) (*AccountCreateResponse, error) {
	if len(rows) != 1 {
		return nil, fmt.Errorf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["status"] == nil {
		return nil, fmt.Errorf("status is not set")
	}
	statusString, ok := (*rows[0]["status"]).(string)
	if !ok {
		return nil, fmt.Errorf("could not convert status to string")
	}
	return ToAccountCreateResponse(statusString)
}

func (c *accounts) Alter(ctx context.Context, opts *AlterAccountOptions) error {
	if opts == nil {
		opts = &AlterAccountOptions{}
	}
	return validateAndExec(c.client, ctx, opts)
}

func (c *accounts) Show(ctx context.Context, opts *ShowAccountOptions) ([]Account, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[accountDBRow](c.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[accountDBRow, Account](dbRows)
}

func (c *accounts) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	accounts, err := c.Show(ctx, &ShowAccountOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(accounts, func(account Account) bool {
		return account.AccountName == id.Name()
	})
}

func (c *accounts) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Account, error) {
	return SafeShowById(c.client, c.ShowByID, ctx, id)
}

func (c *accounts) Drop(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int, opts *DropAccountOptions) error {
	if opts == nil {
		opts = &DropAccountOptions{}
	}
	opts.name = id
	opts.gracePeriodInDays = gracePeriodInDays
	return validateAndExec(c.client, ctx, opts)
}

func (c *accounts) DropSafely(ctx context.Context, id AccountObjectIdentifier, gracePeriodInDays int) error {
	return SafeDrop(c.client, func() error { return c.Drop(ctx, id, gracePeriodInDays, &DropAccountOptions{IfExists: Bool(true)}) }, ctx, id)
}

func (c *accounts) Undrop(ctx context.Context, id AccountObjectIdentifier) error {
	opts := &undropAccountOptions{
		name: id,
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, sql)
	return err
}

func (c *accounts) ShowParameters(ctx context.Context) ([]*Parameter, error) {
	return c.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Account: Bool(true),
		},
	})
}

func (c *accounts) UnsetAllParameters(ctx context.Context) error {
	return c.client.Accounts.Alter(ctx, &AlterAccountOptions{Unset: &AccountUnset{
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
	}})
}

func (row accountDBRow) convert() (*Account, error) {
	acc := &Account{
		OrganizationName:      row.OrganizationName,
		AccountName:           row.AccountName,
		SnowflakeRegion:       row.SnowflakeRegion,
		AccountLocator:        row.AccountLocator,
		IsOrganizationAccount: row.IsOrganizationAccount,
	}
	if row.RegionGroup.Valid {
		acc.RegionGroup = &row.RegionGroup.String
	}
	if row.Edition.Valid {
		acc.Edition = Pointer(AccountEdition(row.Edition.String))
	}
	if row.AccountURL.Valid {
		acc.AccountURL = &row.AccountURL.String
	}
	if row.CreatedOn.Valid {
		acc.CreatedOn = &row.CreatedOn.Time
	}
	if row.Comment.Valid {
		acc.Comment = &row.Comment.String
	}
	if row.AccountLocatorURL.Valid {
		acc.AccountLocatorUrl = &row.AccountLocatorURL.String
	}
	if row.ManagedAccounts.Valid {
		acc.ManagedAccounts = Int(int(row.ManagedAccounts.Int32))
	}
	if row.ConsumptionBillingEntityName.Valid {
		acc.ConsumptionBillingEntityName = &row.ConsumptionBillingEntityName.String
	}
	if row.OldAccountURL.Valid {
		acc.OldAccountURL = &row.OldAccountURL.String
	}
	if row.IsOrgAdmin.Valid {
		acc.IsOrgAdmin = &row.IsOrgAdmin.Bool
	}
	if row.OrganizationOldUrl.Valid {
		acc.OrganizationOldUrl = &row.OrganizationOldUrl.String
	}
	if row.IsEventsAccount.Valid {
		acc.IsEventsAccount = &row.IsEventsAccount.Bool
	}
	if row.MarketplaceConsumerBillingEntityName.Valid {
		acc.MarketplaceConsumerBillingEntityName = &row.MarketplaceConsumerBillingEntityName.String
	}
	if row.MarketplaceProviderBillingEntityName.Valid {
		acc.MarketplaceProviderBillingEntityName = &row.MarketplaceProviderBillingEntityName.String
	}
	if row.AccountOldUrlSavedOn.Valid {
		acc.AccountOldUrlSavedOn = &row.AccountOldUrlSavedOn.Time
	}
	if row.AccountOldUrlLastUsed.Valid {
		acc.AccountOldUrlLastUsed = &row.AccountOldUrlLastUsed.Time
	}
	if row.OrganizationOldUrlSavedOn.Valid {
		acc.OrganizationOldUrlSavedOn = &row.OrganizationOldUrlSavedOn.Time
	}
	if row.OrganizationOldUrlLastUsed.Valid {
		acc.OrganizationOldUrlLastUsed = &row.OrganizationOldUrlLastUsed.Time
	}
	if row.DroppedOn.Valid {
		acc.DroppedOn = &row.DroppedOn.Time
	}
	if row.ScheduledDeletionTime.Valid {
		acc.ScheduledDeletionTime = &row.ScheduledDeletionTime.Time
	}
	if row.RestoredOn.Valid {
		acc.RestoredOn = &row.RestoredOn.Time
	}
	if row.MovedToOrganization.Valid {
		acc.MovedToOrganization = &row.MovedToOrganization.String
	}
	if row.MovedOn.Valid {
		acc.MovedOn = &row.MovedOn.String
	}
	if row.OrganizationUrlExpirationOn.Valid {
		acc.OrganizationUrlExpirationOn = &row.OrganizationUrlExpirationOn.Time
	}
	return acc, nil
}
