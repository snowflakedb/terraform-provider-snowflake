package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	id := storageIntegrationsTestIdAccountObjectIdentifier

	storageIntegrationsTests.Create.
		withDefaultOpts(func() *CreateStorageIntegrationOptions {
			return &CreateStorageIntegrationOptions{
				name: id,
				S3StorageProviderParams: &S3StorageParams{
					Protocol:          RegularS3Protocol,
					StorageAwsRoleArn: "arn:aws:iam::001234567890:role/role",
				},
				Enabled:                 true,
				StorageAllowedLocations: []StorageLocation{{Path: "allowed-loc-1"}, {Path: "allowed-loc-2"}},
			}
		}).
		withExpectedSqlf(
			case_StorageIntegrations_sql_Create_basic,
			`CREATE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::001234567890:role/role' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2')`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Create_all,
			func(opts *CreateStorageIntegrationOptions) {
				opts.IfNotExists = Bool(true)
				opts.S3StorageProviderParams = &S3StorageParams{
					Protocol:               RegularS3Protocol,
					StorageAwsRoleArn:      "arn:aws:iam::001234567890:role/role",
					StorageAwsExternalId:   String("external-id-12345"),
					StorageAwsObjectAcl:    String("bucket-owner-full-control"),
					UsePrivatelinkEndpoint: Bool(true),
				}
				opts.StorageBlockedLocations = []StorageLocation{{Path: "blocked-loc-1"}, {Path: "blocked-loc-2"}}
				opts.Comment = String("some comment")
			},
			`CREATE STORAGE INTEGRATION IF NOT EXISTS %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::001234567890:role/role' STORAGE_AWS_EXTERNAL_ID = 'external-id-12345' STORAGE_AWS_OBJECT_ACL = 'bucket-owner-full-control' USE_PRIVATELINK_ENDPOINT = true ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2') STORAGE_BLOCKED_LOCATIONS = ('blocked-loc-1', 'blocked-loc-2') COMMENT = 'some comment'`,
			id.FullyQualifiedName(),
		)

	storageIntegrationsTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_basic_s3_gov_protocol",
			func(opts *CreateStorageIntegrationOptions) {
				opts.S3StorageProviderParams.Protocol = GovS3Protocol
			},
			`CREATE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3GOV' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::001234567890:role/role' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_basic_s3_china_protocol",
			func(opts *CreateStorageIntegrationOptions) {
				opts.S3StorageProviderParams.Protocol = ChinaS3Protocol
			},
			`CREATE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'S3CHINA' STORAGE_AWS_ROLE_ARN = 'arn:aws:iam::001234567890:role/role' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_all_gcs",
			func(opts *CreateStorageIntegrationOptions) {
				opts.OrReplace = Bool(true)
				opts.S3StorageProviderParams = nil
				opts.GCSStorageProviderParams = new(GCSStorageParams)
				opts.StorageBlockedLocations = []StorageLocation{{Path: "blocked-loc-1"}, {Path: "blocked-loc-2"}}
				opts.Comment = String("some comment")
			},
			`CREATE OR REPLACE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'GCS' ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2') STORAGE_BLOCKED_LOCATIONS = ('blocked-loc-1', 'blocked-loc-2') COMMENT = 'some comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_all_azure",
			func(opts *CreateStorageIntegrationOptions) {
				opts.OrReplace = Bool(true)
				opts.S3StorageProviderParams = nil
				opts.AzureStorageProviderParams = &AzureStorageParams{
					AzureTenantId:          "azure-tenant-id",
					UsePrivatelinkEndpoint: Bool(true),
				}
				opts.StorageBlockedLocations = []StorageLocation{{Path: "blocked-loc-1"}, {Path: "blocked-loc-2"}}
				opts.Comment = String("some comment")
			},
			`CREATE OR REPLACE STORAGE INTEGRATION %s TYPE = EXTERNAL_STAGE STORAGE_PROVIDER = 'AZURE' AZURE_TENANT_ID = 'azure-tenant-id' USE_PRIVATELINK_ENDPOINT = true ENABLED = true STORAGE_ALLOWED_LOCATIONS = ('allowed-loc-1', 'allowed-loc-2') STORAGE_BLOCKED_LOCATIONS = ('blocked-loc-1', 'blocked-loc-2') COMMENT = 'some comment'`,
			id.FullyQualifiedName(),
		)

	storageIntegrationsTests.Alter.
		withModify(
			case_StorageIntegrations_validation_Alter_opts_ConflictingFields,
			func(opts *AlterStorageIntegrationOptions) {
				opts.IfExists = Bool(true)
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("one"),
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Alter_Set,
			func(opts *AlterStorageIntegrationOptions) {
				opts.Set = &StorageIntegrationSet{
					S3Params: &SetS3StorageParams{
						StorageAwsRoleArn:      String("new-aws-role-arn"),
						StorageAwsExternalId:   String("new-external-id"),
						StorageAwsObjectAcl:    String("new-aws-object-acl"),
						UsePrivatelinkEndpoint: Bool(true),
					},
					Enabled:                 Bool(false),
					StorageAllowedLocations: []StorageLocation{{Path: "new-allowed-location"}},
					StorageBlockedLocations: []StorageLocation{{Path: "new-blocked-location"}},
					Comment:                 String("changed comment"),
				}
			},
			`ALTER STORAGE INTEGRATION %s SET STORAGE_AWS_ROLE_ARN = 'new-aws-role-arn' STORAGE_AWS_EXTERNAL_ID = 'new-external-id' STORAGE_AWS_OBJECT_ACL = 'new-aws-object-acl' USE_PRIVATELINK_ENDPOINT = true ENABLED = false STORAGE_ALLOWED_LOCATIONS = ('new-allowed-location') STORAGE_BLOCKED_LOCATIONS = ('new-blocked-location') COMMENT = 'changed comment'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Alter_Unset,
			func(opts *AlterStorageIntegrationOptions) {
				opts.Unset = &StorageIntegrationUnset{
					S3Params: &UnsetS3StorageParams{
						StorageAwsExternalId:   Bool(true),
						StorageAwsObjectAcl:    Bool(true),
						UsePrivatelinkEndpoint: Bool(true),
					},
					Enabled:                 Bool(true),
					StorageBlockedLocations: Bool(true),
					Comment:                 Bool(true),
				}
			},
			`ALTER STORAGE INTEGRATION %s UNSET STORAGE_AWS_EXTERNAL_ID, STORAGE_AWS_OBJECT_ACL, USE_PRIVATELINK_ENDPOINT, ENABLED, STORAGE_BLOCKED_LOCATIONS, COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Alter_SetTags,
			func(opts *AlterStorageIntegrationOptions) {
				opts.IfExists = Bool(true)
				opts.SetTags = []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("name"),
						Value: "value",
					},
					{
						Name:  NewAccountObjectIdentifier("second-name"),
						Value: "second-value",
					},
				}
			},
			`ALTER STORAGE INTEGRATION IF EXISTS %s SET TAG "name" = 'value', "second-name" = 'second-value'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Alter_UnsetTags,
			func(opts *AlterStorageIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("name"),
					NewAccountObjectIdentifier("second-name"),
				}
			},
			`ALTER STORAGE INTEGRATION %s UNSET TAG "name", "second-name"`,
			id.FullyQualifiedName(),
		)

	storageIntegrationsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_azure",
			func(opts *AlterStorageIntegrationOptions) {
				opts.Set = &StorageIntegrationSet{
					AzureParams: &SetAzureStorageParams{
						AzureTenantId:          String("new-azure-tenant-id"),
						UsePrivatelinkEndpoint: Bool(true),
					},
					Enabled:                 Bool(false),
					StorageAllowedLocations: []StorageLocation{{Path: "new-allowed-location"}},
					StorageBlockedLocations: []StorageLocation{{Path: "new-blocked-location"}},
					Comment:                 String("changed comment"),
				}
			},
			`ALTER STORAGE INTEGRATION %s SET AZURE_TENANT_ID = 'new-azure-tenant-id' USE_PRIVATELINK_ENDPOINT = true ENABLED = false STORAGE_ALLOWED_LOCATIONS = ('new-allowed-location') STORAGE_BLOCKED_LOCATIONS = ('new-blocked-location') COMMENT = 'changed comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_gcs",
			func(opts *AlterStorageIntegrationOptions) {
				opts.Set = &StorageIntegrationSet{
					Enabled:                 Bool(false),
					StorageAllowedLocations: []StorageLocation{{Path: "new-allowed-location"}},
					StorageBlockedLocations: []StorageLocation{{Path: "new-blocked-location"}},
					Comment:                 String("changed comment"),
				}
			},
			`ALTER STORAGE INTEGRATION %s SET ENABLED = false STORAGE_ALLOWED_LOCATIONS = ('new-allowed-location') STORAGE_BLOCKED_LOCATIONS = ('new-blocked-location') COMMENT = 'changed comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_azure",
			func(opts *AlterStorageIntegrationOptions) {
				opts.Unset = &StorageIntegrationUnset{
					AzureParams: &UnsetAzureStorageParams{
						UsePrivatelinkEndpoint: Bool(true),
					},
					Enabled:                 Bool(true),
					StorageBlockedLocations: Bool(true),
					Comment:                 Bool(true),
				}
			},
			`ALTER STORAGE INTEGRATION %s UNSET USE_PRIVATELINK_ENDPOINT, ENABLED, STORAGE_BLOCKED_LOCATIONS, COMMENT`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_gcs",
			func(opts *AlterStorageIntegrationOptions) {
				opts.Unset = &StorageIntegrationUnset{
					Enabled:                 Bool(true),
					StorageBlockedLocations: Bool(true),
					Comment:                 Bool(true),
				}
			},
			`ALTER STORAGE INTEGRATION %s UNSET ENABLED, STORAGE_BLOCKED_LOCATIONS, COMMENT`,
			id.FullyQualifiedName(),
		)

	storageIntegrationsTests.Drop.
		withExpectedSqlf(
			case_StorageIntegrations_sql_Drop_basic,
			`DROP STORAGE INTEGRATION %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Drop_all,
			func(opts *DropStorageIntegrationOptions) { opts.IfExists = Bool(true) },
			`DROP STORAGE INTEGRATION IF EXISTS %s`, id.FullyQualifiedName(),
		)

	storageIntegrationsTests.Show.
		withExpectedSql(case_StorageIntegrations_sql_Show_basic, `SHOW STORAGE INTEGRATIONS`).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Show_all,
			func(opts *ShowStorageIntegrationOptions) { opts.Like = &Like{Pattern: String("some-pattern")} },
			`SHOW STORAGE INTEGRATIONS LIKE 'some-pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_StorageIntegrations_sql_Show_Like,
			func(opts *ShowStorageIntegrationOptions) { opts.Like = &Like{Pattern: String("some-pattern")} },
			`SHOW STORAGE INTEGRATIONS LIKE 'some-pattern'`,
		)

	storageIntegrationsTests.Describe.
		withExpectedSqlf(
			case_StorageIntegrations_sql_Describe_basic,
			`DESCRIBE STORAGE INTEGRATION %s`, id.FullyQualifiedName(),
		)
}

func TestToS3Protocol(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected S3Protocol
		Error    string
	}{
		{Input: "S3", Expected: RegularS3Protocol},
		{Input: "s3", Expected: RegularS3Protocol},
		{Input: "S3gov", Expected: GovS3Protocol},
		{Input: "S3GOV", Expected: GovS3Protocol},
		{Input: "S3ChInA", Expected: ChinaS3Protocol},
		{Input: "S3CHINA", Expected: ChinaS3Protocol},
		{Name: "validation: incorrect s3 protocol", Input: "incorrect", Error: "invalid S3 protocol: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "invalid S3 protocol: "},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v s3 protocol", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToS3Protocol(testCase.Input)
			if testCase.Error != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Expected, value)
			}
		})
	}
}
