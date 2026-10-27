package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// CodeBundleAttrs holds reusable code bundle fields.
//
// created_on, name, database_name, schema_name, and owner are SHOW columns only.
// A Readonly attribute must not sit on Attributes.
// DESCRIBE columns come from an uncommitted probe run on 2026-10-26. The user guide has the command and no column list.
// Describe-only strings that came back null are optional in SHOW and in the DB row.
// name stays a plain string because the SHOW row's ID() is built from it.
// created_on is optional only in the DB row, because an optional plain time is not mapped.
// ADD VERSION FROM appends a version. It is an ALTER clause, not an attribute.
// AlterText quotes the source. The user-guide example is unquoted; the integration test checks which form Snowflake accepts.
var CodeBundleAttrs = struct {
	From                            ssot.Attribute
	Comment                         ssot.Attribute
	CreatedOn                       ssot.Attribute
	Name                            ssot.Attribute
	DatabaseName                    ssot.Attribute
	SchemaName                      ssot.Attribute
	Owner                           ssot.Attribute
	DefaultVersion                  ssot.Attribute
	DefaultVersionName              ssot.Attribute
	DefaultVersionAlias             ssot.Attribute
	DefaultVersionLocationUri       ssot.Attribute
	DefaultVersionSourceLocationUri ssot.Attribute
	DefaultVersionGitCommitHash     ssot.Attribute
	LastVersionName                 ssot.Attribute
	LastVersionAlias                ssot.Attribute
	LastVersionLocationUri          ssot.Attribute
	LastVersionSourceLocationUri    ssot.Attribute
	LastVersionGitCommitHash        ssot.Attribute
	LiveVersionLocationUri          ssot.Attribute
	OwnerRole                       ssot.Attribute
	OwnerRoleType                   ssot.Attribute
	ContainerType                   ssot.Attribute
	ContainerDatabase               ssot.Attribute
	ContainerSchema                 ssot.Attribute
	ContainerName                   ssot.Attribute
}{
	From: ssot.Attribute{
		Name:        "from",
		Type:        ssot.AttrString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the source that backs this code bundle: a snow://workspace/... URL, or an internal or temporary stage reference. External stages are not supported.",
		CreateSql:   ssot.StringParameter(ssot.Sql("FROM"), ssot.NoEquals(), ssot.Required()),
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		Type:           ssot.AttrString,
		ForceNew:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Specifies a comment for the code bundle.",
	},
	CreatedOn:                       ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true, OptionalInDb: true},
	Name:                            ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:                    ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:                      ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Owner:                           ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	DefaultVersion:                  codeBundleDescribeText("default_version"),
	DefaultVersionName:              codeBundleDescribeText("default_version_name"),
	DefaultVersionAlias:             codeBundleDescribeText("default_version_alias"),
	DefaultVersionLocationUri:       codeBundleDescribeText("default_version_location_uri"),
	DefaultVersionSourceLocationUri: codeBundleDescribeText("default_version_source_location_uri"),
	DefaultVersionGitCommitHash:     codeBundleDescribeText("default_version_git_commit_hash"),
	LastVersionName:                 codeBundleDescribeText("last_version_name"),
	LastVersionAlias:                codeBundleDescribeText("last_version_alias"),
	LastVersionLocationUri:          codeBundleDescribeText("last_version_location_uri"),
	LastVersionSourceLocationUri:    codeBundleDescribeText("last_version_source_location_uri"),
	LastVersionGitCommitHash:        codeBundleDescribeText("last_version_git_commit_hash"),
	LiveVersionLocationUri:          codeBundleDescribeText("live_version_location_uri"),
	OwnerRole:                       codeBundleDescribeText("owner_role"),
	OwnerRoleType:                   codeBundleDescribeText("owner_role_type"),
	ContainerType:                   codeBundleDescribeText("container_type"),
	ContainerDatabase:               codeBundleDescribeText("container_database"),
	ContainerSchema:                 codeBundleDescribeText("container_schema"),
	ContainerName:                   codeBundleDescribeText("container_name"),
}

func codeBundleDescribeText(name string) ssot.Attribute {
	return ssot.Attribute{
		Name:           name,
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
	}
}

var CodeBundle = ssot.Object{
	Name:                      "CodeBundle",
	NamePlural:                "CodeBundles",
	IdType:                    "SchemaObjectIdentifier",
	IdentityNameForceNew:      true,
	IdentityHierarchyForceNew: true,
	Generators: []ssot.GeneratorOption{
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkShow(
				ssot.SdkShowOptionalExtendedIn(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDExtendedIn),
			),
			ssot.SdkDrop(
				ssot.SdkDropDoc("https://docs.snowflake.com/en/developer-guide/code-bundles/code-bundles#drop-code-bundle"),
				ssot.SdkDropSafelyHook(),
			),
			ssot.SdkAlter(
				ssot.SdkAlterDoc("https://docs.snowflake.com/en/developer-guide/code-bundles/code-bundles#alter-code-bundle"),
				ssot.SdkAlterWithoutIfExists(),
			),
			ssot.SdkDescribe(
				ssot.SdkDescribeDoc("https://docs.snowflake.com/en/developer-guide/code-bundles/code-bundles#describe-code-bundle"),
			),
		),
	},
	Attributes: []ssot.Attribute{
		CodeBundleAttrs.From,
		CodeBundleAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&CodeBundleAttrs.From,
			&CodeBundleAttrs.Comment,
		},
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterText("ADD VERSION FROM"),
	}},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CodeBundleAttrs.CreatedOn,
			&CodeBundleAttrs.Name,
			&CodeBundleAttrs.DatabaseName,
			&CodeBundleAttrs.SchemaName,
			&CodeBundleAttrs.Owner,
			&CodeBundleAttrs.Comment,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&CodeBundleAttrs.Name,
			&CodeBundleAttrs.Comment,
			&CodeBundleAttrs.DefaultVersion,
			&CodeBundleAttrs.DefaultVersionName,
			&CodeBundleAttrs.DefaultVersionAlias,
			&CodeBundleAttrs.DefaultVersionLocationUri,
			&CodeBundleAttrs.DefaultVersionSourceLocationUri,
			&CodeBundleAttrs.DefaultVersionGitCommitHash,
			&CodeBundleAttrs.LastVersionName,
			&CodeBundleAttrs.LastVersionAlias,
			&CodeBundleAttrs.LastVersionLocationUri,
			&CodeBundleAttrs.LastVersionSourceLocationUri,
			&CodeBundleAttrs.LastVersionGitCommitHash,
			&CodeBundleAttrs.LiveVersionLocationUri,
			&CodeBundleAttrs.CreatedOn,
			&CodeBundleAttrs.OwnerRole,
			&CodeBundleAttrs.OwnerRoleType,
			&CodeBundleAttrs.ContainerType,
			&CodeBundleAttrs.ContainerDatabase,
			&CodeBundleAttrs.ContainerSchema,
			&CodeBundleAttrs.ContainerName,
		},
	},
}
