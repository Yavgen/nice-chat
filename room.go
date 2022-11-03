package main

type RoomClient struct {
	connection *Client
	userName   string
}

type Room struct {
	OwnerToken string
	Clients    map[string]*RoomClient
}
