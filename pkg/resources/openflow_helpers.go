package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Mutating Openflow operations are asynchronous, so the provider polls SHOW until the object settles.
// Polling is bounded by the Terraform timeout for the operation, which users can raise with a `timeouts`
// block, less a minute so it fails with its own error rather than being cut off.
const openflowPollTimeoutReserve = time.Minute

func openflowPollTimeout(timeout time.Duration) time.Duration {
	if timeout > openflowPollTimeoutReserve {
		return timeout - openflowPollTimeoutReserve
	}
	return timeout
}

// openflowWait describes when a poll stops. Anything outside done and failed is retried.
type openflowWait[S ~string] struct {
	done   []S
	failed []S
	// goal completes "waiting for ..." in the error a timeout produces.
	goal string
	// missingIsDone treats an object that is gone as done, which the teardown waits need and the waits for
	// a specific status do not. Only ErrObjectNotFound counts: any other error still fails the wait.
	missingIsDone bool
}

// waitForOpenflowObject polls currentStatus until the wait is satisfied. Deployments, runtimes and connectors
// have separate status types but are polled identically, so each supplies only its statuses, a closure
// reading the current one, and the noun for error messages.
func waitForOpenflowObject[S ~string](
	ctx context.Context,
	objectKind string,
	name string,
	timeout time.Duration,
	currentStatus func() (S, error),
	wait openflowWait[S],
) error {
	return retry.RetryContext(ctx, openflowPollTimeout(timeout), func() *retry.RetryError {
		status, err := currentStatus()
		if err != nil {
			if wait.missingIsDone && errors.Is(err, sdk.ErrObjectNotFound) {
				return nil
			}
			return retry.NonRetryableError(fmt.Errorf("waiting for openflow %s %s: %w", objectKind, name, err))
		}
		switch {
		case slices.Contains(wait.done, status):
			return nil
		case slices.Contains(wait.failed, status):
			return retry.NonRetryableError(fmt.Errorf("openflow %s %s entered %s state, waiting for %s", objectKind, name, status, wait.goal))
		default:
			return retry.RetryableError(fmt.Errorf("openflow %s %s is in %s state, waiting for %s", objectKind, name, status, wait.goal))
		}
	})
}

const openflowDeploymentKind = "deployment"

// deploymentStatus reads a deployment's status, erroring when it is gone. Whether that satisfies the wait or
// fails it is the caller's missingIsDone.
func deploymentStatus(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) func() (sdk.OpenflowDeploymentStatus, error) {
	return func() (sdk.OpenflowDeploymentStatus, error) {
		deployment, err := client.OpenflowDeployments.ShowByIDSafely(ctx, id)
		if err != nil {
			return "", err
		}
		return deployment.Status, nil
	}
}

func waitForOpenflowDeploymentStatus(
	ctx context.Context,
	client *sdk.Client,
	id sdk.AccountObjectIdentifier,
	timeout time.Duration,
	targetStatuses []sdk.OpenflowDeploymentStatus,
	failureStatuses []sdk.OpenflowDeploymentStatus,
) error {
	return waitForOpenflowObject(ctx, openflowDeploymentKind, id.Name(), timeout,
		deploymentStatus(ctx, client, id),
		openflowWait[sdk.OpenflowDeploymentStatus]{
			done:   targetStatuses,
			failed: failureStatuses,
			goal:   fmt.Sprintf("one of %v", targetStatuses),
		})
}

// waitForOpenflowDeploymentSettled waits out a create still in flight. A Snowflake-managed deployment goes
// CREATING, PROVISIONING, ACTIVE; BYOC skips PROVISIONING and settles at INACTIVE. A failed deployment counts
// as settled, since it can still be torn down.
func waitForOpenflowDeploymentSettled(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier, timeout time.Duration) error {
	settled := slices.DeleteFunc(slices.Clone(sdk.AllOpenflowDeploymentStatuses), func(status sdk.OpenflowDeploymentStatus) bool {
		return slices.Contains(sdk.OpenflowDeploymentTransientStatuses, status)
	})
	return waitForOpenflowObject(ctx, openflowDeploymentKind, id.Name(), timeout,
		deploymentStatus(ctx, client, id),
		openflowWait[sdk.OpenflowDeploymentStatus]{
			done:          settled,
			goal:          "it to settle",
			missingIsDone: true,
		})
}

// waitForOpenflowDeploymentTerminated waits for TERMINATE to finish, which DROP requires. A terminated
// deployment stays visible in SHOW until it is dropped, so waiting for it to disappear would wait for the
// wrong thing.
func waitForOpenflowDeploymentTerminated(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowDeploymentKind, id.Name(), timeout,
		deploymentStatus(ctx, client, id),
		openflowWait[sdk.OpenflowDeploymentStatus]{
			done:          []sdk.OpenflowDeploymentStatus{sdk.OpenflowDeploymentStatusDeleted},
			failed:        []sdk.OpenflowDeploymentStatus{sdk.OpenflowDeploymentStatusDeleteFailed},
			goal:          string(sdk.OpenflowDeploymentStatusDeleted),
			missingIsDone: true,
		})
}

const openflowRuntimeKind = "runtime"

func runtimeStatus(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) func() (sdk.OpenflowRuntimeStatus, error) {
	return func() (sdk.OpenflowRuntimeStatus, error) {
		runtime, err := client.OpenflowRuntimes.ShowByIDSafely(ctx, id)
		if err != nil {
			return "", err
		}
		return runtime.Status, nil
	}
}

// waitForOpenflowRuntimeActive waits out a create or an alter. Every mutating runtime statement is
// asynchronous, so an update that returns before the runtime settles leaves the following read seeing a
// transient status and reporting drift.
func waitForOpenflowRuntimeActive(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowRuntimeKind, id.Name(), timeout,
		runtimeStatus(ctx, client, id),
		openflowWait[sdk.OpenflowRuntimeStatus]{
			done: []sdk.OpenflowRuntimeStatus{sdk.OpenflowRuntimeStatusActive},
			failed: []sdk.OpenflowRuntimeStatus{
				sdk.OpenflowRuntimeStatusCreateFailed,
				sdk.OpenflowRuntimeStatusUpdateFailed,
				sdk.OpenflowRuntimeStatusActivateFailed,
				sdk.OpenflowRuntimeStatusRestartFailed,
				sdk.OpenflowRuntimeStatusUpgradeFailed,
			},
			goal: string(sdk.OpenflowRuntimeStatusActive),
		})
}

// waitForOpenflowRuntimeUpdated waits out an ALTER. Snowflake accepts a metadata SET on a suspended
// runtime and leaves it SUSPENDED, so waiting for ACTIVE here would hold until the timeout expires.
func waitForOpenflowRuntimeUpdated(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowRuntimeKind, id.Name(), timeout,
		runtimeStatus(ctx, client, id),
		openflowWait[sdk.OpenflowRuntimeStatus]{
			done: []sdk.OpenflowRuntimeStatus{
				sdk.OpenflowRuntimeStatusActive,
				sdk.OpenflowRuntimeStatusSuspended,
			},
			failed: []sdk.OpenflowRuntimeStatus{
				sdk.OpenflowRuntimeStatusUpdateFailed,
				sdk.OpenflowRuntimeStatusActivateFailed,
				sdk.OpenflowRuntimeStatusSuspendFailed,
				sdk.OpenflowRuntimeStatusRestartFailed,
				sdk.OpenflowRuntimeStatusUpgradeFailed,
			},
			goal: "it to finish updating",
		})
}

// waitForOpenflowRuntimeSettled waits out a mutation still in flight. TERMINATE is refused while the
// runtime is in a transient status, so a destroy racing a create has to let it settle first. A failed
// runtime counts as settled, since it can still be torn down.
func waitForOpenflowRuntimeSettled(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	settled := slices.DeleteFunc(slices.Clone(sdk.AllOpenflowRuntimeStatuses), func(status sdk.OpenflowRuntimeStatus) bool {
		return slices.Contains(sdk.OpenflowRuntimeTransientStatuses, status)
	})
	return waitForOpenflowObject(ctx, openflowRuntimeKind, id.Name(), timeout,
		runtimeStatus(ctx, client, id),
		openflowWait[sdk.OpenflowRuntimeStatus]{
			done:          settled,
			goal:          "it to settle",
			missingIsDone: true,
		})
}

// waitForOpenflowRuntimeTerminated waits for TERMINATE to finish, which DROP requires: Snowflake refuses
// DROP from any status other than DELETED.
func waitForOpenflowRuntimeTerminated(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowRuntimeKind, id.Name(), timeout,
		runtimeStatus(ctx, client, id),
		openflowWait[sdk.OpenflowRuntimeStatus]{
			done:          []sdk.OpenflowRuntimeStatus{sdk.OpenflowRuntimeStatusDeleted},
			failed:        []sdk.OpenflowRuntimeStatus{sdk.OpenflowRuntimeStatusDeleteFailed},
			goal:          string(sdk.OpenflowRuntimeStatusDeleted),
			missingIsDone: true,
		})
}

const openflowConnectorKind = "connector"

func connectorStatus(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) func() (sdk.OpenflowConnectorStatus, error) {
	return func() (sdk.OpenflowConnectorStatus, error) {
		connector, err := client.OpenflowConnectors.ShowByIDSafely(ctx, id)
		if err != nil {
			return "", err
		}
		return connector.Status, nil
	}
}

// waitForOpenflowConnectorReady waits out a create or an alter. A connector settles on RUNNING or STOPPED
// depending on whether its version has been committed.
func waitForOpenflowConnectorReady(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowConnectorKind, id.Name(), timeout,
		connectorStatus(ctx, client, id),
		openflowWait[sdk.OpenflowConnectorStatus]{
			done: []sdk.OpenflowConnectorStatus{
				sdk.OpenflowConnectorStatusRunning,
				sdk.OpenflowConnectorStatusStopped,
			},
			// The SDK's failure statuses, plus the troubleshooting ones. done lists only RUNNING and STOPPED,
			// so an unnamed status is retried to the timeout rather than reported as the failure it is.
			failed: append(
				slices.Clone(sdk.OpenflowConnectorFailureStatuses),
				sdk.OpenflowConnectorStatusTroubleshooting,
				sdk.OpenflowConnectorStatusEnterTroubleshootingFailed,
				sdk.OpenflowConnectorStatusExitTroubleshootingFailed,
			),
			goal: "it to be ready",
		})
}

// waitForOpenflowConnectorSettled waits out a mutation still in flight, which TERMINATE requires.
func waitForOpenflowConnectorSettled(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	settled := slices.DeleteFunc(slices.Clone(sdk.AllOpenflowConnectorStatuses), func(status sdk.OpenflowConnectorStatus) bool {
		return slices.Contains(sdk.OpenflowConnectorTransientStatuses, status)
	})
	return waitForOpenflowObject(ctx, openflowConnectorKind, id.Name(), timeout,
		connectorStatus(ctx, client, id),
		openflowWait[sdk.OpenflowConnectorStatus]{
			done:          settled,
			goal:          "it to settle",
			missingIsDone: true,
		})
}

// waitForOpenflowConnectorStopped waits for STOP to finish. Delete needs it because a connector in a failed
// status refuses every other teardown statement.
func waitForOpenflowConnectorStopped(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowConnectorKind, id.Name(), timeout,
		connectorStatus(ctx, client, id),
		openflowWait[sdk.OpenflowConnectorStatus]{
			done:          []sdk.OpenflowConnectorStatus{sdk.OpenflowConnectorStatusStopped},
			failed:        []sdk.OpenflowConnectorStatus{sdk.OpenflowConnectorStatusStopFailed},
			goal:          "it to stop",
			missingIsDone: true,
		})
}

// waitForOpenflowConnectorRunning waits for START to finish. waitForOpenflowConnectorReady would not do,
// since it also accepts the STOPPED a connector reports until START has taken effect. A connector that
// can't start, such as one missing required parameters, goes from STARTING back to STOPPED without
// reporting START_FAILED, so a return to STOPPED after leaving it fails the wait.
func waitForOpenflowConnectorRunning(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	currentStatus := connectorStatus(ctx, client, id)
	leftStopped := false
	return waitForOpenflowObject(ctx, openflowConnectorKind, id.Name(), timeout,
		func() (sdk.OpenflowConnectorStatus, error) {
			status, err := currentStatus()
			if err != nil {
				return status, err
			}
			if status != sdk.OpenflowConnectorStatusStopped {
				leftStopped = true
			} else if leftStopped {
				return status, fmt.Errorf("it returned to %s without running", sdk.OpenflowConnectorStatusStopped)
			}
			return status, nil
		},
		openflowWait[sdk.OpenflowConnectorStatus]{
			done:   []sdk.OpenflowConnectorStatus{sdk.OpenflowConnectorStatusRunning},
			failed: []sdk.OpenflowConnectorStatus{sdk.OpenflowConnectorStatusStartFailed},
			goal:   "it to run",
		})
}

// waitForOpenflowConnectorTerminated waits for TERMINATE to finish, which DROP requires: Snowflake refuses
// DROP from any status other than DELETED.
func waitForOpenflowConnectorTerminated(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, timeout time.Duration) error {
	return waitForOpenflowObject(ctx, openflowConnectorKind, id.Name(), timeout,
		connectorStatus(ctx, client, id),
		openflowWait[sdk.OpenflowConnectorStatus]{
			done:          []sdk.OpenflowConnectorStatus{sdk.OpenflowConnectorStatusDeleted},
			failed:        []sdk.OpenflowConnectorStatus{sdk.OpenflowConnectorStatusDeleteFailed},
			goal:          string(sdk.OpenflowConnectorStatusDeleted),
			missingIsDone: true,
		})
}
