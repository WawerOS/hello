package user

import (
	"encoding/json"
	"fmt"
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

// NewMessage currently makes a new message
func NewMessage(sender string, receiver []string, message string) Message {
	msg := Message{sender, receiver, []byte(message)}
	return msg
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
	err := u.sendJSON.Encode(msg)
	if err != nil {
		fmt.Printf("%s at %s GOT a message with error %s", u.Name, u.recv.RemoteAddr(), err.Error())
	}

}

// Listen provides a channel of messages
func (u *User) Listen() chan Message {
	msgChan := make(chan Message)
	go func() {
		var msg Message
		for {
			err := u.recvJSON.Decode(msg)
			if err != nil {
				fmt.Printf("%s at %s SENT a message with error %s", u.Name, u.recv.RemoteAddr(), err.Error())
			}

			msgChan <- msg
		}
	}()

	return msgChan
}
