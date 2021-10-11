package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/geoirb/tcp/internal/server"
)

var (
	serverName = "server"

	port       = ":1337"
	tcpTimeout = 5 * time.Second
)

func main() {

	logger := log.New(os.Stdout, serverName+" ", log.LUTC)

	ln, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal(err)
	}

	srv := server.New(
		ln,
		tcpTimeout,
		logger,
	)

	srv.Run(context.Background())
}
