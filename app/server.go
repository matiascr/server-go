package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	DELETE RequestMethod = "DELETE"
)

type ResponseCode string

const (
	OK ResponseCode = "HTTP/1.1 200 OK" + CRLF
	NF ResponseCode = "HTTP/1.1 404 Not Found" + CRLF
	CR ResponseCode = "HTTP/1.1 201 Created" + CRLF
)

func responseCode(code ResponseCode) string {
	return string(ResponseCode(code))
}

func contentType(contentType string) string {
	return fmt.Sprintf("Content-Type: %s%s", contentType, CRLF)
}

func contentLength[T string | []byte](content T) string {
	return fmt.Sprintf("Content-Length: %d%s", len(content), CRLF)
}

func contentEncoding(contentEncoding string) string {
	return fmt.Sprintf("Content-Encoding: %s%s", contentEncoding, CRLF)
}

type Request struct {
	method  RequestMethod
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

func respond(conn net.Conn, message ...string) (int, error) {
	return conn.Write([]byte(fmt.Sprint(strings.Join(message, "") + CRLF)))
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
	request := parseData(data)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		os.Exit(1)
	}

	// Calls the HTTP request method
	switch {
	case request.method == GET &&
		request.path == "/":
		respond(conn, responseCode(OK))
	case request.method == GET &&
		strings.HasPrefix(request.path, "/echo"):
		echo(conn, request)
	case request.method == GET &&
		strings.HasPrefix(request.path, "/files"):
		getFiles(conn, request)
	case request.method == POST &&
		strings.HasPrefix(request.path, "/files"):
		postFiles(conn, request)
	default:
		respond(conn, responseCode(NF))
	}

	conn.Close()
}

// Parses the request lines into a Request
func parseData(data []byte) Request {
	request := string(bytes.Trim(data, "\x00"))
	requestLines := strings.Split(request, CRLF)
	parts := strings.Split(requestLines[0], " ")

	return Request{
		method:  RequestMethod(parts[0]),
		path:    parts[1],
		headers: requestLines[1:],
	}
}
