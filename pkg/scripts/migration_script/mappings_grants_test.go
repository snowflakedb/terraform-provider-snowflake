package main

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// TODO(SNOW-2277608): More tests

func TestHandleGrants(t *testing.T) {
	grantOnAccount := sdk.Grant{
		Privilege:   "CREATE DATABASE",
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

	grantOnAccountResourceModel, err := MapGrantToModel(grantOnAccount)
	assert.NoError(t, err)

	grantOnAccountObjectResourceModel, err := MapGrantToModel(grantOnAccountObject)
	assert.NoError(t, err)

	grantOnSchemaResourceModel, err := MapGrantToModel(grantOnSchema)
	assert.NoError(t, err)

	grantOnSchemaObjectResourceModel, err := MapGrantToModel(grantOnSchemaObject)
	assert.NoError(t, err)

	assert.Equal(t, strings.TrimLeft(`
resource "snowflake_grant_privileges_to_account_role" "test_resource_name_on_account" {
  account_role_name = "TEST_ROLE_ON_ACCOUNT"
  on_account = true
  privileges = ["CREATE DATABASE"]
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
  privileges = ["SELECT"]
  with_grant_option = false
}
`, "\n"),
		config.FromModels(t,
			grantOnAccountResourceModel,
			grantOnAccountObjectResourceModel,
			grantOnSchemaResourceModel,
			grantOnSchemaObjectResourceModel,
		),
	)
}
