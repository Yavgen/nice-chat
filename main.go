package main

import (
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
	broadcast  = make(chan []byte)
	register   = make(chan *Client)
	unregister = make(chan *Client)
	clients    = make(map[*Client]bool)
)

func main() {
	go run()
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
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

	connectedResponse := Response{
		Data:   map[string]interface{}{"message": "connected"},
		Status: "ok",
	}

	writeError := connection.WriteJSON(connectedResponse)

	if writeError != nil {
		log.Println(writeError)
		return
	}

	connectedClient := &Client{connection: connection, send: make(chan []byte, 255)}
	register <- connectedClient
	go connectedClient.readPipe()
	go connectedClient.writePipe()
}

func run() {
	for {
		select {
		case client := <-register:
			clients[client] = true
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
