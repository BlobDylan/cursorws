package main

import "golang.org/x/net/websocket"

type Server struct {
	conns     map[*websocket.Conn]bool
	positions map[string]Position
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func NewServer() *Server {
	return &Server{
		conns:     make(map[*websocket.Conn]bool),
		positions: make(map[string]Position),
	}
}
