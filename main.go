package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var indexAddress = ":8080"
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	broadcast       = make(chan *Request)
	register        = make(chan *Client)
	unregister      = make(chan *Client)
	clients         = make(map[*Client]string)
	loginUsers      = make(map[string]*User)
	registeredUsers = make(map[string]*User)
)

const createRoomAction = "createRoom"
const messageAction = "message"
const getUsersAction = "getUsers"

const appendRoomEvent = "appendRoom"
const messageEvent = "message"
const connectedEvent = "connected"
const connectionClosedEvent = "connectionClosed"
const addUsersEvent = "addUsers"

func main() {
	go run()
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/login", loginHandler)
	router.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		ServeWs(writer, request)
	})
	http.Handle("/", router)
	err := http.ListenAndServe(indexAddress, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/login" {
		http.Error(writer, "Not found!", http.StatusNotFound)
	}

	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed!", http.StatusMethodNotAllowed)
	}

	var loginRequest loginRequest
	decoder := json.NewDecoder(request.Body)

	if decodeError := decoder.Decode(&loginRequest); decodeError != nil {
		log.Fatal(decodeError)
	}

	if user, ok := registeredUsers[loginRequest.Name]; ok {
		if user.Password == loginRequest.Password {
			user.Token = generateToken()

			loginResponse := Response{
				Data:   map[string]interface{}{"token": user.Token},
				Status: "ok",
				Event:  "login",
			}

			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(loginResponse)
			loginUsers[user.Token] = user
		} else {
			http.Error(writer, "Password wrong!", http.StatusForbidden)
			return
		}
	} else {
		user = &User{
			Name:     loginRequest.Name,
			Password: loginRequest.Password,
			Token:    generateToken(),
		}
		registeredUsers[loginRequest.Name] = user

		registerResponse := Response{
			Data:   map[string]interface{}{"token": registeredUsers[loginRequest.Name].Token},
			Status: "ok",
			Event:  "register",
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(registerResponse)
		loginUsers[user.Token] = user
	}
}

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Error(writer, "Not found!", http.StatusNotFound)
	}
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed!", http.StatusMethodNotAllowed)
	}
	http.ServeFile(writer, request, "chat.html")
}

func ServeWs(writer http.ResponseWriter, request *http.Request) {
	connection, connectionError := upgrader.Upgrade(writer, request, nil)

	if connectionError != nil {
		log.Println(connectionError)
		return
	}

	_, message, readError := connection.ReadMessage()

	if readError != nil {
		log.Printf("client read error %v", readError)
		http.Error(writer, "Bad request", http.StatusBadRequest)
	}

	var connectionRequest Request
	decoder := json.NewDecoder(bytes.NewReader(message))

	if decodeError := decoder.Decode(&connectionRequest); decodeError != nil {
		log.Fatal(decodeError)
	}

	if _, ok := loginUsers[connectionRequest.Token]; !ok {
		http.Error(writer, "Password wrong!", http.StatusForbidden)
		return
	}

	connectedResponse := Response{
		Data:   map[string]interface{}{"message": "connected"},
		Status: "ok",
		Event:  messageEvent,
	}

	writeError := connection.WriteJSON(connectedResponse)

	if writeError != nil {
		log.Println(writeError)
		return
	}

	connectedClient := &Client{connection: connection, send: make(chan *Request, 255), token: connectionRequest.Token}
	register <- connectedClient
	go connectedClient.readPipe()
	go connectedClient.writePipe()
}

func run() {
	for {
		select {
		case client := <-register:
			clients[client] = client.token
		case client := <-unregister:

			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.send)
			}

		case message := <-broadcast:

			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}

		}
	}
}

func generateToken() string {
	return uuid.NewString()
}
