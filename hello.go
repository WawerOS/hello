package main

import (
  "fmt"
  "io"
  "net"
)

func main() {
  
  fmt.Printf("Server Starting!\n")
  
  ln, err := net.Listen("tcp",":8800")
  if err != nil {
         fmt.Printf("We had a bad thing %s",err)
  }

  accepted := 0
  for {
       conn, err := ln.Accept()
       if err != nil {
         fmt.Printf("We had a bad thing %s",err)
         conn.Close()
       }
       defer conn.Close()
       accepted++
       conn.Write([]byte("Hello!!!"))
       io.Copy(conn,conn)
   }

}
