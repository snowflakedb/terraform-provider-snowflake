package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[AccountRoleRepresentation] = new(AccountRoleCsvRow)

type AccountRoleCsvRow struct {
	AssignedToUsers string `csv:"assigned_to_users"`
	Comment         string `csv:"comment"`
	CreatedOn       string `csv:"created_on"`
	GrantedRoles    string `csv:"granted_roles"`
	GrantedToRoles  string `csv:"granted_to_roles"`
	IsCurrent       string `csv:"is_current"`
	IsDefault       string `csv:"is_default"`
	IsInherited     string `csv:"is_inherited"`
	Name            string `csv:"name"`
	Owner           string `csv:"owner"`
}

type AccountRoleRepresentation struct {
	sdk.Role
}

func (row AccountRoleCsvRow) convert() (*AccountRoleRepresentation, error) {
	accountRoleRepresentation := &AccountRoleRepresentation{
		Role: sdk.Role{
			Name:    row.Name,
			Comment: row.Comment,
			Owner:   row.Owner,
		},
	}

	return accountRoleRepresentation, nil
}
