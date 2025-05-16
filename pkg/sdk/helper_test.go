package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testprofiles"
)

func defaultTestClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func secondaryTestClient(t *testing.T) *Client {
	t.Helper()
	return testClient(t, testprofiles.Secondary)
}

func thirdTestClient(t *testing.T) *Client {
	t.Helper()
	return testClient(t, testprofiles.Third)
}

func fourthTestClient(t *testing.T) *Client {
	t.Helper()
	return testClient(t, testprofiles.Fourth)
}

func testClient(t *testing.T, profile string) *Client {
	t.Helper()

	config, err := ProfileConfig(profile)
	if err != nil {
		t.Skipf("Snowflake %s profile not configured. Must be set in ~/.snowflake/config", profile)
	}
	client, err := NewClient(config)
	if err != nil {
		t.Skipf("Snowflake %s profile not configured. Must be set in ~/.snowflake/config", profile)
	}

	return client
}
