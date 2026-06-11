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

var apiIntegrationAzureApiManagementSchema = func() map[string]*schema.Schema {
	apiIntegrationAzureApiManagement := map[string]*schema.Schema{
		"azure_tenant_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Specifies the ID for your Office 365 tenant that all Azure API Management instances belong to.",
		},
		"azure_ad_application_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The 'Application (client) ID' of the Azure AD app for your Azure API Management instance.",
		},
		"api_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: externalChangesNotDetectedFieldDescription("Specifies the API key (secret) that Snowflake uses to authenticate when making calls to the proxy service. Snowflake returns a masked value for this field in DESCRIBE output, so external changes to it cannot be detected."),
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeAzureApiManagementApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationAzureApiManagement)
}()

func ApiIntegrationAzureApiManagement() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ApiIntegrationAzureApiManagementResource), TrackingCreateWrapper(resources.ApiIntegrationAzureApiManagement, CreateApiIntegrationAzureApiManagement)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ApiIntegrationAzureApiManagementResource), TrackingReadWrapper(resources.ApiIntegrationAzureApiManagement, ReadApiIntegrationAzureApiManagement)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ApiIntegrationAzureApiManagementResource), TrackingUpdateWrapper(resources.ApiIntegrationAzureApiManagement, UpdateApiIntegrationAzureApiManagement)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ApiIntegrationAzureApiManagementResource), TrackingDeleteWrapper(resources.ApiIntegrationAzureApiManagement, deleteFunc)),
		Description:   "Resource used to manage API integration Azure API Management objects. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationAzureApiManagementSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationAzureApiManagement, ImportApiIntegrationAzureApiManagement),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationAzureApiManagementSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationAzureApiManagementSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "azure_tenant_id", "azure_ad_application_id", "api_key"),
		),
	}
}

func ImportApiIntegrationAzureApiManagement(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	details, err := client.ApiIntegrations.DescribeAzureDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe API integration %s during import: %w", id.FullyQualifiedName(), err)
	}

	normalizedProvider, err := sdk.ToApiIntegrationAzureApiProviderType(details.ApiProvider)
	if err != nil || normalizedProvider != sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement {
		return nil, fmt.Errorf("api integration %s has api_provider %s, not compatible with snowflake_api_integration_azure_api_management (expected %s); use the appropriate resource type",
			id.FullyQualifiedName(), details.ApiProvider, sdk.ApiIntegrationAzureApiProviderTypeAzureApiManagement)
	}

	if err = d.Set("name", id.Name()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateApiIntegrationAzureApiManagement(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	azureParams := sdk.NewAzureApiParamsRequest(
		d.Get("azure_tenant_id").(string),
		d.Get("azure_ad_application_id").(string),
	)

	if err = stringAttributeCreate(d, "api_key", &azureParams.ApiKey); err != nil {
		return diag.FromErr(err)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithAzureApiProviderParams(*azureParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Azure API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationAzureApiManagement(ctx, d, meta)
}

func ReadApiIntegrationAzureApiManagement(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query API integration Azure API Management. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration Azure API Management id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	azureDetails, err := client.ApiIntegrations.DescribeAzureDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration Azure API Management (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, azureDetails.AllowedPrefixes, azureDetails.BlockedPrefixes),
		d.Set("azure_tenant_id", azureDetails.AzureTenantId),
		d.Set("azure_ad_application_id", azureDetails.AzureAdApplicationId),
		// api_key intentionally omitted — Snowflake returns a masked value (☺☺☺☺☺☺☺☺☺☺) and external changes cannot be detected
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationAzureApiManagementDetailsToSchema(azureDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationAzureApiManagement(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()
	azureSet := sdk.NewSetAzureApiParamsRequest()
	azureUnset := sdk.NewUnsetAzureApiParamsRequest()

	errs := errors.Join(
		handleApiIntegrationCommonUpdate(d, set, unset),
		stringAttributeUpdateSetOnlyNotEmpty(d, "azure_tenant_id", &azureSet.AzureTenantId),
		stringAttributeUpdateSetOnlyNotEmpty(d, "azure_ad_application_id", &azureSet.AzureAdApplicationId),
		stringAttributeUpdate(d, "api_key", &azureSet.ApiKey, &azureUnset.ApiKey),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*azureSet, *sdk.NewSetAzureApiParamsRequest()) {
		set.WithAzureParams(*azureSet)
	}
	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Azure API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*azureUnset, *sdk.NewUnsetAzureApiParamsRequest()) {
		unset.WithAzureParams(*azureUnset)
	}
	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Azure API integration: %w", err))
		}
	}

	return ReadApiIntegrationAzureApiManagement(ctx, d, meta)
}
