//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_User_RsaPublicKeyWriteOnly(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	key1, key1Fp := random.GenerateRSAPublicKey(t)
	key2, key2Fp := random.GenerateRSAPublicKey(t)

	userModelWithKey1 := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key1).
		WithRsaPublicKeyWoVersion(1)

	// the key changed, but the version did not, so the provider should not touch it
	userModelWithKey2SameVersion := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key2).
		WithRsaPublicKeyWoVersion(1)

	userModelWithKey2 := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key2).
		WithRsaPublicKeyWoVersion(2)

	userModelWithoutKey := model.User("u", userId.Name())

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.WriteOnlyAttributes)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: writeOnlyAttributesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_11_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// SET THE KEY WITH THE WRITE-ONLY FIELD
			{
				Config: config.FromModels(t, providerModel, userModelWithKey1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithKey1.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithKey1.ResourceReference()).
						HasNameString(userId.Name()).
						HasRsaPublicKeyEmpty().
						HasRsaPublicKeyWoVersionString("1"),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(true),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+key1Fp),
				),
			},
			// CHANGE THE KEY WITHOUT BUMPING THE VERSION - NOTHING HAPPENS
			{
				Config: config.FromModels(t, providerModel, userModelWithKey2SameVersion),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+key1Fp),
				),
			},
			// BUMP THE VERSION - THE NEW KEY IS PUSHED
			{
				Config: config.FromModels(t, providerModel, userModelWithKey2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModelWithKey2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithKey2.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithKey2.ResourceReference()).
						HasRsaPublicKeyEmpty().
						HasRsaPublicKeyWoVersionString("2"),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(true),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+key2Fp),
				),
			},
			// REMOVE BOTH FIELDS - THE KEY IS UNSET
			{
				Config: config.FromModels(t, providerModel, userModelWithoutKey),
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithoutKey.ResourceReference()).
						HasRsaPublicKeyEmpty(),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(false),
				),
			},
		},
	})
}

func TestAcc_ServiceUser_RsaPublicKeyWriteOnly(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	key1, key1Fp := random.GenerateRSAPublicKey(t)
	key2, key2Fp := random.GenerateRSAPublicKey(t)
	newKey2, newKey2Fp := random.GenerateRSAPublicKey(t)

	userModelWithBothKeys := model.ServiceUser("u", userId.Name()).
		WithRsaPublicKeyWo(key1).
		WithRsaPublicKeyWoVersion(1).
		WithRsaPublicKey2Wo(key2).
		WithRsaPublicKey2WoVersion(1)

	// only the second key is rotated
	userModelWithRotatedSecondKey := model.ServiceUser("u", userId.Name()).
		WithRsaPublicKeyWo(key1).
		WithRsaPublicKeyWoVersion(1).
		WithRsaPublicKey2Wo(newKey2).
		WithRsaPublicKey2WoVersion(2)

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.WriteOnlyAttributes)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: writeOnlyAttributesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_11_0),
		},
		CheckDestroy: CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// SET BOTH KEYS WITH THE WRITE-ONLY FIELDS
			{
				Config: config.FromModels(t, providerModel, userModelWithBothKeys),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithBothKeys.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue(userModelWithBothKeys.ResourceReference(), tfjsonpath.New("rsa_public_key_2_wo"), knownvalue.Null()),
				},
				Check: assertThat(
					t,
					resourceassert.ServiceUserResource(t, userModelWithBothKeys.ResourceReference()).
						HasNameString(userId.Name()).
						HasRsaPublicKeyEmpty().
						HasRsaPublicKey2Empty().
						HasRsaPublicKeyWoVersionString("1").
						HasRsaPublicKey2WoVersionString("1"),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(true),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+key1Fp).
						HasRsaPublicKey2Fp("SHA256:"+key2Fp),
				),
			},
			// ROTATE ONLY THE SECOND KEY
			{
				Config: config.FromModels(t, providerModel, userModelWithRotatedSecondKey),
				Check: assertThat(
					t,
					resourceassert.ServiceUserResource(t, userModelWithRotatedSecondKey.ResourceReference()).
						HasRsaPublicKeyWoVersionString("1").
						HasRsaPublicKey2WoVersionString("2"),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+key1Fp).
						HasRsaPublicKey2Fp("SHA256:"+newKey2Fp),
				),
			},
		},
	})
}

func TestAcc_User_RsaPublicKeyWriteOnly_Validations(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)

	userModelWithBothVariants := model.User("u", userId.Name()).
		WithRsaPublicKey(key).
		WithRsaPublicKeyWo(key).
		WithRsaPublicKeyWoVersion(1)

	userModelWithoutVersion := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key)

	userModelWithoutKey := model.User("u", userId.Name()).
		WithRsaPublicKeyWoVersion(1)

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.WriteOnlyAttributes)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: writeOnlyAttributesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_11_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, providerModel, userModelWithBothVariants),
				ExpectError: regexp.MustCompile(`(?s)"rsa_public_key":\s+conflicts\s+with\s+rsa_public_key_wo`),
			},
			{
				Config:      config.FromModels(t, providerModel, userModelWithoutVersion),
				ExpectError: regexp.MustCompile(`(?s)"rsa_public_key_wo":\s+all\s+of.*must\s+be\s+specified`),
			},
			{
				Config:      config.FromModels(t, providerModel, userModelWithoutKey),
				ExpectError: regexp.MustCompile(`(?s)"rsa_public_key_wo_version":\s+all\s+of.*must\s+be\s+specified`),
			},
		},
	})
}

func TestAcc_User_RsaPublicKeyWriteOnly_MigrateBetweenRegularAndWriteOnlyFields(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	key, keyFp := random.GenerateRSAPublicKey(t)
	newKey, newKeyFp := random.GenerateRSAPublicKey(t)

	userModelWithRegularField := model.User("u", userId.Name()).
		WithRsaPublicKey(key)

	userModelWithWriteOnlyField := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key).
		WithRsaPublicKeyWoVersion(1)

	userModelWithRegularFieldAndNewKey := model.User("u", userId.Name()).
		WithRsaPublicKey(newKey)

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.WriteOnlyAttributes)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: writeOnlyAttributesProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_11_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// THE KEY IS MANAGED WITH THE REGULAR FIELD AND IS SAVED IN THE STATE
			{
				Config: config.FromModels(t, providerModel, userModelWithRegularField),
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithRegularField.ResourceReference()).
						HasRsaPublicKeyString(key),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+keyFp),
				),
			},
			// AFTER SWITCHING TO THE WRITE-ONLY FIELD THE KEY IS GONE FROM THE STATE BUT STILL SET IN SNOWFLAKE
			{
				Config: config.FromModels(t, providerModel, userModelWithWriteOnlyField),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithWriteOnlyField.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithWriteOnlyField.ResourceReference()).
						HasRsaPublicKeyEmpty().
						HasRsaPublicKeyWoVersionString("1"),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(true),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+keyFp),
				),
			},
			// SWITCHING BACK TO THE REGULAR FIELD PUSHES ITS VALUE AND SAVES IT IN THE STATE AGAIN
			{
				Config: config.FromModels(t, providerModel, userModelWithRegularFieldAndNewKey),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModelWithRegularFieldAndNewKey.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithRegularFieldAndNewKey.ResourceReference()).
						HasRsaPublicKeyString(newKey),
					objectassert.UserDetails(t, userId).
						HasRsaPublicKeyFp("SHA256:"+newKeyFp),
				),
			},
		},
	})
}

func TestAcc_User_RsaPublicKeyWriteOnly_Ephemeral(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()

	ephemeralKeyConfig := `
ephemeral "tls_private_key" "key" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

locals {
  write_only_public_key = trimspace(replace(replace(ephemeral.tls_private_key.key.public_key_pem, "-----BEGIN PUBLIC KEY-----", ""), "-----END PUBLIC KEY-----", ""))
}
`

	userModel := model.User("u", userId.Name()).
		WithRsaPublicKeyWoValue(config.UnquotedWrapperVariable("local.write_only_public_key")).
		WithRsaPublicKeyWoVersion(1)

	providerModel := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.WriteOnlyAttributes)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: writeOnlyAttributesProviderFactory,
		ExternalProviders:        ExternalTlsProvider(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_11_0),
		},
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: ephemeralKeyConfig + config.FromModels(t, providerModel, userModel),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModel.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()).
						HasRsaPublicKeyEmpty().
						HasRsaPublicKeyWoVersionString("1"),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(true),
				),
			},
		},
	})
}

func TestAcc_User_RsaPublicKeyWriteOnly_MissingExperimentFlag(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	key, _ := random.GenerateRSAPublicKey(t)

	expectMissingExperiment := regexp.MustCompile("using `rsa_public_key_wo`, `rsa_public_key_2_wo`, or their version fields requires the .*WRITE_ONLY_ATTRIBUTES.* experiment to be enabled. Add it to the `experimental_features_enabled` list")

	userModel := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key).
		WithRsaPublicKeyWoVersion(1)
	serviceUserModel := model.ServiceUser("u", userId.Name()).
		WithRsaPublicKey2Wo(key).
		WithRsaPublicKey2WoVersion(1)
	legacyServiceUserModel := model.LegacyServiceUser("u", userId.Name()).
		WithRsaPublicKeyWo(key).
		WithRsaPublicKeyWoVersion(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, userModel),
				PlanOnly:    true,
				ExpectError: expectMissingExperiment,
			},
			{
				Config:      config.FromModels(t, serviceUserModel),
				PlanOnly:    true,
				ExpectError: expectMissingExperiment,
			},
			{
				Config:      config.FromModels(t, legacyServiceUserModel),
				PlanOnly:    true,
				ExpectError: expectMissingExperiment,
			},
		},
	})
}
