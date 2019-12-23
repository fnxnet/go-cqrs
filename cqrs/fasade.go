package cqrs

type Facade interface {
    CommandBus() CommandBus
    EventBus() EventBus
}

type FacaceConfig interface{}

type facade struct {
    eventBus   EventBus
    commandBus CommandBus
}

func (f facade) CommandBus() CommandBus {
    return f.commandBus
}

func (f facade) EventBus() EventBus {
    return f.eventBus
}

func NewFacade(config FacaceConfig) (f *facade) {
    eventBus := NewEventBus()
    f = &facade{
        eventBus,
        NewCommandBus(eventBus),
    }
    return
}
