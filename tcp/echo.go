package tcp

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Handle interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// Config stores tcp server properties
type Config struct {
	Address    string        `yaml:"address"`
	MaxConnect uint32        `yaml:"max-connect"`
	Timeout    time.Duration `yaml:"timeout"`
}

//监听并提供服务，并在收到 closeChan 发来的关闭通知后关闭
func ListenAndServer(listen net.Listener, handle Handle, closeChan <-chan struct{}) {
	//监听关闭通知
	go func() {
		<-closeChan
		log.Println("shutting down...")
		// 停止监听，listener.Accept()会立即返回 io.EOF
		_ = listen.Close()
		// 关闭应用层服务器
		_ = handle.Close()
	}()
	// 在异常退出后释放资源
	defer func() {
		// close during unexpected error
		_ = listen.Close()
		_ = handle.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		// 监听端口, 阻塞直到收到新连接或者出现错误
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		// 开启 goroutine 来处理新连接
		log.Println("accept link...")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handle.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}

// ListenAndServeWithSignal 监听中断信号并通过 closeChan 通知服务器关闭
func ListenAndServeWithSignal(cfg *Config, handle Handle) error {
	closeChan := make(chan struct{})
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("bind: %s, start listening..."), cfg.Address)
	ListenAndServer(listen, handle, closeChan)
	return nil
}
