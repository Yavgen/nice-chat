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
			break
		case messageAction:
			if _, ok := request.Data["roomName"].(string); !ok {
				break
			}

			if _, isUserLogin := loginUsers[request.Token]; !isUserLogin {
				break
			}

			if request.Data["roomName"].(string) == publicRoom {
				broadcast <- &request
				break
			}

			roomName := request.Data["roomName"].(string)

			if _, isRoomExist := rooms[roomName]; !isRoomExist {
				break
			}

			if _, isClientInRoom := rooms[roomName].Clients[request.Token]; !isClientInRoom {
				break
			}

			broadcast <- &request
		case createRoomAction:
			if _, ok := request.Data["roomName"].(string); !ok {
				break
			}

			if _, ok := loginUsers[request.Token]; !ok {
				break
			}

			if roomName, ok := request.Data["roomName"]; ok {
				if _, isRoomExist := rooms[roomName.(string)]; isRoomExist {
					break
				}

				roomClients := make(map[string]*RoomClient)

				roomClient := &RoomClient{
					connection: client,
					userName:   loginUsers[request.Token].Name,
				}

				roomClients[request.Token] = roomClient

				room := Room{
					OwnerToken: client.token,
					Clients:    roomClients,
				}

				rooms[roomName.(string)] = &room

				appendRoomResponse := Response{
					Data:   map[string]interface{}{"room": roomName},
					Status: "ok",
					Event:  appendRoomEvent,
				}

				client.connection.WriteJSON(appendRoomResponse)
			}
		case appendUserToRoomAction:
			if _, ok := request.Data["userName"].(string); !ok {
				break
			}

			if _, ok := loginUsers[request.Token]; !ok {
				break
			}

			var userClientToAppend *Client

			if registeredUser, isRegisteredUser := registeredUsers[request.Data["userName"].(string)]; isRegisteredUser {
				if userClient, isUserClient := clients[registeredUser.Token]; isUserClient {
					userClientToAppend = userClient
				} else {
					break
				}
			}

			if roomName, ok := request.Data["roomName"].(string); ok {
				if room, isRoomExist := rooms[request.Data["roomName"].(string)]; isRoomExist {
					roomClient := &RoomClient{
						connection: userClientToAppend,
						userName:   request.Data["userName"].(string),
					}
					room.Clients[userClientToAppend.token] = roomClient
					rooms[roomName] = room

					appendRoomResponse := Response{
						Data:   map[string]interface{}{"room": roomName},
						Status: "ok",
						Event:  appendRoomEvent,
					}

					userClientToAppend.connection.WriteJSON(appendRoomResponse)
				}
			}

		case getUsersAction:
			var usersNames []string

			if _, ok := request.Data["roomName"].(string); !ok {
				break
			}

			if request.Data["roomName"].(string) == publicRoom {
				usersUniqueNames := make(map[string]bool)

				for _, user := range loginUsers {
					usersUniqueNames[user.Name] = true
				}

				for name, _ := range usersUniqueNames {
					usersNames = append(usersNames, name)
				}
			} else {
				if _, ok := request.Data["roomName"].(string); !ok {
					break
				}

				roomName := request.Data["roomName"].(string)

				if room, isRoomExist := rooms[roomName]; isRoomExist {
					for _, roomClient := range room.Clients {
						usersNames = append(usersNames, roomClient.userName)
					}
				}
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
					Data:   map[string]interface{}{"message": "connection closed", "roomName": publicRoom},
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
				Data: map[string]interface{}{
					"message":  request.Data["message"],
					"user":     user.Name,
					"roomName": request.Data["roomName"].(string),
				},
				Status: "ok",
				Event:  messageEvent,
			}

			json.NewEncoder(writer).Encode(messageResponse)

			queueCount := len(client.send)

			for i := 0; i < queueCount; i++ {
				queuedRequest := <-client.send
				queuedUser := loginUsers[queuedRequest.Token]
				queuedMessageResponse := Response{
					Data: map[string]interface{}{
						"message":  queuedRequest.Data["message"],
						"user":     queuedUser.Name,
						"roomName": queuedRequest.Data["roomName"].(string),
					},
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
