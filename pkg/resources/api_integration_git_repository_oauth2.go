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

var apiIntegrationGitRepositoryOauth2Schema = func() map[string]*schema.Schema {
	apiIntegrationGitRepositoryOauth2 := map[string]*schema.Schema{
		"oauth_authorization_endpoint": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The OAuth 2.0 authorization endpoint for the Git repository.",
		},
		"oauth_token_endpoint": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The OAuth 2.0 token endpoint for the Git repository.",
		},
		"oauth_client_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The client ID for the OAuth 2.0 application.",
		},
		"oauth_client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			ForceNew:    true,
			Description: externalChangesNotDetectedFieldDescription("The client secret for the OAuth 2.0 application."),
		},
		"oauth_access_token_validity": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Specifies the validity period (in seconds) for the OAuth 2.0 access token.",
		},
		"oauth_refresh_token_validity": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Specifies the validity period (in seconds) for the OAuth 2.0 refresh token.",
		},
		"oauth_allowed_scopes": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "Specifies a list of scopes to use when making a request from the OAuth by a role with USAGE on the integration. " + enumValuesDescription(sdk.AllApiIntegrationOauthAllowedScopes),
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: sdkValidation(sdk.ToApiIntegrationOauthAllowedScope),
			},
		},
		"oauth_username": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the username to authenticate with the Git repository using OAuth 2.0.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeGitRepositoryOauth2ApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationGitRepositoryOauth2)
}()

func ApiIntegrationGitRepositoryOauth2() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ApiIntegrationGitRepositoryOauth2Resource), TrackingCreateWrapper(resources.ApiIntegrationGitRepositoryOauth2, CreateApiIntegrationGitRepositoryOauth2)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ApiIntegrationGitRepositoryOauth2Resource), TrackingReadWrapper(resources.ApiIntegrationGitRepositoryOauth2, ReadApiIntegrationGitRepositoryOauth2)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ApiIntegrationGitRepositoryOauth2Resource), TrackingUpdateWrapper(resources.ApiIntegrationGitRepositoryOauth2, UpdateApiIntegrationGitRepositoryOauth2)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ApiIntegrationGitRepositoryOauth2Resource), TrackingDeleteWrapper(resources.ApiIntegrationGitRepositoryOauth2, deleteFunc)),
		Description:   "Resource used to manage API integration Git Repository OAuth2 objects. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationGitRepositoryOauth2Schema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationGitRepositoryOauth2, ImportApiIntegrationGitRepositoryOauth2),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryOauth2Schema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryOauth2Schema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment"),
		),
	}
}

func ImportApiIntegrationGitRepositoryOauth2(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	return importApiIntegrationWithDetails(ctx, d, meta,
		func(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationGitHttpsApiDetails, error) {
			return client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
		},
		func(details *sdk.ApiIntegrationGitHttpsApiDetails, id sdk.AccountObjectIdentifier) error {
			if details.UserAuthType != string(sdk.ApiIntegrationUserAuthTypeOauth2) {
				return fmt.Errorf(
					"api integration %s has user auth type %s, not compatible with snowflake_api_integration_git_repository_oauth2 (expected %s); use the appropriate resource type",
					id.FullyQualifiedName(),
					details.UserAuthType,
					sdk.ApiIntegrationUserAuthTypeOauth2,
				)
			}
			return nil
		},
	)
}

func CreateApiIntegrationGitRepositoryOauth2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	authReq := sdk.NewOAuth2GitUserAuthenticationRequest(
		d.Get("oauth_authorization_endpoint").(string),
		d.Get("oauth_token_endpoint").(string),
		d.Get("oauth_client_id").(string),
		d.Get("oauth_client_secret").(string),
	)

	if scopesRaw, ok := d.GetOk("oauth_allowed_scopes"); ok {
		scopes, err := collections.MapErr(scopesRaw.([]any), func(scope any) (sdk.ApiIntegrationOauthAllowedScope, error) {
			return sdk.ToApiIntegrationOauthAllowedScope(scope.(string))
		})
		if err != nil {
			return diag.FromErr(err)
		}

		authReq.WithOauthAllowedScopes(collections.Map(scopes, func(scope sdk.ApiIntegrationOauthAllowedScope) sdk.ApiIntegrationOauthAllowedScopeItem {
			return sdk.ApiIntegrationOauthAllowedScopeItem{Scope: scope}
		}))
	}
	if err = errors.Join(
		intAttributeCreate(d, "oauth_access_token_validity", &authReq.OauthAccessTokenValidity),
		intAttributeCreate(d, "oauth_refresh_token_validity", &authReq.OauthRefreshTokenValidity),
		stringAttributeCreate(d, "oauth_username", &authReq.OauthUsername),
	); err != nil {
		return diag.FromErr(err)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithGitHttpsApiOAuth2ProviderParams(*sdk.NewGitHttpsApiOAuth2ParamsRequest(*authReq))); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Git Repository OAuth2 API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationGitRepositoryOauth2(ctx, d, meta)
}

func ReadApiIntegrationGitRepositoryOauth2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query API integration Git Repository OAuth2. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration Git Repository OAuth2 id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	gitDetails, err := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration Git Repository OAuth2 (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, gitDetails.AllowedPrefixes, gitDetails.BlockedPrefixes),
		d.Set("oauth_authorization_endpoint", gitDetails.OauthAuthorizationEndpoint),
		d.Set("oauth_token_endpoint", gitDetails.OauthTokenEndpoint),
		d.Set("oauth_client_id", gitDetails.OauthClientId),
		// oauth_client_secret intentionally omitted — Snowflake does not return it
		d.Set("oauth_access_token_validity", gitDetails.OauthAccessTokenValidity),
		d.Set("oauth_refresh_token_validity", gitDetails.OauthRefreshTokenValidity),
		d.Set("oauth_allowed_scopes", gitDetails.OauthAllowedScopes),
		d.Set("oauth_username", gitDetails.OauthUsername),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationGitRepositoryOauth2DetailsToSchema(gitDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationGitRepositoryOauth2(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Git Repository OAuth2 API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Git Repository OAuth2 API integration: %w", err))
		}
	}

	return ReadApiIntegrationGitRepositoryOauth2(ctx, d, meta)
}
