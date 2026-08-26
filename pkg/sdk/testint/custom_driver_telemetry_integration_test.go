//go:build non_account_level_tests

package testint

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake/v2"
	"github.com/stretchr/testify/require"
)

// TestInt_Client_CustomDriverTelemetry is the SNOW-4000794 spike: prove
// gosnowflake.SnowflakeConnection.AddTelemetryData from the SDK client.
//
// Flush: the driver sends a batch of N events or on snowflakeConn.Close
// (Client.Close / db.Close).
// Ingestion in Snowflake is not asserted here; search the printed event_id in
// ingested driver telemetry table (there is no table dedicated for TFP yet).
//
// Opt-out is not asserted. gosnowflake enables in-band telemetry from the
// login-response parameter client_telemetry_enabled, not from Config.Params.
// Putting CLIENT_TELEMETRY_ENABLED=false on the DSN is only a login request;
// the server still returns client_telemetry_enabled=true and the custom event
// is flushed. CLIENT_TELEMETRY_ENABLED is also not a customer-settable SQL
// session/user parameter (undocumented; SHOW PARAMETERS does not list it;
// ALTER SESSION / ALTER USER SET is rejected), so the test account cannot
// force the server echo to false. The old driver DisableTelemetry flag was a
// local no-op; that is not equivalent.
func TestInt_Client_CustomDriverTelemetry(t *testing.T) {
	restoreDriverLogLevelAfterTest(t)

	t.Run("enabled: custom event is sent on Client.Close", func(t *testing.T) {
		logs := captureDriverLogs(t)

		client, err := sdk.NewClient(driverConfigForTelemetryTest(t))
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, client.Close())
		})

		eventID := random.UUID()
		t.Logf("event_id=%s", eventID)

		require.NoError(t, addCustomDriverTelemetry(t, client, map[string]string{
			"source":   "terraform_provider",
			"type":     "terraform_provider_telemetry_integration_test",
			"event_id": eventID,
		}))

		logsAfterReturnToPool := logs.String()
		uploadsBeforeClose := strings.Count(logsAfterReturnToPool, "successfully uploaded metrics to telemetry")

		require.NoError(t, client.Close())

		logsAfterClose := logs.String()
		require.Contains(t, logsAfterClose, eventID)
		require.Contains(t, logsAfterClose, `\"source\":\"terraform_provider\"`)
		require.Greater(t, strings.Count(logsAfterClose, "successfully uploaded metrics to telemetry"), uploadsBeforeClose)

		if payload, ok := lastTelemetryPayload(logsAfterClose); ok {
			t.Logf("flushed telemetry payload: %s", payload)
		}
	})
}

func driverConfigForTelemetryTest(t *testing.T) *gosnowflake.Config {
	t.Helper()
	config := sdk.DefaultConfig()
	config.Tracing = "debug"
	return config
}

// captureDriverLogs intercepts the driver logs and returns a buffer containing the logs.
func captureDriverLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	buf := bytes.NewBuffer(nil)
	logger := gosnowflake.GetLogger()
	// Write the logs both to the buffer and stderr.
	logger.SetOutput(io.MultiWriter(buf, os.Stderr))
	require.NoError(t, logger.SetLogLevel("debug"))
	t.Cleanup(func() {
		logger.SetOutput(os.Stderr)
	})
	return buf
}

func addCustomDriverTelemetry(t *testing.T, client *sdk.Client, data map[string]string) error {
	t.Helper()
	ctx := context.Background()
	conn, err := client.GetConn().Conn(ctx)
	if err != nil {
		return err
	}
	addErr := conn.Raw(func(driverConn any) error {
		sc, ok := driverConn.(gosnowflake.SnowflakeConnection)
		if !ok {
			return fmt.Errorf("driver connection is %T, not gosnowflake.SnowflakeConnection", driverConn)
		}
		return sc.AddTelemetryData(ctx, time.Now(), data)
	})
	closeErr := conn.Close()
	if addErr != nil {
		return addErr
	}
	return closeErr
}

func lastTelemetryPayload(logs string) (string, bool) {
	const marker = "telemetry payload being sent: "
	idx := strings.LastIndex(logs, marker)
	if idx < 0 {
		return "", false
	}
	payload := logs[idx+len(marker):]
	if end := strings.IndexByte(payload, '\n'); end >= 0 {
		payload = payload[:end]
	}
	return payload, true
}
