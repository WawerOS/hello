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

	go func() {
		var name string
		for {
			recv, err := ln.Accept()

			if err != nil {
				fmt.Printf("We couldn't talk because %s", err.Error())
				continue
			}

			addr := recv.RemoteAddr()
			send, err := net.Dial(addr.Network(), addr.String())
			if err != nil {
				fmt.Printf("We couldn't talk because %s", err.Error())
				continue
			}

			dec := json.NewDecoder(recv)
			dec.Decode(&name)
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

//UserRespond does things
func UserRespond(msgChan chan user.Message, userChan chan user.User) {
	currentUsers := make([]user.User, 0)
	for {

		select {
		case u := <-userChan:
			currentUsers = append(currentUsers, u)
		case msg := <-msgChan:
			user.MatchAndSend(msg, currentUsers)
		}

	}
}
