package testacc

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func shareConfigOneAccount(shareId sdk.AccountObjectIdentifier, comment string, account string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%s"
	comment        = "%s"
	accounts       = ["%s"]
}
`, shareId.Name(), comment, account)
}
