package config

import (
	"github.com/pelletier/go-toml/v2"
)

var DefaultTomlConfigProvider = NewBasicTomlConfigProvider()

// TomlConfigProvider defines methods to generate TOML configs.
type TomlConfigProvider interface {
	ProviderTomlFromModel(model ProviderModel) ([]byte, error)
}

type basicTomlConfigProvider struct{}

func NewBasicTomlConfigProvider() TomlConfigProvider {
	return &basicTomlConfigProvider{}
}

func (p *basicTomlConfigProvider) ProviderTomlFromModel(model ProviderModel) ([]byte, error) {
	modelToml := map[string]ProviderModel{
		model.ProviderName(): model,
	}

	return toml.Marshal(modelToml)
}
