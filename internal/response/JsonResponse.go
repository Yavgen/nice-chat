package response

const (
	StatusOk = "ok"
)

//TODO вынести RoomName в отдельный параметр

type JsonResponse struct {
	Data   map[string]interface{} `json:"data"`
	Status string                 `json:"status"`
	Event  string                 `json:"event"`
}
