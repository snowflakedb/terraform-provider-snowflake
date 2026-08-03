package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userAuthenticationPolicyAttachmentSchema = map[string]*schema.Schema{
	"user_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedPipesFieldDescription("User name of the user you want to attach the authentication policy to."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"authentication_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the authentication policy."),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

func UserAuthenticationPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Specifies the authentication policy to use for a certain user.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserAuthenticationPolicyAttachmentResource), TrackingCreateWrapper(resources.UserAuthenticationPolicyAttachment, CreateUserAuthenticationPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserAuthenticationPolicyAttachmentResource), TrackingReadWrapper(resources.UserAuthenticationPolicyAttachment, ReadUserAuthenticationPolicyAttachment)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.UserAuthenticationPolicyAttachmentResource), TrackingUpdateWrapper(resources.UserAuthenticationPolicyAttachment, UpdateUserAuthenticationPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.UserAuthenticationPolicyAttachmentResource), TrackingDeleteWrapper(resources.UserAuthenticationPolicyAttachment, DeleteUserAuthenticationPolicyAttachment)),

		Schema: userAuthenticationPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateUserAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName, err := sdk.ParseAccountObjectIdentifier(d.Get("user_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	authenticationPolicy, err := sdk.ParseSchemaObjectIdentifier(d.Get("authentication_policy_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Users.Alter(ctx, sdk.NewAlterUserRequest(userName).
		WithSet(*sdk.NewUserSetRequest().WithAuthenticationPolicy(authenticationPolicy)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating authentication policy attachment, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(userName.FullyQualifiedName(), authenticationPolicy.FullyQualifiedName()))

	return ReadUserAuthenticationPolicyAttachment(ctx, d, meta)
}

func ReadUserAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parts := helpers.ParseResourceIdentifier(d.Id())
	if len(parts) != 2 {
		return diag.FromErr(fmt.Errorf("required id format 'user_name|authentication_policy_name', but got: '%s'", d.Id()))
	}

	userName, err := sdk.ParseAccountObjectIdentifier(parts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	// Note: there is no alphanumeric id for an attachment, so we retrieve the authentication policies attached to a certain user.
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(userName, sdk.PolicyEntityDomainUser))
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to get user policies. Marking the resource as removed.",
					Detail:   fmt.Sprintf("User id: %s, Err: %s", userName.Name(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	authenticationPolicyReferences := make([]sdk.PolicyReference, 0)
	for _, policyReference := range policyReferences {
		if policyReference.PolicyKind == sdk.PolicyKindAuthenticationPolicy {
			authenticationPolicyReferences = append(authenticationPolicyReferences, policyReference)
		}
	}

	// Note: this should never happen, but just in case: so far, Snowflake only allows one Authentication Policy per user.
	if len(authenticationPolicyReferences) > 1 {
		return diag.FromErr(fmt.Errorf("internal error: multiple policy references attached to a user. This should never happen"))
	}

	// Note: this means the resource has been deleted outside of Terraform.
	if len(authenticationPolicyReferences) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find user's authentication policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("User id: %s", userName.Name()),
			},
		}
	}

	if err := d.Set("user_name", userName.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(
		"authentication_policy_name",
		sdk.NewSchemaObjectIdentifier(
			*authenticationPolicyReferences[0].PolicyDb,
			*authenticationPolicyReferences[0].PolicySchema,
			authenticationPolicyReferences[0].PolicyName,
		).FullyQualifiedName(),
	); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateUserAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("authentication_policy_name") {
		userName, err := sdk.ParseAccountObjectIdentifier(d.Get("user_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		newAuthenticationPolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("authentication_policy_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Users.Alter(ctx, sdk.NewAlterUserRequest(userName).
			WithIfExists(true).
			WithSet(*sdk.NewUserSetRequest().
				WithAuthenticationPolicy(newAuthenticationPolicyName).
				WithForce(true))); err != nil {
			return diag.FromErr(fmt.Errorf("error while setting new authentication policy to user %v, err = %w", userName.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(userName.FullyQualifiedName(), newAuthenticationPolicyName.FullyQualifiedName()))
	}

	return ReadUserAuthenticationPolicyAttachment(ctx, d, meta)
}

func DeleteUserAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName, err := sdk.ParseAccountObjectIdentifier(d.Get("user_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Users.Alter(ctx, sdk.NewAlterUserRequest(userName).
		WithUnset(*sdk.NewUserUnsetRequest().WithAuthenticationPolicy(true)))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
