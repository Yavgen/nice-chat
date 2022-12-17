package event

type UserAppendToRoomEvent struct {
	name string
}

func NewUserAppendToRoomEvent() UserAppendToRoomEvent {
	return UserAppendToRoomEvent{name: UserAppendToRoom}
}

func (ur UserAppendToRoomEvent) Name() string {
	return ur.name
}
