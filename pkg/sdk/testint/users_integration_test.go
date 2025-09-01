//go:build account_level_tests

package testint

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1645875]: test setting/unsetting policies
func TestInt_Users(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	randomPrefix := random.AlphaN(6)

	password := random.Password()
	email := random.Email()
	newValue := random.SensitiveAlphanumeric()
	warehouseId := testClientHelper().Ids.WarehouseId()
	schemaId := testClientHelper().Ids.SchemaId()
	var schemaIdObjectIdentifier sdk.ObjectIdentifier = schemaId
	// does not have to exist
	roleId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	key, hash := random.GenerateRSAPublicKey(t)
	key2, hash2 := random.GenerateRSAPublicKey(t)

	user, userCleanup := testClientHelper().User.CreateUserWithPrefix(t, randomPrefix+"_")
	t.Cleanup(userCleanup)

	user2, user2Cleanup := testClientHelper().User.CreateUserWithPrefix(t, randomPrefix)
	t.Cleanup(user2Cleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyNotEmpty(t)
	t.Cleanup(networkPolicyCleanup)

	assertParametersSet := func(userParametersAssert *objectparametersassert.UserParametersAssert) {
		assertThatObject(t, userParametersAssert.
			HasEnableUnredactedQuerySyntaxError(true).
			HasNetworkPolicy(networkPolicy.ID().Name()).
			HasPreventUnloadToInternalStages(true).
			HasAbortDetachedQuery(true).
			HasAutocommit(false).
			HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
			HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
			HasClientMemoryLimit(1024).
			HasClientMetadataRequestUseConnectionCtx(true).
			HasClientPrefetchThreads(2).
			HasClientResultChunkSize(48).
			HasClientResultColumnCaseInsensitive(true).
			HasClientSessionKeepAlive(true).
			HasClientSessionKeepAliveHeartbeatFrequency(2400).
			HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
			HasDateInputFormat("YYYY-MM-DD").
			HasDateOutputFormat("YY-MM-DD").
			HasEnableUnloadPhysicalTypeOptimization(false).
			HasErrorOnNondeterministicMerge(false).
			HasErrorOnNondeterministicUpdate(true).
			HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
			HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
			HasJdbcTreatDecimalAsInt(false).
			HasJdbcTreatTimestampNtzAsUtc(true).
			HasJdbcUseSessionTimezone(false).
			HasJsonIndent(4).
			HasLockTimeout(21222).
			HasLogLevel(sdk.LogLevelError).
			HasMultiStatementCount(0).
			HasNoorderSequenceAsDefault(false).
			HasOdbcTreatDecimalAsInt(true).
			HasQueryTag("some_tag").
			HasQuotedIdentifiersIgnoreCase(true).
			HasRowsPerResultset(2).
			HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
			HasSearchPath("$public, $current").
			HasSimulatedDataSharingConsumer("some_consumer").
			HasStatementQueuedTimeoutInSeconds(10).
			HasStatementTimeoutInSeconds(10).
			HasStrictJsonOutput(true).
			HasTimestampDayIsAlways24h(true).
			HasTimestampInputFormat("YYYY-MM-DD").
			HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
			HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
			HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
			HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
			HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
			HasTimezone("Europe/Warsaw").
			HasTimeInputFormat("HH24:MI").
			HasTimeOutputFormat("HH24:MI").
			HasTraceLevel(sdk.TraceLevelPropagate).
			HasTransactionAbortOnError(true).
			HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
			HasTwoDigitCenturyStart(1980).
			HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
			HasUseCachedResult(false).
			HasWeekOfYearPolicy(1).
			HasWeekStart(1).
			HasBoolParameterValue(sdk.UserParameterUseCachedResult, false),
		)
	}

	t.Run("create: all types of params", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToUpper(random.AlphaN(6))
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.Password()
		loginName := random.SensitiveAlphanumeric()

		opts := &sdk.CreateUserOptions{
			OrReplace: sdk.Bool(true),
			ObjectProperties: &sdk.UserObjectProperties{
				Password:    &password,
				LoginName:   &loginName,
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
			},
			SessionParameters: &sdk.SessionParameters{
				Autocommit: sdk.Bool(true),
			},
			With: sdk.Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)
		assert.Equal(t, defaultRole, userDetails.DefaultRole.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.UserFromObject(t, user).
			HasName(id.Name()).
			HasHasPassword(true).
			HasLoginName(strings.ToUpper(loginName)).
			HasDefaultRole(defaultRole),
		)
	})

	t.Run("create: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		tagValue := random.String()
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		password := random.Password()
		loginName := random.SensitiveAlphanumeric()

		opts := &sdk.CreateUserOptions{
			IfNotExists: sdk.Bool(true),
			ObjectProperties: &sdk.UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
			},
			SessionParameters: &sdk.SessionParameters{
				Autocommit: sdk.Bool(true),
			},
			With: sdk.Bool(true),
			Tags: tags,
		}
		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(loginName), userDetails.LoginName.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.UserFromObject(t, user).
			HasName(id.Name()).
			HasHasPassword(true).
			HasLoginName(strings.ToUpper(loginName)),
		)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		currentRole := testClientHelper().Context.CurrentRole(t)

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(id.Name()), userDetails.LoginName.Value)
		assert.Empty(t, userDetails.Password.Value)
		assert.Empty(t, userDetails.MiddleName.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.UserFromObject(t, user).
			HasDefaults(id.Name()).
			HasDisplayName(id.Name()).
			HasOwner(currentRole.Name()),
		)
	})

	for _, userType := range sdk.AllUserTypes {
		userType := userType
		t.Run(fmt.Sprintf("create: type %s - no options", userType), func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()

			err := client.Users.Create(ctx, id, &sdk.CreateUserOptions{
				ObjectProperties: &sdk.UserObjectProperties{
					Type: sdk.Pointer(userType),
				},
			})
			require.NoError(t, err)
			t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

			userDetails, err := client.Users.Describe(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, id.Name(), userDetails.Name.Value)
			assert.Equal(t, string(userType), userDetails.Type.Value)

			user, err := client.Users.ShowByID(ctx, id)
			require.NoError(t, err)

			assertThatObject(t, objectassert.UserFromObject(t, user).
				HasDefaults(id.Name()).
				HasType(string(userType)),
			)
		})
	}

	t.Run("create: all object properties", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		currentRole := testClientHelper().Context.CurrentRole(t)

		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			Password:              sdk.String(password),
			LoginName:             sdk.String(newValue),
			DisplayName:           sdk.String(newValue),
			FirstName:             sdk.String(newValue),
			MiddleName:            sdk.String(newValue),
			LastName:              sdk.String(newValue),
			Email:                 sdk.String(email),
			MustChangePassword:    sdk.Bool(true),
			Disable:               sdk.Bool(true),
			DaysToExpiry:          sdk.Int(5),
			MinsToUnlock:          sdk.Int(15),
			DefaultWarehouse:      sdk.Pointer(warehouseId),
			DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
			DefaultRole:           sdk.Pointer(roleId),
			DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
			MinsToBypassMFA:       sdk.Int(30),
			RSAPublicKey:          sdk.String(key),
			RSAPublicKey2:         sdk.String(key2),
			Comment:               sdk.String("some comment"),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(newValue), userDetails.LoginName.Value)
		assert.NotEmpty(t, userDetails.Password.Value)
		assert.Equal(t, newValue, userDetails.MiddleName.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasName(user.Name).
			HasCreatedOnNotEmpty().
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(newValue)).
			HasDisplayName(newValue).
			HasFirstName(newValue).
			HasLastName(newValue).
			HasEmail(email).
			HasMinsToUnlock("14").
			HasDaysToExpiryNotEmpty().
			HasComment("some comment").
			HasDisabled(true).
			HasMustChangePassword(true).
			HasSnowflakeLock(false).
			HasDefaultWarehouse(warehouseId.Name()).
			HasDefaultNamespaceId(schemaId).
			HasDefaultRole(roleId.Name()).
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasExtAuthnDuo(false).
			HasExtAuthnUid("").
			HasMinsToBypassMfa("29").
			HasOwner(currentRole.Name()).
			HasLastSuccessLoginEmpty().
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(true).
			HasHasRsaPublicKey(true),
		)
	})

	t.Run("create: all object properties - type service", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		currentRole := testClientHelper().Context.CurrentRole(t)

		// omitting FirstName, MiddleName, LastName, Password, MustChangePassword, and MinsToBypassMFA
		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			LoginName:             sdk.String(newValue),
			DisplayName:           sdk.String(newValue),
			Email:                 sdk.String(email),
			Disable:               sdk.Bool(true),
			DaysToExpiry:          sdk.Int(5),
			MinsToUnlock:          sdk.Int(15),
			DefaultWarehouse:      sdk.Pointer(warehouseId),
			DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
			DefaultRole:           sdk.Pointer(roleId),
			DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
			RSAPublicKey:          sdk.String(key),
			RSAPublicKey2:         sdk.String(key2),
			Comment:               sdk.String("some comment"),
			Type:                  sdk.Pointer(sdk.UserTypeService),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(newValue), userDetails.LoginName.Value)
		assert.Equal(t, newValue, userDetails.DisplayName.Value)
		assert.Equal(t, email, userDetails.Email.Value)
		assert.Equal(t, true, userDetails.Disabled.Value)
		assert.NotEmpty(t, userDetails.DaysToExpiry.Value)
		assert.Equal(t, 14, *userDetails.MinsToUnlock.Value)
		assert.Equal(t, warehouseId.Name(), userDetails.DefaultWarehouse.Value)
		assert.Equal(t, fmt.Sprintf("%s.%s", schemaId.DatabaseName(), schemaId.Name()), userDetails.DefaultNamespace.Value)
		assert.Equal(t, roleId.Name(), userDetails.DefaultRole.Value)
		assert.Equal(t, `["ALL"]`, userDetails.DefaultSecondaryRoles.Value)
		assert.Equal(t, "some comment", userDetails.Comment.Value)
		assert.Equal(t, string(sdk.UserTypeService), userDetails.Type.Value)

		assert.Equal(t, "", userDetails.FirstName.Value)
		assert.Equal(t, "", userDetails.MiddleName.Value)
		assert.Equal(t, "", userDetails.LastName.Value)
		assert.Equal(t, "", userDetails.Password.Value)
		assert.Equal(t, false, userDetails.MustChangePassword.Value)
		assert.Nil(t, userDetails.MinsToBypassMfa.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasName(user.Name).
			HasType(string(sdk.UserTypeService)).
			HasCreatedOnNotEmpty().
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(newValue)).
			HasDisplayName(newValue).
			HasFirstName("").
			HasLastName("").
			HasEmail(email).
			HasMinsToUnlock("14").
			HasDaysToExpiryNotEmpty().
			HasComment("some comment").
			HasDisabled(true).
			HasMustChangePassword(false).
			HasSnowflakeLock(false).
			HasDefaultWarehouse(warehouseId.Name()).
			HasDefaultNamespaceId(schemaId).
			HasDefaultRole(roleId.Name()).
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasExtAuthnDuo(false).
			HasExtAuthnUid("").
			HasMinsToBypassMfa("").
			HasOwner(currentRole.Name()).
			HasLastSuccessLoginEmpty().
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(false).
			HasHasRsaPublicKey(true),
		)
	})

	t.Run("create: all object properties - type legacy service", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		currentRole := testClientHelper().Context.CurrentRole(t)

		// omitting FirstName, MiddleName, LastName, and MinsToBypassMFA
		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			Password:              sdk.String(password),
			MustChangePassword:    sdk.Bool(true),
			LoginName:             sdk.String(newValue),
			DisplayName:           sdk.String(newValue),
			Email:                 sdk.String(email),
			Disable:               sdk.Bool(true),
			DaysToExpiry:          sdk.Int(5),
			MinsToUnlock:          sdk.Int(15),
			DefaultWarehouse:      sdk.Pointer(warehouseId),
			DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
			DefaultRole:           sdk.Pointer(roleId),
			DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
			RSAPublicKey:          sdk.String(key),
			RSAPublicKey2:         sdk.String(key2),
			Comment:               sdk.String("some comment"),
			Type:                  sdk.Pointer(sdk.UserTypeLegacyService),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), userDetails.Name.Value)
		assert.Equal(t, strings.ToUpper(newValue), userDetails.LoginName.Value)
		assert.Equal(t, newValue, userDetails.DisplayName.Value)
		assert.Equal(t, email, userDetails.Email.Value)
		assert.Equal(t, true, userDetails.Disabled.Value)
		assert.NotEmpty(t, userDetails.DaysToExpiry.Value)
		assert.Equal(t, 14, *userDetails.MinsToUnlock.Value)
		assert.Equal(t, warehouseId.Name(), userDetails.DefaultWarehouse.Value)
		assert.Equal(t, fmt.Sprintf("%s.%s", schemaId.DatabaseName(), schemaId.Name()), userDetails.DefaultNamespace.Value)
		assert.Equal(t, roleId.Name(), userDetails.DefaultRole.Value)
		assert.Equal(t, `["ALL"]`, userDetails.DefaultSecondaryRoles.Value)
		assert.Equal(t, "some comment", userDetails.Comment.Value)
		assert.Equal(t, string(sdk.UserTypeLegacyService), userDetails.Type.Value)
		assert.NotEmpty(t, userDetails.Password.Value)
		assert.Equal(t, true, userDetails.MustChangePassword.Value)

		assert.Equal(t, "", userDetails.FirstName.Value)
		assert.Equal(t, "", userDetails.MiddleName.Value)
		assert.Equal(t, "", userDetails.LastName.Value)
		assert.Nil(t, userDetails.MinsToBypassMfa.Value)

		user, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasName(user.Name).
			HasType(string(sdk.UserTypeLegacyService)).
			HasCreatedOnNotEmpty().
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(newValue)).
			HasDisplayName(newValue).
			HasFirstName("").
			HasLastName("").
			HasEmail(email).
			HasMinsToUnlock("14").
			HasDaysToExpiryNotEmpty().
			HasComment("some comment").
			HasDisabled(true).
			HasMustChangePassword(true).
			HasSnowflakeLock(false).
			HasDefaultWarehouse(warehouseId.Name()).
			HasDefaultNamespaceId(schemaId).
			HasDefaultRole(roleId.Name()).
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasExtAuthnDuo(false).
			HasExtAuthnUid("").
			HasMinsToBypassMfa("").
			HasOwner(currentRole.Name()).
			HasLastSuccessLoginEmpty().
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(true).
			HasHasRsaPublicKey(true),
		)
	})

	incorrectObjectPropertiesForServiceType := []struct {
		property             string
		userObjectProperties *sdk.UserObjectProperties
	}{
		{property: "MINS_TO_BYPASS_MFA", userObjectProperties: &sdk.UserObjectProperties{MinsToBypassMFA: sdk.Int(30)}},
		{property: "MUST_CHANGE_PASSWORD", userObjectProperties: &sdk.UserObjectProperties{MustChangePassword: sdk.Bool(true)}},
		{property: "FIRST_NAME", userObjectProperties: &sdk.UserObjectProperties{FirstName: sdk.String(newValue)}},
		{property: "MIDDLE_NAME", userObjectProperties: &sdk.UserObjectProperties{MiddleName: sdk.String(newValue)}},
		{property: "LAST_NAME", userObjectProperties: &sdk.UserObjectProperties{LastName: sdk.String(newValue)}},
		{property: "PASSWORD", userObjectProperties: &sdk.UserObjectProperties{Password: sdk.String(password)}},
	}

	for _, tt := range incorrectObjectPropertiesForServiceType {
		tt := tt
		t.Run(fmt.Sprintf("create: incorrect object property %s - type service", tt.property), func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()

			tt.userObjectProperties.Type = sdk.Pointer(sdk.UserTypeService)
			createOpts := &sdk.CreateUserOptions{ObjectProperties: tt.userObjectProperties}

			err := client.Users.Create(ctx, id, createOpts)
			require.ErrorContains(t, err, fmt.Sprintf("Cannot set %s on users with TYPE=SERVICE.", tt.property))
		})
	}

	incorrectObjectPropertiesForLegacyServiceType := []struct {
		property             string
		userObjectProperties *sdk.UserObjectProperties
	}{
		{property: "MINS_TO_BYPASS_MFA", userObjectProperties: &sdk.UserObjectProperties{MinsToBypassMFA: sdk.Int(30)}},
		{property: "FIRST_NAME", userObjectProperties: &sdk.UserObjectProperties{FirstName: sdk.String(newValue)}},
		{property: "MIDDLE_NAME", userObjectProperties: &sdk.UserObjectProperties{MiddleName: sdk.String(newValue)}},
		{property: "LAST_NAME", userObjectProperties: &sdk.UserObjectProperties{LastName: sdk.String(newValue)}},
	}

	for _, tt := range incorrectObjectPropertiesForLegacyServiceType {
		tt := tt
		t.Run(fmt.Sprintf("create: incorrect object property %s - type legacy service", tt.property), func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()

			tt.userObjectProperties.Type = sdk.Pointer(sdk.UserTypeLegacyService)
			createOpts := &sdk.CreateUserOptions{ObjectProperties: tt.userObjectProperties}

			err := client.Users.Create(ctx, id, createOpts)
			require.ErrorContains(t, err, fmt.Sprintf("Cannot set %s on users with TYPE=LEGACY_SERVICE.", tt.property))
		})
	}

	t.Run("create: set mins to bypass mfa to negative manually", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			MinsToBypassMFA: sdk.Int(-100),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Nil(t, userDetails.MinsToBypassMfa.Value)
	})

	t.Run("create: set mins to bypass mfa to zero manually", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			MinsToBypassMFA: sdk.Int(0),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Nil(t, userDetails.MinsToBypassMfa.Value)
	})

	t.Run("create: set mins to bypass mfa to one manually", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			MinsToBypassMFA: sdk.Int(1),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 0, *userDetails.MinsToBypassMfa.Value)
	})

	t.Run("create and alter: problems with public key fingerprints", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		createOpts := &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			RSAPublicKey:   sdk.String(key),
			RSAPublicKeyFp: sdk.String(hash),
		}}

		err := client.Users.Create(ctx, id, createOpts)
		require.ErrorContains(t, err, "invalid property 'RSA_PUBLIC_KEY_FP' for 'USER'")

		createOpts = &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
			RSAPublicKey2:   sdk.String(key),
			RSAPublicKey2Fp: sdk.String(hash),
		}}

		err = client.Users.Create(ctx, id, createOpts)
		require.ErrorContains(t, err, "invalid property 'RSA_PUBLIC_KEY_2_FP' for 'USER'")

		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		alterOpts := &sdk.AlterUserOptions{Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					RSAPublicKey:   sdk.String(key),
					RSAPublicKeyFp: sdk.String(hash),
				},
			},
		}}

		err = client.Users.Alter(ctx, user.ID(), alterOpts)
		require.ErrorContains(t, err, "invalid property 'RSA_PUBLIC_KEY_FP' for 'USER'")

		alterOpts = &sdk.AlterUserOptions{Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					RSAPublicKey2:   sdk.String(key2),
					RSAPublicKey2Fp: sdk.String(hash2),
				},
			},
		}}

		err = client.Users.Alter(ctx, user.ID(), alterOpts)
		require.ErrorContains(t, err, "invalid property 'RSA_PUBLIC_KEY_2_FP' for 'USER'")
	})

	t.Run("create: default role with hyphen", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToUpper(random.AlphaN(4) + "-" + random.AlphaN(4))

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		assertThatObject(t, objectassert.User(t, id).
			HasDefaultRole(defaultRole),
		)
	})

	t.Run("create: default role in lowercase", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		defaultRole := strings.ToLower(random.AlphaN(6))

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				DefaultRole: sdk.Pointer(sdk.NewAccountObjectIdentifier(defaultRole)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		assertThatObject(t, objectassert.User(t, id).
			HasDefaultRole(defaultRole),
		)
	})

	t.Run("create: client memory limit set to zero", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		opts := &sdk.CreateUserOptions{
			SessionParameters: &sdk.SessionParameters{
				ClientMemoryLimit: sdk.Int(0),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))
	})

	t.Run("create: other params with hyphen and mixed cases", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		randomWithHyphenAndMixedCase := strings.ToUpper(random.AlphaN(4)) + "-" + strings.ToLower(random.SensitiveAlphanumeric())
		var namespaceId sdk.ObjectIdentifier = sdk.NewDatabaseObjectIdentifier(randomWithHyphenAndMixedCase, randomWithHyphenAndMixedCase)

		opts := &sdk.CreateUserOptions{
			ObjectProperties: &sdk.UserObjectProperties{
				LoginName:        sdk.String(randomWithHyphenAndMixedCase),
				DisplayName:      sdk.String(randomWithHyphenAndMixedCase),
				FirstName:        sdk.String(randomWithHyphenAndMixedCase),
				MiddleName:       sdk.String(randomWithHyphenAndMixedCase),
				LastName:         sdk.String(randomWithHyphenAndMixedCase),
				DefaultWarehouse: sdk.Pointer(sdk.NewAccountObjectIdentifier(randomWithHyphenAndMixedCase)),
				DefaultNamespace: sdk.Pointer(namespaceId),
				DefaultRole:      sdk.Pointer(sdk.NewAccountObjectIdentifier(randomWithHyphenAndMixedCase)),
			},
		}

		err := client.Users.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		assertThatObject(t, objectassert.User(t, id).
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(randomWithHyphenAndMixedCase)).
			HasDisplayName(randomWithHyphenAndMixedCase).
			HasFirstName(randomWithHyphenAndMixedCase).
			HasLastName(randomWithHyphenAndMixedCase).
			HasDefaultWarehouse(randomWithHyphenAndMixedCase).
			HasDefaultNamespace(randomWithHyphenAndMixedCase+"."+randomWithHyphenAndMixedCase).
			HasDefaultRole(randomWithHyphenAndMixedCase),
		)

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.MiddleName.Value)
		// login name is always case-insensitive
		assert.Equal(t, strings.ToUpper(randomWithHyphenAndMixedCase), userDetails.LoginName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DisplayName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.FirstName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.LastName.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DefaultWarehouse.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase+"."+randomWithHyphenAndMixedCase, userDetails.DefaultNamespace.Value)
		assert.Equal(t, randomWithHyphenAndMixedCase, userDetails.DefaultRole.Value)
	})

	for _, userType := range sdk.AllUserTypes {
		userType := userType
		t.Run(fmt.Sprintf("create: with all parameters set - type %s", userType), func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()

			opts := &sdk.CreateUserOptions{
				ObjectProperties: &sdk.UserObjectProperties{
					Type: sdk.Pointer(userType),
				},
				SessionParameters: &sdk.SessionParameters{
					AbortDetachedQuery:                       sdk.Bool(true),
					Autocommit:                               sdk.Bool(false),
					BinaryInputFormat:                        sdk.Pointer(sdk.BinaryInputFormatUTF8),
					BinaryOutputFormat:                       sdk.Pointer(sdk.BinaryOutputFormatBase64),
					ClientMemoryLimit:                        sdk.Int(1024),
					ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
					ClientPrefetchThreads:                    sdk.Int(2),
					ClientResultChunkSize:                    sdk.Int(48),
					ClientResultColumnCaseInsensitive:        sdk.Bool(true),
					ClientSessionKeepAlive:                   sdk.Bool(true),
					ClientSessionKeepAliveHeartbeatFrequency: sdk.Int(2400),
					ClientTimestampTypeMapping:               sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
					DateInputFormat:                          sdk.String("YYYY-MM-DD"),
					DateOutputFormat:                         sdk.String("YY-MM-DD"),
					EnableUnloadPhysicalTypeOptimization:     sdk.Bool(false),
					ErrorOnNondeterministicMerge:             sdk.Bool(false),
					ErrorOnNondeterministicUpdate:            sdk.Bool(true),
					GeographyOutputFormat:                    sdk.Pointer(sdk.GeographyOutputFormatWKB),
					GeometryOutputFormat:                     sdk.Pointer(sdk.GeometryOutputFormatWKB),
					JdbcTreatDecimalAsInt:                    sdk.Bool(false),
					JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
					JdbcUseSessionTimezone:                   sdk.Bool(false),
					JsonIndent:                               sdk.Int(4),
					LockTimeout:                              sdk.Int(21222),
					LogLevel:                                 sdk.Pointer(sdk.LogLevelError),
					MultiStatementCount:                      sdk.Int(0),
					NoorderSequenceAsDefault:                 sdk.Bool(false),
					OdbcTreatDecimalAsInt:                    sdk.Bool(true),
					QueryTag:                                 sdk.String("some_tag"),
					QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
					RowsPerResultset:                         sdk.Int(2),
					S3StageVpceDnsName:                       sdk.String("vpce-id.s3.region.vpce.amazonaws.com"),
					SearchPath:                               sdk.String("$public, $current"),
					SimulatedDataSharingConsumer:             sdk.String("some_consumer"),
					StatementQueuedTimeoutInSeconds:          sdk.Int(10),
					StatementTimeoutInSeconds:                sdk.Int(10),
					StrictJsonOutput:                         sdk.Bool(true),
					TimestampDayIsAlways24h:                  sdk.Bool(true),
					TimestampInputFormat:                     sdk.String("YYYY-MM-DD"),
					TimestampLTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
					TimestampNTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
					TimestampOutputFormat:                    sdk.String("YYYY-MM-DD HH24:MI:SS"),
					TimestampTypeMapping:                     sdk.Pointer(sdk.TimestampTypeMappingLtz),
					TimestampTZOutputFormat:                  sdk.String("YYYY-MM-DD HH24:MI:SS"),
					Timezone:                                 sdk.String("Europe/Warsaw"),
					TimeInputFormat:                          sdk.String("HH24:MI"),
					TimeOutputFormat:                         sdk.String("HH24:MI"),
					TraceLevel:                               sdk.Pointer(sdk.TraceLevelPropagate),
					TransactionAbortOnError:                  sdk.Bool(true),
					TransactionDefaultIsolationLevel:         sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
					TwoDigitCenturyStart:                     sdk.Int(1980),
					UnsupportedDDLAction:                     sdk.Pointer(sdk.UnsupportedDDLActionFail),
					UseCachedResult:                          sdk.Bool(false),
					WeekOfYearPolicy:                         sdk.Int(1),
					WeekStart:                                sdk.Int(1),
				},
				ObjectParameters: &sdk.UserObjectParameters{
					EnableUnredactedQuerySyntaxError: sdk.Bool(true),
					NetworkPolicy:                    sdk.Pointer(networkPolicy.ID()),
					PreventUnloadToInternalStages:    sdk.Bool(true),
				},
			}

			err := client.Users.Create(ctx, id, opts)
			require.NoError(t, err)
			t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

			assertParametersSet(objectparametersassert.UserParameters(t, id))

			// check that ShowParameters works too
			parameters, err := client.Users.ShowParameters(ctx, id)
			require.NoError(t, err)
			assertParametersSet(objectparametersassert.UserParametersPrefetched(t, id, parameters))
		})
	}

	t.Run("create: with all parameters default", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		assertThatObject(t, objectparametersassert.UserParameters(t, id).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)

		// check that ShowParameters works too
		parameters, err := client.Users.ShowParameters(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectparametersassert.UserParametersPrefetched(t, id, parameters).
			HasAllDefaults().
			HasAllDefaultsExplicit(),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		alterOptions := &sdk.AlterUserOptions{
			NewName: newID,
		}
		err := client.Users.Alter(ctx, user.ID(), alterOptions)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, newID))

		result, err := client.Users.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), result.Name.Value)
	})

	t.Run("alter: set and unset object properties", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		currentRole := testClientHelper().Context.CurrentRole(t)

		assertThatObject(t, objectassert.UserFromObject(t, user).
			HasDefaults(user.Name).
			HasDisplayName(user.Name).
			HasOwner(currentRole.Name()),
		)

		alterOpts := &sdk.AlterUserOptions{Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					Password:              sdk.String(password),
					LoginName:             sdk.String(newValue),
					DisplayName:           sdk.String(newValue),
					FirstName:             sdk.String(newValue),
					MiddleName:            sdk.String(newValue),
					LastName:              sdk.String(newValue),
					Email:                 sdk.String(email),
					MustChangePassword:    sdk.Bool(true),
					Disable:               sdk.Bool(true),
					DaysToExpiry:          sdk.Int(5),
					MinsToUnlock:          sdk.Int(15),
					DefaultWarehouse:      sdk.Pointer(warehouseId),
					DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
					DefaultRole:           sdk.Pointer(roleId),
					DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
					MinsToBypassMFA:       sdk.Int(30),
					RSAPublicKey:          sdk.String(key),
					RSAPublicKey2:         sdk.String(key2),
					Comment:               sdk.String("some comment"),
				},
				DisableMfa: sdk.Bool(true),
			},
		}}

		err := client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasName(user.Name).
			HasCreatedOnNotEmpty().
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(newValue)).
			HasDisplayName(newValue).
			HasFirstName(newValue).
			HasLastName(newValue).
			HasEmail(email).
			HasMinsToUnlock("14").
			HasDaysToExpiryNotEmpty().
			HasComment("some comment").
			HasDisabled(true).
			HasMustChangePassword(true).
			HasSnowflakeLock(false).
			HasDefaultWarehouse(warehouseId.Name()).
			HasDefaultNamespaceId(schemaId).
			HasDefaultRole(roleId.Name()).
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasExtAuthnDuo(false).
			HasExtAuthnUid("").
			HasMinsToBypassMfa("29").
			HasOwner(currentRole.Name()).
			HasLastSuccessLoginEmpty().
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(true).
			HasHasRsaPublicKey(true),
		)

		alterOpts = &sdk.AlterUserOptions{Unset: &sdk.UserUnset{
			ObjectProperties: &sdk.UserObjectPropertiesUnset{
				Password:              sdk.Bool(true),
				LoginName:             sdk.Bool(true),
				DisplayName:           sdk.Bool(true),
				FirstName:             sdk.Bool(true),
				MiddleName:            sdk.Bool(true),
				LastName:              sdk.Bool(true),
				Email:                 sdk.Bool(true),
				MustChangePassword:    sdk.Bool(true),
				Disable:               sdk.Bool(true),
				DaysToExpiry:          sdk.Bool(true),
				MinsToUnlock:          sdk.Bool(true),
				DefaultWarehouse:      sdk.Bool(true),
				DefaultNamespace:      sdk.Bool(true),
				DefaultRole:           sdk.Bool(true),
				DefaultSecondaryRoles: sdk.Bool(true),
				MinsToBypassMFA:       sdk.Bool(true),
				DisableMfa:            sdk.Bool(true),
				RSAPublicKey:          sdk.Bool(true),
				RSAPublicKey2:         sdk.Bool(true),
				Comment:               sdk.Bool(true),
			},
		}}

		err = client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasDefaults(user.Name).
			HasDisplayName("").
			HasOwner(currentRole.Name()),
		)
	})

	t.Run("alter: set and unset object properties - type service", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateServiceUser(t)
		t.Cleanup(userCleanup)

		currentRole := testClientHelper().Context.CurrentRole(t)

		assertThatObject(t, objectassert.UserFromObject(t, user).
			HasDefaults(user.Name).
			HasDisplayName(user.Name).
			HasOwner(currentRole.Name()),
		)

		// omitting FirstName, MiddleName, LastName, Password, MustChangePassword, MinsToBypassMFA, and DisableMfa
		alterOpts := &sdk.AlterUserOptions{Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					LoginName:             sdk.String(newValue),
					DisplayName:           sdk.String(newValue),
					Email:                 sdk.String(email),
					Disable:               sdk.Bool(true),
					DaysToExpiry:          sdk.Int(5),
					MinsToUnlock:          sdk.Int(15),
					DefaultWarehouse:      sdk.Pointer(warehouseId),
					DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
					DefaultRole:           sdk.Pointer(roleId),
					DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
					RSAPublicKey:          sdk.String(key),
					RSAPublicKey2:         sdk.String(key2),
					Comment:               sdk.String("some comment"),
				},
			},
		}}

		err := client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasName(user.Name).
			HasCreatedOnNotEmpty().
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(newValue)).
			HasDisplayName(newValue).
			HasFirstName("").
			HasLastName("").
			HasEmail(email).
			HasMinsToUnlock("14").
			HasDaysToExpiryNotEmpty().
			HasComment("some comment").
			HasDisabled(true).
			HasMustChangePassword(false).
			HasSnowflakeLock(false).
			HasDefaultWarehouse(warehouseId.Name()).
			HasDefaultNamespaceId(schemaId).
			HasDefaultRole(roleId.Name()).
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasExtAuthnDuo(false).
			HasExtAuthnUid("").
			HasMinsToBypassMfa("").
			HasOwner(currentRole.Name()).
			HasLastSuccessLoginEmpty().
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(false).
			HasHasRsaPublicKey(true),
		)

		alterOpts = &sdk.AlterUserOptions{Unset: &sdk.UserUnset{
			ObjectProperties: &sdk.UserObjectPropertiesUnset{
				LoginName:             sdk.Bool(true),
				DisplayName:           sdk.Bool(true),
				Email:                 sdk.Bool(true),
				Disable:               sdk.Bool(true),
				DaysToExpiry:          sdk.Bool(true),
				MinsToUnlock:          sdk.Bool(true),
				DefaultWarehouse:      sdk.Bool(true),
				DefaultNamespace:      sdk.Bool(true),
				DefaultRole:           sdk.Bool(true),
				DefaultSecondaryRoles: sdk.Bool(true),
				RSAPublicKey:          sdk.Bool(true),
				RSAPublicKey2:         sdk.Bool(true),
				Comment:               sdk.Bool(true),
			},
		}}

		err = client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasDefaults(user.Name).
			HasDisplayName("").
			HasOwner(currentRole.Name()),
		)
	})

	t.Run("alter: set and unset object properties - type legacy service", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateLegacyServiceUser(t)
		t.Cleanup(userCleanup)

		currentRole := testClientHelper().Context.CurrentRole(t)

		assertThatObject(t, objectassert.UserFromObject(t, user).
			HasDefaults(user.Name).
			HasDisplayName(user.Name).
			HasOwner(currentRole.Name()),
		)

		// omitting FirstName, MiddleName, LastName, MinsToBypassMFA, and DisableMfa
		alterOpts := &sdk.AlterUserOptions{Set: &sdk.UserSet{
			ObjectProperties: &sdk.UserAlterObjectProperties{
				UserObjectProperties: sdk.UserObjectProperties{
					Password:              sdk.String(password),
					MustChangePassword:    sdk.Bool(true),
					LoginName:             sdk.String(newValue),
					DisplayName:           sdk.String(newValue),
					Email:                 sdk.String(email),
					Disable:               sdk.Bool(true),
					DaysToExpiry:          sdk.Int(5),
					MinsToUnlock:          sdk.Int(15),
					DefaultWarehouse:      sdk.Pointer(warehouseId),
					DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
					DefaultRole:           sdk.Pointer(roleId),
					DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
					RSAPublicKey:          sdk.String(key),
					RSAPublicKey2:         sdk.String(key2),
					Comment:               sdk.String("some comment"),
				},
			},
		}}

		err := client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasName(user.Name).
			HasCreatedOnNotEmpty().
			// login name is always case-insensitive
			HasLoginName(strings.ToUpper(newValue)).
			HasDisplayName(newValue).
			HasFirstName("").
			HasLastName("").
			HasEmail(email).
			HasMinsToUnlock("14").
			HasDaysToExpiryNotEmpty().
			HasComment("some comment").
			HasDisabled(true).
			HasMustChangePassword(true).
			HasSnowflakeLock(false).
			HasDefaultWarehouse(warehouseId.Name()).
			HasDefaultNamespaceId(schemaId).
			HasDefaultRole(roleId.Name()).
			HasDefaultSecondaryRoles(`["ALL"]`).
			HasExtAuthnDuo(false).
			HasExtAuthnUid("").
			HasMinsToBypassMfa("").
			HasOwner(currentRole.Name()).
			HasLastSuccessLoginEmpty().
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(true).
			HasHasRsaPublicKey(true),
		)

		alterOpts = &sdk.AlterUserOptions{Unset: &sdk.UserUnset{
			ObjectProperties: &sdk.UserObjectPropertiesUnset{
				Password:              sdk.Bool(true),
				MustChangePassword:    sdk.Bool(true),
				LoginName:             sdk.Bool(true),
				DisplayName:           sdk.Bool(true),
				Email:                 sdk.Bool(true),
				Disable:               sdk.Bool(true),
				DaysToExpiry:          sdk.Bool(true),
				MinsToUnlock:          sdk.Bool(true),
				DefaultWarehouse:      sdk.Bool(true),
				DefaultNamespace:      sdk.Bool(true),
				DefaultRole:           sdk.Bool(true),
				DefaultSecondaryRoles: sdk.Bool(true),
				RSAPublicKey:          sdk.Bool(true),
				RSAPublicKey2:         sdk.Bool(true),
				Comment:               sdk.Bool(true),
			},
		}}

		err = client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObject(t, objectassert.User(t, user.ID()).
			HasDefaults(user.Name).
			HasDisplayName("").
			HasOwner(currentRole.Name()),
		)
	})

	incorrectAlterForServiceType := []struct {
		property           string
		alterSet           *sdk.UserAlterObjectProperties
		alterUnset         *sdk.UserObjectPropertiesUnset
		expectNoUnsetError bool
	}{
		{
			property:   "MINS_TO_BYPASS_MFA",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{MinsToBypassMFA: sdk.Int(30)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{MinsToBypassMFA: sdk.Bool(true)},
			// unset for MINS_TO_BYPASS_MFA is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "MUST_CHANGE_PASSWORD",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{MustChangePassword: sdk.Bool(true)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{MustChangePassword: sdk.Bool(true)},
		},
		{
			property:   "FIRST_NAME",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{FirstName: sdk.String(newValue)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{FirstName: sdk.Bool(true)},
			// unset for FIRST_NAME is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "MIDDLE_NAME",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{MiddleName: sdk.String(newValue)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{MiddleName: sdk.Bool(true)},
			// unset for MIDDLE_NAME is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "LAST_NAME",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{LastName: sdk.String(newValue)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{LastName: sdk.Bool(true)},
			// unset for LAST_NAME is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "PASSWORD",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{Password: sdk.String(password)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{Password: sdk.Bool(true)},
			// unset for PASSWORD is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "DISABLE_MFA",
			alterSet:   &sdk.UserAlterObjectProperties{DisableMfa: sdk.Bool(true)},
			alterUnset: &sdk.UserObjectPropertiesUnset{DisableMfa: sdk.Bool(true)},
		},
	}

	for _, tt := range incorrectAlterForServiceType {
		tt := tt
		t.Run(fmt.Sprintf("alter: set and unset incorrect object property %s - type service", tt.property), func(t *testing.T) {
			serviceUser, serviceUserCleanup := testClientHelper().User.CreateServiceUser(t)
			t.Cleanup(serviceUserCleanup)

			alterSet := &sdk.AlterUserOptions{Set: &sdk.UserSet{
				ObjectProperties: tt.alterSet,
			}}

			err := client.Users.Alter(ctx, serviceUser.ID(), alterSet)
			require.ErrorContains(t, err, fmt.Sprintf("Cannot set %s on users with TYPE=SERVICE.", tt.property))

			alterUnset := &sdk.AlterUserOptions{Unset: &sdk.UserUnset{
				ObjectProperties: tt.alterUnset,
			}}

			err = client.Users.Alter(ctx, serviceUser.ID(), alterUnset)
			if tt.expectNoUnsetError {
				require.Nil(t, err)
			} else {
				require.ErrorContains(t, err, fmt.Sprintf("Cannot set %s on users with TYPE=SERVICE.", tt.property))
			}
		})
	}

	incorrectAlterForLegacyServiceType := []struct {
		property           string
		alterSet           *sdk.UserAlterObjectProperties
		alterUnset         *sdk.UserObjectPropertiesUnset
		expectNoUnsetError bool
	}{
		{
			property:   "MINS_TO_BYPASS_MFA",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{MinsToBypassMFA: sdk.Int(30)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{MinsToBypassMFA: sdk.Bool(true)},
			// unset for MINS_TO_BYPASS_MFA is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "FIRST_NAME",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{FirstName: sdk.String(newValue)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{FirstName: sdk.Bool(true)},
			// unset for FIRST_NAME is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "MIDDLE_NAME",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{MiddleName: sdk.String(newValue)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{MiddleName: sdk.Bool(true)},
			// unset for MIDDLE_NAME is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "LAST_NAME",
			alterSet:   &sdk.UserAlterObjectProperties{UserObjectProperties: sdk.UserObjectProperties{LastName: sdk.String(newValue)}},
			alterUnset: &sdk.UserObjectPropertiesUnset{LastName: sdk.Bool(true)},
			// unset for LAST_NAME is not returning an error from Snowflake
			expectNoUnsetError: true,
		},
		{
			property:   "DISABLE_MFA",
			alterSet:   &sdk.UserAlterObjectProperties{DisableMfa: sdk.Bool(true)},
			alterUnset: &sdk.UserObjectPropertiesUnset{DisableMfa: sdk.Bool(true)},
		},
	}

	for _, tt := range incorrectAlterForLegacyServiceType {
		tt := tt
		t.Run(fmt.Sprintf("alter: set and unset incorrect object property %s - type legacy service", tt.property), func(t *testing.T) {
			legacyServiceUser, legacyServiceUserCleanup := testClientHelper().User.CreateLegacyServiceUser(t)
			t.Cleanup(legacyServiceUserCleanup)

			alterSet := &sdk.AlterUserOptions{Set: &sdk.UserSet{
				ObjectProperties: tt.alterSet,
			}}

			err := client.Users.Alter(ctx, legacyServiceUser.ID(), alterSet)
			require.ErrorContains(t, err, fmt.Sprintf("Cannot set %s on users with TYPE=LEGACY_SERVICE.", tt.property))

			alterUnset := &sdk.AlterUserOptions{Unset: &sdk.UserUnset{
				ObjectProperties: tt.alterUnset,
			}}

			err = client.Users.Alter(ctx, legacyServiceUser.ID(), alterUnset)
			if tt.expectNoUnsetError {
				require.Nil(t, err)
			} else {
				require.ErrorContains(t, err, fmt.Sprintf("Cannot set %s on users with TYPE=LEGACY_SERVICE.", tt.property))
			}
		})
	}

	t.Run("set and unset authentication policy", func(t *testing.T) {
		authenticationPolicyTest, authenticationPolicyCleanup := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(authenticationPolicyCleanup)

		err := client.Users.Alter(ctx, user.ID(), &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				AuthenticationPolicy: sdk.Pointer(authenticationPolicyTest.ID()),
			},
		})
		require.NoError(t, err)

		policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, user.ID(), sdk.PolicyEntityDomainUser)
		require.NoError(t, err)

		_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
			return reference.PolicyKind == sdk.PolicyKindAuthenticationPolicy
		})
		require.NoError(t, err)

		err = client.Users.Alter(ctx, user.ID(), &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				AuthenticationPolicy: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		policies, err = testClientHelper().PolicyReferences.GetPolicyReferences(t, user.ID(), sdk.PolicyEntityDomainUser)
		require.NoError(t, err)

		_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
			return reference.PolicyKind == sdk.PolicyKindAuthenticationPolicy
		})
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	for _, userType := range sdk.AllUserTypes {
		userType := userType
		t.Run(fmt.Sprintf("alter: set and unset parameters - type %s", userType), func(t *testing.T) {
			id := testClientHelper().Ids.RandomAccountObjectIdentifier()

			err := client.Users.Create(ctx, id, &sdk.CreateUserOptions{
				ObjectProperties: &sdk.UserObjectProperties{
					Type: sdk.Pointer(userType),
				},
			})
			require.NoError(t, err)
			t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

			alterOpts := &sdk.AlterUserOptions{
				Set: &sdk.UserSet{
					SessionParameters: &sdk.SessionParameters{
						AbortDetachedQuery:                       sdk.Bool(true),
						Autocommit:                               sdk.Bool(false),
						BinaryInputFormat:                        sdk.Pointer(sdk.BinaryInputFormatUTF8),
						BinaryOutputFormat:                       sdk.Pointer(sdk.BinaryOutputFormatBase64),
						ClientMemoryLimit:                        sdk.Int(1024),
						ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
						ClientPrefetchThreads:                    sdk.Int(2),
						ClientResultChunkSize:                    sdk.Int(48),
						ClientResultColumnCaseInsensitive:        sdk.Bool(true),
						ClientSessionKeepAlive:                   sdk.Bool(true),
						ClientSessionKeepAliveHeartbeatFrequency: sdk.Int(2400),
						ClientTimestampTypeMapping:               sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
						DateInputFormat:                          sdk.String("YYYY-MM-DD"),
						DateOutputFormat:                         sdk.String("YY-MM-DD"),
						EnableUnloadPhysicalTypeOptimization:     sdk.Bool(false),
						ErrorOnNondeterministicMerge:             sdk.Bool(false),
						ErrorOnNondeterministicUpdate:            sdk.Bool(true),
						GeographyOutputFormat:                    sdk.Pointer(sdk.GeographyOutputFormatWKB),
						GeometryOutputFormat:                     sdk.Pointer(sdk.GeometryOutputFormatWKB),
						JdbcTreatDecimalAsInt:                    sdk.Bool(false),
						JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
						JdbcUseSessionTimezone:                   sdk.Bool(false),
						JsonIndent:                               sdk.Int(4),
						LockTimeout:                              sdk.Int(21222),
						LogLevel:                                 sdk.Pointer(sdk.LogLevelError),
						MultiStatementCount:                      sdk.Int(0),
						NoorderSequenceAsDefault:                 sdk.Bool(false),
						OdbcTreatDecimalAsInt:                    sdk.Bool(true),
						QueryTag:                                 sdk.String("some_tag"),
						QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
						RowsPerResultset:                         sdk.Int(2),
						S3StageVpceDnsName:                       sdk.String("vpce-id.s3.region.vpce.amazonaws.com"),
						SearchPath:                               sdk.String("$public, $current"),
						SimulatedDataSharingConsumer:             sdk.String("some_consumer"),
						StatementQueuedTimeoutInSeconds:          sdk.Int(10),
						StatementTimeoutInSeconds:                sdk.Int(10),
						StrictJsonOutput:                         sdk.Bool(true),
						TimestampDayIsAlways24h:                  sdk.Bool(true),
						TimestampInputFormat:                     sdk.String("YYYY-MM-DD"),
						TimestampLTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
						TimestampNTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
						TimestampOutputFormat:                    sdk.String("YYYY-MM-DD HH24:MI:SS"),
						TimestampTypeMapping:                     sdk.Pointer(sdk.TimestampTypeMappingLtz),
						TimestampTZOutputFormat:                  sdk.String("YYYY-MM-DD HH24:MI:SS"),
						Timezone:                                 sdk.String("Europe/Warsaw"),
						TimeInputFormat:                          sdk.String("HH24:MI"),
						TimeOutputFormat:                         sdk.String("HH24:MI"),
						TraceLevel:                               sdk.Pointer(sdk.TraceLevelPropagate),
						TransactionAbortOnError:                  sdk.Bool(true),
						TransactionDefaultIsolationLevel:         sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
						TwoDigitCenturyStart:                     sdk.Int(1980),
						UnsupportedDDLAction:                     sdk.Pointer(sdk.UnsupportedDDLActionFail),
						UseCachedResult:                          sdk.Bool(false),
						WeekOfYearPolicy:                         sdk.Int(1),
						WeekStart:                                sdk.Int(1),
					},
					ObjectParameters: &sdk.UserObjectParameters{
						EnableUnredactedQuerySyntaxError: sdk.Bool(true),
						NetworkPolicy:                    sdk.Pointer(networkPolicy.ID()),
						PreventUnloadToInternalStages:    sdk.Bool(true),
					},
				},
			}

			err = client.Users.Alter(ctx, id, alterOpts)
			require.NoError(t, err)

			assertParametersSet(objectparametersassert.UserParameters(t, id))

			// check that ShowParameters works too
			parameters, err := client.Users.ShowParameters(ctx, id)
			require.NoError(t, err)
			assertParametersSet(objectparametersassert.UserParametersPrefetched(t, id, parameters))

			alterOpts = &sdk.AlterUserOptions{
				Unset: &sdk.UserUnset{
					SessionParameters: &sdk.SessionParametersUnset{
						AbortDetachedQuery:                       sdk.Bool(true),
						Autocommit:                               sdk.Bool(true),
						BinaryInputFormat:                        sdk.Bool(true),
						BinaryOutputFormat:                       sdk.Bool(true),
						ClientMemoryLimit:                        sdk.Bool(true),
						ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
						ClientPrefetchThreads:                    sdk.Bool(true),
						ClientResultChunkSize:                    sdk.Bool(true),
						ClientResultColumnCaseInsensitive:        sdk.Bool(true),
						ClientSessionKeepAlive:                   sdk.Bool(true),
						ClientSessionKeepAliveHeartbeatFrequency: sdk.Bool(true),
						ClientTimestampTypeMapping:               sdk.Bool(true),
						DateInputFormat:                          sdk.Bool(true),
						DateOutputFormat:                         sdk.Bool(true),
						EnableUnloadPhysicalTypeOptimization:     sdk.Bool(true),
						ErrorOnNondeterministicMerge:             sdk.Bool(true),
						ErrorOnNondeterministicUpdate:            sdk.Bool(true),
						GeographyOutputFormat:                    sdk.Bool(true),
						GeometryOutputFormat:                     sdk.Bool(true),
						JdbcTreatDecimalAsInt:                    sdk.Bool(true),
						JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
						JdbcUseSessionTimezone:                   sdk.Bool(true),
						JsonIndent:                               sdk.Bool(true),
						LockTimeout:                              sdk.Bool(true),
						LogLevel:                                 sdk.Bool(true),
						MultiStatementCount:                      sdk.Bool(true),
						NoorderSequenceAsDefault:                 sdk.Bool(true),
						OdbcTreatDecimalAsInt:                    sdk.Bool(true),
						QueryTag:                                 sdk.Bool(true),
						QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
						RowsPerResultset:                         sdk.Bool(true),
						S3StageVpceDnsName:                       sdk.Bool(true),
						SearchPath:                               sdk.Bool(true),
						SimulatedDataSharingConsumer:             sdk.Bool(true),
						StatementQueuedTimeoutInSeconds:          sdk.Bool(true),
						StatementTimeoutInSeconds:                sdk.Bool(true),
						StrictJsonOutput:                         sdk.Bool(true),
						TimestampDayIsAlways24h:                  sdk.Bool(true),
						TimestampInputFormat:                     sdk.Bool(true),
						TimestampLTZOutputFormat:                 sdk.Bool(true),
						TimestampNTZOutputFormat:                 sdk.Bool(true),
						TimestampOutputFormat:                    sdk.Bool(true),
						TimestampTypeMapping:                     sdk.Bool(true),
						TimestampTZOutputFormat:                  sdk.Bool(true),
						Timezone:                                 sdk.Bool(true),
						TimeInputFormat:                          sdk.Bool(true),
						TimeOutputFormat:                         sdk.Bool(true),
						TraceLevel:                               sdk.Bool(true),
						TransactionAbortOnError:                  sdk.Bool(true),
						TransactionDefaultIsolationLevel:         sdk.Bool(true),
						TwoDigitCenturyStart:                     sdk.Bool(true),
						UnsupportedDDLAction:                     sdk.Bool(true),
						UseCachedResult:                          sdk.Bool(true),
						WeekOfYearPolicy:                         sdk.Bool(true),
						WeekStart:                                sdk.Bool(true),
					},
					ObjectParameters: &sdk.UserObjectParametersUnset{
						EnableUnredactedQuerySyntaxError: sdk.Bool(true),
						NetworkPolicy:                    sdk.Bool(true),
						PreventUnloadToInternalStages:    sdk.Bool(true),
					},
				},
			}

			err = client.Users.Alter(ctx, id, alterOpts)
			require.NoError(t, err)

			assertThatObject(t, objectparametersassert.UserParameters(t, id).
				HasAllDefaults().
				HasAllDefaultsExplicit(),
			)

			// check that ShowParameters works too
			parameters, err = client.Users.ShowParameters(ctx, id)
			require.NoError(t, err)
			assertThatObject(t, objectparametersassert.UserParametersPrefetched(t, id, parameters).
				HasAllDefaults().
				HasAllDefaultsExplicit(),
			)
		})
	}

	t.Run("alter: set and unset properties and parameters at the same time", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		err := client.Users.Alter(ctx, user.ID(), &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				SessionParameters: &sdk.SessionParameters{
					Autocommit: sdk.Bool(false),
				},
				ObjectParameters: &sdk.UserObjectParameters{
					NetworkPolicy: sdk.Pointer(networkPolicy.ID()),
				},
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						Comment: sdk.String("some comment"),
					},
				},
			},
		})
		require.NoError(t, err)

		err = client.Users.Alter(ctx, user.ID(), &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				SessionParameters: &sdk.SessionParametersUnset{
					Autocommit: sdk.Bool(true),
				},
				ObjectParameters: &sdk.UserObjectParametersUnset{
					NetworkPolicy: sdk.Bool(true),
				},
				ObjectProperties: &sdk.UserObjectPropertiesUnset{
					Comment: sdk.Bool(true),
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("describe: when user exists", func(t *testing.T) {
		userDetails, err := client.Users.Describe(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.Name, userDetails.Name.Value)
	})

	t.Run("describe: when user does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		_, err := client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when user exists", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		id := user.ID()
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		require.NoError(t, err)
		_, err = client.Users.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when user does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		err := client.Users.Drop(ctx, id, &sdk.DropUserOptions{})
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show: with like options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(user.Name),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Equal(t, 1, len(users))
	})

	t.Run("show: with starts with options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			StartsWith: sdk.String(randomPrefix),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Contains(t, users, *user2)
		assert.Equal(t, 2, len(users))
	})

	t.Run("show: with starts with, limit and from options", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit:      sdk.Int(10),
			From:       sdk.String(randomPrefix + "_"),
			StartsWith: sdk.String(randomPrefix),
		}

		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, users, *user)
		assert.Equal(t, 1, len(users))
	})

	t.Run("show: search for a non-existent user", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(NonExistingAccountObjectIdentifier.Name()),
			},
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(users))
	})

	t.Run("show: limit the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowUserOptions{
			Limit: sdk.Int(1),
		}
		users, err := client.Users.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(users))
	})

	// This test proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2817.
	// sql: Scan error on column index 10, name "disabled": sql/driver: couldn't convert "null" into type bool
	t.Run("issue #2817: handle show properly without OWNERSHIP and MANAGE GRANTS", func(t *testing.T) {
		disabledUser, disabledUserCleanup := testClientHelper().User.CreateUserWithOptions(t, testClientHelper().Ids.RandomAccountObjectIdentifier(), &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{Disable: sdk.Bool(true)}})
		t.Cleanup(disabledUserCleanup)

		assertThatObject(t, objectassert.User(t, disabledUser.ID()).
			HasDisabled(true),
		)

		role, roleCleanup := testClientHelper().Role.CreateRoleGrantedToCurrentUser(t)
		t.Cleanup(roleCleanup)

		revertRole := testClientHelper().Role.UseRole(t, role.ID())
		t.Cleanup(revertRole)

		assertThatObject(t, objectassert.User(t, disabledUser.ID()).
			HasDisabled(false),
		)
	})

	t.Run("issue #2817: check the describe behavior", func(t *testing.T) {
		disabledUser, disabledUserCleanup := testClientHelper().User.CreateUserWithOptions(t, testClientHelper().Ids.RandomAccountObjectIdentifier(), &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{Disable: sdk.Bool(true)}})
		t.Cleanup(disabledUserCleanup)

		fetchedDisabledUserDetails, err := client.Users.Describe(ctx, disabledUser.ID())
		require.NoError(t, err)
		require.NotNil(t, fetchedDisabledUserDetails.Disabled)
		require.True(t, fetchedDisabledUserDetails.Disabled.Value)

		role, roleCleanup := testClientHelper().Role.CreateRoleGrantedToCurrentUser(t)
		t.Cleanup(roleCleanup)

		revertRole := testClientHelper().Role.UseRole(t, role.ID())
		t.Cleanup(revertRole)

		fetchedDisabledUserDetails, err = client.Users.Describe(ctx, disabledUser.ID())
		require.ErrorContains(t, err, "Insufficient privileges to operate on user")
		require.Nil(t, fetchedDisabledUserDetails)
	})

	t.Run("issue #2817: check what fields are available when using a user with insufficient privileges to fully inspect another user", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Users.Create(ctx, id, &sdk.CreateUserOptions{
			SessionParameters: &sdk.SessionParameters{
				AbortDetachedQuery:                       sdk.Bool(true),
				Autocommit:                               sdk.Bool(false),
				BinaryInputFormat:                        sdk.Pointer(sdk.BinaryInputFormatUTF8),
				BinaryOutputFormat:                       sdk.Pointer(sdk.BinaryOutputFormatBase64),
				ClientMemoryLimit:                        sdk.Int(1024),
				ClientMetadataRequestUseConnectionCtx:    sdk.Bool(true),
				ClientPrefetchThreads:                    sdk.Int(2),
				ClientResultChunkSize:                    sdk.Int(48),
				ClientResultColumnCaseInsensitive:        sdk.Bool(true),
				ClientSessionKeepAlive:                   sdk.Bool(true),
				ClientSessionKeepAliveHeartbeatFrequency: sdk.Int(2400),
				ClientTimestampTypeMapping:               sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
				DateInputFormat:                          sdk.String("YYYY-MM-DD"),
				DateOutputFormat:                         sdk.String("YY-MM-DD"),
				EnableUnloadPhysicalTypeOptimization:     sdk.Bool(false),
				ErrorOnNondeterministicMerge:             sdk.Bool(false),
				ErrorOnNondeterministicUpdate:            sdk.Bool(true),
				GeographyOutputFormat:                    sdk.Pointer(sdk.GeographyOutputFormatWKB),
				GeometryOutputFormat:                     sdk.Pointer(sdk.GeometryOutputFormatWKB),
				JdbcTreatDecimalAsInt:                    sdk.Bool(false),
				JdbcTreatTimestampNtzAsUtc:               sdk.Bool(true),
				JdbcUseSessionTimezone:                   sdk.Bool(false),
				JsonIndent:                               sdk.Int(4),
				LockTimeout:                              sdk.Int(21222),
				LogLevel:                                 sdk.Pointer(sdk.LogLevelError),
				MultiStatementCount:                      sdk.Int(0),
				NoorderSequenceAsDefault:                 sdk.Bool(false),
				OdbcTreatDecimalAsInt:                    sdk.Bool(true),
				QueryTag:                                 sdk.String("some_tag"),
				QuotedIdentifiersIgnoreCase:              sdk.Bool(true),
				RowsPerResultset:                         sdk.Int(2),
				S3StageVpceDnsName:                       sdk.String("vpce-id.s3.region.vpce.amazonaws.com"),
				SearchPath:                               sdk.String("$public, $current"),
				SimulatedDataSharingConsumer:             sdk.String("some_consumer"),
				StatementQueuedTimeoutInSeconds:          sdk.Int(10),
				StatementTimeoutInSeconds:                sdk.Int(10),
				StrictJsonOutput:                         sdk.Bool(true),
				TimestampDayIsAlways24h:                  sdk.Bool(true),
				TimestampInputFormat:                     sdk.String("YYYY-MM-DD"),
				TimestampLTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
				TimestampNTZOutputFormat:                 sdk.String("YYYY-MM-DD HH24:MI:SS"),
				TimestampOutputFormat:                    sdk.String("YYYY-MM-DD HH24:MI:SS"),
				TimestampTypeMapping:                     sdk.Pointer(sdk.TimestampTypeMappingLtz),
				TimestampTZOutputFormat:                  sdk.String("YYYY-MM-DD HH24:MI:SS"),
				Timezone:                                 sdk.String("Europe/Warsaw"),
				TimeInputFormat:                          sdk.String("HH24:MI"),
				TimeOutputFormat:                         sdk.String("HH24:MI"),
				TraceLevel:                               sdk.Pointer(sdk.TraceLevelPropagate),
				TransactionAbortOnError:                  sdk.Bool(true),
				TransactionDefaultIsolationLevel:         sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
				TwoDigitCenturyStart:                     sdk.Int(1980),
				UnsupportedDDLAction:                     sdk.Pointer(sdk.UnsupportedDDLActionFail),
				UseCachedResult:                          sdk.Bool(false),
				WeekOfYearPolicy:                         sdk.Int(1),
				WeekStart:                                sdk.Int(1),
			},
			ObjectParameters: &sdk.UserObjectParameters{
				EnableUnredactedQuerySyntaxError: sdk.Bool(true),
				NetworkPolicy:                    sdk.Pointer(networkPolicy.ID()),
				PreventUnloadToInternalStages:    sdk.Bool(true),
			},
			ObjectProperties: &sdk.UserObjectProperties{
				Password:              sdk.String(password),
				LoginName:             sdk.String(newValue),
				DisplayName:           sdk.String(newValue),
				FirstName:             sdk.String(newValue),
				MiddleName:            sdk.String(newValue),
				LastName:              sdk.String(newValue),
				Email:                 sdk.String(email),
				MustChangePassword:    sdk.Bool(true),
				Disable:               sdk.Bool(true),
				DaysToExpiry:          sdk.Int(5),
				MinsToUnlock:          sdk.Int(15),
				DefaultWarehouse:      sdk.Pointer(warehouseId),
				DefaultNamespace:      sdk.Pointer(schemaIdObjectIdentifier),
				DefaultRole:           sdk.Pointer(roleId),
				DefaultSecondaryRoles: &sdk.SecondaryRoles{All: sdk.Bool(true)},
				MinsToBypassMFA:       sdk.Int(30),
				RSAPublicKey:          sdk.String(key),
				RSAPublicKey2:         sdk.String(key2),
				Comment:               sdk.String("some comment"),
			},
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		role, roleCleanup := testClientHelper().Role.CreateRoleGrantedToCurrentUser(t)
		t.Cleanup(roleCleanup)

		revertRole := testClientHelper().Role.UseRole(t, role.ID())
		t.Cleanup(revertRole)

		// Both describe and parameters will throw the error
		_, err = testClientHelper().User.Describe(t, id)
		require.ErrorContains(t, err, "Insufficient privileges to operate on user")
		_, err = testClientHelper().Parameter.ShowUserParametersOrError(t, id)
		require.ErrorContains(t, err, "Insufficient privileges to operate on user")

		assertThatObject(t, objectassert.User(t, id).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasLoginName("").
			HasDisplayName("").
			HasFirstName("").
			HasLastName("").
			HasEmail("").
			HasMinsToUnlock("").
			HasDaysToExpiry("").
			HasComment("").
			HasDisabled(false).           // underlying null
			HasMustChangePassword(false). // underlying null
			HasSnowflakeLock(false).      // underlying null
			HasDefaultWarehouse("").
			HasDefaultNamespace("").
			HasDefaultRole("").
			HasDefaultSecondaryRoles("").
			HasExtAuthnDuo(false). // underlying null
			HasExtAuthnUid("").
			HasMinsToBypassMfa("").
			HasOwnerNotEmpty().
			HasLastSuccessLogin(time.Time{}). // underlying null
			HasExpiresAtTimeNotEmpty().
			HasLockedUntilTimeNotEmpty().
			HasHasPassword(false).
			HasHasRsaPublicKey(false).
			HasType(""). // underlying null
			HasHasMfa(false),
		)
	})

	t.Run("login_name and display_name inconsistencies", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// both login_name and display_name were unset so the name is used instead
		assert.Equal(t, id.Name(), userDetails.LoginName.Value)
		assert.Equal(t, id.Name(), userDetails.DisplayName.Value)

		// we unset both values (expecting that it will result in no change)
		unsetBoth := &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				ObjectProperties: &sdk.UserObjectPropertiesUnset{
					LoginName:   sdk.Bool(true),
					DisplayName: sdk.Bool(true),
				},
			},
		}
		err = client.Users.Alter(ctx, id, unsetBoth)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// but login_name is unchanged whereas display_name is nulled out
		assert.Equal(t, id.Name(), userDetails.LoginName.Value)
		assert.Equal(t, "", userDetails.DisplayName.Value)

		// we set both values (expecting that it will result in no change)
		// we use lowercase values on purpose (login_name acts differently than display_name)
		setBoth := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						LoginName:   sdk.String(strings.ToLower(newValue)),
						DisplayName: sdk.String(strings.ToLower(newValue)),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, setBoth)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// they are both set but login_name is uppercase and display_name is lowercase
		assert.Equal(t, strings.ToUpper(newValue), userDetails.LoginName.Value)
		assert.Equal(t, strings.ToLower(newValue), userDetails.DisplayName.Value)

		// we unset both again
		err = client.Users.Alter(ctx, id, unsetBoth)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// and login_name uses name as fallback and display_name does not
		assert.Equal(t, id.Name(), userDetails.LoginName.Value)
		assert.Equal(t, "", userDetails.DisplayName.Value)
	})

	t.Run("default login_name and display_name when the name changes", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// login_name and display_name were not set so the name is used instead
		assert.Equal(t, id.Name(), userDetails.LoginName.Value)
		assert.Equal(t, id.Name(), userDetails.DisplayName.Value)

		// we rename user
		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		rename := &sdk.AlterUserOptions{
			NewName: newId,
		}
		err = client.Users.Alter(ctx, id, rename)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, newId)
		require.NoError(t, err)
		// login_name and display_name are unchanged
		assert.Equal(t, id.Name(), userDetails.LoginName.Value)
		assert.Equal(t, id.Name(), userDetails.DisplayName.Value)

		// we unset both login_name and display_name
		unsetBoth := &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				ObjectProperties: &sdk.UserObjectPropertiesUnset{
					LoginName:   sdk.Bool(true),
					DisplayName: sdk.Bool(true),
				},
			},
		}
		err = client.Users.Alter(ctx, newId, unsetBoth)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, newId)
		require.NoError(t, err)

		// login_name and display_name are changed
		assert.Equal(t, newId.Name(), userDetails.LoginName.Value)
		assert.Equal(t, "", userDetails.DisplayName.Value)
	})

	t.Run("email casing is preserved in Snowflake", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{Email: sdk.String(strings.ToUpper(email))}})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		userShowOutput, err := client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		// email is returned as uppercase both in describe and in show
		assert.Equal(t, strings.ToUpper(email), userDetails.Email.Value)
		assert.Equal(t, strings.ToUpper(email), userShowOutput.Email)

		// we change it to lowercase
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						Email: sdk.String(strings.ToLower(email)),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		userShowOutput, err = client.Users.ShowByID(ctx, id)
		require.NoError(t, err)
		// email is returned as lowercase both in describe and in show
		assert.Equal(t, strings.ToLower(email), userDetails.Email.Value)
		assert.Equal(t, strings.ToLower(email), userShowOutput.Email)
	})

	t.Run("days to expiry setting by hand to a negative value", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		// try to set manually the negative value
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						DaysToExpiry: sdk.Int(-1),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// days to expiry is returned
		assert.NotNil(t, userDetails.DaysToExpiry.Value)
		assert.LessOrEqual(t, *userDetails.DaysToExpiry.Value, float64(-1))
	})

	t.Run("days to expiry set by hand to float value", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		// try to set days to expiry manually to the float value
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`ALTER USER %s SET DAYS_TO_EXPIRY = 1.5`, id.FullyQualifiedName()))
		require.ErrorContains(t, err, "invalid value [1.5] for parameter 'DAYS_TO_EXPIRY'")
	})

	t.Run("days to expiry setting by hand to zero", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		// setting manually to zero
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						DaysToExpiry: sdk.Int(0),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// days to expiry is null
		assert.Nil(t, userDetails.DaysToExpiry.Value)
	})

	t.Run("mins to unlock setting by hand to a negative value", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// mins to unlock is null by default
		assert.Nil(t, userDetails.MinsToUnlock.Value)

		// try to set manually the negative value
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						MinsToUnlock: sdk.Int(-1),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// mins to unlock is returned but not negative but null
		assert.Nil(t, userDetails.MinsToUnlock.Value)
	})

	t.Run("mins to unlock set by hand to float value", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		// try to set mins to unlock manually to the float value
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`ALTER USER %s SET MINS_TO_UNLOCK = 1.5`, id.FullyQualifiedName()))
		require.ErrorContains(t, err, "invalid value [1.5] for parameter 'MINS_TO_UNLOCK'")
	})

	t.Run("mins to unlock setting by hand to zero", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		// setting manually to zero value
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						MinsToUnlock: sdk.Int(0),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// mins to unlock is null
		assert.Nil(t, userDetails.MinsToUnlock.Value)
	})

	t.Run("try to set disable mfa on create", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE USER %s DISABLE_MFA = TRUE`, id.FullyQualifiedName()))
		if err == nil {
			t.Cleanup(testClientHelper().User.DropUserFunc(t, id))
		}
		require.ErrorContains(t, err, "invalid property 'DISABLE_MFA' for 'USER'")
	})

	t.Run("mins to bypass mfa setting by hand to a negative value", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// mins to bypass mfa is null by default
		assert.Nil(t, userDetails.MinsToBypassMfa.Value)

		// try to set manually the negative value
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						MinsToBypassMFA: sdk.Int(-1),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// mins to unlock is returned but not negative but null
		assert.Nil(t, userDetails.MinsToBypassMfa.Value)
	})

	t.Run("mins to bypass mfa setting by hand to zero", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		// setting manually to zero value
		set := &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						MinsToBypassMFA: sdk.Int(0),
					},
				},
			},
		}
		err = client.Users.Alter(ctx, id, set)
		require.NoError(t, err)
		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		// mins to bypass mfa is nil
		require.Nil(t, userDetails.MinsToBypassMfa.Value)
	})

	t.Run("default secondary roles: after BCR 2024_08", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()

		// create, expecting ALL as new default
		err := client.Users.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().User.DropUserFunc(t, id))

		userDetails, err := client.Users.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, `["ALL"]`, userDetails.DefaultSecondaryRoles.Value)

		// set to empty, expecting empty list
		err = client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				ObjectProperties: &sdk.UserAlterObjectProperties{
					UserObjectProperties: sdk.UserObjectProperties{
						DefaultSecondaryRoles: &sdk.SecondaryRoles{None: sdk.Bool(true)},
					},
				},
			},
		})
		require.NoError(t, err)

		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, "[]", userDetails.DefaultSecondaryRoles.Value)

		// unset, expecting ALL
		err = client.Users.Alter(ctx, id, &sdk.AlterUserOptions{
			Unset: &sdk.UserUnset{
				ObjectProperties: &sdk.UserObjectPropertiesUnset{
					DefaultSecondaryRoles: sdk.Bool(true),
				},
			},
		})
		require.NoError(t, err)

		userDetails, err = client.Users.Describe(ctx, id)
		require.NoError(t, err)
		require.Equal(t, `["ALL"]`, userDetails.DefaultSecondaryRoles.Value)
	})
}
