package request

type ChatRequest struct {
	Data   map[string]interface{} `json:"data"`
	Token  string                 `json:"token"`
	Action string                 `json:"action"`
}
