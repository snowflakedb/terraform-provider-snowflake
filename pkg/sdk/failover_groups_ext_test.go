package sdk

import (
	"errors"
)

func init() {
	id := failoverGroupsTestIdAccountObjectIdentifier
	renameTarget := NewAccountObjectIdentifier("myfg1")
	moveTarget := NewAccountObjectIdentifier("fg2")
	allowedAccount := NewAccountIdentifier("MY_ORG", "MY_ACCOUNT")
	tagId := NewSchemaObjectIdentifier("db", "sc", "tag1")

	failoverGroupsTests.Create.
		withDefaultOpts(func() *CreateFailoverGroupOptions {
			return &CreateFailoverGroupOptions{
				IfNotExists:     Bool(true),
				name:            id,
				ObjectTypes:     []FailoverGroupObjectType{FailoverGroupObjectTypeRoles},
				AllowedAccounts: []AccountIdentifier{allowedAccount},
			}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_Create_basic,
			`CREATE FAILOVER GROUP IF NOT EXISTS %s OBJECT_TYPES = ROLES ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Create_all,
			func(opts *CreateFailoverGroupOptions) {
				opts.ObjectTypes = []FailoverGroupObjectType{FailoverGroupObjectTypeShares, FailoverGroupObjectTypeDatabases}
				opts.AllowedDatabases = []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")}
				opts.AllowedShares = []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}
				opts.IgnoreEditionCheck = Bool(true)
				opts.ReplicationSchedule = String("10 MINUTE")
				opts.OptimizedRefresh = Bool(true)
			},
			`CREATE FAILOVER GROUP IF NOT EXISTS %s OBJECT_TYPES = SHARES, DATABASES ALLOWED_DATABASES = "db1" ALLOWED_SHARES = "share1" ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT" IGNORE EDITION CHECK REPLICATION_SCHEDULE = '10 MINUTE' OPTIMIZED_REFRESH = true`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_externalVolumes_tag_errorIntegration",
			func(opts *CreateFailoverGroupOptions) {
				opts.AllowedExternalVolumes = []AccountObjectIdentifier{NewAccountObjectIdentifier("volume1")}
				opts.Tag = []TagAssociation{{Name: NewSchemaObjectIdentifier("db", "sc", "tag1"), Value: "v1"}}
				opts.ErrorIntegration = new(NewAccountObjectIdentifier("notification_integration"))
			},
			`CREATE FAILOVER GROUP IF NOT EXISTS %s OBJECT_TYPES = ROLES ALLOWED_EXTERNAL_VOLUMES = "volume1" ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT" TAG ("db"."sc"."tag1" = 'v1') ERROR_INTEGRATION = "notification_integration"`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.CreateSecondary.
		withDefaultOpts(func() *CreateSecondaryFailoverGroupOptions {
			return &CreateSecondaryFailoverGroupOptions{
				name:        id,
				AsReplicaOf: NewExternalObjectIdentifierFromFullyQualifiedName("myorg.myaccount.fg1"),
			}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_CreateSecondary_basic,
			`CREATE FAILOVER GROUP %s AS REPLICA OF "myorg"."myaccount"."fg1"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_CreateSecondary_all,
			func(opts *CreateSecondaryFailoverGroupOptions) {
				opts.IfNotExists = Bool(true)
			},
			`CREATE FAILOVER GROUP IF NOT EXISTS %s AS REPLICA OF "myorg"."myaccount"."fg1"`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedIntegrationTypes_integrationsRequired",
			func(opts *AlterFailoverGroupOptions) {
				opts.Set = &FailoverGroupSet{
					AllowedIntegrationTypes: []FailoverGroupIntegrationType{FailoverGroupIntegrationTypeSecurityIntegrations},
				}
			},
			errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types"),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_RenameTo,
			func(opts *AlterFailoverGroupOptions) { opts.RenameTo = &renameTarget },
			`ALTER FAILOVER GROUP %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_Set,
			func(opts *AlterFailoverGroupOptions) {
				opts.Set = &FailoverGroupSet{
					ObjectTypes:         []FailoverGroupObjectType{FailoverGroupObjectTypeShares},
					ReplicationSchedule: String("10 MINUTE"),
					OptimizedRefresh:    Bool(true),
				}
			},
			`ALTER FAILOVER GROUP %s SET OBJECT_TYPES = SHARES REPLICATION_SCHEDULE = '10 MINUTE' OPTIMIZED_REFRESH = true`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_Unset,
			func(opts *AlterFailoverGroupOptions) {
				opts.Unset = &FailoverGroupUnset{
					ReplicationSchedule: Bool(true),
				}
			},
			`ALTER FAILOVER GROUP %s UNSET REPLICATION_SCHEDULE`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_Add,
			func(opts *AlterFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{
					AllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			`ALTER FAILOVER GROUP %s ADD "db1" TO ALLOWED_DATABASES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_Move,
			func(opts *AlterFailoverGroupOptions) {
				opts.Move = &FailoverGroupMove{
					Shares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
					To:     moveTarget,
				}
			},
			`ALTER FAILOVER GROUP %s MOVE SHARES "share1" TO FAILOVER GROUP %s`,
			id.FullyQualifiedName(), moveTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_Remove,
			func(opts *AlterFailoverGroupOptions) {
				opts.Remove = &FailoverGroupRemove{
					AllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			`ALTER FAILOVER GROUP %s REMOVE "db1" FROM ALLOWED_DATABASES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_SetTags,
			func(opts *AlterFailoverGroupOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`ALTER FAILOVER GROUP %s SET TAG %s = 'v1'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Alter_UnsetTags,
			func(opts *AlterFailoverGroupOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER FAILOVER GROUP %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	failoverGroupsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Add_AllowedAccounts_ignoreEditionCheck",
			func(opts *AlterFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{AllowedAccounts: []AccountIdentifier{allowedAccount}, IgnoreEditionCheck: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s ADD "MY_ORG"."MY_ACCOUNT" TO ALLOWED_ACCOUNTS IGNORE EDITION CHECK`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_AllowedExternalVolumes",
			func(opts *AlterFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{AllowedExternalVolumes: []AccountObjectIdentifier{NewAccountObjectIdentifier("volume1")}}
			},
			`ALTER FAILOVER GROUP %s ADD "volume1" TO ALLOWED_EXTERNAL_VOLUMES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Move_ExternalVolumes",
			func(opts *AlterFailoverGroupOptions) {
				opts.Move = &FailoverGroupMove{ExternalVolumes: []AccountObjectIdentifier{NewAccountObjectIdentifier("volume1")}, To: moveTarget}
			},
			`ALTER FAILOVER GROUP %s MOVE EXTERNAL VOLUMES "volume1" TO FAILOVER GROUP %s`,
			id.FullyQualifiedName(), moveTarget.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_comment_errorIntegration_readOnly",
			func(opts *AlterFailoverGroupOptions) {
				opts.Set = &FailoverGroupSet{
					Comment:          String("fg comment"),
					ErrorIntegration: new(NewAccountObjectIdentifier("notification_integration")),
					ReadOnly:         Bool(true),
				}
			},
			`ALTER FAILOVER GROUP %s SET COMMENT = 'fg comment' ERROR_INTEGRATION = "notification_integration" READ_ONLY = true`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_multiple",
			func(opts *AlterFailoverGroupOptions) {
				opts.Unset = &FailoverGroupUnset{Comment: Bool(true), ReplicationSchedule: Bool(true), ErrorIntegration: Bool(true), ReadOnly: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s UNSET COMMENT, REPLICATION_SCHEDULE, ERROR_INTEGRATION, READ_ONLY`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_AllowedAccounts",
			func(opts *AlterFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{AllowedAccounts: []AccountIdentifier{allowedAccount}}
			},
			`ALTER FAILOVER GROUP %s ADD "MY_ORG"."MY_ACCOUNT" TO ALLOWED_ACCOUNTS`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_AllowedShares",
			func(opts *AlterFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{AllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}}
			},
			`ALTER FAILOVER GROUP %s ADD "share1" TO ALLOWED_SHARES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Remove_AllowedAccounts",
			func(opts *AlterFailoverGroupOptions) {
				opts.Remove = &FailoverGroupRemove{AllowedAccounts: []AccountIdentifier{allowedAccount}}
			},
			`ALTER FAILOVER GROUP %s REMOVE "MY_ORG"."MY_ACCOUNT" FROM ALLOWED_ACCOUNTS`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Remove_AllowedShares",
			func(opts *AlterFailoverGroupOptions) {
				opts.Remove = &FailoverGroupRemove{AllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}}
			},
			`ALTER FAILOVER GROUP %s REMOVE "share1" FROM ALLOWED_SHARES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_RpoAssurance",
			func(opts *AlterFailoverGroupOptions) {
				opts.Set = &FailoverGroupSet{RpoAssurance: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s SET RPO_ASSURANCE = true`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_OptimizedRefresh",
			func(opts *AlterFailoverGroupOptions) {
				opts.Unset = &FailoverGroupUnset{OptimizedRefresh: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s UNSET OPTIMIZED_REFRESH`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_RpoAssurance",
			func(opts *AlterFailoverGroupOptions) {
				opts.Unset = &FailoverGroupUnset{RpoAssurance: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s UNSET RPO_ASSURANCE`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.AlterSecondary.
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSecondary_Refresh,
			func(opts *AlterSecondaryFailoverGroupOptions) { opts.Refresh = Bool(true) },
			`ALTER FAILOVER GROUP %s REFRESH`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSecondary_Primary,
			func(opts *AlterSecondaryFailoverGroupOptions) { opts.Primary = Bool(true) },
			`ALTER FAILOVER GROUP %s PRIMARY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSecondary_Suspend,
			func(opts *AlterSecondaryFailoverGroupOptions) { opts.Suspend = &FailoverGroupSecondarySuspend{} },
			`ALTER FAILOVER GROUP %s SUSPEND`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSecondary_Resume,
			func(opts *AlterSecondaryFailoverGroupOptions) { opts.Resume = Bool(true) },
			`ALTER FAILOVER GROUP %s RESUME`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterSecondary_Suspend_immediate",
			func(opts *AlterSecondaryFailoverGroupOptions) {
				opts.Suspend = &FailoverGroupSecondarySuspend{Immediate: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s SUSPEND IMMEDIATE`, id.FullyQualifiedName(),
		)

	failoverGroupsTests.Drop.
		withExpectedSqlf(
			case_FailoverGroups_sql_Drop_basic,
			`DROP FAILOVER GROUP %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Drop_all,
			func(opts *DropFailoverGroupOptions) { opts.IfExists = Bool(true) },
			`DROP FAILOVER GROUP IF EXISTS %s`, id.FullyQualifiedName(),
		)

	failoverGroupsTests.Show.
		withExpectedSql(case_FailoverGroups_sql_Show_basic, `SHOW FAILOVER GROUPS`).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Show_all,
			func(opts *ShowFailoverGroupOptions) {
				opts.InAccount = new(NewAccountIdentifierFromAccountLocator("abcd123"))
			},
			`SHOW FAILOVER GROUPS IN ACCOUNT "abcd123"`,
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_Show_InAccount,
			func(opts *ShowFailoverGroupOptions) {
				opts.InAccount = new(NewAccountIdentifierFromAccountLocator("abcd123"))
			},
			`SHOW FAILOVER GROUPS IN ACCOUNT "abcd123"`,
		)

	failoverGroupsTests.ShowFailoverGroupDatabases.
		withDefaultOpts(func() *ShowFailoverGroupDatabasesFailoverGroupOptions {
			return &ShowFailoverGroupDatabasesFailoverGroupOptions{name: id}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupDatabases_basic,
			`SHOW DATABASES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupDatabases_all,
			func(opts *ShowFailoverGroupDatabasesFailoverGroupOptions) { opts.name = id },
			`SHOW DATABASES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		)

	failoverGroupsTests.ShowFailoverGroupShares.
		withDefaultOpts(func() *ShowFailoverGroupSharesFailoverGroupOptions {
			return &ShowFailoverGroupSharesFailoverGroupOptions{name: id}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupShares_basic,
			`SHOW SHARES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupShares_all,
			func(opts *ShowFailoverGroupSharesFailoverGroupOptions) { opts.name = id },
			`SHOW SHARES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		)
}
