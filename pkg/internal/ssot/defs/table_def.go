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
		MinItems:    new(1),
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
				Type:        ssot.AttrList,
				MinItems:    new(1),
				MaxItems:    new(1),
				Description: "Defines the column default value; note due to limitations of Snowflake's ALTER TABLE ADD/MODIFY COLUMN updates to default will not be applied",
				Elem: ssot.ObjectElem([]ssot.Attribute{
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
				}),
			},
			{
				Name:        "identity",
				Type:        ssot.AttrList,
				MinItems:    new(1),
				MaxItems:    new(1),
				Description: "Defines the identity start/step values for a column. **Note** Identity/default are mutually exclusive.",
				Elem: ssot.ObjectElem([]ssot.Attribute{
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
				}),
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
		Type:        ssot.AttrList,
		MaxItems:    new(1),
		Description: "Definitions of primary key constraint to create on table",
		Elem: ssot.ObjectElem([]ssot.Attribute{
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
		}),
	},
	DataRetentionTimeInDays: ssot.Attribute{
		Name:        "data_retention_time_in_days",
		Type:        ssot.AttrInt,
		Min:         new(-1),
		Max:         new(90),
		Description: "Specifies the retention period for the table so that Time Travel actions (SELECT, CLONE, UNDROP) can be performed on historical data in the table. If you wish to inherit the parent schema setting then pass in the schema attribute to this argument or do not fill this parameter at all; the default value for this field is -1, which is a fallback to use Snowflake default - in this case the schema value",
	},
	ChangeTracking: ssot.Attribute{
		Name:        "change_tracking",
		Type:        ssot.AttrBool,
		NativeBool:  true,
		Description: "Specifies whether to enable change tracking on the table. Default false.",
	},
}

// TableScope and TableKind are the CREATE TABLE keywords before TABLE. The SDK Create request holds them,
// and the resource does not expose them.
var (
	TableScope          = ssot.NewEnum("TableScope", "TableScopes", "GLOBAL", "LOCAL")
	TableKind           = ssot.NewEnum("TableKind", "TableKinds", "TEMPORARY", "VOLATILE", "TRANSIENT")
	TableConstraintType = ssot.NewEnum("TableConstraintType", "TableConstraintTypes", "PRIMARY KEY", "UNIQUE", "FOREIGN KEY")
)

// TableCreateAttrs are the CREATE TABLE clauses that are not resource attributes. They are SqlOnly, so the schema,
// resource, model, and assertion generators skip them.
//
// The column list, the stage file format, and the stage copy options name query structs by PreexistingStruct. The SDK
// definition in pkg/sdk/generator/defs/tables_def.go supplies them to CreateOperationFromObject and keeps their fields,
// because other operations and the request types of the resource use them. The row access policy is generated here as
// TableRowAccessPolicy. CREATE AS SELECT, event tables, and external tables refer to that type.
var TableCreateAttrs = struct {
	Scope                      ssot.Attribute
	Kind                       ssot.Attribute
	ColumnsAndConstraints      ssot.Attribute
	EnableSchemaEvolution      ssot.Attribute
	StageFileFormat            ssot.Attribute
	StageCopyOptions           ssot.Attribute
	MaxDataExtensionTimeInDays ssot.Attribute
	DefaultDdlCollation        ssot.Attribute
	CopyGrants                 ssot.Attribute
	RowAccessPolicy            ssot.Attribute
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
		CreateSql: ssot.PreexistingStruct("CreateTableColumnsAndConstraints", ssot.Parentheses(), ssot.Required()),
	},
	EnableSchemaEvolution: ssot.Attribute{
		Name:    "enable_schema_evolution",
		Type:    ssot.AttrBool,
		SqlOnly: true,
	},
	StageFileFormat: ssot.Attribute{
		Name:      "stage_file_format",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.PreexistingStruct("LegacyFileFormat", ssot.Parentheses(), ssot.NoComma(), ssot.Sql("STAGE_FILE_FORMAT =")),
	},
	StageCopyOptions: ssot.Attribute{
		Name:      "stage_copy_options",
		Type:      ssot.AttrObject,
		SqlOnly:   true,
		CreateSql: ssot.PreexistingStruct("LegacyTableCopyOptions", ssot.Parentheses(), ssot.NoComma(), ssot.Sql("STAGE_COPY_OPTIONS =")),
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
	},
	CopyGrants: ssot.Attribute{
		Name:      "copy_grants",
		Type:      ssot.AttrBool,
		SqlOnly:   true,
		CreateSql: ssot.BoolKeyword(ssot.Sql("COPY GRANTS")),
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
				CreateSql: ssot.StringListKeyword(ssot.Parentheses(), ssot.Sql("ON")),
			},
		},
	},
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
	tableShowCreatedOn           = tableShowString("created_on", false, false)
	tableShowName                = tableShowString("name", false, false)
	tableShowSchemaName          = tableShowString("schema_name", false, false)
	tableShowDatabaseName        = tableShowString("database_name", false, false)
	tableShowKind                = tableShowString("kind", false, false)
	tableShowClusterBy           = ssot.Attribute{Name: "show_cluster_by", ShowName: "cluster_by", Type: ssot.AttrString, Readonly: true, OptionalInDb: true}
	tableShowOwner               = tableShowString("owner", false, false)
	tableShowRows                = tableShowInt("rows", false)
	tableShowBytes               = tableShowInt("bytes", true)
	tableShowRetentionTime       = tableShowInt("retention_time", false)
	tableShowDroppedOn           = tableShowString("dropped_on", true, true)
	tableShowAutomaticClustering = tableShowBool("automatic_clustering", "ON")
	// The resource attribute stays a separate field. Sharing its name would drift-suppress change_tracking.
	tableShowChangeTracking = ssot.Attribute{
		Name: "show_change_tracking", ShowName: "change_tracking", Type: ssot.AttrBool, Readonly: true,
		StoredAsText: true, OptionalInDb: true, BoolTrueValue: "ON",
	}
	tableShowSearchOptimization         = tableShowBool("search_optimization", "ON")
	tableShowSearchOptimizationProgress = tableShowString("search_optimization_progress", true, false)
	tableShowSearchOptimizationBytes    = tableShowInt("search_optimization_bytes", true)
	tableShowIsExternal                 = tableShowBool("is_external", "Y")
	tableShowEnableSchemaEvolution      = tableShowBool("enable_schema_evolution", "Y")
	tableShowOwnerRoleType              = tableShowString("owner_role_type", true, false)
	tableShowIsEvent                    = tableShowBool("is_event", "Y")
	tableShowIsDynamic                  = tableShowBool("is_dynamic", "Y")
	tableShowBudget                     = tableShowString("budget", true, true)
)

var Table = ssot.Object{
	Name:                   "Table",
	NamePlural:             "Tables",
	IdType:                 "SchemaObjectIdentifier",
	LegacyIdEncoding:       true,
	PreviewFeature:         true,
	DescribeFailRead:       true,
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
		ssot.WithSdk(
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
		),
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
			&TableCreateAttrs.StageFileFormat,
			&TableCreateAttrs.StageCopyOptions,
			&TableAttrs.DataRetentionTimeInDays,
			&TableCreateAttrs.MaxDataExtensionTimeInDays,
			&TableAttrs.ChangeTracking,
			&TableCreateAttrs.DefaultDdlCollation,
			&TableCreateAttrs.CopyGrants,
			&TableCreateAttrs.RowAccessPolicy,
			ssot.Tags,
			&TableAttrs.Comment,
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
			&tableShowChangeTracking,
			&tableShowSearchOptimization,
			&tableShowSearchOptimizationProgress,
			&tableShowSearchOptimizationBytes,
			&tableShowIsExternal,
			&tableShowEnableSchemaEvolution,
			&tableShowOwnerRoleType,
			&tableShowIsEvent,
			&tableShowIsDynamic,
			&tableShowBudget,
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
	},
	CustomOperations: tableDescribeOperations,
	Attributes: []ssot.Attribute{
		TableAttrs.ClusterBy,
		TableAttrs.Column,
		TableAttrs.Comment,
		TableAttrs.PrimaryKey,
		TableAttrs.DataRetentionTimeInDays,
		TableAttrs.ChangeTracking,
	},
}
