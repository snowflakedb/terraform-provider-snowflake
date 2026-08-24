package sdk

import (
	"errors"
)

func init() {
	id := failoverGroupsTestIdAccountObjectIdentifier
	renameTarget := NewAccountObjectIdentifier("myfg1")
	moveTarget := NewAccountObjectIdentifier("fg2")
	allowedAccount := NewAccountIdentifier("MY_ORG", "MY_ACCOUNT")

	failoverGroupsTests.Create.
		withDefaultOpts(func() *CreateFailoverGroupOptions {
			return &CreateFailoverGroupOptions{
				IfNotExists:     Bool(true),
				name:            id,
				ObjectTypes:     []PluralObjectType{PluralObjectTypeRoles},
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
				opts.ObjectTypes = []PluralObjectType{PluralObjectTypeShares, PluralObjectTypeDatabases}
				opts.AllowedDatabases = []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")}
				opts.AllowedShares = []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}
				opts.IgnoreEditionCheck = Bool(true)
				opts.ReplicationSchedule = String("10 MINUTE")
			},
			`CREATE FAILOVER GROUP IF NOT EXISTS %s OBJECT_TYPES = SHARES, DATABASES ALLOWED_DATABASES = "db1" ALLOWED_SHARES = "share1" ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT" IGNORE EDITION CHECK REPLICATION_SCHEDULE = '10 MINUTE'`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.CreateSecondaryReplicationGroup.
		withDefaultOpts(func() *CreateSecondaryReplicationGroupFailoverGroupOptions {
			return &CreateSecondaryReplicationGroupFailoverGroupOptions{
				name:                 id,
				PrimaryFailoverGroup: NewExternalObjectIdentifierFromFullyQualifiedName("myorg.myaccount.fg1"),
			}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_CreateSecondaryReplicationGroup_basic,
			`CREATE FAILOVER GROUP %s AS REPLICA OF "myorg"."myaccount"."fg1"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_CreateSecondaryReplicationGroup_all,
			func(opts *CreateSecondaryReplicationGroupFailoverGroupOptions) {
				opts.IfNotExists = Bool(true)
			},
			`CREATE FAILOVER GROUP IF NOT EXISTS %s AS REPLICA OF "myorg"."myaccount"."fg1"`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.AlterSource.
		withAdditionalValidationCase(
			"validation_AlterSource_Set_AllowedIntegrationTypes_integrationsRequired",
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Set = &FailoverGroupSet{
					AllowedIntegrationTypes: []IntegrationType{IntegrationTypeSecurityIntegrations},
				}
			},
			errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types"),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSource_RenameTo,
			func(opts *AlterSourceFailoverGroupOptions) { opts.RenameTo = &renameTarget },
			`ALTER FAILOVER GROUP %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSource_Set,
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Set = &FailoverGroupSet{
					ObjectTypes:         []PluralObjectType{PluralObjectTypeShares},
					ReplicationSchedule: String("10 MINUTE"),
				}
			},
			`ALTER FAILOVER GROUP %s SET OBJECT_TYPES = SHARES REPLICATION_SCHEDULE = '10 MINUTE'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSource_Unset,
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Unset = &FailoverGroupUnset{ReplicationSchedule: Bool(true)}
			},
			`ALTER FAILOVER GROUP %s UNSET REPLICATION_SCHEDULE`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSource_Add,
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{
					AllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			`ALTER FAILOVER GROUP %s ADD "db1" TO ALLOWED_DATABASES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSource_Move,
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Move = &FailoverGroupMove{
					Shares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
					To:     moveTarget,
				}
			},
			`ALTER FAILOVER GROUP %s MOVE SHARES "share1" TO FAILOVER GROUP %s`,
			id.FullyQualifiedName(), moveTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterSource_Remove,
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Remove = &FailoverGroupRemove{
					AllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			`ALTER FAILOVER GROUP %s REMOVE "db1" FROM ALLOWED_DATABASES`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.AlterSource.
		withAdditionalSqlCasef(
			"sql_AlterSource_Add_AllowedAccounts",
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{AllowedAccounts: []AccountIdentifier{allowedAccount}}
			},
			`ALTER FAILOVER GROUP %s ADD "MY_ORG"."MY_ACCOUNT" TO ALLOWED_ACCOUNTS`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterSource_Add_AllowedShares",
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Add = &FailoverGroupAdd{AllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}}
			},
			`ALTER FAILOVER GROUP %s ADD "share1" TO ALLOWED_SHARES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterSource_Remove_AllowedAccounts",
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Remove = &FailoverGroupRemove{AllowedAccounts: []AccountIdentifier{allowedAccount}}
			},
			`ALTER FAILOVER GROUP %s REMOVE "MY_ORG"."MY_ACCOUNT" FROM ALLOWED_ACCOUNTS`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterSource_Remove_AllowedShares",
			func(opts *AlterSourceFailoverGroupOptions) {
				opts.Remove = &FailoverGroupRemove{AllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}}
			},
			`ALTER FAILOVER GROUP %s REMOVE "share1" FROM ALLOWED_SHARES`,
			id.FullyQualifiedName(),
		)

	failoverGroupsTests.AlterTarget.
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterTarget_Refresh,
			func(opts *AlterTargetFailoverGroupOptions) { opts.Refresh = Bool(true) },
			`ALTER FAILOVER GROUP %s REFRESH`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterTarget_Primary,
			func(opts *AlterTargetFailoverGroupOptions) { opts.Primary = Bool(true) },
			`ALTER FAILOVER GROUP %s PRIMARY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterTarget_Suspend,
			func(opts *AlterTargetFailoverGroupOptions) { opts.Suspend = Bool(true) },
			`ALTER FAILOVER GROUP %s SUSPEND`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_AlterTarget_Resume,
			func(opts *AlterTargetFailoverGroupOptions) { opts.Resume = Bool(true) },
			`ALTER FAILOVER GROUP %s RESUME`, id.FullyQualifiedName(),
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
				opts.InAccount = NewAccountIdentifierFromAccountLocator("abcd123")
			},
			`SHOW FAILOVER GROUPS IN ACCOUNT "abcd123"`,
		)

	failoverGroupsTests.ShowFailoverGroupDatabases.
		withDefaultOpts(func() *ShowFailoverGroupDatabasesFailoverGroupOptions {
			return &ShowFailoverGroupDatabasesFailoverGroupOptions{In: id}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupDatabases_basic,
			`SHOW DATABASES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupDatabases_all,
			func(opts *ShowFailoverGroupDatabasesFailoverGroupOptions) { opts.In = id },
			`SHOW DATABASES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		)

	failoverGroupsTests.ShowFailoverGroupShares.
		withDefaultOpts(func() *ShowFailoverGroupSharesFailoverGroupOptions {
			return &ShowFailoverGroupSharesFailoverGroupOptions{In: id}
		}).
		withExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupShares_basic,
			`SHOW SHARES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_FailoverGroups_sql_ShowFailoverGroupShares_all,
			func(opts *ShowFailoverGroupSharesFailoverGroupOptions) { opts.In = id },
			`SHOW SHARES IN FAILOVER GROUP %s`, id.FullyQualifiedName(),
		)
}
