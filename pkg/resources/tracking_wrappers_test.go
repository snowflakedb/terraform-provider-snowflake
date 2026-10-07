package resources

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_TrackingWrappers_emitOnEveryOperation(t *testing.T) {
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
	const failedEmit = "[DEBUG] failed to emit resource_op telemetry: client is not connected"

	t.Run("create", func(t *testing.T) {
		assertDiagWrapper(t, &logs, failedEmit, func(impl schema.CreateContextFunc) diag.Diagnostics {
			return TrackingCreateWrapper(resources.Warehouse, impl)(t.Context(), newResourceData(t), meta)
		})
	})
	t.Run("read", func(t *testing.T) {
		assertDiagWrapper(t, &logs, failedEmit, func(impl schema.ReadContextFunc) diag.Diagnostics {
			return TrackingReadWrapper(resources.Warehouse, impl)(t.Context(), newResourceData(t), meta)
		})
	})
	t.Run("update", func(t *testing.T) {
		assertDiagWrapper(t, &logs, failedEmit, func(impl schema.UpdateContextFunc) diag.Diagnostics {
			return TrackingUpdateWrapper(resources.Warehouse, impl)(t.Context(), newResourceData(t), meta)
		})
	})
	t.Run("delete", func(t *testing.T) {
		assertDiagWrapper(t, &logs, failedEmit, func(impl schema.DeleteContextFunc) diag.Diagnostics {
			return TrackingDeleteWrapper(resources.Warehouse, impl)(t.Context(), newResourceData(t), meta)
		})
	})
	t.Run("import", func(t *testing.T) {
		testCases := []struct {
			name      string
			impl      schema.StateContextFunc
			wantError bool
		}{
			{
				name: "success",
				impl: func(_ context.Context, data *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
					data.SetId(`"DB"."PUBLIC"."WH"`)
					return []*schema.ResourceData{data}, nil
				},
			},
			{
				name: "error",
				impl: func(_ context.Context, _ *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
					return nil, fmt.Errorf("boom")
				},
				wantError: true,
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				logs.Reset()
				_, err := TrackingImportWrapper(resources.Warehouse, testCase.impl)(t.Context(), newResourceData(t), meta)
				require.Equal(t, testCase.wantError, err != nil)
				require.Contains(t, logs.String(), failedEmit)
			})
		}
	})
}

func Test_TrackingWrappers_emitAndRepanic(t *testing.T) {
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
	d := newResourceData(t)

	t.Run("create", func(t *testing.T) {
		logs.Reset()
		require.PanicsWithValue(t, "exploded", func() {
			TrackingCreateWrapper(resources.Warehouse, func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
				panic("exploded")
			})(t.Context(), d, meta)
		})
	})
	t.Run("read", func(t *testing.T) {
		logs.Reset()
		require.PanicsWithValue(t, "exploded", func() {
			TrackingReadWrapper(resources.Warehouse, func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
				panic("exploded")
			})(t.Context(), d, meta)
		})
	})
	t.Run("update", func(t *testing.T) {
		logs.Reset()
		require.PanicsWithValue(t, "exploded", func() {
			TrackingUpdateWrapper(resources.Warehouse, func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
				panic("exploded")
			})(t.Context(), d, meta)
		})
	})
	t.Run("delete", func(t *testing.T) {
		logs.Reset()
		require.PanicsWithValue(t, "exploded", func() {
			TrackingDeleteWrapper(resources.Warehouse, func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
				panic("exploded")
			})(t.Context(), d, meta)
		})
	})
	t.Run("import", func(t *testing.T) {
		logs.Reset()
		require.PanicsWithValue(t, "exploded", func() {
			_, _ = TrackingImportWrapper(resources.Warehouse, func(context.Context, *schema.ResourceData, any) ([]*schema.ResourceData, error) {
				panic("exploded")
			})(t.Context(), d, meta)
		})
	})
}

func newResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]any{})
}

func assertDiagWrapper[T ~func(context.Context, *schema.ResourceData, any) diag.Diagnostics](
	t *testing.T,
	logs *bytes.Buffer,
	failedEmit string,
	run func(T) diag.Diagnostics,
) {
	t.Helper()
	testCases := []struct {
		name      string
		impl      T
		wantError bool
	}{
		{
			name: "success",
			impl: T(func(_ context.Context, data *schema.ResourceData, _ any) diag.Diagnostics {
				data.SetId(`"DB"."PUBLIC"."WH"`)
				return nil
			}),
		},
		{
			name: "error",
			impl: T(func(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
				return diag.Errorf("boom")
			}),
			wantError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logs.Reset()
			diags := run(testCase.impl)
			require.Equal(t, testCase.wantError, diags.HasError())
			require.Contains(t, logs.String(), failedEmit)
		})
	}
}
