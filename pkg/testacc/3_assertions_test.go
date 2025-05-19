package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func assertThatObject(t *testing.T, objectAssert assert.InPlaceAssertionVerifier) {
	t.Helper()
	assert.AssertThatObject(t, objectAssert, testClient())
}

func assertThatImport(t *testing.T, fs ...assert.ImportStateCheckFuncProvider) resource.ImportStateCheckFunc {
	t.Helper()
	return assert.AssertThatImport(t, testClient(), fs...)
}

func assertThat(t *testing.T, fs ...assert.TestCheckFuncProvider) resource.TestCheckFunc {
	t.Helper()
	return assert.AssertThat(t, testClient(), fs...)
}
