package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sharedDatabaseSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the database; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"from_share": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: relatedResourceDescription("A fully qualified path to a share from which the database will be created. A fully qualified path follows the format of `\"<organization_name>\".\"<account_name>\".\"<share_name>\"`.", resources.Share),
		// TODO(SNOW-1495079): Add validation when ExternalObjectIdentifier will be available in IsValidIdentifierDescription:      "A fully qualified path to a share from which the database will be created. A fully qualified path follows the format of `\"<organization_name>\".\"<account_name>\".\"<share_name>\"`.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
	// TODO(SNOW-1843347): Add it as an item to discuss and either remove or uncomment (and implement) it
	// "is_transient": {
	//	Type:        schema.TypeBool,
	//	Optional:    true,
	//	ForceNew:    true,
	//	Description: "Specifies the database as transient. Transient databases do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	// },
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func SharedDatabase() *schema.Resource {
	// TODO(SNOW-1818849): unassign network policies inside the database before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Databases.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.SharedDatabase, CreateSharedDatabase),
		UpdateContext: TrackingUpdateWrapper(resources.SharedDatabase, UpdateSharedDatabase),
		ReadContext:   TrackingReadWrapper(resources.SharedDatabase, ReadSharedDatabase),
		DeleteContext: TrackingDeleteWrapper(resources.SharedDatabase, deleteFunc),
		Description:   "A shared database creates a database from a share provided by another Snowflake account. For more information about shares, see [Introduction to Secure Data Sharing](https://docs.snowflake.com/en/user-guide/data-sharing-intro).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.SharedDatabase, customdiff.All(
			ComputedIfAnyAttributeChanged(sharedDatabaseSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: collections.MergeMaps(sharedDatabaseSchema, sharedDatabaseParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.SharedDatabase, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	externalShareId, err := sdk.ParseExternalObjectIdentifier(d.Get("from_share").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &sdk.CreateSharedDatabaseOptions{
		// TODO(SNOW-1843347)
		// Transient:                  GetPropertyAsPointer[bool](d, "is_transient"),
		Comment: GetConfigPropertyAsPointerAllowingZeroValue[string](d, "comment"),
	}
	if parametersCreateDiags := handleSharedDatabaseParametersCreate(d, opts); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}

	err = client.Databases.CreateShared(ctx, id, externalShareId, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadSharedDatabase(ctx, d, meta)
}

func UpdateSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	return ReadSharedDatabase(ctx, d, meta)
}

func ReadSharedDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query shared database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Shared database id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if database.Origin != nil {
		if err := d.Set("from_share", database.Origin.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}
	}

	// TODO(SNOW-1843347)
	// if err := d.Set("is_transient", database.Transient); err != nil {
	//	return diag.FromErr(err)
	// }

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
