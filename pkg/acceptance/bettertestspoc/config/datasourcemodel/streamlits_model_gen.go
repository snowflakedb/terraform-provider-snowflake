// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type StreamlitsModel struct {
	In           tfconfig.Variable `json:"in,omitempty"`
	Like         tfconfig.Variable `json:"like,omitempty"`
	Limit        tfconfig.Variable `json:"limit,omitempty"`
	Streamlits   tfconfig.Variable `json:"streamlits,omitempty"`
	WithDescribe tfconfig.Variable `json:"with_describe,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Streamlits(
	datasourceName string,
) *StreamlitsModel {
	s := &StreamlitsModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.Streamlits)}
	return s
}

func StreamlitsWithDefaultMeta() *StreamlitsModel {
	s := &StreamlitsModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.Streamlits)}
	return s
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (s *StreamlitsModel) MarshalJSON() ([]byte, error) {
	type Alias StreamlitsModel
	return json.Marshal(&struct {
		*Alias
		DependsOn                 []string                      `json:"depends_on,omitempty"`
		SingleAttributeWorkaround config.ReplacementPlaceholder `json:"single_attribute_workaround,omitempty"`
	}{
		Alias:                     (*Alias)(s),
		DependsOn:                 s.DependsOn(),
		SingleAttributeWorkaround: config.SnowflakeProviderConfigSingleAttributeWorkaround,
	})
}

func (s *StreamlitsModel) WithDependsOn(values ...string) *StreamlitsModel {
	s.SetDependsOn(values...)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// in attribute type is not yet supported, so WithIn can't be generated

func (s *StreamlitsModel) WithLike(like string) *StreamlitsModel {
	s.Like = tfconfig.StringVariable(like)
	return s
}

// limit attribute type is not yet supported, so WithLimit can't be generated

// streamlits attribute type is not yet supported, so WithStreamlits can't be generated

func (s *StreamlitsModel) WithWithDescribe(withDescribe bool) *StreamlitsModel {
	s.WithDescribe = tfconfig.BoolVariable(withDescribe)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *StreamlitsModel) WithInValue(value tfconfig.Variable) *StreamlitsModel {
	s.In = value
	return s
}

func (s *StreamlitsModel) WithLikeValue(value tfconfig.Variable) *StreamlitsModel {
	s.Like = value
	return s
}

func (s *StreamlitsModel) WithLimitValue(value tfconfig.Variable) *StreamlitsModel {
	s.Limit = value
	return s
}

func (s *StreamlitsModel) WithStreamlitsValue(value tfconfig.Variable) *StreamlitsModel {
	s.Streamlits = value
	return s
}

func (s *StreamlitsModel) WithWithDescribeValue(value tfconfig.Variable) *StreamlitsModel {
	s.WithDescribe = value
	return s
}
