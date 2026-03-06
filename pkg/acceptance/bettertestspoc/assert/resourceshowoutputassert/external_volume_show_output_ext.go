package resourceshowoutputassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func (e *ExternalVolumeShowOutputAssert) HasCommentEmpty() *ExternalVolumeShowOutputAssert {
	e.AddAssertion(assert.ResourceShowOutputValueSet("comment", ""))
	return e
}
