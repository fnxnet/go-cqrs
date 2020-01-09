package cqrs

type UnsupportedItem interface {
	Error() string
}

type unsupportedItemError struct {
	Message  string `json:"message"`
	ItemName string `json:"item"`
}

func (e unsupportedItemError) Error() string {
	return e.Message
}

func NewUnsupportedItem(i interface{}) UnsupportedItem {
	return unsupportedItemError{
		"cqrs.error.unsupportedItem",
		extractName(i),
	}
}
