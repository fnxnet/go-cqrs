package cqrs

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Handler interface {
	Handle(i interface{}) (err error)
}

type Provider func() Handler

type Definition struct {
	Name     string
	Handler  Handler
	Provider Provider
}

type Bus interface {
	Handle(i ...interface{}) (err []error)
	AddHandler(i interface{}, h Handler)
	AddProvider(i interface{}, p Provider)
	AddRegistry(interface{})
	Add(d *Definition)
	AddMiddleware(m Middleware)
}

type Middleware interface {
	Process(i interface{}) interface{}
}

type bus struct {
	handlers   map[string]*Definition
	middleware []Middleware
}

func (bus *bus) AddMiddleware(m Middleware) {
	bus.middleware = append(bus.middleware, m)
}

func (bus *bus) handleItem(i interface{}) (err error) {
	if i == nil {
		return
	}

	if kind := reflect.TypeOf(i).Kind(); kind != reflect.Struct && kind != reflect.Ptr {
		return
	}

	name := extractName(i)

	for _, middle := range bus.middleware {
		i = middle.Process(i)
	}

	if def, ok := bus.handlers[name]; ok {
		handler := def.Handler
		if handler == nil {
			handler = def.Provider()
			def.Handler = handler
		}

		if err = handler.Handle(i); err != nil {
			return err
		}
		return
	}

	return NewUnsupportedItem(i)
}

func (bus *bus) Handle(i ...interface{}) (err []error) {
	for _, item := range i {
		if e := bus.handleItem(item); e != nil {
			err = append(err, e)
		}
	}
	return
}

func (bus *bus) AddProvider(i interface{}, p Provider) {
	name := extractName(i)
	bus.handlers[name] = &Definition{
		name,
		nil,
		p,
	}
}

func (bus *bus) AddHandler(i interface{}, h Handler) {
	name := extractName(i)
	bus.handlers[name] = &Definition{
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

	p, ok := resValue[0].Interface().(func() Handler)

	if !ok {
		p, ok = resValue[0].Interface().(Provider)
	}

	if p == nil {
		return nil, errors.New("No provider received")
	}

	return
}

func (bus *bus) AddRegistry(r interface{}) {
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

		item := fmt.Sprintf("%s.%s", pkg, method.Name)

		bus.handlers[item] = &Definition{
			Name:     item,
			Provider: provider,
		}
	}
}

func (bus *bus) Add(d *Definition) {
	bus.handlers[d.Name] = d
}

func NewBus() (c Bus) {
	c = &bus{
		make(map[string]*Definition),
		[]Middleware{},
	}

	return
}
