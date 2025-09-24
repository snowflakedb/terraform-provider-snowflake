package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *SemanticViewModel) WithTables(tables []sdk.LogicalTable) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(tables))
	for i, v := range tables {
		m := map[string]tfconfig.Variable{
			"table_name":    tfconfig.StringVariable(v.TableName.Name()),
			"database_name": tfconfig.StringVariable(v.TableName.DatabaseName()),
			"schema_name":   tfconfig.StringVariable(v.TableName.SchemaName()),
		}
		if v.Comment != nil {
			m["comment"] = tfconfig.StringVariable(*v.Comment)
		}
		logicalTableAlias := v.GetLogicalTableAlias()
		if logicalTableAlias != nil {
			m["logical_table_alias"] = tfconfig.StringVariable(logicalTableAlias.LogicalTableAlias)
		}
		primaryKeys := v.GetPrimaryKeys()
		if primaryKeys != nil {
			keys := make([]tfconfig.Variable, len(primaryKeys.PrimaryKey))
			for _, key := range primaryKeys.PrimaryKey {
				keys = append(keys, tfconfig.StringVariable(key.Name))
			}
			m["primary_keys"] = tfconfig.ListVariable(keys...)
		}
		uniqueKeys := v.GetUniqueKeys()
		if uniqueKeys != nil {
			keys := make([]tfconfig.Variable, len(uniqueKeys))
			for _, key := range uniqueKeys {
				uniKeys := make([]tfconfig.Variable, len(key.Unique))
				for _, uniKey := range key.Unique {
					uniKeys = append(uniKeys, tfconfig.StringVariable(uniKey.Name))
				}
				keys = append(keys, tfconfig.ListVariable(uniKeys...))
			}
			m["unique_keys"] = tfconfig.ListVariable(keys...)
		}
		synonyms := v.GetSynonyms()
		if synonyms != nil {
			syns := make([]tfconfig.Variable, len(synonyms.WithSynonyms))
			for _, synonym := range synonyms.WithSynonyms {
				syns = append(syns, tfconfig.StringVariable(synonym.Synonym))
			}
			m["synonyms"] = tfconfig.ListVariable(syns...)
		}
		maps[i] = tfconfig.MapVariable(m)
	}
	s.Tables = tfconfig.ListVariable(maps...)
	return s
}

func LogicalTableWithProps(
	alias string,
	tableName sdk.SchemaObjectIdentifier,
	primaryKeys []sdk.SemanticViewColumn,
	uniqueKeys [][]sdk.SemanticViewColumn,
	synonyms []sdk.Synonym,
	comment string,
) *sdk.LogicalTable {
	table := &sdk.LogicalTable{
		TableName: tableName,
		Comment:   &comment,
	}
	if alias != "" {
		table.SetLogicalTableAlias(alias)
	}
	if primaryKeys != nil {
		table.SetPrimaryKeys(primaryKeys)
	}
	if uniqueKeys != nil {
		table.SetUniqueKeys(uniqueKeys)
	}
	if synonyms != nil {
		table.SetSynonyms(synonyms)
	}
	return table
}
