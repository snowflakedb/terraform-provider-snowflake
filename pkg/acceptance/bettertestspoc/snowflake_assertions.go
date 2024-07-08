package bettertestspoc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

type assertSdk[T any] func(*testing.T, T) error
type objectProvider[T any, I sdk.ObjectIdentifier] func(*testing.T, I) (*T, error)

type SnowflakeObjectAssert[T any, I sdk.ObjectIdentifier] struct {
	assertions []assertSdk[*T]
	id         I
	objectType sdk.ObjectType
	object     *T
	provider   objectProvider[T, I]
	TestCheckFuncProvider
}

func NewSnowflakeObjectAssertWithProvider[T any, I sdk.ObjectIdentifier](objectType sdk.ObjectType, id I, provider objectProvider[T, I]) *SnowflakeObjectAssert[T, I] {
	return &SnowflakeObjectAssert[T, I]{
		assertions: make([]assertSdk[*T], 0),
		id:         id,
		objectType: objectType,
		provider:   provider,
	}
}

func NewSnowflakeObjectAssertWithObject[T any, I sdk.ObjectIdentifier](objectType sdk.ObjectType, id I, object *T) *SnowflakeObjectAssert[T, I] {
	return &SnowflakeObjectAssert[T, I]{
		assertions: make([]assertSdk[*T], 0),
		id:         id,
		objectType: objectType,
		object:     object,
	}
}

func (s *SnowflakeObjectAssert[_, _]) ToTerraformTestCheckFunc(t *testing.T) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		return s.runSnowflakeObjectsAssertions(t)
	}
}

func (s *SnowflakeObjectAssert[_, _]) ToTerraformImportStateCheckFunc(t *testing.T) resource.ImportStateCheckFunc {
	t.Helper()
	return func(_ []*terraform.InstanceState) error {
		return s.runSnowflakeObjectsAssertions(t)
	}
}

func (s *SnowflakeObjectAssert[T, _]) runSnowflakeObjectsAssertions(t *testing.T) error {
	var sdkObject *T
	var err error
	switch {
	case s.object != nil:
		sdkObject = s.object
	case s.provider != nil:
		sdkObject, err = s.provider(t, s.id)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot proceed with object %s[%s] assertion: object or provider must be specified", s.objectType, s.id.FullyQualifiedName())
	}

	var result []error

	for i, assertion := range s.assertions {
		if err = assertion(t, sdkObject); err != nil {
			result = append(result, fmt.Errorf("object %s[%s] assertion [%d/%d]: failed with error: %w", s.objectType, s.id.FullyQualifiedName(), i+1, len(s.assertions), err))
		}
	}

	return errors.Join(result...)
}

func AssertThatObject[T any, I sdk.ObjectIdentifier](t *testing.T, objectAssert *SnowflakeObjectAssert[T, I]) {
	t.Helper()
	err := objectAssert.runSnowflakeObjectsAssertions(t)
	require.NoError(t, err)
}

func (s *SnowflakeObjectAssert[_, _]) CheckAll(t *testing.T) {
	t.Helper()
	err := s.runSnowflakeObjectsAssertions(t)
	require.NoError(t, err)
}
