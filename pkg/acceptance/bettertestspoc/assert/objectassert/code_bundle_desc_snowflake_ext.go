package objectassert

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *CodeBundleDetailsAssert) HasDefaultVersionSourceLocationUriContaining(substring string) *CodeBundleDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CodeBundleDetails) error {
		t.Helper()
		if o.DefaultVersionSourceLocationUri == nil {
			return fmt.Errorf("expected default version source location uri to have value; got: nil")
		}
		if !strings.Contains(*o.DefaultVersionSourceLocationUri, substring) {
			return fmt.Errorf("expected default version source location uri to contain %q; got: %s", substring, *o.DefaultVersionSourceLocationUri)
		}
		return nil
	})
	return c
}

func (c *CodeBundleDetailsAssert) HasLastVersionNameNotEmpty() *CodeBundleDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CodeBundleDetails) error {
		t.Helper()
		if o.LastVersionName == nil {
			return fmt.Errorf("expected last version name to have value; got: nil")
		}
		if *o.LastVersionName == "" {
			return fmt.Errorf("expected last version name to be non-empty")
		}
		return nil
	})
	return c
}

func (c *CodeBundleDetailsAssert) HasLastVersionSourceLocationUriContaining(substring string) *CodeBundleDetailsAssert {
	c.AddAssertion(func(t *testing.T, o *sdk.CodeBundleDetails) error {
		t.Helper()
		if o.LastVersionSourceLocationUri == nil {
			return fmt.Errorf("expected last version source location uri to have value; got: nil")
		}
		if !strings.Contains(*o.LastVersionSourceLocationUri, substring) {
			return fmt.Errorf("expected last version source location uri to contain %q; got: %s", substring, *o.LastVersionSourceLocationUri)
		}
		return nil
	})
	return c
}
