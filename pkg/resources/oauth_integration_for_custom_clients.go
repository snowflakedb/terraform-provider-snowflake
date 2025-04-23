package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

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

var oauthIntegrationForCustomClientsSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the name of the OAuth integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"oauth_client_type": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToOauthSecurityIntegrationClientTypeOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToOauthSecurityIntegrationClientTypeOption),
		Description:      fmt.Sprintf("Specifies the type of client being registered. Snowflake supports both confidential and public clients. Valid options are: %v.", possibleValuesListed(sdk.AllOauthSecurityIntegrationClientTypes)),
	},
	"oauth_redirect_uri": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: "Specifies the client URI. After a user is authenticated, the web browser is redirected to this URI.",
	},
	"enabled": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("enabled"),
		Description:      booleanStringFieldDescription("Specifies whether this OAuth integration is enabled or disabled."),
	},
	"oauth_allow_non_tls_redirect_uri": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("oauth_allow_non_tls_redirect_uri"),
		Description:      booleanStringFieldDescription("If true, allows setting oauth_redirect_uri to a URI not protected by TLS."),
	},
	"oauth_enforce_pkce": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("oauth_enforce_pkce"),
		Description:      booleanStringFieldDescription("Boolean that specifies whether Proof Key for Code Exchange (PKCE) should be required for the integration."),
	},
	"oauth_use_secondary_roles": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption),
		Description:      fmt.Sprintf("Specifies whether default secondary roles set in the user properties are activated by default in the session being opened. Valid options are: %v.", possibleValuesListed(sdk.AllOauthSecurityIntegrationUseSecondaryRoles)),
	},
	"pre_authorized_roles_list": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:    true,
		Description: relatedResourceDescription("A set of Snowflake roles that a user does not need to explicitly consent to using after authenticating.", resources.AccountRole),
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
	"oauth_issue_refresh_tokens": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("oauth_issue_refresh_tokens"),
		Description:      booleanStringFieldDescription("Specifies whether to allow the client to exchange a refresh token for an access token when the current access token has expired."),
	},
	"oauth_refresh_token_validity": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      IntDefault,
		ValidateFunc: validation.IntAtLeast(0),
		Description:  "Specifies how long refresh tokens should be valid (in seconds). OAUTH_ISSUE_REFRESH_TOKENS must be set to TRUE.",
	},
	"network_policy": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("Specifies an existing network policy. This network policy controls network traffic that is attempting to exchange an authorization code for an access or refresh token or to use a refresh token to obtain a new access token.", resources.NetworkPolicy),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"oauth_client_rsa_public_key": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
		Description:      externalChangesNotDetectedFieldDescription("Specifies a Base64-encoded RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers."),
	},
	"oauth_client_rsa_public_key_2": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
		Description:      externalChangesNotDetectedFieldDescription("Specifies a Base64-encoded RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers."),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the OAuth integration.",
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
			Schema: schemas.DescribeOauthIntegrationForCustomClients,
		},
	},
	RelatedParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Parameters related to this security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowOauthForCustomClientsParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func OauthIntegrationForCustomClients() *schema.Resource {
	return &schema.Resource{
		Schema: oauthIntegrationForCustomClientsSchema,

		CreateContext: TrackingCreateWrapper(resources.OauthIntegrationForCustomClients, CreateContextOauthIntegrationForCustomClients),
		ReadContext:   TrackingReadWrapper(resources.OauthIntegrationForCustomClients, ReadContextOauthIntegrationForCustomClients(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OauthIntegrationForCustomClients, UpdateContextOauthIntegrationForCustomClients),
		DeleteContext: TrackingDeleteWrapper(resources.OauthIntegrationForCustomClients, DeleteSecurityIntegration),
		Description:   "Resource used to manage oauth security integration for custom clients objects. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-snowflake).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OauthIntegrationForCustomClients, customdiff.All(
			ComputedIfAnyAttributeChanged(
				oauthIntegrationForCustomClientsSchema,
				ShowOutputAttributeName,
				"enabled",
				"comment",
			),
			ComputedIfAnyAttributeChanged(
				oauthIntegrationForCustomClientsSchema,
				DescribeOutputAttributeName,
				"oauth_client_type",
				"enabled",
				"oauth_allow_non_tls_redirect_uri",
				"oauth_enforce_pkce",
				"oauth_use_secondary_roles",
				"pre_authorized_roles_list",
				"blocked_roles_list",
				"oauth_issue_refresh_tokens",
				"oauth_refresh_token_validity",
				"network_policy",
				"oauth_client_rsa_public_key",
				"oauth_client_rsa_public_key_2",
				"comment",
			),
		)),

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OauthIntegrationForCustomClients, ImportOauthForCustomClientsIntegration),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportOauthForCustomClientsIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

	if allowNonTlsRedirectUri, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ALLOW_NON_TLS_REDIRECT_URI"
	}); err == nil {
		if err = d.Set("oauth_allow_non_tls_redirect_uri", allowNonTlsRedirectUri.Value); err != nil {
			return nil, err
		}
	}

	if enforcePkce, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ENFORCE_PKCE"
	}); err == nil {
		if err = d.Set("oauth_enforce_pkce", enforcePkce.Value); err != nil {
			return nil, err
		}
	}

	if issueRefreshTokens, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ISSUE_REFRESH_TOKENS"
	}); err == nil {
		if err = d.Set("oauth_issue_refresh_tokens", issueRefreshTokens.Value); err != nil {
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

func CreateContextOauthIntegrationForCustomClients(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	oauthClientType, err := sdk.ToOauthSecurityIntegrationClientTypeOption(d.Get("oauth_client_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateOauthForCustomClientsSecurityIntegrationRequest(id, oauthClientType, d.Get("oauth_redirect_uri").(string))

	if v := d.Get("enabled").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithEnabled(parsedBool)
	}

	if v := d.Get("oauth_allow_non_tls_redirect_uri").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthAllowNonTlsRedirectUri(parsedBool)
	}

	if v := d.Get("oauth_enforce_pkce").(string); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthEnforcePkce(parsedBool)
	}

	if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
		oauthUseSecondaryRoles, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithOauthUseSecondaryRoles(oauthUseSecondaryRoles)
	}

	if v, ok := d.GetOk("pre_authorized_roles_list"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		preAuthorizedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			preAuthorizedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		req.WithPreAuthorizedRolesList(sdk.PreAuthorizedRolesListRequest{PreAuthorizedRolesList: preAuthorizedRoles})
	}

	if v, ok := d.GetOk("blocked_roles_list"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		req.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
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

	if v, ok := d.GetOk("network_policy"); ok {
		req.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	if v, ok := d.GetOk("oauth_client_rsa_public_key"); ok {
		req.WithOauthClientRsaPublicKey(v.(string))
	}

	if v, ok := d.GetOk("oauth_client_rsa_public_key_2"); ok {
		req.WithOauthClientRsaPublicKey2(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if err := client.SecurityIntegrations.CreateOauthForCustomClients(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextOauthIntegrationForCustomClients(false)(ctx, d, meta)
}

func ReadContextOauthIntegrationForCustomClients(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		integration, err := client.SecurityIntegrations.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration id: %s, Err: %s", id.FullyQualifiedName(), err),
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

		if err := d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}

		oauthClientType, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_CLIENT_TYPE"
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find oauth client type, err = %w", err))
		}
		if err := d.Set("oauth_client_type", oauthClientType.Value); err != nil {
			return diag.FromErr(err)
		}

		oauthRedirectUri, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_REDIRECT_URI"
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find oauth redirect uri, err = %w", err))
		}
		if err := d.Set("oauth_redirect_uri", oauthRedirectUri.Value); err != nil {
			return diag.FromErr(err)
		}

		preAuthorizedRolesList, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "PRE_AUTHORIZED_ROLES_LIST"
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find pre authorized roles list, err = %w", err))
		}
		var preAuthorizedRoles []string
		if len(preAuthorizedRolesList.Value) > 0 {
			preAuthorizedRoles = strings.Split(preAuthorizedRolesList.Value, ",")
		}
		if err := d.Set("pre_authorized_roles_list", preAuthorizedRoles); err != nil {
			return diag.FromErr(err)
		}

		networkPolicy, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "NETWORK_POLICY"
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find network policy, err = %w", err))
		}
		if err := d.Set("network_policy", sdk.NewAccountObjectIdentifier(networkPolicy.Value).Name()); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"enabled", "enabled", integration.Enabled, booleanStringFromBool(integration.Enabled), nil},
			); err != nil {
				return diag.FromErr(err)
			}

			oauthAllowNonTlsRedirectUri, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_ALLOW_NON_TLS_REDIRECT_URI"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthEnforcePkce, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_ENFORCE_PKCE"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			oauthUseSecondaryRoles, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "OAUTH_USE_SECONDARY_ROLES"
			})
			if err != nil {
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

			blockedRolesList, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "BLOCKED_ROLES_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}

			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"oauth_allow_non_tls_redirect_uri", "oauth_allow_non_tls_redirect_uri", oauthAllowNonTlsRedirectUri.Value, oauthAllowNonTlsRedirectUri.Value, nil},
				describeMapping{"oauth_enforce_pkce", "oauth_enforce_pkce", oauthEnforcePkce.Value, oauthEnforcePkce.Value, nil},
				describeMapping{"oauth_use_secondary_roles", "oauth_use_secondary_roles", oauthUseSecondaryRoles.Value, oauthUseSecondaryRoles.Value, nil},
				describeMapping{"oauth_issue_refresh_tokens", "oauth_issue_refresh_tokens", oauthIssueRefreshTokens.Value, oauthIssueRefreshTokens.Value, nil},
				describeMapping{"oauth_refresh_token_validity", "oauth_refresh_token_validity", oauthRefreshTokenValidity.Value, oauthRefreshTokenValidity.Value, nil},
				describeMapping{"blocked_roles_list", "blocked_roles_list", blockedRolesList.Value, sdk.ParseCommaSeparatedStringArray(blockedRolesList.Value, false), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, oauthIntegrationForCustomClientsSchema, []string{
			"enabled",
			"oauth_allow_non_tls_redirect_uri",
			"oauth_enforce_pkce",
			"oauth_use_secondary_roles",
			"oauth_issue_refresh_tokens",
			"oauth_refresh_token_validity",
			"blocked_roles_list",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.DescribeOauthIntegrationForCustomClientsToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}
		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(RelatedParametersAttributeName, []map[string]any{schemas.OauthForCustomClientsParametersToSchema([]*sdk.Parameter{param})}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextOauthIntegrationForCustomClients(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewOauthForCustomClientsIntegrationSetRequest(), sdk.NewOauthForCustomClientsIntegrationUnsetRequest()

	if d.HasChange("oauth_redirect_uri") {
		set.WithOauthRedirectUri(d.Get("oauth_redirect_uri").(string))
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

	if d.HasChange("oauth_allow_non_tls_redirect_uri") {
		if v := d.Get("oauth_allow_non_tls_redirect_uri").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthAllowNonTlsRedirectUri(parsedBool)
		} else {
			// TODO(SNOW-1515781): No unset available for this field (setting Snowflake default)
			set.WithOauthAllowNonTlsRedirectUri(false)
		}
	}

	if d.HasChange("oauth_enforce_pkce") {
		if v := d.Get("oauth_enforce_pkce").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthEnforcePkce(parsedBool)
		} else {
			// TODO(SNOW-1515781): No unset available for this field (setting Snowflake default)
			set.WithOauthEnforcePkce(false)
		}
	}

	if d.HasChange("oauth_use_secondary_roles") {
		if v, ok := d.GetOk("oauth_use_secondary_roles"); ok {
			oauthUseSecondaryRoles, err := sdk.ToOauthSecurityIntegrationUseSecondaryRolesOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthUseSecondaryRoles(oauthUseSecondaryRoles)
		} else {
			unset.WithOauthUseSecondaryRoles(true)
		}
	}

	if d.HasChange("pre_authorized_roles_list") {
		elems := expandStringList(d.Get("pre_authorized_roles_list").(*schema.Set).List())
		preAuthorizedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			preAuthorizedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		set.WithPreAuthorizedRolesList(sdk.PreAuthorizedRolesListRequest{PreAuthorizedRolesList: preAuthorizedRoles})
	}

	if d.HasChange("blocked_roles_list") {
		elems := expandStringList(d.Get("blocked_roles_list").(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(elems))
		for i := range elems {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(elems[i])
		}
		set.WithBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if d.HasChange("oauth_issue_refresh_tokens") {
		if v := d.Get("oauth_issue_refresh_tokens").(string); v != BooleanDefault {
			parsedBool, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithOauthIssueRefreshTokens(parsedBool)
		} else {
			// TODO(SNOW-1515781): No unset available for this field (setting Snowflake default)
			set.WithOauthIssueRefreshTokens(true)
		}
	}

	if d.HasChange("oauth_refresh_token_validity") {
		if v := d.Get("oauth_refresh_token_validity").(int); v != IntDefault {
			set.WithOauthRefreshTokenValidity(v)
		} else {
			// TODO(SNOW-1515781): No unset available for this field (setting Snowflake default; 90 days in seconds)
			set.WithOauthRefreshTokenValidity(7_776_000)
		}
	}

	if d.HasChange("network_policy") {
		if v, ok := d.GetOk("network_policy"); ok {
			set.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v.(string)))
		} else {
			unset.WithNetworkPolicy(true)
		}
	}

	if d.HasChange("oauth_client_rsa_public_key") {
		if v, ok := d.GetOk("oauth_client_rsa_public_key"); ok {
			set.WithOauthClientRsaPublicKey(v.(string))
		} else {
			unset.WithOauthClientRsaPublicKey(true)
		}
	}

	if d.HasChange("oauth_client_rsa_public_key_2") {
		if v, ok := d.GetOk("oauth_client_rsa_public_key_2"); ok {
			set.WithOauthClientRsaPublicKey2(v.(string))
		} else {
			unset.WithOauthClientRsaPublicKey2(true)
		}
	}

	if d.HasChange("comment") {
		set.WithComment(d.Get("comment").(string))
	}

	if !reflect.DeepEqual(*set, sdk.OauthForCustomClientsIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterOauthForCustomClients(ctx, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, sdk.OauthForCustomClientsIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterOauthForCustomClients(ctx, sdk.NewAlterOauthForCustomClientsSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextOauthIntegrationForCustomClients(false)(ctx, d, meta)
}
