package main

import (
	"bytes"
	"flag"
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	var dirFlag = flag.String("directory", "", "Specify directory of files")
	flag.Parse()
	fmt.Printf("Directory: %v", *dirFlag)

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(c, *dirFlag)
	}
}

func handleConnection(connection net.Conn, directory string) {
	defer connection.Close()

	buf := make([]byte, 1024)
	_, err := connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	request := bytes.Split(buf, []byte("\r\n"))
	start_line := bytes.Split(request[0], []byte(" "))
	path := start_line[1]

	if bytes.Equal(path, []byte("/")) {
		_, err = connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if after, found := bytes.CutPrefix(path, []byte("/echo/")); found {
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-length: %v\r\n\r\n%v", len(after), string(after))
		_, err = connection.Write([]byte(response))
	} else if bytes.Equal(path, []byte("/user-agent")) {
		for i := 1; i < len(request); i++ {
			if after, found := bytes.CutPrefix(request[i], []byte("User-Agent: ")); found {
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-length: %v\r\n\r\n%v", len(after), string(after))
				_, err = connection.Write([]byte(response))
			}
		}
	} else if after, found := bytes.CutPrefix(path, []byte("/files/")); found {
		data, err := os.ReadFile(fmt.Sprintf("%v%v", directory, string(after)))
		if err != nil {
			fmt.Println("Error reading file: ", err.Error())
			_, err = connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-length: %v\r\n\r\n%v", len(data), string(data))
			_, err = connection.Write([]byte(response))
		}
		if err != nil {
			fmt.Println("Error writing data: ", err.Error())
			os.Exit(1)
		}
	} else {
		_, err = connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	if err != nil {
		fmt.Println("Error writing data: ", err.Error())
		os.Exit(1)
	}
}
