package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SemanticViewModel) WithTables(tables []sdk.LogicalTable) *SemanticViewModel {
	maps := make([]tfconfig.Variable, len(tables))
	for i, v := range tables {
		m := map[string]tfconfig.Variable{
			"table_name": tfconfig.StringVariable(v.TableName.FullyQualifiedName()),
		}
		if v.Comment != nil {
			m["comment"] = tfconfig.StringVariable(*v.Comment)
		}
		logicalTableAlias := v.GetLogicalTableAlias()
		if logicalTableAlias != nil {
			m["table_alias"] = tfconfig.StringVariable(logicalTableAlias.LogicalTableAlias)
		}
		primaryKeys := v.GetPrimaryKeys()
		if primaryKeys != nil {
			keys := make([]tfconfig.Variable, len(primaryKeys.PrimaryKey))
			for j, key := range primaryKeys.PrimaryKey {
				keys[j] = tfconfig.StringVariable(key.Name)
			}
			m["primary_key"] = tfconfig.ListVariable(keys...)
		}
		uniqueKeys := v.GetUniqueKeys()
		if uniqueKeys != nil {
			keys := make([]tfconfig.Variable, len(uniqueKeys))
			for j, key := range uniqueKeys {
				uniKeys := make([]tfconfig.Variable, len(key.Unique))
				for k, uniKey := range key.Unique {
					uniKeys[k] = tfconfig.StringVariable(uniKey.Name)
				}
				keys[j] = tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"values": tfconfig.ListVariable(uniKeys...),
				})
			}
			m["unique"] = tfconfig.ListVariable(keys...)
		}
		synonyms := v.GetSynonyms()
		if synonyms != nil {
			syns := make([]tfconfig.Variable, len(synonyms.WithSynonyms))
			for j, synonym := range synonyms.WithSynonyms {
				syns[j] = tfconfig.StringVariable(synonym.Synonym)
			}
			m["synonym"] = tfconfig.SetVariable(syns...)
		}
		maps[i] = tfconfig.ObjectVariable(m)
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
