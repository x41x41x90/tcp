package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Server ...
type Server struct {
	l          net.Listener
	tcpTimeout time.Duration

	connStorage sync.Map
	handler     map[string]func()

	logger *log.Logger
}

// New server.
func New(
	l net.Listener,
	tcpTimeout time.Duration,
	logger *log.Logger,
) (s *Server) {
	s = &Server{
		l:          l,
		tcpTimeout: tcpTimeout,
		logger:     logger,
	}

	s.handler = map[string]func(){
		"help": s.help,
		"info": s.info,
		"send": s.send,
	}

	go s.listen()
	return
}

func (s *Server) Run(ctx context.Context) {
	s.help()
	var command string
	for {
		fmt.Scan(&command)
		if h, isExist := s.handler[command]; isExist {
			h()
			continue
		}
		fmt.Println("command is unknown")
	}
}

func (s *Server) listen() {
	for {
		conn, err := s.l.Accept()
		if err != nil {
			s.logger.Fatal(err)
		}

		parts := strings.Split(conn.RemoteAddr().String(), ":")
		s.connStorage.Store(parts[0], conn)
	}
}
