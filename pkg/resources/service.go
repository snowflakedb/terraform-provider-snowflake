package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var serviceSchema = func() map[string]*schema.Schema {
	serviceSchema := map[string]*schema.Schema{
		"auto_suspend_secs": {
			Type:             schema.TypeInt,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend_secs"),
			Description:      "Specifies the number of seconds of inactivity (service is idle) after which Snowflake automatically suspends the service.",
			Default:          IntDefault,
		},
		"auto_resume": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
			Description:      booleanStringFieldDescription("Specifies whether to automatically resume a service."),
			Default:          BooleanDefault,
		},
		"min_instances": {
			Type:             schema.TypeInt,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_instances"),
			Description:      "Specifies the minimum number of service instances to run.",
		},
		"min_ready_instances": {
			Type:             schema.TypeInt,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("min_ready_instances"),
			Description:      "Indicates the minimum service instances that must be ready for Snowflake to consider the service is ready to process requests.",
		},
		"max_instances": {
			Type:             schema.TypeInt,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("max_instances"),
			Description:      "Specifies the maximum number of service instances to run.",
		},
	}
	return collections.MergeMaps(serviceBaseSchema(false), serviceSchema)
}()

func Service() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Services.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ServiceResource), TrackingCreateWrapper(resources.Service, CreateService)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ServiceResource), TrackingReadWrapper(resources.Service, ReadServiceFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ServiceResource), TrackingUpdateWrapper(resources.Service, UpdateService)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ServiceResource), TrackingDeleteWrapper(resources.Service, deleteFunc)),
		Description:   "Resource used to manage services. For more information, check [services documentation](https://docs.snowflake.com/en/sql-reference/sql/create-service).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Service, customdiff.All(
			ComputedIfAnyAttributeChanged(serviceSchema, ShowOutputAttributeName, "auto_suspend_secs", "auto_resume", "min_instances", "max_instances", "min_ready_instances", "query_warehouse", "comment"),
			ComputedIfAnyAttributeChanged(serviceSchema, DescribeOutputAttributeName, "auto_suspend_secs", "auto_resume", "min_instances", "max_instances", "min_ready_instances", "query_warehouse", "comment"),
		)),

		Schema: serviceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Service, ImportServiceFunc(serviceCustomFieldsHandler)),
		},

		Timeouts: defaultTimeouts,
	}
}

func ReadServiceFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return ReadServiceCommonFunc(withExternalChangesMarking, serviceOutputMappingsFunc, []string{"max_instances", "min_instances", "min_ready_instances", "auto_resume", "auto_suspend_secs"})
}

func CreateService(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	request := sdk.NewCreateServiceRequest(id, computePoolId)
	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "from_specification", request.WithFromSpecification, ToServiceFromSpecificationRequest),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend_secs", request.WithAutoSuspendSecs),
		intAttributeCreateBuilder(d, "min_instances", request.WithMinInstances),
		intAttributeCreateBuilder(d, "max_instances", request.WithMaxInstances),
		intAttributeCreateBuilder(d, "min_ready_instances", request.WithMinReadyInstances),
		accountObjectIdentifierAttributeCreate(d, "query_warehouse", &request.QueryWarehouse),
		booleanStringAttributeCreateBuilder(d, "auto_resume", request.WithAutoResume),
		attributeMappedValueCreateBuilder(d, "external_access_integrations", request.WithExternalAccessIntegrations, ToServiceExternalAccessIntegrationsRequest),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if err := client.Services.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadServiceFunc(false)(ctx, d, meta)
}

func UpdateService(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("from_specification") {
		if v, ok := d.GetOk("from_specification"); ok {
			spec, err := ToServiceFromSpecificationRequest(v)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(id).WithFromSpecification(spec)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	set, unset := sdk.NewServiceSetRequest(), sdk.NewServiceUnsetRequest()
	errs := errors.Join(
		// name, schema, database, and compute_pool are handled by ForceNew.
		intAttributeWithSpecialDefaultUpdate(d, "auto_suspend_secs", &set.AutoSuspendSecs, &unset.AutoSuspendSecs),
		intAttributeUpdate(d, "min_instances", &set.MinInstances, &unset.MinInstances),
		intAttributeUpdate(d, "max_instances", &set.MaxInstances, &unset.MaxInstances),
		intAttributeUpdate(d, "min_ready_instances", &set.MinReadyInstances, &unset.MinReadyInstances),
		accountObjectIdentifierAttributeUpdate(d, "query_warehouse", &set.QueryWarehouse, &unset.QueryWarehouse),
		booleanStringAttributeUpdate(d, "auto_resume", &set.AutoResume, &unset.AutoResume),
		attributeMappedValueUpdate(d, "external_access_integrations", &set.ExternalAccessIntegrations, &unset.ExternalAccessIntegrations, ToServiceExternalAccessIntegrationsRequest),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if (*set != sdk.ServiceSetRequest{}) {
		if err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if (*unset != sdk.ServiceUnsetRequest{}) {
		if err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadServiceFunc(false)(ctx, d, meta)
}

func serviceCustomFieldsHandler(d *schema.ResourceData, service *sdk.Service) error {
	return errors.Join(
		d.Set("max_instances", service.MaxInstances),
		d.Set("min_instances", service.MinInstances),
		d.Set("min_ready_instances", service.MinReadyInstances),
		d.Set("auto_resume", booleanStringFromBool(service.AutoResume)),
		d.Set("auto_suspend_secs", service.AutoSuspendSecs),
	)
}

func serviceOutputMappingsFunc(service *sdk.Service) []outputMapping {
	return []outputMapping{
		{"auto_resume", "auto_resume", service.AutoResume, booleanStringFromBool(service.AutoResume), nil},
		{"auto_suspend_secs", "auto_suspend_secs", service.AutoSuspendSecs, service.AutoSuspendSecs, nil},
		{"min_instances", "min_instances", service.MinInstances, service.MinInstances, nil},
		{"max_instances", "max_instances", service.MaxInstances, service.MaxInstances, nil},
		{"min_ready_instances", "min_ready_instances", service.MinReadyInstances, service.MinReadyInstances, nil},
	}
}
