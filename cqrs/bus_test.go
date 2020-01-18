package cqrs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyItem struct{}

type dummyHandler struct {
	Name        string
	Called      int
	HandledItem interface{}
}

func (d *dummyHandler) Handle(i interface{}) (err error) {
	d.Called++
	d.HandledItem = i
	return
}

type dummyHandlerWithError struct {
	Name           string
	Called         int
	HandledCommand interface{}
}

func (d *dummyHandlerWithError) Handle(i interface{}) (err error) {
	d.Called++
	d.HandledCommand = i
	return errors.New("error")
}

type dummyRegistry struct{}

func (d *dummyRegistry) InvalidItem() {}

func (d *dummyRegistry) InvalidItemWithParam(i int) {}

func (d *dummyRegistry) InvalidItemWithReturnInt() int {
	return 0
}

func (d *dummyRegistry) InvalidItemWithReturnFunc() func() {
	return func() {}
}

func (d *dummyRegistry) InvalidItemWithReturnFuncInt() func() int {
	return func() int { return 0 }
}

func (d *dummyRegistry) InvalidItemWithReturnAndName() (e int) {
	return 0
}

func (d *dummyRegistry) InvalidItemWithReturnAndParam(i int) int {
	return 0
}

func (d *dummyRegistry) InvalidItemWithProviderReturned(i int) (f Provider) {
	f = func() Handler { return &dummyHandler{Name: "invalidHandler"} }
	return
}

func (d *dummyRegistry) InvalidItemWithMultipleOutputs() (f Provider, e error) {
	f = func() Handler { return &dummyHandler{Name: "invalidHandler"} }
	return
}

func (d *dummyRegistry) InvalidItemWithMultipleUnnamedOutputs() (Provider, error) {
	return func() Handler { return &dummyHandler{Name: "invalidHandler"} }, nil
}

func (d *dummyRegistry) DummyCommand() func() Handler {
	return func() Handler { return &dummyHandler{Name: "dummyCommandHandler"} }
}

func (d *dummyRegistry) DummyCommandProvider() Provider {
	return func() Handler { return &dummyHandler{Name: "dummyCommandHandler2"} }
}

type dummyMiddleware struct {
	called int
	item   interface{}
}

func (d *dummyMiddleware) Process(i interface{}) interface{} {
	d.called++
	d.item = i
	return i
}

func TestNewCommandBus(t *testing.T) {
	assert.IsType(t, (*bus)(nil), NewBus())
}

func TestCommandBus_Add(t *testing.T) {
	def := &Definition{Name: "dummy", Provider: func() Handler { return &dummyHandler{} }}
	bus := NewBus().(*bus)
	bus.Add(def)

	assert.Len(t, bus.handlers, 1)
	assert.Contains(t, bus.handlers, "dummy")
	assert.NotNil(t, bus.handlers["dummy"].Provider)
}

func TestCommandBus_AddHandler(t *testing.T) {
	h := &dummyHandler{}
	bus := NewBus().(*bus)
	bus.AddHandler(dummyItem{}, h)

	assert.Len(t, bus.handlers, 1)

	definition := bus.handlers["cqrs.dummyItem"]

	assert.IsType(t, (*Definition)(nil), definition)
	assert.NotNil(t, definition.Handler)
	assert.Nil(t, definition.Provider)
	assert.Equal(t, h, definition.Handler)
	assert.Equal(t, "cqrs.dummyItem", definition.Name)
}

func TestCommandBus_AddProvider(t *testing.T) {
	h := &dummyHandler{}
	f := func() Handler {
		return h
	}

	bus := NewBus().(*bus)
	bus.AddProvider(dummyItem{}, f)

	assert.Len(t, bus.handlers, 1)

	definition := bus.handlers["cqrs.dummyItem"]

	assert.IsType(t, (*Definition)(nil), definition)
	assert.NotNil(t, definition.Provider)
	assert.IsType(t, (Provider)(nil), definition.Provider)
	assert.Equal(t, h, definition.Provider())
	assert.Nil(t, definition.Handler)
	assert.Equal(t, "cqrs.dummyItem", definition.Name)
}

func TestCommandBus_AddRegistry(t *testing.T) {
	bus := NewBus().(*bus)
	bus.AddRegistry(&dummyRegistry{})

	assert.Len(t, bus.handlers, 2)
	if assert.Contains(t, bus.handlers, "cqrs.DummyCommand") {
		definition := bus.handlers["cqrs.DummyCommand"]

		assert.IsType(t, (*Definition)(nil), definition)
		assert.NotNil(t, definition.Provider)
		assert.IsType(t, (Provider)(nil), definition.Provider)
		assert.Equal(t, "dummyCommandHandler", definition.Provider().(*dummyHandler).Name)
		assert.Nil(t, definition.Handler)
		assert.Equal(t, "cqrs.DummyCommand", definition.Name)
	}

	if assert.Contains(t, bus.handlers, "cqrs.DummyCommandProvider") {
		definition := bus.handlers["cqrs.DummyCommandProvider"]

		assert.IsType(t, (*Definition)(nil), definition)
		assert.NotNil(t, definition.Provider)
		assert.IsType(t, (Provider)(nil), definition.Provider)
		assert.Equal(t, "dummyCommandHandler2", definition.Provider().(*dummyHandler).Name)
		assert.Nil(t, definition.Handler)
		assert.Equal(t, "cqrs.DummyCommandProvider", definition.Name)
	}
}

func TestCommandBus_AddRegistryWithoutPointer(t *testing.T) {
	bus := NewBus().(*bus)
	bus.AddRegistry(dummyRegistry{})

	assert.Empty(t, bus.handlers)
}

func TestCommandBus_HandleUnsupportedCommand(t *testing.T) {
	bus := NewBus().(*bus)
	e := bus.Handle(&dummyItem{})

	assert.NotNil(t, e)
}

func TestCommandBus_HandleNilValue(t *testing.T) {
	m := &dummyMiddleware{}
	h := &dummyHandler{}
	bus := NewBus().(*bus)
	pCall := 0
	bus.Add(&Definition{
		Name: "cqrs.dummyItem",
		Provider: func() Handler {
			pCall++
			return h
		},
	})
	bus.AddMiddleware(m)

	assert.Len(t, bus.middleware, 1)
	assert.Len(t, bus.handlers, 1)
	assert.Equal(t, 0, m.called)

	assert.Empty(t, bus.Handle(nil))
	assert.Equal(t, 0, h.Called)
	assert.Equal(t, 0, m.called)
}

func TestCommandBus_HandleSupportedCommand(t *testing.T) {
	m := &dummyMiddleware{}
	h := &dummyHandler{}
	bus := NewBus().(*bus)
	pCall := 0
	bus.Add(&Definition{
		Name: "cqrs.dummyItem",
		Provider: func() Handler {
			pCall++
			return h
		},
	})
	bus.AddMiddleware(m)

	assert.Len(t, bus.middleware, 1)
	assert.Len(t, bus.handlers, 1)
	assert.Equal(t, 0, m.called)

	item := &dummyItem{}

	assert.Nil(t, bus.Handle(item))
	assert.Equal(t, 1, h.Called)
	assert.Same(t, item, h.HandledItem)
	assert.Equal(t, 1, m.called)
	assert.Equal(t, 1, pCall)
	assert.NotNil(t, bus.handlers["cqrs.dummyItem"].Handler)

	// handle second time and see if same handler called
	assert.Nil(t, bus.Handle(item))
	assert.Equal(t, 1, pCall)
	assert.Equal(t, 2, h.Called)
	assert.Equal(t, 2, m.called)

	// handle slice
	assert.Nil(t, bus.Handle(item, item))
	assert.Equal(t, 1, pCall)
	assert.Equal(t, 4, h.Called)
	assert.Equal(t, 4, m.called)

	// handle array - should ignore this
	items := append(make([]interface{}, 0), item, item)
	assert.Nil(t, bus.Handle(items))
	assert.Equal(t, 1, pCall)
	assert.Equal(t, 4, h.Called)
	assert.Equal(t, 4, m.called)

	// handle variadic slice - should handle
	assert.Nil(t, bus.Handle(items[:]...))
	assert.Equal(t, 1, pCall)
	assert.Equal(t, 6, h.Called)
	assert.Equal(t, 6, m.called)
}

func TestCommandBus_HandleSupportedCommandWithError(t *testing.T) {
	handler := &dummyHandlerWithError{}
	bus := NewBus().(*bus)
	bus.Add(&Definition{
		Name: "cqrs.dummyItem",
		Provider: func() Handler {
			return handler
		},
	})

	assert.Len(t, bus.handlers, 1)

	item := &dummyItem{}

	assert.NotNil(t, bus.Handle(item))
	assert.Equal(t, 1, handler.Called)
	assert.Same(t, item, handler.HandledCommand)
}
