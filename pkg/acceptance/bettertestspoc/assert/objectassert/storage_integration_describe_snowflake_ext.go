package objectassert

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type StorageIntegrationPropertiesWrapper struct {
	Properties []sdk.StorageIntegrationProperty
}

type StorageIntegrationPropertiesAssert struct {
	*assert.SnowflakeObjectAssert[StorageIntegrationPropertiesWrapper, sdk.AccountObjectIdentifier]
}

func StorageIntegrationProperties(t *testing.T, id sdk.AccountObjectIdentifier) *StorageIntegrationPropertiesAssert {
	t.Helper()
	return &StorageIntegrationPropertiesAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeStorageIntegration, id, func(testClient *helpers.TestClient) assert.ObjectProvider[StorageIntegrationPropertiesWrapper, sdk.AccountObjectIdentifier] {
			return func(t *testing.T, id sdk.AccountObjectIdentifier) (*StorageIntegrationPropertiesWrapper, error) {
				t.Helper()
				properties, err := testClient.StorageIntegration.Describe(t, id)
				if err != nil {
					return nil, err
				}
				return &StorageIntegrationPropertiesWrapper{Properties: properties}, nil
			}
		}),
	}
}

func (s *StorageIntegrationPropertiesAssert) ContainsProperty(expected sdk.StorageIntegrationProperty) *StorageIntegrationPropertiesAssert {
	s.AddAssertion(func(t *testing.T, o *StorageIntegrationPropertiesWrapper) error {
		t.Helper()
		found := slices.ContainsFunc(o.Properties, func(prop sdk.StorageIntegrationProperty) bool {
			return prop.Name == expected.Name &&
				prop.Type == expected.Type &&
				prop.Value == expected.Value &&
				prop.Default == expected.Default
		})
		if !found {
			return fmt.Errorf("expected storage integration describe output to contain a property row matching %s", storageIntegrationPropertyString(expected))
		}
		return nil
	})
	return s
}

func storageIntegrationPropertyString(p sdk.StorageIntegrationProperty) string {
	return fmt.Sprintf("[name=%s type=%s value=%s default=%s]", p.Name, p.Type, p.Value, p.Default)
}
