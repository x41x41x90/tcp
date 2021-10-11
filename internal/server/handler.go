package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
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
	command, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	command = strings.Trim(command, "\n")
	fmt.Print("Input number of sending: ")
	fmt.Scan(&n)

	if strings.HasPrefix(command, "file") {
		s.file(command, n)
		return
	}

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

func (s *Server) file(command string, n int) {
	commands := strings.Split(command, " ")
	if len(commands) < 3 {
		s.logger.Println("not exist all parts")
		return
	}

	data, err := ioutil.ReadFile(commands[1])
	if err != nil {
		s.logger.Printf("read src file: %s\n", err)
		return
	}

	message := "file " + commands[2]

	i := 0
	s.connStorage.Range(
		func(key, value interface{}) bool {
			conn := value.(net.Conn)

			conn.SetWriteDeadline(time.Now().Add(s.tcpTimeout))
			if _, err = conn.Write([]byte(message + "\n")); err != nil {
				fmt.Printf("addr: %s err: %s\n", conn.RemoteAddr().String(), err)
				s.connStorage.Delete(key)
				return true
			}

			conn.SetReadDeadline(time.Now().Add(s.tcpTimeout))
			message, err = bufio.NewReader(conn).ReadString('\n')
			if message != "" && message != "ready\n" {
				err = fmt.Errorf(message)
			}
			if err != nil {
				fmt.Printf("addr: %s err: %s\n", conn.RemoteAddr().String(), err)
				s.connStorage.Delete(key)
				return true
			}

			for len(data) > 1024 {
				conn.SetWriteDeadline(time.Now().Add(s.tcpTimeout))
				if _, err = conn.Write(data[:1024]); err != nil {
					fmt.Printf("addr: %s err: %s\n", conn.RemoteAddr().String(), err)
					s.connStorage.Delete(key)
					return true
				}
				data = data[1024:]
			}

			fmt.Println(len(data))
			conn.SetWriteDeadline(time.Now().Add(s.tcpTimeout))
			if _, err = conn.Write(data); err != nil {
				fmt.Printf("addr: %s err: %s\n", conn.RemoteAddr().String(), err)
				s.connStorage.Delete(key)
				return true
			}

			i++
			return i < n
		},
	)
	fmt.Printf("command %s sended to %d clients\n", command, i)

}
