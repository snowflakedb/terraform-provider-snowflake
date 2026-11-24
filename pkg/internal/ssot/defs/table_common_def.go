package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// TableClusteringAction is the ALTER TABLE clustering clause: exactly one of CLUSTER BY, RECLUSTER, a recluster state
// change, and DROP CLUSTERING KEY. The block is a plain object clause. The SDK definition builds this query struct, and
// the handwritten ALTER field is a keyword with no SQL prefix, the same shape as AlterKeywordBlock. The rest of ALTER
// stays handwritten.
var TableClusteringAction = ssot.Attribute{
	Name:        "clustering_action",
	Type:        ssot.AttrObject,
	SqlOnly:     true,
	Validations: []ssot.Validation{ssot.ExactlyOneOf()},
	CreateSql:   ssot.ObjectParameter(),
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
				{Name: "max_size", Type: ssot.AttrInt, SqlOnly: true, CreateSql: ssot.IntParameter()},
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
	Name:        "search_optimization_action",
	Type:        ssot.AttrObject,
	SqlOnly:     true,
	Validations: []ssot.Validation{ssot.ExactlyOneOf()},
	CreateSql:   ssot.ObjectParameter(),
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
						Type:        ssot.AttrObject,
						SqlOnly:     true,
						Validations: []ssot.Validation{ssot.ExactlyOneOf()},
						CreateSql:   ssot.ObjectParameter(),
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

// Shared CREATE and ALTER SQL for table policies. RAP ON is CREATE WITH and ALTER ADD; aggregation and join
// identifier/column lists are CREATE and ALTER SET.
var (
	tableRowAccessPolicyOnSql        = ssot.ListParameter("Column", ssot.Sql("ON"), ssot.NoEquals(), ssot.Required(), ssot.FieldName("On"))
	tableAggregationPolicyIdentifier = ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("AGGREGATION POLICY"), ssot.Required(), ssot.FieldName("AggregationPolicy"))
	tableAggregationPolicyEntityKey  = ssot.ListParameter("Column", ssot.Sql("ENTITY KEY"), ssot.NoEquals(), ssot.FieldName("EntityKey"))
	tableJoinPolicyIdentifier        = ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("JOIN POLICY"), ssot.Required(), ssot.FieldName("JoinPolicy"))
	tableJoinPolicyAllowedJoinKeys   = ssot.ListParameter("Column", ssot.Sql("ALLOWED JOIN KEYS"), ssot.NoEquals(), ssot.FieldName("AllowedJoinKeys"))
)

var tableRowAccessPolicyIdentifier = ssot.Attribute{
	Name:           "row_access_policy",
	Type:           ssot.AttrIdentifier,
	IdentifierType: "SchemaObjectIdentifier",
	Required:       true,
	SqlOnly:        true,
	CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("ROW ACCESS POLICY"), ssot.Required(), ssot.FieldName("RowAccessPolicy")),
}

// TableRowAccessPolicy is CREATE ... ROW ACCESS POLICY policy ON (columns). ALTER uses ADD/DROP below, not SET/UNSET.
var TableRowAccessPolicy = ssot.Attribute{
	Name:        "row_access_policy",
	Type:        ssot.AttrObject,
	SkipCreate:  true,
	SkipUpdate:  true,
	Description: "Specifies the row access policy to set on a table.",
	CreateSql:   ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{
			Name:            "policy_name",
			Type:            ssot.AttrIdentifier,
			IdentifierType:  "SchemaObjectIdentifier",
			Required:        true,
			RelatedResource: &ssot.Object{Name: "RowAccessPolicy"},
			Description:     "Row access policy name.",
			CreateSql:       ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("ROW ACCESS POLICY"), ssot.FieldName("Name")),
		},
		{
			Name:        "on",
			Type:        ssot.AttrSet,
			Required:    true,
			Elem:        &ssot.Attribute{Type: ssot.AttrString},
			Description: "Defines which columns are affected by the policy.",
			CreateSql:   tableRowAccessPolicyOnSql,
		},
	},
}

// TableAddRowAccessPolicy is ADD ROW ACCESS POLICY policy ON (columns). Iceberg tables reuse it.
var TableAddRowAccessPolicy = ssot.Attribute{
	Name:      "add_row_access_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{Name: "add_keyword", SqlOnly: true, CreateSql: ssot.StaticSQL("add", "ADD")},
		tableRowAccessPolicyIdentifier,
		{
			Name:      "on",
			Type:      ssot.AttrList,
			Required:  true,
			SqlOnly:   true,
			Elem:      &ssot.Attribute{Type: ssot.AttrString},
			CreateSql: tableRowAccessPolicyOnSql,
		},
	},
}

// TableDropRowAccessPolicy is DROP ROW ACCESS POLICY policy. Iceberg tables reuse it.
var TableDropRowAccessPolicy = ssot.Attribute{
	Name:      "drop_row_access_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{Name: "drop_keyword", SqlOnly: true, CreateSql: ssot.StaticSQL("drop", "DROP")},
		tableRowAccessPolicyIdentifier,
	},
}

// TableDropAndAddRowAccessPolicy replaces one row access policy with another in one statement. Iceberg tables reuse it.
// DROP ALL ROW ACCESS POLICIES stays a statement keyword on ALTER.
var TableDropAndAddRowAccessPolicy = ssot.Attribute{
	Name:      "drop_and_add_row_access_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{
			Name:      "drop_row_access_policy",
			Type:      ssot.AttrObject,
			Required:  true,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.Required(), ssot.FieldName("Drop")),
			Children:  TableDropRowAccessPolicy.Children,
		},
		{
			Name:      "add_row_access_policy",
			Type:      ssot.AttrObject,
			Required:  true,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.Required(), ssot.FieldName("Add")),
			Children:  TableAddRowAccessPolicy.Children,
		},
	},
}

// TableAggregationPolicy is CREATE ... AGGREGATION POLICY and ALTER ... SET/UNSET AGGREGATION POLICY.
// Iceberg CREATE reuses the CREATE struct; iceberg ALTER reuses the SET and UNSET projections.
var TableAggregationPolicy = ssot.Attribute{
	Name:        "aggregation_policy",
	Type:        ssot.AttrObject,
	SkipUpdate:  true,
	Description: "Specifies the aggregation policy to set on a table.",
	CreateSql:   ssot.ObjectParameter(),
	SetSql:      ssot.ObjectParameter(ssot.FieldName("SetAggregationPolicy")),
	UnsetSql:    ssot.ObjectParameter(ssot.FieldName("UnsetAggregationPolicy")),
	Children: []ssot.Attribute{
		{Name: "set", SqlOnly: true, SetSql: ssot.StaticSQL("set", "SET")},
		{
			Name:           "policy_name",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			Description:    "Aggregation policy name.",
			CreateSql:      tableAggregationPolicyIdentifier,
			SetSql:         tableAggregationPolicyIdentifier,
		},
		{
			Name:          "entity_key",
			Type:          ssot.AttrSet,
			Elem:          &ssot.Attribute{Type: ssot.AttrString},
			Description:   "Defines which columns uniquely identify an entity within the table.",
			CreateMapping: ssot.SetMappedValueCreate("toTableColumn"),
			CreateSql:     tableAggregationPolicyEntityKey,
			SetSql:        tableAggregationPolicyEntityKey,
		},
		{
			Name:    "force",
			Type:    ssot.AttrBool,
			SqlOnly: true,
			SetSql:  ssot.BoolKeyword(),
		},
		{Name: "unset_aggregation_policy", SqlOnly: true, UnsetSql: ssot.StaticSQL("unsetAggregationPolicy", "UNSET AGGREGATION POLICY")},
	},
}

// TableJoinPolicy is CREATE ... JOIN POLICY and ALTER ... SET/UNSET JOIN POLICY.
// Iceberg ALTER reuses the SET and UNSET projections.
var TableJoinPolicy = ssot.Attribute{
	Name:        "join_policy",
	Type:        ssot.AttrObject,
	SkipUpdate:  true,
	Description: "Specifies the join policy to set on a table.",
	CreateSql:   ssot.ObjectParameter(),
	SetSql:      ssot.ObjectParameter(ssot.FieldName("SetJoinPolicy")),
	UnsetSql:    ssot.ObjectParameter(ssot.FieldName("UnsetJoinPolicy")),
	Children: []ssot.Attribute{
		{Name: "set", SqlOnly: true, SetSql: ssot.StaticSQL("set", "SET")},
		{
			Name:           "policy_name",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			Description:    "Join policy name.",
			CreateSql:      tableJoinPolicyIdentifier,
			SetSql:         tableJoinPolicyIdentifier,
		},
		{
			Name:          "allowed_join_keys",
			Type:          ssot.AttrSet,
			Elem:          &ssot.Attribute{Type: ssot.AttrString},
			Description:   "Defines which columns can be used as join keys on the table.",
			CreateMapping: ssot.SetMappedValueCreate("toTableColumn"),
			CreateSql:     tableJoinPolicyAllowedJoinKeys,
			SetSql:        tableJoinPolicyAllowedJoinKeys,
		},
		{
			Name:    "force",
			Type:    ssot.AttrBool,
			SqlOnly: true,
			SetSql:  ssot.BoolKeyword(),
		},
		{Name: "unset_join_policy", SqlOnly: true, UnsetSql: ssot.StaticSQL("unsetJoinPolicy", "UNSET JOIN POLICY")},
	},
}
