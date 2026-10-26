package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

var SchemaAttrs = struct {
	WithManagedAccess ssot.Attribute
	IsTransient       ssot.Attribute
	Comment           ssot.Attribute
	CreatedOn         ssot.Attribute
	DroppedOn         ssot.Attribute
	Name              ssot.Attribute
	IsDefault         ssot.Attribute
	IsCurrent         ssot.Attribute
	DatabaseName      ssot.Attribute
	Owner             ssot.Attribute
	Options           ssot.Attribute
	RetentionTime     ssot.Attribute
	OwnerRoleType     ssot.Attribute
	Kind              ssot.Attribute
}{
	WithManagedAccess: ssot.Attribute{
		Name:        "with_managed_access",
		Type:        ssot.AttrBool,
		Description: "Specifies a managed schema. Managed access schemas centralize privilege management with the schema owner.",
		CreateSql:   ssot.BoolKeyword(ssot.Sql("WITH MANAGED ACCESS")),
	},
	IsTransient: ssot.Attribute{
		Name:        "is_transient",
		Type:        ssot.AttrBool,
		ForceNew:    true,
		Description: "Specifies the schema as transient. Transient schemas do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
		CreateSql:   ssot.BoolKeyword(ssot.Sql("TRANSIENT")),
	},
	Comment: ssot.Attribute{
		Name:         "comment",
		Type:         ssot.AttrString,
		Description:  "Specifies a comment for the schema.",
		OptionalInDb: true,
	},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	DroppedOn:     ssot.Attribute{Name: "dropped_on", Type: ssot.AttrTime, OptionalInDb: true, Readonly: true},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	IsDefault:     ssot.Attribute{Name: "is_default", Type: ssot.AttrBool, StoredAsText: true, Readonly: true},
	IsCurrent:     ssot.Attribute{Name: "is_current", Type: ssot.AttrBool, StoredAsText: true, Readonly: true},
	DatabaseName:  ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	Options:       ssot.Attribute{Name: "options", Type: ssot.AttrString, OptionalInShow: true, OptionalInDb: true, Readonly: true},
	RetentionTime: ssot.Attribute{Name: "retention_time", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true},
	Kind:          ssot.Attribute{Name: "kind", Type: ssot.AttrString, Readonly: true},
}

var Schema = ssot.Object{
	Name:                                "Schema",
	NamePlural:                          "Schemas",
	DescribeUsesShowId:                  true,
	HasPartialOnUpdate:                  true,
	AdditionalDescribeOutputDescription: "In order to handle this output, one must grant sufficient privileges, e.g. [grant_ownership](./grant_ownership) on all objects in the schema.",
	IdType:                              "DatabaseObjectIdentifier",
	AdditionalIdentityNameDescription:   "When the name is `PUBLIC`, during creation the provider checks if this schema has already been created and, in such case, `ALTER` is used to match the desired state.",
	AttributesModification:              true,
	Generators: []ssot.GeneratorOption{
		ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow(), ssot.ShowOutputDescribe()),
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
			ssot.ObjectAssertDescribe(ssot.WithNoIdentifiableObject()),
		),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkShow(
				ssot.SdkShowTerse(),
				ssot.SdkShowOptionalSql("HISTORY"),
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalExtendedIn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDExtendedIn, ssot.SdkShowByIDLike),
			),
			ssot.SdkDescribe(
				ssot.SdkDescribeMapping(ssot.SdkDescriptionMappingSlice),
				ssot.SdkDescribeQueryStructName("DescribeSchema"),
			),
			ssot.SdkDrop(ssot.SdkDropCascadeRestrict()),
			ssot.SdkCreate(),
			ssot.SdkAlter(ssot.SdkAlterSetList()),
			ssot.SdkCustomOperations(),
			ssot.SdkCustomInterfaceMethods(),
			ssot.SdkShowParameters(),
			ssot.SdkShowParametersDetails(),
		),
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
	},
	Attributes: []ssot.Attribute{
		SchemaAttrs.WithManagedAccess,
		SchemaAttrs.IsTransient,
		SchemaAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		KeywordsBefore: []*ssot.Attribute{&SchemaAttrs.IsTransient},
		Clauses: []*ssot.Attribute{
			&SchemaAttrs.WithManagedAccess,
			ssot.Parameters,
			&SchemaAttrs.Comment,
			ssot.Tags,
		},
	},
	Set: []*ssot.Attribute{
		ssot.Parameters,
		&SchemaAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		ssot.Parameters,
		&SchemaAttrs.Comment,
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterRenameTo(),
		ssot.AlterSwapWith(),
		ssot.AlterSet(),
		ssot.AlterUnset(),
		ssot.AlterSetTags(),
		ssot.AlterUnsetTags(),
		ssot.AlterKeyword("ENABLE MANAGED ACCESS"),
		ssot.AlterKeyword("DISABLE MANAGED ACCESS"),
	}},
	CustomOperations: []ssot.CustomOperationSpec{
		{
			Kind:                  "Clone",
			Doc:                   "https://docs.snowflake.com/en/sql-reference/sql/create-schema",
			AdditionalValidations: true,
			Clauses: []ssot.CustomOperationClause{
				ssot.CustomStaticSQL("CREATE"),
				ssot.CustomOrReplace(),
				ssot.CustomAttribute(&SchemaAttrs.IsTransient),
				ssot.CustomObjectWords(),
				ssot.CustomIfNotExists(),
				ssot.CustomName(),
				ssot.CustomPredefined("Clone"),
			},
		},
		{
			Kind: "Undrop",
			Clauses: []ssot.CustomOperationClause{
				ssot.CustomStaticSQL("UNDROP"),
				ssot.CustomObjectWords(),
				ssot.CustomName(),
			},
		},
	},
	CustomInterfaceMethods: []ssot.CustomInterfaceMethodSpec{
		{
			Name: "Use",
			Doc:  "Use is based on https://docs.snowflake.com/en/sql-reference/sql/use-schema",
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&SchemaAttrs.CreatedOn,
			&SchemaAttrs.DroppedOn,
			&SchemaAttrs.Name,
			&SchemaAttrs.IsDefault,
			&SchemaAttrs.IsCurrent,
			&SchemaAttrs.DatabaseName,
			&SchemaAttrs.Owner,
			&SchemaAttrs.Comment,
			&SchemaAttrs.Options,
			&SchemaAttrs.RetentionTime,
			&SchemaAttrs.OwnerRoleType,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&SchemaAttrs.CreatedOn,
			&SchemaAttrs.Name,
			&SchemaAttrs.Kind,
		},
	},
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithBeforeCreate(ssot.ExtCall()),
		ssot.WithHierarchyRename(ssot.ExtCall()),
		ssot.WithBeforeAlter(ssot.ExtCall()),
		// Describe is best-effort and its output is a list, so the ext sets it.
		ssot.WithDescribeOutputSet(ssot.Omitted()),
		ssot.WithSetDescribeOutput(ssot.ExtCall()),
		ssot.WithStateUpgraders(ssot.Generated()),
		ssot.WithApplyParametersCreate(ssot.Generated()),
		ssot.WithShowParametersInSdk(ssot.Generated()),
		ssot.WithParametersDetailsFromRaw(ssot.Generated()),
		ssot.WithSetParametersFields(ssot.Generated()),
		ssot.WithParametersOutputSet(ssot.Generated()),
		ssot.WithApplyParametersChanges(ssot.Generated()),
		ssot.WithParametersAttributesSchema(ssot.Generated()),
		ssot.WithParametersOutputSchema(ssot.Generated()),
		ssot.WithParametersFieldNames(ssot.Generated()),
		ssot.WithParametersCustomDiff(ssot.Generated()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
	},
	SchemaStateUpgraders: []ssot.StateUpgraderFunc{
		"v093SchemaStateUpgrader",
		"migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName",
	},
}
