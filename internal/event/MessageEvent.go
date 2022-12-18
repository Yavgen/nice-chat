package event

type MessageEvent struct {
	name string
}

func (m MessageEvent) Name() string {
	return m.name
}

func NewMessageEvent() MessageEvent {
	return MessageEvent{name: Message}
}
