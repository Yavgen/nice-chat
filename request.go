package main

type Request struct {
	Token   string `json:"token"`
	Room    string `json:"room"`
	Message string `json:"message"`
}
