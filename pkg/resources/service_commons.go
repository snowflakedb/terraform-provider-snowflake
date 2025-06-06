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
		// TODO(next PR): add service_type and acc tests for handling ext changes.
	}
	if allFieldsForceNew {
		for _, v := range schema {
			if v.Optional || v.Required {
				v.ForceNew = true
			}
		}
	}
	return schema
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
