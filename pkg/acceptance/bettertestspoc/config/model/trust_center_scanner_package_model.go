package model

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type TrustCenterScannerPackageModel struct {
	ScannerPackageId tfconfig.Variable `json:"scanner_package_id,omitempty"`
	Enabled          tfconfig.Variable `json:"enabled,omitempty"`
	Schedule         tfconfig.Variable `json:"schedule,omitempty"`
	Notification     tfconfig.Variable `json:"notification,omitempty"`
	ShowOutput       tfconfig.Variable `json:"show_output,omitempty"`

	DynamicBlock *config.DynamicBlock `json:"dynamic,omitempty"`

	*config.ResourceModelMeta
}

func TrustCenterScannerPackage(
	resourceName string,
	scannerPackageId string,
) *TrustCenterScannerPackageModel {
	t := &TrustCenterScannerPackageModel{ResourceModelMeta: config.Meta(resourceName, resources.TrustCenterScannerPackage)}
	t.WithScannerPackageId(scannerPackageId)
	return t
}

func TrustCenterScannerPackageWithDefaultMeta(
	scannerPackageId string,
) *TrustCenterScannerPackageModel {
	t := &TrustCenterScannerPackageModel{ResourceModelMeta: config.DefaultMeta(resources.TrustCenterScannerPackage)}
	t.WithScannerPackageId(scannerPackageId)
	return t
}

func (t *TrustCenterScannerPackageModel) MarshalJSON() ([]byte, error) {
	type Alias TrustCenterScannerPackageModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(t),
		DependsOn: t.DependsOn(),
	})
}

func (t *TrustCenterScannerPackageModel) WithDependsOn(values ...string) *TrustCenterScannerPackageModel {
	t.SetDependsOn(values...)
	return t
}

func (t *TrustCenterScannerPackageModel) WithScannerPackageId(scannerPackageId string) *TrustCenterScannerPackageModel {
	t.ScannerPackageId = tfconfig.StringVariable(scannerPackageId)
	return t
}

func (t *TrustCenterScannerPackageModel) WithEnabled(enabled bool) *TrustCenterScannerPackageModel {
	t.Enabled = tfconfig.BoolVariable(enabled)
	return t
}

func (t *TrustCenterScannerPackageModel) WithSchedule(schedule string) *TrustCenterScannerPackageModel {
	t.Schedule = tfconfig.StringVariable(schedule)
	return t
}

func (t *TrustCenterScannerPackageModel) WithScannerPackageIdValue(value tfconfig.Variable) *TrustCenterScannerPackageModel {
	t.ScannerPackageId = value
	return t
}

func (t *TrustCenterScannerPackageModel) WithEnabledValue(value tfconfig.Variable) *TrustCenterScannerPackageModel {
	t.Enabled = value
	return t
}

func (t *TrustCenterScannerPackageModel) WithScheduleValue(value tfconfig.Variable) *TrustCenterScannerPackageModel {
	t.Schedule = value
	return t
}

func (t *TrustCenterScannerPackageModel) WithNotificationValue(value tfconfig.Variable) *TrustCenterScannerPackageModel {
	t.Notification = value
	return t
}
