package network

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Message struct {
	UserID  string
	Message  string
	MessageType  int
}

var uConn *websocket.Conn

var me user = user{
	Conn: nil,
	Id:   "",
	Name: "",
}

func ConnectToServer(ip string, port string, name string, sub chan struct{}) {
	var h string
	
	if port != ""{
		h = fmt.Sprintf("%s:%s", ip, port)
	} else {
		h = ip
	}

	u := url.URL{
		Scheme: "ws",
		Host:   h,
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

	uConn = conn

	me = user{
		Conn: conn,
		Id:   generateUUID(),
		Name: name,
	}

	pkg := packager(messagePayload{
		UserID:  "",
		Data:   name,
		User: 	me,
		MessageType:  2,
	})

	// log.Println(name)

	if err := uConn.WriteMessage(websocket.TextMessage, []byte(pkg)); err != nil {
		log.Println(err)
	}

	HandleMessages(uConn, sub)
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

			Users = newMessage.NewData.Users
		case 3:
			Messages = append(Messages, Message{
				UserID:  newMessage.UserID,
				Message: "left the chat",
				MessageType:  newMessage.MessageType,
			})
		case 4:
			Users = newMessage.NewData.Users
		}

		sub <- struct{}{}
	}
}

func SendMessage(msg string) {
	pkg := packager(messagePayload{
		UserID: me.Id,
		Data:   msg,
		MessageType:  1,
	})

	err := uConn.WriteMessage(websocket.TextMessage, []byte(pkg))
	if err != nil {
		log.Println(err)
	}
}