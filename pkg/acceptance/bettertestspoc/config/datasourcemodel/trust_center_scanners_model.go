package datasourcemodel

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type TrustCenterScannersModel struct {
	ScannerPackageId tfconfig.Variable `json:"scanner_package_id,omitempty"`
	Like             tfconfig.Variable `json:"like,omitempty"`
	Scanners         tfconfig.Variable `json:"scanners,omitempty"`

	*config.DatasourceModelMeta
}

func TrustCenterScanners(
	datasourceName string,
) *TrustCenterScannersModel {
	t := &TrustCenterScannersModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.TrustCenterScanners)}
	return t
}

func TrustCenterScannersWithDefaultMeta() *TrustCenterScannersModel {
	t := &TrustCenterScannersModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.TrustCenterScanners)}
	return t
}

func (t *TrustCenterScannersModel) MarshalJSON() ([]byte, error) {
	type Alias TrustCenterScannersModel
	return json.Marshal(&struct {
		*Alias
		DependsOn                 []string                      `json:"depends_on,omitempty"`
		SingleAttributeWorkaround config.ReplacementPlaceholder `json:"single_attribute_workaround,omitempty"`
	}{
		Alias:                     (*Alias)(t),
		DependsOn:                 t.DependsOn(),
		SingleAttributeWorkaround: config.SnowflakeProviderConfigSingleAttributeWorkaround,
	})
}

func (t *TrustCenterScannersModel) WithDependsOn(values ...string) *TrustCenterScannersModel {
	t.SetDependsOn(values...)
	return t
}

func (t *TrustCenterScannersModel) WithScannerPackageId(scannerPackageId string) *TrustCenterScannersModel {
	t.ScannerPackageId = tfconfig.StringVariable(scannerPackageId)
	return t
}

func (t *TrustCenterScannersModel) WithLike(like string) *TrustCenterScannersModel {
	t.Like = tfconfig.StringVariable(like)
	return t
}

func (t *TrustCenterScannersModel) WithScannerPackageIdValue(value tfconfig.Variable) *TrustCenterScannersModel {
	t.ScannerPackageId = value
	return t
}

func (t *TrustCenterScannersModel) WithLikeValue(value tfconfig.Variable) *TrustCenterScannersModel {
	t.Like = value
	return t
}

func (t *TrustCenterScannersModel) WithScannersValue(value tfconfig.Variable) *TrustCenterScannersModel {
	t.Scanners = value
	return t
}
