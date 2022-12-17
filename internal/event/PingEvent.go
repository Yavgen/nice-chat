package event

type PingEvent struct {
	name string
}

func NewPingEvent() PingEvent {
	return PingEvent{name: Ping}
}

func (pe PingEvent) Name() string {
	return pe.name
}
