package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	openflowPollAttempts    = 60
	openflowPollInterval    = 10 * time.Second
	openflowByocPollAttempts = 180 // 30 minutes for BYOC CloudFormation provisioning
)

func waitForOpenflowDeploymentActive(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		switch deployment.Status {
		case sdk.OpenflowDeploymentStatusActive:
			return nil, true
		case sdk.OpenflowDeploymentStatusCreateFailed:
			return fmt.Errorf("openflow deployment %s entered CREATE_FAILED state", id.Name()), true
		default:
			return nil, false
		}
	})
}

// waitForOpenflowDeploymentActiveByoc waits up to 30 minutes for BYOC deployments
// which require manual CloudFormation provisioning before becoming ACTIVE.
func waitForOpenflowDeploymentActiveByoc(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) error {
	return util.Retry(openflowByocPollAttempts, openflowPollInterval, func() (error, bool) {
		deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		switch deployment.Status {
		case sdk.OpenflowDeploymentStatusActive:
			return nil, true
		case sdk.OpenflowDeploymentStatusCreateFailed:
			return fmt.Errorf("openflow deployment %s entered CREATE_FAILED state", id.Name()), true
		default:
			return nil, false
		}
	})
}

func waitForOpenflowDeploymentDeleted(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		deployment, err := client.OpenflowDeployments.ShowByID(ctx, id)
		if err != nil {
			// Object not found means it's deleted — success.
			return nil, true
		}
		switch deployment.Status {
		case sdk.OpenflowDeploymentStatusDeleted:
			return nil, true
		case sdk.OpenflowDeploymentStatusDeleteFailed:
			return fmt.Errorf("openflow deployment %s entered DELETE_FAILED state", id.Name()), true
		default:
			return nil, false
		}
	})
}

func waitForOpenflowRuntimeActive(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		switch runtime.Status {
		case sdk.OpenflowRuntimeStatusActive:
			return nil, true
		case sdk.OpenflowRuntimeStatusCreateFailed, sdk.OpenflowRuntimeStatusActivateFailed:
			return fmt.Errorf("openflow runtime %s entered failed state: %s", id.Name(), runtime.Status), true
		default:
			return nil, false
		}
	})
}

func waitForOpenflowRuntimeSuspended(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		switch runtime.Status {
		case sdk.OpenflowRuntimeStatusSuspended:
			return nil, true
		case sdk.OpenflowRuntimeStatusSuspendFailed:
			return fmt.Errorf("openflow runtime %s entered SUSPEND_FAILED state", id.Name()), true
		default:
			return nil, false
		}
	})
}

func waitForOpenflowRuntimeDeleted(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		_, err := client.OpenflowRuntimes.ShowByID(ctx, id)
		if err != nil {
			return nil, true
		}
		return nil, false
	})
}

func waitForOpenflowConnectorActive(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		connector, err := client.OpenflowConnectors.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		switch connector.Status {
		case sdk.OpenflowConnectorStatusRunning:
			return nil, true
		case sdk.OpenflowConnectorStatusStartFailed, sdk.OpenflowConnectorStatusCreateFailed:
			return fmt.Errorf("openflow connector %s entered failed state: %s", id.Name(), connector.Status), true
		default:
			return nil, false
		}
	})
}

func waitForOpenflowConnectorStopped(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		connector, err := client.OpenflowConnectors.ShowByID(ctx, id)
		if err != nil {
			return err, false
		}
		switch connector.Status {
		case sdk.OpenflowConnectorStatusStopped:
			return nil, true
		case sdk.OpenflowConnectorStatusStopFailed:
			return fmt.Errorf("openflow connector %s entered STOP_FAILED state", id.Name()), true
		default:
			return nil, false
		}
	})
}

func waitForOpenflowConnectorDeleted(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	return util.Retry(openflowPollAttempts, openflowPollInterval, func() (error, bool) {
		_, err := client.OpenflowConnectors.ShowByID(ctx, id)
		if err != nil {
			return nil, true
		}
		return nil, false
	})
}
