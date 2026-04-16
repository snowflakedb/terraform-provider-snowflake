package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToCatalogIntegrationCatalogSourceType(t *testing.T) {
	type test struct {
		input string
		want  CatalogIntegrationCatalogSourceType
	}

	valid := []test{
		// case insensitive
		{input: "glue", want: CatalogIntegrationCatalogSourceTypeGlue},

		// Supported Values
		{input: "GLUE", want: CatalogIntegrationCatalogSourceTypeGlue},
		{input: "OBJECT_STORE", want: CatalogIntegrationCatalogSourceTypeObjectStore},
		{input: "POLARIS", want: CatalogIntegrationCatalogSourceTypePolaris},
		{input: "ICEBERG_REST", want: CatalogIntegrationCatalogSourceTypeIcebergRest},
		{input: "SAP_BDC", want: CatalogIntegrationCatalogSourceTypeSapBdc},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCatalogIntegrationCatalogSourceType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCatalogIntegrationCatalogSourceType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ToCatalogIntegrationTableFormat(t *testing.T) {
	type test struct {
		input string
		want  CatalogIntegrationTableFormat
	}

	valid := []test{
		// case insensitive
		{input: "iceberg", want: CatalogIntegrationTableFormatIceberg},

		// Supported Values
		{input: "ICEBERG", want: CatalogIntegrationTableFormatIceberg},
		{input: "DELTA", want: CatalogIntegrationTableFormatDelta},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCatalogIntegrationTableFormat(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCatalogIntegrationTableFormat(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ToCatalogIntegrationRestAuthenticationType(t *testing.T) {
	type test struct {
		input string
		want  CatalogIntegrationRestAuthenticationType
	}

	valid := []test{
		// case insensitive
		{input: "oauth", want: CatalogIntegrationRestAuthenticationTypeOauth},

		// Supported Values
		{input: "OAUTH", want: CatalogIntegrationRestAuthenticationTypeOauth},
		{input: "BEARER", want: CatalogIntegrationRestAuthenticationTypeBearer},
		{input: "SIGV4", want: CatalogIntegrationRestAuthenticationTypeSigv4},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCatalogIntegrationRestAuthenticationType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCatalogIntegrationRestAuthenticationType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ToCatalogIntegrationAccessDelegationMode(t *testing.T) {
	type test struct {
		input string
		want  CatalogIntegrationAccessDelegationMode
	}

	valid := []test{
		// case insensitive
		{input: "vended_credentials", want: CatalogIntegrationAccessDelegationModeVendedCredentials},

		// Supported Values
		{input: "VENDED_CREDENTIALS", want: CatalogIntegrationAccessDelegationModeVendedCredentials},
		{input: "EXTERNAL_VOLUME_CREDENTIALS", want: CatalogIntegrationAccessDelegationModeExternalVolumeCredentials},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCatalogIntegrationAccessDelegationMode(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCatalogIntegrationAccessDelegationMode(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ToCatalogIntegrationCatalogApiType(t *testing.T) {
	type test struct {
		input string
		want  CatalogIntegrationCatalogApiType
	}

	valid := []test{
		// case insensitive
		{input: "public", want: CatalogIntegrationCatalogApiTypePublic},

		// Supported Values
		{input: "PUBLIC", want: CatalogIntegrationCatalogApiTypePublic},
		{input: "PRIVATE", want: CatalogIntegrationCatalogApiTypePrivate},
		{input: "AWS_API_GATEWAY", want: CatalogIntegrationCatalogApiTypeAwsApiGateway},
		{input: "AWS_PRIVATE_API_GATEWAY", want: CatalogIntegrationCatalogApiTypeAwsPrivateApiGateway},
		{input: "AWS_GLUE", want: CatalogIntegrationCatalogApiTypeAwsGlue},
		{input: "AWS_PRIVATE_GLUE", want: CatalogIntegrationCatalogApiTypeAwsPrivateGlue},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToCatalogIntegrationCatalogApiType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToCatalogIntegrationCatalogApiType(tc.input)
			require.Error(t, err)
		})
	}
}
