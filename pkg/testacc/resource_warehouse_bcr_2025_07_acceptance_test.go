//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
)

// These tests ensure that the behavior for `generation` attribute is the same as before the 2025_07 bundle.
func TestAcc_Warehouse_Generation_BCR_2025_07(t *testing.T) {
	secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	funcTestAcc_Warehouse_Generation(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary), secondaryTestClient, providerFactoryWithoutCache())
}

func TestAcc_Warehouse_ResourceConstraint_MixedWarehouseTypes_BCR_2025_07(t *testing.T) {
	secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	funcTestAcc_Warehouse_ResourceConstraint_MixedWarehouseTypes(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary), secondaryTestClient, providerFactoryWithoutCache())
}

func TestAcc_Warehouse_Generation_MigrateManuallySetGeneration_BCR_2025_07(t *testing.T) {
	secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	funcTestAcc_Warehouse_Generation_MigrateManuallySetGeneration(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary), secondaryTestClient, providerFactoryWithoutCache())
}

func TestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration_BCR_2025_07(t *testing.T) {
	secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	funcTestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary), secondaryTestClient, providerFactoryWithoutCache())
}

func TestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration_UpdatedExternally_BCR_2025_07(t *testing.T) {
	secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	funcTestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration_UpdatedExternally(t, providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary), secondaryTestClient, providerFactoryWithoutCache())
}
