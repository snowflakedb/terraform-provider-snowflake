package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sharesSchema = map[string]*schema.Schema{
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	"shares": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the shares available in the system.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Identifier for the share.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The comment on the share.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The owner of the share.",
				},
				"kind": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The kind of the share.",
				},
				"to": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "For the OUTBOUND share, list of consumers.",
					Elem:        schema.TypeString,
				},
			},
		},
	},
}

// Shares Snowflake Shares resource.
func Shares() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SharesDatasource), TrackingReadWrapper(datasources.Shares, ReadShares)),
		Schema:      sharesSchema,
	}
}

// ReadShares Reads the database metadata information.
func ReadShares(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	d.SetId("shares_read")
	pattern := d.Get("pattern").(string)
	var opts sdk.ShowShareOptions
	if pattern != "" {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(pattern),
		}
	}
	shares, err := client.Shares.Show(ctx, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	sharesFlatten := []map[string]interface{}{}
	for _, share := range shares {
		m := map[string]interface{}{}
		m["name"] = share.Name.Name()
		m["comment"] = share.Comment
		m["owner"] = share.Owner
		m["kind"] = share.Kind
		var to []string
		for _, consumer := range share.To {
			to = append(to, consumer.Name())
		}
		m["to"] = to
		sharesFlatten = append(sharesFlatten, m)
	}

	if err := d.Set("shares", sharesFlatten); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
