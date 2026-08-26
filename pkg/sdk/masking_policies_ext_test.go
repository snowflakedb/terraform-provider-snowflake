package sdk

func init() {
	id := maskingPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifierInSchema(id.SchemaId())
	tagId := NewAccountObjectIdentifier("123")
	tagId2 := NewAccountObjectIdentifier("456")
	signature := []CreateMaskingPolicySignature{
		{Name: "col1", DataType: dataTypeVarchar},
		{Name: "col2", DataType: dataTypeVarchar},
	}
	expression := "REPLACE('X', 1, 2)"

	maskingPoliciesTests.Create.
		withDefaultOpts(func() *CreateMaskingPolicyOptions {
			return &CreateMaskingPolicyOptions{
				name:      id,
				signature: signature,
				returns:   dataTypeVarchar,
				body:      expression,
			}
		}).
		withExpectedSqlf(
			case_MaskingPolicies_sql_Create_basic,
			`CREATE MASKING POLICY %s AS ("col1" VARCHAR(16777216), "col2" VARCHAR(16777216)) RETURNS %s -> %s`,
			id.FullyQualifiedName(), dataTypeVarchar.ToSql(), expression,
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Create_all,
			func(opts *CreateMaskingPolicyOptions) {
				opts.OrReplace = new(true)
				opts.Comment = new("some comment")
				opts.ExemptOtherPolicies = new(true)
			},
			`CREATE OR REPLACE MASKING POLICY %s AS ("col1" VARCHAR(16777216), "col2" VARCHAR(16777216)) RETURNS %s -> %s COMMENT = 'some comment' EXEMPT_OTHER_POLICIES = true`,
			id.FullyQualifiedName(), dataTypeVarchar.ToSql(), expression,
		)

	maskingPoliciesTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_RenameTo_invalidIdentifier",
			func(opts *AlterMaskingPolicyOptions) {
				opts.RenameTo = new(NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), ""))
			},
			errInvalidIdentifier("AlterMaskingPolicyOptions", "RenameTo"),
		).
		withAdditionalValidationCase(
			"validation_Alter_RenameTo_differentDatabase",
			func(opts *AlterMaskingPolicyOptions) {
				opts.RenameTo = new(NewSchemaObjectIdentifier("other_db", id.SchemaName(), "new_name"))
			},
			ErrDifferentDatabase,
		).
		withAdditionalValidationCase(
			"validation_Alter_RenameTo_differentSchema",
			func(opts *AlterMaskingPolicyOptions) {
				opts.RenameTo = new(NewSchemaObjectIdentifier(id.DatabaseName(), "other_schema", "new_name"))
			},
			ErrDifferentSchema,
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_RenameTo,
			func(opts *AlterMaskingPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER MASKING POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_SetBody,
			func(opts *AlterMaskingPolicyOptions) { opts.SetBody = new("body") },
			"ALTER MASKING POLICY %s SET BODY -> body", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_SetComment,
			func(opts *AlterMaskingPolicyOptions) { opts.SetComment = new("foo") },
			"ALTER MASKING POLICY %s SET COMMENT = 'foo'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_UnsetBody,
			func(opts *AlterMaskingPolicyOptions) { opts.UnsetBody = new(true) },
			"ALTER MASKING POLICY %s UNSET BODY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_UnsetComment,
			func(opts *AlterMaskingPolicyOptions) { opts.UnsetComment = new(true) },
			"ALTER MASKING POLICY %s UNSET COMMENT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_SetTags,
			func(opts *AlterMaskingPolicyOptions) {
				opts.IfExists = new(true)
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value-123"},
					{Name: tagId2, Value: "value-123"},
				}
			},
			`ALTER MASKING POLICY IF EXISTS %s SET TAG %s = 'value-123', %s = 'value-123'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Alter_UnsetTags,
			func(opts *AlterMaskingPolicyOptions) {
				opts.IfExists = new(true)
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER MASKING POLICY IF EXISTS %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	maskingPoliciesTests.Drop.
		withExpectedSqlf(
			case_MaskingPolicies_sql_Drop_basic,
			"DROP MASKING POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Drop_all,
			func(opts *DropMaskingPolicyOptions) { opts.IfExists = new(true) },
			"DROP MASKING POLICY IF EXISTS %s", id.FullyQualifiedName(),
		)

	maskingPoliciesTests.Show.
		withExpectedSql(case_MaskingPolicies_sql_Show_basic, "SHOW MASKING POLICIES").
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Show_all,
			func(opts *ShowMaskingPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
				opts.Limit = &LimitFrom{Rows: new(10), From: new("foo")}
			},
			"SHOW MASKING POLICIES LIKE '%s' IN ACCOUNT LIMIT 10 FROM 'foo'",
			id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Show_Like,
			func(opts *ShowMaskingPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
			},
			"SHOW MASKING POLICIES LIKE '%s'", id.Name(),
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Show_In,
			func(opts *ShowMaskingPolicyOptions) {
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			"SHOW MASKING POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_MaskingPolicies_sql_Show_Limit,
			func(opts *ShowMaskingPolicyOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("foo")}
			},
			"SHOW MASKING POLICIES LIMIT 10 FROM 'foo'",
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Account",
			func(opts *ShowMaskingPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			"SHOW MASKING POLICIES LIKE '%s' IN ACCOUNT", id.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Database",
			func(opts *ShowMaskingPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Database: id.DatabaseId()}}
			},
			"SHOW MASKING POLICIES LIKE '%s' IN DATABASE %s",
			id.Name(), id.DatabaseId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_Like_In_Schema",
			func(opts *ShowMaskingPolicyOptions) {
				opts.Like = &Like{Pattern: new(id.Name())}
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
			},
			"SHOW MASKING POLICIES LIKE '%s' IN SCHEMA %s",
			id.Name(), id.SchemaId().FullyQualifiedName(),
		)

	maskingPoliciesTests.Describe.
		withExpectedSqlf(
			case_MaskingPolicies_sql_Describe_basic,
			"DESCRIBE MASKING POLICY %s", id.FullyQualifiedName(),
		)
}
