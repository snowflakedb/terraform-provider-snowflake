package sdk

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

func RandomUUID(t *testing.T) string {
	t.Helper()
	v, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return v
}

func RandomComment(t *testing.T) string {
	t.Helper()
	return gofakeit.Sentence(10)
}

func RandomBool(t *testing.T) bool {
	t.Helper()
	return gofakeit.Bool()
}

func RandomString(t *testing.T) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, 28)
}

func RandomStringN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, num)
}

func RandomAlphanumericN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, false, false, num)
}

func RandomStringRange(t *testing.T, min, max int) string {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.Password(true, true, true, true, false, RandomIntRange(t, min, max))
}

func RandomIntRange(t *testing.T, min, max int) int {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.IntRange(min, max)
}

func RandomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	return NewSchemaObjectIdentifier(RandomStringN(t, 12), RandomStringN(t, 12), RandomStringN(t, 12))
}

func RandomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	return NewDatabaseObjectIdentifier(RandomStringN(t, 12), RandomStringN(t, 12))
}

func AlphanumericDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
	t.Helper()
	return NewDatabaseObjectIdentifier(RandomAlphanumericN(t, 12), RandomAlphanumericN(t, 12))
}

func RandomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	return NewAccountObjectIdentifier(RandomStringN(t, 12))
}
