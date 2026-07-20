package assert

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

type (
	assertSdk[T any]                                        func(*testing.T, T) error
	ObjectProvider[T any, I sdk.ObjectIdentifier]           func(*testing.T, I) (*T, error)
	testClientObjectProvider[T any, I sdk.ObjectIdentifier] func(client *helpers.TestClient) ObjectProvider[T, I]
)

// SnowflakeObjectAssert is an embeddable struct that should be used to construct new Snowflake object assertions.
// It implements both TestCheckFuncProvider and ImportStateCheckFuncProvider which makes it easy to create new resource assertions.
type SnowflakeObjectAssert[T any, I sdk.ObjectIdentifier] struct {
	assertions               []assertSdk[*T]
	id                       I
	objectType               sdk.ObjectType
	object                   *T
	provider                 ObjectProvider[T, I]
	testClientObjectProvider testClientObjectProvider[T, I]
}

// NewSnowflakeObjectAssertWithTestClientObjectProvider creates a SnowflakeObjectAssert with id and the test client-varying parameters provider.
// Object to check is lazily fetched from Snowflake when the checks are being run.
func NewSnowflakeObjectAssertWithTestClientObjectProvider[T any, I sdk.ObjectIdentifier](objectType sdk.ObjectType, id I, testClientObjectProvider testClientObjectProvider[T, I]) *SnowflakeObjectAssert[T, I] {
	return &SnowflakeObjectAssert[T, I]{
		assertions:               make([]assertSdk[*T], 0),
		id:                       id,
		objectType:               objectType,
		testClientObjectProvider: testClientObjectProvider,
	}
}

// NewSnowflakeObjectAssertWithObject creates a SnowflakeObjectAssert with object that was already fetched from Snowflake.
// All the checks are run against the given object.
func NewSnowflakeObjectAssertWithObject[T any, I sdk.ObjectIdentifier](objectType sdk.ObjectType, id I, object *T) *SnowflakeObjectAssert[T, I] {
	return &SnowflakeObjectAssert[T, I]{
		assertions: make([]assertSdk[*T], 0),
		id:         id,
		objectType: objectType,
		object:     object,
	}
}

// SubStructAssert is used exclusively for composition via the nested assertion adapter.
// It intentionally does NOT implement VerifyAll, ToTerraformTestCheckFunc, or
// ToTerraformImportStateCheckFunc — sub-struct asserters cannot be used standalone.
type SubStructAssert[T any] struct {
	assertions []func(*testing.T, *T) error
}

func NewSubStructAssert[T any]() *SubStructAssert[T] {
	return &SubStructAssert[T]{assertions: make([]func(*testing.T, *T) error, 0)}
}

func (s *SubStructAssert[T]) AddAssertion(assertion func(*testing.T, *T) error) {
	s.assertions = append(s.assertions, assertion)
}

func (s *SubStructAssert[T]) GetAssertions() []func(*testing.T, *T) error {
	return s.assertions
}

func (s *SnowflakeObjectAssert[T, I]) AddAssertion(assertion assertSdk[*T]) {
	s.assertions = append(s.assertions, assertion)
}

// AssertionOnPointerField checks that field is non-nil, then runs assertion against it.
// Used by the nested assertion adapter for optional (pointer) sub-struct fields.
// Returns an error if field is nil so that all sub-assertions consistently report the nil case.
func AssertionOnPointerField[F any](t *testing.T, field *F, fieldName string, assertion func(*testing.T, *F) error) error {
	t.Helper()
	if field == nil {
		return fmt.Errorf("expected %s to have value; got: nil", fieldName)
	}
	return assertion(t, field)
}

func (s *SnowflakeObjectAssert[T, I]) GetId() I {
	return s.id
}

// ToTerraformTestCheckFunc implements TestCheckFuncProvider to allow easier creation of new Snowflake object assertions.
// It goes through all the assertion accumulated earlier and gathers the results of the checks.
func (s *SnowflakeObjectAssert[_, _]) ToTerraformTestCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.TestCheckFunc {
	t.Helper()
	return func(_ *terraform.State) error {
		return s.runSnowflakeObjectsAssertions(t, testClient)
	}
}

// ToTerraformImportStateCheckFunc implements ImportStateCheckFuncProvider to allow easier creation of new Snowflake object assertions.
// It goes through all the assertion accumulated earlier and gathers the results of the checks.
func (s *SnowflakeObjectAssert[_, _]) ToTerraformImportStateCheckFunc(t *testing.T, testClient *helpers.TestClient) resource.ImportStateCheckFunc {
	t.Helper()
	return func(_ []*terraform.InstanceState) error {
		return s.runSnowflakeObjectsAssertions(t, testClient)
	}
}

// VerifyAll implements InPlaceAssertionVerifier to allow easier creation of new Snowflake object assertions.
// It verifies all the assertions accumulated earlier and gathers the results of the checks.
func (s *SnowflakeObjectAssert[_, _]) VerifyAll(t *testing.T, testClient *helpers.TestClient) {
	t.Helper()
	err := s.runSnowflakeObjectsAssertions(t, testClient)
	require.NoError(t, err)
}

func (s *SnowflakeObjectAssert[T, _]) runSnowflakeObjectsAssertions(t *testing.T, testClient *helpers.TestClient) error {
	t.Helper()

	var sdkObject *T
	var err error
	switch {
	case s.object != nil:
		sdkObject = s.object
	case s.testClientObjectProvider != nil:
		if testClient == nil {
			return errors.New("testClient must not be nil")
		}
		sdkObject, err = s.testClientObjectProvider(testClient)(t, s.id)
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
