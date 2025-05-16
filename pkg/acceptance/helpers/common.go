package helpers

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(client *sdk.Client, ctx context.Context) error {
	log.Printf("[DEBUG] Making sure QUOTED_IDENTIFIERS_IGNORE_CASE parameter is set correctly")
	param, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterQuotedIdentifiersIgnoreCase)
	if err != nil {
		return fmt.Errorf("checking QUOTED_IDENTIFIERS_IGNORE_CASE resulted in error: %w", err)
	}
	if param.Value != "false" {
		return fmt.Errorf("parameter QUOTED_IDENTIFIERS_IGNORE_CASE has value %s, expected: false", param.Value)
	}
	return nil
}

func EnsureScimProvisionerRolesExist(client *sdk.Client, ctx context.Context) error {
	log.Printf("[DEBUG] Making sure Scim Provisioner roles exist")
	roleIDs := []sdk.AccountObjectIdentifier{snowflakeroles.GenericScimProvisioner, snowflakeroles.AadProvisioner, snowflakeroles.OktaProvisioner}
	currentRoleID, err := client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		_, err := client.Roles.ShowByID(ctx, roleID)
		if err != nil {
			return err
		}
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				Role: roleID,
			},
		})
		if err != nil {
			return err
		}
		if !hasGranteeName(grants, currentRoleID) {
			return fmt.Errorf("role %s not granted to %s", currentRoleID.Name(), roleID.Name())
		}
	}
	return nil
}

func hasGranteeName(grants []sdk.Grant, role sdk.AccountObjectIdentifier) bool {
	for _, grant := range grants {
		if grant.GranteeName == role {
			return true
		}
	}
	return false
}

// AssertErrorContainsPartsFunc returns a function asserting error message contains each string in parts
func AssertErrorContainsPartsFunc(t *testing.T, parts []string) resource.ErrorCheckFunc {
	t.Helper()
	return func(err error) error {
		for _, part := range parts {
			assert.Contains(t, err.Error(), part)
		}
		return nil
	}
}
