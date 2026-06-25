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

var apiIntegrationExternalMcpDynamicClientSchema = func() map[string]*schema.Schema {
	apiIntegrationExternalMcpDynamicClient := map[string]*schema.Schema{
		"oauth_resource_url": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The URL of the OAuth2 protected resource server. This URL is used by Snowflake to discover OAuth2 provider endpoints via RFC 8414 server metadata.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeExternalMcpDynamicClientApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationExternalMcpDynamicClient)
}()

func ApiIntegrationExternalMcpDynamicClient() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpDynamicClientResource), TrackingCreateWrapper(resources.ApiIntegrationExternalMcpDynamicClient, CreateApiIntegrationExternalMcpDynamicClient)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpDynamicClientResource), TrackingReadWrapper(resources.ApiIntegrationExternalMcpDynamicClient, ReadApiIntegrationExternalMcpDynamicClient)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpDynamicClientResource), TrackingUpdateWrapper(resources.ApiIntegrationExternalMcpDynamicClient, UpdateApiIntegrationExternalMcpDynamicClient)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ApiIntegrationExternalMcpDynamicClientResource), TrackingDeleteWrapper(resources.ApiIntegrationExternalMcpDynamicClient, deleteFunc)),
		Description:   "Resource used to manage API integration External MCP Dynamic Client Registration objects. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationExternalMcpDynamicClientSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationExternalMcpDynamicClient, ImportApiIntegrationExternalMcpDynamicClient),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationExternalMcpDynamicClientSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationExternalMcpDynamicClientSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "oauth_resource_url"),
		),
	}
}

func CreateApiIntegrationExternalMcpDynamicClient(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	auth := sdk.NewDynamicClientMcpUserAuthenticationRequest(d.Get("oauth_resource_url").(string))
	dynamicClientParams := sdk.NewExternalMcpDynamicClientParamsRequest(*auth)

	if err = client.ApiIntegrations.Create(ctx, request.WithExternalMcpDynamicClientProviderParams(*dynamicClientParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating External MCP Dynamic Client API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationExternalMcpDynamicClient(ctx, d, meta)
}

func ImportApiIntegrationExternalMcpDynamicClient(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	return importApiIntegrationWithDetails(ctx, d, meta,
		func(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) (*sdk.ApiIntegrationExternalMcpDetails, error) {
			return client.ApiIntegrations.DescribeExternalMcpDetails(ctx, id)
		},
		func(details *sdk.ApiIntegrationExternalMcpDetails, id sdk.AccountObjectIdentifier) error {
			if _, err := sdk.ToApiIntegrationMcpApiProviderType(details.ApiProvider); err != nil {
				return fmt.Errorf(
					"api integration %s has api_provider %s, not compatible with snowflake_api_integration_external_mcp_dynamic_client (expected external_mcp); use the appropriate resource type",
					id.FullyQualifiedName(),
					details.ApiProvider,
				)
			}
			if details.UserAuthType != string(sdk.ApiIntegrationUserAuthTypeOauthDynamicClient) {
				return fmt.Errorf(
					"api integration %s has user_auth_type %s, not compatible with snowflake_api_integration_external_mcp_dynamic_client (expected %s); use the appropriate resource type",
					id.FullyQualifiedName(),
					details.UserAuthType,
					sdk.ApiIntegrationUserAuthTypeOauthDynamicClient,
				)
			}
			return nil
		},
	)
}

func ReadApiIntegrationExternalMcpDynamicClient(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query API integration External MCP Dynamic Client. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration External MCP Dynamic Client id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.ApiIntegrations.DescribeExternalMcpDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration External MCP Dynamic Client (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, details.AllowedPrefixes, details.BlockedPrefixes),
		d.Set("oauth_resource_url", details.OauthResourceUrl),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationExternalMcpDynamicClientDetailsToSchema(details)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationExternalMcpDynamicClient(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()

	// oauth_resource_url is ForceNew — no provider-specific set/unset params needed
	if err = handleApiIntegrationCommonUpdate(d, set, unset); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating External MCP Dynamic Client API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating External MCP Dynamic Client API integration: %w", err))
		}
	}

	return ReadApiIntegrationExternalMcpDynamicClient(ctx, d, meta)
}
