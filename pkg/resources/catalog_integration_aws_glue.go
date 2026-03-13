package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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

var catalogIntegrationAwsGlueSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier (i.e. name) of the catalog integration; must be unique in your account.",
	},
	"enabled": {
		Type:     schema.TypeBool,
		Required: true,
		Description: joinWithSpace("Specifies whether the catalog integration is available for use for Iceberg tables.",
			"`true` allows users to create new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration function normally.",
			"`false` prevents users from creating new Iceberg tables that reference this integration. Existing Iceberg tables that reference this integration cannot access the catalog in the table definition."),
	},
	"refresh_interval_seconds": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      30,
		Description:  "Specifies the number of seconds to wait between attempts to poll the external Iceberg catalog for metadata updates for automated refresh.",
		ValidateFunc: validation.All(validation.IntAtLeast(30), validation.IntAtMost(86400)),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the catalog integration.",
	},
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
		DiffSuppressFunc: IgnoreAfterCreation,
		ValidateFunc:     validation.StringIsNotEmpty,
	},
	"catalog_namespace": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  "Specifies the default AWS Glue Data Catalog namespace for all Iceberg tables that you associate with the catalog integration.",
		ValidateFunc: validation.StringIsNotEmpty,
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
			Schema: schemas.DescribeCatalogIntegrationAwsGlueDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func CatalogIntegrationAwsGlue() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.CatalogIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingCreateWrapper(resources.CatalogIntegrationAwsGlue, CreateCatalogIntegrationAwsGlue)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingReadWrapper(resources.CatalogIntegrationAwsGlue, ReadCatalogIntegrationAwsGlue)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingUpdateWrapper(resources.CatalogIntegrationAwsGlue, UpdateCatalogIntegrationAwsGlue)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CatalogIntegrationAwsGlueResource), TrackingDeleteWrapper(resources.CatalogIntegrationAwsGlue, deleteFunc)),
		Description:   "Resource used to manage AWS Glue catalog integration objects. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue).",

		Schema: catalogIntegrationAwsGlueSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CatalogIntegrationAwsGlue, ImportCatalogIntegrationAwsGlue),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(catalogIntegrationAwsGlueSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(catalogIntegrationAwsGlueSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment", "glue_aws_role_arn", "glue_catalog_id", "glue_region", "catalog_namespace"),
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
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
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
	return ReadCatalogIntegrationAwsGlue(ctx, d, meta)
}

func ReadCatalogIntegrationAwsGlue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	errs := errors.Join(
		d.Set("name", details.Id.Name()),
		d.Set("enabled", details.Enabled),
		d.Set("refresh_interval_seconds", details.RefreshIntervalSeconds),
		d.Set("comment", details.Comment),
		d.Set("glue_aws_role_arn", details.GlueAwsRoleArn),
		d.Set("glue_catalog_id", details.GlueCatalogId),
		d.Set("glue_region", details.GlueRegion),
		d.Set("catalog_namespace", details.CatalogNamespace),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.CatalogIntegrationToSchema(s)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CatalogIntegrationAwsGlueDetailsToSchema(details)}),
	)

	return diag.FromErr(errs)
}

func UpdateCatalogIntegrationAwsGlue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewCatalogIntegrationSetRequest()

	errs := errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		intAttributeUpdateSetOnly(d, "refresh_interval_seconds", &set.RefreshIntervalSeconds),
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*set, *sdk.NewCatalogIntegrationSetRequest()) {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating AWS Glue catalog integration (%s), err = %w", d.Id(), err))
		}
	}
	return ReadCatalogIntegrationAwsGlue(ctx, d, meta)
}
