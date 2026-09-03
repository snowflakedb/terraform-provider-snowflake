package sdk

func init() {
	id := budgetsTestIdSchemaObjectIdentifier
	notificationIntegrationId := randomAccountObjectIdentifier()
	procedureId := randomSchemaObjectIdentifierWithArguments(DataTypeVARCHAR)
	procedureId2 := randomSchemaObjectIdentifierWithArguments(DataTypeVARCHAR, DataTypeVARCHAR)

	budgetsTests.Create.
		withExpectedSqlf(
			case_Budgets_sql_Create_basic,
			"CREATE SNOWFLAKE.CORE.BUDGET %s ()", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Budgets_sql_Create_all,
			func(opts *CreateBudgetOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("test comment")
			},
			"CREATE SNOWFLAKE.CORE.BUDGET IF NOT EXISTS %s () COMMENT = 'test comment'",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateBudgetOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE SNOWFLAKE.CORE.BUDGET %s ()", id.FullyQualifiedName(),
		)

	budgetsTests.Drop.
		withExpectedSqlf(
			case_Budgets_sql_Drop_basic,
			"DROP SNOWFLAKE.CORE.BUDGET %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Budgets_sql_Drop_all,
			func(opts *DropBudgetOptions) { opts.IfExists = new(true) },
			"DROP SNOWFLAKE.CORE.BUDGET IF EXISTS %s", id.FullyQualifiedName(),
		)

	budgetsTests.SetSpendingLimit.
		withDefaultOpts(func() *SetSpendingLimitBudgetOptions {
			return &SetSpendingLimitBudgetOptions{
				name: id,
				args: BudgetSetSpendingLimitArgs{SpendingLimit: 1000},
			}
		}).
		withExpectedSqlf(
			case_Budgets_sql_SetSpendingLimit_basic,
			"CALL %s!SET_SPENDING_LIMIT (1000)", id.FullyQualifiedName(),
		)

	budgetsTests.GetSpendingLimit.
		withExpectedSqlf(
			case_Budgets_sql_GetSpendingLimit_basic,
			"CALL %s!GET_SPENDING_LIMIT ()", id.FullyQualifiedName(),
		)

	budgetsTests.SetEmailNotifications.
		withDefaultOpts(func() *SetEmailNotificationsBudgetOptions {
			return &SetEmailNotificationsBudgetOptions{
				name: id,
				args: BudgetSetEmailNotificationsArgs{
					Emails: []BudgetEmail{{Email: "test@example.com"}},
				},
			}
		}).
		withExpectedSqlf(
			case_Budgets_sql_SetEmailNotifications_basic,
			"CALL %s!SET_EMAIL_NOTIFICATIONS ('test@example.com')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_SetEmailNotifications_all",
			func(opts *SetEmailNotificationsBudgetOptions) {
				opts.args.NotificationIntegration = &notificationIntegrationId
			},
			`CALL %s!SET_EMAIL_NOTIFICATIONS ('\"%s\"', 'test@example.com')`,
			id.FullyQualifiedName(), notificationIntegrationId.Name(),
		)

	budgetsTests.GetNotificationIntegrations.
		withExpectedSqlf(
			case_Budgets_sql_GetNotificationIntegrations_basic,
			"CALL %s!GET_NOTIFICATION_INTEGRATIONS ()", id.FullyQualifiedName(),
		)

	budgetsTests.GetNotificationEmail.
		withExpectedSqlf(
			case_Budgets_sql_GetNotificationEmail_basic,
			"CALL %s!GET_NOTIFICATION_EMAIL ()", id.FullyQualifiedName(),
		)

	budgetsTests.GetNotificationIntegrationName.
		withExpectedSqlf(
			case_Budgets_sql_GetNotificationIntegrationName_basic,
			"CALL %s!GET_NOTIFICATION_INTEGRATION_NAME ()", id.FullyQualifiedName(),
		)

	budgetsTests.SetCycleStartAction.
		withDefaultOpts(func() *SetCycleStartActionBudgetOptions {
			return &SetCycleStartActionBudgetOptions{
				name: id,
				args: BudgetSetCycleStartActionArgs{
					Procedure: procedureId,
					Arguments: []string{"arg1"},
				},
			}
		}).
		withExpectedSqlf(
			case_Budgets_sql_SetCycleStartAction_basic,
			"CALL %s!SET_CYCLE_START_ACTION (SYSTEM$REFERENCE('PROCEDURE', '%s'), arg1)",
			id.FullyQualifiedName(), procedureId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_SetCycleStartAction_moreArgs",
			func(opts *SetCycleStartActionBudgetOptions) {
				opts.args = BudgetSetCycleStartActionArgs{Procedure: procedureId2, Arguments: []string{"arg1", "arg2"}}
			},
			"CALL %s!SET_CYCLE_START_ACTION (SYSTEM$REFERENCE('PROCEDURE', '%s'), arg1, arg2)",
			id.FullyQualifiedName(), procedureId2.FullyQualifiedName(),
		)

	budgetsTests.GetCycleStartAction.
		withExpectedSqlf(
			case_Budgets_sql_GetCycleStartAction_basic,
			"CALL %s!GET_CYCLE_START_ACTION ()", id.FullyQualifiedName(),
		)
}
