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

var userPasswordPolicyAttachmentSchema = map[string]*schema.Schema{
	"user_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "User name of the user you want to attach the password policy to",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"password_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Fully qualified name of the password policy",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

// UserPasswordPolicyAttachment returns a pointer to the resource representing a user password policy attachment.
func UserPasswordPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Specifies the password policy to use for a certain user.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserPasswordPolicyAttachmentResource), TrackingCreateWrapper(resources.UserPasswordPolicyAttachment, CreateUserPasswordPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserPasswordPolicyAttachmentResource), TrackingReadWrapper(resources.UserPasswordPolicyAttachment, ReadUserPasswordPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.UserPasswordPolicyAttachmentResource), TrackingDeleteWrapper(resources.UserPasswordPolicyAttachment, DeleteUserPasswordPolicyAttachment)),

		Schema: userPasswordPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateUserPasswordPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	passwordPolicy := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("password_policy_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			PasswordPolicy: &passwordPolicy,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(userName.FullyQualifiedName(), passwordPolicy.FullyQualifiedName()))

	return ReadUserPasswordPolicyAttachment(ctx, d, meta)
}

func ReadUserPasswordPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parts := helpers.ParseResourceIdentifier(d.Id())
	if len(parts) != 2 {
		return diag.FromErr(fmt.Errorf("required id format 'user_name|password_policy_name', but got: '%s'", d.Id()))
	}

	// Note: there is no alphanumeric id for an attachment, so we retrieve the password policies attached to a certain user.
	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(userName, sdk.PolicyEntityDomainUser))
	if err != nil {
		return diag.FromErr(err)
	}

	passwordPolicyReferences := make([]sdk.PolicyReference, 0)
	for _, policyReference := range policyReferences {
		if policyReference.PolicyKind == sdk.PolicyKindPasswordPolicy {
			passwordPolicyReferences = append(passwordPolicyReferences, policyReference)
		}
	}

	// Note: this should never happen, but just in case: so far, Snowflake only allows one Password Policy per user.
	if len(passwordPolicyReferences) > 1 {
		return diag.FromErr(fmt.Errorf("internal error: multiple policy references attached to a user. This should never happen"))
	}

	// Note: this means the resource has been deleted outside of Terraform.
	if len(passwordPolicyReferences) == 0 {
		d.SetId("")
		return nil
	}

	if err := d.Set("user_name", userName.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(
		"password_policy_name",
		sdk.NewSchemaObjectIdentifier(
			*passwordPolicyReferences[0].PolicyDb,
			*passwordPolicyReferences[0].PolicySchema,
			passwordPolicyReferences[0].PolicyName,
		).FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

func DeleteUserPasswordPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			PasswordPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
