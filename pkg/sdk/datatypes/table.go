package datatypes

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// TableDataType is based on https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#returning-tabular-data.
// It does not have synonyms.
// It consists of a list of columns; may be empty.
// Each column is either a named column ("arg_name NUMBER") used in return types, or a type-only column ("DATE")
// used in abbreviated argument signatures (e.g. grants and SHOW output for functions with TABLE arguments).
type TableDataType struct {
	columns        []TableDataTypeColumn
	underlyingType string
}

type TableDataTypeColumn struct {
	name     string
	dataType DataType
}

var TableDataTypeSynonyms = []string{"TABLE"}

func (c *TableDataTypeColumn) ColumnName() string {
	return c.name
}

func (c *TableDataTypeColumn) ColumnType() DataType {
	return c.dataType
}

func (c *TableDataTypeColumn) format(formatType func(DataType) string) string {
	if c.name == "" {
		return formatType(c.dataType)
	}
	return fmt.Sprintf("%s %s", c.name, formatType(c.dataType))
}

func (c *TableDataTypeColumn) ToSql() string {
	return c.format(func(dt DataType) string { return dt.ToSql() })
}

func (c *TableDataTypeColumn) ToLegacyDataTypeSql() string {
	return c.format(func(dt DataType) string { return dt.ToLegacyDataTypeSql() })
}

func (c *TableDataTypeColumn) Canonical() string {
	return c.format(func(dt DataType) string { return dt.Canonical() })
}

func (c *TableDataTypeColumn) ToSqlWithoutUnknowns() string {
	return c.format(func(dt DataType) string { return dt.ToSqlWithoutUnknowns() })
}

func (t *TableDataType) ToSql() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return col.ToSql()
	}), ", ")
	return fmt.Sprintf("%s(%s)", t.underlyingType, columns)
}

func (t *TableDataType) ToLegacyDataTypeSql() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return col.ToLegacyDataTypeSql()
	}), ", ")
	return fmt.Sprintf("%s(%s)", TableLegacyDataType, columns)
}

func (t *TableDataType) Canonical() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return col.Canonical()
	}), ", ")
	return fmt.Sprintf("%s(%s)", TableLegacyDataType, columns)
}

func (t *TableDataType) ToSqlWithoutUnknowns() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return col.ToSqlWithoutUnknowns()
	}), ", ")
	return fmt.Sprintf("%s(%s)", t.underlyingType, columns)
}

func (t *TableDataType) Columns() []TableDataTypeColumn {
	return t.columns
}

// splitColumnDefs splits a comma-separated column definition string,
// respecting nested parentheses (e.g. NUMBER(38,0) is kept intact).
func splitColumnDefs(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, ch := range s {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func parseTableDataTypeRaw(raw sanitizedDataTypeRaw) (*TableDataType, error) {
	r := strings.TrimSpace(strings.TrimPrefix(raw.raw, raw.matchedByType))
	if r == "" || (!strings.HasPrefix(r, "(") || !strings.HasSuffix(r, ")")) {
		return nil, fmt.Errorf(`table %s could not be parsed, use "%s(argName argType, ...)" format`, raw.raw, raw.matchedByType)
	}
	onlyArgs := strings.TrimSpace(r[1 : len(r)-1])
	if onlyArgs == "" {
		return &TableDataType{
			columns:        make([]TableDataTypeColumn, 0),
			underlyingType: raw.matchedByType,
		}, nil
	}
	columns, err := collections.MapErr(splitColumnDefs(onlyArgs), func(arg string) (TableDataTypeColumn, error) {
		trimmed := strings.TrimSpace(arg)
		if argDataType, err := ParseDataType(trimmed); err == nil {
			return TableDataTypeColumn{
				name:     "",
				dataType: argDataType,
			}, nil
		}
		argParts := strings.SplitN(trimmed, " ", 2)
		if len(argParts) != 2 {
			return TableDataTypeColumn{}, fmt.Errorf("could not parse table column: %s, it should contain the following format `<arg_name> <arg_type>`; parser failure may be connected to the complex argument names", arg)
		}
		argDataType, err := ParseDataType(argParts[1])
		if err != nil {
			return TableDataTypeColumn{}, err
		}
		return TableDataTypeColumn{
			name:     argParts[0],
			dataType: argDataType,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return &TableDataType{
		columns:        columns,
		underlyingType: raw.matchedByType,
	}, nil
}

func areTableDataTypesTheSame(a, b *TableDataType) bool {
	if len(a.columns) != len(b.columns) {
		return false
	}

	for i := range a.columns {
		aColumn := a.columns[i]
		bColumn := b.columns[i]

		if aColumn.name != bColumn.name || !AreTheSame(aColumn.dataType, bColumn.dataType) {
			return false
		}
	}

	return true
}

// tables are different if:
// - they have different numbers of columns
// - name differs for at least one column
// - data type is different for at least one column
func areTableDataTypesDefinitelyDifferent(a, b *TableDataType) bool {
	if len(a.columns) != len(b.columns) {
		return true
	}

	for i := range a.columns {
		aColumn := a.columns[i]
		bColumn := b.columns[i]

		if aColumn.name != bColumn.name {
			return true
		}
		if AreDefinitelyDifferent(aColumn.dataType, bColumn.dataType) {
			return true
		}
	}

	return false
}
