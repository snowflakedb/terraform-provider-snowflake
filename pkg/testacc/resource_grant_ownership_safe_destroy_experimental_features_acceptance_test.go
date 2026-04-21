//go:build account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_Experimental_GrantOwnership_SafeDestroy_MissingSchema verifies that destroying a
// grant_ownership resource fails when the target schema is deleted externally (default behavior),
// and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses on.all so that Read returns nil immediately (no SHOW check), ensuring Delete is always called.
func TestAcc_Experimental_GrantOwnership_SafeDestroy_MissingSchema(t *testing.T) {
	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantModel := model.GrantOwnership("test", []sdk.OwnershipGrantOn{{
		All: &sdk.GrantOnSchemaObjectIn{
			PluralObjectType: sdk.PluralObjectTypeTables,
			InSchema:         sdk.Pointer(schema.ID()),
		},
	}}).WithAccountRoleName(role.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with the default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the schema externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Schema.DropSchemaFunc(t, schema.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}

// TestAcc_Experimental_GrantOwnership_SafeDestroy_MissingDatabase verifies that destroying a
// grant_ownership resource fails when the target database is deleted externally (default behavior),
// and succeeds when the GRANTS_SAFE_DESTROY experiment is enabled.
// Uses on.all so that Read returns nil immediately (no SHOW check), ensuring Delete is always called.
func TestAcc_Experimental_GrantOwnership_SafeDestroy_MissingDatabase(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabase(t)
	t.Cleanup(databaseCleanup)

	role, roleCleanup := testClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	grantModel := model.GrantOwnership("test", []sdk.OwnershipGrantOn{{
		All: &sdk.GrantOnSchemaObjectIn{
			PluralObjectType: sdk.PluralObjectTypeTables,
			InDatabase:       sdk.Pointer(database.ID()),
		},
	}}).WithAccountRoleName(role.ID().Name())

	experimentProviderModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.GrantsSafeDestroy)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create the grant with the default provider.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, grantModel),
			},
			// Drop the database externally, then try to destroy without experiment — expect failure.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig:                testClient().Database.DropDatabaseFunc(t, database.ID()),
				Config:                   config.FromModels(t, grantModel),
				Destroy:                  true,
				ExpectError:              regexp.MustCompile("does not exist or not authorized"),
			},
			// Destroy with GRANTS_SAFE_DESTROY experiment — succeeds.
			{
				ProtoV6ProviderFactories: grantsSafeDestroyProviderFactory,
				Config:                   config.FromModels(t, experimentProviderModel, grantModel),
				Destroy:                  true,
			},
		},
	})
}
