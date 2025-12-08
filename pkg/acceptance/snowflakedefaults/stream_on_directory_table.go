package snowflakedefaults

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func StageIdentifierOutputFormatForStreamOnDirectoryTable(id sdk.SchemaObjectIdentifier) string {
	if getSnowflakeEnvironmentWithProdDefault() == SnowflakeNonProdEnvironment {
		return fmt.Sprintf(`"%s"."%s".%s`, id.DatabaseName(), id.SchemaName(), id.Name())
	}
	return id.Name()
}
