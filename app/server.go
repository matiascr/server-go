package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	buffer := make([]byte, 1024)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		_, err = conn.Read(buffer)

		msg := string(buffer)

		fmt.Println(msg)
		startLine := strings.Split(msg, "\r\n")[0]

		// method := strings.Split(startLine, " ")[0]
		path := strings.Split(startLine, " ")[1]
		// version := strings.Split(startLine, " ")[2]

		if path == "/" {
			_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			conn.Close()
		} else {
			_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			conn.Close()
		}

		if err != nil {
			fmt.Println("Error reading message: ", err.Error())
			os.Exit(1)
		}

	}
}
