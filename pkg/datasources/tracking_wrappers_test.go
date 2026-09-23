package datasources

import (
	"bytes"
	"context"
	"log"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_TrackingReadWrapper_emitsOnEveryRead(t *testing.T) {
	var logs bytes.Buffer
	originalLogOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() {
		log.SetOutput(originalLogOutput)
	})

	meta := &provider.Context{
		Client: &sdk.Client{},
		SpanID: "span",
	}
	testCases := []struct {
		name      string
		read      schema.ReadContextFunc
		wantError bool
	}{
		{
			name: "success",
			read: func(_ context.Context, data *schema.ResourceData, _ any) diag.Diagnostics {
				data.SetId("databases_read")
				return nil
			},
		},
		{
			name: "error",
			read: func(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
				return diag.Errorf("boom")
			},
			wantError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logs.Reset()
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]any{})

			diags := TrackingReadWrapper(datasources.Databases, testCase.read)(t.Context(), d, meta)

			require.Equal(t, testCase.wantError, diags.HasError())
			require.Contains(t, logs.String(), "[DEBUG] failed to emit datasource_op telemetry: client is not connected")
		})
	}
}
