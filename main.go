package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	log.Println("Listen for TCP traffic on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// handle error
			log.Fatalln("error: ", err.Error())
		}
		log.Println("Accept connection from: ", conn.RemoteAddr())

		lineCh := getLinesChannel(conn)
		for line := range lineCh {
			fmt.Println(line)
		}
		log.Printf("Connection to %s closed \n", conn.RemoteAddr())
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	// Read the file 8 bytes at a time
	buffer := make([]byte, 8)

	currentLineStr := ""
	go func() {
		defer f.Close()
		defer close(ch)
		for {
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineStr != "" {
					ch <- currentLineStr
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			for i := range len(parts) - 1 {
				ch <- fmt.Sprintf("%s%s", currentLineStr, parts[i])
				currentLineStr = ""
			}
			currentLineStr += parts[len(parts)-1]
		}
	}()
	return ch
}
