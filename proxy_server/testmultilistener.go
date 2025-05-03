package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const (
	startPort = 10000
	endPort   = 40000 // 总共30001个端口(10000-40000)
)

func TestMultiListener() {
	// 创建监听器列表
	listeners := make([]net.Listener, 0, endPort-startPort+1)
	var wg sync.WaitGroup

	// 启动所有端口监听
	for port := startPort; port <= endPort; port++ {
		l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			fmt.Printf("无法监听端口 %d: %v\n", port, err)
			continue
		}
		defer l.Close()
		listeners = append(listeners, l)
	}

	fmt.Printf("Echo服务器正在监听端口 %d 到 %d\n", startPort, endPort)

	// 为每个监听器启动一个accept循环
	for i, l := range listeners {
		wg.Add(1)
		go func(port int, listener net.Listener) {
			defer wg.Done()
			acceptConnections(port, listener)
		}(startPort+i, l)
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n接收到中断信号，关闭服务器...")
	wg.Wait()
}

func acceptConnections(port int, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			// 检查是否是正常关闭
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				return
			}
			fmt.Printf("端口 %d 接受连接错误: %v\n", port, err)
			continue
		}

		fmt.Printf("客户端 %s 连接到服务器端口 %d\n", conn.RemoteAddr(), port)
		go handleConnection(conn, port)
	}
}

func handleConnection(conn net.Conn, port int) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// 读取数据
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("端口 %d 连接关闭: %v\n", port, err)
			return
		}

		// 回显消息并打印端口信息
		echoMsg := fmt.Sprintf("[端口%d] %s", port, message)
		_, err = writer.WriteString(echoMsg)
		if err != nil {
			fmt.Printf("端口 %d 写入错误: %v\n", port, err)
			return
		}
		writer.Flush()
	}
}
