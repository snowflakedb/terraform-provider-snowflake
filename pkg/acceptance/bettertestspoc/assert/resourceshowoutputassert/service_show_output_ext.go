package resourceshowoutputassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// ServicesDatasourceShowOutput is a temporary workaround to have better show output assertions in data source acceptance tests.
func ServicesDatasourceShowOutput(t *testing.T, name string) *ServiceShowOutputAssert {
	t.Helper()

	s := ServiceShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "show_output", "services.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (s *ServiceShowOutputAssert) HasDnsNameNotEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("dns_name"))
	return s
}

func (s *ServiceShowOutputAssert) HasExternalAccessIntegrations(expected ...sdk.AccountObjectIdentifier) *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("external_access_integrations.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ResourceShowOutputValueSet(fmt.Sprintf("external_access_integrations.%d", i), v.Name()))
	}
	return s
}

func (a *ServiceShowOutputAssert) HasCreatedOnNotEmpty() *ServiceShowOutputAssert {
	a.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return a
}

func (s *ServiceShowOutputAssert) HasUpdatedOnNotEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("updated_on"))
	return s
}

func (s *ServiceShowOutputAssert) HasResumedOnEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("resumed_on", ""))
	return s
}

func (s *ServiceShowOutputAssert) HasSuspendedOnEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("suspended_on", ""))
	return s
}

func (s *ServiceShowOutputAssert) HasSpecDigestNotEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("spec_digest"))
	return s
}

func (s *ServiceShowOutputAssert) HasManagingObjectDomainEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("managing_object_domain", ""))
	return s
}

func (s *ServiceShowOutputAssert) HasManagingObjectNameEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("managing_object_name", ""))
	return s
}

func (s *ServiceShowOutputAssert) HasQueryWarehouseEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("query_warehouse", ""))
	return s
}

func (s *ServiceShowOutputAssert) HasCommentEmpty() *ServiceShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("comment", ""))
	return s
}
