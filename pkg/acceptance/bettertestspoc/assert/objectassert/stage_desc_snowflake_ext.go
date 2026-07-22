package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

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

func (s *StageDetailsAssert) HasLocationUrl(expected []string) *StageDetailsAssert {
	s.AddAssertion(func(t *testing.T, o *sdk.StageDetails) error {
		t.Helper()
		if o.Location == nil {
			return fmt.Errorf("expected stage location to have value; got: nil")
		}
		if len(o.Location.Url) != len(expected) {
			return fmt.Errorf("expected stage location url: %v; got: %v", expected, o.Location.Url)
		}
		for i, u := range o.Location.Url {
			if u != expected[i] {
				return fmt.Errorf("expected stage location url[%d]: %v; got: %v", i, expected[i], u)
			}
		}
		return nil
	})
	return s
}

func (s *StageDetailsAssert) HasLocationAwsAccessPointArn(expected string) *StageDetailsAssert {
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

func (s *StageDetailsAssert) HasCredentialsAwsKeyId(expected string) *StageDetailsAssert {
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
