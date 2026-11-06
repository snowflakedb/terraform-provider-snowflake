package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FormatAndParseAlertSchedule(t *testing.T) {
	t.Run("interval", func(t *testing.T) {
		formatted := FormatAlertIntervalSchedule(5)
		require.Equal(t, "5 MINUTE", formatted)

		parsed, err := ParseAlertSchedule(formatted)
		require.NoError(t, err)
		require.Equal(t, AlertSchedule{Minutes: 5}, parsed)
	})

	t.Run("cron", func(t *testing.T) {
		formatted := FormatAlertCronSchedule("0 0 * * *", "America/Los_Angeles")
		require.Equal(t, "USING CRON 0 0 * * * America/Los_Angeles", formatted)

		parsed, err := ParseAlertSchedule(formatted)
		require.NoError(t, err)
		require.Equal(t, AlertSchedule{CronExpression: "0 0 * * *", TimeZone: "America/Los_Angeles"}, parsed)
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := ParseAlertSchedule("USING CRON")
		require.Error(t, err)
		_, err = ParseAlertSchedule("daily")
		require.Error(t, err)
		_, err = ParseAlertSchedule("five MINUTE")
		require.Error(t, err)
	})
}

func init() {
	id := alertsTestIdSchemaObjectIdentifier
	applicationId := NewAccountObjectIdentifier("application_name")
	applicationPackageId := NewAccountObjectIdentifier("application_package_name")
	warehouse := randomAccountObjectIdentifier()
	schedule := "1 minute"
	existsCondition := "SELECT 1"
	action := "INSERT INTO FOO VALUES (1)"
	comment := "comment"

	alertsTests.Create.
		withDefaultOpts(func() *CreateAlertOptions {
			return &CreateAlertOptions{
				name:      id,
				condition: AlertCondition{Condition: []string{existsCondition}},
				action:    action,
			}
		}).
		withExpectedSqlf(
			case_Alerts_sql_Create_basic,
			`CREATE ALERT %s IF (EXISTS (%s)) THEN %s`,
			id.FullyQualifiedName(), existsCondition, action,
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Create_all,
			func(opts *CreateAlertOptions) {
				opts.IfNotExists = new(true)
				opts.Warehouse = &warehouse
				opts.Schedule = &schedule
				opts.Comment = &comment
			},
			`CREATE ALERT IF NOT EXISTS %s WAREHOUSE = %s SCHEDULE = '%s' COMMENT = '%s' IF (EXISTS (%s)) THEN %s`,
			id.FullyQualifiedName(), warehouse.FullyQualifiedName(), schedule, comment, existsCondition, action,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateAlertOptions) {
				opts.OrReplace = new(true)
				opts.Warehouse = &warehouse
				opts.Schedule = &schedule
			},
			`CREATE OR REPLACE ALERT %s WAREHOUSE = %s SCHEDULE = '%s' IF (EXISTS (%s)) THEN %s`,
			id.FullyQualifiedName(), warehouse.FullyQualifiedName(), schedule, existsCondition, action,
		)

	alertsTests.Alter.
		withModify(
			case_Alerts_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterAlertOptions) {
				opts.Resume = new(true)
				opts.Set = &AlertSet{Comment: &comment}
			},
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Alter_Resume,
			func(opts *AlterAlertOptions) { opts.Resume = new(true) },
			"ALTER ALERT %s RESUME", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Alter_Suspend,
			func(opts *AlterAlertOptions) { opts.Suspend = new(true) },
			"ALTER ALERT %s SUSPEND", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Alter_Set,
			func(opts *AlterAlertOptions) {
				opts.Set = &AlertSet{
					Warehouse: &warehouse,
					Schedule:  &schedule,
					Comment:   &comment,
				}
			},
			"ALTER ALERT %s SET WAREHOUSE = %s SCHEDULE = '%s' COMMENT = '%s'",
			id.FullyQualifiedName(), warehouse.FullyQualifiedName(), schedule, comment,
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetComment",
			func(opts *AlterAlertOptions) { opts.Set = &AlertSet{Comment: &comment} },
			"ALTER ALERT %s SET COMMENT = '%s'", id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Alter_Unset,
			func(opts *AlterAlertOptions) {
				opts.Unset = &AlertUnset{
					Warehouse: new(true),
					Schedule:  new(true),
					Comment:   new(true),
				}
			},
			"ALTER ALERT %s UNSET WAREHOUSE, SCHEDULE, COMMENT", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetComment",
			func(opts *AlterAlertOptions) { opts.Unset = &AlertUnset{Comment: new(true)} },
			"ALTER ALERT %s UNSET COMMENT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Alter_ModifyCondition,
			func(opts *AlterAlertOptions) { opts.ModifyCondition = []string{"SELECT * FROM FOO"} },
			"ALTER ALERT %s MODIFY CONDITION EXISTS (SELECT * FROM FOO)", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Alter_ModifyAction,
			func(opts *AlterAlertOptions) { opts.ModifyAction = &action },
			"ALTER ALERT %s MODIFY ACTION %s", id.FullyQualifiedName(), action,
		)

	alertsTests.Drop.
		withExpectedSqlf(
			case_Alerts_sql_Drop_basic,
			"DROP ALERT %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Drop_all,
			func(opts *DropAlertOptions) { opts.IfExists = new(true) },
			"DROP ALERT IF EXISTS %s", id.FullyQualifiedName(),
		)

	alertsTests.Show.
		withExpectedSql(case_Alerts_sql_Show_basic, "SHOW ALERTS").
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Show_all,
			func(opts *ShowAlertOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{Account: new(true)}
				opts.StartsWith = new("FOO")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			"SHOW TERSE ALERTS LIKE '%s' IN ACCOUNT STARTS WITH 'FOO' LIMIT 10", id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Show_Like,
			func(opts *ShowAlertOptions) { opts.Like = &Like{Pattern: new(id.Name())} },
			"SHOW ALERTS LIKE '%s'", id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Show_In,
			func(opts *ShowAlertOptions) { opts.In = &ExtendedIn{Account: new(true)} },
			"SHOW ALERTS IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Show_StartsWith,
			func(opts *ShowAlertOptions) { opts.StartsWith = new("FOO") },
			"SHOW ALERTS STARTS WITH 'FOO'",
		).
		withModifyAndExpectedSqlf(
			case_Alerts_sql_Show_Limit,
			func(opts *ShowAlertOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW ALERTS LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_terse",
			func(opts *ShowAlertOptions) { opts.Terse = new(true) },
			"SHOW TERSE ALERTS",
		).
		withAdditionalSqlCasef(
			"sql_Show_likeAndInAccount",
			func(opts *ShowAlertOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{Account: new(true)}
			},
			"SHOW ALERTS LIKE '%s' IN ACCOUNT", id.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Show_likeAndInDatabase",
			func(opts *ShowAlertOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{Database: id.DatabaseId()}
			},
			"SHOW ALERTS LIKE '%s' IN DATABASE %s", id.Name(), id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_likeAndInSchema",
			func(opts *ShowAlertOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{Schema: id.SchemaId()}
			},
			"SHOW ALERTS LIKE '%s' IN SCHEMA %s", id.Name(), id.SchemaId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_inApplication",
			func(opts *ShowAlertOptions) { opts.In = &ExtendedIn{Application: applicationId} },
			"SHOW ALERTS IN APPLICATION %s", applicationId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_inApplicationPackage",
			func(opts *ShowAlertOptions) { opts.In = &ExtendedIn{ApplicationPackage: applicationPackageId} },
			"SHOW ALERTS IN APPLICATION PACKAGE %s", applicationPackageId.FullyQualifiedName(),
		)

	alertsTests.Describe.
		withExpectedSqlf(
			case_Alerts_sql_Describe_basic,
			"DESCRIBE ALERT %s", id.FullyQualifiedName(),
		)
}
