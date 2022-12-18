package auth

import "github.com/google/uuid"

type Token struct {
	value string
}

func (t Token) Value() string {
	return t.value
}

func NewToken() Token {
	return Token{value: generateToken()}
}

func NewTokenFromString(value string) Token {
	return Token{value: value}
}

func generateToken() string {
	return uuid.NewString()
}
