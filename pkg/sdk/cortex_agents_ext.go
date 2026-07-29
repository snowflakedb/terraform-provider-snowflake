package sdk

import (
	"encoding/json"
	"strings"

	"github.com/goccy/go-yaml"
)

func (r *CreateCortexAgentRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type CortexAgentProfile struct {
	DisplayName *string `json:"display_name,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// UnmarshalCortexAgentProfile parses profile JSON into CortexAgentProfile.
func UnmarshalCortexAgentProfile(profileAsJson string) (CortexAgentProfile, error) {
	var profile CortexAgentProfile

	if err := json.Unmarshal([]byte(profileAsJson), &profile); err != nil {
		return profile, err
	}

	return profile, nil
}

// MarshalCortexAgentProfile serializes a CortexAgentProfile to JSON.
func MarshalCortexAgentProfile(profile CortexAgentProfile) (string, error) {
	b, err := json.Marshal(profile)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// NormalizeCortexAgentSpecification parses YAML or JSON agent specifications into a canonical
// JSON string (sorted object keys) so Terraform can compare user YAML with Snowflake JSON responses
// without spurious formatting diffs.
func NormalizeCortexAgentSpecification(spec string) (string, error) {
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
