package datasources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// sharesShowInSdkExt completes the recipients that SHOW SHARES truncates. A share that cannot be completed is a warning.
func sharesShowInSdkExt(ctx context.Context, client *sdk.Client, req *sdk.ShowShareRequest, diags *diag.Diagnostics) ([]sdk.Share, error) {
	shares, err := client.Shares.Show(ctx, req)
	if err != nil {
		return nil, err
	}
	for i := range shares {
		share := &shares[i]
		if err := share.CompleteTo(ctx, client.Shares); err != nil {
			*diags = append(*diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Share %s has %d recipients but SHOW SHARES lists only %d", share.ID().FullyQualifiedName(), share.ToTotal, len(share.To)),
				Detail: "The full list is read from SNOWFLAKE.ACCOUNT_USAGE.SHARES, which failed (the role needs SELECT on that view, for example through IMPORTED PRIVILEGES on the SNOWFLAKE database, and the connection needs a warehouse). " +
					"Recipients that are not listed are missing from the output. Error: " + err.Error(),
			})
		}
	}
	return shares, nil
}

// sharesDescribeInSdkExt describes an inbound share as a consumer and an outbound one as a provider. A share that cannot be
// described is a warning, and its describe output stays empty.
func sharesDescribeInSdkExt(ctx context.Context, client *sdk.Client, share sdk.Share, diags *diag.Diagnostics) ([]sdk.ShareDetails, error) {
	var (
		details []sdk.ShareDetails
		err     error
	)
	if share.Kind == sdk.ShareKindInbound {
		details, err = client.Shares.DescribeConsumer(ctx, share.ExternalID())
	} else {
		details, err = client.Shares.DescribeProvider(ctx, share.ID())
	}
	if err != nil {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Failed to describe share %s", share.Name),
			Detail:   fmt.Sprintf("The describe_output of the share is empty. Error: %s", err),
		})
		return nil, nil
	}
	return details, nil
}
