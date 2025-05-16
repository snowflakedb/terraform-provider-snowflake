package resources

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantDatabaseRoleSchema = map[string]*schema.Schema{
	"database_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      relatedResourceDescription("The fully qualified name of the database role which will be granted to share or parent role.", resources.DatabaseRole),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"parent_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("The fully qualified name of the parent account role which will create a parent-child relationship between the roles.", resources.AccountRole),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"parent_role_name",
			"parent_database_role_name",
			"share_name",
		},
	},
	"parent_database_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("The fully qualified name of the parent database role which will create a parent-child relationship between the roles.", resources.DatabaseRole),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"parent_role_name",
			"parent_database_role_name",
			"share_name",
		},
	},
	"share_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("The fully qualified name of the share on which privileges will be granted.", resources.Share),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"parent_role_name",
			"parent_database_role_name",
			"share_name",
		},
	},
}

func GrantDatabaseRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.GrantDatabaseRole, CreateGrantDatabaseRole),
		ReadContext:   TrackingReadWrapper(resources.GrantDatabaseRole, ReadGrantDatabaseRole),
		DeleteContext: TrackingDeleteWrapper(resources.GrantDatabaseRole, DeleteGrantDatabaseRole),
		Schema:        grantDatabaseRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GrantDatabaseRole, func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := helpers.ParseResourceIdentifier(d.Id())
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified: %v, expected <database_role_name>|<object_type>|<target_identifier>", d.Id())
				}

				databaseRoleId, err := sdk.ParseDatabaseObjectIdentifier(parts[0])
				if err != nil {
					return nil, err
				}
				if err := d.Set("database_role_name", databaseRoleId.FullyQualifiedName()); err != nil {
					return nil, err
				}

				switch parts[1] {
				case "ROLE":
					accountRoleId, err := sdk.ParseAccountObjectIdentifier(parts[2])
					if err != nil {
						return nil, err
					}
					if err := d.Set("parent_role_name", accountRoleId.Name()); err != nil {
						return nil, err
					}
				case "DATABASE ROLE":
					parentDatabaseId, err := sdk.ParseDatabaseObjectIdentifier(parts[2])
					if err != nil {
						return nil, err
					}
					if err := d.Set("parent_database_role_name", parentDatabaseId.FullyQualifiedName()); err != nil {
						return nil, err
					}
				case "SHARE":
					shareId, err := sdk.ParseAccountObjectIdentifier(parts[2])
					if err != nil {
						return nil, err
					}
					if err := d.Set("share_name", shareId.Name()); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid object type specified: %v, expected ROLE, DATABASE ROLE, or SHARE", parts[1])
				}

				return []*schema.ResourceData{d}, nil
			}),
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateGrantDatabaseRole implements schema.CreateFunc.
func CreateGrantDatabaseRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseRoleName := d.Get("database_role_name").(string)
	databaseRoleIdentifier, err := sdk.ParseDatabaseObjectIdentifier(databaseRoleName)
	if err != nil {
		return diag.FromErr(err)
	}
	// format of snowflakeResourceID is <database_role_identifier>|<object type>|<parent_role_name>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier, err := sdk.ParseAccountObjectIdentifier(parentRoleName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		snowflakeResourceID = helpers.EncodeResourceIdentifier(databaseRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantDatabaseRoleRequest(databaseRoleIdentifier).WithAccountRole(parentRoleIdentifier)
		if err := client.DatabaseRoles.Grant(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	} else if parentDatabaseRoleName, ok := d.GetOk("parent_database_role_name"); ok && parentDatabaseRoleName.(string) != "" {
		parentRoleIdentifier, err := sdk.ParseDatabaseObjectIdentifier(parentDatabaseRoleName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		snowflakeResourceID = helpers.EncodeResourceIdentifier(databaseRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeDatabaseRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantDatabaseRoleRequest(databaseRoleIdentifier).WithDatabaseRole(parentRoleIdentifier)
		if err := client.DatabaseRoles.Grant(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	} else if shareName, ok := d.GetOk("share_name"); ok && shareName.(string) != "" {
		shareIdentifier, err := sdk.ParseAccountObjectIdentifier(shareName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		snowflakeResourceID = helpers.EncodeResourceIdentifier(databaseRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeShare.String(), shareIdentifier.FullyQualifiedName())
		req := sdk.NewGrantDatabaseRoleToShareRequest(databaseRoleIdentifier, shareIdentifier)
		if err := client.DatabaseRoles.GrantToShare(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(snowflakeResourceID)
	return ReadGrantDatabaseRole(ctx, d, meta)
}

// ReadGrantDatabaseRole implements schema.ReadFunc.
func ReadGrantDatabaseRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parts := helpers.ParseResourceIdentifier(d.Id())
	databaseRoleName := parts[0]
	databaseRoleIdentifier, err := sdk.ParseDatabaseObjectIdentifier(databaseRoleName)
	if err != nil {
		return diag.FromErr(err)
	}
	objectType := parts[1]
	targetIdentifier := parts[2]
	grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			DatabaseRole: databaseRoleIdentifier,
		},
	})
	if err != nil {
		log.Printf("[DEBUG] database role (%s) not found", databaseRoleIdentifier.FullyQualifiedName())
		d.SetId("")
		return nil
	}

	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == sdk.ObjectType(objectType) &&
			grant.GranteeName.FullyQualifiedName() == targetIdentifier {
			found = true
			break
		}
	}
	if !found {
		log.Printf("[DEBUG] database role grant (%s) not found", d.Id())
		d.SetId("")
	}

	return nil
}

// DeleteGrantDatabaseRole implements schema.DeleteFunc.
func DeleteGrantDatabaseRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parts := helpers.ParseResourceIdentifier(d.Id())
	id, err := sdk.ParseDatabaseObjectIdentifier(parts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	objectType := parts[1]
	granteeName := parts[2]
	switch objectType {
	case "ROLE":
		accountRoleId, err := sdk.ParseAccountObjectIdentifier(granteeName)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.DatabaseRoles.Revoke(ctx, sdk.NewRevokeDatabaseRoleRequest(id).WithAccountRole(accountRoleId)); err != nil {
			return diag.FromErr(err)
		}
	case "DATABASE ROLE":
		databaseRoleId, err := sdk.ParseDatabaseObjectIdentifier(granteeName)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.DatabaseRoles.Revoke(ctx, sdk.NewRevokeDatabaseRoleRequest(id).WithDatabaseRole(databaseRoleId)); err != nil {
			return diag.FromErr(err)
		}
	case "SHARE":
		sharedId, err := sdk.ParseAccountObjectIdentifier(granteeName)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.DatabaseRoles.RevokeFromShare(ctx, sdk.NewRevokeDatabaseRoleFromShareRequest(id, sharedId)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}
