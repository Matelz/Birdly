package network

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type user struct {
	conn *websocket.Conn
	id   int
	Name string
}

type messagePayload struct {
	UserID  int
	Data 	string
	MessageType   int
}

var Users = make(map[int]user)
var Clients = make(map[*websocket.Conn]user)

var Messages []Message

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
				delete(Clients, conn)
				delete(Users, Clients[conn].id)

				broadcastMessage(packager(messagePayload{
					UserID:  Clients[conn].id,
					Data:   "",
					MessageType:  3,
				}))
				return
			}

			pkg := unPackager(msg);

			if pkg.MessageType == 2 {
				Clients[conn] = user{
					conn: conn,
					id:   len(Clients),
					Name: "User",
				}

				Users[len(Clients)] = Clients[conn]
			}

			broadcastMessage(packager(pkg))
		}
	})

	http.ListenAndServe(":8080", nil)
}

func broadcastMessage(msg []byte) {
	for client := range Clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
		}
	}
}

func packager(message messagePayload) []byte {
	pkg, _ := json.Marshal(message)
	return pkg
}

func unPackager(msg []byte) messagePayload {
	pkg := messagePayload{}
	json.Unmarshal(msg, &pkg)
	return pkg
}