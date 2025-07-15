package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var currentOrganizationAccountSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier (i.e. name) for the account. It must be unique within an organization, regardless of which Snowflake Region the account is in and must start with an alphabetic character and cannot contain spaces or special characters except for underscores (_). Note that if the account name includes underscores, features that do not accept account names with underscores (e.g. Okta SSO or SCIM) can reference a version of the account name that substitutes hyphens (-) for the underscores.",
	},
	"resource_monitor": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription("Parameter that specifies the name of the resource monitor used to control all virtual warehouses created in the account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"password_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("Specifies [password policy](https://docs.snowflake.com/en/user-guide/password-authentication#label-using-password-policies) for the current account.", resources.PasswordPolicy),
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

func CurrentOrganizationAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource used to manage an organization account within the organization you are connected to. See [ALTER ORGANIZATION ACCOUNT](https://docs.snowflake.com/en/sql-reference/sql/alter-organization-account) documentation for more information on resource capabilities.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingCreateWrapper(resources.CurrentOrganizationAccount, CreateCurrentOrganizationAccount)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingReadWrapper(resources.CurrentOrganizationAccount, ReadCurrentOrganizationAccount)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingUpdateWrapper(resources.CurrentOrganizationAccount, UpdateCurrentOrganizationAccount)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CurrentOrganizationAccountResource), TrackingDeleteWrapper(resources.CurrentOrganizationAccount, DeleteCurrentOrganizationAccount)),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.CurrentOrganizationAccount, accountParametersCustomDiff),

		Schema: collections.MergeMaps(currentOrganizationAccountSchema, accountParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CurrentOrganizationAccount, ImportName[sdk.AccountObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	organizationAccounts, err := client.OrganizationAccounts.Show(ctx, sdk.NewShowOrganizationAccountRequest())
	if err != nil {
		return diag.FromErr(err)
	}
	if len(organizationAccounts) != 1 {
		return diag.Errorf("expected 1 organization account, found %d", len(organizationAccounts))
	}
	currentOrganizationAccount := organizationAccounts[0]

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if id.Name() != currentOrganizationAccount.AccountName {
		return diag.Errorf("passed name: %s, doesn't match current organization account name: %s, renames can be performed only after resource initialization", id.Name(), currentOrganizationAccount.AccountName)
	}

	if v, ok := d.GetOk("resource_monitor"); ok {
		resourceMonitorId, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitorId))); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithResourceMonitor(true))); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("password_policy"); ok {
		passwordPolicyId, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.OrganizationAccounts.SetPolicySafely(ctx, sdk.PolicyKindPasswordPolicy, passwordPolicyId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindPasswordPolicy); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("session_policy"); ok {
		sessionPolicyId, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.OrganizationAccounts.SetPolicySafely(ctx, sdk.PolicyKindSessionPolicy, sessionPolicyId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindSessionPolicy); err != nil {
			return diag.FromErr(err)
		}
	}

	setParameters := new(sdk.AccountParameters)
	if diags := handleAccountParametersCreate(d, setParameters); diags != nil {
		return diags
	}
	if *setParameters != (sdk.AccountParameters{}) {
		if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithParameters(*setParameters))); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(id.Name())

	return ReadCurrentOrganizationAccount(ctx, d, meta)
}

func ReadCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	attachedPolicies, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(sdk.NewAccountObjectIdentifier(client.GetAccountLocator()), sdk.PolicyEntityDomainAccount))
	if err != nil {
		return diag.FromErr(err)
	}

	for _, policy := range attachedPolicies {
		switch policy.PolicyKind {
		case sdk.PolicyKindPasswordPolicy,
			sdk.PolicyKindSessionPolicy:
			if policy.PolicyDb != nil && policy.PolicySchema != nil {
				if err := d.Set(strings.ToLower(string(policy.PolicyKind)), sdk.NewSchemaObjectIdentifier(*policy.PolicyDb, *policy.PolicySchema, policy.PolicyName).FullyQualifiedName()); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	parameters, err := client.OrganizationAccounts.ShowParameters(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := handleAccountParameterRead(d, parameters); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		id, err := sdk.ParseAccountObjectIdentifier(newName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().
			WithName(sdk.NewAccountObjectIdentifier(oldName.(string))).
			WithRenameTo(*sdk.NewOrganizationAccountRenameRequest(&id))); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(id.Name())
	}

	if d.HasChange("resource_monitor") {
		if v, ok := d.GetOk("resource_monitor"); ok {
			resourceMonitorId, err := sdk.ParseAccountObjectIdentifier(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitorId))); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithResourceMonitor(true))); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("password_policy") {
		if _, newValue := d.GetChange("password_policy"); newValue != nil && newValue.(string) != "" {
			passwordPolicyId, err := sdk.ParseSchemaObjectIdentifier(newValue.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.OrganizationAccounts.SetPolicySafely(ctx, sdk.PolicyKindPasswordPolicy, passwordPolicyId); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindPasswordPolicy); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("session_policy") {
		if _, newValue := d.GetChange("session_policy"); newValue != nil && newValue.(string) != "" {
			sessionPolicyId, err := sdk.ParseSchemaObjectIdentifier(newValue.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.OrganizationAccounts.SetPolicySafely(ctx, sdk.PolicyKindSessionPolicy, sessionPolicyId); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindSessionPolicy); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	setParameters := new(sdk.AccountParameters)
	unsetParameters := new(sdk.AccountParametersUnset)
	if diags := handleAccountParametersUpdate(d, setParameters, unsetParameters); diags != nil {
		return diags
	}
	if *setParameters != (sdk.AccountParameters{}) {
		if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithParameters(*setParameters))); err != nil {
			return diag.FromErr(err)
		}
	}
	if *unsetParameters != (sdk.AccountParametersUnset{}) {
		if err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithParameters(*unsetParameters))); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadCurrentOrganizationAccount(ctx, d, meta)
}

func DeleteCurrentOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.OrganizationAccounts.UnsetAll(ctx); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
