package mediator

type INotifiable interface {
	AddEvent(event interface{})
	ClearEvents()
	GetEvents() []interface{}
}
