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
)

var listingSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the listing identifier (name). It must be unique within the organization, regardless of which Snowflake region the account is located in. Must start with an alphabetic character and cannot contain spaces or special characters except for underscores.",
	},
	"manifest": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_string": {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{"manifest.0.from_stage"},
				},
				"from_stage": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"manifest.0.from_string"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"stage": {
								Type:     schema.TypeList,
								Required: true,
							},
							"version_name": {
								Type:     schema.TypeList,
								Optional: true,
							},
							"location": {
								Type:     schema.TypeList,
								Optional: true,
								// If not passed, the manifest should be in the top-level
							},
						},
					},
				},
			},
		},
	},
	"share": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"application_package"},
	},
	"application_package": {
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"share"},
	},
	"review": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  BooleanDefault,
	},
	"publish": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  BooleanDefault,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW LISTINGS` for the given listing.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaSchema,
		},
	},
}

// TODO: Should we call it free_listing or public_listing?
func Listing() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Listings.DropSafely
		},
	)

	return &schema.Resource{
		Description:   "",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ListingResource), TrackingCreateWrapper(resources.Listing, CreateListing)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ListingResource), TrackingReadWrapper(resources.Listing, ReadListing)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ListingResource), TrackingUpdateWrapper(resources.Listing, UpdateListing)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ListingResource), TrackingDeleteWrapper(resources.Listing, deleteFunc)),

		Schema: listingSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Listing, schema.ImportStatePassthroughContext),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateListingRequest(id)

	if err := client.Listings.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadListing(ctx, d, meta)
}

func UpdateListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_ = client
	_ = id

	return ReadListing(ctx, d, meta)
}

func ReadListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	listing, err := client.Listings.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query listing. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Listing id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	_ = listing

	return nil
}
