package main

import (
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[sdk.Grant] = new(GrantCsvRow)

type GrantCsvRow struct {
	Privilege   string `csv:"privilege"`
	GrantedOn   string `csv:"granted_on"`
	GrantOn     string `csv:"grant_on"`
	Name        string `csv:"name"`
	GrantedTo   string `csv:"granted_to"`
	GrantTo     string `csv:"grant_to"`
	GranteeName string `csv:"grantee_name"`
	GrantOption bool   `csv:"grant_option"`
	GrantedBy   string `csv:"granted_by"`
}

// convert is a copy-paste of the sdk.grantRow convert implementation
func (row GrantCsvRow) convert() (*sdk.Grant, error) {
	grantedTo := sdk.ObjectType(strings.ReplaceAll(row.GrantedTo, "_", " "))
	grantTo := sdk.ObjectType(strings.ReplaceAll(row.GrantTo, "_", " "))
	var grantedOn sdk.ObjectType
	// true for current grants
	if row.GrantedOn != "" {
		grantedOn = sdk.ObjectType(strings.ReplaceAll(row.GrantedOn, "_", " "))
	}
	if row.GrantedOn == "VOLUME" {
		grantedOn = sdk.ObjectTypeExternalVolume
	}
	if row.GrantedOn == "MODULE" {
		grantedOn = sdk.ObjectTypeModel
	}

	var grantOn sdk.ObjectType
	// true for future grants
	if row.GrantOn != "" {
		grantOn = sdk.ObjectType(strings.ReplaceAll(row.GrantOn, "_", " "))
	}
	if row.GrantOn == "VOLUME" {
		grantOn = sdk.ObjectTypeExternalVolume
	}
	if row.GrantOn == "MODULE" {
		grantOn = sdk.ObjectTypeModel
	}

	var name sdk.ObjectIdentifier
	var err error
	// TODO(SNOW-1569535): use a mapper from object type to parsing function
	if sdk.ObjectType(row.GrantedOn).IsWithArguments() {
		name, err = sdk.ParseSchemaObjectIdentifierWithArgumentsAndReturnType(row.Name)
	} else {
		name, err = sdk.ParseObjectIdentifierString(row.Name)
	}
	if err != nil {
		log.Printf("[DEBUG] Failed to parse identifier [%s], err = \"%s\"; falling back to fully qualified name conversion", row.Name, err)
		name = sdk.NewObjectIdentifierFromFullyQualifiedName(row.Name)
	}

	return &sdk.Grant{
		Privilege:   row.Privilege,
		GrantedOn:   grantedOn,
		GrantOn:     grantOn,
		GrantedTo:   grantedTo,
		GrantTo:     grantTo,
		Name:        name,
		GranteeName: sdk.NewAccountObjectIdentifier(row.GranteeName),
		GrantOption: row.GrantOption,
		GrantedBy:   sdk.NewAccountObjectIdentifier(row.GrantedBy),
	}, nil
}
