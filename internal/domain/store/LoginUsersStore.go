package store

import (
	"chat/internal/domain/model"
)

type LoginUsersStore struct {
	loginUsers map[string]model.User
}

func NewLoginUsersStore() *LoginUsersStore {
	return &LoginUsersStore{loginUsers: map[string]model.User{}}
}

func (ls LoginUsersStore) FindUserByToken(token string) (model.User, bool) {
	if user, ok := ls.loginUsers[token]; ok {
		return user, true
	}

	return model.User{}, false
}

func (ls LoginUsersStore) MapUserByToken(user model.User, token string) {
	ls.loginUsers[token] = user
}

func (ls LoginUsersStore) GetAll() map[string]model.User {
	return ls.loginUsers
}

func (ls LoginUsersStore) DeleteByToken(token string) {
	delete(ls.loginUsers, token)
}

func (ls LoginUsersStore) FindUserTokenByName(userName string) (string, bool) {
	for token, user := range ls.loginUsers {
		if userName == user.Name() {
			return token, true
		}
	}

	return "", false
}
