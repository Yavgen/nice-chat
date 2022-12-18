package handler

import "net/http"

type IndexHandler struct {
}

func NewIndexHandler() IndexHandler {
	return IndexHandler{}
}

func (ih IndexHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Error(writer, "Not found!", http.StatusNotFound)
	}
	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed!", http.StatusMethodNotAllowed)
	}
	http.ServeFile(writer, request, "app/chat.html")
}
