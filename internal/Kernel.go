package internal

import (
	"chat/internal/client"
)

type Kernel struct {
	registeredClientsCh   client.RegisteredClientChannel
	unregisteredClientsCh client.UnregisteredClientsChannel
	broadcastCh           client.BroadcastChannel
	clientsStore          *client.ClientsStore
}

func NewKernel(
	registeredClientsCh client.RegisteredClientChannel,
	unregisteredClientsCh client.UnregisteredClientsChannel,
	broadcastCh client.BroadcastChannel,
	clientsStore *client.ClientsStore,
) Kernel {
	return Kernel{
		registeredClientsCh:   registeredClientsCh,
		unregisteredClientsCh: unregisteredClientsCh,
		broadcastCh:           broadcastCh,
		clientsStore:          clientsStore,
	}
}

func (k Kernel) Run() {
	for {
		select {
		case chatClient := <-k.registeredClientsCh.Listen():
			k.clientsStore.MapByToken(chatClient.Token(), chatClient)
		case chatClient := <-k.unregisteredClientsCh.Listen():
			if _, ok := k.clientsStore.FindByToken(chatClient.Token()); ok {
				chatClient.CloseSendCh()
			}
		case chatRequest := <-k.broadcastCh.Listen():
			for _, chatClient := range k.clientsStore.Clients() {
				select {
				case chatClient.Send() <- chatRequest:
				default:
					chatClient.CloseSendCh()
					k.clientsStore.DeleteByToken(chatClient.Token())
				}
			}
		}
	}
}
