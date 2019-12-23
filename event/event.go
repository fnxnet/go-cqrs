package event

type Item struct {
    name string
}

func (command Item) IsHandled() bool {
    return false
}

func (command *Item) SetHandled() {}
