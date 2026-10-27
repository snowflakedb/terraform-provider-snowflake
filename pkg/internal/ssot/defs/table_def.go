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
		Name:        "cluster_by",
		Type:        ssot.AttrList,
		Description: "A list of one or more table columns/expressions to be used as clustering key(s) for the table",
		Elem:        &ssot.Attribute{Type: ssot.AttrString},
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
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the table.",
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

var Table = ssot.Object{
	Name:                   "Table",
	NamePlural:             "Tables",
	IdType:                 "SchemaObjectIdentifier",
	AttributesModification: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithResource(),
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithShowOutputSchema(ssot.Omitted()),
		ssot.WithDescribeOutputSchema(ssot.Omitted()),
	},
	Attributes: []ssot.Attribute{
		TableAttrs.ClusterBy,
		TableAttrs.Column,
		TableAttrs.Comment,
		TableAttrs.PrimaryKey,
		TableAttrs.DataRetentionTimeInDays,
		TableAttrs.ChangeTracking,
	},
}
