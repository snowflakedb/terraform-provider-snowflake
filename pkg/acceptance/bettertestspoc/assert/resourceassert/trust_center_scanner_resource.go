package resourceassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

type TrustCenterScannerResourceAssert struct {
	*assert.ResourceAssert
}

func TrustCenterScannerResource(t *testing.T, name string) *TrustCenterScannerResourceAssert {
	t.Helper()

	return &TrustCenterScannerResourceAssert{
		ResourceAssert: assert.NewResourceAssert(name, "resource"),
	}
}

func ImportedTrustCenterScannerResource(t *testing.T, id string) *TrustCenterScannerResourceAssert {
	t.Helper()

	return &TrustCenterScannerResourceAssert{
		ResourceAssert: assert.NewImportedResourceAssert(id, "imported resource"),
	}
}

func (t *TrustCenterScannerResourceAssert) HasScannerPackageId(expected string) *TrustCenterScannerResourceAssert {
	t.StringValueSet("scanner_package_id", expected)
	return t
}

func (t *TrustCenterScannerResourceAssert) HasScannerId(expected string) *TrustCenterScannerResourceAssert {
	t.StringValueSet("scanner_id", expected)
	return t
}

func (t *TrustCenterScannerResourceAssert) HasEnabled(expected string) *TrustCenterScannerResourceAssert {
	t.AddAssertion(assert.ValueSet("enabled", expected))
	return t
}

func (t *TrustCenterScannerResourceAssert) HasSchedule(expected string) *TrustCenterScannerResourceAssert {
	t.StringValueSet("schedule", expected)
	return t
}

func (t *TrustCenterScannerResourceAssert) HasScannerPackageIdNotEmpty() *TrustCenterScannerResourceAssert {
	t.AddAssertion(assert.ValuePresent("scanner_package_id"))
	return t
}

func (t *TrustCenterScannerResourceAssert) HasScannerIdNotEmpty() *TrustCenterScannerResourceAssert {
	t.AddAssertion(assert.ValuePresent("scanner_id"))
	return t
}

func (t *TrustCenterScannerResourceAssert) HasEnabledNotEmpty() *TrustCenterScannerResourceAssert {
	t.AddAssertion(assert.ValuePresent("enabled"))
	return t
}
