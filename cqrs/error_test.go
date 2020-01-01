package cqrs

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewUnsupportedCommand(t *testing.T) {
    e := NewUnsupportedCommand(dummyCommand{})
    assert.NotNil(t, e)
    assert.IsType(t, (*UnsupportedCommand)(nil), &e)
    assert.Equal(t, "cqrs.error.unsupportedCommand", e.Error())
}

