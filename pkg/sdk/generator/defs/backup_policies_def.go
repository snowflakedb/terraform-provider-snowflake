package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var backupPoliciesDef = g.NewInterface(
	"BackupPolicies",
	"BackupPolicy",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-backup-policy",
		g.NewQueryStruct("CreateBackupPolicy").
			Create().
			OrReplace().
			SQL("BACKUP POLICY").
			IfNotExists().
			Name().
			OptionalTags().
			OptionalSQL("WITH RETENTION LOCK").
			OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
			OptionalNumberAssignment("EXPIRE_AFTER_DAYS", g.ParameterOptions()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.AtLeastOneValueSet, "Schedule", "ExpireAfterDays"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-backup-policy",
		g.NewQueryStruct("AlterBackupPolicy").
			Alter().
			SQL("BACKUP POLICY").
			Name().
			RenameTo().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("BackupPolicySet").
					OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
					OptionalNumberAssignment("EXPIRE_AFTER_DAYS", g.ParameterOptions()).
					OptionalComment().
					WithValidation(g.AtLeastOneValueSet, "Schedule", "ExpireAfterDays", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalSetTags().
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("BackupPolicyUnset").
					OptionalSQL("SCHEDULE").
					OptionalSQL("EXPIRE_AFTER_DAYS").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "Schedule", "ExpireAfterDays", "Comment"),
				g.ListOptions().NoParentheses().SQL("UNSET"),
			).
			OptionalUnsetTags().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "SetTags", "Unset", "UnsetTags"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-backup-policy",
		g.NewQueryStruct("DropBackupPolicy").
			Drop().
			SQL("BACKUP POLICY").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-backup-policies",
		g.StructPair("backupPolicyDBRow", "BackupPolicy").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			Text("owner_role_type").
			OptionalText("comment").
			OptionalText("schedule").
			OptionalNumber("expire_after_days").
			BoolFromText("has_retention_lock", g.WithBoolParsed()),
		g.NewQueryStruct("ShowBackupPolicies").
			Show().
			SQL("BACKUP POLICIES").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimit(),
		g.ShowByIDInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-backup-policy",
		g.StructPair("backupPolicyDetailsRow", "BackupPolicyDetails").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			Text("owner_role_type").
			OptionalText("comment").
			OptionalText("schedule").
			OptionalNumber("expire_after_days").
			BoolFromText("has_retention_lock", g.WithBoolParsed()),
		g.NewQueryStruct("DescribeBackupPolicy").
			Describe().
			SQL("BACKUP POLICY").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
