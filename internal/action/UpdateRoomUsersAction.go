package action

import (
	"chat/internal/client"
	"chat/internal/domain/store"
	"chat/internal/request"
	"chat/internal/response"
	"errors"
)

type UpdateRoomUsersAction struct {
	loginUsersStore *store.LoginUsersStore
	roomsStore      *store.RoomsStore
	clientsStore    *client.ClientsStore
}

//TODO отдавать сторы по ссылке поскольку они могут быть тяжелыми

func NewUpdateRoomUsersAction(
	loginUsersStore *store.LoginUsersStore,
	roomsStore *store.RoomsStore,
	clientsStore *client.ClientsStore,
) UpdateRoomUsersAction {
	return UpdateRoomUsersAction{
		loginUsersStore: loginUsersStore,
		roomsStore:      roomsStore,
		clientsStore:    clientsStore,
	}
}

func (ua UpdateRoomUsersAction) Handle(chatRequest request.ChatRequest) error {
	roomName, isRoomNameValid := chatRequest.Data["roomName"].(string)

	if !isRoomNameValid {
		errors.New("room name invalid")
	}

	if isPublicRoom := chatRequest.Data["roomName"].(string) == "Public"; isPublicRoom {
		var usersNames []string

		usersUniqueNames := make(map[string]bool)
		loginUsers := ua.loginUsersStore.GetAll()

		for _, user := range loginUsers {
			usersUniqueNames[user.Name()] = true
		}

		for name, _ := range usersUniqueNames {
			usersNames = append(usersNames, name)
		}

		addUsersResponse := response.NewUpdateRoomUsersResponse(usersNames, roomName).ToJson()
		chatClient, ok := ua.clientsStore.FindByToken(chatRequest.Token)

		if !ok {
			return errors.New("user not authorized")
		}

		chatClient.Send() <- &addUsersResponse

		return nil
	}

	chatRoom, isRoomExist := ua.roomsStore.FindByName(roomName)

	if !isRoomExist {
		return errors.New("room not exist")
	}

	roomUsers := chatRoom.UserNames()
	usersNames := make([]string, 0, len(roomUsers))

	for k := range roomUsers {
		usersNames = append(usersNames, k)
	}

	addUsersResponse := response.NewUpdateRoomUsersResponse(usersNames, roomName).ToJson()
	chatClient, ok := ua.clientsStore.FindByToken(chatRequest.Token)

	if !ok {
		return errors.New("user not authorized")
	}

	chatClient.Send() <- &addUsersResponse

	return nil
}
