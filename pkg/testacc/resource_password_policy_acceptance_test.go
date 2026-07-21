//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PasswordPolicy_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(id.SchemaId())
	comment := random.Comment()
	externalComment := random.Comment()

	basic := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())

	altered := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), newId.Name()).
		WithMinLength(10).
		WithMaxLength(30).
		WithMinUpperCaseChars(2).
		WithMinLowerCaseChars(3).
		WithMinNumericChars(4).
		WithMinSpecialChars(5).
		WithMinAgeDays(6).
		WithMaxAgeDays(7).
		WithMaxRetries(8).
		WithLockoutTimeMins(9).
		WithHistory(10).
		WithComment(comment)

	allAttributes := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinLength(10).
		WithMaxLength(30).
		WithMinUpperCaseChars(2).
		WithMinLowerCaseChars(3).
		WithMinNumericChars(4).
		WithMinSpecialChars(5).
		WithMinAgeDays(6).
		WithMaxAgeDays(7).
		WithMaxRetries(8).
		WithLockoutTimeMins(9).
		WithHistory(10).
		WithComment(comment)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		resourceassert.PasswordPolicyResource(t, ref).
			HasName(id.Name()).
			HasDatabase(id.DatabaseName()).
			HasSchema(id.SchemaName()).
			HasFullyQualifiedName(id.FullyQualifiedName()).
			HasCommentEmpty(),
		resourceshowoutputassert.PasswordPolicyShowOutput(t, ref).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind(string(sdk.PolicyKindPasswordPolicy)).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.PasswordPolicyDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasPasswordMinLength(14).
			HasPasswordMaxLength(256).
			HasPasswordMinUpperCaseChars(1).
			HasPasswordMinLowerCaseChars(1).
			HasPasswordMinNumericChars(1).
			HasPasswordMinSpecialChars(0).
			HasPasswordMinAgeDays(0).
			HasPasswordMaxAgeDays(90).
			HasPasswordMaxRetries(5).
			HasPasswordLockoutTimeMins(15).
			HasPasswordHistory(5),
	}

	basicAssertionsWithDefaults := append([]assert.TestCheckFuncProvider{
		resourceassert.PasswordPolicyResource(t, ref).
			HasName(id.Name()).
			HasDatabase(id.DatabaseName()).
			HasSchema(id.SchemaName()).
			HasFullyQualifiedName(id.FullyQualifiedName()).
			HasMinLength(0).
			HasMaxLength(0).
			HasMinUpperCaseChars(-1). // IntDefault sentinel
			HasMinLowerCaseChars(-1). // IntDefault sentinel
			HasMinNumericChars(-1).   // IntDefault sentinel
			HasMinSpecialChars(-1).   // IntDefault sentinel
			HasMinAgeDays(-1).        // IntDefault sentinel
			HasMaxAgeDays(-1).        // IntDefault sentinel
			HasMaxRetries(0).
			HasLockoutTimeMins(0).
			HasHistory(-1). // IntDefault sentinel
			HasCommentEmpty(),
	}, basicAssertions[1:]...)

	alteredAssertions := []assert.TestCheckFuncProvider{
		resourceassert.PasswordPolicyResource(t, ref).
			HasName(newId.Name()).
			HasDatabase(newId.DatabaseName()).
			HasSchema(newId.SchemaName()).
			HasMinLength(10).
			HasMaxLength(30).
			HasMinUpperCaseChars(2).
			HasMinLowerCaseChars(3).
			HasMinNumericChars(4).
			HasMinSpecialChars(5).
			HasMinAgeDays(6).
			HasMaxAgeDays(7).
			HasMaxRetries(8).
			HasLockoutTimeMins(9).
			HasHistory(10).
			HasComment(comment),
		resourceshowoutputassert.PasswordPolicyShowOutput(t, ref).
			HasName(newId.Name()).
			HasCreatedOnNotEmpty().
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind(string(sdk.PolicyKindPasswordPolicy)).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.PasswordPolicyDescribeOutput(t, ref).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasPasswordMinLength(10).
			HasPasswordMaxLength(30).
			HasPasswordMinUpperCaseChars(2).
			HasPasswordMinLowerCaseChars(3).
			HasPasswordMinNumericChars(4).
			HasPasswordMinSpecialChars(5).
			HasPasswordMinAgeDays(6).
			HasPasswordMaxAgeDays(7).
			HasPasswordMaxRetries(8).
			HasPasswordLockoutTimeMins(9).
			HasPasswordHistory(10),
	}

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.PasswordPolicyResource(t, ref).
			HasName(id.Name()).
			HasDatabase(id.DatabaseName()).
			HasSchema(id.SchemaName()).
			HasMinLength(10).
			HasMaxLength(30).
			HasMinUpperCaseChars(2).
			HasMinLowerCaseChars(3).
			HasMinNumericChars(4).
			HasMinSpecialChars(5).
			HasMinAgeDays(6).
			HasMaxAgeDays(7).
			HasMaxRetries(8).
			HasLockoutTimeMins(9).
			HasHistory(10).
			HasComment(comment),
		resourceshowoutputassert.PasswordPolicyShowOutput(t, ref).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind(string(sdk.PolicyKindPasswordPolicy)).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.PasswordPolicyDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasPasswordMinLength(10).
			HasPasswordMaxLength(30).
			HasPasswordMinUpperCaseChars(2).
			HasPasswordMinLowerCaseChars(3).
			HasPasswordMinNumericChars(4).
			HasPasswordMinSpecialChars(5).
			HasPasswordMinAgeDays(6).
			HasPasswordMaxAgeDays(7).
			HasPasswordMaxRetries(8).
			HasPasswordLockoutTimeMins(9).
			HasPasswordHistory(10),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PasswordPolicy),
		Steps: []resource.TestStep{
			// Create
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"min_length", "max_length",
					"min_upper_case_chars", "min_lower_case_chars",
					"min_numeric_chars", "min_special_chars",
					"min_age_days", "max_age_days",
					"max_retries", "lockout_time_mins",
					"history",
					"or_replace", "if_not_exists",
				},
			},
			// Set all optional fields
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// External drift (all fields)
			{
				PreConfig: func() {
					testClient().PasswordPolicy.Alter(t, sdk.NewAlterPasswordPolicyRequest(id).WithSet(
						*sdk.NewPasswordPolicySetRequest().
							WithPasswordMinLength(15).
							WithPasswordMaxLength(100).
							WithPasswordMinUpperCaseChars(3).
							WithPasswordMinLowerCaseChars(4).
							WithPasswordMinNumericChars(5).
							WithPasswordMinSpecialChars(6).
							WithPasswordMinAgeDays(7).
							WithPasswordMaxAgeDays(8).
							WithPasswordMaxRetries(9).
							WithPasswordLockoutTimeMins(10).
							WithPasswordHistory(11).
							WithComment(externalComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: func() []plancheck.PlanCheck {
						return []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
							planchecks.ExpectDrift(ref, "min_length", sdk.String("10"), sdk.String("15")),
							planchecks.ExpectChange(ref, "min_length", tfjson.ActionUpdate, sdk.String("15"), sdk.String("10")),
							planchecks.ExpectDrift(ref, "max_length", sdk.String("30"), sdk.String("100")),
							planchecks.ExpectChange(ref, "max_length", tfjson.ActionUpdate, sdk.String("100"), sdk.String("30")),
							planchecks.ExpectDrift(ref, "min_upper_case_chars", sdk.String("2"), sdk.String("3")),
							planchecks.ExpectChange(ref, "min_upper_case_chars", tfjson.ActionUpdate, sdk.String("3"), sdk.String("2")),
							planchecks.ExpectDrift(ref, "min_lower_case_chars", sdk.String("3"), sdk.String("4")),
							planchecks.ExpectChange(ref, "min_lower_case_chars", tfjson.ActionUpdate, sdk.String("4"), sdk.String("3")),
							planchecks.ExpectDrift(ref, "min_numeric_chars", sdk.String("4"), sdk.String("5")),
							planchecks.ExpectChange(ref, "min_numeric_chars", tfjson.ActionUpdate, sdk.String("5"), sdk.String("4")),
							planchecks.ExpectDrift(ref, "min_special_chars", sdk.String("5"), sdk.String("6")),
							planchecks.ExpectChange(ref, "min_special_chars", tfjson.ActionUpdate, sdk.String("6"), sdk.String("5")),
							planchecks.ExpectDrift(ref, "min_age_days", sdk.String("6"), sdk.String("7")),
							planchecks.ExpectChange(ref, "min_age_days", tfjson.ActionUpdate, sdk.String("7"), sdk.String("6")),
							planchecks.ExpectDrift(ref, "max_age_days", sdk.String("7"), sdk.String("8")),
							planchecks.ExpectChange(ref, "max_age_days", tfjson.ActionUpdate, sdk.String("8"), sdk.String("7")),
							planchecks.ExpectDrift(ref, "max_retries", sdk.String("8"), sdk.String("9")),
							planchecks.ExpectChange(ref, "max_retries", tfjson.ActionUpdate, sdk.String("9"), sdk.String("8")),
							planchecks.ExpectDrift(ref, "lockout_time_mins", sdk.String("9"), sdk.String("10")),
							planchecks.ExpectChange(ref, "lockout_time_mins", tfjson.ActionUpdate, sdk.String("10"), sdk.String("9")),
							planchecks.ExpectDrift(ref, "history", sdk.String("10"), sdk.String("11")),
							planchecks.ExpectChange(ref, "history", tfjson.ActionUpdate, sdk.String("11"), sdk.String("10")),
							planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externalComment)),
							planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externalComment), sdk.String(comment)),
						}
					}(),
				},
				Config: config.FromModels(t, allAttributes),
				Check:  assertThat(t, completeAssertions...),
			},
			// Rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, altered),
				Check:  assertThat(t, alteredAssertions...),
			},
			// Unset + rename back
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertionsWithDefaults...),
			},
		},
	})
}

func TestAcc_PasswordPolicy_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	completeModel := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinLength(10).
		WithMaxLength(30).
		WithMinUpperCaseChars(2).
		WithMinLowerCaseChars(3).
		WithMinNumericChars(4).
		WithMinSpecialChars(5).
		WithMinAgeDays(6).
		WithMaxAgeDays(7).
		WithMaxRetries(8).
		WithLockoutTimeMins(9).
		WithHistory(10).
		WithComment(comment)
	ref := completeModel.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		resourceassert.PasswordPolicyResource(t, ref).
			HasName(id.Name()).
			HasDatabase(id.DatabaseName()).
			HasSchema(id.SchemaName()).
			HasMinLength(10).
			HasMaxLength(30).
			HasMinUpperCaseChars(2).
			HasMinLowerCaseChars(3).
			HasMinNumericChars(4).
			HasMinSpecialChars(5).
			HasMinAgeDays(6).
			HasMaxAgeDays(7).
			HasMaxRetries(8).
			HasLockoutTimeMins(9).
			HasHistory(10).
			HasComment(comment),
		resourceshowoutputassert.PasswordPolicyShowOutput(t, ref).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind(string(sdk.PolicyKindPasswordPolicy)).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""),
		resourceshowoutputassert.PasswordPolicyDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasPasswordMinLength(10).
			HasPasswordMaxLength(30).
			HasPasswordMinUpperCaseChars(2).
			HasPasswordMinLowerCaseChars(3).
			HasPasswordMinNumericChars(4).
			HasPasswordMinSpecialChars(5).
			HasPasswordMinAgeDays(6).
			HasPasswordMaxAgeDays(7).
			HasPasswordMaxRetries(8).
			HasPasswordLockoutTimeMins(9).
			HasPasswordHistory(10),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PasswordPolicy),
		Steps: []resource.TestStep{
			// Create with all attributes
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, completeModel),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import
			{
				Config:            config.FromModels(t, completeModel),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"min_length", "max_length",
					"min_upper_case_chars", "min_lower_case_chars",
					"min_numeric_chars", "min_special_chars",
					"min_age_days", "max_age_days",
					"max_retries", "lockout_time_mins",
					"history",
					"or_replace", "if_not_exists",
				},
			},
		},
	})
}

func TestAcc_PasswordPolicy_migrateFromVersion_0_94_1(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	providerConfig := providermodel.V097CompatibleProviderConfig(t)

	basicModel := model.PasswordPolicy("pa", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigWithServiceUserPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.94.1"),
				Config:            providerConfig + config.FromModels(t, basicModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "name", id.Name()),
					resource.TestCheckResourceAttr(ref, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, basicModel),
				Check: assertThat(
					t,
					resourceassert.PasswordPolicyResource(t, ref).
						HasName(id.Name()).
						HasFullyQualifiedName(id.FullyQualifiedName()),
					resourceshowoutputassert.PasswordPolicyShowOutput(t, ref).
						HasName(id.Name()),
					resourceshowoutputassert.PasswordPolicyDescribeOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()),
				),
			},
		},
	})
}

func TestAcc_PasswordPolicy_migrateFromVersion_2_15_0(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.PasswordPolicyResource))

	basicModel := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name())
	ref := basicModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PasswordPolicy),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.15.0"),
				Config:            config.FromModels(t, providerModel, basicModel),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, basicModel),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", id.FullyQualifiedName())),
					resourceassert.PasswordPolicyResource(t, ref).
						HasName(id.Name()).
						HasDatabase(id.DatabaseName()).
						HasSchema(id.SchemaName()).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasCommentEmpty(),
					resourceshowoutputassert.PasswordPolicyShowOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()),
					resourceshowoutputassert.PasswordPolicyDescribeOutput(t, ref).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()),
				),
			},
		},
	})
}

func TestAcc_PasswordPolicy_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	// Fields without IntDefault (minimum 1): min_length, max_length, max_retries, lockout_time_mins
	modelInvalidMinLength := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinLength(0)
	modelInvalidMaxLength := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMaxLength(0)
	modelInvalidMaxRetries := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMaxRetries(0)
	modelInvalidLockoutTimeMins := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithLockoutTimeMins(0)

	// Fields with IntDefault (minimum 0): min_upper_case_chars, min_lower_case_chars, min_numeric_chars, min_special_chars, min_age_days, max_age_days, history
	modelInvalidMinUpperCaseChars := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinUpperCaseChars(-1)
	modelInvalidMinLowerCaseChars := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinLowerCaseChars(-1)
	modelInvalidMinNumericChars := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinNumericChars(-1)
	modelInvalidMinSpecialChars := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinSpecialChars(-1)
	modelInvalidMinAgeDays := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMinAgeDays(-1)
	modelInvalidMaxAgeDays := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithMaxAgeDays(-1)
	modelInvalidHistory := model.PasswordPolicy("test", id.DatabaseName(), id.SchemaName(), id.Name()).
		WithHistory(-1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.PasswordPolicy),
		Steps: []resource.TestStep{
			// Fields with IntAtLeast(1)
			{
				Config:      config.FromModels(t, modelInvalidMinLength),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_length to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMaxLength),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_length to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMaxRetries),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_retries to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelInvalidLockoutTimeMins),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected lockout_time_mins to be at least \(1\), got 0`),
			},
			// Fields with IntAtLeast(0)
			{
				Config:      config.FromModels(t, modelInvalidMinUpperCaseChars),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_upper_case_chars to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMinLowerCaseChars),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_lower_case_chars to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMinNumericChars),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_numeric_chars to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMinSpecialChars),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_special_chars to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMinAgeDays),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_age_days to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidMaxAgeDays),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected max_age_days to be at least \(0\), got -1`),
			},
			{
				Config:      config.FromModels(t, modelInvalidHistory),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected history to be at least \(0\), got -1`),
			},
		},
	})
}
