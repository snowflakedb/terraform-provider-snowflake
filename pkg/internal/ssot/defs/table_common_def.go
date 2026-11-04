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
