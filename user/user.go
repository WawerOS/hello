package user

import (
	"encoding/json"
	"fmt"
	"net"
)

// User is used to keep track of user metadata
type User struct {
	Name string
	recv net.Conn
	send net.Conn
}

// Message gives a message
type Message struct {
	sender   User
	receiver []User
	message  []byte
}

// New is used to create  new User's
func New(name string, recv net.Conn, send net.Conn) User {
	user := User{name, recv, send}
	return user
}

// Send send's a message to the User
func (u *User) Send(msg Message) error {
	// Making Message into json
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Sending json
	n := 0
	for n != len(m) {
		n, err = u.send.Write(m)

		if err != nil {
			return err
		}

	}

	return nil
}

// Listen provides a channel of messages
func (u *User) Listen(msgChan chan Message) error {
	// Making byte channel
	byteChan := make(chan []byte)
	go func(c chan []byte) {
		buff := make([]byte, 255)
		for {
			msgLength, err := u.recv.Read(buff)
			if err != nil {
				fmt.Println(err)
			}
			c <- buff[0 : msgLength-1]
		}
	}(byteChan)

	var msgByte []byte
	var inputBuffer []byte
	var msg Message
	var err error

	for {
		msgByte = <-byteChan
		for i := 0; i < len(msgByte); i++ {

			if msgByte[i] != 0 {
				inputBuffer = append(inputBuffer, msgByte[i])
			} else if msgByte[i] == 0 {
				err = json.Unmarshal(inputBuffer, &msg)
				if err != nil {
					fmt.Println(err)
				}

				msgChan <- msg
			}
		}

	}

}
