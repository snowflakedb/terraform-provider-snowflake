//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"strconv"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StorageLifecyclePolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(id.SchemaId())

	comment := random.Comment()
	externalComment := random.Comment()

	arguments := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar_200,
		},
	}
	newArguments := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar_200,
		},
		{
			Name: "ID",
			Type: testdatatypes.DataTypeVectorFloat_768,
		},
	}

	expectedSignature := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar,
		},
	}
	expectedNewSignature := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar,
		},
		{
			Name: "ID",
			Type: testdatatypes.DataTypeVectorFloat_768,
		},
	}
	importedArguments := expectedSignature

	body := "LENGTH(VAL) > 0"
	newBody := "LENGTH(VAL) > 5"
	externalBody := "LENGTH(VAL) > 10"

	archiveTier := string(sdk.StorageLifecyclePolicyArchiveTierCold)
	externalArchiveTier := string(sdk.StorageLifecyclePolicyArchiveTierCool)

	archiveForDays := 365
	externalArchiveForDays := 200

	basic := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, body)

	complete := model.StorageLifecyclePolicy("t", newId.DatabaseName(), newId.SchemaName(), newId.Name(), arguments, newBody).
		WithArchiveTier(archiveTier).
		WithArchiveForDays(archiveForDays).
		WithComment(comment)

	withArchiveTier := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, body).
		WithArchiveTier(archiveTier)

	withNewArguments := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), newArguments, body)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.StorageLifecyclePolicyResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasArguments(arguments).
			HasBody(body).
			HasArchiveTier("").
			HasArchiveForDays(0).
			HasComment(""),
		resourceshowoutputassert.StorageLifecyclePolicyShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("STORAGE_LIFECYCLE_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(`{"ARCHIVE_FOR_DAYS":null,"ARCHIVE_TIER":"NULL"}`),
		resourceshowoutputassert.StorageLifecyclePolicyDescribeOutput(t, ref).
			HasName(id.Name()).
			HasSignature(expectedSignature...).
			HasReturnType(testdatatypes.DataTypeBoolean).
			HasBody(body).
			HasArchiveTier("").
			HasArchiveForDays(0),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.StorageLifecyclePolicyResource(t, ref).
			HasName(newId.Name()).
			HasSchema(newId.SchemaName()).
			HasDatabase(newId.DatabaseName()).
			HasArguments(arguments).
			HasBody(newBody).
			HasArchiveTier(archiveTier).
			HasArchiveForDays(archiveForDays).
			HasComment(comment),
		resourceshowoutputassert.StorageLifecyclePolicyShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("STORAGE_LIFECYCLE_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(`{"ARCHIVE_FOR_DAYS":365,"ARCHIVE_TIER":"COLD"}`),
		resourceshowoutputassert.StorageLifecyclePolicyDescribeOutput(t, ref).
			HasName(newId.Name()).
			HasSignature(expectedSignature...).
			HasReturnType(testdatatypes.DataTypeBoolean).
			HasBody(newBody).
			HasArchiveTier(archiveTier).
			HasArchiveForDays(archiveForDays),
	}

	withArchiveTierAssertions := []assert.TestCheckFuncProvider{
		resourceassert.StorageLifecyclePolicyResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasArguments(arguments).
			HasBody(body).
			HasArchiveTier(archiveTier).
			HasArchiveForDays(0).
			HasComment(""),
		resourceshowoutputassert.StorageLifecyclePolicyShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("STORAGE_LIFECYCLE_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(`{"ARCHIVE_FOR_DAYS":null,"ARCHIVE_TIER":"COLD"}`),
		resourceshowoutputassert.StorageLifecyclePolicyDescribeOutput(t, ref).
			HasName(id.Name()).
			HasSignature(expectedSignature...).
			HasReturnType(testdatatypes.DataTypeBoolean).
			HasBody(body).
			HasArchiveTier(archiveTier).
			HasArchiveForDays(0),
	}

	withNewArgumentsAssertions := []assert.TestCheckFuncProvider{
		resourceassert.StorageLifecyclePolicyResource(t, ref).
			HasName(id.Name()).
			HasSchema(id.SchemaName()).
			HasDatabase(id.DatabaseName()).
			HasArguments(newArguments).
			HasBody(body).
			HasArchiveTier("").
			HasArchiveForDays(0).
			HasComment(""),
		resourceshowoutputassert.StorageLifecyclePolicyShowOutput(t, ref).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("STORAGE_LIFECYCLE_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(`{"ARCHIVE_FOR_DAYS":null,"ARCHIVE_TIER":"NULL"}`),
		resourceshowoutputassert.StorageLifecyclePolicyDescribeOutput(t, ref).
			HasName(id.Name()).
			HasSignature(expectedNewSignature...).
			HasReturnType(testdatatypes.DataTypeBoolean).
			HasBody(body).
			HasArchiveTier("").
			HasArchiveForDays(0),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.StorageLifecyclePolicy),
		Steps: []resource.TestStep{
			// Create without optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import without optionals
			{
				Config:       config.FromModels(t, basic),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStorageLifecyclePolicyResource(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasSchema(id.SchemaName()).
						HasDatabase(id.DatabaseName()).
						HasArguments(importedArguments).
						HasBody(body).
						HasArchiveTier("").
						HasArchiveForDays(0).
						HasComment(""),
				),
			},
			// Set all optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Unset all optionals (except for archive_tier)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withArchiveTier),
				Check:  assertThat(t, withArchiveTierAssertions...),
			},
			// Unset archive_tier
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Change arguments
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, withNewArguments),
				Check:  assertThat(t, withNewArgumentsAssertions...),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config:  config.FromModels(t, basic),
				Destroy: true,
			},
			// Create with all optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import with all optionals
			{
				Config:       config.FromModels(t, complete),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedStorageLifecyclePolicyResource(t, helpers.EncodeResourceIdentifier(newId)).
						HasName(newId.Name()).
						HasSchema(newId.SchemaName()).
						HasDatabase(newId.DatabaseName()).
						HasArguments(importedArguments).
						HasBody(newBody).
						HasArchiveTier(archiveTier).
						HasArchiveForDays(archiveForDays).
						HasComment(comment),
				),
			},
			// Change all props externally (except for archive_tier)
			{
				PreConfig: func() {
					alterRequest := sdk.NewAlterStorageLifecyclePolicyRequest(newId).WithSetBody(externalBody)
					testClient().StorageLifecyclePolicy.Alter(t, alterRequest)

					alterRequest = sdk.NewAlterStorageLifecyclePolicyRequest(newId).WithSet(
						*sdk.NewStorageLifecyclePolicySetRequest().
							WithArchiveForDays(externalArchiveForDays).
							WithComment(externalComment),
					)
					testClient().StorageLifecyclePolicy.Alter(t, alterRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: func() []plancheck.PlanCheck {
						archiveForDays := new(strconv.Itoa(archiveForDays))
						externalArchiveForDays := new(strconv.Itoa(externalArchiveForDays))
						return []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
							planchecks.ExpectDrift(ref, "body", new(newBody), new(externalBody)),
							planchecks.ExpectChange(ref, "body", tfjson.ActionUpdate, new(externalBody), new(newBody)),
							planchecks.ExpectDrift(ref, "archive_for_days", archiveForDays, externalArchiveForDays),
							planchecks.ExpectChange(ref, "archive_for_days", tfjson.ActionUpdate, externalArchiveForDays, archiveForDays),
							planchecks.ExpectDrift(ref, "comment", new(comment), new(externalComment)),
							planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, new(externalComment), new(comment)),
						}
					}(),
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Change archive_tier and arguments externally
			{
				PreConfig: func() {
					replaceRequest := sdk.NewCreateStorageLifecyclePolicyRequest(newId,
						[]sdk.CreateStorageLifecyclePolicyArgsRequest{{
							Name:     "VAL",
							DataType: testdatatypes.DataTypeNumber,
						}}, newBody).
						WithOrReplace(true).
						WithArchiveTier(sdk.StorageLifecyclePolicyArchiveTierCool).
						WithArchiveForDays(archiveForDays).
						WithComment(comment)
					testClient().StorageLifecyclePolicy.CreateWithRequest(t, newId, replaceRequest)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: func() []plancheck.PlanCheck {
						oldType := new("VARCHAR(200)")
						newType := new("NUMBER")
						return []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
							planchecks.ExpectDrift(ref, "argument.0.type", oldType, newType),
							planchecks.ExpectChange(ref, "argument.0.type", tfjson.ActionDelete, newType, oldType),
							planchecks.ExpectNoChangeOnField(ref, "argument.0.name"),
							planchecks.ExpectDrift(ref, "archive_tier", new(archiveTier), new(externalArchiveTier)),
							planchecks.ExpectChange(ref, "archive_tier", tfjson.ActionDelete, new(externalArchiveTier), new(archiveTier)),
						}
					}(),
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
		},
	})
}

func TestAcc_StorageLifecyclePolicy_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	arguments := []sdk.TableColumnSignature{
		{
			Name: "VAL",
			Type: testdatatypes.DataTypeVarchar_200,
		},
	}
	body := "LENGTH(VAL) > 0"

	archiveForDaysWithoutArchiveTier := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, body).
		WithArchiveForDays(365)

	archiveForDaysBelowMin := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, body).
		WithArchiveTier(string(sdk.StorageLifecyclePolicyArchiveTierCold)).
		WithArchiveForDays(0)

	emptyBody := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, "")

	invalidArchiveTier := model.StorageLifecyclePolicy("t", id.DatabaseName(), id.SchemaName(), id.Name(), arguments, body).
		WithArchiveTier("INVALID")

	invalidDataTypeVariableModel := config.SetMapStringVariable("arguments")
	invalidDataType := model.StorageLifecyclePolicyDynamicArguments("t", id, body)
	invalidDataTypeVariables := tfconfig.Variables{
		"arguments": tfconfig.SetVariable(
			tfconfig.MapVariable(map[string]tfconfig.Variable{
				"name": tfconfig.StringVariable("VAL"),
				"type": tfconfig.StringVariable("invalid-type"),
			}),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, archiveForDaysWithoutArchiveTier),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("\"archive_for_days\": all of `archive_for_days,archive_tier` must be specified"),
			},
			{
				Config:      config.FromModels(t, archiveForDaysBelowMin),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected archive_for_days to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, emptyBody),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "body" to not be an empty string`),
			},
			{
				Config:      config.FromModels(t, invalidArchiveTier),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid storage lifecycle policy archive tier: INVALID`),
			},
			{
				Config:          config.FromModels(t, invalidDataTypeVariableModel, invalidDataType),
				ConfigVariables: invalidDataTypeVariables,
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile(`invalid data type: invalid-type`),
			},
		},
	})
}
