package sdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// SetPackageConfiguration sets configuration for a scanner package.
// It calls the trust_center.set_configuration stored procedure.
func (tc *trustCenter) SetPackageConfiguration(ctx context.Context, req *SetPackageConfigurationRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.ScannerPackageId == "" {
		return fmt.Errorf("scanner_package_id is required")
	}

	// Set ENABLED if specified
	if req.Enabled != nil {
		sql := tc.buildSetConfigurationSQL(req.ScannerPackageId, "", TrustCenterConfigurationLevelEnabled, *req.Enabled)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to set ENABLED configuration: %w", err)
		}
	}

	// Set SCHEDULE if specified
	if req.Schedule != nil {
		sql := tc.buildSetConfigurationSQL(req.ScannerPackageId, "", TrustCenterConfigurationLevelSchedule, *req.Schedule)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to set SCHEDULE configuration: %w", err)
		}
	}

	// Set NOTIFICATION if specified
	if req.Notification != nil {
		notificationJSON, err := req.Notification.ToJSON()
		if err != nil {
			return fmt.Errorf("failed to serialize notification configuration: %w", err)
		}
		sql := tc.buildSetConfigurationSQL(req.ScannerPackageId, "", TrustCenterConfigurationLevelNotification, notificationJSON)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to set NOTIFICATION configuration: %w", err)
		}
	}

	return nil
}

// UnsetPackageConfiguration unsets configuration for a scanner package.
func (tc *trustCenter) UnsetPackageConfiguration(ctx context.Context, req *UnsetPackageConfigurationRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.ScannerPackageId == "" {
		return fmt.Errorf("scanner_package_id is required")
	}

	// Unset ENABLED
	if req.UnsetEnabled {
		sql := tc.buildUnsetConfigurationSQL(req.ScannerPackageId, "", TrustCenterConfigurationLevelEnabled)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to unset ENABLED configuration: %w", err)
		}
	}

	// Unset SCHEDULE
	if req.UnsetSchedule {
		sql := tc.buildUnsetConfigurationSQL(req.ScannerPackageId, "", TrustCenterConfigurationLevelSchedule)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to unset SCHEDULE configuration: %w", err)
		}
	}

	// Unset NOTIFICATION
	if req.UnsetNotification {
		sql := tc.buildUnsetConfigurationSQL(req.ScannerPackageId, "", TrustCenterConfigurationLevelNotification)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to unset NOTIFICATION configuration: %w", err)
		}
	}

	return nil
}

// ShowScannerPackages lists scanner packages from the trust_center.scanner_packages view.
func (tc *trustCenter) ShowScannerPackages(ctx context.Context, req *ShowScannerPackagesRequest) ([]ScannerPackage, error) {
	sql := "SELECT * FROM snowflake.trust_center.scanner_packages"
	if req != nil && req.Like != nil && *req.Like != "" {
		sql += fmt.Sprintf(" WHERE NAME LIKE '%s'", strings.ReplaceAll(*req.Like, "'", "''"))
	}

	var rows []scannerPackageRow
	if err := tc.client.query(ctx, &rows, sql); err != nil {
		return nil, fmt.Errorf("failed to query scanner packages: %w", err)
	}

	packages := make([]ScannerPackage, len(rows))
	for i, row := range rows {
		packages[i] = *row.convert()
	}
	return packages, nil
}

// ShowScannerPackageByID retrieves a single scanner package by ID.
func (tc *trustCenter) ShowScannerPackageByID(ctx context.Context, id string) (*ScannerPackage, error) {
	packages, err := tc.ShowScannerPackages(ctx, &ShowScannerPackagesRequest{Like: &id})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(packages, func(p ScannerPackage) bool { return p.Id == id })
}

// SetScannerConfiguration sets configuration for an individual scanner.
func (tc *trustCenter) SetScannerConfiguration(ctx context.Context, req *SetScannerConfigurationRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.ScannerPackageId == "" {
		return fmt.Errorf("scanner_package_id is required")
	}
	if req.ScannerId == "" {
		return fmt.Errorf("scanner_id is required")
	}

	// Set ENABLED if specified
	if req.Enabled != nil {
		sql := tc.buildSetConfigurationSQL(req.ScannerPackageId, req.ScannerId, TrustCenterConfigurationLevelEnabled, *req.Enabled)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to set ENABLED configuration: %w", err)
		}
	}

	// Set SCHEDULE if specified
	if req.Schedule != nil {
		sql := tc.buildSetConfigurationSQL(req.ScannerPackageId, req.ScannerId, TrustCenterConfigurationLevelSchedule, *req.Schedule)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to set SCHEDULE configuration: %w", err)
		}
	}

	// Set NOTIFICATION if specified
	if req.Notification != nil {
		notificationJSON, err := req.Notification.ToJSON()
		if err != nil {
			return fmt.Errorf("failed to serialize notification configuration: %w", err)
		}
		sql := tc.buildSetConfigurationSQL(req.ScannerPackageId, req.ScannerId, TrustCenterConfigurationLevelNotification, notificationJSON)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to set NOTIFICATION configuration: %w", err)
		}
	}

	return nil
}

// UnsetScannerConfiguration unsets configuration for an individual scanner.
func (tc *trustCenter) UnsetScannerConfiguration(ctx context.Context, req *UnsetScannerConfigurationRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.ScannerPackageId == "" {
		return fmt.Errorf("scanner_package_id is required")
	}
	if req.ScannerId == "" {
		return fmt.Errorf("scanner_id is required")
	}

	// Unset ENABLED
	if req.UnsetEnabled {
		sql := tc.buildUnsetConfigurationSQL(req.ScannerPackageId, req.ScannerId, TrustCenterConfigurationLevelEnabled)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to unset ENABLED configuration: %w", err)
		}
	}

	// Unset SCHEDULE
	if req.UnsetSchedule {
		sql := tc.buildUnsetConfigurationSQL(req.ScannerPackageId, req.ScannerId, TrustCenterConfigurationLevelSchedule)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to unset SCHEDULE configuration: %w", err)
		}
	}

	// Unset NOTIFICATION
	if req.UnsetNotification {
		sql := tc.buildUnsetConfigurationSQL(req.ScannerPackageId, req.ScannerId, TrustCenterConfigurationLevelNotification)
		if _, err := tc.client.exec(ctx, sql); err != nil {
			return fmt.Errorf("failed to unset NOTIFICATION configuration: %w", err)
		}
	}

	return nil
}

// ShowScanners lists scanners from the trust_center.scanners view.
func (tc *trustCenter) ShowScanners(ctx context.Context, req *ShowScannersRequest) ([]Scanner, error) {
	sql := "SELECT * FROM snowflake.trust_center.scanners"

	var conditions []string
	if req != nil {
		if req.ScannerPackageId != nil && *req.ScannerPackageId != "" {
			conditions = append(conditions, fmt.Sprintf("SCANNER_PACKAGE_ID = '%s'", strings.ReplaceAll(*req.ScannerPackageId, "'", "''")))
		}
		if req.Like != nil && *req.Like != "" {
			conditions = append(conditions, fmt.Sprintf("NAME LIKE '%s'", strings.ReplaceAll(*req.Like, "'", "''")))
		}
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	var rows []scannerRow
	if err := tc.client.query(ctx, &rows, sql); err != nil {
		return nil, fmt.Errorf("failed to query scanners: %w", err)
	}

	scanners := make([]Scanner, len(rows))
	for i, row := range rows {
		scanners[i] = *row.convert()
	}
	return scanners, nil
}

// ShowScannerByID retrieves a single scanner by package ID and scanner ID.
func (tc *trustCenter) ShowScannerByID(ctx context.Context, packageId, scannerId string) (*Scanner, error) {
	scanners, err := tc.ShowScanners(ctx, &ShowScannersRequest{ScannerPackageId: &packageId})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(scanners, func(s Scanner) bool { return s.Id == scannerId })
}

// buildSetConfigurationSQL builds the SQL CALL statement for set_configuration.
// Format: CALL snowflake.trust_center.set_configuration(
//
//	'<scanner_package_id>',
//	[SCANNER_ID => '<scanner_id>',]
//	CONFIGURATION_LEVEL => '<level>',
//	VALUE => <value>
//
// )
func (tc *trustCenter) buildSetConfigurationSQL(packageId, scannerId string, level TrustCenterConfigurationLevel, value interface{}) string {
	var sb strings.Builder
	sb.WriteString("CALL snowflake.trust_center.set_configuration(")
	sb.WriteString(fmt.Sprintf("'%s'", strings.ReplaceAll(packageId, "'", "''")))

	if scannerId != "" {
		sb.WriteString(fmt.Sprintf(", SCANNER_ID => '%s'", strings.ReplaceAll(scannerId, "'", "''")))
	}

	sb.WriteString(fmt.Sprintf(", CONFIGURATION_LEVEL => '%s'", level))

	switch v := value.(type) {
	case bool:
		sb.WriteString(fmt.Sprintf(", VALUE => %t", v))
	case string:
		if level == TrustCenterConfigurationLevelNotification {
			// Notification is passed as PARSE_JSON
			sb.WriteString(fmt.Sprintf(", VALUE => PARSE_JSON('%s')", strings.ReplaceAll(v, "'", "''")))
		} else {
			sb.WriteString(fmt.Sprintf(", VALUE => '%s'", strings.ReplaceAll(v, "'", "''")))
		}
	default:
		sb.WriteString(fmt.Sprintf(", VALUE => '%v'", v))
	}

	sb.WriteString(")")
	return sb.String()
}

// buildUnsetConfigurationSQL builds the SQL CALL statement for unset_configuration.
func (tc *trustCenter) buildUnsetConfigurationSQL(packageId, scannerId string, level TrustCenterConfigurationLevel) string {
	var sb strings.Builder
	sb.WriteString("CALL snowflake.trust_center.unset_configuration(")
	sb.WriteString(fmt.Sprintf("'%s'", strings.ReplaceAll(packageId, "'", "''")))

	if scannerId != "" {
		sb.WriteString(fmt.Sprintf(", SCANNER_ID => '%s'", strings.ReplaceAll(scannerId, "'", "''")))
	}

	sb.WriteString(fmt.Sprintf(", CONFIGURATION_LEVEL => '%s'", level))
	sb.WriteString(")")
	return sb.String()
}
