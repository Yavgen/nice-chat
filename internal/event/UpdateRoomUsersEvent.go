package event

type UpdateRoomUsersEvent struct {
	name string
}

func NewUpdateRoomUsersEvent() UpdateRoomUsersEvent {
	return UpdateRoomUsersEvent{name: UpdateRoomUsers}
}

func (uu UpdateRoomUsersEvent) Name() string {
	return uu.name
}
