package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	id := storageLifecyclePoliciesTestIdSchemaObjectIdentifier
	renameTarget := randomSchemaObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")

	storageLifecyclePoliciesTests.Create.
		withDefaultOpts(func() *CreateStorageLifecyclePolicyOptions {
			return &CreateStorageLifecyclePolicyOptions{
				name: id,
				args: []CreateStorageLifecyclePolicyArgs{{
					Name:     "n",
					DataType: dataTypeVarchar,
				}},
				body: "true",
			}
		}).
		withExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Create_basic,
			`CREATE STORAGE LIFECYCLE POLICY %s AS ("n" VARCHAR(16777216)) RETURNS BOOLEAN -> true`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Create_all,
			func(opts *CreateStorageLifecyclePolicyOptions) {
				opts.IfNotExists = new(true)
				opts.ArchiveTier = new(StorageLifecyclePolicyArchiveTierCold)
				opts.ArchiveForDays = new(365)
				opts.Comment = new("some comment")
				opts.Tag = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`CREATE STORAGE LIFECYCLE POLICY IF NOT EXISTS %s AS ("n" VARCHAR(16777216)) RETURNS BOOLEAN -> true ARCHIVE_TIER = COLD ARCHIVE_FOR_DAYS = 365 COMMENT = 'some comment' TAG ("tag1" = 'value1', "tag2" = 'value2')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_twoArgs",
			func(opts *CreateStorageLifecyclePolicyOptions) {
				opts.args = []CreateStorageLifecyclePolicyArgs{{
					Name:     "n1",
					DataType: dataTypeVarchar,
				}, {
					Name:     "n2",
					DataType: dataTypeVarchar,
				}}
			},
			`CREATE STORAGE LIFECYCLE POLICY %s AS ("n1" VARCHAR(16777216), "n2" VARCHAR(16777216)) RETURNS BOOLEAN -> true`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateStorageLifecyclePolicyOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE STORAGE LIFECYCLE POLICY %s AS ("n" VARCHAR(16777216)) RETURNS BOOLEAN -> true`,
			id.FullyQualifiedName(),
		)

	storageLifecyclePoliciesTests.Alter.
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Alter_RenameTo,
			func(opts *AlterStorageLifecyclePolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER STORAGE LIFECYCLE POLICY %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Alter_SetBody,
			func(opts *AlterStorageLifecyclePolicyOptions) { opts.SetBody = new("true") },
			"ALTER STORAGE LIFECYCLE POLICY %s SET BODY -> true",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Alter_Set,
			func(opts *AlterStorageLifecyclePolicyOptions) {
				opts.Set = &StorageLifecyclePolicySet{
					ArchiveTier:    new(StorageLifecyclePolicyArchiveTierCool),
					ArchiveForDays: new(120),
					Comment:        new("some comment"),
				}
			},
			"ALTER STORAGE LIFECYCLE POLICY %s SET ARCHIVE_TIER = COOL ARCHIVE_FOR_DAYS = 120 COMMENT = 'some comment'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Alter_SetTags,
			func(opts *AlterStorageLifecyclePolicyOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER STORAGE LIFECYCLE POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Alter_Unset,
			func(opts *AlterStorageLifecyclePolicyOptions) {
				opts.Unset = &StorageLifecyclePolicyUnset{
					ArchiveForDays: new(true),
					Comment:        new(true),
				}
			},
			"ALTER STORAGE LIFECYCLE POLICY %s UNSET ARCHIVE_FOR_DAYS, COMMENT",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Alter_UnsetTags,
			func(opts *AlterStorageLifecyclePolicyOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER STORAGE LIFECYCLE POLICY %s UNSET TAG "tag1", "tag2"`,
			id.FullyQualifiedName(),
		)

	storageLifecyclePoliciesTests.Drop.
		withExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Drop_basic,
			"DROP STORAGE LIFECYCLE POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Drop_all,
			func(opts *DropStorageLifecyclePolicyOptions) { opts.IfExists = new(true) },
			"DROP STORAGE LIFECYCLE POLICY IF EXISTS %s", id.FullyQualifiedName(),
		)

	storageLifecyclePoliciesTests.Show.
		withExpectedSql(case_StorageLifecyclePolicies_sql_Show_basic, "SHOW STORAGE LIFECYCLE POLICIES").
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Show_Like,
			func(opts *ShowStorageLifecyclePolicyOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
			},
			"SHOW STORAGE LIFECYCLE POLICIES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Show_In,
			func(opts *ShowStorageLifecyclePolicyOptions) {
				opts.In = &In{Account: new(true)}
			},
			"SHOW STORAGE LIFECYCLE POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Show_all,
			func(opts *ShowStorageLifecyclePolicyOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
				opts.In = &In{Account: new(true)}
			},
			"SHOW STORAGE LIFECYCLE POLICIES LIKE 'like-pattern' IN ACCOUNT",
		)

	storageLifecyclePoliciesTests.Describe.
		withExpectedSqlf(
			case_StorageLifecyclePolicies_sql_Describe_basic,
			"DESCRIBE STORAGE LIFECYCLE POLICY %s", id.FullyQualifiedName(),
		)
}

func Test_normalizeStorageLifecyclePolicyArchiveTier(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "NULL literal is normalized to empty string",
			Input:    "NULL",
			Expected: "",
		},
		{
			Name:     "lowercase null is normalized to empty string",
			Input:    "null",
			Expected: "",
		},
		{
			Name:     "empty string stays empty",
			Input:    "",
			Expected: "",
		},
		{
			Name:     "COLD tier is returned as is",
			Input:    "COLD",
			Expected: "COLD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, normalizeStorageLifecyclePolicyArchiveTier(tc.Input))
		})
	}
}
