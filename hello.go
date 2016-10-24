package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	fmt.Printf("Server Starting!\n")

	ln, err := net.Listen("tcp", ":8800")
	if err != nil {
		fmt.Printf("We had a bad thing %s", err)
	}
	ch := clientConn(ln)

	for {
		go handleConn(<-ch)
	}
}

func clientConn(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)

	go func() {
		i := 0
		for {
			client, err := listener.Accept()
			if err != nil {
				fmt.Printf("couldn't accept: " + err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	them := bufio.NewReader(client)
	for {
		line, err := them.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		client.Write(line)
		os.Stdout.Write(line)
	}
}
