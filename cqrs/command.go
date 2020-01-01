package cqrs

import (
    "errors"
    "fmt"
    "reflect"
    "strings"
)

type CQRSCommand struct {
    CommandName string
}

func (c *CQRSCommand) Name() string {
    return c.CommandName
}

type Command interface {
    CommandName() string
}

type CommandHandler interface {
    Handle(i interface{}, bus EventBus) (err error)
}

type Provider func() CommandHandler

type CommandDefinition struct {
    Name           string
    CommandHandler CommandHandler
    Provider       Provider
}

type CommandBus interface {
    Handle(command Command) (err error)
    AddHandler(command Command, handler CommandHandler)
    AddProvider(command Command, p Provider)
    AddRegistry(interface{})
    Add(definition CommandDefinition)
}

type commandBus struct {
    handlers map[string]CommandDefinition
    eventBus EventBus
}

func (bus *commandBus) Handle(c Command) (err error) {
    name := extractName(c)

    if def, ok := bus.handlers[name]; ok {
        handler := def.CommandHandler
        if handler == nil {
            handler = def.Provider()
        }

        if err = handler.Handle(c, bus.eventBus); err != nil {
            return err
        }
        return
    }

    return NewUnsupportedCommand(c)
}

func (bus *commandBus) AddProvider(c Command, p Provider) {
    name := extractName(c)
    bus.handlers[name] = CommandDefinition{
        name,
        nil,
        p,
    }
}

func (bus *commandBus) AddHandler(c Command, h CommandHandler) {
    name := extractName(c)
    bus.handlers[name] = CommandDefinition{
        name,
        h,
        nil,
    }
}

func checkMethod(r interface{}, methodName string) (p Provider, e error) {
    defer func() {
        if r := recover(); r != nil {
            e = errors.New(fmt.Sprintf("Invalid output: %s", r))
            return
        }
    }()

    resValue := reflect.ValueOf(r).MethodByName(methodName).Call([]reflect.Value{})

    if len(resValue) != 1 {
        return nil, errors.New("Too many arguments")
    }

    p, ok := resValue[0].Interface().(func() CommandHandler)

    if !ok {
        p, ok = resValue[0].Interface().(Provider)
    }

    if p == nil {
        return nil, errors.New("No provider received")
    }

    return
}

func (bus *commandBus) AddRegistry(r interface{}) {
    t := reflect.TypeOf(r)

    if t.Kind() != reflect.Ptr {
        return
    }

    pkg := strings.Split(t.String(), ".")[0][1:]

    for i := 0; i < t.NumMethod(); i++ {

        method := t.Method(i)

        provider, e := checkMethod(r, method.Name)

        if provider == nil || e != nil {
            continue
        }

        commandName := fmt.Sprintf("%s.%s", pkg, method.Name)

        bus.handlers[commandName] = CommandDefinition{
            Name:     commandName,
            Provider: provider,
        }
    }
}

func (bus *commandBus) Add(d CommandDefinition) {
    bus.handlers[d.Name] = d
}

func NewCommandBus(eb EventBus) (c CommandBus) {
    c = &commandBus{
        make(map[string]CommandDefinition),
        eb,
    }

    return
}
