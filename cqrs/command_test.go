package cqrs

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
)

type dummyEventBus struct {
    events []Event
}

func (d *dummyEventBus) Emit(events ... Event) {
    d.events = append(d.events, events...)
}

func (dummyEventBus) Register(handler EventHandler, events ... Event) {}

func (dummyEventBus) AddMiddleware(middleware EventMiddleware) {}

type dummyCommand struct{}

func (dummyCommand) CommandName() string {
    return "dummyCommandName"
}

type dummyHandler struct {
    Name           string
    Handled        bool
    HandledCommand interface{}
}

func (d *dummyHandler) Handle(i interface{}, bus EventBus) (err error) {
    d.Handled = true
    d.HandledCommand = i
    return
}

type dummyHandlerWithError struct {
    Name           string
    Handled        bool
    HandledCommand interface{}
}

func (d *dummyHandlerWithError) Handle(i interface{}, bus EventBus) (err error) {
    d.Handled = true
    d.HandledCommand = i
    return errors.New("error")
}

type dummyRegistry struct{}

func (d *dummyRegistry) InvalidCommand() {}

func (d *dummyRegistry) InvalidCommandWithParam(i int) {}

func (d *dummyRegistry) InvalidCommandWithReturnInt() (int) {
    return 0
}

func (d *dummyRegistry) InvalidCommandWithReturnFunc() (func()) {
    return func() {}
}

func (d *dummyRegistry) InvalidCommandWithReturnFuncInt() (func() int) {
    return func() int { return 0}
}

func (d *dummyRegistry) InvalidCommandWithReturnAndName() (e int) {
    return 0
}

func (d *dummyRegistry) InvalidCommandWithReturnAndParam(i int) (int) {
    return 0
}

func (d *dummyRegistry) InvalidCommandWithProviderReturned(i int) (f Provider) {
    f = func() CommandHandler { return &dummyHandler{Name: "invalidHandler"} }
    return
}

func (d *dummyRegistry) InvalidCommandWithMultipleOutputs() (f Provider, e error) {
    f = func() CommandHandler { return &dummyHandler{Name: "invalidHandler"} }
    return
}

func (d *dummyRegistry) InvalidCommandWithMultipleUnnamedOutputs() (Provider, error) {
    return func() CommandHandler { return &dummyHandler{Name: "invalidHandler"} }, nil
}

func (d *dummyRegistry) DummyCommand() func() CommandHandler {
    return func() CommandHandler { return &dummyHandler{Name: "dummyCommandHandler"} }
}

func (d *dummyRegistry) DummyCommandProvider() Provider {
    return func() CommandHandler { return &dummyHandler{Name: "dummyCommandHandler2"} }
}

func TestCQRSCommand_CommandName(t *testing.T) {
    type DummyCommand struct {
        CQRSCommand
    }

    c := DummyCommand{CQRSCommand{CommandName: "test"}}

    assert.Equal(t, "test", c.Name())
}

func TestNewCommandBus(t *testing.T) {
    assert.IsType(t, (*commandBus)(nil), NewCommandBus(&dummyEventBus{}))
}

func TestCommandBus_Add(t *testing.T) {
    def := CommandDefinition{Name: "dummy", Provider: func() CommandHandler { return &dummyHandler{} }}
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    bus.Add(def)

    assert.Len(t, bus.handlers, 1)
    assert.Contains(t, bus.handlers, "dummy")
    assert.NotNil(t, bus.handlers["dummy"].Provider)
}

func TestCommandBus_AddHandler(t *testing.T) {
    h := &dummyHandler{}
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    bus.AddHandler(dummyCommand{}, h)

    assert.Len(t, bus.handlers, 1)
    definition := bus.handlers["cqrs.dummyCommand"]
    assert.IsType(t, (*CommandDefinition)(nil), &definition)
    assert.NotNil(t, definition.CommandHandler)
    assert.Nil(t, definition.Provider)
    assert.Equal(t, h, definition.CommandHandler)
    assert.Equal(t, "cqrs.dummyCommand", definition.Name)
}

func TestCommandBus_AddProvider(t *testing.T) {
    h := &dummyHandler{}
    f := func() CommandHandler {
        return h
    }

    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    bus.AddProvider(dummyCommand{}, f)

    assert.Len(t, bus.handlers, 1)
    definition := bus.handlers["cqrs.dummyCommand"]
    assert.IsType(t, (*CommandDefinition)(nil), &definition)
    assert.NotNil(t, definition.Provider)
    assert.IsType(t, (Provider)(nil), definition.Provider)
    assert.Equal(t, h, definition.Provider())
    assert.Nil(t, definition.CommandHandler)
    assert.Equal(t, "cqrs.dummyCommand", definition.Name)
}

func TestCommandBus_AddRegistry(t *testing.T) {
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    //bus.AddRegistry(dummyRegistry{})
    bus.AddRegistry(&dummyRegistry{})
    assert.Len(t, bus.handlers, 2)
    if assert.Contains(t, bus.handlers, "cqrs.DummyCommand") {
        definition := bus.handlers["cqrs.DummyCommand"]
        assert.IsType(t, (*CommandDefinition)(nil), &definition)
        assert.NotNil(t, definition.Provider)
        assert.IsType(t, (Provider)(nil), definition.Provider)
        assert.Equal(t, "dummyCommandHandler", definition.Provider().(*dummyHandler).Name)
        assert.Nil(t, definition.CommandHandler)
        assert.Equal(t, "cqrs.DummyCommand", definition.Name)
    }

    if assert.Contains(t, bus.handlers, "cqrs.DummyCommandProvider") {
        definition := bus.handlers["cqrs.DummyCommandProvider"]
        assert.IsType(t, (*CommandDefinition)(nil), &definition)
        assert.NotNil(t, definition.Provider)
        assert.IsType(t, (Provider)(nil), definition.Provider)
        assert.Equal(t, "dummyCommandHandler2", definition.Provider().(*dummyHandler).Name)
        assert.Nil(t, definition.CommandHandler)
        assert.Equal(t, "cqrs.DummyCommandProvider", definition.Name)
    }
}

func TestCommandBus_AddRegistryWithoutPointer(t *testing.T) {
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    //bus.AddRegistry(dummyRegistry{})
    bus.AddRegistry(dummyRegistry{})
    assert.Empty(t, bus.handlers)
}

func TestCommandBus_HandleUnsupportedCommand(t *testing.T) {
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    e := bus.Handle(&dummyCommand{})
    assert.NotNil(t, e)
}

func TestCommandBus_HandleSupportedCommand(t *testing.T) {
    handler := &dummyHandler{}
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    bus.Add(CommandDefinition{
        Name: "cqrs.dummyCommand",
        Provider: func() CommandHandler {
            return handler
        },
    })

    assert.Len(t, bus.handlers, 1)

    command := &dummyCommand{}

    assert.Nil(t, bus.Handle(command))
    assert.True(t, handler.Handled)
    assert.Same(t, command, handler.HandledCommand)
}

func TestCommandBus_HandleSupportedCommandWithError(t *testing.T) {
    handler := &dummyHandlerWithError{}
    bus := NewCommandBus(&dummyEventBus{}).(*commandBus)
    bus.Add(CommandDefinition{
        Name: "cqrs.dummyCommand",
        Provider: func() CommandHandler {
            return handler
        },
    })

    assert.Len(t, bus.handlers, 1)

    command := &dummyCommand{}

    assert.NotNil(t, bus.Handle(command))
    assert.True(t, handler.Handled)
    assert.Same(t, command, handler.HandledCommand)
}
