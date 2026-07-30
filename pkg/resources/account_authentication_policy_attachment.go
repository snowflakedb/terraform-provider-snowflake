package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
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
		Description:      "Fully qualified name of the authentication policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

func AccountAuthenticationPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the authentication policy to use for the current account. To set the authentication policy of a different account, use a provider alias.",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingCreateWrapper(resources.AccountAuthenticationPolicyAttachment, CreateAccountAuthenticationPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingReadWrapper(resources.AccountAuthenticationPolicyAttachment, ReadAccountAuthenticationPolicyAttachment)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingUpdateWrapper(resources.AccountAuthenticationPolicyAttachment, UpdateAccountAuthenticationPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingDeleteWrapper(resources.AccountAuthenticationPolicyAttachment, DeleteAccountAuthenticationPolicyAttachment)),

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

	err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
		WithSet(*sdk.NewAccountSetRequest().WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().WithAuthenticationPolicy(authenticationPolicy))))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating authentication policy attachment, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(authenticationPolicy))

	return ReadAccountAuthenticationPolicyAttachment(ctx, d, meta)
}

func ReadAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	currentAccountId, err := sdk.ParseAccountObjectIdentifier(client.GetAccountLocator())
	if err != nil {
		return diag.FromErr(err)
	}

	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(currentAccountId, sdk.PolicyEntityDomainAccount))
	if err != nil {
		return diag.FromErr(err)
	}

	authenticationPolicyReferences := make([]sdk.PolicyReference, 0)
	for _, policyReference := range policyReferences {
		if policyReference.PolicyKind == sdk.PolicyKindAuthenticationPolicy {
			authenticationPolicyReferences = append(authenticationPolicyReferences, policyReference)
		}
	}

	if len(authenticationPolicyReferences) > 1 {
		return diag.FromErr(fmt.Errorf("internal error: multiple authentication policy references attached to an account. This should never happen"))
	}

	if len(authenticationPolicyReferences) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find account's authentication policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Account id: %s", currentAccountId.Name()),
			},
		}
	}

	authenticationPolicyFromRef := sdk.NewSchemaObjectIdentifier(
		*authenticationPolicyReferences[0].PolicyDb,
		*authenticationPolicyReferences[0].PolicySchema,
		authenticationPolicyReferences[0].PolicyName,
	)

	if err := d.Set("authentication_policy", authenticationPolicyFromRef.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(authenticationPolicyFromRef))

	return nil
}

func UpdateAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("authentication_policy") {
		newAuthenticationPolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("authentication_policy").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithUnset(*sdk.NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*sdk.NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true)))); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while unsetting old authentication policy from account, err = %w", err))
		}
		if err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().WithAuthenticationPolicy(newAuthenticationPolicyName)))); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while setting new authentication policy on account, err = %w", err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newAuthenticationPolicyName))
	}

	return ReadAccountAuthenticationPolicyAttachment(ctx, d, meta)
}

func DeleteAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
		WithUnset(*sdk.NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*sdk.NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true))))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
