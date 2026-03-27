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

var catalogIntegrationOpenCatalogSchema = func() map[string]*schema.Schema {
	openCatalogSchema := map[string]*schema.Schema{
		"catalog_namespace": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Specifies the default Open Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"rest_config": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Specifies information about the Open Catalog account and catalog name.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"catalog_uri": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "Specifies Open Catalog account URL.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"catalog_name": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "Specifies the name of the catalog to use in Open Catalog.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"catalog_api_type": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies how Snowflake connects to Open Catalog. " + enumValuesDescription(sdk.AllCatalogIntegrationCatalogApiTypes),
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_config.0.catalog_api_type"),
						ValidateDiagFunc: sdkValidation(sdk.ToCatalogIntegrationCatalogApiType),
					},
					"access_delegation_mode": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies the access delegation mode for accessing Iceberg table files in your external cloud storage. " + enumValuesDescription(sdk.AllCatalogIntegrationAccessDelegationModes),
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_config.0.access_delegation_mode"),
						ValidateDiagFunc: sdkValidation(sdk.ToCatalogIntegrationAccessDelegationMode),
					},
				},
			},
		},
		"rest_authentication": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Specifies authentication details that Snowflake uses to connect to Open Catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"oauth_token_uri": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies URL for the third-party identity provider. If not specified, Snowflake assumes the remote catalog provider is the identity provider.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("rest_authentication.0.oauth_token_uri"),
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"oauth_client_id": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						Description:  "Specifies the client ID of the OAuth2 credential associated with your Open Catalog service connection.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"oauth_client_secret": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						Description:  externalChangesNotDetectedFieldDescription("Specifies the secret of the OAuth2 credential associated with your Open Catalog service connection."),
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"oauth_allowed_scopes": {
						Type:        schema.TypeList,
						Required:    true,
						ForceNew:    true,
						MinItems:    1,
						Description: "Specifies one or more scopes for the OAuth token.",
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}
	return collections.MergeMaps(catalogIntegrationCommonSchema(schemas.DescribeCatalogIntegrationOpenCatalogDetailsSchema), openCatalogSchema)
}()

func CatalogIntegrationOpenCatalog() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationOpenCatalogResource), TrackingCreateWrapper(resources.CatalogIntegrationOpenCatalog, CreateCatalogIntegrationOpenCatalog)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationOpenCatalogResource), TrackingReadWrapper(resources.CatalogIntegrationOpenCatalog, ReadCatalogIntegrationOpenCatalogFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationOpenCatalogResource), TrackingUpdateWrapper(resources.CatalogIntegrationOpenCatalog, UpdateCatalogIntegrationOpenCatalog)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationOpenCatalogResource), TrackingDeleteWrapper(resources.CatalogIntegrationOpenCatalog, deleteCatalogIntegrationFunc())),
		Description:   "Resource used to manage catalog integration objects for Apache Iceberg™ tables that integrate with Snowflake Open Catalog. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-open-catalog).",

		Schema: catalogIntegrationOpenCatalogSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogIntegrationOpenCatalog, ImportCatalogIntegrationOpenCatalog),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(catalogIntegrationOpenCatalogSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(catalogIntegrationOpenCatalogSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		),
	}
}

func ImportCatalogIntegrationOpenCatalog(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypePolaris {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypePolaris, details.CatalogSource)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateCatalogIntegrationOpenCatalog(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}
	enabled := d.Get("enabled").(bool)

	restConfig, err := buildOpenCatalogRestConfigRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	restAuth, err := buildOAuthRestAuthenticationRequest(d, "rest_authentication")
	if err != nil {
		return diag.FromErr(err)
	}

	openCatalogParams := sdk.NewOpenCatalogParamsRequest().
		WithRestConfig(restConfig).
		WithRestAuthentication(restAuth)
	if err := stringAttributeCreateBuilder(d, "catalog_namespace", openCatalogParams.WithCatalogNamespace); err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateCatalogIntegrationRequest(id, enabled).WithOpenCatalogCatalogSourceParams(*openCatalogParams)
	errs := errors.Join(
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.CatalogIntegrations.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Open Catalog catalog integration, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationOpenCatalogFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationOpenCatalogFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		s, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query Open Catalog catalog integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Open Catalog catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe Open Catalog catalog integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			if err = handleExternalChangesToNestedAttrs(d, details); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("name", details.Id.Name()),
			d.Set("enabled", details.Enabled),
			// not reading refresh_interval_seconds on purpose (handled as external change to describe output)
			d.Set("comment", details.Comment),
			d.Set("catalog_namespace", details.CatalogNamespace),
			// not reading rest_config and rest_authentication on purpose
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationOpenCatalogDetailsToSchema(details)}),
		)

		return diag.FromErr(errs)
	}
}

func UpdateCatalogIntegrationOpenCatalog(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewCatalogIntegrationSetRequest()

	if d.HasChange("rest_authentication.0.oauth_client_secret") {
		newSecret := d.Get("rest_authentication.0.oauth_client_secret").(string)
		set.SetOAuthRestAuthentication = sdk.NewSetOAuthRestAuthenticationRequest(newSecret)
	}

	if errs := handleCatalogIntegrationCommonPropsUpdate(d, set); errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*set, *sdk.NewCatalogIntegrationSetRequest()) {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Open Catalog catalog integration (%s), err = %w", d.Id(), err))
		}
	}
	return ReadCatalogIntegrationOpenCatalogFunc(false)(ctx, d, meta)
}

func buildOpenCatalogRestConfigRequest(d *schema.ResourceData) (sdk.OpenCatalogRestConfigRequest, error) {
	v := d.Get("rest_config").([]any)
	block := v[0].(map[string]any)
	catalogUri := block["catalog_uri"].(string)
	catalogName := block["catalog_name"].(string)
	req := sdk.NewOpenCatalogRestConfigRequest(catalogUri, catalogName)

	if v, ok := block["catalog_api_type"].(string); ok && v != "" {
		t, err := sdk.ToCatalogIntegrationCatalogApiType(v)
		if err != nil {
			return sdk.OpenCatalogRestConfigRequest{}, err
		}
		req = req.WithCatalogApiType(t)
	}
	if v, ok := block["access_delegation_mode"].(string); ok && v != "" {
		t, err := sdk.ToCatalogIntegrationAccessDelegationMode(v)
		if err != nil {
			return sdk.OpenCatalogRestConfigRequest{}, err
		}
		req = req.WithAccessDelegationMode(t)
	}
	return *req, nil
}

// handleExternalChangesToNestedAttrs marks drift on rest_config and rest_authentication when
// the prior describe_output differs from the current Snowflake response.
func handleExternalChangesToNestedAttrs(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails) error {
	restConfig := []any{
		map[string]any{
			"catalog_uri":            details.RestConfig.CatalogUri,
			"catalog_name":           details.RestConfig.CatalogName,
			"catalog_api_type":       string(details.RestConfig.CatalogApiType),
			"access_delegation_mode": string(details.RestConfig.AccessDelegationMode),
		},
	}
	err := handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
		outputMapping{"rest_config", "rest_config", restConfig, restConfig, nil},
	)
	if err != nil {
		return err
	}
	return handleExternalChangesToRestAuthentication(d, details)
}

func handleExternalChangesToRestAuthentication(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails) error {
	return handleExternalChangesToOAuthRestAuthentication(
		d,
		"rest_authentication",
		details.RestAuthentication.OauthTokenUri,
		details.RestAuthentication.OauthClientId,
		details.RestAuthentication.OauthAllowedScopes,
	)
}
