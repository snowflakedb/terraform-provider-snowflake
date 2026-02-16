package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
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

var catalogIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the integration; must be unique in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"catalog_source": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the catalog source for the integration. Valid values are: OBJECT_STORE, GLUE, ICEBERG_REST, POLARIS, SAP_BDC.",
		ValidateDiagFunc: StringInSlice([]string{"OBJECT_STORE", "GLUE", "ICEBERG_REST", "POLARIS", "SAP_BDC"}, true),
	},
	"table_format": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the table format. Valid values are: ICEBERG, DELTA.",
		ValidateDiagFunc: StringInSlice([]string{"ICEBERG", "DELTA"}, true),
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether this catalog integration is enabled.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the catalog integration.",
	},
	"catalog_namespace": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the catalog namespace. Used with certain catalog sources.",
	},
	// GLUE-specific parameters
	"glue_aws_role_arn": {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringIsNotEmpty,
		Description:  "Specifies the Amazon Resource Name (ARN) of the AWS identity and access management (IAM) role for AWS Glue catalog. Required when catalog_source is GLUE.",
	},
	"glue_catalog_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the AWS Glue Data Catalog ID. Required when catalog_source is GLUE.",
	},
	"glue_region": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the AWS region where the Glue catalog is located.",
	},
	// REST config parameters (for ICEBERG_REST, POLARIS, SAP_BDC)
	"rest_config_catalog_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the catalog URI for REST-based catalog sources. Required when catalog_source is ICEBERG_REST, POLARIS, or SAP_BDC.",
	},
	"rest_config_catalog_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the catalog name for REST-based catalog sources.",
	},
	"rest_config_catalog_api_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the catalog API type for REST-based catalog sources.",
	},
	"rest_config_prefix": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the prefix for REST-based catalog sources.",
	},
	"rest_config_access_delegation_mode": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the access delegation mode for REST-based catalog sources.",
	},
	// REST authentication parameters (for ICEBERG_REST, POLARIS)
	"rest_auth_type": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the authentication type for REST-based catalog sources. Required when catalog_source is ICEBERG_REST or POLARIS.",
	},
	"rest_auth_oauth_client_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the OAuth client ID for REST authentication.",
	},
	"rest_auth_oauth_client_secret": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Specifies the OAuth client secret for REST authentication.",
	},
	"rest_auth_oauth_token_uri": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the OAuth token URI for REST authentication.",
	},
	"rest_auth_oauth_allowed_scopes": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specifies the OAuth allowed scopes for REST authentication.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"rest_auth_bearer_token": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "Specifies the bearer token for REST authentication.",
	},
	"rest_auth_sigv4_iam_role": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the SIGv4 IAM role for REST authentication.",
	},
	"rest_auth_sigv4_signing_region": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the SIGv4 signing region for REST authentication.",
	},
	"rest_auth_sigv4_external_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the SIGv4 external ID for REST authentication.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CATALOG INTEGRATIONS` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowCatalogIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE CATALOG INTEGRATION` for the given catalog integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeCatalogIntegrationSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func CatalogIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.CatalogIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationResource), TrackingCreateWrapper(resources.CatalogIntegration, CreateCatalogIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationResource), TrackingReadWrapper(resources.CatalogIntegration, ReadCatalogIntegration)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationResource), TrackingUpdateWrapper(resources.CatalogIntegration, UpdateCatalogIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationResource), TrackingDeleteWrapper(resources.CatalogIntegration, deleteFunc)),
		Description:   "Resource used to manage catalog integration objects. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration).",

		Schema: catalogIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogIntegration, ImportCatalogIntegration),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(catalogIntegrationSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(catalogIntegrationSchema, DescribeOutputAttributeName,
				"enabled", "catalog_namespace", "comment",
				"glue_aws_role_arn", "glue_catalog_id", "glue_region",
				"rest_config_catalog_uri", "rest_config_catalog_name", "rest_config_catalog_api_type", "rest_config_prefix", "rest_config_access_delegation_mode",
				"rest_auth_type", "rest_auth_oauth_client_id", "rest_auth_oauth_client_secret", "rest_auth_oauth_token_uri", "rest_auth_oauth_allowed_scopes",
				"rest_auth_bearer_token", "rest_auth_sigv4_iam_role", "rest_auth_sigv4_signing_region", "rest_auth_sigv4_external_id",
			),
		),
	}
}

func ImportCatalogIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	catalogIntegration, err := client.CatalogIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	properties, err := client.CatalogIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to describe catalog integration %s, err: %w", id.FullyQualifiedName(), err)
	}

	propsMap := make(map[string]string)
	for _, prop := range properties {
		propsMap[prop.Name] = prop.Value
	}

	errs := errors.Join(
		d.Set("name", catalogIntegration.Name),
		d.Set("enabled", catalogIntegration.Enabled),
		d.Set("comment", catalogIntegration.Comment),
		d.Set("catalog_source", propsMap["CATALOG_SOURCE"]),
		d.Set("table_format", propsMap["TABLE_FORMAT"]),
		d.Set("catalog_namespace", propsMap["CATALOG_NAMESPACE"]),
		d.Set("glue_aws_role_arn", propsMap["GLUE_AWS_ROLE_ARN"]),
		d.Set("glue_catalog_id", propsMap["GLUE_CATALOG_ID"]),
		d.Set("glue_region", propsMap["GLUE_REGION"]),
		d.Set("rest_config_catalog_uri", propsMap["REST_CONFIG_CATALOG_URI"]),
		d.Set("rest_config_catalog_name", propsMap["REST_CONFIG_CATALOG_NAME"]),
		d.Set("rest_auth_type", propsMap["REST_AUTH_TYPE"]),
		d.Set("rest_auth_oauth_client_id", propsMap["REST_AUTH_OAUTH_CLIENT_ID"]),
		d.Set("rest_auth_oauth_token_uri", propsMap["REST_AUTH_OAUTH_TOKEN_URI"]),
	)
	if errs != nil {
		return nil, errs
	}

	return []*schema.ResourceData{d}, nil
}

func ReadCatalogIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	catalogIntegration, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query catalog integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	properties, err := client.CatalogIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to describe catalog integration %s, err: %w", id.FullyQualifiedName(), err))
	}

	propsMap := make(map[string]string)
	for _, prop := range properties {
		propsMap[prop.Name] = prop.Value
	}

	errs := errors.Join(
		d.Set("enabled", catalogIntegration.Enabled),
		d.Set("comment", catalogIntegration.Comment),
		d.Set("catalog_source", propsMap["CATALOG_SOURCE"]),
		d.Set("table_format", propsMap["TABLE_FORMAT"]),
		d.Set("catalog_namespace", propsMap["CATALOG_NAMESPACE"]),
		d.Set("glue_aws_role_arn", propsMap["GLUE_AWS_ROLE_ARN"]),
		d.Set("glue_catalog_id", propsMap["GLUE_CATALOG_ID"]),
		d.Set("glue_region", propsMap["GLUE_REGION"]),
		d.Set("rest_config_catalog_uri", propsMap["REST_CONFIG_CATALOG_URI"]),
		d.Set("rest_config_catalog_name", propsMap["REST_CONFIG_CATALOG_NAME"]),
		d.Set("rest_config_catalog_api_type", propsMap["REST_CONFIG_CATALOG_API_TYPE"]),
		d.Set("rest_config_prefix", propsMap["REST_CONFIG_PREFIX"]),
		d.Set("rest_config_access_delegation_mode", propsMap["REST_CONFIG_ACCESS_DELEGATION_MODE"]),
		d.Set("rest_auth_type", propsMap["REST_AUTH_TYPE"]),
		d.Set("rest_auth_oauth_client_id", propsMap["REST_AUTH_OAUTH_CLIENT_ID"]),
		d.Set("rest_auth_oauth_token_uri", propsMap["REST_AUTH_OAUTH_TOKEN_URI"]),
		d.Set("rest_auth_sigv4_iam_role", propsMap["REST_AUTH_SIGV4_IAM_ROLE"]),
		d.Set("rest_auth_sigv4_signing_region", propsMap["REST_AUTH_SIGV4_SIGNING_REGION"]),
		d.Set("rest_auth_sigv4_external_id", propsMap["REST_AUTH_SIGV4_EXTERNAL_ID"]),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
	)

	errs = errors.Join(errs,
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationToSchema(catalogIntegration)}),
		d.Set(DescribeOutputAttributeName, schemas.CatalogIntegrationPropertiesToSchema(properties)),
	)

	return diag.FromErr(errs)
}

func CreateCatalogIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(name)
	enabled := d.Get("enabled").(bool)
	catalogSource := d.Get("catalog_source").(string)
	tableFormat := sdk.TableFormat(d.Get("table_format").(string))

	request := sdk.NewCreateCatalogIntegrationRequest(id, enabled)

	switch catalogSource {
	case "OBJECT_STORE":
		request.WithObjectStoreParams(*sdk.NewObjectStoreParamsRequest(tableFormat))
	case "GLUE":
		glueAwsRoleArn := d.Get("glue_aws_role_arn").(string)
		glueCatalogId := d.Get("glue_catalog_id").(string)
		glueParams := sdk.NewGlueParamsRequest(tableFormat, glueAwsRoleArn, glueCatalogId)
		if v, ok := d.GetOk("glue_region"); ok {
			glueParams.WithGlueRegion(v.(string))
		}
		request.WithGlueParams(*glueParams)
	case "ICEBERG_REST":
		icebergRestParams := sdk.NewIcebergRestParamsRequest(tableFormat)
		if rc := buildRestConfig(d); rc != nil {
			icebergRestParams.WithRestConfig(*rc)
		}
		if ra := buildRestAuthentication(d); ra != nil {
			icebergRestParams.WithRestAuthentication(*ra)
		}
		request.WithIcebergRestParams(*icebergRestParams)
	case "POLARIS":
		polarisParams := sdk.NewPolarisParamsRequest(tableFormat)
		if rc := buildRestConfig(d); rc != nil {
			polarisParams.WithRestConfig(*rc)
		}
		if ra := buildRestAuthentication(d); ra != nil {
			polarisParams.WithRestAuthentication(*ra)
		}
		request.WithPolarisParams(*polarisParams)
	case "SAP_BDC":
		sapBdcParams := sdk.NewSapBdcParamsRequest(tableFormat)
		if rc := buildRestConfig(d); rc != nil {
			sapBdcParams.WithRestConfig(*rc)
		}
		request.WithSapBdcParams(*sapBdcParams)
	}

	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		stringAttributeCreateBuilder(d, "catalog_namespace", request.WithCatalogNamespace),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.CatalogIntegrations.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating catalog integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegration(ctx, d, meta)
}

func UpdateCatalogIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewCatalogIntegrationSetRequest()
	unset := sdk.NewCatalogIntegrationUnsetRequest()

	errs := errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		stringAttributeUpdate(d, "catalog_namespace", &set.CatalogNamespace, &unset.CatalogNamespace),
		stringAttributeUpdateSetOnlyNotEmpty(d, "glue_aws_role_arn", &set.GlueAwsRoleArn),
		stringAttributeUpdateSetOnlyNotEmpty(d, "glue_catalog_id", &set.GlueCatalogId),
		stringAttributeUpdateSetOnlyNotEmpty(d, "glue_region", &set.GlueRegion),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if updateRestConfig(d, set) {
		// REST config fields changed; already set on `set`
	}

	if set.Enabled != nil || set.Comment != nil || set.CatalogNamespace != nil ||
		set.GlueAwsRoleArn != nil || set.GlueCatalogId != nil || set.GlueRegion != nil ||
		set.RestConfig != nil || set.RestAuthentication != nil {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating catalog integration, err = %w", err))
		}
	}

	if unset.Comment != nil || unset.CatalogNamespace != nil {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithUnset(*unset)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating catalog integration, err = %w", err))
		}
	}

	return ReadCatalogIntegration(ctx, d, meta)
}

func buildRestConfig(d *schema.ResourceData) *sdk.RestConfigRequest {
	catalogUri, ok := d.GetOk("rest_config_catalog_uri")
	if !ok {
		return nil
	}
	rc := sdk.NewRestConfigRequest(catalogUri.(string))
	if v, ok := d.GetOk("rest_config_catalog_name"); ok {
		rc.WithCatalogName(v.(string))
	}
	if v, ok := d.GetOk("rest_config_catalog_api_type"); ok {
		rc.WithCatalogApiType(v.(string))
	}
	if v, ok := d.GetOk("rest_config_prefix"); ok {
		rc.WithPrefix(v.(string))
	}
	if v, ok := d.GetOk("rest_config_access_delegation_mode"); ok {
		rc.WithAccessDelegationMode(v.(string))
	}
	return rc
}

func buildRestAuthentication(d *schema.ResourceData) *sdk.RestAuthenticationRequest {
	authType, ok := d.GetOk("rest_auth_type")
	if !ok {
		return nil
	}
	ra := sdk.NewRestAuthenticationRequest(authType.(string))
	if v, ok := d.GetOk("rest_auth_oauth_client_id"); ok {
		ra.WithOauthClientId(v.(string))
	}
	if v, ok := d.GetOk("rest_auth_oauth_client_secret"); ok {
		ra.WithOauthClientSecret(v.(string))
	}
	if v, ok := d.GetOk("rest_auth_oauth_token_uri"); ok {
		ra.WithOauthTokenUri(v.(string))
	}
	if v, ok := d.GetOk("rest_auth_oauth_allowed_scopes"); ok {
		raw := v.([]any)
		scopes := make([]string, len(raw))
		for i, s := range raw {
			scopes[i] = s.(string)
		}
		ra.WithOauthAllowedScopes(scopes)
	}
	if v, ok := d.GetOk("rest_auth_bearer_token"); ok {
		ra.WithBearerToken(v.(string))
	}
	if v, ok := d.GetOk("rest_auth_sigv4_iam_role"); ok {
		ra.WithSigv4IamRole(v.(string))
	}
	if v, ok := d.GetOk("rest_auth_sigv4_signing_region"); ok {
		ra.WithSigv4SigningRegion(v.(string))
	}
	if v, ok := d.GetOk("rest_auth_sigv4_external_id"); ok {
		ra.WithSigv4ExternalId(v.(string))
	}
	return ra
}

func updateRestConfig(d *schema.ResourceData, set *sdk.CatalogIntegrationSetRequest) bool {
	restConfigKeys := []string{
		"rest_config_catalog_name", "rest_config_catalog_api_type",
		"rest_config_prefix", "rest_config_access_delegation_mode",
	}
	restAuthKeys := []string{
		"rest_auth_oauth_client_id", "rest_auth_oauth_client_secret",
		"rest_auth_oauth_token_uri", "rest_auth_oauth_allowed_scopes",
		"rest_auth_bearer_token", "rest_auth_sigv4_iam_role",
		"rest_auth_sigv4_signing_region", "rest_auth_sigv4_external_id",
	}

	changed := false
	for _, key := range restConfigKeys {
		if d.HasChange(key) {
			changed = true
			break
		}
	}
	if changed {
		if rc := buildRestConfig(d); rc != nil {
			set.WithRestConfig(*rc)
		}
	}

	authChanged := false
	for _, key := range restAuthKeys {
		if d.HasChange(key) {
			authChanged = true
			break
		}
	}
	if authChanged {
		if ra := buildRestAuthentication(d); ra != nil {
			set.WithRestAuthentication(*ra)
		}
	}

	return changed || authChanged
}
