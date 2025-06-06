package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var jobServiceSchema = func() map[string]*schema.Schema {
	jobServiceSchema := map[string]*schema.Schema{
		"async": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_async_job"),
			Description:      booleanStringFieldDescription("Specifies whether to execute the job service asynchronously."),
			Default:          BooleanDefault,
		},
	}
	return collections.MergeMaps(serviceBaseSchema(true), jobServiceSchema)
}()

func JobService() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Services.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.JobServiceResource), TrackingCreateWrapper(resources.JobService, CreateJobService)),
		// ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.JobServiceResource), TrackingReadWrapper(resources.JobService, ReadJobServiceFunc(true))),
		ReadContext: PreviewFeatureReadContextWrapper(string(previewfeatures.JobServiceResource), TrackingReadWrapper(resources.JobService, ReadServiceGenericFunc(true, jobServiceOutputMappingsFunc, []string{"async"}))),
		// No UpdateContext because altering job service is not supported in Snowflake.
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.JobServiceResource), TrackingDeleteWrapper(resources.JobService, deleteFunc)),
		Description:   "Resource used to manage job services. For more information, check [services documentation](https://docs.snowflake.com/en/sql-reference/sql/execute-job-service).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.JobService, customdiff.All(
			ComputedIfAnyAttributeChanged(jobServiceSchema, ShowOutputAttributeName, "query_warehouse", "comment", "async"),
			ComputedIfAnyAttributeChanged(jobServiceSchema, DescribeOutputAttributeName, "query_warehouse", "comment", "async"),
		)),

		Schema: jobServiceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.JobService, ImportServiceFunc(jobServiceCustomFieldsHandler)),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateJobService(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)
	computePoolRaw := d.Get("compute_pool").(string)
	computePoolId, err := sdk.ParseAccountObjectIdentifier(computePoolRaw)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewExecuteJobServiceServiceRequest(computePoolId, id)
	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "from_specification", request.WithJobServiceFromSpecification, ToJobServiceFromSpecificationRequest),
		accountObjectIdentifierAttributeCreate(d, "query_warehouse", &request.QueryWarehouse),
		attributeMappedValueCreateBuilder(d, "external_access_integrations", request.WithExternalAccessIntegrations, ToServiceExternalAccessIntegrationsRequest),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		booleanStringAttributeCreateBuilder(d, "async", request.WithAsync),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if err := client.Services.ExecuteJobService(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadJobServiceFunc(false)(ctx, d, meta)
}

// TODO: merge and remove
func ReadJobServiceFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		service, err := client.Services.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query service. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Service id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		serviceDetails, err := client.Services.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			var warehouseFullyQualifiedName string
			if service.QueryWarehouse != nil {
				warehouseFullyQualifiedName = service.QueryWarehouse.FullyQualifiedName()
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"async", "async", service.IsAsyncJob, booleanStringFromBool(service.IsAsyncJob), nil},
				outputMapping{"query_warehouse", "query_warehouse", warehouseFullyQualifiedName, warehouseFullyQualifiedName, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, serviceSchema, []string{
			"async",
			"query_warehouse",
		}); err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ServiceToSchema(service)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ServiceDetailsToSchema(serviceDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("compute_pool", service.ComputePool.FullyQualifiedName()),
			d.Set("external_access_integrations", collections.Map(service.ExternalAccessIntegrations, func(id sdk.AccountObjectIdentifier) string { return id.FullyQualifiedName() })),
			d.Set("comment", service.Comment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func jobServiceCustomFieldsHandler(d *schema.ResourceData, service *sdk.Service) error {
	return errors.Join(
		d.Set("async", booleanStringFromBool(service.IsAsyncJob)),
	)
}

func jobServiceOutputMappingsFunc(service *sdk.Service) []outputMapping {
	return []outputMapping{
		{"async", "async", service.IsAsyncJob, booleanStringFromBool(service.IsAsyncJob), nil},
	}
}
