package acceptance

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func ComposeCheckDestroy(t *testing.T, resources ...resources.Resource) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		errs := make([]error, 0)
		for _, resource := range resources {
			checkFunc := CheckDestroy(t, resource)
			errs = append(errs, checkFunc(s))
		}
		return errors.Join(errs...)
	}
}

func CheckDestroy(t *testing.T, resource resources.Resource) func(*terraform.State) error {
	t.Helper()
	// TODO [SNOW-1653619]: use TestClient() here
	client := atc.client
	t.Logf("running check destroy for resource %s", resource)

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resource.String() {
				continue
			}
			t.Logf("found resource %s in state", resource)
			ctx := context.Background()
			id, err := decodeSnowflakeId(rs, resource)
			if err != nil {
				return err
			}
			if id == nil {
				return fmt.Errorf("could not get the id of %s", resource)
			}
			showById, ok := showByIdFunctions[resource]
			if !ok {
				return fmt.Errorf("unsupported show by id in cleanup for %s, with id %v", resource, id.FullyQualifiedName())
			}
			if showById(ctx, client, id) == nil {
				return fmt.Errorf("%s %v still exists", resource, id.FullyQualifiedName())
			} else {
				t.Logf("resource %s (%v) was dropped successfully in Snowflake", resource, id.FullyQualifiedName())
			}
		}
		return nil
	}
}

func decodeSnowflakeId(rs *terraform.ResourceState, resource resources.Resource) (sdk.ObjectIdentifier, error) {
	switch resource {
	case resources.ExternalFunction:
		return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.ID), nil
	case resources.FunctionJava,
		resources.FunctionJavascript,
		resources.FunctionPython,
		resources.FunctionScala,
		resources.FunctionSql:
		return sdk.ParseSchemaObjectIdentifierWithArguments(rs.Primary.ID)
	case resources.ProcedureJava,
		resources.ProcedureJavascript,
		resources.ProcedurePython,
		resources.ProcedureScala,
		resources.ProcedureSql:
		return sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(rs.Primary.ID), nil
	default:
		return helpers.DecodeSnowflakeID(rs.Primary.ID), nil
	}
}

type showByIdFunc func(context.Context, *sdk.Client, sdk.ObjectIdentifier) error

var showByIdFunctions = map[resources.Resource]showByIdFunc{
	resources.Account: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Accounts.ShowByID)
	},
	resources.AccountRole: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Roles.ShowByID)
	},
	resources.Alert: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Alerts.ShowByID)
	},
	resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.ApiAuthenticationIntegrationWithClientCredentials: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.ApiAuthenticationIntegrationWithJwtBearer: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.ApiIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ApiIntegrations.ShowByID)
	},
	resources.AuthenticationPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.AuthenticationPolicies.ShowByID)
	},
	resources.PrimaryConnection: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Connections.ShowByID)
	},
	resources.CortexSearchService: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.CortexSearchServices.ShowByID)
	},
	resources.Database: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Databases.ShowByID)
	},
	resources.DatabaseRole: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.DatabaseRoles.ShowByID)
	},
	resources.DynamicTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.DynamicTables.ShowByID)
	},
	resources.EmailNotificationIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NotificationIntegrations.ShowByID)
	},
	resources.ExternalFunction: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ExternalFunctions.ShowByID)
	},
	resources.ExternalOauthSecurityIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.ExternalTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ExternalTables.ShowByID)
	},
	resources.ExternalVolume: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ExternalVolumes.ShowByID)
	},
	resources.FailoverGroup: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.FailoverGroups.ShowByID)
	},
	resources.FileFormat: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.FileFormats.ShowByID)
	},
	resources.FunctionJava: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Functions.ShowByID)
	},
	resources.FunctionJavascript: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Functions.ShowByID)
	},
	resources.FunctionPython: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Functions.ShowByID)
	},
	resources.FunctionScala: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Functions.ShowByID)
	},
	resources.FunctionSql: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Functions.ShowByID)
	},
	resources.ImageRepository: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ImageRepositories.ShowByID)
	},
	resources.LegacyServiceUser: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Users.ShowByID)
	},
	resources.ManagedAccount: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ManagedAccounts.ShowByID)
	},
	resources.MaskingPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.MaskingPolicies.ShowByID)
	},
	resources.MaterializedView: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.MaterializedViews.ShowByID)
	},
	resources.NetworkPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NetworkPolicies.ShowByID)
	},
	resources.NetworkRule: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NetworkRules.ShowByID)
	},
	resources.NotificationIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.NotificationIntegrations.ShowByID)
	},
	resources.OauthIntegrationForCustomClients: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.OauthIntegrationForPartnerApplications: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.PasswordPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.PasswordPolicies.ShowByID)
	},
	resources.Pipe: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Pipes.ShowByID)
	},
	resources.ProcedureJava: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Procedures.ShowByID)
	},
	resources.ProcedureJavascript: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Procedures.ShowByID)
	},
	resources.ProcedurePython: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Procedures.ShowByID)
	},
	resources.ProcedureScala: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Procedures.ShowByID)
	},
	resources.ProcedureSql: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Procedures.ShowByID)
	},
	resources.ResourceMonitor: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.ResourceMonitors.ShowByID)
	},
	resources.RowAccessPolicy: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.RowAccessPolicies.ShowByID)
	},
	resources.Saml2SecurityIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.Schema: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Schemas.ShowByID)
	},
	resources.ScimSecurityIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.SecurityIntegrations.ShowByID)
	},
	resources.SecondaryConnection: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Connections.ShowByID)
	},
	resources.SecondaryDatabase: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Databases.ShowByID)
	},
	resources.SecretWithAuthorizationCodeGrant: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Secrets.ShowByID)
	},
	resources.SecretWithBasicAuthentication: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Secrets.ShowByID)
	},
	resources.SecretWithClientCredentials: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Secrets.ShowByID)
	},
	resources.SecretWithGenericString: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Secrets.ShowByID)
	},
	resources.Sequence: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Sequences.ShowByID)
	},
	resources.ServiceUser: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Users.ShowByID)
	},
	resources.Share: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Shares.ShowByID)
	},
	resources.SharedDatabase: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Databases.ShowByID)
	},
	resources.Stage: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Stages.ShowByID)
	},
	resources.StorageIntegration: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.StorageIntegrations.ShowByID)
	},
	resources.StreamOnDirectoryTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Streams.ShowByID)
	},
	resources.StreamOnExternalTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Streams.ShowByID)
	},
	resources.StreamOnTable: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Streams.ShowByID)
	},
	resources.StreamOnView: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Streams.ShowByID)
	},
	resources.Streamlit: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Streamlits.ShowByID)
	},
	resources.Table: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Tables.ShowByID)
	},
	resources.Tag: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Tags.ShowByID)
	},
	resources.Task: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Tasks.ShowByID)
	},
	resources.User: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Users.ShowByID)
	},
	resources.View: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Views.ShowByID)
	},
	resources.Warehouse: func(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) error {
		return runShowById(ctx, id, client.Warehouses.ShowByID)
	},
}

func runShowById[T any, U sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier | sdk.SchemaObjectIdentifierWithArguments](ctx context.Context, id sdk.ObjectIdentifier, show func(ctx context.Context, id U) (T, error)) error {
	idCast, err := asId[U](id)
	if err != nil {
		return err
	}
	_, err = show(ctx, *idCast)
	return err
}

func asId[T sdk.AccountObjectIdentifier | sdk.DatabaseObjectIdentifier | sdk.SchemaObjectIdentifier | sdk.TableColumnIdentifier | sdk.SchemaObjectIdentifierWithArguments](id sdk.ObjectIdentifier) (*T, error) {
	if idCast, ok := id.(T); !ok {
		return nil, fmt.Errorf("expected %s identifier type, but got: %T", reflect.TypeOf(new(T)).Elem().Name(), id)
	} else {
		return &idCast, nil
	}
}

// CheckGrantAccountRoleDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckGrantAccountRoleDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_account_role" {
				continue
			}
			parts := strings.Split(rs.Primary.ID, "|")
			roleName := parts[0]
			roleIdentifier := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(roleName)
			objectType := parts[1]
			targetIdentifier := parts[2]
			grants, err := TestClient().Grant.ShowGrantsOfAccountRole(t, roleIdentifier)
			if err != nil {
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
			if found {
				return fmt.Errorf("role grant %v still exists", rs.Primary.ID)
			}
		}
		return nil
	}
}

// CheckGrantDatabaseRoleDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckGrantDatabaseRoleDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_database_role" {
				continue
			}
			id := rs.Primary.ID
			ids := strings.Split(id, "|")
			databaseRoleName := ids[0]
			objectType := ids[1]
			parentRoleName := ids[2]
			grants, err := TestClient().Grant.ShowGrantsOfDatabaseRole(t, sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName))
			if err != nil {
				continue
			}
			for _, grant := range grants {
				if grant.GrantedTo == sdk.ObjectType(objectType) {
					if grant.GranteeName.FullyQualifiedName() == parentRoleName {
						return fmt.Errorf("database role grant %v still exists", grant)
					}
				}
			}
		}
		return nil
	}
}

// CheckAccountRolePrivilegesRevoked is a custom checks that should be later incorporated into generic CheckDestroy
func CheckAccountRolePrivilegesRevoked(t *testing.T) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_account_role" {
				continue
			}

			id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["account_role_name"])
			grants, err := TestClient().Grant.ShowGrantsToAccountRole(t, id)
			if err != nil {
				if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
					continue
				}
				return err
			}
			var grantedPrivileges []string
			for _, grant := range grants {
				grantedPrivileges = append(grantedPrivileges, grant.Privilege)
			}
			if len(grantedPrivileges) > 0 {
				return fmt.Errorf("account role (%s) is still granted, granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}

// CheckDatabaseRolePrivilegesRevoked is a custom checks that should be later incorporated into generic CheckDestroy
func CheckDatabaseRolePrivilegesRevoked(t *testing.T) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_database_role" {
				continue
			}

			id := sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["database_role_name"])
			grants, err := TestClient().Grant.ShowGrantsToDatabaseRole(t, id)
			if err != nil {
				return err
			}
			var grantedPrivileges []string
			for _, grant := range grants {
				// usage is the default privilege available after creation (it won't be revoked)
				if grant.Privilege != "USAGE" {
					grantedPrivileges = append(grantedPrivileges, grant.Privilege)
				}
			}
			if len(grantedPrivileges) > 0 {
				return fmt.Errorf("database role (%s) is still granted, granted privileges %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}

// CheckSharePrivilegesRevoked is a custom checks that should be later incorporated into generic CheckDestroy
func CheckSharePrivilegesRevoked(t *testing.T) func(*terraform.State) error {
	t.Helper()

	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_grant_privileges_to_share" {
				continue
			}

			id := sdk.NewExternalObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["to_share"])
			grants, err := TestClient().Grant.ShowGrantsToShare(t, sdk.NewAccountObjectIdentifier(id.Name()))
			if err != nil {
				return err
			}
			var grantedPrivileges []string
			for _, grant := range grants {
				grantedPrivileges = append(grantedPrivileges, grant.Privilege)
			}
			if len(grantedPrivileges) > 0 {
				return fmt.Errorf("share (%s) is still granted with privileges: %v", id.FullyQualifiedName(), grantedPrivileges)
			}
		}
		return nil
	}
}

// CheckUserPasswordPolicyAttachmentDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckUserPasswordPolicyAttachmentDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_user_password_policy_attachment" {
				continue
			}
			policyReferences, err := TestClient().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["user_name"]), sdk.PolicyEntityDomainUser)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist or not authorized") {
					// Note: this can happen if the Policy Reference or the User has been deleted as well; in this case, ignore the error
					continue
				}
				return err
			}

			_, err = collections.FindFirst(policyReferences, func(reference sdk.PolicyReference) bool {
				return reference.PolicyKind == sdk.PolicyKindPasswordPolicy &&
					sdk.NewSchemaObjectIdentifier(*reference.PolicyDb, *reference.PolicySchema, reference.PolicyName).FullyQualifiedName() == rs.Primary.Attributes["password_policy_name"]
			})
			if err == nil {
				return fmt.Errorf("user password policy attachment %v still exists", policyReferences[0].PolicyName)
			}
		}
		return nil
	}
}

// CheckUserAuthenticationPolicyAttachmentDestroy is a custom checks that should be later incorporated into generic CheckDestroy
func CheckUserAuthenticationPolicyAttachmentDestroy(t *testing.T) func(*terraform.State) error {
	t.Helper()
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_user_authentication_policy_attachment" {
				continue
			}
			policyReferences, err := TestClient().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["user_name"]), sdk.PolicyEntityDomainUser)
			if err != nil {
				if strings.Contains(err.Error(), "does not exist or not authorized") {
					// Note: this can happen if the Policy Reference or the User has been deleted as well; in this case, ignore the error
					continue
				}
				return err
			}
			if len(policyReferences) > 0 {
				return fmt.Errorf("user authentication policy attachment %v still exists", policyReferences[0].PolicyName)
			}
		}
		return nil
	}
}

// CheckResourceTagUnset is a custom check that should be later incorporated into generic CheckDestroy
func CheckResourceTagUnset(t *testing.T) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_tag_association" {
				continue
			}
			objectType := sdk.ObjectType(rs.Primary.Attributes["object_type"])
			tagId, err := sdk.ParseSchemaObjectIdentifier(rs.Primary.Attributes["tag_id"])
			if err != nil {
				return err
			}
			idLen, err := strconv.Atoi(rs.Primary.Attributes["object_identifiers.#"])
			if err != nil {
				return err
			}
			for i := 0; i < idLen; i++ {
				idRaw := rs.Primary.Attributes[fmt.Sprintf("object_identifiers.%d", i)]
				var id sdk.ObjectIdentifier
				// TODO(SNOW-1229218): Use a common mapper to get object id.
				if objectType == sdk.ObjectTypeAccount {
					id, err = sdk.ParseAccountIdentifier(idRaw)
					if err != nil {
						return fmt.Errorf("invalid account id: %w", err)
					}
				} else {
					id, err = sdk.ParseObjectIdentifierString(idRaw)
					if err != nil {
						return fmt.Errorf("invalid object id: %w", err)
					}
				}
				if err := assertTagUnset(t, tagId, id, objectType); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// CheckTagUnset is a custom check that should be later incorporated into generic CheckDestroy
func CheckTagUnset(t *testing.T, tagId sdk.SchemaObjectIdentifier, id sdk.ObjectIdentifier, objectType sdk.ObjectType) func(*terraform.State) error {
	t.Helper()

	return func(s *terraform.State) error {
		return assertTagUnset(t, tagId, id, objectType)
	}
}

func assertTagUnset(t *testing.T, tagId sdk.SchemaObjectIdentifier, id sdk.ObjectIdentifier, objectType sdk.ObjectType) error {
	t.Helper()

	tag, err := TestClient().Tag.GetForObject(t, tagId, id, objectType)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist or not authorized") {
			// Note: this can happen if the referenced object was deleted before; in this case, ignore the error
			t.Logf("could not get tag for %v : %v, continuing...", id.FullyQualifiedName(), err)
			return nil
		}
		return err
	}
	if tag != nil {
		return fmt.Errorf("tag %s for object %s expected to be empty, got %s", tagId.FullyQualifiedName(), id.FullyQualifiedName(), *tag)
	}
	return err
}

func TestAccCheckGrantApplicationRoleDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_grant_application_role" {
			continue
		}
		ctx := context.Background()
		id := rs.Primary.ID
		ids := strings.Split(id, "|")
		applicationRoleName := ids[0]
		objectType := ids[1]
		parentRoleName := ids[2]
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				ApplicationRole: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(applicationRoleName),
			},
		})
		if err != nil {
			continue
		}
		for _, grant := range grants {
			if grant.GrantedTo == sdk.ObjectType(objectType) {
				if grant.GranteeName.FullyQualifiedName() == parentRoleName {
					return fmt.Errorf("application role grant %v still exists", grant)
				}
			}
		}
	}
	return nil
}

func CheckAccountParameterUnset(t *testing.T, paramName sdk.AccountParameter) func(*terraform.State) error {
	t.Helper()
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_account_parameter" {
				continue
			}
			parameter := TestClient().Parameter.ShowAccountParameter(t, paramName)
			if parameter.Level != sdk.ParameterTypeSnowflakeDefault {
				return fmt.Errorf("expected parameter level empty, got %v", parameter.Level)
			}
		}
		return nil
	}
}
