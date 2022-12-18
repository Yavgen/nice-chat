package client

type ClientsStore struct {
	clients map[string]*ChatClient
}

func NewClientsStore() *ClientsStore {
	return &ClientsStore{clients: make(map[string]*ChatClient)}
}

func (cs ClientsStore) Clients() map[string]*ChatClient {
	return cs.clients
}

func (cs ClientsStore) FindByToken(token string) (ChatClient, bool) {
	chatClient, ok := cs.clients[token]

	if !ok {
		return ChatClient{}, false
	}

	return *chatClient, true
}

func (cs ClientsStore) MapByToken(token string, chatClient *ChatClient) {
	cs.clients[token] = chatClient
}

func (cs ClientsStore) DeleteByToken(token string) {
	delete(cs.clients, token)
}
