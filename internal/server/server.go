package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"httpfromtcp.haonguyen.tech/internal/request"
	"httpfromtcp.haonguyen.tech/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("error when accepting connection %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	w := response.NewWriter(conn)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection in handle: %v\n", err)
		}
	}()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("error: parsing request: %v\n", err)
		err := w.WriteStatusLine(response.StatusBadRequest)
		if err != nil {
			log.Printf("error when write status line %v\n", err)
			return
		}
		body := fmt.Appendf(nil, "error: parsing request: %v\n", err)
		err = w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		if err != nil {
			log.Printf("error when write header %v\n", err)
			return
		}
		_, err = w.WriteBody(body)
		if err != nil {
			log.Printf("error when write body %v\n", err)
			return
		}
		return
	}
	s.handler(w, r)
}
