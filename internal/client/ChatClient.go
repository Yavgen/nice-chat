package client

import (
	"chat/internal/response"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type ChatClient struct {
	connection *websocket.Conn
	send       chan *response.JsonResponse
	token      string
}

func (c ChatClient) Send() chan *response.JsonResponse {
	return c.send
}

func (c ChatClient) Token() string {
	return c.token
}

func NewClient(connection *websocket.Conn, send chan *response.JsonResponse, token string) ChatClient {
	return ChatClient{connection: connection, send: send, token: token}
}

func (c ChatClient) CloseSendCh() {
	close(c.send)
}

func (c ChatClient) CloseConnection() {
	err := c.connection.Close()

	//TODO: обработать
	if err != nil {
		return
	}
}

func (c ChatClient) ReadMessage() (messageType int, p []byte, err error) {
	return c.connection.ReadMessage()
}

func (c ChatClient) NextWriter(messageType int) (io.WriteCloser, error) {
	return c.connection.NextWriter(messageType)
}

func (c ChatClient) WriteJSON(json response.JsonResponse) error {
	return c.connection.WriteJSON(json)
}
