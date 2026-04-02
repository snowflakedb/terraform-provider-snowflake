package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var listingsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC LISTING for each listing returned by SHOW LISTINGS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"listings": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all listings details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW LISTINGS.",
					Elem:        &schema.Resource{Schema: schemas.ShowListingSchema},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE LISTING.",
					Elem:        &schema.Resource{Schema: schemas.DescribeListingSchema},
				},
			},
		},
	},
}

func Listings() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ListingsDatasource), TrackingReadWrapper(datasources.Listings, ReadListings)),
		Schema:      listingsSchema,
		Description: "Data source used to get details of filtered listings. Filtering is aligned with the current possibilities for SHOW LISTINGS query (`like`, `starts_with`, and `limit` are supported). The results of SHOW and DESCRIBE are encapsulated in one output collection.",
	}
}

func ReadListings(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := new(sdk.ShowListingRequest)

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	listings, err := client.Listings.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("listings_read")

	flattened := make([]map[string]any, len(listings))
	for i, listing := range listings {
		listing := listing
		var describe []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.Listings.Describe(ctx, sdk.NewDescribeListingRequest(listing.ID()))
			if err != nil {
				return diag.FromErr(err)
			}
			describe = []map[string]any{schemas.ListingDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ListingToSchema(&listing)},
			resources.DescribeOutputAttributeName: describe,
		}
	}
	if err := d.Set("listings", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
