package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func ListenAndServer(address string) {
	//绑定监听地址
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(fmt.Sprintf("listen err:%v", err))
	}
	defer listen.Close()
	log.Printf(fmt.Sprintf("bind address:%d, start listen...", address))

	for {
		//accept 会一直阻塞直到有新的连接建立或者listen中断才会返回
		accept, err := listen.Accept()
		if err != nil {
			log.Fatal(fmt.Sprintf("accept error:%v", err))
		}
		go handle(accept)
	}
}

func handle(conn net.Conn) {
	//使用 bufio 标准库提供的缓冲区功能
	reader := bufio.NewReader(conn)
	for {
		// ReadString 会一直阻塞直到遇到分隔符 '\n'
		// 遇到分隔符后会返回上次遇到分隔符或连接建立后收到的所有数据, 包括分隔符本身
		// 若在遇到分隔符之前遇到异常, ReadString 会返回已收到的数据和错误信息
		readString, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connection close")
			} else {
				log.Println(err)
			}
			return
		}
		b := []byte(readString)
		// 将收到的信息发送给客户端
		conn.Write(b)
	}
}

func main() {
	ListenAndServer(":8000")
}
