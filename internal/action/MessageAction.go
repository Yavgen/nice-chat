package action

import (
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/request"
	"chat/internal/request/data"
	"chat/internal/response"
	"errors"
	"github.com/mitchellh/mapstructure"
)

type MessageAction struct {
	broadcastCh    client.BroadcastChannel
	roomsStore     *store.RoomsStore
	loginUserStore *store.LoginUsersStore
}

func NewMessageAction(
	broadcastCh client.BroadcastChannel,
	roomsStore *store.RoomsStore,
	loginUserStore *store.LoginUsersStore) MessageAction {
	return MessageAction{
		broadcastCh:    broadcastCh,
		roomsStore:     roomsStore,
		loginUserStore: loginUserStore,
	}
}

func (ma MessageAction) Handle(request request.ChatRequest) error {
	messageRequest := data.NewMessageRequest()
	err := mapstructure.Decode(request.Data, &messageRequest)

	if err != nil {
		return err
	}

	user, isUserLogin := ma.loginUserStore.FindUserByToken(request.Token)

	if !isUserLogin {
		return errors.New("user not authorized")
	}

	vErr := messageRequest.Validate()

	if vErr != nil {
		return vErr
	}

	isPublicRoom := messageRequest.RoomName == "Public"
	messageResponse := response.NewMessageResponse(messageRequest.Message, messageRequest.RoomName, user.Name())

	if isPublicRoom {
		ma.broadcastCh.Push(messageResponse.ToJson())

		return nil
	}

	chatRoom, isRoomExist := ma.roomsStore.FindByName(messageRequest.RoomName)

	if !isRoomExist {
		return errors.New("room does not exists")
	}

	if isUserInRoom := chatRoom.IsUserInRoom(user.Name()); !isUserInRoom {
		return errors.New("user not in room")
	}

	ma.broadcastCh.Push(messageResponse.ToJson())

	return nil
}
