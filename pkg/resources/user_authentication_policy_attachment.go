package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userAuthenticationPolicyAttachmentSchema = map[string]*schema.Schema{
	"user_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "User name of the user you want to attach the authentication policy to",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"authentication_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Fully qualified name of the authentication policy",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

// UserAuthenticationPolicyAttachment returns a pointer to the resource representing a user authentication policy attachment.
func UserAuthenticationPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Specifies the authentication policy to use for a certain user.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserAuthenticationPolicyAttachmentResource), TrackingCreateWrapper(resources.UserAuthenticationPolicyAttachment, CreateUserAuthenticationPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserAuthenticationPolicyAttachmentResource), TrackingReadWrapper(resources.UserAuthenticationPolicyAttachment, ReadUserAuthenticationPolicyAttachment)),
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

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	authenticationPolicy := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("authentication_policy_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			AuthenticationPolicy: &authenticationPolicy,
		},
	})
	if err != nil {
		return diag.FromErr(err)
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

	// Note: there is no alphanumeric id for an attachment, so we retrieve the authentication policies attached to a certain user.
	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(userName, sdk.PolicyEntityDomainUser))
	if err != nil {
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
		return nil
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
		).FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

func DeleteUserAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			AuthenticationPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
