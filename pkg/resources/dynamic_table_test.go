//go:build !account_level_tests

package resources_test

import (
	"context"
	"fmt"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("snowflake_dynamic_table", &resource.Sweeper{
		Name: "snowflake_dynamic_table",
		F: func(profile string) error {
			client, err := sdk.NewDefaultClient(sdk.WithUseLegacyTomlFormat(true))
			if err != nil {
				return fmt.Errorf("error getting default client during sweep: %w", err)
			}
			ctx := context.Background()
			dynamicTables, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithIn(&sdk.In{
				Schema: acc.TestClient().Ids.SchemaId(),
			}))
			if err != nil {
				return fmt.Errorf("error getting dynamic tables during sweep: %w", err)
			}
			for _, dynamicTable := range dynamicTables {
				err := client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(dynamicTable.ID()))
				if err != nil {
					return fmt.Errorf("error dropping dynamic table %s during sweep: %w", dynamicTable.ID(), err)
				}
			}
			return nil
		},
	})
}
