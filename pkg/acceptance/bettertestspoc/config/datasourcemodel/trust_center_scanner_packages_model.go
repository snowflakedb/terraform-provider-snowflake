package datasourcemodel

import (
	"encoding/json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

type TrustCenterScannerPackagesModel struct {
	Like            tfconfig.Variable `json:"like,omitempty"`
	ScannerPackages tfconfig.Variable `json:"scanner_packages,omitempty"`

	*config.DatasourceModelMeta
}

func TrustCenterScannerPackages(
	datasourceName string,
) *TrustCenterScannerPackagesModel {
	t := &TrustCenterScannerPackagesModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.TrustCenterScannerPackages)}
	return t
}

func TrustCenterScannerPackagesWithDefaultMeta() *TrustCenterScannerPackagesModel {
	t := &TrustCenterScannerPackagesModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.TrustCenterScannerPackages)}
	return t
}

func (t *TrustCenterScannerPackagesModel) MarshalJSON() ([]byte, error) {
	type Alias TrustCenterScannerPackagesModel
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

func (t *TrustCenterScannerPackagesModel) WithDependsOn(values ...string) *TrustCenterScannerPackagesModel {
	t.SetDependsOn(values...)
	return t
}

func (t *TrustCenterScannerPackagesModel) WithLike(like string) *TrustCenterScannerPackagesModel {
	t.Like = tfconfig.StringVariable(like)
	return t
}

func (t *TrustCenterScannerPackagesModel) WithLikeValue(value tfconfig.Variable) *TrustCenterScannerPackagesModel {
	t.Like = value
	return t
}

func (t *TrustCenterScannerPackagesModel) WithScannerPackagesValue(value tfconfig.Variable) *TrustCenterScannerPackagesModel {
	t.ScannerPackages = value
	return t
}
