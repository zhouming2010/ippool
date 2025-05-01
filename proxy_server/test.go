package main

import (
	"bufio"
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn    net.Conn
	send    chan []byte
	mu      sync.Mutex
	closing bool
}

type TCPServer struct {
	listener   net.Listener
	clients    map[*Client]bool
	mu         sync.RWMutex
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewTCPServer(addr string) (*TCPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := &TCPServer{
		listener:   listener,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go server.run()
	return server, nil
}

func (s *TCPServer) run() {
	// 启动消息队列消费者
	go s.consumeMessageQueue()

	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
			log.Printf("客户端连接: %s (当前连接数: %d)", client.conn.RemoteAddr(), len(s.clients))

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()
			log.Printf("客户端断开: %s (当前连接数: %d)", client.conn.RemoteAddr(), len(s.clients))

		case message := <-s.broadcast:
			s.mu.RLock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					// 防止阻塞，如果客户端发送缓冲区满则关闭连接
					go client.close()
				}
			}
			s.mu.RUnlock()
		}
	}
}

// 模拟从消息队列消费数据
func (s *TCPServer) consumeMessageQueue() {
	for {
		// 这里模拟从消息队列获取数据
		// 实际项目中可以替换为Kafka/RabbitMQ/Redis等的消费者
		time.Sleep(5 * time.Second)
		message := []byte("队列消息: " + time.Now().Format(time.RFC3339))
		s.broadcast <- message
		log.Println("从消息队列获取到新消息:", string(message))
	}
}

func (s *TCPServer) Start() {
	defer s.listener.Close()
	log.Println("TCP服务器启动，监听地址:", s.listener.Addr())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("接受连接错误:", err)
			continue
		}

		client := &Client{
			conn: conn,
			send: make(chan []byte, 100),
		}

		s.register <- client

		go s.handleRead(client)
		go s.handleWrite(client)
	}
}

func (s *TCPServer) handleRead(client *Client) {
	defer func() {
		s.unregister <- client
		client.conn.Close()
	}()

	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		log.Printf("收到来自 %s 的消息: %s", client.conn.RemoteAddr(), string(message))
	}
}

func (s *TCPServer) handleWrite(client *Client) {
	defer client.conn.Close()

	for message := range client.send {
		client.mu.Lock()
		if client.closing {
			client.mu.Unlock()
			return
		}

		_, err := client.conn.Write(append(message, '\n'))
		if err != nil {
			client.mu.Unlock()
			client.close()
			return
		}
		client.mu.Unlock()
	}
}

func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closing {
		c.closing = true
		c.conn.Close()
	}
}

func test() {
	server, err := NewTCPServer(":8080")
	if err != nil {
		log.Fatal("启动服务器失败:", err)
	}

	server.Start()
}
