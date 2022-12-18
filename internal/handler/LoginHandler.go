package handler

import (
	"chat/internal/auth"
	"chat/internal/event"
	"chat/internal/request"
	"chat/internal/response"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type LoginHandler struct {
	authorizer auth.Authorizer
}

func NewLoginHandler(
	authorizer auth.Authorizer,
) LoginHandler {
	return LoginHandler{
		authorizer: authorizer,
	}
}

func (lh LoginHandler) Handle(writer http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.URL.Path != "/login" {
		http.Error(writer, "Not found!", http.StatusNotFound)
	}

	if httpRequest.Method != http.MethodPost {
		http.Error(writer, "Method not allowed!", http.StatusMethodNotAllowed)
	}

	var authRequest request.AuthRequest
	decoder := json.NewDecoder(httpRequest.Body)

	if decodeError := decoder.Decode(&authRequest); decodeError != nil {
		log.Println(decodeError)
		http.Error(writer, "bad request", http.StatusBadRequest)

		return
	}

	authEvent, token, authorizeErr := lh.authorizer.Authorize(authRequest)

	if authorizeErr != nil {
		http.Error(writer, fmt.Sprint(authorizeErr), http.StatusForbidden)

		return
	}

	authResponse, responseErr := lh.makeResponse(authEvent, token)

	if responseErr != nil {
		http.Error(writer, fmt.Sprint(responseErr), http.StatusForbidden)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(authResponse.ToJson())
}

func (lh LoginHandler) makeResponse(chatEvent event.ChatEvent, token auth.Token) (response.ChatResponse, error) {

	_, ok := chatEvent.(event.LoginEvent)

	if ok {
		return response.NewLoginResponse(token.Value()), nil
	}

	_, ok = chatEvent.(event.RegisteredEvent)

	if ok {
		return response.NewRegisteredResponse(token.Value()), nil
	}

	return nil, errors.New("cant handle this event in LoginHandler")
}
