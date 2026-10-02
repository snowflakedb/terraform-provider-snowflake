package sdk

import (
	"encoding/json"
	"strings"

	"github.com/goccy/go-yaml"
)

func (r *CreateGatewayRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (d *GatewayDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}

// NormalizeGatewaySpecification parses YAML or JSON gateway specifications into canonical JSON.
func NormalizeGatewaySpecification(spec string) (string, error) {
	data := strings.TrimSpace(spec)
	if data == "" {
		return "{}", nil
	}

	var m map[string]any
	if err := yaml.Unmarshal([]byte(data), &m); err != nil {
		return "", err
	}
	if m == nil {
		return "{}", nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
