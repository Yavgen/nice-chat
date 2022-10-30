package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	connection *websocket.Conn
	send       chan *Request
	token      string
}

func (client *Client) readPipe() {
	defer func() {
		client.connection.Close()
	}()

	for {
		_, message, readError := client.connection.ReadMessage()
		if readError != nil {
			log.Printf("client read error %v", readError)
			break
		}

		var request Request
		decoder := json.NewDecoder(bytes.NewReader(message))
		if decodeError := decoder.Decode(&request); decodeError != nil {
			log.Fatal(decodeError)
		}

		switch request.Action {
		case messageAction:
			broadcast <- &request
		case createRoomAction:
			log.Println("createRoom")
		case getUsersAction:
			var usersNames []string
			for user := range registeredUsers {
				usersNames = append(usersNames, user)
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
	defer func() {
		client.connection.Close()
	}()

	for {
		select {
		case request, ok := <-client.send:
			if !ok {
				closeResponse := Response{
					Data:   map[string]interface{}{"message": "connection closed"},
					Status: "ok",
					Event:  messageEvent,
				}

				client.connection.WriteJSON(closeResponse)
			}

			writer, writerError := client.connection.NextWriter(websocket.TextMessage)

			if writerError != nil {
				return
			}

			messageResponse := Response{
				Data:   map[string]interface{}{"message": request.Data["message"]},
				Status: "ok",
				Event:  messageEvent,
			}

			json.NewEncoder(writer).Encode(messageResponse)

			queueCount := len(client.send)

			for i := 0; i < queueCount; i++ {
				queuedRequest := <-client.send
				queuedMessageResponse := Response{
					Data:   map[string]interface{}{"message": queuedRequest.Data["message"]},
					Status: "ok",
					Event:  messageEvent,
				}

				json.NewEncoder(writer).Encode(queuedMessageResponse)
			}

			if writerCloseError := writer.Close(); writerCloseError != nil {
				return
			}
		}
	}

}
