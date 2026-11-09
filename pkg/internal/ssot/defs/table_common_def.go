package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// TableClusteringAction is the ALTER TABLE clustering clause: exactly one of CLUSTER BY, RECLUSTER, a recluster state
// change, and DROP CLUSTERING KEY. The block is a plain object clause. The SDK definition builds this query struct, and
// the handwritten ALTER field is a keyword with no SQL prefix, the same shape as AlterKeywordBlock. The rest of ALTER
// stays handwritten.
var TableClusteringAction = ssot.Attribute{
	Name:         "clustering_action",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{
			Name:      "cluster_by",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: ssot.StringListKeyword(ssot.Parentheses(), ssot.Sql("CLUSTER BY")),
		},
		{
			Name:      "recluster_action",
			Type:      ssot.AttrObject,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.Sql("RECLUSTER"), ssot.FieldName("Recluster")),
			Children: []ssot.Attribute{
				{Name: "max_size", Type: ssot.AttrInt, SqlOnly: true, CreateSql: ssot.IntParameter(ssot.Sql("MAX_SIZE"))},
				{Name: "condition", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PlainParameter(ssot.Sql("WHERE"), ssot.NoEquals(), ssot.FieldName("Condition"))},
			},
		},
		{
			Name:      "recluster_change_state",
			Type:      ssot.AttrObject,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.FieldName("ChangeReclusterState")),
			Children: []ssot.Attribute{
				{Name: "state", Type: ssot.AttrEnum, Enum: ReclusterState, SqlOnly: true, CreateSql: ssot.EnumKeyword()},
				{Name: "recluster", SqlOnly: true, CreateSql: ssot.StaticSQL("recluster", "RECLUSTER")},
			},
		},
		{
			Name:      "drop_clustering_key",
			Type:      ssot.AttrBool,
			SqlOnly:   true,
			CreateSql: ssot.BoolKeyword(ssot.Sql("DROP CLUSTERING KEY")),
		},
	},
}

var tableSearchMethodArgs = ssot.Attribute{
	Name:      "search_method_args",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.Parentheses(), ssot.Required(), ssot.Commas(), ssot.FieldName("Args")),
	Children: []ssot.Attribute{
		{
			Name:      "targets",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: ssot.StringListKeyword(ssot.FieldName("Targets")),
		},
		{
			Name:      "analyzer",
			Type:      ssot.AttrString,
			SqlOnly:   true,
			CreateSql: ssot.StringParameter(ssot.ArrowEquals()),
		},
	},
}

var tableSearchMethodWithTargetFields = []ssot.Attribute{
	{
		Name:      "method",
		Type:      ssot.AttrEnum,
		Enum:      TableSearchMethod,
		Required:  true,
		SqlOnly:   true,
		CreateSql: ssot.EnumKeyword(ssot.Required()),
	},
	tableSearchMethodArgs,
}

// TableSearchOptimizationAction is ADD or DROP SEARCH OPTIMIZATION. Iceberg tables reuse the same structs.
var TableSearchOptimizationAction = ssot.Attribute{
	Name:         "search_optimization_action",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{
			Name:      "add_search_optimization",
			Type:      ssot.AttrObject,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.FieldName("Add")),
			Children: []ssot.Attribute{
				{Name: "add_keyword", SqlOnly: true, CreateSql: ssot.StaticSQL("addSearchOptimization", "ADD SEARCH OPTIMIZATION")},
				{
					Name:      "search_method_with_target",
					Type:      ssot.AttrList,
					SqlOnly:   true,
					CreateSql: ssot.ValueKeyword(ssot.Sql("ON"), ssot.FieldName("On")),
					Elem: &ssot.Attribute{
						Type:      ssot.AttrObject,
						SqlOnly:   true,
						CreateSql: ssot.ObjectParameter(),
						Children:  tableSearchMethodWithTargetFields,
					},
				},
			},
		},
		{
			Name:      "drop_search_optimization",
			Type:      ssot.AttrObject,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.FieldName("Drop")),
			Children: []ssot.Attribute{
				{Name: "drop_keyword", SqlOnly: true, CreateSql: ssot.StaticSQL("dropSearchOptimization", "DROP SEARCH OPTIMIZATION")},
				{
					Name:      "drop_search_optimization_on",
					Type:      ssot.AttrList,
					SqlOnly:   true,
					CreateSql: ssot.ValueKeyword(ssot.Sql("ON"), ssot.FieldName("On")),
					Elem: &ssot.Attribute{
						Type:         ssot.AttrObject,
						SqlOnly:      true,
						ChildrenRule: ssot.ChildrenExactlyOne,
						CreateSql:    ssot.ObjectParameter(),
						Children: []ssot.Attribute{
							{
								Name:      "search_method_with_target",
								Type:      ssot.AttrObject,
								SqlOnly:   true,
								CreateSql: ssot.ObjectParameter(),
								Children:  tableSearchMethodWithTargetFields,
							},
							{
								Name:      "column_name",
								Type:      ssot.AttrString,
								SqlOnly:   true,
								CreateSql: ssot.ValueKeyword(),
							},
							{
								Name:      "expression_id",
								Type:      ssot.AttrString,
								SqlOnly:   true,
								CreateSql: ssot.ValueKeyword(),
							},
						},
					},
				},
			},
		},
	},
}

var (
	tableAggregationPolicyIdentifier = ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("AGGREGATION POLICY"), ssot.Required(), ssot.FieldName("AggregationPolicy"))
	tableAggregationPolicyEntityKey  = ssot.ListParameter("Column", ssot.Sql("ENTITY KEY"), ssot.NoEquals())
	tableJoinPolicyIdentifier        = ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("JOIN POLICY"), ssot.Required(), ssot.FieldName("JoinPolicy"))
)

// TableAggregationPolicy is CREATE ... AGGREGATION POLICY and ALTER ... SET/UNSET AGGREGATION POLICY.
// Iceberg CREATE reuses the CREATE struct; iceberg ALTER reuses the SET and UNSET projections.
var TableAggregationPolicy = ssot.Attribute{
	Name:      "aggregation_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	SetSql:    ssot.ObjectParameter(ssot.FieldName("SetAggregationPolicy")),
	UnsetSql:  ssot.ObjectParameter(ssot.FieldName("UnsetAggregationPolicy")),
	Children: []ssot.Attribute{
		{Name: "set", SqlOnly: true, SetSql: ssot.StaticSQL("set", "SET")},
		{
			Name:           "aggregation_policy",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			SqlOnly:        true,
			CreateSql:      tableAggregationPolicyIdentifier,
			SetSql:         tableAggregationPolicyIdentifier,
		},
		{
			Name:      "entity_key",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: tableAggregationPolicyEntityKey,
			SetSql:    tableAggregationPolicyEntityKey,
		},
		{
			Name:    "force",
			Type:    ssot.AttrBool,
			SqlOnly: true,
			SetSql:  ssot.BoolKeyword(ssot.Sql("FORCE")),
		},
		{Name: "unset_aggregation_policy", SqlOnly: true, UnsetSql: ssot.StaticSQL("unsetAggregationPolicy", "UNSET AGGREGATION POLICY")},
	},
}

// TableJoinPolicy is CREATE ... JOIN POLICY and ALTER ... SET/UNSET JOIN POLICY.
// ALLOWED JOIN KEYS is documented on CREATE TABLE and in an ALTER TABLE example, but not in the ALTER TABLE syntax.
// Iceberg ALTER reuses the SET and UNSET projections.
var TableJoinPolicy = ssot.Attribute{
	Name:      "join_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	SetSql:    ssot.ObjectParameter(ssot.FieldName("SetJoinPolicy")),
	UnsetSql:  ssot.ObjectParameter(ssot.FieldName("UnsetJoinPolicy")),
	Children: []ssot.Attribute{
		{Name: "set", SqlOnly: true, SetSql: ssot.StaticSQL("set", "SET")},
		{
			Name:           "join_policy",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			SqlOnly:        true,
			CreateSql:      tableJoinPolicyIdentifier,
			SetSql:         tableJoinPolicyIdentifier,
		},
		{
			Name:      "allowed_join_keys",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: ssot.ListParameter("Column", ssot.Sql("ALLOWED JOIN KEYS"), ssot.NoEquals()),
		},
		{
			Name:    "force",
			Type:    ssot.AttrBool,
			SqlOnly: true,
			SetSql:  ssot.BoolKeyword(ssot.Sql("FORCE")),
		},
		{Name: "unset_join_policy", SqlOnly: true, UnsetSql: ssot.StaticSQL("unsetJoinPolicy", "UNSET JOIN POLICY")},
	},
}
