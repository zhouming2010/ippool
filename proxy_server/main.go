package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
}

func (s *Server) listenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

type Message struct {
	Action  string                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s", conn.RemoteAddr())

	for {
		// 读取消息长度
		var length uint32
		if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
			if err == io.EOF {
				log.Printf("Connection closed by client")
				return
			}
			log.Printf("Error reading message length: %v", err)
			return
		}

		// 限制消息大小 (防止DoS攻击)
		if length > 1024*1024 { // 1MB
			log.Printf("Message too large: %d bytes", length)
			return
		}

		// 读取消息内容
		data := make([]byte, length)
		if _, err := io.ReadFull(conn, data); err != nil {
			log.Printf("Error reading message content: %v", err)
			return
		}

		// 解析JSON消息
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			sendErrorResponse(conn, "invalid message format")
			continue
		}

		log.Printf("Received message: %+v", msg)

		// 处理消息
		response, err := processMessage(msg)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			sendErrorResponse(conn, err.Error())
			continue
		}

		// 发送响应
		if err := sendMessage(conn, response); err != nil {
			log.Printf("Error sending response: %v", err)
			return
		}
	}
}

func processMessage(msg Message) (Message, error) {
	// 根据不同的action处理消息
	switch msg.Action {
	case "hello":
		return Message{
			Action: "hello_response",
			Payload: map[string]interface{}{
				"message": "Hello, client!",
			},
		}, nil
	case "echo":
		return Message{
			Action:  "echo_response",
			Payload: msg.Payload,
		}, nil
	default:
		return Message{}, fmt.Errorf("unknown action: %s", msg.Action)
	}
}

func sendMessage(conn net.Conn, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 先发送消息长度
	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return err
	}

	// 再发送消息内容
	if _, err := conn.Write(data); err != nil {
		return err
	}

	return nil
}

func sendErrorResponse(conn net.Conn, errorMsg string) {
	errMsg := Message{
		Action: "error",
		Payload: map[string]interface{}{
			"message": errorMsg,
		},
	}
	if err := sendMessage(conn, errMsg); err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	//test2()
	//server := Server{}
	//server.listenAndServe("tcp", "0.0.0.0:1080")
	test4()
}

func test2() {
	m := make(map[string]int)
	m["apple"] = 1
	m["pear"] = 2
	for key, value := range m {
		log.Println("key=", key, ",value =", value)
	}

	c := make(chan string)

	go func() {
		c <- "hello"
	}()

	str := <-c
	log.Printf("str = %s", str)

}

func test3() {
	ch := make(chan int)

	go func() {
		time.Sleep(10 * time.Second)
		ch <- 42
	}()

	for {
		select {
		case v := <-ch:
			fmt.Println("Received value:", v)
		default:
			fmt.Println("No value received - channel is empty")
		}
	}
}

func test4() {
	ch := make(chan []byte, 100)

	go func() {
		n := 0
		for {
			time.Sleep(1 * time.Second)
			message := []byte("队列消息: " + time.Now().Format(time.RFC3339))
			ch <- message
			n++
		}
	}()

	for v := range ch {
		s := string(v)
		fmt.Println(s)
	}
}
