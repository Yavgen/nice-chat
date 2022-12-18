package response

import (
	"chat/internal/event"
)

type LoginResponse struct {
	token string
	event event.LoginEvent
}

func NewLoginResponse(token string) LoginResponse {
	return LoginResponse{token: token, event: event.NewLoginEvent()}
}

func (lr LoginResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data: map[string]interface{}{
			"token": lr.token,
		},
		Status: StatusOk,
		Event:  lr.event.Name(),
	}
}
