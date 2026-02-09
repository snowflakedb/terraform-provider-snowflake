package resourceshowoutputassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationGcsDescribeOutputAssert) HasAllowedLocations(expected ...sdk.StorageLocation) *StorageIntegrationGcsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("allowed_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	// TODO [next PRs]: these are sets so the order varies
	// for i, v := range expected {
	// 	s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("allowed_locations.%d", i), v.Path))
	// }
	return s
}

func (s *StorageIntegrationGcsDescribeOutputAssert) HasBlockedLocations(expected ...sdk.StorageLocation) *StorageIntegrationGcsDescribeOutputAssert {
	s.AddAssertion(assert.ResourceDescribeOutputValueSet("blocked_locations.#", strconv.FormatInt(int64(len(expected)), 10)))
	// TODO [next PRs]: these are sets so the order varies
	// for i, v := range expected {
	// 	s.AddAssertion(assert.ResourceDescribeOutputValueSet(fmt.Sprintf("blocked_locations.%d", i), v.Path))
	// }
	return s
}
