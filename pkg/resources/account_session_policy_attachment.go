package resources

import (
	"context"
	"fmt"

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
		Description:      "Fully qualified name of the session policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
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

	sessionPolicyId, err := sdk.ParseSchemaObjectIdentifier(d.Get("session_policy_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			SessionPolicy: &sessionPolicyId,
		},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating session policy attachment, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(sessionPolicyId))

	return ReadAccountSessionPolicyAttachment(ctx, d, meta)
}

func ReadAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	currentAccountId, err := sdk.ParseAccountObjectIdentifier(client.GetAccountLocator())
	if err != nil {
		return diag.FromErr(err)
	}

	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(currentAccountId, sdk.PolicyEntityDomainAccount))
	if err != nil {
		return diag.FromErr(err)
	}

	sessionPolicyReferences := make([]sdk.PolicyReference, 0)
	for _, policyReference := range policyReferences {
		if policyReference.PolicyKind == sdk.PolicyKindSessionPolicy {
			sessionPolicyReferences = append(sessionPolicyReferences, policyReference)
		}
	}

	if len(sessionPolicyReferences) > 1 {
		return diag.FromErr(fmt.Errorf("internal error: multiple session policy references attached to an account. This should never happen"))
	}

	if len(sessionPolicyReferences) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find account's session policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Account id: %s", currentAccountId.Name()),
			},
		}
	}

	sessionPolicyFromRef := sdk.NewSchemaObjectIdentifier(
		*sessionPolicyReferences[0].PolicyDb,
		*sessionPolicyReferences[0].PolicySchema,
		sessionPolicyReferences[0].PolicyName,
	)

	if err := d.Set("session_policy_name", sessionPolicyFromRef.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(sessionPolicyFromRef))

	return nil
}

func UpdateAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("session_policy_name") {
		newSessionPolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("session_policy_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				SessionPolicy: sdk.Bool(true),
			},
		}); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while unsetting old session policy from account, err = %w", err))
		}
		if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				SessionPolicy: &newSessionPolicyName,
			},
		}); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while setting new session policy on account, err = %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newSessionPolicyName))
	}

	return ReadAccountSessionPolicyAttachment(ctx, d, meta)
}

func DeleteAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Unset: &sdk.AccountUnset{
			SessionPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
