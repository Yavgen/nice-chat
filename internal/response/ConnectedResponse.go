package response

import (
	"chat/internal/domain/model"
	"chat/internal/event"
)

type ConnectedResponse struct {
	message  string
	user     model.User
	roomName string
}

func NewConnectedResponse(message string, user model.User, roomName string) ConnectedResponse {
	return ConnectedResponse{message: message, user: user, roomName: roomName}
}

func (cr ConnectedResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data: map[string]interface{}{
			"message": "connected",
			"user":    cr.user.Name(),
			"room":    cr.roomName,
		},
		Status: "ok",
		Event:  event.NewMessageEvent().Name(),
	}
}
