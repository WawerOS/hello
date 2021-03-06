package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Waweros/hello/server"
	"github.com/Waweros/hello/user"
)

func serverRun(port int) {
	portString := ":" + strconv.Itoa(port)

	fmt.Printf("Server Starting on %s!\n", portString)

	// creating a listener to listen for new users
	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Started listening for users
	us, ch := server.ListenForUsers(ln)

	// Makes a list of currentUsers
	currentUsers := make([]user.User, 1)
	currentUsers[1] = us
	for {
		/*
			Waits for a new user and when one comes add to the Register
			and starts to listen to them
		*/
		newPerson := <-ch
		currentUsers := append(currentUsers, newPerson)

		// Collecting everyones name
		usersName := make([]string, 0)
		for _, u := range currentUsers {
			usersName = append(usersName, u.Name)
		}

		// Telling everyone someone new came
		newMember := user.NewMessage("Server", usersName, "")
		user.MatchAndSend(newMember, currentUsers)

		go server.UserRespond(newPerson.Listen(), ch, currentUsers)
	}

}
