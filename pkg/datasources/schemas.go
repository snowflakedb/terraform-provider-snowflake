package datasources

import (
	"context"

	resourceschemas "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemasSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SCHEMA for each schema returned by SHOW SCHEMAS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR SCHEMA for each schema returned by SHOW SCHEMAS. The output of describe is saved to the parameters field as a map. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"starts_with": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-sensitive** characters indicating the beginning of the object name.",
	},
	"limit": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Limits the number of rows returned. If the `limit.from` is set, then the limit will start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rows": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The maximum number of rows to return.",
				},
				"from": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.",
				},
			},
		},
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of streamlits",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account": {
					Type:         schema.TypeBool,
					Optional:     true,
					Description:  "Returns records for the entire account.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.application", "in.0.application_package"},
				},
				"database": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current database in use or for a specified database (db_name).",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.application", "in.0.application_package"},
				},
				"application": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the specified application.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.application", "in.0.application_package"},
				},
				"application_package": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the specified application package.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.application", "in.0.application_package"},
				},
			},
		},
	},
	"schemas": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all SCHEMA details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SCHEMAS.",
					Elem: &schema.Resource{
						Schema: resourceschemas.ShowSchemaSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SCHEMA.",
					Elem: &schema.Resource{
						Schema: resourceschemas.SchemaDescribeSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR SCHEMA.",
					Elem: &schema.Resource{
						Schema: resourceschemas.ShowSchemaParametersSchema,
					},
				},
			},
		},
	},
}

func Schemas() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Schemas, ReadSchemas),
		Schema:      schemasSchema,
		Description: "Data source used to get details of filtered schemas. Filtering is aligned with the current possibilities for [SHOW SCHEMAS](https://docs.snowflake.com/en/sql-reference/sql/show-schemas) query. The results of SHOW, DESCRIBE, and SHOW PARAMETERS IN are encapsulated in one output collection.",
	}
}

func ReadSchemas(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	var opts sdk.ShowSchemaOptions

	if likePattern, ok := d.GetOk("like"); ok {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		}
	}

	if startsWith, ok := d.GetOk("starts_with"); ok {
		opts.StartsWith = sdk.String(startsWith.(string))
	}

	if limit, ok := d.GetOk("limit"); ok && len(limit.([]any)) == 1 {
		limitMap := limit.([]any)[0].(map[string]any)

		rows := limitMap["rows"].(int)
		opts.LimitFrom = &sdk.LimitFrom{
			Rows: &rows,
		}

		if from, ok := limitMap["from"].(string); ok {
			opts.LimitFrom.From = &from
		}
	}

	if v, ok := d.GetOk("in"); ok {
		in := v.([]any)[0].(map[string]any)
		if v, ok := in["account"]; ok {
			if account := v.(bool); account {
				opts.In = &sdk.SchemaIn{Account: sdk.Bool(account)}
			}
		}
		if v, ok := in["database"]; ok {
			if database := v.(string); database != "" {
				opts.In = &sdk.SchemaIn{Name: sdk.NewAccountObjectIdentifier(database), Database: sdk.Pointer(true)}
			}
		}
		if v, ok := in["application"]; ok {
			if application := v.(string); application != "" {
				opts.In = &sdk.SchemaIn{Name: sdk.NewAccountObjectIdentifier(application), Application: sdk.Pointer(true)}
			}
		}
		if v, ok := in["application_package"]; ok {
			if applicationPackage := v.(string); applicationPackage != "" {
				opts.In = &sdk.SchemaIn{Name: sdk.NewAccountObjectIdentifier(applicationPackage), ApplicationPackage: sdk.Pointer(true)}
			}
		}
	}

	schemas, err := client.Schemas.Show(ctx, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("schemas_read")

	flattenedSchemas := make([]map[string]any, len(schemas))

	for i, schema := range schemas {
		schema := schema
		var schemaDescription []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.Schemas.Describe(ctx, schema.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			schemaDescription = resourceschemas.SchemaDescriptionToSchema(describeResult)
		}

		var schemaParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.Schemas.ShowParameters(ctx, schema.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			schemaParameters = []map[string]any{resourceschemas.SchemaParametersToSchema(parameters)}
		}

		flattenedSchemas[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{resourceschemas.SchemaToSchema(&schema)},
			resources.DescribeOutputAttributeName: schemaDescription,
			resources.ParametersAttributeName:     schemaParameters,
		}
	}

	err = d.Set("schemas", flattenedSchemas)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
