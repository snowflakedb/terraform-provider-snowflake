package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var oauthIntegrationForPartnerApplicationsSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the name of the OAuth integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"oauth_client": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Creates an OAuth interface between Snowflake and a partner application. Valid options are: %v.", possibleValuesListed(sdk.AllOauthSecurityIntegrationClients)),
		ValidateDiagFunc: sdkValidation(sdk.ToOauthSecurityIntegrationClientOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToOauthSecurityIntegrationClientOption),
	},
	"oauth_redirect_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Specifies the client URI. After a user is authenticated, the web browser is redirected to this URI. The field should be only set when OAUTH_CLIENT = LOOKER. In any other case the field should be left out empty.",
	},
	"enabled": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("enabled"),
		Description:      booleanStringFieldDescription("Specifies whether this OAuth integration is enabled or disabled."),
	},
	"oauth_issue_refresh_tokens": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_issue_refresh_tokens"),
		Description:      booleanStringFieldDescription("Specifies whether to allow the client to exchange a refresh token for an access token when the current access token has expired."),
	},
	"oauth_refresh_token_validity": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_refresh_token_validity"),
		Description:      "Specifies how long refresh tokens should be valid (in seconds). OAUTH_ISSUE_REFRESH_TOKENS must be set to TRUE.",
	},
	"oauth_use_secondary_roles": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Specifies whether default secondary roles set in the user properties are activated by default in the session being opened. Valid options are: %v.", possibleValuesListed(sdk.AllOauthSecurityIntegrationUseSecondaryRoles)),
		ValidateDiagFunc: sdkValidation(sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption), IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_use_secondary_roles")),
	},
	"blocked_roles_list": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:    true,
		Description: relatedResourceDescription(withPrivilegedRolesDescription("A set of Snowflake roles that a user cannot explicitly consent to using after authenticating.", string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)), resources.AccountRole),
		DiffSuppressFunc: SuppressIfAny(
			IgnoreChangeToCurrentSnowflakeListValueInDescribe("blocked_roles_list"),
			IgnoreValuesFromSetIfParamSet("blocked_roles_list", string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList), privilegedRoles),
			// TODO(SNOW-1517937): uncomment
			// NormalizeAndCompareIdentifiersInSet("blocked_roles_list"),
		),
	},
	"comment": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies a comment for the OAuth integration.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("comment"),
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECURITY INTEGRATION` for the given integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecurityIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECURITY INTEGRATION` for the given integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeOauthIntegrationForPartnerApplications,
		},
	},
	RelatedParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Parameters related to this security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowOauthForPartnerApplicationsParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func OauthIntegrationForPartnerApplications() *schema.Resource {
	return &schema.Resource{
		Schema: oauthIntegrationForPartnerApplicationsSchema,

		CreateContext: TrackingCreateWrapper(resources.OauthIntegrationForPartnerApplications, CreateContextOauthIntegrationForPartnerApplications),
		ReadContext:   TrackingReadWrapper(resources.OauthIntegrationForPartnerApplications, ReadContextOauthIntegrationForPartnerApplications(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OauthIntegrationForPartnerApplications, UpdateContextOauthIntegrationForPartnerApplications),
		DeleteContext: TrackingDeleteWrapper(resources.OauthIntegrationForPartnerApplications, DeleteSecurityIntegration),
		Description:   "Resource used to manage oauth security integration for partner applications objects. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OauthIntegrationForPartnerApplications, customdiff.All(
			ComputedIfAnyAttributeChanged(
				oauthIntegrationForPartnerApplicationsSchema,
				ShowOutputAttributeName,
				"enabled",
				"comment",
			),
			ComputedIfAnyAttributeChanged(
				oauthIntegrationForPartnerApplicationsSchema,
				DescribeOutputAttributeName,
				"oauth_client",
				"enabled",
				"oauth_issue_refresh_tokens",
				"oauth_refresh_token_validity",
				"oauth_use_secondary_roles",
				"blocked_roles_list",
				"comment",
			),
		)),

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OauthIntegrationForPartnerApplications, ImportOauthForPartnerApplicationIntegration),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportOauthForPartnerApplicationIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("name", id.Name()); err != nil {
		return nil, err
	}

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("enabled", booleanStringFromBool(integration.Enabled)); err != nil {
		return nil, err
	}

	if issueRefreshTokens, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ISSUE_REFRESH_TOKENS"
	}); err == nil {
		if err = d.Set("oauth_issue_refresh_tokens", issueRefreshTokens.Value); err != nil {
			return nil, err
		}
	}

	if oauthRedirectUri, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_REDIRECT_URI"
	}); err == nil {
		if err = d.Set("oauth_redirect_uri", oauthRedirectUri.Value); err != nil {
			return nil, err
		}
	}

	if refreshTokenValidity, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
	}); err == nil {
		refreshTokenValidityValue, err := strconv.ParseInt(refreshTokenValidity.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		if err = d.Set("oauth_refresh_token_validity", refreshTokenValidityValue); err != nil {
			return nil, err
		}
	}

	if oauthUseSecondaryRoles, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_USE_SECONDARY_ROLES"
	}); err == nil {
		oauthUseSecondaryRolesValue, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(oauthUseSecondaryRoles.Value)
		if err != nil {
			return nil, err
		}
		if err = d.Set("oauth_use_secondary_roles", oauthUseSecondaryRolesValue); err != nil {
			return nil, err
		}
	}

	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "BLOCKED_ROLES_LIST"
	}); err == nil {
		roles := sdk.ParseCommaSeparatedStringArray(prop.Value, false)
		if err = d.Set("blocked_roles_list", roles); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextOauthIntegrationForPartnerApplications(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	oauthClient, err := sdk.ToOauthSecurityIntegrationClientOption(d.Get("oauth_client").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateOauthForPartnerApplicationsSecurityIntegrationRequest(id, oauthClient)

	if v, ok := d.GetOk("oauth_redirect_uri"); ok {
		req.WithOauthRedirectUri(v.(string))
	}

	if v := d.Get("enabled").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithEnabled(parsedBool)
	}

	if v := d.Get("oauth_issue_refresh_tokens").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthIssueRefreshTokens(parsedBool)
	}

	if v := d.Get("oauth_refresh_token_validity").(int); v != IntDefault {
		req.WithOauthRefreshTokenValidity(v)
	}

	if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
		useSecondaryRolesOption, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthUseSecondaryRoles(useSecondaryRolesOption)
	}

	if v, ok := d.GetOk("blocked_roles_list"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		req.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if err := client.SecurityIntegrations.CreateOauthForPartnerApplications(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextOauthIntegrationForPartnerApplications(false)(ctx, d, meta)
}

func ReadContextOauthIntegrationForPartnerApplications(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if c := integration.Category; c != sdk.SecurityIntegrationCategory {
			return diag.FromErr(fmt.Errorf("expected %v to be a %s integration, got %v", id, sdk.SecurityIntegrationCategory, c))
		}
		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}

		oauthClient, err := integration.SubType()
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("oauth_client", oauthClient); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"enabled", "enabled", integration.Enabled, booleanStringFromBool(integration.Enabled), nil},
			); err != nil {
				return diag.FromErr(err)
			}

			oauthIssueRefreshTokens, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_ISSUE_REFRESH_TOKENS"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthRefreshTokenValidity, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			oauthRefreshTokenValidityValue, err := strconv.ParseInt(oauthRefreshTokenValidity.Value, 10, 64)
			if err != nil {
				return diag.FromErr(err)
			}

			oauthUseSecondaryRoles, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_USE_SECONDARY_ROLES"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			blockedRolesList, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "BLOCKED_ROLES_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			// This has to be handled differently as OAUTH_REDIRECT_URI is only visible for a given OAUTH_CLIENT type.
			var oauthRedirectUri string
			if oauthRedirectUriProp, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_REDIRECT_URI"
			}); err == nil {
				oauthRedirectUri = oauthRedirectUriProp.Value
			}

			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"oauth_issue_refresh_tokens", "oauth_issue_refresh_tokens", oauthIssueRefreshTokens.Value, oauthIssueRefreshTokens.Value, nil},
				describeMapping{"oauth_refresh_token_validity", "oauth_refresh_token_validity", oauthRefreshTokenValidity.Value, oauthRefreshTokenValidityValue, nil},
				describeMapping{"oauth_use_secondary_roles", "oauth_use_secondary_roles", oauthUseSecondaryRoles.Value, oauthUseSecondaryRoles.Value, nil},
				describeMapping{"blocked_roles_list", "blocked_roles_list", blockedRolesList.Value, sdk.ParseCommaSeparatedStringArray(blockedRolesList.Value, false), nil},
				describeMapping{"oauth_redirect_uri", "oauth_redirect_uri", oauthRedirectUri, oauthRedirectUri, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, oauthIntegrationForPartnerApplicationsSchema, []string{
			"enabled",
			"oauth_issue_refresh_tokens",
			"oauth_refresh_token_validity",
			"oauth_use_secondary_roles",
			"blocked_roles_list",
			"oauth_redirect_uri",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.DescribeOauthIntegrationForPartnerApplicationsToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}
		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(RelatedParametersAttributeName, []map[string]any{schemas.OauthForPartnerApplicationsParametersToSchema([]*sdk.Parameter{param})}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateContextOauthIntegrationForPartnerApplications(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewOauthForPartnerApplicationsIntegrationSetRequest(), sdk.NewOauthForPartnerApplicationsIntegrationUnsetRequest()

	if d.HasChange("blocked_roles_list") {
		elems := expandStringList(d.Get("blocked_roles_list").(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		set.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
		// can call SET with an empty list
	}

	if d.HasChange("comment") {
		set.WithComment(d.Get("comment").(string))
		// TODO(SNOW-1515781): No UNSET
	}

	if d.HasChange("enabled") {
		if v := d.Get("enabled").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithEnabled(parsedBool)
		} else {
			unset.WithEnabled(true)
		}
	}

	if d.HasChange("oauth_issue_refresh_tokens") {
		if v := d.Get("oauth_issue_refresh_tokens").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthIssueRefreshTokens(parsedBool)
		} else {
			// TODO(SNOW-1515781): No UNSET
			set.WithOauthIssueRefreshTokens(true)
		}
	}

	if d.HasChange("oauth_redirect_uri") {
		// Field can only be set when oauth_client = LOOKER and is required (shouldn't be UNSET in those cases).
		// With any other case oauth_client, the field shouldn't be set.
		set.WithOauthRedirectUri(d.Get("oauth_redirect_uri").(string))
	}

	if d.HasChange("oauth_refresh_token_validity") {
		if v := d.Get("oauth_refresh_token_validity").(int); v != -1 {
			set.WithOauthRefreshTokenValidity(v)
		} else {
			// TODO(SNOW-1515781): No UNSET
			set.WithOauthRefreshTokenValidity(7776000)
		}
	}

	if d.HasChange("oauth_use_secondary_roles") {
		if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
			value, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthUseSecondaryRoles(value)
		} else {
			unset.WithOauthUseSecondaryRoles(true)
		}
	}

	if !reflect.DeepEqual(*set, sdk.OauthForPartnerApplicationsIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, sdk.OauthForPartnerApplicationsIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterOauthForPartnerApplications(ctx, sdk.NewAlterOauthForPartnerApplicationsSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextOauthIntegrationForPartnerApplications(false)(ctx, d, meta)
}
