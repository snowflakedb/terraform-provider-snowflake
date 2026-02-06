package objectassert

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

func StorageIntegrationPropertiesFromObject(t *testing.T, id sdk.AccountObjectIdentifier, properties []sdk.StorageIntegrationProperty) *StorageIntegrationPropertiesAssert {
	t.Helper()
	return &StorageIntegrationPropertiesAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeStorageIntegration, id, &StorageIntegrationPropertiesWrapper{Properties: properties}),
	}
}

func (s *StorageIntegrationPropertiesAssert) HasDetailsCount(expected int) *StorageIntegrationPropertiesAssert {
	s.AddAssertion(func(t *testing.T, o *StorageIntegrationPropertiesWrapper) error {
		t.Helper()
		if len(o.Properties) != expected {
			return fmt.Errorf("expected %d storage integration details; got: %d", expected, len(o.Properties))
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationPropertiesAssert) ContainsPropertyEqualTo(expected sdk.StorageIntegrationProperty) *StorageIntegrationPropertiesAssert {
	s.AddAssertion(func(t *testing.T, o *StorageIntegrationPropertiesWrapper) error {
		t.Helper()
		prop, err := collections.FindFirst(o.Properties, func(prop sdk.StorageIntegrationProperty) bool {
			return prop.Name == expected.Name
		})
		if err != nil {
			return fmt.Errorf("expected storage integration describe output to contain a property row %s", expected.Name)
		}
		if prop.Type != expected.Type || prop.Value != expected.Value || prop.Default != expected.Default {
			return fmt.Errorf("expected prop: %s, got: %s", storageIntegrationPropertyString(expected), storageIntegrationPropertyString(*prop))
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationPropertiesAssert) ContainsNotEmptyPropertyWithTypeAndDefault(name string, _type string, _default string) *StorageIntegrationPropertiesAssert {
	s.AddAssertion(func(t *testing.T, o *StorageIntegrationPropertiesWrapper) error {
		t.Helper()
		prop, err := collections.FindFirst(o.Properties, func(prop sdk.StorageIntegrationProperty) bool {
			return prop.Name == name
		})
		if err != nil {
			return fmt.Errorf("expected storage integration describe output to contain a property row %s", name)
		}
		if prop.Value == "" {
			return fmt.Errorf("expected property %s to have value, was empty", name)
		}
		if prop.Type != _type || prop.Default != _default {
			return fmt.Errorf("expected property %s to have: [type=%s default=%s], got: [type=%s default=%s]", name, _type, _default, prop.Type, prop.Default)
		}
		return nil
	})
	return s
}

func (s *StorageIntegrationPropertiesAssert) ContainsPropertyNotEqualToWithTypeAndDefault(name string, unexpectedValue string, _type string, _default string) *StorageIntegrationPropertiesAssert {
	s.AddAssertion(func(t *testing.T, o *StorageIntegrationPropertiesWrapper) error {
		t.Helper()
		prop, err := collections.FindFirst(o.Properties, func(prop sdk.StorageIntegrationProperty) bool {
			return prop.Name == name
		})
		if err != nil {
			return fmt.Errorf("expected storage integration describe output to contain a property row %s", name)
		}
		if prop.Value == unexpectedValue {
			return fmt.Errorf("expected property %s not to have value %s", name, unexpectedValue)
		}
		if prop.Type != _type || prop.Default != _default {
			return fmt.Errorf("expected property %s to have: [type=%s default=%s], got: [type=%s default=%s]", name, _type, _default, prop.Type, prop.Default)
		}
		return nil
	})
	return s
}

func storageIntegrationPropertyString(p sdk.StorageIntegrationProperty) string {
	return fmt.Sprintf("[name=%s type=%s value=%s default=%s]", p.Name, p.Type, p.Value, p.Default)
}
