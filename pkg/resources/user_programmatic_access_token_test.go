package resources

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserProgrammaticAccessToken_ExpiryChangeDiff(t *testing.T) {
	for _, rotate := range []bool{false, true} {
		name := "expiry only"
		keeper := "original"
		if rotate {
			name = "expiry and keeper"
			keeper = "rotated"
		}
		t.Run(name, func(t *testing.T) {
			diff, err := UserProgrammaticAccessToken().Diff(
				context.Background(),
				&terraform.InstanceState{
					ID: "existing-token",
					RawPlan: cty.ObjectVal(map[string]cty.Value{
						"keeper": cty.StringVal(keeper),
					}),
					Attributes: map[string]string{
						"user":                             "TEST_USER",
						"name":                             "TEST_TOKEN",
						"keeper":                           "original",
						"disabled":                         BooleanDefault,
						"expire_rotated_token_after_hours": "24",
						"token":                            "secret",
						"rotated_token_name":               "previous-token",
					},
				},
				terraform.NewResourceConfigRaw(map[string]any{
					"user":                             "TEST_USER",
					"name":                             "TEST_TOKEN",
					"keeper":                           keeper,
					"expire_rotated_token_after_hours": 48,
				}),
				&provider.Context{Client: &sdk.Client{}},
			)
			require.NoError(t, err)
			require.NotNil(t, diff)
			assert.False(t, diff.RequiresNew())
			require.Contains(t, diff.Attributes, "expire_rotated_token_after_hours")
			assert.Equal(t, "48", diff.Attributes["expire_rotated_token_after_hours"].New)
			for _, attribute := range []string{"token", "rotated_token_name"} {
				if rotate {
					require.Contains(t, diff.Attributes, attribute)
					assert.True(t, diff.Attributes[attribute].NewComputed)
				} else {
					assert.NotContains(t, diff.Attributes, attribute)
				}
			}
		})
	}
}

func TestShouldRotateToken(t *testing.T) {
	tests := []struct {
		name     string
		old      string
		new      string
		isKnown  bool
		expected bool
	}{
		// Cases where old is empty
		{
			name:     "old empty, new empty, known",
			old:      "",
			new:      "",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old empty, new empty, unknown",
			old:      "",
			new:      "",
			isKnown:  false,
			expected: false,
		},
		// Cases where old is empty and the value was added to the config.
		{
			name:     "old empty, new non-empty, known",
			old:      "",
			new:      "new_value",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old empty, new non-empty, unknown",
			old:      "",
			new:      "new_value",
			isKnown:  false,
			expected: false,
		},

		// Cases where old is non-empty and new is empty (the value is removed from the config)
		{
			name:     "old non-empty, new empty, known",
			old:      "old_value",
			new:      "",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old non-empty, new empty, unknown",
			old:      "old_value",
			new:      "",
			isKnown:  false,
			expected: true,
		},

		// Cases where old and new are the same non-empty values
		{
			name:     "old and new same non-empty, known",
			old:      "same_value",
			new:      "same_value",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old and new same non-empty, unknown",
			old:      "same_value",
			new:      "same_value",
			isKnown:  false,
			expected: true,
		},

		// Cases where old and new are different non-empty values
		{
			name:     "old and new different non-empty, known",
			old:      "old_value",
			new:      "new_value",
			isKnown:  true,
			expected: true,
		},
		{
			name:     "old and new different non-empty, unknown",
			old:      "old_value",
			new:      "new_value",
			isKnown:  false,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldRotateToken(tt.old, tt.new, tt.isKnown)
			assert.Equal(t, tt.expected, result)
		})
	}
}
