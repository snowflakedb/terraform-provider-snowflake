package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountSessionPolicyAttachmentSchema = map[string]*schema.Schema{
	"session_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Fully qualified name of the session policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

func AccountSessionPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the session policy to use for the current account. To set the session policy of a different account, use a provider alias.",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AccountSessionPolicyAttachmentResource), TrackingCreateWrapper(resources.AccountSessionPolicyAttachment, CreateAccountSessionPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AccountSessionPolicyAttachmentResource), TrackingReadWrapper(resources.AccountSessionPolicyAttachment, ReadAccountSessionPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AccountSessionPolicyAttachmentResource), TrackingDeleteWrapper(resources.AccountSessionPolicyAttachment, DeleteAccountSessionPolicyAttachment)),

		Schema: accountSessionPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	sessionPolicyId, ok := sdk.NewObjectIdentifierFromFullyQualifiedName(d.Get("session_policy_name").(string)).(sdk.SchemaObjectIdentifier)
	if !ok {
		return diag.FromErr(fmt.Errorf("session_policy_name %s is not a valid session policy qualified name, expected format: `\"db\".\"schema\".\"policy\"`", d.Get("session_policy_name")))
	}

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
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

func ReadAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	sessionPolicyId, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	currentAccountName, err := client.ContextFunctions.CurrentAccount(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	currentAccountId, err := sdk.ParseAccountObjectIdentifier(currentAccountName)
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
			sessionPolicyReferences = append(sessionPolicyReferences, sdk.PolicyReference{})
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

	if err := d.Set("session_policy_name", sessionPolicyId.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DeleteAccountSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Unset: &sdk.AccountUnset{
			SessionPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
