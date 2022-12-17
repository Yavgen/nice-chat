package response

import "chat/internal/event"

type PingResponse struct {
}

func NewPingResponse() PingResponse {
	return PingResponse{}
}

func (pr PingResponse) ToJson() JsonResponse {
	return JsonResponse{
		Data: map[string]interface{}{
			"message": "ping",
		},
		Status: StatusOk,
		Event:  event.NewPingEvent().Name(),
	}
}
