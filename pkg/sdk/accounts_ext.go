package sdk

import (
	"context"
	"errors"
	"fmt"
)

func (opts *CreateAccountOptions) additionalValidations() error {
	var errs []error
	if opts.AdminName == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "AdminName"))
	}
	if opts.Email == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Email"))
	}
	if opts.Edition == "" {
		errs = append(errs, errNotSet("CreateAccountOptions", "Edition"))
	}
	return errors.Join(errs...)
}

func (opts *AccountSet) additionalValidations() error {
	var errs []error
	if valueSet(opts.Force) &&
		!valueSet(opts.PackagesPolicy) &&
		!valueSet(opts.PasswordPolicy) &&
		!valueSet(opts.SessionPolicySet) &&
		!valueSet(opts.AuthenticationPolicySet) &&
		!valueSet(opts.FeaturePolicySet) {
		errs = append(errs, NewError("force can only be set with PackagesPolicy, PasswordPolicy, SessionPolicy, AuthenticationPolicy, or FeaturePolicy"))
	}
	return errors.Join(errs...)
}

func (opts *AccountLevelParameters) additionalValidations() error {
	var errs []error
	if valueSet(opts.AccountParameters) {
		if err := opts.AccountParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.SessionParameters) {
		if err := opts.SessionParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.ObjectParameters) {
		if err := opts.ObjectParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.UserParameters) {
		if err := opts.UserParameters.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *AccountUnset) additionalValidations() error {
	return nil
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
	if valueSet(opts.Drop) {
		if !exactlyOneValueSet(opts.Drop.OldUrl, opts.Drop.OldOrganizationUrl) {
			errs = append(errs, errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"))
		}
	}
	if valueSet(opts.SaveOldURL) && !valueSet(opts.RenameTo) {
		errs = append(errs, NewError("SaveOldURL can only be set with RenameTo"))
	}
	if valueSet(opts.Drop) || valueSet(opts.RenameTo) {
		if !valueSet(opts.Name) || !ValidObjectIdentifier(opts.Name) {
			errs = append(errs, ErrInvalidObjectIdentifier)
		}
	}
	return errors.Join(errs...)
}

func (row accountDBRow) additionalConvert(result *Account) error {
	if row.ManagedAccounts.Valid {
		result.ManagedAccounts = Int(int(row.ManagedAccounts.Int32))
	}
	return nil
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
			Account: new(true),
		},
	})
}

func (c *accounts) UnsetAllParameters(ctx context.Context) error {
	return c.Alter(ctx, NewAlterAccountRequest().WithUnset(*NewAccountUnsetRequest().WithParameters(
		*NewAccountParametersUnsetRequest().
			WithAbortDetachedQuery(true).
			WithActivePythonProfiler(true).
			WithAllowBindValuesAccess(true).
			WithAllowClientMfaCaching(true).
			WithAllowIdToken(true).
			WithAllowedSpcsWorkloadTypes(true).
			WithAutocommit(true).
			WithBaseLocationPrefix(true).
			WithBinaryInputFormat(true).
			WithBinaryOutputFormat(true).
			WithCatalog(true).
			WithCatalogSync(true).
			WithClientEnableLogInfoStatementParameters(true).
			WithClientEncryptionKeySize(true).
			WithClientMemoryLimit(true).
			WithClientMetadataRequestUseConnectionCtx(true).
			WithClientMetadataUseSessionDatabase(true).
			WithClientPrefetchThreads(true).
			WithClientResultChunkSize(true).
			WithClientResultColumnCaseInsensitive(true).
			WithClientSessionKeepAlive(true).
			WithClientSessionKeepAliveHeartbeatFrequency(true).
			WithClientTimestampTypeMapping(true).
			WithCortexCodeCliDailyEstCreditLimitPerUser(true).
			WithCortexCodeDesktopDailyEstCreditLimitPerUser(true).
			WithCortexCodeSnowsightDailyEstCreditLimitPerUser(true).
			WithCortexEnabledCrossRegion(true).
			WithCortexModelsAllowlist(true).
			WithCsvTimestampFormat(true).
			WithDataMetricSchedule(true).
			WithDataRetentionTimeInDays(true).
			WithDateInputFormat(true).
			WithDateOutputFormat(true).
			WithDefaultDbtVersion(true).
			WithDefaultDdlCollation(true).
			WithDefaultNotebookComputePoolCpu(true).
			WithDefaultNotebookComputePoolGpu(true).
			WithDefaultNullOrdering(true).
			WithDefaultStreamlitComputePool(true).
			WithDefaultStreamlitNotebookWarehouse(true).
			WithDisableUiDownloadButton(true).
			WithDisableUserPrivilegeGrants(true).
			WithDisallowedSpcsWorkloadTypes(true).
			WithEnableAutomaticSensitiveDataClassificationLog(true).
			WithEnableBudgetEventLogging(true).
			WithEnableConsoleOutput(true).
			WithEnableCortexAnalyst(true).
			WithEnableDataCompaction(true).
			WithEnableEgressCostOptimizer(true).
			WithEnableGetDdlUseDataTypeAlias(true).
			WithEnableIcebergMergeOnRead(true).
			WithEnableIdentifierFirstLogin(true).
			WithEnableInternalStagesPrivatelink(true).
			WithEnableNotebookCreationInPersonalDb(true).
			WithEnablePerAccountAppServicePrivatelinkUrl(true).
			WithEnablePersonalDatabase(true).
			WithEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement(true).
			WithEnableTagPropagationEventLogging(true).
			WithEnableTriSecretAndRekeyOptOutForImageRepository(true).
			WithEnableTriSecretAndRekeyOptOutForSpcsBlockStorage(true).
			WithEnableUnhandledExceptionsReporting(true).
			WithEnableUnloadPhysicalTypeOptimization(true).
			WithEnableUnredactedQuerySyntaxError(true).
			WithEnableUnredactedSecureObjectError(true).
			WithEnforceNetworkRulesForInternalStages(true).
			WithErrorOnNondeterministicMerge(true).
			WithErrorOnNondeterministicUpdate(true).
			WithEventTable(true).
			WithExternalOauthAddPrivilegedRolesToBlockedList(true).
			WithExternalVolume(true).
			WithGeographyOutputFormat(true).
			WithGeometryOutputFormat(true).
			WithHybridTableLockTimeout(true).
			WithIcebergVersionDefault(true).
			WithInitialReplicationSizeLimitInTb(true).
			WithJdbcTreatDecimalAsInt(true).
			WithJdbcTreatTimestampNtzAsUtc(true).
			WithJdbcUseSessionTimezone(true).
			WithJsTreatIntegerAsBigint(true).
			WithJsonIndent(true).
			WithListingAutoFulfillmentReplicationRefreshSchedule(true).
			WithLockTimeout(true).
			WithLogEventLevel(true).
			WithLogLevel(true).
			WithMaxConcurrencyLevel(true).
			WithMaxDataExtensionTimeInDays(true).
			WithMetricLevel(true).
			WithMinDataRetentionTimeInDays(true).
			WithMultiStatementCount(true).
			WithNetworkPolicy(true).
			WithNoorderSequenceAsDefault(true).
			WithOauthAddPrivilegedRolesToBlockedList(true).
			WithOdbcTreatDecimalAsInt(true).
			WithPeriodicDataRekeying(true).
			WithPipeExecutionPaused(true).
			WithPreventLoadFromInlineUrl(true).
			WithPreventUnloadToInlineUrl(true).
			WithPreventUnloadToInternalStages(true).
			WithPythonProfilerModules(true).
			WithPythonProfilerTargetStage(true).
			WithQueryTag(true).
			WithQuotedIdentifiersIgnoreCase(true).
			WithReadConsistencyMode(true).
			WithReplaceInvalidCharacters(true).
			WithRequireStorageIntegrationForStageCreation(true).
			WithRequireStorageIntegrationForStageOperation(true).
			WithRowTimestampDefault(true).
			WithRowsPerResultset(true).
			WithS3StageVpceDnsName(true).
			WithSearchPath(true).
			WithServerlessTaskMaxStatementSize(true).
			WithServerlessTaskMinStatementSize(true).
			WithShareRestrictions(true).
			WithSimulatedDataSharingConsumer(true).
			WithSqlTraceQueryText(true).
			WithSsoLoginPage(true).
			WithStatementQueuedTimeoutInSeconds(true).
			WithStatementTimeoutInSeconds(true).
			WithStorageSerializationPolicy(true).
			WithStrictJsonOutput(true).
			WithSuspendTaskAfterNumFailures(true).
			WithTaskAutoRetryAttempts(true).
			WithTimeInputFormat(true).
			WithTimeOutputFormat(true).
			WithTimestampDayIsAlways24H(true).
			WithTimestampInputFormat(true).
			WithTimestampLtzOutputFormat(true).
			WithTimestampNtzOutputFormat(true).
			WithTimestampOutputFormat(true).
			WithTimestampTypeMapping(true).
			WithTimestampTzOutputFormat(true).
			WithTimezone(true).
			WithTraceLevel(true).
			WithTransactionAbortOnError(true).
			WithTransactionDefaultIsolationLevel(true).
			WithTwoDigitCenturyStart(true).
			WithUnsupportedDdlAction(true).
			WithUseCachedResult(true).
			WithUseWorkspacesForSql(true).
			WithUserTaskManagedInitialWarehouseSize(true).
			WithUserTaskMinimumTriggerIntervalInSeconds(true).
			WithUserTaskTimeoutMs(true).
			WithWeekOfYearPolicy(true).
			WithWeekStart(true),
	)))
}

func (c *accounts) UnsetAllPoliciesSafely(ctx context.Context) error {
	return errors.Join(
		c.UnsetPolicySafely(ctx, PolicyKindAuthenticationPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindFeaturePolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPackagesPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindPasswordPolicy),
		c.UnsetPolicySafely(ctx, PolicyKindSessionPolicy),
		// The account-wide unsets above do not cover authentication and session policies attached to a specific user
		// type, so those scoped attachments have to be unset with the same FOR ALL syntax that was used to set them.
		c.unsetPolicySafely(ctx, NewAccountUnsetRequest().WithAuthenticationPolicyUnset(
			*NewAccountAuthenticationPolicyUnsetRequest().
				WithAuthenticationPolicy(true).
				WithForAllPersonUsers(true),
		)),
		c.unsetPolicySafely(ctx, NewAccountUnsetRequest().WithAuthenticationPolicyUnset(
			*NewAccountAuthenticationPolicyUnsetRequest().
				WithAuthenticationPolicy(true).
				WithForAllServiceUsers(true),
		)),
		c.unsetPolicySafely(ctx, NewAccountUnsetRequest().WithSessionPolicyUnset(
			*NewAccountSessionPolicyUnsetRequest().
				WithSessionPolicy(true).
				WithForAllPersonUsers(true),
		)),
		c.unsetPolicySafely(ctx, NewAccountUnsetRequest().WithSessionPolicyUnset(
			*NewAccountSessionPolicyUnsetRequest().
				WithSessionPolicy(true).
				WithForAllServiceUsers(true),
		)),
	)
}

func (c *accounts) UnsetPolicySafely(ctx context.Context, kind PolicyKind) error {
	var unset *AccountUnsetRequest
	switch kind {
	case PolicyKindAuthenticationPolicy:
		unset = NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true))
	case PolicyKindFeaturePolicy:
		unset = NewAccountUnsetRequest().WithFeaturePolicyUnset(*NewAccountFeaturePolicyUnsetRequest().WithFeaturePolicy(true))
	case PolicyKindPackagesPolicy:
		unset = NewAccountUnsetRequest().WithPackagesPolicy(true)
	case PolicyKindPasswordPolicy:
		unset = NewAccountUnsetRequest().WithPasswordPolicy(true)
	case PolicyKindSessionPolicy:
		unset = NewAccountUnsetRequest().WithSessionPolicyUnset(*NewAccountSessionPolicyUnsetRequest().WithSessionPolicy(true))
	default:
		return fmt.Errorf("policy kind %s is not supported for account policies", kind)
	}
	return c.unsetPolicySafely(ctx, unset)
}

func (c *accounts) unsetPolicySafely(ctx context.Context, unset *AccountUnsetRequest) error {
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
		c.Alter(ctx, NewAlterAccountRequest().WithUnset(*NewAccountUnsetRequest().WithResourceMonitor(true))),
	)
}
