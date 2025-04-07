package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"httpfromtcp.haonguyen.tech/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port)
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
