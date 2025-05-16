package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantAccountRoleSchema = map[string]*schema.Schema{
	"role_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      relatedResourceDescription("The fully qualified name of the role which will be granted to the user or parent role.", resources.AccountRole),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
	},
	"user_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("The fully qualified name of the user on which specified role will be granted.", resources.User),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"user_name",
			"parent_role_name",
		},
	},
	"parent_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("The fully qualified name of the parent role which will create a parent-child relationship between the roles.", resources.AccountRole),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		ExactlyOneOf: []string{
			"user_name",
			"parent_role_name",
		},
	},
}

func GrantAccountRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.GrantAccountRole, CreateGrantAccountRole),
		ReadContext:   TrackingReadWrapper(resources.GrantAccountRole, ReadGrantAccountRole),
		DeleteContext: TrackingDeleteWrapper(resources.GrantAccountRole, DeleteGrantAccountRole),
		Schema:        grantAccountRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GrantAccountRole, func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), helpers.IDDelimiter)
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified: %v, expected <role_name>|<grantee_object_type>|<grantee_identifier>", d.Id())
				}
				if err := d.Set("role_name", strings.Trim(parts[0], "\"")); err != nil {
					return nil, err
				}
				switch parts[1] {
				case "ROLE":
					if err := d.Set("parent_role_name", strings.Trim(parts[2], "\"")); err != nil {
						return nil, err
					}
				case "USER":
					if err := d.Set("user_name", strings.Trim(parts[2], "\"")); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid object type specified: %v, expected ROLE or USER", parts[1])
				}

				return []*schema.ResourceData{d}, nil
			}),
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateGrantAccountRole implements schema.CreateFunc.
func CreateGrantAccountRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	roleName := d.Get("role_name").(string)
	roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)
	// format of snowflakeResourceID is <role_identifier>|<object type>|<target_identifier>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parentRoleName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantRoleRequest(roleIdentifier, sdk.GrantRole{
			Role: &parentRoleIdentifier,
		})
		if err := client.Roles.Grant(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	} else if userName, ok := d.GetOk("user_name"); ok && userName.(string) != "" {
		userIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(userName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeUser.String(), userIdentifier.FullyQualifiedName())
		req := sdk.NewGrantRoleRequest(roleIdentifier, sdk.GrantRole{
			User: &userIdentifier,
		})
		if err := client.Roles.Grant(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("invalid role grant specified: both parent_role_name and user_name are empty"))
	}
	d.SetId(snowflakeResourceID)
	return ReadGrantAccountRole(ctx, d, meta)
}

func ReadGrantAccountRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	if len(parts) != 3 {
		return diag.FromErr(fmt.Errorf("invalid ID specified: %v, expected <role_name>|<grantee_object_type>|<grantee_identifier>", d.Id()))
	}
	roleName := parts[0]
	roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)
	objectType := parts[1]
	targetIdentifier := parts[2]
	grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			Role: roleIdentifier,
		},
	})
	if err != nil {
		log.Printf("[DEBUG] role (%s) not found", roleIdentifier.FullyQualifiedName())
		d.SetId("")
		return nil
	}

	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == sdk.ObjectType(objectType) {
			if grant.GranteeName.FullyQualifiedName() == targetIdentifier {
				found = true
				break
			}
		}
	}
	if !found {
		log.Printf("[DEBUG] role grant (%s) not found", d.Id())
		d.SetId("")
	}

	return nil
}

func DeleteGrantAccountRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	if len(parts) != 3 {
		return diag.FromErr(fmt.Errorf("invalid ID specified: %v, expected <role_name>|<grantee_object_type>|<grantee_identifier>", d.Id()))
	}
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	objectType := parts[1]
	granteeName := parts[2]
	granteeIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName)
	switch objectType {
	case "ROLE":
		if err := client.Roles.Revoke(ctx, sdk.NewRevokeRoleRequest(id, sdk.RevokeRole{Role: &granteeIdentifier})); err != nil {
			return diag.FromErr(err)
		}
	case "USER":
		if err := client.Roles.Revoke(ctx, sdk.NewRevokeRoleRequest(id, sdk.RevokeRole{User: &granteeIdentifier})); err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.FromErr(fmt.Errorf("invalid object type specified: %v, expected ROLE or USER", objectType))
	}
	d.SetId("")
	return nil
}
