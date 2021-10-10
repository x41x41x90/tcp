package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func (s *Server) help() {
	fmt.Printf(
		`
	supported command:
		send - for sending message to connections
		info - output all connections
		help - output help message
`,
	)
}

func (s *Server) info() {
	s.connStorage.Range(
		func(key, value interface{}) bool {
			conn := value.(net.Conn)
			conn.SetWriteDeadline(time.Now().Add(s.tcpTimeout))
			if _, err := conn.Write([]byte("test\n")); err != nil {
				s.connStorage.Delete(key)
				return true
			}
			fmt.Printf("%s\n", conn.RemoteAddr().String())
			return true
		},
	)
}

func (s *Server) send() {
	var (
		command, message string
		n                int
		err              error
	)
	fmt.Printf("Input type message and message for sending\ntype message\tmessage\n")
	command, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Print("Input number of sending: ")
	fmt.Scan(&n)

	i := 0
	s.connStorage.Range(
		func(key, value interface{}) bool {
			conn := value.(net.Conn)

			conn.SetWriteDeadline(time.Now().Add(s.tcpTimeout))
			if _, err = conn.Write([]byte(command + "\n")); err == nil {
				conn.SetReadDeadline(time.Now().Add(s.tcpTimeout))
				message, err = bufio.NewReader(conn).ReadString('\n')
			}

			if err == nil {
				fmt.Printf("received from %s message %s\n", conn.RemoteAddr().String(), message)
				i++
			} else {
				fmt.Printf("addr: %s err: %s\n", conn.RemoteAddr().String(), err)
				s.connStorage.Delete(key)
			}

			return i < n
		},
	)
	fmt.Printf("command %s sended to %d clients\n", command, i)
}
