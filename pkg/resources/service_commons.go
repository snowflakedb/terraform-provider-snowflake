package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func serviceFromSpecificationTemplateSchema(allFieldsForceNew, isTemplate bool) map[string]*schema.Schema {
	fieldName := "from_specification"
	if isTemplate {
		fieldName = "from_specification_template"
	}
	objectNameInDescription := "specification"
	if isTemplate {
		objectNameInDescription = "specification template"
	}
	stageFieldName := fmt.Sprintf("%s.0.stage", fieldName)
	fileFieldName := fmt.Sprintf("%s.0.file", fieldName)
	textFieldName := fmt.Sprintf("%s.0.text", fieldName)
	subSchema := map[string]*schema.Schema{
		// Accepted configurations:
		// - stage, file, and optional path
		// - text
		"stage": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
			ForceNew:         allFieldsForceNew,
			RequiredWith:     []string{fileFieldName},
			Description:      fmt.Sprintf("The fully qualified name of the stage containing the service %s file. At symbol (`@`) is added automatically. Example: `\"\\\"<db_name>\\\".\\\"<schema_name>\\\".\\\"<schema_name>\\\"\"`.", objectNameInDescription),
		},
		"path": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     allFieldsForceNew,
			RequiredWith: []string{stageFieldName, fileFieldName},
			Description:  fmt.Sprintf("The path to the service %s file on the given stage. When the path is specified, the `/` character is automatically added as a path prefix. Example: `path/to/spec`.", objectNameInDescription),
		},
		"file": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     allFieldsForceNew,
			RequiredWith: []string{stageFieldName},
			ExactlyOneOf: []string{textFieldName, fileFieldName},
			Description:  fmt.Sprintf("The file name of the service %s. Example: `spec.yaml`.", objectNameInDescription),
		},
		"text": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     allFieldsForceNew,
			Description:  fmt.Sprintf("The embedded text of the service %s.", objectNameInDescription),
			ExactlyOneOf: []string{textFieldName, fileFieldName},
		},
	}
	if isTemplate {
		subSchema["using"] = &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The name of the template variable. The provider wraps it in double quotes by default, so be aware of that while referencing the argument in the spec definition.",
					},
					"value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The value to assign to the variable in the template. The provider wraps it in `$$` by default, so be aware of that while referencing the argument in the spec definition. The value must either be alphanumeric or valid JSON.",
					},
				},
			},
			Required:    true,
			Description: "List of the specified template variables and the values of those variables.",
		}
	}

	return map[string]*schema.Schema{
		fieldName: {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			ForceNew:    allFieldsForceNew,
			Description: fmt.Sprintf("Specifies the service %s to use for the service. Note that external changes on this field and nested fields are not detected. Use correctly formatted YAML files. Watch out for the space/tabs indentation. See [service specification](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/specification-reference#general-guidelines) for more information.", objectNameInDescription),
			Elem: &schema.Resource{
				Schema: subSchema,
			},
			ExactlyOneOf: []string{"from_specification", "from_specification_template"},
		},
	}
}

func serviceBaseSchema(allFieldsForceNew bool) map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
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
			Description:      blocklistedCharactersFieldDescription("Specifies the name of the compute pool in your account on which to run the service. Identifiers with special or lower-case characters are not supported. This limitation in the provider follows the limitation in Snowflake (see [docs](https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool))."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"external_access_integrations": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			ForceNew:    allFieldsForceNew,
			Description: "Specifies the names of the external access integrations that allow your service to access external sites.",
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
				DiffSuppressFunc: suppressIdentifierQuoting,
				ForceNew:         allFieldsForceNew,
			},
			DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("external_access_integrations"),
		},
		"query_warehouse": {
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         allFieldsForceNew,
			Description:      blocklistedCharactersFieldDescription("Warehouse to use if a service container connects to Snowflake to execute a query but does not explicitly specify a warehouse to use."),
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    allFieldsForceNew,
			Description: "Specifies a comment for the service.",
		},
		"service_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies a type for the service. This field is used for checking external changes and recreating the resources if needed.",
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
	return collections.MergeMaps(schema, serviceFromSpecificationTemplateSchema(allFieldsForceNew, false), serviceFromSpecificationTemplateSchema(allFieldsForceNew, true))
}

func ImportServiceFunc(customFieldsHandler func(d *schema.ResourceData, service *sdk.Service) error) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return nil, err
		}

		service, err := client.Services.ShowByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if service.QueryWarehouse != nil {
			if err := d.Set("query_warehouse", service.QueryWarehouse.FullyQualifiedName()); err != nil {
				return nil, err
			}
		}

		errs := errors.Join(
			d.Set("name", service.Name),
			d.Set("schema", service.SchemaName),
			d.Set("database", service.DatabaseName),
			customFieldsHandler(d, service),
		)
		if errs != nil {
			return nil, errs
		}
		return []*schema.ResourceData{d}, nil
	}
}

func ReadServiceCommonFunc(withExternalChangesMarking bool, extraOutputMappingsFunc func(service *sdk.Service) []outputMapping, extraSetStateToValuesFromConfigFields []string) schema.ReadContextFunc {
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
			outputMappings := append(extraOutputMappingsFunc(service), outputMapping{"query_warehouse", "query_warehouse", warehouseFullyQualifiedName, warehouseFullyQualifiedName, nil})
			if err = handleExternalChangesToObjectInShow(d,
				outputMappings...,
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, serviceSchema, append(extraSetStateToValuesFromConfigFields, "query_warehouse")); err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ServiceToSchema(service)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ServiceDetailsToSchema(serviceDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("compute_pool", service.ComputePool.FullyQualifiedName()),
			d.Set("external_access_integrations", collections.Map(service.ExternalAccessIntegrations, func(id sdk.AccountObjectIdentifier) string { return id.FullyQualifiedName() })),
			d.Set("comment", service.Comment),
			d.Set("service_type", service.Type()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
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

func ToJobServiceFromSpecificationRequest(value any) (sdk.JobServiceFromSpecificationRequest, error) {
	spec, err := ToServiceFromSpecificationRequest(value)
	if err != nil {
		return sdk.JobServiceFromSpecificationRequest{}, err
	}
	return sdk.JobServiceFromSpecificationRequest(spec), nil
}

func ToServiceFromSpecificationTemplateRequest(value any) (sdk.ServiceFromSpecificationTemplateRequest, error) {
	base, err := ToServiceFromSpecificationRequest(value)
	if err != nil {
		return sdk.ServiceFromSpecificationTemplateRequest{}, err
	}
	serviceFromSpecificationTemplate := sdk.ServiceFromSpecificationTemplateRequest{
		Location:                  base.Location,
		SpecificationTemplateFile: base.SpecificationFile,
		SpecificationTemplate:     base.Specification,
	}
	for _, v := range value.([]any) {
		fromSpecificationConfig := v.(map[string]any)
		if using := fromSpecificationConfig["using"].([]any); using != nil {
			serviceFromSpecificationTemplate.Using = collections.Map(using, func(v any) sdk.ListItem {
				m := v.(map[string]any)
				var item sdk.ListItem
				if value := m["key"].(string); value != "" {
					item.Key = value
				}
				if value := m["value"].(string); value != "" {
					item.Value = fmt.Sprintf("$$%s$$", value)
				}
				return item
			})
		}
	}
	return serviceFromSpecificationTemplate, nil
}

func ToJobServiceFromSpecificationTemplateRequest(value any) (sdk.JobServiceFromSpecificationTemplateRequest, error) {
	spec, err := ToServiceFromSpecificationTemplateRequest(value)
	if err != nil {
		return sdk.JobServiceFromSpecificationTemplateRequest{}, err
	}
	return sdk.JobServiceFromSpecificationTemplateRequest(spec), nil
}
