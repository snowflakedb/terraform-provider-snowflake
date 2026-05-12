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

func (s *StageDetailsAssert) HasFileFormatName(expected sdk.SchemaObjectIdentifier) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.FileFormatName == nil {
			return fmt.Errorf("expected file format name to have value; got: nil")
		}
		if !reflect.DeepEqual(*o.FileFormatName, expected) {
			return fmt.Errorf("expected file format name: %v; got: %v", expected, *o.FileFormatName)
		}
		return nil
	})
	return s
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

func (s *StageDetailsAssert) HasDirectoryTableNotificationChannelEmpty() *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.DirectoryTable == nil {
			return fmt.Errorf("expected directory table to have value; got: nil")
		}
		if o.DirectoryTable.DirectoryNotificationChannel != nil {
			return fmt.Errorf("expected directory table notification channel to be nil; got: %v", *o.DirectoryTable.DirectoryNotificationChannel)
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
		if o.DirectoryTable.DirectoryNotificationChannel == nil {
			return fmt.Errorf("expected directory table notification channel to have value; got: nil")
		}
		if *o.DirectoryTable.DirectoryNotificationChannel != expected {
			return fmt.Errorf("expected directory notification channel: %v; got: %v", expected, *o.DirectoryTable.DirectoryNotificationChannel)
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
		if o.DirectoryTable.LastRefreshedOn == nil || *o.DirectoryTable.LastRefreshedOn == "" {
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

// Location assertions

func (s *StageDetailsAssert) HasStageLocation(expected sdk.StageLocationDetails) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Location == nil {
			return fmt.Errorf("expected stage location to have value; got: nil")
		}
		if !reflect.DeepEqual(*o.Location, expected) {
			return fmt.Errorf("expected stage location:\n%+v\ngot:\n%+v", expected, *o.Location)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasStageLocationUrl(expected []string) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Location == nil {
			return fmt.Errorf("expected stage location to have value; got: nil")
		}
		if !reflect.DeepEqual(o.Location.Url, expected) {
			return fmt.Errorf("expected stage location url:\n%+v\ngot:\n%+v", expected, o.Location.Url)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasStageLocationAwsAccessPointArn(expected string) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Location == nil {
			return fmt.Errorf("expected stage location to have value; got: nil")
		}
		if o.Location.AwsAccessPointArn != expected {
			return fmt.Errorf("expected stage location aws access point arn: %v; got: %v", expected, o.Location.AwsAccessPointArn)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasStageLocationNil() *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Location != nil {
			return fmt.Errorf("expected stage location to be nil; got: %+v", *o.Location)
		}
		return nil
	})
	return s
}

// Credentials assertions

func (s *StageDetailsAssert) HasStageCredentials(expected sdk.StageCredentials) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Credentials == nil {
			return fmt.Errorf("expected stage credentials to have value; got: nil")
		}
		if !reflect.DeepEqual(*o.Credentials, expected) {
			return fmt.Errorf("expected stage credentials:\n%+v\ngot:\n%+v", expected, *o.Credentials)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasStageCredentialsAwsKeyId(expected string) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Credentials == nil {
			return fmt.Errorf("expected stage credentials to have value; got: nil")
		}
		if o.Credentials.AwsKeyId != expected {
			return fmt.Errorf("expected stage credentials aws key id: %v; got: %v", expected, o.Credentials.AwsKeyId)
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasStageCredentialsNil() *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Credentials != nil {
			return fmt.Errorf("expected stage credentials to be nil; got: %+v", *o.Credentials)
		}
		return nil
	})
	return s
}
