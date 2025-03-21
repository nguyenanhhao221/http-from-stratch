package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// open the file
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatalln(err)
	}

	// Read the file 8 bytes at a time
	b1 := make([]byte, 8)

	for {
		n, err := f.Read(b1)
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		if n == 0 {
			break
		}
		fmt.Printf("read: %s\n", string(b1[:n]))
	}
}
