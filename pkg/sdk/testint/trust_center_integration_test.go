//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_TrustCenter_ShowScannerPackages(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("basic show", func(t *testing.T) {
		packages, err := client.TrustCenter.ShowScannerPackages(ctx, &sdk.ShowScannerPackagesRequest{})
		require.NoError(t, err)
		assert.NotEmpty(t, packages)
	})

	t.Run("with like filter", func(t *testing.T) {
		like := "SECURITY%"
		packages, err := client.TrustCenter.ShowScannerPackages(ctx, &sdk.ShowScannerPackagesRequest{Like: &like})
		require.NoError(t, err)
		assert.NotEmpty(t, packages)
		for _, pkg := range packages {
			assert.Contains(t, pkg.Name, "Security")
		}
	})
}

func TestInt_TrustCenter_ShowScannerPackageByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("found", func(t *testing.T) {
		pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, "SECURITY_ESSENTIALS")
		require.NoError(t, err)
		assert.Equal(t, "SECURITY_ESSENTIALS", pkg.Id)
		assert.NotEmpty(t, pkg.Name)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := client.TrustCenter.ShowScannerPackageByID(ctx, "NON_EXISTENT_PACKAGE")
		require.Error(t, err)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})
}

func TestInt_TrustCenter_SetUnsetPackageConfiguration(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	packageId := "SECURITY_ESSENTIALS"

	// Ensure cleanup: disable after test
	t.Cleanup(func() {
		enabled := false
		_ = client.TrustCenter.SetPackageConfiguration(ctx, &sdk.SetPackageConfigurationRequest{
			ScannerPackageId: packageId,
			Enabled:          &enabled,
		})
		_ = client.TrustCenter.UnsetPackageConfiguration(ctx, &sdk.UnsetPackageConfigurationRequest{
			ScannerPackageId:  packageId,
			UnsetSchedule:     true,
			UnsetNotification: true,
		})
	})

	t.Run("set enabled", func(t *testing.T) {
		enabled := true
		err := client.TrustCenter.SetPackageConfiguration(ctx, &sdk.SetPackageConfigurationRequest{
			ScannerPackageId: packageId,
			Enabled:          &enabled,
		})
		require.NoError(t, err)

		pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, packageId)
		require.NoError(t, err)
		assert.Equal(t, "TRUE", pkg.State)
	})

	t.Run("set schedule", func(t *testing.T) {
		schedule := "USING CRON 0 2 * * * UTC"
		err := client.TrustCenter.SetPackageConfiguration(ctx, &sdk.SetPackageConfigurationRequest{
			ScannerPackageId: packageId,
			Schedule:         &schedule,
		})
		require.NoError(t, err)

		pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, packageId)
		require.NoError(t, err)
		assert.Equal(t, schedule, pkg.Schedule)
	})

	t.Run("set notification", func(t *testing.T) {
		notifyAdmins := true
		severity := "High"
		err := client.TrustCenter.SetPackageConfiguration(ctx, &sdk.SetPackageConfigurationRequest{
			ScannerPackageId: packageId,
			Notification: &sdk.NotificationConfiguration{
				NotifyAdmins:      &notifyAdmins,
				SeverityThreshold: &severity,
			},
		})
		require.NoError(t, err)

		pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, packageId)
		require.NoError(t, err)
		assert.NotEmpty(t, pkg.Notification)
	})

	t.Run("unset schedule and notification", func(t *testing.T) {
		err := client.TrustCenter.UnsetPackageConfiguration(ctx, &sdk.UnsetPackageConfigurationRequest{
			ScannerPackageId:  packageId,
			UnsetSchedule:     true,
			UnsetNotification: true,
		})
		require.NoError(t, err)
	})

	t.Run("unset enabled", func(t *testing.T) {
		err := client.TrustCenter.UnsetPackageConfiguration(ctx, &sdk.UnsetPackageConfigurationRequest{
			ScannerPackageId: packageId,
			UnsetEnabled:     true,
		})
		require.NoError(t, err)

		pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, packageId)
		require.NoError(t, err)
		assert.Equal(t, "FALSE", pkg.State)
	})
}

func TestInt_TrustCenter_ShowScanners(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("basic show", func(t *testing.T) {
		scanners, err := client.TrustCenter.ShowScanners(ctx, &sdk.ShowScannersRequest{})
		require.NoError(t, err)
		assert.NotEmpty(t, scanners)
	})

	t.Run("with package filter", func(t *testing.T) {
		packageId := "SECURITY_ESSENTIALS"
		scanners, err := client.TrustCenter.ShowScanners(ctx, &sdk.ShowScannersRequest{ScannerPackageId: &packageId})
		require.NoError(t, err)
		assert.NotEmpty(t, scanners)
		for _, s := range scanners {
			assert.Equal(t, packageId, s.ScannerPackageId)
		}
	})

	t.Run("with like filter", func(t *testing.T) {
		like := "MFA%"
		_, err := client.TrustCenter.ShowScanners(ctx, &sdk.ShowScannersRequest{Like: &like})
		require.NoError(t, err)
		// May or may not find results depending on naming, just check no error
	})
}

func TestInt_TrustCenter_ShowScannerByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("found", func(t *testing.T) {
		// First get a known scanner from the list
		packageId := "SECURITY_ESSENTIALS"
		scanners, err := client.TrustCenter.ShowScanners(ctx, &sdk.ShowScannersRequest{ScannerPackageId: &packageId})
		require.NoError(t, err)
		require.NotEmpty(t, scanners)

		scanner, err := client.TrustCenter.ShowScannerByID(ctx, packageId, scanners[0].Id)
		require.NoError(t, err)
		assert.Equal(t, scanners[0].Id, scanner.Id)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := client.TrustCenter.ShowScannerByID(ctx, "SECURITY_ESSENTIALS", "NON_EXISTENT_SCANNER")
		require.Error(t, err)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})
}

func TestInt_TrustCenter_SetUnsetScannerConfiguration(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	packageId := "SECURITY_ESSENTIALS"

	// Get a known scanner
	scanners, err := client.TrustCenter.ShowScanners(ctx, &sdk.ShowScannersRequest{ScannerPackageId: &packageId})
	require.NoError(t, err)
	require.NotEmpty(t, scanners)
	scannerId := scanners[0].Id

	// Ensure cleanup: disable after test
	t.Cleanup(func() {
		enabled := false
		_ = client.TrustCenter.SetScannerConfiguration(ctx, &sdk.SetScannerConfigurationRequest{
			ScannerPackageId: packageId,
			ScannerId:        scannerId,
			Enabled:          &enabled,
		})
		_ = client.TrustCenter.UnsetScannerConfiguration(ctx, &sdk.UnsetScannerConfigurationRequest{
			ScannerPackageId:  packageId,
			ScannerId:         scannerId,
			UnsetSchedule:     true,
			UnsetNotification: true,
		})
	})

	t.Run("set enabled", func(t *testing.T) {
		enabled := true
		err := client.TrustCenter.SetScannerConfiguration(ctx, &sdk.SetScannerConfigurationRequest{
			ScannerPackageId: packageId,
			ScannerId:        scannerId,
			Enabled:          &enabled,
		})
		require.NoError(t, err)

		scanner, err := client.TrustCenter.ShowScannerByID(ctx, packageId, scannerId)
		require.NoError(t, err)
		assert.Equal(t, "TRUE", scanner.State)
	})

	t.Run("set schedule", func(t *testing.T) {
		schedule := "USING CRON 0 0 * * * UTC"
		err := client.TrustCenter.SetScannerConfiguration(ctx, &sdk.SetScannerConfigurationRequest{
			ScannerPackageId: packageId,
			ScannerId:        scannerId,
			Schedule:         &schedule,
		})
		require.NoError(t, err)

		scanner, err := client.TrustCenter.ShowScannerByID(ctx, packageId, scannerId)
		require.NoError(t, err)
		assert.Equal(t, schedule, scanner.Schedule)
	})

	t.Run("set notification", func(t *testing.T) {
		notifyAdmins := true
		severity := "High"
		err := client.TrustCenter.SetScannerConfiguration(ctx, &sdk.SetScannerConfigurationRequest{
			ScannerPackageId: packageId,
			ScannerId:        scannerId,
			Notification: &sdk.NotificationConfiguration{
				NotifyAdmins:      &notifyAdmins,
				SeverityThreshold: &severity,
			},
		})
		require.NoError(t, err)

		scanner, err := client.TrustCenter.ShowScannerByID(ctx, packageId, scannerId)
		require.NoError(t, err)
		assert.NotEmpty(t, scanner.Notification)
	})

	t.Run("unset schedule and notification", func(t *testing.T) {
		err := client.TrustCenter.UnsetScannerConfiguration(ctx, &sdk.UnsetScannerConfigurationRequest{
			ScannerPackageId:  packageId,
			ScannerId:         scannerId,
			UnsetSchedule:     true,
			UnsetNotification: true,
		})
		require.NoError(t, err)
	})

	t.Run("unset enabled", func(t *testing.T) {
		err := client.TrustCenter.UnsetScannerConfiguration(ctx, &sdk.UnsetScannerConfigurationRequest{
			ScannerPackageId: packageId,
			ScannerId:        scannerId,
			UnsetEnabled:     true,
		})
		require.NoError(t, err)

		scanner, err := client.TrustCenter.ShowScannerByID(ctx, packageId, scannerId)
		require.NoError(t, err)
		assert.Equal(t, "FALSE", scanner.State)
	})
}
