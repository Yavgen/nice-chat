package pipe

import (
	"bytes"
	"chat/internal/action"
	"chat/internal/auth"
	"chat/internal/client"
	"chat/internal/request"
	"encoding/json"
	"log"
)

type ReadPipe struct {
	actionFactory action.Factory
	authorizer    auth.Authorizer
}

func NewReadPipe(actionFactory action.Factory, authorizer auth.Authorizer) ReadPipe {
	return ReadPipe{actionFactory: actionFactory, authorizer: authorizer}
}

func (rp ReadPipe) Read(client client.ChatClient) {
	defer func() {
		client.CloseConnection()
		rp.authorizer.LogoutChatClient(client)
	}()

	for {
		_, message, readError := client.ReadMessage()
		if readError != nil {
			log.Println(readError)
			break
		}

		var chatRequest request.ChatRequest
		decoder := json.NewDecoder(bytes.NewReader(message))
		if decodeError := decoder.Decode(&chatRequest); decodeError != nil {
			log.Println(decodeError)
			break
		}

		_, _, err := rp.authorizer.GetConnectedUserByRequest(chatRequest)

		if err != nil {
			log.Println(err)
			break
		}

		chatAction, err := rp.actionFactory.MakeAction(chatRequest.Action)

		if err != nil {
			log.Println(err)
			break
		}

		err = chatAction.Handle(chatRequest)

		if err != nil {
			log.Println(err)
			break
		}
	}
}
