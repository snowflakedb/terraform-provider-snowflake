package snowflake

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type Identifier interface {
	QualifiedName() string
}

type SchemaObjectIdentifier struct {
	Database   string
	Schema     string
	ObjectName string `db:"NAME"`
}

func (i *SchemaObjectIdentifier) QualifiedName() string {
	db := sdk.DoubleQuotes.Modify(i.Database)
	schema := sdk.DoubleQuotes.Modify(i.Schema)
	name := sdk.DoubleQuotes.Modify(i.ObjectName)
	return fmt.Sprintf(`%v.%v.%v`, db, schema, name)
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
	db := sdk.DoubleQuotes.Modify(i.Database)
	schema := sdk.DoubleQuotes.Modify(i.Schema)
	name := sdk.DoubleQuotes.Modify(i.ObjectName)
	column := sdk.DoubleQuotes.Modify(i.Column)
	return fmt.Sprintf(`%v.%v.%v.%v`, db, schema, name, column)
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
