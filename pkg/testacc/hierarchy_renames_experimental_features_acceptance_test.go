//go:build non_account_level_tests

package testacc

import (
	"errors"
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

type hierarchyRenameSchemaObjectTestCase struct {
	name     string
	resource resources.Resource
	model    func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel
	exists   func(t *testing.T, id sdk.SchemaObjectIdentifier) error
}

func hierarchyRenameSchemaObjectTestCases(t *testing.T) []hierarchyRenameSchemaObjectTestCase {
	t.Helper()

	table, cleanupTable := testClient().Table.CreateWithPredefinedColumns(t)
	t.Cleanup(cleanupTable)

	return []hierarchyRenameSchemaObjectTestCase{
		{
			name:     "AuthenticationPolicy",
			resource: resources.AuthenticationPolicy,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.AuthenticationPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().AuthenticationPolicy.Show),
		},
		{
			name:     "BackupPolicy",
			resource: resources.BackupPolicy,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.BackupPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name()).
					WithExpireAfterDays(7)
			},
			exists: showError(testClient().BackupPolicy.Show),
		},
		{
			name:     "BackupSet",
			resource: resources.BackupSet,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.BackupSet("t", id.DatabaseName(), id.SchemaName(), id.Name()).
					WithForTable(table.ID().FullyQualifiedName())
			},
			exists: showError(testClient().BackupSet.Show),
		},
		{
			name:     "FileFormatAvro",
			resource: resources.FileFormatAvro,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.FileFormatAvro("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().FileFormat.Show),
		},
		{
			name:     "FileFormatCsv",
			resource: resources.FileFormatCsv,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.FileFormatCsv("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().FileFormat.Show),
		},
		{
			name:     "FileFormatJson",
			resource: resources.FileFormatJson,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.FileFormatJson("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().FileFormat.Show),
		},
		{
			name:     "FileFormatOrc",
			resource: resources.FileFormatOrc,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.FileFormatOrc("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().FileFormat.Show),
		},
		{
			name:     "FileFormatParquet",
			resource: resources.FileFormatParquet,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.FileFormatParquet("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().FileFormat.Show),
		},
		{
			name:     "FileFormatXml",
			resource: resources.FileFormatXml,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.FileFormatXml("t", id.DatabaseName(), id.SchemaName(), id.Name())
			},
			exists: showError(testClient().FileFormat.Show),
		},
		{
			name:     "HybridTable",
			resource: resources.HybridTable,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				columns := []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber}}
				pks := []sdk.TableColumnSignature{{Name: "ID"}}
				return model.HybridTableFromId("t", id, columns, pks)
			},
			exists: showError(testClient().HybridTable.Show),
		},
		{
			name:     "ProjectionPolicy",
			resource: resources.ProjectionPolicy,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.ProjectionPolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), "PROJECTION_CONSTRAINT(ALLOW => false)")
			},
			exists: showError(testClient().ProjectionPolicy.Show),
		},
		{
			name:     "Table",
			resource: resources.Table,
			model: func(id sdk.SchemaObjectIdentifier) accconfig.ResourceModel {
				return model.Table("t", id.DatabaseName(), id.SchemaName(), id.Name(), []sdk.TableColumnSignature{{Name: "ID", Type: testdatatypes.DataTypeNumber}})
			},
			exists: showError(testClient().Table.Show),
		},
	}
}

func TestAcc_Experimental_HierarchyRenames_MoveToAnotherSchema(t *testing.T) {
	providerModel := providermodel.SnowflakeProvider().WithExperimentalFeaturesEnabled(experimentalfeatures.HierarchyRenames)

	for _, tc := range hierarchyRenameSchemaObjectTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			schemaX, cleanupSchemaX := testClient().Schema.CreateSchema(t)
			t.Cleanup(cleanupSchemaX)
			schemaY, cleanupSchemaY := testClient().Schema.CreateSchema(t)
			t.Cleanup(cleanupSchemaY)

			id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaX.ID())
			movedId := testClient().Ids.NewSchemaObjectIdentifierInSchema(id.Name(), schemaY.ID())
			before := tc.model(id)
			after := tc.model(movedId)
			ref := before.ResourceReference()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: experimentalHierarchyRenamesProviderFactory,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, tc.resource),
				Steps: []resource.TestStep{
					{
						Config: accconfig.FromModels(t, providerModel, before),
						Check:  checkHierarchyRenameIdentity(t, tc, ref, id),
					},
					{
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
							},
							PostApplyPostRefresh: []plancheck.PlanCheck{
								plancheck.ExpectEmptyPlan(),
							},
						},
						Config: accconfig.FromModels(t, providerModel, after),
						Check: resource.ComposeTestCheckFunc(
							checkHierarchyRenameIdentity(t, tc, ref, movedId),
							checkHierarchyRenameObjectMissing(t, tc, id),
						),
					},
				},
			})
		})
	}
}

func TestAcc_Experimental_HierarchyRenames_Disabled_ForceRecreation(t *testing.T) {
	for _, tc := range hierarchyRenameSchemaObjectTestCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			schemaX, cleanupSchemaX := testClient().Schema.CreateSchema(t)
			t.Cleanup(cleanupSchemaX)
			schemaY, cleanupSchemaY := testClient().Schema.CreateSchema(t)
			t.Cleanup(cleanupSchemaY)

			id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaX.ID())
			movedId := testClient().Ids.NewSchemaObjectIdentifierInSchema(id.Name(), schemaY.ID())
			before := tc.model(id)
			after := tc.model(movedId)
			ref := before.ResourceReference()

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: CheckDestroy(t, tc.resource),
				Steps: []resource.TestStep{
					{
						Config: accconfig.FromModels(t, before),
						Check:  checkHierarchyRenameIdentity(t, tc, ref, id),
					},
					{
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PostApplyPostRefresh: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
							},
						},
						Config:             accconfig.FromModels(t, after),
						PlanOnly:           true,
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
	}
}

func checkHierarchyRenameIdentity(t *testing.T, tc hierarchyRenameSchemaObjectTestCase, ref string, id sdk.SchemaObjectIdentifier) resource.TestCheckFunc {
	t.Helper()
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(ref, "name", id.Name()),
		resource.TestCheckResourceAttr(ref, "database", id.DatabaseName()),
		resource.TestCheckResourceAttr(ref, "schema", id.SchemaName()),
		resource.TestCheckResourceAttr(ref, "fully_qualified_name", id.FullyQualifiedName()),
		func(_ *terraform.State) error {
			if err := tc.exists(t, id); err != nil {
				return fmt.Errorf("%s %s: %w", tc.name, id.FullyQualifiedName(), err)
			}
			return nil
		},
	)
}

func checkHierarchyRenameObjectMissing(t *testing.T, tc hierarchyRenameSchemaObjectTestCase, id sdk.SchemaObjectIdentifier) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		err := tc.exists(t, id)
		if err == nil {
			return fmt.Errorf("%s %s still exists after the move", tc.name, id.FullyQualifiedName())
		}
		if !errors.Is(err, sdk.ErrObjectNotFound) {
			return fmt.Errorf("%s %s: %w", tc.name, id.FullyQualifiedName(), err)
		}
		return nil
	}
}

func showError[T any](show func(*testing.T, sdk.SchemaObjectIdentifier) (T, error)) func(*testing.T, sdk.SchemaObjectIdentifier) error {
	return func(t *testing.T, id sdk.SchemaObjectIdentifier) error {
		t.Helper()
		_, err := show(t, id)
		return err
	}
}
