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
		ForceNew:         true,
		Description:      "Qualified name (`\"db\".\"schema\".\"policy_name\"`) of the authentication policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

func AccountAuthenticationPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the authentication policy to use for the current account. To set the authentication policy of a different account, use a provider alias.",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingCreateWrapper(resources.AccountAuthenticationPolicyAttachment, CreateAccountAuthenticationPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingReadWrapper(resources.AccountAuthenticationPolicyAttachment, ReadAccountAuthenticationPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AccountAuthenticationPolicyAttachmentResource), TrackingDeleteWrapper(resources.AccountAuthenticationPolicyAttachment, DeleteAccountAuthenticationPolicyAttachment)),

		Schema: accountAuthenticationPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateAccountAuthenticationPolicyAttachment implements schema.CreateFunc.
func CreateAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	authenticationPolicy, ok := sdk.NewObjectIdentifierFromFullyQualifiedName(d.Get("authentication_policy").(string)).(sdk.SchemaObjectIdentifier)
	if !ok {
		return diag.FromErr(fmt.Errorf("authentication_policy %s is not a valid authentication policy qualified name, expected format: `\"db\".\"schema\".\"policy\"`", d.Get("authentication_policy")))
	}

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			AuthenticationPolicy: authenticationPolicy,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(authenticationPolicy))

	return ReadAccountAuthenticationPolicyAttachment(ctx, d, meta)
}

func ReadAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	authenticationPolicy := helpers.DecodeSnowflakeID(d.Id())
	if err := d.Set("authentication_policy", authenticationPolicy.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// DeleteAccountAuthenticationPolicyAttachment implements schema.DeleteFunc.
func DeleteAccountAuthenticationPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Unset: &sdk.AccountUnset{
			AuthenticationPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
