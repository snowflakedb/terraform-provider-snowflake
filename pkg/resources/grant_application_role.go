package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var grantApplicationRoleSchema = map[string]*schema.Schema{
	"application_role_name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the application role to grant.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"parent_account_role_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      relatedResourceDescription("The fully qualified name of the account role on which application role will be granted.", resources.AccountRole),
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"parent_account_role_name",
			"application_name",
		},
	},
	"application_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The fully qualified name of the application on which application role will be granted.",
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		ExactlyOneOf: []string{
			"parent_account_role_name",
			"application_name",
		},
	},
}

func GrantApplicationRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.GrantApplicationRole, CreateContextGrantApplicationRole),
		ReadContext:   TrackingReadWrapper(resources.GrantApplicationRole, ReadContextGrantApplicationRole),
		DeleteContext: TrackingDeleteWrapper(resources.GrantApplicationRole, DeleteContextGrantApplicationRole),
		Schema:        grantApplicationRoleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GrantApplicationRole, func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := helpers.ParseResourceIdentifier(d.Id())
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified: %v, expected <application_role_name>|<object_type>|<target_identifier>", d.Id())
				}
				if err := d.Set("application_role_name", parts[0]); err != nil {
					return nil, err
				}
				switch parts[1] {
				case "ACCOUNT_ROLE":
					accountRoleId, err := sdk.ParseAccountObjectIdentifier(parts[2])
					if err != nil {
						return nil, err
					}
					if err := d.Set("parent_account_role_name", accountRoleId.FullyQualifiedName()); err != nil {
						return nil, err
					}
				case "APPLICATION":
					applicationId, err := sdk.ParseAccountObjectIdentifier(parts[2])
					if err != nil {
						return nil, err
					}
					if err := d.Set("application_name", applicationId.FullyQualifiedName()); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("invalid object type specified: %v, expected ACCOUNT_ROLE, APPLICATION", parts[1])
				}

				return []*schema.ResourceData{d}, nil
			}),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextGrantApplicationRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("application_role_name").(string)
	applicationRoleIdentifier, err := sdk.ParseDatabaseObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}
	// format of snowflakeResourceID is <application_role_identifier>|<object type>|<parent_account_role_name>
	var snowflakeResourceID string
	if parentRoleName, ok := d.GetOk("parent_account_role_name"); ok && parentRoleName.(string) != "" {
		parentRoleIdentifier, err := sdk.ParseAccountObjectIdentifier(parentRoleName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		snowflakeResourceID = helpers.EncodeResourceIdentifier(applicationRoleIdentifier.FullyQualifiedName(), "ACCOUNT_ROLE", parentRoleIdentifier.FullyQualifiedName())
		req := sdk.NewGrantApplicationRoleRequest(applicationRoleIdentifier).WithTo(*sdk.NewKindOfRoleRequest().WithRoleName(&parentRoleIdentifier))
		if err := client.ApplicationRoles.Grant(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	} else if applicationName, ok := d.GetOk("application_name"); ok && applicationName.(string) != "" {
		applicationIdentifier, err := sdk.ParseAccountObjectIdentifier(applicationName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		snowflakeResourceID = helpers.EncodeResourceIdentifier(applicationRoleIdentifier.FullyQualifiedName(), sdk.ObjectTypeApplication.String(), applicationIdentifier.FullyQualifiedName())
		req := sdk.NewGrantApplicationRoleRequest(applicationRoleIdentifier).WithTo(*sdk.NewKindOfRoleRequest().WithApplicationName(&applicationIdentifier))
		if err := client.ApplicationRoles.Grant(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(snowflakeResourceID)
	return ReadContextGrantApplicationRole(ctx, d, meta)
}

func ReadContextGrantApplicationRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	parts := strings.Split(d.Id(), helpers.IDDelimiter)
	applicationRoleName := parts[0]
	applicationRoleIdentifier, err := sdk.ParseDatabaseObjectIdentifier(applicationRoleName)
	if err != nil {
		return diag.FromErr(err)
	}
	objectTypeString := parts[1]
	if objectTypeString == "ACCOUNT_ROLE" {
		objectTypeString = "ROLE"
	}

	// first check if either the target account role or application exists
	targetIdentifier := parts[2]
	objectType := sdk.ObjectType(objectTypeString)
	switch objectType {
	case sdk.ObjectTypeRole:
		roleId, err := sdk.ParseAccountObjectIdentifier(targetIdentifier)
		if err != nil {
			return diag.FromErr(err)
		}
		if _, err := client.Roles.ShowByID(ctx, roleId); err != nil && errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve account role. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s", d.Id()),
				},
			}
		}
	case sdk.ObjectTypeApplication:
		applicationId, err := sdk.ParseAccountObjectIdentifier(targetIdentifier)
		if err != nil {
			return diag.FromErr(err)
		}
		if _, err := client.Applications.ShowByID(ctx, applicationId); err != nil && errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve application. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s", d.Id()),
				},
			}
		}
	}
	// then check if application role exists
	grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			ApplicationRole: applicationRoleIdentifier,
		},
	})
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve application role. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Id: %s", d.Id()),
				},
			}
		} else {
			return diag.FromErr(err)
		}
	}

	// finally check if the grant of the application role to the target (account role / application) exists
	var found bool
	for _, grant := range grants {
		if grant.GrantedTo == objectType {
			if grant.GrantedTo == sdk.ObjectTypeRole || grant.GrantedTo == sdk.ObjectTypeApplication {
				if grant.GranteeName.FullyQualifiedName() == targetIdentifier {
					found = true
					break
				}
			} else {
				// TODO(SNOW-1569535): fix when we'll have data types associated with the correct identifier parser
				/*
					note that grantee_name is not saved as a valid identifier in the
					SHOW GRANTS OF APPLICATION ROLE <application_role_name> command
					for example, "ABC"."test_parent_role" is saved as ABC."test_parent_role"
					or "ABC"."test_parent_role" is saved as ABC.test_parent_role
					and our internal mapper thereby fails to parse it correctly, returning "ABC."test_parent_role"
					so this funny string replacement is needed to make it work
				*/
				s := grant.GranteeName.FullyQualifiedName()
				if !strings.Contains(s, "\"") {
					parts := strings.Split(s, ".")
					s = sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]).FullyQualifiedName()
				} else {
					parts := strings.Split(s, "\".\"")
					if len(parts) < 2 {
						parts = strings.Split(s, "\".")
						if len(parts) < 2 {
							parts = strings.Split(s, ".\"")
							if len(parts) < 2 {
								s = strings.Trim(s, "\"")
								parts = strings.Split(s, ".")
							}
						}
					}
					s = sdk.NewDatabaseObjectIdentifier(parts[0], parts[1]).FullyQualifiedName()
				}
				if s == targetIdentifier {
					found = true
					break
				}
			}
		}
	}
	if !found {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Application role grant not found. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Id: %s", d.Id()),
			},
		}
	}

	return nil
}

func DeleteContextGrantApplicationRole(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parts := helpers.ParseResourceIdentifier(d.Id())
	id, err := sdk.ParseDatabaseObjectIdentifier(parts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	objectType := parts[1]
	granteeName := parts[2]
	switch objectType {
	case "ACCOUNT_ROLE":
		applicationRoleName, err := sdk.ParseAccountObjectIdentifier(granteeName)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.ApplicationRoles.Revoke(ctx, sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*sdk.NewKindOfRoleRequest().WithRoleName(&applicationRoleName))); err != nil {
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return diag.FromErr(err)
			}
		}
	case "APPLICATION":
		applicationName, err := sdk.ParseAccountObjectIdentifier(granteeName)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := client.ApplicationRoles.Revoke(ctx, sdk.NewRevokeApplicationRoleRequest(id).WithFrom(*sdk.NewKindOfRoleRequest().WithApplicationName(&applicationName))); err != nil {
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return diag.FromErr(err)
			}
		}
	}
	d.SetId("")
	return nil
}
