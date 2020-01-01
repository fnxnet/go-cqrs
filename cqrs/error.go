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
