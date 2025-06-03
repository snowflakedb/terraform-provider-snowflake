package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var serviceSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the service."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the service."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the service; must be unique for the schema in which the service is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"compute_pool": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the name of the compute pool in your account on which to run the service."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"from_specification": {
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Description: "Specifies the service specification to use for the service. Note that external changes on this field and nested fields are not detected.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"stage": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The stage containing the service specification file. At symbol (`@`) is added automatically.",
				},
				"path": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The path to the service specification file on the given stage.",
				},
				"file": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The file name of the service specification.",
				},
				"text": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The embedded text of the service specification.",
				},
			},
		},
	},
	// TODO (next PR): add from_specification_template
	"auto_suspend_secs": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend_secs"),
		Description:      "Specifies the number of seconds of inactivity (service is idle) after which Snowflake automatically suspends the service.",
		Default:          IntDefault,
	},
	"external_access_integrations": {
		Type:        schema.TypeSet,
		Optional:    true,
		MinItems:    1,
		Description: "Specifies the names of the external access integrations that allow your service to access external sites.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("external_access_integrations"),
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
	"query_warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      blocklistedCharactersFieldDescription("Warehouse to use if a service container connects to Snowflake to execute a query but does not explicitly specify a warehouse to use."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the service.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SERVICES` for the given service.",
		Elem: &schema.Resource{
			Schema: schemas.ShowServiceSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SERVICE` for the given service.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeServiceSchema,
		},
	},
}

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
			StateContext: TrackingImportWrapper(resources.Service, ImportName[sdk.SchemaObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
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

func ReadServiceFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
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
				outputMapping{"auto_resume", "auto_resume", service.AutoResume, booleanStringFromBool(service.AutoResume), nil},
				outputMapping{"auto_suspend_secs", "auto_suspend_secs", service.AutoSuspendSecs, service.AutoSuspendSecs, nil},
				outputMapping{"min_instances", "min_instances", service.MinInstances, service.MinInstances, nil},
				outputMapping{"max_instances", "max_instances", service.MaxInstances, service.MaxInstances, nil},
				outputMapping{"min_ready_instances", "min_ready_instances", service.MinReadyInstances, service.MinReadyInstances, nil},
				outputMapping{"query_warehouse", "query_warehouse", warehouseFullyQualifiedName, warehouseFullyQualifiedName, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, serviceSchema, []string{
			"auto_resume",
			"auto_suspend_secs",
			"min_instances",
			"max_instances",
			"min_ready_instances",
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

func ToServiceExternalAccessIntegrationsRequest(value any) (sdk.ServiceExternalAccessIntegrationsRequest, error) {
	raw := expandStringList(value.(*schema.Set).List())
	integrations := make([]sdk.AccountObjectIdentifier, len(raw))
	for i, v := range raw {
		integrations[i] = sdk.NewAccountObjectIdentifier(v)
	}
	return sdk.ServiceExternalAccessIntegrationsRequest{
		ExternalAccessIntegrations: integrations,
	}, nil
}

func ToServiceFromSpecificationRequest(value any) (sdk.ServiceFromSpecificationRequest, error) {
	serviceFromSpecification := sdk.ServiceFromSpecificationRequest{}
	for _, v := range value.([]any) {
		fromSpecificationConfig := v.(map[string]any)
		if text := fromSpecificationConfig["text"].(string); text != "" {
			serviceFromSpecification.Specification = &text
		}
		if stageRaw := fromSpecificationConfig["stage"].(string); stageRaw != "" {
			stage, err := sdk.ParseSchemaObjectIdentifier(stageRaw)
			if err != nil {
				return sdk.ServiceFromSpecificationRequest{}, err
			}
			var path string
			if value := fromSpecificationConfig["path"].(string); value != "" {
				path = value
			}
			serviceFromSpecification.Location = sdk.NewStageLocation(stage, path)
		}
		if file := fromSpecificationConfig["file"].(string); file != "" {
			serviceFromSpecification.SpecificationFile = &file
		}
	}
	return serviceFromSpecification, nil
}
