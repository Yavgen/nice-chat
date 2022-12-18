package client

type UnregisteredClientsChannel struct {
	channel chan *ChatClient
}

func NewUnregisteredClientsChannel() UnregisteredClientsChannel {
	return UnregisteredClientsChannel{channel: make(chan *ChatClient)}
}

func (uc UnregisteredClientsChannel) push(chatClient *ChatClient) {
	uc.channel <- chatClient
}

func (uc UnregisteredClientsChannel) Listen() <-chan *ChatClient {
	return uc.channel
}
