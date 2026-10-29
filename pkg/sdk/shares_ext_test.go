package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := sharesTestIdAccountObjectIdentifier

	sharesTests.Create.
		withExpectedSqlf(
			case_Shares_sql_Create_basic,
			`CREATE SHARE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Create_all,
			func(opts *CreateShareOptions) {
				opts.OrReplace = new(true)
				opts.Comment = new("comment")
			},
			`CREATE OR REPLACE SHARE %s COMMENT = 'comment'`, id.FullyQualifiedName(),
		)

	accountId1 := NewAccountIdentifier("my-org", "myaccount")
	accountId2 := NewAccountIdentifier("my-org", "myaccount2")
	tagId := randomSchemaObjectIdentifier()

	sharesTests.Alter.
		withDefaultOpts(func() *AlterShareOptions {
			return &AlterShareOptions{
				IfExists: new(true),
				name:     id,
			}
		}).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Add,
			func(opts *AlterShareOptions) {
				opts.Add = &ShareAdd{Accounts: []AccountIdentifier{accountId1}, ShareRestrictions: new(true)}
			},
			`ALTER SHARE IF EXISTS %s ADD ACCOUNTS = %s SHARE_RESTRICTIONS = true`,
			id.FullyQualifiedName(), accountId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Remove,
			func(opts *AlterShareOptions) {
				opts.Remove = &ShareRemove{Accounts: []AccountIdentifier{accountId1, accountId2}}
			},
			`ALTER SHARE IF EXISTS %s REMOVE ACCOUNTS = %s, %s`,
			id.FullyQualifiedName(), accountId1.FullyQualifiedName(), accountId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Set,
			func(opts *AlterShareOptions) {
				opts.Set = &ShareSet{Accounts: []AccountIdentifier{accountId1}, Comment: new("comment")}
			},
			`ALTER SHARE IF EXISTS %s SET ACCOUNTS = %s COMMENT = 'comment'`,
			id.FullyQualifiedName(), accountId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_Unset,
			func(opts *AlterShareOptions) {
				opts.Unset = &ShareUnset{Comment: new(true)}
			},
			`ALTER SHARE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_SetTags,
			func(opts *AlterShareOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`ALTER SHARE IF EXISTS %s SET TAG %s = 'v1'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Alter_UnsetTags,
			func(opts *AlterShareOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER SHARE IF EXISTS %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	sharesTests.Drop.
		withExpectedSqlf(
			case_Shares_sql_Drop_basic,
			`DROP SHARE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Drop_all,
			func(opts *DropShareOptions) { opts.IfExists = new(true) },
			`DROP SHARE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	sharesTests.Show.
		withExpectedSql(case_Shares_sql_Show_basic, `SHOW SHARES`).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_all,
			func(opts *ShowShareOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("from")}
			},
			`SHOW SHARES LIKE 'pattern' STARTS WITH 'prefix' LIMIT 10 FROM 'from'`,
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_Like,
			func(opts *ShowShareOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW SHARES LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_StartsWith,
			func(opts *ShowShareOptions) { opts.StartsWith = new("prefix") },
			`SHOW SHARES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_Shares_sql_Show_Limit,
			func(opts *ShowShareOptions) { opts.Limit = &LimitFrom{Rows: new(10), From: new("from")} },
			`SHOW SHARES LIMIT 10 FROM 'from'`,
		)

	externalShareId := randomExternalObjectIdentifier()

	sharesTests.Describe.
		withExpectedSqlf(
			case_Shares_sql_Describe_basic,
			`DESCRIBE SHARE %s`, sharesTestIdObjectIdentifier.FullyQualifiedName(),
		).
		// covers the DescribeConsumer path, which passes an ExternalObjectIdentifier through the same options struct
		withAdditionalSqlCasef(
			"sql_Describe_external",
			func(opts *DescribeShareOptions) { opts.name = externalShareId },
			`DESCRIBE SHARE %s`, externalShareId.FullyQualifiedName(),
		)
}

func Test_parseShareRecipients(t *testing.T) {
	testCases := []struct {
		Name   string
		Value  string
		Result []AccountIdentifier
		Total  int
	}{
		{Name: "empty string", Value: "", Result: []AccountIdentifier{}},
		{Name: "single account", Value: "ORG.ACC1", Result: []AccountIdentifier{NewAccountIdentifier("ORG", "ACC1")}},
		{
			Name:  "multiple accounts",
			Value: "ORG.ACC1, ORG.ACC2,ORG.ACC3",
			Result: []AccountIdentifier{
				NewAccountIdentifier("ORG", "ACC1"), NewAccountIdentifier("ORG", "ACC2"), NewAccountIdentifier("ORG", "ACC3"),
			},
		},
		{
			Name:  "truncated list: marker is not a recipient and gives the total",
			Value: "SOME_ORG.ACCOUNT_A,SOME_ORG.ACCOUNT_B,SOME_ORG.ACCOUNT_C,... (4 TOTAL)",
			Result: []AccountIdentifier{
				NewAccountIdentifier("SOME_ORG", "ACCOUNT_A"), NewAccountIdentifier("SOME_ORG", "ACCOUNT_B"), NewAccountIdentifier("SOME_ORG", "ACCOUNT_C"),
			},
			Total: 4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ids, total, err := parseShareRecipients(tc.Value)
			require.NoError(t, err)
			require.Equal(t, tc.Result, ids)
			require.Equal(t, tc.Total, total)
		})
	}

	t.Run("invalid entry is an error", func(t *testing.T) {
		_, _, err := parseShareRecipients("ORG.ACC1,ACC2")
		require.Error(t, err)
	})
}

func Test_shareRow_convert_truncatedRecipients(t *testing.T) {
	share, err := sharesRow{Kind: "OUTBOUND", Name: "S", To: "ORG.ACC1,ORG.ACC2,ORG.ACC3,... (4 TOTAL)"}.convert()
	require.NoError(t, err)
	require.Len(t, share.To, 3)
	require.Equal(t, 4, share.ToTotal)
}

func Test_showShareTargetAccountsOptions_sql(t *testing.T) {
	query, err := structToSQL(&showShareTargetAccountsOptions{name: "SOME_SHARE"})
	require.NoError(t, err)
	require.Equal(t, `SELECT TARGET_ACCOUNTS FROM SNOWFLAKE.ACCOUNT_USAGE.SHARES WHERE DELETED_ON IS NULL AND NAME = 'SOME_SHARE' ORDER BY MODIFIED_ON DESC LIMIT 1`, query)

	query, err = structToSQL(&showShareTargetAccountsOptions{name: `it's`})
	require.NoError(t, err)
	require.Contains(t, query, `NAME = 'it\'s'`)
}
