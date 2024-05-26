package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

type Request struct {
	method  Method
	path    string
	headers []string
}

func main() {
	// Listen to port
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		// Accept connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Allocate request space
	data := make([]byte, 1024)
	_, err := conn.Read(data)

	if err != nil {
		fmt.Println("Error reading message: ", err.Error())
		os.Exit(1)
	}

	// Parse request lines
	request := parseRequest(data)

	switch {
	case request.method == GET &&
		request.path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	case request.method == GET &&
		strings.HasPrefix(request.path, "/echo"):
		echo(conn, request.path)

	case request.method == GET &&
		strings.HasPrefix(request.path, "/files"):
		getFiles(conn, request)

	case request.method == POST &&
		strings.HasPrefix(request.path, "/files"):
		postFiles(conn, request)

	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

	conn.Close()
}

// Parses the request lines into a Request
func parseRequest(data []byte) Request {
	request := string(bytes.Trim(data, "\x00"))
	requestLines := strings.Split(request, CRLF)
	parts := strings.Split(requestLines[0], " ")

	return Request{
		method:  Method(parts[0]),
		path:    parts[1],
		headers: requestLines[1:],
	}
}

// routes

// Implements `/echo` route
func echo(conn net.Conn, path string) {
	body := strings.Split(path, "/echo/")[1]

	respFormat := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s"
	conn.Write([]byte(fmt.Sprintf(respFormat, len(body), body)))
}

// Implements `/getFiles` route
func getFiles(conn net.Conn, request Request) {
	directory := os.Args[2]
	fileName := strings.TrimPrefix(request.path, "/files/")

	data, err := os.ReadFile(directory + fileName)

	if err != nil {
		resp := "HTTP/1.1 404 Not Found\r\n\r\n"
		conn.Write([]byte(resp))
		return
	}

	resp := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s"
	conn.Write([]byte(fmt.Sprintf(resp, len(data), string(data))))

}

// Implements the `/upload` route
func postFiles(conn net.Conn, request Request) {
	directory := os.Args[2]
	fileName := strings.TrimPrefix(request.path, "/files/")

	err := os.WriteFile(directory+fileName, []byte(request.headers[len(request.headers)-1]), os.ModeTemporary)

	if err != nil {
		resp := "HTTP/1.1 404 Not Found\r\n\r\n"
		conn.Write([]byte(resp))
		return
	}

	resp := "HTTP/1.1 201 Created\r\n\r\n"
	conn.Write([]byte(fmt.Sprintf(resp)))
}
