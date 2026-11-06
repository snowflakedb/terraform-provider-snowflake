package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const organizationListingManifestFromConfig = `title: Sales data
description: Daily sales
organization_profile: internal
organization_targets:
  access:
  - account: dsconsumer
    roles:
    - analyst
locations:
  access_regions:
  - name: public.aws_us_west_2
`

// organizationListingManifestFromDescribe mirrors the manifest_yaml column returned by DESCRIBE LISTING for the
// manifest above: values are quoted, identifiers are upper-cased, keys follow the Snowflake order, and Snowflake adds
// the resharing section.
const organizationListingManifestFromDescribe = `title: "Sales data"
description: "Daily sales"
organization_profile: "INTERNAL"
organization_targets:
  access:
  - account: "DSCONSUMER"
    roles:
    - "ANALYST"
locations:
  access_regions:
  - name: "PUBLIC.AWS_US_WEST_2"
resharing:
  enabled: true
  only_within_organization: true
`

func Test_organizationListingManifestContains(t *testing.T) {
	testCases := []struct {
		name     string
		actual   string
		expected string
		want     bool
	}{
		{
			name:     "manifest returned by DESCRIBE LISTING contains the configured one",
			actual:   organizationListingManifestFromDescribe,
			expected: organizationListingManifestFromConfig,
			want:     true,
		},
		{
			name:     "configured manifest does not contain keys added by Snowflake",
			actual:   organizationListingManifestFromConfig,
			expected: organizationListingManifestFromDescribe,
			want:     false,
		},
		{
			name:     "changed value",
			actual:   `title: "Sales data v2"`,
			expected: `title: Sales data`,
			want:     false,
		},
		{
			name:     "case change outside of identifiers",
			actual:   `title: "SALES DATA"`,
			expected: `title: Sales data`,
			want:     false,
		},
		{
			name:     "missing key",
			actual:   `title: "Sales data"`,
			expected: "title: Sales data\ndescription: Daily sales",
			want:     false,
		},
		{
			name:     "removed target account",
			actual:   "organization_targets:\n  access:\n  - account: \"A\"",
			expected: "organization_targets:\n  access:\n  - account: a\n  - account: b",
			want:     false,
		},
		{
			name:     "changed target account",
			actual:   "organization_targets:\n  access:\n  - account: \"A\"",
			expected: "organization_targets:\n  access:\n  - account: b",
			want:     false,
		},
		{
			name:     "region name outside of access regions is case-sensitive",
			actual:   "usage_examples:\n- name: \"EXAMPLE\"",
			expected: "usage_examples:\n- name: example",
			want:     false,
		},
		{
			name:     "boolean and number values",
			actual:   "resharing:\n  enabled: true\ntrial_details:\n  trial_time_limit: 7",
			expected: "resharing:\n  enabled: true\ntrial_details:\n  trial_time_limit: 7",
			want:     true,
		},
		{
			name:     "empty expected value matches a missing key",
			actual:   `title: "Sales data"`,
			expected: "title: Sales data\nsubtitle:",
			want:     true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := organizationListingManifestContains(tc.actual, tc.expected)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("invalid manifest", func(t *testing.T) {
		_, err := organizationListingManifestContains("title: [unclosed", organizationListingManifestFromConfig)
		assert.ErrorContains(t, err, "unable to parse the listing manifest")
	})
}

func Test_organizationListingManifestsEquivalent(t *testing.T) {
	testCases := []struct {
		name  string
		first string
		other string
		want  bool
	}{
		{
			name:  "same manifest with different quoting and key order",
			first: "title: \"Sales data\"\ndescription: \"Daily sales\"",
			other: "description: Daily sales\ntitle: Sales data",
			want:  true,
		},
		{
			name:  "manifest with keys added by Snowflake",
			first: organizationListingManifestFromDescribe,
			other: organizationListingManifestFromConfig,
			want:  false,
		},
		{
			name:  "removed key",
			first: "title: Sales data\ndescription: Daily sales",
			other: "title: Sales data",
			want:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := organizationListingManifestsEquivalent(tc.first, tc.other)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_suppressEquivalentOrganizationListingManifest(t *testing.T) {
	assert.True(t, suppressEquivalentOrganizationListingManifest("", "title: \"Sales data\"", "title: Sales data", nil))
	assert.False(t, suppressEquivalentOrganizationListingManifest("", "title: \"Sales data\"", "title: Sales data v2", nil))
	assert.False(t, suppressEquivalentOrganizationListingManifest("", "", "title: Sales data", nil))
	assert.False(t, suppressEquivalentOrganizationListingManifest("", "title: [unclosed", "title: Sales data", nil))
}
