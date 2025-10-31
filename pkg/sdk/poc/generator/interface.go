package generator

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

type objectIdentifierKind string

const (
	AccountObjectIdentifier             objectIdentifierKind = "AccountObjectIdentifier"
	DatabaseObjectIdentifier            objectIdentifierKind = "DatabaseObjectIdentifier"
	SchemaObjectIdentifier              objectIdentifierKind = "SchemaObjectIdentifier"
	SchemaObjectIdentifierWithArguments objectIdentifierKind = "SchemaObjectIdentifierWithArguments"
)

func ToObjectIdentifierKind(s string) (objectIdentifierKind, error) {
	switch s {
	case "AccountObjectIdentifier":
		return AccountObjectIdentifier, nil
	case "DatabaseObjectIdentifier":
		return DatabaseObjectIdentifier, nil
	case "SchemaObjectIdentifier":
		return SchemaObjectIdentifier, nil
	case "SchemaObjectIdentifierWithArguments":
		return SchemaObjectIdentifierWithArguments, nil
	default:
		return "", fmt.Errorf("invalid string identifier type: %s", s)
	}
}

// Interface groups operations for particular object or objects family (e.g. DATABASE ROLE)
type Interface struct {
	// Name is the interface's name, e.g. "DatabaseRoles"
	Name string
	// NameSingular is the prefix/suffix which can be used to create other structs and methods, e.g. "DatabaseRole"
	NameSingular string
	// Operations contains all operations for given interface
	Operations []*Operation
	// IdentifierKind keeps identifier of the underlying object (e.g. DatabaseObjectIdentifier)
	IdentifierKind string

	*genhelpers.PreambleModel
}

func (i *Interface) ObjectName() string {
	return i.Name
}

func NewInterface(name string, nameSingular string, identifierKind string, operations ...*Operation) *Interface {
	return &Interface{
		Name:           name,
		NameSingular:   nameSingular,
		IdentifierKind: identifierKind,
		Operations:     operations,
	}
}

// NameLowerCased returns interface name starting with a lower case letter
func (i *Interface) NameLowerCased() string {
	return startingWithLowerCase(i.Name)
}

// ObjectIdentifierKind returns the level of the object identifier (e.g. for DatabaseObjectIdentifier, it returns the prefix "Database")
func (i *Interface) ObjectIdentifierPrefix() idPrefix {
	return identifierStringToPrefix(i.IdentifierKind)
}
