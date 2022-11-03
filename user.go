package main

type User struct {
	Name     string
	Password string
	Token    string
	Client   *Client
}
