package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
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
					Type:        schema.TypeString,
					Required:    true,
					Description: blocklistedCharactersFieldDescription("Specifies an identifier for the logical table."),
				},
				"primary_key": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Definitions of primary keys in the logical table."),
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"keys": {
								Type: schema.TypeList,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: blocklistedCharactersFieldDescription("Columns to use in primary key"),
							},
						},
					},
				},
				"unique": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Definitions of unique key combinations in the logical table."),
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"keys": {
								Type: schema.TypeList,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: blocklistedCharactersFieldDescription("Unique key combinations in the logical table"),
							},
						},
					},
				},
				"synonym": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("List of synonyms for the logical table."),
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: blocklistedCharactersFieldDescription("Specifies a comment for the logical table."),
				},
			},
		},
	},
	"relationships": {},
	"facts":         {},
	"dimensions":    {},
	"metrics":       {},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: blocklistedCharactersFieldDescription("Specifies a comment for the semantic view."),
	},
	"copy_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: blocklistedCharactersFieldDescription("Specifies if the access privileges should be copied when OR REPLACE is used in creation of semantic view"),
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func SemanticView() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.SemanticViews.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.SemanticViewResource), TrackingCreateWrapper(resources.SemanticView, CreateSemanticView)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.SemanticViewResource), TrackingReadWrapper(resources.SemanticView, ReadGitRepository)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.SemanticViewResource), TrackingUpdateWrapper(resources.SemanticView, UpdateGitRepository)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.SemanticViewResource), TrackingDeleteWrapper(resources.SemanticView, deleteFunc)),
		Description:   "Resource used to manage git repositories. For more information, check [git repositories documentation](https://docs.snowflake.com/en/sql-reference/sql/create-git-repository).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.GitRepository, customdiff.All(
			ComputedIfAnyAttributeChanged(gitRepositorySchema, ShowOutputAttributeName, "origin", "api_integration", "git_credentials", "comment"),
		)),

		Schema: gitRepositorySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GitRepository, ImportName[sdk.SchemaObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateSemanticView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)
	logicalTableRequests, err := getLogicalTableRequests(d.Get("logical_tables").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	semanticViewName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateSemanticViewRequest(semanticViewName, logicalTableRequests)
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

func getLogicalTableRequest(from interface{}) (*sdk.LogicalTableRequest, error) {
	c := from.(map[string]interface{})
	name := c["name"].(string)
	databaseName := c["database"].(string)
	schemaName := c["schema"].(string)

	logicalTableName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	logicalTableRequest := sdk.NewLogicalTableRequest(logicalTableName)

	return logicalTableRequest.
		WithComment(c["comment"].(string)), nil
}

func getLogicalTableRequests(from interface{}) ([]sdk.LogicalTableRequest, error) {
	cols := from.([]interface{})
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
