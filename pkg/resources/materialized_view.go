package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var materializedViewSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the view; must be unique for the schema in which the view is created.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the view. Don't use the | character.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the view. Don't use the | character.",
		ForceNew:    true,
	},
	"warehouse": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The warehouse name.",
		ForceNew:    true,
	},
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Overwrites the View if it exists.",
	},
	"is_secure": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies that the view is secure.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the view.",
	},
	"statement": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the query used to create the view.",
		ForceNew:         true,
		DiffSuppressFunc: DiffSuppressStatement,
	},
	"tag":                           tagReferenceSchema,
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// MaterializedView returns a pointer to the resource representing a view.
func MaterializedView() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.MaterializedViews.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.MaterializedViewResource), TrackingCreateWrapper(resources.MaterializedView, CreateMaterializedView)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.MaterializedViewResource), TrackingReadWrapper(resources.MaterializedView, ReadMaterializedView)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.MaterializedViewResource), TrackingUpdateWrapper(resources.MaterializedView, UpdateMaterializedView)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.MaterializedViewResource), TrackingDeleteWrapper(resources.MaterializedView, deleteFunc)),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.MaterializedView, customdiff.All(
			ComputedIfAnyAttributeChanged(materializedViewSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: materializedViewSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateMaterializedView implements schema.CreateFunc.
func CreateMaterializedView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	s := d.Get("statement").(string)
	createRequest := sdk.NewCreateMaterializedViewRequest(id, s)

	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		createRequest.WithOrReplace(sdk.Bool(true))
	}

	if v, ok := d.GetOk("is_secure"); ok && v.(bool) {
		createRequest.WithSecure(sdk.Bool(true))
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	warehouseName := d.Get("warehouse").(string)
	// TODO [SNOW-1348355]: this was the old implementation, it's left for now, we will address this with resources rework discussions
	err := client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(warehouseName))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting warehouse %s while creating materialized view %v err = %w", warehouseName, name, err))
	}

	err = client.MaterializedViews.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating materialized view %v err = %w", name, err))
	}

	// TODO [SNOW-1348355]: we have to set tags after creation because existing materialized view extractor is not aware of TAG during CREATE
	// Will be discussed with parser topic during resources redesign.
	if _, ok := d.GetOk("tag"); ok {
		err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(getPropertyTags(d, "tag")))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags on materialized view %v, err = %w", id, err))
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadMaterializedView(ctx, d, meta)
}

// ReadMaterializedView implements schema.ReadFunc.
func ReadMaterializedView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	materializedView, err := client.MaterializedViews.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query materialized view. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Materialized view id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", materializedView.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_secure", materializedView.IsSecure); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", materializedView.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema", materializedView.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", materializedView.DatabaseName); err != nil {
		return diag.FromErr(err)
	}

	// Want to only capture the SELECT part of the query because before that is the CREATE part of the view.
	extractor := snowflake.NewViewSelectStatementExtractor(materializedView.Text)
	substringOfQuery, err := extractor.ExtractMaterializedView()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("statement", substringOfQuery); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// UpdateMaterializedView implements schema.UpdateFunc.
func UpdateMaterializedView(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.MaterializedViews.Alter(ctx, sdk.NewAlterMaterializedViewRequest(id).WithRenameTo(&newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming materialized view %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	var runSetStatement bool
	var runUnsetStatement bool
	setRequest := sdk.NewMaterializedViewSetRequest()
	unsetRequest := sdk.NewMaterializedViewUnsetRequest()

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if comment == "" {
			runUnsetStatement = true
			unsetRequest.WithComment(sdk.Bool(true))
		} else {
			runSetStatement = true
			setRequest.WithComment(sdk.String(comment))
		}
	}
	if d.HasChange("is_secure") {
		if d.Get("is_secure").(bool) {
			runSetStatement = true
			setRequest.WithSecure(sdk.Bool(true))
		} else {
			runUnsetStatement = true
			unsetRequest.WithSecure(sdk.Bool(true))
		}
	}

	if runSetStatement {
		err := client.MaterializedViews.Alter(ctx, sdk.NewAlterMaterializedViewRequest(id).WithSet(setRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating materialized view: %w", err))
		}
	}

	if runUnsetStatement {
		err := client.MaterializedViews.Alter(ctx, sdk.NewAlterMaterializedViewRequest(id).WithUnset(unsetRequest))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating materialized view: %w", err))
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			// TODO [SNOW-1022645]: view is used on purpose here; change after we have an agreement on situations like this in the SDK
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting tags on %v, err = %w", d.Id(), err))
			}
		}

		if len(setTags) > 0 {
			// TODO [SNOW-1022645]: view is used on purpose here; change after we have an agreement on situations like this in the SDK
			err := client.Views.Alter(ctx, sdk.NewAlterViewRequest(id).WithSetTags(setTags))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err))
			}
		}
	}

	return ReadMaterializedView(ctx, d, meta)
}
