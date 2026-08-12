package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	return &schema.Resource{
		Schema: accountRoleSchema,

		CreateContext: TrackingCreateWrapper(resources.AccountRole, CreateAccountRole),
		ReadContext:   TrackingReadWrapper(resources.AccountRole, ReadAccountRole),
		DeleteContext: TrackingDeleteWrapper(resources.AccountRole, DeleteAccountRole),
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
	providerCtx := meta.(*provider.Context)
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountRole, err := showRoleCached(ctx, providerCtx, id)
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
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
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
		// The old identifier no longer resolves to anything after a successful rename; clear
		// any cached lookup for it so nothing in this plan/apply cycle observes a stale hit.
		invalidateRoleShowCache(providerCtx, id)

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

	// Any Alter above (rename or comment change) can make a previously cached lookup for this
	// role stale; invalidate after the mutating SQL has executed, before the trailing Read.
	invalidateRoleShowCache(providerCtx, id)

	return ReadAccountRole(ctx, d, meta)
}

// DeleteAccountRole mirrors ResourceDeleteContextFunc (pkg/resources/resource.go), which is a
// generic helper shared by several resources and therefore not the place to add role-cache
// invalidation. This is a bespoke copy, specific to snowflake_account_role, that additionally
// invalidates the role's cached lookup after a successful drop.
func DeleteAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := providerCtx.Client.Roles.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	invalidateRoleShowCache(providerCtx, id)

	d.SetId("")
	return nil
}

// showRoleCached looks up id via Roles.ShowByIDSafely, transparently caching the result in
// providerCtx.RoleShowCache when the ACCOUNT_ROLE_SHOW_CACHING experiment is enabled. Shared by
// snowflake_account_role, snowflake_grant_application_role, and
// snowflake_grant_privileges_to_account_role.
func showRoleCached(ctx context.Context, providerCtx *provider.Context, id sdk.AccountObjectIdentifier) (*sdk.Role, error) {
	if !experimentalfeatures.IsExperimentEnabled(experimentalfeatures.AccountRoleShowCaching, providerCtx.EnabledExperiments) {
		return providerCtx.Client.Roles.ShowByIDSafely(ctx, id)
	}
	return providerCtx.RoleShowCache.GetOrLoad(ctx, id.FullyQualifiedName(), func(loadCtx context.Context) (*sdk.Role, error) {
		return providerCtx.Client.Roles.ShowByIDSafely(loadCtx, id)
	})
}

// invalidateRoleShowCache invalidates the cached lookup for id, if the ACCOUNT_ROLE_SHOW_CACHING
// experiment is enabled. A no-op (not an error) if the cache was never populated for id — Invalidate
// on an absent key is already a documented no-op.
func invalidateRoleShowCache(providerCtx *provider.Context, id sdk.AccountObjectIdentifier) {
	if !experimentalfeatures.IsExperimentEnabled(experimentalfeatures.AccountRoleShowCaching, providerCtx.EnabledExperiments) {
		return
	}
	providerCtx.RoleShowCache.Invalidate(id.FullyQualifiedName())
}
