package resources

import (
	"context"
	"errors"
	"log"
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

	handlePolicyCreate := func(kind sdk.PolicyKind, policyIdFieldPointerGetter func(*sdk.AccountSet) **sdk.SchemaObjectIdentifier, hasForce bool) error {
		key := strings.ToLower(string(kind))
		set := new(sdk.AccountSet)
		if kind == sdk.PolicyKindFeaturePolicy {
			set.FeaturePolicySet = new(sdk.AccountFeaturePolicySet)
		}

		if hasForce {
			set.Force = sdk.Bool(true)
		} else {
			log.Printf("[DEBUG] Unsetting %s as it doens't support setting with force", kind)
			if err := client.Accounts.UnsetPolicySafely(ctx, kind); err != nil {
				return err
			}
		}

		policySetFieldPointer := policyIdFieldPointerGetter(set)
		log.Printf("[DEBUG] Checking if %s is present in the configuration", key)
		if err := schemaObjectIdentifierAttributeCreate(d, key, policySetFieldPointer); err != nil {
			return err
		}

		if *policySetFieldPointer != nil {
			log.Printf("[DEBUG] Setting %s to the new value", key)
			if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
				return err
			}
		}

		return nil
	}

	if errs := errors.Join(
		handlePolicyCreate(sdk.PolicyKindAuthenticationPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.AuthenticationPolicy }, false),
		handlePolicyCreate(sdk.PolicyKindFeaturePolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.FeaturePolicySet.FeaturePolicy }, true),
		handlePolicyCreate(sdk.PolicyKindPackagesPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PackagesPolicy }, true),
		handlePolicyCreate(sdk.PolicyKindPasswordPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PasswordPolicy }, false),
		handlePolicyCreate(sdk.PolicyKindSessionPolicy, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.SessionPolicy }, false),
	); errs != nil {
		return diag.FromErr(errs)
	}

	setParameters := new(sdk.AccountParameters)
	if diags := handleAccountParametersCreate(d, setParameters); diags != nil {
		return diags
	}
	if *setParameters != (sdk.AccountParameters{}) {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{Parameters: setParameters}}); err != nil {
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

	handlePolicyUpdate := func(kind sdk.PolicyKind, hasForce bool, setFieldGetter func(*sdk.AccountSet) **sdk.SchemaObjectIdentifier) error {
		key := strings.ToLower(string(kind))
		if d.HasChange(key) {
			set, unset := new(sdk.AccountSet), new(sdk.AccountUnset)
			if kind == sdk.PolicyKindFeaturePolicy {
				set.FeaturePolicySet = new(sdk.AccountFeaturePolicySet)
				unset.FeaturePolicyUnset = new(sdk.AccountFeaturePolicyUnset)
			}

			setFieldPointer, unsetBoolFlag := setFieldGetter(set), sdk.Bool(false)
			log.Printf("[DEBUG] Checking for updates in %s", key)
			if err := schemaObjectIdentifierAttributeUpdate(d, key, setFieldPointer, &unsetBoolFlag); err != nil {
				return err
			}

			if *setFieldPointer != nil {
				if hasForce {
					set.Force = sdk.Bool(true)
				} else {
					log.Printf("[DEBUG] Unsetting %s as it doens't support setting with force", kind)
					if err := client.Accounts.UnsetPolicySafely(ctx, kind); err != nil {
						return err
					}
				}

				log.Printf("[DEBUG] Setting %s to the new value", kind)
				if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: set}); err != nil {
					return err
				}
			} else if *unsetBoolFlag {
				log.Printf("[DEBUG] Unsetting %s as it was removed from the configuration", kind)
				if err := client.Accounts.UnsetPolicySafely(ctx, kind); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if errs := errors.Join(
		handlePolicyUpdate(sdk.PolicyKindAuthenticationPolicy, false, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.AuthenticationPolicy }),
		handlePolicyUpdate(sdk.PolicyKindFeaturePolicy, true, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.FeaturePolicySet.FeaturePolicy }),
		handlePolicyUpdate(sdk.PolicyKindPackagesPolicy, true, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PackagesPolicy }),
		handlePolicyUpdate(sdk.PolicyKindPasswordPolicy, false, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.PasswordPolicy }),
		handlePolicyUpdate(sdk.PolicyKindSessionPolicy, false, func(set *sdk.AccountSet) **sdk.SchemaObjectIdentifier { return &set.SessionPolicy }),
	); errs != nil {
		return diag.FromErr(errs)
	}

	setParameters := new(sdk.AccountParameters)
	unsetParameters := new(sdk.AccountParametersUnset)
	if diags := handleAccountParametersUpdate(d, setParameters, unsetParameters); diags != nil {
		return diags
	}
	if *setParameters != (sdk.AccountParameters{}) {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{Parameters: setParameters}}); err != nil {
			return diag.FromErr(err)
		}
	}
	if *unsetParameters != (sdk.AccountParametersUnset{}) {
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{Parameters: unsetParameters}}); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadCurrentAccount(ctx, d, meta)
}

func DeleteCurrentAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.Accounts.UnsetAll(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
