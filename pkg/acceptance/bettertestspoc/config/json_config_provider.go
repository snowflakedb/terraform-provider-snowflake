package config

import (
	"encoding/json"
)

var DefaultJsonConfigProvider = NewBasicJsonConfigProvider()

// JsonConfigProvider defines methods to generate .tf.json configs.
// TODO [SNOW-1501905]: add config builders for other block types (Variable, Output, Localsl, Module, Terraform)
type JsonConfigProvider interface {
	ResourceJsonFromModel(model ResourceModel) ([]byte, error)
	DatasourceJsonFromModel(model DatasourceModel) ([]byte, error)
	ProviderJsonFromModel(model ProviderModel) ([]byte, error)
	TfCommonJsonFromModel(model VariableModel) ([]byte, error)
}

type basicJsonConfigProvider struct{}

func NewBasicJsonConfigProvider() JsonConfigProvider {
	return &basicJsonConfigProvider{}
}

func (p *basicJsonConfigProvider) ResourceJsonFromModel(model ResourceModel) ([]byte, error) {
	modelJson := resourceJson{
		Resource: map[string]map[string]ResourceModel{
			model.Resource().String(): {
				model.ResourceName(): model,
			},
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type resourceJson struct {
	Resource map[string]map[string]ResourceModel `json:"resource"`
}

func (p *basicJsonConfigProvider) DatasourceJsonFromModel(model DatasourceModel) ([]byte, error) {
	modelJson := datasourceJson{
		Datasource: map[string]map[string]DatasourceModel{
			model.Datasource().String(): {
				model.DatasourceName(): model,
			},
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type datasourceJson struct {
	Datasource map[string]map[string]DatasourceModel `json:"data"`
}

func (p *basicJsonConfigProvider) ProviderJsonFromModel(model ProviderModel) ([]byte, error) {
	modelJson := providerJson{
		Provider: map[string]ProviderModel{
			model.ProviderName(): model,
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type providerJson struct {
	Provider map[string]ProviderModel `json:"provider"`
}

func (p *basicJsonConfigProvider) TfCommonJsonFromModel(model VariableModel) ([]byte, error) {
	modelJson := map[string]map[string]VariableModel{
		model.CommonTfType(): {
			model.CommonTfName(): model,
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}
