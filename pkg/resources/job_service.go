package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var jobServiceSchema = func() map[string]*schema.Schema {
	// TODO(SNOW-2129584): add async field, or handle sync jobs in a separate resource/data source
	return serviceBaseSchema(true)
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
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.JobServiceResource), TrackingReadWrapper(resources.JobService, ReadJobServiceFunc(true))),
		// No UpdateContext because altering job service is not supported in Snowflake.
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.JobServiceResource), TrackingDeleteWrapper(resources.JobService, deleteFunc)),
		Description: joinWithSpace(
			"Resource used to manage job services. For more information, check [services documentation](https://docs.snowflake.com/en/sql-reference/sql/execute-job-service).",
			"Executes a Snowpark Container Services service as a job. A service, created using `CREATE SERVICE`, is long-running and you must explicitly stop it when it is no longer needed.",
			"On the other hand, a job, created using EXECUTE JOB SERVICE (with `ASYNC=TRUE` in this resource), returns immediately while the job is running.",
			"See [Working with services](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services) developer guide for more details.",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.JobService, customdiff.All(
			ComputedIfAnyAttributeChanged(jobServiceSchema, ShowOutputAttributeName, "query_warehouse", "comment"),
			ComputedIfAnyAttributeChanged(jobServiceSchema, DescribeOutputAttributeName, "query_warehouse", "comment"),
			RecreateWhenServiceTypeChangedExternally(sdk.ServiceTypeJobService),
		)),

		Schema: jobServiceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.JobService, ImportServiceFunc(jobServiceCustomFieldsHandler)),
		},

		Timeouts: defaultTimeouts,
	}
}

func ReadJobServiceFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return ReadServiceCommonFunc(withExternalChangesMarking, jobServiceOutputMappingsFunc, nil)
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

	request := sdk.NewExecuteJobServiceRequest(computePoolId, id)
	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "from_specification", request.WithJobServiceFromSpecification, ToJobServiceFromSpecificationRequest),
		attributeMappedValueCreateBuilder(d, "from_specification_template", request.WithJobServiceFromSpecificationTemplate, ToJobServiceFromSpecificationTemplateRequest),
		accountObjectIdentifierAttributeCreate(d, "query_warehouse", &request.QueryWarehouse),
		attributeMappedValueCreateBuilder(d, "external_access_integrations", request.WithExternalAccessIntegrations, ToServiceExternalAccessIntegrationsRequest),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	// TODO(SNOW-2129584): adjust default async option
	request.WithAsync(true)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if err := client.Services.ExecuteJob(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadJobServiceFunc(false)(ctx, d, meta)
}

func jobServiceCustomFieldsHandler(d *schema.ResourceData, service *sdk.Service) error {
	// noop, as job service has no custom fields
	return nil
}

func jobServiceOutputMappingsFunc(service *sdk.Service) []outputMapping {
	// noop, as job service has no custom fields
	return nil
}
