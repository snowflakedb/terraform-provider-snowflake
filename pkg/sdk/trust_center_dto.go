package sdk

import (
	"database/sql"
	"encoding/json"
)

// ConfigurationLevel represents the type of configuration being set.
type TrustCenterConfigurationLevel string

const (
	TrustCenterConfigurationLevelEnabled      TrustCenterConfigurationLevel = "ENABLED"
	TrustCenterConfigurationLevelSchedule     TrustCenterConfigurationLevel = "SCHEDULE"
	TrustCenterConfigurationLevelNotification TrustCenterConfigurationLevel = "NOTIFICATION"
)

// ScannerPackage represents a Trust Center scanner package from the scanner_packages view.
type ScannerPackage struct {
	Name                  string
	Id                    string
	Description           string
	DefaultSchedule       string
	State                 string // "TRUE" or "FALSE"
	Schedule              string
	Notification          string // JSON string
	Provider              string
	LastEnabledTimestamp  *string
	LastDisabledTimestamp *string
}

// scannerPackageRow is the database row structure for scanner packages.
type scannerPackageRow struct {
	Name                  sql.NullString `db:"NAME"`
	Id                    sql.NullString `db:"ID"`
	Description           sql.NullString `db:"DESCRIPTION"`
	DefaultSchedule       sql.NullString `db:"DEFAULT_SCHEDULE"`
	State                 sql.NullString `db:"STATE"`
	Schedule              sql.NullString `db:"SCHEDULE"`
	Notification          sql.NullString `db:"NOTIFICATION"`
	Provider              sql.NullString `db:"PROVIDER"`
	LastEnabledTimestamp  sql.NullString `db:"LAST_ENABLED_TIMESTAMP"`
	LastDisabledTimestamp sql.NullString `db:"LAST_DISABLED_TIMESTAMP"`
}

func (r *scannerPackageRow) convert() *ScannerPackage {
	pkg := &ScannerPackage{
		Name:            r.Name.String,
		Id:              r.Id.String,
		Description:     r.Description.String,
		DefaultSchedule: r.DefaultSchedule.String,
		State:           r.State.String,
		Schedule:        r.Schedule.String,
		Notification:    r.Notification.String,
		Provider:        r.Provider.String,
	}
	if r.LastEnabledTimestamp.Valid {
		pkg.LastEnabledTimestamp = &r.LastEnabledTimestamp.String
	}
	if r.LastDisabledTimestamp.Valid {
		pkg.LastDisabledTimestamp = &r.LastDisabledTimestamp.String
	}
	return pkg
}

// Scanner represents a Trust Center scanner from the scanners view.
type Scanner struct {
	Name              string
	Id                string
	ShortDescription  string
	Description       string
	ScannerPackageId  string
	State             string
	Schedule          string
	Notification      string // JSON string
	LastScanTimestamp *string
}

// scannerRow is the database row structure for scanners.
type scannerRow struct {
	Name              sql.NullString `db:"NAME"`
	Id                sql.NullString `db:"ID"`
	ShortDescription  sql.NullString `db:"SHORT_DESCRIPTION"`
	Description       sql.NullString `db:"DESCRIPTION"`
	ScannerPackageId  sql.NullString `db:"SCANNER_PACKAGE_ID"`
	State             sql.NullString `db:"STATE"`
	Schedule          sql.NullString `db:"SCHEDULE"`
	Notification      sql.NullString `db:"NOTIFICATION"`
	LastScanTimestamp sql.NullString `db:"LAST_SCAN_TIMESTAMP"`
}

func (r *scannerRow) convert() *Scanner {
	scanner := &Scanner{
		Name:             r.Name.String,
		Id:               r.Id.String,
		ShortDescription: r.ShortDescription.String,
		Description:      r.Description.String,
		ScannerPackageId: r.ScannerPackageId.String,
		State:            r.State.String,
		Schedule:         r.Schedule.String,
		Notification:     r.Notification.String,
	}
	if r.LastScanTimestamp.Valid {
		scanner.LastScanTimestamp = &r.LastScanTimestamp.String
	}
	return scanner
}

// NotificationConfiguration represents the notification settings for a scanner or package.
type NotificationConfiguration struct {
	NotifyAdmins      *bool    `json:"NOTIFY_ADMINS,omitempty"`
	SeverityThreshold *string  `json:"SEVERITY_THRESHOLD,omitempty"`
	Users             []string `json:"USERS,omitempty"`
}

// ToJSON converts the notification configuration to a JSON string for use in SQL calls.
func (n *NotificationConfiguration) ToJSON() (string, error) {
	if n == nil {
		return "", nil
	}
	bytes, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseNotificationConfiguration parses a JSON string into a NotificationConfiguration.
func ParseNotificationConfiguration(jsonStr string) (*NotificationConfiguration, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var config NotificationConfiguration
	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// SetPackageConfigurationRequest is the request for setting scanner package configuration.
type SetPackageConfigurationRequest struct {
	ScannerPackageId string
	Enabled          *bool
	Schedule         *string
	Notification     *NotificationConfiguration
}

// UnsetPackageConfigurationRequest is the request for unsetting scanner package configuration.
type UnsetPackageConfigurationRequest struct {
	ScannerPackageId   string
	UnsetEnabled       bool
	UnsetSchedule      bool
	UnsetNotification  bool
}

// ShowScannerPackagesRequest is the request for showing scanner packages.
type ShowScannerPackagesRequest struct {
	Like *string
}

// SetScannerConfigurationRequest is the request for setting individual scanner configuration.
type SetScannerConfigurationRequest struct {
	ScannerPackageId string
	ScannerId        string
	Enabled          *bool
	Schedule         *string
	Notification     *NotificationConfiguration
}

// UnsetScannerConfigurationRequest is the request for unsetting individual scanner configuration.
type UnsetScannerConfigurationRequest struct {
	ScannerPackageId  string
	ScannerId         string
	UnsetEnabled      bool
	UnsetSchedule     bool
	UnsetNotification bool
}

// ShowScannersRequest is the request for showing scanners.
type ShowScannersRequest struct {
	ScannerPackageId *string
	Like             *string
}

// TrustCenterScannerPackageId represents the identifier for a scanner package resource.
type TrustCenterScannerPackageId struct {
	Source           string // e.g., "SNOWFLAKE" for first-party
	ScannerPackageId string
}

// String returns the string representation of the scanner package ID.
func (id TrustCenterScannerPackageId) String() string {
	return id.Source + "/" + id.ScannerPackageId
}

// TrustCenterScannerId represents the identifier for a scanner resource.
type TrustCenterScannerId struct {
	Source           string
	ScannerPackageId string
	ScannerId        string
}

// String returns the string representation of the scanner ID.
func (id TrustCenterScannerId) String() string {
	return id.Source + "/" + id.ScannerPackageId + "/" + id.ScannerId
}
