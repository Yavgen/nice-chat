package action

import "chat/internal/request"

type ChatAction interface {
	Handle(request request.ChatRequest) error
}
