package action

import "chat/internal/request"

type PongAction struct {
	name string
}

func NewPongAction() PongAction {
	return PongAction{name: Pong}
}

func (p PongAction) Handle(request request.ChatRequest) error {
	return nil
}
