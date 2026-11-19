package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// TableAttrs holds reusable table fields.
//
// Omitted from Attributes and added in tableAttributesModificationExt:
//   - owner: computed. Readonly cannot sit on Attributes, and there is no computed schema shape.
//   - tag: the shared deprecated tagReferenceSchema.
//   - column.schema_evolution_record: a computed child. Nested attributes have no computed shape.
//
// The resource has no show_output and no describe_output. Read writes SHOW and DESCRIBE
// values onto the attribute fields and the computed fields added in the modification function.
//
// Chrome the generator does not emit is applied in tableAttributesModificationExt:
//   - column.type: IsDataTypeValid and DiffSuppressDataTypes.
//   - column.nullable: Default true. NativeBool keeps TypeBool and does not emit Default.
//   - column.comment, masking_policy, and collate: Default "".
//   - column.identity start_num and step_num: Default 1.
//   - primary_key: Deprecated.
//   - data_retention_time_in_days: Min and Max emit ValidateDiagFunc. The resource uses ValidateFunc.
//   - change_tracking: Default false.
//
// Identity and default on a column are mutually exclusive. Snowflake rejects that combination,
// so the schema does not enforce it. constant, expression, and sequence cannot use ConflictsWith
// because column is a list.
//
// TODO [SNOW-1348114]: there is no such separation of column default kinds on the SDK. Should we keep it in V1?
// TODO [SNOW-1348114]: consider adding fully_qualified_name for columns.
var TableAttrs = struct {
	ClusterBy               ssot.Attribute
	Column                  ssot.Attribute
	Comment                 ssot.Attribute
	PrimaryKey              ssot.Attribute
	DataRetentionTimeInDays ssot.Attribute
	ChangeTracking          ssot.Attribute
}{
	ClusterBy: ssot.Attribute{
		Name:          "cluster_by",
		Type:          ssot.AttrList,
		Description:   "A list of one or more table columns/expressions to be used as clustering key(s) for the table",
		Elem:          &ssot.Attribute{Type: ssot.AttrString},
		CreateMapping: ssot.AttributeMappedValueCreate("toTableClusterBy"),
		CreateSql:     ssot.StringListKeyword(ssot.Parentheses(), ssot.Sql("CLUSTER BY")),
	},
	Column: ssot.Attribute{
		Name:        "column",
		Type:        ssot.AttrList,
		Required:    true,
		Description: "Definitions of a column to create in the table. Minimum one required.",
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{
				Name:        "name",
				Type:        ssot.AttrString,
				Required:    true,
				Description: "Column name",
			},
			{
				Name:        "type",
				Type:        ssot.AttrString,
				Required:    true,
				Description: "Column type, e.g. VARIANT. For a full list of column types, see [Summary of Data Types](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).",
			},
			{
				Name:        "nullable",
				Type:        ssot.AttrBool,
				NativeBool:  true,
				Description: "Whether this column can contain null values. **Note**: Depending on your Snowflake version, the default value will not suffice if this column is used in a primary key constraint.",
			},
			{
				Name:        "default",
				Type:        ssot.AttrObject,
				Description: "Defines the column default value; note due to limitations of Snowflake's ALTER TABLE ADD/MODIFY COLUMN updates to default will not be applied",
				Children: []ssot.Attribute{
					{
						Name:        "constant",
						Type:        ssot.AttrString,
						Description: "The default constant value for the column",
					},
					{
						Name:        "expression",
						Type:        ssot.AttrString,
						Description: "The default expression value for the column",
					},
					{
						Name:           "sequence",
						Type:           ssot.AttrIdentifier,
						IdentifierType: "SchemaObjectIdentifier",
						Description:    "The default sequence to use for the column",
					},
				},
			},
			{
				Name:        "identity",
				Type:        ssot.AttrObject,
				Description: "Defines the identity start/step values for a column. **Note** Identity/default are mutually exclusive.",
				Children: []ssot.Attribute{
					{
						Name:        "start_num",
						Type:        ssot.AttrInt,
						Description: "The number to start incrementing at.",
					},
					{
						Name:        "step_num",
						Type:        ssot.AttrInt,
						Description: "Step size to increment by.",
					},
				},
			},
			{
				Name:        "comment",
				Type:        ssot.AttrString,
				Description: "Column comment",
			},
			{
				Name:        "masking_policy",
				Type:        ssot.AttrString,
				Description: "Masking policy to apply on column. It has to be a fully qualified name.",
			},
			{
				Name:        "collate",
				Type:        ssot.AttrString,
				Description: "Column collation, e.g. utf8",
			},
		}),
	},
	Comment: ssot.Attribute{
		Name:         "comment",
		Type:         ssot.AttrString,
		OptionalInDb: true,
		Description:  "Specifies a comment for the table.",
	},
	PrimaryKey: ssot.Attribute{
		Name:        "primary_key",
		Type:        ssot.AttrObject,
		Description: "Definitions of primary key constraint to create on table",
		Children: []ssot.Attribute{
			{
				Name:        "name",
				Type:        ssot.AttrString,
				Description: "Name of constraint",
			},
			{
				Name:        "keys",
				Type:        ssot.AttrList,
				Required:    true,
				Description: "Columns to use in primary key",
				Elem:        &ssot.Attribute{Type: ssot.AttrString},
			},
		},
	},
	DataRetentionTimeInDays: ssot.Attribute{
		Name:        "data_retention_time_in_days",
		Type:        ssot.AttrInt,
		Min:         new(-1),
		Max:         new(90),
		Description: "Specifies the retention period for the table so that Time Travel actions (SELECT, CLONE, UNDROP) can be performed on historical data in the table. If you wish to inherit the parent schema setting then pass in the schema attribute to this argument or do not fill this parameter at all; the default value for this field is -1, which is a fallback to use Snowflake default - in this case the schema value",
	},
	ChangeTracking: ssot.Attribute{
		Name:          "change_tracking",
		Type:          ssot.AttrBool,
		NativeBool:    true,
		StoredAsText:  true,
		OptionalInDb:  true,
		BoolTrueValue: "ON",
		Description:   "Specifies whether to enable change tracking on the table. Default false.",
	},
}

// TableScope and TableKind are the CREATE TABLE keywords before TABLE. The SDK Create request holds them,
// and the resource does not expose them.
var (
	TableScope          = ssot.NewEnum("TableScope", "TableScopes", "GLOBAL", "LOCAL")
	TableKind           = ssot.NewEnum("TableKind", "TableKinds", "TEMPORARY", "VOLATILE", "TRANSIENT")
	TableConstraintType = ssot.NewEnum("TableConstraintType", "TableConstraintTypes", "PRIMARY KEY", "UNIQUE", "FOREIGN KEY")
	// ReclusterState is RESUME or SUSPEND before the static RECLUSTER keyword.
	ReclusterState    = ssot.NewEnum("ReclusterState", "ReclusterStates", "RESUME", "SUSPEND")
	TableSearchMethod = ssot.NewEnum("TableSearchMethod", "TableSearchMethods", "SUBSTRING", "EQUALITY", "FULL_TEXT")
	// TableCopyOption is COPY = ( COMMENTS | TAGS | GRANTS ).
	TableCopyOption = ssot.NewEnum("TableCopyOption", "TableCopyOptions", "COMMENTS", "TAGS", "GRANTS")
)

// TableStorageLifecyclePolicy is ADD STORAGE LIFECYCLE POLICY policy ON (columns).
// DROP STORAGE LIFECYCLE POLICY stays a statement keyword.
var TableStorageLifecyclePolicy = ssot.Attribute{
	Name:      "storage_lifecycle_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{Name: "storage_lifecycle_policy", Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", Required: true, SqlOnly: true, CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("ADD STORAGE LIFECYCLE POLICY"), ssot.Required(), ssot.FieldName("StorageLifecyclePolicy"))},
		{Name: "on", Type: ssot.AttrList, Required: true, SqlOnly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}, CreateSql: ssot.ListParameter("Column", ssot.Sql("ON"), ssot.NoEquals(), ssot.Required(), ssot.FieldName("On"))},
	},
}

var tableColumnIdentityOrder = ssot.Attribute{Name: "order", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()}

// TableColumnVirtual is [ GENERATED ALWAYS ] AS ( expr ) [ VIRTUAL ] on CREATE TABLE and ALTER TABLE ADD COLUMN.
var TableColumnVirtual = ssot.Attribute{
	Name:      "virtual",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableColumnVirtual")),
	Children: []ssot.Attribute{
		{Name: "generated_always", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("GENERATED ALWAYS"))},
		// Callers include the parentheses in the expression. The ddl renderer does not wrap an unquoted parameter.
		{Name: "expression", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PlainParameter(ssot.Sql("AS"), ssot.NoEquals(), ssot.MustParentheses(), ssot.NoQuotes(), ssot.Required())},
		{Name: "virtual", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
	},
}

// TableColumnDefaultValue is DEFAULT or IDENTITY on a column. CREATE and ALTER ADD COLUMN share this struct, and its
// name stays ColumnDefaultValue because hybrid and iceberg tables refer to that type.
var TableColumnDefaultValue = ssot.Attribute{
	Name:         "default_value",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(ssot.StructName("ColumnDefaultValue")),
	Children: []ssot.Attribute{
		{
			Name:      "expression",
			Type:      ssot.AttrString,
			SqlOnly:   true,
			CreateSql: ssot.PlainParameter(ssot.Sql("DEFAULT"), ssot.NoEquals(), ssot.FieldName("Expression")),
		},
		{
			Name:      "identity",
			Type:      ssot.AttrObject,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.Sql("IDENTITY"), ssot.StructName("ColumnIdentity")),
			Children: []ssot.Attribute{
				{Name: "start", Type: ssot.AttrInt, Required: true, SqlOnly: true, CreateSql: ssot.IntParameter(ssot.NoQuotes(), ssot.NoEquals(), ssot.Required())},
				{Name: "increment", Type: ssot.AttrInt, Required: true, SqlOnly: true, CreateSql: ssot.IntParameter(ssot.NoQuotes(), ssot.NoEquals(), ssot.Required())},
				tableColumnIdentityOrder,
				// ORDER and NOORDER cannot both be set. Either can be absent.
				{Name: "noorder", Type: ssot.AttrBool, SqlOnly: true, ConflictsWith: []ssot.Attribute{tableColumnIdentityOrder}, CreateSql: ssot.BoolKeyword()},
			},
		},
	},
}

// TableRenameColumnAction is ALTER TABLE RENAME COLUMN. The old and new names are double-quoted.
// Table emits this struct with SharedToOpts, and iceberg tables reference it.
var TableRenameColumnAction = ssot.Attribute{
	Name:      "rename_column_action",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableRenameColumnAction"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "rename_column", SqlOnly: true, CreateSql: ssot.StaticSQL("renameColumn", "RENAME COLUMN")},
		{Name: "old_name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required(), ssot.FieldName("OldName"))},
		{Name: "new_name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.Sql("TO"), ssot.NoEquals(), ssot.DoubleQuotes(), ssot.Required(), ssot.FieldName("NewName"))},
	},
}

var (
	tableAlterColumn     = ssot.StaticSQL("alterColumn", "ALTER COLUMN")
	tableAlterColumnName = ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required(), ssot.FieldName("Name"))
)

// TableAlterColumnMaskingPolicy is ALTER COLUMN name SET MASKING POLICY and ALTER COLUMN name UNSET MASKING POLICY.
// SET is TableSetColumnMaskingPolicy and UNSET is TableUnsetColumnMaskingPolicy. Table emits both, and iceberg tables reference them.
var TableAlterColumnMaskingPolicy = ssot.Attribute{
	Name:     "alter_column_masking_policy",
	Type:     ssot.AttrObject,
	SqlOnly:  true,
	SetSql:   ssot.ObjectParameter(ssot.FieldName("SetColumnMaskingPolicy"), ssot.SharedToOpts()),
	UnsetSql: ssot.ObjectParameter(ssot.FieldName("UnsetColumnMaskingPolicy"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "alter_column", SqlOnly: true, SetSql: tableAlterColumn, UnsetSql: tableAlterColumn},
		{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, SetSql: tableAlterColumnName, UnsetSql: tableAlterColumnName},
		{Name: "set", SqlOnly: true, SetSql: ssot.StaticSQL("set", "SET")},
		{Name: "unset", SqlOnly: true, UnsetSql: ssot.StaticSQL("unset", "UNSET")},
		{
			Name:           "masking_policy",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			SqlOnly:        true,
			SetSql:         ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("MASKING POLICY"), ssot.Required(), ssot.FieldName("MaskingPolicy")),
			UnsetSql:       ssot.StaticSQL("maskingPolicy", "MASKING POLICY"),
		},
		{Name: "using", Type: ssot.AttrList, SqlOnly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}, SetSql: ssot.ListParameter("Column", ssot.Sql("USING"), ssot.NoEquals(), ssot.FieldName("Using"))},
		{Name: "force", Type: ssot.AttrBool, SqlOnly: true, SetSql: ssot.BoolKeyword(ssot.Sql("FORCE"))},
	},
}

// TableColumnMaskingPolicy is MASKING POLICY policy [USING (columns)] on a column. The CREATE column and ALTER ADD COLUMN
// embed it, and iceberg table columns reference the struct Table emits.
var TableColumnMaskingPolicy = ssot.Attribute{
	Name:      "column_masking_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableColumnMaskingPolicy"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "masking_policy", Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", Required: true, SqlOnly: true, CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("MASKING POLICY"), ssot.Required(), ssot.FieldName("MaskingPolicy"))},
		{Name: "using", Type: ssot.AttrList, SqlOnly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}, CreateSql: ssot.ListParameter("Column", ssot.Sql("USING"), ssot.NoEquals(), ssot.FieldName("Using"))},
	},
}

// TableColumnProjectionPolicy is PROJECTION POLICY policy on a column. The CREATE column and ADD COLUMN embed it, and
// iceberg table columns reference the struct Table emits.
var TableColumnProjectionPolicy = ssot.Attribute{
	Name:      "column_projection_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableColumnProjectionPolicy"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "projection_policy", Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", Required: true, SqlOnly: true, CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("PROJECTION POLICY"), ssot.Required(), ssot.FieldName("ProjectionPolicy"))},
	},
}

// TableAlterColumnProjectionPolicy is ALTER COLUMN name SET PROJECTION POLICY and ALTER COLUMN name UNSET PROJECTION POLICY.
// SET is TableSetColumnProjectionPolicy and UNSET is TableUnsetColumnProjectionPolicy. Table emits both, and iceberg tables reference them.
var TableAlterColumnProjectionPolicy = ssot.Attribute{
	Name:     "alter_column_projection_policy",
	Type:     ssot.AttrObject,
	SqlOnly:  true,
	SetSql:   ssot.ObjectParameter(ssot.FieldName("SetColumnProjectionPolicy"), ssot.SharedToOpts()),
	UnsetSql: ssot.ObjectParameter(ssot.FieldName("UnsetColumnProjectionPolicy"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "alter_column", SqlOnly: true, SetSql: tableAlterColumn, UnsetSql: tableAlterColumn},
		{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, SetSql: tableAlterColumnName, UnsetSql: tableAlterColumnName},
		{Name: "set", SqlOnly: true, SetSql: ssot.StaticSQL("set", "SET")},
		{Name: "unset", SqlOnly: true, UnsetSql: ssot.StaticSQL("unset", "UNSET")},
		{
			Name:           "projection_policy",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			SqlOnly:        true,
			SetSql:         ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("PROJECTION POLICY"), ssot.Required(), ssot.FieldName("ProjectionPolicy")),
			UnsetSql:       ssot.StaticSQL("projectionPolicy", "PROJECTION POLICY"),
		},
		{Name: "force", Type: ssot.AttrBool, SqlOnly: true, SetSql: ssot.BoolKeyword(ssot.Sql("FORCE"))},
	},
}

// TableAlterColumnTags is ALTER COLUMN name SET TAG and ALTER COLUMN name UNSET TAG.
// SET is TableSetColumnTags and UNSET is TableUnsetColumnTags. Table emits both, and iceberg tables reference them.
var TableAlterColumnTags = ssot.Attribute{
	Name:     "alter_column_tags",
	Type:     ssot.AttrObject,
	SqlOnly:  true,
	SetSql:   ssot.ObjectParameter(ssot.FieldName("SetColumnTags"), ssot.SharedToOpts()),
	UnsetSql: ssot.ObjectParameter(ssot.FieldName("UnsetColumnTags"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "alter_column", SqlOnly: true, SetSql: tableAlterColumn, UnsetSql: tableAlterColumn},
		{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, SetSql: tableAlterColumnName, UnsetSql: tableAlterColumnName},
		{Name: "set_tags", Tags: true, TagForm: ssot.TagFormSet},
		{Name: "unset_tags", Tags: true, TagForm: ssot.TagFormUnset},
	},
}

// TableDropColumnAction is ALTER TABLE DROP COLUMN. Table emits this struct, and iceberg tables reference it.
var TableDropColumnAction = ssot.Attribute{
	Name:      "drop_column_action",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableDropColumnAction"), ssot.SharedToOpts()),
	Children: []ssot.Attribute{
		{Name: "drop_column", SqlOnly: true, CreateSql: ssot.StaticSQL("dropColumn", "DROP COLUMN")},
		{Name: "if_exists", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("IF EXISTS"))},
		{Name: "columns", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PredefinedType("[]Column", ssot.NoEquals(), ssot.Required(), ssot.FieldName("Columns"))},
	},
}

// TableColumn is ALTER TABLE ... column actions: exactly one of add, rename, alter, set or unset masking policy, set or
// unset tags, and drop columns. It generates the TableColumnAction struct. The StructName options keep the public struct
// names, because the nested struct name TableColumn would clash with the CREATE column.
//
// Rename, set or unset masking policy, set or unset projection policy, set or unset tags, and drop columns are structs
// that iceberg tables reference. TableRenameColumnAction, TableSetColumnMaskingPolicy, TableUnsetColumnMaskingPolicy,
// TableSetColumnProjectionPolicy, TableUnsetColumnProjectionPolicy, TableSetColumnTags, TableUnsetColumnTags, and
// TableDropColumnAction are defined above with SharedToOpts. TableColumnInlineConstraint and TableColumnMaskingPolicy
// are defined with SharedToOpts too, and TableColumnProjectionPolicy is embedded as well.
var TableColumn = ssot.Attribute{
	Name:         "column",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(ssot.StructName("TableColumnAction")),
	Children: []ssot.Attribute{
		{
			// ADD COLUMN. NOT NULL is its own flag, which is what the table resource sets.
			Name:      "add",
			Type:      ssot.AttrObject,
			SqlOnly:   true,
			CreateSql: ssot.ObjectParameter(ssot.Sql("ADD"), ssot.NoComma(), ssot.StructName("TableColumnAddAction")),
			Children: []ssot.Attribute{
				{Name: "column", SqlOnly: true, CreateSql: ssot.StaticSQL("column", "COLUMN")},
				{Name: "if_not_exists", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("IF NOT EXISTS"))},
				{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required())},
				{Name: "column_type", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PredefinedType("datatypes.DataType", ssot.NoQuotes(), ssot.NoEquals(), ssot.Required())},
				ssot.WithCreateSql(&TableColumnVirtual, ssot.ObjectParameter(ssot.StructName("TableColumnVirtual"), ssot.NoComma())),
				{Name: "collate", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.NoEquals())},
				TableColumnDefaultValue,
				{Name: "not_null", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("NOT NULL"))},
				{Name: "inline_constraint", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableColumnInlineConstraint, ssot.NoComma())},
				{Name: "masking_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableColumnMaskingPolicy, ssot.NoComma())},
				{Name: "projection_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableColumnProjectionPolicy, ssot.NoComma())},
				{Name: "tags", Tags: true},
				{Name: "comment", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.NoEquals())},
			},
		},
		{Name: "rename", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableRenameColumnAction, ssot.NoComma())},
		{
			// ALTER COLUMN. The name is always present, so it is Required and not one of the exactly-one actions.
			Name:      "alter",
			Type:      ssot.AttrList,
			SqlOnly:   true,
			CreateSql: ssot.ValueKeyword(ssot.StructName("TableColumnAlterAction")),
			Elem: &ssot.Attribute{
				Type:         ssot.AttrObject,
				SqlOnly:      true,
				ChildrenRule: ssot.ChildrenExactlyOne,
				CreateSql:    ssot.ObjectParameter(),
				Children: []ssot.Attribute{
					{Name: "column", SqlOnly: true, CreateSql: ssot.StaticSQL("column", "COLUMN")},
					{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required())},
					{Name: "drop_default", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("DROP DEFAULT"))},
					{Name: "set_default", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PredefinedType("*SequenceName", ssot.FieldName("SetDefault"), ssot.Sql("SET DEFAULT"), ssot.NoEquals())},
					{
						Name:      "not_null_constraint",
						Type:      ssot.AttrObject,
						SqlOnly:   true,
						CreateSql: ssot.ObjectParameter(ssot.StructName("TableColumnNotNullConstraint")),
						Children: []ssot.Attribute{
							{Name: "set", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("SET NOT NULL"), ssot.FieldName("Set"))},
							{Name: "drop", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("DROP NOT NULL"), ssot.FieldName("Drop"))},
						},
					},
					{
						// SET DATA TYPE <type> [COLLATE '<spec>']. COLLATE only accompanies the type, so it is a child of it.
						Name:      "data_type",
						Type:      ssot.AttrObject,
						SqlOnly:   true,
						CreateSql: ssot.ObjectParameter(ssot.Sql("SET DATA TYPE"), ssot.StructName("TableColumnAlterDataType")),
						Children: []ssot.Attribute{
							{Name: "type", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PredefinedType("datatypes.DataType", ssot.FieldName("DataType"), ssot.NoQuotes(), ssot.NoEquals(), ssot.Required())},
							{Name: "collate", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.NoEquals())},
						},
					},
					{Name: "comment", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.NoEquals())},
					{Name: "unset_comment", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("UNSET COMMENT"))},
				},
			},
		},
		{Name: "set_masking_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOfUse(&TableAlterColumnMaskingPolicy, ssot.ClauseSet, ssot.NoComma())},
		{Name: "unset_masking_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOfUse(&TableAlterColumnMaskingPolicy, ssot.ClauseUnset, ssot.NoComma())},
		{Name: "set_projection_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOfUse(&TableAlterColumnProjectionPolicy, ssot.ClauseSet, ssot.NoComma())},
		{Name: "unset_projection_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOfUse(&TableAlterColumnProjectionPolicy, ssot.ClauseUnset, ssot.NoComma())},
		{Name: "set_tags", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOfUse(&TableAlterColumnTags, ssot.ClauseSet, ssot.NoComma())},
		{Name: "unset_tags", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOfUse(&TableAlterColumnTags, ssot.ClauseUnset, ssot.NoComma())},
		{Name: "drop_columns", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableDropColumnAction, ssot.NoComma())},
	},
}

// TableConstraintRenameAction is ALTER TABLE RENAME CONSTRAINT. The names are double-quoted, as ADD writes them.
// It is a child of TableConstraintAction.
var TableConstraintRenameAction = ssot.Attribute{
	Name:      "constraint_rename_action",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.Sql("RENAME CONSTRAINT"), ssot.NoComma(), ssot.FieldName("Rename"), ssot.StructName("TableConstraintRenameAction")),
	Children: []ssot.Attribute{
		{Name: "old_name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required(), ssot.FieldName("OldName"))},
		{Name: "new_name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PlainParameter(ssot.Sql("TO"), ssot.NoEquals(), ssot.DoubleQuotes(), ssot.Required(), ssot.FieldName("NewName"))},
	},
}

// TableConstraintAlterAction is ALTER TABLE ALTER CONSTRAINT. Exactly one of the name and the key kinds is set.
// The columns and the enforcement flags sit beside that choice, so ChildrenRuleFields names the fields the rule covers.
// It is a child of TableConstraintAction.
var TableConstraintAlterAction = ssot.Attribute{
	Name:               "constraint_alter_action",
	Type:               ssot.AttrObject,
	SqlOnly:            true,
	ChildrenRule:       ssot.ChildrenExactlyOne,
	ChildrenRuleFields: []string{"constraint_name", "primary_key", "unique", "foreign_key"},
	CreateSql:          ssot.ObjectParameter(ssot.Sql("ALTER"), ssot.NoComma(), ssot.FieldName("Alter"), ssot.StructName("TableConstraintAlterAction")),
	Children: []ssot.Attribute{
		{Name: "constraint_name", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PlainParameter(ssot.Sql("CONSTRAINT"), ssot.NoEquals(), ssot.DoubleQuotes(), ssot.FieldName("ConstraintName"))},
		{Name: "primary_key", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("PRIMARY KEY"))},
		{Name: "unique", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "foreign_key", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("FOREIGN KEY"))},
		{Name: "columns", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PredefinedType("[]string", ssot.NoEquals(), ssot.Parentheses(), ssot.FieldName("Columns"))},
		{Name: "enforced", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "not_enforced", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("NOT ENFORCED"))},
		{Name: "validate", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "novalidate", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "rely", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "norely", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
	},
}

// TableConstraintDropAction is ALTER TABLE DROP CONSTRAINT. Exactly one of the name and the key kinds is set.
// The columns and the CASCADE and RESTRICT flags sit beside that choice, so ChildrenRuleFields names the fields the rule covers.
// It is a child of TableConstraintAction.
var TableConstraintDropAction = ssot.Attribute{
	Name:               "constraint_drop_action",
	Type:               ssot.AttrObject,
	SqlOnly:            true,
	ChildrenRule:       ssot.ChildrenExactlyOne,
	ChildrenRuleFields: []string{"constraint_name", "primary_key", "unique", "foreign_key"},
	CreateSql:          ssot.ObjectParameter(ssot.Sql("DROP"), ssot.NoComma(), ssot.FieldName("Drop"), ssot.StructName("TableConstraintDropAction")),
	Children: []ssot.Attribute{
		{Name: "constraint_name", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PlainParameter(ssot.Sql("CONSTRAINT"), ssot.NoEquals(), ssot.DoubleQuotes(), ssot.FieldName("ConstraintName"))},
		{Name: "primary_key", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("PRIMARY KEY"))},
		{Name: "unique", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "foreign_key", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("FOREIGN KEY"))},
		{Name: "columns", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.PredefinedType("[]string", ssot.NoEquals(), ssot.Parentheses(), ssot.FieldName("Columns"))},
		{Name: "cascade", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
		{Name: "restrict", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword()},
	},
}

// TableConstraintAction is one ALTER TABLE constraint branch. It composes the rename, alter, and drop attributes above.
// ADD is the out-of-line constraint, so it is the one preexisting struct. Table emits it with SharedToOpts, and iceberg tables reference it.
var TableConstraintAction = ssot.Attribute{
	Name:         "constraint_action",
	Type:         ssot.AttrObject,
	SqlOnly:      true,
	ChildrenRule: ssot.ChildrenExactlyOne,
	CreateSql:    ssot.ObjectParameter(ssot.StructName("TableConstraintAction")),
	Children: []ssot.Attribute{
		{Name: "add", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableOutOfLineConstraint, ssot.NoComma(), ssot.Sql("ADD"))},
		TableConstraintRenameAction,
		TableConstraintAlterAction,
		TableConstraintDropAction,
	},
}

// TableExternalTableColumnAddAction is ALTER TABLE ADD COLUMN with an AS expression, as Snowflake documents it for external
// table columns. The name is not quoted. It is a child of TableExternalTableAction.
var TableExternalTableColumnAddAction = ssot.Attribute{
	Name:      "external_table_column_add_action",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.NoComma(), ssot.FieldName("Add"), ssot.StructName("TableExternalTableColumnAddAction")),
	Children: []ssot.Attribute{
		{Name: "add_column", SqlOnly: true, CreateSql: ssot.StaticSQL("addColumn", "ADD COLUMN")},
		{Name: "if_not_exists", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("IF NOT EXISTS"))},
		{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.Required(), ssot.FieldName("Name"))},
		{Name: "column_type", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PredefinedType("datatypes.DataType", ssot.NoQuotes(), ssot.NoEquals(), ssot.Required(), ssot.FieldName("ColumnType"))},
		{Name: "expression", Type: ssot.AttrList, SqlOnly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}, CreateSql: ssot.ListParameter("string", ssot.Sql("AS"), ssot.NoEquals(), ssot.FieldName("Expression"))},
		{Name: "comment", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.NoEquals())},
	},
}

// TableExternalTableAction is the external table column branch of ALTER TABLE. ADD is the only one: renaming and
// dropping a column use the regular column actions. ChildrenRule needs two alternatives, so the SDK extension
// (pkg/sdk/tables_ext.go) requires ADD.
var TableExternalTableAction = ssot.Attribute{
	Name:      "external_table_action",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableExternalTableAction")),
	Children:  []ssot.Attribute{TableExternalTableColumnAddAction},
}

// TableCreateColumn is one column of CREATE TABLE. It generates the TableColumn struct. The inline constraint, the masking
// policy, and the projection policy are preexisting structs, and iceberg tables reference them. The clauses keep the CREATE
// order, which is not the ADD COLUMN order.
var TableCreateColumn = ssot.Attribute{
	Name:      "create_column",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("TableColumn")),
	Children: []ssot.Attribute{
		{Name: "name", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.TextKeyword(ssot.DoubleQuotes(), ssot.Required())},
		{Name: "column_type", Type: ssot.AttrString, Required: true, SqlOnly: true, CreateSql: ssot.PredefinedType("datatypes.DataType", ssot.NoQuotes(), ssot.NoEquals(), ssot.Required())},
		ssot.WithCreateSql(&TableColumnVirtual, ssot.ObjectParameter(ssot.StructName("TableColumnVirtual"), ssot.KeywordShape())),
		{Name: "inline_constraint", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableColumnInlineConstraint, ssot.KeywordShape())},
		{Name: "not_null", Type: ssot.AttrBool, SqlOnly: true, CreateSql: ssot.BoolKeyword(ssot.Sql("NOT NULL"))},
		{Name: "collate", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.Sql("COLLATE"), ssot.NoEquals())},
		ssot.WithCreateSql(&TableColumnDefaultValue, ssot.ObjectParameter(ssot.StructName("ColumnDefaultValue"), ssot.KeywordShape())),
		{Name: "masking_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableColumnMaskingPolicy, ssot.KeywordShape())},
		{Name: "projection_policy", Type: ssot.AttrObject, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableColumnProjectionPolicy, ssot.KeywordShape())},
		{Name: "tags", Tags: true},
		{Name: "comment", Type: ssot.AttrString, SqlOnly: true, CreateSql: ssot.StringParameter(ssot.Sql("COMMENT"), ssot.NoEquals())},
	},
}

// TableCreateColumnsAndConstraints is the parenthesized column and out-of-line constraint list of CREATE TABLE. It generates
// the CreateTableColumnsAndConstraints struct. The columns are TableCreateColumn, and the out-of-line constraints are the
// struct Table emits, which iceberg tables reference.
var TableCreateColumnsAndConstraints = ssot.Attribute{
	Name:      "create_columns_and_constraints",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(ssot.StructName("CreateTableColumnsAndConstraints")),
	Children: []ssot.Attribute{
		{Name: "columns", Type: ssot.AttrList, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableCreateColumn, ssot.KeywordShape())},
		{Name: "out_of_line_constraint", Type: ssot.AttrList, SqlOnly: true, CreateSql: ssot.PreexistingStructOf(&TableOutOfLineConstraint, ssot.KeywordShape())},
	},
}

// TableCreateAttrs are the CREATE TABLE clauses that are not resource attributes. They are SqlOnly, so the schema,
// resource, model, and assertion generators skip them.
//
// The column list names TableCreateColumnsAndConstraints by PreexistingStructOf. The SDK definition builds that struct
// from this attribute, and the shared structs it names, through the attribute-backed clause. Search optimization,
// clustering, and storage lifecycle policy are SSOT attributes reused through QueryStructFromAttribute. CREATE row access policy is generated
// here as TableRowAccessPolicy. ALTER add, drop, and drop-and-add are TableAddRowAccessPolicy, TableDropRowAccessPolicy,
// and TableDropAndAddRowAccessPolicy in table_common_def.go; iceberg ALTER reuses those structs. Aggregation and join
// policies are TableAggregationPolicy and TableJoinPolicy in table_common_def.go. Iceberg CREATE reuses the CREATE
// structs; iceberg ALTER reuses the SET/UNSET projections.
var TableCreateAttrs = struct {
	Scope                      ssot.Attribute
	Kind                       ssot.Attribute
	ColumnsAndConstraints      ssot.Attribute
	EnableSchemaEvolution      ssot.Attribute
	MaxDataExtensionTimeInDays ssot.Attribute
	DefaultDdlCollation        ssot.Attribute
	ErrorLogging               ssot.Attribute
	Copy                       ssot.Attribute
	IcebergDefaultDdlCollation ssot.Attribute
	RowTimestamp               ssot.Attribute
	Contact                    ssot.Attribute
	StorageLifecyclePolicy     ssot.Attribute
	RowAccessPolicy            ssot.Attribute
	AggregationPolicy          ssot.Attribute
	JoinPolicy                 ssot.Attribute
}{
	Scope: ssot.Attribute{
		Name:      "scope",
		Type:      ssot.AttrEnum,
		Enum:      TableScope,
		SqlOnly:   true,
		CreateSql: ssot.EnumKeyword(),
	},
	Kind: ssot.Attribute{
		Name:      "kind",
		Type:      ssot.AttrEnum,
		Enum:      TableKind,
		SqlOnly:   true,
		CreateSql: ssot.EnumKeyword(),
	},
	ColumnsAndConstraints: ssot.Attribute{
		Name:      "columns_and_constraints",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.PreexistingStructOf(&TableCreateColumnsAndConstraints, ssot.Parentheses(), ssot.Required()),
	},
	EnableSchemaEvolution: ssot.Attribute{
		Name:    "enable_schema_evolution",
		Type:    ssot.AttrBool,
		SqlOnly: true,
	},
	MaxDataExtensionTimeInDays: ssot.Attribute{
		Name:    "max_data_extension_time_in_days",
		Type:    ssot.AttrInt,
		SqlOnly: true,
	},
	DefaultDdlCollation: ssot.Attribute{
		Name:      "default_ddl_collation",
		Type:      ssot.AttrString,
		SqlOnly:   true,
		CreateSql: ssot.StringParameter(ssot.AllowEmpty()),
		SetSql:    ssot.StringParameter(ssot.AllowEmpty()),
	},
	ErrorLogging: ssot.Attribute{
		Name:    "error_logging",
		Type:    ssot.AttrBool,
		SqlOnly: true,
	},
	Copy: ssot.Attribute{
		Name:      "copy",
		Type:      ssot.AttrList,
		SqlOnly:   true,
		Elem:      &ssot.Attribute{Type: ssot.AttrEnum, Enum: TableCopyOption},
		CreateSql: ssot.ListParameter(TableCopyOption.Name, ssot.Parentheses(), ssot.NoQuotes()),
	},
	IcebergDefaultDdlCollation: ssot.Attribute{
		Name:    "iceberg_default_ddl_collation",
		Type:    ssot.AttrString,
		SqlOnly: true,
		SetSql:  ssot.StringParameter(ssot.AllowEmpty()),
	},
	RowTimestamp: ssot.Attribute{
		Name:    "row_timestamp",
		Type:    ssot.AttrBool,
		SqlOnly: true,
	},
	// CREATE is WITH CONTACT ( purpose = contact ). SET is CONTACT purpose = contact and may mix
	// with other SET properties. UNSET CONTACT <purpose> is its own ALTER statement (Snowflake
	// rejects mixing CONTACT with other UNSET properties). The Go type for CREATE/SET is the
	// handwritten TableContact iceberg tables already use; UNSET is the purpose string.
	Contact: ssot.Attribute{
		Name:      "contact",
		Type:      ssot.AttrString,
		SqlOnly:   true,
		CreateSql: ssot.PredefinedType("[]TableContact", ssot.KeywordShape(), ssot.Parentheses(), ssot.Sql("WITH CONTACT"), ssot.FieldName("Contact")),
		SetSql:    ssot.PredefinedType("[]TableContact", ssot.KeywordShape(), ssot.Sql("CONTACT")),
		UnsetSql:  ssot.PlainParameter(ssot.NoEquals()),
	},
	// CREATE is STORAGE LIFECYCLE POLICY policy ON (columns). ALTER ADD stays TableStorageLifecyclePolicy.
	StorageLifecyclePolicy: ssot.Attribute{
		Name:      "storage_lifecycle_policy",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(ssot.StructName("CreateTableStorageLifecyclePolicy")),
		Children: []ssot.Attribute{
			{Name: "storage_lifecycle_policy", Type: ssot.AttrIdentifier, IdentifierType: "SchemaObjectIdentifier", Required: true, SqlOnly: true, CreateSql: ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("STORAGE LIFECYCLE POLICY"), ssot.Required(), ssot.FieldName("StorageLifecyclePolicy"))},
			{Name: "on", Type: ssot.AttrList, Required: true, SqlOnly: true, Elem: &ssot.Attribute{Type: ssot.AttrString}, CreateSql: ssot.ListParameter("Column", ssot.Sql("ON"), ssot.NoEquals(), ssot.Required(), ssot.FieldName("On"))},
		},
	},
	RowAccessPolicy: ssot.Attribute{
		Name:      "row_access_policy",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.ObjectParameter(),
		Children: []ssot.Attribute{
			{
				Name:           "name",
				Type:           ssot.AttrIdentifier,
				IdentifierType: "SchemaObjectIdentifier",
				Required:       true,
				SqlOnly:        true,
				CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("ROW ACCESS POLICY"), ssot.FieldName("Name")),
			},
			{
				Name:      "on",
				Type:      ssot.AttrList,
				Required:  true,
				SqlOnly:   true,
				Elem:      &ssot.Attribute{Type: ssot.AttrString},
				CreateSql: ssot.StringListKeyword(ssot.Parentheses()),
			},
		},
	},
	AggregationPolicy: TableAggregationPolicy,
	JoinPolicy:        TableJoinPolicy,
}

func tableShowString(name string, optionalInDb, optionalInShow bool) ssot.Attribute {
	return ssot.Attribute{Name: name, Type: ssot.AttrString, Readonly: true, OptionalInDb: optionalInDb, OptionalInShow: optionalInShow}
}

func tableShowInt(name string, optionalInShow bool) ssot.Attribute {
	return ssot.Attribute{Name: name, Type: ssot.AttrInt, Readonly: true, OptionalInDb: true, OptionalInShow: optionalInShow}
}

func tableShowBool(name, trueValue string) ssot.Attribute {
	return ssot.Attribute{
		Name: name, Type: ssot.AttrBool, Readonly: true, StoredAsText: true, OptionalInDb: true, BoolTrueValue: trueValue,
	}
}

var (
	tableShowCreatedOn                  = tableShowString("created_on", false, false)
	tableShowName                       = tableShowString("name", false, false)
	tableShowSchemaName                 = tableShowString("schema_name", false, false)
	tableShowDatabaseName               = tableShowString("database_name", false, false)
	tableShowKind                       = tableShowString("kind", false, false)
	tableShowClusterBy                  = ssot.Attribute{Name: "show_cluster_by", ShowName: "cluster_by", Type: ssot.AttrString, Readonly: true, OptionalInDb: true}
	tableShowOwner                      = tableShowString("owner", false, false)
	tableShowRows                       = tableShowInt("rows", false)
	tableShowBytes                      = tableShowInt("bytes", true)
	tableShowRetentionTime              = tableShowInt("retention_time", false)
	tableShowDroppedOn                  = tableShowString("dropped_on", true, true)
	tableShowAutomaticClustering        = tableShowBool("automatic_clustering", "ON")
	tableShowSearchOptimization         = tableShowBool("search_optimization", "ON")
	tableShowSearchOptimizationProgress = tableShowString("search_optimization_progress", true, false)
	tableShowSearchOptimizationBytes    = tableShowInt("search_optimization_bytes", true)
	tableShowIsExternal                 = tableShowBool("is_external", "Y")
	tableShowEnableSchemaEvolution      = tableShowBool("enable_schema_evolution", "Y")
	tableShowOwnerRoleType              = tableShowString("owner_role_type", true, false)
	tableShowIsEvent                    = tableShowBool("is_event", "Y")
	tableShowIsHybrid                   = tableShowBool("is_hybrid", "Y")
	tableShowIsIceberg                  = tableShowBool("is_iceberg", "Y")
	tableShowIsDynamic                  = tableShowBool("is_dynamic", "Y")
	tableShowIsImmutable                = tableShowBool("is_immutable", "Y")
	tableShowIsInteractive              = tableShowBool("is_interactive", "Y")
	tableShowRowTimestamp               = tableShowBool("row_timestamp", "ON")
	tableShowErrorLogging               = tableShowBool("error_logging", "ON")
	tableShowBudget                     = tableShowString("budget", true, true)
)

var Table = ssot.Object{
	Name:       "Table",
	NamePlural: "Tables",
	// These enums are in the nested ALTER blocks and constraints, which the enum walk does not reach.
	AdditionalEnums:        []*ssot.Enum{TableConstraintType, ReclusterState, MatchType, TableSearchMethod, TableCopyOption, TableDataMetricScheduleStatusOperation, TableDataMetricSensitivity},
	IdType:                 "SchemaObjectIdentifier",
	PreviewFeature:         true,
	DescribeFailRead:       true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithObjectAssert(ssot.ObjectAssertShow()),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(ssot.SdkCreateAdditionalValidations()),
			ssot.SdkDrop(
				ssot.SdkDropCascadeRestrict(),
				ssot.SdkDropDoc("https://docs.snowflake.com/en/sql-reference/sql/drop-table"),
			),
			ssot.SdkShow(
				ssot.SdkShowDoc("https://docs.snowflake.com/en/sql-reference/sql/show-tables"),
				ssot.SdkShowTerse(),
				ssot.SdkShowOptionalSql("HISTORY"),
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalExtendedIn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimitFrom(),
				ssot.SdkShowAdditionalValidations(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDExtendedIn, ssot.SdkShowByIDLike),
			),
			ssot.SdkCustomOperations(),
			ssot.SdkAlter(ssot.SdkAlterAdditionalValidations()),
			ssot.SdkShowParameters(),
		),
		ssot.WithDataSource(
			ssot.DataSourceWithTerse(),
			ssot.DataSourceWithHistory(),
			ssot.DataSourceDescribeMethod("DescribeColumns"),
			ssot.DataSourceAdditionalDescribeOutputDescription("The statement is DESCRIBE TABLE ... TYPE = COLUMNS. It is empty for a dropped table that `history` returns."),
			ssot.DataSourceAdditionalParametersOutputDescription("It is empty for a dropped table that `history` returns."),
		),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow()),
		ssot.WithResourceParametersAssert(),
	},
	Create: ssot.CreateSpec{
		KeywordsBefore: []*ssot.Attribute{
			&TableCreateAttrs.Scope,
			&TableCreateAttrs.Kind,
		},
		Clauses: []*ssot.Attribute{
			&TableCreateAttrs.ColumnsAndConstraints,
			&TableAttrs.ClusterBy,
			&TableCreateAttrs.EnableSchemaEvolution,
			&TableAttrs.DataRetentionTimeInDays,
			&TableCreateAttrs.MaxDataExtensionTimeInDays,
			&TableAttrs.ChangeTracking,
			&TableCreateAttrs.DefaultDdlCollation,
			&TableCreateAttrs.ErrorLogging,
			&TableCreateAttrs.Copy,
			&TableCreateAttrs.RowAccessPolicy,
			&TableCreateAttrs.AggregationPolicy,
			&TableCreateAttrs.JoinPolicy,
			&TableCreateAttrs.StorageLifecyclePolicy,
			ssot.Tags,
			&TableCreateAttrs.Contact,
			&TableAttrs.Comment,
			&TableCreateAttrs.RowTimestamp,
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&tableShowCreatedOn,
			&tableShowName,
			&tableShowSchemaName,
			&tableShowDatabaseName,
			&tableShowKind,
			&TableAttrs.Comment,
			&tableShowClusterBy,
			&tableShowRows,
			&tableShowBytes,
			&tableShowOwner,
			&tableShowRetentionTime,
			&tableShowDroppedOn,
			&tableShowAutomaticClustering,
			&TableAttrs.ChangeTracking,
			&tableShowSearchOptimization,
			&tableShowSearchOptimizationProgress,
			&tableShowSearchOptimizationBytes,
			&tableShowIsExternal,
			&tableShowEnableSchemaEvolution,
			&tableShowOwnerRoleType,
			&tableShowIsEvent,
			&tableShowIsHybrid,
			&tableShowIsIceberg,
			&tableShowIsDynamic,
			&tableShowIsImmutable,
			&tableShowIsInteractive,
			&tableShowRowTimestamp,
			&tableShowErrorLogging,
			&tableShowBudget,
		},
	},
	// The branch order is the handwritten ALTER TABLE field order. A keyword block with no SQL
	// leaves the words inside the nested struct. Search optimization, row access policy, and the aggregation and
	// join policy projections are shared with iceberg tables, so Table emits them with SharedToOpts. The aggregation and
	// join policy fields are named like the SetSql and UnsetSql FieldName of their attribute, so they embed that projection.
	// Drop and add row access policy is comma-separated. DROP ALL and DROP STORAGE
	// LIFECYCLE POLICY stay statement keywords.
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSwapWith(),
			ssot.AlterKeywordBlock("ClusteringAction", "", &TableClusteringAction),
			ssot.AlterKeywordBlock("ColumnAction", "", &TableColumn),
			ssot.AlterKeywordBlock("ConstraintAction", "", &TableConstraintAction),
			ssot.AlterKeywordBlock("ExternalTableAction", "", &TableExternalTableAction),
			ssot.AlterKeywordBlock("SearchOptimizationAction", "", &TableSearchOptimizationAction).WithSharedToOpts(),
			// SET and UNSET are the alter-table parameters. UNSET is comma-separated: the alter-table
			// grammar repeats the properties with [ , ... ].
			ssot.AlterSet(
				&TableCreateAttrs.EnableSchemaEvolution,
				&TableAttrs.DataRetentionTimeInDays,
				&TableCreateAttrs.MaxDataExtensionTimeInDays,
				&TableAttrs.ChangeTracking,
				&TableCreateAttrs.DefaultDdlCollation,
				&TableCreateAttrs.IcebergDefaultDdlCollation,
				&TableCreateAttrs.ErrorLogging,
				&TableCreateAttrs.Contact,
				&TableAttrs.Comment,
				&TableCreateAttrs.RowTimestamp,
			),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterUnset(
				ssot.AlterGroup("Contact", &TableCreateAttrs.Contact),
				ssot.AlterGroup(ssot.MainAlterGroup,
					&TableAttrs.DataRetentionTimeInDays,
					&TableCreateAttrs.MaxDataExtensionTimeInDays,
					&TableAttrs.ChangeTracking,
					&TableCreateAttrs.DefaultDdlCollation,
					&TableCreateAttrs.IcebergDefaultDdlCollation,
					&TableCreateAttrs.EnableSchemaEvolution,
					&TableCreateAttrs.ErrorLogging,
					&TableAttrs.Comment,
					&TableCreateAttrs.RowTimestamp,
				),
			),
			ssot.AlterKeywordBlock("AddRowAccessPolicy", "", &TableAddRowAccessPolicy).WithSharedToOpts(),
			ssot.AlterKeywordBlock("DropRowAccessPolicy", "", &TableDropRowAccessPolicy).WithSharedToOpts(),
			ssot.AlterCommaSeparatedBlock("DropAndAddRowAccessPolicy", &TableDropAndAddRowAccessPolicy).WithSharedToOpts(),
			ssot.AlterKeyword("DROP ALL ROW ACCESS POLICIES"),
			ssot.AlterKeywordBlock("AddStorageLifecyclePolicy", "", &TableStorageLifecyclePolicy),
			ssot.AlterKeyword("DROP STORAGE LIFECYCLE POLICY"),
			ssot.AlterKeywordBlock("SetAggregationPolicy", "", &TableAggregationPolicy).WithSharedToOpts(),
			ssot.AlterKeywordBlock("UnsetAggregationPolicy", "", &TableAggregationPolicy).WithSharedToOpts(),
			ssot.AlterKeywordBlock("SetJoinPolicy", "", &TableJoinPolicy).WithSharedToOpts(),
			ssot.AlterKeywordBlock("UnsetJoinPolicy", "", &TableJoinPolicy).WithSharedToOpts(),
			ssot.AlterKeywordBlock("SetDataMetricSchedule", "", &TableSetDataMetricSchedule),
			ssot.AlterKeyword("UNSET DATA_METRIC_SCHEDULE"),
			ssot.AlterKeywordBlock("AddDataMetricFunction", "", &TableAddDataMetricFunction),
			ssot.AlterKeywordBlock("DropDataMetricFunction", "", &TableDropDataMetricFunction),
			ssot.AlterKeywordBlock("ModifyDataMetricFunction", "", &TableModifyDataMetricFunction),
		},
	},
	// Column stays in NewCreateRequestExt: it is a required list block. Tag stays there too, because it is not an attribute.
	// ShowByIdSafely stays in the extension because Read compares schema retention before describing columns.
	// Describe uses DescribeColumns. SetConfigFields copies owner and cluster keys. Column diffs,
	// primary key, and tags stay in the extension. Generated AfterAlter and SetDescribeOutput are unimplemented.
	// The importer is a passthrough, so ShowByIdInSdk stays omitted.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithSetStateToValuesFromConfig(ssot.Omitted()),
		ssot.WithShowOutputSet(ssot.Omitted()),
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.Generated()),
		ssot.WithCustomizeDiff(ssot.Generated()),
		ssot.WithApplySetUnset(ssot.Generated()),
		ssot.WithAlterSetInSdk(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithAfterAlter(ssot.ExtCall()),
		ssot.WithShowByIdInSdk(ssot.Omitted()),
		ssot.WithSetFieldsNotSetByRead(ssot.Omitted()),
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithStateUpgraders(ssot.Generated()),
	},
	// The id changed from the pipe-separated format (db|schema|name) to the fully qualified name.
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName",
	},
	CustomOperations: append(tableCreateVariantOperations, tableDescribeOperations...),
	Attributes: []ssot.Attribute{
		TableAttrs.ClusterBy,
		TableAttrs.Column,
		TableAttrs.Comment,
		TableAttrs.PrimaryKey,
		TableAttrs.DataRetentionTimeInDays,
		TableAttrs.ChangeTracking,
	},
}
