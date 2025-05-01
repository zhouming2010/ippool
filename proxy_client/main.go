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

type Message struct {
	Action  string                 `json:"action"`
	Payload map[string]interface{} `json:"payload"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:1080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.SetOutput(os.Stdout)
	log.Println("Connected to server")

	// 发送hello消息
	helloMsg := Message{
		Action: "hello",
		Payload: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	}
	if err := sendMessage(conn, helloMsg); err != nil {
		log.Fatalf("Failed to send hello message: %v", err)
	}

	// 读取响应
	response, err := readMessage(conn)
	if err != nil {
		log.Fatalf("Failed to read hello response: %v", err)
	}
	log.Printf("Server response: %+v", response)

	// 发送echo消息
	echoMsg := Message{
		Action: "echo",
		Payload: map[string]interface{}{
			"text": "This is a test message",
		},
	}
	if err := sendMessage(conn, echoMsg); err != nil {
		log.Fatalf("Failed to send echo message: %v", err)
	}

	// 读取响应
	response, err = readMessage(conn)
	if err != nil {
		log.Fatalf("Failed to read echo response: %v", err)
	}
	log.Printf("Server response: %+v", response)
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

func readMessage(conn net.Conn) (Message, error) {
	// 读取消息长度
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return Message{}, err
	}

	// 限制消息大小
	if length > 1024*1024 { // 1MB
		return Message{}, fmt.Errorf("message too large: %d bytes", length)
	}

	// 读取消息内容
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return Message{}, err
	}

	// 解析JSON消息
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return Message{}, err
	}

	return msg, nil
}
