package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Waweros/hello/user"
)

const (
	SERVER_SEND_PORT string = ":8800"
	SERVER_RECV_PORT string = ":8880"
)

// ListenForUsers listens for incoming users
func ListenForUsers(ln net.Listener) (user.User, chan user.User) {
	userChan := make(chan user.User)

	// This function accepts and finds places for all incoming users
	go func() {
		var name string
		for {
			// Accepting Users
			recv, err := ln.Accept()

			if err != nil {
				fmt.Printf("We couldn't talk because %s", err.Error())
				continue
			}

			// Finding out who our guest is
			addr := recv.RemoteAddr()
			send, err := net.Dial(addr.Network(), addr.String())
			if err != nil {
				fmt.Printf("We couldn't talk because %s", err.Error())
				continue
			}

			// Finding out what they want to be called
			dec := json.NewDecoder(recv)
			dec.Decode(&name)

			// Registering them and starting the chat
			u := user.NewUser(name, recv, send)
			userChan <- u
		}
	}()
	sendLocal, _ := net.Dial("tcp", SERVER_SEND_PORT)
	recvLocal, _ := net.Dial("tcp", SERVER_RECV_PORT)
	return user.NewUser("Server", recvLocal, sendLocal), userChan
}

// Exit string
var EXIT []byte = []byte("EXIT")

//UserRespond processes messages from a single user
func UserRespond(msgChan chan user.Message, userChan chan user.User, currentUsers []user.User) {
	for {

		select {
		case u := <-userChan:
			currentUsers = append(currentUsers, u)
		case msg := <-msgChan:
			user.MatchAndSend(msg, currentUsers)
		}

	}
}
