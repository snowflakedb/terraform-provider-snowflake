package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var saml2IntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the name of the SAML2 integration. This name follows the rules for Object Identifiers. The name should be unique among security integrations in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"enabled": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("enabled"),
		Description:      booleanStringFieldDescription("Specifies whether this security integration is enabled or disabled."),
	},
	"saml2_issuer": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The string containing the IdP EntityID / Issuer.",
	},
	"saml2_sso_url": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The string containing the IdP SSO URL, where the user should be redirected by Snowflake (the Service Provider) with a SAML AuthnRequest message.",
	},
	"saml2_provider": {
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToSaml2SecurityIntegrationSaml2ProviderOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToSaml2SecurityIntegrationSaml2ProviderOption),
		Description:      fmt.Sprintf("The string describing the IdP. Valid options are: %v.", possibleValuesListed(sdk.AllSaml2SecurityIntegrationSaml2Providers)),
	},
	"saml2_x509_cert": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: "The Base64 encoded IdP signing certificate on a single line without the leading -----BEGIN CERTIFICATE----- and ending -----END CERTIFICATE----- markers.",
	},
	"saml2_sp_initiated_login_page_label": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_sp_initiated_login_page_label"),
		Description:      "The string containing the label to display after the Log In With button on the login page. If this field changes value from non-empty to empty, the whole resource is recreated because of Snowflake limitations.",
	},
	"saml2_enable_sp_initiated": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_enable_sp_initiated"),
		Description:      booleanStringFieldDescription("The Boolean indicating if the Log In With button will be shown on the login page. TRUE: displays the Log in With button on the login page. FALSE: does not display the Log in With button on the login page."),
	},
	"saml2_sign_request": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_sign_request"),
		Description:      booleanStringFieldDescription("The Boolean indicating whether SAML requests are signed. TRUE: allows SAML requests to be signed. FALSE: does not allow SAML requests to be signed."),
	},
	"saml2_requested_nameid_format": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToSaml2SecurityIntegrationSaml2RequestedNameidFormatOption),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToSaml2SecurityIntegrationSaml2RequestedNameidFormatOption), IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_requested_nameid_format")),
		Description:      fmt.Sprintf("The SAML NameID format allows Snowflake to set an expectation of the identifying attribute of the user (i.e. SAML Subject) in the SAML assertion from the IdP to ensure a valid authentication to Snowflake. Valid options are: %v.", possibleValuesListed(sdk.AllSaml2SecurityIntegrationSaml2RequestedNameidFormats)),
	},
	"saml2_post_logout_redirect_url": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_post_logout_redirect_url"),
		Description:      "The endpoint to which Snowflake redirects users after clicking the Log Out button in the classic Snowflake web interface. Snowflake terminates the Snowflake session upon redirecting to the specified endpoint.",
	},
	"saml2_force_authn": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_force_authn"),
		Description:      booleanStringFieldDescription("The Boolean indicating whether users, during the initial authentication flow, are forced to authenticate again to access Snowflake. When set to TRUE, Snowflake sets the ForceAuthn SAML parameter to TRUE in the outgoing request from Snowflake to the identity provider. TRUE: forces users to authenticate again to access Snowflake, even if a valid session with the identity provider exists. FALSE: does not force users to authenticate again to access Snowflake."),
	},
	"saml2_snowflake_issuer_url": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_snowflake_issuer_url"),
		Description:      "The string containing the EntityID / Issuer for the Snowflake service provider. If an incorrect value is specified, Snowflake returns an error message indicating the acceptable values to use. Because Okta does not support underscores in URLs, the underscore in the account name must be converted to a hyphen. See [docs](https://docs.snowflake.com/en/user-guide/organizations-connect#okta-urls).",
	},
	"saml2_snowflake_acs_url": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("saml2_snowflake_acs_url"),
		Description:      "The string containing the Snowflake Assertion Consumer Service URL to which the IdP will send its SAML authentication response back to Snowflake. This property will be set in the SAML authentication request generated by Snowflake when initiating a SAML SSO operation with the IdP. If an incorrect value is specified, Snowflake returns an error message indicating the acceptable values to use. Because Okta does not support underscores in URLs, the underscore in the account name must be converted to a hyphen. See [docs](https://docs.snowflake.com/en/user-guide/organizations-connect#okta-urls).",
	},
	"allowed_user_domains": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "A list of email domains that can authenticate with a SAML2 security integration. If this field changes value from non-empty to empty, the whole resource is recreated because of Snowflake limitations.",
	},
	"allowed_email_patterns": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "A list of regular expressions that email addresses are matched against to authenticate with a SAML2 security integration. If this field changes value from non-empty to empty, the whole resource is recreated because of Snowflake limitations.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
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
			Schema: schemas.DescribeSaml2IntegrationSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func SAML2Integration() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.Saml2SecurityIntegration, CreateContextSAML2Integration),
		ReadContext:   TrackingReadWrapper(resources.Saml2SecurityIntegration, ReadContextSAML2Integration(true)),
		UpdateContext: TrackingUpdateWrapper(resources.Saml2SecurityIntegration, UpdateContextSAML2Integration),
		DeleteContext: TrackingDeleteWrapper(resources.Saml2SecurityIntegration, DeleteSecurityIntegration),
		Description:   "Resource used to manage SAML2 security integration objects. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-saml2).",

		Schema: saml2IntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Saml2SecurityIntegration, ImportSaml2Integration),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Saml2SecurityIntegration, customdiff.All(
			ForceNewIfChangeToEmptySet("allowed_user_domains"),
			ForceNewIfChangeToEmptySet("allowed_email_patterns"),
			ForceNewIfChangeToEmptyString("saml2_snowflake_issuer_url"),
			ForceNewIfChangeToEmptyString("saml2_snowflake_acs_url"),
			ForceNewIfChangeToEmptyString("saml2_sp_initiated_login_page_label"),
			ComputedIfAnyAttributeChanged(saml2IntegrationSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(saml2IntegrationSchema, DescribeOutputAttributeName, "saml2_issuer", "saml2_sso_url", "saml2_provider",
				"saml2_sp_initiated_login_page_label", "saml2_enable_sp_initiated", "saml2_sign_request", "saml2_requtedted_nameid_format",
				"saml2_post_logout_redirect_url", "saml2_force_authn", "saml2_snowflake_issuer_url", "saml2_snowflake_acs_url", "allowed_user_domains",
				"allowed_email_patterns"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func ImportSaml2Integration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

	if err := d.Set("comment", integration.Comment); err != nil {
		return nil, err
	}
	if err := d.Set("enabled", booleanStringFromBool(integration.Enabled)); err != nil {
		return nil, err
	}

	samlIssuer, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_ISSUER" })
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 saml issuer, err = %w", err)
	}
	if err := d.Set("saml2_issuer", samlIssuer.Value); err != nil {
		return nil, err
	}

	ssoUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_SSO_URL" })
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 sso url, err = %w", err)
	}
	if err := d.Set("saml2_sso_url", ssoUrl.Value); err != nil {
		return nil, err
	}

	samlProvider, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_PROVIDER" })
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 provider, err = %w", err)
	}
	samlProviderValue, err := sdk.ToSaml2SecurityIntegrationSaml2ProviderOption(samlProvider.Value)
	if err != nil {
		return nil, err
	}
	if err := d.Set("saml2_provider", samlProviderValue); err != nil {
		return nil, err
	}

	x509Cert, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_X509_CERT" })
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 x509 cert, err = %w", err)
	}
	if err := d.Set("saml2_x509_cert", x509Cert.Value); err != nil {
		return nil, err
	}

	spInitiatedLoginPageLabel, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 sp initiated login page label, err = %w", err)
	}
	if err := d.Set("saml2_sp_initiated_login_page_label", spInitiatedLoginPageLabel.Value); err != nil {
		return nil, err
	}

	enableSpInitiated, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_ENABLE_SP_INITIATED"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 enable sp initiated, err = %w", err)
	}
	if err := d.Set("saml2_enable_sp_initiated", enableSpInitiated.Value); err != nil {
		return nil, err
	}

	signRequest, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_SIGN_REQUEST"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 sign request, err = %w", err)
	}
	if err := d.Set("saml2_sign_request", signRequest.Value); err != nil {
		return nil, err
	}

	requestedNameIdFormat, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_REQUESTED_NAMEID_FORMAT"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 requested nameid format, err = %w", err)
	}
	if err := d.Set("saml2_requested_nameid_format", requestedNameIdFormat.Value); err != nil {
		return nil, err
	}

	postLogoutRedirectUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_POST_LOGOUT_REDIRECT_URL"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 post logout redirect url, err = %w", err)
	}
	if err := d.Set("saml2_post_logout_redirect_url", postLogoutRedirectUrl.Value); err != nil {
		return nil, err
	}

	forceAuthn, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_FORCE_AUTHN"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 force authn, err = %w", err)
	}
	if err := d.Set("saml2_force_authn", forceAuthn.Value); err != nil {
		return nil, err
	}

	snowflakeIssuerUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_SNOWFLAKE_ISSUER_URL"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 snowflake issuer url, err = %w", err)
	}
	if err := d.Set("saml2_snowflake_issuer_url", snowflakeIssuerUrl.Value); err != nil {
		return nil, err
	}

	snowflakeAcsUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "SAML2_SNOWFLAKE_ACS_URL"
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find saml2 snowflake acs url, err = %w", err)
	}
	if err := d.Set("saml2_snowflake_acs_url", snowflakeAcsUrl.Value); err != nil {
		return nil, err
	}

	if err := d.Set("allowed_user_domains", getOptionalListField(integrationProperties, "ALLOWED_USER_DOMAINS")); err != nil {
		return nil, err
	}

	if err := d.Set("allowed_email_patterns", getOptionalListField(integrationProperties, "ALLOWED_EMAIL_PATTERNS")); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func getOptionalListField(props []sdk.SecurityIntegrationProperty, propName string) []string {
	found, err := collections.FindFirst(props, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == propName
	})
	if err != nil {
		log.Printf("[DEBUG] failed to find %s in object properties, err = %v", propName, err)
		return make([]string, 0)
	}
	return sdk.ParseCommaSeparatedStringArray(found.Value, false)
}

func CreateContextSAML2Integration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	samlProvider, err := sdk.ToSaml2SecurityIntegrationSaml2ProviderOption(d.Get("saml2_provider").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateSaml2SecurityIntegrationRequest(
		id,
		d.Get("saml2_issuer").(string),
		d.Get("saml2_sso_url").(string),
		samlProvider,
		d.Get("saml2_x509_cert").(string),
	)

	if v := d.Get("enabled").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithEnabled(parsed)
	}

	if v, ok := d.GetOk("saml2_sp_initiated_login_page_label"); ok {
		req.WithSaml2SpInitiatedLoginPageLabel(v.(string))
	}

	if v := d.Get("saml2_enable_sp_initiated").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithSaml2EnableSpInitiated(parsed)
	}

	if v := d.Get("saml2_sign_request").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithSaml2SignRequest(parsed)
	}

	if v, ok := d.GetOk("saml2_requested_nameid_format"); ok {
		format, err := sdk.ToSaml2SecurityIntegrationSaml2RequestedNameidFormatOption(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithSaml2RequestedNameidFormat(format)
	}

	if v, ok := d.GetOk("saml2_post_logout_redirect_url"); ok {
		req.WithSaml2PostLogoutRedirectUrl(v.(string))
	}

	if v := d.Get("saml2_force_authn").(string); v != BooleanDefault {
		parsed, err := booleanStringToBool(v)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithSaml2ForceAuthn(parsed)
	}

	if v, ok := d.GetOk("saml2_snowflake_issuer_url"); ok {
		req.WithSaml2SnowflakeIssuerUrl(v.(string))
	}

	if v, ok := d.GetOk("saml2_snowflake_acs_url"); ok {
		req.WithSaml2SnowflakeAcsUrl(v.(string))
	}

	if v, ok := d.GetOk("allowed_user_domains"); ok {
		stringAllowedUserDomains := expandStringList(v.(*schema.Set).List())
		allowedUserDomains := make([]sdk.UserDomain, len(stringAllowedUserDomains))
		for i, v := range stringAllowedUserDomains {
			allowedUserDomains[i] = sdk.UserDomain{
				Domain: v,
			}
		}
		req.WithAllowedUserDomains(allowedUserDomains)
	}

	if v, ok := d.GetOk("allowed_email_patterns"); ok {
		stringAllowedEmailPatterns := expandStringList(v.(*schema.Set).List())
		allowedEmailPatterns := make([]sdk.EmailPattern, len(stringAllowedEmailPatterns))
		for i, v := range stringAllowedEmailPatterns {
			allowedEmailPatterns[i] = sdk.EmailPattern{
				Pattern: v,
			}
		}
		req.WithAllowedEmailPatterns(allowedEmailPatterns)
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if err := client.SecurityIntegrations.CreateSaml2(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSAML2Integration(false)(ctx, d, meta)
}

func ReadContextSAML2Integration(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		samlIssuer, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_ISSUER" })
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find saml2 saml issuer, err = %w", err))
		}
		if err := d.Set("saml2_issuer", samlIssuer.Value); err != nil {
			return diag.FromErr(err)
		}

		ssoUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_SSO_URL" })
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find saml2 sso url, err = %w", err))
		}
		if err := d.Set("saml2_sso_url", ssoUrl.Value); err != nil {
			return diag.FromErr(err)
		}

		samlProvider, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_PROVIDER" })
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find saml2 provider, err = %w", err))
		}
		samlProviderValue, err := sdk.ToSaml2SecurityIntegrationSaml2ProviderOption(samlProvider.Value)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("saml2_provider", samlProviderValue); err != nil {
			return diag.FromErr(err)
		}

		x509Cert, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SAML2_X509_CERT" })
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find saml2 x509 cert, err = %w", err))
		}
		if err := d.Set("saml2_x509_cert", x509Cert.Value); err != nil {
			return diag.FromErr(err)
		}

		postLogoutRedirectUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "SAML2_POST_LOGOUT_REDIRECT_URL"
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to find saml2 post logout redirect url, err = %w", err))
		}
		if err := d.Set("saml2_post_logout_redirect_url", postLogoutRedirectUrl.Value); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("allowed_user_domains", getOptionalListField(integrationProperties, "ALLOWED_USER_DOMAINS")); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("allowed_email_patterns", getOptionalListField(integrationProperties, "ALLOWED_EMAIL_PATTERNS")); err != nil {
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

			enableSpInitiated, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_ENABLE_SP_INITIATED"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 enable sp initiated, err = %w", err))
			}

			signRequest, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_SIGN_REQUEST"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 sign request, err = %w", err))
			}

			requestedNameIdFormat, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_REQUESTED_NAMEID_FORMAT"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 requested nameid format, err = %w", err))
			}

			forceAuthn, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_FORCE_AUTHN"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 force authn, err = %w", err))
			}

			snowflakeIssuerUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_SNOWFLAKE_ISSUER_URL"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 snowflake issuer url, err = %w", err))
			}

			snowflakeAcsUrl, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_SNOWFLAKE_ACS_URL"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 snowflake acs url, err = %w", err))
			}

			spInitiatedLoginPageLabel, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool {
				return property.Name == "SAML2_SP_INITIATED_LOGIN_PAGE_LABEL"
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to find saml2 sp initiated login page label, err = %w", err))
			}

			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"saml2_enable_sp_initiated", "saml2_enable_sp_initiated", enableSpInitiated.Value, enableSpInitiated.Value, nil},
				describeMapping{"saml2_sign_request", "saml2_sign_request", signRequest.Value, signRequest.Value, nil},
				describeMapping{"saml2_requested_nameid_format", "saml2_requested_nameid_format", requestedNameIdFormat.Value, requestedNameIdFormat.Value, nil},
				describeMapping{"saml2_force_authn", "saml2_force_authn", forceAuthn.Value, forceAuthn.Value, nil},
				describeMapping{"saml2_snowflake_acs_url", "saml2_snowflake_acs_url", snowflakeAcsUrl.Value, snowflakeAcsUrl.Value, nil},
				describeMapping{"saml2_snowflake_issuer_url", "saml2_snowflake_issuer_url", snowflakeIssuerUrl.Value, snowflakeIssuerUrl.Value, nil},
				describeMapping{"saml2_sp_initiated_login_page_label", "saml2_sp_initiated_login_page_label", spInitiatedLoginPageLabel.Value, spInitiatedLoginPageLabel.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, saml2IntegrationSchema, []string{
			"enabled",
			"saml2_enable_sp_initiated",
			"saml2_sign_request",
			"saml2_requested_nameid_format",
			"saml2_force_authn",
			"saml2_snowflake_acs_url",
			"saml2_snowflake_issuer_url",
			"saml2_sp_initiated_login_page_label",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.DescribeSaml2IntegrationToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextSAML2Integration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewSaml2IntegrationSetRequest(), sdk.NewSaml2IntegrationUnsetRequest()

	if d.HasChange("enabled") {
		if v := d.Get("enabled").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithEnabled(parsed)
		} else {
			// TODO(SNOW-1515781): UNSET is not implemented
			set.WithEnabled(false)
		}
	}

	if d.HasChange("saml2_issuer") {
		set.WithSaml2Issuer(d.Get("saml2_issuer").(string))
	}

	if d.HasChange("saml2_sso_url") {
		set.WithSaml2SsoUrl(d.Get("saml2_sso_url").(string))
	}

	if d.HasChange("saml2_provider") {
		valueRaw := d.Get("saml2_provider").(string)
		value, err := sdk.ToSaml2SecurityIntegrationSaml2ProviderOption(valueRaw)
		if err != nil {
			return diag.FromErr(err)
		}
		set.WithSaml2Provider(value)
	}

	if d.HasChange("saml2_x509_cert") {
		set.WithSaml2X509Cert(d.Get("saml2_x509_cert").(string))
	}

	if d.HasChange("saml2_sp_initiated_login_page_label") {
		// TODO(SNOW-1515781): UNSET is not implemented and SET with empty value is invalid (conditional ForceNew on unset)
		set.WithSaml2SpInitiatedLoginPageLabel(d.Get("saml2_sp_initiated_login_page_label").(string))
	}

	if d.HasChange("saml2_enable_sp_initiated") {
		if v := d.Get("saml2_enable_sp_initiated").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithSaml2EnableSpInitiated(parsed)
		} else {
			// TODO(SNOW-1515781): UNSET is not implemented
			set.WithSaml2EnableSpInitiated(false)
		}
	}

	if d.HasChange("saml2_sign_request") {
		if v := d.Get("saml2_sign_request").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithSaml2SignRequest(parsed)
		} else {
			// TODO(SNOW-1515781): UNSET is not implemented
			set.WithSaml2SignRequest(false)
		}
	}

	if d.HasChange("saml2_requested_nameid_format") {
		if v, ok := d.GetOk("saml2_requested_nameid_format"); ok {
			value, err := sdk.ToSaml2SecurityIntegrationSaml2RequestedNameidFormatOption(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithSaml2RequestedNameidFormat(value)
		} else {
			unset.WithSaml2RequestedNameidFormat(true)
		}
	}

	if d.HasChange("saml2_post_logout_redirect_url") {
		if v, ok := d.GetOk("saml2_post_logout_redirect_url"); ok {
			set.WithSaml2PostLogoutRedirectUrl(v.(string))
		} else {
			unset.WithSaml2PostLogoutRedirectUrl(true)
		}
	}

	if d.HasChange("saml2_force_authn") {
		if v := d.Get("saml2_force_authn").(string); v != BooleanDefault {
			parsed, err := booleanStringToBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithSaml2ForceAuthn(parsed)
		} else {
			// TODO(SNOW-1515781): UNSET is not implemented
			set.WithSaml2ForceAuthn(false)
		}
	}

	if d.HasChange("saml2_snowflake_issuer_url") {
		// TODO(SNOW-1515781): UNSET is not implemented and SET with empty value is invalid (conditional ForceNew on unset)
		set.WithSaml2SnowflakeIssuerUrl(d.Get("saml2_snowflake_issuer_url").(string))
	}

	if d.HasChange("saml2_snowflake_acs_url") {
		// TODO(SNOW-1515781): UNSET is not implemented and SET with empty value is invalid (conditional ForceNew on unset)
		set.WithSaml2SnowflakeAcsUrl(d.Get("saml2_snowflake_acs_url").(string))
	}

	if d.HasChange("allowed_user_domains") {
		// TODO(SNOW-1515781): UNSET is not implemented and SET with empty list is invalid (conditional ForceNew on non-empty to empty set)
		v := d.Get("allowed_user_domains").(*schema.Set).List()
		userDomains := make([]sdk.UserDomain, len(v))
		for i := range v {
			userDomains[i] = sdk.UserDomain{
				Domain: v[i].(string),
			}
		}
		set.WithAllowedUserDomains(userDomains)
	}

	if d.HasChange("allowed_email_patterns") {
		// TODO(SNOW-SNOW-1515781): UNSET is not implemented and SET with empty list is invalid (conditional ForceNew on non-empty to empty set)
		v := d.Get("allowed_email_patterns").(*schema.Set).List()
		emailPatterns := make([]sdk.EmailPattern, len(v))
		for i := range v {
			emailPatterns[i] = sdk.EmailPattern{
				Pattern: v[i].(string),
			}
		}
		set.WithAllowedEmailPatterns(emailPatterns)
	}

	if d.HasChange("comment") {
		if v := d.Get("comment").(string); len(v) > 0 {
			set.WithComment(v)
		} else {
			unset.WithComment(true)
		}
	}

	if !reflect.DeepEqual(*set, sdk.Saml2IntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterSaml2(ctx, sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.Saml2IntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterSaml2(ctx, sdk.NewAlterSaml2SecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSAML2Integration(false)(ctx, d, meta)
}

func DeleteContextSAM2LIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.SecurityIntegrations.Drop(ctx, sdk.NewDropSecurityIntegrationRequest(sdk.NewAccountObjectIdentifier(id.Name())).WithIfExists(true))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error deleting integration",
				Detail:   fmt.Sprintf("id %v err = %v", id.Name(), err),
			},
		}
	}

	d.SetId("")
	return nil
}
