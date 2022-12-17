package response

import (
	"chat/internal/event"
)

type RegisteredResponse struct {
	token string
	event event.RegisteredEvent
}

func NewRegisteredResponse(token string) RegisteredResponse {
	return RegisteredResponse{token: token, event: event.NewRegisteredEvent()}
}

func (rr RegisteredResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data: map[string]interface{}{
			"token": rr.token,
		},
		Status: StatusOk,
		Event:  rr.event.Name(),
	}
}
