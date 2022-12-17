package response

import "chat/internal/event"

type AppendRoomResponse struct {
	room string
}

func NewAppendRoomResponse(room string) AppendRoomResponse {
	return AppendRoomResponse{room: room}
}

func (a AppendRoomResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data:   map[string]interface{}{"room": a.room},
		Status: "ok",
		Event:  event.NewAppendRoomEvent().Name(),
	}
}
