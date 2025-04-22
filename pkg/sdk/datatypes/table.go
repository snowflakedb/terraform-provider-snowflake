package datatypes

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// TableDataType is based on https://docs.snowflake.com/en/developer-guide/stored-procedure/stored-procedures-java#returning-tabular-data.
// It does not have synonyms.
// It consists of a list of column name + column type; may be empty.
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

func (t *TableDataType) ToSql() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return fmt.Sprintf("%s %s", col.name, col.dataType.ToSql())
	}), ", ")
	return fmt.Sprintf("%s(%s)", t.underlyingType, columns)
}

func (t *TableDataType) ToLegacyDataTypeSql() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return fmt.Sprintf("%s %s", col.name, col.dataType.ToLegacyDataTypeSql())
	}), ", ")
	return fmt.Sprintf("%s(%s)", TableLegacyDataType, columns)
}

func (t *TableDataType) Canonical() string {
	columns := strings.Join(collections.Map(t.columns, func(col TableDataTypeColumn) string {
		return fmt.Sprintf("%s %s", col.name, col.dataType.Canonical())
	}), ", ")
	return fmt.Sprintf("%s(%s)", TableLegacyDataType, columns)
}

func (t *TableDataType) ToSqlWithoutUnknowns() string {
	// TODO [this PR]: implement while adding unknown logic to tables
	return t.ToSql()
}

func (t *TableDataType) Columns() []TableDataTypeColumn {
	return t.columns
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
	columns, err := collections.MapErr(strings.Split(onlyArgs, ","), func(arg string) (TableDataTypeColumn, error) {
		argParts := strings.SplitN(strings.TrimSpace(arg), " ", 2)
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

func areTableDataTypesDefinitelyDifferent(a, b *TableDataType) bool {
	// TODO [this PR]: implement!
	return false
}
