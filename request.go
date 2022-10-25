package main

type Request struct {
	Token  string                 `json:"token"`
	Data   map[string]interface{} `json:"data"`
	Action string                 `json:"action"`
}
