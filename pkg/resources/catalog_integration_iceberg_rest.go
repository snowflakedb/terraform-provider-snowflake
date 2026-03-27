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

var catalogIntegrationIcebergRestSchema = func() map[string]*schema.Schema {
	authExactlyOneOf := []string{"oauth_rest_authentication", "bearer_rest_authentication", "sigv4_rest_authentication"}

	icebergRestSchema := map[string]*schema.Schema{
		"catalog_namespace": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Specifies the default namespace for all Iceberg tables that you associate with the catalog integration.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"rest_config": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "Specifies information about REST configuration.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"catalog_uri": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "Specifies the endpoint URL for the catalog REST API.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"prefix": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Specifies an optional prefix appended to all API routes.",
					},
					"catalog_name": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						Description:  "Specifies the catalog or identifier to request from your remote catalog service.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"catalog_api_type": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies the connection type for the catalog API. " + enumValuesDescription(sdk.AllCatalogIntegrationCatalogApiTypes),
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
		"oauth_rest_authentication": {
			Type:         schema.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: authExactlyOneOf,
			Description:  "Specifies OAuth as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"oauth_token_uri": {
						Type:             schema.TypeString,
						Optional:         true,
						ForceNew:         true,
						Description:      "Specifies URL for the third-party identity provider. If not specified, Snowflake assumes the remote catalog provider is the identity provider.",
						DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("oauth_rest_authentication.0.oauth_token_uri"),
						ValidateFunc:     validation.StringIsNotEmpty,
					},
					"oauth_client_id": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						Description:  "Specifies the client ID of the OAuth2 credential.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"oauth_client_secret": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						Description:  externalChangesNotDetectedFieldDescription("Specifies the secret of the OAuth2 credential."),
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
		"bearer_rest_authentication": {
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: authExactlyOneOf,
			Description:  "Specifies a bearer token as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bearer_token": {
						Type:         schema.TypeString,
						Required:     true,
						Sensitive:    true,
						Description:  externalChangesNotDetectedFieldDescription("The bearer token for the identity provider."),
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"sigv4_rest_authentication": {
			Type:         schema.TypeList,
			Optional:     true,
			ForceNew:     true,
			MaxItems:     1,
			ExactlyOneOf: authExactlyOneOf,
			Description:  "Specifies Signature Version 4 as the authentication type for Snowflake to use to connect to the Iceberg REST catalog.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"sigv4_iam_role": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "Specifies the Amazon Resource Name (ARN) for an IAM role that has permission to access your REST API in API Gateway.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sigv4_signing_region": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						Description:  "Specifies the AWS Region associated with your API in API Gateway. If you don’t specify this parameter, Snowflake uses the region in which your Snowflake account is deployed.",
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sigv4_external_id": {
						Type:     schema.TypeString,
						Optional: true,
						ForceNew: true,
						Description: externalChangesNotDetectedFieldDescription(joinWithSpace("Specifies an external ID that Snowflake uses to establish a trust relationship with AWS.",
							"If you don’t specify this parameter, Snowflake automatically generates a unique external ID when you create a catalog integration.")),
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
	return collections.MergeMaps(catalogIntegrationCommonSchema(schemas.DescribeCatalogIntegrationIcebergRestDetailsSchema), icebergRestSchema)
}()

func CatalogIntegrationIcebergRest() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationIcebergRestResource), TrackingCreateWrapper(resources.CatalogIntegrationIcebergRest, CreateCatalogIntegrationIcebergRest)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationIcebergRestResource), TrackingReadWrapper(resources.CatalogIntegrationIcebergRest, ReadCatalogIntegrationIcebergRestFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationIcebergRestResource), TrackingUpdateWrapper(resources.CatalogIntegrationIcebergRest, UpdateCatalogIntegrationIcebergRest)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationIcebergRestResource), TrackingDeleteWrapper(resources.CatalogIntegrationIcebergRest, deleteCatalogIntegrationFunc())),
		Description:   "Resource used to manage catalog integration objects for Apache Iceberg™ tables managed in a remote catalog that complies with the open source Apache Iceberg™ REST OpenAPI specification. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-rest).",

		Schema: catalogIntegrationIcebergRestSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogIntegrationIcebergRest, ImportCatalogIntegrationIcebergRest),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(catalogIntegrationIcebergRestSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(catalogIntegrationIcebergRestSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		),
	}
}

func ImportCatalogIntegrationIcebergRest(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypeIcebergREST {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypeIcebergREST, details.CatalogSource)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateCatalogIntegrationIcebergRest(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}
	enabled := d.Get("enabled").(bool)

	restConfig, err := buildIcebergRestRestConfigRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	params := sdk.NewIcebergRestParamsRequest().WithRestConfig(restConfig)
	if err := stringAttributeCreateBuilder(d, "catalog_namespace", params.WithCatalogNamespace); err != nil {
		return diag.FromErr(err)
	}

	if err := applyIcebergRestAuthenticationToParams(d, params); err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateCatalogIntegrationRequest(id, enabled).WithIcebergRestCatalogSourceParams(*params)
	errs := errors.Join(
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.CatalogIntegrations.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg REST catalog integration, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationIcebergRestFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationIcebergRestFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query Iceberg REST catalog integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Iceberg REST catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe Iceberg REST catalog integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			if err = handleExternalChangesToIcebergRestNestedAttrs(d, details); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("name", details.Id.Name()),
			d.Set("enabled", details.Enabled),
			d.Set("comment", details.Comment),
			d.Set("catalog_namespace", details.CatalogNamespace),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationIcebergRestDetailsToSchema(details)}),
		)

		return diag.FromErr(errs)
	}
}

func UpdateCatalogIntegrationIcebergRest(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewCatalogIntegrationSetRequest()

	if d.HasChange("oauth_rest_authentication.0.oauth_client_secret") {
		newSecret := d.Get("oauth_rest_authentication.0.oauth_client_secret").(string)
		set.SetOAuthRestAuthentication = sdk.NewSetOAuthRestAuthenticationRequest(newSecret)
	}
	if d.HasChange("bearer_rest_authentication.0.bearer_token") {
		newToken := d.Get("bearer_rest_authentication.0.bearer_token").(string)
		set.SetBearerRestAuthentication = sdk.NewSetBearerRestAuthenticationRequest(newToken)
	}

	if errs := handleCatalogIntegrationCommonPropsUpdate(d, set); errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*set, *sdk.NewCatalogIntegrationSetRequest()) {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Iceberg REST catalog integration (%s), err = %w", d.Id(), err))
		}
	}
	return ReadCatalogIntegrationIcebergRestFunc(false)(ctx, d, meta)
}

func buildIcebergRestRestConfigRequest(d *schema.ResourceData) (sdk.IcebergRestRestConfigRequest, error) {
	v := d.Get("rest_config").([]any)
	block := v[0].(map[string]any)
	catalogUri := block["catalog_uri"].(string)
	req := sdk.IcebergRestRestConfigRequest{CatalogUri: catalogUri}
	if p, ok := block["prefix"].(string); ok && p != "" {
		req.Prefix = sdk.String(p)
	}
	if n, ok := block["catalog_name"].(string); ok && n != "" {
		req.CatalogName = sdk.String(n)
	}
	if v, ok := block["catalog_api_type"].(string); ok && v != "" {
		t, err := sdk.ToCatalogIntegrationCatalogApiType(v)
		if err != nil {
			return sdk.IcebergRestRestConfigRequest{}, err
		}
		req.CatalogApiType = &t
	}
	if v, ok := block["access_delegation_mode"].(string); ok && v != "" {
		t, err := sdk.ToCatalogIntegrationAccessDelegationMode(v)
		if err != nil {
			return sdk.IcebergRestRestConfigRequest{}, err
		}
		req.AccessDelegationMode = &t
	}
	return req, nil
}

func applyIcebergRestAuthenticationToParams(d *schema.ResourceData, params *sdk.IcebergRestParamsRequest) error {
	if v, ok := d.GetOk("oauth_rest_authentication"); ok && len(v.([]any)) > 0 {
		oauth, err := buildOAuthRestAuthenticationRequest(d, "oauth_rest_authentication")
		if err != nil {
			return err
		}
		params.WithOAuthRestAuthentication(oauth)
		return nil
	}
	if v, ok := d.GetOk("bearer_rest_authentication"); ok && len(v.([]any)) > 0 {
		block := v.([]any)[0].(map[string]any)
		token := block["bearer_token"].(string)
		params.WithBearerRestAuthentication(*sdk.NewBearerRestAuthenticationRequest(token))
		return nil
	}
	if v, ok := d.GetOk("sigv4_rest_authentication"); ok && len(v.([]any)) > 0 {
		block := v.([]any)[0].(map[string]any)
		role := block["sigv4_iam_role"].(string)
		req := sdk.NewSigV4RestAuthenticationRequest(role)
		if r, ok := block["sigv4_signing_region"].(string); ok && r != "" {
			req = req.WithSigv4SigningRegion(r)
		}
		if e, ok := block["sigv4_external_id"].(string); ok && e != "" {
			req = req.WithSigv4ExternalId(e)
		}
		params.WithSigV4RestAuthentication(*req)
		return nil
	}
	return fmt.Errorf("one of oauth_rest_authentication, bearer_rest_authentication, or sigv4_rest_authentication must be set")
}

func handleExternalChangesToIcebergRestNestedAttrs(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	if err := handleExternalChangesToIcebergRestConfig(d, details); err != nil {
		return err
	}
	if details.OAuthRestAuthentication != nil {
		return handleExternalChangesToOAuthRestAuthenticationIcebergRest(d, details)
	}
	if details.SigV4RestAuthentication != nil {
		return handleExternalChangesToSigV4RestAuthentication(d, details)
	}
	return handleExternalChangesToBearerRestAuthentication(d)
}

func handleExternalChangesToIcebergRestConfig(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	restConfig := []any{
		map[string]any{
			"catalog_uri":            details.RestConfig.CatalogUri,
			"prefix":                 details.RestConfig.Prefix,
			"catalog_name":           details.RestConfig.CatalogName,
			"catalog_api_type":       string(details.RestConfig.CatalogApiType),
			"access_delegation_mode": string(details.RestConfig.AccessDelegationMode),
		},
	}
	return handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
		outputMapping{"rest_config", "rest_config", restConfig, restConfig, nil},
	)
}

func handleExternalChangesToOAuthRestAuthenticationIcebergRest(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	return errors.Join(
		handleExternalChangesToOAuthRestAuthentication(
			d,
			"oauth_rest_authentication",
			details.OAuthRestAuthentication.OauthTokenUri,
			details.OAuthRestAuthentication.OauthClientId,
			details.OAuthRestAuthentication.OauthAllowedScopes,
		),
		d.Set("sigv4_rest_authentication", nil),
	)
}

func handleExternalChangesToSigV4RestAuthentication(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	sigV4RestAuthorization := []any{
		map[string]any{
			"sigv4_iam_role":       details.SigV4RestAuthentication.Sigv4IamRole,
			"sigv4_signing_region": details.SigV4RestAuthentication.Sigv4SigningRegion,
		},
	}
	err := handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
		outputMapping{"sigv4_rest_authentication", "sigv4_rest_authentication", sigV4RestAuthorization, sigV4RestAuthorization, nil},
	)
	return errors.Join(err, d.Set("oauth_rest_authentication", nil))
}

func handleExternalChangesToBearerRestAuthentication(d *schema.ResourceData) error {
	return errors.Join(d.Set("oauth_rest_authentication", nil), d.Set("sigv4_rest_authentication", nil))
}
