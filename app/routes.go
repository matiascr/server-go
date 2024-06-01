package main

import (
	"net"
	"os"
	"strings"
)

// Implements GET `/echo` route
func echo(conn net.Conn, request Request) {
	body := strings.TrimPrefix(request.path, "/echo/")

	encodingHeader := "Accept-Encoding: "

	for _, v := range request.headers {
		if strings.HasPrefix(v, encodingHeader) {
			encoding := strings.TrimPrefix(v, encodingHeader)

			if encoding != "invalid-encoding" {
				respond(conn,
					responseCode(OK),
					contentEncoding(encoding),
					contentType("text/plain"),
					contentLength(body),
					body,
				)
			}
		}
	}

	respond(conn,
		responseCode(OK),
		contentType("text/plain"),
		contentLength(body),
		CRLF,
		body,
	)
}

// Implements GET `/files` route
func getFiles(conn net.Conn, request Request) {
	directory := os.Args[2]
	fileName := strings.TrimPrefix(request.path, "/files/")

	data, err := os.ReadFile(directory + fileName)

	if err != nil {
		respond(conn, responseCode(NF))
	}

	respond(conn,
		responseCode(OK),
		contentType("application/octet-stream"),
		contentLength(data),
		CRLF,
		string(data),
	)

}

// Implements POST `/files` route
func postFiles(conn net.Conn, request Request) {
	directory := os.Args[2]
	fileName := strings.TrimPrefix(request.path, "/files/")

	err := os.WriteFile(
		directory+fileName,
		[]byte(request.headers[len(request.headers)-1]),
		os.ModeTemporary,
	)

	if err != nil {
		respond(conn, responseCode(NF))
		return
	}

	respond(conn, responseCode(CR))
}
