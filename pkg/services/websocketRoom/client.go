package websocketRoom

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     ClientId
	Conn   *websocket.Conn
	Pool   *Pool
	RoomId RoomId
}

type Message struct {
	Type int    `json:"type"` // message type (binary, ascii text)
	Body string `json:"body"` // message body
}

func (client *Client) Read() {
	defer func() {
		client.Pool.UnRegister <- client
		client.Conn.Close()
	}()

	for {
		// Read Message
		msgType, p, err := client.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		mssg := Message{Type: msgType, Body: string(p)}
		// Broadcast message to specific RoomId
		client.Pool.Broadcast <- BroadcastMssg{
			Mssg:   mssg,
			RoomId: client.RoomId,
			Client: client,
		}
	}
}
