package sdk

import (
	"encoding/json"

	"github.com/goccy/go-yaml"
)

func (r *CreateMcpServerRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (d *McpServerDetails) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(d.DatabaseName, d.SchemaName, d.Name)
}

// NormalizeMcpServerSpecification parses YAML or JSON MCP server specifications into a canonical
// JSON string so Terraform can compare user YAML with Snowflake JSON responses without spurious
// formatting diffs. Snowflake always persists a version field, so it is removed when present.
func NormalizeMcpServerSpecification(spec string) (string, error) {
	var m map[string]any
	if err := yaml.Unmarshal([]byte(spec), &m); err != nil {
		return "", err
	}
	delete(m, "version")

	// json.Marshal sorts the keys alphabetically for maps
	json, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(json), nil
}
