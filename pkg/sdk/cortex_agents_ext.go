package sdk

import (
	"encoding/json"
	"strings"
)

func (r *CreateCortexAgentRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type CortexAgentProfile struct {
	DisplayName *string `json:"display_name"`
	Avatar      *string `json:"avatar"`
	Color       *string `json:"color"`
}

// UnmarshalCortexAgentProfile parses profile JSON into CortexAgentProfile.
func UnmarshalCortexAgentProfile(profileAsJson string) (*CortexAgentProfile, error) {
	var profile CortexAgentProfile

	if err := json.Unmarshal([]byte(profileAsJson), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// UnmarshalCortexAgentSpec parses agent specification JSON (as returned by DESCRIBE AGENT) into a map.
func UnmarshalCortexAgentSpec(agentSpec string) (map[string]any, error) {
	data := strings.TrimSpace(agentSpec)
	if data == "" {
		return map[string]any{}, nil
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, err
	}
	if m == nil {
		return map[string]any{}, nil
	}
	return m, nil
}
