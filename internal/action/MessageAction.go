package action

import (
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/request"
	"chat/internal/response"
	"errors"
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

func (ma MessageAction) Handle(chatRequest request.ChatRequest) error {
	user, isUserLogin := ma.loginUserStore.FindUserByToken(chatRequest.Token)

	if !isUserLogin {
		return errors.New("user not authorized")
	}

	roomName, isRoomNameValid := chatRequest.Data["roomName"].(string)

	if !isRoomNameValid {
		return errors.New("invalid room name")
	}

	message, isMessageValid := chatRequest.Data["message"].(string)

	if !isMessageValid {
		return errors.New("invalid message")
	}

	isPublicRoom := chatRequest.Data["roomName"].(string) == "Public"
	messageResponse := response.NewMessageResponse(message, roomName, user.Name())

	if isPublicRoom {
		ma.broadcastCh.Push(messageResponse.ToJson())

		return nil
	}

	chatRoom, isRoomExist := ma.roomsStore.FindByName(roomName)

	if !isRoomExist {
		return errors.New("room does not exists")
	}

	if isUserInRoom := chatRoom.IsUserInRoom(user.Name()); !isUserInRoom {
		return errors.New("user not in room")
	}

	//TODO разделить реквесты по типам
	ma.broadcastCh.Push(messageResponse.ToJson())

	return nil
}
