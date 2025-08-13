package sdk

import "testing"

func TestSemanticViews_Create(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid CreateSemanticViewOptions
	defaultOpts := func() *CreateSemanticViewOptions {
		return &CreateSemanticViewOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSemanticViewOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		logicalTableId1 := randomSchemaObjectIdentifier()
		logicalTableId2 := randomSchemaObjectIdentifier()
		tableAlias1 := "table1"
		tableAlias2 := "table2"
		relationshipAlias1 := "rel1"
		logicalTableComment1 := String("logical table comment 1")
		logicalTableComment2 := String("logical table comment 2")
		tablesObj := []LogicalTable{
			{
				logicalTableAlias: &LogicalTableAlias{LogicalTableAlias: tableAlias1},
				TableName:         logicalTableId1,
				primaryKeys: &PrimaryKeys{PrimaryKey: []SemanticViewColumn{
					{
						Name: "pk1.1",
					},
					{
						Name: "pk1.2",
					},
				}},
				synonyms: &Synonyms{WithSynonyms: []string{"'test1'", "'test2'"}},
				Comment:  logicalTableComment1,
			},
			{
				logicalTableAlias: &LogicalTableAlias{LogicalTableAlias: tableAlias2},
				TableName:         logicalTableId2,
				primaryKeys: &PrimaryKeys{PrimaryKey: []SemanticViewColumn{
					{
						Name: "pk2.1",
					},
					{
						Name: "pk2.2",
					},
				}},
				synonyms: &Synonyms{WithSynonyms: []string{"'test3'", "'test4'"}},
				Comment:  logicalTableComment2,
			},
		}
		relationshipsObj := []SemanticViewRelationship{
			{
				relationshipAlias: &RelationshipAlias{RelationshipAlias: relationshipAlias1},
				tableName:         &RelationshipTableAlias{RelationshipTableAlias: tableAlias1},
				relationshipColumnNames: []SemanticViewColumn{
					{
						Name: "pk1.1",
					},
					{
						Name: "pk1.2",
					},
				},
				refTableName: &RelationshipTableAlias{RelationshipTableAlias: tableAlias2},
				relationshipRefColumnNames: []SemanticViewColumn{
					{
						Name: "pk2.1",
					},
					{
						Name: "pk2.2",
					},
				},
			},
		}
		opts := &CreateSemanticViewOptions{
			name:                      id,
			Comment:                   String("comment"),
			IfNotExists:               Bool(true),
			logicalTables:             tablesObj,
			semanticViewRelationships: relationshipsObj,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE SEMANTIC VIEW IF NOT EXISTS %s TABLES (%s AS %s PRIMARY KEY (pk1.1, pk1.2) WITH SYNONYMS ('test1', 'test2') COMMENT = '%s', %s AS %s PRIMARY KEY (pk2.1, pk2.2) WITH SYNONYMS ('test3', 'test4') COMMENT = '%s') RELATIONSHIPS %s AS %s (pk1.1, pk1.2) REFERENCES %s (pk2.1, pk2.2) COMMENT = '%s'`, id.FullyQualifiedName(), tableAlias1, logicalTableId1.FullyQualifiedName(), *logicalTableComment1, tableAlias2, logicalTableId2.FullyQualifiedName(), *logicalTableComment2, relationshipAlias1, tableAlias1, tableAlias2, "comment")
	})
}

func TestSemanticViews_Drop(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid DropSemanticViewOptions
	defaultOpts := func() *DropSemanticViewOptions {
		return &DropSemanticViewOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestSemanticViews_Describe(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid DescribeSemanticViewOptions
	defaultOpts := func() *DescribeSemanticViewOptions {
		return &DescribeSemanticViewOptions{

			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}

func TestSemanticViews_Show(t *testing.T) {
	// Minimal valid ShowSemanticViewOptions
	defaultOpts := func() *ShowSemanticViewOptions {
		return &ShowSemanticViewOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
