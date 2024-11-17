package config

import (
	"encoding/json"
	"fmt"
)

var DefaultJsonProvider = NewBasicJsonProvider()

type JsonProvider interface {
	ResourceJsonFromModel(model ResourceModel) ([]byte, error)
	ProviderJsonFromModel(model ProviderModel) ([]byte, error)
	// Variable
	// Output
	// Locals
	// Module
	// Terraform
}

type basicJsonProvider struct{}

func NewBasicJsonProvider() JsonProvider {
	return &basicJsonProvider{}
}

func (p *basicJsonProvider) ResourceJsonFromModel(model ResourceModel) ([]byte, error) {
	modelJson := resourceJson{
		Resource: map[string]map[string]ResourceModel{
			fmt.Sprintf("%s", model.Resource()): {
				fmt.Sprintf("%s", model.ResourceName()): model,
			},
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type resourceJson struct {
	Resource map[string]map[string]ResourceModel `json:"resource"`
}

func (p *basicJsonProvider) ProviderJsonFromModel(model ProviderModel) ([]byte, error) {
	modelJson := providerJson{
		Provider: map[string]ProviderModel{
			fmt.Sprintf("%s", model.ProviderName()): model,
		},
	}

	return json.MarshalIndent(modelJson, "", "    ")
}

type providerJson struct {
	Provider map[string]ProviderModel `json:"provider"`
}
