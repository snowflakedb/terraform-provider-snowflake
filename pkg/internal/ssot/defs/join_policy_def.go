package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// JoinPolicyAttrs holds reusable join policy fields.
//
// The SHOW and DESCRIBE pages show the columns only in an example output, with no column table. A probe on 2026-11-24
// returned the same columns in the same order.
// SHOW returned created_on, name, database_name, schema_name, kind, owner, comment, owner_role_type, options.
// comment and options were empty strings when unset, not null. kind was JOIN_POLICY.
// DESCRIBE returned name, signature, return_type, body. It did not return created_on, although the DESCRIBE page says
// it describes the creation date. It did not return database_name or schema_name either.
// DatabaseName and SchemaName are plain-only fields. DescribeDetails sets them from the identifier so ID() can be built.
// signature is always () and return_type is always JOIN_CONSTRAINT.
// CREATE emits one clause, AS () RETURNS JOIN_CONSTRAINT -> <body>. body is the expression.
// CREATE accepts TAG even though the CREATE page omits it. The SDK emits TAG, which Snowflake also accepts as WITH TAG.
// SHOW accepts STARTS WITH and LIMIT after IN. The SHOW page mentions LIMIT only in its usage notes and omits both
// from the syntax block. LIMIT ... FROM is accepted, but FROM did not move the start of the result.
// DROP accepts IF EXISTS, which the DROP page omits. ALTER accepts IF EXISTS on SET TAG and UNSET TAG, which the ALTER
// page shows only on the other clauses.
var JoinPolicyAttrs = struct {
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
		Name:        "body",
		Type:        ssot.AttrString,
		Required:    true,
		Description: "SQL expression that determines the restrictions of a join policy. The expression must call JOIN_CONSTRAINT.",
		CreateSql:   ssot.StringParameter(ssot.Sql("AS () RETURNS JOIN_CONSTRAINT ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.Required(), ssot.FieldName("Body")),
		SetSql:      ssot.StringParameter(ssot.Sql("BODY ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.FieldName("SetBody")),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the join policy.",
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

var JoinPolicy = ssot.Object{
	Name:             "JoinPolicy",
	NamePlural:       "JoinPolicies",
	DescribeFailRead: true,
	IdType:           "SchemaObjectIdentifier",
	Generators: []ssot.GeneratorOption{
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
		ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
		ssot.WithObjectParametersAssert(),
		ssot.WithResourceParametersAssert(),
		ssot.WithSchema(ssot.SchemaShow(), ssot.SchemaDescribe()),
		ssot.WithResource(),
		ssot.WithModel(),
		ssot.WithResourceAssert(),
		ssot.WithShowOutputAssert(ssot.ShowOutputShow(), ssot.ShowOutputDescribe()),
	},
	// body is SET BODY -> and cannot share a SET with comment, so the update sends it in its own SET (the Body group).
	// The request field is SetBody, so the update mapping names it.
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
		JoinPolicyAttrs.Body,
		JoinPolicyAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&JoinPolicyAttrs.Body,
			&JoinPolicyAttrs.Comment,
			ssot.Tags,
		},
	},
	// SET BODY -> and SET COMMENT cannot be combined. A probe on 2026-11-24 rejected both in one SET, in either order.
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterSet(
				ssot.AlterGroup("Body", &JoinPolicyAttrs.Body),
				ssot.AlterGroup(ssot.MainAlterGroup, &JoinPolicyAttrs.Comment),
			),
			ssot.AlterUnset(&JoinPolicyAttrs.Comment),
		},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&JoinPolicyAttrs.CreatedOn,
			&JoinPolicyAttrs.Name,
			&JoinPolicyAttrs.DatabaseName,
			&JoinPolicyAttrs.SchemaName,
			&JoinPolicyAttrs.Kind,
			&JoinPolicyAttrs.Owner,
			&JoinPolicyAttrs.Comment,
			&JoinPolicyAttrs.OwnerRoleType,
			&JoinPolicyAttrs.Options,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&JoinPolicyAttrs.Name,
			&JoinPolicyAttrs.Signature,
			&JoinPolicyAttrs.ReturnType,
			&JoinPolicyAttrs.Body,
		},
	},
}
