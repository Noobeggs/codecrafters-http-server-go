package main

import (
	"bytes"
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	_, err = c.Read(buf)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	request := bytes.Split(buf, []byte("\r\n"))
	start_line := bytes.Split(request[0], []byte(" "))
	path := start_line[1]

	if bytes.Equal(path, []byte("/")) {
		_, err = c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		_, err = c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
		os.Exit(1)
	}
}
