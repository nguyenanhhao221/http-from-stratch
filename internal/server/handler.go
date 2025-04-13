package server

import (
	"io"
	"log"

	"httpfromtcp.haonguyen.tech/internal/request"
	"httpfromtcp.haonguyen.tech/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Message    string
	StatusCode response.StatusCode
}

func (he HandlerError) Write(w io.Writer) {
	e := response.WriteStatusLine(w, he.StatusCode)
	if e != nil {
		log.Printf("error handle connection when writing status line: %v\n", e)
		return
	}
	err := response.WriteHeaders(w, response.GetDefaultHeaders(len(he.Message)))
	if err != nil {
		log.Printf("error writing headers %v", err)
		return
	}
	_, err = w.Write([]byte(he.Message))
	if err != nil {
		log.Printf("error writing body %v", err)
		return
	}
}
