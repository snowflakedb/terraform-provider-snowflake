package sdk

// Validates whether All isn't a non-nil pointer to false.
func (s *SessionPolicySecondaryRoles) validate() error {
	if s == nil {
		return nil
	}
	if s.All != nil && !*s.All {
		return errInvalidValue("SessionPolicySecondaryRoles", "All", "false")
	}
	return nil
}
