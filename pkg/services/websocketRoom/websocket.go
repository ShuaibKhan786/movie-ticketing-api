package websocketRoom

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(pool *Pool, roomId RoomId, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := Client{ID: ClientId(uuid.NewString()), Conn: ws, RoomId: roomId, Pool: pool}
	// Resgister client to regoster chanel
	pool.Register <- &client
	// Listen for message for this client
	go client.Read()
}
