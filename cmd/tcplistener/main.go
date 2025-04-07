package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp.haonguyen.tech/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("err closing listener :%v\n", err)
		}
	}()
	log.Println("Listen for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// handle error
			log.Fatalln("error: ", err.Error())
		}
		log.Println("Accept connection from: ", conn.RemoteAddr())

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("error RequestFromReader: %v\n", err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", request.RequestLine.Method)
		fmt.Println("- Target:", request.RequestLine.RequestTarget)
		fmt.Println("- Version:", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		log.Printf("Connection to %s closed \n", conn.RemoteAddr())
	}
}
