package pipe

import (
	"bytes"
	"chat/internal/action"
	"chat/internal/auth"
	"chat/internal/client"
	"chat/internal/request"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

		validationErr := validation.Errors{
			"action": validation.Validate(chatRequest.Action, validation.Required, validation.Length(5, 20)),
			"data":   validation.Validate(chatRequest.Data, validation.Required, validation.Required),
			"token":  validation.Validate(chatRequest.Token, validation.Required, validation.Length(36, 36)),
		}.Filter()

		if validationErr != nil {
			log.Println(validationErr)
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
