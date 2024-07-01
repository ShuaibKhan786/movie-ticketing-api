# Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var pool = NewPool()

func main() {

	server := http.NewServeMux()
	go pool.Start()

	server.HandleFunc("GET /click/", func(w http.ResponseWriter, r *http.Request) {
		roomId := r.URL.Query().Get("roomId")
		fmt.Println(roomId)
		pool.Broadcast <- BroadcastMssg{
			Mssg:   Message{Type: 1, Body: "Button Click"},
			client: nil,
			RoomId: roomId,
		}
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(fmt.Sprintf("click Now from server\nroomId: %s", roomId)))
	})

	server.HandleFunc("GET /pool/info/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"totalClient": len(pool.Clients),
			"totalRoom":   len(pool.Room),
		})
	})

	server.HandleFunc("GET /chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "chat.html")
		w.Header().Set("Content-Type", "text/html")

	})
	server.HandleFunc("/ws/{roomId}", func(w http.ResponseWriter, r *http.Request) {
		roomId := r.PathValue("roomId")
		fmt.Println("Websocket Link, roomId : ", roomId)
		ServeWS(pool, roomId, w, r)
	})

	log.Fatal(
		http.ListenAndServe(":5000", server),
	)

}

```