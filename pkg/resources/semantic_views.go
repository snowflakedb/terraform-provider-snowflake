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
	"relationships": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"relationship_alias": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies an alias for the relationship."),
				},
				"table_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"table_alias": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"relationship_columns": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"referenced_table_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"referenced_table_alias": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"referenced_relationship_columns": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	"facts": {
		Type:     schema.TypeList,
		Optional: true,
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
					Description: blocklistedCharactersFieldDescription("List of synonyms for the fact."),
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies a comment for the fact."),
				},
			},
		},
	},
	"dimensions": {
		Type:     schema.TypeList,
		Optional: true,
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
					Description: blocklistedCharactersFieldDescription("List of synonyms for the dimension."),
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies a comment for the dimension."),
				},
			},
		},
	},
	"metrics": {
		Type:     schema.TypeList,
		Required: true,
		Description: blocklistedCharactersFieldDescription("Specify a list of metrics for the semantic view. " +
			"Each metric can have either a semantic expression or a window function in its definition."),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// TODO(SNOW-2396311): update the SDK with the newly added/updated fields for semantic expressions, then add them here
				// TODO(SNOW-2396371): add PUBLIC/PRIVATE field
				// TODO(SNOW-2398097): add table_alias
				// TODO(SNOW-2398097): add fact_or_metric
				"semantic_expression": {
					Type:     schema.TypeList,
					Optional: true,
					Description: blocklistedCharactersFieldDescription("Specifies a semantic expression for a metric definition." +
						"Cannot be used in combination with a window function."),
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"qualified_expression_name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: blocklistedCharactersFieldDescription("Specifies a name for the semantic expression"),
							},
							"sql_expression": {
								Type:        schema.TypeString,
								Required:    true,
								Description: blocklistedCharactersFieldDescription("Specifies the sql expression used to compute the metric."),
							},
							"synonym": {
								Type:        schema.TypeSet,
								Optional:    true,
								Description: blocklistedCharactersFieldDescription("List of synonyms for this semantic expression."),
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"comment": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: blocklistedCharactersFieldDescription("Specifies a comment for the semantic expression."),
							},
						},
					},
				},
				// TODO(SNOW-2396397): update the sdk and the model with the newly added/updated fields for window functions
				"window_function": {
					Type:     schema.TypeList,
					Optional: true,
					Description: blocklistedCharactersFieldDescription("Specifies a window function for a metric definition." +
						"Cannot be used in combination with a semantic expression."),
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"window_function": {
								Type:        schema.TypeString,
								Required:    true,
								Description: blocklistedCharactersFieldDescription("Specifies a name for the window function."),
							},
							"metric": {
								Type:        schema.TypeString,
								Required:    true,
								Description: blocklistedCharactersFieldDescription("Specifies a metric expression for this window function."),
							},
							"over_clause": {
								Type:     schema.TypeList,
								Required: true,
								MaxItems: 1,
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

	if d.Get("relationships") != nil {
		relationshipsRequests, err := getRelationshipRequests(d.Get("relationships").([]any))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithSemanticViewRelationships(relationshipsRequests)
	}
	if d.Get("facts") != nil {
		factsRequests, err := getSemanticExpressionRequests(d.Get("facts").([]any))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithSemanticViewFacts(factsRequests)
	}
	if d.Get("dimensions") != nil {
		dimensionsRequests, err := getSemanticExpressionRequests(d.Get("dimensions").([]any))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithSemanticViewDimensions(dimensionsRequests)
	}
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
		d.Set(DescribeOutputAttributeName, schemas.SemanticViewDetailsToSchema(semanticViewDetails)),
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
	if len(c["semantic_expression"].([]any)) > 0 {
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
			if synonyms, ok := semanticExpression["synonym"].(*schema.Set); ok && synonyms.Len() > 0 {
				var syns []sdk.Synonym
				for _, s := range synonyms.List() {
					syns = append(syns, sdk.Synonym{Synonym: s.(string)})
				}
				sRequest := sdk.SynonymsRequest{WithSynonyms: syns}
				semExpRequest = semExpRequest.WithSynonyms(sRequest)
			}
		}
		return metricDefinitionRequest.WithSemanticExpression(*semExpRequest), nil
	} else if len(c["window_function"].([]any)) > 0 {
		windowFunctionDefinition := c["window_function"].([]any)[0].(map[string]any)
		windowFunction := windowFunctionDefinition["window_function"].(string)
		metric := windowFunctionDefinition["metric"].(string)
		windowFuncRequest := sdk.NewWindowFunctionMetricDefinitionRequest(windowFunction, metric)
		if len(windowFunctionDefinition["over_clause"].([]any)) > 0 {
			overClause, ok := windowFunctionDefinition["over_clause"].([]any)[0].(map[string]any)
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

func getSemanticExpressionRequest(from any) (*sdk.SemanticExpressionRequest, error) {
	c := from.(map[string]any)
	qualifiedExpressionName := c["qualified_expression_name"].(string)
	if qualifiedExpressionName == "" {
		return nil, fmt.Errorf("qualified_expression_name is required")
	}
	qualifiedExpNameRequest := sdk.NewQualifiedExpressionNameRequest().
		WithQualifiedExpressionName(qualifiedExpressionName)

	sqlExpression := c["sql_expression"].(string)
	if sqlExpression == "" {
		return nil, fmt.Errorf("sql_expression is required")
	}
	sqlExpRequest := sdk.NewSemanticSqlExpressionRequest().
		WithSqlExpression(sqlExpression)
	semExpRequest := sdk.NewSemanticExpressionRequest(qualifiedExpNameRequest, sqlExpRequest)

	if c["comment"] != nil && c["comment"].(string) != "" {
		semExpRequest = semExpRequest.WithComment(c["comment"].(string))
	}

	if c["synonym"] != nil {
		synonyms, ok := c["synonym"].([]any)
		if ok && len(synonyms) > 0 {
			var syns []sdk.Synonym
			for _, s := range synonyms {
				syns = append(syns, sdk.Synonym{Synonym: s.(string)})
			}
			sRequest := sdk.SynonymsRequest{WithSynonyms: syns}
			semExpRequest = semExpRequest.WithSynonyms(sRequest)
		}
	}
	return semExpRequest, nil
}

func getRelationshipRequest(from any) (*sdk.SemanticViewRelationshipRequest, error) {
	c := from.(map[string]any)
	tableNameOrAliasRequest := sdk.NewRelationshipTableAliasRequest()
	if c["table_name"] != nil && c["table_name"].(string) != "" {
		tableName, err := sdk.ParseSchemaObjectIdentifier(c["table_name"].(string))
		if err != nil {
			return nil, err
		}
		tableNameOrAliasRequest.WithRelationshipTableName(tableName)
	}
	if c["table_alias"] != nil && c["table_alias"].(string) != "" {
		tableNameOrAliasRequest.WithRelationshipTableAlias(c["table_alias"].(string))
	}

	var relationshipColumnRequests []sdk.SemanticViewColumnRequest
	for _, relationshipColumn := range c["relationship_columns"].([]any) {
		relationshipColumnRequests = append(relationshipColumnRequests, *sdk.NewSemanticViewColumnRequest(relationshipColumn.(string)))
	}

	refTableNameOrAliasRequest := sdk.NewRelationshipTableAliasRequest()
	if c["referenced_table_name"] != nil && c["referenced_table_name"].(string) != "" {
		refTableName, err := sdk.ParseSchemaObjectIdentifier(c["referenced_table_name"].(string))
		if err != nil {
			return nil, err
		}
		refTableNameOrAliasRequest.WithRelationshipTableName(refTableName)
	}
	if c["referenced_table_alias"] != nil && c["referenced_table_alias"].(string) != "" {
		refTableNameOrAliasRequest.WithRelationshipTableAlias(c["referenced_table_alias"].(string))
	}

	request := sdk.NewSemanticViewRelationshipRequest(tableNameOrAliasRequest, relationshipColumnRequests, refTableNameOrAliasRequest)

	if c["relationship_alias"] != nil && c["relationship_alias"].(string) != "" {
		relAliasRequest := sdk.NewRelationshipAliasRequest().WithRelationshipAlias(c["relationship_alias"].(string))
		request.WithRelationshipAlias(*relAliasRequest)
	}

	if c["referenced_relationship_columns"] != nil {
		var refRelationshipColumnRequests []sdk.SemanticViewColumnRequest
		for _, refRelationshipColumn := range c["referenced_relationship_columns"].([]any) {
			refRelationshipColumnRequests = append(refRelationshipColumnRequests, *sdk.NewSemanticViewColumnRequest(refRelationshipColumn.(string)))
		}
		request.WithRelationshipRefColumnNames(refRelationshipColumnRequests)
	}

	return request, nil
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

func getSemanticExpressionRequests(from any) ([]sdk.SemanticExpressionRequest, error) {
	cols, ok := from.([]any)
	if !ok {
		return nil, fmt.Errorf("type assertion failure")
	}
	to := make([]sdk.SemanticExpressionRequest, len(cols))
	for i, c := range cols {
		cReq, err := getSemanticExpressionRequest(c)
		if err != nil {
			return nil, err
		}
		to[i] = *cReq
	}
	return to, nil
}

func getRelationshipRequests(from any) ([]sdk.SemanticViewRelationshipRequest, error) {
	cols, ok := from.([]any)
	if !ok {
		return nil, fmt.Errorf("type assertion failure")
	}
	to := make([]sdk.SemanticViewRelationshipRequest, len(cols))
	for i, c := range cols {
		cReq, err := getRelationshipRequest(c)
		if err != nil {
			return nil, err
		}
		to[i] = *cReq
	}
	return to, nil
}
