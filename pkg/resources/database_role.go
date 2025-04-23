package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the database role."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the database role."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database role.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW DATABASE ROLES` for the given database role. Note that this value will be only recomputed whenever comment field changes.",
		Elem: &schema.Resource{
			Schema: schemas.ShowDatabaseRoleSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func DatabaseRole() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseDatabaseObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.DatabaseObjectIdentifier] {
			return client.DatabaseRoles.DropSafely
		},
	)

	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: TrackingCreateWrapper(resources.DatabaseRole, CreateDatabaseRole),
		ReadContext:   TrackingReadWrapper(resources.DatabaseRole, ReadDatabaseRole),
		UpdateContext: TrackingUpdateWrapper(resources.DatabaseRole, UpdateDatabaseRole),
		DeleteContext: TrackingDeleteWrapper(resources.DatabaseRole, deleteFunc),

		Description: "Resource used to manage database roles. For more information, check [database roles documentation](https://docs.snowflake.com/en/sql-reference/sql/create-database-role).",

		Schema: databaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.DatabaseRole, ImportName[sdk.DatabaseObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.DatabaseRole, customdiff.All(
			ComputedIfAnyAttributeChanged(databaseRoleSchema, ShowOutputAttributeName, "comment", "name"),
			ComputedIfAnyAttributeChanged(databaseRoleSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ReadDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	databaseRole, err := client.DatabaseRoles.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Database role not found; marking it as removed",
					Detail:   fmt.Sprintf("Database role id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set("comment", databaseRole.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.DatabaseRoleToSchema(databaseRole)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func CreateDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	roleName := d.Get("name").(string)
	id := sdk.NewDatabaseObjectIdentifier(databaseName, roleName)
	createRequest := sdk.NewCreateDatabaseRoleRequest(id)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(v.(string))
	}

	err := client.DatabaseRoles.Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadDatabaseRole(ctx, d, meta)
}

func UpdateDatabaseRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), d.Get("name").(string))

		err = client.DatabaseRoles.Alter(ctx, sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("comment") {
		newComment := d.Get("comment").(string)
		err := client.DatabaseRoles.Alter(ctx, sdk.NewAlterDatabaseRoleRequest(id).WithSet(*sdk.NewDatabaseRoleSetRequest(newComment)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadDatabaseRole(ctx, d, meta)
}
