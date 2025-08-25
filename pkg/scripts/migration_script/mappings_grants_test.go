package main

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// TODO: More tests

func TestHandleGrants(t *testing.T) {
	grantOnAccount := sdk.Grant{
		Privilege:   "CREATE DATABASE",
		GrantedOn:   sdk.ObjectTypeAccount,
		Name:        sdk.NewAccountObjectIdentifier("TEST_ACCOUNT"),
		GranteeName: sdk.NewAccountObjectIdentifier("TEST_ROLE_ON_ACCOUNT"),
	}
	grantOnAccount2 := sdk.Grant{
		Privilege:   "APPLY MASKING POLICY",
		GrantedOn:   sdk.ObjectTypeAccount,
		Name:        sdk.NewAccountObjectIdentifier("TEST_ACCOUNT"),
		GranteeName: sdk.NewAccountObjectIdentifier("TEST_ROLE_ON_ACCOUNT"),
	}
	grantOnAccountObject := sdk.Grant{
		Privilege:   "USAGE",
		GrantedOn:   sdk.ObjectTypeDatabase,
		Name:        sdk.NewAccountObjectIdentifier("TEST_DATABASE"),
		GranteeName: sdk.NewAccountObjectIdentifier("TEST_ROLE_ON_ACCOUNT_OBJECT"),
	}
	grantOnSchema := sdk.Grant{
		Privilege:   "CREATE TABLE",
		GrantedOn:   sdk.ObjectTypeSchema,
		Name:        sdk.NewDatabaseObjectIdentifier("TEST_DATABASE", "TEST_SCHEMA"),
		GranteeName: sdk.NewAccountObjectIdentifier("TEST_ROLE_ON_SCHEMA"),
	}
	grantOnSchemaObject := sdk.Grant{
		Privilege:   "SELECT",
		GrantedOn:   sdk.ObjectTypeTable,
		Name:        sdk.NewSchemaObjectIdentifier("TEST_DATABASE", "TEST_SCHEMA", "TEST_TABLE"),
		GranteeName: sdk.NewAccountObjectIdentifier("TEST_ROLE_ON_SCHEMA_OBJECT"),
	}
	grantOnSchemaObject2 := sdk.Grant{
		Privilege:   "INSERT",
		GrantedOn:   sdk.ObjectTypeTable,
		Name:        sdk.NewSchemaObjectIdentifier("TEST_DATABASE", "TEST_SCHEMA", "TEST_TABLE"),
		GranteeName: sdk.NewAccountObjectIdentifier("TEST_ROLE_ON_SCHEMA_OBJECT"),
	}

	grantOnAccountResourceModels, err := MapGrantToModel([]sdk.Grant{grantOnAccount, grantOnAccount2})
	assert.NoError(t, err)

	grantOnAccountObjectResourceModels, err := MapGrantToModel([]sdk.Grant{grantOnAccountObject})
	assert.NoError(t, err)

	grantOnSchemaResourceModels, err := MapGrantToModel([]sdk.Grant{grantOnSchema})
	assert.NoError(t, err)

	grantOnSchemaObjectResourceModels, err := MapGrantToModel([]sdk.Grant{grantOnSchemaObject, grantOnSchemaObject2})
	assert.NoError(t, err)

	assert.Equal(t, strings.TrimLeft(`
resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
  account_role_name = "TEST_ROLE_ON_ACCOUNT"
  on_account = true
  privileges = ["CREATE DATABASE", "APPLY MASKING POLICY"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account_object" {
  account_role_name = "TEST_ROLE_ON_ACCOUNT_OBJECT"
  on_account_object {
    object_name = "TEST_DATABASE"
    object_type = "DATABASE"
  }
  privileges = ["USAGE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_schema" {
  account_role_name = "TEST_ROLE_ON_SCHEMA"
  on_schema {
    schema_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\""
  }
  privileges = ["CREATE TABLE"]
  with_grant_option = false
}

resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_schema_object" {
  account_role_name = "TEST_ROLE_ON_SCHEMA_OBJECT"
  on_schema_object {
    object_name = "\"TEST_DATABASE\".\"TEST_SCHEMA\".\"TEST_TABLE\""
    object_type = "TABLE"
  }
  privileges = ["SELECT", "INSERT"]
  with_grant_option = false
}
`, "\n"),
		config.FromModels(t,
			grantOnAccountResourceModels,
			grantOnAccountObjectResourceModels,
			grantOnSchemaResourceModels,
			grantOnSchemaObjectResourceModels,
		),
	)
}
