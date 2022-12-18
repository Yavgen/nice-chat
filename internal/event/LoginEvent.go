package event

type LoginEvent struct {
	name string
}

func NewLoginEvent() LoginEvent {
	return LoginEvent{name: Login}
}

func (le LoginEvent) Name() string {
	return le.name
}
