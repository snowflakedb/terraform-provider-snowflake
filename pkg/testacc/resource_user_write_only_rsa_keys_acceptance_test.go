//go:build account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func checkUserRsaPublicKeyFp(t *testing.T, userId sdk.AccountObjectIdentifier, expectedHash string) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		userDetails, err := testClient().User.Describe(t, userId)
		if err != nil {
			return err
		}
		if userDetails.RsaPublicKeyFp == nil {
			return fmt.Errorf("expected user %s to have RSA_PUBLIC_KEY_FP set", userId.Name())
		}
		if expected := "SHA256:" + expectedHash; userDetails.RsaPublicKeyFp.Value != expected {
			return fmt.Errorf("expected user %s to have RSA_PUBLIC_KEY_FP %s; got: %s", userId.Name(), expected, userDetails.RsaPublicKeyFp.Value)
		}
		return nil
	}
}

func checkUserRsaPublicKey2Fp(t *testing.T, userId sdk.AccountObjectIdentifier, expectedHash string) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		userDetails, err := testClient().User.Describe(t, userId)
		if err != nil {
			return err
		}
		if userDetails.RsaPublicKey2Fp == nil {
			return fmt.Errorf("expected user %s to have RSA_PUBLIC_KEY_2_FP set", userId.Name())
		}
		if expected := "SHA256:" + expectedHash; userDetails.RsaPublicKey2Fp.Value != expected {
			return fmt.Errorf("expected user %s to have RSA_PUBLIC_KEY_2_FP %s; got: %s", userId.Name(), expected, userDetails.RsaPublicKey2Fp.Value)
		}
		return nil
	}
}

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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks:   writeOnlyTerraformVersionChecks(),
		CheckDestroy:             CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// SET THE KEY WITH THE WRITE-ONLY FIELD
			{
				Config: config.FromModels(t, userModelWithKey1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithKey1.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(
						t,
						resourceassert.UserResource(t, userModelWithKey1.ResourceReference()).
							HasNameString(userId.Name()).
							HasRsaPublicKeyNotInState().
							HasRsaPublicKeyWoVersionString("1"),
						objectassert.User(t, userId).
							HasHasRsaPublicKey(true),
					),
					checkUserRsaPublicKeyFp(t, userId, key1Fp),
				),
			},
			// CHANGE THE KEY WITHOUT BUMPING THE VERSION - NOTHING HAPPENS
			{
				Config: config.FromModels(t, userModelWithKey2SameVersion),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					checkUserRsaPublicKeyFp(t, userId, key1Fp),
				),
			},
			// BUMP THE VERSION - THE NEW KEY IS PUSHED
			{
				Config: config.FromModels(t, userModelWithKey2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(userModelWithKey2.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithKey2.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(
						t,
						resourceassert.UserResource(t, userModelWithKey2.ResourceReference()).
							HasRsaPublicKeyNotInState().
							HasRsaPublicKeyWoVersionString("2"),
						objectassert.User(t, userId).
							HasHasRsaPublicKey(true),
					),
					checkUserRsaPublicKeyFp(t, userId, key2Fp),
				),
			},
			// REMOVE BOTH FIELDS - THE KEY IS UNSET
			{
				Config: config.FromModels(t, userModelWithoutKey),
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModelWithoutKey.ResourceReference()).
						HasRsaPublicKeyNotInState(),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks:   writeOnlyTerraformVersionChecks(),
		CheckDestroy:             CheckDestroy(t, resources.ServiceUser),
		Steps: []resource.TestStep{
			// SET BOTH KEYS WITH THE WRITE-ONLY FIELDS
			{
				Config: config.FromModels(t, userModelWithBothKeys),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithBothKeys.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue(userModelWithBothKeys.ResourceReference(), tfjsonpath.New("rsa_public_key_2_wo"), knownvalue.Null()),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(
						t,
						resourceassert.ServiceUserResource(t, userModelWithBothKeys.ResourceReference()).
							HasNameString(userId.Name()).
							HasRsaPublicKeyEmpty().
							HasRsaPublicKey2Empty().
							HasRsaPublicKeyWoVersionString("1").
							HasRsaPublicKey2WoVersionString("1"),
						objectassert.User(t, userId).
							HasHasRsaPublicKey(true),
					),
					checkUserRsaPublicKeyFp(t, userId, key1Fp),
					checkUserRsaPublicKey2Fp(t, userId, key2Fp),
				),
			},
			// ROTATE ONLY THE SECOND KEY
			{
				Config: config.FromModels(t, userModelWithRotatedSecondKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(
						t,
						resourceassert.ServiceUserResource(t, userModelWithRotatedSecondKey.ResourceReference()).
							HasRsaPublicKeyWoVersionString("1").
							HasRsaPublicKey2WoVersionString("2"),
					),
					checkUserRsaPublicKeyFp(t, userId, key1Fp),
					checkUserRsaPublicKey2Fp(t, userId, newKey2Fp),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks:   writeOnlyTerraformVersionChecks(),
		CheckDestroy:             CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, userModelWithBothVariants),
				ExpectError: regexp.MustCompile(`(?s)"rsa_public_key":\s+conflicts\s+with\s+rsa_public_key_wo`),
			},
			{
				Config:      config.FromModels(t, userModelWithoutVersion),
				ExpectError: regexp.MustCompile(`(?s)"rsa_public_key_wo":\s+all\s+of.*must\s+be\s+specified`),
			},
			{
				Config:      config.FromModels(t, userModelWithoutKey),
				ExpectError: regexp.MustCompile(`(?s)"rsa_public_key_wo_version":\s+all\s+of.*must\s+be\s+specified`),
			},
		},
	})
}

func TestAcc_User_RsaPublicKeyWriteOnly_MigrateFromRsaPublicKey(t *testing.T) {
	userId := testClient().Ids.RandomAccountObjectIdentifier()
	key, keyFp := random.GenerateRSAPublicKey(t)

	userModelWithRegularField := model.User("u", userId.Name()).
		WithRsaPublicKey(key)

	userModelWithWriteOnlyField := model.User("u", userId.Name()).
		WithRsaPublicKeyWo(key).
		WithRsaPublicKeyWoVersion(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks:   writeOnlyTerraformVersionChecks(),
		CheckDestroy:             CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// THE KEY IS MANAGED WITH THE REGULAR FIELD AND IS SAVED IN THE STATE
			{
				Config: config.FromModels(t, userModelWithRegularField),
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(
						t,
						resourceassert.UserResource(t, userModelWithRegularField.ResourceReference()).
							HasRsaPublicKeyString(key),
					),
					checkUserRsaPublicKeyFp(t, userId, keyFp),
				),
			},
			// AFTER SWITCHING TO THE WRITE-ONLY FIELD THE KEY IS GONE FROM THE STATE BUT STILL SET IN SNOWFLAKE
			{
				Config: config.FromModels(t, userModelWithWriteOnlyField),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModelWithWriteOnlyField.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					assertThat(
						t,
						resourceassert.UserResource(t, userModelWithWriteOnlyField.ResourceReference()).
							HasRsaPublicKeyNotInState().
							HasRsaPublicKeyWoVersionString("1"),
						objectassert.User(t, userId).
							HasHasRsaPublicKey(true),
					),
					checkUserRsaPublicKeyFp(t, userId, keyFp),
				),
			},
			// THE CONFIGURATION IS STABLE AFTERWARDS
			{
				Config: config.FromModels(t, userModelWithWriteOnlyField),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		ExternalProviders:        ExternalTlsProvider(),
		TerraformVersionChecks:   writeOnlyTerraformVersionChecks(),
		CheckDestroy:             CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: ephemeralKeyConfig + config.FromModels(t, userModel),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(userModel.ResourceReference(), tfjsonpath.New("rsa_public_key_wo"), knownvalue.Null()),
				},
				Check: assertThat(
					t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasNameString(userId.Name()).
						HasRsaPublicKeyNotInState().
						HasRsaPublicKeyWoVersionString("1"),
					objectassert.User(t, userId).
						HasHasRsaPublicKey(true),
				),
			},
			{
				Config: ephemeralKeyConfig + config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
