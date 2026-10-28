package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// BackupPolicyAttrs holds reusable backup policy fields.
var BackupPolicyAttrs = struct {
	HasRetentionLock ssot.Attribute
	Schedule         ssot.Attribute
	ExpireAfterDays  ssot.Attribute
	Comment          ssot.Attribute
	CreatedOn        ssot.Attribute
	Name             ssot.Attribute
	DatabaseName     ssot.Attribute
	SchemaName       ssot.Attribute
	Owner            ssot.Attribute
	OwnerRoleType    ssot.Attribute
}{
	// The pages say Y or N. SHOW and DESCRIBE return true or false.
	HasRetentionLock: ssot.Attribute{
		Name:         "has_retention_lock",
		Type:         ssot.AttrBool,
		ForceNew:     true,
		StoredAsText: true,
		BoolParsed:   true,
		Description:  "Specifies a retention lock. Backups created under the policy cannot be deleted before they expire, including by a privileged user. The lock cannot be changed on an existing backup policy.",
		CreateSql:    ssot.BoolKeyword(ssot.Sql("WITH RETENTION LOCK")),
	},
	Schedule: ssot.Attribute{
		Name:           "schedule",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Specifies how often Snowflake creates a backup. The value is an interval of at least 60 minutes (`<num> MINUTE` or `<num> HOUR`) or `USING CRON <expr> <time_zone>`.",
	},
	// The maximum of 3653 applies only when SCHEDULE is omitted, so Max is unset.
	ExpireAfterDays: ssot.Attribute{
		Name:           "expire_after_days",
		Type:           ssot.AttrInt,
		Min:            new(1),
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Specifies the number of days until a backup expires. The minimum is 1. The maximum is 3653 when SCHEDULE is omitted.",
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Specifies a comment for the backup policy.",
	},
	CreatedOn:     ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Name:          ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName:  ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:    ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	Owner:         ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true},
}

// SHOW and DESCRIBE list owner twice. The second row is the owner of a backup set, not a column.
// The create page documents CREATE OR ALTER. This object does not.
// The DROP page omits IF EXISTS. A probe on 2026-10-27 accepted the clause, so DROP keeps it.
var BackupPolicy = ssot.Object{
	Name:                      "BackupPolicy",
	NamePlural:                "BackupPolicies",
	IdType:                    "SchemaObjectIdentifier",
	IdentityHierarchyForceNew: true,
	AtLeastOneOf: [][]*ssot.Attribute{{
		&BackupPolicyAttrs.Schedule,
		&BackupPolicyAttrs.ExpireAfterDays,
	}},
	Generators: []ssot.GeneratorOption{
		ssot.WithSdk(
			ssot.SdkDefineFromSsot(),
			ssot.SdkCreate(),
			ssot.SdkShow(
				ssot.SdkShowOptionalLike(),
				ssot.SdkShowOptionalIn(),
				ssot.SdkShowOptionalStartsWith(),
				ssot.SdkShowOptionalLimit(),
				ssot.SdkShowByIDFiltering(ssot.SdkShowByIDIn, ssot.SdkShowByIDLike),
			),
			ssot.SdkDrop(),
			ssot.SdkAlter(ssot.SdkAlterWithoutIfExists()),
			ssot.SdkDescribe(),
		),
		ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
	},
	Attributes: []ssot.Attribute{
		BackupPolicyAttrs.HasRetentionLock,
		BackupPolicyAttrs.Schedule,
		BackupPolicyAttrs.ExpireAfterDays,
		BackupPolicyAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			ssot.Tags,
			&BackupPolicyAttrs.HasRetentionLock,
			&BackupPolicyAttrs.Schedule,
			&BackupPolicyAttrs.ExpireAfterDays,
			&BackupPolicyAttrs.Comment,
		},
	},
	Set: []*ssot.Attribute{
		&BackupPolicyAttrs.Comment,
		&BackupPolicyAttrs.Schedule,
		&BackupPolicyAttrs.ExpireAfterDays,
	},
	Unset: []*ssot.Attribute{
		&BackupPolicyAttrs.Comment,
		&BackupPolicyAttrs.Schedule,
		&BackupPolicyAttrs.ExpireAfterDays,
	},
	UnsetConflicts: [][]*ssot.Attribute{{
		&BackupPolicyAttrs.Schedule,
		&BackupPolicyAttrs.ExpireAfterDays,
	}},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterRenameTo(),
		ssot.AlterSet(),
		ssot.AlterUnset(),
		ssot.AlterSetTags(),
		ssot.AlterUnsetTags(),
	}},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&BackupPolicyAttrs.CreatedOn,
			&BackupPolicyAttrs.Name,
			&BackupPolicyAttrs.DatabaseName,
			&BackupPolicyAttrs.SchemaName,
			&BackupPolicyAttrs.Owner,
			&BackupPolicyAttrs.Comment,
			&BackupPolicyAttrs.Schedule,
			&BackupPolicyAttrs.ExpireAfterDays,
			&BackupPolicyAttrs.HasRetentionLock,
			&BackupPolicyAttrs.OwnerRoleType,
		},
	},
	Describe: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&BackupPolicyAttrs.CreatedOn,
			&BackupPolicyAttrs.Name,
			&BackupPolicyAttrs.DatabaseName,
			&BackupPolicyAttrs.SchemaName,
			&BackupPolicyAttrs.Owner,
			&BackupPolicyAttrs.Comment,
			&BackupPolicyAttrs.Schedule,
			&BackupPolicyAttrs.ExpireAfterDays,
			&BackupPolicyAttrs.HasRetentionLock,
			&BackupPolicyAttrs.OwnerRoleType,
		},
	},
}
