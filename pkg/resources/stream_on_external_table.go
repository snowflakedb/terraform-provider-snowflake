package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamOnExternalTableSchema = func() map[string]*schema.Schema {
	streamOnExternalTable := map[string]*schema.Schema{
		"external_table": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      relatedResourceDescription(blocklistedCharactersFieldDescription("Specifies an identifier for the external table the stream will monitor."), resources.ExternalTable),
			DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("table_name")),
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		"insert_only": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("mode", func(x any) any {
				return x.(string) == string(sdk.StreamModeInsertOnly)
			}),
			Description: booleanStringFieldDescription("Specifies whether this is an insert-only stream."),
		},
		AtAttributeName:     atSchema,
		BeforeAttributeName: beforeSchema,
	}
	return collections.MergeMaps(streamCommonSchema, streamOnExternalTable)
}()

func StreamOnExternalTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.StreamOnExternalTable, CreateStreamOnExternalTable(false)),
		ReadContext:   TrackingReadWrapper(resources.StreamOnExternalTable, ReadStreamOnExternalTable(true)),
		UpdateContext: TrackingUpdateWrapper(resources.StreamOnExternalTable, UpdateStreamOnExternalTable),
		DeleteContext: TrackingDeleteWrapper(resources.StreamOnExternalTable, DeleteStreamContext),
		Description:   "Resource used to manage streams on external tables. For more information, check [stream documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stream).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.StreamOnExternalTable, customdiff.All(
			ComputedIfAnyAttributeChanged(streamOnExternalTableSchema, ShowOutputAttributeName, "external_table", "insert_only", "comment"),
			ComputedIfAnyAttributeChanged(streamOnExternalTableSchema, DescribeOutputAttributeName, "external_table", "insert_only", "comment"),
			RecreateWhenStreamIsStale(),
			RecreateWhenStreamTypeChangedExternally(sdk.StreamSourceTypeExternalTable),
		)),

		Schema: streamOnExternalTableSchema,

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.StreamOnExternalTable, ImportStreamOnExternalTable),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportStreamOnExternalTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
	if err := d.Set("insert_only", booleanStringFromBool(v.IsInsertOnly())); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateStreamOnExternalTable(orReplace bool) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		externalTableIdRaw := d.Get("external_table").(string)
		externalTableId, err := sdk.ParseSchemaObjectIdentifier(externalTableIdRaw)
		if err != nil {
			return diag.FromErr(err)
		}

		req := sdk.NewCreateOnExternalTableStreamRequest(id, externalTableId)

		errs := errors.Join(
			copyGrantsAttributeCreate(d, orReplace, &req.OrReplace, &req.CopyGrants),
			booleanStringAttributeCreate(d, "insert_only", &req.InsertOnly),
			stringAttributeCreate(d, "comment", &req.Comment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		streamTimeTravelReq := handleStreamTimeTravel(d)
		if streamTimeTravelReq != nil {
			req.WithOn(*streamTimeTravelReq)
		}

		err = client.Streams.CreateOnExternalTable(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(id))

		return ReadStreamOnExternalTable(false)(ctx, d, meta)
	}
}

func ReadStreamOnExternalTable(withExternalChangesMarking bool) schema.ReadContextFunc {
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

		externalTableId, err := sdk.ParseSchemaObjectIdentifier(*stream.TableName)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to parse external table ID in Read.",
					Detail:   fmt.Sprintf("stream name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		if err := d.Set("external_table", externalTableId.FullyQualifiedName()); err != nil {
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
				outputMapping{"mode", "insert_only", string(mode), booleanStringFromBool(stream.IsInsertOnly()), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, streamOnExternalTableSchema, []string{
			"insert_only",
		}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateStreamOnExternalTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// change on these fields can not be ForceNew because then the object is dropped explicitly and copying grants does not have effect
	// recreate when the stream is stale - see https://community.snowflake.com/s/article/using-tasks-to-avoid-stale-streams-when-incoming-data-is-empty
	if keys := changedKeys(d, "external_table", "insert_only", "at", "before", "stale"); len(keys) > 0 {
		log.Printf("[DEBUG] Detected change on %q, recreating...", keys)
		return CreateStreamOnExternalTable(true)(ctx, d, meta)
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

	return ReadStreamOnExternalTable(false)(ctx, d, meta)
}
