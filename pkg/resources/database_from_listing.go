package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseFromListingSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the database; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"from_listing": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The global name of the listing from which to create the database. This can be an external listing or an organization listing. The listing must not be a paid listing, and listing terms must have been accepted if applicable. For more information, see [About sharing with listings](https://other-docs.snowflake.com/en/collaboration/collaboration-listings-about).",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func DatabaseFromListing() *schema.Resource {
	// TODO(SNOW-1818849): unassign network policies inside the database before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Databases.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.DatabaseFromListing, CreateDatabaseFromListing),
		UpdateContext: TrackingUpdateWrapper(resources.DatabaseFromListing, UpdateDatabaseFromListing),
		ReadContext:   TrackingReadWrapper(resources.DatabaseFromListing, ReadDatabaseFromListing),
		DeleteContext: TrackingDeleteWrapper(resources.DatabaseFromListing, deleteFunc),
		Description:   "A database created from a listing. Supports both external listings and organization listings. For more information about listings, see [About sharing with listings](https://other-docs.snowflake.com/en/collaboration/collaboration-listings-about).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.DatabaseFromListing, customdiff.All(
			ComputedIfAnyAttributeChanged(databaseFromListingSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: collections.MergeMaps(databaseFromListingSchema, sharedDatabaseParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.DatabaseFromListing, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateDatabaseFromListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	listingGlobalName := d.Get("from_listing").(string)

	opts := &sdk.CreateDatabaseFromListingOptions{}

	err = client.Databases.CreateFromListing(ctx, id, listingGlobalName, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	// Comment is set post-creation via ALTER DATABASE since CREATE DATABASE ... FROM LISTING
	// does not accept inline parameters.
	if v, ok := d.GetOk("comment"); ok {
		comment := v.(string)
		if len(comment) > 0 {
			err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					Comment: &comment,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadDatabaseFromListing(ctx, d, meta)
}

func UpdateDatabaseFromListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
				Set: &sdk.DatabaseSet{
					Comment: &comment,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
				Unset: &sdk.DatabaseUnset{
					Comment: sdk.Bool(true),
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadDatabaseFromListing(ctx, d, meta)
}

func ReadDatabaseFromListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	database, err := client.Databases.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query database from listing. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Database from listing id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if database.Origin != nil {
		if err := d.Set("from_listing", database.Origin.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("comment", database.Comment); err != nil {
		return diag.FromErr(err)
	}

	databaseParameters, err := client.Databases.ShowParameters(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := handleDatabaseParameterRead(d, databaseParameters); diags != nil {
		return diags
	}

	return nil
}
