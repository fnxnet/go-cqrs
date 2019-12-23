package bus

import (
    "errors"
    "reflect"
)

type cqrsItem interface {
    IsHandled() bool
    SetHandled()
}

type cqrsBus interface {
    GetHandlers() map[string]reflect.Value
    Register(handler interface{}) *bus
    Handle(in interface{}) (e error)
}

func New() *bus {
    return &bus{make(map[string]reflect.Value)}
}

type bus struct {
    handlers map[string]reflect.Value
}

func (b *bus) GetHandlers() map[string]reflect.Value {
    return b.handlers
}

func (b *bus) Register(handler interface{}) *bus {

    if i, ok := handler.(cqrsBus); ok {
        registerBus(b, i)
    } else {
        registerHandler(b, handler)
    }

    return b
}

func (b *bus) Handle(in interface{}) (e error) {
    item, ok := in.(cqrsItem)

    if !ok {
        return errors.New("Not CQRS Item")
    }

    if item.IsHandled() {
        return errors.New("Already handled")
    }

    t := reflect.TypeOf(in)
    name := t.String()

    if handler, ok := b.handlers[name]; ok {
        inputs := make([]reflect.Value, 1)
        inputs[0] = reflect.ValueOf(in)

        handler.Call(inputs)
        item.SetHandled()
    }

    return
}

func registerBus(bus *bus, i cqrsBus) {
    for k, v := range i.GetHandlers() {
        bus.handlers[k] = v
    }
}

func registerHandler(bus *bus, handler interface{}) {
    t := reflect.TypeOf(handler)
    for i := 0; i < t.NumMethod(); i++ {
        methodType := t.Method(i).Type

        if methodType.NumIn() > 0 {
            actionName := methodType.In(1).String()
            bus.handlers[actionName] = reflect.ValueOf(handler).MethodByName(t.Method(i).Name)
        }
    }
}
