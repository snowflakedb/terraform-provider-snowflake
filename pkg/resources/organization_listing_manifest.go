package resources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Snowflake stores these manifest values in upper case, so they are compared case-insensitively.
var organizationListingManifestCaseInsensitiveKeys = []string{"account", "accounts", "roles", "organization_profile"}

const organizationListingManifestAccessRegionsKey = "access_regions"

// organizationListingManifestContains reports whether every key and value of the expected manifest is present in the
// actual manifest. Keys that only exist in the actual manifest are ignored, because DESCRIBE LISTING returns a manifest
// enriched with values set by Snowflake (for example, resharing).
func organizationListingManifestContains(actualManifest string, expectedManifest string) (bool, error) {
	actual, err := parseOrganizationListingManifest(actualManifest)
	if err != nil {
		return false, err
	}
	expected, err := parseOrganizationListingManifest(expectedManifest)
	if err != nil {
		return false, err
	}
	return organizationListingManifestValueContains(actual, expected, nil), nil
}

// organizationListingManifestsEquivalent reports whether both manifests describe the same YAML document, ignoring
// formatting differences like quoting and key order.
func organizationListingManifestsEquivalent(manifest string, otherManifest string) (bool, error) {
	first, err := parseOrganizationListingManifest(manifest)
	if err != nil {
		return false, err
	}
	second, err := parseOrganizationListingManifest(otherManifest)
	if err != nil {
		return false, err
	}
	return organizationListingManifestValueContains(first, second, nil) && organizationListingManifestValueContains(second, first, nil), nil
}

func suppressEquivalentOrganizationListingManifest(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	}
	equivalent, err := organizationListingManifestsEquivalent(oldValue, newValue)
	return err == nil && equivalent
}

func parseOrganizationListingManifest(manifest string) (any, error) {
	var parsed any
	if err := yaml.Unmarshal([]byte(manifest), &parsed); err != nil {
		return nil, fmt.Errorf("unable to parse the listing manifest: %w", err)
	}
	return parsed, nil
}

func organizationListingManifestValueContains(actual any, expected any, path []string) bool {
	switch expectedValue := expected.(type) {
	case map[string]any:
		actualValue, ok := actual.(map[string]any)
		if !ok {
			return false
		}
		for key, expectedNested := range expectedValue {
			actualNested, ok := actualValue[key]
			if !ok {
				if expectedNested == nil {
					continue
				}
				return false
			}
			if !organizationListingManifestValueContains(actualNested, expectedNested, append(slices.Clone(path), key)) {
				return false
			}
		}
		return true
	case []any:
		actualValue, ok := actual.([]any)
		if !ok || len(actualValue) != len(expectedValue) {
			return false
		}
		for i := range expectedValue {
			if !organizationListingManifestValueContains(actualValue[i], expectedValue[i], path) {
				return false
			}
		}
		return true
	case nil:
		return actual == nil
	default:
		if actual == nil {
			return false
		}
		actualString, expectedString := fmt.Sprint(actual), fmt.Sprint(expectedValue)
		if organizationListingManifestPathIsCaseInsensitive(path) {
			return strings.EqualFold(actualString, expectedString)
		}
		return actualString == expectedString
	}
}

func organizationListingManifestPathIsCaseInsensitive(path []string) bool {
	if len(path) == 0 {
		return false
	}
	key := path[len(path)-1]
	if slices.Contains(organizationListingManifestCaseInsensitiveKeys, key) {
		return true
	}
	return key == "name" && len(path) > 1 && path[len(path)-2] == organizationListingManifestAccessRegionsKey
}
