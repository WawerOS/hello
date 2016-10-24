package server

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Waweros/hello/user"
)

// ListenForUsers listens for incoming users
func ListenForUsers(ln net.Listener) chan user.User {
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

	return userChan
}
