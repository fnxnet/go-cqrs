package cqrs

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewFacade(t *testing.T) {
    f := NewFacade(nil)
    assert.IsType(t, (*facade)(nil), &f)
    assert.IsType(t, (*eventBus)(nil), f.EventBus())
    assert.IsType(t, (*commandBus)(nil), f.CommandBus())
}

func TestNewFacadeWithConfig(t *testing.T) {
    type DummyConfig struct{}

    f := NewFacade(DummyConfig{})
    assert.IsType(t, (*facade)(nil), &f)
}
