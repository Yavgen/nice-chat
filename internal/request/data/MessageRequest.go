package data

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type MessageRequest struct {
	Message  string
	RoomName string
}

func NewMessageRequest() MessageRequest {
	return MessageRequest{}
}

func (mr MessageRequest) Validate() error {
	return validation.ValidateStruct(&mr,
		validation.Field(&mr.Message, validation.Required, validation.Length(1, 1000)),
		validation.Field(&mr.RoomName, validation.Required, validation.Length(1, 20)),
	)
}
