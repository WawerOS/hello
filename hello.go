package main

import (
	"fmt"
	"net"

	"github.com/Waweros/hello/server"
	"github.com/Waweros/hello/user"
)

func serverRun() {

	fmt.Printf("Server Starting!\n")

	// creating a listener to listen for new users
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err.Error())
	}

	// Started listening for users
	us, ch := server.ListenForUsers(ln)

	// Makes a list of currentUsers
	users := make([]user.User, 1)
	users[1] = us

	for {
		/*
			Waits for a new user and when one comes add to the Register
			and starts to listen to them
		*/
		newPerson := <-ch
		users := append(users, newPerson)
		go server.UserRespond(newPerson.Listen(), ch, users)
	}

}
