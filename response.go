package main

type Response struct {
	Data   map[string]interface{} `json:"data"`
	Status string                 `json:"status"`
	Event  string                 `json:"event"`
}
