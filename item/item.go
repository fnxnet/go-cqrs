package item

type CQRSItem struct {
    handled bool
}

func (command CQRSItem) IsHandled() bool {
    if command.handled {
        return true
    }
    return false
}

func (command *CQRSItem) SetHandled() {
    command.handled = true
}
