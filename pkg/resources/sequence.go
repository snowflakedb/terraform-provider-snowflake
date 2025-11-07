package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sequenceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the name for the sequence.",
		ForceNew:    true,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Specifies a comment for the sequence.",
	},
	"increment": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     1,
		Description: "The amount the sequence will increase by each time it is used",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the sequence. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the sequence. Don't use the | character.",
		ForceNew:    true,
	},
	"next_value": {
		Type:        schema.TypeInt,
		Description: "The increment sequence interval.",
		Computed:    true,
		ForceNew:    true,
	},
	"ordering": {
		Type:        schema.TypeString,
		Description: "The ordering of the sequence. Either ORDER or NOORDER. Default is ORDER.",
		Optional:    true,
		Default:     "ORDER",
		ValidateDiagFunc: StringInSlice(
			[]string{
				string(sdk.ValuesBehaviorNoOrder),
				string(sdk.ValuesBehaviorOrder),
			}, false),
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Sequence() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Sequences.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.SequenceResource), TrackingCreateWrapper(resources.Sequence, CreateSequence)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.SequenceResource), TrackingReadWrapper(resources.Sequence, ReadSequence)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.SequenceResource), TrackingUpdateWrapper(resources.Sequence, UpdateSequence)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.SequenceResource), TrackingDeleteWrapper(resources.Sequence, deleteFunc)),

		Schema: sequenceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateSequence(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schema, name)
	req := sdk.NewCreateSequenceRequest(id)

	if v, ok := d.GetOk("increment"); ok {
		req.WithIncrement(v.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}
	if v, ok := d.GetOk("ordering"); ok {
		req.WithValuesBehavior(sdk.ValuesBehavior(v.(string)))
	}
	err := client.Sequences.Create(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeSnowflakeID(database, schema, name))

	return ReadSequence(ctx, d, meta)
}

func ReadSequence(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)
	seq, err := client.Sequences.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query sequence. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Sequence id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set("name", seq.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema", seq.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", seq.DatabaseName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", seq.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("increment", seq.Interval); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("next_value", seq.NextValue); err != nil {
		return diag.FromErr(err)
	}
	if seq.Ordered {
		if err := d.Set("ordering", "ORDER"); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("ordering", "NOORDER"); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateSequence(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("comment") {
		req := sdk.NewAlterSequenceRequest(id)
		req.WithSet(*sdk.NewSequenceSetRequest().WithComment(d.Get("comment").(string)))
		if err := client.Sequences.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("increment") {
		req := sdk.NewAlterSequenceRequest(id)
		req.WithSetIncrement(d.Get("increment").(int))
		if err := client.Sequences.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ordering") {
		req := sdk.NewAlterSequenceRequest(id)
		req.WithSet(*sdk.NewSequenceSetRequest().WithValuesBehavior(sdk.ValuesBehavior(d.Get("ordering").(string))))
		if err := client.Sequences.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadSequence(ctx, d, meta)
}
