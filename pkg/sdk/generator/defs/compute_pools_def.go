package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var ComputePoolBackupInstanceFamilyListItemDef = g.NewQueryStruct("ComputePoolBackupInstanceFamilyListItem").
	Enum("Value", ssotdefs.ComputePoolInstanceFamily, g.KeywordOptions().SingleQuotes().Required())

var computePoolsDef = g.NewInterface(
	"ComputePools",
	"ComputePool",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool",
	g.NewQueryStruct("CreateComputePool").
		Create().
		SQL("COMPUTE POOL").
		// Note: Currently, OR REPLACE is not supported for compute pools.
		IfNotExists().
		Name().
		OptionalIdentifier("ForApplication", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FOR APPLICATION")).
		NumberAssignment("MIN_NODES", g.ParameterOptions().Required()).
		NumberAssignment("MAX_NODES", g.ParameterOptions().Required()).
		EnumAssignment("INSTANCE_FAMILY", ssotdefs.ComputePoolInstanceFamily, g.ParameterOptions().NoQuotes().Required()).
		OptionalBooleanAssignment("AUTO_RESUME", g.ParameterOptions()).
		OptionalBooleanAssignment("INITIALLY_SUSPENDED", g.ParameterOptions()).
		OptionalNumberAssignment("AUTO_SUSPEND_SECS", g.ParameterOptions()).
		OptionalTags().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListAssignment("BACKUP_INSTANCE_FAMILIES", "ComputePoolBackupInstanceFamilyListItem", g.ParameterOptions().Parentheses()).
		WithValidation(g.ValidIdentifier, "name").
		WithAdditionalValidations(),
	ComputePoolBackupInstanceFamilyListItemDef,
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-compute-pool",
	g.NewQueryStruct("AlterComputePool").
		Alter().
		SQL("COMPUTE POOL").
		IfExists().
		Name().
		OptionalSQL("RESUME").
		OptionalSQL("SUSPEND").
		OptionalSQL("STOP ALL").
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ComputePoolSet").
				OptionalNumberAssignment("MIN_NODES", g.ParameterOptions()).
				OptionalNumberAssignment("MAX_NODES", g.ParameterOptions()).
				OptionalBooleanAssignment("AUTO_RESUME", g.ParameterOptions()).
				OptionalNumberAssignment("AUTO_SUSPEND_SECS", g.ParameterOptions()).
				ListAssignment("BACKUP_INSTANCE_FAMILIES", "ComputePoolBackupInstanceFamilyListItem", g.ParameterOptions().Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithAdditionalValidations().
				WithValidation(g.AtLeastOneValueSet, "MinNodes", "MaxNodes", "AutoResume", "AutoSuspendSecs", "BackupInstanceFamilies", "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ComputePoolUnset").
				OptionalSQL("AUTO_RESUME").
				OptionalSQL("AUTO_SUSPEND_SECS").
				OptionalSQL("BACKUP_INSTANCE_FAMILIES").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "AutoResume", "AutoSuspendSecs", "BackupInstanceFamilies", "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Resume", "Suspend", "StopAll", "Set", "Unset", "SetTags", "UnsetTags"),
).DropOperationFromObject(ssotdefs.ComputePool).
	ShowOperationFromObject(ssotdefs.ComputePool).
	DescribeOperationFromObject(ssotdefs.ComputePool).
	WithEnums(
		ssotdefs.ComputePoolState,
		ssotdefs.ComputePoolInstanceFamily,
	)
