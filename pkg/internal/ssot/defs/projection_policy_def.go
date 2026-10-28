package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// ProjectionPolicyAttrs holds reusable projection policy fields.
//
// SHOW and DESCRIBE column lists come from a probe on 2026-10-28. The public pages have no column table.
// SHOW returned created_on, name, database_name, schema_name, kind, owner, comment, owner_role_type, options.
// comment and options were empty strings when unset, not null.
// DESCRIBE returned name, signature, return_type, body. It did not return created_on.
// signature is always () and return_type is always PROJECTION_CONSTRAINT.
// CREATE emits one clause, AS () RETURNS PROJECTION_CONSTRAINT -> <body>. body is the expression.
// CREATE accepts TAG even though the CREATE page omits it. The SDK emits TAG, which Snowflake also accepts as WITH TAG.
// SHOW accepts STARTS WITH and LIMIT after IN. DROP accepts IF EXISTS. Those clauses are omitted from the syntax blocks.
var ProjectionPolicyAttrs = struct {
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
		Description: "SQL expression that determines whether to project a column. The expression must call PROJECTION_CONSTRAINT.",
		CreateSql:   ssot.StringParameter(ssot.Sql("AS () RETURNS PROJECTION_CONSTRAINT ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.Required(), ssot.FieldName("Body")),
		SetSql:      ssot.StringParameter(ssot.Sql("BODY ->"), ssot.NoEquals(), ssot.NoQuotes(), ssot.FieldName("SetBody")),
	},
	Comment: ssot.Attribute{
		Name:        "comment",
		Type:        ssot.AttrString,
		Description: "Specifies a comment for the projection policy.",
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

var ProjectionPolicy = ssot.Object{
	Name:                      "ProjectionPolicy",
	NamePlural:                "ProjectionPolicies",
	IdType:                    "SchemaObjectIdentifier",
	IdentityHierarchyForceNew: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkAlter(),
			ssot.SdkDrop(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalExtendedIn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDLike, ssot.SdkShowByIDExtendedIn),
			),
			ssot.SdkDescribe(),
		),
	},
	Attributes: []ssot.Attribute{
		ProjectionPolicyAttrs.Body,
		ProjectionPolicyAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&ProjectionPolicyAttrs.Body,
			&ProjectionPolicyAttrs.Comment,
			ssot.Tags,
		},
	},
	// SET BODY -> and SET COMMENT cannot be combined. A probe on 2026-10-28 rejected both in one SET.
	Set: []*ssot.Attribute{
		&ProjectionPolicyAttrs.Body,
		&ProjectionPolicyAttrs.Comment,
	},
	Unset: []*ssot.Attribute{&ProjectionPolicyAttrs.Comment},
	Alter: ssot.AlterSpec{
		Clauses: []ssot.AlterClause{
			ssot.AlterRenameTo(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
		},
		SetConflictingFields: [][]*ssot.Attribute{{
			&ProjectionPolicyAttrs.Body,
			&ProjectionPolicyAttrs.Comment,
		}},
	},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ProjectionPolicyAttrs.CreatedOn,
			&ProjectionPolicyAttrs.Name,
			&ProjectionPolicyAttrs.DatabaseName,
			&ProjectionPolicyAttrs.SchemaName,
			&ProjectionPolicyAttrs.Kind,
			&ProjectionPolicyAttrs.Owner,
			&ProjectionPolicyAttrs.Comment,
			&ProjectionPolicyAttrs.OwnerRoleType,
			&ProjectionPolicyAttrs.Options,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&ProjectionPolicyAttrs.Name,
			&ProjectionPolicyAttrs.Signature,
			&ProjectionPolicyAttrs.ReturnType,
			&ProjectionPolicyAttrs.Body,
		},
	},
}
