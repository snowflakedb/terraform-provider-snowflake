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
		Description:      "Fully qualified name of the session policy.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

func UserSessionPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Specifies the session policy to use for a certain user.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingCreateWrapper(resources.UserSessionPolicyAttachment, CreateUserSessionPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingReadWrapper(resources.UserSessionPolicyAttachment, ReadUserSessionPolicyAttachment)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingUpdateWrapper(resources.UserSessionPolicyAttachment, UpdateUserSessionPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.UserSessionPolicyAttachmentResource), TrackingDeleteWrapper(resources.UserSessionPolicyAttachment, DeleteUserSessionPolicyAttachment)),

		Schema: userSessionPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.UserSessionPolicyAttachment, ImportUserSessionPolicyAttachment),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	userName, err := getUserNameForUserSessionPolicyAttachment(d)
	if err != nil {
		return nil, err
	}

	if err := d.Set("user_name", userName.Name()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName, err := sdk.ParseAccountObjectIdentifier(d.Get("user_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	sessionPolicy, err := sdk.ParseSchemaObjectIdentifier(d.Get("session_policy_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Users.Alter(ctx, sdk.NewAlterUserRequest(userName).WithSet(*sdk.NewUserSetRequest().WithSessionPolicy(sessionPolicy)))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating session policy attachment, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(userName.FullyQualifiedName(), sessionPolicy.FullyQualifiedName()))

	return ReadUserSessionPolicyAttachment(ctx, d, meta)
}

func ReadUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName, err := getUserNameForUserSessionPolicyAttachment(d)
	if err != nil {
		return diag.FromErr(err)
	}
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

func UpdateUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("session_policy_name") {
		userName, err := getUserNameForUserSessionPolicyAttachment(d)
		if err != nil {
			return diag.FromErr(err)
		}

		newSessionPolicyName, err := sdk.ParseSchemaObjectIdentifier(d.Get("session_policy_name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Users.Alter(ctx, sdk.NewAlterUserRequest(*userName).WithIfExists(true).WithUnset(*sdk.NewUserUnsetRequest().WithSessionPolicy(true))); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while unsetting old session policy from user %v, err = %w", userName.FullyQualifiedName(), err))
		}
		if err := client.Users.Alter(ctx, sdk.NewAlterUserRequest(*userName).WithIfExists(true).WithSet(*sdk.NewUserSetRequest().WithSessionPolicy(newSessionPolicyName))); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error while setting new session policy to user %v, err = %w", userName.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(userName.FullyQualifiedName(), newSessionPolicyName.FullyQualifiedName()))
	}

	return ReadUserSessionPolicyAttachment(ctx, d, meta)
}

func DeleteUserSessionPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	userName, err := sdk.ParseAccountObjectIdentifier(d.Get("user_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Users.Alter(ctx, sdk.NewAlterUserRequest(userName).WithIfExists(true).WithUnset(*sdk.NewUserUnsetRequest().WithSessionPolicy(true)))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func getUserNameForUserSessionPolicyAttachment(d *schema.ResourceData) (*sdk.AccountObjectIdentifier, error) {
	parts := helpers.ParseResourceIdentifier(d.Id())
	if len(parts) != 2 {
		return nil, fmt.Errorf("required id format '<user_identifier>|<session_policy_fqn>', but got: '%s'", d.Id())
	}

	userName, err := sdk.ParseAccountObjectIdentifier(parts[0])
	if err != nil {
		return nil, err
	}
	return &userName, nil
}
