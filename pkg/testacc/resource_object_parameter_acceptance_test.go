//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ObjectParameter(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigBasic(database.ID(), sdk.DatabaseParameterUserTaskTimeoutMs, "1000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.DatabaseParameterUserTaskTimeoutMs)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "1000"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func objectParameterConfigBasic(databaseId sdk.AccountObjectIdentifier, key sdk.DatabaseParameter, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[2]s"
	value = "%[3]s"
	object_type = "DATABASE"
	object_identifier {
		name = "%[1]s"
	}
}
`, databaseId.Name(), key, value)
}

func TestAcc_ObjectParameterAccount(t *testing.T) {
	// TODO [SNOW-1348325]: Unskip during resource stabilization.
	t.Skip("Skipping temporarily as it messes with the account level setting. The current cleanup is incorrect, so we shouldn't run it even on the account level tests.")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigOnAccount(sdk.AccountParameterDataRetentionTimeInDays, "5"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.AccountParameterDataRetentionTimeInDays)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "5"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "true"),
				),
			},
		},
	})
}

func objectParameterConfigOnAccount(key sdk.AccountParameter, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[1]s"
	value = "%[2]s"
	on_account = true
}
`, key, value)
}

func TestAcc_UserParameter(t *testing.T) {
	user, userCleanup := testClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userParameterConfigBasic(user.ID(), sdk.UserParameterEnableUnredactedQuerySyntaxError, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(sdk.UserParameterEnableUnredactedQuerySyntaxError)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "true"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func userParameterConfigBasic(userId sdk.AccountObjectIdentifier, key sdk.UserParameter, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[2]s"
	value = "%[3]s"
	object_type = "USER"
	object_identifier {
		name = "%[1]s"
	}
}
`, userId.Name(), key, value)
}

func TestAcc_ObjectParameter_ReplicableWithFailoverGroups(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ObjectParameterResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.7.0"),
				Config:            accconfig.FromModels(t, providerModel) + schemaReplicableWithFailoverGroupsConfig(schema.ID(), "NO"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", string(resources.ReplicableWithFailoverGroups)),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "NO"),
				),
			},
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.7.0"),
				Config:            accconfig.FromModels(t, providerModel) + schemaReplicableWithFailoverGroupsConfig(schema.ID(), "NO"),
				Destroy:           true,
				ExpectError:       regexp.MustCompile(`invalid value \[UNSET] for parameter 'REPLICABLE_WITH_FAILOVER_GROUPS'`),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   schemaReplicableWithFailoverGroupsConfig(schema.ID(), "NO"),
				Destroy:                  true,
			},
		},
	})
}

func schemaReplicableWithFailoverGroupsConfig(schemaId sdk.DatabaseObjectIdentifier, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "REPLICABLE_WITH_FAILOVER_GROUPS"
	value = "%[3]s"
	object_type = "SCHEMA"
	object_identifier {
		database = "%[1]s"
		name = "%[2]s"
	}
}
`, schemaId.DatabaseName(), schemaId.Name(), value)
}
