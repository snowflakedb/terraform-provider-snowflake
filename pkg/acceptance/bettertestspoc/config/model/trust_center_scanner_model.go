package model

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type TrustCenterScannerModel struct {
	ScannerPackageId tfconfig.Variable `json:"scanner_package_id,omitempty"`
	ScannerId        tfconfig.Variable `json:"scanner_id,omitempty"`
	Enabled          tfconfig.Variable `json:"enabled,omitempty"`
	Schedule         tfconfig.Variable `json:"schedule,omitempty"`
	Notification     tfconfig.Variable `json:"notification,omitempty"`
	ShowOutput       tfconfig.Variable `json:"show_output,omitempty"`

	DynamicBlock *config.DynamicBlock `json:"dynamic,omitempty"`

	*config.ResourceModelMeta
}

func TrustCenterScanner(
	resourceName string,
	scannerPackageId string,
	scannerId string,
) *TrustCenterScannerModel {
	t := &TrustCenterScannerModel{ResourceModelMeta: config.Meta(resourceName, resources.TrustCenterScanner)}
	t.WithScannerPackageId(scannerPackageId)
	t.WithScannerId(scannerId)
	return t
}

func TrustCenterScannerWithDefaultMeta(
	scannerPackageId string,
	scannerId string,
) *TrustCenterScannerModel {
	t := &TrustCenterScannerModel{ResourceModelMeta: config.DefaultMeta(resources.TrustCenterScanner)}
	t.WithScannerPackageId(scannerPackageId)
	t.WithScannerId(scannerId)
	return t
}

func (t *TrustCenterScannerModel) MarshalJSON() ([]byte, error) {
	type Alias TrustCenterScannerModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(t),
		DependsOn: t.DependsOn(),
	})
}

func (t *TrustCenterScannerModel) WithDependsOn(values ...string) *TrustCenterScannerModel {
	t.SetDependsOn(values...)
	return t
}

func (t *TrustCenterScannerModel) WithScannerPackageId(scannerPackageId string) *TrustCenterScannerModel {
	t.ScannerPackageId = tfconfig.StringVariable(scannerPackageId)
	return t
}

func (t *TrustCenterScannerModel) WithScannerId(scannerId string) *TrustCenterScannerModel {
	t.ScannerId = tfconfig.StringVariable(scannerId)
	return t
}

func (t *TrustCenterScannerModel) WithEnabled(enabled bool) *TrustCenterScannerModel {
	t.Enabled = tfconfig.BoolVariable(enabled)
	return t
}

func (t *TrustCenterScannerModel) WithSchedule(schedule string) *TrustCenterScannerModel {
	t.Schedule = tfconfig.StringVariable(schedule)
	return t
}

func (t *TrustCenterScannerModel) WithScannerPackageIdValue(value tfconfig.Variable) *TrustCenterScannerModel {
	t.ScannerPackageId = value
	return t
}

func (t *TrustCenterScannerModel) WithScannerIdValue(value tfconfig.Variable) *TrustCenterScannerModel {
	t.ScannerId = value
	return t
}

func (t *TrustCenterScannerModel) WithEnabledValue(value tfconfig.Variable) *TrustCenterScannerModel {
	t.Enabled = value
	return t
}

func (t *TrustCenterScannerModel) WithScheduleValue(value tfconfig.Variable) *TrustCenterScannerModel {
	t.Schedule = value
	return t
}

func (t *TrustCenterScannerModel) WithNotificationValue(value tfconfig.Variable) *TrustCenterScannerModel {
	t.Notification = value
	return t
}
