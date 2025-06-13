package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentAccountSchema = map[string]*schema.Schema{
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription("Parameter that specifies the name of the resource monitor used to control all virtual warehouses created in the account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"authentication_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [authentication policy](https://docs.snowflake.com/en/user-guide/authentication-policies) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"feature_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [feature policy](https://docs.snowflake.com/en/developer-guide/native-apps/ui-consumer-feature-policies) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"packages_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [packages policy](https://docs.snowflake.com/en/developer-guide/udf/python/packages-policy) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"password_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [password policy](https://docs.snowflake.com/en/user-guide/password-authentication#label-using-password-policies) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"session_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies [session policy](https://docs.snowflake.com/en/user-guide/session-policies-using) for the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

func CurrentAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource used to manage the account you are currently connected to. This resource is used to set account parameters and other account-level settings. See [ALTER ACCOUNT](https://docs.snowflake.com/en/sql-reference/sql/alter-account) documentation for more information on resource capabilities.",
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
	client := meta.(*provider.Context).Client

	setResourceMonitor := new(sdk.AccountSet)
	if err := accountObjectIdentifierAttributeCreate(d, "resource_monitor", &setResourceMonitor.ResourceMonitor); err != nil {
		return diag.FromErr(err)
	}
	if setResourceMonitor.ResourceMonitor != nil {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: setResourceMonitor}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{ResourceMonitor: sdk.Bool(true)}}); err != nil {
			return diag.FromErr(err)
		}
	}

	handlePolicyCreate := func(kind sdk.PolicyKind, policyIdFieldPointerGetter func(*sdk.AccountSet) **sdk.SchemaObjectIdentifier, hasForce bool, unsetPolicy *sdk.AccountUnset) error {
		set := new(sdk.AccountSet)
		if kind == sdk.PolicyKindFeaturePolicy {
			set.FeaturePolicySet = new(sdk.AccountFeaturePolicySet)
		}

		if hasForce {
			set.Force = sdk.Bool(true)
		} else {
			// To avoid Snowflake errors, every policy has to be first unset before it can be set (unless a policy can be set forcefully).
			if err := unsetAccountPolicySafely(client, ctx, kind, unsetPolicy); err != nil {
				return err
			}
		}

		policySetFieldPointer := policyIdFieldPointerGetter(set)
		if err := schemaObjectIdentifierAttributeCreate(d, strings.ToLower(string(kind)), policySetFieldPointer); err != nil {
			return err
		}

		if *policySetFieldPointer != nil {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
				return err
			}
		}

		return nil
	}

	if errs := errors.Join(
		handlePolicyCreate(sdk.PolicyKindAuthenticationPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.AuthenticationPolicy }, false, &sdk.AccountUnset{AuthenticationPolicy: sdk.Bool(true)}),
		handlePolicyCreate(sdk.PolicyKindFeaturePolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.FeaturePolicySet.FeaturePolicy }, true, &sdk.AccountUnset{FeaturePolicyUnset: &sdk.AccountFeaturePolicyUnset{FeaturePolicy: sdk.Bool(true)}}),
		handlePolicyCreate(sdk.PolicyKindPackagesPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PackagesPolicy }, true, &sdk.AccountUnset{PackagesPolicy: sdk.Bool(true)}),
		handlePolicyCreate(sdk.PolicyKindPasswordPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PasswordPolicy }, false, &sdk.AccountUnset{PasswordPolicy: sdk.Bool(true)}),
		handlePolicyCreate(sdk.PolicyKindSessionPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.SessionPolicy }, false, &sdk.AccountUnset{SessionPolicy: sdk.Bool(true)}),
	); errs != nil {
		return diag.FromErr(errs)
	}

	setParameters := new(sdk.AccountSet)
	if diags := handleAccountParametersCreate(d, setParameters); diags != nil {
		return diags
	}
	if *setParameters != (sdk.AccountSet{}) {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: setParameters}); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("current_account")

	return ReadCurrentAccount(ctx, d, meta)
}

func ReadCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	attachedPolicies, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(sdk.NewAccountObjectIdentifier(client.GetAccountLocator()), sdk.PolicyEntityDomainAccount))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, policy := range attachedPolicies {
		switch policy.PolicyKind {
		case sdk.PolicyKindAuthenticationPolicy,
			sdk.PolicyKindFeaturePolicy,
			sdk.PolicyKindPackagesPolicy,
			sdk.PolicyKindPasswordPolicy,
			sdk.PolicyKindSessionPolicy:
			if err := d.Set(strings.ToLower(string(policy.PolicyKind)), sdk.NewSchemaObjectIdentifier(*policy.PolicyDb, *policy.PolicySchema, policy.PolicyName).FullyQualifiedName()); err != nil {
				return diag.FromErr(err)
			}
		}
	}

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

	if d.HasChange("resource_monitor") {
		set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
		if err := accountObjectIdentifierAttributeUpdate(d, "resource_monitor", &set.ResourceMonitor, &unset.ResourceMonitor); err != nil {
			return diag.FromErr(err)
		}
		if set.ResourceMonitor != nil {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
				return diag.FromErr(err)
			}
		}
		if unset.ResourceMonitor != nil && *unset.ResourceMonitor {
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: unset}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	handlePolicyUpdate := func(kind sdk.PolicyKind, hasForce bool, setFieldGetter func(*sdk.AccountSet) **sdk.SchemaObjectIdentifier, unsetFieldGetter func(*sdk.AccountUnset) **bool) error {
		key := strings.ToLower(string(kind))
		if d.HasChange(key) {
			set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
			if kind == sdk.PolicyKindFeaturePolicy {
				set.FeaturePolicySet = new(sdk.AccountFeaturePolicySet)
				unset.FeaturePolicyUnset = new(sdk.AccountFeaturePolicyUnset)
			}

			setFieldPointer, unsetFieldPointer := setFieldGetter(set), unsetFieldGetter(unset)
			if err := schemaObjectIdentifierAttributeUpdate(d, key, setFieldPointer, unsetFieldPointer); err != nil {
				return err
			}

			if *setFieldPointer != nil {
				if hasForce {
					set.Force = sdk.Bool(true)
				} else {
					*unsetFieldPointer = sdk.Bool(true)
					if err := unsetAccountPolicySafely(client, ctx, kind, unset); err != nil {
						return err
					}
				}

				if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
					return err
				}
			} else if *unsetFieldPointer != nil && **unsetFieldPointer {
				if err := unsetAccountPolicySafely(client, ctx, kind, unset); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if errs := errors.Join(
		handlePolicyUpdate(sdk.PolicyKindAuthenticationPolicy, false, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.AuthenticationPolicy }, func(unset *sdk.AccountUnset) **bool { return &unset.AuthenticationPolicy }),
		handlePolicyUpdate(sdk.PolicyKindFeaturePolicy, true, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.FeaturePolicySet.FeaturePolicy }, func(unset *sdk.AccountUnset) **bool { return &unset.FeaturePolicyUnset.FeaturePolicy }),
		handlePolicyUpdate(sdk.PolicyKindPackagesPolicy, true, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PackagesPolicy }, func(unset *sdk.AccountUnset) **bool { return &unset.PackagesPolicy }),
		handlePolicyUpdate(sdk.PolicyKindPasswordPolicy, false, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PasswordPolicy }, func(unset *sdk.AccountUnset) **bool { return &unset.PasswordPolicy }),
		handlePolicyUpdate(sdk.PolicyKindSessionPolicy, false, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.SessionPolicy }, func(unset *sdk.AccountUnset) **bool { return &unset.SessionPolicy }),
	); errs != nil {
		return diag.FromErr(errs)
	}

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
		unsetAccountPolicySafely(client, ctx, sdk.PolicyKindAuthenticationPolicy, &sdk.AccountUnset{AuthenticationPolicy: sdk.Bool(true)}),
		unsetAccountPolicySafely(client, ctx, sdk.PolicyKindFeaturePolicy, &sdk.AccountUnset{FeaturePolicyUnset: &sdk.AccountFeaturePolicyUnset{FeaturePolicy: sdk.Bool(true)}}),
		unsetAccountPolicySafely(client, ctx, sdk.PolicyKindPackagesPolicy, &sdk.AccountUnset{PackagesPolicy: sdk.Bool(true)}),
		unsetAccountPolicySafely(client, ctx, sdk.PolicyKindPasswordPolicy, &sdk.AccountUnset{PasswordPolicy: sdk.Bool(true)}),
		unsetAccountPolicySafely(client, ctx, sdk.PolicyKindSessionPolicy, &sdk.AccountUnset{SessionPolicy: sdk.Bool(true)}),
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

func unsetAccountPolicySafely(client *sdk.Client, ctx context.Context, kind sdk.PolicyKind, unset *sdk.AccountUnset) error {
	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: unset})
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("Any policy of kind %s is not attached to ACCOUNT", kind)) {
		return nil
	}
	return err
}
