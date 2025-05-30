package snowflake

import (
	"fmt"
	"strings"
)

type Identifier interface {
	QualifiedName() string
}

type SchemaIdentifier struct {
	Database string
	Schema   string
}

func (i *SchemaIdentifier) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"`, i.Database, i.Schema)
}

type SchemaObjectIdentifier struct {
	Database   string
	Schema     string
	ObjectName string `db:"NAME"`
}

func (i *SchemaObjectIdentifier) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"`, i.Database, i.Schema, i.ObjectName)
}

func SchemaObjectIdentifierFromQualifiedName(name string) *SchemaObjectIdentifier {
	parts := strings.Split(name, ".")
	return &SchemaObjectIdentifier{
		Database:   strings.Trim(parts[0], `"`),
		Schema:     strings.Trim(parts[1], `"`),
		ObjectName: strings.Trim(parts[2], `"`),
	}
}

type ColumnIdentifier struct {
	Database   string
	Schema     string
	ObjectName string `db:"NAME"`
	Column     string
}

func (i *ColumnIdentifier) QualifiedName() string {
	return fmt.Sprintf(`"%v"."%v"."%v"."%v"`, i.Database, i.Schema, i.ObjectName, i.Column)
}

func ColumnIdentifierFromQualifiedName(name string) *ColumnIdentifier {
	parts := strings.Split(name, ".")
	return &ColumnIdentifier{
		Database:   strings.Trim(parts[0], `"`),
		Schema:     strings.Trim(parts[1], `"`),
		ObjectName: strings.Trim(parts[2], `"`),
		Column:     strings.Trim(parts[3], `"`),
	}
}
