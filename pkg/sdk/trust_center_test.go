package sdk

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrustCenter_BuildSetConfigurationSQL(t *testing.T) {
	tc := &trustCenter{}

	t.Run("package level enabled", func(t *testing.T) {
		sql := tc.buildSetConfigurationSQL("SECURITY_ESSENTIALS", "", TrustCenterConfigurationLevelEnabled, true)
		assert.Equal(t, "CALL snowflake.trust_center.set_configuration('SECURITY_ESSENTIALS', CONFIGURATION_LEVEL => 'ENABLED', VALUE => true)", sql)
	})

	t.Run("package level disabled", func(t *testing.T) {
		sql := tc.buildSetConfigurationSQL("SECURITY_ESSENTIALS", "", TrustCenterConfigurationLevelEnabled, false)
		assert.Equal(t, "CALL snowflake.trust_center.set_configuration('SECURITY_ESSENTIALS', CONFIGURATION_LEVEL => 'ENABLED', VALUE => false)", sql)
	})

	t.Run("package level schedule", func(t *testing.T) {
		sql := tc.buildSetConfigurationSQL("SECURITY_ESSENTIALS", "", TrustCenterConfigurationLevelSchedule, "USING CRON 0 2 * * * UTC")
		assert.Equal(t, "CALL snowflake.trust_center.set_configuration('SECURITY_ESSENTIALS', CONFIGURATION_LEVEL => 'SCHEDULE', VALUE => 'USING CRON 0 2 * * * UTC')", sql)
	})

	t.Run("package level notification", func(t *testing.T) {
		notification := `{"NOTIFY_ADMINS":true,"SEVERITY_THRESHOLD":"High"}`
		sql := tc.buildSetConfigurationSQL("SECURITY_ESSENTIALS", "", TrustCenterConfigurationLevelNotification, notification)
		assert.Equal(t, `CALL snowflake.trust_center.set_configuration('SECURITY_ESSENTIALS', CONFIGURATION_LEVEL => 'NOTIFICATION', VALUE => PARSE_JSON('{"NOTIFY_ADMINS":true,"SEVERITY_THRESHOLD":"High"}'))`, sql)
	})

	t.Run("scanner level enabled", func(t *testing.T) {
		sql := tc.buildSetConfigurationSQL("SECURITY_ESSENTIALS", "MFA_CHECK", TrustCenterConfigurationLevelEnabled, true)
		assert.Equal(t, "CALL snowflake.trust_center.set_configuration('SECURITY_ESSENTIALS', SCANNER_ID => 'MFA_CHECK', CONFIGURATION_LEVEL => 'ENABLED', VALUE => true)", sql)
	})

	t.Run("scanner level schedule", func(t *testing.T) {
		sql := tc.buildSetConfigurationSQL("SECURITY_ESSENTIALS", "MFA_CHECK", TrustCenterConfigurationLevelSchedule, "USING CRON 0 0 * * * UTC")
		assert.Equal(t, "CALL snowflake.trust_center.set_configuration('SECURITY_ESSENTIALS', SCANNER_ID => 'MFA_CHECK', CONFIGURATION_LEVEL => 'SCHEDULE', VALUE => 'USING CRON 0 0 * * * UTC')", sql)
	})

	t.Run("escape single quotes in package id", func(t *testing.T) {
		sql := tc.buildSetConfigurationSQL("PACKAGE'WITH'QUOTES", "", TrustCenterConfigurationLevelEnabled, true)
		assert.Equal(t, "CALL snowflake.trust_center.set_configuration('PACKAGE''WITH''QUOTES', CONFIGURATION_LEVEL => 'ENABLED', VALUE => true)", sql)
	})
}

func TestTrustCenter_BuildUnsetConfigurationSQL(t *testing.T) {
	tc := &trustCenter{}

	t.Run("package level enabled", func(t *testing.T) {
		sql := tc.buildUnsetConfigurationSQL("SECURITY_ESSENTIALS", "", TrustCenterConfigurationLevelEnabled)
		assert.Equal(t, "CALL snowflake.trust_center.unset_configuration('SECURITY_ESSENTIALS', CONFIGURATION_LEVEL => 'ENABLED')", sql)
	})

	t.Run("package level schedule", func(t *testing.T) {
		sql := tc.buildUnsetConfigurationSQL("SECURITY_ESSENTIALS", "", TrustCenterConfigurationLevelSchedule)
		assert.Equal(t, "CALL snowflake.trust_center.unset_configuration('SECURITY_ESSENTIALS', CONFIGURATION_LEVEL => 'SCHEDULE')", sql)
	})

	t.Run("scanner level enabled", func(t *testing.T) {
		sql := tc.buildUnsetConfigurationSQL("SECURITY_ESSENTIALS", "MFA_CHECK", TrustCenterConfigurationLevelEnabled)
		assert.Equal(t, "CALL snowflake.trust_center.unset_configuration('SECURITY_ESSENTIALS', SCANNER_ID => 'MFA_CHECK', CONFIGURATION_LEVEL => 'ENABLED')", sql)
	})
}

func TestNotificationConfiguration_ToJSON(t *testing.T) {
	t.Run("full configuration", func(t *testing.T) {
		notifyAdmins := true
		severity := "High"
		config := &NotificationConfiguration{
			NotifyAdmins:      &notifyAdmins,
			SeverityThreshold: &severity,
			Users:             []string{"USER1", "USER2"},
		}
		json, err := config.ToJSON()
		require.NoError(t, err)
		assert.Contains(t, json, `"NOTIFY_ADMINS":true`)
		assert.Contains(t, json, `"SEVERITY_THRESHOLD":"High"`)
		assert.Contains(t, json, `"USERS":["USER1","USER2"]`)
	})

	t.Run("partial configuration", func(t *testing.T) {
		severity := "Critical"
		config := &NotificationConfiguration{
			SeverityThreshold: &severity,
		}
		json, err := config.ToJSON()
		require.NoError(t, err)
		assert.Contains(t, json, `"SEVERITY_THRESHOLD":"Critical"`)
		assert.NotContains(t, json, "NOTIFY_ADMINS")
		assert.NotContains(t, json, "USERS")
	})

	t.Run("nil configuration", func(t *testing.T) {
		var config *NotificationConfiguration
		json, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, "", json)
	})
}

func TestParseNotificationConfiguration(t *testing.T) {
	t.Run("full configuration", func(t *testing.T) {
		jsonStr := `{"NOTIFY_ADMINS":true,"SEVERITY_THRESHOLD":"High","USERS":["USER1","USER2"]}`
		config, err := ParseNotificationConfiguration(jsonStr)
		require.NoError(t, err)
		require.NotNil(t, config)
		assert.True(t, *config.NotifyAdmins)
		assert.Equal(t, "High", *config.SeverityThreshold)
		assert.Equal(t, []string{"USER1", "USER2"}, config.Users)
	})

	t.Run("empty string", func(t *testing.T) {
		config, err := ParseNotificationConfiguration("")
		require.NoError(t, err)
		assert.Nil(t, config)
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := ParseNotificationConfiguration("not valid json")
		assert.Error(t, err)
	})
}

func TestTrustCenterScannerPackageId_String(t *testing.T) {
	id := TrustCenterScannerPackageId{
		Source:           "SNOWFLAKE",
		ScannerPackageId: "SECURITY_ESSENTIALS",
	}
	assert.Equal(t, "SNOWFLAKE/SECURITY_ESSENTIALS", id.String())
}

func TestTrustCenterScannerId_String(t *testing.T) {
	id := TrustCenterScannerId{
		Source:           "SNOWFLAKE",
		ScannerPackageId: "SECURITY_ESSENTIALS",
		ScannerId:        "MFA_CHECK",
	}
	assert.Equal(t, "SNOWFLAKE/SECURITY_ESSENTIALS/MFA_CHECK", id.String())
}

func TestScannerPackageRow_Convert(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		row := &scannerPackageRow{
			Name:                  toNullString("Security Essentials"),
			Id:                    toNullString("SECURITY_ESSENTIALS"),
			Description:           toNullString("Core security scanners"),
			DefaultSchedule:       toNullString("USING CRON 0 0 * * * UTC"),
			State:                 toNullString("TRUE"),
			Schedule:              toNullString("USING CRON 0 2 * * * UTC"),
			Notification:          toNullString(`{"NOTIFY_ADMINS":true}`),
			Provider:              toNullString("Snowflake"),
			LastEnabledTimestamp:  toNullString("2024-01-01T00:00:00Z"),
			LastDisabledTimestamp: toNullString("2024-01-02T00:00:00Z"),
		}
		pkg := row.convert()
		assert.Equal(t, "Security Essentials", pkg.Name)
		assert.Equal(t, "SECURITY_ESSENTIALS", pkg.Id)
		assert.Equal(t, "Core security scanners", pkg.Description)
		assert.Equal(t, "TRUE", pkg.State)
		assert.NotNil(t, pkg.LastEnabledTimestamp)
		assert.NotNil(t, pkg.LastDisabledTimestamp)
	})

	t.Run("with null timestamps", func(t *testing.T) {
		row := &scannerPackageRow{
			Name:  toNullString("Test Package"),
			Id:    toNullString("TEST_PACKAGE"),
			State: toNullString("FALSE"),
		}
		pkg := row.convert()
		assert.Nil(t, pkg.LastEnabledTimestamp)
		assert.Nil(t, pkg.LastDisabledTimestamp)
	})
}

func TestScannerRow_Convert(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		row := &scannerRow{
			Name:              toNullString("MFA Check"),
			Id:                toNullString("MFA_CHECK"),
			ShortDescription:  toNullString("Checks MFA status"),
			Description:       toNullString("Detailed MFA check description"),
			ScannerPackageId:  toNullString("SECURITY_ESSENTIALS"),
			State:             toNullString("TRUE"),
			Schedule:          toNullString("USING CRON 0 0 * * * UTC"),
			Notification:      toNullString(`{"NOTIFY_ADMINS":true}`),
			LastScanTimestamp: toNullString("2024-01-01T00:00:00Z"),
		}
		scanner := row.convert()
		assert.Equal(t, "MFA Check", scanner.Name)
		assert.Equal(t, "MFA_CHECK", scanner.Id)
		assert.Equal(t, "SECURITY_ESSENTIALS", scanner.ScannerPackageId)
		assert.NotNil(t, scanner.LastScanTimestamp)
	})

	t.Run("with null timestamp", func(t *testing.T) {
		row := &scannerRow{
			Name:             toNullString("Test Scanner"),
			Id:               toNullString("TEST_SCANNER"),
			ScannerPackageId: toNullString("TEST_PACKAGE"),
			State:            toNullString("FALSE"),
		}
		scanner := row.convert()
		assert.Nil(t, scanner.LastScanTimestamp)
	})
}

// Helper function to create sql.NullString
func toNullString(s string) (ns sql.NullString) {
	if s == "" {
		return
	}
	return sql.NullString{String: s, Valid: true}
}
