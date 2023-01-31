package tcp

import (
	"bufio"
	"context"
	"godis/sync/atomic"
	"godis/sync/wait"
	"io"
	"log"
	"net"
	"sync"
)

//客户端连接的抽象
type Client struct {
	//tcp 连接
	Con  net.Conn
	Wait wait.Wait
}

type EchoHandler struct {
	// 保存所有工作状态client的集合(把map当set用)
	// 需使用并发安全的容器
	activeConn sync.Map

	// 关闭状态标识位
	closing atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	// 关闭中的 handler 不会处理新连接
	if h.closing.Get() {
		conn.Close()
		return
	}
	client := &Client{
		Con: conn,
	}
	h.activeConn.Store(client, struct{}{}) // 记住仍然存活的连接

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connection close")
				h.activeConn.Delete(client)
			} else {
				log.Println(err)
			}
			return
		}
		// 发送数据前先置为waiting状态，阻止连接被关闭
		client.Wait.Add(1)

		// 模拟关闭时未完成发送的情况
		//logger.Info("sleeping")
		//time.Sleep(10 * time.Second)

		b := []byte(msg)
		conn.Write(b)
		// 发送完毕, 结束waiting
		client.Wait.Done()
	}
}
