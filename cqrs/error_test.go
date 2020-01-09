package cqrs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnsupportedCommand(t *testing.T) {
	e := NewUnsupportedItem(dummyItem{})
	assert.NotNil(t, e)
	assert.IsType(t, (*UnsupportedItem)(nil), &e)
	assert.Equal(t, "cqrs.error.unsupportedItem", e.Error())
}
