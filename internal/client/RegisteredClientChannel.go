package client

type RegisteredClientChannel struct {
	channel chan *ChatClient
}

func NewRegisteredClientChannel() RegisteredClientChannel {
	return RegisteredClientChannel{channel: make(chan *ChatClient)}
}

func (rc RegisteredClientChannel) Push(client *ChatClient) {
	rc.channel <- client
}

func (rc RegisteredClientChannel) Listen() <-chan *ChatClient {
	return rc.channel
}
