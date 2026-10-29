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
			name: "explicit value name and empty value",
			enum: NewEnum("ParameterType", "ParameterTypes", "ACCOUNT").
				WithValue("CustomName", "CUSTOM").
				WithEmptyValue("SnowflakeDefault"),
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "ParameterTypeAccount", Value: "ACCOUNT", Aliases: nil},
				{Name: "ParameterTypeCustomName", Value: "CUSTOM", Aliases: nil},
				{Name: "ParameterTypeSnowflakeDefault", Value: "", Aliases: nil},
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
		{
			name: "mixed-case value",
			enum: NewEnum("GeographyOutputFormat", "GeographyOutputFormats", "GeoJSON", "WKT"),
			expectedValueRepresentations: []EnumValueRepresentation{
				{Name: "GeographyOutputFormatGeoJson", Value: "GeoJSON", Aliases: nil},
				{Name: "GeographyOutputFormatWkt", Value: "WKT", Aliases: nil},
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

func TestEnum_HasAliases(t *testing.T) {
	t.Run("no aliases", func(t *testing.T) {
		enum := NewEnum("Status", "Statuses", "ACTIVE", "INACTIVE")
		require.False(t, enum.HasAliases())
	})

	t.Run("with aliases", func(t *testing.T) {
		enum := NewEnum("Size", "Sizes", "XSMALL").WithAliases("XSMALL", "X-SMALL")
		require.True(t, enum.HasAliases())
	})
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

	t.Run("panics on undeclared value", func(t *testing.T) {
		require.Panics(t, func() {
			NewEnum("Size", "Sizes", "XSMALL").WithAliases("UNKNOWN", "alias")
		})
	})
}

func TestEnum_WithValue(t *testing.T) {
	t.Run("empty value delegates to WithValue", func(t *testing.T) {
		withValue := NewEnum("ParameterType", "ParameterTypes", "ACCOUNT").WithValue("SnowflakeDefault", "")
		withEmpty := NewEnum("ParameterType", "ParameterTypes", "ACCOUNT").WithEmptyValue("SnowflakeDefault")
		require.Equal(t, withValue.ValueRepresentations(), withEmpty.ValueRepresentations())
	})

	t.Run("panics on empty field name", func(t *testing.T) {
		require.Panics(t, func() {
			NewEnum("ParameterType", "ParameterTypes", "ACCOUNT").WithValue("", "CUSTOM")
		})
	})

	t.Run("panics on duplicate value", func(t *testing.T) {
		require.Panics(t, func() {
			NewEnum("ParameterType", "ParameterTypes", "ACCOUNT").WithValue("Other", "ACCOUNT")
		})
	})

	t.Run("panics on duplicate constant name", func(t *testing.T) {
		require.Panics(t, func() {
			NewEnum("ParameterType", "ParameterTypes", "ACCOUNT").WithValue("Account", "CUSTOM")
		})
	})
}

func TestEnum_IsLowercase(t *testing.T) {
	t.Run("all uppercase values", func(t *testing.T) {
		enum := NewEnum("Status", "Statuses", "ACTIVE", "INACTIVE")
		require.False(t, enum.IsLowercase())
	})

	t.Run("all lowercase values", func(t *testing.T) {
		enum := NewEnum("Role", "Roles", "snowflake_admin", "application")
		require.True(t, enum.IsLowercase())
	})

	t.Run("mixed case values", func(t *testing.T) {
		enum := NewEnum("Mixed", "Mixeds", "Active", "inactive")
		require.False(t, enum.IsLowercase())
	})
}
