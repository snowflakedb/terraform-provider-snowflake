package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var privilegedRoles = []string{"ACCOUNTADMIN", "ORGADMIN", "SECURITYADMIN"}

var externalOauthIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the name of the External Oath integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"external_oauth_type": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      fmt.Sprintf("Specifies the OAuth 2.0 authorization server to be Okta, Microsoft Azure AD, Ping Identity PingFederate, or a Custom OAuth 2.0 authorization server. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationTypes))),
		ValidateDiagFunc: sdkValidation(sdk.ToExternalOauthSecurityIntegrationTypeOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToExternalOauthSecurityIntegrationTypeOption),
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether to initiate operation of the integration or suspend it.",
	},
	"external_oauth_issuer": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the URL to define the OAuth 2.0 authorization server.",
	},
	"external_oauth_token_user_mapping_claim": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies the access token claim or claims that can be used to map the access token to a Snowflake user record. If removed from the config, the resource is recreated.",
	},
	"external_oauth_snowflake_user_mapping_attribute": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      fmt.Sprintf("Indicates which Snowflake user record attribute should be used to map the access token to a Snowflake user record. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationSnowflakeUserMappingAttributes))),
		ValidateDiagFunc: sdkValidation(sdk.ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption),
	},
	"external_oauth_jws_keys_url": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		MaxItems:      3,
		Optional:      true,
		ConflictsWith: []string{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2"},
		Description:   "Specifies the endpoint or a list of endpoints from which to download public keys or certificates to validate an External OAuth access token. The maximum number of URLs that can be specified in the list is 3. If removed from the config, the resource is recreated.",
	},
	"external_oauth_rsa_public_key": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies a Base64-encoded RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers. If removed from the config, the resource is recreated.",
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
		ConflictsWith:    []string{"external_oauth_jws_keys_url"},
	},
	"external_oauth_rsa_public_key_2": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies a second RSA public key, without the -----BEGIN PUBLIC KEY----- and -----END PUBLIC KEY----- headers. Used for key rotation. If removed from the config, the resource is recreated.",
		DiffSuppressFunc: ignoreTrimSpaceSuppressFunc,
		ConflictsWith:    []string{"external_oauth_jws_keys_url"},
	},
	"external_oauth_blocked_roles_list": {
		Type:             schema.TypeSet,
		Elem:             &schema.Schema{Type: schema.TypeString},
		Optional:         true,
		Description:      relatedResourceDescription(withPrivilegedRolesDescription("Specifies the list of roles that a client cannot set as the primary role.", string(sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList)), resources.AccountRole),
		DiffSuppressFunc: IgnoreValuesFromSetIfParamSet("external_oauth_blocked_roles_list", string(sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList), privilegedRoles),
		ConflictsWith:    []string{"external_oauth_allowed_roles_list"},
	},
	"external_oauth_allowed_roles_list": {
		Type:             schema.TypeSet,
		Elem:             &schema.Schema{Type: schema.TypeString},
		Optional:         true,
		Description:      relatedResourceDescription("Specifies the list of roles that the client can set as the primary role.", resources.AccountRole),
		DiffSuppressFunc: SuppressIfAny(
		// TODO(SNOW-1517937): uncomment
		// NormalizeAndCompareIdentifiersInSet("external_oauth_allowed_roles_list"),
		),
		ConflictsWith: []string{"external_oauth_blocked_roles_list"},
	},
	"external_oauth_audience_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies additional values that can be used for the access token's audience validation on top of using the Customer's Snowflake Account URL ",
	},
	"external_oauth_any_role_mode": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Specifies whether the OAuth client or user can use a role that is not defined in the OAuth access token. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AsStringList(sdk.AllExternalOauthSecurityIntegrationAnyRoleModes)))),
		ValidateDiagFunc: sdkValidation(sdk.ToExternalOauthSecurityIntegrationAnyRoleModeOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToExternalOauthSecurityIntegrationAnyRoleModeOption),
	},
	"external_oauth_scope_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the scope delimiter in the authorization token.",
	},
	"external_oauth_scope_mapping_attribute": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the access token claim to map the access token to an account role. If removed from the config, the resource is recreated.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the OAuth integration.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecurityIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeExternalOauthSecurityIntegrationSchema,
		},
	},
	RelatedParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Parameters related to this security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowExternalOauthParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func ExternalOauthIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ExternalOauthSecurityIntegration, CreateContextExternalOauthIntegration),
		ReadContext:   TrackingReadWrapper(resources.ExternalOauthSecurityIntegration, ReadContextExternalOauthIntegration(true)),
		UpdateContext: TrackingUpdateWrapper(resources.ExternalOauthSecurityIntegration, UpdateContextExternalOauthIntegration),
		DeleteContext: TrackingDeleteWrapper(resources.ExternalOauthSecurityIntegration, DeleteSecurityIntegration),
		Description:   "Resource used to manage external oauth security integration objects. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-oauth-external).",

		Schema: externalOauthIntegrationSchema,
		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalOauthSecurityIntegration, customdiff.All(
			ForceNewIfChangeToEmptyString("external_oauth_rsa_public_key"),
			ForceNewIfChangeToEmptyString("external_oauth_rsa_public_key_2"),
			ForceNewIfChangeToEmptyString("external_oauth_scope_mapping_attribute"),
			ForceNewIfChangeToEmptySet("external_oauth_jws_keys_url"),
			ForceNewIfChangeToEmptySet("external_oauth_token_user_mapping_claim"),
			ComputedIfAnyAttributeChanged(externalOauthIntegrationSchema, ShowOutputAttributeName, "enabled", "external_oauth_type", "comment"),
			ComputedIfAnyAttributeChanged(externalOauthIntegrationSchema, DescribeOutputAttributeName, "enabled", "external_oauth_issuer", "external_oauth_jws_keys_url", "external_oauth_any_role_mode",
				"external_oauth_rsa_public_key", "external_oauth_rsa_public_key_2", "external_oauth_blocked_roles_list", "external_oauth_allowed_roles_list",
				"external_oauth_audience_list", "external_oauth_token_user_mapping_claim", "external_oauth_snowflake_user_mapping_attribute", "external_oauth_scope_delimiter",
				"comment"),
		)),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalOauthSecurityIntegration, ImportExternalOauthIntegration),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v092ExternalOauthIntegrationStateUpgrader,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

	if err = d.Set("enabled", integration.Enabled); err != nil {
		return nil, err
	}
	if oauthType, err := integration.SubType(); err == nil {
		if err = d.Set("external_oauth_type", oauthType); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "EXTERNAL_OAUTH_ISSUER" }); err == nil {
		if err = d.Set("external_oauth_issuer", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_JWS_KEYS_URL"
	}); err == nil {
		if err = d.Set("external_oauth_jws_keys_url", sdk.ParseCommaSeparatedStringArray(prop.Value, false)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_ANY_ROLE_MODE"
	}); err == nil {
		if err = d.Set("external_oauth_any_role_mode", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY"
	}); err == nil {
		if err = d.Set("external_oauth_rsa_public_key", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"
	}); err == nil {
		if err = d.Set("external_oauth_rsa_public_key_2", prop.Value); err != nil {
			return nil, err
		}
	}

	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"
	}); err == nil {
		roles := sdk.ParseCommaSeparatedStringArray(prop.Value, false)
		if err = d.Set("external_oauth_blocked_roles_list", roles); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"
	}); err == nil {
		if err = d.Set("external_oauth_allowed_roles_list", sdk.ParseCommaSeparatedStringArray(prop.Value, false)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_AUDIENCE_LIST"
	}); err == nil {
		if err = d.Set("external_oauth_audience_list", sdk.ParseCommaSeparatedStringArray(prop.Value, false)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"
	}); err == nil {
		if err = d.Set("external_oauth_token_user_mapping_claim", sdk.ParseCommaSeparatedStringArray(prop.Value, true)); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"
	}); err == nil {
		if err = d.Set("external_oauth_snowflake_user_mapping_attribute", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "EXTERNAL_OAUTH_SCOPE_DELIMITER"
	}); err == nil {
		if err = d.Set("external_oauth_scope_delimiter", prop.Value); err != nil {
			return nil, err
		}
	}
	if prop, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "COMMENT" }); err == nil {
		if err = d.Set("comment", prop.Value); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	enabled := d.Get("enabled").(bool)
	externalOauthIssuer := d.Get("external_oauth_issuer").(string)
	externalOauthSnowflakeUserMappingAttributeRaw := d.Get("external_oauth_snowflake_user_mapping_attribute").(string)
	externalOauthTokenUserMappingClaimRaw := expandStringList(d.Get("external_oauth_token_user_mapping_claim").(*schema.Set).List())
	name := d.Get("name").(string)
	integrationTypeRaw := d.Get("external_oauth_type").(string)
	integrationType, err := sdk.ToExternalOauthSecurityIntegrationTypeOption(integrationTypeRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	externalOauthSnowflakeUserMappingAttribute, err := sdk.ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption(externalOauthSnowflakeUserMappingAttributeRaw)
	if err != nil {
		return diag.FromErr(err)
	}
	externalOauthTokenUserMappingClaim := make([]sdk.TokenUserMappingClaim, 0, len(externalOauthTokenUserMappingClaimRaw))
	for _, v := range externalOauthTokenUserMappingClaimRaw {
		externalOauthTokenUserMappingClaim = append(externalOauthTokenUserMappingClaim, sdk.TokenUserMappingClaim{Claim: v})
	}
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateExternalOauthSecurityIntegrationRequest(id, enabled, integrationType, externalOauthIssuer, externalOauthTokenUserMappingClaim, externalOauthSnowflakeUserMappingAttribute)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_allowed_roles_list"); ok {
		vList := expandStringList(v.(*schema.Set).List())
		allowedRoles := make([]sdk.AccountObjectIdentifier, len(vList))
		for i := range vList {
			allowedRoles[i] = sdk.NewAccountObjectIdentifier(vList[i])
		}
		req.WithExternalOauthAllowedRolesList(sdk.AllowedRolesListRequest{AllowedRolesList: allowedRoles})
	}

	if v, ok := d.GetOk("external_oauth_any_role_mode"); ok {
		valueRaw := v.(string)
		value, err := sdk.ToExternalOauthSecurityIntegrationAnyRoleModeOption(valueRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithExternalOauthAnyRoleMode(value)
	}

	if v, ok := d.GetOk("external_oauth_audience_list"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		audienceUrls := make([]sdk.AudienceListItem, len(elems))
		for i := range elems {
			audienceUrls[i] = sdk.AudienceListItem{Item: elems[i]}
		}
		req.WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: audienceUrls})
	}

	if v, ok := d.GetOk("external_oauth_blocked_roles_list"); ok {
		vList := expandStringList(v.(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(vList))
		for i := range vList {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(vList[i])
		}
		req.WithExternalOauthBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if v, ok := d.GetOk("external_oauth_jws_keys_url"); ok {
		elems := expandStringList(v.(*schema.Set).List())
		urls := make([]sdk.JwsKeysUrl, len(elems))
		for i := range elems {
			urls[i] = sdk.JwsKeysUrl{JwsKeyUrl: elems[i]}
		}
		req.WithExternalOauthJwsKeysUrl(urls)
	}

	if v, ok := d.GetOk("external_oauth_rsa_public_key"); ok {
		req.WithExternalOauthRsaPublicKey(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_rsa_public_key_2"); ok {
		req.WithExternalOauthRsaPublicKey2(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_scope_delimiter"); ok {
		req.WithExternalOauthScopeDelimiter(v.(string))
	}

	if v, ok := d.GetOk("external_oauth_scope_mapping_attribute"); ok {
		req.WithExternalOauthScopeMappingAttribute(v.(string))
	}

	if err := client.SecurityIntegrations.CreateExternalOauth(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextExternalOauthIntegration(false)(ctx, d, meta)
}

func ReadContextExternalOauthIntegration(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		if err := d.Set("enabled", integration.Enabled); err != nil {
			return diag.FromErr(err)
		}
		subType, err := integration.SubType()
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("external_oauth_type", subType); err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			externalOauthIssuer, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "EXTERNAL_OAUTH_ISSUER" })
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthJwsKeysUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_JWS_KEYS_URL"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthAnyRoleMode, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_ANY_ROLE_MODE"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthRsaPublicKey, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthRsaPublicKey2, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthBlockedRolesList, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_BLOCKED_ROLES_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthAllowedRolesList, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_ALLOWED_ROLES_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthAudienceList, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_AUDIENCE_LIST"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthTokenUserMappingClaim, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthSnowflakeUserMappingAttribute, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			externalOauthScopeDelimiter, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "EXTERNAL_OAUTH_SCOPE_DELIMITER"
			})
			if err != nil {
				return diag.FromErr(err)
			}
			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"external_oauth_issuer", "external_oauth_issuer", externalOauthIssuer.Value, externalOauthIssuer.Value, nil},
				describeMapping{"external_oauth_jws_keys_url", "external_oauth_jws_keys_url", externalOauthJwsKeysUrl.Value, sdk.ParseCommaSeparatedStringArray(externalOauthJwsKeysUrl.Value, false), nil},
				describeMapping{"external_oauth_any_role_mode", "external_oauth_any_role_mode", externalOauthAnyRoleMode.Value, externalOauthAnyRoleMode.Value, nil},
				describeMapping{"external_oauth_rsa_public_key", "external_oauth_rsa_public_key", externalOauthRsaPublicKey.Value, externalOauthRsaPublicKey.Value, nil},
				describeMapping{"external_oauth_rsa_public_key_2", "external_oauth_rsa_public_key_2", externalOauthRsaPublicKey2.Value, externalOauthRsaPublicKey2.Value, nil},
				describeMapping{"external_oauth_blocked_roles_list", "external_oauth_blocked_roles_list", externalOauthBlockedRolesList.Value, sdk.ParseCommaSeparatedStringArray(externalOauthBlockedRolesList.Value, false), nil},
				describeMapping{"external_oauth_allowed_roles_list", "external_oauth_allowed_roles_list", externalOauthAllowedRolesList.Value, sdk.ParseCommaSeparatedStringArray(externalOauthAllowedRolesList.Value, false), nil},
				describeMapping{"external_oauth_audience_list", "external_oauth_audience_list", externalOauthAudienceList.Value, sdk.ParseCommaSeparatedStringArray(externalOauthAudienceList.Value, false), nil},
				describeMapping{"external_oauth_token_user_mapping_claim", "external_oauth_token_user_mapping_claim", externalOauthTokenUserMappingClaim.Value, sdk.ParseCommaSeparatedStringArray(externalOauthTokenUserMappingClaim.Value, true), nil},
				describeMapping{"external_oauth_snowflake_user_mapping_attribute", "external_oauth_snowflake_user_mapping_attribute", externalOauthSnowflakeUserMappingAttribute.Value, externalOauthSnowflakeUserMappingAttribute.Value, nil},
				describeMapping{"external_oauth_scope_delimiter", "external_oauth_scope_delimiter", externalOauthScopeDelimiter.Value, externalOauthScopeDelimiter.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, externalOauthIntegrationSchema, []string{
			"external_oauth_jws_keys_url",
			"external_oauth_rsa_public_key",
			"external_oauth_rsa_public_key_2",
			"external_oauth_blocked_roles_list",
			"external_oauth_allowed_roles_list",
			"external_oauth_audience_list",
			"external_oauth_any_role_mode",
			"external_oauth_scope_delimiter",
			"external_oauth_scope_mapping_attribute",
			"comment",
		}); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ExternalOauthSecurityIntegrationPropertiesToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}

		param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set(RelatedParametersAttributeName, []map[string]any{schemas.ExternalOauthParametersToSchema([]*sdk.Parameter{param})}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateContextExternalOauthIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewExternalOauthIntegrationSetRequest(), sdk.NewExternalOauthIntegrationUnsetRequest()

	errs := errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if d.HasChange("enabled") {
		// this field is required
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("external_oauth_allowed_roles_list") {
		v := expandStringList(d.Get("external_oauth_allowed_roles_list").(*schema.Set).List())
		allowedRoles := make([]sdk.AccountObjectIdentifier, len(v))
		for i := range v {
			allowedRoles[i] = sdk.NewAccountObjectIdentifier(v[i])
		}
		set.WithExternalOauthAllowedRolesList(sdk.AllowedRolesListRequest{AllowedRolesList: allowedRoles})
	}

	if d.HasChange("external_oauth_any_role_mode") {
		v := d.Get("external_oauth_any_role_mode").(string)
		if len(v) > 0 {
			value, err := sdk.ToExternalOauthSecurityIntegrationAnyRoleModeOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithExternalOauthAnyRoleMode(value)
		} else {
			// TODO(SNOW-1515781): use UNSET
			set.WithExternalOauthAnyRoleMode(sdk.ExternalOauthSecurityIntegrationAnyRoleModeDisable)
		}
	}

	if d.HasChange("external_oauth_audience_list") {
		v := expandStringList(d.Get("external_oauth_audience_list").(*schema.Set).List())
		if len(v) > 0 {
			audienceList := make([]sdk.AudienceListItem, len(v))
			for i := range v {
				audienceList[i] = sdk.AudienceListItem{Item: v[i]}
			}
			set.WithExternalOauthAudienceList(sdk.AudienceListRequest{AudienceList: audienceList})
		} else {
			unset.WithExternalOauthAudienceList(true)
		}
	}

	if d.HasChange("external_oauth_blocked_roles_list") {
		vRaw := d.Get("external_oauth_blocked_roles_list")
		v := expandStringList(vRaw.(*schema.Set).List())
		blockedRoles := make([]sdk.AccountObjectIdentifier, len(v))
		for i := range v {
			blockedRoles[i] = sdk.NewAccountObjectIdentifier(v[i])
		}
		set.WithExternalOauthBlockedRolesList(sdk.BlockedRolesListRequest{BlockedRolesList: blockedRoles})
	}

	if d.HasChange("external_oauth_issuer") {
		// this field is required
		set.WithExternalOauthIssuer(d.Get("external_oauth_issuer").(string))
	}

	if d.HasChange("external_oauth_jws_keys_url") {
		v := expandStringList(d.Get("external_oauth_jws_keys_url").(*schema.Set).List())
		if len(v) > 0 {
			urls := make([]sdk.JwsKeysUrl, len(v))
			for i := range v {
				urls[i] = sdk.JwsKeysUrl{JwsKeyUrl: v[i]}
			}
			set.WithExternalOauthJwsKeysUrl(urls)
		}
		// else: force new
	}

	if d.HasChange("external_oauth_rsa_public_key") {
		set.WithExternalOauthRsaPublicKey(d.Get("external_oauth_rsa_public_key").(string))
		if v, ok := d.GetOk("external_oauth_rsa_public_key"); ok {
			set.WithExternalOauthRsaPublicKey2(v.(string))
		}
		// else: force new
	}

	if d.HasChange("external_oauth_rsa_public_key_2") {
		if v, ok := d.GetOk("external_oauth_rsa_public_key_2"); ok {
			set.WithExternalOauthRsaPublicKey2(v.(string))
		}
		// else: force new
	}

	if d.HasChange("external_oauth_scope_delimiter") {
		if v, ok := d.GetOk("external_oauth_scope_delimiter"); ok {
			set.WithExternalOauthScopeDelimiter(v.(string))
		} else {
			// TODO(SNOW-1515781): use UNSET
			set.WithExternalOauthScopeDelimiter(",")
		}
	}

	if d.HasChange("external_oauth_scope_mapping_attribute") {
		// this field is required
		set.WithExternalOauthScopeMappingAttribute(d.Get("external_oauth_scope_mapping_attribute").(string))
	}

	if d.HasChange("external_oauth_snowflake_user_mapping_attribute") {
		// this field is required
		if v, ok := d.GetOk("external_oauth_snowflake_user_mapping_attribute"); ok {
			value, err := sdk.ToExternalOauthSecurityIntegrationSnowflakeUserMappingAttributeOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithExternalOauthSnowflakeUserMappingAttribute(value)
		}
	}

	if d.HasChange("external_oauth_token_user_mapping_claim") {
		v := expandStringList(d.Get("external_oauth_token_user_mapping_claim").(*schema.Set).List())
		claims := make([]sdk.TokenUserMappingClaim, len(v))
		for i := range v {
			claims[i] = sdk.TokenUserMappingClaim{
				Claim: v[i],
			}
		}
		set.WithExternalOauthTokenUserMappingClaim(claims)
	}

	if d.HasChange("external_oauth_type") {
		// this field is required
		if v, ok := d.GetOk("external_oauth_type"); ok {
			value, err := sdk.ToExternalOauthSecurityIntegrationTypeOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithExternalOauthType(value)
		}
	}

	if !reflect.DeepEqual(*set, sdk.ExternalOauthIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterExternalOauth(ctx, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.ExternalOauthIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterExternalOauth(ctx, sdk.NewAlterExternalOauthSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextExternalOauthIntegration(false)(ctx, d, meta)
}
