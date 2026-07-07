package sdk

import "encoding/json"

type MaskingPolicyOptions struct {
	ExemptOtherPolicies bool `json:"EXEMPT_OTHER_POLICIES"`
}

func ParseMaskingPolicyOptions(s string) (MaskingPolicyOptions, error) {
	var options MaskingPolicyOptions
	err := json.Unmarshal([]byte(s), &options)
	if err != nil {
		return MaskingPolicyOptions{}, err
	}

	return options, nil
}
