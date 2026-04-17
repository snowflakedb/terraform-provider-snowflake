package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userSessionPolicyAttachmentSchema = map[string]*schema.Schema{
	"user_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "User name of the user you want to attach the session policy to.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"session_policy_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Fully qualified name of the session policy.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

func UserSessionPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Specifies the session policy to use for a certain user.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingCreateWrapper(resources.UserSessionPolicyAttachment, CreateUserSessionPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingReadWrapper(resources.UserSessionPolicyAttachment, ReadUserSessionPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingDeleteWrapper(resources.UserSessionPolicyAttachment, DeleteUserSessionPolicyAttachment)),

		Schema: userSessionPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	sessionPolicy := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("session_policy_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			SessionPolicy: &sessionPolicy,
		},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating session policy attachment, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(userName.FullyQualifiedName(), sessionPolicy.FullyQualifiedName()))

	return ReadUserSessionPolicyAttachment(ctx, d, meta)
}

func ReadUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parts := helpers.ParseResourceIdentifier(d.Id())
	if len(parts) != 2 {
		return diag.FromErr(fmt.Errorf("required id format '<user_identifier>|<session_policy_fqn>', but got: '%s'", d.Id()))
	}

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
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

	sessionPolicyReferences := make([]sdk.PolicyReference, 0)
	for _, policyReference := range policyReferences {
		if policyReference.PolicyKind == sdk.PolicyKindSessionPolicy {
			sessionPolicyReferences = append(sessionPolicyReferences, policyReference)
		}
	}

	if len(sessionPolicyReferences) > 1 {
		return diag.FromErr(fmt.Errorf("internal error: multiple session policy references attached to a user. This should never happen"))
	}

	if len(sessionPolicyReferences) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find user's session policy. Marking the resource as removed.",
				Detail:   fmt.Sprintf("User id: %s", userName.Name()),
			},
		}
	}

	if err := d.Set("user_name", userName.Name()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(
		"session_policy_name",
		sdk.NewSchemaObjectIdentifier(
			*sessionPolicyReferences[0].PolicyDb,
			*sessionPolicyReferences[0].PolicySchema,
			sessionPolicyReferences[0].PolicyName,
		).FullyQualifiedName(),
	); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DeleteUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			SessionPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
