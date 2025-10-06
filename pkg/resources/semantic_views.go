package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var semanticViewsSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the semantic view; must be unique for the schema in which the semantic view is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the semantic view."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the semantic view."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"tables": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"table_alias": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies an alias for a logical table in the semantic view."),
				},
				"table_name": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      blocklistedCharactersFieldDescription("Specifies an identifier for the logical table."),
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"primary_key": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Definitions of primary keys in the logical table."),
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"unique": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Definitions of unique key combinations in the logical table."),
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"values": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Description: blocklistedCharactersFieldDescription("Unique key combinations in the logical table"),
							},
						},
					},
				},
				"synonym": {
					Type:        schema.TypeSet,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("List of synonyms for the logical table."),
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies a comment for the logical table.",
				},
			},
		},
	},
	"metrics": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO: update the SDK with the newly added fields for semantic expressions, then add them here
				// TODO: add PUBLIC/PRIVATE field
				// TODO: add table_alias
				// TODO: add fact_or_metric
				"semantic_expression": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies a semantic expression for a metric definition"),
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"qualified_expression_name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"sql_expression": {
								Type:     schema.TypeString,
								Required: true,
							},
							"synonym": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: blocklistedCharactersFieldDescription("List of synonyms for the metric definition."),
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"comment": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: blocklistedCharactersFieldDescription("Specifies a comment for the metric definition."),
							},
						},
					},
					ExactlyOneOf: []string{
						"semantic_expression",
						"window_function",
					},
				},
				"window_function": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies a window function for a metric definition"),
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"window_function": {
								Type:        schema.TypeString,
								Required:    true,
								Description: blocklistedCharactersFieldDescription("The window function for the metric definition"),
							},
							"metric": {
								Type:     schema.TypeString,
								Required: true,
							},
							"over_clause": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"partition_by": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: blocklistedCharactersFieldDescription("Specifies a partition by clause"),
										},
										"order_by": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: blocklistedCharactersFieldDescription("Specifies an order by clause"),
										},
										"window_frame_clause": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: blocklistedCharactersFieldDescription("Specifies a window frame clause"),
										},
									},
								},
							},
						},
					},
					ExactlyOneOf: []string{
						"semantic_expression",
						"window_function",
					},
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: blocklistedCharactersFieldDescription("Specifies a comment for the semantic view."),
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SEMANTIC VIEWS` for the given semantic view.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSemanticViewSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SEMANTIC VIEW` for the given semantic view.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeSemanticViewSchema,
		},
	},
}

func SemanticView() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.SemanticViews.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: CreateSemanticView,
		ReadContext:   ReadSemanticView,
		UpdateContext: UpdateSemanticView,
		DeleteContext: deleteFunc,
		Description:   "Resource used to manage semantic views. For more information, check [semantic views documentation](https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.SemanticView, customdiff.All(
			ComputedIfAnyAttributeChanged(semanticViewsSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(semanticViewsSchema, DescribeOutputAttributeName, "comment"),
		)),

		Schema: semanticViewsSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.SemanticView, ImportName[sdk.SchemaObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateSemanticView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)
	logicalTableRequests, err := getLogicalTableRequests(d.Get("tables").([]any))
	if err != nil {
		return diag.FromErr(err)
	}
	metricDefinitionRequests, err := getMetricDefinitionRequests(d.Get("metrics").([]any))
	if err != nil {
		return diag.FromErr(err)
	}
	semanticViewName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateSemanticViewRequest(semanticViewName, logicalTableRequests).
		WithSemanticViewMetrics(metricDefinitionRequests)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.SemanticViews.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(semanticViewName))
	return ReadSemanticView(ctx, d, meta)
}

func ReadSemanticView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	semanticView, err := client.SemanticViews.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query semantic view. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Semantic View id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	semanticViewDetails, err := client.SemanticViews.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	errs := errors.Join(
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.SemanticViewToSchema(semanticView)}),
		d.Set(DescribeOutputAttributeName, [][]map[string]any{schemas.SemanticViewDetailsToSchema(semanticViewDetails)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("comment", semanticView.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

func UpdateSemanticView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO [SNOW-2108211]: handle rename through ALTER ... RENAME to

	if d.HasChange("comment") {
		if comment := d.Get("comment").(string); comment != "" {
			if err := client.SemanticViews.Alter(ctx, sdk.NewAlterSemanticViewRequest(id).WithSetComment(comment)); err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		} else {
			if err := client.SemanticViews.Alter(ctx, sdk.NewAlterSemanticViewRequest(id).WithUnsetComment(true)); err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		}
	}
	return ReadSemanticView(ctx, d, meta)
}

func getLogicalTableRequest(from any) (*sdk.LogicalTableRequest, error) {
	c := from.(map[string]any)
	qualifiedTableName := c["table_name"].(string)

	logicalTableName, err := sdk.ParseSchemaObjectIdentifier(qualifiedTableName)
	if err != nil {
		return nil, err
	}
	logicalTableRequest := sdk.NewLogicalTableRequest(logicalTableName)

	if c["table_alias"] != nil && c["table_alias"].(string) != "" {
		aliasRequest := sdk.LogicalTableAliasRequest{LogicalTableAlias: c["table_alias"].(string)}
		logicalTableRequest = logicalTableRequest.WithLogicalTableAlias(aliasRequest)
	}

	if c["comment"] != nil && c["comment"].(string) != "" {
		logicalTableRequest = logicalTableRequest.WithComment(c["comment"].(string))
	}

	if c["primary_key"] != nil {
		primaryKeys, ok := c["primary_key"].([]any)
		if ok && len(primaryKeys) > 0 {
			var primaryKeyColumns []sdk.SemanticViewColumn
			for _, pk := range primaryKeys {
				primaryKeyColumns = append(primaryKeyColumns, sdk.SemanticViewColumn{Name: pk.(string)})
			}
			pkRequest := sdk.PrimaryKeysRequest{PrimaryKey: primaryKeyColumns}
			logicalTableRequest = logicalTableRequest.WithPrimaryKeys(pkRequest)
		}
	}

	if c["unique"] != nil {
		uniqueKeys, ok := c["unique"].([]any)
		if ok && len(uniqueKeys) > 0 {
			var ukRequests []sdk.UniqueKeysRequest
			for _, ukSet := range uniqueKeys {
				var uniqueKeyColumns []sdk.SemanticViewColumn
				for _, uk := range ukSet.([]any) {
					uniqueKeyColumns = append(uniqueKeyColumns, sdk.SemanticViewColumn{Name: uk.(string)})
				}
				ukRequest := sdk.UniqueKeysRequest{Unique: uniqueKeyColumns}
				ukRequests = append(ukRequests, ukRequest)
			}
			logicalTableRequest = logicalTableRequest.WithUniqueKeys(ukRequests)
		}
	}

	if c["synonym"] != nil {
		synonyms, ok := c["synonym"].([]any)
		if ok && len(synonyms) > 0 {
			var syns []sdk.Synonym
			for _, s := range synonyms {
				syns = append(syns, sdk.Synonym{Synonym: s.(string)})
			}
			sRequest := sdk.SynonymsRequest{WithSynonyms: syns}
			logicalTableRequest = logicalTableRequest.WithSynonyms(sRequest)
		}
	}

	return logicalTableRequest, nil
}

func getMetricDefinitionRequest(from any) (*sdk.MetricDefinitionRequest, error) {
	c := from.(map[string]any)
	metricDefinitionRequest := sdk.NewMetricDefinitionRequest()
	if c["semantic_expression"] != nil {
		semanticExpression := c["semantic_expression"].([]any)[0].(map[string]any)
		qualifiedExpNameRequest := sdk.NewQualifiedExpressionNameRequest().
			WithQualifiedExpressionName(semanticExpression["qualified_expression_name"].(string))
		sqlExpRequest := sdk.NewSemanticSqlExpressionRequest().
			WithSqlExpression(semanticExpression["sql_expression"].(string))
		semExpRequest := sdk.NewSemanticExpressionRequest(qualifiedExpNameRequest, sqlExpRequest)

		if semanticExpression["comment"] != nil && semanticExpression["comment"].(string) != "" {
			semExpRequest = semExpRequest.WithComment(semanticExpression["comment"].(string))
		}

		if semanticExpression["synonym"] != nil {
			synonyms, ok := semanticExpression["synonym"].([]any)
			if ok && len(synonyms) > 0 {
				var syns []sdk.Synonym
				for _, s := range synonyms {
					syns = append(syns, sdk.Synonym{Synonym: s.(string)})
				}
				sRequest := sdk.SynonymsRequest{WithSynonyms: syns}
				semExpRequest = semExpRequest.WithSynonyms(sRequest)
			}
		}
		return metricDefinitionRequest.WithSemanticExpression(*semExpRequest), nil
	} else if c["window_function"] != nil {
		windowFunctionDefinition := c["window_function"].([]any)[0].(map[string]any)
		windowFunction := windowFunctionDefinition["window_function"].(string)
		metric := windowFunctionDefinition["metric"].(string)
		windowFuncRequest := sdk.NewWindowFunctionMetricDefinitionRequest(windowFunction, metric)
		if windowFunctionDefinition["over_clauses"] != nil {
			overClause, ok := windowFunctionDefinition["over_clauses"].(map[string]any)
			if ok {
				overClauseRequest := sdk.NewWindowFunctionOverClauseRequest()
				if overClause["partition_by"] != nil {
					overClauseRequest = overClauseRequest.WithPartitionBy(overClause["partition_by"].(string))
				}
				if overClause["order_by"] != nil {
					overClauseRequest = overClauseRequest.WithOrderBy(overClause["order_by"].(string))
				}
				if overClause["window_frame_clause"] != nil {
					overClauseRequest = overClauseRequest.WithWindowFrameClause(overClause["window_frame_clause"].(string))
				}
				windowFuncRequest = windowFuncRequest.WithOverClause(*overClauseRequest)
			}
		}
		return metricDefinitionRequest.WithWindowFunctionMetricDefinition(*windowFuncRequest), nil
	} else {
		return nil, fmt.Errorf("either semantic expression or window function is required")
	}
}

func getLogicalTableRequests(from any) ([]sdk.LogicalTableRequest, error) {
	cols, ok := from.([]any)
	if !ok {
		return nil, fmt.Errorf("type assertion failure")
	}
	to := make([]sdk.LogicalTableRequest, len(cols))
	for i, c := range cols {
		cReq, err := getLogicalTableRequest(c)
		if err != nil {
			return nil, err
		}
		to[i] = *cReq
	}
	return to, nil
}

func getMetricDefinitionRequests(from any) ([]sdk.MetricDefinitionRequest, error) {
	cols, ok := from.([]any)
	if !ok {
		return nil, fmt.Errorf("type assertion failure")
	}
	to := make([]sdk.MetricDefinitionRequest, len(cols))
	for i, c := range cols {
		cReq, err := getMetricDefinitionRequest(c)
		if err != nil {
			return nil, err
		}
		to[i] = *cReq
	}
	return to, nil
}
