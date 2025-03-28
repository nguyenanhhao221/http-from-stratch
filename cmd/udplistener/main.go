package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// What it does:
	// net.ResolveUDPAddr takes a network type ("udp") and an address ("localhost:42069").
	// It creates a *net.UDPAddr structure that represents the IP address and port where your UDP messages will be sent.
	// Why we need it:
	// We can’t just pass a string "localhost:42069" to net.DialUDP.
	// ResolveUDPAddr converts the human-readable address into a format Go understands.
	network, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalln(err)
	}

	// This creates a UDP connection from the client to the server.
	// The first argument "udp" specifies that we are using the UDP protocol.
	// The second argument (network) is supposed to be the client’s address (we don’t need to specify it, so it should be nil).
	// The third argument (network) is the server’s address (where we’re sending data).
	conn, err := net.DialUDP("udp", nil, network)
	if err != nil {
		log.Fatalln("error net.DialUDP: ", err.Error())
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error reading from stdin: ", err.Error())
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Println("error when writing to udp connection: ", err.Error())
		}
	}
}
