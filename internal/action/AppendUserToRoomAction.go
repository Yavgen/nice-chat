package action

import (
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/request"
	"chat/internal/response"
	"errors"
)

type AppendUserToRoomAction struct {
	registeredUsersStore *store.RegisteredUsersStore
	roomsStore           *store.RoomsStore
	loginUsersStore      *store.LoginUsersStore
	clientsStore         *client.ClientsStore
}

func NewAppendUserToRoomAction(
	registeredUsersStore *store.RegisteredUsersStore,
	roomsStore *store.RoomsStore,
	loginUsersStore *store.LoginUsersStore,
	clientsStore *client.ClientsStore,
) AppendUserToRoomAction {
	return AppendUserToRoomAction{
		registeredUsersStore: registeredUsersStore,
		roomsStore:           roomsStore,
		loginUsersStore:      loginUsersStore,
		clientsStore:         clientsStore,
	}
}

func (ar AppendUserToRoomAction) Handle(chatRequest request.ChatRequest) error {
	userName, isUserNameValid := chatRequest.Data["userName"].(string)

	if !isUserNameValid {
		return errors.New("invalid username")
	}

	roomName, isRoomNameValid := chatRequest.Data["roomName"].(string)

	if !isRoomNameValid {
		return errors.New("invalid room name")
	}

	chatRoom, isRoomExist := ar.roomsStore.FindByName(roomName)

	if !isRoomExist {
		return errors.New("room not exist")
	}

	user, isRegisteredUser := ar.registeredUsersStore.FindUserByName(userName)

	if !isRegisteredUser {
		return errors.New("user doesn't exist")
	}

	chatRoom.AppendUserToRoom(user.Name())
	ar.roomsStore.MapByName(roomName, chatRoom)
	appendRoomResponse := response.NewAppendRoomResponse(roomName).ToJson()
	appendUserToken, ok := ar.loginUsersStore.FindUserTokenByName(userName)

	if !ok {
		return errors.New("user not found")
	}

	appendUserClient, ok := ar.clientsStore.FindByToken(appendUserToken)

	if !ok {
		return errors.New("user not authorized")
	}

	appendUserClient.Send() <- &appendRoomResponse

	return nil
}
