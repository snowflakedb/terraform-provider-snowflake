package main

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[DatabaseRoleRepresentation] = new(DatabaseRoleCsvRow)

type DatabaseRoleCsvRow struct {
	Comment                string `csv:"comment"`
	CreatedOn              string `csv:"created_on"`
	DatabaseName           string `csv:"database_name"`
	GrantedDatabaseRoles   string `csv:"granted_database_roles"`
	GrantedToDatabaseRoles string `csv:"granted_to_database_roles"`
	GrantedToRoles         string `csv:"granted_to_roles"`
	IsCurrent              string `csv:"is_current"`
	IsDefault              string `csv:"is_default"`
	IsInherited            string `csv:"is_inherited"`
	Name                   string `csv:"name"`
	Owner                  string `csv:"owner"`
	OwnerRoleType          string `csv:"owner_role_type"`
}

type DatabaseRoleRepresentation struct {
	sdk.DatabaseRole
}

func (row DatabaseRoleCsvRow) convert() (*DatabaseRoleRepresentation, error) {
	databaseRoleRepresentation := &DatabaseRoleRepresentation{
		DatabaseRole: sdk.DatabaseRole{
			Name:         row.Name,
			DatabaseName: row.DatabaseName,
			Comment:      row.Comment,
			Owner:        row.Owner,
		},
	}

	return databaseRoleRepresentation, nil
}
