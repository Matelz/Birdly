package network

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

type user struct {
	Conn *websocket.Conn
	Id   string
	Name string
}

type messagePayload struct {
	UserID  string
	Data 	string
	NewData NewData
	User user
	MessageType   int
}

type NewData struct {
	Users map[string]user
	NewSocket *websocket.Conn
}

var Users = make(map[string]user)
var Clients = make(map[*websocket.Conn]user)

var Messages []Message

var Connection *websocket.Conn

func CreateServer(port string) {
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

		// Send the clients and users to the new client
		if len(Users) > 0 {
			err = conn.WriteMessage(websocket.TextMessage, packager(messagePayload{
				UserID:  "server",
				Data:   "",
				NewData: NewData{
					Users: Users,
				},
				MessageType:  4,
			}))
			if err != nil {
				log.Println(err)
				return
			}
		}

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				
				broadcastMessage(packager(messagePayload{
					UserID:  Users[Clients[conn].Id].Id,
					Data:   "",
					MessageType:  3,
				}))

				delete(Clients, conn)
				delete(Users, Clients[conn].Id)
				return
			}

			// log.Println("Received Connection from: ", conn.RemoteAddr().String())

			pkg := unPackager(msg);
			
			if pkg.MessageType == 2 {
				pkg.UserID = pkg.User.Id
				Clients[conn] = pkg.User

				Users[pkg.User.Id] = Clients[conn]

				pkg.NewData = NewData{
					Users: Users,
				}
			}
			broadcastMessage(packager(pkg))
		}
	})

	http.ListenAndServe(":"+port, nil)
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

func generateUUID() string {
	for {
		// Generate a unique id for the user based on the format "XXXX-XXXX"
		x := rand.Intn(10000)
		y := rand.Intn(10000)

		uuid := fmt.Sprintf("%04d-%04d", x, y)
		_, b := Users[uuid]

		if !b {	
			return uuid 
		}
	}
}