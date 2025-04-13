package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"httpfromtcp.haonguyen.tech/internal/request"
	"httpfromtcp.haonguyen.tech/internal/response"
	"httpfromtcp.haonguyen.tech/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, testHandler)
	if err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			log.Fatalf("error closing server: %v\n", err)
		}
	}()
	log.Println("Server started on port:", port)

	// Common pattern to exit the program.
	sigChn := make(chan os.Signal, 1)
	signal.Notify(sigChn, syscall.SIGINT, syscall.SIGTERM)
	<-sigChn
	log.Println("Server gracefully shutdown")
}

func testHandler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			Message:    "Your problem is not my problem\n",
			StatusCode: response.StatusBadRequest,
		}
	case "/myproblem":
		return &server.HandlerError{
			Message:    "Woopsie, my bad\n",
			StatusCode: response.StatusServerInternalError,
		}
	default:
		body := "All good, frfr\n"
		_, err := w.Write([]byte(body))
		if err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusServerInternalError,
				Message:    fmt.Sprintf("error writing body for valid request: %v\n", err),
			}
		}
	}

	return nil
}
