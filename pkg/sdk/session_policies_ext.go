package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
)

func (r showSessionPolicyDBRow) additionalConvert(result *SessionPolicy) error {
	targetScopes, err := parseSessionPolicyTargetScopes(r.Options)
	if err != nil {
		return fmt.Errorf("parsing target scopes for session policy with name %s: %w", r.Name, err)
	}
	result.TargetScopes = targetScopes
	return nil
}

// parseSessionPolicyTargetScopes extracts the "target_scopes" JSON array from the raw options column.
func parseSessionPolicyTargetScopes(options string) ([]SessionPolicyTargetScope, error) {
	if options == "" {
		return nil, nil
	}
	var raw struct {
		TargetScopes []SessionPolicyTargetScope `json:"target_scopes"`
	}
	if err := json.Unmarshal([]byte(options), &raw); err != nil {
		return nil, err
	}
	slices.Sort(raw.TargetScopes)
	return raw.TargetScopes, nil
}

// Validates whether both All and None aren't non-nil pointers to false.
func (s *SessionPolicySecondaryRoles) validate() error {
	if s == nil {
		return nil
	}
	if s.All != nil && !*s.All {
		return errInvalidValue("SessionPolicySecondaryRoles", "All", "false")
	}
	if s.None != nil && !*s.None {
		return errInvalidValue("SessionPolicySecondaryRoles", "None", "false")
	}
	return nil
}

func (d *SessionPolicyDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

func (s *sessionPolicies) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*SessionPolicyDetails, error) {
	properties, err := s.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseSessionPolicyProperties(properties, id)
}

func parseSessionPolicyProperties(properties []SessionPolicyProperty, id SchemaObjectIdentifier) (*SessionPolicyDetails, error) {
	details := &SessionPolicyDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Property {
		case "OWNER":
			details.Owner = prop.Value
		case "OWNER_ROLE_TYPE":
			details.OwnerRoleType = prop.Value
		case "COMMENT":
			details.Comment = emptyIfNull(prop.Value)
		case "SESSION_IDLE_TIMEOUT_MINS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.SessionIdleTimeoutMins = int(val)
			}
		case "SESSION_UI_IDLE_TIMEOUT_MINS":
			if val, err := strconv.ParseInt(prop.Value, 10, 32); err != nil {
				errs = append(errs, err)
			} else {
				details.SessionUiIdleTimeoutMins = int(val)
			}
		case "ALLOWED_SECONDARY_ROLES":
			details.AllowedSecondaryRoles = ParseCommaSeparatedStringArray(prop.Value, false)
		case "BLOCKED_SECONDARY_ROLES":
			details.BlockedSecondaryRoles = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}
	return details, errors.Join(errs...)
}
