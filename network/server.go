package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type user struct {
	conn *websocket.Conn
	id   int
	name string
}

var clients = make(map[*websocket.Conn]user)

var Connection *websocket.Conn

func CreateServer() {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				delete(clients, conn)

				broadcastMessage([]byte("user_disconnected"))
				return
			}

			if string(msg) == "new_user" {
				clients[conn] = user{
					conn: conn,
					id:   len(clients),
					name: "User",
				}
			}

			broadcastMessage(msg)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func broadcastMessage(msg []byte) {
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}