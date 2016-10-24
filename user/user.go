package user

import (
	"encoding/json"
	"net"
)

// User is used to keep track of user metadata
type User struct {
	Name     string
	recv     net.Conn
	send     net.Conn
	recvJSON *json.Decoder
	sendJSON *json.Encoder
}

// Message gives a message
type Message struct {
	Sender   string
	Receiver []string
	Message  []byte
}

// NewUser is used to create  new User's
func NewUser(name string, recv net.Conn, send net.Conn) User {
	user := User{name, recv, send, json.NewDecoder(recv), json.NewEncoder(send)}
	return user
}

// MatchAndSend sends Message to certain Users
func MatchAndSend(msg Message, users []User) {
	for _, u := range users {
		for _, s := range msg.Receiver {
			if s == u.Name {
				u.Send(msg)
			}
		}
	}
}

// Send send's a message to the User
func (u *User) Send(msg Message) {
	u.sendJSON.Encode(msg)
}

// Listen provides a channel of messages
func (u *User) Listen(msgChan chan Message) {

	go func() {
		var msg Message
		for {
			u.recvJSON.Decode(msg)
			msgChan <- msg
		}
	}()

}