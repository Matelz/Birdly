package network

import (
	"example/other"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func ConnectToServer(sub chan struct{}) {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws",
	}

	dialer := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	Connection = conn

	if err := Connection.WriteMessage(websocket.TextMessage, []byte("new_user")); err != nil {
		log.Println(err)
	}

	HandleMessages(Connection, sub)
}

func HandleMessages(conn *websocket.Conn, sub chan struct{}) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		switch string(msg) {
		case "new_user":
			other.Messages = append(other.Messages, "New user joined")
		case "user_disconnected":
			other.Messages = append(other.Messages, "User disconnected")
		default:
			other.Messages = append(other.Messages, string(msg))
		}

		sub <- struct{}{}
	}
}

func SendMessage(msg string) {
	err := Connection.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println(err)
	}
}