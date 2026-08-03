package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountSessionPolicyAttachmentSchema = map[string]*schema.Schema{
	"session_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the session policy to apply to the current account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"for_all_person_users": {
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Description: joinWithSpace("If true, attaches the session policy to all person users in the current account.",
			"Conflicts with `for_all_service_users`. When neither field is set, the policy is attached account-wide."),
		ConflictsWith: []string{"for_all_service_users"},
	},
	"for_all_service_users": {
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Description: joinWithSpace("If true, attaches the session policy to all service users in the current account.",
			"Conflicts with `for_all_person_users`. When neither field is set, the policy is attached account-wide."),
		ConflictsWith: []string{"for_all_person_users"},
	},
}

func AccountSessionPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the session policy to use for the current account. To set the session policy of a different account, use a provider alias.",

		CreateContext: TrackingCreateWrapper(resources.AccountSessionPolicyAttachment, CreateAccountSessionPolicyAttachment),
		ReadContext:   TrackingReadWrapper(resources.AccountSessionPolicyAttachment, ReadAccountSessionPolicyAttachment),
		UpdateContext: TrackingUpdateWrapper(resources.AccountSessionPolicyAttachment, UpdateAccountSessionPolicyAttachment),
		DeleteContext: TrackingDeleteWrapper(resources.AccountSessionPolicyAttachment, DeleteAccountSessionPolicyAttachment),

		Schema: accountSessionPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	sessionPolicy, err := sdk.ParseSchemaObjectIdentifier(d.Get("session_policy_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	setRequest := sdk.NewAccountSessionPolicySetRequest().WithSessionPolicy(sessionPolicy)
	if err := errors.Join(
		boolAttributeCreate(d, "for_all_person_users", &setRequest.ForAllPersonUsers),
		boolAttributeCreate(d, "for_all_service_users", &setRequest.ForAllServiceUsers),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := client.Accounts.Alter(
		ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().WithSessionPolicySet(*setRequest)),
	); err != nil {
		return diag.FromErr(fmt.Errorf("error while creating session policy attachment, err = %w", err))
	}

	scope := accountSessionPolicyScope(d)
	d.SetId(helpers.EncodeResourceIdentifier(sessionPolicy.FullyQualifiedName(), string(scope)))

	return ReadAccountSessionPolicyAttachment(ctx, d, meta)
}

func ReadAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	policies, err := client.SessionPolicies.Show(ctx, sdk.NewShowSessionPolicyRequest().WithOn(sdk.On{Account: new(true)}))
	if err != nil {
		return diag.FromErr(err)
	}

	var attachedPolicy *sdk.SessionPolicy
	expectedScope := accountSessionPolicyScopeForRead(d)
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
				Summary:  "Failed to find account's session policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("No session policy is attached to the current account with the %s target scope.", expectedScope),
			},
		}
	}

	if err := errors.Join(
		d.Set("session_policy_name", attachedPolicy.ID().FullyQualifiedName()),
		d.Set("for_all_person_users", expectedScope == sdk.SessionPolicyTargetScopePersonUsers),
		d.Set("for_all_service_users", expectedScope == sdk.SessionPolicyTargetScopeServiceUsers),
	); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(attachedPolicy.ID().FullyQualifiedName(), string(expectedScope)))

	return nil
}

func UpdateAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("session_policy_name") {
		newSessionPolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("session_policy_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		setRequest := sdk.NewAccountSessionPolicySetRequest().WithSessionPolicy(newSessionPolicyName)
		if err := errors.Join(
			boolAttributeCreate(d, "for_all_person_users", &setRequest.ForAllPersonUsers),
			boolAttributeCreate(d, "for_all_service_users", &setRequest.ForAllServiceUsers),
		); err != nil {
			return diag.FromErr(err)
		}
		if err := client.Accounts.Alter(
			ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().
				WithSessionPolicySet(*setRequest).
				WithForce(true)),
		); err != nil {
			return diag.FromErr(fmt.Errorf("error while setting new session policy on account, err = %w", err))
		}

		scope := accountSessionPolicyScope(d)
		d.SetId(helpers.EncodeResourceIdentifier(newSessionPolicyName.FullyQualifiedName(), string(scope)))
	}

	return ReadAccountSessionPolicyAttachment(ctx, d, meta)
}

func DeleteAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	unsetRequest := sdk.NewAccountSessionPolicyUnsetRequest().WithSessionPolicy(true)
	if err := errors.Join(
		boolAttributeCreate(d, "for_all_person_users", &unsetRequest.ForAllPersonUsers),
		boolAttributeCreate(d, "for_all_service_users", &unsetRequest.ForAllServiceUsers),
	); err != nil {
		return diag.FromErr(err)
	}

	if err := client.Accounts.Alter(
		ctx, sdk.NewAlterAccountRequest().
			WithUnset(*sdk.NewAccountUnsetRequest().WithSessionPolicyUnset(*unsetRequest)),
	); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func accountSessionPolicyScope(d *schema.ResourceData) sdk.SessionPolicyTargetScope {
	switch {
	case d.Get("for_all_person_users").(bool):
		return sdk.SessionPolicyTargetScopePersonUsers
	case d.Get("for_all_service_users").(bool):
		return sdk.SessionPolicyTargetScopeServiceUsers
	default:
		return sdk.SessionPolicyTargetScopeAccount
	}
}

// accountSessionPolicyScopeForRead resolves the target scope managed by this resource instance from its id. Since
// v2.20.0 the scope is encoded as the second id part (`<fully_qualified_name>|<scope>`, two parts). Older ids
// predate scoped attachments and are always account-wide.
func accountSessionPolicyScopeForRead(d *schema.ResourceData) sdk.SessionPolicyTargetScope {
	if parts := helpers.ParseResourceIdentifier(d.Id()); len(parts) == 2 {
		if scope, err := sdk.ToSessionPolicyTargetScope(parts[1]); err == nil {
			return scope
		}
	}
	return sdk.SessionPolicyTargetScopeAccount
}
