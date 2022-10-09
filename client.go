package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	connection *websocket.Conn
	send       chan []byte
}

func (client *Client) readPipe() {
	defer func() {
		client.connection.Close()
	}()

	for {
		_, message, readError := client.connection.ReadMessage()
		if readError != nil {
			log.Printf("client read error %v", readError)
			break
		}

		message = bytes.Replace(message, newline, space, -1)
		broadcast <- message
	}
}

func (client *Client) writePipe() {
	defer func() {
		client.connection.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				closeResponse := Response{
					Data:   map[string]interface{}{"message": message},
					Status: "ok",
				}

				client.connection.WriteJSON(closeResponse)
			}

			writer, writerError := client.connection.NextWriter(websocket.TextMessage)

			if writerError != nil {
				return
			}

			messageResponse := Response{
				Data:   map[string]interface{}{"message": string(message)},
				Status: "ok",
			}

			json.NewEncoder(writer).Encode(messageResponse)

			queueCount := len(client.send)

			for i := 0; i < queueCount; i++ {
				queuedMessageResponse := Response{
					Data:   map[string]interface{}{"message": string(<-client.send)},
					Status: "ok",
				}

				json.NewEncoder(writer).Encode(queuedMessageResponse)
			}

			if writerCloseError := writer.Close(); writerCloseError != nil {
				return
			}
		}
	}

}
