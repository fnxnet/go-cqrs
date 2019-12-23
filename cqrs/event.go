package cqrs

import "fmt"

type CQRSEvent struct {
    EventName   string `json:"name"`
    EventDomain string `json:"name"`
    ID          string `json:"id"`
}

func (event *CQRSEvent) Name() string {
    return event.EventName
}

func (event *CQRSEvent) Domain() string {
    return event.EventDomain
}

func (event *CQRSEvent) Id() string {
    return event.ID
}

type Event interface {
    Name() string
    Domain() string
    Id() string
}

type EventHandler interface {
    Handle(i interface{}, bus EventBus) error
}

type EventMiddleware interface {
    Process(event Event) Event
}

type EventBus interface {
    Emit(events []Event)
    Register(handler EventHandler, events ... Event)
    AddMiddleware(middleware EventMiddleware)
}

type eventBus struct {
    handlers   map[string][]EventHandler
    middleware []EventMiddleware
}

func (bus *eventBus) AddMiddleware(middleware EventMiddleware) {
    bus.middleware = append(bus.middleware, middleware)
}

func (bus *eventBus) Emit(events []Event) {
    for _, event := range events {
        var errors []error
        name := extractName(event)

        for _, middle := range bus.middleware {
            event = middle.Process(event)
        }

        handlers, ok := bus.handlers[name]

        if ok && len(handlers) > 0 {
            for _, handler := range handlers {
                if err := handler.Handle(event, bus); err != nil {
                    errors = append(errors, err)
                }
            }
        } else {
            errors = append(errors, NewUnsupportedEvent(event))
        }

        fmt.Println(errors)
    }
}

func (bus *eventBus) Register(handler EventHandler, events ... Event) {
    for _, event := range events {
        name := extractName(event)
        bus.handlers[name] = append(bus.handlers[name], handler)
    }
}

func NewEventBus() (c EventBus) {
    c = &eventBus{
        make(map[string][]EventHandler),
        []EventMiddleware{},
    }

    return
}
