package response

import "chat/internal/event"

type UserAppendToRoomResponse struct {
	roomName string
}

func NewUserAppendToRoomResponse(roomName string) UserAppendToRoomResponse {
	return UserAppendToRoomResponse{roomName: roomName}
}

func (ur UserAppendToRoomResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data:   map[string]interface{}{"room": ur.roomName},
		Status: "ok",
		Event:  event.NewUserAppendToRoomEvent().Name(),
	}
}
