//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Services(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	spec := testClientHelper().Service.SampleSpecWithContainerName(t, "text-example-original")

	changedSpec := testClientHelper().Service.SampleSpecWithContainerName(t, "text-example-changed")
	specTemplate := testClientHelper().Service.SampleSpecWithContainerName(t, "template-example-{{ container_name }}")
	specTemplateUsing := []sdk.ListItem{
		{Key: "container_name", Value: `'original'`},
	}
	specTemplateUsingChanged := []sdk.ListItem{
		{Key: "container_name", Value: `'changed'`},
	}
	// TODO(SNOW-2129575): We set up a separate database and schema with capitalized ids. Remove this after fix on snowflake side.
	db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStageInSchema(t, schema.ID())
	t.Cleanup(stageCleanup)
	location := sdk.NewStageLocation(stage.ID(), "")

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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))
		// TODO(SNOW-2129575): When we set a stage with lowercase characters, we get the following error:
		// 395069 (23001): Unable to render service spec from given template: Stage '\"INT_TEST_DB_IT_51244746_D66C_B4C5_3632_636F292BB1FB\".\"INT_TEST_SC_IT_51244746_D66C_B4C5_3632_636F292BB1FB\".\"JXRIVBIT_51244746_D66C_B4C5_3632_636F292BB1FB\"' not found for unloading profiler data
		// However, this behavior does not seem to be consistent across accounts. In accounts used in pipelines, it gives no error. Both of the accounts have QUOTED_IDENTIFIERS_IGNORE_CASE set to false.
		require.NoError(t, err)
	})

	t.Run("create - from specification template on stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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

	t.Run("alter: change spec", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		// specification
		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(changedSpec)))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceDetails(t, service.ID()).
			HasName(service.ID().Name()).
			HasSpecThatContains("text-example-changed"),
		)

		// specification on stage
		err = client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithLocation(location).WithSpecificationFile(specFileName)))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceDetails(t, service.ID()).
			HasName(service.ID().Name()).
			HasSpecThatContains("text-example-original"),
		)

		// specification template
		err = client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithFromSpecificationTemplate(*sdk.NewServiceFromSpecificationTemplateRequest(specTemplateUsing).WithSpecificationTemplate(specTemplate)))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceDetails(t, service.ID()).
			HasName(service.ID().Name()).
			HasSpecThatContains("template-example-original"),
		)

		// specification template on stage
		err = client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithFromSpecificationTemplate(*sdk.NewServiceFromSpecificationTemplateRequest(specTemplateUsingChanged).WithLocation(location).WithSpecificationTemplateFile(specTemplateFileName)))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceDetails(t, service.ID()).
			HasName(service.ID().Name()).
			HasSpecThatContains("template-example-changed"),
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
		comment := random.Comment()
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateServiceRequest(id, computePool.ID()).
			WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(spec)).
			WithAutoResume(true).
			WithMinInstances(2).
			WithMaxInstances(3).
			WithQueryWarehouse(testClientHelper().Ids.WarehouseId()).
			WithExternalAccessIntegrations(*sdk.NewServiceExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegrationId})).
			WithComment(comment).
			WithAutoSuspendSecs(3600).
			WithMinReadyInstances(1)

		service, serviceCleanup := testClientHelper().Service.CreateWithRequest(t, request)
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

	t.Run("alter: suspend and resume", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasStatus(sdk.ServiceStatusPending).
			HasNoResumedOn().
			HasNoSuspendedOn(),
		)

		err := client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithSuspend(true))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasStatus(sdk.ServiceStatusSuspending).
			HasNoResumedOn().
			HasSuspendedOnNotEmpty(),
		)

		err = client.Services.Alter(ctx, sdk.NewAlterServiceRequest(service.ID()).WithResume(true))
		require.NoError(t, err)

		service, err = client.Services.ShowByID(ctx, service.ID())
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasStatus(sdk.ServiceStatusPending).
			HasResumedOnNotEmpty().
			HasSuspendedOnNotEmpty(),
		)
	})

	t.Run("execute job service - from specification template on stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewExecuteJobServiceRequest(computePool.ID(), id).
			WithJobServiceFromSpecificationTemplate(*sdk.NewJobServiceFromSpecificationTemplateRequest(specTemplateUsing).WithLocation(location).WithSpecificationTemplateFile(specTemplateFileName)).
			WithAsync(true)

		err := client.Services.ExecuteJob(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(id.Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(0).
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
			HasIsJob(true).
			HasIsAsyncJob(true).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("execute job service - basic, from stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewExecuteJobServiceRequest(computePool.ID(), id).
			WithJobServiceFromSpecification(*sdk.NewJobServiceFromSpecificationRequest().WithLocation(location).WithSpecificationFile(specFileName)).
			WithAsync(true)

		err := client.Services.ExecuteJob(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(0).
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
			HasIsJob(true).
			HasIsAsyncJob(true).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	t.Run("execute job service - complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		comment := random.Comment()
		request := sdk.NewExecuteJobServiceRequest(computePool.ID(), id).
			WithJobServiceFromSpecification(*sdk.NewJobServiceFromSpecificationRequest().WithSpecification(spec)).
			WithAsync(true).
			WithExternalAccessIntegrations(*sdk.NewServiceExternalAccessIntegrationsRequest([]sdk.AccountObjectIdentifier{externalAccessIntegrationId})).
			WithQueryWarehouse(testClientHelper().Ids.WarehouseId()).
			WithComment(comment)

		err := client.Services.ExecuteJob(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Service.DropFunc(t, id))

		service, err := client.Services.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.ServiceFromObject(t, service).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(0).
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
			HasAutoSuspendSecs(0).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasQueryWarehouse(testClientHelper().Ids.WarehouseId()).
			HasIsJob(true).
			HasIsAsyncJob(true).
			HasSpecDigestNotEmpty().
			HasIsUpgrading(false).
			HasNoManagingObjectDomain().
			HasNoManagingObjectName(),
		)
	})

	// TODO(SNOW-2132078): Add an integration test for restoring a service from a snapshot.

	// TODO(SNOW-2132087): Add integration tests for creating and altering services in Native Apps.

	t.Run("describe service", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.CreateWithId(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		assertThatObject(t, objectassert.ServiceDetails(t, service.ID()).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComputePool(computePool.ID()).
			HasSpecThatContains("snowflake/images/snowflake_images/exampleimage:latest").
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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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

	t.Run("show: with like, only jobs", func(t *testing.T) {
		service, serviceCleanup := testClientHelper().Service.ExecuteJobService(t, computePool.ID(), testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID()))
		t.Cleanup(serviceCleanup)

		services, err := client.Services.Show(ctx, sdk.NewShowServiceRequest().
			WithLike(sdk.Like{Pattern: sdk.Pointer(service.ID().Name())}).
			WithJob(true),
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(services))

		assertThatObject(t, objectassert.ServiceFromObject(t, &services[0]).
			HasName(service.ID().Name()).
			HasStatus(sdk.ServiceStatusPending).
			HasDatabaseName(service.ID().DatabaseName()).
			HasSchemaName(service.ID().SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComputePool(computePool.ID()).
			HasDnsNameNotEmpty().
			HasCurrentInstances(0).
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
			HasIsAsyncJob(true).
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
			HasOwner(snowflakeroles.Accountadmin.Name()).
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
}
