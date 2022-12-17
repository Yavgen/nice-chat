package response

import "chat/internal/event"

type UpdateRoomUsersResponse struct {
	roomName  string
	roomUsers []string
}

// TODO: лучше передавать ивент через конструктор
func NewUpdateRoomUsersResponse(roomUsers []string, roomName string) UpdateRoomUsersResponse {
	return UpdateRoomUsersResponse{roomUsers: roomUsers, roomName: roomName}
}

func (ur UpdateRoomUsersResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data: map[string]interface{}{
			"users":    ur.roomUsers,
			"roomName": ur.roomName,
		},
		Status: StatusOk,
		Event:  event.NewUpdateRoomUsersEvent().Name(),
	}
}
