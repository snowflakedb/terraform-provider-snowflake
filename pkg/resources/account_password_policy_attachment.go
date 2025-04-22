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

var accountPasswordPolicyAttachmentSchema = map[string]*schema.Schema{
	"password_policy": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Qualified name (`\"db\".\"schema\".\"policy_name\"`) of the password policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

func AccountPasswordPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the password policy to use for the current account. To set the password policy of a different account, use a provider alias.",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AccountPasswordPolicyAttachmentResource), TrackingCreateWrapper(resources.AccountPasswordPolicyAttachment, CreateAccountPasswordPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AccountPasswordPolicyAttachmentResource), TrackingReadWrapper(resources.AccountPasswordPolicyAttachment, ReadAccountPasswordPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AccountPasswordPolicyAttachmentResource), TrackingDeleteWrapper(resources.AccountPasswordPolicyAttachment, DeleteAccountPasswordPolicyAttachment)),

		Schema: accountPasswordPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateAccountPasswordPolicyAttachment implements schema.CreateFunc.
func CreateAccountPasswordPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	passwordPolicy, ok := sdk.NewObjectIdentifierFromFullyQualifiedName(d.Get("password_policy").(string)).(sdk.SchemaObjectIdentifier)
	if !ok {
		return diag.FromErr(fmt.Errorf("password_policy %s is not a valid password policy qualified name, expected format: `\"db\".\"schema\".\"policy\"`", d.Get("password_policy")))
	}

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			PasswordPolicy: passwordPolicy,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(passwordPolicy))

	return ReadAccountPasswordPolicyAttachment(ctx, d, meta)
}

func ReadAccountPasswordPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	passwordPolicy := helpers.DecodeSnowflakeID(d.Id())
	if err := d.Set("password_policy", passwordPolicy.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// DeleteAccountPasswordPolicyAttachment implements schema.DeleteFunc.
func DeleteAccountPasswordPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Unset: &sdk.AccountUnset{
			PasswordPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
