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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var apiAuthJwtBearerSchema = func() map[string]*schema.Schema {
	apiAuthJwtBearer := map[string]*schema.Schema{
		"oauth_authorization_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the URL for authenticating to the external service.",
		},
		"oauth_assertion_issuer": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return collections.MergeMaps(apiAuthCommonSchema, apiAuthJwtBearer)
}()

func ApiAuthenticationIntegrationWithJwtBearer() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ApiAuthenticationIntegrationWithJwtBearer, CreateContextApiAuthenticationIntegrationWithJwtBearer),
		ReadContext:   TrackingReadWrapper(resources.ApiAuthenticationIntegrationWithJwtBearer, ReadContextApiAuthenticationIntegrationWithJwtBearer(true)),
		UpdateContext: TrackingUpdateWrapper(resources.ApiAuthenticationIntegrationWithJwtBearer, UpdateContextApiAuthenticationIntegrationWithJwtBearer),
		DeleteContext: TrackingDeleteWrapper(resources.ApiAuthenticationIntegrationWithJwtBearer, DeleteSecurityIntegration),
		Description:   "Resource used to manage api authentication security integration objects with jwt bearer. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-api-auth).",

		Schema: apiAuthJwtBearerSchema,
		CustomizeDiff: TrackingCustomDiffWrapper(resources.ApiAuthenticationIntegrationWithJwtBearer, customdiff.All(
			ForceNewIfChangeToEmptyString("oauth_token_endpoint"),
			ForceNewIfChangeToEmptyString("oauth_authorization_endpoint"),
			ForceNewIfChangeToEmptyString("oauth_client_auth_method"),
			ComputedIfAnyAttributeChanged(apiAuthJwtBearerSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiAuthJwtBearerSchema, DescribeOutputAttributeName, "enabled", "comment", "oauth_access_token_validity", "oauth_refresh_token_validity",
				"oauth_client_auth_method", "oauth_authorization_endpoint", "oauth_token_endpoint", "oauth_assertion_issuer")),
		),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiAuthenticationIntegrationWithJwtBearer, ImportApiAuthenticationWithJwtBearer),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportApiAuthenticationWithJwtBearer(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	properties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := handleApiAuthImport(d, integration, properties); err != nil {
		return nil, err
	}
	oauthAuthorizationEndpoint, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
		return property.Name == "OAUTH_AUTHORIZATION_ENDPOINT"
	})
	if err == nil {
		if err = d.Set("oauth_authorization_endpoint", oauthAuthorizationEndpoint.Value); err != nil {
			return nil, err
		}
	}
	oauthAssertionIssuer, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ASSERTION_ISSUER" })
	if err == nil {
		if err = d.Set("oauth_assertion_issuer", oauthAssertionIssuer.Value); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextApiAuthenticationIntegrationWithJwtBearer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	commonCreate, err := handleApiAuthCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := sdk.ParseAccountObjectIdentifier(commonCreate.name)
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id, commonCreate.enabled, d.Get("oauth_assertion_issuer").(string), commonCreate.oauthClientId, commonCreate.oauthClientSecret)
	req.WithOauthGrantJwtBearer(true)
	req.Comment = commonCreate.comment
	req.OauthAccessTokenValidity = commonCreate.oauthAccessTokenValidity
	req.OauthRefreshTokenValidity = commonCreate.oauthRefreshTokenValidity
	req.OauthTokenEndpoint = commonCreate.oauthTokenEndpoint
	req.OauthClientAuthMethod = commonCreate.oauthClientAuthMethod

	if v, ok := d.GetOk("oauth_authorization_endpoint"); ok {
		req.WithOauthAuthorizationEndpoint(v.(string))
	}

	if err := client.SecurityIntegrations.CreateApiAuthenticationWithJwtBearerFlow(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadContextApiAuthenticationIntegrationWithJwtBearer(false)(ctx, d, meta)
}

func ReadContextApiAuthenticationIntegrationWithJwtBearer(withExternalChangesMarking bool) schema.ReadContextFunc {
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
		properties, err := client.SecurityIntegrations.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if c := integration.Category; c != sdk.SecurityIntegrationCategory {
			return diag.FromErr(fmt.Errorf("expected %v to be a %s integration, got %v", id, sdk.SecurityIntegrationCategory, c))
		}
		oauthAuthorizationEndpoint, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool {
			return property.Name == "OAUTH_AUTHORIZATION_ENDPOINT"
		})
		if err != nil {
			return diag.FromErr(err)
		}

		oauthAssertionIssuer, err := collections.FindFirst(properties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "OAUTH_ASSERTION_ISSUER" })
		if err != nil {
			return diag.FromErr(err)
		}
		if err := handleApiAuthRead(d, id, integration, properties, withExternalChangesMarking, []describeMapping{
			{"oauth_authorization_endpoint", "oauth_authorization_endpoint", oauthAuthorizationEndpoint.Value, oauthAuthorizationEndpoint.Value, nil},
			{"oauth_assertion_issuer", "oauth_assertion_issuer", oauthAssertionIssuer.Value, oauthAssertionIssuer.Value, nil},
		}); err != nil {
			return diag.FromErr(err)
		}
		if err := setStateToValuesFromConfig(d, apiAuthJwtBearerSchema, []string{
			"oauth_authorization_endpoint",
			"oauth_assertion_issuer",
		}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func UpdateContextApiAuthenticationIntegrationWithJwtBearer(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	commonSet, commonUnset, err := handleApiAuthUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}
	set := &sdk.ApiAuthenticationWithJwtBearerFlowIntegrationSetRequest{
		Enabled:                   commonSet.enabled,
		OauthTokenEndpoint:        commonSet.oauthTokenEndpoint,
		OauthClientAuthMethod:     commonSet.oauthClientAuthMethod,
		OauthClientId:             commonSet.oauthClientId,
		OauthClientSecret:         commonSet.oauthClientSecret,
		OauthAccessTokenValidity:  commonSet.oauthAccessTokenValidity,
		OauthRefreshTokenValidity: commonSet.oauthRefreshTokenValidity,
		Comment:                   commonSet.comment,
	}
	unset := &sdk.ApiAuthenticationWithJwtBearerFlowIntegrationUnsetRequest{
		Comment: commonUnset.comment,
	}
	if d.HasChange("oauth_authorization_endpoint") {
		set.WithOauthAuthorizationEndpoint(d.Get("oauth_authorization_endpoint").(string))
	}
	if !reflect.DeepEqual(*set, sdk.ApiAuthenticationWithJwtBearerFlowIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterApiAuthenticationWithJwtBearerFlow(ctx, sdk.NewAlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, sdk.ApiAuthenticationWithJwtBearerFlowIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterApiAuthenticationWithJwtBearerFlow(ctx, sdk.NewAlterApiAuthenticationWithJwtBearerFlowSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadContextApiAuthenticationIntegrationWithJwtBearer(false)(ctx, d, meta)
}
