package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountAuthenticationPolicyAttachmentSchema = map[string]*schema.Schema{
	"authentication_policy": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the authentication policy to apply to the current account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"for_all_person_users": {
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Description: joinWithSpace("If true, attaches the authentication policy to all person users in the current account.",
			"Conflicts with `for_all_service_users`. When neither field is set, the policy is attached account-wide."),
		ConflictsWith: []string{"for_all_service_users"},
	},
	"for_all_service_users": {
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Description: joinWithSpace("If true, attaches the authentication policy to all service users in the current account.",
			"Conflicts with `for_all_person_users`. When neither field is set, the policy is attached account-wide."),
		ConflictsWith: []string{"for_all_person_users"},
	},
}

func AccountAuthenticationPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the authentication policy to use for the current account. To set the authentication policy of a different account, use a provider alias.",

		CreateContext: TrackingCreateWrapper(resources.AccountAuthenticationPolicyAttachment, CreateAccountAuthenticationPolicyAttachment),
		ReadContext:   TrackingReadWrapper(resources.AccountAuthenticationPolicyAttachment, ReadAccountAuthenticationPolicyAttachment),
		UpdateContext: TrackingUpdateWrapper(resources.AccountAuthenticationPolicyAttachment, UpdateAccountAuthenticationPolicyAttachment),
		DeleteContext: TrackingDeleteWrapper(resources.AccountAuthenticationPolicyAttachment, DeleteAccountAuthenticationPolicyAttachment),

		Schema: accountAuthenticationPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	authenticationPolicy, err := sdk.ParseSchemaObjectIdentifier(d.Get("authentication_policy").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	setRequest := sdk.NewAccountAuthenticationPolicySetRequest().WithAuthenticationPolicy(authenticationPolicy)
	if err := errors.Join(
		boolAttributeCreate(d, "for_all_person_users", &setRequest.ForAllPersonUsers),
		boolAttributeCreate(d, "for_all_service_users", &setRequest.ForAllServiceUsers),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := client.Accounts.Alter(
		ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().WithAuthenticationPolicySet(*setRequest)),
	); err != nil {
		return diag.FromErr(fmt.Errorf("error while creating authentication policy attachment, err = %w", err))
	}

	scope := accountAuthenticationPolicyScope(d)
	d.SetId(helpers.EncodeResourceIdentifier(authenticationPolicy.FullyQualifiedName(), string(scope)))

	return ReadAccountAuthenticationPolicyAttachment(ctx, d, meta)
}

func ReadAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	policies, err := client.AuthenticationPolicies.Show(ctx, sdk.NewShowAuthenticationPolicyRequest().WithOn(sdk.On{Account: new(true)}))
	if err != nil {
		return diag.FromErr(err)
	}

	var attachedPolicy *sdk.AuthenticationPolicy
	expectedScope := accountAuthenticationPolicyScopeForRead(d)
	for i := range policies {
		if slices.Contains(policies[i].TargetScopes, expectedScope) {
			attachedPolicy = &policies[i]
			break
		}
	}

	if attachedPolicy == nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find account's authentication policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("No authentication policy is attached to the current account with the %s target scope.", expectedScope),
			},
		}
	}

	if err := errors.Join(
		d.Set("authentication_policy", attachedPolicy.ID().FullyQualifiedName()),
		d.Set("for_all_person_users", expectedScope == sdk.AuthenticationPolicyTargetScopePersonUsers),
		d.Set("for_all_service_users", expectedScope == sdk.AuthenticationPolicyTargetScopeServiceUsers),
	); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(attachedPolicy.ID().FullyQualifiedName(), string(expectedScope)))

	return nil
}

func UpdateAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("authentication_policy") {
		newAuthenticationPolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("authentication_policy").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		setRequest := sdk.NewAccountAuthenticationPolicySetRequest().WithAuthenticationPolicy(newAuthenticationPolicyName)
		if err := errors.Join(
			boolAttributeCreate(d, "for_all_person_users", &setRequest.ForAllPersonUsers),
			boolAttributeCreate(d, "for_all_service_users", &setRequest.ForAllServiceUsers),
		); err != nil {
			return diag.FromErr(err)
		}
		if err := client.Accounts.Alter(
			ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().
				WithAuthenticationPolicySet(*setRequest).
				WithForce(true)),
		); err != nil {
			return diag.FromErr(fmt.Errorf("error while setting new authentication policy on account, err = %w", err))
		}

		scope := accountAuthenticationPolicyScope(d)
		d.SetId(helpers.EncodeResourceIdentifier(newAuthenticationPolicyName.FullyQualifiedName(), string(scope)))
	}

	return ReadAccountAuthenticationPolicyAttachment(ctx, d, meta)
}

func DeleteAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	unsetRequest := sdk.NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true)
	if err := errors.Join(
		boolAttributeCreate(d, "for_all_person_users", &unsetRequest.ForAllPersonUsers),
		boolAttributeCreate(d, "for_all_service_users", &unsetRequest.ForAllServiceUsers),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := client.Accounts.Alter(
		ctx, sdk.NewAlterAccountRequest().
			WithUnset(*sdk.NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*unsetRequest)),
	); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func accountAuthenticationPolicyScope(d *schema.ResourceData) sdk.AuthenticationPolicyTargetScope {
	switch {
	case d.Get("for_all_person_users").(bool):
		return sdk.AuthenticationPolicyTargetScopePersonUsers
	case d.Get("for_all_service_users").(bool):
		return sdk.AuthenticationPolicyTargetScopeServiceUsers
	default:
		return sdk.AuthenticationPolicyTargetScopeAccount
	}
}

// accountAuthenticationPolicyScopeForRead resolves the target scope managed by this resource instance from its id.
// Since v2.20.0 the scope is encoded as the second id part (`<fully_qualified_name>|<scope>`, two parts). Older ids
// predate scoped attachments and are always account-wide.
func accountAuthenticationPolicyScopeForRead(d *schema.ResourceData) sdk.AuthenticationPolicyTargetScope {
	if parts := helpers.ParseResourceIdentifier(d.Id()); len(parts) == 2 {
		if scope, err := sdk.ToAuthenticationPolicyTargetScope(parts[1]); err == nil {
			return scope
		}
	}
	return sdk.AuthenticationPolicyTargetScopeAccount
}
