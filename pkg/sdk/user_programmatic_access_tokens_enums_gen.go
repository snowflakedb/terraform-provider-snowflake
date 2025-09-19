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

var AllProgrammaticAccessTokenStatuses = []ProgrammaticAccessTokenStatus{
	ProgrammaticAccessTokenStatusActive,
	ProgrammaticAccessTokenStatusExpired,
	ProgrammaticAccessTokenStatusDisabled,
}

func ToProgrammaticAccessTokenStatus(s string) (ProgrammaticAccessTokenStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllProgrammaticAccessTokenStatuses, ProgrammaticAccessTokenStatus(s)) {
		return "", fmt.Errorf("invalid ProgrammaticAccessTokenStatus: %s", s)
	}
	return ProgrammaticAccessTokenStatus(s), nil
}
