package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type TrustCenterScannerPackageResourceAssert struct {
	*assert.ResourceAssert
}

func TrustCenterScannerPackageResource(t *testing.T, name string) *TrustCenterScannerPackageResourceAssert {
	t.Helper()

	return &TrustCenterScannerPackageResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedTrustCenterScannerPackageResource(t *testing.T, id string) *TrustCenterScannerPackageResourceAssert {
	t.Helper()

	return &TrustCenterScannerPackageResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

func (t *TrustCenterScannerPackageResourceAssert) HasScannerPackageId(expected string) *TrustCenterScannerPackageResourceAssert {
	t.StringValueSet("scanner_package_id", expected)
	return t
}

func (t *TrustCenterScannerPackageResourceAssert) HasEnabled(expected string) *TrustCenterScannerPackageResourceAssert {
	t.AddAssertion(assert.ValueSet("enabled", expected))
	return t
}

func (t *TrustCenterScannerPackageResourceAssert) HasSchedule(expected string) *TrustCenterScannerPackageResourceAssert {
	t.StringValueSet("schedule", expected)
	return t
}

func (t *TrustCenterScannerPackageResourceAssert) HasScannerPackageIdNotEmpty() *TrustCenterScannerPackageResourceAssert {
	t.AddAssertion(assert.ValuePresent("scanner_package_id"))
	return t
}

func (t *TrustCenterScannerPackageResourceAssert) HasEnabledNotEmpty() *TrustCenterScannerPackageResourceAssert {
	t.AddAssertion(assert.ValuePresent("enabled"))
	return t
}
