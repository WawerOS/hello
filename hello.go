package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/Waweros/hello/server"
	"github.com/Waweros/hello/user"
)

var readConsole bufio.Reader

func clientRun(port int, addr net.IP) {
	portString := ":" + strconv.Itoa(port)
	fullAddrString := addr.String() + ":" + strconv.Itoa(port)

	send, err := net.Dial("tcp", fullAddrString)
	if err != nil {
		panic(err.Error())
	}

	ln, err := net.Listen("tcp", portString)
	if err != nil {
		panic(err.Error())
	}

	recv, err := ln.Accept()
	if err != nil {
		panic(err.Error())
	}

	sender := json.NewEncoder(send)
	getter := json.NewDecoder(recv)

	fmt.Print("NickName:")
	name, err := readConsole.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
	}

	sender.Encode(string(name))
	usersUpdate := make(chan []string)

	go func(usersUpdate chan []string) {
		var msg user.Message
		for {

			err := getter.Decode(&msg)
			if err != nil {
				fmt.Printf("Error receiveing message %s", err.Error())
			}

			if msg.Sender == "Server" {
				usersUpdate <- msg.Receiver
			} else {
				fmt.Println(string(msg.Message))
			}

		}
	}(usersUpdate)

	var users []string

	for {
		select {
		case users = <-usersUpdate:
			clientLoop(name, users, sender)
		default:
			clientLoop(name, users, sender)
		}
	}
}

func clientLoop(name string, users []string, sender *json.Encoder) {
	fmt.Print("=>")
	msgText, err := readConsole.ReadString('\n')
	if err != nil {
		fmt.Printf("Error sending message %s", err.Error())
	}

	msg := user.NewMessage(name, users, msgText)
	sender.Encode(msg)
}

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
