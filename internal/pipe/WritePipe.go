package pipe

import (
	"chat/internal/auth"
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/response"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type WritePipe struct {
	authorizer auth.Authorizer
	loginStore *store.LoginUsersStore
}

func NewWritePipe(authorizer auth.Authorizer, loginStore *store.LoginUsersStore) WritePipe {
	return WritePipe{authorizer: authorizer, loginStore: loginStore}
}

func (wp WritePipe) Write(client client.ChatClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		client.CloseConnection()
		wp.authorizer.LogoutChatClient(client)
	}()

	for {
		select {
		case chatResponse, ok := <-client.Send():
			if !ok {
				log.Println("writepipeErr")
				wp.authorizer.LogoutChatClient(client)
				//TODO публичную комнату в константу
				closeResponse := response.NewMessageResponse("connection closed", "Public", "")
				err := client.WriteJSON(closeResponse.ToJson())

				if err != nil {
					log.Println(err)
					break
				}
			}

			writer, writerError := client.NextWriter(websocket.TextMessage)

			if writerError != nil {
				log.Println(writeWait)
				break
			}

			json.NewEncoder(writer).Encode(chatResponse)

			queueCount := len(client.Send())

			for i := 0; i < queueCount; i++ {
				queuedResponse := <-client.Send()

				json.NewEncoder(writer).Encode(queuedResponse)
			}

			if writerCloseError := writer.Close(); writerCloseError != nil {
				return
			}
		case <-ticker.C:
			pingResponse := response.NewPingResponse()
			writer, writerError := client.NextWriter(websocket.TextMessage)

			if writerError != nil {
				log.Println(writerError)
				break
			}

			if pingError := json.NewEncoder(writer).Encode(pingResponse); pingError != nil {
				log.Println(pingError)
				break
			}
		}
	}
}
