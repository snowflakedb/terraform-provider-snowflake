package gen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_ValueRepresentations(t *testing.T) {
	tests := []struct {
		name                         string
		enum                         *Enum
		expectedValueRepresentations []EnumValueRepresentation
	}{
		{
			name: "single value",
			enum: NewEnum("Status", "Statuses", "ACTIVE"),
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "StatusActive", Value: "ACTIVE", Aliases: nil},
			},
		},
		{
			name: "multiple values",
			enum: NewEnum("TokenStatus", "TokenStatuses", "ACTIVE_VALUE", "INACTIVE_VALUE", "EXPIRED_VALUE"),
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "TokenStatusActiveValue", Value: "ACTIVE_VALUE", Aliases: nil},
				{Name: "TokenStatusInactiveValue", Value: "INACTIVE_VALUE", Aliases: nil},
				{Name: "TokenStatusExpiredValue", Value: "EXPIRED_VALUE", Aliases: nil},
			},
		},
		{
			name: "values with aliases",
			enum: NewEnum("Size", "Sizes", "XSMALL", "SMALL", "XXLARGE").
				WithAliases("XSMALL", "X-SMALL").
				WithAliases("XXLARGE", "X2LARGE", "2X-LARGE"),
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "SizeXsmall", Value: "XSMALL", Aliases: []string{"X-SMALL"}},
				{Name: "SizeSmall", Value: "SMALL", Aliases: nil},
				{Name: "SizeXxlarge", Value: "XXLARGE", Aliases: []string{"X2LARGE", "2X-LARGE"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.enum.ValueRepresentations()

			require.Len(t, result, len(tt.expectedValueRepresentations))
			require.Equal(t, tt.expectedValueRepresentations, result)
		})
	}
}

func TestEnum_WithAliases(t *testing.T) {
	t.Run("single alias", func(t *testing.T) {
		enum := NewEnum("Size", "Sizes", "XSMALL", "SMALL").
			WithAliases("XSMALL", "X-SMALL")
		require.Equal(t, []string{"X-SMALL"}, enum.Aliases["XSMALL"])
		require.Nil(t, enum.Aliases["SMALL"])
	})

	t.Run("multiple aliases for same value", func(t *testing.T) {
		enum := NewEnum("Size", "Sizes", "XXLARGE").
			WithAliases("XXLARGE", "X2LARGE", "2X-LARGE")
		require.Equal(t, []string{"X2LARGE", "2X-LARGE"}, enum.Aliases["XXLARGE"])
	})

	t.Run("aliases for multiple values", func(t *testing.T) {
		enum := NewEnum("Size", "Sizes", "XSMALL", "XLARGE").
			WithAliases("XSMALL", "X-SMALL").
			WithAliases("XLARGE", "X-LARGE")
		require.Equal(t, []string{"X-SMALL"}, enum.Aliases["XSMALL"])
		require.Equal(t, []string{"X-LARGE"}, enum.Aliases["XLARGE"])
	})
}
