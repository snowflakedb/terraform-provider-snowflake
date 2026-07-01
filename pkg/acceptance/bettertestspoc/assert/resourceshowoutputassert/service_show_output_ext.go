package resourceshowoutputassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *ServiceShowOutputAssert) HasDnsNameNotEmpty() *ServiceShowOutputAssert {
	s.ValuePresent("dns_name")
	return s
}

func (s *ServiceShowOutputAssert) HasExternalAccessIntegrations(expected ...sdk.AccountObjectIdentifier) *ServiceShowOutputAssert {
	s.StringValueSet("external_access_integrations.#", fmt.Sprintf("%d", len(expected)))
	for _, v := range expected {
		s.SetContainsElem("external_access_integrations", v.Name())
	}
	return s
}

func (s *ServiceShowOutputAssert) HasCreatedOnNotEmpty() *ServiceShowOutputAssert {
	s.ValuePresent("created_on")
	return s
}

func (s *ServiceShowOutputAssert) HasUpdatedOnNotEmpty() *ServiceShowOutputAssert {
	s.ValuePresent("updated_on")
	return s
}

func (s *ServiceShowOutputAssert) HasResumedOnEmpty() *ServiceShowOutputAssert {
	s.StringValueSet("resumed_on", "")
	return s
}

func (s *ServiceShowOutputAssert) HasSuspendedOnEmpty() *ServiceShowOutputAssert {
	s.StringValueSet("suspended_on", "")
	return s
}

func (s *ServiceShowOutputAssert) HasSpecDigestNotEmpty() *ServiceShowOutputAssert {
	s.ValuePresent("spec_digest")
	return s
}

func (s *ServiceShowOutputAssert) HasManagingObjectDomainEmpty() *ServiceShowOutputAssert {
	s.StringValueSet("managing_object_domain", "")
	return s
}

func (s *ServiceShowOutputAssert) HasManagingObjectNameEmpty() *ServiceShowOutputAssert {
	s.StringValueSet("managing_object_name", "")
	return s
}

func (s *ServiceShowOutputAssert) HasQueryWarehouseEmpty() *ServiceShowOutputAssert {
	s.StringValueSet("query_warehouse", "")
	return s
}

func (s *ServiceShowOutputAssert) HasCommentEmpty() *ServiceShowOutputAssert {
	s.StringValueSet("comment", "")
	return s
}
