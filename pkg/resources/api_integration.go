package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var apiIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the API integration. This name follows the rules for Object Identifiers. The name should be unique among api integrations in your account.",
	},
	"api_provider": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"aws_api_gateway", "aws_private_api_gateway", "azure_api_management", "aws_gov_api_gateway", "aws_gov_private_api_gateway", "google_api_gateway"}, false),
		ForceNew:     true,
		Description:  "Specifies the HTTPS proxy service type.",
	},
	"api_aws_role_arn": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "ARN of a cloud platform role.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (API_AWS_IAM_USER_ARN)
	"api_aws_iam_user_arn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Snowflake user that will attempt to assume the AWS role.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (API_AWS_EXTERNAL_ID)
	"api_aws_external_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The external ID that Snowflake will use when assuming the AWS role.",
	},
	"azure_tenant_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies the ID for your Office 365 tenant that all Azure API Management instances belong to.",
	},
	"azure_ad_application_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The 'Application (client) id' of the Azure AD app for your remote service.",
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (AZURE_MULTI_TENANT_APP_NAME)
	"azure_multi_tenant_app_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	// Computed. Info you get by issuing a 'DESCRIBE INTEGRATION <name>' command (AZURE_CONSENT_URL)
	"azure_consent_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"google_audience": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The audience claim when generating the JWT (JSON Web Token) to authenticate to the Google API Gateway.",
	},
	"api_gcp_service_account": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The service account used for communication with the Google API Gateway.",
		Computed:    true,
	},
	"api_allowed_prefixes": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Explicitly limits external functions that use the integration to reference one or more HTTPS proxy service endpoints and resources within those proxies.",
		MinItems:    1,
	},
	"api_blocked_prefixes": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Lists the endpoints and resources in the HTTPS proxy service that are not allowed to be called from Snowflake.",
	},
	"api_key": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "The API key (also called a “subscription key”).",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Specifies whether this API integration is enabled or disabled. If the API integration is disabled, any external function that relies on it will not work.",
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the API integration was created.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func APIIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ApiIntegrationResource), TrackingCreateWrapper(resources.ApiIntegration, CreateAPIIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ApiIntegrationResource), TrackingReadWrapper(resources.ApiIntegration, ReadAPIIntegration)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ApiIntegrationResource), TrackingUpdateWrapper(resources.ApiIntegration, UpdateAPIIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ApiIntegrationResource), TrackingDeleteWrapper(resources.ApiIntegration, deleteFunc)),

		Schema: apiIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func toApiIntegrationEndpointPrefix(paths []string) []sdk.ApiIntegrationEndpointPrefix {
	allowedPrefixes := make([]sdk.ApiIntegrationEndpointPrefix, len(paths))
	for i, prefix := range paths {
		allowedPrefixes[i] = sdk.ApiIntegrationEndpointPrefix{Path: prefix}
	}
	return allowedPrefixes
}

// CreateAPIIntegration implements schema.CreateFunc.
func CreateAPIIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	allowedPrefixesRaw := expandStringList(d.Get("api_allowed_prefixes").([]interface{}))
	allowedPrefixes := toApiIntegrationEndpointPrefix(allowedPrefixesRaw)

	createRequest := sdk.NewCreateApiIntegrationRequest(id, allowedPrefixes, enabled)

	if _, ok := d.GetOk("api_blocked_prefixes"); ok {
		blockedPrefixesRaw := expandStringList(d.Get("api_blocked_prefixes").([]interface{}))
		createRequest.WithApiBlockedPrefixes(toApiIntegrationEndpointPrefix(blockedPrefixesRaw))
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(v.(string))
	}

	apiProvider := d.Get("api_provider").(string)
	switch apiProvider {
	case "aws_api_gateway", "aws_private_api_gateway", "aws_gov_api_gateway", "aws_gov_private_api_gateway":
		roleArn, ok := d.GetOk("api_aws_role_arn")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use AWS api provider you must specify an api_aws_role_arn"))
		}
		awsParams := sdk.NewAwsApiParamsRequest(sdk.ApiIntegrationAwsApiProviderType(apiProvider), roleArn.(string))
		if v, ok := d.GetOk("api_key"); ok {
			awsParams.WithApiKey(v.(string))
		}
		createRequest.WithAwsApiProviderParams(*awsParams)
	case "azure_api_management":
		tenantId, ok := d.GetOk("azure_tenant_id")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use the Azure api provider you must specify an azure_tenant_id"))
		}
		applicationId, ok := d.GetOk("azure_ad_application_id")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use the Azure api provider you must specify an azure_ad_application_id"))
		}
		azureParams := sdk.NewAzureApiParamsRequest(tenantId.(string), applicationId.(string))
		if v, ok := d.GetOk("api_key"); ok {
			azureParams.WithApiKey(v.(string))
		}
		createRequest.WithAzureApiProviderParams(*azureParams)
	case "google_api_gateway":
		audience, ok := d.GetOk("google_audience")
		if !ok {
			return diag.FromErr(fmt.Errorf("if you use GCP api provider you must specify a google_audience"))
		}
		googleParams := sdk.NewGoogleApiParamsRequest(audience.(string))
		createRequest.WithGoogleApiProviderParams(*googleParams)
	default:
		return diag.FromErr(fmt.Errorf("unexpected provider %v", apiProvider))
	}

	err := client.ApiIntegrations.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating api integration: %w", err))
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadAPIIntegration(ctx, d, meta)
}

// ReadAPIIntegration implements schema.ReadFunc.
func ReadAPIIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.ApiIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query api integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Api integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	// Note: category must be API or something is broken
	if c := integration.Category; c != "API" {
		return diag.FromErr(fmt.Errorf("expected %v to be an api integration, got %v", id, c))
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", integration.CreatedOn.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(err)
	}

	// Some properties come from the DESCRIBE INTEGRATION call
	integrationProperties, err := client.ApiIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe api integration: %w", err))
	}

	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value
		switch name {
		case "ENABLED":
			// We set this using the SHOW INTEGRATION call so let's ignore it here
		case "API_ALLOWED_PREFIXES":
			if err := d.Set("api_allowed_prefixes", strings.Split(value, ",")); err != nil {
				return diag.FromErr(err)
			}
		case "API_BLOCKED_PREFIXES":
			if val := value; val != "" {
				if err := d.Set("api_blocked_prefixes", strings.Split(val, ",")); err != nil {
					return diag.FromErr(err)
				}
			}
		case "API_AWS_IAM_USER_ARN":
			if err := d.Set("api_aws_iam_user_arn", value); err != nil {
				return diag.FromErr(err)
			}
		case "API_AWS_ROLE_ARN":
			if err := d.Set("api_aws_role_arn", value); err != nil {
				return diag.FromErr(err)
			}
		case "API_AWS_EXTERNAL_ID":
			if err := d.Set("api_aws_external_id", value); err != nil {
				return diag.FromErr(err)
			}
		case "AZURE_CONSENT_URL":
			if err := d.Set("azure_consent_url", value); err != nil {
				return diag.FromErr(err)
			}
		case "AZURE_MULTI_TENANT_APP_NAME":
			if err := d.Set("azure_multi_tenant_app_name", value); err != nil {
				return diag.FromErr(err)
			}
		case "AZURE_TENANT_ID":
			if err := d.Set("azure_tenant_id", value); err != nil {
				return diag.FromErr(err)
			}
		case "AZURE_AD_APPLICATION_ID":
			if err := d.Set("azure_ad_application_id", value); err != nil {
				return diag.FromErr(err)
			}
		case "GOOGLE_AUDIENCE":
			if err := d.Set("google_audience", value); err != nil {
				return diag.FromErr(err)
			}
		case "API_GCP_SERVICE_ACCOUNT":
			if err := d.Set("api_gcp_service_account", value); err != nil {
				return diag.FromErr(err)
			}
		case "API_PROVIDER":
			if err := d.Set("api_provider", strings.ToLower(value)); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected api integration property %v returned from Snowflake", name)
		}
	}

	return diag.FromErr(err)
}

// UpdateAPIIntegration implements schema.UpdateFunc.
func UpdateAPIIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	setRequest := sdk.NewApiIntegrationSetRequest()
	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("api_allowed_prefixes") {
		runSetStatement = true
		setRequest.WithApiAllowedPrefixes(toApiIntegrationEndpointPrefix(expandStringList(d.Get("api_allowed_prefixes").([]interface{}))))
	}

	if d.HasChange("comment") {
		runSetStatement = true
		setRequest.WithComment(d.Get("comment").(string))
	}

	// We need to UNSET this if we remove all api blocked prefixes.
	if d.HasChange("api_blocked_prefixes") {
		v := d.Get("api_blocked_prefixes").([]interface{})
		if len(v) == 0 {
			err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).WithUnset(*sdk.NewApiIntegrationUnsetRequest().WithApiBlockedPrefixes(true)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting api_blocked_prefixes: %w", err))
			}
		} else {
			runSetStatement = true
			setRequest.WithApiBlockedPrefixes(toApiIntegrationEndpointPrefix(expandStringList(v)))
		}
	}

	apiProvider := d.Get("api_provider").(string)
	switch apiProvider {
	case "aws_api_gateway", "aws_private_api_gateway", "aws_gov_api_gateway", "aws_gov_private_api_gateway":
		awsParams := sdk.NewSetAwsApiParamsRequest()
		if d.HasChange("api_aws_role_arn") {
			awsParams.WithApiAwsRoleArn(d.Get("api_aws_role_arn").(string))
		}
		if d.HasChange("api_key") {
			awsParams.WithApiKey(d.Get("api_key").(string))
		}
		if *awsParams != *sdk.NewSetAwsApiParamsRequest() {
			runSetStatement = true
			setRequest.WithAwsParams(*awsParams)
		}
	case "azure_api_management":
		azureParams := sdk.NewSetAzureApiParamsRequest()
		if d.HasChange("azure_tenant_id") {
			azureParams.WithAzureTenantId(d.Get("azure_tenant_id").(string))
		}
		if d.HasChange("azure_ad_application_id") {
			azureParams.WithAzureAdApplicationId(d.Get("azure_ad_application_id").(string))
		}
		if d.HasChange("api_key") {
			azureParams.WithApiKey(d.Get("api_key").(string))
		}
		if *azureParams != *sdk.NewSetAzureApiParamsRequest() {
			runSetStatement = true
			setRequest.WithAzureParams(*azureParams)
		}
	case "google_api_gateway":
		if d.HasChange("google_audience") {
			runSetStatement = true
			googleParams := sdk.NewSetGoogleApiParamsRequest(d.Get("google_audience").(string))
			setRequest.WithGoogleParams(*googleParams)
		}
	default:
		return diag.FromErr(fmt.Errorf("unexpected provider %v", apiProvider))
	}

	if runSetStatement {
		err := client.ApiIntegrations.Alter(ctx, sdk.NewAlterApiIntegrationRequest(id).WithSet(*setRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating api integration: %w", err))
		}
	}

	return ReadAPIIntegration(ctx, d, meta)
}
