package helpers

import (
	"context"
	"fmt"
	"log"

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

func (c *TestClient) EnsureScimProvisionerRolesExist(ctx context.Context) error {
	log.Printf("[DEBUG] Making sure Scim Provisioner roles exist")
	roleIDs := []sdk.AccountObjectIdentifier{snowflakeroles.GenericScimProvisioner, snowflakeroles.AadProvisioner, snowflakeroles.OktaProvisioner}
	currentRoleID, err := c.context.client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		_, err := c.context.client.Roles.ShowByID(ctx, roleID)
		if err != nil {
			return err
		}
		grants, err := c.context.client.Grants.Show(ctx, &sdk.ShowGrantOptions{
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

func (c *TestClient) EnsureImageRepositoryExist(ctx context.Context) error {
	id := sdk.NewSchemaObjectIdentifier("SNOWFLAKE", "IMAGES", "SNOWFLAKE_IMAGES")
	log.Printf("[DEBUG] Making sure %s image repository exist", id.FullyQualifiedName())
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
