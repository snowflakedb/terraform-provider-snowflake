package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamOnTableSchema = func() map[string]*schema.Schema {
	streamOnTable := map[string]*schema.Schema{
		"table": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      relatedResourceDescription(blocklistedCharactersFieldDescription("Specifies an identifier for the table the stream will monitor."), resources.Table),
			DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("table_name")),
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		"append_only": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("mode", func(x any) any {
				return x.(string) == string(sdk.StreamModeAppendOnly)
			}),
			Description: booleanStringFieldDescription("Specifies whether this is an append-only stream."),
		},
		"show_initial_rows": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      externalChangesNotDetectedFieldDescription(booleanStringFieldDescription("Specifies whether to return all existing rows in the source table as row inserts the first time the stream is consumed.")),
		},
		AtAttributeName:     atSchema,
		BeforeAttributeName: beforeSchema,
	}
	return collections.MergeMaps(streamCommonSchema, streamOnTable)
}()

func StreamOnTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.StreamOnTable, CreateStreamOnTable(false)),
		ReadContext:   TrackingReadWrapper(resources.StreamOnTable, ReadStreamOnTable(true)),
		UpdateContext: TrackingUpdateWrapper(resources.StreamOnTable, UpdateStreamOnTable),
		DeleteContext: TrackingDeleteWrapper(resources.StreamOnTable, DeleteStreamContext),
		Description:   "Resource used to manage streams on tables. For more information, check [stream documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stream).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.StreamOnTable, customdiff.All(
			ComputedIfAnyAttributeChanged(streamOnTableSchema, ShowOutputAttributeName, "table", "append_only", "comment"),
			ComputedIfAnyAttributeChanged(streamOnTableSchema, DescribeOutputAttributeName, "table", "append_only", "comment"),
			RecreateWhenStreamIsStale(),
			RecreateWhenStreamTypeChangedExternally(sdk.StreamSourceTypeTable),
		)),

		Schema: streamOnTableSchema,

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.StreamOnTable, ImportStreamOnTable),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportStreamOnTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting stream import")
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	v, err := client.Streams.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}
	if err := d.Set("append_only", booleanStringFromBool(v.IsAppendOnly())); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateStreamOnTable(orReplace bool) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		tableIdRaw := d.Get("table").(string)
		tableId, err := sdk.ParseSchemaObjectIdentifier(tableIdRaw)
		if err != nil {
			return diag.FromErr(err)
		}

		req := sdk.NewCreateOnTableStreamRequest(id, tableId)

		errs := errors.Join(
			copyGrantsAttributeCreate(d, orReplace, &req.OrReplace, &req.CopyGrants),
			booleanStringAttributeCreate(d, "append_only", &req.AppendOnly),
			booleanStringAttributeCreate(d, "show_initial_rows", &req.ShowInitialRows),
			stringAttributeCreate(d, "comment", &req.Comment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		streamTimeTravelReq := handleStreamTimeTravel(d)
		if streamTimeTravelReq != nil {
			req.WithOn(*streamTimeTravelReq)
		}

		err = client.Streams.CreateOnTable(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(id))

		return ReadStreamOnTable(false)(ctx, d, meta)
	}
}

func ReadStreamOnTable(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		stream, err := client.Streams.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query stream. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Stream id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		tableId, err := sdk.ParseSchemaObjectIdentifier(*stream.TableName)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to parse table ID in Read.",
					Detail:   fmt.Sprintf("stream name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		if err := d.Set("table", tableId.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
		streamDescription, err := client.Streams.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := handleStreamRead(d, id, stream, streamDescription); err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			var mode sdk.StreamMode
			if stream.Mode != nil {
				mode = *stream.Mode
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"mode", "append_only", string(mode), booleanStringFromBool(stream.IsAppendOnly()), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, streamOnTableSchema, []string{
			"append_only",
		}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateStreamOnTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// change on these fields can not be ForceNew because then the object is dropped explicitly and copying grants does not have effect
	// recreate when the stream is stale - see https://community.snowflake.com/s/article/using-tasks-to-avoid-stale-streams-when-incoming-data-is-empty
	if keys := changedKeys(d, "table", "append_only", "at", "before", "show_initial_rows", "stale"); len(keys) > 0 {
		log.Printf("[DEBUG] Detected change on %q, recreating...", keys)
		return CreateStreamOnTable(true)(ctx, d, meta)
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if comment == "" {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithUnsetComment(true))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.Streams.Alter(ctx, sdk.NewAlterStreamRequest(id).WithSetComment(comment))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadStreamOnTable(false)(ctx, d, meta)
}
