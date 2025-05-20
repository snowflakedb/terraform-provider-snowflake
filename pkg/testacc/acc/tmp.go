package acc

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
)

// "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testacc/acc"

var TestDatabaseName = acceptance.TestDatabaseName
var TestSchemaName = acceptance.TestSchemaName
var TestWarehouseName = acceptance.TestWarehouseName
var NonExistingAccountObjectIdentifier = acceptance.NonExistingAccountObjectIdentifier
var NonExistingDatabaseObjectIdentifier = acceptance.NonExistingDatabaseObjectIdentifier
var TestAccProvider = acceptance.TestAccProvider
var TestAccProtoV6ProviderFactories = acceptance.TestAccProtoV6ProviderFactories

func TestClient() *helpers.TestClient {
	return acceptance.TestClient()
}

func SecondaryTestClient() *helpers.TestClient {
	return acceptance.SecondaryTestClient()
}
