package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func assertThatObject(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObject(t, objectAssert, testClientHelper())
}

func assertThatObjectSecondary(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObject(t, objectAssert, secondaryTestClientHelper())
}
