package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func ListenAndServe(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(fmt.Println("Listen err: &v", err))
	}
	defer listener.Close()
	log.Println(fmt.Sprintln("bind: %s, start listening...", address))

	for {
		accept, err := listener.Accept()
		if err != nil {
			log.Fatal(fmt.Sprintln("accept err: %v", err))
		}
		go Handle(accept)
	}
}

func Handle(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		readString, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connect close")
			} else {
				log.Println(err)
			}
			return
		}
		bytes := []byte(readString)
		conn.Write(bytes)
	}
}
func main() {
	ListenAndServe(":8080")
}
