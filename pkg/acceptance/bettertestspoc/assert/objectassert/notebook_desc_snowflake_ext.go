package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (n *NotebookDetailsAssert) HasUrlIdNotEmpty() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.UrlId == "" {
			return fmt.Errorf("expected url id not empty; got empty")
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNonEmptyDefaultPackages() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.DefaultPackages == "" {
			return fmt.Errorf("expected default packages not empty; got empty")
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNonEmptyDefaultVersionLocationUri() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.DefaultVersionLocationUri == "" {
			return fmt.Errorf("expected default version location uri not empty; got empty")
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNonEmptyLastVersionLocationUri() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.LastVersionLocationUri == "" {
			return fmt.Errorf("expected last version location uri not empty; got empty")
		}
		return nil
	})
	return n
}
