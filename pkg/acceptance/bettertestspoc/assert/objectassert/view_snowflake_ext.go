package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (w *ViewAssert) HasCreatedOnNotEmpty() *ViewAssert {
	w.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created on not empty; got: %v", o.CreatedOn)
		}
		return nil
	})
	return w
}

func (v *ViewAssert) HasNonEmptyText() *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		if o.Text == "" {
			return fmt.Errorf("expected non empty text")
		}
		return nil
	})
	return v
}

func (v *ViewAssert) HasNoRowAccessPolicyReferences(client *helpers.TestClient) *ViewAssert {
	return v.hasNoPolicyReference(client, sdk.PolicyKindRowAccessPolicy)
}

func (v *ViewAssert) HasNoAggregationPolicyReferences(client *helpers.TestClient) *ViewAssert {
	return v.hasNoPolicyReference(client, sdk.PolicyKindAggregationPolicy)
}

func (v *ViewAssert) HasNoMaskingPolicyReferences(client *helpers.TestClient) *ViewAssert {
	return v.hasNoPolicyReference(client, sdk.PolicyKindMaskingPolicy)
}

func (v *ViewAssert) HasNoProjectionPolicyReferences(client *helpers.TestClient) *ViewAssert {
	return v.hasNoPolicyReference(client, sdk.PolicyKindProjectionPolicy)
}

func (v *ViewAssert) hasNoPolicyReference(client *helpers.TestClient, kind sdk.PolicyKind) *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		refs, err := client.PolicyReferences.GetPolicyReferences(t, o.ID(), sdk.PolicyEntityDomainView)
		if err != nil {
			return err
		}
		refs = slices.DeleteFunc(refs, func(reference sdk.PolicyReference) bool {
			return reference.PolicyKind != kind
		})
		if len(refs) > 0 {
			return fmt.Errorf("expected no %s policy references; got: %v", kind, refs)
		}
		return nil
	})
	return v
}

func (v *ViewAssert) HasRowAccessPolicyReferences(client *helpers.TestClient, n int) *ViewAssert {
	return v.hasPolicyReference(client, sdk.PolicyKindRowAccessPolicy, n)
}

func (v *ViewAssert) HasAggregationPolicyReferences(client *helpers.TestClient, n int) *ViewAssert {
	return v.hasPolicyReference(client, sdk.PolicyKindAggregationPolicy, n)
}

func (v *ViewAssert) HasMaskingPolicyReferences(client *helpers.TestClient, n int) *ViewAssert {
	return v.hasPolicyReference(client, sdk.PolicyKindMaskingPolicy, n)
}

func (v *ViewAssert) HasProjectionPolicyReferences(client *helpers.TestClient, n int) *ViewAssert {
	return v.hasPolicyReference(client, sdk.PolicyKindProjectionPolicy, n)
}

func (v *ViewAssert) hasPolicyReference(client *helpers.TestClient, kind sdk.PolicyKind, n int) *ViewAssert {
	v.AddAssertion(func(t *testing.T, o *sdk.View) error {
		t.Helper()
		refs, err := client.PolicyReferences.GetPolicyReferences(t, o.ID(), sdk.PolicyEntityDomainView)
		if err != nil {
			return err
		}
		refs = slices.DeleteFunc(refs, func(reference sdk.PolicyReference) bool {
			return reference.PolicyKind != kind
		})
		if len(refs) != n {
			return fmt.Errorf("expected %d %s policy references; got: %d, %v", n, kind, len(refs), refs)
		}
		return nil
	})
	return v
}
