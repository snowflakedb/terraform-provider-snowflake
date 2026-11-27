package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// AggregationPolicyAttrs holds reusable aggregation policy fields.
//
// SHOW and DESCRIBE column lists come from a probe on 2026-11-24. The public pages have no column table.
// SHOW returned created_on, name, database_name, schema_name, kind, owner, comment, owner_role_type, options.
// comment and options were empty strings when unset, not null. kind is always AGGREGATION_POLICY.
// DESCRIBE returned name, signature, return_type, body. It did not return created_on, database_name, or schema_name,
// although the DESCRIBE page says it describes the creation date.
// DatabaseName and SchemaName are plain-only fields. DescribeDetails sets them from the identifier so ID() can be built.
// signature is always () and return_type is always AGGREGATION_CONSTRAINT.
// CREATE emits one clause, AS () RETURNS AGGREGATION_CONSTRAINT -> <body>. body is the expression.
// DESCRIBE returns body as it was sent, without changing case or whitespace.
// CREATE accepts TAG and WITH TAG even though the CREATE page omits them. The SDK emits TAG.
// SHOW accepts STARTS WITH and LIMIT after IN. LIMIT ... FROM is accepted, but FROM is ignored.
// Those clauses are omitted from the syntax block. The usage notes mention LIMIT only.
// DROP accepts IF EXISTS, and ALTER accepts IF EXISTS on SET TAG and UNSET TAG. The syntax blocks omit both.
var AggregationPolicyAttrs = struct {
	Body          ssot.Attribute
	Comment       ssot.Attribute
	CreatedOn     ssot.Attribute
	Name          ssot.Attribute
	DatabaseName  ssot.Attribute
	SchemaName    ssot.Attribute
	Kind          ssot.Attribute
	Owner         ssot.Attribute
	OwnerRoleType ssot.Attribute
	Options       ssot.Attribute
	Signature     ssot.Attribute
	ReturnType    ssot.Attribute
}{
	Body: ssot.Attribute{
		Name:         "body",
		Type:         ssot.AttrString,
		Required:     true,
		Description:  "SQL expression that determines the restrictions of an aggregation policy. The expression must call AGGREGATION_CONSTRAINT or NO_AGGREGATION_CONSTRAINT.",
		DiffSuppress: ssot.SuppressStatement,
		CreateSql:    ssot.StringParameter(ssot.Sql("AS () RETURNS AGGREGATION_CONSTRAINT ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.Required(), ssot.FieldName("Body")),
		SetSql:       ssot.StringParameter(ssot.Sql("BODY ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.FieldName("SetBody")),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the aggregation policy.",
	},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:  ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:    ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Kind:          ssot.Attribute{Name: "kind", Type: ssot.AttrString, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true},
	Options:       ssot.Attribute{Name: "options", Type: ssot.AttrString, Readonly: true},
	Signature:     ssot.Attribute{Name: "signature", Type: ssot.AttrString, Readonly: true},
	ReturnType:    ssot.Attribute{Name: "return_type", Type: ssot.AttrString, Readonly: true},
}

var AggregationPolicy = ssot.Object{
	Name:             "AggregationPolicy",
	NamePlural:       "AggregationPolicies",
	DescribeFailRead: true,
	IdType:           "SchemaObjectIdentifier",
	PreviewFeature:   true,
	Generators: []ssot.GeneratorOption{
		ssot.WithAllGenerators(),
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkAlter(),
			ssot.SdkDrop(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalIn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike, ssot.SdkShowByIDIn),
			),
			ssot.SdkDescribe(
				ssot.SdkDescribePlainOnlyField("DatabaseName", "string"),
				ssot.SdkDescribePlainOnlyField("SchemaName", "string"),
				ssot.SdkDescribeDetailsFromIdentifier(),
			),
		),
	},
	// body is SET BODY -> and cannot share a SET with comment, so the update sends it in its own SET (the Body group).
	// SHOW does not return body. Read takes it from DESCRIBE.
	ResourceHooks: []ssot.ResourceHookOption{
		ssot.WithParseIdFromConfig(ssot.Generated()),
		ssot.WithParseId(ssot.Generated()),
		ssot.WithNewCreateRequest(ssot.Generated()),
		ssot.WithApplyCreateOptionals(ssot.Generated()),
		ssot.WithCreateInSdk(ssot.Generated()),
		ssot.WithShowByIdInSdk(ssot.Generated()),
		ssot.WithShowByIdSafelyInSdk(ssot.Generated()),
		ssot.WithDescribeInSdk(ssot.Generated()),
		ssot.WithDropSafelyInSdk(ssot.Generated()),
		ssot.WithSetConfigFields(ssot.Generated()),
		ssot.WithSetDescribeOutput(ssot.Generated()),
		ssot.WithSetFieldsNotSetByRead(ssot.Generated()),
		ssot.WithHierarchyRename(ssot.Generated()),
		ssot.WithApplySetUnset(ssot.Generated()),
		ssot.WithAlterSetInSdk(ssot.Generated()),
		ssot.WithAlterUnsetInSdk(ssot.Generated()),
		ssot.WithHandleExternalChanges(ssot.Omitted()),
		ssot.WithCustomizeDiff(ssot.Generated()),
	},
	Attributes: []ssot.Attribute{
		AggregationPolicyAttrs.Body,
		AggregationPolicyAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&AggregationPolicyAttrs.Body,
			&AggregationPolicyAttrs.Comment,
			ssot.Tags,
		},
	},
	// SET BODY -> and SET COMMENT cannot be combined. A probe on 2026-11-24 rejected both orders in one SET.
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterSet(
				ssot.AlterGroup("Body", &AggregationPolicyAttrs.Body),
				ssot.AlterGroup(ssot.MainAlterGroup, &AggregationPolicyAttrs.Comment),
			),
			ssot.AlterUnset(&AggregationPolicyAttrs.Comment),
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AggregationPolicyAttrs.CreatedOn,
			&AggregationPolicyAttrs.Name,
			&AggregationPolicyAttrs.DatabaseName,
			&AggregationPolicyAttrs.SchemaName,
			&AggregationPolicyAttrs.Kind,
			&AggregationPolicyAttrs.Owner,
			&AggregationPolicyAttrs.Comment,
			&AggregationPolicyAttrs.OwnerRoleType,
			&AggregationPolicyAttrs.Options,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&AggregationPolicyAttrs.Name,
			&AggregationPolicyAttrs.Signature,
			&AggregationPolicyAttrs.ReturnType,
			&AggregationPolicyAttrs.Body,
		},
	},
}
