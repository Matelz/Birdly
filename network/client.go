package network

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Message struct {
	UserID  int
	Message  string
	MessageType  int
}

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

	pkg := packager(messagePayload{
		UserID:  Clients[Connection].id,
		Data:   "",
		MessageType:  2,
	})

	if err := Connection.WriteMessage(websocket.TextMessage, []byte(pkg)); err != nil {
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

		newMessage := unPackager(msg)

		switch newMessage.MessageType {
		case 1:
			Messages = append(Messages, Message{
				UserID:  newMessage.UserID,
				Message: newMessage.Data,
				MessageType:  newMessage.MessageType,
			})
		case 2:
		
			
			Messages = append(Messages, Message{
				UserID:  newMessage.UserID,
				Message: "joined the chat",
				MessageType:  newMessage.MessageType,
			})
		case 3:
			Messages = append(Messages, Message{
				UserID:  newMessage.UserID,
				Message: "left the chat",
				MessageType:  newMessage.MessageType,
			})
		}


		sub <- struct{}{}
	}
}

func SendMessage(msg string) {
	pkg := packager(messagePayload{
		UserID: Clients[Connection].id,
		Data:   msg,
		MessageType:  1,
	})

	err := Connection.WriteMessage(websocket.TextMessage, []byte(pkg))
	if err != nil {
		log.Println(err)
	}
}