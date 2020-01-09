package main

type Polo interface {
	CommandName() string
}

type P struct {
}

func (P) CommandName() string {
	return "P"
}

type Handler func(i interface{}) error

type BusA struct {
	h Handler
}

func (b BusA) handle(p Polo) error {
	return b.h(p)
}

type BusB struct {
}

func (b BusB) handle(p Polo) error {
	return nil
}
