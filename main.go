package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func (s *Server) HandleWS(conn *websocket.Conn) {
	remoteAddr := conn.Request().RemoteAddr
	fmt.Println("New connection from", remoteAddr)
	s.conns[conn] = true
	defer func() {
		delete(s.conns, conn)
		conn.Close()
	}()
	s.readLoop(conn)
}

func (s *Server) readLoop(conn *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error reading message:", err)
			}
			break
		}
		if n == 0 {
			continue
		}
		msg := buf[:n]
		var position Position
		if err := json.Unmarshal(msg, &position); err == nil {
			s.positions[conn.Request().RemoteAddr] = position
		}
		fmt.Printf("current positions: %v\n", s.positions)
		positionsJSON, err := json.Marshal(s.positions)
		if err != nil {
			fmt.Println("Error marshalling positions:", err)
			continue
		}

		// Broadcast the updated positions to all clients
		s.broadcast(positionsJSON)
	}
}

func (s *Server) broadcast(msg []byte) {
	for conn := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Error broadcasting message:", err)
				ws.Close()
			}
		}(conn)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.HandleWS))

	fs := http.FileServer(http.Dir("client"))
	http.Handle("/", fs)

	ip := os.Getenv("SERVER_IP")
	if ip == "" {
		ip = "localhost"
	}

	fmt.Printf("Listening on ws://%s:8080/ws and serving static files on http://%s:8080\n", ip, ip)
	http.ListenAndServe(":8080", nil)
}
