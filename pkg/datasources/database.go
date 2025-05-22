package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return its metadata.",
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"origin": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"options": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

// Database the Snowflake Database resource.
func Database() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.DatabaseDatasource), TrackingReadWrapper(datasources.Database, ReadDatabase)),
		Schema:      databaseSchema,
	}
}

// ReadDatabase read the database meta-data information.
func ReadDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	database, err := client.Databases.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(database.ID()))
	if err := d.Set("name", database.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", database.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner", database.Owner); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_default", database.IsDefault); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_current", database.IsCurrent); err != nil {
		return diag.FromErr(err)
	}
	var origin string
	if database.Origin != nil {
		origin = database.Origin.FullyQualifiedName()
	}
	if err := d.Set("origin", origin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("retention_time", database.RetentionTime); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_on", database.CreatedOn.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("options", database.Options); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
