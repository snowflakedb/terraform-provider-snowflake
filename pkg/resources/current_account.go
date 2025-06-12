package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO
// deprecate account_parameter resource
// in the next pr support policies and organization user group

var currentAccountSchema = map[string]*schema.Schema{
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription("Parameter that specifies the name of the resource monitor used to control all virtual warehouses created in the account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	/*
		"authentication_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "",
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"password_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "",
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"session_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "",
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"feature_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "",
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"packages_policy": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "",
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"organization_user_group": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			Optional:    true,
			Description: "The list of organization user groups imported into the account.",
		},
	*/
	// TODO: Tags are done by tags_association resource
}

func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "TODO",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CurrentAccountResource), TrackingCreateWrapper(resources.CurrentAccount, CreateCurrentAccount)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CurrentAccountResource), TrackingReadWrapper(resources.CurrentAccount, ReadCurrentAccount)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CurrentAccountResource), TrackingUpdateWrapper(resources.CurrentAccount, UpdateCurrentAccount)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CurrentAccountResource), TrackingDeleteWrapper(resources.CurrentAccount, DeleteCurrentAccount)),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.CurrentAccount, accountParametersCustomDiff),

		Schema: collections.MergeMaps(currentAccountSchema, accountParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CurrentAccount, schema.ImportStatePassthroughContext),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("current_account")
	return UpdateCurrentAccount(ctx, d, meta)
}

func ReadCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parameters, err := client.Accounts.ShowParameters(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := handleAccountParameterRead(d, parameters); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	alterIfIdentifierAttributeChanged := func(set *sdk.AccountSet, unset *sdk.AccountUnset, setId sdk.ObjectIdentifier, unsetBool *bool) diag.Diagnostics {
		if setId != nil {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
				return diag.FromErr(err)
			}
		}
		if unsetBool != nil && *unsetBool {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: unset}); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}

	if d.HasChange("resource_monitor") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.ResourceMonitor, unset.ResourceMonitor); diags != nil {
			return diags
		}
	}

	/* TODO(next prs): implement policies
	if d.HasChange("authentication_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "authentication_policy", &set.AuthenticationPolicy, &unset.AuthenticationPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.AuthenticationPolicy, unset.AuthenticationPolicy); diags != nil {
			return diags
		}
	}

	if d.HasChange("password_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "password_policy", &set.PasswordPolicy, &unset.PasswordPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.PasswordPolicy, unset.PasswordPolicy); diags != nil {
			return diags
		}
	}

	if d.HasChange("session_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "session_policy", &set.SessionPolicy, &unset.SessionPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.SessionPolicy, unset.SessionPolicy); diags != nil {
			return diags
		}
	}

	if d.HasChange("packages_policy") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := schemaObjectIdentifierAttributeUpdate(d, "packages_policy", &set.PackagesPolicy, &unset.PackagesPolicy); err != nil {
			return diag.FromErr(err)
		}
		if diags := alterIfIdentifierAttributeChanged(set, unset, set.PackagesPolicy, unset.PackagesPolicy); diags != nil {
			return diags
		}
	}
	*/

	setParameters := new(sdk.AccountSet)
	unsetParameters := new(sdk.AccountUnset)
	if diags := handleAccountParametersUpdate(d, setParameters, unsetParameters); diags != nil {
		return diags
	}
	if *setParameters != (sdk.AccountSet{}) {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: setParameters}); err != nil {
			return diag.FromErr(err)
		}
	}
	if *unsetParameters != (sdk.AccountUnset{}) {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: unsetParameters}); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadCurrentAccount(ctx, d, meta)
}

func DeleteCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if errs := errors.Join(
		client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{ResourceMonitor: sdk.Bool(true)}}),
		client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{
			Parameters: &sdk.AccountParametersUnset{
				AbortDetachedQuery:                               sdk.Bool(true),
				ActivePythonProfiler:                             sdk.Bool(true),
				AllowClientMFACaching:                            sdk.Bool(true),
				AllowIDToken:                                     sdk.Bool(true),
				Autocommit:                                       sdk.Bool(true),
				BaseLocationPrefix:                               sdk.Bool(true),
				BinaryInputFormat:                                sdk.Bool(true),
				BinaryOutputFormat:                               sdk.Bool(true),
				Catalog:                                          sdk.Bool(true),
				CatalogSync:                                      sdk.Bool(true),
				ClientEnableLogInfoStatementParameters:           sdk.Bool(true),
				ClientEncryptionKeySize:                          sdk.Bool(true),
				ClientMemoryLimit:                                sdk.Bool(true),
				ClientMetadataRequestUseConnectionCtx:            sdk.Bool(true),
				ClientMetadataUseSessionDatabase:                 sdk.Bool(true),
				ClientPrefetchThreads:                            sdk.Bool(true),
				ClientResultChunkSize:                            sdk.Bool(true),
				ClientResultColumnCaseInsensitive:                sdk.Bool(true),
				ClientSessionKeepAlive:                           sdk.Bool(true),
				ClientSessionKeepAliveHeartbeatFrequency:         sdk.Bool(true),
				ClientTimestampTypeMapping:                       sdk.Bool(true),
				CortexEnabledCrossRegion:                         sdk.Bool(true),
				CortexModelsAllowlist:                            sdk.Bool(true),
				CsvTimestampFormat:                               sdk.Bool(true),
				DataRetentionTimeInDays:                          sdk.Bool(true),
				DateInputFormat:                                  sdk.Bool(true),
				DateOutputFormat:                                 sdk.Bool(true),
				DefaultDDLCollation:                              sdk.Bool(true),
				DefaultNotebookComputePoolCpu:                    sdk.Bool(true),
				DefaultNotebookComputePoolGpu:                    sdk.Bool(true),
				DefaultNullOrdering:                              sdk.Bool(true),
				DefaultStreamlitNotebookWarehouse:                sdk.Bool(true),
				DisableUiDownloadButton:                          sdk.Bool(true),
				DisableUserPrivilegeGrants:                       sdk.Bool(true),
				EnableAutomaticSensitiveDataClassificationLog:    sdk.Bool(true),
				EnableEgressCostOptimizer:                        sdk.Bool(true),
				EnableIdentifierFirstLogin:                       sdk.Bool(true),
				EnableInternalStagesPrivatelink:                  sdk.Bool(true),
				EnableTriSecretAndRekeyOptOutForImageRepository:  sdk.Bool(true),
				EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: sdk.Bool(true),
				EnableUnhandledExceptionsReporting:               sdk.Bool(true),
				EnableUnloadPhysicalTypeOptimization:             sdk.Bool(true),
				EnableUnredactedQuerySyntaxError:                 sdk.Bool(true),
				EnableUnredactedSecureObjectError:                sdk.Bool(true),
				EnforceNetworkRulesForInternalStages:             sdk.Bool(true),
				ErrorOnNondeterministicMerge:                     sdk.Bool(true),
				ErrorOnNondeterministicUpdate:                    sdk.Bool(true),
				EventTable:                                       sdk.Bool(true),
				ExternalOAuthAddPrivilegedRolesToBlockedList:     sdk.Bool(true),
				ExternalVolume:                                   sdk.Bool(true),
				GeographyOutputFormat:                            sdk.Bool(true),
				GeometryOutputFormat:                             sdk.Bool(true),
				HybridTableLockTimeout:                           sdk.Bool(true),
				InitialReplicationSizeLimitInTB:                  sdk.Bool(true),
				JdbcTreatDecimalAsInt:                            sdk.Bool(true),
				JdbcTreatTimestampNtzAsUtc:                       sdk.Bool(true),
				JdbcUseSessionTimezone:                           sdk.Bool(true),
				JsonIndent:                                       sdk.Bool(true),
				JsTreatIntegerAsBigInt:                           sdk.Bool(true),
				ListingAutoFulfillmentReplicationRefreshSchedule: sdk.Bool(true),
				LockTimeout:                                      sdk.Bool(true),
				LogLevel:                                         sdk.Bool(true),
				MaxConcurrencyLevel:                              sdk.Bool(true),
				MaxDataExtensionTimeInDays:                       sdk.Bool(true),
				MetricLevel:                                      sdk.Bool(true),
				MinDataRetentionTimeInDays:                       sdk.Bool(true),
				MultiStatementCount:                              sdk.Bool(true),
				NetworkPolicy:                                    sdk.Bool(true),
				NoorderSequenceAsDefault:                         sdk.Bool(true),
				OAuthAddPrivilegedRolesToBlockedList:             sdk.Bool(true),
				OdbcTreatDecimalAsInt:                            sdk.Bool(true),
				PeriodicDataRekeying:                             sdk.Bool(true),
				PipeExecutionPaused:                              sdk.Bool(true),
				PreventUnloadToInlineURL:                         sdk.Bool(true),
				PreventUnloadToInternalStages:                    sdk.Bool(true),
				PythonProfilerModules:                            sdk.Bool(true),
				PythonProfilerTargetStage:                        sdk.Bool(true),
				QueryTag:                                         sdk.Bool(true),
				QuotedIdentifiersIgnoreCase:                      sdk.Bool(true),
				ReplaceInvalidCharacters:                         sdk.Bool(true),
				RequireStorageIntegrationForStageCreation:        sdk.Bool(true),
				RequireStorageIntegrationForStageOperation:       sdk.Bool(true),
				RowsPerResultset:                                 sdk.Bool(true),
				S3StageVpceDnsName:                               sdk.Bool(true),
				SamlIdentityProvider:                             sdk.Bool(true),
				SearchPath:                                       sdk.Bool(true),
				ServerlessTaskMaxStatementSize:                   sdk.Bool(true),
				ServerlessTaskMinStatementSize:                   sdk.Bool(true),
				SimulatedDataSharingConsumer:                     sdk.Bool(true),
				SsoLoginPage:                                     sdk.Bool(true),
				StatementQueuedTimeoutInSeconds:                  sdk.Bool(true),
				StatementTimeoutInSeconds:                        sdk.Bool(true),
				StorageSerializationPolicy:                       sdk.Bool(true),
				StrictJsonOutput:                                 sdk.Bool(true),
				SuspendTaskAfterNumFailures:                      sdk.Bool(true),
				TaskAutoRetryAttempts:                            sdk.Bool(true),
				TimestampDayIsAlways24h:                          sdk.Bool(true),
				TimestampInputFormat:                             sdk.Bool(true),
				TimestampLtzOutputFormat:                         sdk.Bool(true),
				TimestampNtzOutputFormat:                         sdk.Bool(true),
				TimestampOutputFormat:                            sdk.Bool(true),
				TimestampTypeMapping:                             sdk.Bool(true),
				TimestampTzOutputFormat:                          sdk.Bool(true),
				Timezone:                                         sdk.Bool(true),
				TimeInputFormat:                                  sdk.Bool(true),
				TimeOutputFormat:                                 sdk.Bool(true),
				TraceLevel:                                       sdk.Bool(true),
				TransactionAbortOnError:                          sdk.Bool(true),
				TransactionDefaultIsolationLevel:                 sdk.Bool(true),
				TwoDigitCenturyStart:                             sdk.Bool(true),
				UnsupportedDdlAction:                             sdk.Bool(true),
				UserTaskManagedInitialWarehouseSize:              sdk.Bool(true),
				UserTaskMinimumTriggerIntervalInSeconds:          sdk.Bool(true),
				UserTaskTimeoutMs:                                sdk.Bool(true),
				UseCachedResult:                                  sdk.Bool(true),
				WeekOfYearPolicy:                                 sdk.Bool(true),
				WeekStart:                                        sdk.Bool(true),
			},
		}}),
	); errs != nil {
		return diag.FromErr(errs)
	}

	d.SetId("")
	return nil
}
