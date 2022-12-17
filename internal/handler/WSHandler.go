package handler

import (
	"bytes"
	"chat/internal/auth"
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/pipe"
	"chat/internal/request"
	"chat/internal/response"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WSHandler struct {
	upgrader           websocket.Upgrader
	loginUsersStore    *store.LoginUsersStore
	registeredClientCh client.RegisteredClientChannel
	authorizer         auth.Authorizer
	readPipe           pipe.ReadPipe
	writePipe          pipe.WritePipe
	broadcastCh        client.BroadcastChannel
}

func NewWSHandler(
	loginUsersStore *store.LoginUsersStore,
	registeredClientCh client.RegisteredClientChannel,
	authorizer auth.Authorizer,
	readPipe pipe.ReadPipe,
	writePipe pipe.WritePipe,
	broadcastCh client.BroadcastChannel,
) WSHandler {
	return WSHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		loginUsersStore:    loginUsersStore,
		registeredClientCh: registeredClientCh,
		authorizer:         authorizer,
		readPipe:           readPipe,
		writePipe:          writePipe,
		broadcastCh:        broadcastCh,
	}
}

func (wh WSHandler) Handle(writer http.ResponseWriter, httpRequest *http.Request) {
	connection, connectionError := wh.upgrader.Upgrade(writer, httpRequest, nil)

	if connectionError != nil {
		log.Println(connectionError)
		return
	}

	_, message, readError := connection.ReadMessage()

	if readError != nil {
		log.Printf("client read error %v", readError)
		return
	}

	var connRequest request.ChatRequest
	decoder := json.NewDecoder(bytes.NewReader(message))

	if decodeError := decoder.Decode(&connRequest); decodeError != nil {
		log.Println(decodeError)
		return
	}

	user, token, verifyErr := wh.authorizer.GetConnectedUserByRequest(connRequest)

	if verifyErr != nil {
		log.Println(verifyErr)
		return
	}

	connectedClient := client.NewClient(connection, make(chan *response.JsonResponse, 255), token.Value())
	//TODO: определить название публичной комнаты в константу
	connectedResponse := response.NewConnectedResponse("connected", user, "Public")
	writeError := connection.WriteJSON(connectedResponse.ToJson())

	if writeError != nil {
		log.Println(writeError)
		return
	}

	wh.broadcastCh.Push(connectedResponse.ToJson())
	wh.registeredClientCh.Push(&connectedClient)

	go wh.readPipe.Read(connectedClient)
	go wh.writePipe.Write(connectedClient)
}
