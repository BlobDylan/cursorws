package main

import (
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns     map[*websocket.Conn]bool
	positions map[string]Position
	mu        sync.Mutex
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
