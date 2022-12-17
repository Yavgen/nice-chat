package event

type AppendRoomEvent struct {
	name string
}

func NewAppendRoomEvent() AppendRoomEvent {
	return AppendRoomEvent{name: AppendRoom}
}

func (a AppendRoomEvent) Name() string {
	return a.name
}
