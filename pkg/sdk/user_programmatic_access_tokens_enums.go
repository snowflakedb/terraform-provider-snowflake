package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type ProgrammaticAccessTokenStatus string

const (
	ProgrammaticAccessTokenStatusActive   ProgrammaticAccessTokenStatus = "ACTIVE"
	ProgrammaticAccessTokenStatusExpired  ProgrammaticAccessTokenStatus = "EXPIRED"
	ProgrammaticAccessTokenStatusDisabled ProgrammaticAccessTokenStatus = "DISABLED"
)

var allProgrammaticAccessTokenStatuses = []ProgrammaticAccessTokenStatus{
	ProgrammaticAccessTokenStatusActive,
	ProgrammaticAccessTokenStatusExpired,
	ProgrammaticAccessTokenStatusDisabled,
}

func toProgrammaticAccessTokenStatus(s string) (ProgrammaticAccessTokenStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allProgrammaticAccessTokenStatuses, ProgrammaticAccessTokenStatus(s)) {
		return "", fmt.Errorf("invalid programmatic access token status: %s", s)
	}
	return ProgrammaticAccessTokenStatus(s), nil
}
