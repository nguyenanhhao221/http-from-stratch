package response

import (
	"fmt"

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

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	defaultHeaders := map[string]string{
		"Content-Length": fmt.Sprintf("%d", contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}

	for k, v := range defaultHeaders {
		h.Set(k, v)
	}
	return h
}
