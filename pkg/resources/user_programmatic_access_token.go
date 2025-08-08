package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var userProgrammaticAccessTokenSchema = map[string]*schema.Schema{
	"user": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The name of the user that the token is associated with. A user cannot use another user's programmatic access token to authenticate."),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the name for the programmatic access token; must be unique for the user."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"role_restriction": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The name of the role used for privilege evaluation and object creation. This must be one of the roles that has already been granted to the user."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"days_to_expiry": {
		Type:             schema.TypeInt,
		Optional:         true,
		ForceNew:         true,
		Description:      externalChangesNotDetectedFieldDescription("The number of days that the programmatic access token can be used for authentication. This field cannot be altered after the token is created. Instead, you must rotate the token with the `keeper` field."),
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
	},
	"mins_to_bypass_network_policy_requirement": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription("The number of minutes during which a user can use this token to access Snowflake without being subject to an active network policy."),
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
	},
	"disabled": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Disables or enables the programmatic access token."),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("status", func(x any) any {
			return x.(string) == string(sdk.ProgrammaticAccessTokenStatusDisabled)
		}),
	},
	"expire_rotated_token_after_hours": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      "This field is only used when the token is rotated by changing the `keeper` field. Sets the expiration time of the existing token secret to expire after the specified number of hours. You can set this to a value of 0 to expire the current token secret immediately.",
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
	},
	"keeper": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Arbitrary string that, if and only if, changed from a non-empty to a different non-empty value (or known after apply), will trigger a key to be rotated. When you add this field to the configuration, or remove it from the configuration, the rotation is not triggered. When the token is rotated, the `token` and `rotated_token_name` fields are marked as computed.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Descriptive comment about the programmatic access token.",
	},
	"token": {
		Type:        schema.TypeString,
		Computed:    true,
		Sensitive:   true,
		Description: "The token itself. Use this to authenticate to an endpoint. The data in this field is updated only when the token is created or rotated. In this case, the field is marked as computed.",
	},
	"rotated_token_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the token that represents the prior secret. This field is updated only when the token is rotated. In this case, the field is marked as computed.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW USER PROGRAMMATIC ACCESS TOKENS` for the given user programmatic access token.",
		Elem: &schema.Resource{
			Schema: schemas.ShowProgrammaticAccessTokenSchema,
		},
	},
}

func UserProgrammaticAccessToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserProgrammaticAccessTokenResource), TrackingCreateWrapper(resources.UserProgrammaticAccessToken, CreateUserProgrammaticAccessToken)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserProgrammaticAccessTokenResource), TrackingReadWrapper(resources.UserProgrammaticAccessToken, ReadUserProgrammaticAccessToken(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.UserProgrammaticAccessTokenResource), TrackingUpdateWrapper(resources.UserProgrammaticAccessToken, UpdateUserProgrammaticAccessToken)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.UserProgrammaticAccessTokenResource), TrackingDeleteWrapper(resources.UserProgrammaticAccessToken, DeleteUserProgrammaticAccessToken)),
		Description: joinWithSpace(
			"Resource used to manage user programmatic access tokens. For more information, check [user programmatic access tokens documentation](https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token).",
			"A programmatic access token is a token that can be used to authenticate to an endpoint.",
			"See [Using programmatic access tokens for authentication](https://docs.snowflake.com/en/user-guide/programmatic-access-tokens) user guide for more details.",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.UserProgrammaticAccessToken, customdiff.All(
			ComputedIfAnyAttributeChanged(userProgrammaticAccessTokenSchema, ShowOutputAttributeName, "disabled", "mins_to_bypass_network_policy_requirement", "comment"),
			func(_ context.Context, diff *schema.ResourceDiff, _ any) error {
				o, n := diff.GetChange("keeper")
				// If the key is being rotated, mark the `token` and `rotated_token_name` as computed to inform that these values will change.
				if shouldRotateToken(o.(string), n.(string), diff.GetRawPlan().AsValueMap()["keeper"].IsKnown()) {
					errs := errors.Join(
						diff.SetNewComputed("token"),
						diff.SetNewComputed("rotated_token_name"),
					)
					if errs != nil {
						return errs
					}
				}

				return nil
			},
		),
		),

		Schema: userProgrammaticAccessTokenSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.UserProgrammaticAccessToken, ImportUserProgrammaticAccessToken),
		},

		Timeouts: defaultTimeouts,
	}
}

// Rotate the token, but only when the change is from a non-empty to a non-empty value or the new value is unknown (known after apply).
// When the value was empty, but is now known after apply, it means that the token is not being rotated.
func shouldRotateToken(old, new string, isKnown bool) bool {
	return old != "" && (new != "" && old != new || !isKnown)
}

func ImportUserProgrammaticAccessToken(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := userProgrammaticAccessTokenIdFromData(d)
	if err != nil {
		return nil, err
	}

	token, err := client.Users.ShowProgrammaticAccessTokenByNameSafely(ctx, id.userName, id.tokenName)
	if err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("name", token.Name),
		d.Set("user", id.userName.Name()),
		// not reading mins_to_bypass_network_policy_requirement on purpose (it always changes)
		d.Set("disabled", booleanStringFromBool(token.Status == sdk.ProgrammaticAccessTokenStatusDisabled)),
	)
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateUserProgrammaticAccessToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	user := d.Get("user").(string)
	name := d.Get("name").(string)
	resourceId := userProgrammaticAccessTokenId{
		userName:  sdk.NewAccountObjectIdentifier(user),
		tokenName: sdk.NewAccountObjectIdentifier(name),
	}

	request := sdk.NewAddUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName)
	errs := errors.Join(
		accountObjectIdentifierAttributeCreate(d, "role_restriction", &request.RoleRestriction),
		intAttributeCreateBuilder(d, "days_to_expiry", request.WithDaysToExpiry),
		intAttributeCreateBuilder(d, "mins_to_bypass_network_policy_requirement", request.WithMinsToBypassNetworkPolicyRequirement),
		// disabled is handled separately
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	token, err := client.Users.AddProgrammaticAccessToken(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId.String())

	if v := d.Get("disabled").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		request := sdk.NewModifyProgrammaticAccessTokenSetRequest().WithDisabled(parsed)
		if err := client.Users.ModifyProgrammaticAccessToken(ctx, sdk.NewModifyUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName).WithSet(*request)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	err = errors.Join(
		d.Set("token", token.TokenSecret),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	return ReadUserProgrammaticAccessToken(false)(ctx, d, meta)
}

func ReadUserProgrammaticAccessToken(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		resourceId, err := userProgrammaticAccessTokenIdFromData(d)
		if err != nil {
			return diag.FromErr(err)
		}
		token, err := client.Users.ShowProgrammaticAccessTokenByNameSafely(ctx, resourceId.userName, resourceId.tokenName)
		if err != nil {
			if errors.Is(err, sdk.ErrPatNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query user programmatic access token. Marking the resource as removed.",
						Detail:   fmt.Sprintf("User programmatic access token name: %s for user: %s, Err: %s", resourceId.tokenName.FullyQualifiedName(), resourceId.userName.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"status", "disabled", string(token.Status), booleanStringFromBool(token.Status == sdk.ProgrammaticAccessTokenStatusDisabled), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, userProgrammaticAccessTokenSchema, []string{
			"disabled",
		}); err != nil {
			return diag.FromErr(err)
		}

		roleRestriction := ""
		if token.RoleRestriction != nil {
			roleRestriction = token.RoleRestriction.Name()
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ProgrammaticAccessTokenToSchema(token)}),
			d.Set("role_restriction", roleRestriction),
			d.Set("comment", token.Comment),
			// not reading mins_to_bypass_network_policy_requirement on purpose (it always changes)
			// not reading days_to_expiry on purpose (Snowflake returns expires_at which is a timestamp)
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateUserProgrammaticAccessToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	resourceId, err := userProgrammaticAccessTokenIdFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		err := client.Users.ModifyProgrammaticAccessToken(ctx, sdk.NewModifyUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming user programmatic access token %v err = %w", d.Id(), err))
		}

		resourceId.tokenName = newId
		d.SetId(resourceId.String())
	}

	setRequest := sdk.NewModifyProgrammaticAccessTokenSetRequest()
	unsetRequest := sdk.NewModifyProgrammaticAccessTokenUnsetRequest()
	errs := errors.Join(
		booleanStringAttributeUpdate(d, "disabled", &setRequest.Disabled, &unsetRequest.Disabled),
		stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment),
		intAttributeUpdate(d, "mins_to_bypass_network_policy_requirement", &setRequest.MinsToBypassNetworkPolicyRequirement, &unsetRequest.MinsToBypassNetworkPolicyRequirement),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if (*setRequest != sdk.ModifyProgrammaticAccessTokenSetRequest{}) {
		if err := client.Users.ModifyProgrammaticAccessToken(ctx, sdk.NewModifyUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName).WithSet(*setRequest)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}
	if (*unsetRequest != sdk.ModifyProgrammaticAccessTokenUnsetRequest{}) {
		if err := client.Users.ModifyProgrammaticAccessToken(ctx, sdk.NewModifyUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName).WithUnset(*unsetRequest)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	o, n := d.GetChange("keeper")
	if shouldRotateToken(o.(string), n.(string), d.GetRawPlan().AsValueMap()["keeper"].IsKnown()) {
		request := sdk.NewRotateUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName)
		if v := d.Get("expire_rotated_token_after_hours").(int); v != IntDefault {
			request.WithExpireRotatedTokenAfterHours(v)
		}

		token, err := client.Users.RotateProgrammaticAccessToken(ctx, request)
		if err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set("token", token.TokenSecret),
			d.Set("rotated_token_name", token.RotatedTokenName),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
	}

	return ReadUserProgrammaticAccessToken(false)(ctx, d, meta)
}

func DeleteUserProgrammaticAccessToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	resourceId, err := userProgrammaticAccessTokenIdFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Users.RemoveProgrammaticAccessTokenSafely(ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(resourceId.userName, resourceId.tokenName))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

type userProgrammaticAccessTokenId struct {
	userName  sdk.AccountObjectIdentifier
	tokenName sdk.AccountObjectIdentifier
}

func (id *userProgrammaticAccessTokenId) String() string {
	return helpers.EncodeResourceIdentifier(id.userName.FullyQualifiedName(), id.tokenName.FullyQualifiedName())
}

func userProgrammaticAccessTokenIdFromData(d *schema.ResourceData) (userProgrammaticAccessTokenId, error) {
	idRaw := helpers.ParseResourceIdentifier(d.Id())
	if len(idRaw) != 2 {
		return userProgrammaticAccessTokenId{}, fmt.Errorf("invalid resource id: %s", d.Id())
	}
	return userProgrammaticAccessTokenId{
		userName:  sdk.NewAccountObjectIdentifier(idRaw[0]),
		tokenName: sdk.NewAccountObjectIdentifier(idRaw[1]),
	}, nil
}
