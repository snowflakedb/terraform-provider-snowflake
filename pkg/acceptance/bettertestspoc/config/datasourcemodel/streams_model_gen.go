// Code generated by config model builder generator; DO NOT EDIT.

package datasourcemodel

import (
	"encoding/json"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
)

type StreamsModel struct {
	In           tfconfig.Variable `json:"in,omitempty"`
	Like         tfconfig.Variable `json:"like,omitempty"`
	Limit        tfconfig.Variable `json:"limit,omitempty"`
	StartsWith   tfconfig.Variable `json:"starts_with,omitempty"`
	Streams      tfconfig.Variable `json:"streams,omitempty"`
	WithDescribe tfconfig.Variable `json:"with_describe,omitempty"`

	*config.DatasourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func Streams(
	datasourceName string,
) *StreamsModel {
	s := &StreamsModel{DatasourceModelMeta: config.DatasourceMeta(datasourceName, datasources.Streams)}
	return s
}

func StreamsWithDefaultMeta() *StreamsModel {
	s := &StreamsModel{DatasourceModelMeta: config.DatasourceDefaultMeta(datasources.Streams)}
	return s
}

///////////////////////////////////////////////////////
// set proper json marshalling and handle depends on //
///////////////////////////////////////////////////////

func (s *StreamsModel) MarshalJSON() ([]byte, error) {
	type Alias StreamsModel
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

func (s *StreamsModel) WithDependsOn(values ...string) *StreamsModel {
	s.SetDependsOn(values...)
	return s
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

// in attribute type is not yet supported, so WithIn can't be generated

func (s *StreamsModel) WithLike(like string) *StreamsModel {
	s.Like = tfconfig.StringVariable(like)
	return s
}

// limit attribute type is not yet supported, so WithLimit can't be generated

func (s *StreamsModel) WithStartsWith(startsWith string) *StreamsModel {
	s.StartsWith = tfconfig.StringVariable(startsWith)
	return s
}

// streams attribute type is not yet supported, so WithStreams can't be generated

func (s *StreamsModel) WithWithDescribe(withDescribe bool) *StreamsModel {
	s.WithDescribe = tfconfig.BoolVariable(withDescribe)
	return s
}

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (s *StreamsModel) WithInValue(value tfconfig.Variable) *StreamsModel {
	s.In = value
	return s
}

func (s *StreamsModel) WithLikeValue(value tfconfig.Variable) *StreamsModel {
	s.Like = value
	return s
}

func (s *StreamsModel) WithLimitValue(value tfconfig.Variable) *StreamsModel {
	s.Limit = value
	return s
}

func (s *StreamsModel) WithStartsWithValue(value tfconfig.Variable) *StreamsModel {
	s.StartsWith = value
	return s
}

func (s *StreamsModel) WithStreamsValue(value tfconfig.Variable) *StreamsModel {
	s.Streams = value
	return s
}

func (s *StreamsModel) WithWithDescribeValue(value tfconfig.Variable) *StreamsModel {
	s.WithDescribe = value
	return s
}
