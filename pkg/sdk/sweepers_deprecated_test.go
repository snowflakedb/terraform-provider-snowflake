package sdk_test

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func getAccountPolicyAttachmentsSweeper(client *sdk.Client) func() error {
	return func() error {
		log.Printf("[DEBUG] Unsetting password and session policies set on the account level")
		ctx := context.Background()
		_ = client.Accounts.UnsetPolicySafely(ctx, sdk.PolicyKindPasswordPolicy)
		_ = client.Accounts.UnsetPolicySafely(ctx, sdk.PolicyKindSessionPolicy)
		return nil
	}
}

func getResourceMonitorSweeper(client *sdk.Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping resource monitors with suffix %s", suffix)
		ctx := context.Background()

		rms, err := client.ResourceMonitors.Show(ctx, nil)
		if err != nil {
			return fmt.Errorf("sweeping resource monitor ended with error, err = %w", err)
		}
		for _, rm := range rms {
			if strings.HasSuffix(rm.Name, suffix) {
				log.Printf("[DEBUG] Dropping resource monitor %s", rm.ID().FullyQualifiedName())
				if err := client.ResourceMonitors.Drop(ctx, rm.ID(), &sdk.DropResourceMonitorOptions{IfExists: sdk.Bool(true)}); err != nil {
					return fmt.Errorf("sweeping resource monitor %s ended with error, err = %w", rm.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping resource monitor %s", rm.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}

func getFailoverGroupSweeper(client *sdk.Client, suffix string) func() error {
	return func() error {
		log.Printf("[DEBUG] Sweeping failover groups with suffix %s", suffix)
		ctx := context.Background()

		currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
		if err != nil {
			return fmt.Errorf("sweeping failover groups ended with error, err = %w", err)
		}
		opts := &sdk.ShowFailoverGroupOptions{
			InAccount: sdk.NewAccountIdentifierFromAccountLocator(currentAccount),
		}
		fgs, err := client.FailoverGroups.Show(ctx, opts)
		if err != nil {
			return fmt.Errorf("sweeping failover groups ended with error, err = %w", err)
		}
		for _, fg := range fgs {
			if strings.HasSuffix(fg.Name, suffix) && fg.AccountLocator == currentAccount {
				log.Printf("[DEBUG] Dropping failover group %s", fg.ID().FullyQualifiedName())
				if err := client.FailoverGroups.Drop(ctx, fg.ID(), nil); err != nil {
					return fmt.Errorf("sweeping failover group %s ended with error, err = %w", fg.ID().FullyQualifiedName(), err)
				}
			} else {
				log.Printf("[DEBUG] Skipping failover group %s", fg.ID().FullyQualifiedName())
			}
		}
		return nil
	}
}
