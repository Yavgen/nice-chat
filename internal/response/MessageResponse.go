package response

import "chat/internal/event"

type MessageResponse struct {
	message  string
	roomName string
	userName string
}

func NewMessageResponse(message string, roomName string, userName string) MessageResponse {
	return MessageResponse{message: message, roomName: roomName, userName: userName}
}

func (mr MessageResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data: map[string]interface{}{
			"message": mr.message,
			"room":    mr.roomName,
			"user":    mr.userName,
		},
		Status: StatusOk,
		Event:  event.NewMessageEvent().Name(),
	}
}
