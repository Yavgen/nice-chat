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

func main() {
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

	writeError := connection.WriteMessage(1, []byte("connected"))
	if writeError != nil {
		log.Println(writeError)
		return
	}
}
