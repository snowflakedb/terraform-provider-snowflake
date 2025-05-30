//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Services(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	specTemplate := `
spec:
  containers:
  - name: {{ container_name }}
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	specTemplateUsing := []sdk.ListItem{
		{Key: "container_name", Value: `'example'`},
	}
	// TODO(SNOW-2129575): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
	t.Cleanup(stageCleanup)

	revertParameter := testClientHelper().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterPythonProfilerTargetStage, stage.ID().FullyQualifiedName())
	t.Cleanup(revertParameter)

	specFileName := "spec.yaml"
	testClientHelper().Stage.PutInLocationWithContent(t, stage.Location(), specFileName, spec)

	specTemplateFileName := "spec_template.yaml"
	testClientHelper().Stage.PutInLocationWithContent(t, stage.Location(), specTemplateFileName, specTemplate)

	computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
	t.Cleanup(computePoolCleanup)

	networkRule, networkRuleCleanup := testClientHelper().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)

	externalAccessIntegrationId, externalAccessIntegrationCleanup := testClientHelper().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	testClientHelper().Grant.GrantPrivilegesOnComputePoolToAccountRole(t, role.ID(), computePool.ID(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	testClientHelper().Grant.GrantPrivilegesOnIntegrationToAccountRole(t, role.ID(), externalAccessIntegrationId, []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	testClientHelper().Grant.GrantPrivilegesOnWarehouseToAccountRole(t, role.ID(), testClientHelper().Ids.WarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	testClientHelper().Grant.GrantPrivilegesOnSchemaObjectToAccountRole(t, role.ID(), sdk.ObjectTypeStage, stage.ID(), []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeRead, sdk.SchemaObjectPrivilegeWrite}, false)
	testClientHelper().Grant.GrantAllOnDatabaseToAccountRole(t, db.ID(), role.ID())
	testClientHelper().Grant.GrantAllOnSchemaToAccountRole(t, schema.ID(), role.ID())
	testClientHelper().Grant.GrantAllOnDatabaseToAccountRole(t, testClientHelper().Ids.DatabaseId(), role.ID())
	testClientHelper().Grant.GrantAllOnSchemaToAccountRole(t, testClientHelper().Ids.SchemaId(), role.ID())
	testClientHelper().Role.GrantRoleToCurrentRole(t, role.ID())

	useRoleCleanup := testClientHelper().Role.UseRole(t, role.ID())
	t.Cleanup(useRoleCleanup)

	t.Run("create - from specification", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(spec))

		err := client.Services.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(role.ID().Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("create - from specification on stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		location := sdk.NewStageLocation(stage.ID(), "")
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithLocation(location).WithSpecificationFile(specFileName))

		err := client.Services.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(role.ID().Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("create - from specification template", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecificationTemplate(*sdk.NewServiceFromSpecificationTemplateRequest(specTemplateUsing).WithSpecificationTemplate(specTemplate))

		err := client.Services.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(role.ID().Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("create - from specification template with lowercased PYTHON_PROFILER_TARGET_STAGE fails", func(t *testing.T) {
		lowercasedStage, lowercasedStageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(lowercasedStageCleanup)

		revertParameter := testClientHelper().Parameter.UpdateAccountParameterTemporarily(t, sdk.AccountParameterPythonProfilerTargetStage, lowercasedStage.ID().FullyQualifiedName())
		t.Cleanup(revertParameter)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecificationTemplate(*sdk.NewServiceFromSpecificationTemplateRequest(specTemplateUsing).WithSpecificationTemplate(specTemplate))

		err := client.Services.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(role.ID().Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("create - from specification template on stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		location := sdk.NewStageLocation(stage.ID(), "")
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecificationTemplate(*sdk.NewServiceFromSpecificationTemplateRequest(specTemplateUsing).WithLocation(location).WithSpecificationTemplateFile(specTemplateFileName))

		err := client.Services.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(role.ID().Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("create - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		comment := random.Comment()
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(spec)).
			WithAutoSuspendSecs(3600).
			WithExternalAccessIntegrations(*sdk.NewServiceExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegrationId})).
			WithAutoResume(true).
			WithMinInstances(1).
			WithMinReadyInstances(1).
			WithMaxInstances(1).
			WithQueryWarehouse(testClientHelper().Ids.WarehouseId()).
			WithComment(comment)

		err := client.Services.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(role.ID().Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasExternalAccessIntegrations(externalAccessIntegrationId).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(3600).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasQueryWarehouse(testClientHelper().Ids.WarehouseId()).
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("alter: set", func(t *testing.T) {
		comment := random.Comment()
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithSet(sdk.ServiceSetRequest{
			MinReadyInstances:          sdk.Pointer(1),
			MinInstances:               sdk.Pointer(2),
			MaxInstances:               sdk.Pointer(3),
			AutoSuspendSecs:            sdk.Pointer(3600),
			QueryWarehouse:             sdk.Pointer(testClientHelper().Ids.WarehouseId()),
			AutoResume:                 sdk.Pointer(true),
			ExternalAccessIntegrations: sdk.NewServiceExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegrationId}),
			Comment:                    sdk.Pointer(comment),
		}))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasMinReadyInstances(1).
			HasMinInstances(2).
			HasMaxInstances(3).
			HasAutoResume(true).
			HasQueryWarehouse(testClientHelper().Ids.WarehouseId()).
			HasExternalAccessIntegrations(externalAccessIntegrationId).
			HasComment(comment).
			HasAutoSuspendSecs(3600),
		)
	})

	t.Run("alter: unset", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithUnset(sdk.ServiceUnsetRequest{
			AutoResume:                 sdk.Pointer(true),
			MinInstances:               sdk.Pointer(true),
			MaxInstances:               sdk.Pointer(true),
			QueryWarehouse:             sdk.Pointer(true),
			ExternalAccessIntegrations: sdk.Pointer(true),
			Comment:                    sdk.Pointer(true),
			AutoSuspendSecs:            sdk.Pointer(true),
			MinReadyInstances:          sdk.Pointer(true),
		}))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasAutoResume(true).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasNoQueryWarehouse().
			HasNoExternalAccessIntegrations().
			HasNoComment().
			HasAutoSuspendSecs(0).
			HasMinReadyInstances(1),
		)
	})

	t.Run("alter: suspend", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithSuspend(true))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasStatus(sdk.ServiceStatusSuspending).
			HasSuspendedOnNotEmpty(),
		)
	})

	t.Run("alter: resume", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithResume(true))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasStatus(sdk.ServiceStatusPending),
		)
	})

	t.Run("alter: restore", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithIdWithBlockVolume(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithSuspend(true))
		require.NoError(t, err)

		snapshotId, snapshotCleanup := testClientHelper().Snapshot.Create(t, service.ID(), "block-volume")
		t.Cleanup(snapshotCleanup)

		restoreRequest := sdk.NewRestoreRequest("block-volume", []int{0}, snapshotId)

		err = client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithRestore(*restoreRequest))
		require.NoError(t, err)
	})

	t.Run("describe service", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		assertThatObject(t, objectassert.ServiceDetails(t, service.ID()).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(role.Name).
			HasComputePool(computePool.ID()).
			HasSpecNotEmpty().
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("show: with like, exclude jobs", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		services, err := client.Services.Show(ctx, sdk.NewShowServiceRequest().
			WithLike(sdk.Like{Pattern: sdk.Pointer(service.ID().Name())}).
			WithExcludeJobs(true),
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(services))

		assertThatObject(t, objectassert.ServiceFromObject(t, &services[0]).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(role.Name).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("show: in compute pool", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		services, err := client.Services.Show(ctx, sdk.NewShowServiceRequest().
			WithIn(sdk.ServiceIn{ComputePool: computePool.ID()}))
		require.NoError(t, err)
		require.Equal(t, 1, len(services))

		assertThatObject(t, objectassert.ServiceFromObject(t, &services[0]).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(role.Name).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(1).
			HasTargetInstances(1).
			HasMinReadyInstances(1).
			HasMinInstances(1).
			HasMaxInstances(1).
			HasAutoResume(true).
			HasNoExternalAccessIntegrations().
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasNoResumedOn().
			HasNoSuspendedOn().
			HasAutoSuspendSecs(0).
			HasNoComment().
			HasOwnerRoleType("ROLE").
			HasNoQueryWarehouse().
			HasIsJob(false).
			HasIsAsyncJob(false).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})
	// TODO (next PRs): add a test for EXCLUDE JOBS
}
