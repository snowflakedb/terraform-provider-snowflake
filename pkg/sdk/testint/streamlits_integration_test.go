package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Streamlits(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)

	cleanupStreamlitHandle := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Streamlits.Drop(ctx, sdk.NewDropStreamlitRequest(id).WithIfExists(sdk.Bool(true)))
			require.NoError(t, err)
		}
	}

	createStreamlitHandle := func(t *testing.T, stage *sdk.Stage, mainFile string) *sdk.Streamlit {
		t.Helper()

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile)
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
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)

		comment := random.StringN(4)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile).WithComment(&comment)
		err := client.Streamlits.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupStreamlitHandle(id))

		assertStreamlit(t, id, comment, "")
	})

	// TODO [SNOW-1272222]: fix the test when it starts working on Snowflake side
	t.Run("grant privilege to streamlits", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)

		role, roleCleanup := createRole(t, client)
		t.Cleanup(roleCleanup)

		comment := random.StringN(4)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(4))
		mainFile := "manifest.yml"
		request := sdk.NewCreateStreamlitRequest(id, stage.Location(), mainFile).WithComment(&comment)
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
		assert.Equal(t, 1, len(grants))
		assert.Equal(t, sdk.SchemaObjectPrivilegeUsage.String(), grants[0].Privilege)
		assert.Equal(t, id.FullyQualifiedName(), grants[0].Name.FullyQualifiedName())

		on = &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				Future: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeStreamlits,
					InDatabase:       sdk.Pointer(sdk.NewAccountObjectIdentifier(TestDatabaseName)),
				},
			},
		}
		err = client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
		require.Error(t, err)
		require.ErrorContains(t, err, "Unsupported feature 'STREAMLIT'")

		on = &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypeStreamlits,
					InDatabase:       sdk.Pointer(sdk.NewAccountObjectIdentifier(TestDatabaseName)),
				},
			},
		}
		err = client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
		require.NoError(t, err)
	})

	t.Run("alter streamlit: set", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)
		manifest := "manifest.yml"
		e := createStreamlitHandle(t, stage, manifest)

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		comment := random.StringN(4)
		set := sdk.NewStreamlitSetRequest(sdk.String(stage.Location()), &manifest).WithComment(&comment)
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithSet(set))
		require.NoError(t, err)
		assertStreamlit(t, id, comment, "")
	})

	t.Run("alter function: rename", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml")

		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)
		nid := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, random.StringN(3))
		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithRenameTo(&nid))
		if err != nil {
			t.Cleanup(cleanupStreamlitHandle(id))
		} else {
			t.Cleanup(cleanupStreamlitHandle(nid))
		}
		require.NoError(t, err)

		_, err = client.Streamlits.ShowByID(ctx, id)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		o, err := client.Streamlits.ShowByID(ctx, nid)
		require.NoError(t, err)
		require.Equal(t, nid.Name(), o.Name)
	})

	t.Run("show streamlit: with like", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml")

		streamlits, err := client.Streamlits.Show(ctx, sdk.NewShowStreamlitRequest().WithLike(&sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(streamlits))
		require.Equal(t, *e, streamlits[0])
	})

	t.Run("show streamlit: terse with like", func(t *testing.T) {
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)
		e := createStreamlitHandle(t, stage, "manifest.yml")

		streamlits, err := client.Streamlits.Show(ctx, sdk.NewShowStreamlitRequest().WithTerse(sdk.Bool(true)).WithLike(&sdk.Like{Pattern: &e.Name}))
		require.NoError(t, err)
		require.Equal(t, 1, len(streamlits))
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
		stage, cleanupStage := createStage(t, client, sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(4)))
		t.Cleanup(cleanupStage)

		mainFile := "manifest.yml"
		e := createStreamlitHandle(t, stage, mainFile)
		id := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, e.Name)

		detail, err := client.Streamlits.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, e.Name, detail.Name)
		require.Equal(t, e.UrlId, detail.UrlId)
		require.Equal(t, mainFile, detail.MainFile)
		require.Equal(t, stage.Location(), detail.RootLocation)
		require.Empty(t, detail.Title)
		require.Empty(t, detail.QueryWarehouse)
	})
}
