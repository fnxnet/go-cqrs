package cqrs

type CQRSCommand struct {
    Name     string
}

func (command *CQRSCommand) CommandName() string {
    return extractName(command)
}

type Command interface {
    CommandName() string
}

type CommandHandler interface {
    Handle(i interface{}, bus EventBus) (err error)
}

type CommandBus interface {
    Handle(command Command) (err error)
    Register(command Command, handler CommandHandler)
}

type commandBus struct {
    handlers map[string]CommandHandler
    eventBus EventBus
}

func (bus *commandBus) Handle(command Command) (err error) {
    name := extractName(command)

    if handler, ok := bus.handlers[name]; ok {
        if err = handler.Handle(command, bus.eventBus); err != nil {
            return err
        }
        return
    }

    return NewUnsupportedCommand(command)
}

func (bus *commandBus) Register(command Command, handler CommandHandler) {
    name := extractName(command)
    bus.handlers[name] = handler
}

func NewCommandBus(eb EventBus) (c CommandBus) {
    c = &commandBus{
        make(map[string]CommandHandler),
        eb,
    }

    return
}
