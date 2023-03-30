package parser

import (
	"bufio"
	"bytes"
	"godis/redis"
	"godis/redis/protocol"
	"io"
	"log"
	"runtime/debug"
)

// Payload stores redis.Reply or error
type Payload struct {
	Data redis.Reply
	Err  error
}

// ParseStream reads data from io.Reader and send payloads through channel
//func ParseStream(reader io.Reader) <-chan *Payload {
//ch := make(chan *Payload)
//}

func Parse0(rawReader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			log.Panicln(err, string(debug.Stack()))
		}
	}()
	reader := bufio.NewReader(rawReader)
	for {
		// 上文中我们提到 RESP 是以行为单位的
		// 因为行分为简单字符串和二进制安全的BulkString
		line, err := reader.ReadBytes('\n')
		if err != nil {
			// 处理错误
			ch <- &Payload{Err: err}
			close(ch)
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			// there are some empty lines within replication traffic, ignore this error
			continue
		}
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		case '+':
			content := string(line[1:])
			ch <- &Payload{
				Data: protocol.MakeStatusReply(content),
			}
		}
	}
}
