package resources

import (
	"context"
	"errors"
	"fmt"

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

var catalogIntegrationAwsGlueSchema = func() map[string]*schema.Schema {
	awsGlueSchema := map[string]*schema.Schema{
		"glue_aws_role_arn": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Specifies the Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role to assume.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"glue_catalog_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "Specifies the ID of your AWS account.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"glue_region": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Description: joinWithSpace("Specifies the AWS region of your AWS Glue Data Catalog.",
				"You must specify a value for this attribute if your Snowflake account is not hosted on AWS.",
				"Otherwise, the default region is the Snowflake deployment region for the account."),
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("glue_region"),
			ValidateFunc:     validation.StringIsNotEmpty,
		},
		"catalog_namespace": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Description:  "Specifies the default AWS Glue Data Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	return collections.MergeMaps(
		catalogIntegrationCommonSchema(schemas.DescribeCatalogIntegrationAwsGlueDetailsSchema), awsGlueSchema)
}()

func CatalogIntegrationAwsGlue() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingCreateWrapper(resources.CatalogIntegrationAwsGlue, CreateCatalogIntegrationAwsGlue)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingReadWrapper(resources.CatalogIntegrationAwsGlue, ReadCatalogIntegrationAwsGlueFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingUpdateWrapper(resources.CatalogIntegrationAwsGlue, UpdateCatalogIntegrationAwsGlue)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingDeleteWrapper(resources.CatalogIntegrationAwsGlue, deleteCatalogIntegrationFunc())),
		Description:   "Resource used to manage AWS Glue catalog integration objects. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue).",

		Schema: catalogIntegrationAwsGlueSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogIntegrationAwsGlue, ImportCatalogIntegrationAwsGlue),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(catalogIntegrationAwsGlueSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(catalogIntegrationAwsGlueSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		),
	}
}

func ImportCatalogIntegrationAwsGlue(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.CatalogIntegrations.DescribeAwsGlueDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypeAWSGlue {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypeAWSGlue, details.CatalogSource)
	}

	return []*schema.ResourceData{d}, nil
}

func CreateCatalogIntegrationAwsGlue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(name)
	enabled := d.Get("enabled").(bool)
	glueAwsRoleArn := d.Get("glue_aws_role_arn").(string)
	glueCatalogId := d.Get("glue_catalog_id").(string)

	request := sdk.NewCreateCatalogIntegrationRequest(id, enabled)
	awsGlueRequest := sdk.NewAwsGlueParamsRequest(glueAwsRoleArn, glueCatalogId)
	errs := errors.Join(
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		stringAttributeCreateBuilder(d, "glue_region", awsGlueRequest.WithGlueRegion),
		stringAttributeCreateBuilder(d, "catalog_namespace", awsGlueRequest.WithCatalogNamespace),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.CatalogIntegrations.Create(ctx, request.WithAwsGlueCatalogSourceParams(*awsGlueRequest)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating AWS Glue catalog integration, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCatalogIntegrationAwsGlueFunc(false)(ctx, d, meta)
}

func ReadCatalogIntegrationAwsGlueFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
						Summary:  "Failed to query AWS Glue catalog integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("AWS Glue catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.CatalogIntegrations.DescribeAwsGlueDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe AWS Glue catalog integration (%s), err = %w", d.Id(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
				outputMapping{"glue_region", "glue_region", details.GlueRegion, details.GlueRegion, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set("name", details.Id.Name()),
			d.Set("enabled", details.Enabled),
			// not reading refresh_interval_seconds on purpose (handled as external change to describe output)
			d.Set("comment", details.Comment),
			d.Set("glue_aws_role_arn", details.GlueAwsRoleArn),
			d.Set("glue_catalog_id", details.GlueCatalogId),
			// not reading glue_region on purpose (handled as external change to describe output)
			d.Set("catalog_namespace", details.CatalogNamespace),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationAwsGlueDetailsToSchema(details)}),
		)

		return diag.FromErr(errs)
	}
}

func UpdateCatalogIntegrationAwsGlue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := handleCatalogIntegrationUpdate(ctx, d, meta); err != nil {
		return diag.FromErr(err)
	}
	return ReadCatalogIntegrationAwsGlueFunc(false)(ctx, d, meta)
}
