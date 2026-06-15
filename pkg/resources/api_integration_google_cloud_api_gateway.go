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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiIntegrationGoogleCloudApiGatewaySchema = func() map[string]*schema.Schema {
	apiIntegrationGoogleCloudApiGateway := map[string]*schema.Schema{
		"google_audience": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the audience claim used by Snowflake when generating the JWT to authenticate with the Google Cloud API Gateway.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeGoogleCloudApiGatewayApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationGoogleCloudApiGateway)
}()

func ApiIntegrationGoogleCloudApiGateway() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ApiIntegrationGoogleCloudApiGateway, CreateApiIntegrationGoogleCloudApiGateway),
		ReadContext:   TrackingReadWrapper(resources.ApiIntegrationGoogleCloudApiGateway, ReadApiIntegrationGoogleCloudApiGateway),
		UpdateContext: TrackingUpdateWrapper(resources.ApiIntegrationGoogleCloudApiGateway, UpdateApiIntegrationGoogleCloudApiGateway),
		DeleteContext: TrackingDeleteWrapper(resources.ApiIntegrationGoogleCloudApiGateway, deleteFunc),
		Description:   "Resource used to manage API integration Google Cloud API Gateway objects. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationGoogleCloudApiGatewaySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationGoogleCloudApiGateway, ImportApiIntegrationGoogleCloudApiGateway),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationGoogleCloudApiGatewaySchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationGoogleCloudApiGatewaySchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "google_audience"),
		),
	}
}

func CreateApiIntegrationGoogleCloudApiGateway(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	googleParams := sdk.NewGoogleApiParamsRequest(d.Get("google_audience").(string))

	if err = client.ApiIntegrations.Create(ctx, request.WithGoogleApiProviderParams(*googleParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Google Cloud API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationGoogleCloudApiGateway(ctx, d, meta)
}

func ImportApiIntegrationGoogleCloudApiGateway(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.ApiIntegrations.DescribeGoogleDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe API integration %s during import: %w", id.FullyQualifiedName(), err)
	}

	if _, err := sdk.ToApiIntegrationGoogleApiProviderType(details.ApiProvider); err != nil {
		return nil, fmt.Errorf(
			"api integration %s has api_provider %s, not compatible with snowflake_api_integration_google_cloud_api_gateway (expected one of %s); use the appropriate resource type",
			id.FullyQualifiedName(),
			details.ApiProvider,
			possibleValuesListed(sdk.AsStringList(sdk.AllApiIntegrationGoogleApiProviderTypes)),
		)
	}

	return ImportName[sdk.AccountObjectIdentifier](ctx, d, meta)
}

func ReadApiIntegrationGoogleCloudApiGateway(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query API integration Google Cloud API Gateway. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration Google Cloud API Gateway id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	googleDetails, err := client.ApiIntegrations.DescribeGoogleDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration Google Cloud API Gateway (%s): %w", d.Id(), err))
	}

	if _, err := sdk.ToApiIntegrationGoogleApiProviderType(googleDetails.ApiProvider); err != nil {
		return diag.FromErr(fmt.Errorf("could not normalize api_provider value (%s): %w", googleDetails.ApiProvider, err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, googleDetails.AllowedPrefixes, googleDetails.BlockedPrefixes),
		d.Set("google_audience", googleDetails.GoogleAudience),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationGoogleCloudApiGatewayDetailsToSchema(googleDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationGoogleCloudApiGateway(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	if d.HasChange("google_audience") {
		googleSet := sdk.NewSetGoogleApiParamsRequest(d.Get("google_audience").(string))
		set.WithGoogleParams(*googleSet)
	}

	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Google Cloud API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Google Cloud API integration: %w", err))
		}
	}

	return ReadApiIntegrationGoogleCloudApiGateway(ctx, d, meta)
}
