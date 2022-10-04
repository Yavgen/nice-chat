package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	connection *websocket.Conn

	send chan []byte
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
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
			}

			writer, writerError := client.connection.NextWriter(websocket.TextMessage)

			if writerError != nil {
				return
			}

			writer.Write(message)

			queueCount := len(client.send)

			for i := 0; i < queueCount; i++ {
				writer.Write(newline)
				writer.Write(<-client.send)
			}

			if writerCloseError := writer.Close(); writerCloseError != nil {
				return
			}
		}
	}

}
