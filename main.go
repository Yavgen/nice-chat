package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var indexAddress = ":8080"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
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
