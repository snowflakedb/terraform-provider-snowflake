package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (l *ListingAssert) HasNoReviewState() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.ReviewState != nil {
			return fmt.Errorf("expected review state to be nil, but got %q", *o.ReviewState)
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasNoComment() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be nil, but got %q", *o.Comment)
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasNoRegions() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.Regions != nil {
			return fmt.Errorf("expected regions to be nil, but got %q", *o.Regions)
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasNoUniformListingLocator() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.UniformListingLocator != nil {
			return fmt.Errorf("expected uniform listing locator to be nil, but got %q", *o.UniformListingLocator)
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasNoDetailedTargetAccounts() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.DetailedTargetAccounts != nil {
			return fmt.Errorf("expected detailed target accounts to be nil, but got %q", *o.DetailedTargetAccounts)
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasNoPublishedOn() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.PublishedOn != nil {
			return fmt.Errorf("expected published_on to be nil, but got %q", *o.PublishedOn)
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasGlobalNameNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.GlobalName == "" {
			return fmt.Errorf("expected global_name to be not empty")
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasCreatedOnNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.CreatedOn == "" {
			return fmt.Errorf("expected created_on to be not empty")
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasUpdatedOnNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if o.UpdatedOn == "" {
			return fmt.Errorf("expected updated_on to be not empty")
		}
		return nil
	})
	return l
}

func (l *ListingAssert) HasDetailedTargetAccountsNotEmpty() *ListingAssert {
	l.AddAssertion(func(t *testing.T, o *sdk.Listing) error {
		t.Helper()
		if *o.DetailedTargetAccounts == "" {
			return fmt.Errorf("expected detailed_target_accounts to be not empty")
		}
		return nil
	})
	return l
}
