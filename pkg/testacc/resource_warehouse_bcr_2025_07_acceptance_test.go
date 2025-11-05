//go:build account_level_tests

package testacc

import (
	"testing"
)

// These tests ensure that the behavior for `generation` attribute is the same as before the 2025_07 bundle.
func TestAcc_Warehouse_Generation_BCR_2025_07(t *testing.T) {
	testClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	TestAcc_Warehouse_Generation(t)
}

func TestAcc_Warehouse_ResourceConstraint_MixedWarehouseTypes_BCR_2025_07(t *testing.T) {
	testClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	TestAcc_Warehouse_ResourceConstraint_MixedWarehouseTypes(t)
}

func TestAcc_Warehouse_Generation_MigrateManuallySetGeneration_BCR_2025_07(t *testing.T) {
	testClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	TestAcc_Warehouse_Generation_MigrateManuallySetGeneration(t)
}

func TestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration_BCR_2025_07(t *testing.T) {
	testClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	TestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration(t)
}

func TestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration_UpdatedExternally_BCR_2025_07(t *testing.T) {
	testClient().BcrBundles.EnableBcrBundle(t, "2025_07")
	TestAcc_Warehouse_Generation_MigrateStandardWithoutGeneration_UpdatedExternally(t)
}
