package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
)

func assertThatObject(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObject(t, objectAssert, testClientHelper())
}
