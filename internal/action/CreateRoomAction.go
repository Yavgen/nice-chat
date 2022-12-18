package action

import (
	"chat/internal/client"
	"chat/internal/domain/model"
	"chat/internal/domain/store"
	"chat/internal/request"
	"chat/internal/response"
	"errors"
)

type CreateRoomAction struct {
	roomsStore     *store.RoomsStore
	loginUserStore *store.LoginUsersStore
	clientsStore   *client.ClientsStore
}

func NewCreateRoomAction(
	roomsStore *store.RoomsStore,
	loginUserStore *store.LoginUsersStore,
	clientsStore *client.ClientsStore,
) CreateRoomAction {
	return CreateRoomAction{
		roomsStore:     roomsStore,
		loginUserStore: loginUserStore,
		clientsStore:   clientsStore,
	}
}

func (ca CreateRoomAction) Handle(chatRequest request.ChatRequest) error {
	//TODO: вынести валидацию
	if _, ok := chatRequest.Data["roomName"].(string); !ok {
		return errors.New("invalid room name")
	}

	roomName := chatRequest.Data["roomName"].(string)
	_, isRoomExist := ca.roomsStore.FindByName(roomName)

	if isRoomExist {
		return errors.New("room already exist")
	}

	if isPublicRoom := chatRequest.Data["roomName"].(string) == "Public"; isPublicRoom {
		return errors.New("room already exist")
	}

	user, isUserLogin := ca.loginUserStore.FindUserByToken(chatRequest.Token)

	if !isUserLogin {
		return errors.New("user not authorized")
	}

	chatRoom := model.NewRoom(
		roomName,
		user.Name(),
		map[string]bool{user.Name(): true},
	)

	ca.roomsStore.MapByName(roomName, chatRoom)
	appendRoomResponse := response.NewAppendRoomResponse(roomName).ToJson()
	chatClient, ok := ca.clientsStore.FindByToken(chatRequest.Token)

	if !ok {
		return errors.New("user not authorized")
	}

	chatClient.Send() <- &appendRoomResponse

	return nil
}
