package sdk

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errObjectNotFoundOnNoRows(t *testing.T) {
	assert.ErrorIs(t, errObjectNotFoundOnNoRows(sql.ErrNoRows), ErrObjectNotFound)
	assert.NotErrorIs(t, errObjectNotFoundOnNoRows(sql.ErrNoRows), sql.ErrNoRows)

	other := errors.New("boom")
	assert.ErrorIs(t, errObjectNotFoundOnNoRows(other), other)
	assert.NotErrorIs(t, errObjectNotFoundOnNoRows(other), ErrObjectNotFound)
}
