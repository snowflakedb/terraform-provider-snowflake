package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountRoleSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		// TODO(SNOW-1495079): Uncomment once better identifier validation will be implemented
		// ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Description: blocklistedCharactersFieldDescription("Identifier for the role; must be unique for your account."),
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW ROLES` for the given role.",
		Elem: &schema.Resource{
			Schema: schemas.ShowRoleSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func AccountRole() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] { return client.Roles.DropSafely },
	)

	return &schema.Resource{
		Schema: accountRoleSchema,

		CreateContext: TrackingCreateWrapper(resources.AccountRole, CreateAccountRole),
		ReadContext:   TrackingReadWrapper(resources.AccountRole, ReadAccountRole),
		DeleteContext: TrackingDeleteWrapper(resources.AccountRole, deleteFunc),
		UpdateContext: TrackingUpdateWrapper(resources.AccountRole, UpdateAccountRole),
		Description:   "The resource is used for role management, where roles can be assigned privileges and, in turn, granted to users and other roles. When granted to roles they can create hierarchies of privilege structures. For more details, refer to the [official documentation](https://docs.snowflake.com/en/user-guide/security-access-control-overview).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.AccountRole, customdiff.All(
			ComputedIfAnyAttributeChanged(accountRoleSchema, ShowOutputAttributeName, "comment", "name"),
			ComputedIfAnyAttributeChanged(accountRoleSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.AccountRole, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	req := sdk.NewCreateRoleRequest(id)

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	err = client.Roles.Create(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to create account role",
				Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.Name(), err),
			},
		}
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadAccountRole(ctx, d, meta)
}

func ReadAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountRole, err := client.Roles.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Account role not found; marking it as removed",
					Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.Name(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to show account role by id",
				Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.Name(), err),
			},
		}
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", accountRole.Comment); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to set account role comment",
				Detail:   fmt.Sprintf("Account role name: %s, comment: %s, err: %s", accountRole.ID().FullyQualifiedName(), accountRole.Comment, err),
			},
		}
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.RoleToSchema(accountRole)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

		if err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithRenameTo(newId)); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to rename account role name",
					Detail:   fmt.Sprintf("Previous account role name: %s, new account role name: %s, err: %s", id.Name(), newId.Name(), err),
				},
			}
		}

		id = newId
		d.SetId(helpers.EncodeResourceIdentifier(newId))
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			if err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithSetComment(v.(string))); err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to set account role comment",
						Detail:   fmt.Sprintf("Account role name: %s, comment: %s, err: %s", id.Name(), v, err),
					},
				}
			}
		} else {
			err := client.Roles.Alter(ctx, sdk.NewAlterRoleRequest(id).WithUnsetComment(true))
			if err != nil {
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to unset account role comment",
						Detail:   fmt.Sprintf("Account role name: %s, err: %s", id.Name(), err),
					},
				}
			}
		}
	}

	return ReadAccountRole(ctx, d, meta)
}
