package client

import (
	"chat/internal/response"
)

type BroadcastChannel struct {
	channel chan *response.JsonResponse
}

func NewBroadcastChannel() BroadcastChannel {
	return BroadcastChannel{channel: make(chan *response.JsonResponse)}
}

func (bc BroadcastChannel) Push(chatResponse response.JsonResponse) {
	bc.channel <- &chatResponse
}

func (bc BroadcastChannel) Listen() <-chan *response.JsonResponse {
	return bc.channel
}
