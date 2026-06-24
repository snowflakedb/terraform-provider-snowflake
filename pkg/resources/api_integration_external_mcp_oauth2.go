package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

var apiIntegrationExternalMcpOAuth2Schema = func() map[string]*schema.Schema {
	apiIntegrationExternalMcpOAuth2 := map[string]*schema.Schema{
		"oauth_client_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the OAuth 2.0 client ID for the MCP server.",
		},
		"oauth_client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: externalChangesNotDetectedFieldDescription("Specifies the OAuth 2.0 client secret for the MCP server."),
		},
		"oauth_token_endpoint": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the OAuth 2.0 token endpoint URL for the MCP server.",
		},
		"oauth_authorization_endpoint": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the OAuth 2.0 authorization endpoint URL for the MCP server.",
		},
		"oauth_client_auth_method": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: sdkValidation(sdk.ToApiIntegrationOauthClientAuthMethod),
			DiffSuppressFunc: NormalizeAndCompare(sdk.ToApiIntegrationOauthClientAuthMethod),
			Description:      fmt.Sprintf("Specifies the OAuth 2.0 client authentication method. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AllApiIntegrationOauthClientAuthMethods))),
		},
		// oauth_discovery_url is intentionally omitted: the field is documented but does not work in Snowflake.
		"oauth_refresh_token_validity": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specifies the validity period (in seconds) for refresh tokens issued by the MCP server.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeExternalMcpOAuth2ApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationExternalMcpOAuth2)
}()

func ApiIntegrationExternalMcpOAuth2() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpOAuth2Resource), TrackingCreateWrapper(resources.ApiIntegrationExternalMcpOAuth2, CreateApiIntegrationExternalMcpOAuth2)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpOAuth2Resource), TrackingReadWrapper(resources.ApiIntegrationExternalMcpOAuth2, ReadApiIntegrationExternalMcpOAuth2)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpOAuth2Resource), TrackingUpdateWrapper(resources.ApiIntegrationExternalMcpOAuth2, UpdateApiIntegrationExternalMcpOAuth2)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpOAuth2Resource), TrackingDeleteWrapper(resources.ApiIntegrationExternalMcpOAuth2, deleteFunc)),
		Description:   "Resource used to manage API integration for external MCP (Model Context Protocol) servers using OAuth 2.0 authentication. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationExternalMcpOAuth2Schema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationExternalMcpOAuth2, ImportApiIntegrationExternalMcpOAuth2),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationExternalMcpOAuth2Schema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationExternalMcpOAuth2Schema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "oauth_client_id", "oauth_client_secret", "oauth_token_endpoint", "oauth_authorization_endpoint", "oauth_client_auth_method", "oauth_refresh_token_validity"),
			// Snowflake retains unspecified optional auth fields in ALTER, so removing them requires recreate.
			customdiff.ForceNewIf("oauth_client_auth_method", func(_ context.Context, d *schema.ResourceDiff, _ any) bool {
				old, new := d.GetChange("oauth_client_auth_method")
				return old.(string) != "" && new.(string) == ""
			}),
			customdiff.ForceNewIf("oauth_refresh_token_validity", func(_ context.Context, d *schema.ResourceDiff, _ any) bool {
				old, new := d.GetChange("oauth_refresh_token_validity")
				return old.(int) != 0 && new.(int) == 0
			}),
		),
	}
}

func buildExternalMcpOAuth2Auth(d *schema.ResourceData) (*sdk.OAuth2McpUserAuthenticationRequest, error) {
	auth := sdk.NewOAuth2McpUserAuthenticationRequest(
		d.Get("oauth_client_id").(string),
		d.Get("oauth_client_secret").(string),
		d.Get("oauth_token_endpoint").(string),
		d.Get("oauth_authorization_endpoint").(string),
	)
	if errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "oauth_client_auth_method", auth.WithOauthClientAuthMethod, sdk.ToApiIntegrationOauthClientAuthMethod),
		intAttributeCreateBuilder(d, "oauth_refresh_token_validity", auth.WithOauthRefreshTokenValidity),
	); errs != nil {
		return nil, errs
	}
	return auth, nil
}

func ImportApiIntegrationExternalMcpOAuth2(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	return importApiIntegrationWithDetails(ctx, d, meta,
		func(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationExternalMcpDetails, error) {
			return client.ApiIntegrations.DescribeExternalMcpDetails(ctx, id)
		},
		func(details *sdk.ApiIntegrationExternalMcpDetails, id sdk.AccountObjectIdentifier) error {
			if _, err := sdk.ToApiIntegrationMcpApiProviderType(details.ApiProvider); err != nil {
				return fmt.Errorf(
					"api integration %s has api_provider %q, not compatible with snowflake_api_integration_external_mcp_oauth2; use the appropriate resource type",
					id.FullyQualifiedName(),
					details.ApiProvider,
				)
			}
			if details.UserAuthType != string(sdk.ApiIntegrationUserAuthTypeOauth2) {
				return fmt.Errorf(
					"api integration %s has user_auth_type %q, not compatible with snowflake_api_integration_external_mcp_oauth2; use the appropriate resource type",
					id.FullyQualifiedName(),
					details.UserAuthType,
				)
			}
			return nil
		},
	)
}

func CreateApiIntegrationExternalMcpOAuth2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	auth, err := buildExternalMcpOAuth2Auth(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithExternalMcpOAuth2ProviderParams(*sdk.NewExternalMcpOAuth2ParamsRequest(*auth))); err != nil {
		return diag.FromErr(fmt.Errorf("error creating external MCP OAuth2 API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationExternalMcpOAuth2(ctx, d, meta)
}

func ReadApiIntegrationExternalMcpOAuth2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	s, err := client.ApiIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query API integration external MCP OAuth2. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration external MCP OAuth2 id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	mcpDetails, err := client.ApiIntegrations.DescribeExternalMcpDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration external MCP OAuth2 (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, mcpDetails.AllowedPrefixes, mcpDetails.BlockedPrefixes),
		d.Set("oauth_client_id", mcpDetails.OauthClientId),
		// oauth_client_secret intentionally omitted: Snowflake returns a masked value and external changes cannot be detected
		d.Set("oauth_token_endpoint", mcpDetails.OauthTokenEndpoint),
		d.Set("oauth_authorization_endpoint", mcpDetails.OauthAuthorizationEndpoint),
		d.Set("oauth_client_auth_method", mcpDetails.OauthClientAuthMethod),
		d.Set("oauth_refresh_token_validity", mcpDetails.OauthRefreshTokenValidity),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationExternalMcpOAuth2DetailsToSchema(mcpDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationExternalMcpOAuth2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()

	if err := handleApiIntegrationCommonUpdate(d, set, unset); err != nil {
		return diag.FromErr(err)
	}

	// If any OAuth2 auth field changed, rebuild the entire ApiUserAuthentication block.
	// There is no UnsetExternalMcpOAuth2Params in the SDK, so all auth changes go through Set.
	if d.HasChanges("oauth_client_id", "oauth_client_secret", "oauth_token_endpoint", "oauth_authorization_endpoint", "oauth_client_auth_method", "oauth_refresh_token_validity") {
		auth, err := buildExternalMcpOAuth2Auth(d)
		if err != nil {
			return diag.FromErr(err)
		}
		mcpSet := sdk.NewSetExternalMcpOAuth2ParamsRequest(*auth)
		set.WithExternalMcpOAuth2Params(*mcpSet)
	}

	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating external MCP OAuth2 API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating external MCP OAuth2 API integration: %w", err))
		}
	}

	return ReadApiIntegrationExternalMcpOAuth2(ctx, d, meta)
}
