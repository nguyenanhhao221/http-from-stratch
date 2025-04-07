package response

import (
	"fmt"
	"io"
	"log"

	"httpfromtcp.haonguyen.tech/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusServerInternalError StatusCode = 500
)

var statusCodeMap = map[StatusCode]string{
	StatusOK:                  "HTTP/1.1 200 OK\r\n",
	StatusBadRequest:          "HTTP/1.1 400 Bad Request\r\n",
	StatusServerInternalError: "HTTP/1.1 500 Internal Server Error\r\n",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine, ok := statusCodeMap[statusCode]
	if !ok {
		log.Println("cannot find default reason phrase status")
		statusLine = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}

	if _, err := w.Write([]byte(statusLine)); err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	defaultHeaders := map[string]string{
		"Content-Length": fmt.Sprintf("%d", contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
	return defaultHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}
	// write empty line by the end of headers
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}
