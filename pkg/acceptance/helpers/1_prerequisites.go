package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
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

func (c *TestClient) EnsureEnableIdentifierFirstLoginIsSetToTrue(ctx context.Context) error {
	log.Printf("[DEBUG] Making sure ENABLE_IDENTIFIER_FIRST_LOGIN parameter is set correctly")
	param, err := c.context.client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterEnableIdentifierFirstLogin)
	if err != nil {
		return fmt.Errorf("checking ENABLE_IDENTIFIER_FIRST_LOGIN resulted in error: %w", err)
	}
	if param.Value != "true" {
		return fmt.Errorf("parameter ENABLE_IDENTIFIER_FIRST_LOGIN has value %s, expected: true", param.Value)
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
			return fmt.Errorf("showing role %s: %w", roleGrant.RoleID.Name(), err)
		}
		grants, err := c.context.client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				Role: roleGrant.RoleID,
			},
		})
		if err != nil {
			return fmt.Errorf("showing grants for role %s: %w", roleGrant.RoleID.Name(), err)
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

func (c *TestClient) EnsureOpenflowPostgresCdcDefinitionExists(ctx context.Context) error {
	log.Printf("[DEBUG] Making sure %s connector definition exists", PostgresCdcDefinitionName)
	definitions, err := c.context.client.OpenflowConnectorDefinitions.Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest().WithLike(sdk.Like{Pattern: sdk.String(PostgresCdcDefinitionName)}))
	if err != nil {
		return fmt.Errorf("checking %s connector definition resulted in error: %w", PostgresCdcDefinitionName, err)
	}
	if len(definitions) == 0 {
		return fmt.Errorf("connector definition %s is not available on this account", PostgresCdcDefinitionName)
	}
	return nil
}

// EnsureOpenflowDeploymentIsProvisioned checks the deployment named by TEST_SF_TF_OPENFLOW_DEPLOYMENT exists
// and is ACTIVE, which runtimes require of the deployment they are created in.
func (c *TestClient) EnsureOpenflowDeploymentIsProvisioned(ctx context.Context) error {
	name := os.Getenv(string(testenvs.OpenflowDeployment))
	log.Printf("[DEBUG] Making sure the Openflow deployment named by %s is provisioned", testenvs.OpenflowDeployment)
	if name == "" {
		return fmt.Errorf("%s must name an ACTIVE Openflow deployment; see the comment on testenvs.OpenflowDeployment", testenvs.OpenflowDeployment)
	}

	deployment, err := c.context.client.OpenflowDeployments.ShowByID(ctx, sdk.NewAccountObjectIdentifier(name))
	if err != nil {
		return fmt.Errorf("checking Openflow deployment %s resulted in error: %w", name, err)
	}
	if deployment.Status != sdk.OpenflowDeploymentStatusActive {
		return fmt.Errorf("openflow deployment %s has status %s, expected %s", name, deployment.Status, sdk.OpenflowDeploymentStatusActive)
	}
	return nil
}

// EnsureOpenflowRuntimeIsProvisioned checks the runtime named by TEST_SF_TF_OPENFLOW_RUNTIME exists. See
// testenvs.OpenflowRuntime for the statements that provision it.
func (c *TestClient) EnsureOpenflowRuntimeIsProvisioned(ctx context.Context) error {
	raw := os.Getenv(string(testenvs.OpenflowRuntime))
	log.Printf("[DEBUG] Making sure the Openflow runtime named by %s is provisioned", testenvs.OpenflowRuntime)
	if raw == "" {
		return fmt.Errorf("%s must name an ACTIVE Openflow runtime by fully qualified name; see the comment on testenvs.OpenflowRuntime", testenvs.OpenflowRuntime)
	}

	id, err := sdk.ParseSchemaObjectIdentifier(raw)
	if err != nil {
		return fmt.Errorf("%s must be a fully qualified name of the form <database>.<schema>.<runtime>, got %q: %w", testenvs.OpenflowRuntime, raw, err)
	}
	if _, err := c.context.client.OpenflowRuntimes.ShowByID(ctx, id); err != nil {
		return fmt.Errorf("checking Openflow runtime %s resulted in error: %w", raw, err)
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
	switch {
	case err != nil:
		t.Errorf("Failed to show organization accounts, err = %v.", err)
	case len(organizationAccounts) != 1:
		t.Errorf("Wrong number of organization accounts returned. Expected one, got = %d.", len(organizationAccounts))
	case organizationAccounts[0].AccountLocator != nonProdModifiableAccountLocator:
		t.Skipf("The TEST_SF_TF_NON_PROD_MODIFIABLE_ACCOUNT_LOCATOR does not match the organization account's locator, please adjust the environment variable.")
	}
}
