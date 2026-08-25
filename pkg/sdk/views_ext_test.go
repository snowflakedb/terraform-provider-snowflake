package sdk

func init() {
	id := viewsTestIdSchemaObjectIdentifier
	sql := "SELECT id FROM t"
	renameTarget := randomSchemaObjectIdentifier()
	rowAccessPolicyId := randomSchemaObjectIdentifier()
	rowAccessPolicy2Id := randomSchemaObjectIdentifier()
	aggregationPolicyId := randomSchemaObjectIdentifier()
	tag1Id := randomSchemaObjectIdentifier()
	tag2Id := randomSchemaObjectIdentifier()
	maskingPolicyId := randomSchemaObjectIdentifier()
	projectionPolicyId := randomSchemaObjectIdentifier()
	dmfId := randomSchemaObjectIdentifier()
	alterTag1Id := NewAccountObjectIdentifier("tag1")
	alterTag2Id := NewAccountObjectIdentifier("tag2")
	applicationId := randomAccountObjectIdentifier()
	applicationPackageId := randomAccountObjectIdentifier()
	databaseId := randomAccountObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifier()

	viewsTests.Create.
		withDefaultOpts(func() *CreateViewOptions {
			return &CreateViewOptions{
				name: id,
				sql:  sql,
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_Columns_MaskingPolicy_ValidIdentifier",
			func(opts *CreateViewOptions) {
				opts.Columns = []ViewColumn{{
					Name: "foo",
					MaskingPolicy: &ViewColumnMaskingPolicy{
						MaskingPolicy: emptySchemaObjectIdentifier,
					},
				}}
			},
			errInvalidIdentifier("CreateViewOptions.Columns[0]", "MaskingPolicy"),
		).
		withAdditionalValidationCase(
			"validation_Create_Columns_ProjectionPolicy_ValidIdentifier",
			func(opts *CreateViewOptions) {
				opts.Columns = []ViewColumn{{
					Name: "foo",
					ProjectionPolicy: &ViewColumnProjectionPolicy{
						ProjectionPolicy: emptySchemaObjectIdentifier,
					},
				}}
			},
			errInvalidIdentifier("CreateViewOptions.Columns[0]", "ProjectionPolicy"),
		).
		withExpectedSqlf(
			case_Views_sql_Create_basic,
			"CREATE VIEW %s AS %s", id.FullyQualifiedName(), sql,
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Create_all,
			func(opts *CreateViewOptions) {
				opts.OrReplace = new(true)
				opts.Secure = new(true)
				opts.Temporary = new(true)
				opts.Recursive = new(true)
				opts.Columns = []ViewColumn{
					{Name: "column_without_comment"},
					{Name: "column_with_comment", Comment: new("column 2 comment")},
					{
						Name: "column",
						MaskingPolicy: &ViewColumnMaskingPolicy{
							MaskingPolicy: maskingPolicyId,
							Using:         []Column{{"a"}, {"b"}},
						},
						Tag: []TagAssociation{{Name: tag1Id, Value: "v1"}},
					},
					{
						Name: "column 2",
						ProjectionPolicy: &ViewColumnProjectionPolicy{
							ProjectionPolicy: projectionPolicyId,
						},
					},
				}
				opts.CopyGrants = new(true)
				opts.Comment = new("comment")
				opts.RowAccessPolicy = &ViewRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicyId,
					On:              []Column{{"c"}, {"d"}},
				}
				opts.AggregationPolicy = &ViewAggregationPolicy{
					AggregationPolicy: aggregationPolicyId,
					EntityKey:         []Column{{"column_with_comment"}},
				}
				opts.Tag = []TagAssociation{{Name: tag2Id, Value: "v2"}}
			},
			`CREATE OR REPLACE SECURE TEMPORARY RECURSIVE VIEW %s `+
				`("column_without_comment", "column_with_comment" COMMENT 'column 2 comment', "column" MASKING POLICY %s USING ("a", "b") TAG (%s = 'v1'), "column 2" PROJECTION POLICY %s) COPY GRANTS COMMENT = 'comment' ROW ACCESS POLICY %s ON ("c", "d") AGGREGATION POLICY %s ENTITY KEY ("column_with_comment") TAG (%s = 'v2') AS %s`,
			id.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName(), tag1Id.FullyQualifiedName(), projectionPolicyId.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), aggregationPolicyId.FullyQualifiedName(), tag2Id.FullyQualifiedName(), sql,
		)

	viewsTests.Alter.
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_RenameTo,
			func(opts *AlterViewOptions) { opts.RenameTo = &renameTarget },
			"ALTER VIEW %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetComment,
			func(opts *AlterViewOptions) { opts.SetComment = new("comment") },
			"ALTER VIEW %s SET COMMENT = 'comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetComment,
			func(opts *AlterViewOptions) { opts.UnsetComment = new(true) },
			"ALTER VIEW %s UNSET COMMENT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetSecure,
			func(opts *AlterViewOptions) { opts.SetSecure = new(true) },
			"ALTER VIEW %s SET SECURE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetChangeTracking,
			func(opts *AlterViewOptions) { opts.SetChangeTracking = new(true) },
			"ALTER VIEW %s SET CHANGE_TRACKING = true", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetSecure,
			func(opts *AlterViewOptions) { opts.UnsetSecure = new(true) },
			"ALTER VIEW %s UNSET SECURE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetTags,
			func(opts *AlterViewOptions) {
				opts.SetTags = []TagAssociation{
					{Name: alterTag1Id, Value: "value1"},
					{Name: alterTag2Id, Value: "value2"},
				}
			},
			`ALTER VIEW %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetTags,
			func(opts *AlterViewOptions) {
				opts.UnsetTags = []ObjectIdentifier{alterTag1Id, alterTag2Id}
			},
			`ALTER VIEW %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_AddDataMetricFunction,
			func(opts *AlterViewOptions) {
				opts.AddDataMetricFunction = &ViewAddDataMetricFunction{
					DataMetricFunction: []ViewDataMetricFunction{{DataMetricFunction: dmfId, On: []Column{{"foo"}}}},
				}
			},
			"ALTER VIEW %s ADD DATA METRIC FUNCTION %s ON (\"foo\")", id.FullyQualifiedName(), dmfId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_DropDataMetricFunction,
			func(opts *AlterViewOptions) {
				opts.DropDataMetricFunction = &ViewDropDataMetricFunction{
					DataMetricFunction: []ViewDataMetricFunction{{DataMetricFunction: dmfId, On: []Column{{"foo"}}}},
				}
			},
			"ALTER VIEW %s DROP DATA METRIC FUNCTION %s ON (\"foo\")", id.FullyQualifiedName(), dmfId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_ModifyDataMetricFunction,
			func(opts *AlterViewOptions) {
				opts.ModifyDataMetricFunction = &ViewModifyDataMetricFunctions{
					DataMetricFunction: []ViewModifyDataMetricFunction{{DataMetricFunction: dmfId, On: []Column{{"foo"}}, ViewDataMetricScheduleStatusOperationOption: ViewDataMetricScheduleStatusOperationOptionSuspend}},
				}
			},
			"ALTER VIEW %s MODIFY DATA METRIC FUNCTION %s ON (\"foo\") SUSPEND", id.FullyQualifiedName(), dmfId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetDataMetricSchedule,
			func(opts *AlterViewOptions) {
				opts.SetDataMetricSchedule = &ViewSetDataMetricSchedule{DataMetricSchedule: "5 MINUTE"}
			},
			"ALTER VIEW %s SET DATA_METRIC_SCHEDULE = '5 MINUTE'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetDataMetricSchedule,
			func(opts *AlterViewOptions) { opts.UnsetDataMetricSchedule = &ViewUnsetDataMetricSchedule{} },
			"ALTER VIEW %s UNSET DATA_METRIC_SCHEDULE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_AddRowAccessPolicy,
			func(opts *AlterViewOptions) {
				opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicyId,
					On:              []Column{{"a"}, {"b"}},
				}
			},
			"ALTER VIEW %s ADD ROW ACCESS POLICY %s ON (\"a\", \"b\")", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_DropRowAccessPolicy,
			func(opts *AlterViewOptions) {
				opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId}
			},
			"ALTER VIEW %s DROP ROW ACCESS POLICY %s", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_DropAndAddRowAccessPolicy,
			func(opts *AlterViewOptions) {
				opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{
					Drop: ViewDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId},
					Add: ViewAddRowAccessPolicy{
						RowAccessPolicy: rowAccessPolicy2Id,
						On:              []Column{{"a"}, {"b"}},
					},
				}
			},
			"ALTER VIEW %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (\"a\", \"b\")", id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), rowAccessPolicy2Id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_DropAllRowAccessPolicies,
			func(opts *AlterViewOptions) { opts.DropAllRowAccessPolicies = new(true) },
			"ALTER VIEW %s DROP ALL ROW ACCESS POLICIES", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetAggregationPolicy,
			func(opts *AlterViewOptions) {
				opts.SetAggregationPolicy = &ViewSetAggregationPolicy{
					AggregationPolicy: aggregationPolicyId,
					EntityKey:         []Column{{"a"}, {"b"}},
					Force:             new(true),
				}
			},
			"ALTER VIEW %s SET AGGREGATION POLICY %s ENTITY KEY (\"a\", \"b\") FORCE", id.FullyQualifiedName(), aggregationPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetAggregationPolicy,
			func(opts *AlterViewOptions) { opts.UnsetAggregationPolicy = &ViewUnsetAggregationPolicy{} },
			"ALTER VIEW %s UNSET AGGREGATION POLICY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetMaskingPolicyOnColumn,
			func(opts *AlterViewOptions) {
				opts.SetMaskingPolicyOnColumn = &ViewSetColumnMaskingPolicy{
					Name:          "column",
					MaskingPolicy: maskingPolicyId,
					Using:         []Column{{"a"}, {"b"}},
					Force:         new(true),
				}
			},
			"ALTER VIEW %s ALTER COLUMN \"column\" SET MASKING POLICY %s USING (\"a\", \"b\") FORCE", id.FullyQualifiedName(), maskingPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetMaskingPolicyOnColumn,
			func(opts *AlterViewOptions) {
				opts.UnsetMaskingPolicyOnColumn = &ViewUnsetColumnMaskingPolicy{Name: "column"}
			},
			"ALTER VIEW %s ALTER COLUMN \"column\" UNSET MASKING POLICY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetProjectionPolicyOnColumn,
			func(opts *AlterViewOptions) {
				opts.SetProjectionPolicyOnColumn = &ViewSetProjectionPolicy{
					Name:             "column",
					ProjectionPolicy: projectionPolicyId,
					Force:            new(true),
				}
			},
			"ALTER VIEW %s ALTER COLUMN \"column\" SET PROJECTION POLICY %s FORCE", id.FullyQualifiedName(), projectionPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetProjectionPolicyOnColumn,
			func(opts *AlterViewOptions) {
				opts.UnsetProjectionPolicyOnColumn = &ViewUnsetProjectionPolicy{Name: "column"}
			},
			"ALTER VIEW %s ALTER COLUMN \"column\" UNSET PROJECTION POLICY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_SetTagsOnColumn,
			func(opts *AlterViewOptions) {
				opts.SetTagsOnColumn = &ViewSetColumnTags{
					Name: "column",
					SetTags: []TagAssociation{
						{Name: alterTag1Id, Value: "value1"},
						{Name: alterTag2Id, Value: "value2"},
					},
				}
			},
			`ALTER VIEW %s ALTER COLUMN "column" SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Alter_UnsetTagsOnColumn,
			func(opts *AlterViewOptions) {
				opts.UnsetTagsOnColumn = &ViewUnsetColumnTags{
					Name:      "column",
					UnsetTags: []ObjectIdentifier{alterTag1Id, alterTag2Id},
				}
			},
			`ALTER VIEW %s ALTER COLUMN "column" UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetChangeTracking_false",
			func(opts *AlterViewOptions) { opts.SetChangeTracking = new(false) },
			"ALTER VIEW %s SET CHANGE_TRACKING = false", id.FullyQualifiedName(),
		)

	viewsTests.Drop.
		withExpectedSqlf(
			case_Views_sql_Drop_basic,
			"DROP VIEW %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Drop_all,
			func(opts *DropViewOptions) { opts.IfExists = new(true) },
			"DROP VIEW IF EXISTS %s", id.FullyQualifiedName(),
		)

	viewsTests.Show.
		withExpectedSql(
			case_Views_sql_Show_basic,
			"SHOW VIEWS",
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Show_all,
			func(opts *ShowViewOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new("myaccount")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
				opts.StartsWith = new("abc")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			"SHOW TERSE VIEWS LIKE 'myaccount' IN ACCOUNT STARTS WITH 'abc' LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Show_Like,
			func(opts *ShowViewOptions) { opts.Like = &Like{Pattern: new("myaccount")} },
			"SHOW VIEWS LIKE 'myaccount'",
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Show_In,
			func(opts *ShowViewOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			"SHOW VIEWS IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Show_StartsWith,
			func(opts *ShowViewOptions) { opts.StartsWith = new("abc") },
			"SHOW VIEWS STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_Views_sql_Show_Limit,
			func(opts *ShowViewOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW VIEWS LIMIT 10",
		).
		withAdditionalSqlCasef(
			"sql_Show_inDatabase",
			func(opts *ShowViewOptions) { opts.In = &ExtendedIn{In: In{Database: databaseId}} },
			"SHOW VIEWS IN DATABASE %s", databaseId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_inSchema",
			func(opts *ShowViewOptions) { opts.In = &ExtendedIn{In: In{Schema: schemaId}} },
			"SHOW VIEWS IN SCHEMA %s", schemaId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_inApplication",
			func(opts *ShowViewOptions) { opts.In = &ExtendedIn{Application: applicationId} },
			"SHOW VIEWS IN APPLICATION %s", applicationId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_inApplicationPackage",
			func(opts *ShowViewOptions) { opts.In = &ExtendedIn{ApplicationPackage: applicationPackageId} },
			"SHOW VIEWS IN APPLICATION PACKAGE %s", applicationPackageId.FullyQualifiedName(),
		)

	viewsTests.Describe.
		withExpectedSqlf(
			case_Views_sql_Describe_basic,
			"DESCRIBE VIEW %s", id.FullyQualifiedName(),
		)
}
