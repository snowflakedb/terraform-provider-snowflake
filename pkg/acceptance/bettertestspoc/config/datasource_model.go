package config

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
)

// DatasourceModel is the base interface all of our datasource config models will implement.
// To allow easy implementation, DatasourceModelMeta can be embedded inside the struct (and the struct will automatically implement it).
// TODO [SNOW-1501905]: consider merging ResourceModel with DatasourceModel (currently the implementation is really similar)
type DatasourceModel interface {
	Datasource() datasources.Datasource
	DatasourceName() string
	SetDatasourceName(name string)
	DatasourceReference() string
	DependsOn() []string
	SetDependsOn(values ...string)
}

type DatasourceModelMeta struct {
	name       string
	datasource datasources.Datasource
	dependsOn  []string
}

func (m *DatasourceModelMeta) Datasource() datasources.Datasource {
	return m.datasource
}

func (m *DatasourceModelMeta) DatasourceName() string {
	return m.name
}

func (m *DatasourceModelMeta) SetDatasourceName(name string) {
	m.name = name
}

func (m *DatasourceModelMeta) DatasourceReference() string {
	return fmt.Sprintf(`data.%s.%s`, m.datasource, m.name)
}

func (m *DatasourceModelMeta) DependsOn() []string {
	return m.dependsOn
}

func (m *DatasourceModelMeta) SetDependsOn(values ...string) {
	m.dependsOn = values
}

// DefaultDatasourceName is exported to allow assertions against the resources using the default name.
const DefaultDatasourceName = "test"

func DatasourceDefaultMeta(datasource datasources.Datasource) *DatasourceModelMeta {
	return &DatasourceModelMeta{name: DefaultDatasourceName, datasource: datasource}
}

func DatasourceMeta(resourceName string, datasource datasources.Datasource) *DatasourceModelMeta {
	return &DatasourceModelMeta{name: resourceName, datasource: datasource}
}
