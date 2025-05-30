// Code generated by assertions generator; DO NOT EDIT.

package objectassert

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type ImageRepositoryAssert struct {
	*assert.SnowflakeObjectAssert[sdk.ImageRepository, sdk.SchemaObjectIdentifier]
}

func ImageRepository(t *testing.T, id sdk.SchemaObjectIdentifier) *ImageRepositoryAssert {
	t.Helper()
	return &ImageRepositoryAssert{
		assert.NewSnowflakeObjectAssertWithTestClientObjectProvider(sdk.ObjectTypeImageRepository, id, func(testClient *helpers.TestClient) assert.ObjectProvider[sdk.ImageRepository, sdk.SchemaObjectIdentifier] {
			return testClient.ImageRepository.Show
		}),
	}
}

func ImageRepositoryFromObject(t *testing.T, imageRepository *sdk.ImageRepository) *ImageRepositoryAssert {
	t.Helper()
	return &ImageRepositoryAssert{
		assert.NewSnowflakeObjectAssertWithObject(sdk.ObjectTypeImageRepository, imageRepository.ID(), imageRepository),
	}
}

func (i *ImageRepositoryAssert) HasCreatedOn(expected time.Time) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.CreatedOn != expected {
			return fmt.Errorf("expected created on: %v; got: %v", expected, o.CreatedOn)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasName(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.Name != expected {
			return fmt.Errorf("expected name: %v; got: %v", expected, o.Name)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasDatabaseName(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.DatabaseName != expected {
			return fmt.Errorf("expected database name: %v; got: %v", expected, o.DatabaseName)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasSchemaName(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.SchemaName != expected {
			return fmt.Errorf("expected schema name: %v; got: %v", expected, o.SchemaName)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasRepositoryUrl(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.RepositoryUrl != expected {
			return fmt.Errorf("expected repository url: %v; got: %v", expected, o.RepositoryUrl)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasOwner(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.Owner != expected {
			return fmt.Errorf("expected owner: %v; got: %v", expected, o.Owner)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasOwnerRoleType(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.OwnerRoleType != expected {
			return fmt.Errorf("expected owner role type: %v; got: %v", expected, o.OwnerRoleType)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasComment(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.Comment != expected {
			return fmt.Errorf("expected comment: %v; got: %v", expected, o.Comment)
		}
		return nil
	})
	return i
}

func (i *ImageRepositoryAssert) HasPrivatelinkRepositoryUrl(expected string) *ImageRepositoryAssert {
	i.AddAssertion(func(t *testing.T, o *sdk.ImageRepository) error {
		t.Helper()
		if o.PrivatelinkRepositoryUrl != expected {
			return fmt.Errorf("expected privatelink repository url: %v; got: %v", expected, o.PrivatelinkRepositoryUrl)
		}
		return nil
	})
	return i
}
