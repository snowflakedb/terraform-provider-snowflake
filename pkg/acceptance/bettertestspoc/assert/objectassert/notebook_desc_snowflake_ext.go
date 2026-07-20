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

func (n *NotebookDetailsAssert) HasNoDefaultVersionAlias() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.DefaultVersionAlias != nil {
			return fmt.Errorf("expected default version alias to be nil; got: %v", *o.DefaultVersionAlias)
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

func (n *NotebookDetailsAssert) HasNoDefaultVersionGitCommitHash() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.DefaultVersionGitCommitHash != nil {
			return fmt.Errorf("expected default version git commit hash to be nil; got: %v", *o.DefaultVersionGitCommitHash)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoLastVersionAlias() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.LastVersionAlias != nil {
			return fmt.Errorf("expected last version alias to be nil; got: %v", *o.LastVersionAlias)
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

func (n *NotebookDetailsAssert) HasNoLastVersionGitCommitHash() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.LastVersionGitCommitHash != nil {
			return fmt.Errorf("expected last version git commit hash to be nil; got: %v", *o.LastVersionGitCommitHash)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoLiveVersionLocationUri() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.LiveVersionLocationUri != nil {
			return fmt.Errorf("expected live version location uri to be nil; got: %v", *o.LiveVersionLocationUri)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoQueryWarehouse() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.QueryWarehouse != nil {
			return fmt.Errorf("expected query warehouse to be nil; got: %v", o.QueryWarehouse)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoComment() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.Comment != nil {
			return fmt.Errorf("expected comment to be empty")
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoTitle() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.Title != nil {
			return fmt.Errorf("expected title to be nil; got: %v", *o.Title)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoComputePool() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.ComputePool != nil {
			return fmt.Errorf("expected compute pool to be nil; got: %v", o.ComputePool)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoDefaultVersionSourceLocationUri() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.DefaultVersionSourceLocationUri != nil {
			return fmt.Errorf("expected default version source location uri to be nil; got: %v", *o.DefaultVersionSourceLocationUri)
		}
		return nil
	})
	return n
}

func (n *NotebookDetailsAssert) HasNoLastVersionSourceLocationUri() *NotebookDetailsAssert {
	n.AddAssertion(func(t *testing.T, o *sdk.NotebookDetails) error {
		t.Helper()
		if o.LastVersionSourceLocationUri != nil {
			return fmt.Errorf("expected last version source location uri to be nil; got: %v", *o.LastVersionSourceLocationUri)
		}
		return nil
	})
	return n
}
