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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var mcpServerSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the MCP server."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the MCP server."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the MCP server."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"specification": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      joinWithSpace("Specifies a YAML object containing the MCP server tool definitions.", doubleDollarQuotesDescription()),
		DiffSuppressFunc: NormalizeAndCompare(sdk.NormalizeMcpServerSpecification),
		ValidateDiagFunc: validation.AllDiag(validation.ToDiagFunc(validation.StringIsNotEmpty), forbidDoubleDollarQuotes),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the MCP server.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW MCP SERVERS` for this MCP server.",
		Elem: &schema.Resource{
			Schema: schemas.ShowMcpServerSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE MCP SERVER` for this MCP server.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeMcpServerDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func McpServer() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.McpServers.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.McpServerResource), TrackingCreateWrapper(resources.McpServer, CreateMcpServer)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.McpServerResource), TrackingReadWrapper(resources.McpServer, ReadMcpServer)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.McpServerResource), TrackingUpdateWrapper(resources.McpServer, UpdateMcpServer)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.McpServerResource), TrackingDeleteWrapper(resources.McpServer, deleteFunc)),
		Description:   "Resource used to manage MCP server objects. For more information, check [MCP server documentation](https://docs.snowflake.com/en/sql-reference/sql/create-mcp-server).",

		Schema: mcpServerSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.McpServer, ImportName[sdk.SchemaObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateMcpServer(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	spec := d.Get("specification").(string)

	request := sdk.NewCreateMcpServerRequest(id, spec)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.McpServers.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating MCP server: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadMcpServer(ctx, d, meta)
}

func UpdateMcpServer(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return nil
}

func ReadMcpServer(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mcpServer, err := client.McpServers.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "Failed to query MCP server. Marking the resource as removed.",
					Detail:   fmt.Sprintf("MCP server id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.McpServers.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	errs := errors.Join(
		d.Set("specification", details.ServerSpec),
		d.Set("comment", details.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.McpServerToSchema(mcpServer)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.McpServerDetailsToSchema(details)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
	)
	return diag.FromErr(errs)
}
