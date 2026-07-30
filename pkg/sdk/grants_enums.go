package sdk

import (
	"fmt"
	"strings"
)

type GrantInheritedFrom string

const (
	GrantInheritedFromAccount  GrantInheritedFrom = "ACCOUNT"
	GrantInheritedFromDatabase GrantInheritedFrom = "DATABASE"
	GrantInheritedFromSchema   GrantInheritedFrom = "SCHEMA"
)

var AllGrantInheritedFroms = []GrantInheritedFrom{
	GrantInheritedFromAccount,
	GrantInheritedFromDatabase,
	GrantInheritedFromSchema,
}

func ToGrantInheritedFrom(s string) (GrantInheritedFrom, error) {
	s = strings.ToUpper(s)
	switch s {
	case string(GrantInheritedFromAccount):
		return GrantInheritedFromAccount, nil
	case string(GrantInheritedFromDatabase):
		return GrantInheritedFromDatabase, nil
	case string(GrantInheritedFromSchema):
		return GrantInheritedFromSchema, nil
	default:
		return "", fmt.Errorf("invalid grant inherited from: %s", s)
	}
}
