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

var apiIntegrationAmazonApiGatewaySchema = func() map[string]*schema.Schema {
	apiIntegrationAmazonApiGateway := map[string]*schema.Schema{
		"api_provider": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: NormalizeAndCompare(sdk.ToApiIntegrationAwsApiProviderType),
			ValidateDiagFunc: sdkValidation(sdk.ToApiIntegrationAwsApiProviderType),
			Description:      fmt.Sprintf("Specifies the type of AWS gateway. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AsStringList(sdk.AllApiIntegrationAwsApiProviderTypes))),
		},
		"api_aws_role_arn": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Amazon Resource Name (ARN) of the IAM role that grants Snowflake permission to call the API endpoint.",
		},
		"api_key": {
			Type:             schema.TypeString,
			Optional:         true,
			Sensitive:        true,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("api_key"),
			Description:      "Specifies the API key (secret) that Snowflake uses to authenticate when making calls to the proxy service.",
		},
		"api_aws_iam_user_arn": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The Snowflake-created IAM user ARN used for the AWS integration. Grant this ARN permissions in the IAM role policy.",
		},
		"api_aws_external_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The external ID used by Snowflake to establish a trust relationship with the AWS IAM role.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeAmazonApiGatewayApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationAmazonApiGateway)
}()

func ApiIntegrationAmazonApiGateway() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ApiIntegrationAmazonApiGatewayResource), TrackingCreateWrapper(resources.ApiIntegrationAmazonApiGateway, CreateApiIntegrationAmazonApiGateway)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ApiIntegrationAmazonApiGatewayResource), TrackingReadWrapper(resources.ApiIntegrationAmazonApiGateway, GetReadApiIntegrationAmazonApiGatewayFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ApiIntegrationAmazonApiGatewayResource), TrackingUpdateWrapper(resources.ApiIntegrationAmazonApiGateway, UpdateApiIntegrationAmazonApiGateway)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ApiIntegrationAmazonApiGatewayResource), TrackingDeleteWrapper(resources.ApiIntegrationAmazonApiGateway, deleteFunc)),
		Description:   "Resource used to manage API integration Amazon API Gateway objects. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationAmazonApiGatewaySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationAmazonApiGateway, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationAmazonApiGatewaySchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationAmazonApiGatewaySchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "api_aws_role_arn", "api_key"),
		),
	}
}

func CreateApiIntegrationAmazonApiGateway(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	providerType, err := sdk.ToApiIntegrationAwsApiProviderType(d.Get("api_provider").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	awsParams := sdk.NewAwsApiParamsRequest(providerType, d.Get("api_aws_role_arn").(string))

	if err = stringAttributeCreate(d, "api_key", &awsParams.ApiKey); err != nil {
		return diag.FromErr(err)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithAwsApiProviderParams(*awsParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating AWS API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return GetReadApiIntegrationAmazonApiGatewayFunc(false)(ctx, d, meta)
}

func GetReadApiIntegrationAmazonApiGatewayFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
						Summary:  "Failed to query API integration Amazon API Gateway. Marking the resource as removed.",
						Detail:   fmt.Sprintf("API integration Amazon API Gateway id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		awsDetails, err := client.ApiIntegrations.DescribeAwsDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe API integration Amazon API Gateway (%s): %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"api_key", "api_key", awsDetails.ApiKey, awsDetails.ApiKey, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			handleApiIntegrationCommonRead(d, id, s, awsDetails.AllowedPrefixes, awsDetails.BlockedPrefixes),
			d.Set("api_provider", awsDetails.ApiProvider),
			d.Set("api_aws_role_arn", awsDetails.ApiAwsRoleArn),
			d.Set("api_aws_iam_user_arn", awsDetails.ApiAwsIamUserArn),
			d.Set("api_aws_external_id", awsDetails.ApiAwsExternalId),
			// api_key intentionally omitted — handled by external changes marking above
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationAmazonApiGatewayDetailsToSchema(awsDetails)}),
		)
		return diag.FromErr(errs)
	}
}

func UpdateApiIntegrationAmazonApiGateway(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()
	awsSet := sdk.NewSetAwsApiParamsRequest()
	awsUnset := sdk.NewUnsetAwsApiParamsRequest()

	errs := errors.Join(
		handleApiIntegrationCommonUpdate(d, set, unset),
		stringAttributeUpdateSetOnlyNotEmpty(d, "api_aws_role_arn", &awsSet.ApiAwsRoleArn),
		stringAttributeUpdate(d, "api_key", &awsSet.ApiKey, &awsUnset.ApiKey),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*awsSet, *sdk.NewSetAwsApiParamsRequest()) {
		set.WithAwsParams(*awsSet)
	}
	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating AWS API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*awsUnset, *sdk.NewUnsetAwsApiParamsRequest()) {
		unset.WithAwsParams(*awsUnset)
	}
	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating AWS API integration: %w", err))
		}
	}

	return GetReadApiIntegrationAmazonApiGatewayFunc(false)(ctx, d, meta)
}
