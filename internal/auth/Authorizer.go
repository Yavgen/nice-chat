package auth

import (
	"chat/internal/client"
	"chat/internal/domain/model"
	"chat/internal/domain/store"
	"chat/internal/event"
	"chat/internal/request"
	"errors"
)

type Authorizer struct {
	registeredUsersStore *store.RegisteredUsersStore
	loginUsersStore      *store.LoginUsersStore
	clientsStore         *client.ClientsStore
	roomsStore           *store.RoomsStore
}

func NewAuthorizer(
	registeredUsersStore *store.RegisteredUsersStore,
	loginUsersStore *store.LoginUsersStore,
	clientsStore *client.ClientsStore,
	roomsStore *store.RoomsStore,
) Authorizer {
	return Authorizer{
		registeredUsersStore: registeredUsersStore,
		loginUsersStore:      loginUsersStore,
		clientsStore:         clientsStore,
		roomsStore:           roomsStore,
	}
}

func (a Authorizer) Authorize(authRequest request.AuthRequest) (event.ChatEvent, Token, error) {
	user, ok := a.registeredUsersStore.FindUserByName(authRequest.Name)

	if ok {
		token, err := a.login(user, authRequest)

		return event.NewLoginEvent(), token, err
	}

	token := a.register(authRequest)

	return event.NewRegisteredEvent(), token, nil
}

func (a Authorizer) login(user model.User, authRequest request.AuthRequest) (Token, error) {
	ok := user.PasswordIsValid(authRequest.Password)

	if !ok {
		return Token{}, errors.New("password wrong")
	}

	token := NewToken()
	a.loginUsersStore.MapUserByToken(user, token.Value())

	return token, nil
}

func (a Authorizer) register(authRequest request.AuthRequest) Token {
	user := model.NewUser(authRequest.Name, authRequest.Password)
	token := NewToken()

	a.registeredUsersStore.MapUserByName(user)
	a.loginUsersStore.MapUserByToken(user, token.Value())

	return token
}

func (a Authorizer) GetConnectedUserByRequest(connRequest request.ChatRequest) (model.User, Token, error) {
	token := NewTokenFromString(connRequest.Token)
	user, ok := a.loginUsersStore.FindUserByToken(token.Value())

	if !ok {
		return model.User{}, Token{}, errors.New("user dont login")
	}

	return user, token, nil
}
func (a Authorizer) LogoutChatClient(client client.ChatClient) {
	user, ok := a.loginUsersStore.FindUserByToken(client.Token())
	if !ok {
		return
	}

	a.loginUsersStore.DeleteByToken(client.Token())
	a.roomsStore.DeleteUserByName(user.Name())
	a.clientsStore.DeleteByToken(client.Token())
}
