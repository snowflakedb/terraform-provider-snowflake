package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// BackupSetAttrs holds reusable backup set fields.
//
// There is no handwritten SDK definition for this object.
// CREATE OR ALTER is on the create page. This object does not add it.
// DROP and ALTER pages omit IF EXISTS. The SDK step probes those clauses before leaving them out.
// The parser accepts UNSET BACKUP POLICY, FOR ICEBERG TABLE, and FOR DYNAMIC ICEBERG TABLE. No public page shows them.
//
// comment and the backup policy columns are optional because the create page allows omitting both.
// object_database_name and object_schema_name are optional because the pages do not say which columns are
// nullable. A database target has no schema, and those columns are empty for that target. A null scan must not fail.
var BackupSetAttrs = struct {
	ForTable                 ssot.Attribute
	ForDynamicTable          ssot.Attribute
	ForSchema                ssot.Attribute
	ForDatabase              ssot.Attribute
	BackupPolicy             ssot.Attribute
	Comment                  ssot.Attribute
	CreatedOn                ssot.Attribute
	Name                     ssot.Attribute
	DatabaseName             ssot.Attribute
	SchemaName               ssot.Attribute
	ObjectKind               ssot.Attribute
	ObjectName               ssot.Attribute
	ObjectDatabaseName       ssot.Attribute
	ObjectSchemaName         ssot.Attribute
	BackupPolicyName         ssot.Attribute
	BackupPolicyDatabaseName ssot.Attribute
	BackupPolicySchemaName   ssot.Attribute
	BackupPolicyState        ssot.Attribute
	OwnerRole                ssot.Attribute
	OwnerRoleType            ssot.Attribute
}{
	ForTable: ssot.Attribute{
		Name:           "for_table",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "SchemaObjectIdentifier",
		ForceNew:       true,
		Description:    "Specifies the table this backup set backs up (`FOR TABLE`). Exactly one of the four targets is set, and the target cannot be changed.",
		CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("FOR TABLE")),
	},
	ForDynamicTable: ssot.Attribute{
		Name:           "for_dynamic_table",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "SchemaObjectIdentifier",
		ForceNew:       true,
		Description:    "Specifies the dynamic table this backup set backs up (`FOR DYNAMIC TABLE`). Exactly one of the four targets is set, and the target cannot be changed.",
		CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("FOR DYNAMIC TABLE")),
	},
	ForSchema: ssot.Attribute{
		Name:           "for_schema",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "DatabaseObjectIdentifier",
		ForceNew:       true,
		Description:    "Specifies the schema this backup set backs up (`FOR SCHEMA`). Exactly one of the four targets is set, and the target cannot be changed.",
		CreateSql:      ssot.IdentifierKeyword("DatabaseObjectIdentifier", ssot.Sql("FOR SCHEMA")),
	},
	ForDatabase: ssot.Attribute{
		Name:           "for_database",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "AccountObjectIdentifier",
		ForceNew:       true,
		Description:    "Specifies the database this backup set backs up (`FOR DATABASE`). Exactly one of the four targets is set, and the target cannot be changed.",
		CreateSql:      ssot.IdentifierKeyword("AccountObjectIdentifier", ssot.Sql("FOR DATABASE")),
	},
	// The clause text is WITH BACKUP POLICY, so the CREATE field is WithBackupPolicy.
	// Changing it is APPLY BACKUP POLICY, not SET. An attached policy cannot be unset.
	BackupPolicy: ssot.Attribute{
		Name:           "backup_policy",
		Type:           ssot.AttrIdentifier,
		IdentifierType: "SchemaObjectIdentifier",
		Description:    "Specifies the backup policy for the set. It can be omitted at creation and attached later. FORCE overwrites an existing policy only when that policy has no retention lock. An attached policy cannot be removed.",
		CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("WITH BACKUP POLICY")),
	},
	Comment: ssot.Attribute{
		Name:           "comment",
		Type:           ssot.AttrString,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Specifies a comment for the backup set.",
	},
	CreatedOn:    ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true},
	Name:         ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	DatabaseName: ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true},
	SchemaName:   ssot.Attribute{Name: "schema_name", Type: ssot.AttrString, Readonly: true},
	ObjectKind: ssot.Attribute{
		Name:        "object_kind",
		Type:        ssot.AttrString,
		Readonly:    true,
		Description: "Type of the object that the backup set backs up.",
	},
	ObjectName: ssot.Attribute{
		Name:        "object_name",
		Type:        ssot.AttrString,
		Readonly:    true,
		Description: "Name of the object that the backup set backs up.",
	},
	ObjectDatabaseName: ssot.Attribute{
		Name:           "object_database_name",
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Name of the database that contains the object that this backup set backs up.",
	},
	ObjectSchemaName: ssot.Attribute{
		Name:           "object_schema_name",
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Name of the schema that contains the object that this backup set backs up.",
	},
	BackupPolicyName: ssot.Attribute{
		Name:           "backup_policy_name",
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Name of the backup policy attached to this backup set.",
	},
	BackupPolicyDatabaseName: ssot.Attribute{
		Name:           "backup_policy_database_name",
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Name of the database that contains the backup policy.",
	},
	BackupPolicySchemaName: ssot.Attribute{
		Name:           "backup_policy_schema_name",
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Name of the schema that contains the backup policy.",
	},
	BackupPolicyState: ssot.Attribute{
		Name:           "backup_policy_state",
		Type:           ssot.AttrString,
		Readonly:       true,
		OptionalInShow: true,
		OptionalInDb:   true,
		Description:    "Current state of the backup policy attached to this backup set.",
	},
	OwnerRole:     ssot.Attribute{Name: "owner_role", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType: ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, Readonly: true},
}

// APPLY BACKUP POLICY <policy> [FORCE]. The policy identifier is the alter block's own field.
// The resource attribute stays backup_policy; its CREATE clause is WITH BACKUP POLICY.
var backupSetApplyPolicy = ssot.Attribute{
	Name:      "apply_backup_policy",
	Type:      ssot.AttrObject,
	SqlOnly:   true,
	CreateSql: ssot.ObjectParameter(),
	Children: []ssot.Attribute{
		{
			Name:           "policy_name",
			Type:           ssot.AttrIdentifier,
			IdentifierType: "SchemaObjectIdentifier",
			Required:       true,
			CreateSql:      ssot.IdentifierKeyword("SchemaObjectIdentifier", ssot.Sql("BACKUP POLICY"), ssot.NoEquals(), ssot.Required()),
		},
		{
			Name:      "force",
			Type:      ssot.AttrBool,
			CreateSql: ssot.BoolKeyword(ssot.Sql("FORCE")),
		},
	},
}

// MODIFY BACKUP IDENTIFIER '<id>' followed by one action. A list element is how an alter block gets
// variants: the block field is that element struct, not a slice.
var (
	backupSetModifyBackupID = ssot.Attribute{
		Name:        "backup_id",
		Type:        ssot.AttrString,
		Required:    true,
		Description: "Identifier of one backup in the set, a UUID in the form returned by UUID_STRING.",
		CreateSql:   ssot.StringParameter(ssot.Sql("BACKUP IDENTIFIER"), ssot.NoEquals(), ssot.Required()),
	}
	backupSetModifySetComment = ssot.Attribute{
		Name:      "comment",
		Type:      ssot.AttrString,
		Required:  true,
		CreateSql: ssot.StringParameter(ssot.Sql("SET COMMENT"), ssot.Required()),
	}
	backupSetModifyAddLegalHold    = ssot.Attribute{Name: "add_legal_hold", SqlOnly: true, CreateSql: ssot.StaticSQL("addLegalHold", "ADD LEGAL HOLD")}
	backupSetModifyRemoveLegalHold = ssot.Attribute{Name: "remove_legal_hold", SqlOnly: true, CreateSql: ssot.StaticSQL("removeLegalHold", "REMOVE LEGAL HOLD")}
	backupSetModifyUnsetComment    = ssot.Attribute{Name: "unset_comment", SqlOnly: true, CreateSql: ssot.StaticSQL("unsetComment", "UNSET COMMENT")}
	backupSetModifyBackup          = ssot.Attribute{
		Name:    "modify_backup",
		Type:    ssot.AttrList,
		SqlOnly: true,
		Elem: &ssot.Attribute{
			Type: ssot.AttrObject,
			CreateSql: ssot.ObjectParameter(ssot.WithVariants(
				ssot.Variant("AddLegalHold", &backupSetModifyAddLegalHold),
				ssot.Variant("RemoveLegalHold", &backupSetModifyRemoveLegalHold),
				ssot.Variant("SetComment", &backupSetModifySetComment),
				ssot.Variant("UnsetComment", &backupSetModifyUnsetComment),
			)),
			Children: []ssot.Attribute{
				backupSetModifyBackupID,
				backupSetModifySetComment,
			},
		},
	}
)

// BackupSet is one schema object. The four FOR targets are mutually exclusive and permanent.
// SHOW and DESCRIBE publish the same columns, in this order.
var BackupSet = ssot.Object{
	Name:                      "BackupSet",
	NamePlural:                "BackupSets",
	IdType:                    "SchemaObjectIdentifier",
	IdentityHierarchyForceNew: true,
	// DROP and ALTER pages omit IF EXISTS. Both stay on until a probe shows Snowflake rejects the clause.
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
			ssot.SdkAlter(),
			ssot.SdkDescribe(),
		),
	},
	ExactlyOneOf: [][]*ssot.Attribute{{
		&BackupSetAttrs.ForTable,
		&BackupSetAttrs.ForDynamicTable,
		&BackupSetAttrs.ForSchema,
		&BackupSetAttrs.ForDatabase,
	}},
	Attributes: []ssot.Attribute{
		BackupSetAttrs.ForTable,
		BackupSetAttrs.ForDynamicTable,
		BackupSetAttrs.ForSchema,
		BackupSetAttrs.ForDatabase,
		BackupSetAttrs.BackupPolicy,
		BackupSetAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/create-backup-set",
		Clauses: []*ssot.Attribute{
			&BackupSetAttrs.ForTable,
			&BackupSetAttrs.ForDynamicTable,
			&BackupSetAttrs.ForSchema,
			&BackupSetAttrs.ForDatabase,
			ssot.Tags,
			&BackupSetAttrs.BackupPolicy,
			&BackupSetAttrs.Comment,
		},
	},
	Set: []*ssot.Attribute{
		&BackupSetAttrs.Comment,
	},
	Unset: []*ssot.Attribute{
		&BackupSetAttrs.Comment,
	},
	Alter: ssot.AlterSpec{
		Doc: "https://docs.snowflake.com/en/sql-reference/sql/alter-backup-set",
		Clauses: []ssot.AlterClause{
			ssot.AlterKeyword("ADD BACKUP"),
			ssot.AlterKeywordBlock("ApplyBackupPolicy", "APPLY", &backupSetApplyPolicy),
			ssot.AlterKeyword("SUSPEND BACKUP POLICY"),
			ssot.AlterKeyword("SUSPEND BACKUP CREATION POLICY"),
			ssot.AlterKeyword("SUSPEND BACKUP EXPIRATION POLICY"),
			ssot.AlterKeyword("RESUME BACKUP POLICY"),
			ssot.AlterKeyword("RESUME BACKUP CREATION POLICY"),
			ssot.AlterKeyword("RESUME BACKUP EXPIRATION POLICY"),
			ssot.AlterText("DELETE BACKUP IDENTIFIER"),
			ssot.AlterKeywordBlock("ModifyBackup", "MODIFY", &backupSetModifyBackup),
			ssot.AlterRenameTo(),
			ssot.AlterSet(),
			ssot.AlterUnset(),
			ssot.AlterSetTags(),
			ssot.AlterUnsetTags(),
		},
	},
	Show: ssot.OutputSpec{Columns: backupSetColumns()},
	Describe: ssot.OutputSpec{
		Columns: backupSetColumns(),
	},
}

func backupSetColumns() []*ssot.Attribute {
	return []*ssot.Attribute{
		&BackupSetAttrs.CreatedOn,
		&BackupSetAttrs.Name,
		&BackupSetAttrs.DatabaseName,
		&BackupSetAttrs.SchemaName,
		&BackupSetAttrs.ObjectKind,
		&BackupSetAttrs.ObjectName,
		&BackupSetAttrs.ObjectDatabaseName,
		&BackupSetAttrs.ObjectSchemaName,
		&BackupSetAttrs.BackupPolicyName,
		&BackupSetAttrs.BackupPolicyDatabaseName,
		&BackupSetAttrs.BackupPolicySchemaName,
		&BackupSetAttrs.BackupPolicyState,
		&BackupSetAttrs.OwnerRole,
		&BackupSetAttrs.OwnerRoleType,
		&BackupSetAttrs.Comment,
	}
}
