// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type DatabaseRoleModel struct {
	Comment  tfconfig.Variable `json:"comment,omitempty"`
	Database tfconfig.Variable `json:"database,omitempty"`
	Name     tfconfig.Variable `json:"name,omitempty"`
	Owner    tfconfig.Variable `json:"owner,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func DatabaseRole(
	datasourceName string,
	database string,
	name string,
) *DatabaseRoleModel {
	d := &DatabaseRoleModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.DatabaseRole)}
	d.WithDatabase(database)
	d.WithName(name)
	return d
}

func DatabaseRoleWithDefaultMeta(
	database string,
	name string,
) *DatabaseRoleModel {
	d := &DatabaseRoleModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.DatabaseRole)}
	d.WithDatabase(database)
	d.WithName(name)
	return d
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (d *DatabaseRoleModel) MarshalJSON() ([]byte, error) {
	type Alias DatabaseRoleModel
	return json.Marshal(&struct {
		*Alias
		DependsOn []string `json:"depends_on,omitempty"`
	}{
		Alias:     (*Alias)(d),
		DependsOn: d.DependsOn(),
	})
}

func (d *DatabaseRoleModel) WithDependsOn(values ...string) *DatabaseRoleModel {
	d.SetDependsOn(values...)
	return d
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (d *DatabaseRoleModel) WithComment(comment string) *DatabaseRoleModel {
	d.Comment = tfconfig.StringVariable(comment)
	return d
}

func (d *DatabaseRoleModel) WithDatabase(database string) *DatabaseRoleModel {
	d.Database = tfconfig.StringVariable(database)
	return d
}

func (d *DatabaseRoleModel) WithName(name string) *DatabaseRoleModel {
	d.Name = tfconfig.StringVariable(name)
	return d
}

func (d *DatabaseRoleModel) WithOwner(owner string) *DatabaseRoleModel {
	d.Owner = tfconfig.StringVariable(owner)
	return d
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (d *DatabaseRoleModel) WithCommentValue(value tfconfig.Variable) *DatabaseRoleModel {
	d.Comment = value
	return d
}

func (d *DatabaseRoleModel) WithDatabaseValue(value tfconfig.Variable) *DatabaseRoleModel {
	d.Database = value
	return d
}

func (d *DatabaseRoleModel) WithNameValue(value tfconfig.Variable) *DatabaseRoleModel {
	d.Name = value
	return d
}

func (d *DatabaseRoleModel) WithOwnerValue(value tfconfig.Variable) *DatabaseRoleModel {
	d.Owner = value
	return d
}
