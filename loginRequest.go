package main

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
