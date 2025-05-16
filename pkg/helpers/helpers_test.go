package helpers

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSnowflakeParameterID(t *testing.T) {
	testCases := map[string]struct {
		id                 string
		fullyQualifiedName string
	}{
		"decodes quoted account object identifier": {
			id:                 `"test.name"`,
			fullyQualifiedName: `"test.name"`,
		},
		"decodes quoted database object identifier": {
			id:                 `"db"."test.name"`,
			fullyQualifiedName: `"db"."test.name"`,
		},
		"decodes quoted schema object identifier": {
			id:                 `"db"."schema"."test.name"`,
			fullyQualifiedName: `"db"."schema"."test.name"`,
		},
		"decodes quoted table column identifier": {
			id:                 `"db"."schema"."table.name"."test.name"`,
			fullyQualifiedName: `"db"."schema"."table.name"."test.name"`,
		},
		"decodes unquoted account object identifier": {
			id:                 `name`,
			fullyQualifiedName: `"name"`,
		},
		"decodes unquoted database object identifier": {
			id:                 `db.name`,
			fullyQualifiedName: `"db"."name"`,
		},
		"decodes unquoted schema object identifier": {
			id:                 `db.schema.name`,
			fullyQualifiedName: `"db"."schema"."name"`,
		},
		"decodes unquoted table column identifier": {
			id:                 `db.schema.table.name`,
			fullyQualifiedName: `"db"."schema"."table"."name"`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			id, err := DecodeSnowflakeParameterID(tc.id)
			require.NoError(t, err)
			require.Equal(t, tc.fullyQualifiedName, id.FullyQualifiedName())
		})
	}

	t.Run("identifier with too many parts", func(t *testing.T) {
		id := `this.identifier.is.too.long.to.be.decoded`
		_, err := DecodeSnowflakeParameterID(id)
		require.ErrorContains(t, err, fmt.Sprintf("unable to classify identifier: %s", id))
	})

	t.Run("incompatible empty identifier", func(t *testing.T) {
		id := ""
		_, err := DecodeSnowflakeParameterID(id)
		require.ErrorContains(t, err, fmt.Sprintf("incompatible identifier: %s", id))
	})

	t.Run("incompatible multiline identifier", func(t *testing.T) {
		id := "db.\nname"
		_, err := DecodeSnowflakeParameterID(id)
		require.ErrorContains(t, err, fmt.Sprintf("unable to read identifier: %s", id))
	})
}

// TODO: add tests for non object identifiers
func TestEncodeSnowflakeID(t *testing.T) {
	testCases := map[string]struct {
		identifier        sdk.ObjectIdentifier
		expectedEncodedID string
	}{
		"encodes account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("database"),
			expectedEncodedID: `database`,
		},
		"encodes quoted account object identifier": {
			identifier:        sdk.NewAccountObjectIdentifier("\"database\""),
			expectedEncodedID: `database`,
		},
		"encodes account object identifier with a dot": {
			identifier:        sdk.NewAccountObjectIdentifier("data.base"),
			expectedEncodedID: `data.base`,
		},
		"encodes pointer to account object identifier": {
			identifier:        sdk.Pointer(sdk.NewAccountObjectIdentifier("database")),
			expectedEncodedID: `database`,
		},
		"encodes database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("database", "schema"),
			expectedEncodedID: `database|schema`,
		},
		"encodes quoted database object identifier": {
			identifier:        sdk.NewDatabaseObjectIdentifier("\"database\"", "\"schema\""),
			expectedEncodedID: `database|schema`,
		},
		"encodes database object identifier with dots": {
			identifier:        sdk.NewDatabaseObjectIdentifier("data.base", "sche.ma"),
			expectedEncodedID: `data.base|sche.ma`,
		},
		"encodes pointer to database object identifier": {
			identifier:        sdk.Pointer(sdk.NewDatabaseObjectIdentifier("database", "schema")),
			expectedEncodedID: `database|schema`,
		},
		"encodes schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("database", "schema", "table"),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes quoted schema object identifier": {
			identifier:        sdk.NewSchemaObjectIdentifier("\"database\"", "\"schema\"", "\"table\""),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes schema object identifier with dots": {
			identifier:        sdk.NewSchemaObjectIdentifier("data.base", "sche.ma", "tab.le"),
			expectedEncodedID: `data.base|sche.ma|tab.le`,
		},
		"encodes pointer to schema object identifier": {
			identifier:        sdk.Pointer(sdk.NewSchemaObjectIdentifier("database", "schema", "table")),
			expectedEncodedID: `database|schema|table`,
		},
		"encodes table column identifier": {
			identifier:        sdk.NewTableColumnIdentifier("database", "schema", "table", "column"),
			expectedEncodedID: `database|schema|table|column`,
		},
		"encodes quoted table column identifier": {
			identifier:        sdk.NewTableColumnIdentifier("\"database\"", "\"schema\"", "\"table\"", "\"column\""),
			expectedEncodedID: `database|schema|table|column`,
		},
		"encodes table column identifier with dots": {
			identifier:        sdk.NewTableColumnIdentifier("data.base", "sche.ma", "tab.le", "col.umn"),
			expectedEncodedID: `data.base|sche.ma|tab.le|col.umn`,
		},
		"encodes pointer to table column identifier": {
			identifier:        sdk.Pointer(sdk.NewTableColumnIdentifier("database", "schema", "table", "column")),
			expectedEncodedID: `database|schema|table|column`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			encodedID := EncodeSnowflakeID(tc.identifier)
			require.Equal(t, tc.expectedEncodedID, encodedID)
		})
	}

	t.Run("panics for unsupported object identifier", func(t *testing.T) {
		id := unsupportedObjectIdentifier{}
		require.PanicsWithValue(t, fmt.Sprintf("Unsupported object identifier: %v", id), func() {
			EncodeSnowflakeID(id)
		})
	})

	nilTestCases := []any{
		(*sdk.AccountObjectIdentifier)(nil),
		(*sdk.DatabaseObjectIdentifier)(nil),
		(*sdk.SchemaObjectIdentifier)(nil),
		(*sdk.TableColumnIdentifier)(nil),
	}

	for i, tt := range nilTestCases {
		t.Run(fmt.Sprintf("handle nil pointer to object identifier %d", i), func(t *testing.T) {
			require.PanicsWithValue(t, "Nil object identifier received", func() {
				EncodeSnowflakeID(tt)
			})
		})
	}
}

type unsupportedObjectIdentifier struct{}

func (i unsupportedObjectIdentifier) Name() string {
	return "name"
}

func (i unsupportedObjectIdentifier) FullyQualifiedName() string {
	return "fully qualified name"
}

func Test_DecodeSnowflakeAccountIdentifier(t *testing.T) {
	t.Run("decodes account identifier", func(t *testing.T) {
		id, err := DecodeSnowflakeAccountIdentifier("abc.def")

		require.NoError(t, err)
		require.Equal(t, sdk.NewAccountIdentifier("abc", "def"), id)
	})

	t.Run("does not accept account locator", func(t *testing.T) {
		_, err := DecodeSnowflakeAccountIdentifier("ABC12345")

		require.ErrorContains(t, err, "identifier: ABC12345 seems to be account locator and these are not allowed - please use <organization_name>.<account_name>")
	})

	t.Run("identifier with too many parts", func(t *testing.T) {
		id := `this.identifier.is.too.long.to.be.decoded`
		_, err := DecodeSnowflakeAccountIdentifier(id)

		require.ErrorContains(t, err, fmt.Sprintf("unable to classify account identifier: %s", id))
	})

	t.Run("empty identifier", func(t *testing.T) {
		id := ""
		_, err := DecodeSnowflakeAccountIdentifier(id)

		require.ErrorContains(t, err, fmt.Sprintf("incompatible identifier: %s", id))
	})

	t.Run("multiline identifier", func(t *testing.T) {
		id := "db.\nname"
		_, err := DecodeSnowflakeAccountIdentifier(id)

		require.ErrorContains(t, err, fmt.Sprintf("unable to read identifier: %s", id))
	})
}

func TestParseRootLocation(t *testing.T) {
	tests := []struct {
		name         string
		location     string
		expectedId   string
		expectedPath string
		expectedErr  string
	}{
		{
			name:        "empty",
			location:    ``,
			expectedErr: "incompatible identifier",
		},
		{
			name:       "unquoted",
			location:   `@a.b.c`,
			expectedId: `"a"."b"."c"`,
		},
		{
			name:         "unquoted with path",
			location:     `@a.b.c/foo`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo`,
		},
		{
			name:       "partially quoted",
			location:   `@"a".b.c`,
			expectedId: `"a"."b"."c"`,
		},
		{
			name:         "partially quoted with path",
			location:     `@"a".b.c/foo`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo`,
		},
		{
			name:       "quoted",
			location:   `@"a"."b"."c"`,
			expectedId: `"a"."b"."c"`,
		},
		{
			name:         "quoted with path",
			location:     `@"a"."b"."c"/foo`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo`,
		},
		{
			name:         "unquoted with path with dots",
			location:     `@a.b.c/foo.d`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo.d`,
		},
		{
			name:         "quoted with path with dots",
			location:     `@"a"."b"."c"/foo.d`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo.d`,
		},
		{
			name:         "quoted with complex path",
			location:     `@"a"."b"."c"/foo.a/bar.b//hoge.c`,
			expectedId:   `"a"."b"."c"`,
			expectedPath: `foo.a/bar.b/hoge.c`,
		},
		{
			name:        "invalid location",
			location:    `@foo`,
			expectedErr: "expected 3 parts for location foo, got 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotPath, gotErr := ParseRootLocation(tt.location)
			if len(tt.expectedErr) > 0 {
				assert.ErrorContains(t, gotErr, tt.expectedErr)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.expectedId, gotId.FullyQualifiedName())
				assert.Equal(t, tt.expectedPath, gotPath)
			}
		})
	}
}

func Test_ContainsIdentifierIgnoreQuotes(t *testing.T) {
	testCases := []struct {
		Name          string
		Ids           []string
		Id            string
		ShouldContain bool
	}{
		{
			Name: "validation: nil Ids",
			Id:   "id",
		},
		{
			Name: "validation: empty Id",
			Ids:  []string{"id"},
			Id:   "",
		},
		{
			Name: "validation: Ids with too many parts",
			Ids:  []string{"this.id.has.too.many.parts"},
			Id:   "id",
		},
		{
			Name: "validation: Id with too many parts",
			Ids:  []string{"id"},
			Id:   "this.id.has.too.many.parts",
		},
		{
			Name: "validation: account object identifier in Ids ignore quotes with upper cased Id",
			Ids:  []string{"object", "db.schema", "db.schema.object"},
			Id:   "\"OBJECT\"",
		},
		{
			Name: "validation: account object identifier in Ids ignore quotes with upper cased id in Ids",
			Ids:  []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:   "\"object\"",
		},
		{
			Name:          "account object identifier in Ids",
			Ids:           []string{"object", "db.schema", "db.schema.object"},
			Id:            "\"object\"",
			ShouldContain: true,
		},
		{
			Name:          "database object identifier in Ids",
			Ids:           []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:            "\"db\".\"schema\"",
			ShouldContain: true,
		},
		{
			Name:          "schema object identifier in Ids",
			Ids:           []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:            "\"db\".\"schema\".\"object\"",
			ShouldContain: true,
		},
		{
			Name:          "account object identifier in Ids upper-cased",
			Ids:           []string{"OBJECT", "db.schema", "db.schema.object"},
			Id:            "\"OBJECT\"",
			ShouldContain: true,
		},
		{
			Name:          "database object identifier in Ids upper-cased",
			Ids:           []string{"object", "DB.SCHEMA", "db.schema.object"},
			Id:            "\"DB\".\"SCHEMA\"",
			ShouldContain: true,
		},
		{
			Name:          "schema object identifier in Ids upper-cased",
			Ids:           []string{"object", "db.schema", "DB.SCHEMA.OBJECT"},
			Id:            "\"DB\".\"SCHEMA\".\"OBJECT\"",
			ShouldContain: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.ShouldContain, ContainsIdentifierIgnoringQuotes(tc.Ids, tc.Id))
		})
	}
}

// External volume helper tests

// Generate input to the ParseExternalVolumeDescribedInput, useful for testing purposes
func GenerateParseExternalVolumeDescribedInput(comment string, allowWrites string, storageLocations []string, active string) []sdk.ExternalVolumeProperty {
	storageLocationProperties := make([]sdk.ExternalVolumeProperty, len(storageLocations))
	allowWritesProperty := sdk.ExternalVolumeProperty{
		Parent:  "",
		Name:    "ALLOW_WRITES",
		Type:    "Boolean",
		Value:   allowWrites,
		Default: "true",
	}

	commentProperty := sdk.ExternalVolumeProperty{
		Parent:  "",
		Name:    "COMMENT",
		Type:    "String",
		Value:   comment,
		Default: "",
	}

	activeProperty := sdk.ExternalVolumeProperty{
		Parent:  "STORAGE_LOCATIONS",
		Name:    "ACTIVE",
		Type:    "String",
		Value:   active,
		Default: "",
	}

	for i, property := range storageLocations {
		storageLocationProperties[i] = sdk.ExternalVolumeProperty{
			Parent:  "STORAGE_LOCATIONS",
			Name:    fmt.Sprintf("STORAGE_LOCATION_%s", strconv.Itoa(i+1)),
			Type:    "String",
			Value:   property,
			Default: "",
		}
	}

	return append(append([]sdk.ExternalVolumeProperty{allowWritesProperty, commentProperty}, storageLocationProperties...), activeProperty)
}

func Test_GenerateParseExternalVolumeDescribedInput(t *testing.T) {
	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_TENANT_ID":"%s","AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
		azureTenantId,
	)

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionTypeNone := "NONE"
	gcsStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":""}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeNone,
	)

	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3StorageAwsExternalId := "123456789"
	s3EncryptionTypeNone := "NONE"

	s3StorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeNone,
	)

	allowWritesTrue := "true"
	comment := "some comment"
	cases := []struct {
		TestName         string
		Comment          string
		AllowWrites      string
		StorageLocations []string
		Active           string
		ExpectedOutput   []sdk.ExternalVolumeProperty
	}{
		{
			TestName:         "Generate input",
			Comment:          comment,
			AllowWrites:      allowWritesTrue,
			StorageLocations: []string{s3StorageLocationStandard},
			Active:           "",
			ExpectedOutput: []sdk.ExternalVolumeProperty{
				{
					Parent:  "",
					Name:    "ALLOW_WRITES",
					Type:    "Boolean",
					Value:   allowWritesTrue,
					Default: "true",
				},
				{
					Parent:  "",
					Name:    "COMMENT",
					Type:    "String",
					Value:   comment,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   "",
					Default: "",
				},
			},
		},
		{
			TestName:         "Generate input - multiple locations and active set",
			Comment:          comment,
			AllowWrites:      allowWritesTrue,
			StorageLocations: []string{s3StorageLocationStandard, azureStorageLocationStandard, gcsStorageLocationStandard},
			Active:           s3StorageLocationName,
			ExpectedOutput: []sdk.ExternalVolumeProperty{
				{
					Parent:  "",
					Name:    "ALLOW_WRITES",
					Type:    "Boolean",
					Value:   allowWritesTrue,
					Default: "true",
				},
				{
					Parent:  "",
					Name:    "COMMENT",
					Type:    "String",
					Value:   comment,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_2",
					Type:    "String",
					Value:   azureStorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_3",
					Type:    "String",
					Value:   gcsStorageLocationStandard,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   s3StorageLocationName,
					Default: "",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			generatedInput := GenerateParseExternalVolumeDescribedInput(
				tc.Comment,
				tc.AllowWrites,
				tc.StorageLocations,
				tc.Active,
			)

			assert.Equal(t, len(tc.ExpectedOutput), len(generatedInput))
			for i := range generatedInput {
				assert.Equal(t, tc.ExpectedOutput[i].Parent, generatedInput[i].Parent)
				assert.Equal(t, tc.ExpectedOutput[i].Name, generatedInput[i].Name)
				assert.Equal(t, tc.ExpectedOutput[i].Type, generatedInput[i].Type)
				assert.Equal(t, tc.ExpectedOutput[i].Value, generatedInput[i].Value)
				assert.Equal(t, tc.ExpectedOutput[i].Default, generatedInput[i].Default)
			}
		})
	}
}

func Test_ParseExternalVolumeDescribed(t *testing.T) {
	azureStorageLocationName := "azureTest"
	azureStorageProvider := "AZURE"
	azureStorageBaseUrl := "azure://123456789.blob.core.windows.net/my_example_container"
	azureTenantId := "123456789"
	azureEncryptionTypeNone := "NONE"
	azureStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_TENANT_ID":"%s","AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
		azureTenantId,
	)

	azureStorageLocationWithExtraFields := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_TENANT_ID":"%s","AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":"","EXTRA_FIELD_ONE":"testing","EXTRA_FIELD_TWO":"123456"}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
		azureTenantId,
	)

	azureStorageLocationMissingTenantId := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["azure://123456789.blob.core.windows.net/my_example_container"],"AZURE_MULTI_TENANT_APP_NAME":"test12","AZURE_CONSENT_URL":"https://login.microsoftonline.com/123456789/oauth2/authorize?client_id=test&response_type=test","ENCRYPTION_TYPE":"NONE","ENCRYPTION_KMS_KEY_ID":""}`,
		azureStorageLocationName,
		azureStorageProvider,
		azureStorageBaseUrl,
	)

	gcsStorageLocationName := "gcsTest"
	gcsStorageProvider := "GCS"
	gcsStorageBaseUrl := "gcs://my_example_bucket"
	gcsEncryptionTypeNone := "NONE"
	gcsEncryptionTypeSseKms := "GCS_SSE_KMS"
	gcsEncryptionKmsKeyId := "123456789"
	gcsStorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":""}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeNone,
	)

	gcsStorageLocationWithExtraFields := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"","EXTRA_FIELD_ONE":"testing","EXTRA_FIELD_TWO":"123456"}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeNone,
	)

	gcsStorageLocationKmsEncryption := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s"}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsStorageBaseUrl,
		gcsEncryptionTypeSseKms,
		gcsEncryptionKmsKeyId,
	)

	gcsStorageLocationMissingBaseUrl := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_ALLOWED_LOCATIONS":["gcs://my_example_bucket/*"],"STORAGE_GCP_SERVICE_ACCOUNT":"test@test.iam.test.com","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":""}`,
		gcsStorageLocationName,
		gcsStorageProvider,
		gcsEncryptionTypeNone,
	)

	s3StorageLocationName := "s3Test"
	s3StorageProvider := "S3"
	s3StorageBaseUrl := "s3://my_example_bucket"
	s3StorageAwsRoleArn := "arn:aws:iam::123456789012:role/myrole"
	s3StorageAwsExternalId := "123456789"
	s3EncryptionTypeNone := "NONE"
	s3EncryptionTypeSseS3 := "AWS_SSE_S3"
	s3EncryptionTypeSseKms := "AWS_SSE_KMS"
	s3EncryptionKmsKeyId := "123456789"

	s3StorageLocationStandard := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeNone,
	)

	s3StorageLocationWithExtraFields := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s","ENCRYPTION_KMS_KEY_ID":"%s","EXTRA_FIELD_ONE":"testing","EXTRA_FIELD_TWO":"123456"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeSseKms,
		s3EncryptionKmsKeyId,
	)

	s3StorageLocationSseS3Encryption := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeSseS3,
	)

	s3StorageLocationSseKmsEncryption := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_ROLE_ARN":"%s","STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s", "ENCRYPTION_KMS_KEY_ID":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsRoleArn,
		s3StorageAwsExternalId,
		s3EncryptionTypeSseKms,
		s3EncryptionKmsKeyId,
	)

	s3StorageLocationMissingRoleArn := fmt.Sprintf(
		`{"NAME":"%s","STORAGE_PROVIDER":"%s","STORAGE_BASE_URL":"%s","STORAGE_ALLOWED_LOCATIONS":["s3://my_example_bucket/*"],"STORAGE_AWS_IAM_USER_ARN":"arn:aws:iam::123456789:user/a11b0000-s","STORAGE_AWS_EXTERNAL_ID":"%s","ENCRYPTION_TYPE":"%s"}`,
		s3StorageLocationName,
		s3StorageProvider,
		s3StorageBaseUrl,
		s3StorageAwsExternalId,
		s3EncryptionTypeNone,
	)
	allowWritesTrue := "true"
	allowWritesFalse := "false"
	comment := "some comment"
	validCases := []struct {
		Name                 string
		DescribeOutput       []sdk.ExternalVolumeProperty
		ParsedDescribeOutput ParsedExternalVolumeDescribed
	}{
		{
			Name:           "Volume with azure storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesFalse, []string{azureStorageLocationStandard}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 azureStorageLocationName,
						StorageProvider:      azureStorageProvider,
						StorageBaseUrl:       azureStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       azureEncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        azureTenantId,
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesFalse,
			},
		},
		{
			Name:           "Volume with azure storage location, with extra fields",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesFalse, []string{azureStorageLocationWithExtraFields}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 azureStorageLocationName,
						StorageProvider:      azureStorageProvider,
						StorageBaseUrl:       azureStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       azureEncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        azureTenantId,
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesFalse,
			},
		},
		{
			Name:           "Volume with gcs storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationStandard}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 gcsStorageLocationName,
						StorageProvider:      gcsStorageProvider,
						StorageBaseUrl:       gcsStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       gcsEncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with gcs storage location, with extra fields",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationWithExtraFields}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 gcsStorageLocationName,
						StorageProvider:      gcsStorageProvider,
						StorageBaseUrl:       gcsStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       gcsEncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with gcs storage location, sse kms encryption",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationKmsEncryption}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 gcsStorageLocationName,
						StorageProvider:      gcsStorageProvider,
						StorageBaseUrl:       gcsStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       gcsEncryptionTypeSseKms,
						EncryptionKmsKeyId:   gcsEncryptionKmsKeyId,
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationStandard}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 s3StorageLocationName,
						StorageProvider:      s3StorageProvider,
						StorageBaseUrl:       s3StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageAwsExternalId,
						EncryptionType:       s3EncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location, with extra fields",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationWithExtraFields}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 s3StorageLocationName,
						StorageProvider:      s3StorageProvider,
						StorageBaseUrl:       s3StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageAwsExternalId,
						EncryptionType:       s3EncryptionTypeSseKms,
						EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location, sse s3 encryption",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationSseS3Encryption}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 s3StorageLocationName,
						StorageProvider:      s3StorageProvider,
						StorageBaseUrl:       s3StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageAwsExternalId,
						EncryptionType:       s3EncryptionTypeSseS3,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name:           "Volume with s3 storage location, sse kms encryption",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationSseKmsEncryption}, ""),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 s3StorageLocationName,
						StorageProvider:      s3StorageProvider,
						StorageBaseUrl:       s3StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageAwsExternalId,
						EncryptionType:       s3EncryptionTypeSseKms,
						EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
						AzureTenantId:        "",
					},
				},
				Active:      "",
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name: "Volume with multiple storage locations and active set",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(
				comment,
				allowWritesTrue,
				[]string{s3StorageLocationStandard, gcsStorageLocationStandard, azureStorageLocationStandard},
				s3StorageLocationName,
			),
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 s3StorageLocationName,
						StorageProvider:      s3StorageProvider,
						StorageBaseUrl:       s3StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageAwsExternalId,
						EncryptionType:       s3EncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        "",
					},
					{
						Name:                 gcsStorageLocationName,
						StorageProvider:      gcsStorageProvider,
						StorageBaseUrl:       gcsStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       gcsEncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        "",
					},
					{
						Name:                 azureStorageLocationName,
						StorageProvider:      azureStorageProvider,
						StorageBaseUrl:       azureStorageBaseUrl,
						StorageAwsRoleArn:    "",
						StorageAwsExternalId: "",
						EncryptionType:       azureEncryptionTypeNone,
						EncryptionKmsKeyId:   "",
						AzureTenantId:        azureTenantId,
					},
				},
				Active:      s3StorageLocationName,
				Comment:     comment,
				AllowWrites: allowWritesTrue,
			},
		},
		{
			Name: "Volume with s3 storage location that has no comment set (in this case describe doesn't contain a comment property)",
			DescribeOutput: []sdk.ExternalVolumeProperty{
				{
					Parent:  "",
					Name:    "ALLOW_WRITES",
					Type:    "Boolean",
					Value:   allowWritesTrue,
					Default: "true",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationSseKmsEncryption,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   s3StorageLocationName,
					Default: "",
				},
			},
			ParsedDescribeOutput: ParsedExternalVolumeDescribed{
				StorageLocations: []StorageLocation{
					{
						Name:                 s3StorageLocationName,
						StorageProvider:      s3StorageProvider,
						StorageBaseUrl:       s3StorageBaseUrl,
						StorageAwsRoleArn:    s3StorageAwsRoleArn,
						StorageAwsExternalId: s3StorageAwsExternalId,
						EncryptionType:       s3EncryptionTypeSseKms,
						EncryptionKmsKeyId:   s3EncryptionKmsKeyId,
						AzureTenantId:        "",
					},
				},
				Active:      s3StorageLocationName,
				Comment:     "",
				AllowWrites: allowWritesTrue,
			},
		},
	}

	invalidCases := []struct {
		Name           string
		DescribeOutput []sdk.ExternalVolumeProperty
	}{
		{
			Name:           "Volume with s3 storage location, missing STORAGE_AWS_ROLE_ARN",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{s3StorageLocationMissingRoleArn}, ""),
		},
		{
			Name:           "Volume with azure storage location, missing AZURE_TENANT_ID",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{azureStorageLocationMissingTenantId}, ""),
		},
		{
			Name:           "Volume with gcs storage location, missing STORAGE_BASE_URL",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{gcsStorageLocationMissingBaseUrl}, ""),
		},
		{
			Name:           "Volume with no storage locations",
			DescribeOutput: GenerateParseExternalVolumeDescribedInput(comment, allowWritesTrue, []string{}, ""),
		},
		{
			Name: "Volume with no allow writes",
			DescribeOutput: []sdk.ExternalVolumeProperty{
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "STORAGE_LOCATION_1",
					Type:    "String",
					Value:   s3StorageLocationSseKmsEncryption,
					Default: "",
				},
				{
					Parent:  "STORAGE_LOCATIONS",
					Name:    "ACTIVE",
					Type:    "String",
					Value:   s3StorageLocationName,
					Default: "",
				},
			},
		},
	}

	for _, tc := range validCases {
		t.Run(tc.Name, func(t *testing.T) {
			parsed, err := ParseExternalVolumeDescribed(tc.DescribeOutput)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(tc.ParsedDescribeOutput, parsed))
		})
	}

	for _, tc := range invalidCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseExternalVolumeDescribed(tc.DescribeOutput)
			require.Error(t, err)
		})
	}
}
