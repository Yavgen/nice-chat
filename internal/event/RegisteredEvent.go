package event

type RegisteredEvent struct {
	name string
}

func NewRegisteredEvent() RegisteredEvent {
	return RegisteredEvent{name: Registered}
}

func (le RegisteredEvent) Name() string {
	return le.name
}
