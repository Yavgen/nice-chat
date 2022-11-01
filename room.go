package main

type Room struct {
	OwnerToken string
	Clients    []*Client
}
