package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiAuthCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier (i.e. name) for the integration. This value must be unique in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether this security integration is enabled or disabled.",
	},
	"oauth_client_id": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: "Specifies the client ID for the OAuth application in the external service.",
	},
	"oauth_client_secret": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: externalChangesNotDetectedFieldDescription("Specifies the client secret for the OAuth application in the ServiceNow instance from the previous step. The connector uses this to request an access token from the ServiceNow instance."),
	},
	"oauth_token_endpoint": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the token endpoint used by the client to obtain an access token by presenting its authorization grant or refresh token. The token endpoint is used with every authorization grant except for the implicit grant type (since an access token is issued directly). If removed from the config, the resource is recreated.",
	},
	"oauth_client_auth_method": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption), IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_client_auth_method")),
		Description:      fmt.Sprintf("Specifies that POST is used as the authentication method to the external service. If removed from the config, the resource is recreated. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AllApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption))),
	},
	"oauth_access_token_validity": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateFunc:     validation.IntAtLeast(0),
		Default:          IntDefault,
		Description:      "Specifies the default lifetime of the OAuth access token (in seconds) issued by an OAuth server.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_access_token_validity"),
	},
	"oauth_refresh_token_validity": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateFunc:     validation.IntAtLeast(1),
		Description:      "Specifies the value to determine the validity of the refresh token obtained from the OAuth server.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("oauth_refresh_token_validity"),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
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
			Schema: schemas.DescribeApiAuthSecurityIntegrationSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

type commonApiAuthSet struct {
	enabled                   *bool
	oauthClientId             *string
	oauthClientSecret         *string
	comment                   *string
	oauthAccessTokenValidity  *int
	oauthClientAuthMethod     *sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption
	oauthTokenEndpoint        *string
	oauthRefreshTokenValidity *int
}

type commonApiAuthUnset struct {
	comment *bool
}

func handleApiAuthUpdate(d *schema.ResourceData) (commonApiAuthSet, commonApiAuthUnset, error) {
	set, unset := commonApiAuthSet{}, commonApiAuthUnset{}

	if d.HasChange("enabled") {
		// required field
		set.enabled = sdk.Pointer(d.Get("enabled").(bool))
	}

	if d.HasChange("oauth_client_id") {
		// required field
		set.oauthClientId = sdk.Pointer(d.Get("oauth_client_id").(string))
	}

	if d.HasChange("oauth_client_secret") {
		// required field
		set.oauthClientSecret = sdk.Pointer(d.Get("oauth_client_secret").(string))
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.comment = sdk.Pointer(v.(string))
		} else {
			unset.comment = sdk.Pointer(true)
		}
	}

	if d.HasChange("oauth_access_token_validity") {
		if v := d.Get("oauth_access_token_validity").(int); v != IntDefault {
			set.oauthAccessTokenValidity = sdk.Pointer(v)
		} else {
			// TODO(SNOW-1515781): use UNSET
			set.oauthAccessTokenValidity = sdk.Pointer(0)
		}
	}

	if d.HasChange("oauth_client_auth_method") {
		v := d.Get("oauth_client_auth_method").(string)
		if len(v) > 0 {
			value, err := sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(v)
			if err != nil {
				return commonApiAuthSet{}, commonApiAuthUnset{}, err
			}
			set.oauthClientAuthMethod = sdk.Pointer(value)
		}
		// else: force new
	}

	if d.HasChange("oauth_refresh_token_validity") {
		if v, ok := d.GetOk("oauth_refresh_token_validity"); ok {
			set.oauthRefreshTokenValidity = sdk.Pointer(v.(int))
		} else {
			// TODO(SNOW-1515781): use UNSET
			set.oauthRefreshTokenValidity = sdk.Pointer(7776000)
		}
	}

	if d.HasChange("oauth_token_endpoint") {
		if v, ok := d.GetOk("oauth_token_endpoint"); ok {
			set.oauthTokenEndpoint = sdk.Pointer(v.(string))
		}
		// else: force new
	}
	return set, unset, nil
}

type commonApiAuthCreate struct {
	name                      string
	enabled                   bool
	oauthClientId             string
	oauthClientSecret         string
	comment                   *string
	oauthAccessTokenValidity  *int
	oauthClientAuthMethod     *sdk.ApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption
	oauthTokenEndpoint        *string
	oauthRefreshTokenValidity *int
}

func handleApiAuthCreate(d *schema.ResourceData) (commonApiAuthCreate, error) {
	create := commonApiAuthCreate{
		enabled:           d.Get("enabled").(bool),
		name:              d.Get("name").(string),
		oauthClientId:     d.Get("oauth_client_id").(string),
		oauthClientSecret: d.Get("oauth_client_secret").(string),
	}
	if v, ok := d.GetOk("comment"); ok {
		create.comment = sdk.Pointer(v.(string))
	}

	if v := d.Get("oauth_access_token_validity").(int); v != IntDefault {
		create.oauthAccessTokenValidity = sdk.Pointer(v)
	}
	if v, ok := d.GetOk("oauth_refresh_token_validity"); ok {
		create.oauthRefreshTokenValidity = sdk.Pointer(v.(int))
	}
	if v, ok := d.GetOk("oauth_token_endpoint"); ok {
		create.oauthTokenEndpoint = sdk.Pointer(v.(string))
	}
	if v, ok := d.GetOk("oauth_client_auth_method"); ok {
		value, err := sdk.ToApiAuthenticationSecurityIntegrationOauthClientAuthMethodOption(v.(string))
		if err != nil {
			return commonApiAuthCreate{}, err
		}
		create.oauthClientAuthMethod = sdk.Pointer(value)
	}

	return create, nil
}

func handleApiAuthImport(d *schema.ResourceData, integration *sdk.SecurityIntegration,
	properties []sdk.SecurityIntegrationProperty,
) error {
	if _, err := ImportName[sdk.AccountObjectIdentifier](context.Background(), d, nil); err != nil {
		return err
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return err
	}
	if err := d.Set("comment", integration.Comment); err != nil {
		return err
	}

	oauthAccessTokenValidity, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_ACCESS_TOKEN_VALIDITY"
	})
	if err == nil {
		value, err := strconv.Atoi(oauthAccessTokenValidity.Value)
		if err != nil {
			return err
		}
		if err = d.Set("oauth_access_token_validity", value); err != nil {
			return err
		}
	}
	oauthRefreshTokenValidity, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
	})
	if err == nil {
		value, err := strconv.Atoi(oauthRefreshTokenValidity.Value)
		if err != nil {
			return err
		}
		if err = d.Set("oauth_refresh_token_validity", value); err != nil {
			return err
		}
	}
	oauthClientId, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_CLIENT_ID" })
	if err == nil {
		if err = d.Set("oauth_client_id", oauthClientId.Value); err != nil {
			return err
		}
	}
	oauthClientAuthMethod, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_CLIENT_AUTH_METHOD"
	})
	if err == nil {
		if err = d.Set("oauth_client_auth_method", oauthClientAuthMethod.Value); err != nil {
			return err
		}
	}
	oauthTokenEndpoint, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_TOKEN_ENDPOINT" })
	if err == nil {
		if err = d.Set("oauth_token_endpoint", oauthTokenEndpoint.Value); err != nil {
			return err
		}
	}

	return nil
}

func handleApiAuthRead(d *schema.ResourceData,
	id sdk.AccountObjectIdentifier,
	integration *sdk.SecurityIntegration,
	properties []sdk.SecurityIntegrationProperty,
	withExternalChangesMarking bool,
	extraFieldsDescribeMappings []describeMapping,
) error {
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return err
	}
	if err := d.Set("comment", integration.Comment); err != nil {
		return err
	}
	if err := d.Set("enabled", integration.Enabled); err != nil {
		return err
	}
	if withExternalChangesMarking {
		oauthAccessTokenValidity, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_ACCESS_TOKEN_VALIDITY"
		})
		if err != nil {
			return err
		}

		oauthRefreshTokenValidity, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_REFRESH_TOKEN_VALIDITY"
		})
		if err != nil {
			return err
		}

		oauthClientId, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_CLIENT_ID" })
		if err != nil {
			return err
		}

		oauthClientAuthMethod, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_CLIENT_AUTH_METHOD"
		})
		if err != nil {
			return err
		}

		oauthTokenEndpoint, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_TOKEN_ENDPOINT" })
		if err != nil {
			return err
		}

		oauthAccessTokenValidityInt, err := strconv.Atoi(oauthAccessTokenValidity.Value)
		if err != nil {
			return err
		}
		oauthRefreshTokenValidityInt, err := strconv.Atoi(oauthRefreshTokenValidity.Value)
		if err != nil {
			return err
		}

		if err = handleExternalChangesToObjectInDescribe(d,
			append(extraFieldsDescribeMappings,
				describeMapping{"oauth_access_token_validity", "oauth_access_token_validity", oauthAccessTokenValidity.Value, oauthAccessTokenValidityInt, nil},
				describeMapping{"oauth_refresh_token_validity", "oauth_refresh_token_validity", oauthRefreshTokenValidity.Value, oauthRefreshTokenValidityInt, nil},
				describeMapping{"oauth_client_id", "oauth_client_id", oauthClientId.Value, oauthClientId.Value, nil},
				describeMapping{"oauth_client_auth_method", "oauth_client_auth_method", oauthClientAuthMethod.Value, oauthClientAuthMethod.Value, nil},
				describeMapping{"oauth_token_endpoint", "oauth_token_endpoint", oauthTokenEndpoint.Value, oauthTokenEndpoint.Value, nil},
			)...,
		); err != nil {
			return err
		}
	}
	if err := setStateToValuesFromConfig(d, apiAuthCommonSchema, []string{
		"oauth_access_token_validity",
		"oauth_refresh_token_validity",
		"oauth_client_id",
		"oauth_client_auth_method",
		"oauth_token_endpoint",
	}); err != nil {
		return err
	}

	if err := d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
		return err
	}

	if err := d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiAuthSecurityIntegrationPropertiesToSchema(properties)}); err != nil {
		return err
	}
	return nil
}
