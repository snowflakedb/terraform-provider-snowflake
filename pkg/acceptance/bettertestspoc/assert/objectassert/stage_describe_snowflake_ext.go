package objectassert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type StageDetailsAssert struct {
	*assert.SnowflakeObjectAssert[sdk.StageDetails, sdk.SchemaObjectIdentifier]
}

func StageDetails(t *testing.T, id sdk.SchemaObjectIdentifier) *StageDetailsAssert {
	t.Helper()
	return &StageDetailsAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(
			sdk.ObjectType("STAGE_DETAILS"),
			id,
			func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.StageDetails, sdk.SchemaObjectIdentifier] {
				return testClient.Stage.DescribeDetails
			}),
	}
}

func (s *StageDetailsAssert) HasFileFormatCsv(expected sdk.FileFormatCsv) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.FileFormatCsv == nil {
			return fmt.Errorf("expected file format to be CSV; got: nil")
		}
		if !reflect.DeepEqual(*o.FileFormatCsv, expected) {
			return fmt.Errorf("expected file format csv:\n%+v\ngot:\n%+v", expected, *o.FileFormatCsv)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasDirectoryTable(expected sdk.StageDirectoryTable) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if !reflect.DeepEqual(*o.DirectoryTable, expected) {
			return fmt.Errorf("expected directory table:\n%+v\ngot:\n%+v", expected, *o.DirectoryTable)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasDirectoryTableEnable(expected bool) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if o.DirectoryTable.Enable != expected {
			return fmt.Errorf("expected directory table enable: %v; got: %v", expected, o.DirectoryTable.Enable)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasDirectoryTableAutoRefresh(expected bool) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if o.DirectoryTable.AutoRefresh != expected {
			return fmt.Errorf("expected directory table auto refresh: %v; got: %v", expected, o.DirectoryTable.AutoRefresh)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasDirectoryTableNotificationChannel(expected string) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if o.DirectoryTable.DirectoryNotificationChannel != expected {
			return fmt.Errorf("expected directory notification channel: %v; got: %v", expected, o.DirectoryTable.DirectoryNotificationChannel)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasDirectoryTableLastRefreshedOnNotEmpty() *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if o.DirectoryTable.LastRefreshedOn == nil || o.DirectoryTable.LastRefreshedOn.IsZero() {
			return fmt.Errorf("expected directory table last refreshed on to not be empty")
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasDirectoryTableLastRefreshedOnNil() *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if o.DirectoryTable.LastRefreshedOn != nil {
			return fmt.Errorf("expected directory table last refreshed on to be nil; got: %v", *o.DirectoryTable.LastRefreshedOn)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasPrivateLinkUsePrivatelinkEndpoint(expected bool) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.PrivateLink == nil {
			return fmt.Errorf("expected private link to have value; got: nil")
		}
		if o.PrivateLink.UsePrivatelinkEndpoint != expected {
			return fmt.Errorf("expected private link use privatelink endpoint: %v; got: %v", expected, o.PrivateLink.UsePrivatelinkEndpoint)
		}
		return nil
	})
	return s
}
