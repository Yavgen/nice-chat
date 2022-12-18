package model

type User struct {
	name     string
	password string
}

func NewUser(name string, password string) User {
	return User{
		name:     name,
		password: password,
	}
}

func (u User) PasswordIsValid(password string) bool {
	return u.password == password
}

func (u User) Name() string {
	return u.name
}
