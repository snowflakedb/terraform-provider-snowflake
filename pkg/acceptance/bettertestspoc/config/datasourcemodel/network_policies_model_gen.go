// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type NetworkPoliciesModel struct {
	Like            tfconfig.Variable `json:"like,omitempty"`
	NetworkPolicies tfconfig.Variable `json:"network_policies,omitempty"`
	WithDescribe    tfconfig.Variable `json:"with_describe,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func NetworkPolicies(
	datasourceName string,
) *NetworkPoliciesModel {
	n := &NetworkPoliciesModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.NetworkPolicies)}
	return n
}

func NetworkPoliciesWithDefaultMeta() *NetworkPoliciesModel {
	n := &NetworkPoliciesModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.NetworkPolicies)}
	return n
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (n *NetworkPoliciesModel) MarshalJSON() ([]byte, error) {
	type Alias NetworkPoliciesModel
	return json.Marshal(&struct {
		*Alias
		DependsOn                 []string                      `json:"depends_on,omitempty"`
		SingleAttributeWorkaround config.ReplacementPlaceholder `json:"single_attribute_workaround,omitempty"`
	}{
		Alias:                     (*Alias)(n),
		DependsOn:                 n.DependsOn(),
		SingleAttributeWorkaround: config.SnowflakeProviderConfigSingleAttributeWorkaround,
	})
}

func (n *NetworkPoliciesModel) WithDependsOn(values ...string) *NetworkPoliciesModel {
	n.SetDependsOn(values...)
	return n
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (n *NetworkPoliciesModel) WithLike(like string) *NetworkPoliciesModel {
	n.Like = tfconfig.StringVariable(like)
	return n
}

// network_policies attribute type is not yet supported, so WithNetworkPolicies can't be generated

func (n *NetworkPoliciesModel) WithWithDescribe(withDescribe bool) *NetworkPoliciesModel {
	n.WithDescribe = tfconfig.BoolVariable(withDescribe)
	return n
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (n *NetworkPoliciesModel) WithLikeValue(value tfconfig.Variable) *NetworkPoliciesModel {
	n.Like = value
	return n
}

func (n *NetworkPoliciesModel) WithNetworkPoliciesValue(value tfconfig.Variable) *NetworkPoliciesModel {
	n.NetworkPolicies = value
	return n
}

func (n *NetworkPoliciesModel) WithWithDescribeValue(value tfconfig.Variable) *NetworkPoliciesModel {
	n.WithDescribe = value
	return n
}
