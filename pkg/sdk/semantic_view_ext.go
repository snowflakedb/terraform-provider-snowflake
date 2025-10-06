package sdk

func (s *CreateSemanticViewRequest) GetName() SchemaObjectIdentifier {
	return s.name
}

func (l *LogicalTable) GetLogicalTableAlias() *LogicalTableAlias {
	return l.logicalTableAlias
}

func (l *LogicalTable) GetPrimaryKeys() *PrimaryKeys {
	return l.primaryKeys
}

func (l *LogicalTable) GetUniqueKeys() []UniqueKeys {
	return l.uniqueKeys
}

func (l *LogicalTable) GetSynonyms() *Synonyms {
	return l.synonyms
}

func (l *LogicalTable) SetLogicalTableAlias(alias string) {
	l.logicalTableAlias = &LogicalTableAlias{LogicalTableAlias: alias}
}

func (l *LogicalTable) SetPrimaryKeys(keys []SemanticViewColumn) {
	l.primaryKeys = &PrimaryKeys{
		PrimaryKey: keys,
	}
}

func (l *LogicalTable) SetUniqueKeys(keys [][]SemanticViewColumn) {
	for _, key := range keys {
		l.uniqueKeys = append(l.uniqueKeys, UniqueKeys{Unique: key})
	}
}

func (l *LogicalTable) SetSynonyms(synonyms []Synonym) {
	l.synonyms = &Synonyms{
		WithSynonyms: synonyms,
	}
}
