package websocketRoom

import "fmt"

type RoomId string
type ClientId string

const GlobalRoomId = "0"

type BroadcastMssg struct {
	Mssg   Message // message to boradcast
	Client *Client // boradcaster client info
	RoomId RoomId  // RoomId to boradcast
}

type Pool struct {
	Register   chan *Client                    // chanel to Broadcast register client
	UnRegister chan *Client                    // chanel to Broadcast unRegister client
	Broadcast  chan BroadcastMssg              // message boradcast chanel
	Clients    map[*Client]bool                // clients present in the pool
	Room       map[RoomId]map[ClientId]*Client // room structure
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Broadcast:  make(chan BroadcastMssg),
		Clients:    make(map[*Client]bool),
		Room:       make(map[RoomId]map[ClientId]*Client),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			_, ok := pool.Room[client.RoomId]
			if !ok {
				// Create a new room if not exist in the pool
				pool.Room[client.RoomId] = map[ClientId]*Client{}
			}
			pool.Room[client.RoomId][client.ID] = client
			// transmit new client join message to all client in the roomid
			for _, IClient := range pool.Room[client.RoomId] {
				IClient.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("New client joined, id: %s", client.ID)})
			}
			break
		case client := <-pool.UnRegister:
			delete(pool.Clients, client)                // delete client from pool
			delete(pool.Room[client.RoomId], client.ID) // delete client from Room
			// delete room if no user exist in that Room
			if len(pool.Room[client.RoomId]) == 0 {
				delete(pool.Room, client.RoomId)
			}
			// Transmit client disconnect message to Room
			for _, IClient := range pool.Room[client.RoomId] {
				IClient.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("Cliient: %s disconnected", client.ID)})
			}
			break
		case mssg := <-pool.Broadcast:
			if mssg.RoomId == GlobalRoomId {
				for iClient, _ := range pool.Clients {
					iClient.Conn.WriteJSON(mssg.Mssg)
				}
			} else {
				for _, iClient := range pool.Room[mssg.RoomId] {
					err := iClient.Conn.WriteJSON(mssg.Mssg)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
			break
		}
	}
}
