package proxyserver

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"sync"
)

type ProxyListener struct {
	listener net.Listener
	port     int
}

type ProxyServer struct {
	listeners []ProxyListener
}

var (
	instance *ProxyServer
	once     sync.Once
)

func GetInstance() *ProxyServer {
	once.Do(func() {
		instance = &ProxyServer{}
		instance.listeners = make([]ProxyListener, 0, 500000)
	})
	return instance
}

func (s *ProxyServer) Start() {
	var wg sync.WaitGroup
	startPort := 10000
	endPort := 20000

	for port := startPort; port <= endPort; port++ {
		l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			log.Printf("无法监听端口 %d: %v\n", port, err)
			continue
		}
		defer l.Close()

		var pl ProxyListener
		pl.listener = l
		pl.port = port
		s.listeners = append(s.listeners, pl)

		wg.Add(1)
		go func(l ProxyListener) {
			defer wg.Done()
			s.acceptConnections(&pl)
		}(pl)
	}

	log.Println("ProxyServer Started.")
	wg.Wait()
}

func (s *ProxyServer) acceptConnections(pl *ProxyListener) {
	sessionId := 0
	for {
		conn, err := pl.listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				return
			}
			log.Printf("端口 %d 接受连接错误: %v\n", pl.port, err)
			continue
		}

		log.Printf("客户端 %s 连接到服务器端口 %d\n", conn.RemoteAddr(), pl.port)
		go s.handleConnection(conn, sessionId, pl)
	}
}

func (s *ProxyServer) handleConnection(conn net.Conn, sessionId int, pl *ProxyListener) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	for {
		binary.BigEndian.PutUint32(buffer[0:], uint32(pl.port))
		binary.BigEndian.PutUint32(buffer[4:], uint32(sessionId))
		n, err := conn.Read(buffer[8:])
		if err != nil {
			log.Printf("端口 %d 连接关闭: %v\n", pl.port, err)
			return
		}
		log.Println(n)

	}
}

/*
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




*/
