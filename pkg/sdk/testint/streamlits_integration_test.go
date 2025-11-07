//go:build non_account_level_tests

package testint

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Streamlits(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupStreamlitHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Streamlits.Drop(ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(true))
			require.NoError(t, err)
		}
	}

	createStreamlitHandle := func(t *testing.T, stage *sdk.Stage, mainFile string, opts func(*sdk.CreateStreamlitRequest)) *sdk.Streamlit {
		t.Helper()

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile)
		if opts != nil {
			opts(request)
		}
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		e, err := client.Streamlits.ShowByID(ctx, id)
		require.NoError(t, err)
		return e
	}

	assertStreamlit := func(t *testing.T, id sdk.SchemaObjectIdentifier, comment string, warehouse string) {
		t.Helper()

		e, err := client.Streamlits.ShowByID(ctx, id)
		require.NoError(t, err)

		require.NotEmpty(t, e.CreatedOn)
		require.Equal(t, id.Name(), e.Name)
		require.Equal(t, id.DatabaseName(), e.DatabaseName)
		require.Equal(t, id.SchemaName(), e.SchemaName)
		require.Empty(t, e.Title)
		require.Equal(t, "ACCOUNTADMIN", e.Owner)
		require.Equal(t, comment, e.Comment)
		require.Equal(t, warehouse, e.QueryWarehouse)
		require.NotEmpty(t, e.UrlId)
		require.Equal(t, "ROLE", e.OwnerRoleType)
	}

	t.Run("create streamlit", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile).WithComment(comment)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		assertStreamlit(t, id, comment, "")
	})

	// TODO(SNOW-1541938): remove this after fix on snowflake side
	t.Run("create streamlit with lower case warehouse fails", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix("lowercase"), nil)
		t.Cleanup(warehouseCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile).WithQueryWarehouse(warehouse.ID())
		err := client.Streamlits.Create(ctx, request)
		require.ErrorContains(t, err, fmt.Sprintf("The specified warehouse %s does not exist", strings.ToUpper(warehouse.ID().Name())))
	})

	t.Run("grant privilege to streamlits to role", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		role, roleCleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(roleCleanup)

		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile).WithComment(comment)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		assertStreamlit(t, id, comment, "")

		privileges := &sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeUsage},
		}
		on := &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: sdk.ObjectTypeStreamlit,
					Name:       id,
				},
			},
		}
		err = client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
		require.NoError(t, err)

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				Role: role.ID(),
			},
		})
		require.NoError(t, err)
		assert.Len(t, grants, 1)
		assert.Equal(t, sdk.SchemaObjectPrivilegeUsage.String(), grants[0].Privilege)
		assert.Equal(t, id.FullyQualifiedName(), grants[0].Name.FullyQualifiedName())

		on = &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeStreamlits,
					InDatabase:       sdk.Pointer(testClientHelper().Ids.DatabaseId()),
				},
			},
		}
		err = client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
		require.NoError(t, err)

		on = &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeStreamlits,
					InDatabase:       sdk.Pointer(testClientHelper().Ids.DatabaseId()),
				},
			},
		}
		err = client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
		require.NoError(t, err)
	})

	t.Run("grant privilege to streamlits to database role", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
		t.Cleanup(databaseRoleCleanup)

		databaseRoleId := testClientHelper().Ids.NewDatabaseObjectIdentifier(databaseRole.Name)

		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile).WithComment(comment)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		assertStreamlit(t, id, comment, "")

		privileges := &sdk.DatabaseRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeUsage},
		}
		on := &sdk.DatabaseRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: sdk.ObjectTypeStreamlit,
					Name:       id,
				},
			},
		}
		err = client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			To: &sdk.ShowGrantsTo{
				DatabaseRole: databaseRoleId,
			},
		})
		require.NoError(t, err)
		// Expecting two grants because database role has usage on database by default
		require.Len(t, grants, 2)

		on = &sdk.DatabaseRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeStreamlits,
					InDatabase:       sdk.Pointer(testClientHelper().Ids.DatabaseId()),
				},
			},
		}
		err = client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)

		on = &sdk.DatabaseRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeStreamlits,
					InDatabase:       sdk.Pointer(testClientHelper().Ids.DatabaseId()),
				},
			},
		}
		err = client.Grants.GrantPrivilegesToDatabaseRole(ctx, privileges, on, databaseRoleId, nil)
		require.NoError(t, err)
	})

	t.Run("alter streamlit: set", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		manifest := "manifest.yml"
		e := createStreamlitHandle(t, stage, manifest, nil)

		id := e.ID()
		comment := random.Comment()
		set := sdk.NewStreamlitSetRequest().
			WithRootLocation(stage.Location()).
			WithQueryWarehouse(warehouse.ID()).
			WithMainFile(manifest).
			WithComment(comment)
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithSet(*set))
		require.NoError(t, err)
		assertStreamlit(t, id, comment, warehouse.ID().Name())
	})

	t.Run("alter streamlit: unset", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		manifest := "manifest.yml"
		e := createStreamlitHandle(t, stage, manifest, func(r *sdk.CreateStreamlitRequest) {
			r.WithComment("foo").WithTitle("foo").WithQueryWarehouse(warehouse.ID())
		})

		id := e.ID()
		unset := sdk.NewStreamlitUnsetRequest().
			WithTitle(true).
			WithQueryWarehouse(true).
			WithComment(true)
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithUnset(*unset))
		require.NoError(t, err)
		assertStreamlit(t, id, "", "")
	})

	t.Run("alter function: rename", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml", nil)
		id := e.ID()
		t.Cleanup(cleanupStreamlitHandle(id))

		nid := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithRenameTo(nid))
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(nid))

		_, err = client.Streamlits.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		o, err := client.Streamlits.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), o.Name)
	})

	t.Run("show streamlit: with like", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml", nil)

		streamlits, err := client.Streamlits.Show(ctx, sdk.NewShowStreamlitRequest().WithLike(sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Len(t, streamlits, 1)
		require.Equal(t, *e, streamlits[0])
	})

	t.Run("show streamlit: terse with like", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml", nil)

		streamlits, err := client.Streamlits.Show(ctx, sdk.NewShowStreamlitRequest().WithTerse(true).WithLike(sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Len(t, streamlits, 1)
		sl := streamlits[0]
		require.Equal(t, e.Name, sl.Name)
		require.Equal(t, e.DatabaseName, sl.DatabaseName)
		require.Equal(t, e.SchemaName, sl.SchemaName)
		require.Equal(t, e.UrlId, sl.UrlId)
		require.Equal(t, e.CreatedOn, sl.CreatedOn)
		require.Empty(t, sl.Title)
		require.Empty(t, sl.Owner)
		require.Empty(t, sl.Comment)
		require.Empty(t, sl.QueryWarehouse)
		require.Empty(t, sl.OwnerRoleType)
	})

	t.Run("describe streamlit", func(t *testing.T) {
		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)

		mainFile := "manifest.yml"
		e := createStreamlitHandle(t, stage, mainFile, nil)
		id := e.ID()

		detail, err := client.Streamlits.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, e.Name, detail.Name)
		require.Equal(t, e.UrlId, detail.UrlId)
		require.Equal(t, mainFile, detail.MainFile)
		// TODO [SNOW-1569516]: make nicer during the identifiers rework follow up
		require.Equal(t, stage.ID().FullyQualifiedName(), sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(detail.RootLocation[1:]).FullyQualifiedName())
		require.Empty(t, detail.Title)
		require.Empty(t, detail.QueryWarehouse)
	})
}

func TestInt_StreamlitsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupStreamlitHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Streamlits.Drop(ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(true))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createStreamlitHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier, stage *sdk.Stage, mainFile string) {
		t.Helper()

		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		stage, cleanupStage := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(cleanupStage)
		manifest := "manifest.yml"

		createStreamlitHandle(t, id1, stage, manifest)
		createStreamlitHandle(t, id2, stage, manifest)

		e1, err := client.Streamlits.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Streamlits.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
