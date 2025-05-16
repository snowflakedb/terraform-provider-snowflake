package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func v092DatabaseStateUpgrader(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	client := meta.(*provider.Context).Client

	if rawState == nil {
		return rawState, nil
	}

	if v, ok := rawState["from_share"]; ok && v != nil && len(v.(map[string]any)) > 0 {
		return nil, fmt.Errorf("failed to upgrade the state with database created from share, please use snowflake_shared_database instead")
	}

	if v, ok := rawState["from_replica"]; ok && v != nil && len(v.(string)) > 0 {
		return nil, fmt.Errorf("failed to upgrade the state with database created from replica, please use snowflake_secondary_database instead")
	}

	if v, ok := rawState["from_database"]; ok && v != nil && len(v.(string)) > 0 {
		return nil, fmt.Errorf("failed to upgrade the state with database created from database, please use snowflake_database instead. Disclaimer: Right now, database cloning is not supported. They can be imported into the mentioned resource, but any difference in behavior from standard database won't be handled (and can result in errors)")
	}

	if replicationConfigurations, ok := rawState["replication_configuration"]; ok && len(replicationConfigurations.([]any)) == 1 {
		replicationConfiguration := replicationConfigurations.([]any)[0].(map[string]any)
		replication := make(map[string]any)
		replication["ignore_edition_check"] = replicationConfiguration["ignore_edition_check"].(bool)

		accountLocators := replicationConfiguration["accounts"].([]any)
		enableForAccounts := make([]map[string]any, len(accountLocators))

		if len(accountLocators) > 0 {
			replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
			if err != nil {
				return nil, err
			}

			for i, accountLocator := range accountLocators {
				replicationAccount, err := collections.FindFirst(replicationAccounts, func(account *sdk.ReplicationAccount) bool {
					return account.AccountLocator == accountLocator
				})
				if err != nil {
					return nil, fmt.Errorf("couldn't find replication account locator '%s', err = %w", accountLocator, err)
				}
				foundReplicationAccount := *replicationAccount
				enableForAccounts[i] = map[string]any{
					"account_identifier": sdk.NewAccountIdentifier(foundReplicationAccount.OrganizationName, foundReplicationAccount.AccountName),
				}
			}
		}
	}

	return rawState, nil
}
