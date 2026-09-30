//go:build account_level_tests

package testint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The tests below pin down what SHOW GRANTS reports for objects whose specific object type differs from the
// generic one users put in their configuration, with the 2026_06 bundle (BCR-2371) disabled and enabled.
// Background: https://github.com/snowflakedb/terraform-provider-snowflake/issues/5200 and
// https://docs.snowflake.com/en/release-notes/bcr-bundles/2026_06/bcr-2371.
//
// The provider consequence - a permadiff that surfaces on refresh - is covered by the acceptance tests in
// pkg/testacc/resource_grant_privileges_to_database_role_bcr_202606_account_level_acceptance_test.go.
//
// Toggling a behavior change bundle is account-wide, so these run against the secondary account.

type grantedOnTestTarget struct {
	specificType                sdk.ObjectType
	genericType                 sdk.ObjectType
	grantedOnWhenBundleDisabled sdk.ObjectType
	id                          sdk.SchemaObjectIdentifier
}

type grantObservation struct {
	err         error
	found       bool
	grantedOn   sdk.ObjectType
	grantOn     sdk.ObjectType
	grantedBy   string
	granteeName string
	granteeFqn  string
}

func (o grantObservation) String() string {
	if o.err != nil {
		return fmt.Sprintf("error: %v", o.err)
	}
	if !o.found {
		return "grant not present in output"
	}
	return fmt.Sprintf("granted_on=%q grant_on=%q granted_by=%q grantee_name=%q (fqn %q)",
		o.grantedOn, o.grantOn, o.grantedBy, o.granteeName, o.granteeFqn)
}

type grantObservations struct {
	onGeneric      grantObservation
	onSpecific     grantObservation
	toDatabaseRole grantObservation
}

// observeGrant grants SELECT on the target using grantAs, records how each SHOW GRANTS form reports it, and
// revokes it again using revokeAs. The returned error is the revoke error, if any.
func observeGrant(
	t *testing.T,
	ctx context.Context,
	client *sdk.Client,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	target grantedOnTestTarget,
	grantAs sdk.ObjectType,
	revokeAs sdk.ObjectType,
) (grantObservations, error) {
	t.Helper()

	privileges := sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect})
	on := func(objectType sdk.ObjectType) *sdk.DatabaseRoleGrantOnRequest {
		return sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithSchemaObject(sdk.Object{ObjectType: objectType, Name: target.id}))
	}

	err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, sdk.NewGrantPrivilegesToDatabaseRoleRequest(databaseRoleId).WithPrivileges(*privileges).WithOn(*on(grantAs)))
	require.NoErrorf(t, err, "GRANT SELECT ON %s %s failed", grantAs, target.id.FullyQualifiedName())

	observations := grantObservations{
		onGeneric:      observeShowGrantsOnObject(t, ctx, client, target.genericType, target.id, databaseRoleId),
		onSpecific:     observeShowGrantsOnObject(t, ctx, client, target.specificType, target.id, databaseRoleId),
		toDatabaseRole: observeShowGrantsToDatabaseRole(t, target, databaseRoleId),
	}
	t.Logf("granted as %q:", grantAs)
	t.Logf("  SHOW GRANTS ON %s   -> %s", target.genericType, observations.onGeneric)
	t.Logf("  SHOW GRANTS ON %s   -> %s", target.specificType, observations.onSpecific)
	t.Logf("  SHOW GRANTS TO DATABASE ROLE -> %s", observations.toDatabaseRole)

	revokeErr := client.Grants.RevokePrivilegesFromDatabaseRole(ctx, sdk.NewRevokePrivilegesFromDatabaseRoleRequest(databaseRoleId).WithPrivileges(*privileges).WithOn(*on(revokeAs)))
	if revokeErr != nil {
		t.Logf("  REVOKE SELECT ON %s %s failed: %v", revokeAs, target.id.FullyQualifiedName(), revokeErr)
		// Fall back to the specific type so the next case starts from a clean slate.
		require.NoError(t, client.Grants.RevokePrivilegesFromDatabaseRole(ctx, sdk.NewRevokePrivilegesFromDatabaseRoleRequest(databaseRoleId).WithPrivileges(*privileges).WithOn(*on(target.specificType))))
	}
	return observations, revokeErr
}

func observeShowGrantsOnObject(
	t *testing.T,
	ctx context.Context,
	client *sdk.Client,
	objectType sdk.ObjectType,
	id sdk.SchemaObjectIdentifier,
	databaseRoleId sdk.DatabaseObjectIdentifier,
) grantObservation {
	t.Helper()

	grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{ObjectType: objectType, Name: id})))
	if err != nil {
		return grantObservation{err: err}
	}
	for _, grant := range grants {
		if grant.Privilege != string(sdk.SchemaObjectPrivilegeSelect) {
			continue
		}
		if grant.GrantTo != sdk.ObjectTypeDatabaseRole && grant.GrantedTo != sdk.ObjectTypeDatabaseRole {
			continue
		}
		// Matched loosely on purpose: whether the grantee comes back prefixed with its database is one of
		// the things being measured.
		if !strings.Contains(grant.GranteeName.FullyQualifiedName(), databaseRoleId.Name()) {
			continue
		}
		return grantObservation{
			found:       true,
			grantedOn:   grant.GrantedOn,
			grantOn:     grant.GrantOn,
			grantedBy:   grant.GrantedBy.Name(),
			granteeName: grant.GranteeName.Name(),
			granteeFqn:  grant.GranteeName.FullyQualifiedName(),
		}
	}
	return grantObservation{}
}

func observeShowGrantsToDatabaseRole(t *testing.T, target grantedOnTestTarget, databaseRoleId sdk.DatabaseObjectIdentifier) grantObservation {
	t.Helper()

	grants, err := secondaryTestClientHelper().Grant.ShowGrantsToDatabaseRole(t, databaseRoleId)
	if err != nil {
		return grantObservation{err: err}
	}
	for _, grant := range grants {
		if grant.Name.FullyQualifiedName() == target.id.FullyQualifiedName() && grant.Privilege == string(sdk.SchemaObjectPrivilegeSelect) {
			return grantObservation{
				found:       true,
				grantedOn:   grant.GrantedOn,
				grantOn:     grant.GrantOn,
				grantedBy:   grant.GrantedBy.Name(),
				granteeName: grant.GranteeName.Name(),
				granteeFqn:  grant.GranteeName.FullyQualifiedName(),
			}
		}
	}
	return grantObservation{}
}

func assertGrantedOnBundleDisabled(
	t *testing.T,
	ctx context.Context,
	client *sdk.Client,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	target grantedOnTestTarget,
) {
	t.Helper()

	t.Run("granted with the generic object type", func(t *testing.T) {
		observations, revokeErr := observeGrant(t, ctx, client, databaseRoleId, target, target.genericType, target.genericType)
		assert.NoError(t, revokeErr)

		assert.NoError(t, observations.onGeneric.err)
		assert.True(t, observations.onGeneric.found)
		assert.Equal(t, target.grantedOnWhenBundleDisabled, observations.onGeneric.grantedOn)
		assert.NotEmpty(t, observations.onGeneric.grantedBy)
		assert.Equal(t, target.specificType, observations.toDatabaseRole.grantedOn)
	})

	t.Run("granted with the specific object type", func(t *testing.T) {
		observations, revokeErr := observeGrant(t, ctx, client, databaseRoleId, target, target.specificType, target.specificType)
		assert.NoError(t, revokeErr)

		assert.NoError(t, observations.onSpecific.err)
		assert.True(t, observations.onSpecific.found)
		assert.Equal(t, target.grantedOnWhenBundleDisabled, observations.onSpecific.grantedOn)
		assert.Equal(t, target.specificType, observations.toDatabaseRole.grantedOn)
	})
}

func assertGrantedOnBundleEnabled(
	t *testing.T,
	ctx context.Context,
	client *sdk.Client,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	target grantedOnTestTarget,
) {
	t.Helper()

	t.Run("granted with the generic object type", func(t *testing.T) {
		observations, revokeErr := observeGrant(t, ctx, client, databaseRoleId, target, target.genericType, target.genericType)
		assert.NoError(t, revokeErr)

		assert.NoError(t, observations.onGeneric.err)
		assert.True(t, observations.onGeneric.found)
		assert.Equal(t, target.specificType, observations.onGeneric.grantedOn)
		assert.NotEmpty(t, observations.onGeneric.grantedBy)
		assert.Equal(t, target.specificType, observations.toDatabaseRole.grantedOn)
	})

	t.Run("granted with the specific object type", func(t *testing.T) {
		observations, revokeErr := observeGrant(t, ctx, client, databaseRoleId, target, target.specificType, target.specificType)
		assert.NoError(t, revokeErr)

		assert.NoError(t, observations.onSpecific.err)
		assert.True(t, observations.onSpecific.found)
		assert.Equal(t, target.specificType, observations.onSpecific.grantedOn)
		assert.Equal(t, target.specificType, observations.toDatabaseRole.grantedOn)
	})
}

func TestInt_ShowGrants_bcr2026_06_grantedOnSpecificObjectType(t *testing.T) {
	client := testSecondaryClient(t)
	ctx := context.Background()

	databaseRole, databaseRoleCleanup := secondaryTestClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	table, tableCleanup := secondaryTestClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	view, viewCleanup := secondaryTestClientHelper().View.CreateView(t, fmt.Sprintf(`select "ID" from %s`, table.ID().FullyQualifiedName()))
	t.Cleanup(viewCleanup)

	dynamicTable, dynamicTableCleanup := secondaryTestClientHelper().DynamicTable.CreateDynamicTable(t, table.ID())
	t.Cleanup(dynamicTableCleanup)

	targets := []grantedOnTestTarget{
		{specificType: sdk.ObjectTypeView, genericType: sdk.ObjectTypeTable, grantedOnWhenBundleDisabled: sdk.ObjectTypeView, id: view.ID()},
		{specificType: sdk.ObjectTypeDynamicTable, genericType: sdk.ObjectTypeTable, grantedOnWhenBundleDisabled: sdk.ObjectTypeDynamicTable, id: dynamicTable.ID()},
	}

	t.Run("bundle disabled", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.DisableBcrBundle(t, "2026_06")

		for _, target := range targets {
			t.Run(target.specificType.String(), func(t *testing.T) {
				assertGrantedOnBundleDisabled(t, ctx, client, databaseRole.ID(), target)
			})
		}
	})

	t.Run("bundle enabled", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2026_06")

		for _, target := range targets {
			t.Run(target.specificType.String(), func(t *testing.T) {
				assertGrantedOnBundleEnabled(t, ctx, client, databaseRole.ID(), target)
			})
		}
	})
}

// TestInt_ShowGrants_bcr2026_06_grantedOnIcebergTable is the Iceberg-table variant, matching the object type from the reported issue.
func TestInt_ShowGrants_bcr2026_06_grantedOnIcebergTable(t *testing.T) {
	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	client := testSecondaryClient(t)
	ctx := context.Background()

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	externalVolumeId, externalVolumeCleanup := secondaryTestClientHelper().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	// The external volume is set at the database level so the Iceberg table does not have to repeat it.
	db, dbCleanup := secondaryTestClientHelper().Database.CreateDatabaseWithRequest(t,
		secondaryTestClientHelper().Database.TestParametersSet(secondaryTestClientHelper().Ids.RandomAccountObjectIdentifier()).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaId := sdk.NewDatabaseObjectIdentifier(db.ID().Name(), "PUBLIC")

	databaseRole, databaseRoleCleanup := secondaryTestClientHelper().DatabaseRole.CreateDatabaseRoleInDatabase(t, db.ID())
	t.Cleanup(databaseRoleCleanup)

	icebergTableId := secondaryTestClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	_, icebergTableCleanup := secondaryTestClientHelper().IcebergTable.CreateWithRequest(t,
		sdk.NewCreateIcebergTableRequest(icebergTableId, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{{Name: "ID", ColumnType: testdatatypes.DataTypeNumber}},
		}).
			WithCatalog(sdk.IcebergTableCatalogSnowflake).
			WithExternalVolume(externalVolumeId).
			WithBaseLocation("bcr_2371_grants_test"))
	t.Cleanup(icebergTableCleanup)

	target := grantedOnTestTarget{
		specificType:                sdk.ObjectTypeIcebergTable,
		genericType:                 sdk.ObjectTypeTable,
		grantedOnWhenBundleDisabled: sdk.ObjectTypeTable,
		id:                          icebergTableId,
	}

	t.Run("bundle disabled", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.DisableBcrBundle(t, "2026_06")

		assertGrantedOnBundleDisabled(t, ctx, client, databaseRole.ID(), target)
	})

	t.Run("bundle enabled", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2026_06")

		assertGrantedOnBundleEnabled(t, ctx, client, databaseRole.ID(), target)
	})
}
