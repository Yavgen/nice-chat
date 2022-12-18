package action

import (
	"chat/internal/client"
	"chat/internal/domain/store"
	"errors"
)

const (
	Pong             = "Pong"
	Message          = "Message"
	CreateRoom       = "CreateRoom"
	AppendUserToRoom = "AppendUserToRoom"
	UpdateRoomUsers  = "UpdateRoomUsers"
)

type Factory struct {
	actions              map[string]bool
	broadcastCh          client.BroadcastChannel
	loginUsersStore      *store.LoginUsersStore
	registeredUsersStore *store.RegisteredUsersStore
	roomsStore           *store.RoomsStore
	clientsStore         *client.ClientsStore
}

func NewFactory(
	broadcastCh client.BroadcastChannel,
	loginUsersStore *store.LoginUsersStore,
	registeredUsersStore *store.RegisteredUsersStore,
	roomsStore *store.RoomsStore,
	clientsStore *client.ClientsStore,
) Factory {
	return Factory{
		actions:              make(map[string]bool),
		broadcastCh:          broadcastCh,
		loginUsersStore:      loginUsersStore,
		registeredUsersStore: registeredUsersStore,
		roomsStore:           roomsStore,
		clientsStore:         clientsStore,
	}
}

func (f Factory) MakeAction(name string) (ChatAction, error) {
	if name == Pong {
		return NewPongAction(), nil
	}

	if name == Message {
		return NewMessageAction(f.broadcastCh, f.roomsStore, f.loginUsersStore), nil
	}

	if name == CreateRoom {
		return NewCreateRoomAction(f.roomsStore, f.loginUsersStore, f.clientsStore), nil
	}

	if name == AppendUserToRoom {
		return NewAppendUserToRoomAction(
			f.registeredUsersStore,
			f.roomsStore,
			f.loginUsersStore,
			f.clientsStore,
		), nil
	}

	if name == UpdateRoomUsers {
		return NewUpdateRoomUsersAction(f.loginUsersStore, f.roomsStore, f.clientsStore), nil
	}

	return nil, errors.New("unknown action")
}
