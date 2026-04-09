package sdk

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
