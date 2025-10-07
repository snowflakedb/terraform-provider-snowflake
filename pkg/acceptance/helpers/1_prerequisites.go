package helpers

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *TestClient) EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(ctx context.Context) error {
	log.Printf("[DEBUG] Making sure QUOTED_IDENTIFIERS_IGNORE_CASE parameter is set correctly")
	param, err := c.context.client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterQuotedIdentifiersIgnoreCase)
	if err != nil {
		return fmt.Errorf("checking QUOTED_IDENTIFIERS_IGNORE_CASE resulted in error: %w", err)
	}
	if param.Value != "false" {
		return fmt.Errorf("parameter QUOTED_IDENTIFIERS_IGNORE_CASE has value %s, expected: false", param.Value)
	}
	return nil
}

func (c *TestClient) EnsureEssentialRolesExist(ctx context.Context) error {
	log.Printf("[DEBUG] Making sure essential roles exist")
	type RoleGrant struct {
		RoleID          sdk.AccountObjectIdentifier
		ShouldBeGranted bool
	}
	roleGrants := []RoleGrant{
		{RoleID: snowflakeroles.GenericScimProvisioner, ShouldBeGranted: true},
		{RoleID: snowflakeroles.AadProvisioner, ShouldBeGranted: true},
		{RoleID: snowflakeroles.OktaProvisioner, ShouldBeGranted: true},
		{RoleID: snowflakeroles.Restricted, ShouldBeGranted: false},
	}
	currentRoleID, err := c.context.client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}
	for _, roleGrant := range roleGrants {
		_, err := c.context.client.Roles.ShowByID(ctx, roleGrant.RoleID)
		if err != nil {
			return err
		}
		grants, err := c.context.client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				Role: roleGrant.RoleID,
			},
		})
		if err != nil {
			return err
		}
		isGranted := hasGranteeName(grants, currentRoleID)
		if roleGrant.ShouldBeGranted && !isGranted {
			return fmt.Errorf("role %s should be granted to %s, but is not", roleGrant.RoleID.Name(), currentRoleID.Name())
		}
		if !roleGrant.ShouldBeGranted && isGranted {
			return fmt.Errorf("role %s should not be granted to %s, but is", roleGrant.RoleID.Name(), currentRoleID.Name())
		}
	}
	return nil
}

func (c *TestClient) EnsureImageRepositoryExist(ctx context.Context) error {
	id := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "IMAGES", "SNOWFLAKE_IMAGES")
	log.Printf("[DEBUG] Making sure %s image repository exists", id.FullyQualifiedName())
	_, err := c.context.client.ImageRepositories.ShowByID(ctx, id)
	if err != nil {
		return err
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

func (c *TestClient) EnsureValidNonProdAccountIsUsed(t *testing.T) {
	t.Helper()
	testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)
	nonProdModifiableAccountLocator := testenvs.GetOrSkipTest(t, testenvs.TestNonProdModifiableAccountLocator)
	if c.GetAccountLocator() != nonProdModifiableAccountLocator {
		t.Skipf("Current client account locator does not match the required non-prod modifiable account's locator set in %s env variable. Skipping test.", testenvs.TestNonProdModifiableAccountLocator)
	}
}

func (c *TestClient) EnsureValidNonProdOrganizationAccountIsUsed(t *testing.T) {
	t.Helper()
	nonProdModifiableAccountLocator := testenvs.GetOrSkipTest(t, testenvs.TestNonProdModifiableAccountLocator)
	if c.GetAccountLocator() != nonProdModifiableAccountLocator {
		t.Skipf("Current client account locator does not match the required non-prod modifiable account's locator set in %s env variable. Skipping test.", testenvs.TestNonProdModifiableAccountLocator)
	}
	organizationAccounts, err := c.context.client.OrganizationAccounts.Show(context.Background(), sdk.NewShowOrganizationAccountRequest())
	if err != nil {
		t.Errorf("Failed to show organization accounts, err = %v.", err)
	}
	if len(organizationAccounts) != 1 {
		t.Errorf("Wrong number of organization accounts returned. Expected one, got = %d.", len(organizationAccounts))
	}
	if organizationAccounts[0].AccountLocator != nonProdModifiableAccountLocator {
		t.Skipf("The TEST_SF_TF_NON_PROD_MODIFIABLE_ACCOUNT_LOCATOR does not match the organization account's locator, please adjust the environment variable.")
	}
}
