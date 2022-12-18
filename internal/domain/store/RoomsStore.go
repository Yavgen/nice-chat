package store

import (
	"chat/internal/domain/model"
)

type RoomsStore struct {
	rooms map[string]model.Room
}

func NewRoomsStore() *RoomsStore {
	return &RoomsStore{rooms: make(map[string]model.Room)}
}

func (rs RoomsStore) FindByName(roomName string) (model.Room, bool) {
	room, ok := rs.rooms[roomName]

	if !ok {
		return model.Room{}, false
	}

	return room, true
}

func (rs RoomsStore) MapByName(roomName string, chatRoom model.Room) {
	rs.rooms[roomName] = chatRoom
}

// TODO добавить к пользователю связь на комнаты чтобы оптимизировать
func (rs RoomsStore) DeleteUserByName(userName string) {
	for _, chatRoom := range rs.rooms {
		chatRoom.DeleteUserFromRoom(userName)
	}
}
