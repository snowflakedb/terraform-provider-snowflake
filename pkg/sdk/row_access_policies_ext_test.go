package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := rowAccessPoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")

	rowAccessPoliciesTests.Create.
		withDefaultOpts(func() *CreateRowAccessPolicyOptions {
			return &CreateRowAccessPolicyOptions{
				name: id,
				args: []CreateRowAccessPolicyArgs{{
					Name:     "n",
					DataType: dataTypeVarchar,
				}},
				body: "true",
			}
		}).
		withExpectedSqlf(
			case_RowAccessPolicies_sql_Create_basic,
			`CREATE ROW ACCESS POLICY %s AS ("n" VARCHAR(16777216)) RETURNS BOOLEAN -> true`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Create_all,
			func(opts *CreateRowAccessPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE ROW ACCESS POLICY IF NOT EXISTS %s AS ("n" VARCHAR(16777216)) RETURNS BOOLEAN -> true COMMENT = 'some comment'`,
			id.FullyQualifiedName(),
		)

	rowAccessPoliciesTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_twoArgs",
			func(opts *CreateRowAccessPolicyOptions) {
				opts.args = []CreateRowAccessPolicyArgs{{
					Name:     "n",
					DataType: dataTypeVarchar,
				}, {
					Name:     "h",
					DataType: dataTypeVarchar,
				}}
			},
			`CREATE ROW ACCESS POLICY %s AS ("n" VARCHAR(16777216), "h" VARCHAR(16777216)) RETURNS BOOLEAN -> true`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateRowAccessPolicyOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE ROW ACCESS POLICY %s AS ("n" VARCHAR(16777216)) RETURNS BOOLEAN -> true`,
			id.FullyQualifiedName(),
		)

	rowAccessPoliciesTests.Alter.
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Alter_RenameTo,
			func(opts *AlterRowAccessPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER ROW ACCESS POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Alter_SetBody,
			func(opts *AlterRowAccessPolicyOptions) { opts.SetBody = new("true") },
			"ALTER ROW ACCESS POLICY %s SET BODY -> true", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Alter_SetTags,
			func(opts *AlterRowAccessPolicyOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER ROW ACCESS POLICY %s SET TAG %s = 'value1', %s = 'value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Alter_UnsetTags,
			func(opts *AlterRowAccessPolicyOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER ROW ACCESS POLICY %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Alter_SetComment,
			func(opts *AlterRowAccessPolicyOptions) { opts.SetComment = new("comment") },
			"ALTER ROW ACCESS POLICY %s SET COMMENT = 'comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Alter_UnsetComment,
			func(opts *AlterRowAccessPolicyOptions) { opts.UnsetComment = new(true) },
			"ALTER ROW ACCESS POLICY %s UNSET COMMENT", id.FullyQualifiedName(),
		)

	rowAccessPoliciesTests.Drop.
		withExpectedSqlf(
			case_RowAccessPolicies_sql_Drop_basic,
			"DROP ROW ACCESS POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Drop_all,
			func(opts *DropRowAccessPolicyOptions) { opts.IfExists = new(true) },
			"DROP ROW ACCESS POLICY IF EXISTS %s", id.FullyQualifiedName(),
		)

	rowAccessPoliciesTests.Show.
		withExpectedSql(case_RowAccessPolicies_sql_Show_basic, "SHOW ROW ACCESS POLICIES").
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Show_all,
			func(opts *ShowRowAccessPolicyOptions) {
				opts.Like = &Like{Pattern: new("myaccount")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
				opts.Limit = &LimitFrom{Rows: new(10), From: new("foo")}
			},
			"SHOW ROW ACCESS POLICIES LIKE 'myaccount' IN ACCOUNT LIMIT 10 FROM 'foo'",
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Show_Like,
			func(opts *ShowRowAccessPolicyOptions) {
				opts.Like = &Like{Pattern: new("myaccount")}
			},
			"SHOW ROW ACCESS POLICIES LIKE 'myaccount'",
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Show_In,
			func(opts *ShowRowAccessPolicyOptions) {
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			"SHOW ROW ACCESS POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_RowAccessPolicies_sql_Show_Limit,
			func(opts *ShowRowAccessPolicyOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("foo")}
			},
			"SHOW ROW ACCESS POLICIES LIMIT 10 FROM 'foo'",
		)

	rowAccessPoliciesTests.Describe.
		withExpectedSqlf(
			case_RowAccessPolicies_sql_Describe_basic,
			"DESCRIBE ROW ACCESS POLICY %s", id.FullyQualifiedName(),
		)
}

func TestRowAccessPolicyDescription_Signature(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		want      []TableColumnSignature
	}{
		{
			name:      "signature with 1 arg",
			signature: "(A VARCHAR)",
			want: []TableColumnSignature{
				{
					Name: "A",
					Type: dataTypeVarchar,
				},
			},
		},
		{
			name:      "signature with multiple args",
			signature: "(A VARCHAR, B BOOLEAN)",
			want: []TableColumnSignature{
				{
					Name: "A",
					Type: dataTypeVarchar,
				},
				{
					Name: "B",
					Type: dataTypeBoolean,
				},
			},
		},
		{
			name:      "signature with complex name",
			signature: "(a B VARCHAR)",
			want: []TableColumnSignature{
				{
					Name: "a B",
					Type: dataTypeVarchar,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &describeRowAccessPolicyDBRow{
				Signature: tt.signature,
			}
			got, err := d.convert()
			assert.NoError(t, err)
			require.Equal(t, tt.want, got.Signature)
		})
	}
}
