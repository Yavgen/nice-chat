package main

import (
	"bytes"
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

type Client struct {
	connection *websocket.Conn
	send       chan *Request
	token      string
}

func (client *Client) readPipe() {
	defer func() {
		client.connection.Close()
		delete(loginUsers, client.token)
	}()

	for {
		_, message, readError := client.connection.ReadMessage()
		if readError != nil {
			delete(loginUsers, client.token)
			log.Printf("client read error %v", readError)
			break
		}

		var request Request
		decoder := json.NewDecoder(bytes.NewReader(message))
		if decodeError := decoder.Decode(&request); decodeError != nil {
			log.Fatal(decodeError)
		}

		switch request.Action {
		case pongAction:
		case messageAction:
			broadcast <- &request
		case createRoomAction:
			log.Println("createRoom")
		case getUsersAction:
			usersUniqueNames := make(map[string]bool)

			for _, user := range loginUsers {
				usersUniqueNames[user.Name] = true
			}

			var usersNames []string

			for name, _ := range usersUniqueNames {
				usersNames = append(usersNames, name)
			}

			addUsersResponse := Response{
				Data:   map[string]interface{}{"users": usersNames},
				Status: "ok",
				Event:  addUsersEvent,
			}

			client.connection.WriteJSON(addUsersResponse)
		default:
			break
		}
	}
}

func (client *Client) writePipe() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		client.connection.Close()
		delete(loginUsers, client.token)
	}()

	for {
		select {
		case request, ok := <-client.send:
			if !ok {
				delete(loginUsers, client.token)
				closeResponse := Response{
					Data:   map[string]interface{}{"message": "connection closed"},
					Status: "ok",
					Event:  messageEvent,
				}

				client.connection.WriteJSON(closeResponse)
			}

			writer, writerError := client.connection.NextWriter(websocket.TextMessage)

			if writerError != nil {
				delete(loginUsers, client.token)
				return
			}

			user := loginUsers[request.Token]

			messageResponse := Response{
				Data:   map[string]interface{}{"message": request.Data["message"], "user": user.Name},
				Status: "ok",
				Event:  messageEvent,
			}

			json.NewEncoder(writer).Encode(messageResponse)

			queueCount := len(client.send)

			for i := 0; i < queueCount; i++ {
				queuedRequest := <-client.send
				queuedUser := loginUsers[queuedRequest.Token]
				queuedMessageResponse := Response{
					Data:   map[string]interface{}{"message": queuedRequest.Data["message"], "user": queuedUser.Name},
					Status: "ok",
					Event:  messageEvent,
				}

				json.NewEncoder(writer).Encode(queuedMessageResponse)
			}

			if writerCloseError := writer.Close(); writerCloseError != nil {
				return
			}
		case <-ticker.C:
			pingResponse := Response{
				Data:   map[string]interface{}{"message": "ping"},
				Status: "ok",
				Event:  pingEvent,
			}

			writer, writerError := client.connection.NextWriter(websocket.TextMessage)

			if writerError != nil {
				delete(loginUsers, client.token)
				return
			}

			if pingError := json.NewEncoder(writer).Encode(pingResponse); pingError != nil {
				delete(loginUsers, client.token)
				return
			}
		}
	}
}
