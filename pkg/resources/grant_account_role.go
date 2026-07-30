package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
			StateContext: TrackingImportWrapper(resources.GrantAccountRole, func(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
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
func CreateGrantAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	roleName := d.Get("role_name").(string)
	roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)

	safePublicRole := experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantAccountRoleSafePublicRole, providerCtx.EnabledExperiments) &&
		roleIdentifier.Name() == snowflakeroles.Public.Name()

	// format of snowflakeResourceID is <role_identifier>|<object type>|<target_identifier>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parentRoleName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeRole.String(), parentRoleIdentifier.FullyQualifiedName())
		if !safePublicRole {
			req := sdk.NewGrantRoleRequest(roleIdentifier, *sdk.NewGrantRoleToRequest().WithRole(parentRoleIdentifier))
			if err := client.Roles.Grant(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	} else if userName, ok := d.GetOk("user_name"); ok && userName.(string) != "" {
		userIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(userName.(string))
		snowflakeResourceID = helpers.EncodeSnowflakeID(roleIdentifier.FullyQualifiedName(), sdk.ObjectTypeUser.String(), userIdentifier.FullyQualifiedName())
		if !safePublicRole {
			req := sdk.NewGrantRoleRequest(roleIdentifier, *sdk.NewGrantRoleToRequest().WithUser(userIdentifier))
			if err := client.Roles.Grant(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		return diag.FromErr(fmt.Errorf("invalid role grant specified: both parent_role_name and user_name are empty"))
	}
	d.SetId(snowflakeResourceID)
	if safePublicRole {
		log.Printf("[DEBUG] skipping SHOW GRANTS for PUBLIC role grant (%s) — experiment %s enabled", snowflakeResourceID, experimentalfeatures.GrantAccountRoleSafePublicRole)
		return nil
	}
	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantAccountRoleShowCaching, providerCtx.EnabledExperiments) {
		// The trailing Read only re-confirms the grant we just created; this resource has no
		// computed/server-default fields to populate, so the Read is redundant here. With caching
		// enabled we skip it to avoid an extra SHOW GRANTS during apply. We still invalidate so any
		// later Read for this role in the same plan observes the new grant.
		providerCtx.GrantShowOfRoleCache.Invalidate(roleIdentifier.FullyQualifiedName())
		log.Printf("[DEBUG] skipping trailing SHOW GRANTS read after create (%s) — experiment %s enabled", snowflakeResourceID, experimentalfeatures.GrantAccountRoleShowCaching)
		return nil
	}
	return ReadGrantAccountRole(ctx, d, meta)
}

func ReadGrantAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	if len(parts) != 3 {
		return diag.FromErr(fmt.Errorf("invalid ID specified: %v, expected <role_name>|<grantee_object_type>|<grantee_identifier>", d.Id()))
	}
	roleName := parts[0]
	roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)

	// PUBLIC is always implicitly granted; SHOW GRANTS won't list it as an explicit grant.
	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantAccountRoleSafePublicRole, providerCtx.EnabledExperiments) &&
		roleIdentifier.Name() == snowflakeroles.Public.Name() {
		log.Printf("[DEBUG] skipping SHOW GRANTS for PUBLIC role grant (%s) — experiment %s enabled", d.Id(), experimentalfeatures.GrantAccountRoleSafePublicRole)
		return nil
	}

	objectType, err := sdk.ToObjectType(parts[1])
	if err != nil {
		return diag.FromErr(err)
	}
	targetIdentifier := parts[2]

	var grants []sdk.Grant
	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantAccountRoleShowCaching, providerCtx.EnabledExperiments) {
		cacheKey := roleIdentifier.FullyQualifiedName()
		grants, err = providerCtx.GrantShowOfRoleCache.GetOrLoad(cacheKey, func(loadCtx context.Context) ([]sdk.Grant, error) {
			return client.Grants.Show(loadCtx, &sdk.ShowGrantOptions{
				Of: &sdk.ShowGrantsOf{Role: roleIdentifier},
			})
		})
	} else {
		grants, err = client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{Role: roleIdentifier},
		})
	}
	if err != nil {
		log.Printf("[DEBUG] role (%s) not found", roleIdentifier.FullyQualifiedName())
		d.SetId("")
		return nil
	}

	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == objectType {
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

func DeleteGrantAccountRole(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	if len(parts) != 3 {
		return diag.FromErr(fmt.Errorf("invalid ID specified: %v, expected <role_name>|<grantee_object_type>|<grantee_identifier>", d.Id()))
	}
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(parts[0])
	objectType := parts[1]
	granteeName := parts[2]
	granteeIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(granteeName)

	// PUBLIC is always implicitly granted and cannot be explicitly revoked.
	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantAccountRoleSafePublicRole, providerCtx.EnabledExperiments) &&
		id.Name() == snowflakeroles.Public.Name() {
		log.Printf("[DEBUG] skipping REVOKE for PUBLIC role grant (%s) — experiment %s enabled", d.Id(), experimentalfeatures.GrantAccountRoleSafePublicRole)
		d.SetId("")
		return nil
	}

	revokeFunc := client.Roles.Revoke
	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantsSafeDestroy, providerCtx.EnabledExperiments) {
		revokeFunc = client.Roles.RevokeSafely
	}
	var err error
	switch objectType {
	case "ROLE":
		err = revokeFunc(ctx, sdk.NewRevokeRoleRequest(id, *sdk.NewRevokeRoleFromRequest().WithRole(granteeIdentifier)))
	case "USER":
		err = revokeFunc(ctx, sdk.NewRevokeRoleRequest(id, *sdk.NewRevokeRoleFromRequest().WithUser(granteeIdentifier)))
	default:
		return diag.FromErr(fmt.Errorf("invalid object type specified: %v, expected ROLE or USER", objectType))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.GrantAccountRoleShowCaching, providerCtx.EnabledExperiments) {
		providerCtx.GrantShowOfRoleCache.Invalidate(id.FullyQualifiedName())
	}
	d.SetId("")
	return nil
}
