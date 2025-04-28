package main

import (
	"fmt"
	"io"
	"log"
	"net"
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

func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("error reading from connection: %v", err)
			}
			break
		}

		fmt.Print(string(buf[:n]))
	}

	fmt.Print("endof handleconnection")
	return nil
}

func main() {
	server := Server{}
	server.listenAndServe("tcp", "0.0.0.0:1080")
}
