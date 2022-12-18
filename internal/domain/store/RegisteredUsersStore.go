package store

import (
	"chat/internal/domain/model"
)

type RegisteredUsersStore struct {
	registeredUsers map[string]model.User
}

func NewRegisteredUsersStore() *RegisteredUsersStore {
	return &RegisteredUsersStore{registeredUsers: map[string]model.User{}}
}

func (rs RegisteredUsersStore) FindUserByName(name string) (model.User, bool) {
	if user, ok := rs.registeredUsers[name]; ok {
		return user, true
	}

	return model.User{}, false
}

func (rs RegisteredUsersStore) MapUserByName(user model.User) {
	rs.registeredUsers[user.Name()] = user
}
