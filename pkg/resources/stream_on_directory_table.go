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

var streamOnDirectoryTableSchema = func() map[string]*schema.Schema {
	streamOnDirectoryTable := map[string]*schema.Schema{
		"stage": {
			Type:        schema.TypeString,
			Required:    true,
			Description: relatedResourceDescription(blocklistedCharactersFieldDescription("Specifies an identifier for the stage the stream will monitor. Due to Snowflake limitations, the provider can not read the stage's database and schema. For stages, Snowflake returns only partially qualified name instead of fully qualified name. Please use stages located in the same schema as the stream."), resources.Stage),
			// TODO (SNOW-1733130): the returned value is not a fully qualified name
			DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuotingPartiallyQualifiedName, IgnoreChangeToCurrentSnowflakeValueInShow("stage")),
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
	}
	return collections.MergeMaps(streamCommonSchema, streamOnDirectoryTable)
}()

func StreamOnDirectoryTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.StreamOnDirectoryTable, CreateStreamOnDirectoryTable(false)),
		ReadContext:   TrackingReadWrapper(resources.StreamOnDirectoryTable, ReadStreamOnDirectoryTable(true)),
		UpdateContext: TrackingUpdateWrapper(resources.StreamOnDirectoryTable, UpdateStreamOnDirectoryTable),
		DeleteContext: TrackingDeleteWrapper(resources.StreamOnDirectoryTable, DeleteStreamContext),
		Description:   "Resource used to manage streams on directory tables. For more information, check [stream documentation](https://docs.snowflake.com/en/sql-reference/sql/create-stream).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.StreamOnDirectoryTable, customdiff.All(
			ComputedIfAnyAttributeChanged(streamOnDirectoryTableSchema, ShowOutputAttributeName, "stage", "comment"),
			ComputedIfAnyAttributeChanged(streamOnDirectoryTableSchema, DescribeOutputAttributeName, "stage", "comment"),
			RecreateWhenStreamIsStale(),
			RecreateWhenStreamTypeChangedExternally(sdk.StreamSourceTypeStage),
		)),

		Schema: streamOnDirectoryTableSchema,

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.StreamOnDirectoryTable, ImportName[sdk.SchemaObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateStreamOnDirectoryTable(orReplace bool) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		stageIdRaw := d.Get("stage").(string)
		stageId, err := sdk.ParseSchemaObjectIdentifier(stageIdRaw)
		if err != nil {
			return diag.FromErr(err)
		}

		req := sdk.NewCreateOnDirectoryTableStreamRequest(id, stageId)

		errs := errors.Join(
			copyGrantsAttributeCreate(d, orReplace, &req.OrReplace, &req.CopyGrants),
			stringAttributeCreate(d, "comment", &req.Comment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		err = client.Streams.CreateOnDirectoryTable(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(id))

		return ReadStreamOnDirectoryTable(false)(ctx, d, meta)
	}
}

func ReadStreamOnDirectoryTable(withDirectoryChangesMarking bool) schema.ReadContextFunc {
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
		// TODO (SNOW-1733130): the returned value is not a fully qualified name
		if stream.TableName == nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Could not parse stage id",
					Detail:   fmt.Sprintf("stream name: %s", id.FullyQualifiedName()),
				},
			}
		}
		if err := d.Set("stage", *stream.TableName); err != nil {
			return diag.FromErr(err)
		}
		streamDescription, err := client.Streams.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := handleStreamRead(d, id, stream, streamDescription); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateStreamOnDirectoryTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// change on these fields can not be ForceNew because then the object is dropped explicitly and copying grants does not have effect
	// recreate when the stream is stale - see https://community.snowflake.com/s/article/using-tasks-to-avoid-stale-streams-when-incoming-data-is-empty
	if keys := changedKeys(d, "stage", "stale"); len(keys) > 0 {
		log.Printf("[DEBUG] Detected change on %q, recreating...", keys)
		return CreateStreamOnDirectoryTable(true)(ctx, d, meta)
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

	return ReadStreamOnDirectoryTable(false)(ctx, d, meta)
}
