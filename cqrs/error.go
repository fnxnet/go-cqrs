package cqrs

type UnsupportedCommand interface {
    Error() string
}

type UnsupportedEvent interface {
    Error() string
}

type unsupportedCommandError struct {
    Message     string `json:"message"`
    CommandName string `json:"command"`
}

func (e unsupportedCommandError) Error() string {
    return e.Message
}

func NewUnsupportedCommand(command Command) UnsupportedCommand {
    return unsupportedCommandError{
        "cqrs.error.unsupportedCommand",
        command.CommandName(),
    }
}

type unsupportedEventError struct {
    Message     string `json:"message"`
    EventName string `json:"event"`
}

func (e unsupportedEventError) Error() string {
    return e.Message
}

func NewUnsupportedEvent(event Event) UnsupportedEvent {
    return unsupportedEventError{
        "cqrs.error.unsupportedEvent",
        event.Name(),
    }
}

type UnsupportedItem interface {
    Error() string
}

type unsupportedItemError struct {
    Message string `json:"message"`
    Handler string `json:"handler"`
    Item    string `json:"item"`
}

func (e unsupportedItemError) Error() string {
    return e.Message
}

func NewUnsupportedItem(handler interface{}, i interface{}) UnsupportedItem {
    return unsupportedItemError{
        "cqrs.error.unsupportedItem",
        extractName(handler),
        extractName(i),
    }
}
