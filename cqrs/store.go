package cqrs

type EventStore interface {
    Get(domain string, id string) ([]Event)
    Store(event Event)
    StoreEvents(events []Event)
}
